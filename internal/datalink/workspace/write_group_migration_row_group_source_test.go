package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLegacyRowGroupMigrationBlocksTagWithAnotherPointSource(t *testing.T) {
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newRowGroupMigrationFixture(t.Context(), t, db)
	_, otherPointID := seedWriteGroupDevice(t.Context(), t, db)
	_, err := db.ExecContext(t.Context(), `INSERT INTO mappings
		(id,point_id,tag_id,transform_pipeline,status,enabled)
		VALUES ('competing-tag-source',?,?,'[]','active',1)`, otherPointID, fixture.tagIDs[0])
	require.NoError(t, err)
	preview, err := service.PreviewRowGroupMigration(t.Context(), fixture.workspaceID, []string{fixture.rowGroupID})
	require.NoError(t, err)
	require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, preview.Items[0].Status)
	require.Nil(t, preview.Items[0].CandidateGroup)
	requireFindingCode(t, preview.Items[0].Issues, "multiple-enabled-source-mappings")
	_, err = service.ReviewRowGroupMigration(t.Context(), rowGroupReviewRequest(preview))
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count(*) FROM write_groups`).Scan(&count))
	require.Zero(t, count)
}
