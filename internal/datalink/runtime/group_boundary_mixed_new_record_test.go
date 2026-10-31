package runtime

import (
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestMixedIndependentRecordDoesNotInheritDurableBaseline(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shared-journal.db")
	db := openDurableDB(t, path)
	old := mixedCounterBoundaryFixture(t, false)
	store := mixedDurableFixture(t, old, db)
	a, err := NewGroupBoundaryContext(t.Context(), old.config)
	require.NoError(t, err)
	old.accept(t, a, old.envelope(0, "old-base", boundaryAt(14), uint64(100)))
	require.NoError(t, a.Tick(t.Context(), boundaryAt(15)))
	old.accept(t, a, old.envelope(0, "old-tail", boundaryAt(29), uint64(108)))
	require.NoError(t, a.Tick(t.Context(), boundaryAt(30)))
	var oldPayload, oldDigest string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT payload,payload_digest FROM wg_delivery_counter_baselines WHERE group_id=?`, old.config.Group.ID).Scan(&oldPayload, &oldDigest))
	fresh := mixedCounterBoundaryFixture(t, false)
	fresh.config.Group.ID = "group-new-record"
	fresh.config.Group.RecordID = "98338991-859d-4a59-ab66-de64539db56c"
	// A new independently accepted group requires its own sealed binding. It
	// cannot attach the old group's proof or baseline just because source IDs match.
	for _, previous := range fresh.config.CounterBindings {
		facts := previous.CounterBindingFacts
		facts.Member.GroupID = fresh.config.Group.ID
		proof, err := workspace.NewCounterBindingProof(facts)
		require.NoError(t, err)
		fresh.config.CounterBindings = workspace.CounterBindingProofs{proof.Member.Key(): proof}
	}
	mixedDurableFixture(t, fresh, db)
	key := groupdelivery.GroupKey{WorkspaceID: fresh.config.Group.WorkspaceID, GroupID: fresh.config.Group.ID, GroupRevision: fresh.config.Group.AppliedRevision}
	empty, err := store.Restore(t.Context(), key)
	require.NoError(t, err)
	require.Empty(t, empty.Baselines)
	b, err := NewGroupBoundaryContext(t.Context(), fresh.config)
	require.NoError(t, err)
	fresh.accept(t, b, fresh.envelope(0, "new-first", boundaryAt(14), uint64(1000)))
	require.NoError(t, b.Tick(t.Context(), boundaryAt(15)))
	var rows int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count(*) FROM wg_delivery_outbox WHERE group_id=?`, key.GroupID).Scan(&rows))
	require.Zero(t, rows, "new first window is baseline-only, not 1000 minus old108")
	require.NoError(t, db.Close())
	db = openDurableDB(t, path)
	mixedDurableFixture(t, fresh, db)
	b, err = NewGroupBoundaryContext(t.Context(), fresh.config)
	require.NoError(t, err)
	fresh.accept(t, b, fresh.envelope(0, "new-second", boundaryAt(29), uint64(1003)))
	require.NoError(t, b.Tick(t.Context(), boundaryAt(30)))
	var raw []byte
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_outbox WHERE group_id=?`, key.GroupID).Scan(&raw))
	row, err := groupdelivery.DecodeRowPayload(raw)
	require.NoError(t, err)
	var value any
	for _, cell := range row.Cells {
		if cell.Column == "value_decimal" {
			value = cell.Value
		}
	}
	require.Equal(t, "3", value, "only the new record's own first reading may seed delta")
	var afterPayload, afterDigest string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT payload,payload_digest FROM wg_delivery_counter_baselines WHERE group_id=?`, old.config.Group.ID).Scan(&afterPayload, &afterDigest))
	require.Equal(t, oldPayload, afterPayload)
	require.Equal(t, oldDigest, afterDigest)
}
