package handlers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type studioV2WorkspaceDatabaseConfigRequest struct {
	ConnectorID               string                       `json:"connector_id"`
	ExpectedConnectorRevision string                       `json:"expected_connector_revision"`
	Kind                      schema.DatabaseConnectorKind `json:"kind"`
	Name                      string                       `json:"name"`
	Host                      string                       `json:"host"`
	Port                      int                          `json:"port"`
	Database                  string                       `json:"database"`
	Username                  string                       `json:"username"`
	Password                  string                       `json:"password"`
	ClearPassword             bool                         `json:"clear_password"`
	Schema                    string                       `json:"schema"`
	Table                     string                       `json:"table"`
	WriteMode                 schema.DatabaseWriteMode     `json:"write_mode"`
	WriteIntervalSeconds      int                          `json:"write_interval_seconds"`
	TimestampColumn           string                       `json:"timestamp_column"`
	RowGroups                 []workspace.DatabaseRowGroup `json:"row_groups"`
	ExpectedSetupRevision     string                       `json:"expected_setup_revision"`
}

func workspaceDatabasePasswordRequired(kind schema.DatabaseConnectorKind) bool {
	return kind != schema.DatabaseConnectorKindSQLite
}

type studioV2WorkspaceDatabaseConfigResponse struct {
	ID                     string                         `json:"id"`
	IdentityRevision       string                         `json:"identity_revision"`
	SetupRevision          string                         `json:"setup_revision,omitempty"`
	WorkspaceID            string                         `json:"workspace_id"`
	Kind                   schema.DatabaseConnectorKind   `json:"kind"`
	Name                   string                         `json:"name"`
	Host                   string                         `json:"host"`
	Port                   int                            `json:"port"`
	Database               string                         `json:"database"`
	Username               string                         `json:"username"`
	Schema                 string                         `json:"schema"`
	Table                  string                         `json:"table"`
	WriteMode              schema.DatabaseWriteMode       `json:"write_mode"`
	WriteIntervalSeconds   int                            `json:"write_interval_seconds"`
	TimestampColumn        string                         `json:"timestamp_column"`
	Status                 schema.DatabaseConnectorStatus `json:"status"`
	LastSchemaEnsureAt     *time.Time                     `json:"last_schema_ensure_at,omitempty"`
	LastSchemaEnsureStatus string                         `json:"last_schema_ensure_status,omitempty"`
	LastSchemaEnsureError  string                         `json:"last_schema_ensure_error,omitempty"`
	LastWriteAt            *time.Time                     `json:"last_write_at,omitempty"`
	LastWriteStatus        string                         `json:"last_write_status,omitempty"`
	LastWriteError         string                         `json:"last_write_error,omitempty"`
	LastFlushAt            *time.Time                     `json:"last_flush_at,omitempty"`
	LastFlushStatus        string                         `json:"last_flush_status,omitempty"`
	LastFlushError         string                         `json:"last_flush_error,omitempty"`
	SaveState              string                         `json:"save_state"`
	RuntimeApplyStatus     string                         `json:"runtime_apply_status,omitempty"`
	RuntimeApplyMessage    string                         `json:"runtime_apply_message,omitempty"`
	RuntimeApplyIssues     []workspace.ReadinessIssue     `json:"runtime_apply_issues,omitempty"`
	RowGroups              []workspace.DatabaseRowGroup   `json:"row_groups,omitempty"`
	CreatedAt              any                            `json:"created_at"`
	UpdatedAt              any                            `json:"updated_at"`
}

type studioV2WorkspaceDatabaseTargetRequest struct {
	ColumnName            string `json:"column_name"`
	Enabled               bool   `json:"enabled"`
	RowGroupID            string `json:"row_group_id"`
	ExpectedSetupRevision string `json:"expected_setup_revision"`
}

