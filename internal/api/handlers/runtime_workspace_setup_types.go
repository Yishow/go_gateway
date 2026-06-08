package handlers

import (
	"context"

	"go-gateway/internal/datalink/dbtarget"
	mappingpkg "go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

type runtimeWorkspaceSourceRuleReader interface {
	ListByDeviceIDs(ctx context.Context, deviceIDs []string) ([]*schema.SourceRule, error)
	ListLinks(ctx context.Context, ruleID string) ([]*schema.SourceRuleLink, error)
}

type runtimeWorkspaceTagReader interface {
	GetByID(ctx context.Context, id string) (*schema.Tag, error)
}

type runtimeWorkspaceMappingReader interface {
	GetByID(ctx context.Context, id string) (*schema.Mapping, error)
	List(ctx context.Context, filter mappingpkg.ListFilter) ([]*schema.Mapping, error)
}

type runtimeWorkspaceDatabaseConnectorReader interface {
	GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error)
}

type runtimeWorkspaceDatabaseTargetReader interface {
	List(ctx context.Context, filter dbtarget.TargetMappingListFilter) ([]*schema.DatabaseTargetMapping, error)
}

type runtimeWorkspaceSetupContextResponse struct {
	ReadinessSummary *workspace.ReadinessSummary                  `json:"readiness_summary,omitempty"`
	SourceRules      []runtimeWorkspaceSetupSourceRuleResponse    `json:"source_rules"`
	Mappings         []runtimeWorkspaceSetupMappingResponse       `json:"mappings"`
	DatabaseConfig   *studioV2WorkspaceDatabaseConfigResponse     `json:"database_config,omitempty"`
	DatabaseTargets  []runtimeWorkspaceSetupDatabaseTargetMapping `json:"database_targets"`
}

type runtimeWorkspaceSetupSourceRuleResponse struct {
	ID             string           `json:"id"`
	DeviceID       string           `json:"device_id"`
	StartAddress   string           `json:"start_address"`
	Count          int              `json:"count"`
	DataType       schema.DataType  `json:"data_type"`
	NamingPrefix   string           `json:"naming_prefix"`
	Enabled        bool             `json:"enabled"`
	RevisionID     string           `json:"revision_id,omitempty"`
	TargetDataType *schema.DataType `json:"target_data_type,omitempty"`
}

type runtimeWorkspaceSetupMappingResponse struct {
	ID          string               `json:"id,omitempty"`
	RuleID      string               `json:"rule_id"`
	DeviceID    string               `json:"device_id"`
	PointID     string               `json:"point_id"`
	Address     string               `json:"address"`
	TagID       string               `json:"tag_id,omitempty"`
	TagKey      string               `json:"tag_key,omitempty"`
	DisplayName string               `json:"display_name,omitempty"`
	Unit        string               `json:"unit,omitempty"`
	TargetType  schema.DataType      `json:"target_type,omitempty"`
	Enabled     bool                 `json:"enabled"`
	Status      schema.MappingStatus `json:"status,omitempty"`
}

type runtimeWorkspaceSetupDatabaseTargetMapping struct {
	ID                   string                   `json:"id"`
	TagID                string                   `json:"tag_id"`
	ConnectorID          string                   `json:"connector_id"`
	TableSchema          string                   `json:"table_schema"`
	TableName            string                   `json:"table_name"`
	ColumnName           string                   `json:"column_name"`
	WriteMode            schema.DatabaseWriteMode `json:"write_mode"`
	TimestampColumn      *string                  `json:"timestamp_column,omitempty"`
	GroupKey             *string                  `json:"group_key,omitempty"`
	WriteIntervalSeconds *int                     `json:"write_interval_seconds,omitempty"`
	Enabled              bool                     `json:"enabled"`
}
