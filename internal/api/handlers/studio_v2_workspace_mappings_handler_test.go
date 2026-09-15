package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func decodeWorkspaceMappingBody(t *testing.T, resp *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	return body
}

func seedWorkspaceMappingDevice(t *testing.T, repo *device.MemoryRepository, id string) {
	t.Helper()

	require.NoError(t, repo.Create(context.Background(), &schema.Device{
		ID:               id,
		Name:             id,
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":5}`,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}))
}

func newWorkspaceMappingHandler(t *testing.T) (
	handlerResult *StudioV2WorkspaceMappingsHandler,
	mappingService *mapping.Service,
	ruleService *sourcerule.Service,
	tagService *tag.Service,
) {
	t.Helper()

	deviceRepo := device.NewMemoryRepository()
	seedWorkspaceMappingDevice(t, deviceRepo, "dev-A")
	seedWorkspaceMappingDevice(t, deviceRepo, "dev-B")

	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	deviceSvc := device.NewService(deviceRepo, nil)
	ruleSvc := sourcerule.NewService(sourcerule.NewMemoryRepository(), deviceSvc, pointSvc, nil)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	_, err := workspaceSvc.AttachDevice(context.Background(), "dev-A")
	require.NoError(t, err)
	_, err = workspaceSvc.AttachDevice(context.Background(), "dev-B")
	require.NoError(t, err)

	_, err = ruleSvc.Create(context.Background(), sourcerule.CreateRuleRequest{
		ID:           "rule-A",
		DeviceID:     "dev-A",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "A_",
		Enabled:      true,
	})
	require.NoError(t, err)
	_, err = ruleSvc.Create(context.Background(), sourcerule.CreateRuleRequest{
		ID:           "rule-B",
		DeviceID:     "dev-B",
		StartAddress: "40011",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "B_",
		Enabled:      true,
	})
	require.NoError(t, err)

	handler := NewStudioV2WorkspaceMappingsHandler(workspaceSvc, deviceSvc, ruleSvc, pointSvc, tagSvc, mappingSvc)
	return handler, mappingSvc, ruleSvc, tagSvc
}

func TestStudioV2WorkspaceMappingsHandler_CreateReturnsValidationErrorMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, _, _, _ := newWorkspaceMappingHandler(t)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/mappings", strings.NewReader(`{
		"rule_id":"rule-A",
		"address":"40001",
		"tag_key":"",
		"display_name":"Line A Temp",
		"unit":"C",
		"target_type":"float64",
		"scale":1,
		"offset":0,
		"enabled":true
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Create(c)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Body.String(), "tag_key is required")
}

func TestStudioV2WorkspaceMappingsHandler_UpdateRejectsOwnershipMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, _, _, _ := newWorkspaceMappingHandler(t)

	createReq := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/mappings", strings.NewReader(`{
		"rule_id":"rule-A",
		"address":"40001",
		"tag_key":"line.a.temp",
		"display_name":"Line A Temp",
		"unit":"C",
		"target_type":"float64",
		"scale":1,
		"offset":0,
		"enabled":true
	}`))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	createCtx, _ := gin.CreateTestContext(createResp)
	createCtx.Request = createReq
	handler.Create(createCtx)
	require.Equal(t, http.StatusCreated, createResp.Code)

	createBody := decodeWorkspaceMappingBody(t, createResp)
	mappingID := createBody["data"].(map[string]any)["id"].(string)

	updateReq := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/api/v1/datalink/studio-v2/workspace/mappings/"+mappingID, bytes.NewBufferString(`{
		"rule_id":"rule-B",
		"address":"40011",
		"tag_key":"line.b.temp",
		"display_name":"Line B Temp",
		"unit":"C",
		"target_type":"float64",
		"scale":1,
		"offset":0,
		"enabled":true
	}`))
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp := httptest.NewRecorder()
	updateCtx, _ := gin.CreateTestContext(updateResp)
	updateCtx.Request = updateReq
	updateCtx.Params = gin.Params{{Key: "id", Value: mappingID}}

	handler.Update(updateCtx)

	require.Equal(t, http.StatusBadRequest, updateResp.Code)
	require.Contains(t, updateResp.Body.String(), "ownership mismatch")
}

func TestStudioV2WorkspaceMappingsHandler_CreateRejectsExistingUnmanagedTagKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, _, _, tagSvc := newWorkspaceMappingHandler(t)

	_, err := tagSvc.Create(context.Background(), tag.CreateTagRequest{
		Key:         "line.a.temp",
		DisplayName: "shared tag",
		DataType:    schema.DataTypeFloat64,
		Unit:        "C",
	})
	require.NoError(t, err)

	createReq := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/mappings", strings.NewReader(`{
		"rule_id":"rule-A",
		"address":"40001",
		"tag_key":"line.a.temp",
		"display_name":"Line A Temp",
		"unit":"C",
		"target_type":"float64",
		"scale":1,
		"offset":0,
		"enabled":true
	}`))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	createCtx, _ := gin.CreateTestContext(createResp)
	createCtx.Request = createReq

	handler.Create(createCtx)

	require.Equal(t, http.StatusBadRequest, createResp.Code)
	require.Contains(t, createResp.Body.String(), "ownership mismatch")
}

func TestStudioV2WorkspaceMappingsHandler_ListRecoversExistingPointMappingWhenRuleLinkLost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, mappingSvc, ruleSvc, tagSvc := newWorkspaceMappingHandler(t)
	ctx := context.Background()
	links, err := ruleSvc.ListLinks(ctx, "rule-A")
	require.NoError(t, err)
	require.Len(t, links, 1)
	require.Nil(t, links[0].MappingID)

	tagRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "line.a.saved",
		DisplayName: "Line A Temp",
		Unit:        "C",
		DataType:    schema.DataTypeFloat64,
		Labels: map[string]string{
			ruleManagedTagLabelSource:  ruleManagedTagLabelValue,
			ruleManagedTagLabelRuleID:  "rule-A",
			ruleManagedTagLabelAddress: "40001",
		},
	})
	require.NoError(t, err)
	mappingRecord, err := mappingSvc.Create(ctx, mapping.CreateMappingRequest{
		PointID:           links[0].PointID,
		TagID:             tagRecord.ID,
		Enabled:           boolPtr(true),
		TransformPipeline: buildWorkspaceMappingPipeline(schema.DataTypeInt16, schema.DataTypeFloat64, 1, 0),
	})
	require.NoError(t, err)

	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/mappings", http.NoBody)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.List(c)

	require.Equal(t, http.StatusOK, resp.Code)
	body := decodeWorkspaceMappingBody(t, resp)
	rows := body["data"].([]any)
	require.Len(t, rows, 1)
	row := rows[0].(map[string]any)
	require.Equal(t, mappingRecord.ID, row["id"])
	require.Equal(t, "line.a.saved", row["tag_key"])

	links, err = ruleSvc.ListLinks(ctx, "rule-A")
	require.NoError(t, err)
	require.NotNil(t, links[0].MappingID)
	require.Equal(t, mappingRecord.ID, *links[0].MappingID)
	require.NotNil(t, links[0].TagID)
	require.Equal(t, tagRecord.ID, *links[0].TagID)
}

func TestStudioV2WorkspaceMappingsHandler_CreateAdoptsExistingPointMappingWhenRuleLinkLost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, mappingSvc, ruleSvc, tagSvc := newWorkspaceMappingHandler(t)
	ctx := context.Background()
	links, err := ruleSvc.ListLinks(ctx, "rule-A")
	require.NoError(t, err)
	require.Len(t, links, 1)
	require.Nil(t, links[0].MappingID)

	tagRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "line.a.saved",
		DisplayName: "Line A Temp",
		Unit:        "C",
		DataType:    schema.DataTypeFloat64,
		Labels: map[string]string{
			ruleManagedTagLabelSource:  ruleManagedTagLabelValue,
			ruleManagedTagLabelRuleID:  "rule-A",
			ruleManagedTagLabelAddress: "40001",
		},
	})
	require.NoError(t, err)
	mappingRecord, err := mappingSvc.Create(ctx, mapping.CreateMappingRequest{
		PointID:           links[0].PointID,
		TagID:             tagRecord.ID,
		Enabled:           boolPtr(true),
		TransformPipeline: buildWorkspaceMappingPipeline(schema.DataTypeInt16, schema.DataTypeFloat64, 1, 0),
	})
	require.NoError(t, err)

	createReq := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/mappings", strings.NewReader(`{
		"rule_id":"rule-A",
		"address":"40001",
		"tag_key":"line.a.saved",
		"display_name":"Line A Temp",
		"unit":"C",
		"target_type":"float64",
		"scale":2,
		"offset":1,
		"enabled":true
	}`))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	createCtx, _ := gin.CreateTestContext(createResp)
	createCtx.Request = createReq

	handler.Create(createCtx)

	require.Equal(t, http.StatusOK, createResp.Code)
	body := decodeWorkspaceMappingBody(t, createResp)
	row := body["data"].(map[string]any)
	require.Equal(t, mappingRecord.ID, row["id"])
	require.Equal(t, float64(2), row["scale"])
	require.Equal(t, float64(1), row["offset"])

	allMappings, err := mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, err)
	require.Len(t, allMappings, 1)

	links, err = ruleSvc.ListLinks(ctx, "rule-A")
	require.NoError(t, err)
	require.NotNil(t, links[0].MappingID)
	require.Equal(t, mappingRecord.ID, *links[0].MappingID)
}
