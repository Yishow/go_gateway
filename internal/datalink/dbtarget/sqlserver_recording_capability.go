package dbtarget

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
)

const (
	sqlServerRecordingCapabilityProjectionVersion = "sqlserver-recording-capability-v1"
	sqlServerRecordingMetadataVerified            = "verified"
	sqlServerRecordingMetadataForbidden           = "forbidden"
	sqlServerRecordingMetadataUnconfirmed         = "unconfirmed"
	sqlServerRecordingMetadataMissing             = "missing"
)

const (
	sqlServerRecordingMetadataQuery             = "SELECT CONVERT(nvarchar(128), SERVERPROPERTY('ProductVersion')), CONVERT(nvarchar(128), SERVERPROPERTY('Edition')), CONVERT(int, SERVERPROPERTY('EngineEdition')), DB_NAME(), (SELECT compatibility_level FROM sys.databases WHERE name = DB_NAME())"
	sqlServerRecordingVisibilityPermissionQuery = "SELECT HAS_PERMS_BY_NAME(DB_NAME(), 'DATABASE', 'VIEW DEFINITION'), HAS_PERMS_BY_NAME(@p1, 'SCHEMA', 'VIEW DEFINITION')"
	sqlServerRecordingSchemaQuery               = "SELECT name FROM sys.schemas WHERE name = @p1"
	sqlServerRecordingPermissionQuery           = "SELECT HAS_PERMS_BY_NAME(DB_NAME(), 'DATABASE', 'CREATE TABLE'), HAS_PERMS_BY_NAME(@p1, 'SCHEMA', 'ALTER'), HAS_PERMS_BY_NAME(@p1, 'SCHEMA', 'SELECT'), HAS_PERMS_BY_NAME(@p1, 'SCHEMA', 'INSERT'), HAS_PERMS_BY_NAME(@p1, 'SCHEMA', 'DELETE'), HAS_PERMS_BY_NAME(DB_NAME(), 'DATABASE', 'ALTER ANY DATASPACE')"
)

// SQLServerRecordingCapability is the backend-derived projection for a saved
// SQL Server connector. No operation is enabled before real-server release
// evidence is available.
type SQLServerRecordingCapability struct {
	ConnectorID             string                        `json:"connector_id"`
	ConnectorRevision       string                        `json:"connector_revision"`
	Database                string                        `json:"database"`
	Schema                  string                        `json:"schema"`
	ServerVersion           string                        `json:"server_version"`
	ServerEdition           string                        `json:"server_edition"`
	EngineEdition           int64                         `json:"engine_edition"`
	CompatibilityLevel      int64                         `json:"compatibility_level"`
	Connected               bool                          `json:"connected"`
	MetadataStatus          string                        `json:"metadata_status"`
	Permissions             SQLServerRecordingPermissions `json:"permissions"`
	SupportsManagedLong     bool                          `json:"supports_managed_long"`
	SupportsNativePartition bool                          `json:"supports_native_partition"`
	SupportsGroupTestWrite  bool                          `json:"supports_group_test_write"`
	SupportsMissingRepair   bool                          `json:"supports_missing_repair"`
	CapabilityRevision      string                        `json:"capability_revision"`
	Reasons                 []string                      `json:"reasons"`
}

// SQLServerRecordingPermissions contains read-only results from SQL Server's
// effective permission checks. These facts never enable a recording operation.
type SQLServerRecordingPermissions struct {
	Verified               bool `json:"verified"`
	DatabaseViewDefinition bool `json:"database_view_definition"`
	SchemaViewDefinition   bool `json:"schema_view_definition"`
	CreateTable            bool `json:"create_table"`
	AlterSchema            bool `json:"alter_schema"`
	SelectSchema           bool `json:"select_schema"`
	InsertSchema           bool `json:"insert_schema"`
	DeleteSchema           bool `json:"delete_schema"`
	AlterAnyDataspace      bool `json:"alter_any_dataspace"`
}

type sqlServerRecordingMetadata struct {
	serverVersion      string
	serverEdition      string
	engineEdition      int64
	engineEditionValid bool
	compatibilityLevel int64
	compatibilityValid bool
	database           string
	actualDatabase     string
	schema             string
	actualSchema       string
	metadataStatus     string
	permissions        SQLServerRecordingPermissions
	reasons            []string
}

