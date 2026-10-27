package dbtarget

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func sqlServerLongTestRef() GroupLongTestRowRef {
	const effectKey = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	return GroupLongTestRowRef{
		Kind: schema.DatabaseConnectorKindSQLServer, SchemaName: "dbo", TableName: "gw_long_0123456789abcdef",
		OwnerColumn: groupLongTestOwner, OwnerValue: "gateway", ProvenanceColumn: "provenance",
		WorkspaceID: "ws-1", GroupID: "group-1", GroupRevision: "rev-1", OperationOwner: "gw-test-op1",
		RecordIDs: []string{strings.Repeat("a", 64)}, SampleIDs: []string{"gw-test-op1-0"},
		EffectKeys: []string{effectKey}, EffectDigests: map[string]string{effectKey: "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"},
	}
}

func TestSQLServerLongTestRowRefAcceptsSafeOperationIdentity(t *testing.T) {
	ref := sqlServerLongTestRef()
	require.NoError(t, ref.validate(), "a SQL Server operation-owned long-row proof should be valid")

	foreignSample := ref
	foreignSample.SampleIDs = []string{"gw-test-other-0"}
	require.ErrorIs(t, foreignSample.validate(), ErrOwnedRowRefInvalid)

	malformedSample := ref
	malformedSample.SampleIDs = []string{"gw-test-op1-not-an-index"}
	require.ErrorIs(t, malformedSample.validate(), ErrOwnedRowRefInvalid)

	malformedEffect := ref
	malformedEffect.EffectKeys = []string{"not-an-exact-effect-key"}
	malformedEffect.EffectDigests = map[string]string{"not-an-exact-effect-key": strings.Repeat("c", 64)}
	require.ErrorIs(t, malformedEffect.validate(), ErrOwnedRowRefInvalid)
}

type sqlServerLongTestQueryScript struct {
	query            string
	args             []any
	foreignWorkspace bool
}

type sqlServerLongTestQueryConnector struct{ script *sqlServerLongTestQueryScript }

func (c sqlServerLongTestQueryConnector) Connect(context.Context) (driver.Conn, error) {
	return &sqlServerLongTestQueryConn{script: c.script}, nil
}

func (sqlServerLongTestQueryConnector) Driver() driver.Driver { return sqlServerLongTestQueryDriver{} }

type sqlServerLongTestQueryDriver struct{}

func (sqlServerLongTestQueryDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("open is not used by Connector")
}

type sqlServerLongTestQueryConn struct{ script *sqlServerLongTestQueryScript }

func (*sqlServerLongTestQueryConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare is not used")
}
func (*sqlServerLongTestQueryConn) Close() error { return nil }
func (*sqlServerLongTestQueryConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transaction is not used")
}
func (c *sqlServerLongTestQueryConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	c.script.query = query
	c.script.args = make([]any, len(args))
	for index, arg := range args {
		c.script.args[index] = arg.Value
	}
	columns := []string{longColRecordID, longColWorkspaceID, longColGroupID, longColGroupRev, longColProvenance, groupLongTestOwner}
	if c.script.foreignWorkspace && strings.Count(strings.ToLower(query), "[workspace_id]") > 1 {
		return &sqlServerLongTestQueryRows{columns: columns, done: true}, nil
	}
	workspaceID := "ws-1"
	if c.script.foreignWorkspace {
		workspaceID = "another-workspace"
	}
	return &sqlServerLongTestQueryRows{
		columns: columns,
		values: []driver.Value{
			strings.Repeat("a", 64), workspaceID, "group-1", "rev-1",
			`{"sample_id":"gw-test-op1-0","member":"test"}`, "gateway",
		},
	}, nil
}

type sqlServerLongTestQueryRows struct {
	columns []string
	values  []driver.Value
	done    bool
}

func (r *sqlServerLongTestQueryRows) Columns() []string { return r.columns }
func (*sqlServerLongTestQueryRows) Close() error        { return nil }
func (r *sqlServerLongTestQueryRows) Next(dest []driver.Value) error {
	if r.done {
		return io.EOF
	}
	r.done = true
	copy(dest, r.values)
	return nil
}

