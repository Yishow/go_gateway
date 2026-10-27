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

const sqlServerNativeRepairPartitionStatementCount = 4

func executeSQLServerManagedRepairBatch(
	ctx context.Context,
	db *sql.DB,
	capability *SQLServerRecordingCapability,
	request SQLServerManagedRepairRequest,
	finish func(context.Context, *sql.Tx) error,
) (*SchemaStatementExecution, error) {
	_, native, err := parseSQLServerNativeRepairDefinition(capability, request)
	if err != nil {
		return &SchemaStatementExecution{RolledBack: true}, err
	}
	if !native {
		return executeSchemaBatchWithFinish(ctx, db, schema.DatabaseConnectorKindSQLServer, capability.Schema, request.Statements, finish)
	}
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return &SchemaStatementExecution{RolledBack: true}, fmt.Errorf("begin sqlserver repair transaction: %w", err)
	}
	statements, err := sqlServerManagedRepairStatementsInTransaction(ctx, tx, capability, request)
	if err != nil {
		return rollbackSchemaBatch(tx, fmt.Errorf("recheck sqlserver repair catalog: %w", err))
	}
	for _, statement := range statements {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return rollbackSchemaBatch(tx, fmt.Errorf("execute sqlserver repair statement: %w", err))
		}
	}
	if finish != nil {
		if err := finish(ctx, tx); err != nil {
			return rollbackSchemaBatch(tx, fmt.Errorf("finish sqlserver repair transaction: %w", err))
		}
	}
	if err := tx.Commit(); err != nil {
		return &SchemaStatementExecution{}, fmt.Errorf("commit sqlserver repair transaction: %w", err)
	}
	return &SchemaStatementExecution{Committed: len(statements)}, nil
}

func sqlServerManagedRepairStatementsInTransaction(
	ctx context.Context,
	tx *sql.Tx,
	capability *SQLServerRecordingCapability,
	request SQLServerManagedRepairRequest,
) ([]string, error) {
	desired, native, err := parseSQLServerNativeRepairDefinition(capability, request)
	if err != nil || !native {
		return slices.Clone(request.Statements), err
	}
	metadata, err := verifySQLServerManagedRepairMetadataInTransaction(ctx, tx, capability)
	if err != nil {
		return nil, err
	}
	if _, err := readSQLServerManagedObject(ctx, tx, capability.Schema, request.TableName); !errors.Is(err, sql.ErrNoRows) {
		if err != nil {
			return nil, err
		}
		return nil, ErrSQLServerManagedScopeChanged
	}
	actual, inspected, err := readSQLServerManagedPartitions(ctx, tx, capability.Schema, request.TableName, 0)
	if err != nil {
		return nil, err
	}
	if !inspected {
		return nil, ErrSQLServerManagedTableUnconfirmed
	}
	if actual == nil {
		if !metadata.permissions.AlterAnyDataspace {
			return nil, ErrSchemaExecutionPermissionDenied
		}
		return slices.Clone(request.Statements), nil
	}
	if !actual.FunctionExists || !actual.Exclusive || !sameSQLServerManagedPartitionDefinition(actual, desired) {
		return nil, ErrSQLServerManagedTableUnconfirmed
	}
	if actual.SchemeExists {
		if actual.Filegroup != desired.Filegroup {
			return nil, ErrSQLServerManagedTableUnconfirmed
		}
		return slices.Clone(request.Statements[sqlServerNativeRepairPartitionStatementCount:]), nil
	}
	if !metadata.permissions.AlterAnyDataspace {
		return nil, ErrSchemaExecutionPermissionDenied
	}
	return slices.Clone(request.Statements[2:]), nil
}

func verifySQLServerManagedRepairMetadataInTransaction(
	ctx context.Context,
	tx *sql.Tx,
	capability *SQLServerRecordingCapability,
) (sqlServerRecordingMetadata, error) {
	metadata := inspectSQLServerRecordingMetadataWithQuerier(ctx, tx, capability.Database, capability.Schema)
	switch metadata.metadataStatus {
	case sqlServerRecordingMetadataForbidden:
		return metadata, errors.Join(ErrSchemaExecutionPermissionDenied, ErrSQLServerManagedTableForbidden)
	case sqlServerRecordingMetadataVerified:
	default:
		if slices.Contains(metadata.reasons, "database_scope_mismatch") || slices.Contains(metadata.reasons, "schema_scope_mismatch") {
			return metadata, ErrSQLServerManagedScopeChanged
		}
		return metadata, ErrSQLServerManagedTableUnconfirmed
	}
	connector := &schema.DatabaseConnector{
		ID: capability.ConnectorID, Kind: schema.DatabaseConnectorKindSQLServer,
		IdentityRevision: capability.ConnectorRevision,
	}
	actual := projectSQLServerRecordingCapability(connector, metadata, true)
	if metadata.actualDatabase != capability.Database || metadata.actualSchema != capability.Schema ||
		!sqlServerRecordingCapabilityScopeMatches(capability, actual) {
		return metadata, ErrSQLServerManagedScopeChanged
	}
	if !sqlServerRecordingPermissionsAllowManagedLong(metadata.permissions) {
		return metadata, ErrSchemaExecutionPermissionDenied
	}
	return metadata, nil
}

