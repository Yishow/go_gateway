package dbtarget

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	mssql "github.com/microsoft/go-mssqldb"
	"github.com/stretchr/testify/require"
)

func TestSQLServerSchemaErrorNumbersSelectOnlyShapeFailures(t *testing.T) {
	for _, test := range []struct {
		number       int32
		shapeFailure bool
		missingTable bool
	}{
		{number: 207, shapeFailure: true},
		{number: 208, shapeFailure: true, missingTable: true},
		{number: 229},
		{number: 1205},
	} {
		t.Run(fmt.Sprint(test.number), func(t *testing.T) {
			err := fmt.Errorf("SQL Server statement failed: %w", mssql.Error{Number: test.number})
			require.Equal(t, test.shapeFailure, IsSchemaShapeError(err))
			require.Equal(t, test.missingTable, IsMissingRelationError(err))
		})
	}
}

func TestSQLServerOpenManagedTableRequiresVisibleMatchingScope(t *testing.T) {
	t.Run("existing table without reopening", func(t *testing.T) {
		script := sqlServerManagedCatalogFixture()
		db, capability := sqlServerOpenManagedInspectionTestDB(t, script)
		openCalls := forbidSQLServerManagedInspectionOpens(t)

		result, err := InspectOpenSQLServerManagedTable(t.Context(), db, capability, "circuit_long")
		require.NoError(t, err)
		require.Equal(t, TableInspectionExists, result.Inspection.Status)
		require.True(t, result.Evidence.MetadataVerified)
		require.Zero(t, script.beginCalls)
		require.Zero(t, script.execCalls)
		require.Zero(t, *openCalls)
		require.True(t, sqlServerOpenManagedHasCatalogQuery(script))
	})

	t.Run("missing table in visible schema", func(t *testing.T) {
		script := sqlServerManagedCatalogFixture()
		script.missingObject = true
		db, capability := sqlServerOpenManagedInspectionTestDB(t, script)

		result, err := InspectOpenSQLServerManagedTable(t.Context(), db, capability, "circuit_long")
		require.NoError(t, err)
		require.Equal(t, TableInspectionMissing, result.Inspection.Status)
		require.True(t, result.Inspection.IndexesInspected)
		require.NotNil(t, result.Indexes)
		require.Empty(t, result.Indexes)
		require.True(t, result.Evidence.MetadataVerified)
	})

	t.Run("schema missing is not table missing", func(t *testing.T) {
		script := sqlServerManagedCatalogFixture()
		db, capability := sqlServerOpenManagedInspectionTestDB(t, script)
		script.metadata.visibleSchema = ""

		result, err := InspectOpenSQLServerManagedTable(t.Context(), db, capability, "circuit_long")
		require.ErrorIs(t, err, ErrSQLServerManagedTableUnconfirmed)
		require.Equal(t, TableInspectionFailed, result.Inspection.Status)
		require.Equal(t, sqlServerManagedSchemaMissingReason, result.Inspection.Reason)
		require.False(t, sqlServerOpenManagedHasCatalogQuery(script))
	})

	t.Run("hidden metadata is not table missing", func(t *testing.T) {
		script := sqlServerManagedCatalogFixture()
		db, capability := sqlServerOpenManagedInspectionTestDB(t, script)
		script.metadata.viewDefinition = false
		script.metadata.schemaViewDefinition = false

		result, err := InspectOpenSQLServerManagedTable(t.Context(), db, capability, "circuit_long")
		require.ErrorIs(t, err, ErrSQLServerManagedTableForbidden)
		require.Equal(t, TableInspectionForbidden, result.Inspection.Status)
		require.False(t, sqlServerOpenManagedHasCatalogQuery(script))
	})

	t.Run("revision change blocks catalog lookup", func(t *testing.T) {
		script := sqlServerManagedCatalogFixture()
		db, capability := sqlServerOpenManagedInspectionTestDB(t, script)
		capability.ConnectorRevision = "revision-changed"

		_, err := InspectOpenSQLServerManagedTable(t.Context(), db, capability, "circuit_long")
		require.ErrorIs(t, err, ErrSQLServerManagedScopeChanged)
		require.False(t, sqlServerOpenManagedHasCatalogQuery(script))
	})

	t.Run("foreign table occupies name", func(t *testing.T) {
		script := sqlServerManagedCatalogFixture()
		db, capability := sqlServerOpenManagedInspectionTestDB(t, script)
		script.objectRows[0][3] = "VIEW"
		script.objectRows[0][4] = false

		result, err := InspectOpenSQLServerManagedTable(t.Context(), db, capability, "circuit_long")
		require.NoError(t, err)
		require.Equal(t, TableInspectionExists, result.Inspection.Status)
		require.Equal(t, "VIEW", result.Evidence.TableType)
		require.False(t, result.Evidence.Plain)
	})

	t.Run("typed catalog permission denial is forbidden", func(t *testing.T) {
		script := sqlServerManagedCatalogFixture()
		db, capability := sqlServerOpenManagedInspectionTestDB(t, script)
		script.failContains = "sys.objects"
		script.queryError = mssql.Error{Number: 229}

		result, err := InspectOpenSQLServerManagedTable(t.Context(), db, capability, "circuit_long")
		require.ErrorIs(t, err, ErrSQLServerManagedTableForbidden)
		require.Equal(t, TableInspectionForbidden, result.Inspection.Status)
	})
}

