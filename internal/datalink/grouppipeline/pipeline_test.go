package grouppipeline

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var base = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func at(seconds int) time.Time { return base.Add(time.Duration(seconds) * time.Second) }

type fakeGroups struct {
	groups    []*workspace.WriteGroup
	snapshots map[string][]*workspace.WriteGroupAppliedSnapshot // by group, ascending EffectiveAt
	listErr   error
}

func (f *fakeGroups) List(context.Context) (*workspace.WriteGroupListResult, error) {
	return &workspace.WriteGroupListResult{Groups: f.groups}, f.listErr
}

func (f *fakeGroups) ResolveAppliedAt(_ context.Context, id string, when time.Time) (*workspace.WriteGroupAppliedSnapshot, error) {
	var chosen *workspace.WriteGroupAppliedSnapshot
	for _, snap := range f.snapshots[id] {
		if !snap.EffectiveAt.After(when) {
			chosen = snap
		}
	}
	if chosen == nil {
		return nil, errors.New("no snapshot")
	}
	return chosen, nil
}

type fakeTags struct{ missing map[string]bool }

func (f fakeTags) GetByID(_ context.Context, id string) (*schema.Tag, error) {
	if f.missing[id] {
		return nil, errors.New("gone")
	}
	if id == "tag-p" {
		return &schema.Tag{ID: id, DataType: schema.DataTypeInt64}, nil
	}
	return &schema.Tag{ID: id, DataType: schema.DataTypeFloat64}, nil
}

type fakeDestinations struct {
	connector *schema.DatabaseConnector
	err       error
}

func (f *fakeDestinations) GetByID(context.Context, string) (*schema.DatabaseConnector, error) {
	return f.connector, f.err
}

func (f *fakeDestinations) OpenDestination(context.Context, string, string) (*dbtarget.OpenedDestination, error) {
	return nil, errors.New("not used in these tests")
}

type fakeInspector struct{ inspection *dbtarget.TableInspection }

func (f *fakeInspector) InspectTable(context.Context, string, string, string) (*dbtarget.TableInspection, error) {
	return f.inspection, nil
}

type fixture struct {
	groups  *fakeGroups
	dest    *fakeDestinations
	inspect *fakeInspector
	tags    fakeTags
	store   *groupdelivery.Store
	clock   time.Time
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "gateway.db")+"?_pragma=busy_timeout(5000)")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	group := groupAt("rev-1")
	return &fixture{
		groups: &fakeGroups{
			groups:    []*workspace.WriteGroup{liveGroup(group, workspace.WriteGroupStatusReady)},
			snapshots: map[string][]*workspace.WriteGroupAppliedSnapshot{group.ID: {snapshotOf(group, at(0))}},
		},
		dest: &fakeDestinations{connector: &schema.DatabaseConnector{ID: "connector-1", Kind: schema.DatabaseConnectorKindSQLite, IdentityRevision: "crev-1", Enabled: true}},
		inspect: &fakeInspector{inspection: &dbtarget.TableInspection{
			Status:  dbtarget.TableInspectionExists,
			Columns: []dbtarget.ColumnInfo{{Name: "temperature", DataType: "REAL"}, {Name: "pressure", DataType: "INTEGER"}},
		}},
		tags:  fakeTags{missing: map[string]bool{}},
		store: groupdelivery.NewStore(db),
		clock: at(1),
	}
}

func groupAt(revision string) *workspace.WriteGroup {
	return &workspace.WriteGroup{
		ID: "group-G", WorkspaceID: "workspace-A", Revision: revision, AppliedRevision: revision, Status: workspace.WriteGroupStatusReady,
		Members: []workspace.WriteGroupMember{
			{DeviceID: "device-1", PointID: "point-1", TagID: "tag-t", TargetColumn: "temperature", Required: true, SourceRevision: "src-" + revision, MappingRevision: "map"},
			{DeviceID: "device-1", PointID: "point-2", TagID: "tag-p", TargetColumn: "pressure", Required: true, SourceRevision: "src-" + revision, MappingRevision: "map"},
		},
		Destination: workspace.WriteGroupDestination{ConnectorID: "connector-1", ConnectorRevision: "crev-1", TableSchema: "main", TableName: "readings"},
		RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 10},
	}
}

