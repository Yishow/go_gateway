package recordingplan

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerLongSchemaTypes(t *testing.T) {
	const (
		testSchema = "dbo"
		testTable  = "gw_long_0123456789abcdef"
		testOwner  = "_gw_owner_0123456789abcdef01234567"
	)

	t.Run("generic_unsupported_path", func(t *testing.T) {
		// Calling ManagedLongGroupLayout with "sqlserver" dialect remains unsupported
		// because SQL Server requires an explicit schema argument and specific options.
		_, err := ManagedLongGroupLayout("sqlserver", testTable, testOwner)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrUnsupportedSchemaDialect), "generic constructor must report ErrUnsupportedSchemaDialect")
	})

	t.Run("explicit_constructor_types", func(t *testing.T) {
		options := SQLServerGroupSchemaOptions{
			Schema:    testSchema,
			Collation: sqlserverCollation,
		}
		layout, err := ManagedSQLServerLongGroupLayout(testTable, testOwner, options)
		require.NoError(t, err)
		require.Equal(t, testTable, layout.TableName)
		require.Equal(t, testOwner, layout.OwnerColumn)
		require.NotNil(t, layout.SQLServer)
		require.Equal(t, testSchema, layout.SQLServer.Schema)
		require.Equal(t, sqlserverCollation, layout.SQLServer.Collation)

		// 28 long columns + 1 owner column = 29 columns
		require.Len(t, layout.Columns, 29)

		colTypes := make(map[string]string, len(layout.Columns))
		colNullable := make(map[string]bool, len(layout.Columns))
		colPK := make(map[string]bool, len(layout.Columns))
		for _, col := range layout.Columns {
			colTypes[col.Name] = col.SQLType
			colNullable[col.Name] = col.Nullable
			colPK[col.Name] = col.PrimaryKey

			// Must never use rowversion or SQL Server's legacy timestamp keyword
			require.False(t, strings.Contains(strings.ToLower(col.SQLType), "rowversion"), "rowversion must never be used")
			require.False(t, strings.EqualFold(col.SQLType, "timestamp"), "timestamp keyword in SQL Server means rowversion and must not be used")
		}

		// Verify UTC datetime2(6)
		for _, colName := range []string{"bucket_start", "bucket_end", "observed_at", "received_at", "delta_from", "delta_to"} {
			require.Equal(t, "DATETIME2(6)", colTypes[colName], "column %s must be DATETIME2(6)", colName)
		}

		// Verify bounded Unicode sourceIDs NVARCHAR(191)
		for _, colName := range []string{
			"workspace_id", "group_id", "group_revision", "circuit_id", "device_id",
			"point_id", "tag_id", "measurement_id", "series_epoch", "profile_kind",
			"value_type", "unit", "status", "quality",
		} {
			require.Equal(t, "NVARCHAR(191)", colTypes[colName], "column %s must be NVARCHAR(191)", colName)
		}

		// Verify hash, int, decimal, float, bit, text, provenance and owner
		require.Equal(t, "CHAR(64)", colTypes["record_id"])
		require.Equal(t, "BIGINT", colTypes["value_int"])
		require.Equal(t, "DECIMAL(38,18)", colTypes["value_decimal"])
		require.Equal(t, "DECIMAL(38,18)", colTypes["delta_value"])
		require.Equal(t, "FLOAT(53)", colTypes["value_float"])
		require.Equal(t, "BIT", colTypes["value_bool"])
		require.Equal(t, "NVARCHAR(MAX)", colTypes["value_text"])
		require.Equal(t, "NVARCHAR(MAX)", colTypes["provenance"])
		require.Equal(t, "NVARCHAR(16)", colTypes[testOwner])

		// Primary Key columns: (bucket_start, record_id)
		require.True(t, colPK["bucket_start"])
		require.True(t, colPK["record_id"])
		require.False(t, colNullable["bucket_start"])
		require.False(t, colNullable["record_id"])

		// Nullability checks
		require.False(t, colNullable["workspace_id"])
		require.False(t, colNullable["tag_id"])
		require.True(t, colNullable["circuit_id"])
		require.True(t, colNullable["value_decimal"])
		require.True(t, colNullable["delta_value"])
		require.False(t, colNullable["provenance"])
		require.False(t, colNullable[testOwner])
	})

	t.Run("options_validation_and_mismatches", func(t *testing.T) {
		options := SQLServerGroupSchemaOptions{Schema: testSchema}

		// Empty/missing schema name is strictly rejected (no default fallback)
		_, err := ManagedSQLServerLongGroupLayout(testTable, testOwner, SQLServerGroupSchemaOptions{})
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)

		_, err = ManagedSQLServerLongGroupLayout(testTable, testOwner, SQLServerGroupSchemaOptions{Schema: ""})
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)

		_, err = ManagedSQLServerLongGroupLayout(testTable, testOwner, SQLServerGroupSchemaOptions{Schema: "   "})
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)

		// Invalid schema identifier
		_, err = ManagedSQLServerLongGroupLayout(testTable, testOwner, SQLServerGroupSchemaOptions{Schema: "invalid schema"})
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)

		_, err = ManagedSQLServerLongGroupLayout(testTable, testOwner, SQLServerGroupSchemaOptions{Schema: "dbo;DROP TABLE x;"})
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)

		// Invalid owner column pattern
		_, err = ManagedSQLServerLongGroupLayout(testTable, "invalid_owner", options)
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)

		// Invalid table name
		_, err = ManagedSQLServerLongGroupLayout("invalid table", testOwner, options)
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)

		// Unsupported collation
		_, err = ManagedSQLServerLongGroupLayout(testTable, testOwner, SQLServerGroupSchemaOptions{
			Schema:    testSchema,
			Collation: "SQL_Latin1_General_CP1_CI_AS",
		})
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)

		// Plain only: generic partition metadata rejected
		layout, err := ManagedSQLServerLongGroupLayout(testTable, testOwner, options)
		require.NoError(t, err)
		layout.Partition = &GroupSchemaPartition{Kind: PartitionKindTimescale, Column: "bucket_start", ChunkInterval: "7 days"}
		_, err = ValidateGroupSchemaLayout(layout)
		require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)

		// Dialect mismatch in GenerateGroupSchemaDDL
		layout, err = ManagedSQLServerLongGroupLayout(testTable, testOwner, options)
		require.NoError(t, err)
		_, err = GenerateGroupSchemaDDL("postgres", layout)
		require.ErrorIs(t, err, ErrUnsupportedSchemaDialect)

		_, err = GenerateGroupSchemaDDL("mysql", layout)
		require.ErrorIs(t, err, ErrUnsupportedSchemaDialect)
	})
}

