package grouppipeline

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"
)

func TestIntakeObservationDoesNotActivateFutureOrDrainOnlyRevision(t *testing.T) {
	f := newFixture(t)
	f.clock = at(12)
	first, second := groupAt("rev-1"), groupAt("rev-2")
	f.groups.groups = []*workspace.WriteGroup{liveGroup(second, workspace.WriteGroupStatusReady)}
	f.groups.snapshots[first.ID] = []*workspace.WriteGroupAppliedSnapshot{snapshotOf(first, at(0)), snapshotOf(second, at(20))}
	p := f.pipeline()
	var out []IntakeQuotaObservation
	reads := 0
	intakeCapture(p, &out, &reads)
	require.NoError(t, p.Reconcile(t.Context()))
	applied := intakeKind(out, IntakeGroupApplied)
	require.Empty(t, applied, "native Status is retiring once its future end is installed")
	retiring := intakeKind(out, IntakeGroupRetiring)
	require.Len(t, retiring, 1)
	require.Equal(t, "rev-1", retiring[0].Key.GroupRevision)
	out = nil
	p.boundaries[boundaryKey(first.ID, "rev-1")].drainOnly = true
	require.NoError(t, p.Reconcile(t.Context()))
	require.Empty(t, intakeKind(out, IntakeGroupApplied))
}
func TestIntakeObservationDisableRetainsHistoricalBacklog(t *testing.T) {
	f := newFixture(t)
	p := f.pipeline()
	require.NoError(t, p.Reconcile(t.Context()))
	group := f.groups.groups[0]
	require.NoError(t, p.AcceptSample(t.Context(), envelope(group, "accepted", at(1), 1.0)))
	before, err := f.store.PendingSampleCount(t.Context(), groupdelivery.GroupKey{WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.AppliedRevision})
	require.NoError(t, err)
	var out []IntakeQuotaObservation
	reads := 0
	intakeCapture(p, &out, &reads)
	f.groups.groups = []*workspace.WriteGroup{liveGroup(group, workspace.WriteGroupStatusDisabled)}
	require.NoError(t, p.Reconcile(t.Context()))
	paused := intakeKind(out, IntakeGroupPaused)
	require.Len(t, paused, 1)
	require.Equal(t, "not_running", paused[0].IntakeState)
	require.NotEmpty(t, intakeKind(out, IntakeGroupRetiring))
	after, err := f.store.PendingSampleCount(t.Context(), paused[0].Key)
	require.NoError(t, err)
	require.Equal(t, before, after)
	out = nil
	f.groups.groups = nil
	require.NoError(t, p.Reconcile(t.Context()))
	require.Empty(t, intakeKind(out, IntakeGroupPaused), "absence is not pause proof")
}
func TestIntakeObservationRevisionBoundaryKeepsOldAndNewProofSeparate(t *testing.T) {
	f := newFixture(t)
	p := f.pipeline()
	require.NoError(t, p.Reconcile(t.Context()))
	first, second := groupAt("rev-1"), groupAt("rev-2")
	f.groups.groups = []*workspace.WriteGroup{liveGroup(second, workspace.WriteGroupStatusReady)}
	f.groups.snapshots[first.ID] = []*workspace.WriteGroupAppliedSnapshot{snapshotOf(first, at(0)), snapshotOf(second, at(20))}
	var out []IntakeQuotaObservation
	reads := 0
	intakeCapture(p, &out, &reads)
	f.clock = at(21)
	require.NoError(t, p.Reconcile(t.Context()))
	applied := intakeKind(out, IntakeGroupApplied)
	require.Len(t, applied, 1)
	require.Equal(t, "rev-2", applied[0].Key.GroupRevision)
	retired := intakeKind(out, IntakeGroupRetiring)
	require.Len(t, retired, 1)
	require.Equal(t, "rev-1", retired[0].Key.GroupRevision)
}
func TestIntakeObservationMissingChangedOrForeignProofIsUnavailable(t *testing.T) {
	for _, mode := range []string{"missing", "changed", "foreign", "panic"} {
		t.Run(mode, func(t *testing.T) {
			f := newFixture(t)
			p := f.pipeline()
			var out []IntakeQuotaObservation
			reads := 0
			intakeCapture(p, &out, &reads)
			p.deps.IntakeQuotaObservation.ScopeRevision = func(context.Context, groupdelivery.GroupKey) (string, bool) {
				reads++
				switch mode {
				case "missing":
					return "", true
				case "foreign":
					return "w", false
				case "panic":
					panic("private-password")
				case "changed":
					if reads > 1 {
						return "next", true
					}
				}
				return "workspace-revision", true
			}
			require.NoError(t, p.Reconcile(t.Context()))
			require.NotEmpty(t, out)
			for _, o := range out {
				require.Equal(t, IntakeQuotaObservation{}, o)
			}
		})
	}
}
func TestQuotaObservationPreservesOriginalGlobalAndGroupRefusal(t *testing.T) {
	for _, scope := range []groupdelivery.QuotaScope{groupdelivery.QuotaScopeGlobal, groupdelivery.QuotaScopeGroup} {
		t.Run(string(scope), func(t *testing.T) {
			f, p, out, _ := intakeQuotaFixture(t, scope)
			err := p.AcceptSample(t.Context(), envelope(f.groups.groups[0], "refused", at(2), 2.0))
			require.ErrorIs(t, err, groupdelivery.ErrQuotaExceeded)
			var q *groupdelivery.QuotaError
			require.True(t, errors.As(err, &q))
			observed := intakeKind(*out, IntakeQuotaRefused)
			require.Len(t, observed, 1)
			require.Equal(t, q.Scope, observed[0].Quota.Scope)
			require.Equal(t, q.UsedBytes, observed[0].Quota.UsedBytes)
			require.Equal(t, q.MaxBytes, observed[0].Quota.MaxBytes)
			require.Greater(t, q.UsedBytes, q.MaxBytes, "native accepted payload overshoot remains observed")
			require.Empty(t, observed[0].Quota.LossRiskNotice)
		})
	}
}
func TestQuotaObservationRecoveryRequiresExactScopeMeasurement(t *testing.T) {
	f, p, out, reads := intakeQuotaFixture(t, groupdelivery.QuotaScopeGroup)
	p.deps.IntakeQuotaObservation.IntakeEnabled = false
	hint := QuotaRecoveryHint{Key: groupdelivery.GroupKey{WorkspaceID: "workspace-A", GroupID: "group-G", GroupRevision: "rev-1"}, WorkspaceRevision: "workspace-revision", ConnectorID: "connector-1", ConnectorRevision: "crev-1", ReferenceID: "original-ref", Scope: groupdelivery.QuotaScopeGroup}
	var recovered []IntakeQuotaObservation
	hint.Observe = func(o IntakeQuotaObservation) { recovered = append(recovered, o) }
	p.deps.IntakeQuotaObservation.TakeQuotaRecoveryHint = func(context.Context) (QuotaRecoveryHint, bool) { return hint, true }
	require.NoError(t, p.Reconcile(t.Context()))
	require.Empty(t, recovered, "hard limit is not recovered")
	_, err := f.store.DB().Exec(`UPDATE wg_delivery_samples SET consumed=1`)
	require.NoError(t, err)
	require.NoError(t, p.Reconcile(t.Context()))
	require.Len(t, recovered, 1)
	require.Equal(t, IntakeQuotaRecovered, recovered[0].Kind)
	require.Zero(t, recovered[0].Quota.UsedBytes)
	require.Equal(t, groupdelivery.QuotaScopeGroup, recovered[0].Quota.Scope)
	require.False(t, recovered[0].Quota.IntakeRefused)
	hint.WorkspaceRevision = "old"
	recovered = nil
	require.NoError(t, p.Reconcile(t.Context()))
	require.Empty(t, recovered)
	require.Contains(t, *out, IntakeQuotaObservation{})
	*reads = 0
	p.deps.IntakeQuotaObservation.TakeQuotaRecoveryHint = func(context.Context) (QuotaRecoveryHint, bool) { return QuotaRecoveryHint{}, false }
	require.NoError(t, p.Reconcile(t.Context()))
	require.Zero(t, *reads)
}
func TestQuotaObservationDuplicateOrIgnoredSampleCannotRecover(t *testing.T) {
	f, p, out, reads := intakeQuotaFixture(t, groupdelivery.QuotaScopeGroup)
	hints := 0
	p.deps.IntakeQuotaObservation.TakeQuotaRecoveryHint = func(context.Context) (QuotaRecoveryHint, bool) { hints++; return QuotaRecoveryHint{}, false }
	require.NoError(t, p.AcceptSample(t.Context(), envelope(f.groups.groups[0], "first", at(1), 1.0)))
	require.Empty(t, *out)
	require.Zero(t, *reads)
	require.Zero(t, hints)
}
func TestIntakeQuotaObservationCallbacksRunAfterOuterLocks(t *testing.T) {
	f := newFixture(t)
	p := f.pipeline()
	var out []IntakeQuotaObservation
	reads := 0
	intakeCapture(p, &out, &reads)
	entered := false
	p.deps.IntakeQuotaObservation.Observe = func(o IntakeQuotaObservation) {
		out = append(out, o)
		if entered {
			return
		}
		entered = true
		require.True(t, p.reconcileMu.TryLock())
		p.reconcileMu.Unlock()
		require.NoError(t, p.Reconcile(t.Context()))
		panic("capture failure")
	}
	require.NoError(t, p.Reconcile(t.Context()))
	require.True(t, entered)
	require.Equal(t, StateActive, statusOf(p).State)
}

