package grouppipeline

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"
)

func TestOwnerObservationOriginalCounterScopeAndExactInputs(t *testing.T) {
	for _, reason := range []string{"counter-replay-proof-missing", "counter-replay-proof-conflict", "counter-replay-proof-unsupported"} {
		t.Run(reason, func(t *testing.T) {
			f := ownerCounterFixture(t)
			require.NoError(t, f.reading(0, "second-input", at(2), 105))
			var out []OwnerObservation
			reads := 0
			ownerCapture(f.pipe, &out, &reads)
			damageCounterProof(t, f, reason)
			before := counterDurableState(t, f.store)
			f.pipe.boundaries = map[string]*managed{}
			// A mutable/latest group cannot relabel the persisted failed revision.
			f.groups.groups[0].AppliedRevision = "future-revision"
			require.NoError(t, f.pipe.Reconcile(t.Context()))
			require.NotEmpty(t, out)
			for _, o := range out {
				require.True(t, o.Available)
				require.Equal(t, groupKey(f.snapshot), o.Key)
				require.Equal(t, "observed-workspace-revision", o.WorkspaceRevision)
				require.Equal(t, "connector-1", o.ConnectorID)
				require.Equal(t, "crev-1", o.ConnectorRevision)
				require.NotEmpty(t, o.ObservationID)
				require.Equal(t, OwnerCounterQuarantine, o.Kind)
				require.Equal(t, reason, o.Reason)
				require.Equal(t, uint64(2), o.PendingSamples)
			}
			require.Positive(t, reads)
			require.Equal(t, before, counterDurableState(t, f.store))
		})
	}
}

func TestOwnerObservationLiveProofFailureRunsAfterLocks(t *testing.T) {
	f := ownerCounterFixture(t)
	m := f.pipe.boundaries[boundaryKey(f.group.ID, "rev-1")]
	var out []OwnerObservation
	reads := 0
	ownerCapture(f.pipe, &out, &reads)
	var entered atomic.Bool
	f.pipe.deps.OwnerObservation.Observe = func(o OwnerObservation) {
		out = append(out, o)
		if !entered.CompareAndSwap(false, true) {
			return
		}
		require.True(t, m.proofGate.TryLock())
		m.proofGate.Unlock()
		require.True(t, f.pipe.mu.TryLock())
		f.pipe.mu.Unlock()
		require.NoError(t, f.pipe.Reconcile(t.Context()), "actual callback reentry must not deadlock")
	}
	damageCounterProof(t, f, "counter-replay-proof-conflict")
	err := f.reading(0, "refused", at(2), 105)
	require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict)
	require.True(t, entered.Load())
	require.NotEmpty(t, out)
}

func TestOwnerObservationReconcileCallbackReentersAfterOuterLock(t *testing.T) {
	f := ownerCounterFixture(t)
	var out []OwnerObservation
	reads := 0
	ownerCapture(f.pipe, &out, &reads)
	var entered atomic.Bool
	f.pipe.deps.OwnerObservation.Observe = func(o OwnerObservation) {
		out = append(out, o)
		if entered.CompareAndSwap(false, true) {
			require.True(t, f.pipe.reconcileMu.TryLock())
			f.pipe.reconcileMu.Unlock()
			require.NoError(t, f.pipe.Reconcile(t.Context()))
		}
	}
	damageCounterProof(t, f, "counter-replay-proof-missing")
	f.pipe.boundaries = map[string]*managed{}
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	require.True(t, entered.Load())
	require.NotEmpty(t, out)
}

