package recordingplan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

const (
	sqlserverReceiptEffectColumn = "effect_key"
	sqlserverReceiptDigestColumn = "payload_digest"
	sqlserverReceiptTimeColumn   = "committed_at"
)

// TargetSQLServerTableEvidence is the saved destination's actual SQL Server catalog evidence.
// Missing tables still require server version, edition, compatibility, database, and metadata proof.
type TargetSQLServerTableEvidence struct {
	ServerVersion         string             `json:"server_version,omitempty"`
	ServerEdition         string             `json:"server_edition,omitempty"`
	CapabilityRevision    string             `json:"capability_revision,omitempty"`
	Database              string             `json:"database,omitempty"`
	Schema                string             `json:"schema,omitempty"`
	EngineEdition         int64              `json:"engine_edition,omitzero"`
	CompatibilityLevel    int64              `json:"compatibility_level,omitzero"`
	MetadataVerified      bool               `json:"metadata_verified"`
	TableType             string             `json:"table_type,omitempty"` // USER_TABLE
	ObjectID              int64              `json:"object_id,omitzero"`
	Plain                 bool               `json:"plain"`
	PrimaryKey            []string           `json:"primary_key,omitempty"`
	ColumnCollations      map[string]string  `json:"column_collations,omitempty"`
	ColumnDefaults        map[string]*string `json:"column_defaults,omitempty"`
	ProvenanceJSONChecked bool               `json:"provenance_json_checked"`
	RepairOperationID     string             `json:"repair_operation_id,omitempty"`
	// NativeMonthly is the inspected owned partition projection of a group
	// table. It stays nil for plain tables and for a verified absence of the
	// scoped function/scheme pair; PartitionInspected distinguishes that from
	// catalog evidence that was not read.
	NativeMonthly      *TargetSQLServerPartitionEvidence `json:"native_monthly,omitempty"`
	PartitionInspected bool                              `json:"partition_inspected,omitzero"`
}

// TargetSQLServerPartitionEvidence is the live catalog projection of the
// scoped SQL Server partition function and scheme. Exclusive requires the
// canonical owner marker on both objects and no foreign references.
type TargetSQLServerPartitionEvidence struct {
	Method         string                       `json:"method"`
	Column         string                       `json:"column"`
	FunctionName   string                       `json:"function_name"`
	SchemeName     string                       `json:"scheme_name"`
	Filegroup      string                       `json:"filegroup"`
	Partitions     []SQLServerPartitionBoundary `json:"partitions"`
	FunctionExists bool                         `json:"function_exists"`
	SchemeExists   bool                         `json:"scheme_exists"`
	Exclusive      bool                         `json:"exclusive"`
	IndexesAligned bool                         `json:"indexes_aligned"`
}

// inspectSQLServerGroupTables inspects the group table and receipt table from the target.
func inspectSQLServerGroupTables(ctx context.Context, layout GroupSchemaLayout, database string, inspect TargetInspector) (groupInspection, error) {
	var found groupInspection
	database = strings.TrimSpace(database)
	if inspect == nil || database == "" {
		return found, ErrTargetInspectionUnconfirmed
	}
	if layout.SQLServer == nil {
		return found, ErrInvalidGroupSchemaLayout
	}
	if err := validateSQLServerPreviewScope(database, layout.SQLServer.Schema, dialectSQLServer, *layout.SQLServer); err != nil {
		return found, err
	}
	if _, err := GenerateGroupSchemaDDL(dialectSQLServer, layout); err != nil {
		return found, err
	}
	group, err := inspect(ctx, layout.TableName)
	if err != nil {
		return found, fmt.Errorf("%w: %s: %w", ErrTargetInspectionUnconfirmed, layout.TableName, err)
	}
	receipt, err := inspect(ctx, managedEffectReceiptTable)
	if err != nil {
		return found, fmt.Errorf("%w: %s: %w", ErrTargetInspectionUnconfirmed, managedEffectReceiptTable, err)
	}
	return inspectSQLServerGroupSnapshot(layout, database, group, receipt)
}

