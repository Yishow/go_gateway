package sourcerule

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type rollbackTestConnector struct {
	conn *rollbackTestConn
}

func (c rollbackTestConnector) Connect(context.Context) (driver.Conn, error) {
	return c.conn, nil
}

func (c rollbackTestConnector) Driver() driver.Driver {
	return rollbackTestDriver{}
}

type rollbackTestDriver struct{}

func (rollbackTestDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("rollback test driver requires OpenDB")
}

type rollbackTestConn struct {
	execErr       error
	rollbackErr   error
	rollbackCalls int
	tx            *rollbackTestTx
}

func (c *rollbackTestConn) Prepare(query string) (driver.Stmt, error) {
	return &rollbackTestStmt{conn: c, query: query}, nil
}

func (c *rollbackTestConn) PrepareContext(_ context.Context, query string) (driver.Stmt, error) {
	return c.Prepare(query)
}

func (c *rollbackTestConn) Close() error {
	return nil
}

func (c *rollbackTestConn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

func (c *rollbackTestConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	c.tx = &rollbackTestTx{conn: c}
	return c.tx, nil
}

func (c *rollbackTestConn) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	return c.execResult()
}

func (c *rollbackTestConn) Exec(string, []driver.Value) (driver.Result, error) {
	return c.execResult()
}

func (c *rollbackTestConn) execResult() (driver.Result, error) {
	if c.execErr != nil {
		return nil, c.execErr
	}
	return driver.RowsAffected(1), nil
}

type rollbackTestTx struct {
	conn *rollbackTestConn
}

func (tx *rollbackTestTx) Commit() error {
	return nil
}

func (tx *rollbackTestTx) Rollback() error {
	tx.conn.rollbackCalls++
	return tx.conn.rollbackErr
}

type rollbackTestStmt struct {
	conn  *rollbackTestConn
	query string
}

func (s *rollbackTestStmt) Close() error {
	return nil
}

func (s *rollbackTestStmt) NumInput() int {
	return -1
}

func (s *rollbackTestStmt) Exec([]driver.Value) (driver.Result, error) {
	return s.conn.execResult()
}

func (s *rollbackTestStmt) ExecContext(context.Context, []driver.NamedValue) (driver.Result, error) {
	return s.conn.execResult()
}

func (s *rollbackTestStmt) Query([]driver.Value) (driver.Rows, error) {
	return nil, errors.New("rollback test driver does not support queries")
}

func TestSQLRepository_CreateLinksJoinsNonDoneRollbackError(t *testing.T) {
	mainErr := errors.New("injected link insert failure")
	rollbackErr := errors.New("injected transaction rollback failure")
	conn := &rollbackTestConn{execErr: mainErr, rollbackErr: rollbackErr}
	db := sql.OpenDB(rollbackTestConnector{conn: conn})
	defer db.Close()

	err := NewSQLRepository(db).CreateLinks(context.Background(), []*schema.SourceRuleLink{{
		ID:        "link-1",
		RuleID:    "rule-1",
		Address:   "40001",
		PointID:   "point-1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}})

	require.Error(t, err)
	assert.ErrorIs(t, err, mainErr)
	assert.ErrorIs(t, err, rollbackErr)
	assert.Equal(t, 1, conn.rollbackCalls)
}

func TestSQLRepository_CreateLinksIgnoresErrTxDoneAfterCommit(t *testing.T) {
	conn := &rollbackTestConn{}
	db := sql.OpenDB(rollbackTestConnector{conn: conn})
	defer db.Close()

	err := NewSQLRepository(db).CreateLinks(context.Background(), []*schema.SourceRuleLink{{
		ID:        "link-1",
		RuleID:    "rule-1",
		Address:   "40001",
		PointID:   "point-1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}})

	require.NoError(t, err)
	assert.Zero(t, conn.rollbackCalls)
}

func TestSQLRepository_CreateLinksIgnoresDriverErrTxDone(t *testing.T) {
	mainErr := errors.New("injected link insert failure")
	conn := &rollbackTestConn{execErr: mainErr, rollbackErr: sql.ErrTxDone}
	db := sql.OpenDB(rollbackTestConnector{conn: conn})
	defer db.Close()

	err := NewSQLRepository(db).CreateLinks(context.Background(), []*schema.SourceRuleLink{{
		ID:        "link-1",
		RuleID:    "rule-1",
		Address:   "40001",
		PointID:   "point-1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}})

	require.Error(t, err)
	assert.ErrorIs(t, err, mainErr)
	assert.NotErrorIs(t, err, sql.ErrTxDone)
	assert.Equal(t, 1, conn.rollbackCalls)
}

func TestSQLRepository_ReplaceCandidateSnapshotsJoinsRollbackError(t *testing.T) {
	mainErr := errors.New("injected snapshot write failure")
	rollbackErr := errors.New("injected snapshot rollback failure")
	conn := &rollbackTestConn{execErr: mainErr, rollbackErr: rollbackErr}
	db := sql.OpenDB(rollbackTestConnector{conn: conn})
	defer db.Close()

	err := NewSQLRepository(db).ReplaceCandidateSnapshots(context.Background(), []*schema.SourceRuleCandidateSnapshot{{
		SourceRuleID:  "rule-1",
		RevisionID:    "revision-1",
		CandidateType: schema.SourceRuleCandidateTypeTags,
		Payload:       `{"candidates":[]}`,
		GeneratedAt:   time.Now(),
	}})

	require.Error(t, err)
	assert.ErrorIs(t, err, mainErr)
	assert.ErrorIs(t, err, rollbackErr)
	assert.Equal(t, 1, conn.rollbackCalls)
}