func TestSQLServerLongTestRowsUseBoundExactSQLServerQuery(t *testing.T) {
	ref := sqlServerLongTestRef()
	script := &sqlServerLongTestQueryScript{}
	db := sql.OpenDB(sqlServerLongTestQueryConnector{script: script})
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	rows, err := ReadGroupLongTestRows(t.Context(), db, ref, []string{longColRecordID})
	require.NoError(t, err)
	require.Equal(t, [][]any{{strings.Repeat("a", 64)}}, rows)
	query := strings.ToLower(script.query)
	require.Contains(t, query, "select top (257)")
	require.Contains(t, query, "[dbo].[gw_long_0123456789abcdef]")
	require.Contains(t, query, "latin1_general_100_bin2")
	where := strings.SplitN(query, " where ", 2)[1]
	require.Contains(t, where, "[record_id]")
	require.NotContains(t, where, "workspace_id")
	require.NotContains(t, where, "group_id")
	require.NotContains(t, where, "group_revision")
	require.NotContains(t, where, "provenance")
	require.Equal(t, []any{strings.Repeat("a", 64)}, script.args)
	require.NotContains(t, query, "limit ")
	require.NotContains(t, query, "hex(")
	require.NotContains(t, query, "json_extract")
	require.NotContains(t, query, "?")
}

func TestSQLServerLongTestReadDoesNotHideForeignScopeWithSameRecordID(t *testing.T) {
	ref := sqlServerLongTestRef()
	script := &sqlServerLongTestQueryScript{foreignWorkspace: true}
	db := sql.OpenDB(sqlServerLongTestQueryConnector{script: script})
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	_, err := ReadGroupLongTestRows(t.Context(), db, ref, []string{longColRecordID})
	require.ErrorIs(t, err, ErrOwnedRowRefInvalid, "the exact record ID must be fetched before scope proof is checked")
}

type sqlServerLongTestCleanupScript struct {
	queries          []string
	queryArgs        [][]any
	execs            []string
	execArgs         [][]any
	rowPresent       bool
	receiptPresent   bool
	foreignWorkspace bool
	foreignSample    bool
	foreignReceipt   bool
}

type sqlServerLongTestCleanupConnector struct {
	script *sqlServerLongTestCleanupScript
}

func (c sqlServerLongTestCleanupConnector) Connect(context.Context) (driver.Conn, error) {
	return &sqlServerLongTestCleanupConn{script: c.script}, nil
}

func (sqlServerLongTestCleanupConnector) Driver() driver.Driver {
	return sqlServerLongTestQueryDriver{}
}

type sqlServerLongTestCleanupConn struct {
	script *sqlServerLongTestCleanupScript
}

func (*sqlServerLongTestCleanupConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare is not used")
}
func (*sqlServerLongTestCleanupConn) Close() error { return nil }
func (*sqlServerLongTestCleanupConn) Begin() (driver.Tx, error) {
	return nil, errors.New("BeginTx is required")
}
func (c *sqlServerLongTestCleanupConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return sqlServerLongTestCleanupTx{}, nil
}
func (c *sqlServerLongTestCleanupConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	c.script.queries = append(c.script.queries, query)
	c.script.queryArgs = append(c.script.queryArgs, sqlServerTestNamedValues(args))
	lower := strings.ToLower(query)
	switch {
	case strings.Contains(lower, "count(*)"):
		return newSQLServerTestRows([]string{"count"}, [][]driver.Value{{int64(0)}}), nil
	case strings.Contains(lower, "select [effect_key]"):
		if c.script.receiptPresent {
			ref := sqlServerLongTestRef()
			key := ref.EffectKeys[0]
			digest := ref.EffectDigests[key]
			if c.script.foreignReceipt {
				digest = strings.Repeat("d", 64)
			}
			return newSQLServerTestRows([]string{"effect_key", "payload_digest"}, [][]driver.Value{{key, digest}}), nil
		}
		return newSQLServerTestRows([]string{"effect_key", "payload_digest"}, nil), nil
	case strings.Contains(lower, "select top"):
		columns := []string{longColRecordID, longColWorkspaceID, longColGroupID, longColGroupRev, longColProvenance, groupLongTestOwner}
		if c.script.rowPresent {
			workspaceID, sampleID := "ws-1", "gw-test-op1-0"
			if c.script.foreignWorkspace {
				workspaceID = "another-workspace"
			}
			if c.script.foreignSample {
				sampleID = "gw-test-op1-other"
			}
			values := []driver.Value{
				strings.Repeat("a", 64), workspaceID, "group-1", "rev-1",
				`{"sample_id":"` + sampleID + `","member":"test"}`, "gateway",
			}
			return newSQLServerTestRows(columns, [][]driver.Value{values}), nil
		}
		return newSQLServerTestRows(columns, nil), nil
	default:
		return nil, errors.New("unexpected scripted query")
	}
}
func (c *sqlServerLongTestCleanupConn) ExecContext(_ context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	c.script.execs = append(c.script.execs, query)
	c.script.execArgs = append(c.script.execArgs, sqlServerTestNamedValues(args))
	switch {
	case strings.Contains(strings.ToLower(query), "gw_long_0123456789abcdef"):
		c.script.rowPresent = false
	case strings.Contains(strings.ToLower(query), "gw_effect_receipts"):
		c.script.receiptPresent = false
	default:
		return nil, errors.New("unexpected scripted exec")
	}
	return driver.RowsAffected(1), nil
}

