package grouppipeline

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

type ownershipPipelineFixture struct {
	*longPipeline
	source   *workspace.WriteGroupAppliedSnapshot
	target   *workspace.WriteGroupAppliedSnapshot
	boundary time.Time
}

// Reuse the supported managed target and real disposable SQLite stores. Only
// reconciliation is fake/controlled; commit owns a real local SQL transaction.
// No Start/worker goroutine or wall-clock sleep participates in these tests.
func newOwnershipPipelineFixture(t *testing.T, shape averageShape) *ownershipPipelineFixture {
	t.Helper()
	f := newAveragePipeline(t, shape)
	boundary := at(22).Add(500 * time.Millisecond)
	var target workspace.WriteGroupAppliedSnapshot
	raw, err := f.store.RuntimeVersion(t.Context(), groupKey(f.snapshot))
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, &target))
	target.EffectiveAt = boundary
	raw, err = json.Marshal(target)
	require.NoError(t, err)
	// Initial fixture construction: no input has ever been accepted by this target.
	_, err = f.local.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE group_id=?`, string(raw), target.GroupID)
	require.NoError(t, err)
	var old workspace.WriteGroup
	raw, err = json.Marshal(target.Group)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, &old))
	old.ID, old.Revision, old.AppliedRevision = "legacy-source", "source-rev", "source-rev"
	old.LayoutVersion, old.StorageMode, old.RecordingProfile = "", "", nil
	old.Destination.TableName, old.Destination.StorageStrategy = "legacy_readings", workspace.WriteGroupStorageStrategyCustom
	old.RowPolicy = workspace.WriteGroupRowPolicy{IntervalSeconds: 10}
	old.WritePolicy = workspace.WriteGroupWritePolicy{Mode: "append", DedupeCapability: groupdelivery.DedupeNone}
	var columns []dbtarget.ColumnInfo
	for i := range old.Members {
		old.Members[i].ValueRole = ""
		old.Members[i].Required = true
		old.Members[i].SourceValueType, old.Members[i].ResultValueType = "", ""
		columns = append(columns, dbtarget.ColumnInfo{Name: old.Members[i].TargetColumn, DataType: "TEXT", Nullable: true})
	}
	source := snapshotOf(&old, at(0))
	source.RuntimeLayout = &workspace.WriteGroupRuntimeLayout{Dialect: dbtarget.SQLDialectSQLite, Columns: columns, TagTypes: map[string]schema.DataType{longKwhFirst: schema.DataTypeUint64, longKwhSecond: schema.DataTypeUint64}}
	f.groups = &fakeGroups{groups: []*workspace.WriteGroup{&old}, snapshots: map[string][]*workspace.WriteGroupAppliedSnapshot{old.ID: {source}}}
	f.pipe = New(Dependencies{Groups: f.groups, Destinations: f.dest, Store: f.store}, Config{Now: func() time.Time { return f.clock }, ReconcileInterval: 5 * time.Second})
	require.NoError(t, f.pipe.saveRuntimeSnapshot(t.Context(), source))
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	require.Len(t, f.pipe.boundaries, 1, "%+v", f.pipe.Status())
	return &ownershipPipelineFixture{longPipeline: f, source: source, target: &target, boundary: boundary}
}

func persistOwnershipVersion(t *testing.T, tx *sql.Tx, s *workspace.WriteGroupAppliedSnapshot) {
	t.Helper()
	g := s.Group
	_, err := tx.Exec(`INSERT INTO write_groups(id,workspace_id,revision,applied_revision,name,status,destination_connector_id,destination_connector_revision,destination_table_name) VALUES(?,?,?,?,?,'ready',?,?,?)`, g.ID, g.WorkspaceID, g.Revision, g.AppliedRevision, g.ID, g.Destination.ConnectorID, g.Destination.ConnectorRevision, g.Destination.TableName)
	require.NoError(t, err)
	raw, err := json.Marshal(g)
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_group_versions(group_id,group_revision,effective_at,payload) VALUES(?,?,?,?)`, g.ID, g.Revision, s.EffectiveAt, string(raw))
	require.NoError(t, err)
}