func sqlServerOpenManagedHasCatalogQuery(script *sqlServerManagedCatalogScript) bool {
	for _, query := range script.queries {
		if strings.Contains(query.query, "sys.objects") {
			return true
		}
	}
	return false
}

func sqlServerOpenManagedInspectionTestDB(
	t *testing.T,
	script *sqlServerManagedCatalogScript,
) (*sql.DB, *SQLServerRecordingCapability) {
	t.Helper()
	session := openSQLServerManagedCatalogTestSession(t, script)
	connector := &schema.DatabaseConnector{
		ID: "sqlserver-open-inspection", Kind: schema.DatabaseConnectorKindSQLServer,
		IdentityRevision: "revision-1",
	}
	metadata := inspectSQLServerRecordingMetadata(t.Context(), session.db, "fixture_db", "dbo")
	return session.db, projectSQLServerRecordingCapability(connector, metadata, true)
}

func forbidSQLServerManagedInspectionOpens(t *testing.T) *int {
	t.Helper()
	openCalls := 0
	originalSessionOpen := openSQLServerManagedSessionFunc
	originalManagerOpen := openExternalDBManagerFunc
	openSQLServerManagedSessionFunc = func(schema.DatabaseConnectorKind, ConnectionConfig) (sqlServerManagedSession, error) {
		openCalls++
		return nil, errors.New("unexpected SQL Server inspection open")
	}
	openExternalDBManagerFunc = func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
		openCalls++
		return nil, errors.New("unexpected SQL Server destination open")
	}
	t.Cleanup(func() {
		openSQLServerManagedSessionFunc = originalSessionOpen
		openExternalDBManagerFunc = originalManagerOpen
	})
	return &openCalls
}

func TestSQLServerManagedRepairBatchRejectsClosedCapabilities(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*SQLServerRecordingCapability)
		want   error
	}{
		{name: "managed recording disabled", mutate: func(capability *SQLServerRecordingCapability) {
			capability.SupportsManagedLong = false
		}, want: ErrSchemaExecutionUnsupported},
		{name: "missing repair disabled", mutate: func(capability *SQLServerRecordingCapability) {
			capability.SupportsMissingRepair = false
		}, want: ErrSchemaExecutionUnsupported},
		{name: "scope unverified", mutate: func(capability *SQLServerRecordingCapability) {
			capability.MetadataStatus = sqlServerRecordingMetadataUnconfirmed
		}, want: ErrSQLServerManagedTableUnconfirmed},
	} {
		t.Run(test.name, func(t *testing.T) {
			db, capability, script := sqlServerRepairTestDB(t)
			test.mutate(capability)
			script.queries = nil
			request := sqlServerManagedRepairTestRequest([]string{"CREATE TABLE [dbo].[circuit_long] ([record_id] char(64))"})

			result, err := RunSQLServerManagedRepairBatch(t.Context(), db, capability, request)
			require.ErrorIs(t, err, test.want)
			require.True(t, result.RolledBack)
			require.Zero(t, script.queryCalls)
			require.Zero(t, script.beginCalls)
			require.Zero(t, script.execCalls)
		})
	}

	t.Run("receipt is not a repair target", func(t *testing.T) {
		db, capability, script := sqlServerRepairTestDB(t)
		request := sqlServerManagedRepairTestRequest([]string{"CREATE TABLE [dbo].[gw_effect_receipts_v1] ([id] char(64))"})
		request.TableName = EffectReceiptTable

		result, err := RunSQLServerManagedRepairBatch(t.Context(), db, capability, request)
		require.ErrorIs(t, err, ErrSchemaExecutionUnsupported)
		require.True(t, result.RolledBack)
		require.Zero(t, script.beginCalls)
		require.Zero(t, script.execCalls)
	})

	db, capability, script := sqlServerRepairTestDB(t)
	request := sqlServerManagedRepairTestRequest([]string{"CREATE TABLE [dbo].[circuit_long] ([record_id] char(64))"})
	result, err := RunSQLServerManagedRepairBatch(t.Context(), nil, capability, request)
	require.ErrorIs(t, err, ErrSQLServerManagedTableUnconfirmed)
	require.True(t, result.RolledBack)
	require.Zero(t, script.beginCalls)
	require.Zero(t, script.execCalls)
	_ = db
}

