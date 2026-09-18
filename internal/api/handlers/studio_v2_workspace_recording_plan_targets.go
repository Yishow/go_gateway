package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

const (
	recordingDialectSQLite   = "sqlite"
	recordingDialectPostgres = "postgres"
	recordingDialectMySQL    = "mysql"
)

var (
	errRecordingPlanNotFound                  = errors.New("recording plan not found")
	errRecordingPlanDestinationRequired       = errors.New("recording plan requires a saved database destination")
	errRecordingPlanDestinationScope          = errors.New("recording destination is outside the workspace")
	errRecordingPlanConnectorRevisionMissing  = errors.New("recording connector revision is required")
	errRecordingPlanConnectorRevisionMismatch = errors.New("recording connector revision does not match the plan destination")
	errRecordingPlanConnectorUnavailable      = errors.New("recording connector is unavailable")
	errRecordingPlanResolverUnavailable       = errors.New("recording connector resolver is unavailable")
	errRecordingPlanDialectMismatch           = errors.New("recording connector dialect does not match the request")
)

const (
	recordingPlanNotFoundCode          = "RECORDING_PLAN_NOT_FOUND"
	recordingPlanDestinationCode       = "RECORDING_PLAN_DESTINATION_INVALID"
	recordingPlanConnectorCode         = "RECORDING_CONNECTOR_UNAVAILABLE"
	recordingPlanConnectorRevisionCode = "RECORDING_CONNECTOR_REVISION_CONFLICT"
	recordingPlanPreviewFailedCode     = "SCHEMA_PREVIEW_FAILED"
	recordingPlanUnavailableCode       = "RECORDING_PLAN_UNAVAILABLE"
	recordingPlanNotFoundMessage       = "recording plan not found"
	recordingPlanDestinationMessage    = "recording plan destination is invalid"
	recordingPlanConnectorMessage      = "recording connector is unavailable"
	recordingPlanConnectorRevisionMsg  = "recording connector revision is stale"
	recordingPlanPreviewFailedMessage  = "recording schema preview could not be generated"
	recordingPlanUnavailableMessage    = "recording plan is temporarily unavailable"
)

func (h *StudioV2WorkspaceRecordingPlansHandler) validatePlanDestinations(ctx context.Context, record *workspace.Record, destinations []recordingplan.PlanDestination) error {
	if len(destinations) == 0 {
		return errRecordingPlanDestinationRequired
	}
	if record == nil || strings.TrimSpace(record.DatabaseConnectorID) == "" {
		return errRecordingPlanConnectorUnavailable
	}

	workspaceConnectorID := strings.TrimSpace(record.DatabaseConnectorID)
	for _, destination := range destinations {
		if strings.TrimSpace(destination.ConnectorID) != workspaceConnectorID {
			return errRecordingPlanDestinationScope
		}
		if strings.TrimSpace(destination.ConnectorRevision) == "" {
			return errRecordingPlanConnectorRevisionMissing
		}
	}

	for _, destination := range destinations {
		destinationRevision := strings.TrimSpace(destination.ConnectorRevision)
		connector, err := h.resolveSavedConnector(ctx, workspaceConnectorID, destinationRevision)
		if err != nil {
			return err
		}
		if connector == nil || strings.TrimSpace(connector.ID) != workspaceConnectorID {
			return errRecordingPlanDestinationScope
		}
		if strings.TrimSpace(connector.IdentityRevision) != destinationRevision {
			return errRecordingPlanConnectorRevisionMismatch
		}
		if !connector.Enabled {
			return errRecordingPlanConnectorUnavailable
		}
	}
	return nil
}

func (h *StudioV2WorkspaceRecordingPlansHandler) resolvePreviewConnector(ctx context.Context, record *workspace.Record, plan *recordingplan.RecordingPlan, connectorID, expectedConnectorRevision string) (*schema.DatabaseConnector, error) {
	if record == nil || strings.TrimSpace(record.DatabaseConnectorID) != connectorID {
		return nil, errRecordingPlanDestinationScope
	}

	selectedRevision := strings.TrimSpace(expectedConnectorRevision)
	selected := false
	for _, destination := range plan.Destinations {
		if strings.TrimSpace(destination.ConnectorID) != connectorID {
			continue
		}
		selected = true
		destinationRevision := strings.TrimSpace(destination.ConnectorRevision)
		if destinationRevision == "" {
			return nil, errRecordingPlanConnectorRevisionMissing
		}
		if destinationRevision != selectedRevision {
			return nil, errRecordingPlanConnectorRevisionMismatch
		}
	}
	if !selected {
		return nil, errRecordingPlanDestinationScope
	}

	connector, err := h.resolveSavedConnector(ctx, connectorID, selectedRevision)
	if err != nil {
		return nil, err
	}
	if connector == nil || strings.TrimSpace(connector.ID) != connectorID {
		return nil, errRecordingPlanDestinationScope
	}
	if strings.TrimSpace(connector.IdentityRevision) != selectedRevision {
		return nil, errRecordingPlanConnectorRevisionMismatch
	}
	if !connector.Enabled {
		return nil, errRecordingPlanConnectorUnavailable
	}
	return connector, nil
}

