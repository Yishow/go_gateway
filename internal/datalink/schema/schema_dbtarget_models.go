package schema

import "time"

type DatabaseConnectorKind string

const (
	DatabaseConnectorKindSQLite    DatabaseConnectorKind = "sqlite"
	DatabaseConnectorKindPostgres  DatabaseConnectorKind = "postgres"
	DatabaseConnectorKindMySQL     DatabaseConnectorKind = "mysql"
	DatabaseConnectorKindSQLServer DatabaseConnectorKind = "sqlserver"
)

type DatabaseConnectorStatus string

const (
	DatabaseConnectorStatusReady       DatabaseConnectorStatus = "ready"
	DatabaseConnectorStatusUnreachable DatabaseConnectorStatus = "unreachable"
	DatabaseConnectorStatusAuthFailed  DatabaseConnectorStatus = "auth_failed"
	DatabaseConnectorStatusError       DatabaseConnectorStatus = "error"
)

type DatabaseWriteMode string

const (
	DatabaseWriteModeInsert DatabaseWriteMode = "insert"
	DatabaseWriteModeUpsert DatabaseWriteMode = "upsert"
)

type DatabaseConnector struct {
	ID                          string                  `json:"id" db:"id"`
	Name                        string                  `json:"name" db:"name"`
	Kind                        DatabaseConnectorKind   `json:"kind" db:"kind"`
	ConnectionConfig            string                  `json:"connection_config" db:"connection_config"`
	IdentityRevision            string                  `json:"identity_revision" db:"identity_revision"`
	Status                      DatabaseConnectorStatus `json:"status" db:"status"`
	LastCheckAt                 *time.Time              `json:"last_check_at,omitempty" db:"last_check_at"`
	LastCheckError              string                  `json:"last_check_error,omitempty" db:"last_check_error"`
	Enabled                     bool                    `json:"enabled" db:"enabled"`
	DefaultWriteIntervalSeconds int                     `json:"default_write_interval_seconds" db:"default_write_interval_seconds"`
	LastSchemaEnsureAt          *time.Time              `json:"last_schema_ensure_at,omitempty" db:"last_schema_ensure_at"`
	LastSchemaEnsureStatus      string                  `json:"last_schema_ensure_status,omitempty" db:"last_schema_ensure_status"`
	LastSchemaEnsureError       string                  `json:"last_schema_ensure_error,omitempty" db:"last_schema_ensure_error"`
	LastWriteAt                 *time.Time              `json:"last_write_at,omitempty" db:"last_write_at"`
	LastWriteStatus             string                  `json:"last_write_status,omitempty" db:"last_write_status"`
	LastWriteError              string                  `json:"last_write_error,omitempty" db:"last_write_error"`
	LastFlushAt                 *time.Time              `json:"last_flush_at,omitempty" db:"last_flush_at"`
	LastFlushStatus             string                  `json:"last_flush_status,omitempty" db:"last_flush_status"`
	LastFlushError              string                  `json:"last_flush_error,omitempty" db:"last_flush_error"`
	CreatedAt                   time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt                   time.Time               `json:"updated_at" db:"updated_at"`
}

type DatabaseTargetMapping struct {
	ID                   string            `json:"id" db:"id"`
	TagID                string            `json:"tag_id" db:"tag_id"`
	ConnectorID          string            `json:"connector_id" db:"connector_id"`
	TableSchema          string            `json:"table_schema" db:"table_schema"`
	TableName            string            `json:"table_name" db:"table_name"`
	ColumnName           string            `json:"column_name" db:"column_name"`
	WriteMode            DatabaseWriteMode `json:"write_mode" db:"write_mode"`
	TimestampColumn      *string           `json:"timestamp_column,omitempty" db:"timestamp_column"`
	GroupKey             *string           `json:"group_key,omitempty" db:"group_key"`
	WriteIntervalSeconds *int              `json:"write_interval_seconds,omitempty" db:"write_interval_seconds"`
	Enabled              bool              `json:"enabled" db:"enabled"`
	CreatedAt            time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time         `json:"updated_at" db:"updated_at"`
}
