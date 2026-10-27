package dbtarget

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestSQLServerPartitionRollover(t *testing.T) {
	db, capability, before, script := sqlServerNativePartitionMaintenanceFixture(t)
	after := sqlServerNativePartitionAfterTwoMonths(before.Evidence.NativeMonthly)
	require.False(t, after.FunctionExists)
	require.False(t, after.SchemeExists)
	require.False(t, after.Exclusive)
	require.False(t, after.IndexesAligned)
	statements := []string{
		"ALTER PARTITION SCHEME [ps_circuit_long] NEXT USED [PRIMARY]",
		"ALTER PARTITION FUNCTION [pf_circuit_long]() SPLIT RANGE ('2027-03-01T00:00:00.000000')",
		"ALTER PARTITION SCHEME [ps_circuit_long] NEXT USED [PRIMARY]",
		"ALTER PARTITION FUNCTION [pf_circuit_long]() SPLIT RANGE ('2027-04-01T00:00:00.000000')",
	}
	script.expectedStatements = slices.Clone(statements)
	script.splitBoundaries = []string{"2027-03-01T00:00:00.000000", "2027-04-01T00:00:00.000000"}
	script.resetTrace()

	err := SplitOpenSQLServerNativeMonthly(t.Context(), db, capability, before, after)
	require.NoError(t, err)
	require.Equal(t, statements, script.execStatements)
	require.Equal(t, []int{1, 1, 1, 1}, script.execTransactionIDs)
	require.Len(t, script.native.bounds, 7)
	require.Equal(t, 1, script.beginCalls)
	require.Equal(t, 1, script.commitCalls)
	require.Zero(t, script.rollbackCalls)
	for _, transactionID := range script.queryTransactionIDs {
		require.Equal(t, 1, transactionID, "before and after catalog facts must be checked in the DDL transaction")
	}
}

func TestSQLServerPartitionDowntime(t *testing.T) {
	db, capability, before, script := sqlServerNativePartitionMaintenanceFixture(t)
	additions := []SQLServerPartitionBoundary{
		{Name: "p_gap_202810", LessThan: "2028-10-01T00:00:00.000000", Ordinal: 6},
		{Name: "p_202810", LessThan: "2028-11-01T00:00:00.000000", Ordinal: 7},
		{Name: "p_202811", LessThan: "2028-12-01T00:00:00.000000", Ordinal: 8},
		{Name: "p_202812", LessThan: "2029-01-01T00:00:00.000000", Ordinal: 9},
		{Name: "p_202901", LessThan: "2029-02-01T00:00:00.000000", Ordinal: 10},
	}
	after, statements, boundaries := sqlServerNativePartitionMaintenancePlan(before.Evidence.NativeMonthly, additions)
	script.expectedStatements = statements
	script.splitBoundaries = boundaries
	script.resetTrace()

	require.NoError(t, SplitOpenSQLServerNativeMonthly(t.Context(), db, capability, before, after))
	require.Len(t, script.execStatements, 2*sqlServerNativePartitionMaxSplits)
	require.Len(t, script.native.bounds, 10)
	gapCount := 0
	for _, boundary := range after.Partitions {
		if strings.HasPrefix(boundary.Name, "p_gap_") {
			gapCount++
		}
	}
	require.Equal(t, 1, gapCount, "long downtime advances with one gap, not one split per missed month")
}

func TestSQLServerNativePartitionProofAcceptsExtendedGapProjection(t *testing.T) {
	req, _, catalog, db := nativeSQLServerInsertFixture(t)
	before := req.SQLServer.Group.Evidence.NativeMonthly
	after := *before
	after.Partitions = slices.Clone(before.Partitions[:len(before.Partitions)-1])
	after.Partitions = append(after.Partitions,
		SQLServerPartitionBoundary{Name: "p_gap_202608", LessThan: "2026-08-01T00:00:00.000000", Ordinal: 6},
		SQLServerPartitionBoundary{Name: "p_202608", LessThan: "2026-09-01T00:00:00.000000", Ordinal: 7},
		SQLServerPartitionBoundary{Name: "p_future", LessThan: "MAXVALUE", Ordinal: 8},
	)
	req.SQLServer.Group.Evidence.NativeMonthly = &after

	catalog.bounds = nil
	for i, value := range []string{"2026-01-01", "2026-02-01", "2026-03-01", "2026-04-01", "2026-05-01", "2026-08-01", "2026-09-01"} {
		at, err := time.Parse("2006-01-02", value)
		require.NoError(t, err)
		catalog.bounds = append(catalog.bounds, []driver.Value{int64(i + 1), int64(1), at})
	}
	catalog.destinations = append(catalog.destinations,
		[]driver.Value{int64(7), "PRIMARY", "FG"},
		[]driver.Value{int64(8), "PRIMARY", "FG"},
	)

	result, err := InsertGroupRow(t.Context(), db, req)
	require.NoError(t, err)
	require.False(t, result.AlreadyCommitted)
}

