package grouppipeline

import (
	"database/sql"
	"testing"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func reopenStore(t *testing.T, f *fixture) {
	t.Helper()
	var seq int
	var name, path string
	require.NoError(t, f.store.DB().QueryRowContext(t.Context(), "PRAGMA database_list").Scan(&seq, &name, &path))
	require.NoError(t, f.store.DB().Close())
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	f.store = groupdelivery.NewStore(db)
}

func acceptComplete(t *testing.T, f *fixture, p *Pipeline, group *workspace.WriteGroup, id string, second int) {
	t.Helper()
	f.clock = at(second)
	require.NoError(t, p.AcceptSample(t.Context(), envelope(group, id+"-t", at(second), 21.5)))
	pressure := envelope(group, id+"-p", at(second), int64(42))
	pressure.PointID, pressure.TagID = group.Members[1].PointID, group.Members[1].TagID
	pressure.SourceRevision = group.Members[1].SourceRevision
	require.NoError(t, p.AcceptSample(t.Context(), pressure))
}

func TestDraftSaveKeepsAppliedIntakeAcrossBucketsAndRestart(t *testing.T) {
	f := newFixture(t)
	p := f.pipeline()
	require.NoError(t, p.Reconcile(t.Context()))
	applied := f.groups.snapshots["group-G"][0].Group
	draft := liveGroup(applied, workspace.WriteGroupStatusDraft)
	draft.Revision, draft.Name = "draft-2", "Renamed"
	draft.Members = nil
	f.groups.groups = []*workspace.WriteGroup{draft}
	require.NoError(t, p.Reconcile(t.Context()))
	f.clock = at(21)
	p.TickAll(t.Context())
	acceptComplete(t, f, p, applied, "still-applied", 22)
	reopenStore(t, f)
	restarted := f.pipeline()
	require.NoError(t, restarted.Reconcile(t.Context()))
	require.True(t, restarted.WantsSample("device-1", "point-1", "tag-t"))
	acceptComplete(t, f, restarted, applied, "after-restart", 32)
	var samples int
	require.NoError(t, rowCount(f.store, &samples, "wg_delivery_samples"))
	require.Equal(t, 4, samples)
}

func TestRestartUsesEachAppliedSnapshotsInterval(t *testing.T) {
	f := newFixture(t)
	old, next := groupAt("rev-1"), groupAt("rev-2")
	old.RowPolicy.IntervalSeconds, next.RowPolicy.IntervalSeconds = 7, 10
	f.groups.snapshots["group-G"] = []*workspace.WriteGroupAppliedSnapshot{snapshotOf(old, at(7)), snapshotOf(next, at(70))}
	f.groups.groups = []*workspace.WriteGroup{next}
	f.clock = at(62)
	p := f.pipeline()
	require.NoError(t, p.Reconcile(t.Context()))
	require.Len(t, p.Status(), 2)
	require.Contains(t, p.Status(), GroupStatus{GroupID: old.ID, Revision: old.Revision, State: StateRetiring})
	acceptComplete(t, f, p, old, "old-period", 63)
	var samples int
	require.NoError(t, rowCountWhere(f.store, &samples, "wg_delivery_samples", "group_revision = 'rev-1'"))
	require.Equal(t, 2, samples)
}

func TestRestartDrainsAcceptedHistoricalJournalWithoutIntake(t *testing.T) {
	for _, state := range []workspace.WriteGroupStatus{workspace.WriteGroupStatusDisabled, workspace.WriteGroupStatusDeleted, workspace.WriteGroupStatusReady} {
		t.Run(string(state), func(t *testing.T) {
			f := newFixture(t)
			p := f.pipeline()
			require.NoError(t, p.Reconcile(t.Context()))
			old := f.groups.snapshots["group-G"][0].Group
			acceptComplete(t, f, p, old, "accepted", 2)
			if state == workspace.WriteGroupStatusReady {
				next := groupAt("rev-2")
				f.groups.snapshots["group-G"] = append(f.groups.snapshots["group-G"], snapshotOf(next, at(10)))
				f.groups.groups = []*workspace.WriteGroup{next}
			} else {
				f.groups.groups = []*workspace.WriteGroup{liveGroup(old, state)}
			}
			f.clock = at(12)
			f.inspect.inspection = nil // recovery must not depend on the remote DB
			reopenStore(t, f)
			restarted := f.pipeline()
			require.NoError(t, restarted.Reconcile(t.Context()))
			require.NoError(t, restarted.AcceptSample(t.Context(), envelope(old, "must-ignore", at(3), 25.0)))
			restarted.TickAll(t.Context())
			var open, outbox int
			require.NoError(t, rowCountWhere(f.store, &open, "wg_delivery_samples", "consumed = 0 AND group_revision = 'rev-1'"))
			require.Zero(t, open)
			require.NoError(t, rowCountWhere(f.store, &outbox, "wg_delivery_outbox", "group_revision = 'rev-1' AND connector_revision = 'crev-1'"))
			require.Equal(t, 1, outbox)
			restarted.TickAll(t.Context())
			require.NoError(t, rowCountWhere(f.store, &outbox, "wg_delivery_outbox", "group_revision = 'rev-1'"))
			require.Equal(t, 1, outbox)
		})
	}
}

func TestVerifiedAppliedLayoutRestoresIntakeWhileRemoteIsOffline(t *testing.T) {
	f := newFixture(t)
	require.NoError(t, f.pipeline().Reconcile(t.Context()))
	f.inspect.inspection = nil
	f.tags.missing["tag-p"] = true // types belong to the applied version
	reopenStore(t, f)
	restarted := f.pipeline()
	require.NoError(t, restarted.Reconcile(t.Context()))
	require.True(t, restarted.WantsSample("device-1", "point-1", "tag-t"))
	acceptComplete(t, f, restarted, f.groups.snapshots["group-G"][0].Group, "offline", 2)
	f.clock = at(12)
	restarted.TickAll(t.Context())
	var rows int
	require.NoError(t, rowCount(f.store, &rows, "wg_delivery_outbox"))
	require.Equal(t, 1, rows)
	f.dest.connector.IdentityRevision = "crev-2"
	changed := f.pipeline()
	require.NoError(t, changed.Reconcile(t.Context()))
	require.False(t, changed.WantsSample("device-1", "point-1", "tag-t"))
	require.Equal(t, reasonDestinationRevision, statusOf(changed).Reason)
}

func TestMissingLegacyRecoveryDescriptorRetainsAcceptedJournalWhenOffline(t *testing.T) {
	f := newFixture(t)
	p := f.pipeline()
	require.NoError(t, p.Reconcile(t.Context()))
	old := f.groups.snapshots["group-G"][0].Group
	acceptComplete(t, f, p, old, "legacy", 2)
	_, err := f.store.DB().ExecContext(t.Context(), `DELETE FROM wg_runtime_versions`)
	require.NoError(t, err)
	f.groups.groups = []*workspace.WriteGroup{liveGroup(old, workspace.WriteGroupStatusDeleted)}
	f.inspect.inspection = nil
	f.clock = at(12)
	reopenStore(t, f)
	restarted := f.pipeline()
	require.NoError(t, restarted.Reconcile(t.Context()))
	restarted.TickAll(t.Context())
	require.False(t, restarted.WantsSample("device-1", "point-1", "tag-t"))
	require.Equal(t, StateBlocked, statusOf(restarted).State)
	var open int
	require.NoError(t, rowCountWhere(f.store, &open, "wg_delivery_samples", "consumed = 0"))
	require.Equal(t, 2, open, "unverifiable legacy data stays durable until metadata can be verified")
}
