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
	ID               string                  `json:"id" db:"id"`
	Name             string                  `json:"name" db:"name"`
	Kind             DatabaseConnectorKind   `json:"kind" db:"kind"`
	ConnectionConfig string                  `json:"connection_config" db:"connection_config"`
	Status           DatabaseConnectorStatus `json:"status" db:"status"`
	LastCheckAt      *time.Time              `json:"last_check_at,omitempty" db:"last_check_at"`
	LastCheckError   string                  `json:"last_check_error,omitempty" db:"last_check_error"`
	Enabled          bool                    `json:"enabled" db:"enabled"`
	CreatedAt        time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at" db:"updated_at"`
}

type DatabaseTargetMapping struct {
	ID              string            `json:"id" db:"id"`
	TagID           string            `json:"tag_id" db:"tag_id"`
	ConnectorID     string            `json:"connector_id" db:"connector_id"`
	TableSchema     string            `json:"table_schema" db:"table_schema"`
	TableName       string            `json:"table_name" db:"table_name"`
	ColumnName      string            `json:"column_name" db:"column_name"`
	WriteMode       DatabaseWriteMode `json:"write_mode" db:"write_mode"`
	TimestampColumn *string           `json:"timestamp_column,omitempty" db:"timestamp_column"`
	Enabled         bool              `json:"enabled" db:"enabled"`
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" db:"updated_at"`
}
