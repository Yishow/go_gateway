package dbtarget

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
)

// ExecuteSQLServerManagedStatementsAtRevision executes only the server-derived
// confirmed SQL Server schema batch for the exact saved connector revision and
// capability. It is separate from the closed generic ExecuteSchemaStatements
// path and does not expose arbitrary SQL through HTTP.
func (i *ManagedTableInspector) ExecuteSQLServerManagedStatementsAtRevision(
	ctx context.Context,
	connectorID, expectedRevision, database, schemaName, expectedCapabilityRevision string,
	native bool, statements []string,
) (*SchemaStatementExecution, error) {
	reject := func(err error) (*SchemaStatementExecution, error) {
		return &SchemaStatementExecution{RolledBack: true}, err
	}
	if i == nil || i.service == nil || i.service.repo == nil {
		return reject(managedDestinationInvalid("saved connector service is unavailable"))
	}
	database, schemaName, expectedCapabilityRevision = strings.TrimSpace(database), strings.TrimSpace(schemaName), strings.TrimSpace(expectedCapabilityRevision)
	if database == "" || schemaName == "" || expectedCapabilityRevision == "" || len(statements) == 0 {
		return reject(ErrSQLServerManagedTableUnconfirmed)
	}
	connector, err := i.service.ResolveSavedTarget(ctx, strings.TrimSpace(connectorID), strings.TrimSpace(expectedRevision))
	if err != nil {
		return reject(err)
	}
	if connector == nil || connector.Kind != schema.DatabaseConnectorKindSQLServer || !connector.Enabled {
		return reject(managedDestinationInvalid("saved sqlserver connector is unavailable"))
	}
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return reject(managedDestinationInvalid("saved sqlserver connector configuration is invalid"))
	}
	settings, err := sqlServerConnectionSettingsFor(config)
	if err != nil || settings.Database != database || settings.Schema != schemaName {
		return reject(ErrSQLServerManagedScopeChanged)
	}
	session, err := openSQLServerManagedSessionFunc(connector.Kind, config)
	if err != nil || session == nil || session.DB() == nil {
		if session != nil {
			_ = session.Close()
		}
		return reject(ErrSQLServerManagedTableUnconfirmed)
	}
	defer session.Close()
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	metadata := inspectSQLServerRecordingMetadata(queryCtx, session.DB(), database, schemaName)
	projected := projectSQLServerRecordingCapability(connector, metadata, true)
	if metadata.metadataStatus != sqlServerRecordingMetadataVerified || metadata.actualDatabase != database || metadata.actualSchema != schemaName ||
		projected.CapabilityRevision != expectedCapabilityRevision || !projected.SupportsManagedLong || (native && !projected.SupportsNativePartition) ||
		!sqlServerRecordingPermissionsAllowManagedLong(projected.Permissions) {
		return reject(ErrSQLServerManagedScopeChanged)
	}
	execution, err := executeSchemaBatch(ctx, session.DB(), schema.DatabaseConnectorKindSQLServer, schemaName, statements)
	if isSQLServerPermissionDeniedNumber(err) && err != nil {
		return execution, fmt.Errorf("%w: %w", ErrSchemaExecutionPermissionDenied, err)
	}
	return execution, err
}
