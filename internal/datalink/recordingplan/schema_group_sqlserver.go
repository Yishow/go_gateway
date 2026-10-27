package recordingplan

import (
	"cmp"
	"fmt"
	"strings"
)

const (
	dialectSQLServer          = "sqlserver"
	sqlserverNotNullClause    = " NOT NULL"
	sqlserverBucketColumn     = "bucket_start"
	sqlserverWorkspaceColumn  = "workspace_id"
	sqlserverTagColumn        = "tag_id"
	sqlserverProvenanceColumn = "provenance"

	sqlserverCollation = "Latin1_General_100_BIN2"

	sqlserverTimeSQLType    = "DATETIME2(6)"
	sqlserverBigIntSQLType  = "BIGINT"
	sqlserverDecimalSQLType = "DECIMAL(38,18)"
	sqlserverFloatSQLType   = "FLOAT(53)"
	sqlserverBitSQLType     = "BIT"
	sqlserverTextSQLType    = "NVARCHAR(MAX)"
	sqlserverIDSQLType      = "NVARCHAR(191)"
	sqlserverHashSQLType    = "CHAR(64)"
	sqlserverOwnerSQLType   = "NVARCHAR(16)"
)

var allowedSQLServerGroupSQLTypes = map[string]struct{}{
	sqlserverTimeSQLType:    {},
	sqlserverBigIntSQLType:  {},
	sqlserverDecimalSQLType: {},
	sqlserverFloatSQLType:   {},
	sqlserverBitSQLType:     {},
	sqlserverTextSQLType:    {},
	sqlserverIDSQLType:      {},
	sqlserverHashSQLType:    {},
	sqlserverOwnerSQLType:   {},
}

// SQLServerGroupSchemaOptions binds a managed long table to SQL Server specific
// layout parameters, including the explicit database, schema, case-exact collation,
// and verified server capability facts.
type SQLServerGroupSchemaOptions struct {
	Database           string `json:"database,omitempty"`
	Schema             string `json:"schema"`
	Collation          string `json:"collation,omitempty"`
	ServerVersion      string `json:"server_version,omitempty"`
	ServerEdition      string `json:"server_edition,omitempty"`
	CapabilityRevision string `json:"capability_revision,omitempty"`
	EngineEdition      int64  `json:"engine_edition,omitzero"`
	CompatibilityLevel int64  `json:"compatibility_level,omitzero"`
	// NativeMonthly freezes the initial owned SQL Server partition layout.
	// Nil preserves the plain SQL Server layout and its existing JSON bytes.
	NativeMonthly *SQLServerPartitionEvidence `json:"native_monthly,omitempty"`
}

// ManagedSQLServerLongGroupLayout returns the canonical SQL Server managed long layout.
// It iterates the shared longColumns contract so names, order, nullability, and the
// composite primary key (bucket_start, record_id) are preserved across database dialects.
// Options must specify the explicit SQL Server schema.
func ManagedSQLServerLongGroupLayout(tableName, ownerColumn string, options SQLServerGroupSchemaOptions) (GroupSchemaLayout, error) {
	normalized, err := validateSQLServerGroupSchemaOptions(options)
	if err != nil {
		return GroupSchemaLayout{}, err
	}
	layout := GroupSchemaLayout{
		TableName:   tableName,
		OwnerColumn: ownerColumn,
		SQLServer:   &normalized,
		Columns:     make([]GroupSchemaColumn, 0, len(longColumns)+1),
		Indexes: []GroupSchemaIndex{{
			Name:    "idx_" + tableName + "_src_time",
			Columns: []string{sqlserverWorkspaceColumn, sqlserverTagColumn, sqlserverBucketColumn},
		}},
	}
	for _, col := range longColumns {
		sqlType := sqlserverColumnType(col)
		layout.Columns = append(layout.Columns, GroupSchemaColumn{
			Name:       col.name,
			SQLType:    sqlType,
			Nullable:   col.nullable,
			PrimaryKey: col.primaryKey,
		})
	}
	layout.Columns = append(layout.Columns, GroupSchemaColumn{
		Name:    ownerColumn,
		SQLType: sqlserverOwnerSQLType,
	})
	return ValidateGroupSchemaLayout(layout)
}

func sqlserverColumnType(col longColumn) string {
	switch col.name {
	case sqlserverBucketColumn, "bucket_end", "observed_at", "received_at", "delta_from", "delta_to":
		return sqlserverTimeSQLType
	case "record_id":
		return sqlserverHashSQLType
	case "value_int":
		return sqlserverBigIntSQLType
	case "value_decimal", "delta_value":
		return sqlserverDecimalSQLType
	case "value_float":
		return sqlserverFloatSQLType
	case "value_bool":
		return sqlserverBitSQLType
	case "value_text", sqlserverProvenanceColumn:
		return sqlserverTextSQLType
	default:
		return sqlserverIDSQLType
	}
}

