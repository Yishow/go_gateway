package grouppipeline

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/schema"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"
)

func reconciledHandoffCell(t *testing.T, f *ownershipPipelineFixture, groupID string, bucket time.Time, column string) []any {
	t.Helper()
	rows, err := f.local.Query(`SELECT bucket_start,payload FROM wg_delivery_outbox WHERE group_id=? ORDER BY effect_key`, groupID)
	require.NoError(t, err)
	defer rows.Close()
	var values []any
	for rows.Next() {
		var when, raw string
		require.NoError(t, rows.Scan(&when, &raw))
		at, err := time.Parse(time.RFC3339Nano, when)
		require.NoError(t, err)
		if !at.Equal(bucket) {
			continue
		}
		row, err := groupdelivery.DecodeRowPayload([]byte(raw))
		require.NoError(t, err)
		for _, cell := range row.Cells {
			if cell.Column == column && cell.Value != nil {
				values = append(values, cell.Value)
			}
		}
	}
	require.NoError(t, rows.Err())
	return values
}

func TestCutoverReconcileRebuildsBothEligibleSets(t *testing.T) {
	for _, shape := range []averageShape{shapeMeanLong, shapeMeanWide} {
		t.Run(string(shape), func(t *testing.T) {
			f := newHandoffPipelineFixture(t, shape)
			_, err := f.plainTarget().Exec(`CREATE TABLE legacy_readings(kwh_1 TEXT,kwh_2 TEXT)`)
			require.NoError(t, err)
			healthy := f.readyHealthy()
			for _, input := range []struct {
				id       string
				observed time.Time
				value    uint64
			}{
				{"after", at(23), 7}, {"exact", f.boundary, 7}, {"before", f.boundary.Add(-time.Millisecond), 6},
			} {
				require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, input.id, input.observed, input.value)))
			}
			second := *f.source.Group
			second.Members = append([]workspace.WriteGroupMember(nil), second.Members[1:]...)
			require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(&second, "required-second", f.boundary.Add(-time.Millisecond), uint64(600))))
			original := handoffOriginalRows(t, f)
			f.commit()
			f.source.EndsAt = f.boundary
			f.groups.groups = append(f.groups.groups, f.target.Group)
			f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
			empty, reason := f.pipe.build(t.Context(), f.target)
			require.Empty(t, reason)
			f.pipe.installManaged(empty)
			old := f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)]
			require.Error(t, old.boundary.Tick(t.Context(), at(40)), "stale source cannot consume the transferred observation")
			require.Error(t, empty.boundary.Tick(t.Context(), at(40)), "empty target cannot close ahead of accepted input")
			require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(healthy.Group, "healthy", at(23), uint64(9))))
			f.clock = at(40)
			require.NoError(t, f.pipe.Reconcile(t.Context()))
			source, err := f.store.Restore(t.Context(), groupKey(f.source))
			require.NoError(t, err)
			require.Len(t, source.Samples, 2, "live reconcile transfers exact/after boundary while preserving both eligible source members")
			target, err := f.store.Restore(t.Context(), groupKey(f.target))
			require.NoError(t, err)
			require.Len(t, target.Samples, 2)
			require.Equal(t, original, handoffOriginalRows(t, f))
			require.NotSame(t, old, f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)])
			require.NotSame(t, empty, f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)])
			require.NoError(t, f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)].boundary.Tick(t.Context(), f.clock))
			require.NoError(t, f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)].boundary.Tick(t.Context(), f.clock))
			f.pipe.TickAll(t.Context())
			require.Equal(t, []any{"6"}, reconciledHandoffCell(t, f, f.source.GroupID, at(20), "kwh_1"))
			column := "value_decimal"
			if shape.wide() {
				column = "kwh_1"
			}
			require.Equal(t, []any{"7"}, reconciledHandoffCell(t, f, f.target.GroupID, at(15), column))
			var sourceConsumed, targetConsumed int
			require.NoError(t, f.local.QueryRow(`SELECT SUM(consumed) FROM wg_delivery_samples WHERE group_id=? AND json_extract(payload,'$.sample_id') IN ('after','exact')`, f.source.GroupID).Scan(&sourceConsumed))
			require.NoError(t, f.local.QueryRow(`SELECT SUM(consumed) FROM wg_delivery_samples WHERE group_id=?`, f.target.GroupID).Scan(&targetConsumed))
			require.Zero(t, sourceConsumed)
			require.Equal(t, 2, targetConsumed)
			require.Equal(t, []any{"9"}, reconciledHandoffCell(t, f, healthy.GroupID, at(20), "kwh_1"))
			f.deliver()
			var oldValue, newValue string
			require.NoError(t, f.plainTarget().QueryRow(`SELECT kwh_1 FROM legacy_readings WHERE kwh_2='600'`).Scan(&oldValue))
			query := `SELECT CAST(value_decimal AS TEXT) FROM ` + f.table + ` WHERE tag_id='tag-kwh-1' AND value_decimal IS NOT NULL`
			if shape.wide() {
				query = `SELECT kwh_1 FROM ` + f.table + ` WHERE kwh_1 IS NOT NULL`
			}
			require.NoError(t, f.plainTarget().QueryRow(query).Scan(&newValue))
			require.Equal(t, "6", oldValue)
			require.Equal(t, "7", newValue)
		})
	}
}

