package groupdelivery

import (
	"fmt"
	"testing"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

func TestCounterHandoffClosurePreservesOriginalProof(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		for _, damaged := range []bool{false, true} {
			t.Run(fmt.Sprintf("legacy=%t/damaged=%t", legacy, damaged), func(t *testing.T) {
				store, _, _, targetProof := handoffCounterProofFixture(t, legacy)
				db := store.DB()
				seedOwnershipCutover(t, db, "counter-transfer", testKey, otherKey, at(2))
				batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
				require.NoError(t, err)
				require.Equal(t, 1, batch.Accepted)
				target, err := store.Restore(t.Context(), otherKey)
				require.NoError(t, err)
				require.Len(t, target.Samples, 1)
				ref, err := targetProof.Reference()
				require.NoError(t, err)
				require.True(t, target.AuthorizesHandoff(target.Samples[0], ref))

				if damaged {
					_, err = db.ExecContext(t.Context(), `UPDATE wg_runtime_versions SET payload=REPLACE(payload, 'kWh', 'MWh') WHERE group_id=?`, otherKey.GroupID)
					require.NoError(t, err)
				}
				bucket := silentBucket("", t0)
				bucket.Outcome.WindowStart, bucket.Outcome.WindowEnd = at(2), at(10)
				bucket.Outcome.Truncated = "start"
				closure := Closure{OwnershipView: target.OwnershipView, HandoffView: target.HandoffView,
					Destination: testDestination, NextClose: at(10), Buckets: []ClosedBucket{bucket},
					CounterProofs: []counterproof.CounterBindingReference{ref}}
				err = store.CommitClosure(t.Context(), otherKey, closure)
				if damaged {
					require.Error(t, err)
					require.Zero(t, openCount(t, db, "wg_delivery_samples", "consumed=1"))
					require.Zero(t, countRows(t, db, "wg_delivery_checkpoints"))
					return
				}
				require.NoError(t, err)
				var consumed int
				require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_samples WHERE consumed=1 AND group_id=?`, otherKey.GroupID).Scan(&consumed))
				require.Equal(t, 1, consumed)
			})
		}
	}
}

func TestCounterDuplicateUsesNativeLogicalIdentityAfterHandoff(t *testing.T) {
	store, _, _, targetProof := handoffCounterProofFixture(t, false)
	db := store.DB()
	seedOwnershipCutover(t, db, "counter-transfer", testKey, otherKey, at(2))
	_, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	target, err := store.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	require.Len(t, target.Samples, 1)
	var handoffAddress string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT sample_id FROM wg_delivery_samples WHERE group_id=? AND acceptance_kind='handoff'`, otherKey.GroupID).Scan(&handoffAddress))
	ref, err := targetProof.Reference()
	require.NoError(t, err)
	native := target.Samples[0]
	native.AcceptanceOrigin = snapshot.AcceptanceOrigin{}
	native.SampleID, native.CounterProof, native.ObservedAt = handoffAddress, ref, at(3)
	require.NoError(t, store.AppendSample(t.Context(), otherKey, t0, native))
	require.NoError(t, store.VerifyCounterDuplicate(t.Context(), otherKey, native))
	require.Equal(t, 3, countRows(t, db, "wg_delivery_samples"), "source, handoff and native samples remain distinct")
}
