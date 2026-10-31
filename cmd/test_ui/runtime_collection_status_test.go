package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	appconfig "go-gateway/internal/config"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/modbusshare"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// The only fake is the local wire peer. Diagnosed draft, source rule, mapping,
// schema confirmation, scoped recording start and status use production wiring.
func TestProductionRuntimeCollectionStatusUsesSelectedScope(t *testing.T) {
	for _, mappings := range []int{1, 8} {
		t.Run(fmt.Sprintf("%d_mapped_one_selected", mappings), func(t *testing.T) {
			f, selected := newCollectionStatusFixture(t, mappings)
			require.Eventually(t, func() bool {
				p, err := f.services.point.GetByID(t.Context(), selected)
				return err == nil && p.LastReadAt != nil && p.LastError == ""
			}, 5*time.Second, 10*time.Millisecond)
			status := collectionHTTPStatus(t, f.router)
			require.True(t, status.Running)
			require.Len(t, status.Collectors, 1)
			c := status.Collectors[0]
			require.Equal(t, 1, c.PointsTotal, "seven configured but unselected points are not collector scope")
			require.Equal(t, 1, c.PointsHealthy)
			require.Zero(t, c.PointsStale)
			require.True(t, c.Running)
			require.Equal(t, "running", c.Status)
			code, body := httpJSON(t, f.router, http.MethodGet, appconfig.Get().API.BasePath+"/datalink/points?device_id=status-plc", nil)
			require.Equal(t, http.StatusOK, code, "%+v", body)
			require.Len(t, body["data"].([]any), 8, "configuration inventory remains intact")
			require.NoError(t, f.services.runtime.Stop(t.Context()))
			stopped := collectionHTTPStatus(t, f.router)
			require.False(t, stopped.Running)
			require.False(t, stopped.Collectors[0].Running, "old good data and counters cannot revive a stopped collector")
			require.Equal(t, 1, stopped.Collectors[0].PointsTotal)
		})
	}
}

func collectionHTTPStatus(t *testing.T, router http.Handler) datalinkruntime.RuntimeStatusSnapshot {
	t.Helper()
	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, appconfig.Get().API.BasePath+"/datalink/runtime/status?device_id=status-plc", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var result struct {
		Data datalinkruntime.RuntimeStatusSnapshot `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &result))
	return result.Data
}

type collectionStatusFixture struct {
	services gatewayServices
	router   http.Handler
}

func newCollectionStatusFixture(t *testing.T, mappingCount int) (collectionStatusFixture, string) {
	t.Helper()
	_, config := wordWireServer(t, schema.ProtocolModbusTCP)
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "status.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	_, err = db.Exec(`INSERT INTO devices(id,name,description,protocol,status,connection_config,readiness_status) VALUES('status-plc','PLC','','modbus_tcp','draft',?,'{}')`, config)
	require.NoError(t, err)
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	services, err := wireGatewayServices(db, manager)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, services.runtime.Stop(context.Background())) })
	record, err := services.workspace.AttachDevice(t.Context(), "status-plc")
	require.NoError(t, err)
	router := wordProductionRouter(services, db)
	require.NoError(t, services.runtime.Start(t.Context()))
	base := appconfig.Get().API.BasePath + "/datalink"
	code, body := httpJSON(t, router, http.MethodPost, base+"/devices/status-plc/test", nil)
	require.Equal(t, http.StatusOK, code, "%+v", body)
	code, body = httpJSON(t, router, http.MethodPost, base+"/studio-v2/workspace/source-rules", map[string]any{
		"id": "status-rule", "device_id": "status-plc", "rule_type": "range", "start_address": "40001", "count": 8, "data_type": "int16", "prefix": "STATUS_", "enabled": true, "share_enabled": false,
	})
	require.Equal(t, http.StatusCreated, code, "%+v", body)
	var selected, tagID string
	for i := range mappingCount {
		code, body = httpJSON(t, router, http.MethodPost, base+"/studio-v2/workspace/mappings", map[string]any{
			"rule_id": "status-rule", "address": fmt.Sprint(40001 + i), "tag_key": fmt.Sprintf("status.selected.%d", i), "display_name": "Selected", "target_type": "int16", "scale": 1, "offset": 0, "enabled": true,
		})
		require.Contains(t, []int{http.StatusOK, http.StatusCreated}, code, "%+v", body)
		data := body["data"].(map[string]any)
		if i == 0 {
			selected, tagID = data["point_id"].(string), data["tag_id"].(string)
		}
	}

	targetFile := filepath.Join(t.TempDir(), "destination.db")
	targetConfig, err := json.Marshal(map[string]string{"dsn": targetFile})
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO database_connectors(id,name,kind,connection_config,identity_revision,status,enabled) VALUES('connector-A','Target','sqlite',?,'connector-1','ready',1)`, string(targetConfig))
	require.NoError(t, err)
	env := &outageEnv{db: db, services: services, workspID: record.ID, targets: map[string]string{"connector-A": targetFile}, schemas: map[string]string{"A": "main"}}
	target := &managedPipelineTarget{env: env, file: targetFile, namespace: "main", kind: schema.DatabaseConnectorKindSQLite}
	group := target.prepare(t, []workspace.WriteGroupMember{{DeviceID: "status-plc", PointID: selected, TagID: tagID, Required: true}}, workspace.WriteGroupRowPolicy{IntervalSeconds: 1})
	record, err = services.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	revision, _, err := modbusshare.NewSQLWorkspaceRevisionStore(db).GetRevision(t.Context(), record.ID)
	require.NoError(t, err)
	settings, err := services.modbusShare.GetSettings(t.Context())
	require.NoError(t, err)
	code, body = httpJSON(t, router, http.MethodPost, base+"/studio-v2/workspace/recording-start", workspace.RecordingStartRequest{
		RequestID: "status-start", WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision,
		WorkspaceRevision: revision, SettingsRevision: settings.SettingsRevision, DeviceIDs: []string{"status-plc"},
		Groups: []workspace.RecordingStartGroupIntent{{GroupID: group.ID, ExpectedGroupRevision: group.Revision, ExpectedConnectorRevision: "connector-1"}},
	})
	require.Equal(t, http.StatusOK, code, "%+v", body)
	require.Equal(t, "succeeded", body["data"].(map[string]any)["status"], "%+v", body)
	return collectionStatusFixture{services, router}, selected
}