func validateSQLServerGroupSchemaOptions(options SQLServerGroupSchemaOptions) (SQLServerGroupSchemaOptions, error) {
	invalid := func(reason string) (SQLServerGroupSchemaOptions, error) {
		return SQLServerGroupSchemaOptions{}, fmt.Errorf("%w: sqlserver %s", ErrInvalidGroupSchemaLayout, reason)
	}
	normalized := options
	normalized.Database = strings.TrimSpace(normalized.Database)
	if strings.ContainsRune(normalized.Database, '\x00') {
		return invalid("database cannot contain NUL")
	}
	normalized.Schema = strings.TrimSpace(normalized.Schema)
	if normalized.Schema == "" || !validGroupIdentifier(normalized.Schema) {
		return invalid("schema")
	}
	normalized.Collation = cmp.Or(strings.TrimSpace(normalized.Collation), sqlserverCollation)
	if !strings.EqualFold(normalized.Collation, sqlserverCollation) {
		return invalid("collation")
	}
	normalized.Collation = sqlserverCollation
	normalized.ServerVersion = strings.TrimSpace(normalized.ServerVersion)
	normalized.ServerEdition = strings.TrimSpace(normalized.ServerEdition)
	normalized.CapabilityRevision = strings.TrimSpace(normalized.CapabilityRevision)
	normalized.NativeMonthly = cloneSQLServerPartitionEvidence(normalized.NativeMonthly)
	return normalized, nil
}

// validateSQLServerPreviewScope enforces strict database, schema, dialect and
// verified server capability facts before target inspection.
func validateSQLServerPreviewScope(database, schemaName, dialect string, options SQLServerGroupSchemaOptions) error {
	if normalizePreviewDialect(dialect) != dialectSQLServer {
		return fmt.Errorf("%w: a sqlserver layout needs the sqlserver dialect", ErrUnsupportedSchemaDialect)
	}
	database = strings.TrimSpace(database)
	if database == "" || strings.ContainsRune(database, '\x00') {
		return fmt.Errorf("%w: sqlserver database is required and cannot contain NUL", ErrInvalidGroupSchemaLayout)
	}
	schemaName = strings.TrimSpace(schemaName)
	if schemaName == "" || !validGroupIdentifier(schemaName) {
		return fmt.Errorf("%w: sqlserver schema is required and must be a valid identifier", ErrInvalidGroupSchemaLayout)
	}
	if schemaName != options.Schema {
		return fmt.Errorf("%w: scope schema %q does not match sqlserver options schema %q", ErrInvalidGroupSchemaLayout, schemaName, options.Schema)
	}
	if options.Database == "" || strings.ContainsRune(options.Database, '\x00') {
		return fmt.Errorf("%w: sqlserver options database is required and cannot contain NUL", ErrInvalidGroupSchemaLayout)
	}
	if options.Database != database {
		return fmt.Errorf("%w: scope database %q does not match sqlserver options database %q", ErrInvalidGroupSchemaLayout, database, options.Database)
	}
	return validateSQLServerPreviewFacts(options)
}

func validateSQLServerPreviewFacts(options SQLServerGroupSchemaOptions) error {
	invalid := func(reason string) error {
		return fmt.Errorf("%w: sqlserver facts %s", ErrInvalidGroupSchemaLayout, reason)
	}
	version16 := strings.HasPrefix(options.ServerVersion, "16.") || options.ServerVersion == "16"
	version17 := strings.HasPrefix(options.ServerVersion, "17.") || options.ServerVersion == "17"
	if options.ServerVersion == "" || (!version16 && !version17) {
		return invalid("server version must be SQL Server 2022 baseline (16.*) or 2025 (17.*)")
	}
	if options.ServerEdition == "" {
		return invalid("server edition is required")
	}
	if options.CapabilityRevision == "" {
		return invalid("capability revision is required")
	}
	if options.EngineEdition != 2 && options.EngineEdition != 3 && options.EngineEdition != 4 {
		return invalid("engine edition must be 2 (Standard), 3 (Enterprise), or 4 (Express)")
	}
	if (version16 && options.CompatibilityLevel != 160) || (version17 && options.CompatibilityLevel != 170) || (!version16 && !version17) {
		return invalid("compatibility level must be 160 (SQL Server 2022) or 170 (SQL Server 2025)")
	}
	return nil
}