func (h *StudioV2WorkspaceRecordingPlansHandler) resolveSavedConnector(ctx context.Context, connectorID, expectedRevision string) (*schema.DatabaseConnector, error) {
	if h.dbTargetSvc == nil {
		return nil, errRecordingPlanResolverUnavailable
	}
	if connectorSvc, ok := h.dbTargetSvc.(*dbtarget.ConnectorService); ok && connectorSvc == nil {
		return nil, errRecordingPlanResolverUnavailable
	}
	return h.dbTargetSvc.ResolveSavedTarget(ctx, strings.TrimSpace(connectorID), strings.TrimSpace(expectedRevision))
}

// getPlanForWorkspace reads a plan through the workspace-scoped repository
// lookup; an absent plan and another workspace's plan are indistinguishable.
func (h *StudioV2WorkspaceRecordingPlansHandler) getPlanForWorkspace(ctx context.Context, workspaceID, planID string) (*recordingplan.RecordingPlan, error) {
	if strings.TrimSpace(planID) == "" {
		return nil, errRecordingPlanNotFound
	}
	plan, err := h.planSvc.GetPlanByWorkspace(ctx, planID, workspaceID)
	if errors.Is(err, recordingplan.ErrPlanNotFound) {
		return nil, errRecordingPlanNotFound
	}
	return plan, err
}

func recordingDialectMatches(kind, requested string) bool {
	requested = strings.TrimSpace(requested)
	if requested == "" {
		return true
	}
	return canonicalRecordingDialect(kind) == canonicalRecordingDialect(requested)
}

func canonicalRecordingDialect(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case recordingDialectSQLite, "sqlite3":
		return recordingDialectSQLite
	case recordingDialectPostgres, "postgresql", "pgx":
		return recordingDialectPostgres
	case recordingDialectMySQL:
		return recordingDialectMySQL
	default:
		return strings.ToLower(strings.TrimSpace(value))
	}
}

func renderRecordingPlanTargetError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errRecordingPlanDestinationRequired):
		renderRecordingPlanError(c, http.StatusUnprocessableEntity, recordingPlanDestinationCode, recordingPlanDestinationMessage, false, "configure a saved database destination")
	case errors.Is(err, errRecordingPlanDestinationScope):
		renderRecordingPlanError(c, http.StatusNotFound, recordingPlanNotFoundCode, recordingPlanNotFoundMessage, false, "review the selected recording plan")
	case errors.Is(err, errRecordingPlanConnectorRevisionMissing):
		renderRecordingPlanError(c, http.StatusBadRequest, "RECORDING_CONNECTOR_REVISION_REQUIRED", "expected connector revision is required", false, "refresh the saved connector")
	case errors.Is(err, errRecordingPlanConnectorRevisionMismatch):
		renderRecordingPlanError(c, http.StatusConflict, recordingPlanConnectorRevisionCode, recordingPlanConnectorRevisionMsg, false, "refresh the saved connector")
	case errors.Is(err, errRecordingPlanConnectorUnavailable):
		renderRecordingPlanError(c, http.StatusUnprocessableEntity, recordingPlanConnectorCode, recordingPlanConnectorMessage, false, "enable the saved connector")
	case errors.Is(err, errRecordingPlanDialectMismatch):
		renderRecordingPlanError(c, http.StatusUnprocessableEntity, "RECORDING_CONNECTOR_DIALECT_MISMATCH", "recording connector dialect does not match the request", false, "use the saved connector dialect")
	case errors.Is(err, errRecordingPlanResolverUnavailable):
		renderRecordingPlanError(c, http.StatusServiceUnavailable, recordingPlanUnavailableCode, recordingPlanUnavailableMessage, true, "retry after the saved connector service is available")
	case errors.Is(err, dbtarget.ErrConnectorDisabled):
		renderRecordingPlanError(c, http.StatusUnprocessableEntity, recordingPlanConnectorCode, recordingPlanConnectorMessage, false, "enable the saved connector")
	case errors.Is(err, dbtarget.ErrConnectorNotFound):
		renderRecordingPlanError(c, http.StatusNotFound, recordingPlanNotFoundCode, recordingPlanNotFoundMessage, false, "review the selected recording plan")
	case errors.Is(err, dbtarget.ErrConnectorRevisionConflict):
		renderRecordingPlanError(c, http.StatusConflict, recordingPlanConnectorRevisionCode, recordingPlanConnectorRevisionMsg, false, "refresh the saved connector")
	default:
		renderRecordingPlanError(c, http.StatusServiceUnavailable, recordingPlanUnavailableCode, recordingPlanUnavailableMessage, true, "retry after the saved connector service is available")
	}
}

func renderRecordingPlanReadError(c *gin.Context, err error) {
	if errors.Is(err, errRecordingPlanNotFound) {
		renderRecordingPlanError(c, http.StatusNotFound, recordingPlanNotFoundCode, recordingPlanNotFoundMessage, false, "review the recording plan")
		return
	}
	renderRecordingPlanError(c, http.StatusInternalServerError, recordingPlanUnavailableCode, recordingPlanUnavailableMessage, true, "retry loading recording plans")
}

func renderRecordingPlanError(c *gin.Context, status int, code, message string, retryable bool, action string) {
	c.JSON(status, gin.H{
		apiResponseSuccessKey: false,
		apiResponseErrorKey: TypedAPIErrorEnvelope{
			Code:      code,
			Message:   message,
			Retryable: retryable,
			RequestID: getOrGenerateRequestID(c),
			Action:    action,
		},
	})
}
