package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicGroupReadiness(t *testing.T) {
	for _, strategy := range []WriteGroupStorageStrategy{
		WriteGroupStorageStrategyManaged,
		WriteGroupStorageStrategyCustom,
	} {
		t.Run(string(strategy), func(t *testing.T) {
			ctx := t.Context()
			db := openWorkspaceTestDB(t, ":memory:")
			defer db.Close()
			fixture := seedWriteGroupFixture(ctx, t, db)
			workspaceSvc := NewService(NewSQLRepository(db))
			service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
			current, err := workspaceSvc.GetOrCreate(ctx)
			require.NoError(t, err)
			group := newBasicWriteGroup(fixture)
			group.Destination.StorageStrategy = strategy
			group.RowPolicy.IntervalSeconds = 15
			group.WritePolicy.Mode = "append"
			created, err := service.Create(ctx, WriteGroupMutation{
				WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
				ExpectedConnectorRevision: "connector-revision-1", Group: group,
			})
			require.NoError(t, err)

			readiness, err := service.Readiness(ctx, created.Group.ID)
			require.NoError(t, err)
			require.Equal(t, fixture.workspaceID, readiness.WorkspaceID)
			require.Equal(t, created.WorkspaceRevision, readiness.WorkspaceRevision)
			require.Equal(t, created.Group.ID, readiness.GroupID)
			require.Equal(t, created.Group.Revision, readiness.GroupRevision)
			require.Empty(t, readiness.AppliedRevision)
			require.True(t, readiness.ConfigReady)
			require.False(t, readiness.SchemaReady)
			require.False(t, readiness.Ready)
			require.Empty(t, readiness.SchemaDigest)
			require.NotEmpty(t, readiness.Issues)
		})
	}
}
