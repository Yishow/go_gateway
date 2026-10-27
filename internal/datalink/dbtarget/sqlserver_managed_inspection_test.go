package dbtarget

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
	"testing"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type sqlServerManagedCatalogScript struct {
	metadata          *sqlServerRecordingMetadataScript
	objectRows        [][]driver.Value
	columnRows        [][]driver.Value
	indexRows         [][]driver.Value
	checkRows         [][]driver.Value
	queries           []sqlServerManagedQueryCall
	queryHasDeadline  bool
	beginCalls        int
	execCalls         int
	closed            bool
	missingObject     bool
	repairOperationID string
	failContains      string
	queryError        error
}

type sqlServerManagedQueryCall struct {
	query string
	args  []driver.NamedValue
}

type sqlServerManagedCatalogConnector struct {
	script *sqlServerManagedCatalogScript
}

func (c sqlServerManagedCatalogConnector) Connect(context.Context) (driver.Conn, error) {
	return &sqlServerManagedCatalogConn{script: c.script}, nil
}

func (sqlServerManagedCatalogConnector) Driver() driver.Driver {
	return sqlServerManagedCatalogDriver{}
}

type sqlServerManagedCatalogDriver struct{}

func (sqlServerManagedCatalogDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("scripted SQL Server catalog connector must be used")
}

type sqlServerManagedCatalogConn struct {
	script *sqlServerManagedCatalogScript
}

func (*sqlServerManagedCatalogConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepared statements are not scripted")
}

func (*sqlServerManagedCatalogConn) Close() error { return nil }

func (c *sqlServerManagedCatalogConn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

func (c *sqlServerManagedCatalogConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	c.script.beginCalls++
	return sqlServerManagedCatalogTx{}, nil
}

func (c *sqlServerManagedCatalogConn) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	c.script.execCalls++
	return nil, errors.New("readonly SQL Server inspection must not execute statements")
}

func (c *sqlServerManagedCatalogConn) QueryContext(
	ctx context.Context,
	query string,
	args []driver.NamedValue,
) (driver.Rows, error) {
	if _, ok := ctx.Deadline(); ok {
		c.script.queryHasDeadline = true
	}
	normalized := strings.Join(strings.Fields(query), " ")
	c.script.queries = append(c.script.queries, sqlServerManagedQueryCall{query: normalized, args: append([]driver.NamedValue(nil), args...)})
	if c.script.failContains != "" && strings.Contains(normalized, c.script.failContains) {
		return nil, c.script.queryError
	}
	switch normalized {
	case sqlServerRecordingMetadataQuery:
		return sqlServerRecordingRows(
			[]string{"version", "edition", "engine_edition", "database", "compatibility"},
			[]driver.Value{"16.0.1000.6", "Developer Edition (64-bit)", int64(3), "fixture_db", int64(160)},
		), nil
	case sqlServerRecordingVisibilityPermissionQuery:
		return sqlServerRecordingRows(
			[]string{"database_view_definition", "schema_view_definition"},
			[]driver.Value{c.script.metadata.viewDefinition, c.script.metadata.schemaViewDefinition},
		), nil
	case sqlServerRecordingSchemaQuery:
		if c.script.metadata.visibleSchema == "" {
			return sqlServerRecordingRows([]string{"name"}), nil
		}
		return sqlServerRecordingRows([]string{"name"}, []driver.Value{c.script.metadata.visibleSchema}), nil
	case sqlServerRecordingPermissionQuery:
		return sqlServerRecordingRows(
			[]string{"create_table", "alter_schema", "select_schema", "insert_schema", "delete_schema", "alter_any_dataspace"},
			[]driver.Value{
				c.script.metadata.createTable, c.script.metadata.alterSchema, c.script.metadata.selectSchema,
				c.script.metadata.insertSchema, c.script.metadata.deleteSchema, c.script.metadata.alterAnyDataspace,
			},
		), nil
	}
	switch {
	case strings.Contains(normalized, "gw_native_partition_objects"):
		return sqlServerRecordingRows([]string{
			"function_id", "right", "function_type", "parameter_type", "scale", "parameter_count",
			"scheme_id", "scheme_function_id", "function_owner", "scheme_owner", "scheme_count", "visible",
		}, []driver.Value{nil, nil, nil, nil, nil, int64(0), nil, nil, nil, nil, int64(0), true}), nil
	case strings.Contains(normalized, "sys.extended_properties"):
		if len(args) != 3 || args[0].Value != sqlServerManagedRepairOperationPropertyName {
			return nil, errors.New("repair operation marker lookup was not bound to the inspected table")
		}
		if _, ok := args[1].Value.(string); !ok {
			return nil, errors.New("repair operation marker schema was not bound as text")
		}
		if _, ok := args[2].Value.(string); !ok {
			return nil, errors.New("repair operation marker table was not bound as text")
		}
		if c.script.repairOperationID == "" {
			return sqlServerRecordingRows([]string{"repair_operation_id"}), nil
		}
		return sqlServerRecordingRows([]string{"repair_operation_id"}, []driver.Value{c.script.repairOperationID}), nil
	case strings.Contains(normalized, "sys.objects"):
		if c.script.missingObject {
			return sqlServerRecordingRows([]string{"object facts"}), nil
		}
		return sqlServerRecordingRows(sqlServerManagedObjectColumns(), c.script.objectRows...), nil
	case strings.Contains(normalized, "sys.indexes"):
		return sqlServerRecordingRows(sqlServerManagedIndexColumns(), c.script.indexRows...), nil
	case strings.Contains(normalized, "sys.columns"):
		return sqlServerRecordingRows(sqlServerManagedColumnColumns(), c.script.columnRows...), nil
	case strings.Contains(normalized, "sys.check_constraints"):
		return sqlServerRecordingRows(sqlServerManagedCheckColumns(), c.script.checkRows...), nil
	default:
		return nil, errors.New("unexpected SQL Server catalog query")
	}
}

