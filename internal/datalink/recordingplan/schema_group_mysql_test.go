package recordingplan

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	mysqlTestVersion            = "8.4.5"
	mysqlTestCapabilityRevision = "mysql-unverified-1"
	mysqlTestCollation          = "utf8mb4_0900_bin"
	mariadbTestVersion          = "11.4.2-MariaDB-1:11.4.2+maria~ubu2204"
	mariadbTestCollation        = "utf8mb4_nopad_bin"
)

func mysqlTestOptions() MySQLGroupSchemaOptions {
	return MySQLGroupSchemaOptions{
		ServerFamily:       "mysql",
		ServerVersion:      mysqlTestVersion,
		CapabilityRevision: mysqlTestCapabilityRevision,
		Engine:             "InnoDB",
		Collation:          mysqlTestCollation,
	}
}

func mariadbTestOptions() MySQLGroupSchemaOptions {
	return MySQLGroupSchemaOptions{
		ServerFamily:       "mariadb",
		ServerVersion:      mariadbTestVersion,
		CapabilityRevision: "mariadb-unverified-1",
		Engine:             "InnoDB",
		Collation:          mariadbTestCollation,
	}
}

// mysqlExpectedLongColumns is the fixed managed_long_v1 contract in its MySQL
// expression. It is deliberately a literal: the test must not read the
// production column list to decide what it expects.
func mysqlExpectedLongColumns() []GroupSchemaColumn {
	return []GroupSchemaColumn{
		{Name: "bucket_start", SQLType: "DATETIME(6)", PrimaryKey: true},
		{Name: "bucket_end", SQLType: "DATETIME(6)"},
		{Name: "record_id", SQLType: "CHAR(64)", PrimaryKey: true},
		{Name: "workspace_id", SQLType: "VARCHAR(191)"},
		{Name: "group_id", SQLType: "VARCHAR(191)"},
		{Name: "group_revision", SQLType: "VARCHAR(191)"},
		{Name: "circuit_id", SQLType: "VARCHAR(191)", Nullable: true},
		{Name: "device_id", SQLType: "VARCHAR(191)"},
		{Name: "point_id", SQLType: "VARCHAR(191)"},
		{Name: "tag_id", SQLType: "VARCHAR(191)"},
		{Name: "measurement_id", SQLType: "VARCHAR(191)", Nullable: true},
		{Name: "series_epoch", SQLType: "VARCHAR(191)", Nullable: true},
		{Name: "profile_kind", SQLType: "VARCHAR(191)"},
		{Name: "value_type", SQLType: "VARCHAR(191)"},
		{Name: "unit", SQLType: "VARCHAR(191)", Nullable: true},
		{Name: "status", SQLType: "VARCHAR(191)"},
		{Name: "quality", SQLType: "VARCHAR(191)", Nullable: true},
		{Name: "observed_at", SQLType: "DATETIME(6)", Nullable: true},
		{Name: "received_at", SQLType: "DATETIME(6)", Nullable: true},
		{Name: "value_int", SQLType: "BIGINT", Nullable: true},
		{Name: "value_decimal", SQLType: "DECIMAL(65,30)", Nullable: true},
		{Name: "value_float", SQLType: "DOUBLE", Nullable: true},
		{Name: "value_bool", SQLType: "BOOLEAN", Nullable: true},
		{Name: "value_text", SQLType: "LONGTEXT", Nullable: true},
		{Name: "delta_value", SQLType: "DECIMAL(65,30)", Nullable: true},
		{Name: "delta_from", SQLType: "DATETIME(6)", Nullable: true},
		{Name: "delta_to", SQLType: "DATETIME(6)", Nullable: true},
		{Name: "provenance", SQLType: "JSON"},
		{Name: longTestOwner, SQLType: "VARCHAR(64)"},
	}
}

