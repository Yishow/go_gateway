package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// SQLServerManagedRowProof freezes the complete physical group/receipt catalog
// under one saved connector identity. It carries no connection configuration.
type SQLServerManagedRowProof struct {
	ConnectorID       string                      `json:"connector_id"`
	ConnectorRevision string                      `json:"connector_revision"`
	Group             *SQLServerManagedInspection `json:"group"`
	Receipt           *SQLServerManagedInspection `json:"receipt"`
}

// ErrSQLServerManagedPhysicalMismatch means a frozen SQL Server proof no longer matches catalog facts.
var ErrSQLServerManagedPhysicalMismatch = errors.New("sqlserver managed physical proof mismatch")

func insertSQLServerGroupRowInTx(ctx context.Context, tx *sql.Tx, req GroupInsertRequest) (GroupInsertResult, error) {
	proof := req.SQLServer
	group, receipt := proof.Group, proof.Receipt
	committedAt, err := parseSQLServerCommittedAt(req.CommittedAt)
	if err != nil {
		return GroupInsertResult{}, fmt.Errorf("%w: committed_at must represent UTC microseconds", ErrGroupInsertUnsupported)
	}
	receiptTable := qualifiedTableName(schema.DatabaseConnectorKindSQLServer, req.SchemaName, receipt.Inspection.Table)
	var digest, receiptOwner sql.NullString
	err = tx.QueryRowContext(ctx, sqlServerReceiptLockStatement(receiptTable), req.Row.EffectKey).Scan(&digest, &receiptOwner)
	receiptExists := err == nil
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return GroupInsertResult{}, fmt.Errorf("read SQL Server effect receipt: %w", err)
	}

	existingRow, groupRowExists, err := readSQLServerManagedGroupRow(ctx, tx, req)
	if err != nil {
		return GroupInsertResult{}, fmt.Errorf("read SQL Server managed row: %w", err)
	}
	if err := verifySQLServerManagedTransactionScope(ctx, tx, group.Evidence); err != nil {
		return GroupInsertResult{}, err
	}
	if err := verifySQLServerManagedSchema(ctx, tx, group.Evidence.Schema); err != nil {
		return GroupInsertResult{}, err
	}
	if err := verifySQLServerManagedCatalog(ctx, tx, group); err != nil {
		return GroupInsertResult{}, err
	}
	if err := verifySQLServerManagedCatalog(ctx, tx, receipt); err != nil {
		return GroupInsertResult{}, err
	}

	if receiptExists {
		if !digest.Valid || digest.String != req.PayloadDigest {
			return GroupInsertResult{}, ErrReceiptDigestMismatch
		}
		if !receiptOwner.Valid || receiptOwner.String != sqlServerManagedOwnerMarkerValue {
			return GroupInsertResult{}, fmt.Errorf("%w: SQL Server receipt owner marker differs", ErrGroupInsertUnsupported)
		}
		if !groupRowExists || !sqlServerManagedRowMatches(req, existingRow) {
			return GroupInsertResult{}, ErrExistingRowDiffers
		}
		return GroupInsertResult{AlreadyCommitted: true}, nil
	}
	if groupRowExists {
		return GroupInsertResult{}, ErrExistingRowDiffers
	}

	groupTable := qualifiedTableName(schema.DatabaseConnectorKindSQLServer, req.SchemaName, group.Inspection.Table)
	columns := make([]string, 0, len(req.Row.Cells))
	args := make([]any, 0, len(req.Row.Cells))
	for _, cell := range req.Row.Cells {
		columns = append(columns, quoteIdentifier(schema.DatabaseConnectorKindSQLServer, cell.Column))
		args = append(args, normalizeDBValue(cell.Value))
	}
	if _, err := tx.ExecContext(ctx, groupInsertStatement(req, groupTable, columns), args...); err != nil {
		return GroupInsertResult{}, fmt.Errorf("insert SQL Server managed row: %w", err)
	}

	receiptColumns := []string{sqlServerManagedReceiptEffectKeyColumn, "payload_digest", "committed_at", sqlServerManagedReceiptOwnerColumn}
	quotedReceiptColumns := make([]string, len(receiptColumns))
	for index, column := range receiptColumns {
		quotedReceiptColumns[index] = quoteIdentifier(schema.DatabaseConnectorKindSQLServer, column)
	}
	receiptArgs := []any{req.Row.EffectKey, req.PayloadDigest, committedAt, sqlServerManagedOwnerMarkerValue}
	if _, err := tx.ExecContext(ctx, sqlServerReceiptInsertStatement(receiptTable, quotedReceiptColumns), receiptArgs...); err != nil {
		return GroupInsertResult{}, fmt.Errorf("insert SQL Server effect receipt: %w", err)
	}
	return GroupInsertResult{}, nil
}

