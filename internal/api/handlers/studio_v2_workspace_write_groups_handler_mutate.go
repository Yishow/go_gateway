package handlers

import (
	"errors"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

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
	// Create and ensure address no existing group, so they carry no group revision.
	create := action == writeGroupCreate || action == writeGroupEnsureCircuitProfile
	requireGroup := create || action == writeGroupUpdate
	bindErr := c.ShouldBindJSON(&request)
	if errors.Is(bindErr, workspace.ErrWriteGroupProfileIntervalInvalid) || errors.Is(bindErr, workspace.ErrWriteGroupMemberMetadataInvalid) {
		renderWriteGroupError(c, bindErr)
		return
	}
	if bindErr != nil ||
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
	case writeGroupEnsureCircuitProfile:
		result, err = h.groups.EnsureCircuitProfile(c.Request.Context(), mutation)
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
		if action == writeGroupDelete && renderPermanentDeleteError(c, err) {
			return
		}
		if action != writeGroupDelete && errors.Is(err, workspace.ErrWriteGroupLifecycleBlocked) {
			renderWriteGroupTyped(c, http.StatusConflict, "WRITE_GROUP_LIFECYCLE_BLOCKED", "write group cannot be changed in its current lifecycle state", false, "reload")
			return
		}
		renderWriteGroupError(c, err)
		return
	}
	if action == writeGroupDelete {
		c.JSON(status, WriteGroupDeleteResponse{Success: true, Data: WriteGroupDeleteData{
			Deleted: result.Deleted, GroupID: result.GroupID, LineageID: result.LineageID, DeletedGenerations: result.DeletedGenerations,
		}})
		return
	}
	c.JSON(status, WriteGroupSaveResponse{Success: true, Data: result})
}
