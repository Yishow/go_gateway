package recordingplan

import (
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
	managedNumericSQLType      = "NUMERIC"

	sqlTypeInteger         = "INTEGER"
	sqlTypeBoolean         = "BOOLEAN"
	sqlTypeBigint          = "BIGINT"
	sqlTypeDoublePrecision = "DOUBLE PRECISION"
	sqlTypeTimestamptz     = "TIMESTAMPTZ"
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
	managedTextSQLType:     {},
	sqlTypeInteger:         {},
	"REAL":                 {},
	sqlTypeBoolean:         {},
	sqlTypeBigint:          {},
	sqlTypeDoublePrecision: {},
	"NUMERIC(20,0)":        {},
	sqlTypeTimestamptz:     {},
	// The long layout keeps decimals and deltas unbounded and provenance as JSON.
	managedNumericSQLType: {},
	"JSONB":               {},
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
	if layout.MySQL != nil && layout.SQLServer != nil {
		return GroupSchemaLayout{}, fmt.Errorf("%w: layout cannot specify both mysql and sqlserver options", ErrInvalidGroupSchemaLayout)
	}
	if layout.MySQL != nil {
		options, err := validateMySQLGroupSchemaOptions(*layout.MySQL)
		if err != nil {
			return GroupSchemaLayout{}, err
		}
		layout.MySQL = &options
		if layout.Partition != nil {
			return GroupSchemaLayout{}, fmt.Errorf("%w: mysql layouts require native_monthly evidence, not generic partition metadata", ErrInvalidGroupSchemaLayout)
		}
	}
	if layout.SQLServer != nil {
		options, err := validateSQLServerGroupSchemaOptions(*layout.SQLServer)
		if err != nil {
			return GroupSchemaLayout{}, err
		}
		if err := validateSQLServerPartitionEvidence(options.NativeMonthly, layout.TableName); err != nil {
			return GroupSchemaLayout{}, err
		}
		layout.SQLServer = &options
		if layout.Partition != nil {
			return GroupSchemaLayout{}, fmt.Errorf("%w: sqlserver plain layout does not support generic partition metadata", ErrInvalidGroupSchemaLayout)
		}
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
		allowedTypes := allowedGroupSQLTypes
		if layout.MySQL != nil {
			allowedTypes = allowedMySQLGroupSQLTypes
		} else if layout.SQLServer != nil {
			allowedTypes = allowedSQLServerGroupSQLTypes
		}
		if _, ok := allowedTypes[column.SQLType]; !ok {
			return GroupSchemaLayout{}, fmt.Errorf("%w: unsupported SQL type %s", ErrInvalidGroupSchemaLayout, column.SQLType)
		}
		if column.PrimaryKey && column.Nullable {
			return GroupSchemaLayout{}, fmt.Errorf("%w: primary key cannot be nullable", ErrInvalidGroupSchemaLayout)
		}
		if column.Name == layout.OwnerColumn {
			ownerFound = true
			ownerType := managedTextSQLType
			if layout.MySQL != nil {
				ownerType = mysqlOwnerSQLType
			} else if layout.SQLServer != nil {
				ownerType = sqlserverOwnerSQLType
			}
			if column.SQLType != ownerType || column.Nullable || column.PrimaryKey {
				return GroupSchemaLayout{}, fmt.Errorf("%w: owner column must be non-null %s", ErrInvalidGroupSchemaLayout, ownerType)
			}
		}
		columns[i] = column
	}
	if !ownerFound {
		return GroupSchemaLayout{}, fmt.Errorf("%w: owner column is not in columns", ErrInvalidGroupSchemaLayout)
	}
	layout.Columns = columns
	if layout.MySQL != nil && layout.MySQL.NativeMonthly != nil {
		if err := validateMySQLNativeMonthlyGroupColumns(layout.Columns); err != nil {
			return GroupSchemaLayout{}, err
		}
	}
	if layout.SQLServer != nil && layout.SQLServer.NativeMonthly != nil {
		if err := validateSQLServerNativeMonthlyGroupKey(layout.Columns); err != nil {
			return GroupSchemaLayout{}, err
		}
	}
	indexes, err := validateGroupSchemaIndexes(layout, seen)
	if err != nil {
		return GroupSchemaLayout{}, err
	}
	layout.Indexes = indexes
	if layout.Partition != nil {
		if err := validateGroupPartition(layout); err != nil {
			return GroupSchemaLayout{}, err
		}
	}
	return layout, nil
}

