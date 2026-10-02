package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// WriteGroupMutationRequest carries explicit CAS values, including an empty
// workspace revision for a workspace that has not saved database setup yet.
type WriteGroupMutationRequest struct {
	WorkspaceID               string                `json:"workspace_id"`
	ExpectedWorkspaceRevision *string               `json:"expected_workspace_revision"`
	ExpectedGroupRevision     *string               `json:"expected_group_revision,omitempty"`
	ExpectedConnectorRevision *string               `json:"expected_connector_revision"`
	Group                     *workspace.WriteGroup `json:"group,omitempty"`
}

// WriteGroupSaveResponse is the canonical save/read result and its CAS revision.
type WriteGroupSaveResponse struct {
	Success bool                            `json:"success"`
	Data    *workspace.WriteGroupSaveResult `json:"data"`
}

// WriteGroupListResponse includes the workspace revision for the next save.
type WriteGroupListResponse struct {
	Success bool                            `json:"success"`
	Data    *workspace.WriteGroupListResult `json:"data"`
}

// WriteGroupReadinessResponse reports configuration and schema evidence.
type WriteGroupReadinessResponse struct {
	Success bool                           `json:"success"`
	Data    *workspace.WriteGroupReadiness `json:"data"`
}

// WriteGroupDeliveryResponse reports where a group's accepted data stands.
type WriteGroupDeliveryResponse struct {
	Success bool                        `json:"success"`
	Data    *grouppipeline.DeliveryView `json:"data"`
}

// WriteGroupDeliveryReader is the durable delivery truth source of a group.
type WriteGroupDeliveryReader interface {
	Delivery(ctx context.Context, groupID string) (*grouppipeline.DeliveryView, error)
}

// StudioV2WorkspaceWriteGroupsHandler exposes the canonical group service.
type StudioV2WorkspaceWriteGroupsHandler struct {
	groups   *workspace.WriteGroupService
	delivery WriteGroupDeliveryReader
}

// WithDelivery enables the delivery status route.
func (h *StudioV2WorkspaceWriteGroupsHandler) WithDelivery(reader WriteGroupDeliveryReader) *StudioV2WorkspaceWriteGroupsHandler {
	h.delivery = reader
	return h
}

// NewStudioV2WorkspaceWriteGroupsHandler uses the already wired local store.
func NewStudioV2WorkspaceWriteGroupsHandler(groups *workspace.WriteGroupService) *StudioV2WorkspaceWriteGroupsHandler {
	return &StudioV2WorkspaceWriteGroupsHandler{groups: groups}
}

// List returns a coherent group and workspace revision snapshot.
// @Summary List canonical workspace write groups
// @Tags Studio V2 Write Groups
// @Produce json
// @Success 200 {object} WriteGroupListResponse
// @Failure 500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups [get]
func (h *StudioV2WorkspaceWriteGroupsHandler) List(c *gin.Context) {
	if !h.available(c) {
		return
	}
	result, err := h.groups.List(c.Request.Context())
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, WriteGroupListResponse{Success: true, Data: result})
}

// Get returns one current-workspace group, including tombstones.
// @Summary Get a canonical workspace write group
// @Tags Studio V2 Write Groups
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} WriteGroupSaveResponse
// @Failure 404,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/{id} [get]
func (h *StudioV2WorkspaceWriteGroupsHandler) Get(c *gin.Context) {
	if !h.available(c) {
		return
	}
	result, err := h.groups.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, WriteGroupSaveResponse{Success: true, Data: result})
}

// Readiness checks a saved group without applying it or altering its target.
// @Summary Check canonical write-group configuration and schema
// @Tags Studio V2 Write Groups
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} WriteGroupReadinessResponse
// @Failure 404,409,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/{id}/readiness [get]
func (h *StudioV2WorkspaceWriteGroupsHandler) Readiness(c *gin.Context) {
	if !h.available(c) {
		return
	}
	result, err := h.groups.Readiness(c.Request.Context(), c.Param("id"))
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, WriteGroupReadinessResponse{Success: true, Data: result})
}

