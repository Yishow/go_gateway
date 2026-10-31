package main

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/api/handlers"
	appconfig "go-gateway/internal/config"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"modernc.org/sqlite"
)

type recoveryFaultConnector struct {
	path     string
	mu       sync.Mutex
	mode     string
	callback func()
	calls    int
}

func (c *recoveryFaultConnector) Driver() driver.Driver { return &sqlite.Driver{} }
func (c *recoveryFaultConnector) Connect(context.Context) (driver.Conn, error) {
	raw, err := c.Driver().Open(c.path)
	if err != nil {
		return nil, err
	}
	return &recoveryFaultConn{Conn: raw, owner: c}, nil
}
func (c *recoveryFaultConnector) arm(mode string, callback func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.mode = mode
	c.callback = callback
}
func (c *recoveryFaultConnector) count() int { c.mu.Lock(); defer c.mu.Unlock(); return c.calls }

type recoveryFaultConn struct {
	driver.Conn
	owner *recoveryFaultConnector
}

func (c *recoveryFaultConn) Begin() (driver.Tx, error) {
	raw, err := c.Conn.Begin()
	if err != nil {
		return nil, err
	}
	return &recoveryFaultTx{Tx: raw, owner: c.owner}, nil
}
func (c *recoveryFaultConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	raw, err := c.Conn.(driver.ConnBeginTx).BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &recoveryFaultTx{Tx: raw, owner: c.owner}, nil
}

type recoveryFaultTx struct {
	driver.Tx
	owner *recoveryFaultConnector
}

func (tx *recoveryFaultTx) Commit() error {
	tx.owner.mu.Lock()
	mode, callback := tx.owner.mode, tx.owner.callback
	tx.owner.mode = ""
	tx.owner.callback = nil
	if mode != "" {
		tx.owner.calls++
	}
	tx.owner.mu.Unlock()
	if mode == "before" {
		if err := tx.Tx.Rollback(); err != nil {
			return err
		}
		return errors.New("PRIVATE_COMMIT_FAILURE")
	}
	err := tx.Tx.Commit()
	if err != nil {
		return err
	}
	if mode != "" {
		if callback != nil {
			callback()
		}
		return errors.New("PRIVATE_LOST_ACK")
	}
	return nil
}

func newRecoveryFaultProductionFixture(t *testing.T) (wordProductionFixture, *recoveryFaultConnector) {
	t.Helper()
	plc, config := wordWireServer(t, schema.ProtocolModbusTCP)
	fault := &recoveryFaultConnector{path: filepath.Join(t.TempDir(), "recovery-fault.db")}
	db := sql.OpenDB(fault)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	_, err := db.Exec(`INSERT INTO devices(id,name,description,protocol,status,connection_config,readiness_status) VALUES('word-device','PLC','','modbus_tcp','draft',?,'{}')`, config)
	require.NoError(t, err)
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	services, err := wireGatewayServices(db, manager)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, services.runtime.Stop(context.Background())) })
	require.NoError(t, services.device.ProbeAndActivate(t.Context(), "word-device"))
	_, err = services.workspace.AttachDevice(t.Context(), "word-device")
	require.NoError(t, err)
	router := wordProductionRouter(services, db)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, httptest.NewRequest("POST", appconfig.Get().API.BasePath+"/datalink/devices/word-device/activate", nil))
	require.Equal(t, 200, res.Code, res.Body.String())
	return wordProductionFixture{services, db, plc, manager, router}, fault
}

func recoveryFaultSave(t *testing.T, f wordProductionFixture) point.FormatRecoveryRequest {
	t.Helper()
	require.NoError(t, f.services.modbusShare.ApplySettings(t.Context(), modbusshare.Settings{Enabled: false, BindAddress: "127.0.0.1", Port: 15020, SlaveID: 1, CapacityRegisters: 64}))
	rule, p, _ := f.createRule(t, schema.ProtocolModbusTCP, 0, schema.DataTypeUint32, "ABCD")
	revisions := modbusshare.NewSQLWorkspaceRevisionStore(f.db)
	hydrateStartupShare(t.Context(), f.services, revisions, modbusshare.NewReconciler(f.services.modbusShare, revisions))
	require.NoError(t, f.services.runtime.Start(t.Context()))
	rule, p = enableRecoveryRuleShare(t, f, rule, p)
	record, err := f.services.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	read, err := f.services.point.ReadFormat(t.Context(), point.FormatReadScope{WorkspaceID: record.ID, DeviceID: p.DeviceID, PointID: p.ID, RuleID: rule.ID})
	require.NoError(t, err)
	save := point.FormatMutationRequest{FormatScope: point.FormatScope{WorkspaceID: record.ID, DeviceID: p.DeviceID, PointID: p.ID, RuleID: rule.ID, ExpectedWorkspaceRevision: read.WorkspaceRevision, ExpectedPointRevision: &read.PointRevision, ExpectedRuleRevision: read.RuleRevision, RequestID: "recovery-fault"}, FormatMode: schema.PointFormatOverride, ConfiguredFormat: "CDAB"}
	raw, err := json.Marshal(save)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	request := httptest.NewRequest("PUT", appconfig.Get().API.BasePath+"/datalink/points/"+p.ID+"/format", bytes.NewReader(raw))
	request.Header.Set("Content-Type", "application/json")
	f.router.ServeHTTP(res, request)
	require.Equal(t, 200, res.Code, res.Body.String())
	var saved handlers.PointFormatMutationResponse
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &saved))
	req := point.FormatRecoveryRequest{FormatScope: save.FormatScope, ExpectedConnectionRevision: read.ConnectionRevision, ExpectedSettingsRevision: f.services.modbusShare.Settings().SettingsRevision}
	req.ExpectedPointRevision = &saved.Data.Point.RevisionID
	req.ExpectedWorkspaceRevision = saved.Data.WorkspaceRevision
	return req
}

