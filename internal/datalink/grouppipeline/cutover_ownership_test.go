package grouppipeline

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/workspace"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCutoverOwnershipDelayedReconcileRefusesUnownedAck(t *testing.T) {
	for _, shape := range []averageShape{shapeMeanLong, shapeMeanWide} {
		t.Run(string(shape), func(t *testing.T) {
			f := newOwnershipPipelineFixture(t, shape)
			f.clock = at(21)
			f.commit()
			f.clock = at(24)
			observations := []struct {
				id    string
				at    time.Time
				value uint64
			}{
				{"after", at(23), 7}, {"exact", f.boundary, 7}, {"before", at(22).Add(499 * time.Millisecond), 6},
			}
			for i, input := range observations {
				err := f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, input.id, input.at, input.value))
				if i < 2 {
					require.ErrorContains(t, err, "recording-ownership-", "not-ready target must not ACK through old worker")
				} else {
					require.NoError(t, err)
				}
			}
			old, err := f.store.Restore(t.Context(), groupKey(f.source))
			require.NoError(t, err)
			require.Len(t, old.Samples, 1)
			require.Equal(t, "before", old.Samples[0].SampleID)
			require.Equal(t, measurement.NewUint64(6), old.Samples[0].Value)
			f.source.EndsAt = f.boundary // Canonical lookup now observes the persisted revision end.
			f.readyTarget()
			for _, input := range observations {
				require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, input.id, input.at, input.value)))
			}
			old, err = f.store.Restore(t.Context(), groupKey(f.source))
			require.NoError(t, err)
			require.Len(t, old.Samples, 1)
			next, err := f.store.Restore(t.Context(), groupKey(f.target))
			require.NoError(t, err)
			require.Len(t, next.Samples, 2)
			for _, s := range next.Samples {
				require.True(t, !s.ObservedAt.Before(f.boundary))
				require.Equal(t, measurement.NewUint64(7), s.Value)
			}
		})
	}
}

func TestCutoverOwnershipNoWorkersCannotAcknowledge(t *testing.T) {
	f := newOwnershipPipelineFixture(t, shapeMeanWide)
	f.commit()
	f.clock = at(24)
	// All workers can be absent during restart or a failed activation. The
	// durable cutover still owns its source; empty cache is not an ACK.
	f.pipe.boundaries = map[string]*managed{}
	err := f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "unaccepted", at(23), uint64(7)))
	require.ErrorContains(t, err, "recording-ownership-unavailable")
	var count int
	require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples`).Scan(&count))
	require.Zero(t, count)
}

func TestCutoverOwnershipCachedCutoffAndQuarantineCannotAck(t *testing.T) {
	for _, mode := range []string{"cached cutoff", "quarantined worker"} {
		t.Run(mode, func(t *testing.T) {
			f := newOwnershipPipelineFixture(t, shapeMeanWide)
			f.commit()
			f.clock = at(24)
			f.source.EndsAt = f.boundary // Refresh may not pair the new owner view with stale bounds.
			f.readyTarget()
			m := f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)]
			if mode == "cached cutoff" {
				m.boundary.SetUntil(f.boundary)
			} else {
				m.proofQuarantined = true
			}
			err := f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "ignored", at(23), uint64(7)))
			var owner *groupdelivery.OwnershipError
			require.ErrorAs(t, err, &owner)
			var count int
			require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples`).Scan(&count))
			require.Zero(t, count)
		})
	}
}

