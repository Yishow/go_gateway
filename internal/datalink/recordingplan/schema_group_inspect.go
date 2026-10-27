package recordingplan

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

// TargetTableInspection is the actual state of one managed table.
type TargetTableInspection struct {
	Status           string
	Columns          []string
	ColumnTypes      map[string]string
	ColumnNullable   map[string]bool
	ColumnPrimaryKey map[string]bool
	// Indexes lists the non-primary-key indexes read from the catalog. It is
	// required only when the expected layout declares indexes.
	Indexes []TargetIndexInspection
	// IndexesInspected is true when Indexes was read from the catalog; an empty
	// list is then a real answer rather than missing evidence.
	IndexesInspected bool
	// Timescale is the read-only capability and hypertable evidence of the
	// group table. It is read only for a layout that asks for a partition.
	Timescale *TargetTimescaleEvidence
	// MySQL includes precision, engine, collation and ownership catalog facts.
	MySQL *TargetMySQLTableEvidence
	// SQLServer includes SQL Server engine edition, compatibility, collation and metadata evidence.
	SQLServer *TargetSQLServerTableEvidence `json:",omitempty"`
}

// TargetTimescaleEvidence is what the destination's catalog says about
// TimescaleDB for one table: whether preparing a hypertable is possible, and
// what the table already is.
type TargetTimescaleEvidence struct {
	Available       bool
	Reasons         []string
	ExtensionSchema string
	Version         string
	// Hypertable is nil when the table is not a hypertable.
	Hypertable *TargetHypertable
}

// TargetHypertable is the catalog evidence of one hypertable.
type TargetHypertable struct {
	TimeColumn         string
	ChunkInterval      time.Duration
	Dimensions         int
	CompressionEnabled bool
	Jobs               int
}

// TargetIndexInspection is one catalog index with its ordered key columns.
type TargetIndexInspection struct {
	Name    string
	Columns []string
	Unique  bool
}

func groupSchemaLayoutsEqual(left, right GroupSchemaLayout) bool {
	left, leftErr := ValidateGroupSchemaLayout(left)
	right, rightErr := ValidateGroupSchemaLayout(right)
	if leftErr != nil || rightErr != nil {
		return false
	}
	if !sameMySQLGroupSchemaOptions(left.MySQL, right.MySQL) ||
		!sameSQLServerGroupSchemaOptions(left.SQLServer, right.SQLServer) ||
		!samePartition(left.Partition, right.Partition) ||
		left.TableName != right.TableName || left.OwnerColumn != right.OwnerColumn || len(left.Columns) != len(right.Columns) ||
		!slices.EqualFunc(left.Indexes, right.Indexes, func(a, b GroupSchemaIndex) bool {
			return a.Name == b.Name && slices.Equal(a.Columns, b.Columns)
		}) {
		return false
	}
	for i := range left.Columns {
		if left.Columns[i] != right.Columns[i] {
			return false
		}
	}
	return true
}

func sameSQLServerGroupSchemaOptions(left, right *SQLServerGroupSchemaOptions) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return left.Database == right.Database && left.Schema == right.Schema && left.Collation == right.Collation &&
		left.ServerVersion == right.ServerVersion && left.ServerEdition == right.ServerEdition &&
		left.CapabilityRevision == right.CapabilityRevision && left.EngineEdition == right.EngineEdition &&
		left.CompatibilityLevel == right.CompatibilityLevel && sameSQLServerPartitionEvidence(left.NativeMonthly, right.NativeMonthly)
}

// samePartition compares what is partitioned and how, not where the extension
// happens to be installed: that is probed evidence, checked at confirmation.
func samePartition(left, right *GroupSchemaPartition) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return left.Kind == right.Kind && left.Column == right.Column && left.ChunkInterval == right.ChunkInterval
}

