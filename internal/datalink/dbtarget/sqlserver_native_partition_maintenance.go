package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	"go-gateway/internal/datalink/schema"
)

const (
	sqlServerNativePartitionMaintenanceTimeout = 30 * time.Second
	sqlServerNativePartitionReconcileTimeout   = 5 * time.Second
	sqlServerNativePartitionMaxSplits          = 5
)

var (
	// ErrSQLServerNativePartitionMaintenanceBlocked marks a maintenance attempt that could not be safely completed.
	ErrSQLServerNativePartitionMaintenanceBlocked = errors.New("sqlserver native partition maintenance blocked")
	// ErrSQLServerNativePartitionOutcomeUnknown marks catalog state outside the exact before and after proofs.
	ErrSQLServerNativePartitionOutcomeUnknown = errors.New("sqlserver native partition outcome unknown")
)

// SplitOpenSQLServerNativeMonthly appends at most five verified monthly
// boundaries to an exclusively owned SQL Server partition function. It derives
// expected after-proof metadata from the verified before proof and the requested
// after layout; actual before and after catalogs are checked inside the DDL transaction.
func SplitOpenSQLServerNativeMonthly(
	ctx context.Context,
	db *sql.DB,
	capability *SQLServerRecordingCapability,
	beforeGroup *SQLServerManagedInspection,
	after *SQLServerManagedPartitionEvidence,
) error {
	blocked := func(err error) error {
		return fmt.Errorf("%w: %w", ErrSQLServerNativePartitionMaintenanceBlocked, err)
	}
	if ctx == nil || db == nil {
		return blocked(ErrSQLServerManagedTableUnconfirmed)
	}
	if err := validateSQLServerRecordingScope(capability); err != nil {
		return blocked(err)
	}
	if !capability.SupportsManagedLong || !capability.SupportsNativePartition {
		return blocked(ErrSchemaExecutionUnsupported)
	}
	if !sqlServerNativePartitionPermissionsAllowMaintenance(capability.Permissions) {
		return blocked(ErrSchemaExecutionPermissionDenied)
	}
	if !sqlServerNativePartitionBeforeProofMatchesCapability(capability, beforeGroup) {
		return blocked(ErrSQLServerManagedScopeChanged)
	}
	if !sqlServerManagedInspectionReady(beforeGroup) {
		return blocked(ErrSQLServerManagedTableUnconfirmed)
	}
	if _, ok := sqlServerManagedGroupShapeMatches(beforeGroup); !ok {
		return blocked(ErrSQLServerManagedPhysicalMismatch)
	}
	expectedAfter := sqlServerNativePartitionExpectedAfter(beforeGroup.Inspection.Table, beforeGroup.Evidence.NativeMonthly, after)
	if expectedAfter == nil {
		return blocked(ErrSQLServerManagedPhysicalMismatch)
	}

	afterGroup := sqlServerManagedInspectionWithNativeProof(beforeGroup, expectedAfter)
	if _, ok := sqlServerManagedGroupShapeMatches(afterGroup); !ok {
		return blocked(ErrSQLServerManagedPhysicalMismatch)
	}
	statements := sqlServerNativePartitionSplitStatements(beforeGroup.Evidence.NativeMonthly, expectedAfter)
	queryCtx, cancel := context.WithTimeout(ctx, sqlServerNativePartitionMaintenanceTimeout)
	defer cancel()

	preflightPassed, executeDDL := false, false
	before := func(txCtx context.Context, tx *sql.Tx) (bool, error) {
		if err := verifySQLServerNativePartitionMaintenanceScope(txCtx, tx, capability, beforeGroup); err != nil {
			return false, err
		}
		facts, err := inspectSQLServerManagedCatalog(txCtx, tx, beforeGroup.Inspection.Schema, beforeGroup.Inspection.Table)
		if err != nil {
			return false, err
		}
		if sqlServerManagedCatalogMatches(afterGroup, facts) {
			preflightPassed = true
			return false, nil
		}
		if !sqlServerManagedCatalogMatches(beforeGroup, facts) {
			return false, ErrSQLServerManagedPhysicalMismatch
		}
		preflightPassed = true
		executeDDL = len(statements) > 0
		return executeDDL, nil
	}
	finish := func(txCtx context.Context, tx *sql.Tx) error {
		facts, err := inspectSQLServerManagedCatalog(txCtx, tx, afterGroup.Inspection.Schema, afterGroup.Inspection.Table)
		if err != nil {
			return err
		}
		if !sqlServerManagedCatalogMatches(afterGroup, facts) {
			return ErrSQLServerManagedPhysicalMismatch
		}
		return nil
	}

	_, err := executeSchemaBatchWithHooks(
		queryCtx, db, schema.DatabaseConnectorKindSQLServer, capability.Schema, statements, before, finish,
	)
	if err == nil {
		return nil
	}
	if !preflightPassed {
		if errors.Is(err, ErrSQLServerManagedPhysicalMismatch) {
			return errors.Join(blocked(err), ErrSQLServerNativePartitionOutcomeUnknown)
		}
		return blocked(err)
	}
	if !executeDDL && errors.Is(err, ErrSchemaExecutionPermissionDenied) {
		return blocked(err)
	}
	return resolveSQLServerNativePartitionOutcome(ctx, db, capability, beforeGroup, afterGroup, err, blocked)
}

