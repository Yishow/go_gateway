package recordingplan

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// restoreFixture is a confirmed managed long table in a real SQLite
// destination, its receipt table beside it, and the local ledger and lineage a
// confirmation leaves behind. Nothing here is a fake of the code under test:
// the destination is a database file, the ledger is the production SQL
// repository, and the only seam is the statement runner the test may break.
type restoreFixture struct {
	t        *testing.T
	local    *sql.DB
	target   *sql.DB
	svc      *Service
	layout   GroupSchemaLayout
	scope    StorageScope
	request  StorageRestoreRequest
	digest   string
	ddlRuns  atomic.Int32
	runner   func(ctx context.Context, statements []string) (SchemaExecution, error)
	canWrite func(ctx context.Context) (bool, error)
	inspect  func(ctx context.Context, table string) (TargetTableInspection, error)
}

func newRestoreFixture(t *testing.T) *restoreFixture {
	t.Helper()
	local, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "local.db")+"?_pragma=busy_timeout(5000)")
	require.NoError(t, err)
	t.Cleanup(func() { _ = local.Close() })
	applyRecordingDDL(t, local)
	for _, name := range []string{"032_write_group_storage_scopes_sqlite.up.sql", "034_write_group_storage_lineage_sqlite.up.sql", "035_write_group_storage_events_sqlite.up.sql"} {
		ddl, err := os.ReadFile(filepath.Join("..", "schema", "migrations", name))
		require.NoError(t, err)
		_, err = local.ExecContext(t.Context(), string(ddl))
		require.NoError(t, err, name)
	}
	target, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "target.db")+"?_pragma=busy_timeout(5000)")
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })

	f := &restoreFixture{t: t, local: local, target: target, svc: NewService(NewSQLRepository(local)), scope: lineageScope}
	f.layout = seedLineage(t, local, f.scope)
	statements, err := GenerateGroupSchemaDDL("sqlite", f.layout)
	require.NoError(t, err)
	for _, statement := range statements {
		_, err := target.ExecContext(t.Context(), statement)
		require.NoError(t, err, statement)
	}
	f.inspect = f.sqliteInspector
	inspected, err := inspectGroupTables(t.Context(), "sqlite", f.layout, f.sqliteInspector)
	require.NoError(t, err)
	f.digest = schemaTablesDigest(inspected.tables)
	require.NoError(t, RecordStorageConfirmation(t.Context(), local, f.scope, StorageConfirmation{
		Layout: f.layout, ConnectorRevision: "crev-1", OperationID: "op-confirm", VerifiedDigest: f.digest, CreatedTable: true, GroupID: "g1",
	}, time.Now()))
	f.request = StorageRestoreRequest{
		Scope: f.scope, Dialect: "sqlite", ConnectorRevision: "crev-1",
		TableName: f.layout.TableName, OwnerColumn: f.layout.OwnerColumn, GroupID: "g1",
	}
	f.runner = f.realRunner
	return f
}

func (f *restoreFixture) sqliteInspector(ctx context.Context, table string) (TargetTableInspection, error) {
	info, err := dbtarget.InspectOpenManagedTable(ctx, f.target, schema.DatabaseConnectorKindSQLite, "main", table, table == f.layout.TableName)
	if err != nil {
		return TargetTableInspection{}, err
	}
	state := TargetTableInspection{
		Status: string(info.Status), ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{}, ColumnPrimaryKey: map[string]bool{},
		IndexesInspected: info.IndexesInspected,
	}
	for _, column := range info.Columns {
		state.Columns = append(state.Columns, column.Name)
		state.ColumnTypes[column.Name] = column.DataType
		state.ColumnNullable[column.Name] = column.Nullable
		state.ColumnPrimaryKey[column.Name] = column.PrimaryKey
	}
	for _, index := range info.Indexes {
		state.Indexes = append(state.Indexes, TargetIndexInspection{Name: index.Name, Columns: index.Columns, Unique: index.Unique})
	}
	return state, nil
}

