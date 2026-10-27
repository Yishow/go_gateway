package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"

	mysqldriver "github.com/go-sql-driver/mysql"
)

// MySQLManagedRowProof freezes the exact catalog evidence a MySQL managed row
// insert is allowed to target. Group and Receipt are complete inspections, not
// a second metadata model; the producer snapshots them before delivery.
type MySQLManagedRowProof struct {
	ConnectorID       string                  `json:"connector_id"`
	ConnectorRevision string                  `json:"connector_revision"`
	Group             *MySQLManagedInspection `json:"group"`
	Receipt           *MySQLManagedInspection `json:"receipt"`
}

var mysqlCommittedAtPattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{6}(Z|[+-]\d{2}:\d{2})$`)

// ErrMySQLManagedPhysicalMismatch marks an actual catalog drift from the
// frozen managed physical proof. It is not used for a missing table or an
// untrusted/incomplete proof.
var ErrMySQLManagedPhysicalMismatch = errors.New("mysql managed physical proof mismatch")

// ErrMySQLManagedPartitionChanged marks a group-table partition-only change
// between proof resolution and the insert. It is transient: the outbox may
// retry and resolve the fresh authorized proof after maintenance settles.
var ErrMySQLManagedPartitionChanged = errors.New("mysql managed partition changed")

func validateMySQLManagedRowProof(req GroupInsertRequest) error {
	unsupported := func(reason string) error {
		return fmt.Errorf("%w: %s", ErrGroupInsertUnsupported, reason)
	}
	proof := req.MySQL
	if proof == nil || proof.Group == nil || proof.Receipt == nil {
		return unsupported("mysql managed proof is required")
	}
	if strings.TrimSpace(proof.ConnectorID) == "" || strings.TrimSpace(proof.ConnectorRevision) == "" {
		return unsupported("mysql managed proof connector binding is required")
	}
	if req.Strategy != GroupEffectReceipt {
		return unsupported("mysql managed insert requires the receipt strategy")
	}
	group, receipt := proof.Group, proof.Receipt
	if group.Inspection == nil || receipt.Inspection == nil || group.Evidence == nil || receipt.Evidence == nil {
		return unsupported("mysql managed proof is incomplete")
	}
	if group.Inspection.Status != TableInspectionExists || receipt.Inspection.Status != TableInspectionExists {
		return unsupported("mysql managed proof tables must exist")
	}
	if group.Indexes == nil || receipt.Indexes == nil {
		return unsupported("mysql managed proof indexes were not inspected")
	}
	database := strings.TrimSpace(req.SchemaName)
	if database == "" || database != strings.TrimSpace(group.Evidence.Database) || database != strings.TrimSpace(receipt.Evidence.Database) {
		return unsupported("mysql database scope does not match the proof")
	}
	if strings.TrimSpace(req.TableName) == "" || req.TableName != group.Inspection.Table {
		return unsupported("mysql group table does not match the proof")
	}
	if receipt.Inspection.Table != EffectReceiptTable {
		return unsupported("mysql receipt table does not match the proof")
	}
	for _, evidence := range []*MySQLManagedTableEvidence{group.Evidence, receipt.Evidence} {
		if evidence.ServerFamily != group.Evidence.ServerFamily ||
			evidence.ServerVersion != group.Evidence.ServerVersion ||
			evidence.CapabilityRevision != group.Evidence.CapabilityRevision {
			return unsupported("mysql proof server bindings differ")
		}
		if evidence.ServerFamily != mysqlManagedFamilyMySQL && evidence.ServerFamily != mysqlManagedFamilyMariaDB {
			return unsupported("mysql proof server family is unknown")
		}
		if !MySQLRecordingVersionVerified(evidence.ServerFamily, evidence.ServerVersion) {
			return unsupported("mysql proof server version is not a verified recording version")
		}
		if !strings.EqualFold(evidence.Engine, mysqlManagedEngineInnoDB) || !strings.EqualFold(evidence.TableType, "BASE TABLE") {
			return unsupported("mysql proof table is not an InnoDB base table")
		}
		if len(evidence.PrimaryKey) == 0 {
			return unsupported("mysql proof primary key is missing")
		}
	}
	if strings.TrimSpace(group.Evidence.CapabilityRevision) == "" {
		return unsupported("mysql proof capability revision is empty")
	}
	if len(group.Inspection.Columns) == 0 || len(receipt.Inspection.Columns) == 0 {
		return unsupported("mysql proof columns are missing")
	}
	return nil
}

func insertMySQLGroupRowInTx(ctx context.Context, tx *sql.Tx, req GroupInsertRequest) (GroupInsertResult, error) {
	if err := validateMySQLManagedRowProof(req); err != nil {
		return GroupInsertResult{}, err
	}
	proof := req.MySQL
	group, receipt := proof.Group, proof.Receipt
	committedAt, err := parseMySQLCommittedAt(req.CommittedAt)
	if err != nil {
		return GroupInsertResult{}, fmt.Errorf("%w: committed_at must be UTC microseconds", ErrGroupInsertUnsupported)
	}
	database := strings.TrimSpace(group.Evidence.Database)
	groupTable := group.Inspection.Table
	receiptTable := receipt.Inspection.Table

	for _, table := range []string{groupTable, receiptTable} {
		if err := holdMySQLMDL(ctx, tx, database, table); err != nil {
			return GroupInsertResult{}, fmt.Errorf("acquire mysql metadata lock on %s: %w", table, err)
		}
	}
	family := group.Evidence.ServerFamily
	actualGroup, err := readMySQLManagedTableFacts(ctx, tx, family, database, groupTable)
	if err != nil {
		return GroupInsertResult{}, fmt.Errorf("read mysql group catalog: %w", err)
	}
	if err := matchMySQLManagedInspection(group, actualGroup); err != nil {
		return GroupInsertResult{}, err
	}
	actualReceipt, err := readMySQLManagedTableFacts(ctx, tx, family, database, receiptTable)
	if err != nil {
		return GroupInsertResult{}, fmt.Errorf("read mysql receipt catalog: %w", err)
	}
	if err := matchMySQLManagedInspection(receipt, actualReceipt); err != nil {
		return GroupInsertResult{}, err
	}

	receipts := qualifiedTableName(schema.DatabaseConnectorKindMySQL, database, receiptTable)
	_, err = tx.ExecContext(ctx, receiptInsertStatement(schema.DatabaseConnectorKindMySQL, receipts),
		req.Row.EffectKey, req.PayloadDigest, committedAt)
	switch {
	case err == nil:
	case isMySQLDuplicateKey(err):
		var existing string
		if readErr := tx.QueryRowContext(ctx, mysqlReceiptLookupForUpdateStatement(receipts), req.Row.EffectKey).Scan(&existing); readErr != nil {
			return GroupInsertResult{}, fmt.Errorf("read destination receipt: %w", readErr)
		}
		if existing == req.PayloadDigest {
			return GroupInsertResult{AlreadyCommitted: true}, nil
		}
		return GroupInsertResult{}, ErrReceiptDigestMismatch
	default:
		return GroupInsertResult{}, fmt.Errorf("insert destination receipt: %w", err)
	}

	table := qualifiedTableName(schema.DatabaseConnectorKindMySQL, database, groupTable)
	columns := make([]string, 0, len(req.Row.Cells))
	args := make([]any, 0, len(req.Row.Cells))
	for _, cell := range req.Row.Cells {
		columns = append(columns, quoteIdentifier(schema.DatabaseConnectorKindMySQL, cell.Column))
		args = append(args, normalizeDBValue(cell.Value))
	}
	_, err = tx.ExecContext(ctx, groupInsertStatement(req, table, columns), args...)
	if err != nil {
		if isMySQLDuplicateKey(err) {
			return GroupInsertResult{}, ErrExistingRowDiffers
		}
		return GroupInsertResult{}, fmt.Errorf("insert destination row: %w", err)
	}
	return GroupInsertResult{}, nil
}

func holdMySQLMDL(ctx context.Context, tx *sql.Tx, database, table string) error {
	rows, err := tx.QueryContext(ctx, "SELECT 1 FROM "+qualifiedTableName(schema.DatabaseConnectorKindMySQL, database, table)+" LIMIT 0") // #nosec G202 -- Both identifiers are quoted by qualifiedTableName; no values are interpolated.
	if err != nil {
		return err
	}
	return rows.Close()
}

func matchMySQLManagedInspection(expected *MySQLManagedInspection, facts mysqlManagedTableFacts) error {
	if facts.status != TableInspectionExists {
		return fmt.Errorf("%w: mysql managed table is %s", ErrGroupInsertUnsupported, facts.status)
	}
	server := mysqlManagedServerFacts{
		family: expected.Evidence.ServerFamily, version: expected.Evidence.ServerVersion,
		capabilityRevision: expected.Evidence.CapabilityRevision, database: expected.Evidence.Database,
		metadataStatus: mysqlRecordingMetadataVerified, transactionalEngine: true,
	}
	actual, err := buildMySQLManagedInspection(expected.Evidence.Database, expected.Inspection.Table, server, facts)
	if err != nil {
		if errors.Is(err, ErrMySQLManagedTableUnconfirmed) || errors.Is(err, ErrMySQLManagedTableForbidden) {
			return fmt.Errorf("%w: verify mysql managed catalog: %w", ErrGroupInsertUnsupported, err)
		}
		return fmt.Errorf("%w: %w: verify mysql managed catalog: %w", ErrMySQLManagedPhysicalMismatch, ErrGroupInsertUnsupported, err)
	}
	if !mysqlManagedEvidenceEqual(expected.Inspection, actual.Inspection) ||
		!mysqlManagedTableEvidenceEqual(expected.Evidence, actual.Evidence) ||
		!mysqlManagedIndexesEqual(expected.Indexes, actual.Indexes) {
		if mysqlManagedPartitionOnlyChange(expected, actual) {
			return fmt.Errorf("%w: mysql managed partition-only change", ErrMySQLManagedPartitionChanged)
		}
		return fmt.Errorf("%w: %w: mysql managed table physical proof changed", ErrMySQLManagedPhysicalMismatch, ErrGroupInsertUnsupported)
	}
	return nil
}

func mysqlManagedPartitionOnlyChange(expected, actual *MySQLManagedInspection) bool {
	if expected == nil || actual == nil || expected.Inspection == nil || actual.Inspection == nil ||
		expected.Evidence == nil || actual.Evidence == nil {
		return false
	}
	if !strings.HasPrefix(expected.Inspection.Table, "gw_long_") || expected.Inspection.Table == EffectReceiptTable {
		return false
	}
	if expected.Evidence.Partitions == nil || actual.Evidence.Partitions == nil ||
		sameMySQLPartitionEvidence(expected.Evidence.Partitions, actual.Evidence.Partitions) {
		return false
	}
	if !mysqlManagedEvidenceEqual(expected.Inspection, actual.Inspection) ||
		!mysqlManagedIndexesEqual(expected.Indexes, actual.Indexes) {
		return false
	}
	return mysqlManagedEvidenceExceptPartitionsEqual(expected.Evidence, actual.Evidence)
}

func mysqlManagedEvidenceExceptPartitionsEqual(left, right *MySQLManagedTableEvidence) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	leftCopy, rightCopy := *left, *right
	leftCopy.Partitions, rightCopy.Partitions = nil, nil
	return mysqlManagedTableEvidenceEqual(&leftCopy, &rightCopy)
}

func mysqlManagedEvidenceEqual(left, right *TableInspection) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return left.Status == right.Status && left.Schema == right.Schema && left.Table == right.Table &&
		slices.Equal(left.Columns, right.Columns)
}

func mysqlManagedTableEvidenceEqual(left, right *MySQLManagedTableEvidence) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	if left.ServerFamily != right.ServerFamily || left.ServerVersion != right.ServerVersion ||
		left.CapabilityRevision != right.CapabilityRevision || left.Database != right.Database ||
		left.Engine != right.Engine || left.Charset != right.Charset ||
		left.Collation != right.Collation || left.TableType != right.TableType ||
		!slices.Equal(left.PrimaryKey, right.PrimaryKey) ||
		!maps.Equal(left.ColumnCollations, right.ColumnCollations) ||
		!maps.Equal(left.JSONColumns, right.JSONColumns) ||
		!sameMySQLPartitionEvidence(left.Partitions, right.Partitions) ||
		len(left.ColumnDefaults) != len(right.ColumnDefaults) {
		return false
	}
	for name, want := range left.ColumnDefaults {
		got := right.ColumnDefaults[name]
		if (want == nil) != (got == nil) || (want != nil && *want != *got) {
			return false
		}
	}
	return true
}

func mysqlManagedIndexesEqual(left, right []IndexInfo) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i].Name != right[i].Name || left[i].Unique != right[i].Unique ||
			!slices.Equal(left[i].Columns, right[i].Columns) {
			return false
		}
	}
	return true
}

func parseMySQLCommittedAt(raw string) (time.Time, error) {
	value := strings.TrimSpace(raw)
	if !mysqlCommittedAtPattern.MatchString(value) {
		return time.Time{}, fmt.Errorf("invalid committed_at")
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, err
	}
	if parsed.Nanosecond()%1000 != 0 {
		return time.Time{}, fmt.Errorf("committed_at exceeds microseconds")
	}
	return parsed.UTC(), nil
}

func mysqlReceiptLookupForUpdateStatement(receipts string) string {
	return fmt.Sprintf("SELECT payload_digest FROM %s WHERE effect_key = ? FOR UPDATE", receipts)
}

func isMySQLDuplicateKey(err error) bool {
	var mysqlErr *mysqldriver.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

func mysqlDeliveryDatabase(config ConnectionConfig) (string, error) {
	dbConfig, err := buildExternalDBConfig(schema.DatabaseConnectorKindMySQL, config)
	if err != nil {
		return "", err
	}
	driverConfig, err := mysqldriver.ParseDSN(dbConfig.DSN)
	if err != nil {
		return "", err
	}
	if !mysqlRecordingTransportVerified(driverConfig) {
		return "", fmt.Errorf("mysql delivery requires verified TLS")
	}
	database := strings.TrimSpace(driverConfig.DBName)
	if database == "" {
		return "", fmt.Errorf("mysql delivery requires a database")
	}
	return database, nil
}

func mysqlDeliveryCapability(ctx context.Context, connector *schema.DatabaseConnector, db *sql.DB, config ConnectionConfig) (*MySQLRecordingCapability, error) {
	database, err := mysqlDeliveryDatabase(config)
	if err != nil {
		return nil, err
	}
	metadata := inspectMySQLRecordingMetadata(ctx, db, database)
	capability := projectMySQLRecordingCapability(connector, metadata, true)
	if capability.MetadataStatus != mysqlRecordingMetadataVerified || !capability.TransactionalEngine {
		return nil, fmt.Errorf("mysql delivery metadata is not verified")
	}
	if !MySQLRecordingVersionVerified(capability.ServerFamily, capability.ServerVersion) {
		return nil, fmt.Errorf("mysql delivery server is not a verified recording version")
	}
	return capability, nil
}