func quoteSQLServerIdentifier(s string) string {
	return "[" + strings.ReplaceAll(s, "]", "]]") + "]"
}

func sqlserverCollatedSQLType(sqlType string) bool {
	return strings.HasPrefix(sqlType, "NVARCHAR(") || strings.HasPrefix(sqlType, "CHAR(")
}

// generateSQLServerGroupDDL emits SQL Server CREATE TABLE and CREATE INDEX statements.
// It uses bracket quoting, explicit schema qualification, exact Unicode types,
// and standard ISJSON check constraint on provenance.
func generateSQLServerGroupDDL(layout GroupSchemaLayout) ([]string, error) {
	if layout.SQLServer == nil {
		return nil, fmt.Errorf("%w: missing sqlserver options", ErrInvalidGroupSchemaLayout)
	}
	options := *layout.SQLServer
	collation := options.Collation
	schema := quoteSQLServerIdentifier(options.Schema)
	table := quoteSQLServerIdentifier(layout.TableName)
	fullTable := schema + "." + table
	fullReceipt := schema + "." + quoteSQLServerIdentifier(managedEffectReceiptTable)

	var keyColumns []string
	for _, col := range layout.Columns {
		if col.PrimaryKey {
			keyColumns = append(keyColumns, quoteSQLServerIdentifier(col.Name))
		}
	}

	columnDefinitions := make([]string, 0, len(layout.Columns)+1)
	for _, col := range layout.Columns {
		definition := quoteSQLServerIdentifier(col.Name) + " " + col.SQLType
		if sqlserverCollatedSQLType(col.SQLType) {
			definition += " COLLATE " + collation
		}
		if !col.Nullable {
			definition += sqlserverNotNullClause
		} else {
			definition += " NULL"
		}
		if col.Name == sqlserverProvenanceColumn {
			definition += " CHECK (ISJSON(" + quoteSQLServerIdentifier(col.Name) + ") = 1)"
		}
		if col.Name == layout.OwnerColumn {
			definition += " DEFAULT N'" + managedOwnerDefault + "'"
		}
		columnDefinitions = append(columnDefinitions, definition)
	}
	if len(keyColumns) > 0 {
		keyType := "PRIMARY KEY"
		if options.NativeMonthly != nil {
			keyType = "PRIMARY KEY CLUSTERED"
		}
		columnDefinitions = append(columnDefinitions, keyType+" ("+strings.Join(keyColumns, ", ")+")")
	}

	statements := make([]string, 0, len(layout.Indexes)+6)
	if options.NativeMonthly != nil {
		partitionStatements, err := sqlServerPartitionDDL(layout, options.NativeMonthly)
		if err != nil {
			return nil, err
		}
		statements = append(statements, partitionStatements...)
	}
	tableStatement := "CREATE TABLE " + fullTable + " (" + strings.Join(columnDefinitions, ", ") + ")"
	if options.NativeMonthly != nil {
		tableStatement += " ON " + quoteSQLServerIdentifier(options.NativeMonthly.SchemeName) + " ([" + sqlserverPartitionColumn + "])"
	}
	statements = append(statements, tableStatement+";")

	for _, index := range layout.Indexes {
		quoted := make([]string, len(index.Columns))
		for i, col := range index.Columns {
			quoted[i] = quoteSQLServerIdentifier(col)
		}
		indexStatement := "CREATE INDEX " + quoteSQLServerIdentifier(index.Name) + " ON " + fullTable + " (" + strings.Join(quoted, ", ") + ")"
		if options.NativeMonthly != nil {
			indexStatement += " ON " + quoteSQLServerIdentifier(options.NativeMonthly.SchemeName) + " ([" + sqlserverPartitionColumn + "])"
		}
		statements = append(statements, indexStatement+";")
	}

	receiptDefinitions := []string{
		quoteSQLServerIdentifier("effect_key") + " " + sqlserverHashSQLType + " COLLATE " + collation + sqlserverNotNullClause,
		quoteSQLServerIdentifier("payload_digest") + " " + sqlserverHashSQLType + " COLLATE " + collation + sqlserverNotNullClause,
		quoteSQLServerIdentifier("committed_at") + " " + sqlserverTimeSQLType + sqlserverNotNullClause,
		quoteSQLServerIdentifier(managedEffectReceiptMarker) + " " + sqlserverOwnerSQLType + " COLLATE " + collation + " NOT NULL DEFAULT N'" + managedOwnerDefault + "'",
		"PRIMARY KEY (" + quoteSQLServerIdentifier("effect_key") + ")",
	}
	statements = append(statements, "CREATE TABLE "+fullReceipt+" ("+strings.Join(receiptDefinitions, ", ")+");")
	return statements, nil
}
