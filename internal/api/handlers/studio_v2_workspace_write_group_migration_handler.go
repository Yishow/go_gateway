package handlers

import (
	"net/http"
	"strings"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// WriteGroupMigrationPreviewRequest selects legacy mappings for a read-only
// review in the current workspace. It grants no mutation or activation.
type WriteGroupMigrationPreviewRequest struct {
	WorkspaceID string   `json:"workspace_id"`
	SourceIDs   []string `json:"source_ids"`
}

// WriteGroupMigrationPreviewResponse contains current revisions and explicit
// differences; a preview never represents an applied migration.
type WriteGroupMigrationPreviewResponse struct {
	Success bool                                  `json:"success"`
	Data    *workspace.WriteGroupMigrationPreview `json:"data"`
}

// PreviewSingleMappingMigration reviews legacy output intent without saving
// groups, changing the source mapping or opening an external database.
// @Summary Preview legacy single-mapping migration differences
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param request body WriteGroupMigrationPreviewRequest true "Current workspace and legacy mapping IDs"
// @Success 200 {object} WriteGroupMigrationPreviewResponse
// @Failure 400,404,422,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/migrations/single-mappings/preview [post]
func (h *StudioV2WorkspaceWriteGroupsHandler) PreviewSingleMappingMigration(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var request WriteGroupMigrationPreviewRequest
	if err := c.ShouldBindJSON(&request); err != nil || strings.TrimSpace(request.WorkspaceID) == "" || len(request.SourceIDs) == 0 {
		renderWriteGroupTyped(c, http.StatusBadRequest, "WRITE_GROUP_INVALID_REQUEST", "migration preview requires workspace and source identities", false, "review_request")
		return
	}
	for _, id := range request.SourceIDs {
		if strings.TrimSpace(id) == "" {
			renderWriteGroupTyped(c, http.StatusBadRequest, "WRITE_GROUP_INVALID_REQUEST", "migration preview source identity is missing", false, "review_request")
			return
		}
	}
	result, err := h.groups.PreviewSingleMappingMigration(c.Request.Context(), request.WorkspaceID, request.SourceIDs)
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, WriteGroupMigrationPreviewResponse{Success: true, Data: result})
}

// WriteGroupMigrationReviewRequest binds explicit confirmation to a previously
// displayed migration preview. Missing CAS fields differ from empty revisions.
type WriteGroupMigrationReviewRequest struct {
	WorkspaceID               string   `json:"workspace_id"`
	ExpectedWorkspaceRevision *string  `json:"expected_workspace_revision"`
	ExpectedConnectorRevision string   `json:"expected_connector_revision"`
	ReviewDigest              string   `json:"review_digest"`
	SourceIDs                 []string `json:"source_ids"`
	ConfirmSnapshotConversion *bool    `json:"confirm_snapshot_conversion"`
}

// WriteGroupMigrationReviewResponse returns saved canonical drafts, not an
// activation result or evidence of external SQL writing.
type WriteGroupMigrationReviewResponse struct {
	Success bool                                       `json:"success"`
	Data    *workspace.WriteGroupMigrationReviewResult `json:"data"`
}

// ReviewSingleMappingMigration saves explicitly reviewed snapshot drafts.
// @Summary Save confirmed legacy single-mapping snapshot conversion
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param request body WriteGroupMigrationReviewRequest true "Preview digest, expected revisions and explicit snapshot confirmation"
// @Success 200 {object} WriteGroupMigrationReviewResponse
// @Failure 400,404,409,422,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/migrations/single-mappings/review [post]
func (h *StudioV2WorkspaceWriteGroupsHandler) ReviewSingleMappingMigration(c *gin.Context) {
	if !h.available(c) {
		return
	}
	var request WriteGroupMigrationReviewRequest
	if err := c.ShouldBindJSON(&request); err != nil || strings.TrimSpace(request.WorkspaceID) == "" ||
		request.ExpectedWorkspaceRevision == nil || strings.TrimSpace(request.ExpectedConnectorRevision) == "" ||
		strings.TrimSpace(request.ReviewDigest) == "" || len(request.SourceIDs) == 0 || request.ConfirmSnapshotConversion == nil {
		renderWriteGroupTyped(c, http.StatusBadRequest, "WRITE_GROUP_INVALID_REQUEST", "migration review requires preview, revisions and explicit confirmation", false, "review_request")
		return
	}
	for _, id := range request.SourceIDs {
		if strings.TrimSpace(id) == "" {
			renderWriteGroupTyped(c, http.StatusBadRequest, "WRITE_GROUP_INVALID_REQUEST", "migration review source identity is missing", false, "review_request")
			return
		}
	}
	result, err := h.groups.ReviewSingleMappingMigration(c.Request.Context(), workspace.WriteGroupMigrationReviewRequest{
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
