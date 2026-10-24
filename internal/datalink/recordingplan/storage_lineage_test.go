package recordingplan

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func lineageDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "lineage.db")+"?_pragma=busy_timeout(5000)")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	for _, name := range []string{"032_write_group_storage_scopes_sqlite.up.sql", "034_write_group_storage_lineage_sqlite.up.sql", "035_write_group_storage_events_sqlite.up.sql"} {
		ddl, err := os.ReadFile(filepath.Join("..", "schema", "migrations", name))
		require.NoError(t, err)
		_, err = db.ExecContext(t.Context(), string(ddl))
		require.NoError(t, err, name)
	}
	return db
}

var lineageScope = StorageScope{WorkspaceID: "ws", ConnectorID: "conn", Database: "db", TableSchema: "main", StorageMode: "plain", LayoutVersion: "managed_long_v1"}

func seedLineage(t *testing.T, db *sql.DB, scope StorageScope) GroupSchemaLayout {
	t.Helper()
	layout, err := ManagedLongGroupLayout("sqlite", "gw_long_0123456789abcdef", "_gw_owner_0123456789abcdef01234567")
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_storage_scopes (workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version,
		table_name, owner_column, layout_digest, generation) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1)`,
		scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode, scope.LayoutVersion,
		layout.TableName, layout.OwnerColumn, "digest")
	require.NoError(t, err)
	return layout
}

// TestStorageLineageRecordsConfirmationsAndRebuilds pins the local evidence a
// recovery relies on: the first confirmation freezes the layout and revisions
// of generation one, further confirmations of the scope do not change the
// generation, and only a table that had to be created again after a
// confirmation, or a rebuild by the repair operation, starts the next one and
// says history may be lost.
func TestStorageLineageRecordsConfirmationsAndRebuilds(t *testing.T) {
	db := lineageDB(t)
	layout := seedLineage(t, db, lineageScope)
	now := time.Date(2026, 10, 24, 1, 2, 3, 456000000, time.UTC)

	before, err := ReadStorageLineage(t.Context(), db, lineageScope)
	require.NoError(t, err)
	require.False(t, before.Confirmed(), "a saved scope is not yet a confirmed table")
	require.Equal(t, 1, before.Generation)

	first := StorageConfirmation{Layout: layout, ConnectorRevision: "crev-1", OperationID: "op-1", VerifiedDigest: "digest-1", CreatedTable: true, GroupID: "g1"}
	require.NoError(t, RecordStorageConfirmation(t.Context(), db, lineageScope, first, now))
	lineage, err := ReadStorageLineage(t.Context(), db, lineageScope)
	require.NoError(t, err)
	require.True(t, lineage.Confirmed())
	require.Equal(t, 1, lineage.Generation, "creating the table for the first time is generation one")
	require.False(t, lineage.HistoryMayBeLost)
	require.Equal(t, layout, *lineage.Layout, "the layout is frozen exactly")
	require.Equal(t, [3]string{"crev-1", "op-1", "digest-1"}, [3]string{lineage.ConnectorRevision, lineage.ConfirmedOperationID, lineage.ConfirmedDigest})

	t.Run("a second group's confirmation or a repeat does not move the generation", func(t *testing.T) {
		require.NoError(t, RecordStorageConfirmation(t.Context(), db, lineageScope, first, now), "the same operation again is a no-op")
		second := first
		second.OperationID, second.GroupID, second.CreatedTable = "op-2", "g2", false
		require.NoError(t, RecordStorageConfirmation(t.Context(), db, lineageScope, second, now))
		lineage, err := ReadStorageLineage(t.Context(), db, lineageScope)
		require.NoError(t, err)
		require.Equal(t, 1, lineage.Generation)
		require.Equal(t, "op-2", lineage.ConfirmedOperationID, "the freshest verification is the evidence")
		require.False(t, lineage.HistoryMayBeLost)
	})

	t.Run("a table created again by a confirmation starts the next generation", func(t *testing.T) {
		again := StorageConfirmation{Layout: layout, ConnectorRevision: "crev-1", OperationID: "op-3", VerifiedDigest: "digest-1", CreatedTable: true, GroupID: "g1"}
		require.NoError(t, RecordStorageConfirmation(t.Context(), db, lineageScope, again, now))
		lineage, err := ReadStorageLineage(t.Context(), db, lineageScope)
		require.NoError(t, err)
		require.Equal(t, 2, lineage.Generation)
		require.True(t, lineage.HistoryMayBeLost)
	})

	t.Run("a rebuild applies once per generation", func(t *testing.T) {
		var wg sync.WaitGroup
		applied := make([]bool, 6)
		for i := range applied {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ok, err := RecordStorageReconstruction(t.Context(), db, lineageScope, 2, "repair-1", "g1", now)
				require.NoError(t, err)
				applied[i] = ok
			}()
		}
		wg.Wait()
		count := 0
		for _, ok := range applied {
			if ok {
				count++
			}
		}
		require.Equal(t, 1, count, "concurrent completions of one repair move the generation once")
		lineage, err := ReadStorageLineage(t.Context(), db, lineageScope)
		require.NoError(t, err)
		require.Equal(t, 3, lineage.Generation)
		require.Equal(t, 1, lineage.Reconstructions)
		require.Equal(t, "repair-1", lineage.LastReconstructionOperationID)
		require.True(t, lineage.LastReconstructedAt.Equal(now))
		require.True(t, lineage.HistoryMayBeLost)
		require.Equal(t, layout, *lineage.Layout, "the confirmed layout is not rewritten by a rebuild")
		require.Equal(t, "op-3", lineage.ConfirmedOperationID, "and no confirmation is made on anyone's behalf")
	})

	t.Run("a blocked repair stays until a confirmation checks the table again", func(t *testing.T) {
		require.NoError(t, BlockStorageRepair(t.Context(), db, lineageScope, 3, "unknown-ddl-result", "repair-2", "g1", now))
		lineage, err := ReadStorageLineage(t.Context(), db, lineageScope)
		require.NoError(t, err)
		require.Equal(t, [2]string{StorageRepairStateBlocked, "unknown-ddl-result"}, [2]string{lineage.RepairState, lineage.RepairReason})

		verified := StorageConfirmation{Layout: layout, ConnectorRevision: "crev-1", OperationID: "op-4", VerifiedDigest: "digest-1", GroupID: "g1"}
		require.NoError(t, RecordStorageConfirmation(t.Context(), db, lineageScope, verified, now))
		lineage, err = ReadStorageLineage(t.Context(), db, lineageScope)
		require.NoError(t, err)
		require.Empty(t, lineage.RepairState)
		require.Empty(t, lineage.RepairReason)
	})

	t.Run("diagnostics are durable, value free and counted", func(t *testing.T) {
		require.NoError(t, RecordStorageRefusal(t.Context(), db, lineageScope, 3, "g1", "storage_owner_mismatch", "table-occupied", now))
		require.NoError(t, RecordStorageRefusal(t.Context(), db, lineageScope, 3, "g1", "storage_owner_mismatch", "table-occupied", now.Add(time.Second)))
		events, err := ReadStorageEvents(t.Context(), db, lineageScope, "g1")
		require.NoError(t, err)
		var refusals []StorageEvent
		for _, event := range events {
			if event.Kind == StorageEventRepairBlocked && event.Reason == "table-occupied" {
				refusals = append(refusals, event)
			}
		}
		require.Len(t, refusals, 1)
		require.Equal(t, 2, refusals[0].Occurrences)
		require.Equal(t, "storage_owner_mismatch", refusals[0].Code)
		kinds := map[string]bool{}
		for _, event := range events {
			kinds[event.Kind] = true
		}
		require.True(t, kinds[StorageEventConfirmed] && kinds[StorageEventRecreatedByConfirm] && kinds[StorageEventReconstructed])
	})

	t.Run("a scope that was never saved has no lineage", func(t *testing.T) {
		other := lineageScope
		other.ConnectorID = "elsewhere"
		_, err := ReadStorageLineage(t.Context(), db, other)
		require.ErrorIs(t, err, ErrStorageLineageNotFound)
		require.ErrorIs(t, RecordStorageConfirmation(t.Context(), db, other, first, now), ErrStorageLineageNotFound)
	})
}