// TestMySQLLongExactValues pins the MySQL type contract that keeps the accepted
// exact values representable: DECIMAL(65,30) for uint64/exact decimals,
// DATETIME(6) for UTC microseconds, bounded VARCHAR IDs with a case-exact
// collation, and native/alias JSON provenance.
func TestMySQLLongExactValues(t *testing.T) {
	layout, err := ManagedMySQLLongGroupLayout(longTestTable, longTestOwner, mysqlTestOptions())
	require.NoError(t, err)
	require.Equal(t, longTestTable, layout.TableName)
	require.Equal(t, longTestOwner, layout.OwnerColumn)
	require.Equal(t, mysqlExpectedLongColumns(), layout.Columns)
	require.Equal(t, []GroupSchemaIndex{{
		Name: "idx_" + longTestTable + "_src_time", Columns: []string{"workspace_id", "tag_id", "bucket_start"},
	}}, layout.Indexes)
	require.NotNil(t, layout.MySQL)
	require.Nil(t, layout.Partition)
	require.Equal(t, "mysql", layout.MySQL.ServerFamily)
	require.Equal(t, mysqlTestVersion, layout.MySQL.ServerVersion)
	require.Equal(t, mysqlTestCapabilityRevision, layout.MySQL.CapabilityRevision)
	require.Equal(t, "InnoDB", layout.MySQL.Engine)
	require.Equal(t, mysqlTestCollation, layout.MySQL.Collation)
}

func TestMySQLLongDDL(t *testing.T) {
	layout, err := ManagedMySQLLongGroupLayout(longTestTable, longTestOwner, mysqlTestOptions())
	require.NoError(t, err)
	statements, err := GenerateGroupSchemaDDL(dialectMySQL, layout)
	require.NoError(t, err)
	require.Len(t, statements, 3)

	create := statements[0]
	require.True(t, strings.HasPrefix(create, "CREATE TABLE `"+longTestTable+"` ("))
	for _, fragment := range []string{
		"`bucket_start` DATETIME(6) NOT NULL",
		"`record_id` CHAR(64) COLLATE " + mysqlTestCollation + " NOT NULL",
		"`workspace_id` VARCHAR(191) COLLATE " + mysqlTestCollation + " NOT NULL",
		"`value_int` BIGINT",
		"`value_decimal` DECIMAL(65,30)",
		"`value_float` DOUBLE",
		"`value_bool` BOOLEAN",
		"`value_text` LONGTEXT COLLATE " + mysqlTestCollation,
		"`provenance` JSON NOT NULL",
		"PRIMARY KEY (`bucket_start`, `record_id`)",
		"ENGINE=InnoDB",
		"DEFAULT CHARSET=utf8mb4",
		"COLLATE " + mysqlTestCollation,
	} {
		require.Contains(t, create, fragment)
	}
	require.Contains(t, create, "`"+longTestOwner+"` VARCHAR(64) COLLATE "+mysqlTestCollation+" NOT NULL DEFAULT 'gateway'")
	require.NotContains(t, create, "TIMESTAMP")
	require.NotContains(t, create, "NUMERIC")

	require.Equal(t, "CREATE INDEX `idx_"+longTestTable+"_src_time` ON `"+longTestTable+
		"` (`workspace_id`, `tag_id`, `bucket_start`);", statements[1])

	receipt := statements[2]
	require.True(t, strings.HasPrefix(receipt, "CREATE TABLE `gw_effect_receipts` ("))
	require.Contains(t, receipt, "`effect_key` VARCHAR(191) COLLATE "+mysqlTestCollation+" NOT NULL")
	require.Contains(t, receipt, "`payload_digest` CHAR(64)")
	require.Contains(t, receipt, "`committed_at` DATETIME(6) NOT NULL")
	require.Contains(t, receipt, "`"+managedEffectReceiptMarker+"` VARCHAR(64)")
	require.Contains(t, receipt, "PRIMARY KEY (`effect_key`)")
	require.Contains(t, receipt, "ENGINE=InnoDB")
	require.NotContains(t, receipt, "`effect_key` TEXT")

	upper := strings.ToUpper(strings.Join(statements, "\n"))
	require.NotContains(t, upper, "IF NOT EXISTS")
	require.NotContains(t, upper, "INSERT IGNORE")
	require.NotContains(t, upper, "REPLACE")
	require.NotContains(t, upper, "FOREIGN KEY")
}