// Delivery reports the group's durable delivery truth: intake state, how many
// rows are collecting, queued, retrying, blocked, quarantined, unknown or
// confirmed by the destination, the backlog per group and destination
// revision, and capacity. A committed count only ever comes from destination
// confirmation; buffered or acknowledged data is never reported as written.
// @Summary Read a write group's delivery status
// @Tags Studio V2 Write Groups
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} WriteGroupDeliveryResponse
// @Failure 404,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/{id}/delivery [get]
func (h *StudioV2WorkspaceWriteGroupsHandler) Delivery(c *gin.Context) {
	if !h.available(c) {
		return
	}
	if _, err := h.groups.Get(c.Request.Context(), c.Param("id")); err != nil {
		renderWriteGroupError(c, err)
		return
	}
	if h.delivery == nil {
		renderWriteGroupError(c, workspace.ErrWriteGroupServiceUnavailable)
		return
	}
	view, err := h.delivery.Delivery(c.Request.Context(), c.Param("id"))
	if err != nil {
		renderWriteGroupError(c, workspace.ErrWriteGroupServiceUnavailable)
		return
	}
	c.JSON(http.StatusOK, WriteGroupDeliveryResponse{Success: true, Data: view})
}

// Create saves a draft without external DDL or activation.
// @Summary Save a canonical write-group draft
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param request body WriteGroupMutationRequest true "Draft and expected revisions"
// @Success 201 {object} WriteGroupSaveResponse
// @Failure 400,404,409,422,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups [post]
func (h *StudioV2WorkspaceWriteGroupsHandler) Create(c *gin.Context) {
	h.mutate(c, http.StatusCreated, writeGroupCreate)
}

// Update preserves applied state while saving a CAS-checked draft.
// @Summary Update a canonical write-group draft
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body WriteGroupMutationRequest true "Draft and expected revisions"
// @Success 200 {object} WriteGroupSaveResponse
// @Failure 400,404,409,422,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/{id} [put]
func (h *StudioV2WorkspaceWriteGroupsHandler) Update(c *gin.Context) {
	h.mutate(c, http.StatusOK, writeGroupUpdate)
}

// Delete tombstones a group while preserving its persisted identities.
// @Summary Tombstone a canonical write group
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body WriteGroupMutationRequest true "Expected revisions"
// @Success 200 {object} WriteGroupSaveResponse
// @Failure 400,404,409,422,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/{id} [delete]
func (h *StudioV2WorkspaceWriteGroupsHandler) Delete(c *gin.Context) {
	h.mutate(c, http.StatusOK, writeGroupDelete)
}

// Disable stops new intake while retaining accepted data and applied snapshots.
// @Summary Disable a canonical write group
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body WriteGroupMutationRequest true "Expected revisions"
// @Success 200 {object} WriteGroupSaveResponse
// @Failure 400,404,409,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/{id}/disable [post]
func (h *StudioV2WorkspaceWriteGroupsHandler) Disable(c *gin.Context) {
	h.mutate(c, http.StatusOK, writeGroupDisable)
}

// Apply schedules the saved draft for the next UTC bucket after readiness
// passes. It performs no DDL and no connection to the destination; a draft
// that is not ready, or any stale revision, changes nothing.
// @Summary Apply a canonical write group
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body WriteGroupMutationRequest true "Expected revisions (no group payload)"
// @Success 200 {object} WriteGroupSaveResponse
// @Failure 400,404,409,422,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/{id}/apply [post]
func (h *StudioV2WorkspaceWriteGroupsHandler) Apply(c *gin.Context) {
	h.mutate(c, http.StatusOK, writeGroupApply)
}

type writeGroupMutationAction string

const (
	writeGroupCreate  writeGroupMutationAction = "create"
	writeGroupUpdate  writeGroupMutationAction = "update"
	writeGroupDelete  writeGroupMutationAction = "delete"
	writeGroupDisable writeGroupMutationAction = "disable"
	writeGroupApply   writeGroupMutationAction = "apply"
)

