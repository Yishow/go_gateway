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

func TestManagedGroupSchemaRestartReconcilesRealLostDDLReply(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			for _, committed := range []int{0, 1, 2} {
				t.Run(strconv.Itoa(committed), func(t *testing.T) {
					f := newManagedGroupFixture(t, kind)
					token := f.preview(t)
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
					request := f.schemaRequest()
					request["token"], request["operation_id"] = token.Token, token.OperationID
					want := recordingplan.SchemaOperationUnknown
					if committed == 2 {
						want = recordingplan.SchemaOperationSucceeded
					}
					for range 2 {
						response := performJSONRequest(t, f.router, http.MethodPost, f.confirmPath(), request)
						require.Equal(t, http.StatusOK, response.Code, response.Body.String())
						var body struct {
							Data recordingplan.SchemaOperation `json:"data"`
						}
						require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
						require.Equal(t, token.OperationID, body.Data.OperationID)
						require.Equal(t, want, body.Data.Status)
					}
					var count int
					query, args := `SELECT count(*) FROM sqlite_master WHERE name IN ('raw_values','gw_effect_receipts')`, []any{}
					if kind == "postgres" {
						query, args = `SELECT count(*) FROM information_schema.tables WHERE table_schema=$1`, []any{f.namespace}
					}
					require.NoError(t, target.QueryRowContext(t.Context(), query, args...).Scan(&count))
					require.Equal(t, committed, count, "recovery only inspects evidence; it never runs the missing DDL")
					require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM managed_schema_operations WHERE operation_id=?`, token.OperationID).Scan(&count))
					require.Equal(t, 1, count)
				})
			}
		})
	}
}