type sqlServerManagedCatalogTx struct{}

func (sqlServerManagedCatalogTx) Commit() error   { return nil }
func (sqlServerManagedCatalogTx) Rollback() error { return nil }

type sqlServerManagedCatalogTestSession struct {
	db     *sql.DB
	script *sqlServerManagedCatalogScript
}

func (s *sqlServerManagedCatalogTestSession) DB() *sql.DB { return s.db }

func (s *sqlServerManagedCatalogTestSession) Close() error {
	s.script.closed = true
	return s.db.Close()
}

func openSQLServerManagedCatalogTestSession(t *testing.T, script *sqlServerManagedCatalogScript) *sqlServerManagedCatalogTestSession {
	t.Helper()
	db := sql.OpenDB(sqlServerManagedCatalogConnector{script: script})
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	return &sqlServerManagedCatalogTestSession{db: db, script: script}
}

func sqlServerManagedCatalogFixture() *sqlServerManagedCatalogScript {
	metadata := sqlServerRecordingMetadataFixture()
	return &sqlServerManagedCatalogScript{
		metadata: metadata,
		objectRows: [][]driver.Value{{
			int64(91), "dbo", "circuit_long", "USER_TABLE", true, false, int64(0), false, false, false, false,
		}},
		columnRows: [][]driver.Value{
			{int64(1), "bucket_start", "datetime2", "sys", false, int64(8), int64(27), int64(6), false, false, false, int64(0), false, false, false, nil, nil},
			{int64(2), "record_id", "char", "sys", false, int64(64), int64(0), int64(0), false, false, false, int64(0), false, false, false, "Latin1_General_100_BIN2", nil},
			{int64(3), "workspace_id", "nvarchar", "sys", false, int64(382), int64(0), int64(0), false, false, false, int64(0), false, false, false, "Latin1_General_100_BIN2", nil},
			{int64(4), "tag_id", "nvarchar", "sys", false, int64(382), int64(0), int64(0), false, false, false, int64(0), false, false, false, "Latin1_General_100_BIN2", nil},
			{int64(5), "value_decimal", "decimal", "sys", false, int64(17), int64(38), int64(18), true, false, false, int64(0), false, false, false, nil, nil},
			{int64(6), "value_text", "nvarchar", "sys", false, int64(-1), int64(0), int64(0), true, false, false, int64(0), false, false, false, "Latin1_General_100_BIN2", nil},
			{int64(7), "provenance", "nvarchar", "sys", false, int64(-1), int64(0), int64(0), false, false, false, int64(0), false, false, false, "Latin1_General_100_BIN2", nil},
			{int64(8), "_gw_owner_0123456789abcdef01234567", "nvarchar", "sys", false, int64(32), int64(0), int64(0), false, false, false, int64(0), false, false, false, "Latin1_General_100_BIN2", "((N'gateway'))"},
		},
		indexRows: [][]driver.Value{
			{int64(1), "PK_circuit_long", int64(1), true, true, false, false, false, false, "FG", int64(1), int64(1), false, false, "bucket_start", false},
			{int64(1), "PK_circuit_long", int64(1), true, true, false, false, false, false, "FG", int64(2), int64(2), false, false, "record_id", false},
			{int64(2), "idx_circuit_long_src_time", int64(2), false, false, false, false, false, false, "FG", int64(1), int64(1), false, false, "workspace_id", false},
			{int64(2), "idx_circuit_long_src_time", int64(2), false, false, false, false, false, false, "FG", int64(2), int64(2), false, false, "tag_id", false},
			{int64(2), "idx_circuit_long_src_time", int64(2), false, false, false, false, false, false, "FG", int64(3), int64(3), false, false, "bucket_start", false},
		},
		checkRows: [][]driver.Value{{"CK_circuit_long_provenance", "([ISJSON]([provenance])=(1))", false, false, int64(0)}},
	}
}

