package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type StudioV2WorkspaceMeasurementsHandler struct {
	workspaceSvc   *workspace.Service
	deviceSvc      *device.Service
	measurementSvc *measurement.Service
}

func NewStudioV2WorkspaceMeasurementsHandler(
	workspaceSvc *workspace.Service,
	deviceSvc *device.Service,
	measurementSvc *measurement.Service,
) *StudioV2WorkspaceMeasurementsHandler {
	return &StudioV2WorkspaceMeasurementsHandler{
		workspaceSvc:   workspaceSvc,
		deviceSvc:      deviceSvc,
		measurementSvc: measurementSvc,
	}
}

func (h *StudioV2WorkspaceMeasurementsHandler) List(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	deviceID := c.Query("device_id")
	var list []measurement.MeasurementDefinition
	if deviceID != "" {
		list, err = h.measurementSvc.ListByDevice(c.Request.Context(), record.ID, deviceID)
	} else {
		list, err = h.measurementSvc.ListByWorkspace(c.Request.Context(), record.ID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MEASUREMENT_LIST_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    list,
	})
}

func (h *StudioV2WorkspaceMeasurementsHandler) ListTemplates(c *gin.Context) {
	templates := make([]measurement.MeasurementTemplate, 0, len(measurement.BuiltinTemplates))
	for _, tmpl := range measurement.BuiltinTemplates {
		templates = append(templates, tmpl)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    templates,
	})
}

type previewTemplateRequest struct {
	TemplateID  string `json:"template_id" binding:"required"`
	DeviceID    string `json:"device_id" binding:"required"`
	DeviceName  string `json:"device_name"`
	BaseAddress string `json:"base_address"`
}

func (h *StudioV2WorkspaceMeasurementsHandler) PreviewTemplate(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	var req previewTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}

	if !workspaceOwnsDevice(record, req.DeviceID) {
		renderStudioV2WorkspaceValidationError(c, fmt.Errorf("device_id does not belong to workspace: %s", req.DeviceID))
		return
	}

	devName := req.DeviceName
	if devName == "" {
		dev, err := h.deviceSvc.GetByID(c.Request.Context(), req.DeviceID)
		if err == nil && dev != nil {
			devName = dev.Name
		}
	}

	preview, err := measurement.GenerateTemplatePreview(req.TemplateID, req.DeviceID, devName, req.BaseAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "TEMPLATE_PREVIEW_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    preview,
	})
}

func (h *StudioV2WorkspaceMeasurementsHandler) ApplyTemplate(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	var preview measurement.TemplateApplyPreview
	if err := c.ShouldBindJSON(&preview); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}

	if !workspaceOwnsDevice(record, preview.DeviceID) {
		renderStudioV2WorkspaceValidationError(c, fmt.Errorf("device_id does not belong to workspace: %s", preview.DeviceID))
		return
	}

	if err := h.measurementSvc.ApplyTemplate(c.Request.Context(), record.ID, &preview); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "TEMPLATE_APPLY_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"applied": true,
			"count":   len(preview.Definitions),
		},
	})
}

func (h *StudioV2WorkspaceMeasurementsHandler) Create(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	var def measurement.MeasurementDefinition
	if err := c.ShouldBindJSON(&def); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}

	if !workspaceOwnsDevice(record, def.DeviceID) {
		renderStudioV2WorkspaceValidationError(c, fmt.Errorf("device_id does not belong to workspace: %s", def.DeviceID))
		return
	}

	def.WorkspaceID = record.ID
	if err := h.measurementSvc.CreateMeasurement(c.Request.Context(), &def); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MEASUREMENT_CREATE_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    def,
	})
}

func (h *StudioV2WorkspaceMeasurementsHandler) Update(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	id := strings.TrimSpace(c.Param("id"))
	var def measurement.MeasurementDefinition
	if err := c.ShouldBindJSON(&def); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}
	def.ID = id
	def.WorkspaceID = record.ID

	if !workspaceOwnsDevice(record, def.DeviceID) {
		renderStudioV2WorkspaceValidationError(c, fmt.Errorf("device_id does not belong to workspace: %s", def.DeviceID))
		return
	}

	epochTransitioned, err := h.measurementSvc.UpdateMeasurement(c.Request.Context(), &def)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MEASUREMENT_UPDATE_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"measurement":        def,
			"epoch_transitioned": epochTransitioned,
		},
	})
}

func (h *StudioV2WorkspaceMeasurementsHandler) Delete(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if err := h.measurementSvc.DeleteMeasurement(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MEASUREMENT_DELETE_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"deleted": true,
			"id":      id,
		},
	})
}