// prepareGroupSchemaPreview returns the statements and tables of a preview and
// the effective layout: a timescale layout gets the extension schema the
// read-only probe found, so the confirmed statements and digest carry it.
func prepareGroupSchemaPreview(ctx context.Context, scope SchemaPreviewScope, inspect TargetInspector) ([]string, []SchemaPreviewTable, *GroupSchemaLayout, error) {
	layout := *scope.GroupLayout
	if layout.MySQL != nil {
		found, err := inspectMySQLGroupTables(ctx, layout, scope.Database, inspect)
		if err != nil {
			return nil, nil, nil, err
		}
		statements, err := groupSchemaStatements(scope.Dialect, layout, found.tables)
		return statements, found.tables, &layout, err
	}
	if layout.SQLServer != nil {
		found, err := inspectSQLServerGroupTables(ctx, layout, scope.Database, inspect)
		if err != nil {
			if layout.SQLServer.NativeMonthly != nil && errors.Is(err, ErrTargetPermissionDenied) {
				return nil, nil, nil, sqlServerPartitionUnavailable("partition metadata permission denied")
			}
			return nil, nil, nil, err
		}
		statements, err := groupSchemaStatements(scope.Dialect, layout, found.tables)
		return statements, found.tables, &layout, err
	}
	inspected, err := inspectGroupTables(ctx, scope.Dialect, layout, inspect)
	if err != nil {
		return nil, nil, nil, err
	}
	if layout.Partition != nil && layout.Partition.ExtensionSchema == "" && inspected.extensionSchema != "" {
		partition := *layout.Partition
		partition.ExtensionSchema = inspected.extensionSchema
		layout.Partition = &partition
		inspected.tables[0].Partition = &partition
	}
	statements, err := groupSchemaStatements(scope.Dialect, layout, inspected.tables)
	return statements, inspected.tables, &layout, err
}

type groupInspection struct {
	pendingGroup, pendingReceipt bool
	tables                       []SchemaPreviewTable
	// extensionSchema is where the probed TimescaleDB functions live.
	extensionSchema string
}

func inspectGroupManagedTables(ctx context.Context, dialect string, layout GroupSchemaLayout, inspect TargetInspector) (pendingGroup, pendingReceipt bool, tables []SchemaPreviewTable, err error) {
	inspected, err := inspectGroupTables(ctx, dialect, layout, inspect)
	return inspected.pendingGroup, inspected.pendingReceipt, inspected.tables, err
}

// checkTimescaleEvidence applies the timescale rules to the group table. A
// table that has to be created needs a capable destination whose extension
// schema is still the confirmed one; a table that exists must already be the
// owned hypertable of the layout. Nothing existing is converted, whether or not
// it holds rows: a plain table of the same name is simply incompatible.
func checkTimescaleEvidence(layout GroupSchemaLayout, state TargetTableInspection) (extensionSchema string, err error) {
	evidence := state.Timescale
	if evidence == nil {
		return "", fmt.Errorf("%w: no timescale evidence for %s", ErrTargetInspectionUnconfirmed, layout.TableName)
	}
	switch state.Status {
	case inspectionMissing:
		if !evidence.Available {
			reasons := evidence.Reasons
			if len(reasons) == 0 {
				reasons = []string{"capability-unconfirmed"}
			}
			return "", &TimescaleUnavailableError{Reasons: slices.Clone(reasons)}
		}
		if confirmed := layout.Partition.ExtensionSchema; confirmed != "" && confirmed != evidence.ExtensionSchema {
			return "", fmt.Errorf("%w: %w", ErrPreviewTokenStale, ErrSchemaTargetChanged)
		}
		return evidence.ExtensionSchema, nil
	case inspectionExists:
		if !hypertableMatches(layout.Partition, evidence.Hypertable) {
			return "", ErrIncompatibleExistingTable
		}
		return evidence.ExtensionSchema, nil
	}
	return "", nil
}

// hypertableMatches requires the catalog to show exactly the confirmed
// partitioning: one time dimension on the partition column at the confirmed
// interval, with no compression and no background job bound to it.
func hypertableMatches(partition *GroupSchemaPartition, hypertable *TargetHypertable) bool {
	return hypertable != nil && hypertable.TimeColumn == partition.Column && hypertable.Dimensions == 1 &&
		hypertable.ChunkInterval == 7*24*time.Hour && !hypertable.CompressionEnabled && hypertable.Jobs == 0
}

