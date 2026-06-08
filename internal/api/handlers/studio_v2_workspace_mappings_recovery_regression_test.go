package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestStudioV2WorkspaceMappingsHandlerRecoverWorkspaceLinkMappingRejectsAmbiguousPointMappings(t *testing.T) {
	handler, mappingRepo, mappingSvc, ruleSvc, tagSvc := newWorkspaceMappingHandlerWithRepo(t)
	ctx := context.Background()

	links, err := ruleSvc.ListLinks(ctx, "rule-A")
	require.NoError(t, err)
	require.Len(t, links, 1)

	primaryTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "line.a.first",
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

	primaryMapping, err := mappingSvc.Create(ctx, mapping.CreateMappingRequest{
		PointID:           links[0].PointID,
		TagID:             primaryTag.ID,
		Enabled:           boolPtr(true),
		TransformPipeline: buildWorkspaceMappingPipeline(schema.DataTypeInt16, schema.DataTypeFloat64, 1, 0),
	})
	require.NoError(t, err)

	shadowTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "line.a.second",
		DisplayName: "Line A Temp Shadow",
		Unit:        "C",
		DataType:    schema.DataTypeFloat64,
		Labels: map[string]string{
			ruleManagedTagLabelSource:  ruleManagedTagLabelValue,
			ruleManagedTagLabelRuleID:  "rule-A",
			ruleManagedTagLabelAddress: "40001",
		},
	})
	require.NoError(t, err)
	require.NoError(t, mappingRepo.Create(ctx, &schema.Mapping{
		ID:                "mapping-shadow",
		PointID:           links[0].PointID,
		TagID:             shadowTag.ID,
		TransformPipeline: primaryMapping.TransformPipeline,
		Status:            schema.MappingStatusActive,
		Enabled:           true,
	}))

	mappingRecord, found, changed, err := handler.recoverWorkspaceLinkMapping(ctx, links[0])
	require.Nil(t, mappingRecord)
	require.False(t, found)
	require.False(t, changed)
	require.Error(t, err)
	require.ErrorContains(t, err, "multiple persisted mappings")
	require.Nil(t, links[0].MappingID)
	require.Nil(t, links[0].TagID)
}

func TestStudioV2WorkspaceMappingsHandlerRecoverWorkspaceLinkMappingRejectsMismatchedOwnedMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, _, mappingSvc, ruleSvc, tagSvc := newWorkspaceMappingHandlerWithRepo(t)
	ctx := context.Background()

	links, err := ruleSvc.ListLinks(ctx, "rule-A")
	require.NoError(t, err)
	require.Len(t, links, 1)

	tagRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "line.b.temp",
		DisplayName: "Line B Temp",
		Unit:        "C",
		DataType:    schema.DataTypeFloat64,
		Labels: map[string]string{
			ruleManagedTagLabelSource:  ruleManagedTagLabelValue,
			ruleManagedTagLabelRuleID:  "rule-B",
			ruleManagedTagLabelAddress: "40011",
		},
	})
	require.NoError(t, err)

	_, err = mappingSvc.Create(ctx, mapping.CreateMappingRequest{
		PointID:           links[0].PointID,
		TagID:             tagRecord.ID,
		Enabled:           boolPtr(true),
		TransformPipeline: buildWorkspaceMappingPipeline(schema.DataTypeInt16, schema.DataTypeFloat64, 1, 0),
	})
	require.NoError(t, err)

	mappingRecord, found, changed, err := handler.recoverWorkspaceLinkMapping(ctx, links[0])
	require.Nil(t, mappingRecord)
	require.False(t, found)
	require.False(t, changed)
	require.Error(t, err)
	require.ErrorContains(t, err, "belong to another source rule link")
	require.Nil(t, links[0].MappingID)
	require.Nil(t, links[0].TagID)
}

func TestStudioV2WorkspaceMappingsHandlerListSkipsRecoveredMappingWhenTagIsMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, _, mappingSvc, ruleSvc, tagSvc := newWorkspaceMappingHandlerWithRepo(t)
	ctx := context.Background()

	links, err := ruleSvc.ListLinks(ctx, "rule-A")
	require.NoError(t, err)
	require.Len(t, links, 1)

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
	_, err = mappingSvc.Create(ctx, mapping.CreateMappingRequest{
		PointID:           links[0].PointID,
		TagID:             tagRecord.ID,
		Enabled:           boolPtr(true),
		TransformPipeline: buildWorkspaceMappingPipeline(schema.DataTypeInt16, schema.DataTypeFloat64, 1, 0),
	})
	require.NoError(t, err)
	require.NoError(t, tagSvc.Delete(ctx, tagRecord.ID))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace/mappings", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.List(c)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	body := decodeWorkspaceMappingBody(t, resp)
	rows := body["data"].([]any)
	require.Len(t, rows, 0)
}

func TestStudioV2WorkspaceMappingsHandlerCreateRepairsRecoveredMappingWhenTagIsMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, _, mappingSvc, ruleSvc, tagSvc := newWorkspaceMappingHandlerWithRepo(t)
	ctx := context.Background()

	links, err := ruleSvc.ListLinks(ctx, "rule-A")
	require.NoError(t, err)
	require.Len(t, links, 1)

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
	require.NoError(t, tagSvc.Delete(ctx, tagRecord.ID))

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/mappings", bytes.NewBufferString(`{
		"rule_id":"rule-A",
		"address":"40001",
		"tag_key":"line.a.repaired",
		"display_name":"Line A Repaired",
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

	require.Equal(t, http.StatusOK, createResp.Code, createResp.Body.String())
	body := decodeWorkspaceMappingBody(t, createResp)
	row := body["data"].(map[string]any)
	require.Equal(t, mappingRecord.ID, row["id"])
	require.Equal(t, "line.a.repaired", row["tag_key"])
	require.Equal(t, float64(2), row["scale"])
	require.Equal(t, float64(1), row["offset"])

	allMappings, err := mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, err)
	require.Len(t, allMappings, 1)
}

func newWorkspaceMappingHandlerWithRepo(t *testing.T) (*StudioV2WorkspaceMappingsHandler, *mapping.MemoryRepository, *mapping.Service, *sourcerule.Service, *tag.Service) {
	t.Helper()

	deviceRepo := device.NewMemoryRepository()
	seedWorkspaceMappingDevice(t, deviceRepo, "dev-A")
	seedWorkspaceMappingDevice(t, deviceRepo, "dev-B")

	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingRepo := mapping.NewMemoryRepository()
	mappingSvc := mapping.NewServiceWithTagResolver(mappingRepo, tagSvc.GetByID)
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
	return handler, mappingRepo, mappingSvc, ruleSvc, tagSvc
}