func (f *ownershipPipelineFixture) commit() {
	t := f.t
	t.Helper()
	tx, err := f.local.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	persistOwnershipVersion(t, tx, f.source)
	persistOwnershipVersion(t, tx, f.target)
	moved := make([]workspace.CutoverMemberRef, 0, len(f.source.Group.Members))
	for _, m := range f.source.Group.Members {
		moved = append(moved, workspace.CutoverMemberRef{DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID})
	}
	detail, err := json.Marshal(map[string]any{"ended": []workspace.CutoverEnd{{GroupID: f.source.GroupID, Revision: f.source.AppliedRevision}}, "started": []workspace.CutoverStart{{GroupID: f.target.GroupID, Revision: f.target.AppliedRevision}}, "moved": moved, "ownership_fence_revision": "cutover-f2"})
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_group_cutovers(id,workspace_id,kind,source_group_id,boundary,review_digest,detail,intent) VALUES('cutover-f2',?,'forward',?,?, 'review-f2',?,'{}')`, f.source.WorkspaceID, f.source.GroupID, f.boundary.Format(time.RFC3339Nano), string(detail))
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_group_revision_ends(workspace_id,group_id,group_revision,ends_at,cutover_id) VALUES(?,?,?,?,'cutover-f2')`, f.source.WorkspaceID, f.source.GroupID, f.source.AppliedRevision, f.boundary.Format(time.RFC3339Nano))
	require.NoError(t, err)
	for direction, s := range map[string]*workspace.WriteGroupAppliedSnapshot{"ended": f.source, "started": f.target} {
		g := s.Group
		kind := ""
		if g.RecordingProfile != nil {
			kind = string(g.RecordingProfile.Kind)
		}
		for _, m := range g.Members {
			_, err = tx.Exec(`INSERT INTO write_group_ownership_fences(fence_revision,workspace_id,direction,group_id,group_revision,device_id,point_id,tag_id,profile_kind,connector_id,connector_revision,database_name,table_schema,table_name) VALUES('cutover-f2',?,?,?,?,?,?,?,?,?,?,?,?,?)`, g.WorkspaceID, direction, g.ID, s.AppliedRevision, m.DeviceID, m.PointID, m.TagID, kind, g.Destination.ConnectorID, g.Destination.ConnectorRevision, g.Destination.Database, g.Destination.TableSchema, g.Destination.TableName)
			require.NoError(t, err)
		}
	}
	require.NoError(t, tx.Commit())
}

func (f *ownershipPipelineFixture) readyTarget() {
	f.t.Helper()
	f.groups.groups = append(f.groups.groups, f.target.Group)
	f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
	require.NoError(f.t, f.pipe.Reconcile(f.t.Context()))
	require.NotNil(f.t, f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)])
}

func (f *ownershipPipelineFixture) readyHealthy() *workspace.WriteGroupAppliedSnapshot {
	t := f.t
	t.Helper()
	raw, err := json.Marshal(f.source)
	require.NoError(t, err)
	var healthy workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(raw, &healthy))
	healthy.GroupID, healthy.Group.ID = "healthy-group", "healthy-group"
	healthy.Group.Members = healthy.Group.Members[:1]
	healthy.Group.Members[0].DeviceID, healthy.Group.Members[0].PointID, healthy.Group.Members[0].TagID = "healthy-device", "healthy-point", "healthy-tag"
	healthy.RuntimeLayout.TagTypes = map[string]schema.DataType{"healthy-tag": schema.DataTypeUint64}
	healthy.RuntimeLayout.Columns = healthy.RuntimeLayout.Columns[:1]
	require.NoError(t, f.pipe.saveRuntimeSnapshot(t.Context(), &healthy))
	f.groups.groups = append(f.groups.groups, healthy.Group)
	f.groups.snapshots[healthy.GroupID] = []*workspace.WriteGroupAppliedSnapshot{&healthy}
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	require.NotNil(t, f.pipe.boundaries[boundaryKey(healthy.GroupID, healthy.AppliedRevision)])
	return &healthy
}
