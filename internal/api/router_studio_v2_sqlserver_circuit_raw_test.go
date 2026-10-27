package api

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

// TestSQLServerCircuitRawIntegration uses API-created profiles and the actual
// journal/outbox/sender path. The fixture and operation guards precede DDL.
func TestSQLServerCircuitRawIntegration(t *testing.T) {
	fixture := SQLServerCircuitFixtureOrSkip(t)
	for _, tc := range []struct {
		name, tagType, column string
		kind                  measurement.ExactType
		value                 any
	}{
		{"uint64_max", "uint64", "value_decimal", measurement.ExactUint64, ^uint64(0)},
		{"bit", "bool", "value_bool", measurement.ExactBool, true},
		{"unicode", "string", "value_text", measurement.ExactText, "運轉中🟢"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newSQLServerCircuitRawFixture(t, fixture, tc.tagType)
			group, err := f.groups.Get(t.Context(), f.instantID)
			require.NoError(t, err)
			member := group.Group.Members[0]
			observed := f.base.Add(9*time.Second + 123456*time.Microsecond)
			f.clock.Store(observed.Add(time.Second).UnixNano())
			require.NoError(t, f.pipe.AcceptSample(t.Context(), measurement.SampleEnvelope{
				SampleID: "sql-raw-" + tc.name, WorkspaceID: f.record.ID,
				DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID,
				SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
				ValueType: string(tc.kind), Value: tc.value, ObservedAt: observed,
				ReceivedAt: observed.Add(30 * time.Millisecond), TimeOrigin: "gateway", Quality: schema.QualityGood,
			}))
			var journal string
			require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_samples WHERE group_id=? AND sample_id=?`,
				f.instantID, "sql-raw-"+tc.name).Scan(&journal))
			require.True(t, json.Valid([]byte(journal)), "the sample must be durable before closing its window")
			f.closeAt(t, 15*time.Second)
			var payload, digest, key string
			require.Eventually(t, func() bool {
				return f.db.QueryRowContext(t.Context(), `SELECT payload,payload_digest,effect_key FROM wg_delivery_outbox WHERE group_id=? AND state='sql_committed'`,
					f.instantID).Scan(&payload, &digest, &key) == nil
			}, 20*time.Second, 25*time.Millisecond)
			row, err := groupdelivery.DecodeRowPayload([]byte(payload))
			require.NoError(t, err)
			actualDigest, err := groupdelivery.RowPayloadDigest(row)
			require.NoError(t, err)
			require.Equal(t, digest, actualDigest)
			require.Equal(t, key, row.EffectKey)
			var raw any
			var gotObserved time.Time
			var provenance string
			require.NoError(t, f.target.QueryRowContext(t.Context(), "SELECT ["+tc.column+"],observed_at,provenance FROM "+
				sqlServerQualified(f.namespace, group.Group.Destination.TableName)+" WHERE record_id=@p1 AND workspace_id=@p2 AND group_id=@p3",
				row.RecordID, f.record.ID, f.instantID).Scan(&raw, &gotObserved, &provenance))
			declared := map[string]string{"value_decimal": "DECIMAL(38,18)", "value_bool": "BIT", "value_text": "NVARCHAR(MAX)"}[tc.column]
			decoded, err := dbtarget.DecodeExactValue(dbtarget.SQLDialectSQLServer, declared, tc.kind, raw)
			require.NoError(t, err)
			expected, err := measurement.ExactFromGo(tc.kind, tc.value)
			require.NoError(t, err)
			require.True(t, expected.Equal(decoded), "typed SQL readback must retain the exact source value")
			require.True(t, observed.Equal(gotObserved), "DATETIME2(6) must retain UTC microseconds")
			require.True(t, json.Valid([]byte(provenance)))
			var receipts int
			require.NoError(t, f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM "+sqlServerQualified(f.namespace, dbtarget.EffectReceiptTable)+
				" WHERE effect_key=@p1 AND payload_digest=@p2", key, digest).Scan(&receipts))
			require.Equal(t, 1, receipts)
			cleanupSQLServerCircuitEffects(t, f.target, f.namespace, group.Group.Destination.TableName, []string{row.RecordID}, map[string]string{key: digest}, map[string]dbtarget.EncodedRow{row.RecordID: row})
		})
	}
}

func newSQLServerCircuitRawFixture(t *testing.T, fixture *SQLServerCircuitFixture, tagType string, prepareOnly ...bool) *sqlServerCircuitPlainFixture {
	t.Helper()
	f := newManagedDestinationFixture(t, "sqlite", tagType)
	_, err := f.db.ExecContext(t.Context(), `UPDATE database_connectors SET kind='sqlserver',connection_config=?,identity_revision='connector-1' WHERE id='connector-A'`, fixture.ConnectionConfigJSON())
	require.NoError(t, err)
	f.record, err = f.workspace.BindDatabaseConnector(t.Context(), "connector-A")
	require.NoError(t, err)
	f.namespace = fixture.Manifest.Schema
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	capability, err := targets.ProbeSQLServerRecordingAtRevision(t.Context(), "connector-A", "connector-1")
	require.NoError(t, err)
	require.NoError(t, fixture.ValidateCapability(capability))
	if !capability.SupportsManagedLong {
		t.Skip("NOT_RUN: SQL Server managed-long operation gate is closed")
	}
	opened, err := targets.OpenDestination(t.Context(), "connector-A", "connector-1")
	require.NoError(t, err)
	if opened.Close != nil {
		t.Cleanup(func() { require.NoError(t, opened.Close()) })
	}
	f.target = opened.DB
	request := sqlServerCircuitRawProfileRequest(t, f, fixture.Manifest.Database, tagType)
	created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
	require.Equal(t, http.StatusOK, created.Code, created.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, created))
	id := saved["id"].(string)
	clock := &atomic.Int64{}
	base := time.Now().UTC().Truncate(15 * time.Second).Add(30 * time.Second)
	clock.Store(base.Add(-time.Second).UnixNano())
	f.groups.WithClock(func() time.Time { return time.Unix(0, clock.Load()).UTC() })
	if len(prepareOnly) == 1 && prepareOnly[0] {
		return &sqlServerCircuitPlainFixture{
			managedGroupFixture: f, fixture: fixture, target: f.target, targetService: targets,
			instantID: id, base: base, clock: clock,
		}
	}
	token := previewLong(t, f, id)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, id, token).Status)
	registerSQLServerCircuitCreatedTableCleanup(t, f, targets, token)
	apply := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
	require.Equal(t, http.StatusOK, apply.Code, apply.Body.String())
	snapshot, err := f.groups.ResolveAppliedAt(t.Context(), id, base.Add(time.Hour))
	require.NoError(t, err)
	base = snapshot.EffectiveAt
	clock.Store(base.UnixNano())
	tags := tag.NewService(tag.NewSQLRepository(f.db))
	store := groupdelivery.NewStore(f.db)
	config := grouppipeline.Config{
		NodeID: "sqlserver-raw", Owner: "sqlserver-raw/test",
		Now:          func() time.Time { return time.Unix(0, clock.Load()).UTC() },
		TickInterval: time.Hour, WorkerInterval: 20 * time.Millisecond, ReconcileInterval: time.Hour,
		MaxFutureSkew: time.Minute,
	}
	pipe := grouppipeline.New(grouppipeline.Dependencies{
		Groups: f.groups, Tags: tags, Destinations: targets,
		Inspector: dbtarget.NewReadOnlyTableInspector(targets), Store: store,
	}, config)
	require.NoError(t, pipe.Start(t.Context()))
	t.Cleanup(func() { _ = pipe.Stop(time.Second) })
	return &sqlServerCircuitPlainFixture{
		managedGroupFixture: f, fixture: fixture, target: f.target, targetService: targets,
		instantID: id, base: base, clock: clock, pipe: pipe, tags: tags, store: store, pipeConfig: config,
	}
}

func sqlServerCircuitRawProfileRequest(t *testing.T, f managedGroupFixture, database, tagType string) map[string]any {
	t.Helper()
	request := profileRequest(t, f.writeGroupRouterFixture, "instantaneous", 15, "second", "sql-raw-"+tagType, "")
	g := request["group"].(map[string]any)
	g["layout_version"], g["storage_mode"] = "managed_long_v1", "plain"
	destination := g["destination"].(map[string]any)
	destination["database"], destination["table_schema"] = database, f.namespace
	g["members"].([]map[string]any)[0]["value_role"] = "state"
	return request
}

// This admission regression uses local SQLite, independently of real SQL
// availability, so invalid source roles cannot hide behind a fixture skip.
func TestSQLServerCircuitRawProfileAdmission(t *testing.T) {
	for _, tagType := range []string{"uint64", "bool", "string"} {
		t.Run(tagType, func(t *testing.T) {
			f := newManagedDestinationFixture(t, "sqlite", tagType)
			response := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath,
				sqlServerCircuitRawProfileRequest(t, f, "", tagType))
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			_, saved := groupSaveData(t, decodeJSONBody(t, response))
			group, err := f.groups.Get(t.Context(), saved["id"].(string))
			require.NoError(t, err)
			require.Equal(t, schema.DataType(tagType), group.Group.Members[0].SourceValueType)
			require.Nil(t, group.Group.Members[0].MeasurementID)
		})
	}
}
