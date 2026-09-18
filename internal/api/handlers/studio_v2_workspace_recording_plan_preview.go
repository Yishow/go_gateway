package handlers

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

var (
	errRecordingSchemaPreviewStale    = errors.New("recording schema preview scope is stale")
	errRecordingSchemaPrefixAmbiguous = errors.New("recording plan destination has several table prefixes")
	errRecordingTargetInspectionEmpty = errors.New("target inspection returned no result")
)

const (
	recordingSchemaPreviewStaleCode = "RECORDING_SCHEMA_PREVIEW_STALE"
	connectorConfigDatabaseKey      = "database"
)

// expectedPlanRevisions binds the workspace, plan and connector revisions.
type expectedPlanRevisions struct {
	ExpectedWorkspaceRevision string `json:"expected_workspace_revision" binding:"required"`
	ExpectedPlanRevision      string `json:"expected_plan_revision" binding:"required"`
	ExpectedConnectorRevision string `json:"expected_connector_revision" binding:"required"`
}

func (r *expectedPlanRevisions) normalize() {
	r.ExpectedWorkspaceRevision = strings.TrimSpace(r.ExpectedWorkspaceRevision)
	r.ExpectedPlanRevision = strings.TrimSpace(r.ExpectedPlanRevision)
	r.ExpectedConnectorRevision = strings.TrimSpace(r.ExpectedConnectorRevision)
}

func (r *expectedPlanRevisions) complete() bool {
	return r.ExpectedWorkspaceRevision != "" && r.ExpectedPlanRevision != "" && r.ExpectedConnectorRevision != ""
}

type schemaPreviewRequest struct {
	PlanID      string `json:"plan_id" binding:"required"`
	ConnectorID string `json:"connector_id" binding:"required"`
	expectedPlanRevisions
	TablePrefix string `json:"table_prefix"`
	Dialect     string `json:"dialect"`
}

func (r *schemaPreviewRequest) normalize() {
	r.PlanID = strings.TrimSpace(r.PlanID)
	r.ConnectorID = strings.TrimSpace(r.ConnectorID)
	r.expectedPlanRevisions.normalize()
	r.TablePrefix = strings.TrimSpace(r.TablePrefix)
}

func (r *schemaPreviewRequest) complete() bool {
	return r.PlanID != "" && r.ConnectorID != "" && r.expectedPlanRevisions.complete()
}

// SchemaPreview inspects the saved target and issues a revision-bound preview.
// Stale workspace or plan revisions are rejected before the target is read,
// the connector decides the dialect, and the preview never writes to the target.
func (h *StudioV2WorkspaceRecordingPlansHandler) SchemaPreview(c *gin.Context) {
	var req schemaPreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}
	req.normalize()
	if !req.complete() {
		renderStudioV2WorkspaceValidationError(c, errors.New("plan_id, connector_id and the expected connector, workspace and plan revisions are required"))
		return
	}

	ctx := c.Request.Context()
	record, err := h.workspaceSvc.GetOrCreate(ctx)
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}
	plan, err := h.getPlanForWorkspace(ctx, record.ID, req.PlanID)
	if err != nil {
		renderRecordingPlanReadError(c, err)
		return
	}
	if record.DatabaseSetupRevision != req.ExpectedWorkspaceRevision || plan.Revision != req.ExpectedPlanRevision {
		renderRecordingSchemaPreviewError(c, errRecordingSchemaPreviewStale)
		return
	}
	tablePrefix, err := savedPreviewTablePrefix(plan, req.ConnectorID, req.TablePrefix)
	if err != nil {
		renderRecordingSchemaPreviewError(c, err)
		return
	}
	connector, err := h.resolvePreviewConnector(ctx, record, plan, req.ConnectorID, req.ExpectedConnectorRevision)
	if err != nil {
		renderRecordingPlanTargetError(c, err)
		return
	}
	if !recordingDialectMatches(string(connector.Kind), req.Dialect) {
		renderRecordingPlanTargetError(c, errRecordingPlanDialectMismatch)
		return
	}

	scope := recordingPreviewScope(record, plan, connector, tablePrefix)
	token, err := h.planSvc.PrepareSchemaPreview(ctx, scope, h.recordingTargetInspector(connector.ID, scope.Schema))
	if err != nil {
		renderRecordingSchemaPreviewError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: token})
}

// savedPreviewTablePrefix returns the table prefix saved on the plan
// destination for the connector. A request naming another prefix is stale;
// without a saved prefix the requested one, or the default, is used.
func savedPreviewTablePrefix(plan *recordingplan.RecordingPlan, connectorID, requested string) (string, error) {
	saved := ""
	for _, destination := range plan.Destinations {
		prefix := strings.TrimSpace(destination.TablePrefix)
		if strings.TrimSpace(destination.ConnectorID) != connectorID || prefix == "" {
			continue
		}
		if saved != "" && saved != prefix {
			return "", errRecordingSchemaPrefixAmbiguous
		}
		saved = prefix
	}
	switch {
	case saved == "":
		return requested, nil
	case requested != "" && requested != saved:
		return "", errRecordingSchemaPreviewStale
	default:
		return saved, nil
	}
}

func recordingPreviewScope(record *workspace.Record, plan *recordingplan.RecordingPlan, connector *schema.DatabaseConnector, tablePrefix string) recordingplan.SchemaPreviewScope {
	dialect := canonicalRecordingDialect(string(connector.Kind))
	return recordingplan.SchemaPreviewScope{
		WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision,
		PlanID: plan.ID, PlanRevision: plan.Revision,
		ConnectorID: connector.ID, ConnectorRevision: strings.TrimSpace(connector.IdentityRevision),
		Dialect: dialect, Database: recordingTargetDatabase(connector, dialect), Schema: recordingTargetSchema(connector, dialect),
		TablePrefix: tablePrefix,
	}
}

