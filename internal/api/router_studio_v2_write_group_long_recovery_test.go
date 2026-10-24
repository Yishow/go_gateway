package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

// TestSchemaOperationLostResponse drives the real ledger and executor: the
// long-layout DDL commits (wholly or partly) but its reply is lost, the
// gateway restarts, and recovery must read the same operation's real catalog
// evidence instead of issuing any DDL or claiming success from runtime state.
func TestSchemaOperationLostResponse(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			// table, index and receipts: a partial commit leaves some, never all.
			for _, committed := range []int{0, 1, 2, 3} {
				t.Run(strconv.Itoa(committed), func(t *testing.T) {
					f := newManagedDestinationFixture(t, kind)
					group := ensureLong(t, f, "instantaneous", "second", 15, "")
					groupID := group["id"].(string)
					table := destinationOf(group)["table_name"].(string)
					token := previewLong(t, f, groupID)
					require.Len(t, token.Statements, 3)

					ledger := recordingplan.NewService(recordingplan.NewSQLRepository(f.db))
					_, outcome, err := ledger.ClaimSchemaApply(t.Context(), &token)
					require.NoError(t, err)
					require.Equal(t, recordingplan.ClaimAcquired, outcome)

					lose := &atomic.Bool{}
					lose.Store(true)
					var target *sql.DB
					if kind == "sqlite" {
						target = lostcommit.Open(f.file, lose)
					} else {
						target = lostcommit.Wrap(stdlib.GetDefaultDriver(), os.Getenv("POSTGRES_DSN"), lose)
					}
					t.Cleanup(func() { require.NoError(t, target.Close()) })
					tx, err := target.BeginTx(t.Context(), nil)
					require.NoError(t, err)
					if kind == "postgres" {
						_, err = tx.ExecContext(t.Context(), `SET LOCAL search_path TO "`+f.namespace+`"`)
						require.NoError(t, err)
					}
					for _, statement := range token.Statements[:committed] {
						_, err := tx.ExecContext(t.Context(), statement)
						require.NoError(t, err)
					}
					require.ErrorIs(t, tx.Commit(), lostcommit.ErrResponseLost)

					old := time.Now().UTC().Add(-recordingplan.SchemaOperationLease - time.Minute)
					_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at=? WHERE operation_id=?`, old, token.OperationID)
					require.NoError(t, err)
					_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_preview_tokens SET expires_at=? WHERE token=?`, old, token.Token)
					require.NoError(t, err)
					var seq int
					var dbName, file string
					require.NoError(t, f.db.QueryRowContext(t.Context(), `PRAGMA database_list`).Scan(&seq, &dbName, &file))
					require.NoError(t, f.db.Close())
					f.db, err = sql.Open("sqlite", file)
					require.NoError(t, err)
					f.db.SetMaxOpenConns(1)
					t.Cleanup(func() { require.NoError(t, f.db.Close()) })
					f.workspace = workspace.NewService(workspace.NewSQLRepository(f.db))
					targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
					f.groups = workspace.NewWriteGroupService(f.workspace, workspace.NewSQLWriteGroupRepository(f.db)).WithTableInspector(dbtarget.NewReadOnlyTableInspector(targets))
					f.groups.WithManagedTableInspector(f.groups.ManagedTableInspector(targets))
					f.router = NewRouter(&DatalinkServices{Workspace: f.workspace, WriteGroups: f.groups, DBTarget: targets, RecordingPlan: recordingplan.NewService(recordingplan.NewSQLRepository(f.db))})

					request := scopeBody(t, f, groupID)
					request["token"], request["operation_id"] = token.Token, token.OperationID
					want := recordingplan.SchemaOperationUnknown
					if committed == 3 {
						want = recordingplan.SchemaOperationSucceeded
					}
					for range 2 {
						response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/schema-apply", request)
						require.Equal(t, http.StatusOK, response.Code, response.Body.String())
						var body struct {
							Data recordingplan.SchemaOperation `json:"data"`
						}
						require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
						require.Equal(t, token.OperationID, body.Data.OperationID)
						require.Equal(t, want, body.Data.Status)
						request = scopeBody(t, f, groupID)
						request["token"], request["operation_id"] = token.Token, token.OperationID
					}

					var objects int
					query, args := `SELECT count(*) FROM sqlite_master WHERE name IN (?, ?, 'gw_effect_receipts')`, []any{table, "idx_" + table + "_src_time"}
					if kind == "postgres" {
						query = `SELECT count(*) FROM pg_catalog.pg_class c JOIN pg_catalog.pg_namespace n ON n.oid=c.relnamespace
							WHERE n.nspname=$1 AND c.relname IN ($2, $3, 'gw_effect_receipts')`
						args = []any{f.namespace, table, "idx_" + table + "_src_time"}
					}
					require.NoError(t, target.QueryRowContext(t.Context(), query, args...).Scan(&objects))
					require.Equal(t, committed, objects, "recovery only inspects evidence; it never runs the missing DDL")
					var operations int
					require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM managed_schema_operations WHERE operation_id=?`, token.OperationID).Scan(&operations))
					require.Equal(t, 1, operations, "the same operation identity is read, never a second one")
					saved := getGroup(t, f, groupID)
					if want == recordingplan.SchemaOperationSucceeded {
						require.Equal(t, token.OperationID, destinationOf(saved)["schema_revision"])
					} else {
						require.Empty(t, destinationOf(saved)["schema_revision"], "an unconfirmed result must not be saved as confirmed")
					}
				})
			}
		})
	}
}