func TestQuotaObservationZeroOvershootAndDiskFullRemainTruthful(t *testing.T) {
	f := newFixture(t)
	p := f.pipeline()
	require.NoError(t, p.Reconcile(t.Context()))
	var out []IntakeQuotaObservation
	reads := 0
	intakeCapture(p, &out, &reads)
	m := p.boundaries[boundaryKey("group-G", "rev-1")]
	ctx, finish := p.ownerOperation(t.Context())
	expected := &groupdelivery.QuotaError{Reason: groupdelivery.ReasonDiskFull, Scope: groupdelivery.QuotaScopeGlobal, UsedBytes: 123}
	require.Same(t, expected, p.observedProofGate(ctx, m, func() error { return expected }))
	finish()
	require.Equal(t, []IntakeQuotaObservation{{}}, out, "max0 disk-full is not configured quota evidence")
	require.Equal(t, StateActive, statusOf(p).State)
}

func TestIntakeObservationPostScopeFrozenAndNativeChangesFailClosed(t *testing.T) {
	for _, mode := range []string{"frozen", "pointer", "state"} {
		t.Run(mode, func(t *testing.T) {
			f := newFixture(t)
			p := f.pipeline()
			var out []IntakeQuotaObservation
			reads := 0
			intakeCapture(p, &out, &reads)
			p.deps.IntakeQuotaObservation.ScopeRevision = func(ctx context.Context, key groupdelivery.GroupKey) (string, bool) {
				reads++
				if reads == 2 {
					switch mode {
					case "frozen":
						raw, err := f.store.RuntimeVersion(ctx, key)
						require.NoError(t, err)
						_, err = f.store.DB().Exec(`UPDATE wg_runtime_versions SET payload=?`, strings.Replace(string(raw), "crev-1", "crev-2", 1))
						require.NoError(t, err)
					case "pointer":
						p.mu.Lock()
						p.boundaries[boundaryKey(key.GroupID, key.GroupRevision)] = &managed{}
						p.mu.Unlock()
					case "state":
						p.boundaries[boundaryKey(key.GroupID, key.GroupRevision)].boundary.SetUntil(at(10))
					}
				}
				return "workspace-revision", true
			}
			require.NoError(t, p.Reconcile(t.Context()))
			require.Equal(t, []IntakeQuotaObservation{{}}, out)
		})
	}
}