// ProbeSQLServerRecordingAtRevision inspects only the exact saved SQL Server
// connector revision. It opens that destination directly and never uses a
// probe or provisioning path that could create a database or schema.
func (s *ConnectorService) ProbeSQLServerRecordingAtRevision(
	ctx context.Context,
	connectorID, expectedRevision string,
) (*SQLServerRecordingCapability, error) {
	if s == nil || s.repo == nil {
		return nil, managedDestinationInvalid("saved connector service is unavailable")
	}
	connector, err := s.ResolveSavedTarget(ctx, strings.TrimSpace(connectorID), strings.TrimSpace(expectedRevision))
	if err != nil {
		return nil, err
	}
	if connector.Kind != schema.DatabaseConnectorKindSQLServer {
		return nil, managedDestinationInvalid("saved connector kind is not sqlserver")
	}

	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return nil, managedDestinationInvalid("saved sqlserver connector configuration is invalid")
	}
	metadata := sqlServerRecordingMetadata{
		database:       strings.TrimSpace(stringConfigValue(config, "database")),
		schema:         strings.TrimSpace(stringConfigValue(config, "schema")),
		metadataStatus: sqlServerRecordingMetadataUnconfirmed,
		reasons:        []string{},
	}
	settings, err := sqlServerConnectionSettingsFor(config)
	if err != nil {
		metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, "connection_settings_unverified")
		return projectSQLServerRecordingCapability(connector, metadata, false), nil
	}
	metadata.database, metadata.schema = settings.Database, settings.Schema
	if settings.DSN == "" || settings.Database == "" || settings.Schema == "" {
		metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, "connection_settings_unverified")
		return projectSQLServerRecordingCapability(connector, metadata, false), nil
	}

	// Do not call probeConnector: generic probing may have provider-specific
	// side effects and is not the recording capability path.
	manager, err := openExternalDBManagerFunc(connector.Kind, config)
	if err != nil {
		metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, sqlServerConnectionFailureReason(err))
		return projectSQLServerRecordingCapability(connector, metadata, false), nil
	}
	if manager == nil || manager.DB() == nil {
		if manager != nil {
			_ = manager.Close()
		}
		metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, inspectionReasonConnectionFailed)
		return projectSQLServerRecordingCapability(connector, metadata, false), nil
	}
	defer manager.Close()

	metadata = inspectSQLServerRecordingMetadata(ctx, manager.DB(), settings.Database, settings.Schema)
	return projectSQLServerRecordingCapability(connector, metadata, true), nil
}

func inspectSQLServerRecordingMetadata(
	ctx context.Context,
	db *sql.DB,
	expectedDatabase, expectedSchema string,
) sqlServerRecordingMetadata {
	if db == nil {
		return sqlServerRecordingMetadata{
			database: expectedDatabase, schema: expectedSchema,
			metadataStatus: sqlServerRecordingMetadataUnconfirmed, reasons: []string{"metadata_unconfirmed"},
		}
	}
	return inspectSQLServerRecordingMetadataWithQuerier(ctx, db, expectedDatabase, expectedSchema)
}