func liveGroup(group *workspace.WriteGroup, status workspace.WriteGroupStatus) *workspace.WriteGroup {
	live := *group
	live.Status = status
	return &live
}

func snapshotOf(group *workspace.WriteGroup, effective time.Time) *workspace.WriteGroupAppliedSnapshot {
	return &workspace.WriteGroupAppliedSnapshot{
		WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision, AppliedRevision: group.AppliedRevision,
		EffectiveAt: effective, Group: group,
	}
}

func (f *fixture) pipeline() *Pipeline {
	return New(Dependencies{Groups: f.groups, Tags: f.tags, Destinations: f.dest, Inspector: f.inspect, Store: f.store},
		Config{NodeID: "node-1", Owner: "node-1/test", Now: func() time.Time { return f.clock }, ReconcileInterval: 5 * time.Second})
}

func envelope(group *workspace.WriteGroup, id string, observed time.Time, value any) measurement.SampleEnvelope {
	m := group.Members[0]
	return measurement.SampleEnvelope{
		SampleID: id, WorkspaceID: group.WorkspaceID, DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID,
		SourceRevision: m.SourceRevision, MappingRevision: m.MappingRevision, ObservedAt: observed, ReceivedAt: observed,
		TimeOrigin: "gateway", Value: value, Quality: schema.QualityGood,
	}
}

func statusOf(p *Pipeline) GroupStatus {
	for _, status := range p.Status() {
		if status.GroupID == "group-G" {
			return status
		}
	}
	return GroupStatus{}
}

func TestProductionGroupOutageRecoveryPipelineBlocksWithSafeReasonsAndRecovers(t *testing.T) {
	cases := map[string]struct {
		mutate func(*fixture)
		reason string
	}{
		"destination missing": {func(f *fixture) { f.dest.err = dbtarget.ErrConnectorNotFound }, reasonDestinationMissing},
		"destination offline": {func(f *fixture) { f.dest.err = errors.New("dial tcp secret-host:5432: refused") }, reasonDestinationUnavailable},
		"revision changed":    {func(f *fixture) { f.dest.connector.IdentityRevision = "crev-2" }, reasonDestinationRevision},
		"unsupported kind":    {func(f *fixture) { f.dest.connector.Kind = schema.DatabaseConnectorKindMySQL }, reasonDestinationUnsupported},
		"table missing": {func(f *fixture) {
			f.inspect.inspection = &dbtarget.TableInspection{Status: dbtarget.TableInspectionMissing}
		}, reasonTableMissing},
		"inspection forbidden": {func(f *fixture) {
			f.inspect.inspection = &dbtarget.TableInspection{Status: dbtarget.TableInspectionForbidden}
		}, reasonTableUnavailable},
		"tag gone": {func(f *fixture) { f.tags.missing["tag-p"] = true }, reasonTagUnavailable},
		"column type unsupported": {func(f *fixture) {
			// an int64 tag cannot be stored exactly in a REAL column
			f.inspect.inspection.Columns[1] = dbtarget.ColumnInfo{Name: "pressure", DataType: "REAL"}
		}, "layout-blocked"},
	}
	for name, c := range cases {
		f := newFixture(t)
		c.mutate(f)
		p := f.pipeline()
		require.NoError(t, p.Reconcile(t.Context()), name)
		status := statusOf(p)
		require.Equal(t, StateBlocked, status.State, name)
		require.Equal(t, c.reason, status.Reason, name)
		require.NotContains(t, status.Reason, "secret", "a blocked reason is a safe code, never driver text")
		require.False(t, p.Owns("connector-1", "tag-t"), "a group that cannot run owns nothing, so the legacy writer keeps working: "+name)
	}

	f := newFixture(t)
	f.dest.err = errors.New("offline")
	p := f.pipeline()
	require.NoError(t, p.Reconcile(t.Context()))
	require.Equal(t, StateBlocked, statusOf(p).State)
	f.dest.err = nil
	require.NoError(t, p.Reconcile(t.Context()))
	require.Equal(t, StateActive, statusOf(p).State, "the next reconcile picks the group up once its destination is back")
	require.True(t, p.Owns("connector-1", "tag-t"))
}