func TestSQLServerManagedRepairBatchRechecksRevisionAndPermissions(t *testing.T) {
	t.Run("stale capability revision", func(t *testing.T) {
		db, capability, script := sqlServerRepairTestDB(t)
		capability.CapabilityRevision = "sqlservercap-v1-stale"
		request := sqlServerManagedRepairTestRequest([]string{"CREATE TABLE [dbo].[circuit_long] ([record_id] char(64))"})

		result, err := RunSQLServerManagedRepairBatch(t.Context(), db, capability, request)
		require.ErrorIs(t, err, ErrSQLServerManagedScopeChanged)
		require.True(t, result.RolledBack)
		require.Zero(t, script.beginCalls)
		require.Zero(t, script.execCalls)
	})

	t.Run("live create permission revoked", func(t *testing.T) {
		db, capability, script := sqlServerRepairTestDB(t)
		script.metadata.createTable = false
		request := sqlServerManagedRepairTestRequest([]string{"CREATE TABLE [dbo].[circuit_long] ([record_id] char(64))"})

		result, err := RunSQLServerManagedRepairBatch(t.Context(), db, capability, request)
		require.ErrorIs(t, err, ErrSchemaExecutionPermissionDenied)
		require.True(t, result.RolledBack)
		require.Zero(t, script.beginCalls)
		require.Zero(t, script.execCalls)
	})
}

type sqlServerRepairScript struct {
	metadata                    *sqlServerRecordingMetadataScript
	queryCalls                  int
	queries                     []string
	txMetadata                  *sqlServerRecordingMetadataScript
	metadataQueryTransactionIDs []int
	beginCalls                  int
	txOptions                   []driver.TxOptions
	execCalls                   int
	commitCalls                 int
	rollbackCalls               int
	statements                  []string
	execStatements              []string
	execTransactionIDs          []int
	native                      *sqlServerNativeRepairTestState
	markerArgs                  []driver.NamedValue
	pendingRepairOperationID    string
	repairOperationID           string
	activeTransaction           int
	nextTransactionID           int
	failExecAt                  int
	execError                   error
	rollbackError               error
	commitError                 error
}

type sqlServerRepairConnector struct{ script *sqlServerRepairScript }
type sqlServerRepairDriver struct{}
type sqlServerRepairConn struct{ script *sqlServerRepairScript }
type sqlServerRepairTx struct {
	script *sqlServerRepairScript
	id     int
}

func (connector sqlServerRepairConnector) Connect(context.Context) (driver.Conn, error) {
	return &sqlServerRepairConn{script: connector.script}, nil
}

