package recordingplan

import (
	"cmp"
	"fmt"
	"regexp"
	"strings"
)

const (
	dialectMySQL = "mysql"

	mysqlFamilyMySQL   = "mysql"
	mysqlFamilyMariaDB = "mariadb"
	mysqlEngineInnoDB  = "InnoDB"
	mysqlCharsetUTF8   = "utf8mb4"

	// The two family collations were verified by a read-only fixture SELECT:
	// both exist on charset utf8mb4 and keep case and trailing spaces distinct
	// (case_eq=0, trail_eq=0, NO PAD). They are a fixed family mapping; this
	// package never guesses or discovers another collation.
	mysqlCollationMySQL   = "utf8mb4_0900_bin"
	mysqlCollationMariaDB = "utf8mb4_nopad_bin"

	mysqlOwnerSQLType    = "VARCHAR(64)"
	mysqlHashSQLType     = "CHAR(64)"
	mysqlJSONSQLType     = "JSON"
	mysqlIDSQLType       = "VARCHAR(191)"
	mysqlTimeSQLType     = "DATETIME(6)"
	mysqlDecimalSQLType  = "DECIMAL(65,30)"
	mysqlDoubleSQLType   = "DOUBLE"
	mysqlLongTextSQLType = "LONGTEXT"
	mysqlSourceTagColumn = "tag_id"
	mysqlWorkspaceColumn = "workspace_id"
)

