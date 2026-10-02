package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReviewSingleMappingMigrationSavesDraftAndPreservesLegacyWriter(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	workspaceBefore, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)

	preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
	result, err := service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(preview))
	require.NoError(t, err)
	require.Equal(t, fixture.workspaceID, result.WorkspaceID)
	require.NotEqual(t, workspaceBefore.DatabaseSetupRevision, result.WorkspaceRevision)
	require.Equal(t, preview.ConnectorRevision, result.ConnectorRevision)
	require.Len(t, result.Groups, 1)
	group := result.Groups[0]
	require.NotEmpty(t, group.ID)
	require.Equal(t, WriteGroupStatusDraft, group.Status)
	require.Empty(t, group.AppliedRevision)
	require.Equal(t, writeGroupMigrationReviewResultConfirmed, group.Migration.ReviewResult)
	require.Equal(t, []string{"legacy-A"}, group.Migration.SourceIDs)
	require.Equal(t, 15, group.RowPolicy.IntervalSeconds)
	require.NotNil(t, group.Members[0].MaxAgeSeconds)
	require.Equal(t, 15, *group.Members[0].MaxAgeSeconds)

	var tableSchema, tableName, columnName, writeMode, timestampColumn string
	var enabled bool
	require.NoError(t, db.QueryRowContext(ctx, `
		SELECT table_schema, table_name, column_name, write_mode, timestamp_column, enabled
		FROM database_target_mappings WHERE id = ?
	`, "legacy-A").Scan(&tableSchema, &tableName, &columnName, &writeMode, &timestampColumn, &enabled))
	require.Equal(t, "main", tableSchema)
	require.Equal(t, "raw_values", tableName)
	require.Equal(t, "temperature", columnName)
	require.Equal(t, "insert", writeMode)
	require.Equal(t, "observed_at", timestampColumn)
	require.True(t, enabled)

	var sourceRevision, beforeIntent, reviewDigest, adapterVersion string
	require.NoError(t, db.QueryRowContext(ctx, `
		SELECT source_revision, before_intent, review_digest, adapter_version
		FROM write_group_migration_maps
		WHERE workspace_id = ? AND source_kind = ? AND source_id = ?
	`, fixture.workspaceID, "legacy-single-mapping", "legacy-A").Scan(
		&sourceRevision, &beforeIntent, &reviewDigest, &adapterVersion))
	require.Equal(t, preview.Items[0].SourceRevision, sourceRevision)
	var intent map[string]any
	require.NoError(t, json.Unmarshal([]byte(beforeIntent), &intent))
	require.Equal(t, "legacy-A", intent["source_id"])
	require.Equal(t, "insert", intent["write_mode"])
	require.Equal(t, preview.ReviewDigest, reviewDigest)
	require.Equal(t, preview.AdapterVersion, adapterVersion)
}

