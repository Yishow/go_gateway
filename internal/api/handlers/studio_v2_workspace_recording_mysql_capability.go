package handlers

import (
	"context"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
)

const (
	// mysqlRecordingScopeInvalidCode is the stable validation code for a
	// saved-scope capability query that is missing a required parameter. It
	// never reports which connector the caller named.
	mysqlRecordingScopeInvalidCode    = "RECORDING_CONNECTOR_SCOPE_INVALID"
	mysqlRecordingScopeInvalidMessage = "saved connector scope is incomplete"
)

// mysqlRecordingScopeRequested reports whether the caller asked for a
// saved-scope MySQL capability. A bare kind query stays on the legacy static
// matrix; the presence of either saved-scope key commits the caller to the
// probe branch, even when its value is empty or whitespace.
func mysqlRecordingScopeRequested(c *gin.Context) bool {
	query := c.Request.URL.Query()
	for _, key := range []string{"connector_id", "expected_connector_revision"} {
		if _, present := query[key]; present {
			return true
		}
	}
	return false
}

// mysqlRecordingCapabilities answers a saved-scope MySQL capability query from
// the workspace-bound connector only. The request names a saved identity and
// revision; every capability fact comes from the server-side readonly probe, so
// a client cannot declare its own family, version or verified operations.
func (h *StudioV2WorkspaceRecordingPlansHandler) mysqlRecordingCapabilities(c *gin.Context) {
	kind := strings.TrimSpace(c.Query("kind"))
	connectorID := strings.TrimSpace(c.Query("connector_id"))
	expectedRevision := strings.TrimSpace(c.Query("expected_connector_revision"))

	if canonicalRecordingDialect(kind) != recordingDialectMySQL {
		// A saved scope is only defined for MySQL. Refuse before resolving any
		// connector so a foreign kind cannot be used to enumerate destinations.
		renderRecordingPlanTargetError(c, errRecordingPlanDialectMismatch)
		return
	}
	if connectorID == "" || expectedRevision == "" {
		renderRecordingPlanError(c, http.StatusUnprocessableEntity, mysqlRecordingScopeInvalidCode, mysqlRecordingScopeInvalidMessage, false, "select a saved MySQL connector and its current revision")
		return
	}

	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}
	if record == nil || strings.TrimSpace(record.DatabaseConnectorID) == "" {
		renderRecordingPlanTargetError(c, errRecordingPlanConnectorUnavailable)
		return
	}
	if connectorID != strings.TrimSpace(record.DatabaseConnectorID) {
		renderRecordingPlanTargetError(c, errRecordingPlanDestinationScope)
		return
	}

	connector, err := h.resolveSavedConnector(c.Request.Context(), connectorID, expectedRevision)
	if err != nil {
		renderRecordingPlanTargetError(c, err)
		return
	}
	if connector == nil || strings.TrimSpace(connector.ID) != connectorID {
		renderRecordingPlanTargetError(c, errRecordingPlanDestinationScope)
		return
	}
	if strings.TrimSpace(connector.IdentityRevision) != expectedRevision {
		renderRecordingPlanTargetError(c, errRecordingPlanConnectorRevisionMismatch)
		return
	}
	if !connector.Enabled {
		renderRecordingPlanTargetError(c, errRecordingPlanConnectorUnavailable)
		return
	}
	if connector.Kind != schema.DatabaseConnectorKindMySQL {
		renderRecordingPlanTargetError(c, errRecordingPlanDialectMismatch)
		return
	}

	// The readonly family/operation probe is an optional narrow interface on the
	// saved connector service, so existing resolver fakes and services do not
	// grow a method they never implement. Without it the answer stays a safe
	// unavailable: no capability is inferred from the driver, the kind or the
	// request, and no DDL or write path is opened.
	probe, ok := h.dbTargetSvc.(recordingMySQLCapabilityProbe)
	if !ok || probe == nil {
		renderRecordingPlanTargetError(c, errRecordingPlanResolverUnavailable)
		return
	}
	capability, err := probe.ProbeMySQLRecordingAtRevision(c.Request.Context(), connectorID, expectedRevision)
	if err != nil {
		renderRecordingPlanTargetError(c, err)
		return
	}
	if capability == nil {
		renderRecordingPlanTargetError(c, errRecordingPlanResolverUnavailable)
		return
	}
	// The probe result must bind the saved connector scope it was asked about.
	// A mismatched answer is an adapter inconsistency, never a client success.
	if strings.TrimSpace(capability.ConnectorID) != connectorID ||
		strings.TrimSpace(capability.ConnectorRevision) != expectedRevision {
		renderRecordingPlanTargetError(c, errRecordingPlanResolverUnavailable)
		return
	}

	// The static matrix keeps kind=mysql connected while its managed and group
	// mutations stay masked; the probe carries only server-derived facts and
	// operation gates. Client query claims are never copied into either.
	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey: mysqlRecordingCapabilityResponse{
			ConnectorCapability:      h.planSvc.GetConnectorCapability(recordingDialectMySQL),
			MySQLRecordingCapability: capability,
		},
	})
}

// recordingMySQLCapabilityProbe is the narrow optional interface a saved
// connector service may implement. It is intentionally separate from
// recordingSavedConnectorResolver so existing fakes and callers keep working.
type recordingMySQLCapabilityProbe interface {
	ProbeMySQLRecordingAtRevision(ctx context.Context, connectorID, expectedRevision string) (*dbtarget.MySQLRecordingCapability, error)
}

// mysqlRecordingCapabilityResponse flattens the static kind matrix and the
// server-derived probe into one object, so a client can see the connector is
// connected and MySQL-kind while every mutating operation stays gated.
type mysqlRecordingCapabilityResponse struct {
	recordingplan.ConnectorCapability
	*dbtarget.MySQLRecordingCapability
}