type sqlServerLongTestCleanupTx struct{}

func (sqlServerLongTestCleanupTx) Commit() error   { return nil }
func (sqlServerLongTestCleanupTx) Rollback() error { return nil }

func sqlServerTestNamedValues(values []driver.NamedValue) []any {
	result := make([]any, len(values))
	for index, value := range values {
		result[index] = value.Value
	}
	return result
}

type sqlServerTestRows struct {
	columns []string
	values  [][]driver.Value
	next    int
}

func newSQLServerTestRows(columns []string, values [][]driver.Value) *sqlServerTestRows {
	return &sqlServerTestRows{columns: columns, values: values}
}
func (r *sqlServerTestRows) Columns() []string { return r.columns }
func (*sqlServerTestRows) Close() error        { return nil }
func (r *sqlServerTestRows) Next(dest []driver.Value) error {
	if r.next >= len(r.values) {
		return io.EOF
	}
	copy(dest, r.values[r.next])
	r.next++
	return nil
}

func TestSQLServerLongTestCleanupRejectsForeignProofBeforeDeleting(t *testing.T) {
	for name, options := range map[string]func(*sqlServerLongTestCleanupScript){
		"foreign scope":          func(script *sqlServerLongTestCleanupScript) { script.foreignWorkspace = true },
		"foreign sample":         func(script *sqlServerLongTestCleanupScript) { script.foreignSample = true },
		"foreign receipt digest": func(script *sqlServerLongTestCleanupScript) { script.foreignReceipt = true },
	} {
		t.Run(name, func(t *testing.T) {
			script := &sqlServerLongTestCleanupScript{rowPresent: true, receiptPresent: true}
			options(script)
			db := sql.OpenDB(sqlServerLongTestCleanupConnector{script: script})
			db.SetMaxOpenConns(1)
			t.Cleanup(func() { require.NoError(t, db.Close()) })

			_, err := RemoveGroupLongTestRows(t.Context(), db, sqlServerLongTestRef())
			require.ErrorIs(t, err, ErrOwnedRowRefInvalid)
			require.Empty(t, script.execs, "preflight conflicts must stop before any DELETE")
			require.True(t, script.rowPresent)
			require.True(t, script.receiptPresent)
		})
	}
}

func TestSQLServerLongTestCleanupUsesExactSQLServerProofAndReceipt(t *testing.T) {
	script := &sqlServerLongTestCleanupScript{rowPresent: true, receiptPresent: true}
	db := sql.OpenDB(sqlServerLongTestCleanupConnector{script: script})
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	removed, err := RemoveGroupLongTestRows(t.Context(), db, sqlServerLongTestRef())
	require.NoError(t, err)
	require.EqualValues(t, 1, removed)
	require.Len(t, script.execs, 2)
	for _, query := range append(append([]string{}, script.queries...), script.execs...) {
		lower := strings.ToLower(query)
		require.NotContains(t, lower, "?")
		require.NotContains(t, lower, "limit ")
		require.NotContains(t, lower, "hex(")
		require.NotContains(t, lower, "json_extract")
	}
	rowDelete := strings.ToLower(script.execs[0])
	require.Contains(t, rowDelete, "delete from [dbo].[gw_long_0123456789abcdef]")
	require.Contains(t, rowDelete, "json_value")
	require.Contains(t, rowDelete, "latin1_general_100_bin2")
	receiptDelete := strings.ToLower(script.execs[1])
	require.Contains(t, receiptDelete, "[effect_key]")
	require.Contains(t, receiptDelete, "[payload_digest]")
	var allArgs []any
	for _, args := range append(script.queryArgs, script.execArgs...) {
		allArgs = append(allArgs, args...)
	}
	for _, expected := range []any{"gateway", "ws-1", "group-1", "rev-1", strings.Repeat("a", 64), "gw-test-op1-0",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"} {
		require.Contains(t, allArgs, expected)
	}
	require.False(t, script.rowPresent)
	require.False(t, script.receiptPresent)
}