func inspectGroupTables(ctx context.Context, dialect string, layout GroupSchemaLayout, inspect TargetInspector) (groupInspection, error) {
	if layout.SQLServer != nil {
		if normalizePreviewDialect(dialect) != dialectSQLServer {
			return groupInspection{}, fmt.Errorf("%w: a sqlserver layout needs the sqlserver dialect", ErrUnsupportedSchemaDialect)
		}
		if layout.SQLServer.Database == "" {
			return groupInspection{}, fmt.Errorf("%w: sqlserver database is required in frozen layout", ErrInvalidGroupSchemaLayout)
		}
		return inspectSQLServerGroupTables(ctx, layout, layout.SQLServer.Database, inspect)
	}
	var found groupInspection
	groupState, inspectErr := inspect(ctx, layout.TableName)
	if inspectErr != nil {
		return found, fmt.Errorf("%w: %s: %w", ErrTargetInspectionUnconfirmed, layout.TableName, inspectErr)
	}
	if layout.Partition != nil {
		if found.extensionSchema, inspectErr = checkTimescaleEvidence(layout, groupState); inspectErr != nil {
			return groupInspection{}, inspectErr
		}
	}
	var added []GroupSchemaColumn
	if groupState.Status == inspectionExists {
		added, inspectErr = inspectGroupExtension(groupState, layout)
		found.pendingGroup = len(added) > 0
		if inspectErr == nil {
			inspectErr = validateGroupIndexes(groupState, layout)
		}
	} else {
		found.pendingGroup, inspectErr = inspectGroupTable(layout.TableName, groupState, layout.Columns)
	}
	if inspectErr != nil {
		return groupInspection{}, inspectErr
	}

	receiptState, inspectErr := inspect(ctx, managedEffectReceiptTable)
	if inspectErr != nil {
		return groupInspection{}, fmt.Errorf("%w: %s: %w", ErrTargetInspectionUnconfirmed, managedEffectReceiptTable, inspectErr)
	}
	found.pendingReceipt, inspectErr = inspectGroupTable(managedEffectReceiptTable, receiptState, managedEffectReceiptColumns(dialect))
	if inspectErr != nil {
		return groupInspection{}, inspectErr
	}

	found.tables = []SchemaPreviewTable{
		{Name: layout.TableName, Action: SchemaTableActionUnchanged},
		{Name: managedEffectReceiptTable, Action: SchemaTableActionUnchanged},
	}
	if len(added) > 0 {
		found.tables[0] = SchemaPreviewTable{Name: layout.TableName, Action: SchemaTableActionAddColumns, Columns: groupColumnNames(added)}
	} else if found.pendingGroup {
		found.tables[0] = SchemaPreviewTable{Name: layout.TableName, Action: SchemaTableActionCreate, Columns: groupColumnNames(layout.Columns)}
	}
	if found.pendingReceipt {
		found.tables[1] = SchemaPreviewTable{Name: managedEffectReceiptTable, Action: SchemaTableActionCreate, Columns: groupColumnNames(managedEffectReceiptColumns(dialect))}
	}
	if layout.Partition != nil {
		partition := *layout.Partition
		found.tables[0].Partition = &partition
	}
	return found, nil
}

func inspectGroupTable(table string, state TargetTableInspection, expected []GroupSchemaColumn) (bool, error) {
	switch state.Status {
	case inspectionMissing:
		return true, nil
	case inspectionForbidden:
		return false, fmt.Errorf("%w: %w: %s", ErrTargetInspectionUnconfirmed, ErrTargetPermissionDenied, table)
	case inspectionExists:
		if err := validateGroupInspection(state, expected); err != nil {
			return false, fmt.Errorf("%s: %w", table, err)
		}
		return false, nil
	default:
		return false, fmt.Errorf("%w: %s is %s", ErrTargetInspectionUnconfirmed, table, state.Status)
	}
}

