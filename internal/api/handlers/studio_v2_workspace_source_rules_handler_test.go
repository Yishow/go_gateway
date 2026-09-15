package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func seedWorkspaceRuleDevice(t *testing.T, repo *device.MemoryRepository, id string) {
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

func seedRunningWorkspaceRuleDevice(t *testing.T, repo *device.MemoryRepository, id string) {
	t.Helper()

	lastTestSuccess := true
	require.NoError(t, repo.Create(context.Background(), &schema.Device{
		ID:               id,
		Name:             id,
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":5}`,
		LastTestSuccess:  &lastTestSuccess,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}))
}

func TestStudioV2WorkspaceSourceRulesHandler_UpdateRejectsOwnershipMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	deviceRepo := device.NewMemoryRepository()
	seedWorkspaceRuleDevice(t, deviceRepo, "dev-A")
	seedWorkspaceRuleDevice(t, deviceRepo, "dev-B")

	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
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

	handler := NewStudioV2WorkspaceSourceRulesHandler(workspaceSvc, deviceSvc, ruleSvc)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/api/v1/datalink/studio-v2/workspace/source-rules/rule-A", bytes.NewBufferString(`{
		"device_id":"dev-B",
		"start_address":"40005",
		"count":1,
		"data_type":"int16",
		"naming_prefix":"A_"
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "rule-A"}}

	handler.Update(c)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Body.String(), "ownership mismatch")
	require.Contains(t, resp.Body.String(), "device_id=dev-B")
}

func TestStudioV2WorkspaceSourceRulesHandler_CreateReturnsValidationErrorMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	deviceRepo := device.NewMemoryRepository()
	seedWorkspaceRuleDevice(t, deviceRepo, "dev-A")

	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	ruleSvc := sourcerule.NewService(sourcerule.NewMemoryRepository(), deviceSvc, pointSvc, nil)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	_, err := workspaceSvc.AttachDevice(context.Background(), "dev-A")
	require.NoError(t, err)

	handler := NewStudioV2WorkspaceSourceRulesHandler(workspaceSvc, deviceSvc, ruleSvc)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/source-rules", strings.NewReader(`{
		"id":"rule-A",
		"device_id":"dev-A",
		"start_address":"",
		"count":1,
		"data_type":"int16",
		"naming_prefix":"A_",
		"enabled":true
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Create(c)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Body.String(), "start_address is required")
}

func TestStudioV2WorkspaceSourceRulesHandler_CreateRejectsUnknownDevice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	deviceRepo := device.NewMemoryRepository()
	seedWorkspaceRuleDevice(t, deviceRepo, "dev-A")

	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	ruleSvc := sourcerule.NewService(sourcerule.NewMemoryRepository(), deviceSvc, pointSvc, nil)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	_, err := workspaceSvc.AttachDevice(context.Background(), "dev-A")
	require.NoError(t, err)

	handler := NewStudioV2WorkspaceSourceRulesHandler(workspaceSvc, deviceSvc, ruleSvc)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/source-rules", strings.NewReader(`{
		"id":"rule-unknown",
		"device_id":"deleted-device-id",
		"start_address":"40001",
		"count":1,
		"data_type":"int16",
		"naming_prefix":"U_",
		"enabled":true
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Create(c)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Body.String(), "device_id does not belong to workspace")
	require.Contains(t, resp.Body.String(), "device_id=deleted-device-id")
}

func TestStudioV2WorkspaceSourceRulesHandler_CreateDefersLiveApplyWhenReadinessBlocks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	deviceRepo := device.NewMemoryRepository()
	seedRunningWorkspaceRuleDevice(t, deviceRepo, "dev-A")

	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	ruleSvc := sourcerule.NewService(sourcerule.NewMemoryRepository(), deviceSvc, pointSvc, nil)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository()).WithReadinessServices(deviceSvc, ruleSvc, nil, nil)
	_, err := workspaceSvc.AttachDevice(context.Background(), "dev-A")
	require.NoError(t, err)

	handler := NewStudioV2WorkspaceSourceRulesHandler(workspaceSvc, deviceSvc, ruleSvc)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/source-rules", strings.NewReader(`{
		"id":"rule-A",
		"device_id":"dev-A",
		"start_address":"40001",
		"count":1,
		"data_type":"int16",
		"naming_prefix":"A_",
		"enabled":true
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Create(c)

	require.Equal(t, http.StatusCreated, resp.Code, resp.Body.String())
	require.Contains(t, resp.Body.String(), `"runtime_apply_status":"deferred"`)
	require.Contains(t, resp.Body.String(), "tag-missing")
}

func TestStudioV2WorkspaceSourceRulesHandler_UpdateReturnsRuntimeReconcileOutcome(t *testing.T) {
	fixture := newWorkspaceSourceRuleReconcileFixture(t)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/api/v1/datalink/studio-v2/workspace/source-rules/rule-A", bytes.NewBufferString(`{
		"naming_prefix":"LINE_"
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "rule-A"}}

	fixture.handler.Update(c)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	require.Contains(t, resp.Body.String(), `"runtime_apply_status":"aligned"`)
	require.Contains(t, resp.Body.String(), `"runtime_apply_message":"projection aligned"`)
	require.Len(t, fixture.runtimeSync.reconcileRequests, 1)
	require.Equal(t, sourcerule.RuntimeReconcileOperationUpdate, fixture.runtimeSync.reconcileRequests[0].Operation)
	require.Equal(t, "rule-A", fixture.runtimeSync.reconcileRequests[0].Scope.RuleID)
	require.Equal(t, "dev-A", fixture.runtimeSync.reconcileRequests[0].Scope.DeviceID)
}