// validateGroupPartition accepts exactly the documented timescale partition: the
// seven-day interval on the non-null bucket_start primary-key column, so the
// unique key always includes the dimension TimescaleDB requires. Anything else
// is a layout nobody confirmed.
func validateGroupPartition(layout GroupSchemaLayout) error {
	partition := layout.Partition
	invalid := func(reason string) error { return fmt.Errorf("%w: partition %s", ErrInvalidGroupSchemaLayout, reason) }
	if partition.Kind != PartitionKindTimescale {
		return invalid("kind")
	}
	if partition.ChunkInterval != TimescaleChunkInterval {
		return invalid("chunk interval")
	}
	if partition.ExtensionSchema != "" && !validGroupIdentifier(partition.ExtensionSchema) {
		return invalid("extension schema")
	}
	for _, column := range layout.Columns {
		if column.Name != partition.Column {
			continue
		}
		if partition.Column != TimescalePartitionColumn || column.Nullable || !column.PrimaryKey || column.SQLType != sqlTypeTimestamptz {
			return invalid("column")
		}
		return nil
	}
	return invalid("column")
}

// validateGroupSchemaIndexes normalizes the declared indexes. Names share the
// table namespace of the target, so they must not collide with a table either.
func validateGroupSchemaIndexes(layout GroupSchemaLayout, columns map[string]struct{}) ([]GroupSchemaIndex, error) {
	if len(layout.Indexes) == 0 {
		return nil, nil
	}
	indexes := make([]GroupSchemaIndex, len(layout.Indexes))
	names := map[string]struct{}{strings.ToLower(layout.TableName): {}, managedEffectReceiptTable: {}}
	for i, index := range layout.Indexes {
		index.Name = strings.TrimSpace(index.Name)
		if !validGroupIdentifier(index.Name) {
			return nil, fmt.Errorf("%w: index name", ErrInvalidGroupSchemaLayout)
		}
		if _, exists := names[strings.ToLower(index.Name)]; exists {
			return nil, fmt.Errorf("%w: duplicate index name %s", ErrInvalidGroupSchemaLayout, index.Name)
		}
		names[strings.ToLower(index.Name)] = struct{}{}
		if len(index.Columns) == 0 {
			return nil, fmt.Errorf("%w: index %s has no columns", ErrInvalidGroupSchemaLayout, index.Name)
		}
		seen := make(map[string]struct{}, len(index.Columns))
		normalized := make([]string, len(index.Columns))
		for j, column := range index.Columns {
			column = strings.TrimSpace(column)
			key := strings.ToLower(column)
			if _, known := columns[key]; !known {
				return nil, fmt.Errorf("%w: index %s uses an unknown column", ErrInvalidGroupSchemaLayout, index.Name)
			}
			if _, duplicate := seen[key]; duplicate {
				return nil, fmt.Errorf("%w: index %s repeats a column", ErrInvalidGroupSchemaLayout, index.Name)
			}
			seen[key] = struct{}{}
			normalized[j] = column
		}
		index.Columns = normalized
		indexes[i] = index
	}
	return indexes, nil
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
		return sqlTypeBoolean
	case "INT", "INT4":
		return sqlTypeInteger
	case "INT8":
		return sqlTypeBigint
	case "FLOAT8", "DOUBLE":
		return sqlTypeDoublePrecision
	case "DECIMAL":
		return managedNumericSQLType
	case "TIMESTAMP WITH TIME ZONE":
		return sqlTypeTimestamptz
	default:
		return value
	}
}

