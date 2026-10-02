package workspace

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSingleMappingMigrationRetainsReviewedDatabaseAcrossReload(t *testing.T) {
	ctx := t.Context()
	dbPath := filepath.Join(t.TempDir(), "review-database.db")
	db := openWorkspaceTestDB(t, dbPath)
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
	require.Equal(t, "persisted-db", preview.Items[0].BeforeIntent.Database)
	require.Equal(t, "persisted-db", preview.Items[0].CandidateGroup.Destination.Database)
	result, err := service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(preview))
	require.NoError(t, err)
	require.NoError(t, db.Close())

	reloaded := openWorkspaceTestDB(t, dbPath)
	defer reloaded.Close()
	var raw string
	require.NoError(t, reloaded.QueryRowContext(ctx, `SELECT before_intent FROM write_group_migration_maps
		WHERE workspace_id = ? AND source_kind = ? AND source_id = ?`,
		fixture.workspaceID, "legacy-single-mapping", "legacy-A").Scan(&raw))
	var intent WriteGroupMigrationIntent
	require.NoError(t, json.Unmarshal([]byte(raw), &intent))
	require.Equal(t, *preview.Items[0].BeforeIntent, intent)
	require.Equal(t, "persisted-db", intent.Database)

	restarted := NewWriteGroupService(NewService(NewSQLRepository(reloaded)), NewSQLWriteGroupRepository(reloaded))
	current := mustSingleMappingMigrationPreview(t, restarted, fixture.workspaceID)
	_, err = reloaded.ExecContext(ctx, `UPDATE database_connectors SET connection_config = ? WHERE id = ?`,
		`{"database":"new-database","schema":"main"}`, fixture.connectorID)
	require.NoError(t, err)
	_, err = restarted.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(current))
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	group, err := restarted.Get(ctx, result.Groups[0].ID)
	require.NoError(t, err)
	require.Equal(t, "persisted-db", group.Group.Destination.Database)
	require.NoError(t, reloaded.QueryRowContext(ctx, `SELECT before_intent FROM write_group_migration_maps
		WHERE workspace_id = ? AND source_kind = ? AND source_id = ?`,
		fixture.workspaceID, "legacy-single-mapping", "legacy-A").Scan(&raw))
	require.NoError(t, json.Unmarshal([]byte(raw), &intent))
	require.Equal(t, "persisted-db", intent.Database)
}
