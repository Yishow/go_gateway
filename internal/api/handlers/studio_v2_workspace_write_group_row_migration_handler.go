package handlers

import (
	"net/http"
	"strings"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// PreviewRowGroupMigration preserves the persisted row layout for review.
// @Summary Preview legacy row-group migration differences
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param request body WriteGroupMigrationPreviewRequest true "Current workspace and persisted row-group IDs"
// @Success 200 {object} WriteGroupMigrationPreviewResponse
// @Failure 400,404,422,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/migrations/row-groups/preview [post]
func (h *StudioV2WorkspaceWriteGroupsHandler) PreviewRowGroupMigration(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var request WriteGroupMigrationPreviewRequest
	if err := c.ShouldBindJSON(&request); err != nil || strings.TrimSpace(request.WorkspaceID) == "" || !validMigrationSourceIDs(request.SourceIDs) {
		renderWriteGroupTyped(c, http.StatusBadRequest, "WRITE_GROUP_INVALID_REQUEST", "migration preview requires workspace and source identities", false, "review_request")
		return
	}
	result, err := h.groups.PreviewRowGroupMigration(c.Request.Context(), request.WorkspaceID, request.SourceIDs)
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, WriteGroupMigrationPreviewResponse{Success: true, Data: result})
}

// ReviewRowGroupMigration saves confirmed drafts without activating a writer.
// @Summary Save confirmed legacy row-group snapshot drafts
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param request body WriteGroupMigrationReviewRequest true "Preview digest, expected revisions and explicit snapshot confirmation"
// @Success 200 {object} WriteGroupMigrationReviewResponse
// @Failure 400,404,409,422,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/migrations/row-groups/review [post]
func (h *StudioV2WorkspaceWriteGroupsHandler) ReviewRowGroupMigration(c *gin.Context) {
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
	result, err := h.groups.ReviewRowGroupMigration(c.Request.Context(), workspace.WriteGroupMigrationReviewRequest{
		WorkspaceID: request.WorkspaceID, ExpectedWorkspaceRevision: *request.ExpectedWorkspaceRevision,
		ExpectedConnectorRevision: request.ExpectedConnectorRevision, ReviewDigest: request.ReviewDigest,
		SourceIDs: request.SourceIDs, ConfirmSnapshotConversion: *request.ConfirmSnapshotConversion,
	})
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, WriteGroupMigrationReviewResponse{Success: true, Data: result})
}

func validMigrationSourceIDs(ids []string) bool {
	if len(ids) == 0 {
		return false
	}
	for _, id := range ids {
		if strings.TrimSpace(id) == "" {
			return false
		}
	}
	return true
}
