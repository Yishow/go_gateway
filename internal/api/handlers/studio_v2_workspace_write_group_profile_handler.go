package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const writeGroupEnsureCircuitProfile writeGroupMutationAction = "ensure_circuit_profile"

// EnsureCircuitProfile creates or reuses the canonical group of one circuit
// profile. A resent request returns the saved group (200) instead of adding a
// second writer; a different interval, source set or owner is a 409 conflict.
// The interval is an integer in the profile's own unit and is never defaulted.
// @Summary Ensure the canonical group of a circuit recording profile
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param request body WriteGroupMutationRequest true "Profile group and expected revisions"
// @Success 200 {object} WriteGroupSaveResponse
// @Failure 400,404,409,422,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/circuit-profiles/ensure [post]
func (h *StudioV2WorkspaceWriteGroupsHandler) EnsureCircuitProfile(c *gin.Context) {
	h.mutate(c, http.StatusOK, writeGroupEnsureCircuitProfile)
}