func sqlServerNativePartitionMaintenancePlan(
	before *SQLServerManagedPartitionEvidence,
	additions []SQLServerPartitionBoundary,
) (planned *SQLServerManagedPartitionEvidence, statements, boundaries []string) {
	after := *before
	after.Partitions = slices.Clone(before.Partitions[:len(before.Partitions)-1])
	after.Partitions = append(after.Partitions, additions...)
	after.Partitions = append(after.Partitions, SQLServerPartitionBoundary{
		Name: "p_future", LessThan: "MAXVALUE", Ordinal: len(after.Partitions) + 1,
	})
	after.FunctionExists, after.SchemeExists, after.Exclusive, after.IndexesAligned = false, false, false, false
	statements = make([]string, 0, len(additions)*2)
	boundaries = make([]string, len(additions))
	for index, boundary := range additions {
		statements = append(statements,
			"ALTER PARTITION SCHEME [ps_circuit_long] NEXT USED [PRIMARY]",
			fmt.Sprintf("ALTER PARTITION FUNCTION [pf_circuit_long]() SPLIT RANGE ('%s')", boundary.LessThan),
		)
		boundaries[index] = boundary.LessThan
	}
	return &after, statements, boundaries
}

func sqlServerNativePartitionAfterTwoMonths(before *SQLServerManagedPartitionEvidence) *SQLServerManagedPartitionEvidence {
	after := *before
	after.Partitions = slices.Clone(before.Partitions[:len(before.Partitions)-1])
	after.Partitions = append(after.Partitions,
		SQLServerPartitionBoundary{Name: "p_202702", LessThan: "2027-03-01T00:00:00.000000", Ordinal: 6},
		SQLServerPartitionBoundary{Name: "p_202703", LessThan: "2027-04-01T00:00:00.000000", Ordinal: 7},
		SQLServerPartitionBoundary{Name: "p_future", LessThan: "MAXVALUE", Ordinal: 8},
	)
	after.FunctionExists, after.SchemeExists, after.Exclusive, after.IndexesAligned = false, false, false, false
	return &after
}

type sqlServerNativePartitionMaintenanceScript struct {
	mu                          sync.Mutex
	gate                        chan struct{}
	native                      *nativeSQLServerCatalogScript
	expectedStatements          []string
	splitBoundaries             []string
	failExecAt                  int
	execError                   error
	commitError                 error
	applyCommitOnError          bool
	rollbackError               error
	applyOnRollbackError        bool
	foreignOwnerAfterFinalSplit bool
	execAttempts                int
	beginCalls                  int
	commitCalls                 int
	rollbackCalls               int
	execStatements              []string
	execTransactionIDs          []int
	queryTransactionIDs         []int
	nextTransactionID           int
	splitCount                  int
}

type sqlServerNativePartitionMaintenanceConnector struct {
	script *sqlServerNativePartitionMaintenanceScript
}

type sqlServerNativePartitionMaintenanceConn struct {
	script *sqlServerNativePartitionMaintenanceScript
	tx     *sqlServerNativePartitionMaintenanceTx
}

type sqlServerNativePartitionMaintenanceTx struct {
	conn       *sqlServerNativePartitionMaintenanceConn
	id         int
	bounds     [][]driver.Value
	execIndex  int
	splitCount int
	armed      bool
	closed     bool
}

var (
	_ driver.Connector      = sqlServerNativePartitionMaintenanceConnector{}
	_ driver.Conn           = (*sqlServerNativePartitionMaintenanceConn)(nil)
	_ driver.ConnBeginTx    = (*sqlServerNativePartitionMaintenanceConn)(nil)
	_ driver.QueryerContext = (*sqlServerNativePartitionMaintenanceConn)(nil)
	_ driver.ExecerContext  = (*sqlServerNativePartitionMaintenanceConn)(nil)
)