func TestSingleMappingPreviewExplainsSnapshotBoundariesAndEffectiveFreshness(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `UPDATE database_target_mappings SET write_interval_seconds = NULL, timestamp_column = ? WHERE id = ?`, "observed_at", "legacy-A")
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE database_connectors SET default_write_interval_seconds = 0 WHERE id = ?`, fixture.connectorID)
	require.NoError(t, err)

	preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
	item := preview.Items[0]
	require.Equal(t, WriteGroupMigrationPreviewStatusNeedsReview, item.Status)
	require.NotNil(t, item.CandidateGroup)
	require.Equal(t, defaultSingleMappingSnapshotIntervalSeconds, item.CandidateGroup.RowPolicy.IntervalSeconds)
	require.NotNil(t, item.CandidateGroup.Members[0].MaxAgeSeconds)
	require.Equal(t, defaultSingleMappingSnapshotIntervalSeconds, *item.CandidateGroup.Members[0].MaxAgeSeconds)
	for _, code := range []string{
		"legacy-every-sample-vs-periodic-snapshot",
		"snapshot-bucket-max-observed-tie-sample-id",
		"snapshot-late-sample-does-not-reopen",
		"snapshot-missing-or-stale-skips-row",
		"legacy-timestamp-column-ignored",
		"snapshot-interval-defaulted",
	} {
		requireFindingCode(t, item.Differences, code)
	}
}

func TestSingleMappingPreviewShowsNonPositiveLegacyIntervalFallback(t *testing.T) {
	for _, test := range []struct {
		name             string
		connectorDefault int
		expectedInterval int
		expectedSource   string
		expectedMessage  string
	}{
		{name: "connector default", connectorDefault: 30, expectedInterval: 30, expectedSource: "connector-default", expectedMessage: "connector default"},
		{name: "design default", connectorDefault: 0, expectedInterval: defaultSingleMappingSnapshotIntervalSeconds, expectedSource: "design-default", expectedMessage: "design default"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()
			db := openWorkspaceTestDB(t, ":memory:")
			defer db.Close()
			service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
			_, err := db.ExecContext(ctx, `UPDATE database_target_mappings SET write_interval_seconds = 0 WHERE id = ?`, "legacy-A")
			require.NoError(t, err)
			_, err = db.ExecContext(ctx, `UPDATE database_connectors SET default_write_interval_seconds = ? WHERE id = ?`, test.connectorDefault, fixture.connectorID)
			require.NoError(t, err)

			preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
			item := preview.Items[0]
			require.Equal(t, WriteGroupMigrationPreviewStatusNeedsReview, item.Status)
			require.NotNil(t, item.CandidateGroup)
			require.Equal(t, test.expectedInterval, item.CandidateGroup.RowPolicy.IntervalSeconds)
			require.Equal(t, test.expectedInterval, *item.CandidateGroup.Members[0].MaxAgeSeconds)
			require.Equal(t, test.expectedSource, item.BeforeIntent.IntervalSource)
			require.NotNil(t, item.BeforeIntent.WriteIntervalSeconds)
			require.Equal(t, 0, *item.BeforeIntent.WriteIntervalSeconds)
			requireFindingMessageContains(t, item.Differences, "legacy-interval-nonpositive-fallback", test.expectedMessage)
		})
	}
}

func TestReviewSingleMappingMigrationIsIdempotentAcrossReloadAndCanonicalEdit(t *testing.T) {
	ctx := t.Context()
	dbPath := filepath.Join(t.TempDir(), "reviewed-single-mapping.db")
	db := openWorkspaceTestDB(t, dbPath)
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	firstPreview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
	first, err := service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(firstPreview))
	require.NoError(t, err)
	firstGroup := first.Groups[0]
	require.NoError(t, db.Close())

	restartedDB := openWorkspaceTestDB(t, dbPath)
	defer restartedDB.Close()
	restartedWorkspace := NewService(NewSQLRepository(restartedDB))
	restartedService := NewWriteGroupService(restartedWorkspace, NewSQLWriteGroupRepository(restartedDB))
	current, err := restartedWorkspace.GetOrCreate(ctx)
	require.NoError(t, err)
	secondPreview := mustSingleMappingMigrationPreview(t, restartedService, current.ID)
	edited := cloneWriteGroup(firstGroup)
	edited.Name = "canonical operator edit"
	updated, err := restartedService.Update(ctx, firstGroup.ID, WriteGroupMutation{
		WorkspaceID:               current.ID,
		ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedGroupRevision:     firstGroup.Revision,
		ExpectedConnectorRevision: secondPreview.ConnectorRevision,
		Group:                     edited,
	})
	require.NoError(t, err)

	thirdPreview := mustSingleMappingMigrationPreview(t, restartedService, current.ID)
	second, err := restartedService.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(thirdPreview))
	require.NoError(t, err)
	require.Len(t, second.Groups, 1)
	require.Equal(t, firstGroup.ID, second.Groups[0].ID)
	require.Equal(t, updated.WorkspaceRevision, second.WorkspaceRevision)
	require.Equal(t, "canonical operator edit", second.Groups[0].Name)
	require.Equal(t, updated.Group.Revision, second.Groups[0].Revision)

	var mapCount int
	require.NoError(t, restartedDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_group_migration_maps`).Scan(&mapCount))
	require.Equal(t, 1, mapCount)
}

func TestReviewSingleMappingMigrationRejectsStaleUnconfirmedBlockedAndForeign(t *testing.T) {
	t.Run("unconfirmed", func(t *testing.T) {
		ctx := t.Context()
		db := openWorkspaceTestDB(t, ":memory:")
		defer db.Close()
		service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
		preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
		request := reviewRequestFromPreview(preview)
		request.ConfirmSnapshotConversion = false
		_, err := service.ReviewSingleMappingMigration(ctx, request)
		require.ErrorIs(t, err, ErrWriteGroupValidation)
		requireMigrationReviewCounts(ctx, t, db, 0, 0)
	})

	t.Run("stale digest", func(t *testing.T) {
		ctx := t.Context()
		db := openWorkspaceTestDB(t, ":memory:")
		defer db.Close()
		service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
		preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
		_, err := db.ExecContext(ctx, `UPDATE points SET name = ? WHERE id = ?`, "stale-source", fixture.pointID)
		require.NoError(t, err)
		_, err = service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(preview))
		require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
		requireMigrationReviewCounts(ctx, t, db, 0, 0)
	})

	t.Run("blocked", func(t *testing.T) {
		ctx := t.Context()
		db := openWorkspaceTestDB(t, ":memory:")
		defer db.Close()
		service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
		_, err := db.ExecContext(ctx, `UPDATE database_target_mappings SET group_key = ? WHERE id = ?`, "grouped", "legacy-A")
		require.NoError(t, err)
		preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
		require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, preview.Items[0].Status)
		_, err = service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(preview))
		require.ErrorIs(t, err, ErrWriteGroupValidation)
		requireMigrationReviewCounts(ctx, t, db, 0, 0)
	})

	t.Run("foreign workspace", func(t *testing.T) {
		ctx := t.Context()
		db := openWorkspaceTestDB(t, ":memory:")
		defer db.Close()
		service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
		_, err := service.ReviewSingleMappingMigration(ctx, WriteGroupMigrationReviewRequest{
			WorkspaceID:               "foreign-workspace",
			ExpectedWorkspaceRevision: "workspace-revision",
			ExpectedConnectorRevision: "connector-revision-1",
			ReviewDigest:              "digest",
			SourceIDs:                 []string{"legacy-A"},
			ConfirmSnapshotConversion: true,
		})
		require.ErrorIs(t, err, ErrWriteGroupNotFound)
		require.ErrorIs(t, err, ErrNotFound)
		requireMigrationReviewCounts(ctx, t, db, 0, 0)
		_ = fixture
	})
}