func inspectSQLServerGroupSnapshot(layout GroupSchemaLayout, database string, group, receipt TargetTableInspection) (groupInspection, error) {
	var found groupInspection
	database = strings.TrimSpace(database)
	if database == "" || layout.SQLServer == nil {
		return found, ErrTargetInspectionUnconfirmed
	}
	if err := validateSQLServerPreviewScope(database, layout.SQLServer.Schema, dialectSQLServer, *layout.SQLServer); err != nil {
		return found, err
	}

	// 1. Check Permissions First
	if group.Status == inspectionForbidden {
		return found, fmt.Errorf("%w: %w: %s", ErrTargetInspectionUnconfirmed, ErrTargetPermissionDenied, layout.TableName)
	}
	if receipt.Status == inspectionForbidden {
		return found, fmt.Errorf("%w: %w: %s", ErrTargetInspectionUnconfirmed, ErrTargetPermissionDenied, managedEffectReceiptTable)
	}

	// 2. Validate Group Table Binding & Status
	if err := checkSQLServerCatalogBinding(database, layout, group.SQLServer); err != nil {
		return found, err
	}
	groupTable := SchemaPreviewTable{Name: layout.TableName, Action: SchemaTableActionUnchanged}
	switch group.Status {
	case inspectionMissing:
		groupTable.Action = SchemaTableActionCreate
		groupTable.Columns = groupColumnNames(layout.Columns)
		partitionPreview, err := sqlServerPartitionPreviewForMissingGroup(layout, group.SQLServer)
		if err != nil {
			return found, err
		}
		groupTable.SQLServerPartition = partitionPreview
		found.pendingGroup = true
	case inspectionExists:
		if err := validateSQLServerCatalogTable(group, layout.Columns, layout.OwnerColumn, *layout.SQLServer); err != nil {
			return found, err
		}
		if err := validateSQLServerCatalogIndexes(group, layout.Indexes); err != nil {
			return found, err
		}
		partitionPreview, err := sqlServerPartitionPreviewForExistingGroup(layout, group.SQLServer)
		if err != nil {
			return found, err
		}
		groupTable.SQLServerPartition = partitionPreview
		digest, err := sqlserverCatalogDigest(group)
		if err != nil {
			return found, err
		}
		groupTable.CatalogDigest = digest
	default:
		return found, fmt.Errorf("%w: %s is %s", ErrTargetInspectionUnconfirmed, layout.TableName, group.Status)
	}

	// 3. Validate Receipt Table Binding & Status
	if err := checkSQLServerCatalogBinding(database, layout, receipt.SQLServer); err != nil {
		return found, err
	}
	receiptColumns := sqlserverReceiptColumns()
	receiptTable := SchemaPreviewTable{Name: managedEffectReceiptTable, Action: SchemaTableActionUnchanged}
	switch receipt.Status {
	case inspectionMissing:
		receiptTable.Action = SchemaTableActionCreate
		receiptTable.Columns = groupColumnNames(receiptColumns)
		if err := validateSQLServerReceiptHasNoPartition(receipt.SQLServer); err != nil {
			return found, err
		}
		found.pendingReceipt = true
	case inspectionExists:
		if !receipt.IndexesInspected {
			return found, ErrTargetInspectionUnconfirmed
		}
		if len(receipt.Indexes) > 0 {
			return found, ErrIncompatibleExistingTable
		}
		if err := validateSQLServerCatalogTable(receipt, receiptColumns, managedEffectReceiptMarker, *layout.SQLServer); err != nil {
			return found, err
		}
		if err := validateSQLServerReceiptHasNoPartition(receipt.SQLServer); err != nil {
			return found, err
		}
		digest, err := sqlserverCatalogDigest(receipt)
		if err != nil {
			return found, err
		}
		receiptTable.CatalogDigest = digest
	default:
		return found, fmt.Errorf("%w: %s is %s", ErrTargetInspectionUnconfirmed, managedEffectReceiptTable, receipt.Status)
	}

	found.tables = []SchemaPreviewTable{groupTable, receiptTable}
	return found, nil
}

