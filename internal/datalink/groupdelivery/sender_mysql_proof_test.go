package groupdelivery

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestSenderPassesMySQLFrozenProofAndMicrosecondReceiptTime(t *testing.T) {
	fixture := newDeliveryFixture(t, plainTargetTable, true)
	destination := testDestination
	destination.Database, destination.TableSchema, destination.TableName = "gwtest", "gwtest", "readings"
	destination.DedupeCapability = DedupeReceipt
	bucket := rowBucket("", t0, 7)
	require.NoError(t, fixture.store.CommitClosure(t.Context(), testKey, Closure{
		Destination: destination, NextClose: t0.Add(10 * time.Second), Buckets: []ClosedBucket{bucket},
	}))

	var queries atomic.Int32
	targetDB := sql.OpenDB(mysqlProofQueryConnector{queries: &queries})
	t.Cleanup(func() { _ = targetDB.Close() })
	fixture.resolver.target = Target{
		DB: targetDB, Kind: schema.DatabaseConnectorKindMySQL,
		MySQL: senderTestMySQLProof(),
	}

	result, err := fixture.sender.Deliver(t.Context(), bucket.Outcome.EffectKey)

	require.NoError(t, err)
	require.EqualValues(t, 1, queries.Load(), "a frozen proof and a MySQL DATETIME(6) receipt timestamp must reach the insert transaction")
	require.NotEqual(t, codeUnsupported, result.Code)
}

func senderTestMySQLProof() *dbtarget.MySQLManagedRowProof {
	groupEvidence := senderTestMySQLEvidence([]string{"record_id"})
	receiptEvidence := senderTestMySQLEvidence([]string{"effect_key"})
	groupIndexes := []dbtarget.IndexInfo{}
	receiptIndexes := []dbtarget.IndexInfo{}
	return &dbtarget.MySQLManagedRowProof{
		ConnectorID: "connector-1", ConnectorRevision: "crev-1",
		Group: &dbtarget.MySQLManagedInspection{
			Inspection: &dbtarget.TableInspection{
				Status: dbtarget.TableInspectionExists, Schema: "gwtest", Table: "readings",
				Columns: []dbtarget.ColumnInfo{{Name: "record_id", DataType: "char(64)"}},
				Indexes: groupIndexes, IndexesInspected: true,
			},
			Evidence: groupEvidence, Indexes: groupIndexes,
		},
		Receipt: &dbtarget.MySQLManagedInspection{
			Inspection: &dbtarget.TableInspection{
				Status: dbtarget.TableInspectionExists, Schema: "gwtest", Table: dbtarget.EffectReceiptTable,
				Columns: []dbtarget.ColumnInfo{{Name: "effect_key", DataType: "varchar(191)"}},
				Indexes: receiptIndexes, IndexesInspected: true,
			},
			Evidence: receiptEvidence, Indexes: receiptIndexes,
		},
	}
}

func senderTestMySQLEvidence(primaryKey []string) *dbtarget.MySQLManagedTableEvidence {
	return &dbtarget.MySQLManagedTableEvidence{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1", Database: "gwtest",
		Engine: "InnoDB", TableType: "BASE TABLE", PrimaryKey: primaryKey,
	}
}

type mysqlProofQueryConnector struct{ queries *atomic.Int32 }

func (c mysqlProofQueryConnector) Connect(context.Context) (driver.Conn, error) {
	return mysqlProofQueryConn(c), nil
}

func (mysqlProofQueryConnector) Driver() driver.Driver { return mysqlProofQueryDriver{} }

type mysqlProofQueryDriver struct{}

func (mysqlProofQueryDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("database/sql must use the supplied connector")
}

type mysqlProofQueryConn struct{ queries *atomic.Int32 }

func (mysqlProofQueryConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unexpected prepared statement")
}

func (mysqlProofQueryConn) Close() error { return nil }

func (mysqlProofQueryConn) Begin() (driver.Tx, error) { return mysqlProofQueryTx{}, nil }

func (mysqlProofQueryConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return mysqlProofQueryTx{}, nil
}

func (c mysqlProofQueryConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	c.queries.Add(1)
	return nil, errors.New("stop after the MySQL catalog read is reached")
}

type mysqlProofQueryTx struct{}

func (mysqlProofQueryTx) Commit() error   { return nil }
func (mysqlProofQueryTx) Rollback() error { return nil }
