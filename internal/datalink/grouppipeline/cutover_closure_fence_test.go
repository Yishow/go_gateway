package grouppipeline

import (
	"database/sql"
	"encoding/json"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"
)

func TestCutoverClosureFenceRejectsCapturedPlanAndRefreshes(t *testing.T) {
	for _, shape := range []averageShape{shapeMeanLong, shapeMeanWide} {
		for _, empty := range []bool{false, true} {
			name := "preboundary-only"
			if empty {
				name = "empty"
			}
			t.Run(string(shape)+"/"+name, func(t *testing.T) {
				f, _ := reductionHandoffFixture(t, shape, "latest")
				if !empty {
					require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "old-only", at(9), uint64(20))))
				}
				old := f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)]
				entered, release := make(chan struct{}), make(chan struct{})
				var once atomic.Bool
				var released sync.Once
				t.Cleanup(func() { released.Do(func() { close(release) }); f.localFault.Reject.Store(nil) })
				reject := func(query string) error {
					if strings.Contains(query, "UPDATE wg_delivery_checkpoints SET updated_at = updated_at WHERE 0") && once.CompareAndSwap(false, true) {
						close(entered) // The stale plan exists; it has not acquired write authority.
						<-release
					}
					return nil
				}
				f.localFault.Reject.Store(&reject)
				closed := make(chan error, 1)
				go func() { closed <- old.boundary.Tick(t.Context(), at(15)) }()
				<-entered
				f.commit() // Durable cutover wins while no receipt/pending successor input exists.
				before := conflictPairState(t, f)
				released.Do(func() { close(release) })
				err := <-closed
				f.localFault.Reject.Store(nil)
				require.ErrorIs(t, err, groupdelivery.ErrOwnershipConflict)
				require.Equal(t, before, conflictPairState(t, f), "stale closure advances no journal, checkpoint, baseline or outbox")
				require.Equal(t, at(0), old.boundary.NextClose(), "failed durable closure cannot commit the in-memory plan")
				require.Zero(t, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_sample_handoffs`), "the epoch fence cannot depend on receipts")

				f.source.EndsAt, f.clock = f.boundary, at(15)
				f.groups.groups = append(f.groups.groups, f.target.Group)
				f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
				require.NoError(t, f.pipe.Reconcile(t.Context()))
				fresh := f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)]
				require.NotSame(t, old, fresh, "reconciliation refreshes owner evidence even without handoff participants")
				require.NoError(t, fresh.boundary.Tick(t.Context(), at(15)), "legal refreshed partial window must not remain permanently blocked")
				require.Equal(t, at(15), fresh.boundary.NextClose(), "the aligned close cursor may exceed actual ownership end 10")
				f.deliver()
				var value sql.NullString
				var raw string
				require.NoError(t, f.plainTarget().QueryRow(`SELECT kwh_1,provenance FROM `+f.source.Group.Destination.TableName+` WHERE group_id=?`, f.source.GroupID).Scan(&value, &raw))
				require.Equal(t, !empty, value.Valid)
				if !empty {
					require.Equal(t, "20", value.String)
				}
				var entries []reductionRowEvidence
				require.NoError(t, json.Unmarshal([]byte(raw), &entries))
				require.NotEmpty(t, entries)
				require.Equal(t, "end", entries[0].Window.Truncated)
				require.Equal(t, at(10).Format("2006-01-02T15:04:05.000000Z"), entries[0].Window.End)
				digests := f.outboxDigests()
				require.NoError(t, fresh.boundary.Tick(t.Context(), at(30)))
				require.Equal(t, digests, f.outboxDigests(), "refresh cannot reopen finalized output")
				require.NoError(t, f.pipe.Reconcile(t.Context()))
				require.Same(t, fresh, f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)], "an unchanged owner view does not rebuild the journal again")
			})
		}
	}
}

func TestCutoverClosureFenceNewRestoreCannotAuthorizeStaleCanonicalBounds(t *testing.T) {
	f, _ := reductionHandoffFixture(t, shapeMeanLong, "latest")
	canonical, err := f.groups.ResolveAppliedAt(t.Context(), f.source.GroupID, f.clock)
	require.NoError(t, err)
	stale := *canonical // Captured before the cutover, with the old open-ended bounds.
	entered, release := make(chan struct{}), make(chan struct{})
	var once atomic.Bool
	var released sync.Once
	t.Cleanup(func() { released.Do(func() { close(release) }); f.localFault.Reject.Store(nil) })
	reject := func(query string) error {
		if strings.Contains(query, "SELECT payload FROM wg_runtime_versions\n") && once.CompareAndSwap(false, true) {
			close(entered)
			<-release
		}
		return nil
	}
	f.localFault.Reject.Store(&reject)
	type result struct {
		worker *managed
		reason string
	}
	built := make(chan result, 1)
	go func() { worker, reason := f.pipe.build(t.Context(), &stale); built <- result{worker, reason} }()
	<-entered
	f.commit()
	before := conflictPairState(t, f)
	released.Do(func() { close(release) })
	got := <-built
	f.localFault.Reject.Store(nil)
	require.Nil(t, got.worker, "a newly read owner view cannot be attached to the stale canonical configuration")
	require.Equal(t, groupdelivery.ErrOwnershipConflict.Error(), got.reason)
	require.Equal(t, before, conflictPairState(t, f))
	f.source.EndsAt = f.boundary
	fresh, reason := f.pipe.build(t.Context(), f.source)
	require.Empty(t, reason)
	require.NotNil(t, fresh)
	require.NoError(t, fresh.boundary.Tick(t.Context(), at(15)), "a fresh canonical read restores legal partial closure")
}

func TestCutoverClosureFencePreservesEarlierDisabledStop(t *testing.T) {
	for _, shape := range []averageShape{shapeMeanLong, shapeMeanWide} {
		t.Run(string(shape), func(t *testing.T) {
			f := newOwnershipPipelineFixture(t, shape)
			f.clock = at(7)
			_, err := f.plainTarget().Exec(`CREATE TABLE legacy_readings(kwh_1 TEXT,kwh_2 TEXT)`)
			require.NoError(t, err)
			for _, member := range f.source.Group.Members {
				input := envelope(f.source.Group, "before-disable-"+member.TagID, at(9), uint64(20))
				input.DeviceID, input.PointID, input.TagID = member.DeviceID, member.PointID, member.TagID
				input.SourceRevision, input.MappingRevision = member.SourceRevision, member.MappingRevision
				require.NoError(t, f.pipe.AcceptSample(t.Context(), input))
			}
			old := f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)]
			f.commit() // The future 22.5 cutover is already durable before disable.
			f.source.EndsAt = f.boundary
			// Current lifecycle status is separate from the immutable applied descriptor.
			disabled := *f.source.Group
			disabled.Status = workspace.WriteGroupStatusDisabled
			f.groups.groups = []*workspace.WriteGroup{&disabled}
			err = f.pipe.Reconcile(t.Context())
			require.Equal(t, at(10), old.boundary.ReplayState().Until)
			require.NoError(t, err, "an earlier disabled drain is inside the new owner's bounds")
			fresh := f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)]
			require.NotSame(t, old, fresh)
			require.Equal(t, at(10), fresh.boundary.ReplayState().Until, "refresh must not extend the disabled stop to 30")
			require.NoError(t, fresh.boundary.Tick(t.Context(), at(30)))
			require.Equal(t, at(10), fresh.boundary.NextClose())
			f.deliver()
			var value string
			require.NoError(t, f.plainTarget().QueryRow(`SELECT kwh_1 FROM legacy_readings`).Scan(&value))
			require.Equal(t, "20", value)
			digests := f.outboxDigests()
			require.NoError(t, fresh.boundary.Tick(t.Context(), at(40)))
			require.Equal(t, digests, f.outboxDigests())
		})
	}
}
