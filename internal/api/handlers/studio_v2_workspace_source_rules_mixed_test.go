package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type mixedRuleHandlerFixture struct {
	db     *sql.DB
	rules  *sourcerule.Service
	points *point.Service
	router *gin.Engine
}

func newMixedRuleHandlerFixture(t *testing.T) *mixedRuleHandlerFixture {
	t.Helper()
	gin.SetMode(gin.TestMode)
	ctx := t.Context()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "configuration.db")+"?_pragma=busy_timeout(15000)")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, datalinkbase.NewMigrator().Migrate(db))
	success := true
	dr := device.NewSQLRepository(db)
	require.NoError(t, dr.Create(ctx, &schema.Device{
		ID: "mc-dev", Name: "MC", Protocol: schema.ProtocolMC3E, Status: schema.DeviceStatusActive, LastTestSuccess: &success,
		ConnectionConfig: `{"host":"127.0.0.1","port":5000,"network_no":0,"pc_no":255,"acquisition":{"version":1,"profile":"mitsubishi-q-iqr-3e-binary-v1"}}`,
		CreatedAt:        time.Now(), UpdatedAt: time.Now(),
	}))
	ds := device.NewService(dr, nil)
	ps := point.NewService(point.NewSQLRepository(db), nil)
	rs := sourcerule.NewService(sourcerule.NewSQLRepository(db), ds, ps, nil)
	ws := workspace.NewService(workspace.NewSQLRepository(db))
	_, err = ws.AttachDevice(ctx, "mc-dev")
	require.NoError(t, err)
	h := NewStudioV2WorkspaceSourceRulesHandler(ws, ds, rs)
	router := gin.New()
	router.GET("/workspace/source-rules", h.List)
	router.POST("/workspace/source-rules", h.Create)
	router.PUT("/workspace/source-rules/:id", h.Update)
	return &mixedRuleHandlerFixture{db: db, rules: rs, points: ps, router: router}
}

func (f *mixedRuleHandlerFixture) do(t *testing.T, method, path, body string) (status int, payload map[string]any) {
	t.Helper()
	response := httptest.NewRecorder()
	request := httptest.NewRequestWithContext(t.Context(), method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	f.router.ServeHTTP(response, request)
	var envelope map[string]any
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &envelope), response.Body.String())
	return response.Code, envelope
}

const mixedRuleCreateBody = `{"id":"mixed-rule","device_id":"mc-dev","naming_prefix":"MIX","enabled":true,"layout_mode":"mixed","mixed_items":[
	{"item_id":"d39","name":"Speed","address":"D39","data_type":"int16","semantic_kind":"gauge"},
	{"item_id":"d100","address":"D100","data_type":"int32","semantic_kind":"counter"},
	{"item_id":"m5","address":"M5","data_type":"bool","semantic_kind":"state"}]}`

func mixedItemsSummary(t *testing.T, raw any) [][3]string {
	t.Helper()
	items, ok := raw.([]any)
	require.True(t, ok, "mixed_items must be a JSON array: %#v", raw)
	summary := make([][3]string, 0, len(items))
	for _, entry := range items {
		item := entry.(map[string]any)
		summary = append(summary, [3]string{item["item_id"].(string), item["address"].(string), item["data_type"].(string)})
	}
	return summary
}

func TestWorkspaceSourceRuleMixedRoundTrip(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	code, envelope := f.do(t, http.MethodPost, "/workspace/source-rules", mixedRuleCreateBody)
	require.Equal(t, http.StatusCreated, code, envelope)
	want := [][3]string{{"d39", "D39", "int16"}, {"d100", "D100", "int32"}, {"m5", "M5", "bool"}}
	created := envelope["data"].(map[string]any)
	require.Equal(t, "mixed", created["layout_mode"])
	require.Equal(t, want, mixedItemsSummary(t, created["mixed_items"]))

	code, envelope = f.do(t, http.MethodGet, "/workspace/source-rules", "")
	require.Equal(t, http.StatusOK, code, envelope)
	listed := envelope["data"].([]any)
	require.Len(t, listed, 1)
	rule := listed[0].(map[string]any)
	require.Equal(t, "mixed", rule["layout_mode"])
	require.Equal(t, want, mixedItemsSummary(t, rule["mixed_items"]))
	require.Equal(t, "Speed", rule["mixed_items"].([]any)[0].(map[string]any)["name"])

	links, err := f.rules.ListLinks(t.Context(), "mixed-rule")
	require.NoError(t, err)
	require.Len(t, links, 3)
	gotPoints := map[string]schema.DataType{}
	for _, link := range links {
		pointRecord, getErr := f.points.GetByID(t.Context(), link.PointID)
		require.NoError(t, getErr)
		require.Equal(t, link.Address, pointRecord.Address)
		gotPoints[pointRecord.Address] = pointRecord.DataType
	}
	require.Equal(t, map[string]schema.DataType{"D39": schema.DataTypeInt16, "D100": schema.DataTypeInt32, "M5": schema.DataTypeBool}, gotPoints)

	update := `{"layout_mode":"mixed","mixed_items":[
		{"item_id":"d39","address":"D39","data_type":"int16","semantic_kind":"gauge"},
		{"item_id":"d200","address":"D200","data_type":"float32","semantic_kind":"gauge"}]}`
	code, envelope = f.do(t, http.MethodPut, "/workspace/source-rules/mixed-rule", update)
	require.Equal(t, http.StatusOK, code, envelope)
	require.Equal(t, [][3]string{{"d39", "D39", "int16"}, {"d200", "D200", "float32"}}, mixedItemsSummary(t, envelope["data"].(map[string]any)["mixed_items"]))
	links, err = f.rules.ListLinks(t.Context(), "mixed-rule")
	require.NoError(t, err)
	addresses := make([]string, 0, len(links))
	for _, link := range links {
		addresses = append(addresses, link.Address)
	}
	require.ElementsMatch(t, []string{"D39", "D200"}, addresses)
}

func TestWorkspaceSourceRuleMixedOverlapIsFieldError(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	before := addressConflictDatabaseSnapshot(t, f.db)
	body := `{"id":"overlap-rule","device_id":"mc-dev","naming_prefix":"MIX","enabled":true,"layout_mode":"mixed","mixed_items":[
		{"item_id":"d3","address":"D3","data_type":"uint32","semantic_kind":"gauge"},
		{"item_id":"d4","address":"D4","data_type":"int16","semantic_kind":"gauge"}]}`
	code, envelope := f.do(t, http.MethodPost, "/workspace/source-rules", body)
	require.Equal(t, http.StatusBadRequest, code, envelope)
	problem := envelope["error"].(map[string]any)
	require.Equal(t, apiValidationErrorCode, problem["code"])
	require.Equal(t, "mixed_items", problem["field"])
	require.Equal(t, []any{"d3", "d4"}, problem["item_ids"])
	require.Contains(t, problem["message"], "D3")
	require.Contains(t, problem["message"], "D4")
	require.Equal(t, before, addressConflictDatabaseSnapshot(t, f.db), "rejected mixed layout changed configuration")
}