func TestCutoverOwnershipAlreadyAcceptedFutureIsNotOrdinaryTargetAppend(t *testing.T) {
	f := newOwnershipPipelineFixture(t, shapeMeanWide)
	f.clock = at(21)
	env := envelope(f.source.Group, "pending-transfer", at(23), uint64(7))
	require.NoError(t, f.pipe.AcceptSample(t.Context(), env))
	var original string
	require.NoError(t, f.local.QueryRow(`SELECT payload FROM wg_delivery_samples WHERE sample_id='pending-transfer'`).Scan(&original))
	f.commit()
	f.clock = at(24)
	// Keep this ordinary-append guard independent from the now-live Reconcile
	// handoff adapter: create a target worker without transferring the source.
	f.groups.groups = append(f.groups.groups, f.target.Group)
	f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
	empty, reason := f.pipe.build(t.Context(), f.target)
	require.Empty(t, reason)
	f.pipe.installManaged(empty)
	err := f.pipe.AcceptSample(t.Context(), env)
	require.ErrorIs(t, err, groupdelivery.ErrOwnershipConflict)
	var after string
	require.NoError(t, f.local.QueryRow(`SELECT payload FROM wg_delivery_samples WHERE sample_id='pending-transfer'`).Scan(&after))
	require.Equal(t, original, after)
	var count int
	require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples`).Scan(&count))
	require.Equal(t, 1, count)
}

func TestCutoverOwnershipCommitBetweenRoutingAndAppend(t *testing.T) {
	f := newOwnershipPipelineFixture(t, shapeMeanLong)
	f.clock = at(21)
	reached, resume := make(chan struct{}), make(chan struct{})
	var once sync.Once
	defer once.Do(func() { close(resume) })
	var hit atomic.Bool
	hook := func(query string) error {
		if strings.Contains(query, "UPDATE wg_delivery_checkpoints SET updated_at = updated_at WHERE 0") && hit.CompareAndSwap(false, true) {
			close(reached)
			<-resume
		}
		return nil
	}
	f.localFault.Reject.Store(&hook)
	result := make(chan error, 1)
	go func() {
		result <- f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "stale-plan", at(23), uint64(7)))
	}()
	select {
	case <-reached:
	case err := <-result:
		t.Fatalf("append did not reach transaction barrier: %v", err)
	}
	f.commit()
	once.Do(func() { close(resume) })
	err := <-result
	f.localFault.Reject.Store(nil)
	require.ErrorIs(t, err, groupdelivery.ErrOwnershipConflict)
	for _, table := range []string{"wg_delivery_samples", "wg_delivery_checkpoints"} {
		var count int
		require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM `+table).Scan(&count))
		require.Zero(t, count)
	}
}

func TestCutoverOwnershipUnavailableDoesNotStopHealthyScope(t *testing.T) {
	f := newOwnershipPipelineFixture(t, shapeMeanWide)
	healthy := f.readyHealthy()
	f.commit()
	f.clock = at(24)
	_, err := f.local.Exec(`DELETE FROM write_group_ownership_fences WHERE direction='started'`)
	require.NoError(t, err)
	err = f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "blocked", at(23), uint64(7)))
	require.ErrorIs(t, err, groupdelivery.ErrOwnershipUnavailable)
	for range 2 {
		require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(healthy.Group, "healthy", at(23), uint64(9))))
	}
	var count int
	require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id='healthy-group'`).Scan(&count))
	require.Equal(t, 1, count)
	require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id<>'healthy-group'`).Scan(&count))
	require.Zero(t, count)
}

func TestCutoverOwnershipCommitBeforeCachedIgnoreCannotAck(t *testing.T) {
	f := newOwnershipPipelineFixture(t, shapeMeanWide)
	f.clock = at(21)
	reached, resume := make(chan struct{}), make(chan struct{})
	var once sync.Once
	defer once.Do(func() { close(resume) })
	var hit atomic.Bool
	// Freeze after the advisory read transaction really ends, before its caller
	// sees the old no-cutover result. No SQLite lock or sleep stands in for this.
	hook := func() {
		if hit.CompareAndSwap(false, true) {
			close(reached)
			<-resume
		}
	}
	f.localFault.AfterCommit.Store(&hook)
	result := make(chan error, 1)
	go func() {
		result <- f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "ignored-stale-plan", at(23), uint64(7)))
	}()
	select {
	case <-reached:
	case err := <-result:
		t.Fatalf("routing did not reach commit barrier: %v", err)
	}
	f.commit()
	// Reconciliation can update the old cache before a new worker is ready.
	f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)].boundary.SetUntil(f.boundary)
	once.Do(func() { close(resume) })
	err := <-result
	f.localFault.AfterCommit.Store(nil)
	require.ErrorIs(t, err, groupdelivery.ErrOwnershipConflict)
	var count int
	require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples`).Scan(&count))
	require.Zero(t, count)
}

func TestCutoverOwnershipOutsideAllDurableBoundsRefused(t *testing.T) {
	f := newOwnershipPipelineFixture(t, shapeMeanWide)
	f.commit()
	f.clock = at(24)
	err := f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "before-any-owner", at(-1), uint64(7)))
	require.ErrorIs(t, err, groupdelivery.ErrOwnershipConflict)
	var count int
	require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples`).Scan(&count))
	require.Zero(t, count)
}

