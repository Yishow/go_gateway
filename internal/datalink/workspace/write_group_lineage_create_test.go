package workspace

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCanonicalLineageFreshCreateReplayAndRestart(t *testing.T) {
	for _, kind := range []string{"basic", "profile", "dual"} {
		t.Run(kind, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "lineage.db")
			h := newCircuitProfileHarness(t, path)
			h.db.SetMaxOpenConns(1)
			sequence := 0
			h.service.repo.newID = func() string {
				sequence++
				if sequence == 1 {
					return "G1"
				}
				return fmt.Sprintf("R%d", sequence-1)
			}
			mutation := h.mutation(t.Context(), lineageTestDraft(h, kind))
			first, err := ensureLineageKind(h, kind, mutation)
			require.NoError(t, err)
			require.Equal(t, "G1", first.Group.ID)
			require.Equal(t, "R1", first.Group.Revision, "lineage UUID generation must not consume group revisions")
			lineage := requireLineageGenesis(t, h.db, first.Group)
			_, err = uuid.Parse(lineage)
			require.NoError(t, err, "lineage uses an independent UUID")
			for _, table := range lineageKeyTables(kind) {
				var stored string
				require.NoError(t, h.db.QueryRow("SELECT lineage_id FROM "+table+" WHERE group_id='G1'").Scan(&stored))
				require.Equal(t, lineage, stored)
			}
			before := lifecycleReadState(t, h.db)
			// Keep the stale pre-create workspace revision: the first response was lost.
			replayed, err := ensureLineageKind(h, kind, mutation)
			require.NoError(t, err)
			require.Equal(t, first.Group.ID, replayed.Group.ID)
			require.Equal(t, lineage, requireLineageGenesis(t, h.db, replayed.Group))
			require.Equal(t, before, lifecycleReadState(t, h.db))
			require.NoError(t, h.db.Close())
			reopened := newCircuitProfileHarnessOn(t, openWorkspaceTestDB(t, path), h)
			reopened.db.SetMaxOpenConns(1)
			replayed, err = ensureLineageKind(reopened, kind, mutation)
			require.NoError(t, err)
			require.Equal(t, "G1", replayed.Group.ID)
			require.Equal(t, lineage, requireLineageGenesis(t, reopened.db, replayed.Group))
		})
	}
}

func TestCanonicalLineageCreateRollbackIsAtomic(t *testing.T) {
	for _, failure := range []string{"caller_rollback", "basic_key", "profile_key"} {
		t.Run(failure, func(t *testing.T) {
			h := newCircuitProfileHarness(t, ":memory:")
			h.db.SetMaxOpenConns(1)
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			group, _, err := circuitProfileCandidate(h.mutation(ctx, lineageTestDraft(h, "dual")))
			require.NoError(t, err)
			if failure != "caller_rollback" {
				table := "write_group_" + failure + "s"
				_, err = h.db.Exec("CREATE TRIGGER fail_key BEFORE INSERT ON " + table + " BEGIN SELECT RAISE(ABORT,'injected key failure'); END")
				require.NoError(t, err)
				_, err = ensureLineageKind(h, "dual", h.mutation(ctx, group))
				require.ErrorContains(t, err, "injected key failure")
			} else {
				tx, err := h.db.BeginTx(ctx, nil)
				require.NoError(t, err)
				require.NoError(t, h.service.repo.CreateInTx(ctx, tx, group))
				require.NoError(t, h.service.repo.insertBasicManagedKeyInTx(ctx, tx, h.fixture.workspaceID, h.fixture.deviceID, group.ID))
				var count int
				require.NoError(t, tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM write_group_lineages").Scan(&count))
				require.Equal(t, 1, count)
				require.NoError(t, tx.Rollback())
			}
			for _, table := range []string{"write_groups", "write_group_basic_keys", "write_group_profile_keys", "write_group_lineages", "write_group_generations"} {
				require.Zero(t, h.countRows("SELECT COUNT(*) FROM "+table), table)
			}
		})
	}
}

func TestCanonicalLineageDeletedEnsureNeverCreatesSuccessor(t *testing.T) {
	for _, kind := range []string{"basic", "profile"} {
		t.Run(kind, func(t *testing.T) {
			h := newCircuitProfileHarness(t, ":memory:")
			h.db.SetMaxOpenConns(1)
			mutation := h.mutation(t.Context(), lineageTestDraft(h, kind))
			first, err := ensureLineageKind(h, kind, mutation)
			require.NoError(t, err)
			_, err = h.db.Exec("UPDATE write_groups SET status='deleted' WHERE id=?", first.Group.ID)
			require.NoError(t, err)
			before := lifecycleReadState(t, h.db)
			_, err = ensureLineageKind(h, kind, mutation)
			require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
			require.Equal(t, before, lifecycleReadState(t, h.db))
			require.Equal(t, 1, h.countRows("SELECT COUNT(*) FROM write_groups"))
		})
	}
}

func lineageTestDraft(h *circuitProfileHarness, kind string) *WriteGroup {
	if kind == "basic" {
		return newBasicManagedDraft(h.fixture)
	}
	group := h.draft(WriteGroupProfileInstantaneous, "circuit-1", "Circuit 1", h.primary())
	// Basic/profile dual aliases use the profile's persisted interval from the
	// start; interval normalization is not part of this lineage change.
	group.RowPolicy.IntervalSeconds = group.RecordingProfile.IntervalSeconds()
	return group
}

func ensureLineageKind(h *circuitProfileHarness, kind string, mutation WriteGroupMutation) (*WriteGroupSaveResult, error) {
	if kind == "profile" {
		return h.service.EnsureCircuitProfile(h.t.Context(), mutation)
	}
	return h.service.EnsureBasicManaged(h.t.Context(), h.fixture.deviceID, mutation)
}

func lineageKeyTables(kind string) []string {
	if kind == "dual" {
		return []string{"write_group_basic_keys", "write_group_profile_keys"}
	}
	return []string{"write_group_" + kind + "_keys"}
}

func requireLineageGenesis(t *testing.T, db *sql.DB, group *WriteGroup) string {
	t.Helper()
	var lineage, workspace, current, genesis string
	var revision, generation int64
	var reserved, operation, intent, predecessor sql.NullString
	require.NoError(t, db.QueryRow(`SELECT l.lineage_id,l.workspace_id,l.current_group_id,l.revision,
		l.reserved_group_id,l.reservation_operation_id,l.reservation_intent_digest,g.group_id,g.generation,g.predecessor_group_id
		FROM write_group_lineages l JOIN write_group_generations g ON g.lineage_id=l.lineage_id WHERE g.group_id=?`, group.ID).
		Scan(&lineage, &workspace, &current, &revision, &reserved, &operation, &intent, &genesis, &generation, &predecessor))
	require.NotEmpty(t, lineage)
	require.Equal(t, group.WorkspaceID, workspace)
	require.Equal(t, group.ID, current)
	require.Equal(t, group.ID, genesis)
	require.EqualValues(t, 1, revision)
	require.EqualValues(t, 1, generation)
	for _, absent := range []sql.NullString{reserved, operation, intent, predecessor} {
		require.False(t, absent.Valid)
	}
	return lineage
}