func TestSQLServerLongDDL(t *testing.T) {
	const (
		testSchema = "custom_schema"
		testTable  = "gw_long_7f8e9d"
		testOwner  = "_gw_owner_0123456789abcdef01234567"
	)

	t.Run("canonical_ddl_generation", func(t *testing.T) {
		options := SQLServerGroupSchemaOptions{Schema: testSchema}
		layout, err := ManagedSQLServerLongGroupLayout(testTable, testOwner, options)
		require.NoError(t, err)

		statements, err := GenerateGroupSchemaDDL("sqlserver", layout)
		require.NoError(t, err)
		require.Len(t, statements, 3, "should generate main table, 1 index, and 1 receipt table")

		tableStmt := statements[0]
		indexStmt := statements[1]
		receiptStmt := statements[2]

		// 1. Main Table DDL
		require.True(t, strings.HasPrefix(tableStmt, "CREATE TABLE [custom_schema].[gw_long_7f8e9d] ("))
		require.True(t, strings.HasSuffix(tableStmt, ");"))

		// Check PK
		require.Contains(t, tableStmt, "PRIMARY KEY ([bucket_start], [record_id])")

		// Check Unicode and collations
		require.Contains(t, tableStmt, "[workspace_id] NVARCHAR(191) COLLATE Latin1_General_100_BIN2 NOT NULL")
		require.Contains(t, tableStmt, "[record_id] CHAR(64) COLLATE Latin1_General_100_BIN2 NOT NULL")
		require.Contains(t, tableStmt, "[value_text] NVARCHAR(MAX) COLLATE Latin1_General_100_BIN2")

		// Check exact decimal and types
		require.Contains(t, tableStmt, "[value_decimal] DECIMAL(38,18)")
		require.Contains(t, tableStmt, "[delta_value] DECIMAL(38,18)")
		require.Contains(t, tableStmt, "[value_float] FLOAT(53)")
		require.Contains(t, tableStmt, "[value_bool] BIT")
		require.Contains(t, tableStmt, "[bucket_start] DATETIME2(6) NOT NULL")

		// Check provenance with ISJSON constraint
		require.Contains(t, tableStmt, "[provenance] NVARCHAR(MAX) COLLATE Latin1_General_100_BIN2 NOT NULL CHECK (ISJSON([provenance]) = 1)")

		// Check owner default
		require.Contains(t, tableStmt, "["+testOwner+"] NVARCHAR(16) COLLATE Latin1_General_100_BIN2 NOT NULL DEFAULT N'gateway'")

		// 2. Index DDL
		require.Equal(t, "CREATE INDEX [idx_gw_long_7f8e9d_src_time] ON [custom_schema].[gw_long_7f8e9d] ([workspace_id], [tag_id], [bucket_start]);", indexStmt)

		// Verify index key width <= 900 bytes:
		// workspace_id NVARCHAR(191) = 382 bytes
		// tag_id NVARCHAR(191) = 382 bytes
		// bucket_start datetime2(6) = 8 bytes
		// Total: 772 bytes <= 900 bytes
		indexKeyBytes := (191 * 2) + (191 * 2) + 8
		require.LessOrEqual(t, indexKeyBytes, 900, "source/time index key size must not exceed 900 bytes limit")

		// 3. Receipt Table DDL
		require.True(t, strings.HasPrefix(receiptStmt, "CREATE TABLE [custom_schema].[gw_effect_receipts] ("))
		require.True(t, strings.HasSuffix(receiptStmt, ");"))
		require.Contains(t, receiptStmt, "[effect_key] CHAR(64) COLLATE Latin1_General_100_BIN2 NOT NULL")
		require.Contains(t, receiptStmt, "[payload_digest] CHAR(64) COLLATE Latin1_General_100_BIN2 NOT NULL")
		require.Contains(t, receiptStmt, "[committed_at] DATETIME2(6) NOT NULL")
		require.Contains(t, receiptStmt, "[_gw_effect_receipts_v1] NVARCHAR(16) COLLATE Latin1_General_100_BIN2 NOT NULL DEFAULT N'gateway'")
		require.Contains(t, receiptStmt, "PRIMARY KEY ([effect_key])")
	})

	t.Run("bracket_escaping_and_injection_defense", func(t *testing.T) {
		require.Equal(t, "[table]", quoteSQLServerIdentifier("table"))
		require.Equal(t, "[my]]table]", quoteSQLServerIdentifier("my]table"))
		require.Equal(t, "[foo]]]]bar]", quoteSQLServerIdentifier("foo]]bar"))
		require.Equal(t, "[col[name]]]", quoteSQLServerIdentifier("col[name]"))
	})

	t.Run("forbidden_syntax_omitted", func(t *testing.T) {
		options := SQLServerGroupSchemaOptions{Schema: "dbo"}
		layout, err := ManagedSQLServerLongGroupLayout(testTable, testOwner, options)
		require.NoError(t, err)

		statements, err := GenerateGroupSchemaDDL("sqlserver", layout)
		require.NoError(t, err)
		joined := strings.Join(statements, "\n")

		require.NotContains(t, joined, "`", "must not use backtick quoting")
		require.NotContains(t, joined, "$1", "must not use PostgreSQL parameter placeholders")
		require.NotContains(t, joined, "IF NOT EXISTS", "must omit IF NOT EXISTS")
		require.NotContains(t, joined, "MERGE", "must not use MERGE")
		require.NotContains(t, joined, "DROP", "must not use DROP")
		require.NotContains(t, joined, "CREATE DATABASE", "must not CREATE DATABASE")
		require.NotContains(t, joined, "CREATE SCHEMA", "must not CREATE SCHEMA")
		require.NotContains(t, joined, "FILEGROUP", "must not use FILEGROUP")
		require.NotContains(t, joined, "rowversion", "must not use rowversion")
	})
}

func TestValidateSQLServerPreviewFacts2025(t *testing.T) {
	options := SQLServerGroupSchemaOptions{
		Database:           "gw_circuit_acceptance",
		Schema:             "dbo",
		ServerVersion:      "17.0.5035.5",
		ServerEdition:      "Enterprise Developer Edition (64-bit)",
		CapabilityRevision: "sqlservercap-v1-test",
		EngineEdition:      3,
		CompatibilityLevel: 170,
	}
	err := validateSQLServerPreviewFacts(options)
	require.NoError(t, err)
}
