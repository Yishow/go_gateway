package dbtarget

import (
	"context"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/schema"
)

var ErrValidation = errors.New("database target validation failed")

func validationError(message string) error {
	return fmt.Errorf("%w: %s", ErrValidation, message)
}

type ConnectionConfig map[string]any

type ConnectorListFilter struct {
	Enabled *bool
}

type TargetMappingListFilter struct {
	ConnectorID *string
	TagID       *string
	Enabled     *bool
}

type CreateConnectorRequest struct {
	Name                        string                       `json:"name"`
	Kind                        schema.DatabaseConnectorKind `json:"kind"`
	ConnectionConfig            ConnectionConfig             `json:"connection_config"`
	Enabled                     *bool                        `json:"enabled,omitempty"`
	DefaultWriteIntervalSeconds *int                         `json:"default_write_interval_seconds,omitempty"`
}

type UpdateConnectorRequest struct {
	Name                        *string                       `json:"name,omitempty"`
	ConnectionConfig            *ConnectionConfig             `json:"connection_config,omitempty"`
	ClearPassword               *bool                         `json:"clear_password,omitempty"`
	Enabled                     *bool                         `json:"enabled,omitempty"`
	Kind                        *schema.DatabaseConnectorKind `json:"kind,omitempty"`
	DefaultWriteIntervalSeconds *int                          `json:"default_write_interval_seconds,omitempty"`
}

type CreateTargetMappingRequest struct {
	TagID                string                   `json:"tag_id"`
	ConnectorID          string                   `json:"connector_id"`
	TableSchema          string                   `json:"table_schema"`
	TableName            string                   `json:"table_name"`
	ColumnName           string                   `json:"column_name"`
	WriteMode            schema.DatabaseWriteMode `json:"write_mode"`
	TimestampColumn      *string                  `json:"timestamp_column,omitempty"`
	GroupKey             *string                  `json:"group_key,omitempty"`
	WriteIntervalSeconds *int                     `json:"write_interval_seconds,omitempty"`
	Enabled              *bool                    `json:"enabled,omitempty"`
}

type UpdateTargetMappingRequest struct {
	TableSchema          *string                   `json:"table_schema,omitempty"`
	TableName            *string                   `json:"table_name,omitempty"`
	ColumnName           *string                   `json:"column_name,omitempty"`
	WriteMode            *schema.DatabaseWriteMode `json:"write_mode,omitempty"`
	TimestampColumn      *string                   `json:"timestamp_column,omitempty"`
	GroupKey             *string                   `json:"group_key,omitempty"`
	WriteIntervalSeconds *int                      `json:"write_interval_seconds,omitempty"`
	Enabled              *bool                     `json:"enabled,omitempty"`
}

type ColumnInfo struct {
	Name       string `json:"name"`
	DataType   string `json:"data_type"`
	Nullable   bool   `json:"nullable"`
	PrimaryKey bool   `json:"primary_key"`
	Unique     bool   `json:"unique"`
}

type TableInfo struct {
	Schema  string       `json:"schema"`
	Name    string       `json:"name"`
	Columns []ColumnInfo `json:"columns"`
}

type ValidationIssue struct {
	Severity  string `json:"severity"`
	MappingID string `json:"mapping_id,omitempty"`
	TagID     string `json:"tag_id,omitempty"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

type ValidationResult struct {
	Ready  bool              `json:"ready"`
	Issues []ValidationIssue `json:"issues"`
}

type ConnectorRepository interface {
	Create(ctx context.Context, connector *schema.DatabaseConnector) error
	Update(ctx context.Context, connector *schema.DatabaseConnector) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error)
	List(ctx context.Context, filter ConnectorListFilter) ([]*schema.DatabaseConnector, error)
}

type TargetMappingRepository interface {
	Create(ctx context.Context, mapping *schema.DatabaseTargetMapping) error
	Update(ctx context.Context, mapping *schema.DatabaseTargetMapping) error
	Delete(ctx context.Context, id string) error
	DeleteByConnectorID(ctx context.Context, connectorID string) error
	GetByID(ctx context.Context, id string) (*schema.DatabaseTargetMapping, error)
	List(ctx context.Context, filter TargetMappingListFilter) ([]*schema.DatabaseTargetMapping, error)
}