func sqlServerReceiptLockStatement(receiptTable string) string {
	return fmt.Sprintf(
		"SELECT %s, %s FROM %s WITH (UPDLOCK, HOLDLOCK) WHERE %s = @p1",
		quoteIdentifier(schema.DatabaseConnectorKindSQLServer, "payload_digest"),
		quoteIdentifier(schema.DatabaseConnectorKindSQLServer, sqlServerManagedReceiptOwnerColumn), receiptTable,
		quoteIdentifier(schema.DatabaseConnectorKindSQLServer, sqlServerManagedReceiptEffectKeyColumn),
	)
}

func readSQLServerManagedGroupRow(ctx context.Context, tx *sql.Tx, req GroupInsertRequest) (values []any, exists bool, err error) {
	columns := make([]string, len(req.Row.Cells))
	for index, cell := range req.Row.Cells {
		columns[index] = quoteIdentifier(schema.DatabaseConnectorKindSQLServer, cell.Column)
	}
	table := qualifiedTableName(schema.DatabaseConnectorKindSQLServer, req.SchemaName, req.TableName)
	// #nosec G201 -- validated managed identifiers are quoted; lookup values use bound @p1/@p2 parameters.
	statement := fmt.Sprintf(
		"SELECT %s FROM %s WITH (UPDLOCK, HOLDLOCK) WHERE %s = @p1 AND %s = @p2",
		strings.Join(columns, ", "), table,
		quoteIdentifier(schema.DatabaseConnectorKindSQLServer, longColBucketStart),
		quoteIdentifier(schema.DatabaseConnectorKindSQLServer, longColRecordID),
	)
	values = make([]any, len(columns))
	destinations := make([]any, len(values))
	for index := range values {
		destinations[index] = &values[index]
	}
	scanErr := tx.QueryRowContext(ctx, statement, req.Row.BucketStart, req.Row.RecordID).Scan(destinations...)
	if errors.Is(scanErr, sql.ErrNoRows) {
		return nil, false, nil
	}
	if scanErr != nil {
		return nil, false, scanErr
	}
	return values, true, nil
}

func verifySQLServerManagedTransactionScope(ctx context.Context, tx *sql.Tx, expected *SQLServerManagedTableEvidence) error {
	var version, edition, database sql.NullString
	var engineEdition, compatibility sql.NullInt64
	if err := tx.QueryRowContext(ctx, sqlServerRecordingMetadataQuery).Scan(
		&version, &edition, &engineEdition, &database, &compatibility,
	); err != nil {
		return fmt.Errorf("verify SQL Server transaction scope: %w", err)
	}
	if !version.Valid || !edition.Valid || !engineEdition.Valid || !database.Valid || !compatibility.Valid {
		return fmt.Errorf("%w: %w", ErrSQLServerManagedTableUnconfirmed, ErrGroupInsertUnsupported)
	}
	if version.String != expected.ServerVersion || edition.String != expected.ServerEdition ||
		engineEdition.Int64 != expected.EngineEdition || database.String != expected.Database ||
		compatibility.Int64 != expected.CompatibilityLevel {
		return fmt.Errorf("%w: %w", ErrSQLServerManagedScopeChanged, ErrGroupInsertUnsupported)
	}
	return nil
}

func verifySQLServerManagedSchema(ctx context.Context, tx *sql.Tx, schemaName string) error {
	var schemaID sql.NullInt64
	if err := tx.QueryRowContext(ctx, "SELECT SCHEMA_ID(@p1)", schemaName).Scan(&schemaID); err != nil {
		return fmt.Errorf("verify SQL Server schema scope: %w", err)
	}
	if !schemaID.Valid || schemaID.Int64 <= 0 {
		return fmt.Errorf("%w: %w", ErrSQLServerManagedTableUnconfirmed, ErrGroupInsertUnsupported)
	}
	return nil
}

func verifySQLServerManagedCatalog(ctx context.Context, tx *sql.Tx, expected *SQLServerManagedInspection) error {
	facts, err := inspectSQLServerManagedCatalog(ctx, tx, expected.Inspection.Schema, expected.Inspection.Table)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %w", ErrSQLServerManagedTableUnconfirmed, ErrGroupInsertUnsupported)
	}
	if errors.Is(err, ErrSQLServerManagedScopeChanged) {
		return fmt.Errorf("%w: %w", ErrSQLServerManagedScopeChanged, ErrGroupInsertUnsupported)
	}
	if err != nil {
		return fmt.Errorf("inspect SQL Server managed catalog: %w", err)
	}
	if !sqlServerManagedCatalogMatches(expected, facts) {
		return sqlServerManagedPhysicalMismatch()
	}
	return nil
}

