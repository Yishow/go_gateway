package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/runtimeview"
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

// List returns saved measurement definitions without changing their revisions.
// @Summary Read workspace measurements
// @Tags Studio V2 Recording Plans
// @Produce json
// @Param device_id query string false "Optional owned device filter"
// @Param locator_workspace_id query string false "Locator-origin workspace identity; requires revision pair"
// @Param locator_workspace_revision query string false "Exact observed revision; present empty only matches legacy empty"
// @Success 200 {array} measurement.MeasurementDefinition
// @Failure 400 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 409 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Failure 503 {object} APIErrorResponse
// @Router /v1/datalink/studio-v2/workspace/measurements [get]
func (h *StudioV2WorkspaceMeasurementsHandler) List(c *gin.Context) {
	record, locatorScope, ok := beginStudioLocatorNativeRead(c, h.workspaceSvc, func() { renderStudioV2WorkspaceBootstrapError(c) })
	if !ok {
		return
	}
	var err error

	deviceID := c.Query("device_id")
	if locatorScope != nil && deviceID != "" && (!runtimeDeliveryIDs(deviceID) || !workspaceOwnsDevice(record, deviceID)) {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonForeign))
		return
	}
	var list []measurement.MeasurementDefinition
	if deviceID != "" {
		list, err = h.measurementSvc.ListByDevice(c.Request.Context(), record.ID, deviceID)
	} else {
		list, err = h.measurementSvc.ListByWorkspace(c.Request.Context(), record.ID)
	}

	if !finishStudioLocatorNativeRead(c, h.workspaceSvc, locatorScope) {
		return
	}
	if err != nil && locatorScope != nil {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonUnavailable))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey: gin.H{
				apiResponseCodeKey:    "MEASUREMENT_LIST_FAILED",
				apiResponseMessageKey: err.Error(),
			},
		})
		return
	}

	if locatorScope != nil {
		for _, definition := range list {
			if definition.WorkspaceID != record.ID || (deviceID != "" && definition.DeviceID != deviceID) {
				renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonForeign))
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    list,
	})
}

func (h *StudioV2WorkspaceMeasurementsHandler) ListTemplates(c *gin.Context) {
	templates := make([]measurement.MeasurementTemplate, 0, len(measurement.BuiltinTemplates))
	for _, tmpl := range measurement.BuiltinTemplates {
		templates = append(templates, tmpl)
	}

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    templates,
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
			apiResponseSuccessKey: false,
			apiResponseErrorKey: gin.H{
				apiResponseCodeKey:    "TEMPLATE_PREVIEW_FAILED",
				apiResponseMessageKey: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    preview,
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
			apiResponseSuccessKey: false,
			apiResponseErrorKey: gin.H{
				apiResponseCodeKey:    "TEMPLATE_APPLY_FAILED",
				apiResponseMessageKey: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey: gin.H{
			apiResponseAppliedKey: true,
			apiResponseCountKey:   len(preview.Definitions),
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
			apiResponseSuccessKey: false,
			apiResponseErrorKey: gin.H{
				apiResponseCodeKey:    "MEASUREMENT_CREATE_FAILED",
				apiResponseMessageKey: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    def,
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
			apiResponseSuccessKey: false,
			apiResponseErrorKey: gin.H{
				apiResponseCodeKey:    "MEASUREMENT_UPDATE_FAILED",
				apiResponseMessageKey: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey: gin.H{
			"measurement":        def,
			"epoch_transitioned": epochTransitioned,
		},
	})
}

func (h *StudioV2WorkspaceMeasurementsHandler) Delete(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if err := h.measurementSvc.DeleteMeasurement(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey: gin.H{
				apiResponseCodeKey:    "MEASUREMENT_DELETE_FAILED",
				apiResponseMessageKey: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey: gin.H{
			apiResponseDeletedKey: true,
			"id":                  id,
		},
	})
}
