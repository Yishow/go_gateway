package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	datalinkapi "go-gateway/internal/datalink/api"
	"go-gateway/internal/datalink/connector"
	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// Production SQL repositories + real HTTP handlers. A revisioned proposal is
// not the applied workspace mapping projection; clients must join exact IDs
// through the current source-rule link projection, not invent proposal IDs.
func TestRecordingEditorRangeReadbackUsesLiveWorkspaceMappingIdentity(t *testing.T) {
	f := newWorkspaceConcurrentFixture(t, 1, 2, nil)
	router := gin.New()
	router.GET("/rules/:id/candidates", NewSourceRuleHandler(f.rules).Candidates)
	router.GET("/workspace/mappings", f.h.List)
	router.GET("/points/:id", NewPointHandler(f.h.pointSvc).Get)
	router.GET("/tags/:id", NewTagHandler(f.tags).Get)
	get := func(path string) map[string]any {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequestWithContext(f.ctx, http.MethodGet, path, nil))
		require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		var body map[string]any
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
		return body
	}
	candidates := get("/rules/diag-rule/candidates")["data"].(map[string]any)["tags"].(map[string]any)["candidates"].([]any)
	mappings := get("/workspace/mappings")["data"].([]any)
	require.Len(t, candidates, 2)
	require.Len(t, mappings, 2)
	for _, raw := range candidates {
		proposal := raw.(map[string]any)
		require.Empty(t, proposal["mapping_id"])
		require.Empty(t, proposal["tag_id"])
		require.Equal(t, "draft", proposal["status"])
		var exact map[string]any
		for _, row := range mappings {
			applied := row.(map[string]any)
			if applied["point_id"] == proposal["point_id"] && applied["address"] == proposal["address"] {
				exact = applied
			}
		}
		require.NotNil(t, exact)
		require.Equal(t, "diag-rule", exact["rule_id"])
		require.Equal(t, "diag-dev", exact["device_id"])
		require.Contains(t, f.ids, exact["id"])
		point := get("/points/" + exact["point_id"].(string))["data"].(map[string]any)
		tag := get("/tags/" + exact["tag_id"].(string))["data"].(map[string]any)
		require.Equal(t, exact["point_id"], point["id"])
		require.Equal(t, exact["device_id"], point["device_id"])
		require.Equal(t, exact["tag_id"], tag["id"])
		require.Equal(t, exact["tag_key"], tag["key"])
		require.Equal(t, "", exact["unit"], "workspace mapping serializes an empty unit")
		require.NotContains(t, tag, "unit", "Tag GET omits empty unit")
		wire, err := json.Marshal(tag)
		require.NoError(t, err)
		t.Logf("actual-empty-unit-tag-wire=%s", wire)
	}
}

