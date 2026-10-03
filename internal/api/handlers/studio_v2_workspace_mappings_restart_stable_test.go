package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// A write group binds to the revision of the mapping Step 3 saved. The sync a
// restart runs must find that mapping already settled, or every group built on
// it goes stale and its data is refused until someone re-saves it.
func TestStudioV2WorkspaceMappingsHandler_SavedMappingRevisionSurvivesRestartSync(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	deviceRepo := device.NewMemoryRepository()
	seedWorkspaceMappingDevice(t, deviceRepo, "dev-A")
	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	deviceSvc := device.NewService(deviceRepo, nil)
	ruleSvc := sourcerule.NewService(sourcerule.NewMemoryRepository(), deviceSvc, pointSvc, nil)
	ruleSvc.SetTagMappingServices(tagSvc, mappingSvc)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	_, err := workspaceSvc.AttachDevice(ctx, "dev-A")
	require.NoError(t, err)
	one, zero := 1.0, 0.0
	_, err = ruleSvc.Create(ctx, sourcerule.CreateRuleRequest{
		ID: "rule-A", DeviceID: "dev-A", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt16,
		NamingPrefix: "A_", Enabled: true, ScaleMultiplier: &one, ScaleOffset: &zero,
	})
	require.NoError(t, err)
	handler := NewStudioV2WorkspaceMappingsHandler(workspaceSvc, deviceSvc, ruleSvc, pointSvc, tagSvc, mappingSvc)

	req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/mappings", strings.NewReader(`{
		"rule_id":"rule-A","address":"40001","tag_key":"line.a.temp","display_name":"Temp","unit":"","target_type":"int16","scale":1,"offset":0,"enabled":true
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	handler.Create(c)
	require.Equal(t, http.StatusCreated, resp.Code, resp.Body.String())

	saved, err := mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, err)
	require.Len(t, saved, 1)
	before, err := measurement.MappingRevision(*saved[0])
	require.NoError(t, err)

	require.NoError(t, ruleSvc.SyncDerivedPointState(ctx)) // what a restart does

	after, err := mappingSvc.GetByID(ctx, saved[0].ID)
	require.NoError(t, err)
	afterRevision, err := measurement.MappingRevision(*after)
	require.NoError(t, err)
	require.Equal(t, before, afterRevision, "restart sync changed the mapping revision a write group is bound to")
	require.Equal(t, saved[0].Status, after.Status)
}
