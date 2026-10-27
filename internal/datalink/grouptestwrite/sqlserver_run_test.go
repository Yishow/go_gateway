package grouptestwrite

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"maps"
	"slices"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// This local driver exercises the public service and SQL adapter contracts. It
// is deliberately not evidence of SQL Server engine or transaction behavior.
type sqlServerRunScript struct {
	*sqlServerExecutionScript
	denyDelete      bool
	ambiguousCommit bool
	inserts         int
	deletes         int
}

type sqlServerRunConnector struct{ script *sqlServerRunScript }

func (c sqlServerRunConnector) Connect(context.Context) (driver.Conn, error) {
	return &sqlServerRunConn{script: c.script}, nil
}
func (sqlServerRunConnector) Driver() driver.Driver { return sqlServerExecutionDriver{} }

type sqlServerRunConn struct{ script *sqlServerRunScript }

func (*sqlServerRunConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("local script does not prepare statements")
}
func (*sqlServerRunConn) Close() error { return nil }
func (c *sqlServerRunConn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}
func (c *sqlServerRunConn) BeginTx(_ context.Context, _ driver.TxOptions) (driver.Tx, error) {
	return &sqlServerRunTx{script: c.script, rows: slices.Clone(c.script.rows), receipts: maps.Clone(c.script.receipts)}, nil
}
func (c *sqlServerRunConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if strings.Contains(query, "[effect_key], [payload_digest]") {
		var rows [][]driver.Value
		for _, arg := range args {
			key, _ := arg.Value.(string)
			if digest, ok := c.script.receipts[key]; ok {
				rows = append(rows, []driver.Value{key, digest})
			}
		}
		return &sqlServerExecutionRows{columns: []string{"effect_key", "payload_digest"}, rows: rows}, nil
	}
	return (&sqlServerExecutionConn{script: c.script.sqlServerExecutionScript}).QueryContext(ctx, query, args)
}
func (c *sqlServerRunConn) ExecContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Result, error) {
	if c.script.denyDelete {
		return nil, errors.New("local script denies cleanup")
	}
	c.script.deletes++
	if strings.Contains(query, "[gw_effect_receipts]") {
		removed := len(c.script.receipts)
		clear(c.script.receipts)
		return driver.RowsAffected(removed), nil
	}
	removed := len(c.script.rows)
	c.script.rows = nil
	return driver.RowsAffected(removed), nil
}

type sqlServerRunTx struct {
	script   *sqlServerRunScript
	rows     [][]driver.Value
	receipts map[string]string
}

func (tx *sqlServerRunTx) Commit() error {
	if tx.script.ambiguousCommit {
		return errors.New("local script loses commit acknowledgement")
	}
	return nil
}
func (tx *sqlServerRunTx) Rollback() error {
	tx.script.rows, tx.script.receipts = tx.rows, tx.receipts
	return nil
}

type sqlServerRunDestinations struct {
	sqlServerTestWriteDestinations
	db *sql.DB
}

func (d sqlServerRunDestinations) OpenDestination(context.Context, string, string) (*dbtarget.OpenedDestination, error) {
	return &dbtarget.OpenedDestination{DB: d.db, Kind: schema.DatabaseConnectorKindSQLServer}, nil
}

func TestSQLServerPublicConfirmKeepsWriteAndCleanupResultsSeparate(t *testing.T) {
	for _, native := range []bool{false, true} {
		for _, cleanup := range []string{CleanupCleaned, CleanupFailed, CleanupUnknown} {
			name := "plain/" + cleanup
			if native {
				name = "native/" + cleanup
			}
			t.Run(name, func(t *testing.T) {
				group, layout, tagTypes := sqlServerTestWriteFixture(t)
				if native {
					group, _, layout, tagTypes = sqlServerNativeTestWriteFixture(t)
				}
				columns := make([]string, 0, len(layout.Columns))
				for _, column := range layout.Columns {
					columns = append(columns, column.Name)
				}
				script := &sqlServerRunScript{
					sqlServerExecutionScript: &sqlServerExecutionScript{columns: columns, receipts: map[string]string{}},
					denyDelete:               cleanup == CleanupFailed, ambiguousCommit: cleanup == CleanupUnknown,
				}
				db := sql.OpenDB(sqlServerRunConnector{script: script})
				t.Cleanup(func() { require.NoError(t, db.Close()) })
				groups := &sqlServerTestWriteGroups{
					fakeGroups: &fakeGroups{group: group}, layout: layout, testOnlySupportsGroupTestWrite: true,
				}
				svc := New(Dependencies{
					Groups: groups, Tags: fakeTags(tagTypes), Destinations: sqlServerRunDestinations{db: db},
					Ledger: recordingplan.NewService(recordingplan.NewMemoryRepository()),
				}, Config{StepTimeout: time.Second})
				svc.insertRow = func(_ context.Context, _ *sql.DB, req dbtarget.GroupInsertRequest) (dbtarget.GroupInsertResult, error) {
					script.inserts++
					digest, err := groupdelivery.RowPayloadDigest(req.Row)
					require.NoError(t, err)
					require.Equal(t, digest, req.PayloadDigest)
					require.NotNil(t, req.SQLServer)
					script.receipts[req.Row.EffectKey] = digest
					// Preserve the actual encoded values, including NULLs and UTC times.
					values := make([]driver.Value, len(columns))
					for i, column := range columns {
						for _, cell := range req.Row.Cells {
							if cell.Column == column {
								values[i] = cell.Value
							}
						}
					}
					script.rows = append(script.rows, values)
					return dbtarget.GroupInsertResult{}, nil
				}
				ws := Workspace{ID: group.WorkspaceID, Revision: testRevision}
				preview, err := svc.Preview(t.Context(), ws, group.ID)
				require.NoError(t, err)
				require.Zero(t, script.inserts)
				req := Confirmation{Token: preview.Token, OperationID: preview.OperationID, GroupID: group.ID}
				outcome, err := svc.Confirm(t.Context(), ws, req)
				require.NoError(t, err)
				require.Equal(t, WriteVerified, outcome.Operation.WriteOutcome, outcome.Operation.Reason)
				require.Equal(t, cleanup, outcome.Operation.CleanupStatus, outcome.Operation.CleanupReason)
				require.Equal(t, len(group.Members), script.inserts)
				if cleanup == CleanupFailed {
					require.Len(t, script.rows, len(group.Members))
					require.Len(t, script.receipts, len(group.Members))
				} else {
					require.Empty(t, script.rows)
					require.Empty(t, script.receipts)
				}
				deletes := script.deletes
				again, err := svc.Confirm(t.Context(), ws, req)
				require.NoError(t, err)
				require.Equal(t, recordingplan.ClaimCompleted, again.Claim)
				require.Equal(t, WriteVerified, again.Operation.WriteOutcome)
				require.Equal(t, cleanup, again.Operation.CleanupStatus)
				require.Equal(t, len(group.Members), script.inserts)
				require.Equal(t, deletes, script.deletes, "retained result must not repeat target operations")
			})
		}
	}
}