func sqlServerNativePartitionMaintenanceFixture(
	t *testing.T,
) (*sql.DB, *SQLServerRecordingCapability, *SQLServerManagedInspection, *sqlServerNativePartitionMaintenanceScript) {
	t.Helper()
	native := nativeSQLServerCatalogFixture()
	group := sqlServerInsertTestInspection("circuit_long", false, 91)
	group.Evidence.Plain = false
	native.base = sqlServerInsertCatalogScript(group)
	for _, row := range native.base.indexRows {
		row[9] = "PS"
	}
	script := &sqlServerNativePartitionMaintenanceScript{native: native, gate: make(chan struct{}, 1)}
	script.gate <- struct{}{}
	db := sql.OpenDB(sqlServerNativePartitionMaintenanceConnector{script: script})
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	connector := &schema.DatabaseConnector{
		ID: "sqlserver-native-maintenance-test", Kind: schema.DatabaseConnectorKindSQLServer,
		IdentityRevision: "revision-1",
	}
	metadata := inspectSQLServerRecordingMetadata(t.Context(), db, "fixture_db", "dbo")
	capability := projectSQLServerRecordingCapability(connector, metadata, true)
	capability.SupportsManagedLong = true
	capability.SupportsNativePartition = true
	before, err := InspectOpenSQLServerManagedTable(t.Context(), db, capability, "circuit_long")
	require.NoError(t, err)
	require.NotNil(t, before.Evidence.NativeMonthly)
	return db, capability, before, script
}

func (script *sqlServerNativePartitionMaintenanceScript) resetTrace() {
	script.mu.Lock()
	defer script.mu.Unlock()
	script.execAttempts = 0
	script.beginCalls = 0
	script.commitCalls = 0
	script.rollbackCalls = 0
	script.execStatements = nil
	script.execTransactionIDs = nil
	script.queryTransactionIDs = nil
	script.splitCount = 0
	script.native.nativeReads = 0
}

func (connector sqlServerNativePartitionMaintenanceConnector) Connect(context.Context) (driver.Conn, error) {
	return &sqlServerNativePartitionMaintenanceConn{script: connector.script}, nil
}

func (sqlServerNativePartitionMaintenanceConnector) Driver() driver.Driver {
	return sqlServerManagedCatalogDriver{}
}

func (*sqlServerNativePartitionMaintenanceConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("native maintenance statements must use context-aware driver calls")
}

func (*sqlServerNativePartitionMaintenanceConn) Close() error { return nil }

func (conn *sqlServerNativePartitionMaintenanceConn) Begin() (driver.Tx, error) {
	return conn.BeginTx(context.Background(), driver.TxOptions{})
}

func (conn *sqlServerNativePartitionMaintenanceConn) BeginTx(ctx context.Context, _ driver.TxOptions) (driver.Tx, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-conn.script.gate:
	}
	conn.script.mu.Lock()
	conn.script.nextTransactionID++
	conn.script.beginCalls++
	tx := &sqlServerNativePartitionMaintenanceTx{
		conn: conn, id: conn.script.nextTransactionID, bounds: cloneSQLServerNativeBoundaryRows(conn.script.native.bounds),
	}
	conn.script.mu.Unlock()
	conn.tx = tx
	return tx, nil
}

func (conn *sqlServerNativePartitionMaintenanceConn) QueryContext(
	ctx context.Context,
	query string,
	args []driver.NamedValue,
) (driver.Rows, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	txID := 0
	if conn.tx != nil {
		txID = conn.tx.id
	}
	conn.script.mu.Lock()
	conn.script.queryTransactionIDs = append(conn.script.queryTransactionIDs, txID)
	conn.script.mu.Unlock()
	switch {
	case strings.Contains(query, "gw_native_partition_bounds"):
		bounds := conn.currentBounds()
		return sqlServerRecordingRows([]string{"ordinal", "parameter", "boundary"}, bounds...), nil
	case strings.Contains(query, "gw_native_partition_destinations"):
		count := len(conn.currentBounds()) + 1
		rows := make([][]driver.Value, count)
		for index := range count {
			rows[index] = []driver.Value{int64(index + 1), "PRIMARY", "FG"}
		}
		return sqlServerRecordingRows([]string{"ordinal", "name", "kind"}, rows...), nil
	case strings.Contains(query, "gw_native_partition_"):
		nativeConn := &nativeSQLServerCatalogConn{
			sqlServerManagedCatalogConn: sqlServerManagedCatalogConn{script: conn.script.native.base},
			native:                      conn.script.native,
		}
		return nativeConn.QueryContext(ctx, query, args)
	default:
		return (&sqlServerManagedCatalogConn{script: conn.script.native.base}).QueryContext(ctx, query, args)
	}
}