func sqlServerNativePartitionPermissionsAllowMaintenance(permissions SQLServerRecordingPermissions) bool {
	return permissions.Verified && permissions.DatabaseViewDefinition && permissions.AlterAnyDataspace
}

func sqlServerNativePartitionBeforeProofMatchesCapability(
	capability *SQLServerRecordingCapability,
	before *SQLServerManagedInspection,
) bool {
	if capability == nil || before == nil || before.Inspection == nil || before.Evidence == nil {
		return false
	}
	evidence := before.Evidence
	return before.Inspection.Schema == capability.Schema && evidence.Schema == capability.Schema &&
		evidence.Database == capability.Database && evidence.CapabilityRevision == capability.CapabilityRevision &&
		evidence.ServerVersion == capability.ServerVersion && evidence.ServerEdition == capability.ServerEdition &&
		evidence.EngineEdition == capability.EngineEdition && evidence.CompatibilityLevel == capability.CompatibilityLevel &&
		evidence.MetadataVerified && evidence.PartitionInspected && evidence.NativeMonthly != nil
}

func verifySQLServerNativePartitionMaintenanceScope(
	ctx context.Context,
	tx *sql.Tx,
	capability *SQLServerRecordingCapability,
	before *SQLServerManagedInspection,
) error {
	metadata := inspectSQLServerRecordingMetadataWithQuerier(ctx, tx, capability.Database, capability.Schema)
	switch metadata.metadataStatus {
	case sqlServerRecordingMetadataForbidden:
		return errors.Join(ErrSchemaExecutionPermissionDenied, ErrSQLServerManagedTableForbidden)
	case sqlServerRecordingMetadataVerified:
	default:
		if slices.Contains(metadata.reasons, "database_scope_mismatch") || slices.Contains(metadata.reasons, "schema_scope_mismatch") {
			return ErrSQLServerManagedScopeChanged
		}
		return ErrSQLServerManagedTableUnconfirmed
	}
	if !sqlServerNativePartitionPermissionsAllowMaintenance(metadata.permissions) {
		return ErrSchemaExecutionPermissionDenied
	}
	connector := &schema.DatabaseConnector{
		ID: capability.ConnectorID, Kind: schema.DatabaseConnectorKindSQLServer,
		IdentityRevision: capability.ConnectorRevision,
	}
	actual := projectSQLServerRecordingCapability(connector, metadata, true)
	if metadata.actualDatabase != capability.Database || metadata.actualSchema != capability.Schema ||
		!sqlServerRecordingCapabilityScopeMatches(capability, actual) ||
		before.Evidence.CapabilityRevision != actual.CapabilityRevision {
		return ErrSQLServerManagedScopeChanged
	}
	return nil
}

func sqlServerNativePartitionExpectedAfter(
	tableName string,
	before, desired *SQLServerManagedPartitionEvidence,
) *SQLServerManagedPartitionEvidence {
	if !sqlServerManagedNativeMonthlyEvidenceMatchesForTable(tableName, before) ||
		!sqlServerManagedNativeDefinitionPropertiesEqual(before, desired) ||
		!sqlServerManagedNativeMonthlyDefinitionMatchesForTable(tableName, desired) {
		return nil
	}
	beforeFinite := before.Partitions[:len(before.Partitions)-1]
	desiredFinite := desired.Partitions[:len(desired.Partitions)-1]
	newBoundaryCount := len(desiredFinite) - len(beforeFinite)
	if newBoundaryCount < 0 || newBoundaryCount > sqlServerNativePartitionMaxSplits ||
		len(desiredFinite) < len(beforeFinite) || !slices.Equal(desiredFinite[:len(beforeFinite)], beforeFinite) {
		return nil
	}

	expected := cloneSQLServerManagedPartitionEvidence(desired)
	expected.FunctionExists = before.FunctionExists
	expected.SchemeExists = before.SchemeExists
	expected.Exclusive = before.Exclusive
	expected.IndexesAligned = before.IndexesAligned
	return expected
}

