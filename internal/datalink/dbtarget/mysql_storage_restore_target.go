package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// InspectOpenMySQLManagedTable inspects one MySQL managed table on a delivery
// connection that is already open. It never opens a manager or provisions a
// database.
func InspectOpenMySQLManagedTable(ctx context.Context, db *sql.DB, expected *MySQLRecordingCapability, database, table string) (*MySQLManagedInspection, error) {
	if db == nil || expected == nil {
		return nil, ErrMySQLManagedServerMismatch
	}
	actual, err := reprojectOpenMySQLCapability(ctx, db, expected, database)
	if err != nil {
		if actual != nil {
			return openMySQLBindingInspection(actual, database, table, err), err
		}
		return nil, err
	}
	facts, inspectErr := readMySQLManagedTableFacts(ctx, db, actual.ServerFamily, database, table)
	if inspectErr != nil {
		result := openMySQLBindingInspection(actual, database, table, inspectErr)
		if errors.Is(inspectErr, ErrMySQLManagedTableForbidden) {
			return result, fmt.Errorf("inspect mysql managed table: %w", inspectErr)
		}
		if isPermissionDeniedError(inspectErr) {
			return result, fmt.Errorf("%w: inspect mysql managed table: %w", ErrMySQLManagedTableForbidden, inspectErr)
		}
		return result, fmt.Errorf("%w: %w", ErrMySQLManagedTableUnconfirmed, inspectErr)
	}
	server := mysqlManagedServerFacts{
		family: actual.ServerFamily, version: actual.ServerVersion, capabilityRevision: actual.CapabilityRevision,
		database: actual.Database, metadataStatus: mysqlRecordingMetadataVerified, transactionalEngine: true,
	}
	return buildMySQLManagedInspection(database, table, server, facts)
}

// ExecuteOpenMySQLManagedStatement executes one confirmed MySQL DDL statement
// on an already opened delivery connection. It never starts a row transaction
// and never provisions a database.
func ExecuteOpenMySQLManagedStatement(ctx context.Context, db *sql.DB, expected *MySQLRecordingCapability, database, statement string) error {
	if err := validateMySQLManagedStatement(statement); err != nil {
		return err
	}
	if db == nil || expected == nil {
		return ErrMySQLManagedServerMismatch
	}
	if _, err := reprojectOpenMySQLCapability(ctx, db, expected, database); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, statement); err != nil {
		return fmt.Errorf("execute mysql managed statement: %w", err)
	}
	return nil
}

func reprojectOpenMySQLCapability(ctx context.Context, db *sql.DB, expected *MySQLRecordingCapability, database string) (*MySQLRecordingCapability, error) {
	database = strings.TrimSpace(database)
	if strings.TrimSpace(expected.ConnectorID) == "" || strings.TrimSpace(expected.ConnectorRevision) == "" || database == "" {
		return nil, ErrMySQLManagedServerMismatch
	}
	if strings.TrimSpace(expected.Database) != database || expected.ServerFamily == "" || expected.ServerVersion == "" ||
		expected.CapabilityRevision == "" {
		return nil, ErrMySQLManagedServerMismatch
	}
	connector := &schema.DatabaseConnector{
		ID: expected.ConnectorID, IdentityRevision: expected.ConnectorRevision, Kind: schema.DatabaseConnectorKindMySQL,
	}
	metadata := inspectMySQLRecordingMetadata(ctx, db, database)
	actual := projectMySQLRecordingCapability(connector, metadata, true)
	switch metadata.metadataStatus {
	case mysqlRecordingMetadataVerified:
	case mysqlRecordingMetadataForbidden:
		return actual, ErrMySQLManagedTableForbidden
	default:
		return actual, ErrMySQLManagedTableUnconfirmed
	}
	if !actual.TransactionalEngine {
		return actual, fmt.Errorf("%w: InnoDB is unverified", ErrMySQLManagedServerMismatch)
	}
	if actual.ServerFamily != expected.ServerFamily || actual.ServerVersion != expected.ServerVersion ||
		actual.Database != database {
		return actual, ErrMySQLManagedServerMismatch
	}
	if actual.CapabilityRevision != expected.CapabilityRevision {
		return actual, ErrMySQLManagedCapabilityChanged
	}
	return actual, nil
}

func openMySQLBindingInspection(actual *MySQLRecordingCapability, database, table string, err error) *MySQLManagedInspection {
	if isPermissionDeniedError(err) && !errors.Is(err, ErrMySQLManagedTableForbidden) {
		err = fmt.Errorf("%w: %w", ErrMySQLManagedTableForbidden, err)
	}
	server := mysqlManagedServerFacts{
		family: actual.ServerFamily, version: actual.ServerVersion, capabilityRevision: actual.CapabilityRevision,
		database: actual.Database, metadataStatus: mysqlRecordingMetadataVerified, transactionalEngine: actual.TransactionalEngine,
	}
	return &MySQLManagedInspection{
		Inspection: mysqlManagedFailedInspection(database, table, err),
		Evidence:   mysqlManagedBindingEvidence(server, database),
	}
}