func sqlServerPartitionPreviewForMissingGroup(layout GroupSchemaLayout, evidence *TargetSQLServerTableEvidence) (*SQLServerPartitionPreview, error) {
	if layout.SQLServer.NativeMonthly == nil {
		if evidence != nil && evidence.NativeMonthly != nil {
			return nil, ErrIncompatibleExistingTable
		}
		return nil, nil
	}
	if evidence == nil || !evidence.PartitionInspected {
		return nil, ErrTargetInspectionUnconfirmed
	}
	if evidence.NativeMonthly == nil {
		return sqlServerPartitionPreview(layout, SchemaTableActionCreate), nil
	}
	if !evidence.NativeMonthly.FunctionExists || !evidence.NativeMonthly.SchemeExists {
		return nil, sqlServerPartitionUnavailable("expected function and scheme are not both present")
	}
	if !evidence.NativeMonthly.Exclusive {
		return nil, sqlServerPartitionUnavailable("expected partition objects are shared or lack the owner marker")
	}
	if !sameSQLServerTargetPartition(evidence.NativeMonthly, layout.SQLServer.NativeMonthly) {
		return nil, sqlServerPartitionUnavailable("expected partition objects differ from the confirmed layout")
	}
	return sqlServerPartitionPreview(layout, SchemaTableActionUnchanged), nil
}

func sqlServerPartitionPreviewForExistingGroup(layout GroupSchemaLayout, evidence *TargetSQLServerTableEvidence) (*SQLServerPartitionPreview, error) {
	if layout.SQLServer.NativeMonthly == nil {
		if evidence != nil && evidence.NativeMonthly != nil {
			return nil, ErrIncompatibleExistingTable
		}
		return nil, nil
	}
	if evidence == nil || evidence.NativeMonthly == nil || !evidence.PartitionInspected {
		return nil, sqlServerPartitionUnavailable("partition catalog evidence is missing")
	}
	if !evidence.NativeMonthly.FunctionExists || !evidence.NativeMonthly.SchemeExists {
		return nil, sqlServerPartitionUnavailable("expected function and scheme are not both present")
	}
	if !evidence.NativeMonthly.Exclusive {
		return nil, sqlServerPartitionUnavailable("expected partition objects are shared or lack the owner marker")
	}
	if !evidence.NativeMonthly.IndexesAligned {
		return nil, sqlServerPartitionUnavailable("one or more group indexes are not aligned to the partition scheme")
	}
	if !sameSQLServerTargetPartition(evidence.NativeMonthly, layout.SQLServer.NativeMonthly) {
		return nil, sqlServerPartitionUnavailable("expected partition objects differ from the confirmed layout")
	}
	return sqlServerPartitionPreview(layout, SchemaTableActionUnchanged), nil
}

func validateSQLServerReceiptHasNoPartition(evidence *TargetSQLServerTableEvidence) error {
	if evidence != nil && evidence.NativeMonthly != nil {
		return fmt.Errorf("%w: receipt table must remain plain", ErrIncompatibleExistingTable)
	}
	return nil
}

func sameSQLServerTargetPartition(actual *TargetSQLServerPartitionEvidence, expected *SQLServerPartitionEvidence) bool {
	if actual == nil || expected == nil {
		return actual == nil && expected == nil
	}
	return actual.Method == expected.Method && actual.Column == expected.Column &&
		actual.FunctionName == expected.FunctionName && actual.SchemeName == expected.SchemeName &&
		actual.Filegroup == expected.Filegroup && slices.Equal(actual.Partitions, expected.Partitions)
}

func sqlServerPartitionPreview(layout GroupSchemaLayout, action string) *SQLServerPartitionPreview {
	native := layout.SQLServer.NativeMonthly
	if native == nil {
		return nil
	}
	return &SQLServerPartitionPreview{
		Action: action, FunctionName: native.FunctionName, SchemeName: native.SchemeName,
		Filegroup: native.Filegroup, Partitions: slices.Clone(native.Partitions),
	}
}

func sqlServerPartitionUnavailable(reason string) error {
	return &SQLServerPartitionUnavailableError{Reasons: []string{reason}}
}