func (conn *sqlServerNativePartitionMaintenanceConn) currentBounds() [][]driver.Value {
	if conn.tx != nil {
		return cloneSQLServerNativeBoundaryRows(conn.tx.bounds)
	}
	conn.script.mu.Lock()
	defer conn.script.mu.Unlock()
	return cloneSQLServerNativeBoundaryRows(conn.script.native.bounds)
}

func (conn *sqlServerNativePartitionMaintenanceConn) ExecContext(
	ctx context.Context,
	query string,
	_ []driver.NamedValue,
) (driver.Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if conn.tx == nil {
		return nil, errors.New("partition DDL escaped its transaction")
	}
	script, tx := conn.script, conn.tx
	script.mu.Lock()
	script.execAttempts++
	script.execStatements = append(script.execStatements, strings.Join(strings.Fields(query), " "))
	script.execTransactionIDs = append(script.execTransactionIDs, tx.id)
	fail := script.failExecAt == script.execAttempts
	execError := script.execError
	script.mu.Unlock()
	if fail {
		return nil, execError
	}
	if tx.execIndex >= len(script.expectedStatements) || query != script.expectedStatements[tx.execIndex] {
		return nil, errors.New("unexpected native partition DDL statement")
	}
	tx.execIndex++
	if strings.HasPrefix(query, "ALTER PARTITION SCHEME ") {
		if tx.armed {
			return nil, errors.New("NEXT USED was repeated before SPLIT")
		}
		tx.armed = true
		return driver.RowsAffected(0), nil
	}
	if !strings.HasPrefix(query, "ALTER PARTITION FUNCTION ") || !strings.Contains(query, " SPLIT RANGE (") || !tx.armed {
		return nil, errors.New("only NEXT USED followed by SPLIT RANGE is allowed")
	}
	if tx.splitCount >= len(script.splitBoundaries) {
		return nil, errors.New("unexpected split boundary")
	}
	at, err := time.Parse(sqlServerNativeBoundaryTimeLayout, script.splitBoundaries[tx.splitCount])
	if err != nil {
		return nil, err
	}
	tx.bounds = append(tx.bounds, []driver.Value{int64(len(tx.bounds) + 1), int64(1), at})
	tx.splitCount++
	tx.armed = false
	if script.foreignOwnerAfterFinalSplit && tx.splitCount == len(script.splitBoundaries) {
		script.mu.Lock()
		script.native.objects[8] = "dbo.foreign"
		script.mu.Unlock()
	}
	return driver.RowsAffected(0), nil
}

func (tx *sqlServerNativePartitionMaintenanceTx) Commit() error {
	if tx.closed {
		return sql.ErrTxDone
	}
	tx.closed = true
	script := tx.conn.script
	script.mu.Lock()
	script.commitCalls++
	if script.commitError == nil || script.applyCommitOnError {
		script.native.bounds = cloneSQLServerNativeBoundaryRows(tx.bounds)
	}
	err := script.commitError
	script.mu.Unlock()
	tx.conn.tx = nil
	script.gate <- struct{}{}
	return err
}

func (tx *sqlServerNativePartitionMaintenanceTx) Rollback() error {
	if tx.closed {
		return sql.ErrTxDone
	}
	tx.closed = true
	script := tx.conn.script
	script.mu.Lock()
	script.rollbackCalls++
	if script.rollbackError != nil && script.applyOnRollbackError {
		script.native.bounds = cloneSQLServerNativeBoundaryRows(tx.bounds)
	}
	err := script.rollbackError
	script.mu.Unlock()
	tx.conn.tx = nil
	script.gate <- struct{}{}
	return err
}

func cloneSQLServerNativeBoundaryRows(rows [][]driver.Value) [][]driver.Value {
	cloned := make([][]driver.Value, len(rows))
	for index, row := range rows {
		cloned[index] = slices.Clone(row)
	}
	return cloned
}