// recordingTargetDatabase names the saved target database without secrets:
// the file location for SQLite and the database name for server adapters.
func recordingTargetDatabase(connector *schema.DatabaseConnector, dialect string) string {
	keys := []string{connectorConfigDatabaseKey}
	if dialect == recordingDialectSQLite {
		keys = []string{"path", "dsn", connectorConfigDatabaseKey}
	}
	for _, key := range keys {
		value := connectorConfigString(connector, key)
		if index := strings.IndexByte(value, '?'); index >= 0 {
			value = value[:index]
		}
		if value != "" {
			if key == "dsn" {
				value = sanitizeDSNValue(value)
			}
			return value
		}
	}
	return ""
}

var dsnCredentialPattern = regexp.MustCompile(`(?i)\b(?:password|pwd)=\S*`)

// sanitizeDSNValue strips credentials from a connection string, so a database
// name derived from a DSN can never carry them into a token or response.
func sanitizeDSNValue(value string) string {
	value = dsnCredentialPattern.ReplaceAllString(value, "")
	if schemeEnd := strings.Index(value, "://"); schemeEnd >= 0 {
		authority := value[schemeEnd+3:]
		rest := ""
		if slash := strings.IndexAny(authority, "/?"); slash >= 0 {
			rest = authority[slash:]
			authority = authority[:slash]
		}
		if at := strings.LastIndex(authority, "@"); at >= 0 {
			value = value[:schemeEnd+3] + authority[at+1:] + rest
		}
	}
	return strings.Join(strings.Fields(value), " ")
}

// recordingTargetSchema returns the saved schema or the adapter default that
// table inspection uses when none is configured.
func recordingTargetSchema(connector *schema.DatabaseConnector, dialect string) string {
	if value := connectorConfigString(connector, "schema"); value != "" {
		return value
	}
	switch dialect {
	case recordingDialectSQLite:
		return "main"
	case recordingDialectPostgres:
		return "public"
	default:
		return ""
	}
}

// recordingTargetInspector reads managed tables through the saved target
// service; table existence, column names and data types reach the preview.
func (h *StudioV2WorkspaceRecordingPlansHandler) recordingTargetInspector(connectorID, schemaName string) recordingplan.TargetInspector {
	return func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
		inspection, err := h.dbTargetSvc.InspectTable(ctx, connectorID, schemaName, table)
		if err != nil {
			return recordingplan.TargetTableInspection{}, err
		}
		if inspection == nil {
			return recordingplan.TargetTableInspection{}, errRecordingTargetInspectionEmpty
		}
		columns := make([]string, 0, len(inspection.Columns))
		columnTypes := make(map[string]string, len(inspection.Columns))
		for _, column := range inspection.Columns {
			columns = append(columns, column.Name)
			columnTypes[strings.ToLower(strings.TrimSpace(column.Name))] = column.DataType
		}
		return recordingplan.TargetTableInspection{
			Status:      string(inspection.Status),
			Columns:     columns,
			ColumnTypes: columnTypes,
		}, nil
	}
}

func renderRecordingSchemaPreviewError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errRecordingSchemaPreviewStale):
		renderRecordingPlanError(c, http.StatusConflict, recordingSchemaPreviewStaleCode, "recording schema preview scope is stale", false, "reload the saved setup and preview again")
	case errors.Is(err, errRecordingSchemaPrefixAmbiguous):
		renderRecordingPlanError(c, http.StatusUnprocessableEntity, recordingPlanDestinationCode, recordingPlanDestinationMessage, false, "keep one table prefix per database destination")
	case errors.Is(err, recordingplan.ErrInvalidTablePrefix):
		renderRecordingPlanError(c, http.StatusUnprocessableEntity, "RECORDING_SCHEMA_PREFIX_INVALID", "managed table prefix is invalid", false, "use up to 40 letters, digits or underscores, not starting with a digit")
	case errors.Is(err, recordingplan.ErrIncompatibleExistingTable):
		renderRecordingPlanError(c, http.StatusUnprocessableEntity, "RECORDING_SCHEMA_INCOMPATIBLE", "an existing managed table is incompatible", false, "choose another table prefix or inspect the existing table manually")
	case errors.Is(err, recordingplan.ErrTargetPermissionDenied):
		renderRecordingPlanError(c, http.StatusUnprocessableEntity, "RECORDING_SCHEMA_PERMISSION_DENIED", "the saved connection cannot read the managed tables", false, "grant table read permission, then preview again")
	case errors.Is(err, recordingplan.ErrTargetInspectionUnconfirmed):
		renderRecordingPlanError(c, http.StatusServiceUnavailable, "RECORDING_SCHEMA_TARGET_UNAVAILABLE", "the target tables could not be checked", true, "check the database connection, then preview again")
	case errors.Is(err, recordingplan.ErrUnsupportedSchemaDialect):
		renderRecordingPlanError(c, http.StatusUnprocessableEntity, recordingPlanPreviewFailedCode, recordingPlanPreviewFailedMessage, false, "review the selected connector")
	default:
		renderRecordingPlanError(c, http.StatusServiceUnavailable, recordingPlanUnavailableCode, recordingPlanUnavailableMessage, true, "retry the schema preview")
	}
}
