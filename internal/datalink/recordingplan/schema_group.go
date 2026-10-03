package recordingplan

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	managedEffectReceiptTable  = "gw_effect_receipts"
	managedEffectReceiptMarker = "_gw_effect_receipts_v1"
	managedOwnerDefault        = "gateway"
	managedTextSQLType         = "TEXT"
)

// ManagedEffectReceiptTable is the shared receipt table created by a
// canonical managed group schema preview.
const ManagedEffectReceiptTable = managedEffectReceiptTable

var (
	ErrInvalidGroupSchemaLayout = errors.New("invalid managed group schema layout")
	groupIdentifierPattern      = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	groupOwnerColumnPattern     = regexp.MustCompile(`^_gw_owner_[0-9a-f]{24}$`)
)

var allowedGroupSQLTypes = map[string]struct{}{
	managedTextSQLType: {},
	"INTEGER":          {},
	"REAL":             {},
	"BOOLEAN":          {},
	"BIGINT":           {},
	"DOUBLE PRECISION": {},
	"NUMERIC(20,0)":    {},
	"TIMESTAMPTZ":      {},
}

// ValidateGroupSchemaLayout normalizes and validates the server-provided
// canonical group layout before it can participate in SQL generation.
func ValidateGroupSchemaLayout(layout GroupSchemaLayout) (GroupSchemaLayout, error) {
	layout.TableName = strings.TrimSpace(layout.TableName)
	layout.OwnerColumn = strings.TrimSpace(layout.OwnerColumn)
	if !validGroupIdentifier(layout.TableName) || strings.EqualFold(layout.TableName, managedEffectReceiptTable) {
		return GroupSchemaLayout{}, fmt.Errorf("%w: table name", ErrInvalidGroupSchemaLayout)
	}
	if !validGroupIdentifier(layout.OwnerColumn) || !groupOwnerColumnPattern.MatchString(layout.OwnerColumn) {
		return GroupSchemaLayout{}, fmt.Errorf("%w: owner column", ErrInvalidGroupSchemaLayout)
	}
	if len(layout.Columns) == 0 {
		return GroupSchemaLayout{}, fmt.Errorf("%w: no columns", ErrInvalidGroupSchemaLayout)
	}

	columns := make([]GroupSchemaColumn, len(layout.Columns))
	seen := make(map[string]struct{}, len(layout.Columns))
	ownerFound := false
	for i, column := range layout.Columns {
		column.Name = strings.TrimSpace(column.Name)
		column.SQLType = normalizeGroupSQLType(column.SQLType)
		if !validGroupIdentifier(column.Name) {
			return GroupSchemaLayout{}, fmt.Errorf("%w: column name", ErrInvalidGroupSchemaLayout)
		}
		key := strings.ToLower(column.Name)
		if _, exists := seen[key]; exists {
			return GroupSchemaLayout{}, fmt.Errorf("%w: duplicate column %s", ErrInvalidGroupSchemaLayout, column.Name)
		}
		seen[key] = struct{}{}
		if _, ok := allowedGroupSQLTypes[column.SQLType]; !ok {
			return GroupSchemaLayout{}, fmt.Errorf("%w: unsupported SQL type %s", ErrInvalidGroupSchemaLayout, column.SQLType)
		}
		if column.PrimaryKey && column.Nullable {
			return GroupSchemaLayout{}, fmt.Errorf("%w: primary key cannot be nullable", ErrInvalidGroupSchemaLayout)
		}
		if column.Name == layout.OwnerColumn {
			ownerFound = true
			if column.SQLType != managedTextSQLType || column.Nullable || column.PrimaryKey {
				return GroupSchemaLayout{}, fmt.Errorf("%w: owner column must be non-null TEXT", ErrInvalidGroupSchemaLayout)
			}
		}
		columns[i] = column
	}
	if !ownerFound {
		return GroupSchemaLayout{}, fmt.Errorf("%w: owner column is not in columns", ErrInvalidGroupSchemaLayout)
	}
	layout.Columns = columns
	return layout, nil
}

func validGroupIdentifier(value string) bool {
	return len(value) <= 63 && groupIdentifierPattern.MatchString(value)
}

func normalizeGroupSQLType(raw string) string {
	value := strings.ToUpper(strings.Join(strings.Fields(strings.TrimSpace(raw)), " "))
	value = strings.ReplaceAll(value, " (", "(")
	value = strings.ReplaceAll(value, "( ", "(")
	value = strings.ReplaceAll(value, " ,", ",")
	value = strings.ReplaceAll(value, ", ", ",")
	value = strings.ReplaceAll(value, " )", ")")
	return value
}

func observedGroupSQLType(raw string) string {
	value := normalizeGroupSQLType(raw)
	switch value {
	case "BOOL":
		return "BOOLEAN"
	case "INT", "INT4":
		return "INTEGER"
	case "INT8":
		return "BIGINT"
	case "FLOAT8", "DOUBLE":
		return "DOUBLE PRECISION"
	case "DECIMAL":
		return "NUMERIC"
	case "TIMESTAMP WITH TIME ZONE":
		return "TIMESTAMPTZ"
	default:
		return value
	}
}

func groupSQLTypesCompatible(expected, actual string) bool {
	expected = normalizeGroupSQLType(expected)
	actual = observedGroupSQLType(actual)
	if expected == "NUMERIC(20,0)" && actual == "NUMERIC" {
		// Keep PostgreSQL's unconstrained NUMERIC equivalent for the exact-value
		// codec; bounded precision and scale are preserved by the inspector.
		return true
	}
	return expected == actual
}