func sqlServerManagedNativeDefinitionPropertiesEqual(before, desired *SQLServerManagedPartitionEvidence) bool {
	return before != nil && desired != nil && before.Method == desired.Method && before.Column == desired.Column &&
		before.FunctionName == desired.FunctionName && before.SchemeName == desired.SchemeName && before.Filegroup == desired.Filegroup
}

func sqlServerManagedInspectionWithNativeProof(
	before *SQLServerManagedInspection,
	after *SQLServerManagedPartitionEvidence,
) *SQLServerManagedInspection {
	result := *before
	evidence := *before.Evidence
	evidence.NativeMonthly = cloneSQLServerManagedPartitionEvidence(after)
	evidence.PartitionInspected = true
	result.Evidence = &evidence
	return &result
}

func sqlServerNativePartitionSplitStatements(
	before, after *SQLServerManagedPartitionEvidence,
) []string {
	beforeFiniteCount := len(before.Partitions) - 1
	newBoundaries := after.Partitions[beforeFiniteCount : len(after.Partitions)-1]
	statements := make([]string, 0, len(newBoundaries)*2)
	for _, boundary := range newBoundaries {
		statements = append(statements,
			"ALTER PARTITION SCHEME "+quoteIdentifier(schema.DatabaseConnectorKindSQLServer, after.SchemeName)+
				" NEXT USED "+quoteIdentifier(schema.DatabaseConnectorKindSQLServer, sqlServerNativePrimaryFilegroup),
			"ALTER PARTITION FUNCTION "+quoteIdentifier(schema.DatabaseConnectorKindSQLServer, after.FunctionName)+
				"() SPLIT RANGE ('"+boundary.LessThan+"')",
		)
	}
	return statements
}

func resolveSQLServerNativePartitionOutcome(
	ctx context.Context,
	db *sql.DB,
	capability *SQLServerRecordingCapability,
	before, after *SQLServerManagedInspection,
	attemptErr error,
	blocked func(error) error,
) error {
	reconcileCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sqlServerNativePartitionReconcileTimeout)
	defer cancel()
	actual, inspectErr := InspectOpenSQLServerManagedTable(reconcileCtx, db, capability, before.Inspection.Table)
	if inspectErr == nil {
		facts := sqlServerManagedFactsFromInspection(actual)
		if sqlServerManagedCatalogMatches(after, facts) {
			return nil
		}
		if sqlServerManagedCatalogMatches(before, facts) {
			return blocked(attemptErr)
		}
	}
	return errors.Join(
		blocked(attemptErr),
		fmt.Errorf("%w: live catalog matches neither exact before nor exact after proof", ErrSQLServerNativePartitionOutcomeUnknown),
		inspectErr,
	)
}

func sqlServerManagedFactsFromInspection(inspection *SQLServerManagedInspection) *sqlServerManagedCatalogFacts {
	if inspection == nil || inspection.Inspection == nil || inspection.Evidence == nil {
		return nil
	}
	evidence := inspection.Evidence
	return &sqlServerManagedCatalogFacts{
		object: sqlServerManagedObjectFacts{
			objectID: evidence.ObjectID, schema: evidence.Schema, name: inspection.Inspection.Table,
			tableType: evidence.TableType, userTable: evidence.TableType == sqlServerManagedUserTableType,
		},
		columns: inspection.Inspection.Columns, columnCollations: evidence.ColumnCollations,
		columnDefaults: evidence.ColumnDefaults, primaryKey: evidence.PrimaryKey, indexes: inspection.Indexes,
		plain: evidence.Plain, provenanceJSONChecked: evidence.ProvenanceJSONChecked,
		repairOperationID: evidence.RepairOperationID, nativeMonthly: evidence.NativeMonthly,
		partitionInspected: evidence.PartitionInspected,
	}
}