func groupSQLTypesCompatible(expected, actual string) bool {
	expected = normalizeGroupSQLType(expected)
	actual = observedGroupSQLType(actual)
	if (expected == "NUMERIC(20,0)" || expected == managedNumericSQLType) && actual == managedNumericSQLType {
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
	if layout.MySQL != nil {
		if normalizePreviewDialect(dialect) != dialectMySQL {
			return nil, fmt.Errorf("%w: a mysql layout needs the mysql dialect", ErrUnsupportedSchemaDialect)
		}
		return generateMySQLGroupDDL(layout)
	}
	if layout.SQLServer != nil {
		if normalizePreviewDialect(dialect) != dialectSQLServer {
			return nil, fmt.Errorf("%w: a sqlserver layout needs the sqlserver dialect", ErrUnsupportedSchemaDialect)
		}
		return generateSQLServerGroupDDL(layout)
	}
	switch normalizePreviewDialect(dialect) {
	case dialectSQLite, dialectPostgres:
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedSchemaDialect, dialect)
	}
	if layout.Partition != nil {
		if normalizePreviewDialect(dialect) != dialectPostgres {
			return nil, fmt.Errorf("%w: a timescale partition needs PostgreSQL", ErrUnsupportedSchemaDialect)
		}
		if layout.Partition.ExtensionSchema == "" {
			return nil, fmt.Errorf("%w: partition extension schema is unknown", ErrInvalidGroupSchemaLayout)
		}
	}

	// One key column stays inline (the legacy form). Several key columns become a
	// single table-level composite key; repeating inline PRIMARY KEY is invalid.
	var keyColumns []string
	for _, column := range layout.Columns {
		if column.PrimaryKey {
			keyColumns = append(keyColumns, quoteGroupIdentifier(column.Name))
		}
	}
	inlineKey := len(keyColumns) == 1
	columnDefinitions := make([]string, 0, len(layout.Columns)+1)
	for _, column := range layout.Columns {
		definition := quoteGroupIdentifier(column.Name) + " " + column.SQLType
		if !column.Nullable {
			definition += " NOT NULL"
		}
		if column.PrimaryKey && inlineKey {
			definition += " PRIMARY KEY"
		}
		if column.Name == layout.OwnerColumn {
			definition += " DEFAULT '" + managedOwnerDefault + "'"
		}
		columnDefinitions = append(columnDefinitions, definition)
	}
	if len(keyColumns) > 1 {
		columnDefinitions = append(columnDefinitions, "PRIMARY KEY ("+strings.Join(keyColumns, ", ")+")")
	}
	receiptDefinitions := []string{
		quoteGroupIdentifier("effect_key") + " " + managedTextSQLType + " PRIMARY KEY",
		quoteGroupIdentifier("payload_digest") + " " + managedTextSQLType + " NOT NULL",
		quoteGroupIdentifier("committed_at") + " " + managedTextSQLType + " NOT NULL",
		quoteGroupIdentifier(managedEffectReceiptMarker) + " " + managedTextSQLType + " NOT NULL DEFAULT '" + managedOwnerDefault + "'",
	}
	// Order: group table, its indexes, then the shared receipt table. Callers use
	// the first and last entries and the indexes in between.
	statements := []string{"CREATE TABLE " + quoteGroupIdentifier(layout.TableName) + " (" + strings.Join(columnDefinitions, ", ") + ");"}
	if layout.Partition != nil {
		// The table was created empty in this same batch, so there is nothing to
		// migrate; if_not_exists is false so an existing hypertable can never turn
		// into a silent success, and the default indexes are off so the table keeps
		// exactly the confirmed key and index.
		schema := quoteGroupIdentifier(layout.Partition.ExtensionSchema)
		create, byRange := schema+"."+quoteGroupIdentifier("create_hypertable"), schema+"."+quoteGroupIdentifier("by_range")
		statements = append(statements, "SELECT "+create+"('"+quoteGroupIdentifier(layout.TableName)+"'::regclass, "+
			byRange+"('"+layout.Partition.Column+"'::name, INTERVAL '"+layout.Partition.ChunkInterval+"'), "+
			"create_default_indexes => false, if_not_exists => false, migrate_data => false);")
	}
	for _, index := range layout.Indexes {
		quoted := make([]string, len(index.Columns))
		for i, column := range index.Columns {
			quoted[i] = quoteGroupIdentifier(column)
		}
		statements = append(statements, "CREATE INDEX "+quoteGroupIdentifier(index.Name)+" ON "+quoteGroupIdentifier(layout.TableName)+" ("+strings.Join(quoted, ", ")+");")
	}
	return append(statements, "CREATE TABLE "+quoteGroupIdentifier(managedEffectReceiptTable)+" ("+strings.Join(receiptDefinitions, ", ")+");"), nil
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