func TestOwnerObservationMissingChangedOrUnreadableProofUnavailable(t *testing.T) {
	for _, mode := range []string{"missing revision", "changed revision", "scope panic", "foreign workspace", "count error", "snapshot missing", "connector missing"} {
		t.Run(mode, func(t *testing.T) {
			f := ownerCounterFixture(t)
			var out []OwnerObservation
			reads := 0
			ownerCapture(f.pipe, &out, &reads)
			f.pipe.deps.OwnerObservation.ScopeRevision = func(_ context.Context, key groupdelivery.GroupKey) (string, bool) {
				reads++
				switch mode {
				case "missing revision":
					return "", true
				case "changed revision":
					if reads > 1 {
						return "changed", true
					}
				case "scope panic":
					panic("private-secret")
				case "foreign workspace":
					return "foreign", false
				}
				return "observed-workspace-revision", key.WorkspaceID == "workspace-A"
			}
			if mode == "snapshot missing" {
				_, err := f.local.Exec(`DELETE FROM wg_runtime_versions`)
				require.NoError(t, err)
			}
			if mode == "connector missing" {
				raw, err := f.store.RuntimeVersion(t.Context(), groupKey(f.snapshot))
				require.NoError(t, err)
				var v map[string]any
				require.NoError(t, json.Unmarshal(raw, &v))
				v["group"].(map[string]any)["destination"].(map[string]any)["connector_revision"] = ""
				raw, err = json.Marshal(v)
				require.NoError(t, err)
				_, err = f.local.Exec(`UPDATE wg_runtime_versions SET payload=?`, string(raw))
				require.NoError(t, err)
			}
			if mode == "count error" {
				reject := func(q string) error {
					if strings.Contains(q, "COUNT(*) FROM wg_delivery_samples WHERE workspace_id") {
						return errors.New("private count error")
					}
					return nil
				}
				f.localFault.Reject.Store(&reject)
				defer f.localFault.Reject.Store(nil)
			}
			require.True(t, f.pipe.setCounterQuarantine(groupKey(f.snapshot), "counter-replay-proof-missing"))
			requireOwnerUnavailable(t, out)
		})
	}
}

func TestOwnerObservationTokenlessHandoffRefusalUnavailable(t *testing.T) {
	f, _, _ := originalConflictFixture(t, shapeMeanLong, "epoch")
	f.clock = at(7)
	require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "future-refused", at(11), uint64(100))))
	f.commit()
	f.source.EndsAt = f.boundary
	f.groups.groups = append(f.groups.groups, f.target.Group)
	f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
	var out []OwnerObservation
	reads := 0
	ownerCapture(f.pipe, &out, &reads)
	f.pipe.deps.OwnerObservation.ScopeRevision = func(context.Context, groupdelivery.GroupKey) (string, bool) {
		reads++
		return "observed-workspace-revision", true
	}
	before := handoffOriginalRows(t, f)
	err := f.pipe.PrepareCutoverHandoffs(t.Context(), f.source.WorkspaceID)
	require.ErrorIs(t, err, groupdelivery.ErrHandoffProofConflict)
	require.NotEmpty(t, out)
	requireOwnerUnavailable(t, out)
	require.Equal(t, before, handoffOriginalRows(t, f))
}

func TestOwnerObservationNilAndPanicPreserveNativeBehavior(t *testing.T) {
	for _, observer := range []bool{false, true} {
		t.Run(map[bool]string{false: "nil", true: "panic"}[observer], func(t *testing.T) {
			f := ownerCounterFixture(t)
			reads := 0
			f.pipe.deps.OwnerObservation.RecordLimit = 1024
			f.pipe.deps.OwnerObservation.ScopeRevision = func(context.Context, groupdelivery.GroupKey) (string, bool) {
				reads++
				return "observed-workspace-revision", true
			}
			if observer {
				f.pipe.deps.OwnerObservation.Observe = func(OwnerObservation) { panic("private-secret") }
			}
			damageCounterProof(t, f, "counter-replay-proof-conflict")
			before := counterDurableState(t, f.store)
			require.ErrorIs(t, f.reading(0, "refused", at(2), 105), counterproof.ErrCounterBindingProofConflict)
			require.Equal(t, before, counterDurableState(t, f.store))
			if !observer {
				require.Zero(t, reads)
			} else {
				require.Positive(t, reads)
			}
		})
	}
}

