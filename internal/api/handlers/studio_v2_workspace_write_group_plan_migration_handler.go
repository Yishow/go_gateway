package handlers

import (
	"errors"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// PreviewRecordingPlanMigration explains why original plans cannot be
// automatically converted into basic periodic snapshots.
// @Summary Preview recording-plan migration eligibility
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param request body WriteGroupMigrationPreviewRequest true "Current workspace and persisted plan IDs"
// @Success 200 {object} WriteGroupMigrationPreviewResponse
// @Failure 400,404,422,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/migrations/recording-plans/preview [post]
func (h *StudioV2WorkspaceWriteGroupsHandler) PreviewRecordingPlanMigration(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var request WriteGroupMigrationPreviewRequest
	if err := c.ShouldBindJSON(&request); err != nil || strings.TrimSpace(request.WorkspaceID) == "" || !validMigrationSourceIDs(request.SourceIDs) {
		renderWriteGroupTyped(c, http.StatusBadRequest, "WRITE_GROUP_INVALID_REQUEST", "migration preview requires workspace and source identities", false, "review_request")
		return
	}
	result, err := h.groups.PreviewRecordingPlanMigration(c.Request.Context(), request.WorkspaceID, request.SourceIDs)
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, WriteGroupMigrationPreviewResponse{Success: true, Data: result})
}

// ReviewRecordingPlanMigration validates the current preview and rejects an
// unsupported conversion while preserving the original persisted plan.
// @Summary Reject unsupported recording-plan conversion after revision checks
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param request body WriteGroupMigrationReviewRequest true "Preview digest, expected revisions and explicit confirmation"
// @Failure 400,404,409,422,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/migrations/recording-plans/review [post]
func (h *StudioV2WorkspaceWriteGroupsHandler) ReviewRecordingPlanMigration(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var request WriteGroupMigrationReviewRequest
	if err := c.ShouldBindJSON(&request); err != nil || strings.TrimSpace(request.WorkspaceID) == "" ||
		request.ExpectedWorkspaceRevision == nil || strings.TrimSpace(request.ExpectedConnectorRevision) == "" ||
		strings.TrimSpace(request.ReviewDigest) == "" || !validMigrationSourceIDs(request.SourceIDs) || request.ConfirmSnapshotConversion == nil {
		renderWriteGroupTyped(c, http.StatusBadRequest, "WRITE_GROUP_INVALID_REQUEST", "migration review requires preview, revisions and explicit confirmation", false, "review_request")
		return
	}
	_, err := h.groups.ReviewRecordingPlanMigration(c.Request.Context(), workspace.WriteGroupMigrationReviewRequest{
		WorkspaceID: request.WorkspaceID, ExpectedWorkspaceRevision: *request.ExpectedWorkspaceRevision,
		ExpectedConnectorRevision: request.ExpectedConnectorRevision, ReviewDigest: request.ReviewDigest,
		SourceIDs: request.SourceIDs, ConfirmSnapshotConversion: *request.ConfirmSnapshotConversion,
	})
	if errors.Is(err, workspace.ErrWriteGroupPlanMigrationBlocked) {
		renderWriteGroupTyped(c, http.StatusUnprocessableEntity, "WRITE_GROUP_PLAN_MIGRATION_BLOCKED", "the original recording plan is preserved; configure a canonical write group for basic snapshots", false, "open_write_groups")
		return
	}
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	// This adapter has no equivalent candidate. An unexpected nil error is
	// not evidence that a conversion was saved or a writer was activated.
	renderWriteGroupTyped(c, http.StatusInternalServerError, "WRITE_GROUP_UNAVAILABLE", "recording-plan migration result is unavailable", true, "reload")
}
