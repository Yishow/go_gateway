package point

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/sqltx"
	"modernc.org/sqlite"
	"path/filepath"
	"testing"
)

type pointLostAckConnector struct {
	path string
	fail bool
}

func (c *pointLostAckConnector) Driver() driver.Driver { return &sqlite.Driver{} }
func (c *pointLostAckConnector) Connect(context.Context) (driver.Conn, error) {
	conn, err := c.Driver().Open(c.path)
	if err != nil {
		return nil, err
	}
	return &pointLostAckConn{Conn: conn, owner: c}, nil
}

type pointLostAckConn struct {
	driver.Conn
	owner *pointLostAckConnector
}

func (c *pointLostAckConn) Begin() (driver.Tx, error) {
	tx, err := c.Conn.Begin()
	if err != nil {
		return nil, err
	}
	return &pointLostAckTx{Tx: tx, owner: c.owner}, nil
}

func (c *pointLostAckConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	tx, err := c.Conn.(driver.ConnBeginTx).BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &pointLostAckTx{Tx: tx, owner: c.owner}, nil
}

type pointLostAckTx struct {
	driver.Tx
	owner *pointLostAckConnector
}

func (t *pointLostAckTx) Commit() error {
	err := t.Tx.Commit()
	if err == nil && t.owner.fail {
		t.owner.fail = false
		return errors.New("private lost commit acknowledgment")
	}
	return err
}
func TestFormatUnknownCommitStopsAttemptButDoesNotFinalizeReceipt(t *testing.T) {
	memory, _, req := formatFixture(t)
	path := filepath.Join(t.TempDir(), "lost-ack.db")
	_, err := memory.Exec(`VACUUM INTO ?`, path)
	require.NoError(t, err)
	connector := &pointLostAckConnector{path: path, fail: true}
	db := sql.OpenDB(connector)
	defer db.Close()
	db.SetMaxOpenConns(1)
	svc := NewService(NewSQLRepository(db), nil)
	spy := &recordingFormatInvalidator{}
	svc.SetFormatInvalidator(spy)
	result, err := svc.MutateFormat(t.Context(), req)
	require.Nil(t, result)
	require.ErrorIs(t, err, sqltx.ErrCommitOutcomeUnknown)
	require.NotContains(t, err.Error(), "private lost")
	require.Equal(t, 1, spy.calls)
	require.Equal(t, *req.ExpectedPointRevision, spy.seen.ExpectedPointRevision)
	inspected, err := svc.InspectFormatMutation(t.Context(), req.FormatScope)
	require.NoError(t, err)
	require.Equal(t, "pending", inspected.InvalidationStatus)
	require.Equal(t, inspected.Point.RevisionID, spy.seen.PointRevision)
	next, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	require.NoError(t, next.Rollback())
}