var (
	mysqlServerVersionPattern      = regexp.MustCompile(`^\d+\.\d+(\.\d+)?[A-Za-z0-9._+~:-]*$`)
	mysqlCapabilityRevisionPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:+-]{0,127}$`)
	mysqlFamilyCollations          = map[string]string{
		mysqlFamilyMySQL:   mysqlCollationMySQL,
		mysqlFamilyMariaDB: mysqlCollationMariaDB,
	}
	allowedMySQLGroupSQLTypes = map[string]struct{}{
		mysqlIDSQLType:       {},
		mysqlHashSQLType:     {},
		"VARCHAR(64)":        {},
		mysqlTimeSQLType:     {},
		"BIGINT":             {},
		mysqlDecimalSQLType:  {},
		mysqlDoubleSQLType:   {},
		sqlTypeBoolean:       {},
		mysqlLongTextSQLType: {},
		mysqlJSONSQLType:     {},
	}
)

// MySQLGroupSchemaOptions is the backend-resolved MySQL family binding of one
// managed long table. It is persisted with the layout, so the family, exact
// server version, capability revision, storage engine and case-exact collation
// are part of the confirmed layout fingerprint. NativeMonthly is the optional,
// frozen initial partition layout; nil preserves the plain InnoDB contract.
type MySQLGroupSchemaOptions struct {
	ServerFamily       string                  `json:"server_family"`
	ServerVersion      string                  `json:"server_version"`
	CapabilityRevision string                  `json:"capability_revision"`
	Engine             string                  `json:"engine"`
	Collation          string                  `json:"collation"`
	Charset            string                  `json:"charset,omitempty"`
	NativeMonthly      *MySQLPartitionEvidence `json:"native_monthly,omitempty"`
}

// ManagedMySQLLongGroupLayout is the canonical MySQL long layout. It iterates
// the shared longColumns contract, so names, order, nullability, the composite
// key (bucket_start, record_id) and source/time index are stable across server
// families. NativeMonthly optionally freezes the initial native partitioning.
func ManagedMySQLLongGroupLayout(tableName, ownerColumn string, options MySQLGroupSchemaOptions) (GroupSchemaLayout, error) {
	normalized, err := validateMySQLGroupSchemaOptions(options)
	if err != nil {
		return GroupSchemaLayout{}, err
	}
	layout := GroupSchemaLayout{
		TableName: tableName, OwnerColumn: ownerColumn, MySQL: &normalized,
		Columns: make([]GroupSchemaColumn, 0, len(longColumns)+1),
		Indexes: []GroupSchemaIndex{{
			Name: "idx_" + tableName + "_src_time", Columns: []string{mysqlWorkspaceColumn, mysqlSourceTagColumn, TimescalePartitionColumn},
		}},
	}
	for _, column := range longColumns {
		if column.mysql == "" {
			// A family-less shared column is a programming error, never a table
			// that silently loses a value.
			return GroupSchemaLayout{}, fmt.Errorf("%w: mysql type for column %s", ErrInvalidGroupSchemaLayout, column.name)
		}
		layout.Columns = append(layout.Columns, GroupSchemaColumn{
			Name: column.name, SQLType: column.mysql, Nullable: column.nullable, PrimaryKey: column.primaryKey,
		})
	}
	layout.Columns = append(layout.Columns, GroupSchemaColumn{Name: ownerColumn, SQLType: mysqlOwnerSQLType})
	return ValidateGroupSchemaLayout(layout)
}

func validateMySQLGroupSchemaOptions(options MySQLGroupSchemaOptions) (MySQLGroupSchemaOptions, error) {
	invalid := func(reason string) (MySQLGroupSchemaOptions, error) {
		return MySQLGroupSchemaOptions{}, fmt.Errorf("%w: mysql %s", ErrInvalidGroupSchemaLayout, reason)
	}
	normalized := options
	normalized.ServerFamily = strings.ToLower(strings.TrimSpace(normalized.ServerFamily))
	normalized.ServerVersion = strings.TrimSpace(normalized.ServerVersion)
	normalized.CapabilityRevision = strings.TrimSpace(normalized.CapabilityRevision)
	normalized.Engine = strings.TrimSpace(normalized.Engine)
	normalized.Collation = strings.ToLower(strings.TrimSpace(normalized.Collation))
	normalized.Charset = cmp.Or(strings.ToLower(strings.TrimSpace(normalized.Charset)), mysqlCharsetUTF8)

	expectedCollation, ok := mysqlFamilyCollations[normalized.ServerFamily]
	if !ok {
		return invalid("server family")
	}
	if !mysqlServerVersionPattern.MatchString(normalized.ServerVersion) {
		return invalid("server version")
	}
	if !mysqlCapabilityRevisionPattern.MatchString(normalized.CapabilityRevision) {
		return invalid("capability revision")
	}
	if !strings.EqualFold(normalized.Engine, mysqlEngineInnoDB) {
		return invalid("storage engine")
	}
	normalized.Engine = mysqlEngineInnoDB
	if normalized.Charset != mysqlCharsetUTF8 {
		return invalid("charset")
	}
	if normalized.Collation != expectedCollation {
		return invalid("collation")
	}
	nativeMonthly, err := normalizeMySQLPartitionEvidence(normalized.NativeMonthly)
	if err != nil {
		return MySQLGroupSchemaOptions{}, err
	}
	normalized.NativeMonthly = nativeMonthly
	return normalized, nil
}

// generateMySQLGroupDDL emits MySQL CREATE statements from a validated layout.
// NativeMonthly affects only the group table; the receipt table remains plain.
// It never uses IF NOT EXISTS, INSERT IGNORE or REPLACE, so an unknown same-name
// table fails rather than becoming a silent success.
func generateMySQLGroupDDL(layout GroupSchemaLayout) ([]string, error) {
	options := *layout.MySQL
	collation := options.Collation
	quote := quoteMySQLIdentifier

	var keyColumns []string
	for _, column := range layout.Columns {
		if column.PrimaryKey {
			keyColumns = append(keyColumns, quote(column.Name))
		}
	}
	columnDefinitions := make([]string, 0, len(layout.Columns)+1)
	for _, column := range layout.Columns {
		definition := quote(column.Name) + " " + column.SQLType
		if mysqlCollatedSQLType(column.SQLType) {
			definition += " COLLATE " + collation
		}
		if !column.Nullable {
			definition += " NOT NULL"
		}
		if column.Name == layout.OwnerColumn {
			definition += " DEFAULT '" + managedOwnerDefault + "'"
		}
		columnDefinitions = append(columnDefinitions, definition)
	}
	if len(keyColumns) > 0 {
		columnDefinitions = append(columnDefinitions, "PRIMARY KEY ("+strings.Join(keyColumns, ", ")+")")
	}

	tableOptions := " ENGINE=InnoDB DEFAULT CHARSET=" + options.Charset + " COLLATE " + collation
	groupStatement := "CREATE TABLE " + quote(layout.TableName) + " (" + strings.Join(columnDefinitions, ", ") + ")" + tableOptions
	if options.NativeMonthly != nil {
		partitionDDL, err := mysqlNativeMonthlyPartitionDDL(options.NativeMonthly)
		if err != nil {
			return nil, err
		}
		groupStatement += partitionDDL
	}
	statements := make([]string, 0, len(layout.Indexes)+2)
	statements = append(statements, groupStatement+";")

	for _, index := range layout.Indexes {
		quoted := make([]string, len(index.Columns))
		for i, column := range index.Columns {
			quoted[i] = quote(column)
		}
		statements = append(statements, "CREATE INDEX "+quote(index.Name)+" ON "+quote(layout.TableName)+" ("+strings.Join(quoted, ", ")+");")
	}

	// The receipt table stays plain InnoDB so partition maintenance never owns
	// deduplication. Its primary key is bounded VARCHAR, never TEXT.
	receiptDefinitions := []string{
		quote("effect_key") + " VARCHAR(191) COLLATE " + collation + " NOT NULL",
		quote("payload_digest") + " CHAR(64) COLLATE " + collation + " NOT NULL",
		quote("committed_at") + " DATETIME(6) NOT NULL",
		quote(managedEffectReceiptMarker) + " VARCHAR(64) COLLATE " + collation + " NOT NULL DEFAULT '" + managedOwnerDefault + "'",
		"PRIMARY KEY (" + quote("effect_key") + ")",
	}
	statements = append(statements, "CREATE TABLE "+quote(managedEffectReceiptTable)+" ("+strings.Join(receiptDefinitions, ", ")+")"+tableOptions+";")
	return statements, nil
}

// validateMySQLPreviewScope enforces the single MySQL namespace: database and
// schema must both be explicit and exactly equal. It never synthesizes main or
// public; those defaults belong to SQLite and PostgreSQL.
func validateMySQLPreviewScope(database, schemaName, dialect string) error {
	if normalizePreviewDialect(dialect) != dialectMySQL {
		return fmt.Errorf("%w: mysql layout requires the mysql dialect", ErrInvalidGroupSchemaLayout)
	}
	if strings.TrimSpace(database) == "" || strings.TrimSpace(schemaName) == "" {
		return fmt.Errorf("%w: mysql database and schema are required", ErrInvalidGroupSchemaLayout)
	}
	if database != schemaName {
		return fmt.Errorf("%w: mysql database and schema must be the same scope", ErrInvalidGroupSchemaLayout)
	}
	return nil
}

// sameMySQLGroupSchemaOptions compares the family and native-partition
// evidence that participate in the confirmed layout. Both sides are validated
// before this runs; nested evidence is compared by value, and nil preserves
// plain layouts.
func sameMySQLGroupSchemaOptions(left, right *MySQLGroupSchemaOptions) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return left.ServerFamily == right.ServerFamily && left.ServerVersion == right.ServerVersion &&
		left.CapabilityRevision == right.CapabilityRevision && left.Engine == right.Engine &&
		left.Collation == right.Collation && left.Charset == right.Charset &&
		sameMySQLPartitionEvidence(left.NativeMonthly, right.NativeMonthly)
}

func mysqlCollatedSQLType(sqlType string) bool {
	return strings.HasPrefix(sqlType, "VARCHAR(") || strings.HasPrefix(sqlType, "CHAR(") || sqlType == mysqlLongTextSQLType
}

func quoteMySQLIdentifier(identifier string) string {
	return "`" + strings.ReplaceAll(identifier, "`", "``") + "`"
}