func parseSQLServerNativeRepairDefinition(
	capability *SQLServerRecordingCapability,
	request SQLServerManagedRepairRequest,
) (*SQLServerManagedPartitionEvidence, bool, error) {
	containsPartitionDDL := false
	for _, statement := range request.Statements {
		containsPartitionDDL = containsPartitionDDL || strings.Contains(statement, "PARTITION FUNCTION") ||
			strings.Contains(statement, "PARTITION SCHEME")
	}
	if !containsPartitionDDL {
		return nil, false, nil
	}
	if capability == nil || strings.TrimSpace(capability.Schema) == "" || request.TableName == "" ||
		len(request.Statements) < sqlServerNativeRepairPartitionStatementCount {
		return nil, true, ErrSchemaExecutionUnsupported
	}
	functionName, schemeName := "pf_"+request.TableName, "ps_"+request.TableName
	functionIdentifier := quoteIdentifier(schema.DatabaseConnectorKindSQLServer, functionName)
	schemeIdentifier := quoteIdentifier(schema.DatabaseConnectorKindSQLServer, schemeName)
	prefix := "CREATE PARTITION FUNCTION " + functionIdentifier + " (datetime2(6)) AS RANGE RIGHT FOR VALUES ("
	statement := request.Statements[0]
	if !strings.HasPrefix(statement, prefix) || !strings.HasSuffix(statement, ");") {
		return nil, true, ErrSchemaExecutionUnsupported
	}
	rawValues := strings.TrimSuffix(strings.TrimPrefix(statement, prefix), ");")
	values := strings.Split(rawValues, ", ")
	if len(values) < 5 || len(values) > 15000 {
		return nil, true, ErrSchemaExecutionUnsupported
	}
	partitions := make([]SQLServerPartitionBoundary, 0, len(values)+1)
	var previous time.Time
	for index, value := range values {
		literal, ok := strings.CutPrefix(value, "N'")
		if !ok {
			return nil, true, ErrSchemaExecutionUnsupported
		}
		literal, ok = strings.CutSuffix(literal, "'")
		if !ok {
			return nil, true, ErrSchemaExecutionUnsupported
		}
		at, err := time.ParseInLocation(sqlServerNativeBoundaryTimeLayout, literal, time.UTC)
		if err != nil || at.Format(sqlServerNativeBoundaryTimeLayout) != literal || at.Day() != 1 ||
			at.Hour() != 0 || at.Minute() != 0 || at.Second() != 0 || at.Nanosecond() != 0 ||
			(!previous.IsZero() && !at.After(previous)) {
			return nil, true, ErrSchemaExecutionUnsupported
		}
		name := sqlServerNativePastPartitionName
		if !previous.IsZero() {
			name = "p_" + previous.Format("200601")
			if !at.Equal(previous.AddDate(0, 1, 0)) {
				name = "p_gap_" + at.Format("200601")
			}
		}
		partitions = append(partitions, SQLServerPartitionBoundary{
			Name: name, LessThan: literal, Ordinal: index + 1,
		})
		previous = at
	}
	partitions = append(partitions, SQLServerPartitionBoundary{
		Name: sqlServerNativeFuturePartitionName, LessThan: "MAXVALUE", Ordinal: len(partitions) + 1,
	})
	desired := &SQLServerManagedPartitionEvidence{
		Method: sqlServerNativeRangeRight, Column: "bucket_start", FunctionName: functionName, SchemeName: schemeName,
		Filegroup: sqlServerNativePrimaryFilegroup, Partitions: partitions,
	}
	owner := strings.ReplaceAll(capability.Schema+"."+request.TableName, "'", "''")
	marker := "gateway.recording.partition_owner"
	want := []string{
		statement,
		"EXEC sys.sp_addextendedproperty @name=N'" + marker + "', @value=N'" + owner +
			"', @level0type=N'PARTITION FUNCTION', @level0name=N'" + functionName + "';",
		"CREATE PARTITION SCHEME " + schemeIdentifier + " AS PARTITION " + functionIdentifier + " ALL TO ([PRIMARY]);",
		"EXEC sys.sp_addextendedproperty @name=N'" + marker + "', @value=N'" + owner +
			"', @level0type=N'PARTITION SCHEME', @level0name=N'" + schemeName + "';",
	}
	if !slices.Equal(request.Statements[:sqlServerNativeRepairPartitionStatementCount], want) {
		return nil, true, ErrSchemaExecutionUnsupported
	}
	return desired, true, nil
}

func sameSQLServerManagedPartitionDefinition(actual, desired *SQLServerManagedPartitionEvidence) bool {
	return actual != nil && desired != nil && actual.Method == desired.Method && actual.Column == desired.Column &&
		actual.FunctionName == desired.FunctionName && actual.SchemeName == desired.SchemeName &&
		slices.Equal(actual.Partitions, desired.Partitions)
}