func quoteGroupIdentifier(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}

// GenerateGroupSchemaDDL emits only server-generated CREATE statements. It
// intentionally omits IF NOT EXISTS so a same-name unknown table cannot be
// silently accepted after the preview race.
func GenerateGroupSchemaDDL(dialect string, layout GroupSchemaLayout) ([]string, error) {
	layout, err := ValidateGroupSchemaLayout(layout)
	if err != nil {
		return nil, err
	}
	switch normalizePreviewDialect(dialect) {
	case dialectSQLite, dialectPostgres:
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedSchemaDialect, dialect)
	}

	columnDefinitions := make([]string, 0, len(layout.Columns))
	for _, column := range layout.Columns {
		definition := quoteGroupIdentifier(column.Name) + " " + column.SQLType
		if !column.Nullable {
			definition += " NOT NULL"
		}
		if column.PrimaryKey {
			definition += " PRIMARY KEY"
		}
		if column.Name == layout.OwnerColumn {
			definition += " DEFAULT '" + managedOwnerDefault + "'"
		}
		columnDefinitions = append(columnDefinitions, definition)
	}
	receiptDefinitions := []string{
		quoteGroupIdentifier("effect_key") + " " + managedTextSQLType + " PRIMARY KEY",
		quoteGroupIdentifier("payload_digest") + " " + managedTextSQLType + " NOT NULL",
		quoteGroupIdentifier("committed_at") + " " + managedTextSQLType + " NOT NULL",
		quoteGroupIdentifier(managedEffectReceiptMarker) + " " + managedTextSQLType + " NOT NULL DEFAULT '" + managedOwnerDefault + "'",
	}
	return []string{
		"CREATE TABLE " + quoteGroupIdentifier(layout.TableName) + " (" + strings.Join(columnDefinitions, ", ") + ");",
		"CREATE TABLE " + quoteGroupIdentifier(managedEffectReceiptTable) + " (" + strings.Join(receiptDefinitions, ", ") + ");",
	}, nil
}

func managedEffectReceiptColumns(dialect string) []GroupSchemaColumn {
	effectKeyNullable := normalizePreviewDialect(dialect) != dialectPostgres
	return []GroupSchemaColumn{
		{Name: "effect_key", SQLType: managedTextSQLType, Nullable: effectKeyNullable, PrimaryKey: true},
		{Name: "payload_digest", SQLType: managedTextSQLType, Nullable: false},
		{Name: "committed_at", SQLType: managedTextSQLType, Nullable: false},
		{Name: managedEffectReceiptMarker, SQLType: managedTextSQLType, Nullable: false},
	}
}

func prepareGroupSchemaPreview(ctx context.Context, scope SchemaPreviewScope, inspect TargetInspector) ([]string, []SchemaPreviewTable, error) {
	statements, err := GenerateGroupSchemaDDL(scope.Dialect, *scope.GroupLayout)
	if err != nil {
		return nil, nil, err
	}
	groupPending, receiptPending, tables, err := inspectGroupManagedTables(ctx, scope.Dialect, *scope.GroupLayout, inspect)
	if err != nil {
		return nil, nil, err
	}
	kept := make([]string, 0, 2)
	if groupPending {
		kept = append(kept, statements[0])
	}
	if receiptPending {
		kept = append(kept, statements[1])
	}
	return kept, tables, nil
}

func inspectGroupManagedTables(ctx context.Context, dialect string, layout GroupSchemaLayout, inspect TargetInspector) (pendingGroup, pendingReceipt bool, tables []SchemaPreviewTable, err error) {
	groupState, inspectErr := inspect(ctx, layout.TableName)
	if inspectErr != nil {
		return false, false, nil, fmt.Errorf("%w: %s: %w", ErrTargetInspectionUnconfirmed, layout.TableName, inspectErr)
	}
	pendingGroup, inspectErr = inspectGroupTable(layout.TableName, groupState, layout.Columns)
	if inspectErr != nil {
		return false, false, nil, inspectErr
	}

	receiptState, inspectErr := inspect(ctx, managedEffectReceiptTable)
	if inspectErr != nil {
		return false, false, nil, fmt.Errorf("%w: %s: %w", ErrTargetInspectionUnconfirmed, managedEffectReceiptTable, inspectErr)
	}
	pendingReceipt, inspectErr = inspectGroupTable(managedEffectReceiptTable, receiptState, managedEffectReceiptColumns(dialect))
	if inspectErr != nil {
		return false, false, nil, inspectErr
	}

	tables = []SchemaPreviewTable{
		{Name: layout.TableName, Action: SchemaTableActionUnchanged},
		{Name: managedEffectReceiptTable, Action: SchemaTableActionUnchanged},
	}
	if pendingGroup {
		tables[0] = SchemaPreviewTable{Name: layout.TableName, Action: SchemaTableActionCreate, Columns: groupColumnNames(layout.Columns)}
	}
	if pendingReceipt {
		tables[1] = SchemaPreviewTable{Name: managedEffectReceiptTable, Action: SchemaTableActionCreate, Columns: groupColumnNames(managedEffectReceiptColumns(dialect))}
	}
	return pendingGroup, pendingReceipt, tables, nil
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

// ValidateManagedGroupTableInspection verifies a live table against the exact
// canonical layout, including the managed owner marker.
func ValidateManagedGroupTableInspection(layout GroupSchemaLayout, state TargetTableInspection) error {
	normalized, err := ValidateGroupSchemaLayout(layout)
	if err != nil {
		return err
	}
	return validateGroupInspection(state, normalized.Columns)
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