func TestOwnerObservationIgnoresUnclassifiedOrIncomingClaims(t *testing.T) {
	f := ownerCounterFixture(t)
	var out []OwnerObservation
	reads := 0
	ownerCapture(f.pipe, &out, &reads)
	require.False(t, f.pipe.setCounterQuarantine(groupKey(f.snapshot), "private-endpoint-password"))
	bad := envelope(f.group, "wrong-epoch", at(2), uint64(105))
	bad.SeriesEpoch = "foreign-secret"
	require.Error(t, f.pipe.AcceptSample(t.Context(), bad))
	require.Empty(t, out)
}

func TestOwnerObservationBatchOverflowIsBoundedAndLossAfterUnlock(t *testing.T) {
	f := ownerCounterFixture(t)
	var out []OwnerObservation
	reads := 0
	ownerCapture(f.pipe, &out, &reads)
	f.pipe.deps.OwnerObservation.RecordLimit = 1
	calls := 0
	var lost uint64
	f.pipe.deps.OwnerObservation.OnDrop = func(n uint64) {
		calls++
		lost += n
		require.True(t, f.pipe.mu.TryLock())
		f.pipe.mu.Unlock()
		require.True(t, f.pipe.reconcileMu.TryLock())
		f.pipe.reconcileMu.Unlock()
	}
	before := counterDurableState(t, f.store)
	ctx, finish := f.pipe.ownerOperation(t.Context())
	for range 5 {
		require.True(t, f.pipe.observedCounterQuarantine(ctx, groupKey(f.snapshot), "counter-replay-proof-missing"))
	}
	require.Empty(t, out)
	require.Zero(t, reads)
	b := ctx.Value(ownerBatchContextKey{}).(*ownerObservationBatch)
	require.Len(t, b.descriptors, 1)
	require.Equal(t, uint64(4), b.dropped)
	finish()
	require.Len(t, out, 1)
	require.Equal(t, 1, calls)
	require.Equal(t, uint64(4), lost)
	require.Equal(t, before, counterDurableState(t, f.store))
	b.dropped = ^uint64(0)
	b.drop()
	require.Equal(t, ^uint64(0), b.dropped)
}

func TestOwnerObservationChangedFrozenOrQuarantineCannotBorrowProof(t *testing.T) {
	for _, mode := range []string{"frozen", "recovered"} {
		t.Run(mode, func(t *testing.T) {
			f := ownerCounterFixture(t)
			var out []OwnerObservation
			reads := 0
			ownerCapture(f.pipe, &out, &reads)
			f.pipe.deps.OwnerObservation.ScopeRevision = func(context.Context, groupdelivery.GroupKey) (string, bool) {
				reads++
				if reads == 2 {
					if mode == "recovered" {
						f.pipe.mu.Lock()
						delete(f.pipe.quarantines, groupKey(f.snapshot))
						f.pipe.mu.Unlock()
					} else {
						_, err := f.local.Exec(`UPDATE wg_runtime_versions SET payload=replace(payload,'crev-1','crev-2')`)
						require.NoError(t, err)
					}
				}
				return "observed-workspace-revision", true
			}
			require.True(t, f.pipe.setCounterQuarantine(groupKey(f.snapshot), "counter-replay-proof-missing"))
			requireOwnerUnavailable(t, out)
		})
	}
}

func TestOwnerObservationHealthyAcquisitionDoesNotReadDiagnosticMetadata(t *testing.T) {
	f := ownerCounterFixture(t)
	var out []OwnerObservation
	reads := 0
	ownerCapture(f.pipe, &out, &reads)
	require.NoError(t, f.reading(0, "healthy-input", at(2), 105))
	require.Zero(t, reads)
	require.Empty(t, out)
}