func (sqlServerRepairConnector) Driver() driver.Driver { return sqlServerRepairDriver{} }
func (sqlServerRepairDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("scripted SQL Server repair connector must be used")
}
func (*sqlServerRepairConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepared repair statements are not scripted")
}
func (*sqlServerRepairConn) Close() error { return nil }
func (conn *sqlServerRepairConn) Begin() (driver.Tx, error) {
	return conn.BeginTx(context.Background(), driver.TxOptions{})
}
func (conn *sqlServerRepairConn) BeginTx(ctx context.Context, options driver.TxOptions) (driver.Tx, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	conn.script.beginCalls++
	conn.script.txOptions = append(conn.script.txOptions, options)
	conn.script.nextTransactionID++
	conn.script.activeTransaction = conn.script.nextTransactionID
	return &sqlServerRepairTx{script: conn.script, id: conn.script.activeTransaction}, nil
}
func (conn *sqlServerRepairConn) QueryContext(
	ctx context.Context,
	query string,
	args []driver.NamedValue,
) (driver.Rows, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	conn.script.queryCalls++
	normalized := strings.Join(strings.Fields(query), " ")
	conn.script.queries = append(conn.script.queries, normalized)
	metadata := conn.script.metadata
	switch normalized {
	case sqlServerRecordingMetadataQuery, sqlServerRecordingVisibilityPermissionQuery, sqlServerRecordingSchemaQuery, sqlServerRecordingPermissionQuery:
		if conn.script.activeTransaction > 0 {
			conn.script.metadataQueryTransactionIDs = append(conn.script.metadataQueryTransactionIDs, conn.script.activeTransaction)
			if conn.script.txMetadata != nil {
				metadata = conn.script.txMetadata
			}
		}
	}
	if strings.Contains(normalized, "gw_native_partition_") || strings.Contains(normalized, "FROM sys.objects AS o") {
		if conn.script.native == nil {
			return nil, fmt.Errorf("unexpected SQL Server repair query: %s", normalized)
		}
		return conn.script.native.query(ctx, normalized, args, conn.script.activeTransaction)
	}
	switch normalized {
	case sqlServerRecordingMetadataQuery:
		return sqlServerRecordingRows(
			[]string{"version", "edition", "engine_edition", "database", "compatibility"},
			[]driver.Value{
				metadata.serverVersion, metadata.serverEdition,
				metadata.engineEdition, metadata.database,
				metadata.compatibilityLevel,
			},
		), nil
	case sqlServerRecordingVisibilityPermissionQuery:
		if len(args) != 1 || args[0].Value != metadata.expectedSchema {
			return nil, errors.New("unexpected SQL Server visibility scope")
		}
		return sqlServerRecordingRows(
			[]string{"database_view_definition", "schema_view_definition"},
			[]driver.Value{metadata.viewDefinition, metadata.schemaViewDefinition},
		), nil
	case sqlServerRecordingSchemaQuery:
		if len(args) != 1 || args[0].Value != metadata.expectedSchema {
			return nil, errors.New("unexpected SQL Server schema scope")
		}
		if metadata.visibleSchema == "" {
			return sqlServerRecordingRows([]string{"name"}), nil
		}
		return sqlServerRecordingRows([]string{"name"}, []driver.Value{metadata.visibleSchema}), nil
	case sqlServerRecordingPermissionQuery:
		if len(args) != 1 || args[0].Value != metadata.expectedSchema {
			return nil, errors.New("unexpected SQL Server permission scope")
		}
		return sqlServerRecordingRows(
			[]string{"create_table", "alter_schema", "select_schema", "insert_schema", "delete_schema", "alter_any_dataspace"},
			[]driver.Value{
				metadata.createTable, metadata.alterSchema, metadata.selectSchema,
				metadata.insertSchema, metadata.deleteSchema, metadata.alterAnyDataspace,
			},
		), nil
	default:
		return nil, fmt.Errorf("unexpected SQL Server repair query: %s", normalized)
	}
}
func (conn *sqlServerRepairConn) ExecContext(
	ctx context.Context,
	query string,
	args []driver.NamedValue,
) (driver.Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	conn.script.execCalls++
	conn.script.execStatements = append(conn.script.execStatements, query)
	conn.script.execTransactionIDs = append(conn.script.execTransactionIDs, conn.script.activeTransaction)
	if conn.script.activeTransaction == 0 {
		return nil, errors.New("SQL Server repair DDL outside transaction")
	}
	if conn.script.failExecAt == conn.script.execCalls {
		return nil, conn.script.execError
	}
	if query == sqlServerRepairMarkerAddStatement {
		if len(args) != 4 {
			return nil, errors.New("repair marker parameters were not bound")
		}
		conn.script.markerArgs = append([]driver.NamedValue(nil), args...)
		operationID, ok := args[1].Value.(string)
		if !ok || operationID == "" {
			return nil, errors.New("repair marker operation ID was not text")
		}
		conn.script.pendingRepairOperationID = operationID
	} else {
		conn.script.statements = append(conn.script.statements, query)
	}
	return driver.RowsAffected(1), nil
}
func (tx *sqlServerRepairTx) Commit() error {
	tx.script.commitCalls++
	tx.script.repairOperationID = tx.script.pendingRepairOperationID
	tx.script.pendingRepairOperationID = ""
	tx.script.activeTransaction = 0
	return tx.script.commitError
}
func (tx *sqlServerRepairTx) Rollback() error {
	tx.script.rollbackCalls++
	if tx.script.rollbackError != nil {
		return tx.script.rollbackError
	}
	tx.script.pendingRepairOperationID = ""
	tx.script.activeTransaction = 0
	return nil
}

func sqlServerManagedRepairTestRequest(statements []string) SQLServerManagedRepairRequest {
	return SQLServerManagedRepairRequest{
		TableName: "circuit_long", OperationID: "repair-operation-test", Statements: slices.Clone(statements),
	}
}

func sqlServerRepairTestDB(t *testing.T) (*sql.DB, *SQLServerRecordingCapability, *sqlServerRepairScript) {
	t.Helper()
	script := &sqlServerRepairScript{metadata: sqlServerRecordingMetadataFixture()}
	db := sql.OpenDB(sqlServerRepairConnector{script: script})
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	connector := &schema.DatabaseConnector{
		ID: "sqlserver-repair-test", Kind: schema.DatabaseConnectorKindSQLServer,
		IdentityRevision: "revision-1",
	}
	metadata := inspectSQLServerRecordingMetadata(t.Context(), db, "fixture_db", "dbo")
	capability := projectSQLServerRecordingCapability(connector, metadata, true)
	capability.SupportsManagedLong = true
	capability.SupportsMissingRepair = true
	script.queryCalls = 0
	script.queries = nil
	return db, capability, script
}