// Blocking destination metadata holds the rebuild after receipt commit and
// after the old participant gates are acquired, without sleeps or extra writers.
type cutoverRebuildBarrier struct {
	Destinations
	entered, resume chan struct{}
	once            atomic.Bool
}

func (b *cutoverRebuildBarrier) GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error) {
	if b.once.CompareAndSwap(false, true) {
		close(b.entered)
		select {
		case <-b.resume:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return b.Destinations.GetByID(ctx, id)
}

func prepareLiveHandoffFixture(t *testing.T) (*ownershipPipelineFixture, *managed, *managed) {
	t.Helper()
	f := newHandoffPipelineFixture(t, shapeMeanLong)
	require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "buffered", at(23), uint64(7))))
	f.commit()
	f.source.EndsAt = f.boundary
	f.groups.groups = append(f.groups.groups, f.target.Group)
	f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
	target, reason := f.pipe.build(t.Context(), f.target)
	require.Empty(t, reason)
	f.pipe.installManaged(target)
	return f, f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)], target
}

func TestCutoverReconcileCapturedWorkersCannotRunAfterReplacement(t *testing.T) {
	f := newHandoffPipelineFixture(t, shapeMeanLong)
	healthy := f.readyHealthy()
	require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "buffered", at(23), uint64(7))))
	f.commit()
	f.source.EndsAt = f.boundary
	f.groups.groups = append(f.groups.groups, f.target.Group)
	f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
	target, reason := f.pipe.build(t.Context(), f.target)
	require.Empty(t, reason)
	f.pipe.installManaged(target)
	source := f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)]
	f.clock = at(22)
	barrier := &cutoverRebuildBarrier{Destinations: f.pipe.deps.Destinations, entered: make(chan struct{}), resume: make(chan struct{})}
	f.pipe.deps.Destinations = barrier
	var release sync.Once
	defer release.Do(func() { close(barrier.resume) })
	rebuilt := make(chan error, 1)
	go func() { rebuilt <- f.pipe.Reconcile(t.Context()) }()
	<-barrier.entered
	// These are the same captured worker pointers used by AcceptSample/TickAll.
	accepted, ticked := make(chan error, 1), make(chan error, 1)
	go func() {
		accepted <- f.pipe.acceptWithOwnership(t.Context(), envelope(f.target.Group, "new-native", at(24), uint64(9)), []*managed{target})
	}()
	go func() {
		ticked <- f.pipe.withCounterProofGate(source, func() error { return source.boundary.Tick(t.Context(), at(40)) })
	}()
	independent := make(chan error, 1)
	go func() {
		independent <- f.pipe.AcceptSample(t.Context(), envelope(healthy.Group, "independent", at(23), uint64(9)))
	}()
	deadline, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	select {
	case err := <-independent:
		require.NoError(t, err)
	case <-deadline.Done():
		t.Fatal("unrelated worker waited for cutover rebuild")
	}
	release.Do(func() { close(barrier.resume) })
	require.NoError(t, <-rebuilt)
	require.ErrorIs(t, <-accepted, groupdelivery.ErrOwnershipUnavailable)
	require.ErrorIs(t, <-ticked, groupdelivery.ErrOwnershipUnavailable)
	var checkpoints int
	require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_checkpoints`).Scan(&checkpoints))
	require.Zero(t, checkpoints, "captured stale Tick never reaches closure")
	require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.target.Group, "new-native", at(24), uint64(9))))
	restored, err := f.store.Restore(t.Context(), groupKey(f.target))
	require.NoError(t, err)
	require.Len(t, restored.Samples, 2)
	f.clock = at(40)
	require.NoError(t, f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)].boundary.Tick(t.Context(), f.clock))
	require.Equal(t, []any{"8"}, reconciledHandoffCell(t, f, f.target.GroupID, at(15), "value_decimal"))
}

func TestCutoverReconcilePreservesDisabledDrainStop(t *testing.T) {
	f, _, target := prepareLiveHandoffFixture(t)
	target.boundary.SetUntil(at(30))
	target.drainOnly = true
	disabled := *f.target.Group
	disabled.Status = workspace.WriteGroupStatusDisabled
	for i, g := range f.groups.groups {
		if g.ID == disabled.ID {
			f.groups.groups[i] = &disabled
		}
	}
	f.clock = at(40)
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	fresh := f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)]
	require.NotSame(t, target, fresh)
	require.True(t, fresh.drainOnly)
	require.Equal(t, at(30), fresh.boundary.ReplayState().Until, "rebuild cannot replace original stop with now's next boundary")
	require.NoError(t, fresh.boundary.Tick(t.Context(), f.clock))
	require.Equal(t, at(30), fresh.boundary.NextClose())
	require.True(t, fresh.boundary.Retired())
}

type cutoverSecondBuildFailure struct {
	Destinations
	calls atomic.Int32
}

func (f *cutoverSecondBuildFailure) GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error) {
	if f.calls.Add(1) == 2 {
		return nil, errors.New("private destination metadata failure")
	}
	return f.Destinations.GetByID(ctx, id)
}

func TestCutoverReconcileBuildFailureDoesNotPartiallyPublish(t *testing.T) {
	f, source, target := prepareLiveHandoffFixture(t)
	target.boundary.SetUntil(at(30))
	failing := &cutoverSecondBuildFailure{Destinations: f.pipe.deps.Destinations}
	f.pipe.deps.Destinations = failing
	err := f.pipe.Reconcile(t.Context())
	require.ErrorIs(t, err, groupdelivery.ErrOwnershipUnavailable)
	require.NotContains(t, err.Error(), "private destination")
	require.Equal(t, int32(2), failing.calls.Load(), "first replacement was prepared before the second failed")
	require.Same(t, source, f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)])
	require.Same(t, target, f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)])
	require.Equal(t, at(30), target.boundary.ReplayState().Until)
	require.Error(t, target.boundary.Tick(t.Context(), at(40)), "unpublished candidate cannot release the stale view guard")
	var checkpoints int
	require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_checkpoints`).Scan(&checkpoints))
	require.Zero(t, checkpoints)
	f.pipe.deps.Destinations = failing.Destinations
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	fresh := f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)]
	require.NotSame(t, target, fresh)
	require.Equal(t, at(30), fresh.boundary.ReplayState().Until)
	require.NoError(t, fresh.boundary.Tick(t.Context(), at(40)))
	require.Equal(t, []any{"7"}, reconciledHandoffCell(t, f, f.target.GroupID, at(15), "value_decimal"))
}

func TestCutoverReconcileDoesNotTakeHandoffWriteLocksForUnrelatedInput(t *testing.T) {
	f := newHandoffPipelineFixture(t, shapeMeanLong)
	require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "ordinary", at(23), uint64(7))))
	writes := 0
	observe := func(query string) error {
		if strings.Contains(query, "UPDATE wg_delivery_checkpoints SET updated_at = updated_at WHERE 0") {
			writes++
		}
		return nil
	}
	f.localFault.Reject.Store(&observe)
	err := f.pipe.Reconcile(t.Context())
	f.localFault.Reject.Store(nil)
	require.NoError(t, err)
	require.Zero(t, writes, "ordinary unrelated accepted input is not cutover migration work")
}

func TestCutoverReconcileSupersededGateRefusesBeforeInvoking(t *testing.T) {
	f, source, target := prepareLiveHandoffFixture(t)
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	for _, captured := range []*managed{source, target} {
		invoked := false
		err := f.pipe.withCounterProofGate(captured, func() error {
			invoked = true
			return nil
		})
		require.ErrorIs(t, err, groupdelivery.ErrOwnershipUnavailable)
		require.False(t, invoked, "captured Tick or Accept cannot reach any downstream work")
	}
}