func (f *restoreFixture) realRunner(ctx context.Context, statements []string) (SchemaExecution, error) {
	execution, err := dbtarget.RunManagedSchemaBatch(ctx, f.target, schema.DatabaseConnectorKindSQLite, "main", statements)
	if execution == nil {
		return SchemaExecution{RolledBack: true}, err
	}
	return SchemaExecution{Committed: execution.Committed, RolledBack: execution.RolledBack}, err
}

func (f *restoreFixture) restore(ctx context.Context) (StorageRestoreResult, error) {
	return f.svc.RestoreManagedStorage(ctx, f.request, StorageRestoreTarget{
		Inspect: f.inspect,
		Execute: func(ctx context.Context, statements []string) (SchemaExecution, error) {
			f.ddlRuns.Add(1)
			return f.runner(ctx, statements)
		},
		CanCreate: f.canWrite,
	})
}

func (f *restoreFixture) exists(table string) bool {
	f.t.Helper()
	var n int
	require.NoError(f.t, f.target.QueryRowContext(f.t.Context(), `SELECT COUNT(*) FROM sqlite_master WHERE name = ?`, table).Scan(&n))
	return n > 0
}

func (f *restoreFixture) drop(table string) {
	f.t.Helper()
	_, err := f.target.ExecContext(f.t.Context(), `DROP TABLE "`+table+`"`)
	require.NoError(f.t, err)
}

func (f *restoreFixture) lineage() *StorageLineage {
	f.t.Helper()
	lineage, err := ReadStorageLineage(f.t.Context(), f.local, f.scope)
	require.NoError(f.t, err)
	return lineage
}

func (f *restoreFixture) repairOperations() int {
	f.t.Helper()
	var n int
	require.NoError(f.t, f.local.QueryRowContext(f.t.Context(), `SELECT COUNT(*) FROM managed_schema_operations WHERE action = ?`, StorageRepairAction).Scan(&n))
	return n
}

func (f *restoreFixture) receipts() int {
	f.t.Helper()
	var n int
	require.NoError(f.t, f.target.QueryRowContext(f.t.Context(), `SELECT COUNT(*) FROM gw_effect_receipts`).Scan(&n))
	return n
}

// TestStorageRestoreRebuildsTheConfirmedTable: the confirmed layout comes back
// exactly, once, verified from the catalog; the receipts are not touched; the
// lineage moves to the next generation and says history may be lost.
func TestStorageRestoreRebuildsTheConfirmedTable(t *testing.T) {
	f := newRestoreFixture(t)
	_, err := f.target.ExecContext(t.Context(), `INSERT INTO gw_effect_receipts (effect_key, payload_digest, committed_at) VALUES ('e1', 'd1', 't1')`)
	require.NoError(t, err)
	f.drop(f.layout.TableName)

	result, err := f.restore(t.Context())
	require.NoError(t, err)
	require.Equal(t, StorageRestored, result.Outcome, "%+v", result)
	require.True(t, f.exists(f.layout.TableName))
	require.Equal(t, 1, int(f.ddlRuns.Load()))
	require.Equal(t, SchemaOperationSucceeded, result.Operation.Status)
	require.Equal(t, f.digest, result.Operation.VerifiedDigest, "the rebuilt table verifies to the digest the confirmation recorded")
	require.Equal(t, 2, result.Operation.ExecutedStatements, "the table and its index; never the receipt table")
	require.Equal(t, 1, f.receipts(), "the receipt history is untouched")

	lineage := f.lineage()
	require.Equal(t, 2, lineage.Generation)
	require.True(t, lineage.HistoryMayBeLost)
	require.Equal(t, 1, lineage.Reconstructions)
	require.Equal(t, result.Operation.OperationID, lineage.LastReconstructionOperationID)
	require.Equal(t, "op-confirm", lineage.ConfirmedOperationID, "nobody's confirmation was made on their behalf")
	require.Equal(t, f.layout, *lineage.Layout)

	again, err := f.restore(t.Context())
	require.NoError(t, err)
	require.Equal(t, StorageNotMissing, again.Outcome, "a table that is there is nothing to restore")
	require.Equal(t, 1, int(f.ddlRuns.Load()))
	require.Equal(t, 1, f.repairOperations())
}