func TestStudioV2WorkspaceSourceRulesHandler_DeleteReturnsRuntimeReconcileOutcome(t *testing.T) {
	fixture := newWorkspaceSourceRuleReconcileFixture(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/api/v1/datalink/studio-v2/workspace/source-rules/rule-A", http.NoBody)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "rule-A"}}

	fixture.handler.Delete(c)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	require.Contains(t, resp.Body.String(), `"runtime_apply_status":"aligned"`)
	require.Contains(t, resp.Body.String(), `"runtime_apply_message":"projection aligned"`)
	require.Len(t, fixture.runtimeSync.reconcileRequests, 1)
	require.Equal(t, sourcerule.RuntimeReconcileOperationDelete, fixture.runtimeSync.reconcileRequests[0].Operation)
	require.Equal(t, "rule-A", fixture.runtimeSync.reconcileRequests[0].Scope.RuleID)
	require.Equal(t, "dev-A", fixture.runtimeSync.reconcileRequests[0].Scope.DeviceID)
}

type workspaceSourceRuleReconcileFixture struct {
	handler     *StudioV2WorkspaceSourceRulesHandler
	runtimeSync *sourceRuleRuntimeReconcileRecorder
}

func newWorkspaceSourceRuleReconcileFixture(t *testing.T) workspaceSourceRuleReconcileFixture {
	t.Helper()
	gin.SetMode(gin.TestMode)

	deviceRepo := device.NewMemoryRepository()
	seedRunningWorkspaceRuleDevice(t, deviceRepo, "dev-A")

	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	runtimeSync := &sourceRuleRuntimeReconcileRecorder{
		outcome: sourcerule.RuntimeReconcileOutcome{
			Status:  sourcerule.RuntimeReconcileStatusAligned,
			Message: "projection aligned",
		},
	}
	ruleSvc := sourcerule.NewService(sourcerule.NewMemoryRepository(), deviceSvc, pointSvc, runtimeSync)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository()).WithReadinessServices(deviceSvc, ruleSvc, nil, nil)
	_, err := workspaceSvc.AttachDevice(context.Background(), "dev-A")
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

	return workspaceSourceRuleReconcileFixture{
		handler:     NewStudioV2WorkspaceSourceRulesHandler(workspaceSvc, deviceSvc, ruleSvc),
		runtimeSync: runtimeSync,
	}
}

type sourceRuleRuntimeReconcileRecorder struct {
	outcome           sourcerule.RuntimeReconcileOutcome
	reconcileRequests []sourcerule.RuntimeReconcileRequest
}

func (r *sourceRuleRuntimeReconcileRecorder) UpsertPoint(*schema.Point) {}

func (r *sourceRuleRuntimeReconcileRecorder) RemovePoint(string) {}

func (r *sourceRuleRuntimeReconcileRecorder) ReconcileSourceRule(_ context.Context, req sourcerule.RuntimeReconcileRequest) sourcerule.RuntimeReconcileOutcome {
	r.reconcileRequests = append(r.reconcileRequests, req)
	if r.outcome.Scope.RuleID == "" {
		r.outcome.Scope = req.Scope
	}
	return r.outcome
}

func TestStudioV2WorkspaceSourceRulesHandler_CreateReturnsAddressContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	deviceRepo := device.NewMemoryRepository()
	seedWorkspaceRuleDevice(t, deviceRepo, "dev-address-create")
	deviceSvc := device.NewService(deviceRepo, nil)
	ruleSvc := sourcerule.NewService(sourcerule.NewMemoryRepository(), deviceSvc, point.NewService(point.NewMemoryRepository(), nil), nil)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	_, err := workspaceSvc.AttachDevice(context.Background(), "dev-address-create")
	require.NoError(t, err)
	handler := NewStudioV2WorkspaceSourceRulesHandler(workspaceSvc, deviceSvc, ruleSvc)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/source-rules", strings.NewReader(`{
		"id":"rule-address-create",
		"device_id":"dev-address-create",
		"start_address":"40O01",
		"count":1,
		"data_type":"int16",
		"naming_prefix":"SRC_",
		"enabled":true
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Create(c)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	for _, fragment := range []string{"rule_id=rule-address-create", "start_address=40O01", "protocol=modbus_tcp"} {
		require.Contains(t, resp.Body.String(), fragment)
	}
}

func TestStudioV2WorkspaceSourceRulesHandler_UpdateReturnsAddressContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	deviceRepo := device.NewMemoryRepository()
	seedWorkspaceRuleDevice(t, deviceRepo, "dev-address-update")
	deviceSvc := device.NewService(deviceRepo, nil)
	ruleSvc := sourcerule.NewService(sourcerule.NewMemoryRepository(), deviceSvc, point.NewService(point.NewMemoryRepository(), nil), nil)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	_, err := workspaceSvc.AttachDevice(context.Background(), "dev-address-update")
	require.NoError(t, err)
	_, err = ruleSvc.Create(context.Background(), sourcerule.CreateRuleRequest{
		ID:           "rule-address-update",
		DeviceID:     "dev-address-update",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC_",
		Enabled:      true,
	})
	require.NoError(t, err)
	handler := NewStudioV2WorkspaceSourceRulesHandler(workspaceSvc, deviceSvc, ruleSvc)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/api/v1/datalink/studio-v2/workspace/source-rules/rule-address-update", strings.NewReader(`{
		"start_address":"40O01"
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "rule-address-update"}}

	handler.Update(c)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	for _, fragment := range []string{"rule_id=rule-address-update", "start_address=40O01", "protocol=modbus_tcp"} {
		require.Contains(t, resp.Body.String(), fragment)
	}
}