func TestIntakeQuotaObservationSharedCapFiltersBeforeAllocation(t *testing.T) {
	f := newFixture(t)
	p := f.pipeline()
	require.NoError(t, p.Reconcile(t.Context()))
	var out []IntakeQuotaObservation
	reads := 0
	intakeCapture(p, &out, &reads)
	p.deps.IntakeQuotaObservation.RecordLimit = 1
	p.deps.IntakeQuotaObservation.IntakeEnabled = false
	var dropped uint64
	calls := 0
	p.deps.IntakeQuotaObservation.OnDrop = func(n uint64) { calls++; dropped += n }
	p.deps.OwnerObservation = OwnerObservationConfig{RecordLimit: 1, Observe: func(OwnerObservation) {}, OnDrop: func(n uint64) { calls++; dropped += n }}
	ctx, finish := p.ownerOperation(t.Context())
	p.recordIntakeObservation(ctx, intakeDescriptor{kind: IntakeGroupApplied})
	p.recordOwnerObservation(ctx, ownerDescriptor{})
	p.recordIntakeObservation(ctx, intakeDescriptor{kind: IntakeQuotaRefused})
	b := ctx.Value(ownerBatchContextKey{}).(*ownerObservationBatch)
	require.Len(t, b.descriptors, 1)
	finish()
	require.Equal(t, uint64(1), dropped)
	require.Equal(t, 1, calls, "one shared loss owner")
	require.Zero(t, reads)
	require.Empty(t, out)
}