func TestCutoverOwnershipInterestSurvivesMissingWorkers(t *testing.T) {
	f := newOwnershipPipelineFixture(t, shapeMeanWide)
	f.commit()
	f.clock = at(24)
	f.pipe.boundaries = map[string]*managed{}
	m := f.source.Group.Members[0]
	require.True(t, f.pipe.WantsSample(m.DeviceID, m.PointID, m.TagID), "sinkWants must not silently drop the source before explicit acceptance refusal")
	require.ErrorIs(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "no-worker-interest", at(23), uint64(7))), groupdelivery.ErrOwnershipUnavailable)
	require.False(t, f.pipe.WantsSample("unconfigured-device", "unconfigured-point", "unconfigured-tag"))
	_, err := f.local.Exec(`UPDATE write_groups SET status='disabled' WHERE id=?`, f.target.GroupID)
	require.NoError(t, err)
	for _, arrival := range []time.Time{at(24), at(20)} {
		f.clock = arrival
		require.True(t, f.pipe.WantsSample(m.DeviceID, m.PointID, m.TagID), "disabled target does not erase pre-boundary source interest, before or after boundary wall time")
		require.ErrorIs(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "before-end-no-worker", f.boundary.Add(-time.Nanosecond), uint64(6))), groupdelivery.ErrOwnershipUnavailable)
	}
	_, err = f.local.Exec(`UPDATE write_groups SET status='disabled' WHERE id=?`, f.source.GroupID)
	require.NoError(t, err)
	require.False(t, f.pipe.WantsSample(m.DeviceID, m.PointID, m.TagID), "historical fence cannot reactivate known disabled owners")
}

type ownershipInterestGroups struct {
	*fakeGroups
	workspace string
	inspect   func(context.Context)
}

func (g *ownershipInterestGroups) List(ctx context.Context) (*workspace.WriteGroupListResult, error) {
	if g.inspect != nil {
		g.inspect(ctx)
	}
	return &workspace.WriteGroupListResult{WorkspaceID: g.workspace}, nil
}

func TestCutoverOwnershipInterestIsScopedBoundedAndConservative(t *testing.T) {
	f := newOwnershipPipelineFixture(t, shapeMeanWide)
	f.commit()
	f.clock = at(24)
	member := f.source.Group.Members[0]
	f.pipe.deps.Groups = &ownershipInterestGroups{fakeGroups: f.groups, inspect: func(context.Context) { t.Fatal("cache hit must not perform canonical I/O") }}
	require.True(t, f.pipe.WantsSample(member.DeviceID, member.PointID, member.TagID))
	f.pipe.boundaries = map[string]*managed{}
	inspect := func(ctx context.Context) {
		deadline, ok := ctx.Deadline()
		require.True(t, ok)
		require.LessOrEqual(t, time.Until(deadline), time.Second)
		require.True(t, f.pipe.mu.TryLock(), "canonical I/O must not hold pipeline mutex")
		f.pipe.mu.Unlock()
	}
	f.pipe.deps.Groups = &ownershipInterestGroups{fakeGroups: f.groups, workspace: "foreign-workspace", inspect: inspect}
	require.False(t, f.pipe.WantsSample(member.DeviceID, member.PointID, member.TagID), "bare source IDs cannot activate a foreign workspace")
	f.pipe.deps.Groups = &ownershipInterestGroups{fakeGroups: f.groups, workspace: f.source.WorkspaceID, inspect: inspect}
	require.True(t, f.pipe.WantsSample(member.DeviceID, member.PointID, member.TagID))
	_, err := f.local.Exec(`DROP TABLE write_group_ownership_fences`)
	require.NoError(t, err)
	require.True(t, f.pipe.WantsSample(member.DeviceID, member.PointID, member.TagID), "read failure opts into explicit refusal, not silent drop")
	require.ErrorIs(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "authority-unavailable", at(23), uint64(7))), groupdelivery.ErrOwnershipUnavailable)
	empty := New(Dependencies{}, Config{})
	require.False(t, empty.WantsSample(member.DeviceID, member.PointID, member.TagID), "an unconfigured integration does not claim ordinary input")
}

func TestCutoverOwnershipInterestCancelsCanonicalRead(t *testing.T) {
	f := newOwnershipPipelineFixture(t, shapeMeanLong)
	f.pipe.boundaries = map[string]*managed{}
	f.pipe.deps.Groups = &ownershipInterestGroups{fakeGroups: f.groups, workspace: f.source.WorkspaceID, inspect: func(ctx context.Context) {
		<-ctx.Done()
		require.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
	}}
	member := f.source.Group.Members[0]
	require.True(t, f.pipe.WantsSample(member.DeviceID, member.PointID, member.TagID), "canceled authority read must reach explicit acceptance refusal")
}
