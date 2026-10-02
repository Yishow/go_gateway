package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReviewRowGroupMigrationBlocksOverlappingSourceIDsAsOneBatch(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newRowGroupMigrationFixture(ctx, t, db)

	workspaceSvc := NewService(NewSQLRepository(db))
	workspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	workspace.DatabaseRowGroups = append(workspace.DatabaseRowGroups, DatabaseRowGroup{
		ID:               "legacy-g-duplicate",
		ConnectorID:      fixture.connectorID,
		TableSchema:      "main",
		TableName:        "raw_values",
		MemberPointIDs:   []string{fixture.pointIDs[0]},
		GroupKeyColumns:  []string{"entity_id"},
		UniqueKeyColumns: []string{"entity_id", "observed_at"},
	})
	_, err = workspaceSvc.SaveDatabaseRowGroups(ctx, fixture.connectorID, "main", "raw_values", workspace.DatabaseRowGroups)
	require.NoError(t, err)

	beforeWorkspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	var beforeGroups, beforeMaps int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups`).Scan(&beforeGroups))
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_group_migration_maps`).Scan(&beforeMaps))

	preview, err := service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID, "legacy-g-duplicate"})
	require.NoError(t, err)
	require.Len(t, preview.Items, 2)
	for _, item := range preview.Items {
		require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, item.Status)
		requireFindingCode(t, item.Issues, "row-group-source-overlap")
		require.Nil(t, item.CandidateGroup)
	}

	_, err = service.ReviewRowGroupMigration(ctx, WriteGroupMigrationReviewRequest{
		WorkspaceID:               preview.WorkspaceID,
		ExpectedWorkspaceRevision: preview.WorkspaceRevision,
		ExpectedConnectorRevision: preview.ConnectorRevision,
		ReviewDigest:              preview.ReviewDigest,
		SourceIDs:                 []string{fixture.rowGroupID, "legacy-g-duplicate"},
		ConfirmSnapshotConversion: true,
	})
	require.ErrorIs(t, err, ErrWriteGroupValidation)

	afterWorkspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Equal(t, beforeWorkspace, afterWorkspace)
	var afterGroups, afterMaps int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups`).Scan(&afterGroups))
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_group_migration_maps`).Scan(&afterMaps))
	require.Equal(t, beforeGroups, afterGroups)
	require.Equal(t, beforeMaps, afterMaps)
}
