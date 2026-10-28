package groupdelivery

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func lifecycleRead(t *testing.T, db *sql.DB, key GroupKey, destination FrozenDestination) (LifecyclePending, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	require.NoError(t, err)
	defer tx.Rollback()
	return ReadLifecyclePending(ctx, tx, key, destination)
}

func TestLifecyclePendingMixedUnitsRemainSeparate(t *testing.T) {
	store, db, _ := newTestStore(t)
	enqueueFor(t, store, testKey, testDestination, "first", t0, 1)
	enqueueFor(t, store, testKey, testDestination, "second", at(10), 2)
	_, err := db.Exec(`UPDATE wg_delivery_outbox SET state='unknown' WHERE entity_key='second'`)
	require.NoError(t, err)
	for _, id := range []string{"one", "two", "three"} {
		require.NoError(t, store.AppendSample(t.Context(), testKey, at(20), sampleAt(id, "temp", 23, 1)))
	}
	db.SetMaxOpenConns(1)
	got, err := lifecycleRead(t, db, testKey, testDestination)
	require.NoError(t, err)
	require.Equal(t, 3, got.Samples)
	require.Equal(t, 2, got.Outbox.Total)
	require.Equal(t, 2, got.Outbox.Outstanding)
	require.Equal(t, 1, got.Outbox.Queued)
	require.Equal(t, 1, got.Outbox.Unknown)
	require.Empty(t, got.Handoffs)
	require.Equal(t, 2, countRows(t, db, "wg_delivery_outbox"))
	require.Equal(t, 3, countRows(t, db, "wg_delivery_samples"))
}

func TestLifecyclePendingHandoffPreservesOriginAndConsumingOwner(t *testing.T) {
	store := acceptedMetadataHandoff(t)
	store.DB().SetMaxOpenConns(1)
	for _, key := range []GroupKey{testKey, otherKey} {
		got, err := lifecycleRead(t, store.DB(), key, testDestination)
		require.NoError(t, err)
		want := 0
		if key == otherKey {
			want = 1
		}
		require.Equal(t, want, got.Samples)
		require.Equal(t, []LifecycleHandoff{{Origin: testKey, Source: testKey, Target: otherKey, State: "accepted", Samples: 1}}, got.Handoffs)
	}
}

func TestLifecyclePendingCorruptOrForeignEvidenceFailsClosed(t *testing.T) {
	for _, mutation := range []string{
		`UPDATE wg_delivery_outbox SET destination_scope='foreign'`,
		`UPDATE wg_delivery_outbox SET connector_revision='foreign'`,
		`UPDATE wg_delivery_outbox SET workspace_id='foreign'`,
		`UPDATE wg_delivery_outbox SET state='private-unknown-state'`,
	} {
		t.Run(mutation, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			enqueueFor(t, store, testKey, testDestination, "e", t0, 1)
			_, err := db.Exec(mutation)
			require.NoError(t, err)
			got, err := lifecycleRead(t, db, testKey, testDestination)
			require.ErrorIs(t, err, ErrOwnershipUnavailable)
			require.Equal(t, LifecyclePending{}, got)
		})
	}
	t.Run("handoff", func(t *testing.T) {
		store := acceptedMetadataHandoff(t)
		_, err := store.DB().Exec(`UPDATE wg_delivery_sample_handoffs SET record='{}'`)
		require.NoError(t, err)
		got, err := lifecycleRead(t, store.DB(), testKey, testDestination)
		require.ErrorIs(t, err, ErrOwnershipUnavailable)
		require.Equal(t, LifecyclePending{}, got)
	})
}

func TestLifecyclePendingUnresolvedHandoffsRemainVisible(t *testing.T) {
	for _, state := range []string{"pending", "quarantined"} {
		t.Run(state, func(t *testing.T) {
			store, _ := seedPreparedHandoff(t)
			third := addHandoffFanoutTarget(t, store)
			if state == "quarantined" {
				descriptor := handoffProofDescriptor(t, third, nil, true)
				descriptor = []byte(strings.ReplaceAll(string(descriptor), "src-1", "changed-source"))
				require.NoError(t, store.SaveRuntimeVersion(t.Context(), third, descriptor))
			}
			_, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
			require.Error(t, err)
			got, err := lifecycleRead(t, store.DB(), testKey, testDestination)
			require.NoError(t, err)
			require.Zero(t, got.Samples, "transferred-source exclusion must not hide unresolved audit")
			require.Len(t, got.Handoffs, 2)
			byTarget := map[GroupKey]LifecycleHandoff{}
			for _, edge := range got.Handoffs {
				byTarget[edge.Target] = edge
			}
			require.Equal(t, LifecycleHandoff{Origin: testKey, Source: testKey, Target: third, State: state, Samples: 1}, byTarget[third])
			require.Equal(t, "accepted", byTarget[otherKey].State)
		})
	}
}

func TestLifecyclePendingFinalizedAuditSurvivesSampleReclaim(t *testing.T) {
	store := acceptedMetadataHandoff(t)
	target, err := store.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	require.NoError(t, store.CommitClosure(t.Context(), otherKey, Closure{OwnershipView: target.OwnershipView, HandoffView: target.HandoffView, Destination: testDestination, NextClose: at(30)}))
	_, err = store.Reclaim(t.Context(), time.Hour)
	require.NoError(t, err)
	require.Zero(t, countRows(t, store.DB(), "wg_delivery_samples"))
	got, err := lifecycleRead(t, store.DB(), testKey, testDestination)
	require.NoError(t, err)
	require.Zero(t, got.Samples)
	require.Equal(t, []LifecycleHandoff{{Origin: testKey, Source: testKey, Target: otherKey, State: "finalized", Samples: 1}}, got.Handoffs)
}