func TestOwnerObservationHandoffOwnerChangedBeforeFlushUnavailable(t *testing.T) {
	f, _, _ := originalConflictFixture(t, shapeMeanLong, "epoch")
	f.clock = at(7)
	require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "future-owner-change", at(11), uint64(100))))
	f.commit()
	f.source.EndsAt = f.boundary
	var out []OwnerObservation
	reads := 0
	ownerCapture(f.pipe, &out, &reads)
	f.pipe.deps.OwnerObservation.ScopeRevision = func(context.Context, groupdelivery.GroupKey) (string, bool) {
		return "observed-workspace-revision", true
	}
	original, err := f.store.OwnershipView(t.Context(), groupKey(f.source))
	require.NoError(t, err)
	ctx, finish := f.pipe.ownerOperation(t.Context())
	_, err = f.pipe.prepareWorkspaceHandoffs(ctx, f.source.WorkspaceID)
	require.ErrorIs(t, err, groupdelivery.ErrHandoffProofConflict)
	require.Empty(t, out)
	_, err = f.local.Exec(`UPDATE write_group_versions SET payload=json_set(payload,'$.row_policy.interval_seconds',120) WHERE group_id=? AND group_revision=?`, f.source.GroupID, f.source.AppliedRevision)
	require.NoError(t, err)
	changed, err := f.store.OwnershipView(t.Context(), groupKey(f.source))
	require.NoError(t, err)
	require.NotEqual(t, original, changed)
	finish()
	requireOwnerUnavailable(t, out)
}

func TestOwnerObservationNativeTokenBackedRefusalPreservesReason(t *testing.T) {
	f, _, env := ownerRefusalFixture(t)
	var out []OwnerObservation
	reads := 0
	ownerCapture(f.pipe, &out, &reads)
	f.pipe.deps.OwnerObservation.ScopeRevision = func(context.Context, groupdelivery.GroupKey) (string, bool) {
		reads++
		return "observed-workspace-revision", true
	}
	before := counterDurableState(t, f.store)
	require.ErrorIs(t, f.pipe.AcceptSample(t.Context(), env), groupdelivery.ErrOwnershipConflict)
	require.Len(t, out, 1)
	o := out[0]
	require.True(t, o.Available)
	require.Equal(t, OwnerRecordingRefusal, o.Kind)
	require.Equal(t, groupKey(f.target), o.Key)
	require.Equal(t, groupdelivery.ErrOwnershipConflict.Error(), o.Reason)
	require.NotEmpty(t, o.ObservationID)
	require.Equal(t, uint64(0), o.PendingSamples)
	require.Equal(t, f.target.Group.Destination.ConnectorID, o.ConnectorID)
	require.Equal(t, f.target.Group.Destination.ConnectorRevision, o.ConnectorRevision)
	require.Equal(t, before, counterDurableState(t, f.store))
}

func TestOwnerObservationHeldNativeTokenChangedBeforeFlushUnavailable(t *testing.T) {
	f, target, env := ownerRefusalFixture(t)
	var out []OwnerObservation
	reads := 0
	ownerCapture(f.pipe, &out, &reads)
	f.pipe.deps.OwnerObservation.ScopeRevision = func(context.Context, groupdelivery.GroupKey) (string, bool) {
		return "observed-workspace-revision", true
	}
	ctx, finish := f.pipe.ownerOperation(t.Context())
	require.ErrorIs(t, f.pipe.acceptWithOwnership(ctx, env, []*managed{target}), groupdelivery.ErrOwnershipConflict)
	require.Empty(t, out)
	_, err := f.local.Exec(`UPDATE write_group_versions SET payload=json_set(payload,'$.row_policy.interval_seconds',120) WHERE group_id=? AND group_revision=?`, f.target.GroupID, f.target.AppliedRevision)
	require.NoError(t, err)
	changed, err := f.store.OwnershipView(t.Context(), groupKey(f.target))
	require.NoError(t, err)
	require.NotEqual(t, target.boundary.ReplayState().OwnershipView, changed)
	finish()
	requireOwnerUnavailable(t, out)
}