func validateGroupInspection(state TargetTableInspection, expected []GroupSchemaColumn) error {
	if state.ColumnTypes == nil || state.ColumnNullable == nil || state.ColumnPrimaryKey == nil {
		return ErrTargetInspectionUnconfirmed
	}
	actualNames := make(map[string]struct{}, len(state.Columns))
	for _, name := range state.Columns {
		key := strings.ToLower(strings.TrimSpace(name))
		if key == "" {
			return ErrIncompatibleExistingTable
		}
		if _, exists := actualNames[key]; exists {
			return ErrIncompatibleExistingTable
		}
		actualNames[key] = struct{}{}
	}
	if len(actualNames) != len(expected) || len(state.ColumnTypes) != len(expected) || len(state.ColumnNullable) != len(expected) || len(state.ColumnPrimaryKey) != len(expected) {
		return ErrIncompatibleExistingTable
	}
	for _, column := range expected {
		key := strings.ToLower(column.Name)
		if _, exists := actualNames[key]; !exists {
			return ErrIncompatibleExistingTable
		}
		actualType, ok := lookupGroupType(state.ColumnTypes, column.Name)
		if !ok || strings.TrimSpace(actualType) == "" {
			return ErrTargetInspectionUnconfirmed
		}
		if !groupSQLTypesCompatible(column.SQLType, actualType) {
			return ErrIncompatibleExistingTable
		}
		actualNullable, ok := lookupGroupBool(state.ColumnNullable, column.Name)
		if !ok || actualNullable != column.Nullable {
			return ErrIncompatibleExistingTable
		}
		actualPrimaryKey, ok := lookupGroupBool(state.ColumnPrimaryKey, column.Name)
		if !ok || actualPrimaryKey != column.PrimaryKey {
			return ErrIncompatibleExistingTable
		}
	}
	return nil
}

// ValidateManagedGroupTableInspection verifies the canonical layout and owner,
// permitting retained nullable managed value columns from prior memberships.
func ValidateManagedGroupTableInspection(layout GroupSchemaLayout, state TargetTableInspection) error {
	normalized, err := ValidateGroupSchemaLayout(layout)
	if err != nil {
		return err
	}
	missing, err := inspectGroupExtension(state, normalized)
	if err != nil {
		return err
	}
	if len(missing) > 0 {
		return ErrIncompatibleExistingTable
	}
	if normalized.Partition != nil {
		// A timescale table is confirmed only while its catalog still shows the
		// owned hypertable; a plain table of that name never is.
		if state.Timescale == nil {
			return ErrTargetInspectionUnconfirmed
		}
		if !hypertableMatches(normalized.Partition, state.Timescale.Hypertable) {
			return ErrIncompatibleExistingTable
		}
	}
	return validateGroupIndexes(state, normalized)
}

// validateGroupIndexes requires every declared index to exist with the same
// ordered, non-unique key. A layout with indexes cannot be confirmed from an
// inspection that did not read the catalog's indexes.
func validateGroupIndexes(state TargetTableInspection, layout GroupSchemaLayout) error {
	if len(layout.Indexes) == 0 {
		return nil
	}
	if !state.IndexesInspected {
		return ErrTargetInspectionUnconfirmed
	}
	for _, expected := range layout.Indexes {
		matched := slices.ContainsFunc(state.Indexes, func(actual TargetIndexInspection) bool {
			return strings.EqualFold(actual.Name, expected.Name) && !actual.Unique &&
				slices.EqualFunc(actual.Columns, expected.Columns, strings.EqualFold)
		})
		if !matched {
			return ErrIncompatibleExistingTable
		}
	}
	return nil
}

// ValidateManagedReceiptTableInspection verifies the complete shared receipt
// table contract, including its ownership marker.
func ValidateManagedReceiptTableInspection(dialect string, state TargetTableInspection) error {
	return validateGroupInspection(state, managedEffectReceiptColumns(dialect))
}

func lookupGroupType(values map[string]string, name string) (string, bool) {
	for key, value := range values {
		if strings.EqualFold(strings.TrimSpace(key), name) {
			return value, true
		}
	}
	return "", false
}

func lookupGroupBool(values map[string]bool, name string) (matched, found bool) {
	for key, value := range values {
		if strings.EqualFold(strings.TrimSpace(key), name) {
			return value, true
		}
	}
	return false, false
}

func groupColumnNames(columns []GroupSchemaColumn) []string {
	names := make([]string, len(columns))
	for i, column := range columns {
		names[i] = column.Name
	}
	return names
}