func sqlServerManagedObjectColumns() []string {
	return []string{"object_id", "schema_name", "object_name", "type_desc", "is_user_table", "is_memory_optimized", "temporal_type", "is_filetable", "is_external", "is_node", "is_edge"}
}

func sqlServerManagedColumnColumns() []string {
	return []string{"column_id", "name", "type_name", "type_schema", "is_user_defined", "max_length", "precision", "scale", "is_nullable", "is_computed", "is_identity", "generated_always_type", "is_hidden", "is_sparse", "is_rowguidcol", "collation_name", "default_definition"}
}

func sqlServerManagedIndexColumns() []string {
	return []string{"index_id", "name", "type", "is_unique", "is_primary_key", "has_filter", "is_disabled", "ignore_dup_key", "is_hypothetical", "data_space_type", "key_ordinal", "index_column_id", "is_included_column", "is_descending_key", "column_name", "is_computed"}
}

func sqlServerManagedCheckColumns() []string {
	return []string{"name", "definition", "is_disabled", "is_not_trusted", "parent_column_id"}
}

func TestSQLServerManagedInspectionReportsActualCatalogFacts(t *testing.T) {
	internalDB := openMigratedTestDB(t)
	repo := NewSQLConnectorRepository(internalDB)
	service := NewConnectorService(repo)
	config := sqlServerConnectorConfigForTest()
	connector := createSQLServerRecordingSavedConnector(t, repo, "sqlserver-managed-catalog", schema.DatabaseConnectorKindSQLServer, config, true)
	script := sqlServerManagedCatalogFixture()
	session := openSQLServerManagedCatalogTestSession(t, script)

	originalSessionOpen := openSQLServerManagedSessionFunc
	originalManagerOpen := openExternalDBManagerFunc
	openCalls, managerOpenCalls := 0, 0
	var openedKind schema.DatabaseConnectorKind
	var openedConfig ConnectionConfig
	openSQLServerManagedSessionFunc = func(kind schema.DatabaseConnectorKind, got ConnectionConfig) (sqlServerManagedSession, error) {
		openCalls++
		openedKind, openedConfig = kind, got
		return session, nil
	}
	openExternalDBManagerFunc = func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
		managerOpenCalls++
		return nil, errors.New("unexpected generic open")
	}
	t.Cleanup(func() {
		openSQLServerManagedSessionFunc = originalSessionOpen
		openExternalDBManagerFunc = originalManagerOpen
	})

	result, err := serviceManagedInspector(service, internalDB).InspectSQLServerManagedTableAtRevision(
		t.Context(), connector.ID, connector.IdentityRevision, "fixture_db", "dbo",
		sqlServerManagedExpectedCapabilityRevision(connector), "circuit_long",
	)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, TableInspectionExists, result.Inspection.Status)
	require.Equal(t, "dbo", result.Inspection.Schema)
	require.Equal(t, "circuit_long", result.Inspection.Table)
	require.True(t, result.Inspection.IndexesInspected)
	require.Equal(t, 1, openCalls)
	require.Zero(t, managerOpenCalls, "managed catalog inspection must not use generic probe/open paths")
	require.Equal(t, schema.DatabaseConnectorKindSQLServer, openedKind)
	require.Equal(t, "fixture_db", stringConfigValue(openedConfig, "database"))
	require.Equal(t, "dbo", stringConfigValue(openedConfig, "schema"))
	require.True(t, script.queryHasDeadline)
	require.Zero(t, script.beginCalls, "catalog reads must not start a read-only transaction")
	require.Zero(t, script.execCalls, "catalog inspection must not execute DDL or DML")
	require.True(t, session.script.closed)

	require.Equal(t, "USER_TABLE", result.Evidence.TableType)
	require.Equal(t, int64(91), result.Evidence.ObjectID)
	require.True(t, result.Evidence.MetadataVerified)
	require.True(t, result.Evidence.Plain)
	require.Equal(t, []string{"bucket_start", "record_id"}, result.Evidence.PrimaryKey)
	require.Equal(t, "Latin1_General_100_BIN2", result.Evidence.ColumnCollations["record_id"])
	require.Equal(t, "gateway", *result.Evidence.ColumnDefaults["_gw_owner_0123456789abcdef01234567"])
	require.True(t, result.Evidence.ProvenanceJSONChecked)
	require.Equal(t, []string{"DATETIME2(6)", "CHAR(64)", "NVARCHAR(191)", "NVARCHAR(191)", "DECIMAL(38,18)", "NVARCHAR(MAX)", "NVARCHAR(MAX)", "NVARCHAR(16)"}, sqlServerManagedColumnTypes(result.Inspection.Columns))
	require.Equal(t, []IndexInfo{{Name: "idx_circuit_long_src_time", Columns: []string{"workspace_id", "tag_id", "bucket_start"}}}, result.Indexes)
	require.Equal(t, result.Indexes, result.Inspection.Indexes)

	objectCall := sqlServerManagedQueryContaining(t, script, "sys.objects")
	require.Contains(t, objectCall.query, "@p1")
	require.Contains(t, objectCall.query, "@p2")
	require.Equal(t, []driver.Value{"dbo", "circuit_long"}, sqlServerManagedQueryValues(objectCall.args))
	require.NotContains(t, objectCall.query, "circuit_long", "table names must remain query parameters")
	require.NotContains(t, objectCall.query, "CREATE")
	require.NotContains(t, objectCall.query, "ALTER")
	require.NotContains(t, objectCall.query, "DROP")
	require.NotContains(t, objectCall.query, "EXEC")
}

