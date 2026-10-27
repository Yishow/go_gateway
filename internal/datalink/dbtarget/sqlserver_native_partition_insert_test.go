package dbtarget

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type nativeSQLServerInsertConnector struct {
	script  *sqlServerGroupInsertScript
	catalog *nativeSQLServerCatalogScript
}

func (c nativeSQLServerInsertConnector) Connect(context.Context) (driver.Conn, error) {
	return &nativeSQLServerInsertConn{
		sqlServerGroupInsertConn: &sqlServerGroupInsertConn{script: c.script},
		native:                   &nativeSQLServerCatalogConn{native: c.catalog},
	}, nil
}

func (nativeSQLServerInsertConnector) Driver() driver.Driver { return sqlServerManagedCatalogDriver{} }

type nativeSQLServerInsertConn struct {
	*sqlServerGroupInsertConn
	native *nativeSQLServerCatalogConn
}

func (c *nativeSQLServerInsertConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if strings.Contains(query, "gw_native_partition_") {
		if c.tx == nil {
			return nil, errors.New("native catalog must be checked inside the row and receipt transaction")
		}
		return c.native.QueryContext(ctx, query, args)
	}
	return c.sqlServerGroupInsertConn.QueryContext(ctx, query, args)
}

func nativeSQLServerInsertFixture(t *testing.T) (GroupInsertRequest, *sqlServerGroupInsertScript, *nativeSQLServerCatalogScript, *sql.DB) {
	t.Helper()
	req := sqlServerGroupInsertTestRequest()
	req.SQLServer.Group.Evidence.Plain = false
	req.SQLServer.Group.Evidence.PartitionInspected = true
	req.SQLServer.Group.Evidence.NativeMonthly = sqlServerInsertTestNativeMonthlyEvidence(req.TableName)
	script := newSQLServerGroupInsertScript(req)
	for _, row := range script.groupCatalog.indexRows {
		row[9] = "PS"
	}
	catalog := nativeSQLServerCatalogFixture()
	catalog.table = req.TableName
	catalog.objects[8], catalog.objects[9] = "dbo."+req.TableName, "dbo."+req.TableName
	catalog.bounds = nil
	for i, value := range []string{"2026-01-01", "2026-02-01", "2026-03-01", "2026-04-01", "2026-05-01"} {
		at, err := time.Parse("2006-01-02", value)
		require.NoError(t, err)
		catalog.bounds = append(catalog.bounds, []driver.Value{int64(i + 1), int64(1), at})
	}
	for _, row := range catalog.indexes {
		row[0] = req.SQLServer.Group.Evidence.ObjectID
	}
	db := sql.OpenDB(nativeSQLServerInsertConnector{script, catalog})
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	return req, script, catalog, db
}

func TestSQLServerNativeInsertCommitsAndReplaysWithLivePartitionProof(t *testing.T) {
	req, script, catalog, db := nativeSQLServerInsertFixture(t)
	first, err := InsertGroupRow(t.Context(), db, req)
	require.NoError(t, err)
	require.False(t, first.AlreadyCommitted)
	replayed, err := InsertGroupRow(t.Context(), db, req)
	require.NoError(t, err)
	require.True(t, replayed.AlreadyCommitted)
	require.Equal(t, 1, script.commitCalls)
	require.Equal(t, 2, catalog.nativeReads, "even receipt replay must recheck the live partition proof")
	require.Len(t, script.groupInsertArgs, 1)
	require.Len(t, script.receiptInsertArgs, 1)
	require.True(t, script.rowExists)
	require.True(t, script.receiptExists)
}

func TestSQLServerNativeInsertBlocksPartitionDriftBeforeReceiptReplay(t *testing.T) {
	req, script, catalog, db := nativeSQLServerInsertFixture(t)
	_, err := InsertGroupRow(t.Context(), db, req)
	require.NoError(t, err)
	catalog.objects[8] = "dbo.foreign"
	_, err = InsertGroupRow(t.Context(), db, req)
	require.ErrorIs(t, err, ErrGroupInsertUnsupported)
	require.Equal(t, 1, script.commitCalls)
	require.Equal(t, 1, script.rollbackCalls)
	require.Len(t, script.groupInsertArgs, 1)
	require.Len(t, script.receiptInsertArgs, 1)
	require.True(t, script.rowExists)
	require.True(t, script.receiptExists)
}