// validateSQLServerCatalogTable performs strict contract verification on an existing table.
func validateSQLServerCatalogTable(state TargetTableInspection, columns []GroupSchemaColumn, owner string, options SQLServerGroupSchemaOptions) error {
	evidence := state.SQLServer
	if evidence == nil || !evidence.MetadataVerified || state.ColumnTypes == nil || state.ColumnNullable == nil || state.ColumnPrimaryKey == nil {
		return ErrTargetInspectionUnconfirmed
	}
	if owner == managedEffectReceiptMarker && evidence.RepairOperationID != "" {
		return ErrIncompatibleExistingTable
	}
	if evidence.TableType != "USER_TABLE" || evidence.ObjectID <= 0 {
		return ErrIncompatibleExistingTable
	}
	native := state.SQLServer != nil && state.SQLServer.NativeMonthly != nil
	if native == evidence.Plain {
		// A partitioned table reports Plain=false; a plain table reports true.
		// A native partition projection must not claim a plain table and vice versa.
		return ErrIncompatibleExistingTable
	}
	if evidence.Database != options.Database || evidence.Schema != options.Schema {
		return ErrIncompatibleExistingTable
	}
	if evidence.ServerVersion != options.ServerVersion ||
		evidence.ServerEdition != options.ServerEdition ||
		evidence.CapabilityRevision != options.CapabilityRevision ||
		evidence.EngineEdition != options.EngineEdition ||
		evidence.CompatibilityLevel != options.CompatibilityLevel {
		return ErrIncompatibleExistingTable
	}
	if len(state.Columns) != len(columns) {
		return ErrIncompatibleExistingTable
	}

	var primary []string
	hasProvenance := false
	for i, column := range columns {
		if state.Columns[i] != column.Name {
			return ErrIncompatibleExistingTable
		}
		actualType, hasType := state.ColumnTypes[column.Name]
		nullable, hasNullable := state.ColumnNullable[column.Name]
		primaryKey, hasKey := state.ColumnPrimaryKey[column.Name]
		if !hasType || !hasNullable || !hasKey || nullable != column.Nullable || primaryKey != column.PrimaryKey {
			return ErrIncompatibleExistingTable
		}
		if !sqlserverCatalogTypeMatches(column.SQLType, actualType) {
			return ErrIncompatibleExistingTable
		}
		if column.PrimaryKey {
			primary = append(primary, column.Name)
		}
		if sqlserverCollatedSQLType(column.SQLType) {
			if evidence.ColumnCollations == nil || evidence.ColumnCollations[column.Name] != options.Collation {
				return ErrIncompatibleExistingTable
			}
		}
		def := evidence.ColumnDefaults[column.Name]
		if column.Name == owner {
			if def == nil || *def != managedOwnerDefault {
				return ErrIncompatibleExistingTable
			}
		} else {
			if def != nil {
				return ErrIncompatibleExistingTable
			}
		}
		if column.Name == sqlserverProvenanceColumn {
			hasProvenance = true
		}
	}
	if !slices.Equal(primary, evidence.PrimaryKey) {
		return ErrIncompatibleExistingTable
	}
	if hasProvenance && !evidence.ProvenanceJSONChecked {
		return ErrIncompatibleExistingTable
	}
	return nil
}

func validateSQLServerCatalogIndexes(state TargetTableInspection, expected []GroupSchemaIndex) error {
	if !state.IndexesInspected {
		return ErrTargetInspectionUnconfirmed
	}
	if len(state.Indexes) != len(expected) {
		return ErrIncompatibleExistingTable
	}
	seen := make(map[string]bool, len(state.Indexes))
	for _, actual := range state.Indexes {
		if seen[actual.Name] {
			return ErrIncompatibleExistingTable
		}
		seen[actual.Name] = true
		i := slices.IndexFunc(expected, func(idx GroupSchemaIndex) bool {
			return idx.Name == actual.Name
		})
		if i < 0 || actual.Unique || !slices.Equal(actual.Columns, expected[i].Columns) {
			return ErrIncompatibleExistingTable
		}
	}
	return nil
}