func serviceManagedInspector(service *ConnectorService, internalDB *sql.DB) *ManagedTableInspector {
	return NewManagedTableInspector(service, internalDB)
}

func sqlServerManagedExpectedCapabilityRevision(connector *schema.DatabaseConnector) string {
	metadata := sqlServerRecordingMetadata{
		serverVersion: "16.0.1000.6", serverEdition: "Developer Edition (64-bit)",
		engineEdition: 3, engineEditionValid: true, compatibilityLevel: 160, compatibilityValid: true,
		database: "fixture_db", actualDatabase: "fixture_db", schema: "dbo", actualSchema: "dbo",
		metadataStatus: sqlServerRecordingMetadataVerified,
		permissions: SQLServerRecordingPermissions{
			Verified: true, DatabaseViewDefinition: true, SchemaViewDefinition: true,
			CreateTable: true, AlterSchema: true, SelectSchema: true, InsertSchema: true,
			DeleteSchema: true, AlterAnyDataspace: true,
		},
	}
	return sqlServerRecordingCapabilityRevision(connector, metadata, true)
}

func sqlServerManagedColumnTypes(columns []ColumnInfo) []string {
	result := make([]string, len(columns))
	for index, column := range columns {
		result[index] = column.DataType
	}
	return result
}

func sqlServerManagedQueryContaining(t *testing.T, script *sqlServerManagedCatalogScript, fragment string) sqlServerManagedQueryCall {
	t.Helper()
	for _, call := range script.queries {
		if strings.Contains(call.query, fragment) {
			return call
		}
	}
	t.Fatalf("query containing %q was not executed", fragment)
	return sqlServerManagedQueryCall{}
}

func sqlServerManagedQueryValues(args []driver.NamedValue) []driver.Value {
	values := make([]driver.Value, len(args))
	for index, arg := range args {
		values[index] = arg.Value
	}
	return values
}

var _ driver.QueryerContext = (*sqlServerManagedCatalogConn)(nil)
var _ driver.ExecerContext = (*sqlServerManagedCatalogConn)(nil)
var _ driver.ConnBeginTx = (*sqlServerManagedCatalogConn)(nil)
var _ driver.Conn = (*sqlServerManagedCatalogConn)(nil)
var _ driver.Tx = sqlServerManagedCatalogTx{}
