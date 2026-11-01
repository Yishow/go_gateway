package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"slices"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type acquisitionPreviewQuery struct{ workspaceID, deviceID, workspaceRevision, previewRevision string }

func parseAcquisitionPreviewQuery(raw string) (acquisitionPreviewQuery, error) {
	q, err := url.ParseQuery(raw)
	bad := runtimeEvidenceError(runtimeview.ReasonUnavailable)
	if err != nil {
		return acquisitionPreviewQuery{}, bad
	}
	allowed := map[string]bool{"workspace_id": true, "device_id": true, expectedWorkspaceRevisionField: true, "expected_preview_revision": true}
	for key, values := range q {
		if !allowed[key] || len(values) != 1 || (values[0] == "" && key != expectedWorkspaceRevisionField) || (values[0] != "" && !runtimeDeliveryIDs(values[0])) {
			return acquisitionPreviewQuery{}, bad
		}
	}
	for _, key := range []string{"workspace_id", "device_id", expectedWorkspaceRevisionField} {
		if _, ok := q[key]; !ok {
			return acquisitionPreviewQuery{}, bad
		}
	}
	return acquisitionPreviewQuery{q.Get("workspace_id"), q.Get("device_id"), q.Get(expectedWorkspaceRevisionField), q.Get("expected_preview_revision")}, nil
}

// AcquisitionPreviewReader has no mutation, probe, or installation methods.
type AcquisitionPreviewReader interface {
	ReadAcquisitionPreview(context.Context, string, string, string) (collector.AcquisitionPreview, func() bool, error)
}

// StudioV2AcquisitionPreviewHandler scopes the native read before accessing it.
type StudioV2AcquisitionPreviewHandler struct {
	workspace *workspace.Service
	reader    AcquisitionPreviewReader
}

// NewStudioV2AcquisitionPreviewHandler wires existing native installed authority.
func NewStudioV2AcquisitionPreviewHandler(svc *workspace.Service, reader AcquisitionPreviewReader) *StudioV2AcquisitionPreviewHandler {
	return &StudioV2AcquisitionPreviewHandler{svc, reader}
}

// AcquisitionPreviewResponse is the scoped read-only API envelope.
type AcquisitionPreviewResponse struct {
	Success bool                         `json:"success"`
	Data    collector.AcquisitionPreview `json:"data"`
}

// Get serves a coherent, revision-bound installed acquisition preview.
// @Summary Read the installed acquisition plan and matching latest scheduled scan
// @Tags Datalink Studio v2
// @Produce json
// @Param workspace_id query string true "Exact owning workspace"
// @Param device_id query string true "Exact workspace device"
// @Param expected_workspace_revision query string true "Required presence; empty is valid only for a legacy-empty workspace revision"
// @Param expected_preview_revision query string false "Omit on first coherent read; if provided, must be nonempty and match the current composite head"
// @Success 200 {object} AcquisitionPreviewResponse
// @Failure 400 {object} TypedAPIErrorEnvelope "Malformed, duplicate, unknown or missing query input"
// @Failure 403 {object} TypedAPIErrorEnvelope "Foreign workspace or device"
// @Failure 409 {object} TypedAPIErrorEnvelope "Stale revision or capture changed"
// @Failure 503 {object} TypedAPIErrorEnvelope "Scoped evidence unavailable"
// @Router /datalink/studio-v2/workspace/acquisition-plan-preview [get]
func (h *StudioV2AcquisitionPreviewHandler) Get(c *gin.Context) {
	if h == nil || h.workspace == nil {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonUnavailable))
		return
	}
	query, err := parseAcquisitionPreviewQuery(c.Request.URL.RawQuery)
	if err != nil {
		renderRuntimeHandlerError(c, http.StatusBadRequest, string(runtimeview.ReasonUnavailable), false)
		return
	}
	record, err := h.workspace.ReadSnapshot(c.Request.Context())
	if err != nil || record == nil {
		renderRuntimeEvidenceError(c, studioLocatorReadError(err))
		return
	}
	if record.ID != query.workspaceID || !slices.Contains(record.OrderedDeviceIDs, query.deviceID) {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonForeign))
		return
	}
	if record.DatabaseSetupRevision != query.workspaceRevision {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonStale))
		return
	}
	unavailable := collector.AcquisitionPreview{Status: "unavailable", Reason: "installation_unavailable"}
	view := unavailable
	var verify func() bool
	var before *workspace.RuntimeProjection
	var sourceHead string
	if h.reader != nil {
		before, err = h.workspace.ReadRuntimeProjection(c.Request.Context())
		if err == nil && before != nil && before.WorkspaceID == query.workspaceID {
			graph, ok := before.AcquisitionSourceGraph(query.deviceID)
			if ok {
				sourceHead = graph.IdentityRevision()
				view, verify, err = h.reader.ReadAcquisitionPreview(c.Request.Context(), query.workspaceID, query.deviceID, sourceHead)
			}
		}
	}
	if errors.Is(err, collector.ErrAcquisitionPreviewStale) {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonStale))
		return
	}
	if err != nil {
		view = unavailable
	}
	if view.Status == runtimeEvidenceAvailable {
		after, readErr := h.workspace.ReadRuntimeProjection(c.Request.Context())
		if readErr != nil || after == nil {
			renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonUnavailable))
			return
		}
		graph, ok := after.AcquisitionSourceGraph(query.deviceID)
		if !ok || after.WorkspaceID != query.workspaceID || graph.IdentityRevision() != sourceHead || !previewCadencesMatch(before, view) || !previewCadencesMatch(after, view) || verify == nil || !verify() {
			renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonStale))
			return
		}
		revision, hashErr := measurement.HashStringTuple(query.workspaceID, query.workspaceRevision, view.PreviewRevision)
		if hashErr != nil {
			renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonUnavailable))
			return
		}
		view.PreviewRevision = revision
	}
	afterRecord, err := h.workspace.ReadSnapshot(c.Request.Context())
	if err != nil || afterRecord == nil {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonUnavailable))
		return
	}
	if !sameStudioLocatorRecord(record, afterRecord) {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonStale))
		return
	}
	if query.previewRevision != "" && (view.Status != runtimeEvidenceAvailable || query.previewRevision != view.PreviewRevision) {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonStale))
		return
	}
	if view.Status == runtimeEvidenceAvailable && (verify == nil || !verify()) {
		renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonStale))
		return
	}
	view.WorkspaceID, view.DeviceID, view.WorkspaceRevision = query.workspaceID, query.deviceID, query.workspaceRevision
	c.JSON(http.StatusOK, AcquisitionPreviewResponse{Success: true, Data: view})
}
func previewCadencesMatch(projection *workspace.RuntimeProjection, view collector.AcquisitionPreview) bool {
	for _, group := range view.Groups {
		found := false
		for _, source := range projection.PollingGroups {
			if source != nil && source.ID == group.PollingGroupID && source.Enabled && int64(source.IntervalMs) == group.PollIntervalMS {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
