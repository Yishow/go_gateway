package handlers

import (
	"errors"
	"net/http"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// renderStudioV2WorkspaceDatabaseTyped writes one database error case through
// the shared success/error wrapper with a typed envelope.
func renderStudioV2WorkspaceDatabaseTyped(c *gin.Context, status int, envelope TypedAPIErrorEnvelope) {
	c.JSON(status, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: envelope})
}

// renderStudioV2WorkspaceDatabaseMessage writes the shared wrapper for cases
// whose entire safe payload is one message without a typed code.
func renderStudioV2WorkspaceDatabaseMessage(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		apiResponseSuccessKey: false,
		apiResponseErrorKey:   gin.H{apiResponseMessageKey: message},
	})
}

// studioV2WorkspaceDatabaseFailureMessage is the safe fallback message shared
// by every Studio V2 workspace database 500 response.
const studioV2WorkspaceDatabaseFailureMessage = "Studio V2 database operation failed"

func renderStudioV2WorkspaceDatabaseError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, workspace.ErrSetupRevisionConflict):
		renderStudioV2WorkspaceDatabaseTyped(c, http.StatusConflict, TypedAPIErrorEnvelope{
			Code: "revision_mismatch", Message: "database setup changed; reload the saved configuration",
			RequestID: getOrGenerateRequestID(c), Action: "reload",
		})
	case errors.Is(err, workspace.ErrDatabaseSetupUnavailable), errors.Is(err, dbtarget.ErrSetupTransactionUnavailable):
		renderStudioV2WorkspaceDatabaseTyped(c, http.StatusServiceUnavailable, TypedAPIErrorEnvelope{
			Code: apiInternalErrorCode, Message: "Studio V2 database setup is unavailable", Retryable: true,
			RequestID: getOrGenerateRequestID(c), Action: "retry",
		})
	case errors.Is(err, dbtarget.ErrConnectorDisabled):
		renderStudioV2WorkspaceDatabaseMessage(c, http.StatusUnprocessableEntity, "database connection is disabled; enable it before selecting it")
	case errors.Is(err, dbtarget.ErrConnectorNotFound):
		renderStudioV2WorkspaceDatabaseMessage(c, http.StatusNotFound, "database connection is not available")
	case errors.Is(err, dbtarget.ErrConnectorRevisionConflict):
		renderStudioV2WorkspaceDatabaseMessage(c, http.StatusConflict, "database connection changed; reload the saved configuration")
	case errors.Is(err, workspace.ErrValidation), errors.Is(err, sourcerule.ErrValidation), errors.Is(err, dbtarget.ErrValidation):
		renderStudioV2WorkspaceValidationError(c, err)
	default:
		renderStudioV2WorkspaceDatabaseMessage(c, http.StatusInternalServerError, studioV2WorkspaceDatabaseFailureMessage)
	}
}

func optionalConnectorRevision(revision string) *string {
	if revision == "" {
		return nil
	}
	return &revision
}

var (
	// errWorkspaceSchemaPreparationRequired marks a workspace whose saved target
	// is missing required tables; activation reports it instead of creating them.
	errWorkspaceSchemaPreparationRequired = errors.New("database schema preparation is required before activation")
	// errSchemaMutationNeedsConfirmation marks a non-preview schema request that
	// arrived without a confirmed preview.
	errSchemaMutationNeedsConfirmation = errors.New("schema creation requires a confirmed preview")
)

const (
	workspaceSchemaPreparationRequiredCode = "WORKSPACE_SCHEMA_PREPARATION_REQUIRED"
	schemaConfirmationRequiredCode         = "SCHEMA_CONFIRMATION_REQUIRED"
)

// renderSchemaConfirmationRequired refuses a non-preview schema request that
// carries no confirmed preview. The read-only plan stays available.
func renderSchemaConfirmationRequired(c *gin.Context) {
	renderStudioV2WorkspaceDatabaseTyped(c, http.StatusConflict, TypedAPIErrorEnvelope{
		Code:      schemaConfirmationRequiredCode,
		Message:   errSchemaMutationNeedsConfirmation.Error(),
		RequestID: getOrGenerateRequestID(c),
		Action:    "preview the managed schema and confirm it before creating tables",
	})
}

// requireSchemaConfirmation is the shared dry_run gate of every schema
// tooling entry: only a read-only dry run passes without a confirmation.
func requireSchemaConfirmation(c *gin.Context, dryRun bool) bool {
	if dryRun {
		return true
	}
	renderSchemaConfirmationRequired(c)
	return false
}

func renderStudioV2WorkspaceSchemaEnsureError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errWorkspaceSchemaPreparationRequired):
		renderStudioV2WorkspaceDatabaseTyped(c, http.StatusUnprocessableEntity, TypedAPIErrorEnvelope{
			Code:      workspaceSchemaPreparationRequiredCode,
			Message:   err.Error(),
			RequestID: getOrGenerateRequestID(c),
			Action:    "prepare the database schema, then activate again",
		})
	case errors.Is(err, workspace.ErrValidation), errors.Is(err, sourcerule.ErrValidation), errors.Is(err, dbtarget.ErrValidation):
		renderStudioV2WorkspaceDatabaseTyped(c, http.StatusUnprocessableEntity, TypedAPIErrorEnvelope{
			Code:      apiValidationErrorCode,
			Message:   err.Error(),
			RequestID: getOrGenerateRequestID(c),
		})
	default:
		renderStudioV2WorkspaceDatabaseMessage(c, http.StatusInternalServerError, studioV2WorkspaceDatabaseFailureMessage)
	}
}