func TestStorageRestoreIsDoneOnceWhenMeetingConcurrently(t *testing.T) {
	f := newRestoreFixture(t)
	f.drop(f.layout.TableName)

	const callers = 8
	results := make([]StorageRestoreResult, callers)
	errs := make([]error, callers)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			results[i], errs[i] = f.restore(t.Context())
		}()
	}
	close(start)
	wg.Wait()

	restored := 0
	for i, result := range results {
		require.NoError(t, errs[i])
		switch result.Outcome {
		case StorageRestored:
			restored++
		case StorageNotMissing, StorageInProgress:
		default:
			t.Fatalf("unexpected outcome %+v", result)
		}
	}
	require.Equal(t, 1, restored, "exactly one caller restored the table")
	require.Equal(t, 1, int(f.ddlRuns.Load()), "one effective reconstruction")
	require.Equal(t, 2, f.lineage().Generation)
	require.Equal(t, 1, f.lineage().Reconstructions)
	require.True(t, f.exists(f.layout.TableName))
	// A latecomer that claimed after the table was back may leave a no-op
	// operation behind, but exactly one operation ever ran statements.
	var executing, unfinished int
	require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM managed_schema_operations WHERE action = ? AND executed_statements > 0`, StorageRepairAction).Scan(&executing))
	require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM managed_schema_operations WHERE action = ? AND status NOT IN ('succeeded')`, StorageRepairAction).Scan(&unfinished))
	require.Equal(t, 1, executing)
	require.Zero(t, unfinished, "every operation ended, and none failed")
}