func TestIntakeQuotaObservationNilDisabledAndPanicsPreserveNativeData(t *testing.T) {
	f := newFixture(t)
	p := f.pipeline()
	var out []IntakeQuotaObservation
	reads := 0
	intakeCapture(p, &out, &reads)
	p.deps.IntakeQuotaObservation.IntakeEnabled = false
	p.deps.IntakeQuotaObservation.QuotaEnabled = false
	p.deps.IntakeQuotaObservation.TakeQuotaRecoveryHint = func(context.Context) (QuotaRecoveryHint, bool) { panic("must not run") }
	require.NoError(t, p.Reconcile(t.Context()))
	require.Zero(t, reads)
	require.Empty(t, out)
	p.deps.IntakeQuotaObservation.IntakeEnabled = true
	p.deps.IntakeQuotaObservation.Observe = func(IntakeQuotaObservation) { panic("capture panic") }
	p.deps.IntakeQuotaObservation.OnDrop = func(uint64) { panic("loss panic") }
	require.NoError(t, p.Reconcile(t.Context()))
	require.Equal(t, StateActive, statusOf(p).State)
	require.NoError(t, p.AcceptSample(t.Context(), envelope(f.groups.groups[0], "native", at(1), 1.0)))
	count, err := f.store.PendingSampleCount(t.Context(), groupdelivery.GroupKey{WorkspaceID: "workspace-A", GroupID: "group-G", GroupRevision: "rev-1"})
	require.NoError(t, err)
	require.EqualValues(t, 1, count)
}

func TestIntakeObservationAbsentGroupIsNotPausedAndDeletedRequiresTransition(t *testing.T) {
	for _, deleted := range []bool{false, true} {
		t.Run(map[bool]string{false: "absent", true: "deleted"}[deleted], func(t *testing.T) {
			f := newFixture(t)
			p := f.pipeline()
			require.NoError(t, p.Reconcile(t.Context()))
			var out []IntakeQuotaObservation
			reads := 0
			intakeCapture(p, &out, &reads)
			if deleted {
				f.groups.groups[0] = liveGroup(f.groups.groups[0], workspace.WriteGroupStatusDeleted)
			} else {
				f.groups.groups = nil
			}
			require.NoError(t, p.Reconcile(t.Context()))
			require.Len(t, intakeKind(out, IntakeGroupRetiring), 1)
			if deleted {
				require.Len(t, intakeKind(out, IntakeGroupPaused), 1)
			} else {
				require.Empty(t, intakeKind(out, IntakeGroupPaused))
			}
			out = nil
			require.NoError(t, p.Reconcile(t.Context()))
			require.Empty(t, intakeKind(out, IntakeGroupPaused), "no new transition is not a new paused fact")
		})
	}
}