func TestReviewSingleMappingMigrationRejectsChangedSourceAfterPriorReview(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	firstPreview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
	first, err := service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(firstPreview))
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `UPDATE points SET name = ? WHERE id = ?`, "source-changed", fixture.pointID)
	require.NoError(t, err)
	current, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	secondPreview := mustSingleMappingMigrationPreview(t, service, current.ID)
	_, err = service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(secondPreview))
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	requireMigrationReviewCounts(ctx, t, db, 1, 1)
	got, err := service.Get(ctx, first.Groups[0].ID)
	require.NoError(t, err)
	require.Equal(t, first.Groups[0].Name, got.Group.Name)
}

func TestReviewSingleMappingMigrationRollsBackCanonicalGroupAndMapOnFailure(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	before, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `
		CREATE TRIGGER fail_write_group_migration_map
		AFTER INSERT ON write_group_migration_maps
		BEGIN SELECT RAISE(FAIL, 'forced migration-map failure'); END
	`)
	require.NoError(t, err)
	preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
	_, err = service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(preview))
	require.Error(t, err)
	requireMigrationReviewCounts(ctx, t, db, 0, 0)
	after, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Equal(t, before.DatabaseSetupRevision, after.DatabaseSetupRevision)
}

func mustSingleMappingMigrationPreview(t *testing.T, service *WriteGroupService, workspaceID string) *WriteGroupMigrationPreview {
	t.Helper()
	preview, err := service.PreviewSingleMappingMigration(t.Context(), workspaceID, []string{"legacy-A"})
	require.NoError(t, err)
	return preview
}

func reviewRequestFromPreview(preview *WriteGroupMigrationPreview) WriteGroupMigrationReviewRequest {
	ids := make([]string, 0, len(preview.Items))
	for _, item := range preview.Items {
		ids = append(ids, item.SourceID)
	}
	return WriteGroupMigrationReviewRequest{
		WorkspaceID:               preview.WorkspaceID,
		ExpectedWorkspaceRevision: preview.WorkspaceRevision,
		ExpectedConnectorRevision: preview.ConnectorRevision,
		ReviewDigest:              preview.ReviewDigest,
		SourceIDs:                 ids,
		ConfirmSnapshotConversion: true,
	}
}

func requireMigrationReviewCounts(ctx context.Context, t *testing.T, db *sql.DB, expectedGroups, expectedMaps int) {
	t.Helper()
	var groups, maps int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups`).Scan(&groups))
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_group_migration_maps`).Scan(&maps))
	require.Equal(t, expectedGroups, groups)
	require.Equal(t, expectedMaps, maps)
}

func requireFindingMessageContains(t *testing.T, findings []WriteGroupMigrationFinding, code, fragment string) {
	t.Helper()
	for _, finding := range findings {
		if finding.Code == code {
			require.Contains(t, finding.Message, fragment)
			return
		}
	}
	require.Failf(t, "missing finding", "expected finding code %q in %#v", code, findings)
}

func TestPrepareWriteGroupMigrationReviewAllowsEmptyInitialWorkspaceRevision(t *testing.T) {
	request, ids, err := prepareWriteGroupMigrationReviewRequest(WriteGroupMigrationReviewRequest{
		WorkspaceID:               "workspace",
		ExpectedConnectorRevision: "connector",
		ReviewDigest:              "digest",
		SourceIDs:                 []string{"source"},
		ConfirmSnapshotConversion: true,
	})
	require.NoError(t, err)
	require.Empty(t, request.ExpectedWorkspaceRevision)
	require.Equal(t, []string{"source"}, ids)
}
