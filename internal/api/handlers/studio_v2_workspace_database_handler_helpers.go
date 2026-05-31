package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type studioV2WorkspaceDatabaseConfigRequest struct {
	Kind                 schema.DatabaseConnectorKind `json:"kind"`
	Name                 string                       `json:"name"`
	Host                 string                       `json:"host"`
	Port                 int                          `json:"port"`
	Database             string                       `json:"database"`
	Username             string                       `json:"username"`
	Schema               string                       `json:"schema"`
	Table                string                       `json:"table"`
	WriteMode            schema.DatabaseWriteMode     `json:"write_mode"`
	WriteIntervalSeconds int                          `json:"write_interval_seconds"`
	TimestampColumn      string                       `json:"timestamp_column"`
}

type studioV2WorkspaceDatabaseConfigResponse struct {
	ID                   string                         `json:"id"`
	WorkspaceID          string                         `json:"workspace_id"`
	Kind                 schema.DatabaseConnectorKind   `json:"kind"`
	Name                 string                         `json:"name"`
	Host                 string                         `json:"host"`
	Port                 int                            `json:"port"`
	Database             string                         `json:"database"`
	Username             string                         `json:"username"`
	Schema               string                         `json:"schema"`
	Table                string                         `json:"table"`
	WriteMode            schema.DatabaseWriteMode       `json:"write_mode"`
	WriteIntervalSeconds int                            `json:"write_interval_seconds"`
	TimestampColumn      string                         `json:"timestamp_column"`
	Status               schema.DatabaseConnectorStatus `json:"status"`
	SaveState            string                         `json:"save_state"`
	RuntimeApplyStatus   string                         `json:"runtime_apply_status,omitempty"`
	RuntimeApplyMessage  string                         `json:"runtime_apply_message,omitempty"`
	CreatedAt            any                            `json:"created_at"`
	UpdatedAt            any                            `json:"updated_at"`
}

type studioV2WorkspaceDatabaseTargetRequest struct {
	ColumnName string `json:"column_name"`
	Enabled    bool   `json:"enabled"`
}

type studioV2WorkspaceDatabaseTargetResponse struct {
	ID                  string `json:"id"`
	WorkspaceID         string `json:"workspace_id"`
	PointID             string `json:"point_id"`
	TagID               string `json:"tag_id"`
	ColumnName          string `json:"column_name"`
	Enabled             bool   `json:"enabled"`
	SaveState           string `json:"save_state"`
	RuntimeApplyStatus  string `json:"runtime_apply_status,omitempty"`
	RuntimeApplyMessage string `json:"runtime_apply_message,omitempty"`
	CreatedAt           any    `json:"created_at"`
	UpdatedAt           any    `json:"updated_at"`
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

func (h *StudioV2WorkspaceDatabaseHandler) saveWorkspaceConnector(c *gin.Context, record *workspace.Record, req studioV2WorkspaceDatabaseConfigRequest) (*schema.DatabaseConnector, error) {
	connectionConfig := dbtarget.ConnectionConfig{
		"host":             strings.TrimSpace(req.Host),
		"port":             req.Port,
		"user":             strings.TrimSpace(req.Username),
		"database":         strings.TrimSpace(req.Database),
		"schema":           strings.TrimSpace(req.Schema),
		"table":            strings.TrimSpace(req.Table),
		"write_mode":       string(req.WriteMode),
		"timestamp_column": strings.TrimSpace(req.TimestampColumn),
	}
	if req.Kind == schema.DatabaseConnectorKindSQLite {
		connectionConfig["dsn"] = strings.TrimSpace(req.Database)
	}

	if strings.TrimSpace(record.DatabaseConnectorID) == "" {
		saved, err := h.connectorSvc.Create(c.Request.Context(), dbtarget.CreateConnectorRequest{
			Name:                        strings.TrimSpace(req.Name),
			Kind:                        req.Kind,
			ConnectionConfig:            connectionConfig,
			DefaultWriteIntervalSeconds: workspaceOptionalInt(req.WriteIntervalSeconds),
		})
		if err != nil {
			return nil, err
		}
		if _, err := h.workspaceSvc.BindDatabaseConnector(c.Request.Context(), saved.ID); err != nil {
			if cleanupErr := h.connectorSvc.Delete(c.Request.Context(), saved.ID); cleanupErr != nil {
				return nil, fmt.Errorf("bind workspace database connector failed and cleanup failed: bind=%w cleanup=%v", err, cleanupErr)
			}
			return nil, err
		}
		return saved, nil
	}

	return h.connectorSvc.Update(c.Request.Context(), record.DatabaseConnectorID, dbtarget.UpdateConnectorRequest{
		Name:                        workspaceOptionalString(strings.TrimSpace(req.Name)),
		Kind:                        &req.Kind,
		ConnectionConfig:            &connectionConfig,
		DefaultWriteIntervalSeconds: workspaceOptionalInt(req.WriteIntervalSeconds),
	})
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

func buildWorkspaceDatabaseConfigResponse(workspaceID string, connector *schema.DatabaseConnector) studioV2WorkspaceDatabaseConfigResponse {
	return studioV2WorkspaceDatabaseConfigResponse{
		ID:                   connector.ID,
		WorkspaceID:          workspaceID,
		Kind:                 connector.Kind,
		Name:                 connector.Name,
		Host:                 connectorConfigString(connector, "host"),
		Port:                 connectorConfigInt(connector, "port"),
		Database:             connectorConfigString(connector, "database"),
		Username:             connectorConfigString(connector, "user"),
		Schema:               connectorConfigString(connector, "schema"),
		Table:                connectorConfigString(connector, "table"),
		WriteMode:            connectorConfigWriteMode(connector),
		WriteIntervalSeconds: connector.DefaultWriteIntervalSeconds,
		TimestampColumn:      connectorConfigString(connector, "timestamp_column"),
		Status:               connector.Status,
		SaveState:            "saved",
		CreatedAt:            connector.CreatedAt,
		UpdatedAt:            connector.UpdatedAt,
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
	value, _ := strconv.Atoi(stringValueFromAny(response.ConnectionConfig[key]))
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
		data, _ := json.Marshal(typed)
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

func renderStudioV2WorkspaceDatabaseError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, workspace.ErrValidation), errors.Is(err, sourcerule.ErrValidation), errors.Is(err, dbtarget.ErrValidation):
		renderStudioV2WorkspaceValidationError(c, err)
	default:
		c.JSON(500, gin.H{"success": false, "error": gin.H{"message": "Studio V2 database operation failed"}})
	}
}
