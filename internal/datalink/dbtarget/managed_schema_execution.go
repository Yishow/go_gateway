package dbtarget

import (
	"context"

	"go-gateway/internal/datalink/schema"
)

// ExecuteSchemaStatements binds managed DDL to the confirmed connector revision.
func (i *ManagedTableInspector) ExecuteSchemaStatements(ctx context.Context, connectorID, expectedRevision, schemaName string, statements []string) (*SchemaStatementExecution, error) {
	connector, err := i.ValidateDestination(ctx, connectorID, expectedRevision)
	if err != nil {
		return &SchemaStatementExecution{RolledBack: true}, err
	}
	if connector.Kind == schema.DatabaseConnectorKindSQLite && schemaName != managedSQLiteSchemaName {
		return &SchemaStatementExecution{RolledBack: true}, managedDestinationInvalid("managed SQLite schema must be main")
	}
	return executeSchemaStatementsOnConnector(ctx, connector, schemaName, statements)
}