func checkSQLServerCatalogBinding(database string, layout GroupSchemaLayout, evidence *TargetSQLServerTableEvidence) error {
	if evidence == nil || !evidence.MetadataVerified {
		return ErrTargetInspectionUnconfirmed
	}
	options := layout.SQLServer
	if options == nil {
		return ErrInvalidGroupSchemaLayout
	}
	if database == "" || evidence.Database != database || evidence.Database != options.Database || evidence.Schema != options.Schema {
		return fmt.Errorf("%w: %w", ErrPreviewTokenStale, ErrSchemaTargetChanged)
	}
	if evidence.ServerVersion != options.ServerVersion ||
		evidence.ServerEdition != options.ServerEdition ||
		evidence.CapabilityRevision != options.CapabilityRevision ||
		evidence.EngineEdition != options.EngineEdition ||
		evidence.CompatibilityLevel != options.CompatibilityLevel {
		return fmt.Errorf("%w: %w", ErrPreviewTokenStale, ErrSchemaTargetChanged)
	}
	return nil
}

func sqlserverReceiptColumns() []GroupSchemaColumn {
	return []GroupSchemaColumn{
		{Name: sqlserverReceiptEffectColumn, SQLType: sqlserverHashSQLType, PrimaryKey: true},
		{Name: sqlserverReceiptDigestColumn, SQLType: sqlserverHashSQLType},
		{Name: sqlserverReceiptTimeColumn, SQLType: sqlserverTimeSQLType},
		{Name: managedEffectReceiptMarker, SQLType: sqlserverOwnerSQLType},
	}
}

func sqlserverCatalogTypeMatches(expected, actual string) bool {
	expected = strings.ToUpper(strings.TrimSpace(expected))
	actual = strings.ToUpper(strings.TrimSpace(actual))
	if expected == actual {
		return true
	}
	if (expected == "DECIMAL(38,18)" || expected == "NUMERIC(38,18)") && (actual == "DECIMAL(38,18)" || actual == "NUMERIC(38,18)") {
		return true
	}
	if (expected == "FLOAT(53)" || expected == "FLOAT") && (actual == "FLOAT(53)" || actual == "FLOAT") {
		return true
	}
	return false
}

func sqlserverCatalogDigest(state TargetTableInspection) (string, error) {
	payload, err := json.Marshal(state)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

func sqlserverPendingSchemaStatements(layout GroupSchemaLayout, tables []SchemaPreviewTable, generated []string) ([]string, error) {
	var statements []string
	if len(generated) == 0 {
		return nil, ErrInvalidGroupSchemaLayout
	}
	partitionPrefix := 0
	if layout.SQLServer != nil && layout.SQLServer.NativeMonthly != nil {
		partitionPrefix = 4
	}
	if partitionPrefix > len(generated)-1 {
		return nil, ErrInvalidGroupSchemaLayout
	}
	partitionStatements := generated[:partitionPrefix]
	groupStatements := generated[partitionPrefix : len(generated)-1]
	receiptStmt := generated[len(generated)-1]
	for _, table := range tables {
		switch {
		case table.Name == layout.TableName && table.Action == SchemaTableActionCreate:
			if partitionPrefix > 0 {
				if table.SQLServerPartition == nil || table.SQLServerPartition.Action == SchemaTableActionCreate {
					statements = append(statements, partitionStatements...)
				} else if table.SQLServerPartition.Action != SchemaTableActionUnchanged {
					return nil, fmt.Errorf("%w: sqlserver partition action %q", ErrInvalidGroupSchemaLayout, table.SQLServerPartition.Action)
				}
			}
			statements = append(statements, groupStatements...)
		case table.Name == managedEffectReceiptTable && table.Action == SchemaTableActionCreate:
			statements = append(statements, receiptStmt)
		case table.Action == SchemaTableActionUnchanged:
			// Unchanged table produces no statements
		default:
			return nil, fmt.Errorf("%w: sqlserver plain layout does not support action %q", ErrInvalidGroupSchemaLayout, table.Action)
		}
	}
	return statements, nil
}
