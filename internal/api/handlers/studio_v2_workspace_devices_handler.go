package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"

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

type deviceRuntimeProjectionMarker interface {
	MarkDeviceProjectionRestartRequired(ctx context.Context, deviceID string, message string)
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
	// 採集設定：間隔由設備專屬輪詢群組持有，門檻與自動調整存於設備。
	AcquisitionIntervalMs   *int `json:"acquisition_interval_ms,omitempty"`
	AcquisitionAutoInterval bool `json:"acquisition_auto_interval"`
	StaleThresholdMs        *int `json:"stale_threshold_ms"`
	FrozenThresholdMs       int  `json:"frozen_threshold_ms"`
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
				apiResponseSuccessKey: false,
				apiResponseErrorKey: gin.H{
					apiResponseMessageKey: "無法讀取 Studio V2 裝置列表，請稍後重試。",
				},
			})
			return
		}
		devices = append(devices, mapStudioV2WorkspaceDeviceResponse(savedDevice, "", ""))
	}

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    devices,
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
		if cleanupErr := h.deviceSvc.Delete(c.Request.Context(), savedDevice.ID); cleanupErr != nil {
			renderStudioV2WorkspaceDeviceError(c, fmt.Errorf("attach workspace device and cleanup failed: attach=%w cleanup=%w", err, cleanupErr))
			return
		}
		renderStudioV2WorkspaceDeviceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    mapStudioV2WorkspaceDeviceResponse(savedDevice, runtimeApplyNotRunningStatus, ""),
	})
}

func (h *StudioV2WorkspaceDevicesHandler) Update(c *gin.Context) {
	if !h.requireWorkspaceDevice(c) {
		return
	}

	var req device.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}

	// 以實際儲存前後的連線設定判斷是否需要重啟；表單每次都會帶 connection_config，
	// 只改採集設定時不得誤判為連線變更。
	previous, err := h.deviceSvc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		renderStudioV2WorkspaceDeviceError(c, err)
		return
	}
	previousProtocol, previousConnection := previous.Protocol, previous.ConnectionConfig

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

	var runtimeApplyStatus, runtimeApplyMessage string
	savedDevice, err = h.deviceSvc.ApplyCurrentRuntimeDevice(c.Request.Context(), savedDevice.ID, func(current *schema.Device) error {
		connectionChanged := current.Protocol != previousProtocol ||
			!sameTransportConnectionConfig(previousConnection, current.ConnectionConfig)
		runtimeApplyStatus, runtimeApplyMessage = h.applyRuntimeDeviceUpdate(c.Request.Context(), current, connectionChanged)
		return nil
	})
	if err != nil {
		renderStudioV2WorkspaceDeviceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    mapStudioV2WorkspaceDeviceResponse(savedDevice, runtimeApplyStatus, runtimeApplyMessage),
	})
}

func (h *StudioV2WorkspaceDevicesHandler) UpdateAvailability(c *gin.Context) {
	if !h.requireWorkspaceDevice(c) {
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
		apiResponseSuccessKey: true,
		apiResponseDataKey:    mapStudioV2WorkspaceDeviceResponse(savedDevice, "", ""),
	})
}

func (h *StudioV2WorkspaceDevicesHandler) Delete(c *gin.Context) {
	if !h.requireWorkspaceDevice(c) {
		return
	}

	deviceID := c.Param("id")
	if _, err := h.workspaceSvc.DetachDevice(c.Request.Context(), deviceID); err != nil {
		renderStudioV2WorkspaceDeviceError(c, err)
		return
	}
	if err := h.deviceSvc.Delete(c.Request.Context(), deviceID); err != nil {
		if _, rollbackErr := h.workspaceSvc.AttachDevice(c.Request.Context(), deviceID); rollbackErr != nil {
			renderStudioV2WorkspaceDeviceError(c, fmt.Errorf("delete workspace device failed and rollback failed: delete=%w rollback=%w", err, rollbackErr))
			return
		}
		renderStudioV2WorkspaceDeviceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true})
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
		apiResponseSuccessKey: true,
		apiResponseDataKey:    record,
	})
}

func (h *StudioV2WorkspaceDevicesHandler) requireWorkspaceDevice(c *gin.Context) bool {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return false
	}

	deviceID := c.Param("id")
	if slices.Contains(record.OrderedDeviceIDs, deviceID) {
		return true
	}

	c.JSON(http.StatusNotFound, gin.H{
		apiResponseSuccessKey: false,
		apiResponseErrorKey: gin.H{
			apiResponseMessageKey: "Studio V2 device not found",
		},
	})
	return false
}

func renderStudioV2WorkspaceBootstrapError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		apiResponseSuccessKey: false,
		apiResponseErrorKey: gin.H{
			apiResponseMessageKey: "無法建立 Studio V2 工作區，請檢查資料庫狀態後重試。",
		},
	})
}

func renderStudioV2WorkspaceValidationError(c *gin.Context, err error) {
	renderDeviceValidationError(c, err)
}

func renderStudioV2WorkspaceDeviceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, device.ErrValidation), errors.Is(err, workspace.ErrValidation):
		renderStudioV2WorkspaceValidationError(c, err)
	case errors.Is(err, workspace.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey: gin.H{
				apiResponseMessageKey: "Studio V2 device not found",
			},
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey: gin.H{
				apiResponseMessageKey: "Studio V2 device operation failed",
			},
		})
	}
}

func (h *StudioV2WorkspaceDevicesHandler) applyRuntimeDeviceUpdate(ctx context.Context, savedDevice *schema.Device, connectionChanged bool) (runtimeStatus, runtimeMessage string) {
	if savedDevice == nil || savedDevice.Status != schema.DeviceStatusActive {
		return runtimeApplyNotRunningStatus, ""
	}
	if connectionChanged {
		message := "device connection change requires runtime restart"
		if marker, ok := h.runtimeSync.(deviceRuntimeProjectionMarker); ok {
			marker.MarkDeviceProjectionRestartRequired(ctx, savedDevice.ID, message)
		}
		return "restart-required", message
	}
	if h.runtimeSync == nil {
		return runtimeApplyAppliedStatus, ""
	}
	if err := h.runtimeSync.UpsertDevice(ctx, savedDevice); err != nil {
		return runtimeApplyFailedStatus, "runtime device projection could not be applied"
	}
	return runtimeApplyAppliedStatus, ""
}

func mapStudioV2WorkspaceDeviceResponse(savedDevice *schema.Device, runtimeApplyStatus, runtimeApplyMessage string) studioV2WorkspaceDeviceResponse {
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

		AcquisitionIntervalMs:   savedDevice.AcquisitionIntervalMs,
		AcquisitionAutoInterval: savedDevice.AcquisitionAutoInterval,
		StaleThresholdMs:        savedDevice.StaleThresholdMs,
		FrozenThresholdMs:       savedDevice.FrozenThresholdMs,
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

// sameTransportConnectionConfig 比對連線設定時忽略 acquisition 宣告：取值優化設定檔
// 由 runtime 重新同步設備即可生效，不屬於需要重啟的傳輸連線變更。無法解析時視為不同。
func sameTransportConnectionConfig(previous, next string) bool {
	if previous == next {
		return true
	}
	var before, after map[string]json.RawMessage
	if json.Unmarshal([]byte(previous), &before) != nil || json.Unmarshal([]byte(next), &after) != nil {
		return false
	}
	delete(before, "acquisition")
	delete(after, "acquisition")
	if len(before) != len(after) {
		return false
	}
	for key, value := range before {
		other, ok := after[key]
		if !ok || !bytes.Equal(value, other) {
			return false
		}
	}
	return true
}