func inspectSQLServerRecordingMetadataWithQuerier(
	ctx context.Context,
	db interface {
		QueryRowContext(context.Context, string, ...any) *sql.Row
	},
	expectedDatabase, expectedSchema string,
) (metadata sqlServerRecordingMetadata) {
	metadata = sqlServerRecordingMetadata{
		database:       expectedDatabase,
		schema:         expectedSchema,
		metadataStatus: sqlServerRecordingMetadataUnconfirmed,
		reasons:        []string{},
	}
	if expectedDatabase == "" || expectedSchema == "" {
		metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, "destination_scope_unconfirmed")
		return metadata
	}

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var engineEdition, compatibilityLevel sql.NullInt64
	if err := db.QueryRowContext(queryCtx, sqlServerRecordingMetadataQuery).Scan(
		&metadata.serverVersion,
		&metadata.serverEdition,
		&engineEdition,
		&metadata.actualDatabase,
		&compatibilityLevel,
	); err != nil {
		setSQLServerRecordingMetadataFailure(&metadata, err)
		return metadata
	}
	if engineEdition.Valid {
		metadata.engineEdition = engineEdition.Int64
		metadata.engineEditionValid = true
	}
	if compatibilityLevel.Valid {
		metadata.compatibilityLevel = compatibilityLevel.Int64
		metadata.compatibilityValid = true
	}
	if metadata.actualDatabase == "" {
		metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, "database_scope_unconfirmed")
		return metadata
	}
	if metadata.actualDatabase != expectedDatabase {
		metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, "database_scope_mismatch")
		return metadata
	}
	if metadata.serverVersion == "" || metadata.serverEdition == "" || !metadata.engineEditionValid || !metadata.compatibilityValid {
		metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, "server_metadata_unconfirmed")
		return metadata
	}

	var databaseViewDefinition, schemaViewDefinition sql.NullBool
	if err := db.QueryRowContext(queryCtx, sqlServerRecordingVisibilityPermissionQuery, expectedSchema).Scan(
		&databaseViewDefinition, &schemaViewDefinition,
	); err != nil {
		setSQLServerRecordingMetadataFailure(&metadata, err)
		return metadata
	}
	metadata.permissions.DatabaseViewDefinition = databaseViewDefinition.Bool
	metadata.permissions.SchemaViewDefinition = schemaViewDefinition.Bool

	var visibleSchema string
	err := db.QueryRowContext(queryCtx, sqlServerRecordingSchemaQuery, expectedSchema).Scan(&visibleSchema)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		if databaseViewDefinition.Valid && databaseViewDefinition.Bool {
			metadata.metadataStatus = sqlServerRecordingMetadataMissing
			metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, "schema_missing")
		} else {
			metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, "schema_visibility_unconfirmed")
		}
		return metadata
	case err != nil:
		setSQLServerRecordingMetadataFailure(&metadata, err)
		return metadata
	case visibleSchema != expectedSchema:
		metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, "schema_scope_mismatch")
		return metadata
	}
	databaseHasMetadata := databaseViewDefinition.Valid && databaseViewDefinition.Bool
	schemaHasMetadata := schemaViewDefinition.Valid && schemaViewDefinition.Bool
	if !databaseHasMetadata && !schemaHasMetadata {
		if databaseViewDefinition.Valid && schemaViewDefinition.Valid {
			metadata.metadataStatus = sqlServerRecordingMetadataForbidden
			metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, "metadata_forbidden")
		} else {
			metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, "schema_visibility_unconfirmed")
		}
		return metadata
	}

	var createTable, alterSchema, selectSchema, insertSchema, deleteSchema, alterAnyDataspace sql.NullBool
	if err := db.QueryRowContext(queryCtx, sqlServerRecordingPermissionQuery, expectedSchema).Scan(
		&createTable,
		&alterSchema,
		&selectSchema,
		&insertSchema,
		&deleteSchema,
		&alterAnyDataspace,
	); err != nil {
		setSQLServerRecordingMetadataFailure(&metadata, err)
		return metadata
	}
	permissionValues := []sql.NullBool{
		createTable, alterSchema, selectSchema, insertSchema, deleteSchema, alterAnyDataspace,
	}
	for _, value := range permissionValues {
		if !value.Valid {
			metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, "operation_permissions_unconfirmed")
			return metadata
		}
	}
	metadata.permissions = SQLServerRecordingPermissions{
		Verified:               true,
		DatabaseViewDefinition: databaseViewDefinition.Bool,
		SchemaViewDefinition:   schemaViewDefinition.Bool,
		CreateTable:            createTable.Bool,
		AlterSchema:            alterSchema.Bool,
		SelectSchema:           selectSchema.Bool,
		InsertSchema:           insertSchema.Bool,
		DeleteSchema:           deleteSchema.Bool,
		AlterAnyDataspace:      alterAnyDataspace.Bool,
	}
	metadata.actualSchema = visibleSchema
	metadata.metadataStatus = sqlServerRecordingMetadataVerified
	if !sqlServerRecordingPermissionsAllowManagedLong(metadata.permissions) {
		metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, "operation_permissions_denied")
	}
	if !metadata.permissions.DatabaseViewDefinition || !metadata.permissions.AlterAnyDataspace {
		metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, "native_partition_permissions_denied")
	}
	return metadata
}

func setSQLServerRecordingMetadataFailure(metadata *sqlServerRecordingMetadata, err error) {
	if isSQLServerPermissionDeniedError(err) {
		metadata.metadataStatus = sqlServerRecordingMetadataForbidden
		metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, "metadata_forbidden")
		return
	}
	metadata.metadataStatus = sqlServerRecordingMetadataUnconfirmed
	metadata.reasons = appendSQLServerRecordingReason(metadata.reasons, "metadata_unconfirmed")
}

func isSQLServerPermissionDeniedError(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return isPermissionDeniedError(err) || strings.Contains(message, "permission was denied") ||
		strings.Contains(message, "is denied") || strings.Contains(message, "not authorized")
}

