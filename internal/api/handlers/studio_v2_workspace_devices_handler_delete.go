package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"

	"go-gateway/internal/datalink/audit"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// WithWriteGroups provides the shared configuration transaction for deletion.
func (h *StudioV2WorkspaceDevicesHandler) WithWriteGroups(groups *workspace.WriteGroupService) *StudioV2WorkspaceDevicesHandler {
	h.writeGroups = groups
	return h
}

// WithAudit enables workspace device deletion audit entries.
func (h *StudioV2WorkspaceDevicesHandler) WithAudit(svc *audit.Service) *StudioV2WorkspaceDevicesHandler {
	if svc != nil {
		h.auditSvc = svc
	}
	return h
}

// DeviceDeleteResponse is the device deletion success envelope.
type DeviceDeleteResponse struct {
	Success bool `json:"success"`
}

// Delete removes a device only after all recording references are resolved.
// @Summary Delete a workspace device
// @Tags Studio V2 Devices
// @Produce json
// @Param id path string true "Device ID"
// @Success 200 {object} DeviceDeleteResponse
// @Failure 404,409,500,503 {object} APIErrorResponse
// @Router /v1/datalink/studio-v2/workspace/devices/{id} [delete]
func (h *StudioV2WorkspaceDevicesHandler) Delete(c *gin.Context) {
	if h.writeGroups == nil || h.deviceSvc == nil {
		renderWriteGroupTyped(c, http.StatusServiceUnavailable, "DEVICE_DELETE_UNAVAILABLE", "device deletion is unavailable", true, "reload")
		return
	}
	ctx := c.Request.Context()
	deviceID := c.Param("id")
	result, err := h.writeGroups.DeleteWorkspaceDevice(ctx, deviceID, h.deviceSvc.DeleteWorkspaceDeviceInTx)
	if result != nil {
		closeErr := h.deviceSvc.RetireDeletedWorkspaceDevice(deviceID)
		h.recordDeletedDevice(ctx, result)
		if closeErr != nil {
			log.Printf("close deleted workspace device connection failed")
			if err == nil {
				renderWriteGroupTyped(c, http.StatusInternalServerError, "DEVICE_DELETED_CONNECTION_CLOSE_FAILED", "device was deleted but its connection could not be closed", false, "reload")
				return
			}
		}
	}
	if errors.Is(err, workspace.ErrDeviceScopeRefresh) {
		log.Printf("refresh share scope after device deletion failed")
		renderWriteGroupTyped(c, http.StatusInternalServerError, "DEVICE_DELETED_SHARE_SCOPE_REFRESH_FAILED", "device was deleted but Share scope refresh failed", false, "reload")
		return
	}
	var references *workspace.DeviceRecordingsReferencedError
	if errors.As(err, &references) {
		c.JSON(http.StatusConflict, gin.H{"success": false, "error": gin.H{
			"code": "DEVICE_REFERENCED_BY_RECORDINGS", "message": "remove the device's recordings first", "retryable": false,
			"request_id": getOrGenerateRequestID(c), "action": "remove_recordings_first", "details": references,
		}})
		return
	}
	if err != nil {
		if renderPermanentDeleteError(c, err) {
			return
		}
		if errors.Is(err, workspace.ErrWriteGroupServiceUnavailable) {
			renderWriteGroupTyped(c, http.StatusServiceUnavailable, "DEVICE_DELETE_UNAVAILABLE", "device deletion is unavailable", true, "reload")
			return
		}
		renderStudioV2WorkspaceDeviceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *StudioV2WorkspaceDevicesHandler) recordDeletedDevice(ctx context.Context, result *workspace.DeviceDeleteResult) {
	if h.auditSvc == nil {
		log.Printf("record device deletion audit unavailable")
		return
	}
	if err := h.auditSvc.Record(ctx, audit.RecordEvent{WorkspaceID: result.WorkspaceID, EventType: audit.EventTypeDeviceDeleted, Result: audit.ResultSuccess,
		Scope: "workspace", ReferenceID: result.DeviceID, Details: map[string]any{"name": result.DeviceName, "protocol": result.Protocol}}); err != nil {
		log.Printf("record device deletion audit failed")
	}
}
