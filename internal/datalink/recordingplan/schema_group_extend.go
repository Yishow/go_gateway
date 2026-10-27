package recordingplan

import (
	"fmt"
	"regexp"
	"strings"
)

var managedValueColumnPattern = regexp.MustCompile(`^v_[0-9a-f]{24}$`)

// inspectGroupExtension accepts only additive, nullable value columns. The
// current owner marker and all structural columns must already match exactly.
// Removed members' nullable managed columns remain physically present.
func inspectGroupExtension(state TargetTableInspection, layout GroupSchemaLayout) ([]GroupSchemaColumn, error) {
	if state.ColumnTypes == nil || state.ColumnNullable == nil || state.ColumnPrimaryKey == nil {
		return nil, ErrTargetInspectionUnconfirmed
	}
	names := make(map[string]bool, len(state.Columns))
	for _, name := range state.Columns {
		names[strings.ToLower(strings.TrimSpace(name))] = true
	}
	expected := make(map[string]bool, len(layout.Columns))
	existing := make([]GroupSchemaColumn, 0, len(state.Columns))
	var missing []GroupSchemaColumn
	for _, column := range layout.Columns {
		name := strings.ToLower(column.Name)
		expected[name] = true
		if names[name] {
			existing = append(existing, column)
			continue
		}
		if !managedValueColumnPattern.MatchString(column.Name) || !column.Nullable || column.PrimaryKey {
			return nil, ErrIncompatibleExistingTable
		}
		missing = append(missing, column)
	}
	for _, name := range state.Columns {
		if expected[strings.ToLower(strings.TrimSpace(name))] {
			continue
		}
		sqlType, ok := lookupGroupType(state.ColumnTypes, name)
		if !ok || !managedValueColumnPattern.MatchString(name) {
			return nil, ErrIncompatibleExistingTable
		}
		if _, allowed := allowedGroupSQLTypes[normalizeGroupSQLType(sqlType)]; !allowed && observedGroupSQLType(sqlType) != managedNumericSQLType {
			return nil, ErrIncompatibleExistingTable
		}
		existing = append(existing, GroupSchemaColumn{Name: name, SQLType: sqlType, Nullable: true})
	}
	// Exact validation also rejects duplicate names, incomplete metadata, changed
	// types, non-null historical columns, additional keys, and missing ownership.
	if err := validateGroupInspection(state, existing); err != nil {
		return nil, err
	}
	return missing, nil
}

func groupSchemaStatements(dialect string, layout GroupSchemaLayout, tables []SchemaPreviewTable) ([]string, error) {
	create, err := GenerateGroupSchemaDDL(dialect, layout)
	if err != nil {
		return nil, err
	}
	if layout.MySQL != nil {
		return mysqlPendingSchemaStatements(layout, tables, create)
	}
	if layout.SQLServer != nil {
		return sqlserverPendingSchemaStatements(layout, tables, create)
	}
	// create is the group table, then its indexes, then the shared receipt table.
	receipt := create[len(create)-1]
	var statements []string
	for _, table := range tables {
		switch {
		case table.Name == layout.TableName && table.Action == SchemaTableActionCreate:
			statements = append(statements, create[:len(create)-1]...)
		case table.Name == layout.TableName && table.Action == SchemaTableActionAddColumns:
			for _, name := range table.Columns {
				found := false
				for _, column := range layout.Columns {
					if column.Name != name {
						continue
					}
					if !managedValueColumnPattern.MatchString(name) || !column.Nullable || column.PrimaryKey {
						return nil, ErrIncompatibleExistingTable
					}
					statements = append(statements, "ALTER TABLE "+quoteGroupIdentifier(layout.TableName)+" ADD COLUMN "+quoteGroupIdentifier(name)+" "+column.SQLType+";")
					found = true
					break
				}
				if !found {
					return nil, fmt.Errorf("%w: unknown managed value column", ErrInvalidGroupSchemaLayout)
				}
			}
		case table.Name == managedEffectReceiptTable && table.Action == SchemaTableActionCreate:
			statements = append(statements, receipt)
		}
	}
	return statements, nil
}
