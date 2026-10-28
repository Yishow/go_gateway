package grouppipeline

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// Establish these original descriptors before any sample is accepted. Both
// owners are supported latest profiles; this is not a replay-time proof edit.
func latestHandoffFixture(t *testing.T, shape averageShape) *ownershipPipelineFixture {
	t.Helper()
	f := newHandoffPipelineFixture(t, shape)
	f.boundary, f.clock = at(12).Add(500*time.Millisecond), at(11)
	f.target.EffectiveAt = f.boundary
	f.target.Group.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReductionLatest
	for i := range f.target.Group.Members {
		f.target.Group.Members[i].ResultValueType = measurement.ExactUint64
	}
	raw, err := json.Marshal(f.target)
	require.NoError(t, err)
	var source workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(raw, &source))
	source.GroupID, source.AppliedRevision, source.EffectiveAt = f.source.GroupID, f.source.AppliedRevision, at(0)
	source.Group.ID, source.Group.Revision, source.Group.AppliedRevision = source.GroupID, source.AppliedRevision, source.AppliedRevision
	if shape.wide() {
		source.Group.Destination.TableName = "source_latest_wide"
		statements, err := recordingplan.GenerateGroupSchemaDDL("sqlite", wideAverageLayout(source.Group))
		require.NoError(t, err)
		// The shared receipt table already exists from the target fixture.
		_, err = f.plainTarget().Exec(statements[0])
		require.NoError(t, err)
		for i := range source.RuntimeLayout.Columns {
			if source.RuntimeLayout.Columns[i].Name == canonicalManagedOwnerColumn(f.target.Group) {
				source.RuntimeLayout.Columns[i].Name = canonicalManagedOwnerColumn(source.Group)
			}
		}
	}
	f.source = &source
	for _, s := range []*workspace.WriteGroupAppliedSnapshot{f.source, f.target} {
		raw, err = json.Marshal(s)
		require.NoError(t, err)
		_, err = f.local.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE group_id=?`, string(raw), s.GroupID)
		require.NoError(t, err)
	}
	f.groups.groups = []*workspace.WriteGroup{f.source.Group}
	f.groups.snapshots = map[string][]*workspace.WriteGroupAppliedSnapshot{f.source.GroupID: {f.source}}
	f.pipe = New(f.pipe.deps, f.pipe.config)
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	require.NotNil(t, f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)], "%+v", f.pipe.Status())
	return f
}

func TestCutoverReconcileLatestPreservesEarlier100AndSuccessor200(t *testing.T) {
	for _, shape := range []averageShape{shapeMeanLong, shapeMeanWide} {
		t.Run(string(shape), func(t *testing.T) {
			f := latestHandoffFixture(t, shape)
			for _, input := range []struct {
				id    string
				at    time.Time
				value uint64
			}{{"future", at(13), 200}, {"eligible", f.boundary.Add(-time.Millisecond), 100}, {"earlier", at(12), 90}} {
				require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, input.id, input.at, input.value)))
			}
			original := handoffOriginalRows(t, f)
			f.commit()
			f.source.EndsAt = f.boundary
			f.groups.groups = append(f.groups.groups, f.target.Group)
			f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
			f.clock = at(40)
			require.NoError(t, f.pipe.Reconcile(t.Context()))
			require.Equal(t, original, handoffOriginalRows(t, f))
			for _, s := range []*workspace.WriteGroupAppliedSnapshot{f.source, f.target} {
				require.NoError(t, f.pipe.boundaries[boundaryKey(s.GroupID, s.AppliedRevision)].boundary.Tick(t.Context(), f.clock))
			}
			column := "value_decimal"
			if shape.wide() {
				column = "kwh_1"
			}
			require.Equal(t, []any{"100"}, reconciledHandoffCell(t, f, f.source.GroupID, at(0), column))
			require.Equal(t, []any{"200"}, reconciledHandoffCell(t, f, f.target.GroupID, at(0), column))
			var oldConsumed, newConsumed int
			require.NoError(t, f.local.QueryRow(`SELECT SUM(consumed) FROM wg_delivery_samples WHERE group_id=?`, f.source.GroupID).Scan(&oldConsumed))
			require.NoError(t, f.local.QueryRow(`SELECT SUM(consumed) FROM wg_delivery_samples WHERE group_id=?`, f.target.GroupID).Scan(&newConsumed))
			require.Equal(t, 2, oldConsumed)
			require.Equal(t, 1, newConsumed)
			f.deliver()
			for _, owner := range []struct {
				snapshot   *workspace.WriteGroupAppliedSnapshot
				value      string
				start, end time.Time
			}{{f.source, "100", at(0), f.boundary}, {f.target, "200", f.boundary, at(15)}} {
				query := `SELECT CAST(value_decimal AS TEXT),provenance FROM ` + owner.snapshot.Group.Destination.TableName + ` WHERE group_id=? AND tag_id='tag-kwh-1' AND value_decimal IS NOT NULL`
				if shape.wide() {
					query = `SELECT kwh_1,provenance FROM ` + owner.snapshot.Group.Destination.TableName + ` WHERE group_id=? AND kwh_1 IS NOT NULL`
				}
				var value, provenance string
				require.NoError(t, f.plainTarget().QueryRow(query, owner.snapshot.GroupID).Scan(&value, &provenance))
				require.Equal(t, owner.value, value, "real disposable SQLite destination readback")
				var entries []struct{ Window struct{ Start, End string } }
				if !shape.wide() {
					provenance = "[" + provenance + "]"
				}
				require.NoError(t, json.Unmarshal([]byte(provenance), &entries))
				require.NotEmpty(t, entries)
				for _, entry := range entries {
					start, err := time.Parse(time.RFC3339Nano, entry.Window.Start)
					require.NoError(t, err)
					end, err := time.Parse(time.RFC3339Nano, entry.Window.End)
					require.NoError(t, err)
					require.Equal(t, owner.start, start)
					require.Equal(t, owner.end, end)
				}
			}
		})
	}
}

func TestCutoverReconcilePendingSiblingDoesNotBlockAcceptedBranch(t *testing.T) {
	f, source, target := prepareLiveHandoffFixture(t)
	// The existing independent cumulative fanout has no frozen runtime proof.
	// It must stay pending, not become a required member of the accepted pair.
	const pending = "pending-cumulative"
	_, err := f.local.Exec(`INSERT INTO write_group_versions(group_id,group_revision,effective_at,payload) SELECT ?,group_revision,effective_at,json_set(payload,'$.id',?,'$.recording_profile.kind','cumulative','$.destination.table_name',?) FROM write_group_versions WHERE group_id=?`, pending, pending, pending, f.target.GroupID)
	require.NoError(t, err)
	_, err = f.local.Exec(`UPDATE write_group_cutovers SET detail=json_insert(detail,'$.started[#]',json_object('group_id',?,'revision',?))`, pending, f.target.AppliedRevision)
	require.NoError(t, err)
	_, err = f.local.Exec(`INSERT INTO write_group_ownership_fences SELECT fence_revision,workspace_id,direction,?,group_revision,device_id,point_id,tag_id,'cumulative',connector_id,connector_revision,database_name,table_schema,? FROM write_group_ownership_fences WHERE direction='started'`, pending, pending)
	require.NoError(t, err)
	f.clock = at(40)
	require.ErrorIs(t, f.pipe.Reconcile(t.Context()), groupdelivery.ErrOwnershipUnavailable)
	fresh := f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)]
	require.NotSame(t, target, fresh)
	require.NotSame(t, source, f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)])
	require.Nil(t, f.pipe.boundaries[boundaryKey(pending, f.target.AppliedRevision)])
	require.NoError(t, fresh.boundary.Tick(t.Context(), f.clock))
	require.Equal(t, []any{"7"}, reconciledHandoffCell(t, f, f.target.GroupID, at(15), "value_decimal"))
	var pendingCount, sourceConsumed int
	require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_sample_handoffs WHERE state='pending'`).Scan(&pendingCount))
	require.Equal(t, 1, pendingCount)
	require.NoError(t, f.local.QueryRow(`SELECT SUM(consumed) FROM wg_delivery_samples WHERE group_id=?`, f.source.GroupID).Scan(&sourceConsumed))
	require.Zero(t, sourceConsumed)
}

