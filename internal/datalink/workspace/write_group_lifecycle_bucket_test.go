package workspace

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGroupLifecycleIntervalReductionWaitsForOldBucket(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture, created := newReadinessGroup(ctx, t, db)
	service.WithTableInspector(writeGroupLifecycleInspector{})
	now := time.Date(2026, 10, 2, 12, 0, 7, 0, time.UTC)
	service.WithClock(func() time.Time { return now })

	workspaceBeforeFirstEdit, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	firstConfig := cloneWriteGroup(created.Group)
	firstConfig.RowPolicy.IntervalSeconds = 60
	firstDraft, err := service.Update(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforeFirstEdit.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1", Group: firstConfig,
	})
	require.NoError(t, err)
	workspaceBeforeFirstApply, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	first, err := service.Apply(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforeFirstApply.DatabaseSetupRevision,
		ExpectedGroupRevision: firstDraft.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	firstEffectiveAt := time.Date(2026, 10, 2, 12, 1, 0, 0, time.UTC)
	firstSnapshot, err := service.ResolveAppliedAt(ctx, first.Group.ID, firstEffectiveAt)
	require.NoError(t, err)
	require.Equal(t, first.Group.Revision, firstSnapshot.GroupRevision)

	now = time.Date(2026, 10, 2, 12, 1, 10, 0, time.UTC)
	workspaceBeforeReduction, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	reducedConfig := cloneWriteGroup(first.Group)
	reducedConfig.RowPolicy.IntervalSeconds = 15
	reducedDraft, err := service.Update(ctx, first.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforeReduction.DatabaseSetupRevision,
		ExpectedGroupRevision: first.Group.Revision, ExpectedConnectorRevision: "connector-revision-1", Group: reducedConfig,
	})
	require.NoError(t, err)
	workspaceBeforeReductionApply, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	second, err := service.Apply(ctx, first.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforeReductionApply.DatabaseSetupRevision,
		ExpectedGroupRevision: reducedDraft.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)

	oldBucketEnd := time.Date(2026, 10, 2, 12, 2, 0, 0, time.UTC)
	oldAtCutover, err := service.ResolveAppliedAt(ctx, second.Group.ID, oldBucketEnd.Add(-time.Nanosecond))
	require.NoError(t, err)
	require.Equal(t, first.Group.Revision, oldAtCutover.GroupRevision)
	newAtCutover, err := service.ResolveAppliedAt(ctx, second.Group.ID, oldBucketEnd)
	require.NoError(t, err)
	require.Equal(t, second.Group.Revision, newAtCutover.GroupRevision)
	require.Equal(t, oldBucketEnd, newAtCutover.EffectiveAt)
}

func TestWriteGroupIntervalArithmeticRejectsDurationOverflow(t *testing.T) {
	_, err := writeGroupCommonInterval(
		time.Duration(2_000_000_000)*time.Second,
		time.Duration(2_000_000_001)*time.Second,
	)
	require.ErrorIs(t, err, ErrWriteGroupIntervalOverflow)

	if int64(^uint(0)>>1) > maxWriteGroupDurationSeconds {
		_, err = writeGroupIntervalDuration(int(maxWriteGroupDurationSeconds + 1))
		require.ErrorIs(t, err, ErrWriteGroupIntervalOverflow)
	}
}
