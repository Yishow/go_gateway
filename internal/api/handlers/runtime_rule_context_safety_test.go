package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

func TestRuntimeRuleContextSafetyRemoteConnector(t *testing.T) {
	for _, kind := range []schema.DatabaseConnectorKind{schema.DatabaseConnectorKindSQLite, schema.DatabaseConnectorKindPostgres} {
		t.Run(string(kind), func(t *testing.T) {
			ws := workspace.NewService(workspace.NewMemoryRepository())
			_, err := ws.BindDatabaseConnector(t.Context(), "connector-C")
			require.NoError(t, err)
			h := NewRuntimeHandler(device.NewService(device.NewMemoryRepository(), nil), nil, nil, nil, nil, ws)
			h.dbConnectorReader = runtimeSetupConnectorReaderStub{record: &schema.DatabaseConnector{
				ID: "connector-C", Name: "Safe destination", Kind: kind, Status: schema.DatabaseConnectorStatusReady,
				ConnectionConfig: `{"host":"private.example.internal","port":6543,"database":"/private/fixture-secret.sqlite","user":"fixture-user","password":"fixture-password","schema":"public","table":"meter_readings","timestamp_column":"recorded_at","extra":{"token":"fixture-nested"}}`,
				LastWriteError:   "SELECT secret FROM credentials; driver fixture-driver",
			}}
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace/runtime-context", nil)
			c.Request.RemoteAddr = "203.0.113.9:4567"
			h.WorkspaceContext(c)
			require.Equal(t, 200, recorder.Code)
			for _, secret := range []string{"private.example", "6543", "fixture-secret", "fixture-user", "fixture-password", "fixture-nested", "fixture-driver", "SELECT", `"host"`, `"database"`, `"username"`, `"row_groups"`} {
				require.NotContains(t, recorder.Body.String(), secret)
			}
			require.Contains(t, recorder.Body.String(), `"table":"meter_readings"`)
			require.Contains(t, recorder.Body.String(), `"name":"Safe destination"`)
			require.Contains(t, recorder.Body.String(), `"evidence_state":"saved-only"`)
		})
	}
}

type runtimeChangingRules struct{ calls int }

func (r *runtimeChangingRules) ListByDeviceIDs(context.Context, []string) ([]*schema.SourceRule, error) {
	r.calls++
	rev := "R1"
	if r.calls > 1 {
		rev = "R2"
	}
	return []*schema.SourceRule{{ID: "rule", DeviceID: "device", RevisionID: rev, StartAddress: "40001", DataType: schema.DataTypeInt16}}, nil
}
func (*runtimeChangingRules) ListLinks(context.Context, string) ([]*schema.SourceRuleLink, error) {
	return nil, nil
}
func TestRuntimeRuleContextSafetyMixedRevision(t *testing.T) {
	ws := workspace.NewService(workspace.NewMemoryRepository())
	_, err := ws.AttachDevice(t.Context(), "device")
	require.NoError(t, err)
	repo := device.NewMemoryRepository()
	require.NoError(t, repo.Create(t.Context(), &schema.Device{ID: "device", Name: "Device", Protocol: schema.ProtocolModbusTCP}))
	h := NewRuntimeHandler(device.NewService(repo, nil), nil, nil, nil, nil, ws)
	h.sourceRuleReader = &runtimeChangingRules{}
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/runtime-context", nil)
	h.WorkspaceContext(c)
	require.Equal(t, http.StatusConflict, rec.Code, rec.Body.String())
	require.Contains(t, rec.Body.String(), "stale-reference")
	require.NotContains(t, rec.Body.String(), `"source_rules"`)
}
func TestRuntimeRuleContextSafetyHostileLabelsAndAddresses(t *testing.T) {
	for _, value := range []string{"postgres://fixture-user:fixture-password@private.internal/db", "/private/fixture.sqlite", "SELECT secret FROM users", "192.0.2.1:502"} {
		rule := mapRuntimeSetupSourceRule(&schema.SourceRule{ID: "rule", DeviceID: "device", NamingPrefix: value, StartAddress: value})
		raw, err := json.Marshal(rule)
		require.NoError(t, err)
		require.NotContains(t, string(raw), value)
	}
	rule := mapRuntimeSetupSourceRule(&schema.SourceRule{ID: "rule", DeviceID: "device", NamingPrefix: "Meter", StartAddress: "40001"})
	require.Equal(t, "40001", rule.StartAddress)
	require.Equal(t, "Meter", rule.NamingPrefix)
}

func TestRuntimeRuleContextSafetyReadinessAndTargets(t *testing.T) {
	ws := workspace.NewService(workspace.NewMemoryRepository())
	repo := device.NewMemoryRepository()
	require.NoError(t, repo.Create(t.Context(), &schema.Device{ID: "device", Name: "Device", Protocol: schema.ProtocolModbusTCP, ReadinessStatus: `{"availability_status":"unavailable","availability_reason":"driver /private/fixture.sqlite password=fixture-password"}`}))
	_, err := ws.AttachDevice(t.Context(), "device")
	require.NoError(t, err)
	devices := device.NewService(repo, nil)
	ws.WithReadinessServices(devices, nil, nil, nil)
	h := NewRuntimeHandler(devices, nil, nil, nil, nil, ws)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/runtime-context", nil)
	h.WorkspaceContext(c)
	require.NotContains(t, rec.Body.String(), "fixture")
	target := mapRuntimeSetupDatabaseTarget(&schema.DatabaseTargetMapping{ID: "target", TagID: "tag", ConnectorID: "connector", TableSchema: "/private/fixture.sqlite", TableName: "SELECT fixture FROM secrets", ColumnName: "value"})
	raw, err := json.Marshal(target)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "fixture")
	require.Contains(t, string(raw), `"column_name":"value"`)
}

