package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type StudioV2WorkspaceDevicesHandler struct {
	workspaceSvc *workspace.Service
	deviceSvc    *device.Service
	runtimeSync  deviceRuntimeSyncer
}

type studioV2WorkspaceDeviceOrderRequest struct {
	OrderedDeviceIDs []string `json:"ordered_device_ids"`
}

type studioV2WorkspaceDeviceResponse struct {
	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	Protocol            schema.ProtocolType `json:"protocol"`
	Status              schema.DeviceStatus `json:"status"`
	ConnectionConfig    string              `json:"connection_config"`
	LastTestAt          any                 `json:"last_test_at"`
	LastTestSuccess     *bool               `json:"last_test_success"`
	LastTestError       string              `json:"last_test_error"`
	CreatedAt           any                 `json:"created_at"`
	UpdatedAt           any                 `json:"updated_at"`
	AvailabilityStatus  string              `json:"availability_status"`
	AvailabilityReason  string              `json:"availability_reason,omitempty"`
	Running             bool                `json:"running"`
	RuntimeApplyStatus  string              `json:"runtime_apply_status,omitempty"`
	RuntimeApplyMessage string              `json:"runtime_apply_message,omitempty"`
}

type studioV2WorkspaceDeviceAvailabilityRequest struct {
	AvailabilityStatus string `json:"availability_status"`
	AvailabilityReason string `json:"availability_reason"`
}

func NewStudioV2WorkspaceDevicesHandler(workspaceSvc *workspace.Service, deviceSvc *device.Service, runtimeSync ...deviceRuntimeSyncer) *StudioV2WorkspaceDevicesHandler {
	var syncer deviceRuntimeSyncer
	if len(runtimeSync) > 0 {
		syncer = runtimeSync[0]
	}
	return &StudioV2WorkspaceDevicesHandler{
		workspaceSvc: workspaceSvc,
		deviceSvc:    deviceSvc,
		runtimeSync:  syncer,
	}
}

func (h *StudioV2WorkspaceDevicesHandler) List(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	devices := make([]any, 0, len(record.OrderedDeviceIDs))
	for _, deviceID := range record.OrderedDeviceIDs {
		savedDevice, err := h.deviceSvc.GetByID(c.Request.Context(), deviceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"message": "無法讀取 Studio V2 裝置列表，請稍後重試。",
				},
			})
			return
		}
		devices = append(devices, mapStudioV2WorkspaceDeviceResponse(savedDevice, "", ""))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    devices,
	})
}

func (h *StudioV2WorkspaceDevicesHandler) Create(c *gin.Context) {
	var req device.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}

	savedDevice, err := h.deviceSvc.Create(c.Request.Context(), req)
	if err != nil {
		renderStudioV2WorkspaceDeviceError(c, err)
		return
	}

	if _, err := h.workspaceSvc.AttachDevice(c.Request.Context(), savedDevice.ID); err != nil {
		_ = h.deviceSvc.Delete(c.Request.Context(), savedDevice.ID)
		renderStudioV2WorkspaceDeviceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    mapStudioV2WorkspaceDeviceResponse(savedDevice, "not_running", ""),
	})
}

func (h *StudioV2WorkspaceDevicesHandler) Update(c *gin.Context) {
	if _, ok := h.requireWorkspaceDevice(c); !ok {
		return
	}

	var req device.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}

	savedDevice, err := h.deviceSvc.Update(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		renderStudioV2WorkspaceDeviceError(c, err)
		return
	}
	savedDevice, err = h.deviceSvc.SetAvailability(c.Request.Context(), savedDevice.ID, device.AvailabilityStatusAvailable, "")
	if err != nil {
		renderStudioV2WorkspaceDeviceError(c, err)
		return
	}

	runtimeApplyStatus, runtimeApplyMessage := h.applyRuntimeDeviceUpdate(c.Request.Context(), savedDevice)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    mapStudioV2WorkspaceDeviceResponse(savedDevice, runtimeApplyStatus, runtimeApplyMessage),
	})
}

func (h *StudioV2WorkspaceDevicesHandler) UpdateAvailability(c *gin.Context) {
	if _, ok := h.requireWorkspaceDevice(c); !ok {
		return
	}

	var req studioV2WorkspaceDeviceAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}

	savedDevice, err := h.deviceSvc.SetAvailability(c.Request.Context(), c.Param("id"), req.AvailabilityStatus, req.AvailabilityReason)
	if err != nil {
		renderStudioV2WorkspaceDeviceError(c, err)
		return
	}
	if req.AvailabilityStatus == device.AvailabilityStatusUnavailable && h.runtimeSync != nil {
		h.runtimeSync.RemoveDevice(savedDevice.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    mapStudioV2WorkspaceDeviceResponse(savedDevice, "", ""),
	})
}