func mysqlTestScope(layout GroupSchemaLayout) SchemaPreviewScope {
	return SchemaPreviewScope{
		WorkspaceID: "ws-mysql", WorkspaceRevision: "setup-1", PlanID: "plan-mysql", PlanRevision: "rev-1",
		ConnectorID: "conn-mysql", ConnectorRevision: "identity-rev-1", Dialect: dialectMySQL,
		Database: "gw", Schema: "gw", TablePrefix: layout.TableName,
		GroupLayout: &layout, SourceDigest: "source-digest", SchemaRevision: "schema-rev-1", SchemaDigest: "schema-digest-1",
	}
}

func mysqlTestToken(layout GroupSchemaLayout) SchemaPreviewToken {
	return SchemaPreviewToken{
		Token: "tok-mysql-scope", OperationID: "op-mysql-scope", Action: SchemaApplyAction,
		WorkspaceID: "ws-mysql", WorkspaceRevision: "setup-1", PlanID: "plan-mysql", PlanRevision: "rev-1",
		ConnectorID: "conn-mysql", ConnectorRevision: "identity-rev-1", Dialect: dialectMySQL,
		Database: "gw", Schema: "gw", TablePrefix: layout.TableName,
		GroupLayout: &layout, SourceDigest: "source-digest", SchemaRevision: "schema-rev-1", SchemaDigest: "schema-digest-1",
	}
}

func TestMySQLLongOptionsArePartOfPreviewScope(t *testing.T) {
	layout, err := ManagedMySQLLongGroupLayout(longTestTable, longTestOwner, mysqlTestOptions())
	require.NoError(t, err)
	scope := mysqlTestScope(layout)
	token := mysqlTestToken(layout)
	require.Empty(t, previewScopeMismatch(&token, scope))

	mutations := map[string]func(*MySQLGroupSchemaOptions){
		"family": func(o *MySQLGroupSchemaOptions) {
			o.ServerFamily = mysqlFamilyMariaDB
			o.Collation = mariadbTestCollation
		},
		"version":             func(o *MySQLGroupSchemaOptions) { o.ServerVersion = "8.4.6" },
		"capability revision": func(o *MySQLGroupSchemaOptions) { o.CapabilityRevision = "mysql-unverified-2" },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			changed := layout
			options := *layout.MySQL
			changed.MySQL = &options
			mutate(changed.MySQL)
			changedScope := mysqlTestScope(changed)
			require.Equal(t, "group layout", previewScopeMismatch(&token, changedScope))
		})
	}
}

func TestMySQLLongTokenRepositoryCopiesOptions(t *testing.T) {
	repo := NewMemoryRepository()
	layout, err := ManagedMySQLLongGroupLayout(longTestTable, longTestOwner, mysqlTestOptions())
	require.NoError(t, err)
	token := mysqlTestToken(layout)
	require.NoError(t, repo.SavePreviewToken(t.Context(), &token))

	// Mutating the caller's options after Save must not reach the repository.
	token.GroupLayout.MySQL.CapabilityRevision = "mutated-after-save"
	saved, err := repo.GetPreviewToken(t.Context(), token.Token)
	require.NoError(t, err)
	require.Equal(t, mysqlTestCapabilityRevision, saved.GroupLayout.MySQL.CapabilityRevision)

	// Mutating a returned copy must not reach the repository either.
	saved.GroupLayout.MySQL.ServerVersion = "9.9.9-mutated"
	again, err := repo.GetPreviewToken(t.Context(), token.Token)
	require.NoError(t, err)
	require.Equal(t, mysqlTestVersion, again.GroupLayout.MySQL.ServerVersion)
}

