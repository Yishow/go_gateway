package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
)

func sqlServerSQLErrorNumber(err error) (int32, bool) {
	var numbered interface{ SQLErrorNumber() int32 }
	if !errors.As(err, &numbered) {
		return 0, false
	}
	return numbered.SQLErrorNumber(), true
}

func isSQLServerSchemaShapeError(err error) bool {
	number, ok := sqlServerSQLErrorNumber(err)
	return ok && (number == 207 || number == 208)
}

func isSQLServerMissingRelationError(err error) bool {
	number, ok := sqlServerSQLErrorNumber(err)
	return ok && number == 208
}

func isSQLServerPermissionDeniedNumber(err error) bool {
	number, ok := sqlServerSQLErrorNumber(err)
	return ok && number == 229
}

const sqlServerRepairMarkerAddStatement = `EXEC sys.sp_addextendedproperty
	@name = @p1, @value = @p2,
	@level0type = N'SCHEMA', @level0name = @p3,
	@level1type = N'TABLE', @level1name = @p4`

// SQLServerManagedRepairRequest binds the reconstructed group table and repair
// operation to the confirmed schema statements.
type SQLServerManagedRepairRequest struct {
	TableName   string
	OperationID string
	Statements  []string
}

func validateSQLServerRecordingScope(capability *SQLServerRecordingCapability) error {
	if capability == nil || !capability.Connected || capability.MetadataStatus != sqlServerRecordingMetadataVerified ||
		!capability.Permissions.Verified || strings.TrimSpace(capability.ConnectorID) == "" ||
		strings.TrimSpace(capability.ConnectorRevision) == "" || strings.TrimSpace(capability.Database) == "" ||
		strings.TrimSpace(capability.Schema) == "" || strings.TrimSpace(capability.ServerVersion) == "" ||
		strings.TrimSpace(capability.ServerEdition) == "" || capability.EngineEdition <= 0 ||
		capability.CompatibilityLevel <= 0 || strings.TrimSpace(capability.CapabilityRevision) == "" {
		return ErrSQLServerManagedTableUnconfirmed
	}
	return nil
}

func sqlServerRecordingCapabilityScopeMatches(expected, actual *SQLServerRecordingCapability) bool {
	return expected != nil && actual != nil && expected.ConnectorID == actual.ConnectorID &&
		expected.ConnectorRevision == actual.ConnectorRevision && expected.Database == actual.Database &&
		expected.Schema == actual.Schema && expected.ServerVersion == actual.ServerVersion &&
		expected.ServerEdition == actual.ServerEdition && expected.EngineEdition == actual.EngineEdition &&
		expected.CompatibilityLevel == actual.CompatibilityLevel && expected.Connected == actual.Connected &&
		expected.MetadataStatus == actual.MetadataStatus && expected.Permissions == actual.Permissions &&
		expected.CapabilityRevision == actual.CapabilityRevision
}

// InspectOpenSQLServerManagedTable reads the supplied SQL Server destination
// without opening another connection. The capability binds the saved identity
// and scope; live metadata and catalog facts are rechecked before returning.
func InspectOpenSQLServerManagedTable(
	ctx context.Context,
	db *sql.DB,
	capability *SQLServerRecordingCapability,
	tableName string,
) (*SQLServerManagedInspection, error) {
	if db == nil || validateSQLServerRecordingScope(capability) != nil || strings.TrimSpace(tableName) == "" {
		return nil, ErrSQLServerManagedTableUnconfirmed
	}
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	connector := &schema.DatabaseConnector{
		ID: capability.ConnectorID, Kind: schema.DatabaseConnectorKindSQLServer,
		IdentityRevision: capability.ConnectorRevision,
	}
	metadata := inspectSQLServerRecordingMetadata(queryCtx, db, capability.Database, capability.Schema)
	projected := projectSQLServerRecordingCapability(connector, metadata, true)
	return inspectSQLServerManagedTable(
		queryCtx, db, capability.Schema, tableName, capability.CapabilityRevision, metadata, projected, capability,
	)
}

// RunSQLServerManagedRepairBatch runs confirmed, schema-qualified repair
// statements and the operation marker in one SQL Server transaction. Callers
// must authorize the exact group-table and index statements from the ledger.
func RunSQLServerManagedRepairBatch(
	ctx context.Context,
	db *sql.DB,
	capability *SQLServerRecordingCapability,
	request SQLServerManagedRepairRequest,
) (*SchemaStatementExecution, error) {
	reject := func(err error) (*SchemaStatementExecution, error) {
		return &SchemaStatementExecution{RolledBack: true}, err
	}
	if strings.TrimSpace(request.TableName) == "" || strings.TrimSpace(request.OperationID) == "" ||
		strings.EqualFold(strings.TrimSpace(request.TableName), EffectReceiptTable) || len(request.Statements) == 0 {
		return reject(ErrSchemaExecutionUnsupported)
	}
	if db == nil {
		return reject(ErrSQLServerManagedTableUnconfirmed)
	}
	if err := validateSQLServerRecordingScope(capability); err != nil {
		return reject(err)
	}
	if !capability.SupportsManagedLong || !capability.SupportsMissingRepair {
		return reject(ErrSchemaExecutionUnsupported)
	}

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	metadata := inspectSQLServerRecordingMetadata(queryCtx, db, capability.Database, capability.Schema)
	connector := &schema.DatabaseConnector{
		ID: capability.ConnectorID, Kind: schema.DatabaseConnectorKindSQLServer,
		IdentityRevision: capability.ConnectorRevision,
	}
	projected := projectSQLServerRecordingCapability(connector, metadata, true)
	if metadata.metadataStatus == sqlServerRecordingMetadataForbidden {
		return reject(fmt.Errorf("%w: %w", ErrSchemaExecutionPermissionDenied, ErrSQLServerManagedTableForbidden))
	}
	if metadata.metadataStatus == sqlServerRecordingMetadataMissing {
		return reject(ErrSQLServerManagedTableUnconfirmed)
	}
	if metadata.metadataStatus != sqlServerRecordingMetadataVerified ||
		metadata.actualDatabase != capability.Database || metadata.actualSchema != capability.Schema {
		if slices.Contains(metadata.reasons, "database_scope_mismatch") || slices.Contains(metadata.reasons, "schema_scope_mismatch") {
			return reject(ErrSQLServerManagedScopeChanged)
		}
		return reject(ErrSQLServerManagedTableUnconfirmed)
	}
	if !metadata.permissions.Verified {
		return reject(ErrSQLServerManagedTableUnconfirmed)
	}
	if !sqlServerRecordingPermissionsAllowManagedLong(metadata.permissions) {
		return reject(ErrSchemaExecutionPermissionDenied)
	}
	if !sqlServerRecordingCapabilityScopeMatches(capability, projected) {
		return reject(ErrSQLServerManagedScopeChanged)
	}

	finish := func(finishCtx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(
			finishCtx,
			sqlServerRepairMarkerAddStatement,
			sqlServerManagedRepairOperationPropertyName, request.OperationID, capability.Schema, request.TableName,
		)
		return err
	}
	execution, err := executeSQLServerManagedRepairBatch(ctx, db, capability, request, finish)
	if isSQLServerPermissionDeniedNumber(err) && !errors.Is(err, ErrSchemaExecutionPermissionDenied) {
		return execution, fmt.Errorf("%w: %w", ErrSchemaExecutionPermissionDenied, err)
	}
	return execution, err
}