type studioV2WorkspaceDatabaseTargetResponse struct {
	ID                  string                     `json:"id"`
	WorkspaceID         string                     `json:"workspace_id"`
	PointID             string                     `json:"point_id"`
	TagID               string                     `json:"tag_id"`
	ColumnName          string                     `json:"column_name"`
	Enabled             bool                       `json:"enabled"`
	RowGroupID          string                     `json:"row_group_id,omitempty"`
	SaveState           string                     `json:"save_state"`
	SetupRevision       string                     `json:"setup_revision,omitempty"`
	RuntimeApplyStatus  string                     `json:"runtime_apply_status,omitempty"`
	RuntimeApplyMessage string                     `json:"runtime_apply_message,omitempty"`
	RuntimeApplyIssues  []workspace.ReadinessIssue `json:"runtime_apply_issues,omitempty"`
	CreatedAt           any                        `json:"created_at"`
	UpdatedAt           any                        `json:"updated_at"`
}

type workspacePointBinding struct {
	PointID string
	TagID   string
}

func (h *StudioV2WorkspaceDatabaseHandler) parseWorkspaceDatabaseConfigRequest(c *gin.Context) (*workspace.Record, studioV2WorkspaceDatabaseConfigRequest, bool) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return nil, studioV2WorkspaceDatabaseConfigRequest{}, false
	}

	var req studioV2WorkspaceDatabaseConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return nil, studioV2WorkspaceDatabaseConfigRequest{}, false
	}
	if err := validateWorkspaceDatabaseConfigRequest(req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return nil, studioV2WorkspaceDatabaseConfigRequest{}, false
	}
	return record, req, true
}

func (h *StudioV2WorkspaceDatabaseHandler) resolveWorkspaceDatabaseTarget(c *gin.Context) (*workspace.Record, *schema.DatabaseConnector, workspacePointBinding, studioV2WorkspaceDatabaseTargetRequest, bool) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return nil, nil, workspacePointBinding{}, studioV2WorkspaceDatabaseTargetRequest{}, false
	}
	if strings.TrimSpace(record.DatabaseConnectorID) == "" {
		renderStudioV2WorkspaceValidationError(c, errors.New("workspace database config not found"))
		return nil, nil, workspacePointBinding{}, studioV2WorkspaceDatabaseTargetRequest{}, false
	}

	var req studioV2WorkspaceDatabaseTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return nil, nil, workspacePointBinding{}, studioV2WorkspaceDatabaseTargetRequest{}, false
	}
	if strings.TrimSpace(req.ColumnName) == "" {
		renderStudioV2WorkspaceValidationError(c, errors.New("column_name is required"))
		return nil, nil, workspacePointBinding{}, studioV2WorkspaceDatabaseTargetRequest{}, false
	}

	connector, err := h.connectorSvc.GetByID(c.Request.Context(), record.DatabaseConnectorID)
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return nil, nil, workspacePointBinding{}, studioV2WorkspaceDatabaseTargetRequest{}, false
	}
	bindings, err := h.listWorkspacePointBindings(c)
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return nil, nil, workspacePointBinding{}, studioV2WorkspaceDatabaseTargetRequest{}, false
	}
	pointID := c.Param("point_id")
	binding, ok := bindings[pointID]
	if !ok || binding.TagID == "" {
		renderStudioV2WorkspaceValidationError(c, errors.New("workspace database point is not bound to a persisted tag"))
		return nil, nil, workspacePointBinding{}, studioV2WorkspaceDatabaseTargetRequest{}, false
	}
	return record, connector, binding, req, true
}

func (h *StudioV2WorkspaceDatabaseHandler) listWorkspacePointBindings(c *gin.Context) (map[string]workspacePointBinding, error) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		return nil, err
	}
	rules, err := h.ruleSvc.ListByDeviceIDs(c.Request.Context(), record.OrderedDeviceIDs)
	if err != nil {
		return nil, err
	}

	bindings := make(map[string]workspacePointBinding)
	for _, rule := range rules {
		links, err := h.ruleSvc.ListLinks(c.Request.Context(), rule.ID)
		if err != nil {
			return nil, err
		}
		for _, link := range links {
			if link != nil && link.TagID != nil && strings.TrimSpace(*link.TagID) != "" {
				bindings[link.PointID] = workspacePointBinding{PointID: link.PointID, TagID: *link.TagID}
			}
		}
	}
	return bindings, nil
}