func TestProductionGroupOutageRecoveryPipelineDisableEndsIntakeAtTheNextBoundary(t *testing.T) {
	f := newFixture(t)
	p := f.pipeline()
	require.NoError(t, p.Reconcile(t.Context()))
	require.Equal(t, StateActive, statusOf(p).State)
	group := f.groups.snapshots["group-G"][0].Group

	f.clock = at(12)
	require.NoError(t, p.AcceptSample(t.Context(), envelope(group, "s-before", at(11), 21.5)))
	f.groups.groups = []*workspace.WriteGroup{liveGroup(group, workspace.WriteGroupStatusDisabled)}
	require.NoError(t, p.Reconcile(t.Context()))
	require.Equal(t, StateRetiring, statusOf(p).State)
	require.True(t, p.Owns("connector-1", "tag-t"), "the group still owns its output while its last bucket closes")

	require.NoError(t, p.AcceptSample(t.Context(), envelope(group, "s-after", at(21), 22.5)), "intake after the end boundary is ignored, not an error")
	f.clock = at(25)
	p.TickAll(t.Context())
	require.Empty(t, p.Status(), "once its closing buckets are done the group is gone from the pipeline")
	require.False(t, p.Owns("connector-1", "tag-t"))
	var buckets int
	require.NoError(t, rowCount(f.store, &buckets, "wg_delivery_buckets"))
	require.Positive(t, buckets, "the buckets before the disable were closed and recorded")
}

func TestProductionGroupOutageRecoveryPipelineSupersededRevisionHandsOverAtTheEffectiveBoundary(t *testing.T) {
	f := newFixture(t)
	first := groupAt("rev-1")
	second := groupAt("rev-2")
	f.groups.groups = []*workspace.WriteGroup{liveGroup(second, workspace.WriteGroupStatusReady)}
	f.groups.snapshots["group-G"] = []*workspace.WriteGroupAppliedSnapshot{snapshotOf(first, at(0)), snapshotOf(second, at(20))}
	f.clock = at(12)
	p := f.pipeline()
	require.NoError(t, p.Reconcile(t.Context()))
	require.Len(t, p.Status(), 2, "the new revision is prepared before it becomes effective")

	// A sample in [10,20) belongs to rev-1; one in [20,30) belongs to rev-2.
	require.NoError(t, p.AcceptSample(t.Context(), envelope(first, "old", at(15), 1.0)))
	f.clock = at(22)
	require.NoError(t, p.AcceptSample(t.Context(), envelope(second, "new", at(21), 2.0)))
	var journaled int
	require.NoError(t, rowCount(f.store, &journaled, "wg_delivery_samples"))
	require.Equal(t, 2, journaled)
	var rev1, rev2 int
	require.NoError(t, rowCountWhere(f.store, &rev1, "wg_delivery_samples", "group_revision = 'rev-1'"))
	require.NoError(t, rowCountWhere(f.store, &rev2, "wg_delivery_samples", "group_revision = 'rev-2'"))
	require.Equal(t, 1, rev1)
	require.Equal(t, 1, rev2, "each sample is journaled under exactly one revision")

	f.clock = at(40)
	p.TickAll(t.Context())
	statuses := p.Status()
	require.Len(t, statuses, 1, "the superseded revision retires after closing its last bucket")
	require.Equal(t, "rev-2", statuses[0].Revision)
}