func TestQuotaObservationRecoveryRejectsWrongScopeAndFrozenIdentity(t *testing.T) {
	for _, mode := range []string{"unconfigured-global", "foreign-connector", "missing-frozen", "changed-after-scope", "scope-panic"} {
		t.Run(mode, func(t *testing.T) {
			f, p, out, _ := intakeQuotaFixture(t, groupdelivery.QuotaScopeGroup)
			p.deps.IntakeQuotaObservation.IntakeEnabled = false
			_, err := f.store.DB().Exec(`UPDATE wg_delivery_samples SET consumed=1`)
			require.NoError(t, err)
			hint := QuotaRecoveryHint{Key: groupdelivery.GroupKey{WorkspaceID: "workspace-A", GroupID: "group-G", GroupRevision: "rev-1"}, WorkspaceRevision: "workspace-revision", ConnectorID: "connector-1", ConnectorRevision: "crev-1", ReferenceID: "original-ref", Scope: groupdelivery.QuotaScopeGroup, Observe: func(IntakeQuotaObservation) { t.Error("unproven recovery must not clear original token") }}
			switch mode {
			case "unconfigured-global":
				hint.Scope = groupdelivery.QuotaScopeGlobal
			case "foreign-connector":
				hint.ConnectorID = "foreign"
			case "missing-frozen":
				_, err = f.store.DB().Exec(`DELETE FROM wg_runtime_versions`)
				require.NoError(t, err)
			}
			reads := 0
			p.deps.IntakeQuotaObservation.ScopeRevision = func(ctx context.Context, key groupdelivery.GroupKey) (string, bool) {
				reads++
				if mode == "scope-panic" {
					panic("secret")
				}
				if mode == "changed-after-scope" && reads == 2 {
					_, err = f.store.DB().Exec(`DELETE FROM wg_runtime_versions`)
					require.NoError(t, err)
				}
				return "workspace-revision", true
			}
			p.deps.IntakeQuotaObservation.TakeQuotaRecoveryHint = func(context.Context) (QuotaRecoveryHint, bool) { return hint, true }
			require.NoError(t, p.Reconcile(t.Context()))
			require.Contains(t, *out, IntakeQuotaObservation{})
		})
	}
}

func TestQuotaObservationInvalidHintScopeIsRejectedBeforeAllocationAndReads(t *testing.T) {
	f := newFixture(t)
	p := f.pipeline()
	require.NoError(t, p.Reconcile(t.Context()))
	var out []IntakeQuotaObservation
	reads := 0
	intakeCapture(p, &out, &reads)
	p.deps.IntakeQuotaObservation.IntakeEnabled = false
	hint := QuotaRecoveryHint{Key: groupdelivery.GroupKey{WorkspaceID: "workspace-A", GroupID: "group-G", GroupRevision: "rev-1"}, WorkspaceRevision: "workspace-revision", ConnectorID: "connector-1", ConnectorRevision: "crev-1", ReferenceID: "original-ref", Scope: groupdelivery.QuotaScope(strings.Repeat("private", 2048)), Observe: func(IntakeQuotaObservation) { t.Error("invalid scope cannot recover") }}
	p.deps.IntakeQuotaObservation.TakeQuotaRecoveryHint = func(context.Context) (QuotaRecoveryHint, bool) { return hint, true }
	ctx, finish := p.ownerOperation(t.Context())
	p.observeQuotaRecovery(ctx)
	b := ctx.Value(ownerBatchContextKey{}).(*ownerObservationBatch)
	require.Len(t, b.descriptors, 1)
	retained := b.descriptors[0].intake.hint
	finish()
	require.Zero(t, reads, "invalid closed scope must be rejected before proof reads")
	require.Nil(t, retained, "shared bounded queue must not retain arbitrary scope bytes")
	require.Equal(t, []IntakeQuotaObservation{{}}, out)
}