type afterFirstRestoreBarrier struct {
	Destinations
	calls           atomic.Int32
	entered, resume chan struct{}
}

func (b *afterFirstRestoreBarrier) GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error) {
	if b.calls.Add(1) == 2 {
		close(b.entered)
		select {
		case <-b.resume:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return b.Destinations.GetByID(ctx, id)
}

func TestCutoverReconcileRestoreToPublishIntakeIsNotLost(t *testing.T) {
	f, _, target := prepareLiveHandoffFixture(t)
	require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.target.Group, "native-before", at(24), uint64(9))))
	f.clock = at(22)
	barrier := &afterFirstRestoreBarrier{Destinations: f.pipe.deps.Destinations, entered: make(chan struct{}), resume: make(chan struct{})}
	f.pipe.deps.Destinations = barrier
	var release sync.Once
	defer release.Do(func() { close(barrier.resume) })
	rebuilt := make(chan error, 1)
	go func() { rebuilt <- f.pipe.Reconcile(t.Context()) }()
	<-barrier.entered // group-A sorts first; its verified Restore has completed.
	accepted := make(chan error, 1)
	go func() {
		accepted <- f.pipe.acceptWithOwnership(t.Context(), envelope(f.target.Group, "native-during", at(25), uint64(11)), []*managed{target})
	}()
	release.Do(func() { close(barrier.resume) })
	require.NoError(t, <-rebuilt)
	require.ErrorIs(t, <-accepted, groupdelivery.ErrOwnershipUnavailable, "captured worker cannot ACK an input absent from the replacement")
	require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.target.Group, "native-during", at(25), uint64(11))))
	restored, err := f.store.Restore(t.Context(), groupKey(f.target))
	require.NoError(t, err)
	require.Len(t, restored.Samples, 3)
	f.clock = at(40)
	fresh := f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)]
	require.NoError(t, fresh.boundary.Tick(t.Context(), f.clock))
	require.Equal(t, []any{"9"}, reconciledHandoffCell(t, f, f.target.GroupID, at(15), "value_decimal"))
}

func TestCutoverReconcileCapturedRetirementCannotRemoveReplacement(t *testing.T) {
	f, captured, _ := prepareLiveHandoffFixture(t)
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	key := boundaryKey(f.source.GroupID, f.source.AppliedRevision)
	fresh := f.pipe.boundaries[key]
	require.NotSame(t, captured, fresh)
	captured.boundary.SetUntil(at(0))
	fresh.boundary.SetUntil(at(0))
	require.True(t, captured.boundary.Retired())
	f.pipe.pruneRetired([]*managed{captured})
	require.Same(t, fresh, f.pipe.boundaries[key], "a captured TickAll cleanup only owns its original pointer")
	f.pipe.pruneRetired([]*managed{fresh})
	require.NotContains(t, f.pipe.boundaries, key)
}