func TestProductionGroupOutageRecoveryPipelineFanOutJoinsRefusalsAndStaysIsolated(t *testing.T) {
	f := newFixture(t)
	p := f.pipeline()
	require.NoError(t, p.Reconcile(t.Context()))
	group := f.groups.snapshots["group-G"][0].Group
	f.clock = at(30)
	p.TickAll(t.Context()) // closes the buckets before t=30

	other := envelope(group, "x", at(29), 1.0)
	other.TagID = "tag-unrelated"
	require.NoError(t, p.AcceptSample(t.Context(), other), "a sample for no group is not an error")
	late := envelope(group, "late", at(1), 1.0)
	err := p.AcceptSample(t.Context(), late)
	require.Error(t, err)
	var refused interface{ Error() string }
	require.ErrorAs(t, err, &refused)
	require.Contains(t, err.Error(), "late-after-close", "an expired sample is reported as late")
}

func rowCount(store *groupdelivery.Store, into *int, table string) error {
	return rowCountWhere(store, into, table, "1 = 1")
}

func rowCountWhere(store *groupdelivery.Store, into *int, table, where string) error {
	return store.DB().QueryRowContext(context.Background(), `SELECT COUNT(*) FROM `+table+` WHERE `+where).Scan(into)
}

func TestProductionGroupOutageRecoveryBacklogGuardRefusesWhenOwnershipCannotBeProven(t *testing.T) {
	f := newFixture(t)
	db := f.store.DB()
	ctx := t.Context()
	group := &workspace.WriteGroup{ID: "group-G"}
	check := func() error {
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer func() { _ = tx.Rollback() }()
		return BacklogGuard{}.CheckWriteGroupBacklog(ctx, tx, group)
	}
	require.NoError(t, check(), "no backlog means nothing can be orphaned")

	_, err := db.ExecContext(ctx, `INSERT INTO write_groups (id, workspace_id, revision, name, destination_connector_id, destination_connector_revision, destination_table_name)
		VALUES ('group-G', 'workspace-A', 'rev-1', 'g', 'connector-1', 'crev-1', 'readings')`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO write_group_versions (group_id, group_revision, effective_at, payload) VALUES ('group-G', 'rev-1', '2026-01-01T00:00:00Z', '{}')`)
	require.NoError(t, err)
	insertRow := func(effect, revision, connectorRevision, state string) {
		_, err := db.ExecContext(ctx, `INSERT INTO wg_delivery_outbox (effect_key, record_id, workspace_id, group_id, group_revision, bucket_start, partition_key,
			destination_scope, connector_id, connector_revision, table_name, payload, payload_digest, state, next_retry_at, created_at, updated_at)
			VALUES (?, 'r', 'workspace-A', 'group-G', ?, '2026-01-01T00:00:00Z', 'p', 's', 'connector-1', ?, 'readings', '{}', 'd', ?, 'x', 'x', 'x')`,
			effect, revision, connectorRevision, state)
		require.NoError(t, err)
	}
	insertRow("e1", "rev-1", "crev-1", "pending")
	require.NoError(t, check(), "a backlog row that names its immutable revision and frozen destination is owned")

	insertRow("e2", "rev-9", "crev-1", "retrying")
	require.ErrorIs(t, check(), ErrBacklogOwnershipUnproven, "a revision without an immutable version cannot be tied to an owner")
	_, err = db.ExecContext(ctx, `DELETE FROM wg_delivery_outbox WHERE effect_key = 'e2'`)
	require.NoError(t, err)

	insertRow("e3", "rev-1", "", "blocked")
	require.ErrorIs(t, check(), ErrBacklogOwnershipUnproven, "a row without its frozen destination revision is refused")
	_, err = db.ExecContext(ctx, `DELETE FROM wg_delivery_outbox WHERE effect_key = 'e3'`)
	require.NoError(t, err)

	insertRow("e4", "rev-9", "crev-1", "sql_committed")
	require.NoError(t, check(), "finished rows no longer need an owner")
	require.Error(t, BacklogGuard{}.CheckWriteGroupBacklog(ctx, nil, group))
}