func TestRuntimeRuleContextSafetyUnknownEnums(t *testing.T) {
	hostile := schema.DataType("postgres://fixture-password@internal/db")
	rule := mapRuntimeSetupSourceRule(&schema.SourceRule{ID: "rule", DeviceID: "device", DataType: hostile, TargetDataType: &hostile})
	raw, err := json.Marshal(rule)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "fixture-password")
	config := buildRuntimeDatabaseConfig("workspace", &schema.DatabaseConnector{ID: "C", Name: "Safe destination", Kind: "future-driver", ConnectionConfig: `{"schema":{"password":"fixture-password"},"table":"safe_table","future":{"dsn":"fixture-password"}}`})
	raw, err = json.Marshal(config)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "fixture-password")
	require.Contains(t, string(raw), `"kind":"unavailable"`)
}

func TestRuntimeRuleContextSafetyMappingFields(t *testing.T) {
	tagID := "tag"
	h := &RuntimeHandler{tagReader: runtimeSetupTagReaderStub{records: map[string]*schema.Tag{"tag": {ID: "tag", Key: "postgres://fixture-password@internal/db", DisplayName: "Safe Tag", Unit: "password=fixture-password", DataType: "fixture-password"}}}, mappingReader: runtimeSetupMappingReaderStub{records: map[string]*schema.Mapping{"mapping": {ID: "mapping", PointID: "point", TagID: "tag", Status: "fixture-password"}}}}
	mappingID := "mapping"
	item, err := h.mapRuntimeSetupMapping(t.Context(), &schema.SourceRuleLink{RuleID: "rule", PointID: "point", Address: "40001", TagID: &tagID, MappingID: &mappingID}, map[string]*schema.SourceRule{"rule": {ID: "rule", DeviceID: "device", DataType: schema.DataTypeInt16}})
	require.NoError(t, err)
	raw, err := json.Marshal(item)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "fixture-password")
}

type runtimeChangingConnector struct{ calls int }

func (r *runtimeChangingConnector) GetByID(context.Context, string) (*schema.DatabaseConnector, error) {
	r.calls++
	rev := "C1"
	if r.calls > 1 {
		rev = "C2"
	}
	return &schema.DatabaseConnector{ID: "C", IdentityRevision: rev, Kind: schema.DatabaseConnectorKindSQLite, ConnectionConfig: `{"table":"safe_table"}`}, nil
}
func TestRuntimeRuleContextSafetyConnectorRevisionChanges(t *testing.T) {
	ws := workspace.NewService(workspace.NewMemoryRepository())
	_, err := ws.BindDatabaseConnector(t.Context(), "C")
	require.NoError(t, err)
	h := NewRuntimeHandler(device.NewService(device.NewMemoryRepository(), nil), nil, nil, nil, nil, ws)
	h.dbConnectorReader = &runtimeChangingConnector{}
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/runtime-context", nil)
	h.WorkspaceContext(c)
	require.Equal(t, http.StatusConflict, rec.Code, rec.Body.String())
	require.Contains(t, rec.Body.String(), "stale-reference")
}

type runtimeMutatingRules struct{ rule *schema.SourceRule }

func (r *runtimeMutatingRules) ListByDeviceIDs(context.Context, []string) ([]*schema.SourceRule, error) {
	return []*schema.SourceRule{r.rule}, nil
}
func (r *runtimeMutatingRules) ListLinks(context.Context, string) ([]*schema.SourceRuleLink, error) {
	r.rule.RevisionID = "R2"
	return nil, nil
}
func TestRuntimeRuleContextSafetyMutableReadSnapshot(t *testing.T) {
	ws := workspace.NewService(workspace.NewMemoryRepository())
	_, err := ws.AttachDevice(t.Context(), "device")
	require.NoError(t, err)
	repo := device.NewMemoryRepository()
	require.NoError(t, repo.Create(t.Context(), &schema.Device{ID: "device", Name: "Device", Protocol: schema.ProtocolModbusTCP}))
	h := NewRuntimeHandler(device.NewService(repo, nil), nil, nil, nil, nil, ws)
	h.sourceRuleReader = &runtimeMutatingRules{rule: &schema.SourceRule{ID: "rule", DeviceID: "device", RevisionID: "R1"}}
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/runtime-context", nil)
	h.WorkspaceContext(c)
	require.Equal(t, http.StatusConflict, rec.Code, rec.Body.String())
	require.Contains(t, rec.Body.String(), "stale-reference")
}

func TestRuntimeRuleContextSafetyMappingIdentitiesAndTargetMode(t *testing.T) {
	h := &RuntimeHandler{}
	item, err := h.mapRuntimeSetupMapping(t.Context(), &schema.SourceRuleLink{RuleID: "postgres://fixture-password@internal/db", PointID: "/private/fixture.sqlite", Address: "40001"}, nil)
	require.NoError(t, err)
	raw, err := json.Marshal(item)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "fixture")
	target := mapRuntimeSetupDatabaseTarget(&schema.DatabaseTargetMapping{ID: "target", WriteMode: "password=fixture-password"})
	raw, err = json.Marshal(target)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "fixture")
}

func TestRuntimeRuleContextSafetyMalformedConfigUnavailable(t *testing.T) {
	config := buildRuntimeDatabaseConfig("workspace", &schema.DatabaseConnector{ID: "C", Kind: schema.DatabaseConnectorKindSQLite, ConnectionConfig: `{"table":"safe_table","write_mode":{"password":"fixture"}}`})
	require.Equal(t, "unavailable", config.WriteMode)
	require.Empty(t, config.Table)
}
