package workspace

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func counterCutoverFixture(t *testing.T) (counterApplyFixture, *WriteGroup, CutoverRequest) {
	t.Helper()
	f := newCounterApplyFixture(t, WriteGroupLayoutManagedWideProfileV1, WriteGroupProfileCumulative, schema.DataTypeInt64)
	source, err := f.service.Apply(t.Context(), f.group.ID, f.mutation)
	require.NoError(t, err)
	f.group = source.Group
	draft := cloneWriteGroup(f.group)
	draft.Name, draft.LayoutVersion, draft.StorageMode = "counter_raw target", WriteGroupLayoutManagedLongV1, WriteGroupStorageModePlain
	draft.Destination.TableName = ""
	draft.RecordingProfile.Kind, draft.RecordingProfile.Interval, draft.RecordingProfile.Unit = WriteGroupProfileInstantaneous, 15, WriteGroupIntervalUnitSecond
	draft.RowPolicy.IntervalSeconds = 15
	draft.Members[0].ValueRole = WriteGroupValueRoleCounterRaw
	created, err := f.service.Create(t.Context(), WriteGroupMutation{WorkspaceID: f.group.WorkspaceID, ExpectedWorkspaceRevision: source.WorkspaceRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: draft})
	require.NoError(t, err)
	saveManagedReadinessProof(t, f.db, created.Group.ID)
	target, err := f.service.Get(t.Context(), created.Group.ID)
	require.NoError(t, err)
	layout, err := managedSchemaLayout(target.Group, map[string]schema.DataType{target.Group.Members[0].TagID: schema.DataTypeInt64}, "postgres")
	require.NoError(t, err)
	f.inspector.otherLayouts = map[string]*recordingplan.GroupSchemaLayout{layout.TableName: layout}
	req := CutoverRequest{Kind: WriteGroupCutoverForward, WorkspaceID: f.group.WorkspaceID, ExpectedWorkspaceRevision: created.WorkspaceRevision,
		SourceGroupID: f.group.ID, ExpectedSourceRevision: f.group.Revision, ExpectedConnectorRevision: "connector-revision-1",
		Targets:  []CutoverTargetRef{{GroupID: target.Group.ID, ExpectedRevision: target.Group.Revision}},
		Boundary: time.Date(2026, 10, 1, 0, 20, 0, 0, time.UTC)}
	review, err := f.service.ReviewCutover(t.Context(), req)
	require.NoError(t, err)
	require.True(t, review.Applicable, "%v", review.Issues)
	req.ReviewDigest = review.Digest
	return f, target.Group, req
}

func reviewCounterRollback(t *testing.T, f counterApplyFixture, forward *WriteGroupCutover) CutoverRequest {
	t.Helper()
	current, err := f.service.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	req := CutoverRequest{Kind: WriteGroupCutoverRollback, WorkspaceID: f.group.WorkspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		SourceGroupID: f.group.ID, ExpectedSourceRevision: f.group.Revision, ExpectedConnectorRevision: "connector-revision-1",
		CutoverID: forward.ID, Boundary: forward.Boundary.Add(10 * time.Minute)}
	review, err := f.service.ReviewCutover(t.Context(), req)
	require.NoError(t, err)
	require.True(t, review.Applicable, "%v", review.Issues)
	req.ReviewDigest = review.Digest
	return req
}

func TestCounterBindingCutoverReviewedDefinitionEditIsStale(t *testing.T) {
	f, target, req := counterCutoverFixture(t)
	_, original := readCounterApplySnapshot(t, f.db.QueryRowContext, f.group.ID)
	_, err := f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET definition_revision='def-2', unit='MWh', series_epoch='epoch-2' WHERE id='measurement-1'`)
	require.NoError(t, err)
	_, err = f.service.ApplyCutover(t.Context(), req)
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	loaded, err := f.service.Get(t.Context(), target.ID)
	require.NoError(t, err)
	require.Empty(t, loaded.Group.AppliedRevision)
	_, after := readCounterApplySnapshot(t, f.db.QueryRowContext, f.group.ID)
	require.Equal(t, original, after)
}

func TestCounterBindingCutoverRollbackPersistsOnlyNewRevisionProof(t *testing.T) {
	f, target, req := counterCutoverFixture(t)
	_, original := readCounterApplySnapshot(t, f.db.QueryRowContext, f.group.ID)
	forward, err := f.service.ApplyCutover(t.Context(), req)
	require.NoError(t, err)
	targetSnapshot, _ := readCounterApplySnapshot(t, f.db.QueryRowContext, target.ID)
	require.Len(t, targetSnapshot.CounterBindings, 1, "forward uses the shared apply proof transaction")
	_, err = f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET definition_revision='def-2', unit='MWh', series_epoch='epoch-2' WHERE id='measurement-1'`)
	require.NoError(t, err)
	rollback := reviewCounterRollback(t, f, forward)
	restored, err := f.service.ApplyCutover(t.Context(), rollback)
	require.NoError(t, err)
	require.Len(t, restored.Started, 1)
	require.NotEqual(t, f.group.AppliedRevision, restored.Started[0].Revision)
	var oldPayload, newPayload string
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT payload FROM wg_runtime_versions WHERE group_id=? AND group_revision=?`, f.group.ID, f.group.AppliedRevision).Scan(&oldPayload))
	require.Equal(t, original, oldPayload)
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT payload FROM wg_runtime_versions WHERE group_id=? AND group_revision=?`, f.group.ID, restored.Started[0].Revision).Scan(&newPayload))
	var snapshot WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal([]byte(newPayload), &snapshot))
	require.Len(t, snapshot.CounterBindings, 1)
	for _, proof := range snapshot.CounterBindings {
		require.Equal(t, restored.Started[0].Revision, proof.Member.AppliedRevision)
		require.Equal(t, "def-2", proof.DefinitionRevision)
		require.Equal(t, "epoch-2", proof.SeriesEpoch)
		require.Equal(t, "MWh", proof.Unit)
	}
}

func TestCounterBindingCutoverRollbackCASAfterReview(t *testing.T) {
	f, _, req := counterCutoverFixture(t)
	forward, err := f.service.ApplyCutover(t.Context(), req)
	require.NoError(t, err)
	rollback := reviewCounterRollback(t, f, forward)
	// Rollback reads source and target before its two readiness snapshots.
	store := &counterApplyRaceRepository{SQLRepository: NewSQLRepository(f.db), triggerAt: 5, beforeApply: func() {
		_, err := f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET source_binding_revision='binding-2' WHERE id='measurement-1'`)
		require.NoError(t, err)
	}}
	f.service.workspaceSvc.repo = store
	_, err = f.service.ApplyCutover(t.Context(), rollback)
	require.Equal(t, 5, store.begins)
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	source, err := f.service.Get(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Equal(t, f.group.AppliedRevision, source.Group.AppliedRevision)
	var count int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_runtime_versions WHERE group_id=?`, f.group.ID).Scan(&count))
	require.Equal(t, 1, count)
}
