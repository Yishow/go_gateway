package api

import (
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// TestLongStorageMissingRecovery removes the confirmed shared long table of a
// running, applied group (the receipts stay) and requires the production
// delivery path to restore exactly that table from the confirmed layout and
// deliver the rows that were still waiting, in plain SQLite, plain PostgreSQL
// and a TimescaleDB hypertable. What was stored before the table vanished is
// not restored, and nothing already committed is sent again.
func TestLongStorageMissingRecovery(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres", "timescale"} {
		t.Run(kind, func(t *testing.T) {
			r := newRecoveryRig(t, kind, "aux")
			r.start()

			r.record(10.5)
			r.committed(2, "the first window is delivered into the confirmed table for both shared-scope groups")
			before := r.receipts()
			require.Len(t, before, 2)
			require.Equal(t, 2, r.count(r.table))

			r.drop(r.table)
			require.False(t, r.tableExists(r.table))
			r.record(11.5)
			r.record(12.5)
			r.committed(6, "the rows waiting when the table vanished are delivered after it is restored")

			require.True(t, r.tableExists(r.table), "the confirmed table is back")
			require.Equal(t, 4, r.count(r.table), "only the rows that were still waiting; the earlier rows are history that was lost with the table")
			after := r.receipts()
			require.Len(t, after, 6, "one receipt per delivered effect")
			for key, value := range before {
				require.Equal(t, value, after[key], "an existing receipt is neither deleted nor rewritten")
			}
			require.Equal(t, map[string]int{"sql_committed": 6}, r.outbox(), "no row is blocked, retried forever or sent twice")

			lineage := r.lineage()
			require.Equal(t, 2, lineage.Generation)
			require.True(t, lineage.HistoryMayBeLost)
			require.Equal(t, 1, lineage.Reconstructions, "concurrent groups sharing a scope rebuild the table once")
			require.Empty(t, lineage.RepairState)

			for _, g := range r.groups {
				storage := r.delivery(g.id)["storage"].(map[string]any)
				require.EqualValues(t, 2, storage["generation"])
				require.Equal(t, true, storage["history_may_be_lost"])
				require.EqualValues(t, 1, storage["reconstructions"])
				require.Equal(t, true, storage["recreate_missing"])
			}

			var executing int
			require.NoError(t, r.f.db.QueryRowContext(t.Context(),
				`SELECT COUNT(*) FROM managed_schema_operations WHERE action = ? AND executed_statements > 0`,
				recordingplan.StorageRepairAction).Scan(&executing))
			require.Equal(t, 1, executing, "only one repair operation executed DDL")

			if kind == "timescale" {
				var interval int64
				require.NoError(t, r.dest.QueryRowContext(t.Context(), `
					SELECT d.interval_length FROM _timescaledb_catalog.hypertable h
					JOIN _timescaledb_catalog.dimension d ON d.hypertable_id = h.id
					WHERE h.schema_name = $1 AND h.table_name = $2`, r.f.namespace, r.table).Scan(&interval))
				require.Equal(t, int64(7*24*3600*1_000_000), interval, "the rebuilt hypertable keeps the 7-day chunk interval")
			}
		})
	}

	t.Run("restart_resolves_abandoned_repair_from_catalog", func(t *testing.T) {
		r := newRecoveryRig(t, "sqlite")
		scope := workspace.StorageScopeOf(r.groups[0].snapshot.Group)
		lineage := r.lineage()
		expired := time.Now().UTC().Add(-recordingplan.SchemaOperationLease - time.Minute)
		scopeKey := (&recordingplan.SchemaPreviewToken{
			WorkspaceID: scope.WorkspaceID, ConnectorID: scope.ConnectorID, Database: scope.Database,
			Schema: scope.TableSchema, GroupLayout: lineage.Layout,
		}).ScopeKey()
		detail := fmt.Sprintf(`{"kind":"storage_repair","generation":1,"to_generation":2,"group_id":%q,"connector_id":%q,"table":%q,"scope":{"workspace_id":%q,"connector_id":%q,"database":%q,"table_schema":%q,"storage_mode":%q,"layout_version":%q},"dialect":"sqlite","owner_column":%q}`,
			r.groups[0].id, scope.ConnectorID, r.table, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode, scope.LayoutVersion, lineage.Layout.OwnerColumn)
		_, _, err := recordingplan.NewSQLRepository(r.f.db).ClaimSchemaOperation(t.Context(), &recordingplan.SchemaOperation{
			OperationID: "repair-stale-1", Token: "repair:stale:1:0", WorkspaceID: scope.WorkspaceID,
			ScopeKey: scopeKey, Owner: "dead-owner", Action: recordingplan.StorageRepairAction,
			Status: recordingplan.SchemaOperationRunning, Detail: detail, CreatedAt: expired, UpdatedAt: expired,
		})
		require.NoError(t, err)

		// A started pipeline resolves the abandoned repair from the catalog
		// without executing DDL again, and delivers new rows cleanly.
		r.start()
		r.record(10.5)
		r.committed(1, "window delivered after stale repair resolution")
		r.waitFor(func() bool {
			current, err := recordingplan.ReadStorageLineage(t.Context(), r.f.db, scope)
			return err == nil && current.Generation == 2
		}, "the abandoned repair is resolved from the catalog")
		after := r.lineage()
		require.Equal(t, 2, after.Generation)
		require.True(t, after.HistoryMayBeLost)
		require.Equal(t, "repair-stale-1", after.LastReconstructionOperationID)
	})
}