func (h *StudioV2WorkspaceDevicesHandler) Delete(c *gin.Context) {
	if _, ok := h.requireWorkspaceDevice(c); !ok {
		return
	}

	deviceID := c.Param("id")
	if _, err := h.workspaceSvc.DetachDevice(c.Request.Context(), deviceID); err != nil {
		renderStudioV2WorkspaceDeviceError(c, err)
		return
	}
	if err := h.deviceSvc.Delete(c.Request.Context(), deviceID); err != nil {
		if _, rollbackErr := h.workspaceSvc.AttachDevice(c.Request.Context(), deviceID); rollbackErr != nil {
			renderStudioV2WorkspaceDeviceError(c, fmt.Errorf("delete workspace device failed and rollback failed: delete=%w rollback=%v", err, rollbackErr))
			return
		}
		renderStudioV2WorkspaceDeviceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *StudioV2WorkspaceDevicesHandler) UpdateOrder(c *gin.Context) {
	var req studioV2WorkspaceDeviceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}

	record, err := h.workspaceSvc.ReplaceDeviceOrder(c.Request.Context(), req.OrderedDeviceIDs)
	if err != nil {
		renderStudioV2WorkspaceDeviceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    record,
	})
}

func (h *StudioV2WorkspaceDevicesHandler) requireWorkspaceDevice(c *gin.Context) (*workspace.Record, bool) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return nil, false
	}

	deviceID := c.Param("id")
	for _, attachedDeviceID := range record.OrderedDeviceIDs {
		if attachedDeviceID == deviceID {
			return record, true
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"error": gin.H{
			"message": "Studio V2 device not found",
		},
	})
	return nil, false
}

func renderStudioV2WorkspaceBootstrapError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error": gin.H{
			"message": "無法建立 Studio V2 工作區，請檢查資料庫狀態後重試。",
		},
	})
}

func renderStudioV2WorkspaceValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "validation",
			"message": err.Error(),
		},
	})
}

func renderStudioV2WorkspaceDeviceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, device.ErrValidation), errors.Is(err, workspace.ErrValidation):
		renderStudioV2WorkspaceValidationError(c, err)
	case errors.Is(err, workspace.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"message": "Studio V2 device not found",
			},
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"message": "Studio V2 device operation failed",
			},
		})
	}
}

func (h *StudioV2WorkspaceDevicesHandler) applyRuntimeDeviceUpdate(ctx context.Context, savedDevice *schema.Device) (string, string) {
	if savedDevice == nil || savedDevice.Status != schema.DeviceStatusActive {
		return "not_running", ""
	}
	if h.runtimeSync == nil {
		return "applied", ""
	}
	if err := h.runtimeSync.UpsertDevice(ctx, savedDevice); err != nil {
		return "apply_failed", err.Error()
	}
	return "applied", ""
}

func mapStudioV2WorkspaceDeviceResponse(savedDevice *schema.Device, runtimeApplyStatus string, runtimeApplyMessage string) studioV2WorkspaceDeviceResponse {
	if savedDevice == nil {
		return studioV2WorkspaceDeviceResponse{}
	}
	return studioV2WorkspaceDeviceResponse{
		ID:                  savedDevice.ID,
		Name:                savedDevice.Name,
		Description:         savedDevice.Description,
		Protocol:            savedDevice.Protocol,
		Status:              savedDevice.Status,
		ConnectionConfig:    savedDevice.ConnectionConfig,
		LastTestAt:          savedDevice.LastTestAt,
		LastTestSuccess:     savedDevice.LastTestSuccess,
		LastTestError:       savedDevice.LastTestError,
		CreatedAt:           savedDevice.CreatedAt,
		UpdatedAt:           savedDevice.UpdatedAt,
		AvailabilityStatus:  availabilityStatusOf(savedDevice),
		AvailabilityReason:  availabilityReasonOf(savedDevice),
		Running:             isWorkspaceDeviceRunning(savedDevice),
		RuntimeApplyStatus:  runtimeApplyStatus,
		RuntimeApplyMessage: runtimeApplyMessage,
	}
}

func availabilityStatusOf(savedDevice *schema.Device) string {
	status, _ := device.AvailabilityOf(savedDevice)
	return status
}

func availabilityReasonOf(savedDevice *schema.Device) string {
	_, reason := device.AvailabilityOf(savedDevice)
	return reason
}

func isWorkspaceDeviceRunning(savedDevice *schema.Device) bool {
	if savedDevice == nil || savedDevice.Status != schema.DeviceStatusActive {
		return false
	}
	status, _ := device.AvailabilityOf(savedDevice)
	return status == device.AvailabilityStatusAvailable
}
