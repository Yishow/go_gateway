package handlers

import (
	"errors"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
	"net/http"

	"github.com/gin-gonic/gin"
)

func renderStudioV2WorkspaceMappingError(c *gin.Context, err error) {
	status, code, message, action := http.StatusInternalServerError, "workspace_mapping_save_failed", "Mapping could not be saved; retry explicitly", "retry"
	switch {
	case errors.Is(err, mapping.ErrWorkspaceConfirmationConflict):
		status, code, message, action = http.StatusConflict, "workspace_mapping_conflict", "Source rule or mapping changed; review and explicitly reapply the source candidate", "review_source"
	case errors.Is(err, mapping.ErrValidation), errors.Is(err, sourcerule.ErrValidation), errors.Is(err, workspace.ErrValidation):
		renderStudioV2WorkspaceValidationError(c, err)
		return
	case errors.Is(err, mapping.ErrMappingNotFound), errors.Is(err, point.ErrPointNotFound), errors.Is(err, tag.ErrTagNotFound), errors.Is(err, sourcerule.ErrSourceRuleNotFound):
		status, code, message, action = http.StatusNotFound, "workspace_mapping_not_found", studioV2MappingRowNotFoundMessage, "reload"
	}
	c.JSON(status, APIErrorResponse{Error: TypedAPIErrorEnvelope{Code: code, Message: message, Retryable: true, RequestID: getOrGenerateRequestID(c), Action: action}})
}