func TestRecordingEditorNewUnprobedDevicePersistsDisabledRangeAndMappingsWithoutActivation(t *testing.T) {
	f := newWorkspaceConcurrentFixture(t, 1, 1, nil)
	bank := memory.NewMemoryBank(256)
	require.NoError(t, bank.WriteWord(0, 123))
	server := virtualmodbus.NewServer(bank)
	require.NoError(t, server.Start(0))
	defer server.Stop()
	diagnosticService := device.NewService(device.NewSQLRepository(f.db), connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()))
	router := gin.New()
	router.POST("/devices/:id/test", NewDeviceHandler(diagnosticService).TestConnection)
	router.POST("/devices/:id/readiness", func(c *gin.Context) {
		datalinkapi.NewDeviceHandler(diagnosticService, nil).CheckReadiness(c.Writer, c.Request, c.Param("id"))
	})
	router.PUT("/rules/:id", NewStudioV2WorkspaceSourceRulesHandler(f.h.workspaceSvc, f.h.deviceSvc, f.rules).Update)
	router.PUT("/mappings/:id", f.h.Update)
	router.GET("/points/:id", NewPointHandler(f.h.pointSvc).Get)
	router.GET("/tags/:id", NewTagHandler(f.tags).Get)
	router.POST("/devices", NewStudioV2WorkspaceDevicesHandler(f.h.workspaceSvc, f.h.deviceSvc).Create)
	router.POST("/rules", NewStudioV2WorkspaceSourceRulesHandler(f.h.workspaceSvc, f.h.deviceSvc, f.rules).Create)
	router.POST("/mappings", f.h.Create)
	call := func(method, path, body string, status int) map[string]any {
		response := httptest.NewRecorder()
		request := httptest.NewRequestWithContext(f.ctx, method, path, strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(response, request)
		require.Equal(t, status, response.Code, response.Body.String())
		var envelope map[string]any
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &envelope))
		return envelope["data"].(map[string]any)
	}
	post := func(path, body string) map[string]any { return call(http.MethodPost, path, body, http.StatusCreated) }
	savedDevice := post("/devices", fmt.Sprintf(`{"id":"offline-B","name":"Offline B","protocol":"modbus_tcp","connection_config":{"host":"127.0.0.1","port":%d,"slave_id":1,"timeout":2}}`, server.Port()))
	require.Nil(t, savedDevice["last_test_success"])
	require.Equal(t, "draft", savedDevice["status"])
	savedRule := post("/rules", `{"id":"offline-range","device_id":"offline-B","start_address":"40101","count":2,"data_type":"int16","naming_prefix":"OFFLINE","enabled":true}`)
	require.Equal(t, false, savedRule["enabled"], "server converts activation intent to its canonical effective disabled state")
	require.NotEmpty(t, savedRule["revision_id"])
	links, err := f.rules.ListLinks(f.ctx, "offline-range")
	require.NoError(t, err)
	require.Len(t, links, 2)
	savedMappings := make([]map[string]any, 0, len(links))
	for _, link := range links {
		point, err := f.h.pointSvc.GetByID(f.ctx, link.PointID)
		require.NoError(t, err)
		require.False(t, point.Enabled)
		require.Nil(t, link.MappingID)
		require.Nil(t, link.TagID)
		saved := post("/mappings", fmt.Sprintf(`{"rule_id":"offline-range","address":"%s","tag_key":"offline.%s","display_name":"Offline %s","unit":"","target_type":"int16","scale":1,"offset":0,"enabled":false}`, link.Address, link.Address, link.Address))
		savedMappings = append(savedMappings, saved)
		require.Equal(t, false, saved["enabled"])
		require.Equal(t, link.PointID, saved["point_id"])
		require.NotEmpty(t, saved["tag_id"])
		mapping, err := f.mappings.GetByID(f.ctx, saved["id"].(string))
		require.NoError(t, err)
		require.False(t, mapping.Enabled)
		pointWire := call(http.MethodGet, "/points/"+saved["point_id"].(string), "", http.StatusOK)
		tagWire := call(http.MethodGet, "/tags/"+saved["tag_id"].(string), "", http.StatusOK)
		require.Equal(t, false, pointWire["enabled"])
		require.Equal(t, saved["point_id"], pointWire["id"])
		require.Equal(t, saved["device_id"], pointWire["device_id"])
		require.Equal(t, saved["tag_id"], tagWire["id"])
		require.Equal(t, saved["tag_key"], tagWire["key"])
		require.NotContains(t, tagWire, "unit")
	}
	savedDeviceRecord, err := f.h.deviceSvc.GetByID(f.ctx, "offline-B")
	require.NoError(t, err)
	require.Nil(t, savedDeviceRecord.LastTestSuccess)
	require.EqualValues(t, "draft", savedDeviceRecord.Status)
	// Only this explicit diagnostic performs protocol I/O. No activation route is
	// registered, and the actual SQL device must remain draft throughout.
	tested := call(http.MethodPost, "/devices/offline-B/test", "", http.StatusOK)
	require.Equal(t, true, tested["success"])
	require.Equal(t, true, tested["can_activate"])
	ready := call(http.MethodPost, "/devices/offline-B/readiness", "", http.StatusOK)
	require.Equal(t, "offline-B", ready["device_id"])
	require.Equal(t, true, ready["activation_allowed"])
	updatedRule := call(http.MethodPut, "/rules/offline-range", `{"enabled":true}`, http.StatusOK)
	require.Equal(t, savedRule["id"], updatedRule["id"])
	require.Equal(t, true, updatedRule["enabled"])
	for _, saved := range savedMappings {
		body := fmt.Sprintf(`{"rule_id":"offline-range","address":"%s","tag_key":"%s","display_name":"%s","unit":"","target_type":"int16","scale":1,"offset":0,"enabled":true}`, saved["address"], saved["tag_key"], saved["display_name"])
		updated := call(http.MethodPut, "/mappings/"+saved["id"].(string), body, http.StatusOK)
		for _, key := range []string{"id", "rule_id", "device_id", "point_id", "tag_id"} {
			require.Equal(t, saved[key], updated[key], key)
		}
		require.Equal(t, true, updated["enabled"])
		pointWire := call(http.MethodGet, "/points/"+updated["point_id"].(string), "", http.StatusOK)
		tagWire := call(http.MethodGet, "/tags/"+updated["tag_id"].(string), "", http.StatusOK)
		require.Equal(t, true, pointWire["enabled"])
		require.Equal(t, saved["point_id"], pointWire["id"])
		require.Equal(t, saved["device_id"], pointWire["device_id"])
		require.Equal(t, saved["tag_id"], tagWire["id"])
		require.Equal(t, saved["tag_key"], tagWire["key"])
	}
	finalLinks, err := f.rules.ListLinks(f.ctx, "offline-range")
	require.NoError(t, err)
	require.Len(t, finalLinks, len(links))
	savedDeviceRecord, err = diagnosticService.GetByID(f.ctx, "offline-B")
	require.NoError(t, err)
	require.NotNil(t, savedDeviceRecord.LastTestSuccess)
	require.True(t, *savedDeviceRecord.LastTestSuccess)
	require.EqualValues(t, "draft", savedDeviceRecord.Status, "preparing configuration does not activate collection")
}