func buildWorkspaceDatabaseConfigResponse(workspaceID string, connector *schema.DatabaseConnector, rowGroups []workspace.DatabaseRowGroup) studioV2WorkspaceDatabaseConfigResponse {
	return studioV2WorkspaceDatabaseConfigResponse{
		ID:                     connector.ID,
		IdentityRevision:       connector.IdentityRevision,
		WorkspaceID:            workspaceID,
		Kind:                   connector.Kind,
		Name:                   connector.Name,
		Host:                   connectorConfigString(connector, "host"),
		Port:                   connectorConfigInt(connector, "port"),
		Database:               connectorConfigString(connector, "database"),
		Username:               connectorConfigString(connector, "user"),
		Schema:                 connectorConfigString(connector, "schema"),
		Table:                  connectorConfigString(connector, "table"),
		WriteMode:              connectorConfigWriteMode(connector),
		WriteIntervalSeconds:   connector.DefaultWriteIntervalSeconds,
		TimestampColumn:        connectorConfigString(connector, "timestamp_column"),
		Status:                 connector.Status,
		LastSchemaEnsureAt:     connector.LastSchemaEnsureAt,
		LastSchemaEnsureStatus: connector.LastSchemaEnsureStatus,
		LastSchemaEnsureError:  connector.LastSchemaEnsureError,
		LastWriteAt:            connector.LastWriteAt,
		LastWriteStatus:        connector.LastWriteStatus,
		LastWriteError:         connector.LastWriteError,
		LastFlushAt:            connector.LastFlushAt,
		LastFlushStatus:        connector.LastFlushStatus,
		LastFlushError:         connector.LastFlushError,
		SaveState:              workspaceSaveStateSaved,
		RowGroups:              rowGroups,
		CreatedAt:              connector.CreatedAt,
		UpdatedAt:              connector.UpdatedAt,
	}
}

func validateWorkspaceDatabaseConfigRequest(req studioV2WorkspaceDatabaseConfigRequest) error {
	switch {
	case strings.TrimSpace(req.Name) == "":
		return errors.New("name is required")
	case strings.TrimSpace(req.Database) == "":
		return errors.New("database is required")
	case strings.TrimSpace(req.Table) == "":
		return errors.New("table is required")
	case req.WriteIntervalSeconds <= 0:
		return errors.New("write_interval_seconds must be greater than zero")
	}
	if req.Kind == schema.DatabaseConnectorKindPostgres && strings.TrimSpace(req.Username) == "" {
		return errors.New("username is required")
	}
	return nil
}

func connectorConfigString(connector *schema.DatabaseConnector, key string) string {
	response := toDatabaseConnectorResponse(connector)
	return strings.TrimSpace(stringValueFromAny(response.ConnectionConfig[key]))
}

func connectorConfigInt(connector *schema.DatabaseConnector, key string) int {
	response := toDatabaseConnectorResponse(connector)
	value, err := strconv.Atoi(stringValueFromAny(response.ConnectionConfig[key]))
	if err != nil {
		return 0
	}
	return value
}

func connectorConfigWriteMode(connector *schema.DatabaseConnector) schema.DatabaseWriteMode {
	value := connectorConfigString(connector, "write_mode")
	if value == "" {
		return schema.DatabaseWriteModeInsert
	}
	return schema.DatabaseWriteMode(value)
}

func stringValueFromAny(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return strconv.Itoa(int(typed))
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	case json.Number:
		return typed.String()
	case nil:
		return ""
	default:
		data, err := json.Marshal(typed)
		if err != nil {
			return ""
		}
		return strings.Trim(string(data), "\"")
	}
}

func workspaceOptionalString(value string) *string {
	copyValue := value
	return &copyValue
}

func workspaceOptionalInt(value int) *int {
	copyValue := value
	return &copyValue
}

func workspaceOptionalWriteMode(value schema.DatabaseWriteMode) *schema.DatabaseWriteMode {
	copyValue := value
	return &copyValue
}