func sqlServerManagedCatalogMatches(expected *SQLServerManagedInspection, actual *sqlServerManagedCatalogFacts) bool {
	if expected == nil || expected.Inspection == nil || expected.Evidence == nil || actual == nil {
		return false
	}
	inspection, evidence := expected.Inspection, expected.Evidence
	if actual.object.objectID != evidence.ObjectID || actual.object.schema != inspection.Schema ||
		actual.object.name != inspection.Table || actual.object.tableType != evidence.TableType ||
		actual.object.tableType != sqlServerManagedUserTableType || !actual.object.userTable ||
		!sqlServerManagedCatalogPartitionMatches(evidence, actual) || actual.indexUnusable ||
		!slices.Equal(actual.columns, inspection.Columns) || !slices.Equal(actual.primaryKey, evidence.PrimaryKey) ||
		!sqlServerManagedIndexesEqual(actual.indexes, expected.Indexes) ||
		!maps.Equal(actual.columnCollations, evidence.ColumnCollations) ||
		!sqlServerManagedDefaultsEqual(actual.columnDefaults, evidence.ColumnDefaults) ||
		actual.provenanceJSONChecked != evidence.ProvenanceJSONChecked ||
		actual.repairOperationID != evidence.RepairOperationID {
		return false
	}
	return true
}

func sqlServerManagedCatalogPartitionMatches(expected *SQLServerManagedTableEvidence, actual *sqlServerManagedCatalogFacts) bool {
	if expected == nil || actual == nil {
		return false
	}
	if expected.NativeMonthly == nil {
		if !actual.plain || actual.nativeMonthly != nil {
			return false
		}
		return !expected.PartitionInspected || actual.partitionInspected
	}
	return expected.PartitionInspected && actual.partitionInspected && !actual.plain &&
		sqlServerManagedPartitionEvidenceEqual(expected.NativeMonthly, actual.nativeMonthly)
}

func sqlServerManagedPartitionEvidenceEqual(left, right *SQLServerManagedPartitionEvidence) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return left.Method == right.Method && left.Column == right.Column &&
		left.FunctionName == right.FunctionName && left.SchemeName == right.SchemeName &&
		left.Filegroup == right.Filegroup && slices.Equal(left.Partitions, right.Partitions) &&
		left.FunctionExists == right.FunctionExists && left.SchemeExists == right.SchemeExists &&
		left.Exclusive == right.Exclusive && left.IndexesAligned == right.IndexesAligned
}

func sqlServerManagedDefaultsEqual(actual, expected map[string]*string) bool {
	if len(actual) != len(expected) {
		return false
	}
	for name, want := range expected {
		got, exists := actual[name]
		if !exists || (got == nil) != (want == nil) || (got != nil && *got != *want) {
			return false
		}
	}
	return true
}

func sqlServerManagedRowMatches(req GroupInsertRequest, actual []any) bool {
	if len(actual) != len(req.Row.Cells) {
		return false
	}
	for index, cell := range req.Row.Cells {
		want, got := normalizeDBValue(cell.Value), actual[index]
		if want == nil {
			if got != nil {
				return false
			}
			continue
		}
		if cell.Column == longColValueDecimal || cell.Column == longColDeltaValue {
			wantText, wantOK := want.(string)
			gotText, gotOK := readbackString(got)
			if !wantOK || !gotOK {
				return false
			}
			wantDecimal, wantErr := measurement.NewDecimal(wantText)
			gotDecimal, gotErr := measurement.NewDecimal(gotText)
			if wantErr != nil || gotErr != nil || !wantDecimal.Equal(gotDecimal) {
				return false
			}
			continue
		}
		if expectedTime, ok := want.(time.Time); ok {
			actualTime, valid := got.(time.Time)
			if !valid || !sqlServerUTCMicrosecond(actualTime) || !actualTime.Equal(expectedTime) {
				return false
			}
			continue
		}
		if !groupCellEqual(want, got) {
			return false
		}
	}
	return true
}

func sqlServerReceiptInsertStatement(receiptTable string, columns []string) string {
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", receiptTable, strings.Join(columns, ", "),
		strings.Join(buildPlaceholders(schema.DatabaseConnectorKindSQLServer, len(columns)), ", "))
}
