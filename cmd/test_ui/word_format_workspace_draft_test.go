package main

import (
	"bytes"
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	appconfig "go-gateway/internal/config"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestProductionWordWorkspaceDiagnosedDraftCreate(t *testing.T) {
	for _, diagnosed := range []bool{false, true} {
		name := "offline-unprobed"
		if diagnosed {
			name = "diagnosed-draft"
		}
		t.Run(name, func(t *testing.T) { testWordWorkspaceDraftCreate(t, diagnosed, false) })
	}
	t.Run("late-foreign-activation", func(t *testing.T) { testWordWorkspaceDraftCreate(t, true, true) })
}
func testWordWorkspaceDraftCreate(t *testing.T, diagnosed, foreignActivate bool) {
	_, config := wordWireServer(t, schema.ProtocolModbusTCP)
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "draft.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	_, err = db.Exec(`INSERT INTO devices(id,name,description,protocol,status,connection_config,readiness_status) VALUES('draft-plc','PLC','','modbus_tcp','draft',?,'{}')`, config)
	require.NoError(t, err)
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { _ = manager.CloseAll() })
	services, err := wireGatewayServices(db, manager)
	require.NoError(t, err)
	t.Cleanup(func() { _ = services.runtime.Stop(context.Background()) })
	_, err = services.workspace.AttachDevice(t.Context(), "draft-plc")
	require.NoError(t, err)
	router := wordProductionRouter(services, db)
	base := appconfig.Get().API.BasePath + "/datalink"
	response := httptest.NewRecorder()
	if diagnosed {
		router.ServeHTTP(response, httptest.NewRequest("POST", base+"/devices/draft-plc/test", nil))
		require.Equal(t, 200, response.Code, response.Body.String())
		require.Contains(t, response.Body.String(), `"can_activate":true`)
	}
	dev, err := services.device.GetByID(t.Context(), "draft-plc")
	require.NoError(t, err)
	require.Equal(t, schema.DeviceStatusDraft, dev.Status)
	foreignLeases := map[string]uint64{}
	if foreignActivate {
		reader := &wordDraftProjectionRace{reader: services.workspace}
		reader.hook = func(projection *workspace.RuntimeProjection) {
			if len(projection.Points) == 0 || len(projection.RuleLinks) == 0 {
				return
			}
			reader.hook = nil
			require.NoError(t, services.device.ProbeAndActivate(t.Context(), "draft-plc"))
			active, e := services.device.GetByID(t.Context(), "draft-plc")
			require.NoError(t, e)
			require.NoError(t, services.runtime.UpsertDevice(t.Context(), active))
			for _, p := range projection.Points {
				foreignLeases[p.ID] = services.runtime.CapturePointLease(p.ID, p.RevisionID)
				require.NotZero(t, foreignLeases[p.ID])
			}
		}
		services.runtime.SetWorkspaceProjectionReader(reader)
	}
	body := []byte(`{"id":"draft-rule","device_id":"draft-plc","name":"BLOCK3","rule_type":"range","start_address":"40009","count":8,"data_type":"int16","prefix":"BLOCK3_","enabled":true,"skipped":[],"scale_factor":1,"scale_offset":0,"data_format":"","share_enabled":false}`)
	response = httptest.NewRecorder()
	request := httptest.NewRequest("POST", base+"/studio-v2/workspace/source-rules", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)
	var rules, points int
	require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM source_rules`).Scan(&rules))
	require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM points`).Scan(&points))
	t.Logf("status=%d rules=%d points=%d body=%s", response.Code, rules, points, response.Body.String())
	if foreignActivate {
		require.Equal(t, 503, response.Code, response.Body.String())
		require.Len(t, foreignLeases, 8)
		for id, lease := range foreignLeases {
			p, e := services.point.GetByID(t.Context(), id)
			require.NoError(t, e)
			require.Equal(t, lease, services.runtime.CapturePointLease(id, p.RevisionID), "late draft callback must preserve newer active installation")
		}
		return
	}
	require.Equal(t, 201, response.Code, response.Body.String())
	var firstID string
	require.NoError(t, db.QueryRow(`SELECT id FROM points ORDER BY id LIMIT 1`).Scan(&firstID))
	stored, err := services.point.GetByID(t.Context(), firstID)
	require.NoError(t, err)
	require.Equal(t, diagnosed, stored.Enabled, "unprobed draft preserves the existing disabled rule contract")
	require.Zero(t, services.runtime.CapturePointLease(firstID, stored.RevisionID), "saving a draft must not install or invent a lease")
	require.Empty(t, services.scheduler.PollNowContext(t.Context(), []string{firstID}))
	if diagnosed {
		response = httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest("POST", base+"/devices/draft-plc/activate", nil))
		require.Equal(t, 200, response.Code, response.Body.String())
		lease := services.runtime.CapturePointLease(firstID, stored.RevisionID)
		require.NotZero(t, lease)
		profile, ok := services.runtime.CapturePointReadProfile(firstID, lease)
		require.True(t, ok)
		require.True(t, profile.Present())
		values := services.scheduler.PollNowContext(t.Context(), []string{firstID})
		require.Len(t, values, 1)
		require.Equal(t, schema.QualityGood, values[0].Quality)
	}

}

type wordDraftProjectionRace struct {
	reader *workspace.Service
	hook   func(*workspace.RuntimeProjection)
}

func (r *wordDraftProjectionRace) RuntimeProjection(ctx context.Context) (*workspace.RuntimeProjection, error) {
	p, err := r.reader.RuntimeProjection(ctx)
	if err == nil && r.hook != nil {
		r.hook(p)
	}
	return p, err
}
