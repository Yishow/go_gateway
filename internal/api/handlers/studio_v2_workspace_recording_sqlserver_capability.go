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
	// sqlserverRecordingScopeInvalidCode is the stable validation code for a
	// saved-scope SQL Server capability query that is missing a required
	// parameter. It never reports which connector the caller named.
	sqlServerRecordingScopeInvalidCode    = "RECORDING_CONNECTOR_SCOPE_INVALID"
	sqlServerRecordingScopeInvalidMessage = "saved connector scope is incomplete"
	sqlServerConnectorIDQueryKey          = "connector_id"
)

// sqlserverRecordingScopeRequested reports whether the caller asked for a
// saved-scope SQL Server capability. A bare kind query stays on the legacy
// static matrix; the presence of either saved-scope key commits the caller to
// the strict scope branch, even when its value is empty or whitespace.
func sqlserverRecordingScopeRequested(c *gin.Context) bool {
	// The kind gate comes first: MySQL saved-scope queries keep their own
	// branch, and only kind=sqlserver commits the caller to this one.
	if canonicalRecordingDialect(strings.TrimSpace(c.Query("kind"))) != string(schema.DatabaseConnectorKindSQLServer) {
		return false
	}
	if studioLocatorNativeRequested(c) {
		return true
	}
	query := c.Request.URL.Query()
	for _, key := range []string{sqlServerConnectorIDQueryKey, "expected_connector_revision"} {
		if _, present := query[key]; present {
			return true
		}
	}
	return false
}

// sqlserverRecordingCapabilities answers a saved-scope SQL Server capability
// query from the workspace-bound connector only. The request names a saved
// identity and revision; every capability fact comes from the server-side
// readonly probe, so a client cannot declare its own family, version or
// verified operations. A connected, metadata-verified probe still reports
// every operation flag false until the real acceptance matrix verifies them.
func (h *StudioV2WorkspaceRecordingPlansHandler) sqlserverRecordingCapabilities(c *gin.Context) {
	if !validateStudioLocatorNativeQuery(c) {
		return
	}
	kind := strings.TrimSpace(c.Query("kind"))
	connectorID := strings.TrimSpace(c.Query("connector_id"))
	expectedRevision := strings.TrimSpace(c.Query("expected_connector_revision"))

	if canonicalRecordingDialect(kind) != string(schema.DatabaseConnectorKindSQLServer) {
		// A saved scope is only defined for SQL Server. Refuse before resolving
		// any connector so a foreign kind cannot enumerate destinations.
		renderRecordingPlanTargetError(c, errRecordingPlanDialectMismatch)
		return
	}
	if connectorID == "" || expectedRevision == "" {
		renderRecordingPlanError(c, http.StatusUnprocessableEntity, sqlServerRecordingScopeInvalidCode, sqlServerRecordingScopeInvalidMessage, false, "select a saved SQL Server connector and its current revision")
		return
	}

	record, locatorScope, ok := beginStudioLocatorNativeRead(c, h.workspaceSvc, func() { renderStudioV2WorkspaceBootstrapError(c) })
	if !ok {
		return
	}
	var err error
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
	if connector.Kind != schema.DatabaseConnectorKindSQLServer {
		renderRecordingPlanTargetError(c, errRecordingPlanDialectMismatch)
		return
	}

	// The readonly scope probe is an optional narrow interface on the saved
	// connector service, so existing resolver fakes and services do not grow a
	// method they never implement. Without it the answer stays a safe
	// unavailable: no capability is inferred from the driver, the kind or the
	// request, and no DDL or write path is opened.
	probe, ok := h.dbTargetSvc.(recordingSQLServerCapabilityProbe)
	if !ok || probe == nil {
		renderRecordingPlanTargetError(c, errRecordingPlanResolverUnavailable)
		return
	}
	capability, err := probe.ProbeSQLServerRecordingAtRevision(c.Request.Context(), connectorID, expectedRevision)
	if !finishStudioLocatorNativeRead(c, h.workspaceSvc, locatorScope) {
		return
	}
	if err != nil {
		// Raw database and TLS failures render as a typed, generic unavailable;
		// the operator response never carries the driver's error text.
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

	// The static matrix keeps kind=sqlserver unsupported while every mutating
	// operation stays masked; the probe carries only server-derived facts and
	// operation gates. Client query claims are never copied into either.
	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey: sqlServerRecordingCapabilityResponse{
			ConnectorCapability:          h.planSvc.GetConnectorCapability(string(schema.DatabaseConnectorKindSQLServer)),
			SQLServerRecordingCapability: capability,
		},
	})
}

// recordingSQLServerCapabilityProbe is the narrow optional interface a saved
// connector service may implement. It is intentionally separate from
// recordingSavedConnectorResolver so existing fakes and callers keep working.
type recordingSQLServerCapabilityProbe interface {
	ProbeSQLServerRecordingAtRevision(ctx context.Context, connectorID, expectedRevision string) (*dbtarget.SQLServerRecordingCapability, error)
}

// sqlServerRecordingCapabilityResponse flattens the static kind matrix and the
// server-derived probe into one object, so a client can see the connector is
// connected and SQL Server-kind while every mutating operation stays gated.
type sqlServerRecordingCapabilityResponse struct {
	recordingplan.ConnectorCapability
	*dbtarget.SQLServerRecordingCapability
}