func TestMySQLLongScope(t *testing.T) {
	layout, err := ManagedMySQLLongGroupLayout(longTestTable, longTestOwner, mysqlTestOptions())
	require.NoError(t, err)

	// MariaDB is the same logical layout with its own verified collation and
	// JSON alias. The family is part of the persisted layout, not inferred from
	// the dialect or a version substring.
	mariadb, err := ManagedMySQLLongGroupLayout(longTestTable, longTestOwner, mariadbTestOptions())
	require.NoError(t, err)
	require.Equal(t, "mariadb", mariadb.MySQL.ServerFamily)
	require.Equal(t, mariadbTestVersion, mariadb.MySQL.ServerVersion)
	require.Equal(t, mariadbTestCollation, mariadb.MySQL.Collation)
	statements, err := GenerateGroupSchemaDDL(dialectMySQL, mariadb)
	require.NoError(t, err)
	require.Contains(t, statements[0], "`provenance` JSON")
	require.Contains(t, statements[0], "COLLATE "+mariadbTestCollation)
	require.Contains(t, statements[2], "ENGINE=InnoDB")

	// Old layouts carry no MySQL options and keep their JSON shape.
	legacy, err := ManagedLongGroupLayout(dialectPostgres, longTestTable, longTestOwner)
	require.NoError(t, err)
	legacyJSON, err := json.Marshal(legacy)
	require.NoError(t, err)
	require.NotContains(t, string(legacyJSON), `"mysql"`)

	mysqlJSON, err := json.Marshal(layout)
	require.NoError(t, err)
	require.Contains(t, string(mysqlJSON), `"mysql":{"server_family":"mysql"`)
	require.Contains(t, string(mysqlJSON), `"capability_revision":"`+mysqlTestCapabilityRevision+`"`)
	require.Contains(t, string(mysqlJSON), `"collation":"`+mysqlTestCollation+`"`)

	// Capability facts change the confirmation digest, while both operations
	// still contend for the same physical table's DDL lease.
	changed := layout
	changedOptions := *layout.MySQL
	changedOptions.CapabilityRevision = "mysql-unverified-2"
	changed.MySQL = &changedOptions
	storedToken := SchemaPreviewToken{GroupLayout: &layout}
	changedToken := SchemaPreviewToken{GroupLayout: &changed}
	require.NotEmpty(t, storedToken.ScopeKey())
	require.Equal(t, storedToken.ScopeKey(), changedToken.ScopeKey())
	storedDigest, err := previewDigest(&storedToken)
	require.NoError(t, err)
	changedDigest, err := previewDigest(&changedToken)
	require.NoError(t, err)
	require.NotEqual(t, storedDigest, changedDigest)

	// A MySQL layout cannot borrow PostgreSQL/Timescale partitioning.
	partitioned := layout
	partitioned.Partition = &GroupSchemaPartition{Kind: PartitionKindTimescale, Column: TimescalePartitionColumn, ChunkInterval: TimescaleChunkInterval, ExtensionSchema: "public"}
	_, err = ValidateGroupSchemaLayout(partitioned)
	require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
	_, err = GenerateGroupSchemaDDL(dialectMySQL, partitioned)
	require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)

	// A MySQL layout is never rendered by another dialect, and the legacy
	// generic MySQL managed DDL refusal stays in force.
	_, err = GenerateGroupSchemaDDL(dialectPostgres, layout)
	require.ErrorIs(t, err, ErrUnsupportedSchemaDialect)
	_, err = GenerateManagedSchemaDDL(dialectMySQL, "gw_")
	require.ErrorIs(t, err, ErrUnsupportedSchemaDialect)
	require.False(t, ManagedSchemaExecutionVerified(dialectMySQL))

	// MySQL has one namespace: database and schema are the same explicit scope.
	// No main/public default is ever synthesized, but a real database named
	// public is valid when both values are exactly equal.
	t.Run("mysql database and schema scope", func(t *testing.T) {
		valid := mysqlTestScope(layout)
		_, err := normalizePreviewScope(valid)
		require.NoError(t, err)

		public := mysqlTestScope(layout)
		public.Database, public.Schema = "public", "public"
		_, err = normalizePreviewScope(public)
		require.NoError(t, err, "a real database named public is allowed when both values match")

		cases := []struct {
			name     string
			database string
			schema   string
			dialect  string
		}{
			{"empty database", "", "gw", dialectMySQL},
			{"empty schema", "gw", "", dialectMySQL},
			{"whitespace database", "   ", "gw", dialectMySQL},
			{"database differs", "gw", "public", dialectMySQL},
			{"another database", "other", "gw", dialectMySQL},
			{"kind mismatch", "gw", "gw", dialectPostgres},
		}
		for _, testCase := range cases {
			t.Run(testCase.name, func(t *testing.T) {
				scope := mysqlTestScope(layout)
				scope.Database, scope.Schema, scope.Dialect = testCase.database, testCase.schema, testCase.dialect
				_, err := normalizePreviewScope(scope)
				require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
			})
		}
	})

	// Other dialects keep their existing scope handling; the MySQL guard must
	// never be applied to them.
	for _, legacyCase := range []struct{ dialect, schema string }{
		{dialectSQLite, "main"},
		{dialectPostgres, "public"},
	} {
		legacyLayout, err := ManagedLongGroupLayout(legacyCase.dialect, longTestTable, longTestOwner)
		require.NoError(t, err)
		scope := SchemaPreviewScope{Dialect: legacyCase.dialect, Database: "gw", Schema: legacyCase.schema, GroupLayout: &legacyLayout}
		_, err = normalizePreviewScope(scope)
		require.NoError(t, err, legacyCase.dialect)
	}
}

