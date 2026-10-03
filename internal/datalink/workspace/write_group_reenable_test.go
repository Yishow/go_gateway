package workspace

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGroupLifecycleReenableSameSemanticsCreatesFreshAppliedRevision(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `
		UPDATE database_connectors SET kind = ?, connection_config = ? WHERE id = ?
	`, "postgres", `{"database":"persisted-db","schema":"persisted-schema"}`, fixture.connectorID)
	require.NoError(t, err)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	service.WithTableInspector(writeGroupLifecycleInspector{})
	now := time.Date(2026, 10, 4, 12, 0, 7, 0, time.UTC)
	service.WithClock(func() time.Time { return now })
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	group := newBasicWriteGroup(fixture)
	group.RowPolicy.IntervalSeconds = 15
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: group,
	})
	require.NoError(t, err)

	beforeApply, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	first, err := service.Apply(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeApply.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	now = time.Date(2026, 10, 4, 12, 0, 16, 0, time.UTC)

	beforeDisable, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	disabled, err := service.Disable(ctx, first.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeDisable.DatabaseSetupRevision,
		ExpectedGroupRevision: first.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	require.Equal(t, first.Group.AppliedRevision, disabled.Group.AppliedRevision)

	beforeReenable, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	reenabled, err := service.Apply(ctx, disabled.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeReenable.DatabaseSetupRevision,
		ExpectedGroupRevision: disabled.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	require.Equal(t, WriteGroupStatusReady, reenabled.Group.Status)
	require.NotEqual(t, disabled.Group.AppliedRevision, reenabled.Group.AppliedRevision,
		"re-enabling must create a fresh intake revision instead of reviving the retired boundary")
	require.Equal(t, reenabled.Group.Revision, reenabled.Group.AppliedRevision)
}
