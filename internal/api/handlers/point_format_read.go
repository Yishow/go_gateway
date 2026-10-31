package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

// PointFormatReadResponse returns only scoped configuration evidence.
type PointFormatReadResponse struct {
	Success bool                    `json:"success"`
	Data    *point.FormatReadResult `json:"data"`
}

// PointFormatPoint excludes volatile readings, diagnostics and connection secrets.
type PointFormatPoint struct {
	ID             string                 `json:"id"`
	DeviceID       string                 `json:"device_id"`
	RevisionID     string                 `json:"revision_id"`
	DataFormat     string                 `json:"data_format"`
	FormatMode     schema.PointFormatMode `json:"format_mode"`
	FormatOverride string                 `json:"format_override"`
	FormatVersion  int                    `json:"format_version"`
}

type PointFormatMutationData struct {
	Point              PointFormatPoint         `json:"point"`
	Resolution         decoding.PointResolution `json:"resolution"`
	WorkspaceRevision  string                   `json:"workspace_revision"`
	RuntimeApplyStatus string                   `json:"runtime_apply_status"`
	InvalidationStatus string                   `json:"invalidation_status"`
}

type PointFormatMutationResponse struct {
	Success bool                    `json:"success"`
	Data    PointFormatMutationData `json:"data"`
}

type PointFormatErrorResponse struct {
	Success bool             `json:"success"`
	Error   PointFormatError `json:"error"`
}

type PointFormatError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func formatMutationResponse(result *point.FormatMutationResult) PointFormatMutationResponse {
	p := result.Point
	return PointFormatMutationResponse{Success: true, Data: PointFormatMutationData{
		Point:      PointFormatPoint{ID: p.ID, DeviceID: p.DeviceID, RevisionID: p.RevisionID, DataFormat: p.DataFormat, FormatMode: p.FormatMode, FormatOverride: p.FormatOverride, FormatVersion: p.FormatVersion},
		Resolution: result.Resolution, WorkspaceRevision: result.WorkspaceRevision, RuntimeApplyStatus: result.RuntimeApplyStatus, InvalidationStatus: result.InvalidationStatus,
	}}
}

// ReadFormat godoc
// @Summary Read scoped point format intent and effective policy without writes
// @Description workspace_revision is the source/Share authority, not database setup. Empty point_revision is an observed legacy revision. editable=false means source reconciliation is pending. Save must use this editing-start evidence, never silently refresh it.
// @Tags points
// @Produce json
// @Param id path string true "Point ID"
// @Param workspace_id query string true "Source workspace ID"
// @Param device_id query string true "Owning device ID"
// @Param rule_id query string false "Owning rule ID; omitted discovers the unique link"
// @Success 200 {object} PointFormatReadResponse
// @Failure 403,422,503 {object} PointFormatErrorResponse
// @Router /datalink/points/{id}/format [get]
func (h *PointHandler) ReadFormat(c *gin.Context) {
	scope := point.FormatReadScope{WorkspaceID: c.Query("workspace_id"), DeviceID: c.Query("device_id"), PointID: c.Param("id"), RuleID: c.Query("rule_id")}
	result, err := h.svc.ReadFormat(c.Request.Context(), scope)
	if err != nil {
		renderPointFormatError(c, err)
		return
	}
	c.JSON(http.StatusOK, PointFormatReadResponse{Success: true, Data: result})
}
