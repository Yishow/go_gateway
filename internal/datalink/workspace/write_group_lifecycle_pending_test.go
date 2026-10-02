package workspace

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGroupLifecyclePendingApplyCannotReactivateSupersededPolicy(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture, created := newReadinessGroup(ctx, t, db)
	service.WithTableInspector(writeGroupLifecycleInspector{})
	now := time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC)
	service.WithClock(func() time.Time { return now })
	edit := func(group *WriteGroup, workspaceRevision string, interval int) *WriteGroupSaveResult {
		t.Helper()
		candidate := cloneWriteGroup(group)
		candidate.RowPolicy.IntervalSeconds = interval
		result, err := service.Update(ctx, group.ID, WriteGroupMutation{
			WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceRevision,
			ExpectedGroupRevision: group.Revision, ExpectedConnectorRevision: "connector-revision-1", Group: candidate,
		})
		require.NoError(t, err)
		return result
	}
	apply := func(saved *WriteGroupSaveResult) (*WriteGroupSaveResult, error) {
		return service.Apply(ctx, saved.Group.ID, WriteGroupMutation{
			WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: saved.WorkspaceRevision,
			ExpectedGroupRevision: saved.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
		})
	}
	first, err := apply(edit(created.Group, created.WorkspaceRevision, 6))
	require.NoError(t, err)
	now = now.Add(7 * time.Second)
	second, err := apply(edit(first.Group, first.WorkspaceRevision, 10))
	require.NoError(t, err)
	now = now.Add(time.Second)
	thirdDraft := edit(second.Group, second.WorkspaceRevision, 5)
	_, err = apply(thirdDraft)
	require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
	var count int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_group_versions WHERE group_id = ?`, first.Group.ID).Scan(&count))
	require.Equal(t, 2, count)
	resolved, err := service.ResolveAppliedAt(ctx, first.Group.ID, time.Date(2026, 10, 2, 12, 0, 31, 0, time.UTC))
	require.NoError(t, err)
	require.Equal(t, second.Group.AppliedRevision, resolved.GroupRevision)
	require.Equal(t, 10, resolved.Group.RowPolicy.IntervalSeconds)
}