func recoveryFaultDurable(t *testing.T, db *sql.DB, id string) []string {
	t.Helper()
	var intent string
	require.NoError(t, db.QueryRow(`SELECT revision_id||'|'||format_mode||'|'||format_override||'|'||data_format FROM points WHERE id=?`, id).Scan(&intent))
	result := []string{intent}
	rows, err := db.Query(`SELECT key||'|'||value FROM system_settings WHERE key LIKE 'point_format_receipt:%' OR key LIKE 'modbus_share_workspace_revision:%' ORDER BY key`)
	require.NoError(t, err)
	defer rows.Close()
	for rows.Next() {
		var value string
		require.NoError(t, rows.Scan(&value))
		result = append(result, value)
	}
	require.NoError(t, rows.Err())
	return result
}

func TestProductionFormatRecoveryCommitFaultCleanup(t *testing.T) {
	for _, mode := range []string{"before", "lost-ack", "foreign-after-commit"} {
		t.Run(mode, func(t *testing.T) {
			f, fault := newRecoveryFaultProductionFixture(t)
			req := recoveryFaultSave(t, f)
			f.plc.set(0, []uint16{0x3344, 0x1122})
			before := recoveryFaultDurable(t, f.db, req.PointID)
			evidence, err := f.services.point.ReadFormatRecovery(t.Context(), req)
			require.NoError(t, err)
			prepared, err := datalinkruntime.PrepareFormatRecoveryBindings(req.WorkspaceID, evidence.Point, evidence.Profile, evidence.TargetInputs.Device, evidence.TargetInputs.Mappings, evidence.TargetInputs.Tags)
			require.NoError(t, err)
			var foreignLease uint64
			var callbackErr error
			var callback func()
			if mode == "foreign-after-commit" {
				callback = func() {
					bindingProof, bindingOK := f.services.runtime.CaptureFormatRecoveryBindings(req.PointID)
					proof, ok := f.services.runtime.CaptureOrdinaryPoint(req.PointID)
					if !ok || !bindingOK {
						callbackErr = errors.New("cannot capture foreign install")
						return
					}
					foreignLease, ok = f.services.runtime.InstallPreparedFormatRecovery(&evidence.Point, proof, evidence.Profile)
					_, activated := f.services.runtime.ActivateFormatRecoveryBindings(req.PointID, foreignLease, evidence.Profile, bindingProof, prepared)
					if !ok || !activated {
						callbackErr = errors.New("cannot activate foreign install")
						return
					}
					state, e := f.services.modbusShare.CheckHydration(t.Context())
					if e != nil {
						callbackErr = e
						return
					}
					state.ReadinessToken = "foreign-ready-token"
					state.Readiness = true
					state.State = modbusshare.HydrationStateReady
					f.services.modbusShare.SetHydrationState(state)
				}
			}
			fault.arm(mode, callback)
			raw, err := json.Marshal(map[string]any{"format_recovery": req})
			require.NoError(t, err)
			res := httptest.NewRecorder()
			request := httptest.NewRequest("POST", appconfig.Get().API.BasePath+"/datalink/modbus-share/reconcile", bytes.NewReader(raw))
			request.Header.Set("Content-Type", "application/json")
			f.router.ServeHTTP(res, request)
			require.Equal(t, 503, res.Code, res.Body.String())
			require.NotContains(t, res.Body.String(), "PRIVATE_")
			require.Contains(t, res.Body.String(), "point_format_unavailable")
			require.Equal(t, 1, fault.count())
			require.NoError(t, callbackErr)
			state, err := f.services.modbusShare.CheckHydration(t.Context())
			require.NoError(t, err)
			if mode == "foreign-after-commit" {
				require.NotZero(t, foreignLease)
				require.True(t, f.services.runtime.MatchesFormatRecoveryBoundLease(req.PointID, foreignLease, prepared))
				require.Equal(t, "foreign-ready-token", state.ReadinessToken)
				values := f.services.scheduler.PollNowContext(t.Context(), []string{req.PointID})
				require.Len(t, values, 1)
				require.Equal(t, uint32(0x11223344), values[0].Value)
			} else {
				current, ok := f.services.runtime.CaptureOrdinaryPoint(req.PointID)
				require.True(t, ok)
				require.Zero(t, current.Lease, "failed attempt must leave no owned staged or live slot")
				require.False(t, f.services.runtime.MatchesFormatRecoveryBindings(prepared), "cleanup must remove only this attempt binding")
				require.False(t, state.Readiness)
				require.Empty(t, state.ReadinessToken)
				require.Empty(t, f.services.scheduler.PollNowContext(t.Context(), []string{req.PointID}))
			}
			require.Equal(t, before, recoveryFaultDurable(t, f.db, req.PointID))
			tx, err := f.db.Begin()
			require.NoError(t, err, "cleanup must leave SQL pool usable")
			require.NoError(t, tx.Rollback())
		})
	}
}