func TestMySQLLongScopeRejectsInvalidBindings(t *testing.T) {
	base, err := ManagedMySQLLongGroupLayout(longTestTable, longTestOwner, mysqlTestOptions())
	require.NoError(t, err)

	mutate := func(change func(*GroupSchemaLayout)) GroupSchemaLayout {
		layout := base
		options := *base.MySQL
		layout.MySQL = &options
		columns := slices.Clone(base.Columns)
		layout.Columns = columns
		change(&layout)
		return layout
	}

	cases := []struct {
		name   string
		layout GroupSchemaLayout
	}{
		{"missing options", mutate(func(l *GroupSchemaLayout) { l.MySQL = nil })},
		{"foreign family", mutate(func(l *GroupSchemaLayout) { l.MySQL.ServerFamily = "postgres" })},
		{"empty version", mutate(func(l *GroupSchemaLayout) { l.MySQL.ServerVersion = "" })},
		{"unsafe version", mutate(func(l *GroupSchemaLayout) { l.MySQL.ServerVersion = "8.4.5; DROP" })},
		{"empty capability revision", mutate(func(l *GroupSchemaLayout) { l.MySQL.CapabilityRevision = "" })},
		{"non-innodb engine", mutate(func(l *GroupSchemaLayout) { l.MySQL.Engine = "MyISAM" })},
		{"unknown collation", mutate(func(l *GroupSchemaLayout) { l.MySQL.Collation = "utf8mb4_bin" })},
		{"family collation mismatch", mutate(func(l *GroupSchemaLayout) { l.MySQL.Collation = mariadbTestCollation })},
		{"non-utf8mb4 charset", mutate(func(l *GroupSchemaLayout) { l.MySQL.Charset = "latin1" })},
		{"unbounded decimal", mutate(func(l *GroupSchemaLayout) {
			l.Columns[20].SQLType = "DECIMAL(65,31)"
		})},
		{"unbounded id", mutate(func(l *GroupSchemaLayout) {
			l.Columns[2].SQLType = "VARCHAR(192)"
		})},
		{"text primary key", mutate(func(l *GroupSchemaLayout) {
			l.Columns[2].SQLType = "TEXT"
		})},
		{"text owner", mutate(func(l *GroupSchemaLayout) {
			l.Columns[len(l.Columns)-1].SQLType = "TEXT"
		})},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := ValidateGroupSchemaLayout(testCase.layout)
			require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
			_, err = GenerateGroupSchemaDDL(dialectMySQL, testCase.layout)
			require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
		})
	}
}
