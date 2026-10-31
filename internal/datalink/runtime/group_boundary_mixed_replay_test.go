package runtime

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
)

type capturedMixedClosureLedger struct {
	*groupdelivery.Ledger
	last groupdelivery.Closure
}

func (l *capturedMixedClosureLedger) CommitClosure(ctx context.Context, c groupdelivery.Closure) error {
	if err := l.Ledger.CommitClosure(ctx, c); err != nil {
		return err
	}
	l.last = c
	return nil
}

func TestMixedDurableAcceptedPayloadReplayConflict(t *testing.T) {
	f := mixedCounterBoundaryFixture(t, true)
	db := openDurableDB(t, filepath.Join(t.TempDir(), "journal.db"))
	store := mixedDurableFixture(t, f, db)
	g := f.config.Group
	key := groupdelivery.GroupKey{WorkspaceID: g.WorkspaceID, GroupID: g.ID, GroupRevision: g.AppliedRevision}
	capture := &capturedMixedClosureLedger{Ledger: store.Ledger(key)}
	f.config.Ledger = capture
	b, err := NewGroupBoundaryContext(t.Context(), f.config)
	require.NoError(t, err)
	f.accept(t, b, f.envelope(0, "a", boundaryAt(14), uint64(100)))
	require.NoError(t, b.Tick(t.Context(), boundaryAt(15)))
	f.accept(t, b, f.envelope(0, "b", boundaryAt(29), uint64(108)))
	require.NoError(t, b.Tick(t.Context(), boundaryAt(30)))
	require.NoError(t, store.CommitClosure(t.Context(), key, capture.last), "original accepted operation is idempotent")
	raw, err := json.Marshal(capture.last)
	require.NoError(t, err)
	var changed groupdelivery.Closure
	require.NoError(t, json.Unmarshal(raw, &changed))
	require.Len(t, changed.Buckets, 1)
	effect := changed.Buckets[0].Row.EffectKey
	before, err := store.GetOutbox(t.Context(), effect)
	require.NoError(t, err)
	changedValue := false
	for i, cell := range changed.Buckets[0].Row.Cells {
		if cell.Column == "value_decimal" {
			changed.Buckets[0].Row.Cells[i].Value = "999"
			changedValue = true
		}
	}
	require.True(t, changedValue)
	require.Error(t, store.CommitClosure(t.Context(), key, changed), "same effect cannot accept a different finalized payload")
	after, err := store.GetOutbox(t.Context(), effect)
	require.NoError(t, err)
	require.Equal(t, before, after)
	require.Equal(t, 1, countTable(t, db, "wg_delivery_outbox"))
	require.NoError(t, store.CommitClosure(t.Context(), key, capture.last), "conflicting replay must not poison original replay")
}