func TestStorageRestoreRefusals(t *testing.T) {
	cases := map[string]struct {
		prepare func(*restoreFixture)
		code    string
		reason  string
	}{
		"a foreign table holds the name": {func(f *restoreFixture) {
			f.drop(f.layout.TableName)
			_, err := f.target.ExecContext(t.Context(), `CREATE TABLE "`+f.layout.TableName+`" (note TEXT)`)
			require.NoError(t, err)
		}, StorageCodeOwnerMismatch, StorageReasonTableOccupied},
		"a view holds the name": {func(f *restoreFixture) {
			f.drop(f.layout.TableName)
			_, err := f.target.ExecContext(t.Context(), `CREATE VIEW "`+f.layout.TableName+`" AS SELECT 1 AS note`)
			require.NoError(t, err)
		}, StorageCodeOwnerMismatch, StorageReasonTableOccupied},
		"the receipts are gone too": {func(f *restoreFixture) {
			f.drop(f.layout.TableName)
			f.drop("gw_effect_receipts")
		}, StorageCodeRepairBlocked, StorageReasonReceiptsMissing},
		"the connector revision changed": {func(f *restoreFixture) {
			f.drop(f.layout.TableName)
			f.request.ConnectorRevision = "crev-2"
		}, StorageCodeRepairBlocked, StorageReasonConnectorRevision},
		"the group derives another table": {func(f *restoreFixture) {
			f.drop(f.layout.TableName)
			f.request.TableName = "gw_long_other"
		}, StorageCodeOwnerMismatch, StorageReasonOwnershipMismatch},
		"the group derives another owner marker": {func(f *restoreFixture) {
			f.drop(f.layout.TableName)
			f.request.OwnerColumn = "_gw_owner_other"
		}, StorageCodeOwnerMismatch, StorageReasonOwnershipMismatch},
		"the table was never confirmed": {func(f *restoreFixture) {
			f.drop(f.layout.TableName)
			_, err := f.local.ExecContext(t.Context(), `UPDATE write_group_storage_scopes SET confirmed_operation_id = '', confirmed_layout = ''`)
			require.NoError(t, err)
			// The frozen group still derives its layout, which is enough to look
			// at the catalog and never enough to restore.
			f.request.FallbackLayout = &f.layout
		}, StorageCodeOwnerMismatch, StorageReasonConfirmationMissing},
		"a restore was left blocked": {func(f *restoreFixture) {
			f.drop(f.layout.TableName)
			require.NoError(t, BlockStorageRepair(t.Context(), f.local, f.scope, 1, StorageReasonUnknownResult, "repair-x", "g1", time.Now()))
		}, StorageCodeRepairBlocked, StorageReasonRepairNeedsConfirm},
		"the role may not create tables": {func(f *restoreFixture) {
			f.drop(f.layout.TableName)
			f.canWrite = func(context.Context) (bool, error) { return false, nil }
		}, StorageCodeRepairBlocked, StorageReasonPermissionDenied},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			f := newRestoreFixture(t)
			tc.prepare(f)
			before := f.lineage()

			result, err := f.restore(t.Context())
			require.NoError(t, err)
			require.Equal(t, StorageBlocked, result.Outcome, "%+v", result)
			require.Equal(t, [2]string{tc.code, tc.reason}, [2]string{result.Code, result.Reason})
			require.Zero(t, f.ddlRuns.Load(), "a refusal runs no statement")
			require.Zero(t, f.repairOperations(), "and claims no operation")
			after := f.lineage()
			require.Equal(t, before.Generation, after.Generation)
			require.Equal(t, before.HistoryMayBeLost, after.HistoryMayBeLost)

			events, err := ReadStorageEvents(t.Context(), f.local, f.scope, "g1")
			require.NoError(t, err)
			var recorded bool
			for _, event := range events {
				recorded = recorded || (event.Kind == StorageEventRepairBlocked && event.Code == tc.code && event.Reason == tc.reason)
			}
			require.True(t, recorded, "the refusal is a durable diagnostic: %+v", events)
		})
	}

	t.Run("an inspection that cannot be confirmed decides nothing", func(t *testing.T) {
		f := newRestoreFixture(t)
		f.drop(f.layout.TableName)
		f.inspect = func(context.Context, string) (TargetTableInspection, error) {
			return TargetTableInspection{}, errors.New("catalog unavailable")
		}
		_, err := f.restore(t.Context())
		require.Error(t, err)
		require.Zero(t, f.ddlRuns.Load())
		require.False(t, f.exists(f.layout.TableName))
	})

	t.Run("a hypertable layout needs the extension where it was confirmed", func(t *testing.T) {
		for name, tc := range map[string]struct {
			evidence *TargetTimescaleEvidence
			reason   string
		}{
			"unavailable": {&TargetTimescaleEvidence{Reasons: []string{"extension-not-installed"}}, StorageReasonTimescaleUnavailable},
			"moved":       {&TargetTimescaleEvidence{Available: true, ExtensionSchema: "elsewhere", Version: "2.31.0"}, StorageReasonTimescaleMoved},
		} {
			t.Run(name, func(t *testing.T) {
				f := newRestoreFixture(t)
				layout, err := ManagedLongTimescaleLayout(f.layout.TableName, f.layout.OwnerColumn)
				require.NoError(t, err)
				layout.Partition.ExtensionSchema = "public"
				_, err = f.local.ExecContext(t.Context(), `UPDATE write_group_storage_scopes SET confirmed_layout = ?`, mustJSON(t, layout))
				require.NoError(t, err)
				f.request.Dialect = "postgres"
				f.inspect = func(_ context.Context, table string) (TargetTableInspection, error) {
					if table == f.layout.TableName {
						return TargetTableInspection{Status: inspectionMissing, Timescale: tc.evidence}, nil
					}
					state := groupReceiptCompatibleState()
					state.ColumnNullable["effect_key"] = false
					return state, nil
				}
				result, err := f.restore(t.Context())
				require.NoError(t, err)
				require.Equal(t, StorageBlocked, result.Outcome)
				require.Equal(t, [2]string{StorageCodeRepairBlocked, tc.reason}, [2]string{result.Code, result.Reason})
				require.Zero(t, f.ddlRuns.Load(), "a destination that cannot host the confirmed hypertable gets no statement")
			})
		}
	})
}