func (h *StudioV2WorkspaceWriteGroupsHandler) mutate(c *gin.Context, status int, action writeGroupMutationAction) {
	if !h.available(c) {
		return
	}
	var request WriteGroupMutationRequest
	create := action == writeGroupCreate
	requireGroup := create || action == writeGroupUpdate
	if err := c.ShouldBindJSON(&request); err != nil ||
		strings.TrimSpace(request.WorkspaceID) == "" || request.ExpectedWorkspaceRevision == nil ||
		request.ExpectedConnectorRevision == nil || strings.TrimSpace(*request.ExpectedConnectorRevision) == "" ||
		(!create && (request.ExpectedGroupRevision == nil || strings.TrimSpace(*request.ExpectedGroupRevision) == "")) ||
		(requireGroup && request.Group == nil) {
		renderWriteGroupTyped(c, http.StatusBadRequest, "WRITE_GROUP_INVALID_REQUEST", "write-group request is incomplete", false, "review_request")
		return
	}
	mutation := workspace.WriteGroupMutation{
		WorkspaceID: request.WorkspaceID, ExpectedWorkspaceRevision: *request.ExpectedWorkspaceRevision,
		ExpectedConnectorRevision: *request.ExpectedConnectorRevision, Group: request.Group,
	}
	if request.ExpectedGroupRevision != nil {
		mutation.ExpectedGroupRevision = *request.ExpectedGroupRevision
	}
	var result *workspace.WriteGroupSaveResult
	var err error
	switch action {
	case writeGroupCreate:
		result, err = h.groups.Create(c.Request.Context(), mutation)
	case writeGroupDelete:
		result, err = h.groups.Delete(c.Request.Context(), c.Param("id"), mutation)
	case writeGroupDisable:
		result, err = h.groups.Disable(c.Request.Context(), c.Param("id"), mutation)
	case writeGroupUpdate:
		result, err = h.groups.Update(c.Request.Context(), c.Param("id"), mutation)
	case writeGroupApply:
		result, err = h.groups.Apply(c.Request.Context(), c.Param("id"), mutation)
	default:
		err = workspace.ErrWriteGroupValidation
	}
	if err != nil {
		if action != writeGroupDelete && errors.Is(err, workspace.ErrWriteGroupLifecycleBlocked) {
			renderWriteGroupTyped(c, http.StatusConflict, "WRITE_GROUP_LIFECYCLE_BLOCKED", "write group cannot be changed in its current lifecycle state", false, "reload")
			return
		}
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(status, WriteGroupSaveResponse{Success: true, Data: result})
}

func (h *StudioV2WorkspaceWriteGroupsHandler) available(c *gin.Context) bool {
	if h.groups != nil {
		return true
	}
	renderWriteGroupError(c, workspace.ErrDatabaseSetupUnavailable)
	return false
}

func renderWriteGroupError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, workspace.ErrWriteGroupNotFound), errors.Is(err, workspace.ErrWriteGroupConnectorNotFound), errors.Is(err, workspace.ErrNotFound):
		renderWriteGroupTyped(c, http.StatusNotFound, "WRITE_GROUP_NOT_FOUND", "write group or resource is not available", false, "reload")
	case errors.Is(err, workspace.ErrSetupRevisionConflict), errors.Is(err, workspace.ErrWriteGroupRevisionConflict):
		renderWriteGroupTyped(c, http.StatusConflict, "revision_mismatch", "saved configuration changed; reload before saving", false, "reload")
	case errors.Is(err, workspace.ErrWriteGroupLifecycleBlocked):
		renderWriteGroupTyped(c, http.StatusConflict, "WRITE_GROUP_LIFECYCLE_BLOCKED", "write group cannot be deleted while accepted data ownership is unresolved", false, "disable_group")
	case errors.Is(err, workspace.ErrWriteGroupApplyNotReady):
		renderWriteGroupTyped(c, http.StatusConflict, "WRITE_GROUP_NOT_READY", "write group is not ready to apply; review its readiness issues", false, "review_group")
	case errors.Is(err, workspace.ErrWriteGroupValidation):
		renderWriteGroupTyped(c, http.StatusUnprocessableEntity, "WRITE_GROUP_INVALID", "write-group configuration is invalid", false, "review_group")
	case errors.Is(err, workspace.ErrDatabaseSetupUnavailable), errors.Is(err, workspace.ErrWriteGroupServiceUnavailable), errors.Is(err, workspace.ErrWriteGroupReadinessUnavailable):
		renderWriteGroupTyped(c, http.StatusServiceUnavailable, "WRITE_GROUP_UNAVAILABLE", "write-group service is unavailable", true, "retry")
	default:
		renderWriteGroupTyped(c, http.StatusInternalServerError, "WRITE_GROUP_UNAVAILABLE", "write group could not be saved or read", true, "reload")
	}
}

func renderWriteGroupTyped(c *gin.Context, status int, code, message string, retryable bool, action string) {
	renderStudioV2WorkspaceDatabaseTyped(c, status, TypedAPIErrorEnvelope{
		Code: code, Message: message, Retryable: retryable,
		RequestID: getOrGenerateRequestID(c), Action: action,
	})
}