func projectSQLServerRecordingCapability(
	connector *schema.DatabaseConnector,
	metadata sqlServerRecordingMetadata,
	connected bool,
) *SQLServerRecordingCapability {
	capability := &SQLServerRecordingCapability{
		Database:           metadata.database,
		Schema:             metadata.schema,
		ServerVersion:      metadata.serverVersion,
		ServerEdition:      metadata.serverEdition,
		EngineEdition:      metadata.engineEdition,
		CompatibilityLevel: metadata.compatibilityLevel,
		Connected:          connected,
		MetadataStatus:     metadata.metadataStatus,
		Permissions:        metadata.permissions,
		Reasons:            []string{},
	}
	if connector != nil {
		capability.ConnectorID = connector.ID
		capability.ConnectorRevision = connector.IdentityRevision
	}
	for _, reason := range metadata.reasons {
		capability.Reasons = appendSQLServerRecordingReason(capability.Reasons, reason)
	}
	versionOK := sqlServerRecordingVersionSupported(metadata.serverVersion)
	if !versionOK {
		capability.Reasons = appendSQLServerRecordingReason(capability.Reasons, "server_version_unverified")
	}
	editionOK := strings.TrimSpace(metadata.serverEdition) != "" && metadata.engineEditionValid &&
		(metadata.engineEdition == 2 || metadata.engineEdition == 3 || metadata.engineEdition == 4)
	if !editionOK {
		capability.Reasons = appendSQLServerRecordingReason(capability.Reasons, "server_edition_unverified")
	}
	compatOK := metadata.compatibilityValid && sqlServerRecordingCompatibilitySupported(metadata.serverVersion, metadata.compatibilityLevel)
	if !compatOK {
		capability.Reasons = appendSQLServerRecordingReason(capability.Reasons, "compatibility_unverified")
	}

	scopeOK := connected && metadata.metadataStatus == sqlServerRecordingMetadataVerified &&
		metadata.database != "" && metadata.database == metadata.actualDatabase &&
		metadata.schema != "" && metadata.schema == metadata.actualSchema

	managedLongAllowed := scopeOK && versionOK && editionOK && compatOK &&
		sqlServerRecordingPermissionsAllowManagedLong(metadata.permissions)

	if managedLongAllowed {
		capability.SupportsManagedLong = true
		capability.SupportsMissingRepair = true
		capability.SupportsGroupTestWrite = true
		if metadata.permissions.DatabaseViewDefinition && metadata.permissions.AlterAnyDataspace {
			capability.SupportsNativePartition = true
		} else {
			capability.Reasons = appendSQLServerRecordingReason(capability.Reasons, "native_partition_permissions_denied")
		}
	} else {
		capability.Reasons = appendSQLServerRecordingReason(capability.Reasons,
			"adapter_unverified",
			"managed_long_unverified",
			"native_partition_unverified",
			"group_test_write_unverified",
			"missing_repair_unverified",
		)
	}
	capability.CapabilityRevision = sqlServerRecordingCapabilityRevision(connector, metadata, connected)
	return capability
}

func sqlServerRecordingVersionSupported(version string) bool {
	return strings.HasPrefix(version, "16.") || strings.HasPrefix(version, "17.")
}

func sqlServerRecordingCompatibilitySupported(version string, compatibilityLevel int64) bool {
	switch {
	case strings.HasPrefix(version, "16."):
		return compatibilityLevel == 160
	case strings.HasPrefix(version, "17."):
		return compatibilityLevel == 170
	default:
		return false
	}
}

func sqlServerRecordingPermissionsAllowManagedLong(permissions SQLServerRecordingPermissions) bool {
	return permissions.Verified && permissions.CreateTable && permissions.AlterSchema &&
		permissions.SelectSchema && permissions.InsertSchema && permissions.DeleteSchema
}

func sqlServerRecordingCapabilityRevision(
	connector *schema.DatabaseConnector,
	metadata sqlServerRecordingMetadata,
	connected bool,
) string {
	connectorID, connectorRevision := "", ""
	if connector != nil {
		connectorID = connector.ID
		connectorRevision = connector.IdentityRevision
	}
	permissions := metadata.permissions
	canonical := fmt.Appendf(nil,
		"%q|%q|%q|%q|%q|%q|%d|%t|%d|%t|%q|%q|%q|%q|%t|%t|%t|%t|%t|%t|%t|%t|%t|%t",
		sqlServerRecordingCapabilityProjectionVersion,
		connectorID,
		connectorRevision,
		metadata.serverVersion,
		metadata.serverEdition,
		metadata.actualDatabase,
		metadata.engineEdition,
		metadata.engineEditionValid,
		metadata.compatibilityLevel,
		metadata.compatibilityValid,
		metadata.database,
		metadata.schema,
		metadata.actualSchema,
		metadata.metadataStatus,
		permissions.Verified,
		permissions.DatabaseViewDefinition,
		permissions.SchemaViewDefinition,
		permissions.CreateTable,
		permissions.AlterSchema,
		permissions.SelectSchema,
		permissions.InsertSchema,
		permissions.DeleteSchema,
		permissions.AlterAnyDataspace,
		connected,
	)
	digest := sha256.Sum256(canonical)
	return "sqlservercap-v1-" + hex.EncodeToString(digest[:])
}

func appendSQLServerRecordingReason(reasons []string, additions ...string) []string {
	for _, addition := range additions {
		if addition != "" && !slices.Contains(reasons, addition) {
			reasons = append(reasons, addition)
		}
	}
	return reasons
}
