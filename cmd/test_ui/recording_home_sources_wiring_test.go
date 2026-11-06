package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"go-gateway/internal/api"
	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestProductionRecordingHomeSourcesReadOnlyGraph(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "configuration.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	services, err := wireGatewayServices(db, connector.NewConnectionManager(connector.ConnectionManagerConfig{}))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, services.runtime.Stop(context.Background())) })
	for _, q := range []string{
		`INSERT INTO devices(id,name,description,protocol,status,connection_config,readiness_status) VALUES('D1','PLC','','modbus_tcp','draft','{"host":"private.host","password":"private-secret"}','{}')`,
		`INSERT INTO points(id,device_id,name,address,function,data_type,mode,enabled,last_value) VALUES('P1','D1','Temperature','40001','FC03','float32','read',1,'private-value')`,
		`INSERT INTO tags(id,key,key_lower,display_name,data_type,status,labels) VALUES('T1','temperature','temperature',' Temperature ','float32','active','{"source":"source-rule","source_rule_id":"R1","source_rule_address":"40001"}')`,
		`INSERT INTO mappings(id,point_id,tag_id,transform_pipeline,status,enabled) VALUES('M1','P1','T1','[]','active',1)`,
		`INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix,enabled,origin,revision_id,skipped_addresses) VALUES('R1','D1','40001',1,'float32','Temperature',1,'manual','revision-1','[]')`,
		`INSERT INTO source_rule_links(id,rule_id,address,point_id,tag_id,mapping_id) VALUES('L1','R1','40001','P1','T1','M1')`,
	} {
		_, err = db.ExecContext(t.Context(), q)
		require.NoError(t, err)
	}
	record, err := services.workspace.AttachDevice(t.Context(), "D1")
	require.NoError(t, err)
	_, err = services.device.GetByID(t.Context(), "D1")
	require.NoError(t, err)
	_, err = services.point.List(t.Context(), point.ListFilter{})
	require.NoError(t, err)
	_, err = services.sourceRule.ListByDeviceIDs(t.Context(), []string{"D1"})
	require.NoError(t, err)
	_, err = services.sourceRule.ListLinks(t.Context(), "R1")
	require.NoError(t, err)
	_, err = services.mapping.List(t.Context(), mapping.ListFilter{})
	require.NoError(t, err)
	_, err = services.mapping.GetByID(t.Context(), "M1")
	require.NoError(t, err)
	_, err = services.tag.GetByID(t.Context(), "T1")
	require.NoError(t, err)
	// These are the exact six source authorities supplied by main.go to NewRouter.
	router := api.NewRouter(&api.DatalinkServices{Workspace: services.workspace, Device: services.device, SourceRule: services.sourceRule, Point: services.point, Mapping: services.mapping, Tag: services.tag})
	var before int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT total_changes()`).Scan(&before))
	var digest string
	for range 3 {
		out := httptest.NewRecorder()
		router.ServeHTTP(out, httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/recording-home/sources?workspace_id="+url.QueryEscape(record.ID), http.NoBody))
		require.Equal(t, http.StatusOK, out.Code, out.Body.String())
		var body handlers.RecordingHomeSourcesEnvelope
		require.NoError(t, json.Unmarshal(out.Body.Bytes(), &body))
		require.True(t, body.Success)
		require.Equal(t, "available", body.Data.State)
		require.Equal(t, []workspace.RecordingHomeCandidate{{Key: "P1|T1", DeviceID: "D1", PointID: "P1", TagID: "T1", TagKey: "temperature", Label: "Temperature", DeviceName: "PLC", Address: "40001", TargetType: "float32"}}, body.Data.Candidates)
		require.Len(t, body.Data.Devices, 1)
		require.Empty(t, body.Data.Excluded)
		require.Len(t, body.Data.SourceDigest, 64)
		if digest != "" {
			require.Equal(t, digest, body.Data.SourceDigest)
		}
		digest = body.Data.SourceDigest
		for _, secret := range []string{"private.host", "private-secret", "private-value", "connection_config", "last_value", "transform_pipeline"} {
			require.NotContains(t, out.Body.String(), secret)
		}
	}
	var after int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT total_changes()`).Scan(&after))
	require.Equal(t, before, after, "GET must not mutate any SQLite table")
	current, err := services.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	require.Equal(t, record, current)
	// A stale link remains unknown on this route; GET must not invoke legacy repair.
	_, err = db.ExecContext(t.Context(), `UPDATE source_rule_links SET mapping_id=NULL WHERE id='L1'`)
	require.NoError(t, err)
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT total_changes()`).Scan(&before))
	out := httptest.NewRecorder()
	router.ServeHTTP(out, httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/recording-home/sources?workspace_id="+url.QueryEscape(record.ID), http.NoBody))
	require.Equal(t, http.StatusOK, out.Code, out.Body.String())
	var unknown struct {
		Data map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(out.Body.Bytes(), &unknown))
	require.Equal(t, "unknown", unknown.Data["state"])
	require.Equal(t, "source-inconsistent", unknown.Data["reason_code"])
	for _, field := range []string{"devices", "candidates", "excluded", "source_digest"} {
		require.NotContains(t, unknown.Data, field)
	}
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT total_changes()`).Scan(&after))
	require.Equal(t, before, after)

}
