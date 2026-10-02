package workspace

import (
	"context"
	"database/sql"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPreviewSingleMappingMigrationIsReadOnlyAndNeedsReview(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)

	beforeWorkspace, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `SELECT id FROM database_target_mappings WHERE id = ?`, "legacy-A")
	require.NoError(t, err)

	preview, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
	require.NoError(t, err)
	require.Equal(t, fixture.workspaceID, preview.WorkspaceID)
	require.Equal(t, beforeWorkspace.DatabaseSetupRevision, preview.WorkspaceRevision)
	require.Equal(t, "connector-revision-1", preview.ConnectorRevision)
	require.Equal(t, writeGroupMigrationPreviewAdapterVersion, preview.AdapterVersion)
	require.NotEmpty(t, preview.ReviewDigest)
	require.Len(t, preview.Items, 1)
	item := preview.Items[0]
	require.Equal(t, "legacy-A", item.SourceID)
	require.Equal(t, WriteGroupMigrationPreviewStatusNeedsReview, item.Status)
	require.NotEmpty(t, item.SourceRevision)
	require.NotEmpty(t, item.Differences)
	require.Empty(t, item.Issues)
	require.NotNil(t, item.BeforeIntent)
	require.Equal(t, "insert", item.BeforeIntent.WriteMode)
	require.Equal(t, "observed_at", *item.BeforeIntent.TimestampColumn)
	require.Equal(t, 15, *item.BeforeIntent.WriteIntervalSeconds)
	require.NotNil(t, item.CandidateGroup)
	require.Equal(t, fixture.workspaceID, item.CandidateGroup.WorkspaceID)
	require.Equal(t, WriteGroupStatusDraft, item.CandidateGroup.Status)
	require.Equal(t, fixture.pointID, item.CandidateGroup.Members[0].PointID)

	second, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
	require.NoError(t, err)
	require.Equal(t, preview, second)

	afterWorkspace, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Equal(t, beforeWorkspace, afterWorkspace)
	var groups int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups`).Scan(&groups))
	require.Zero(t, groups)
	var mappingCount int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM database_target_mappings WHERE id = ?`, "legacy-A").Scan(&mappingCount))
	require.Equal(t, 1, mappingCount)
}

func TestPreviewSingleMappingMigrationDigestChangesWithSource(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)

	first, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE points SET name = ? WHERE id = ?`, "temperature-renamed", fixture.pointID)
	require.NoError(t, err)
	second, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
	require.NoError(t, err)
	require.NotEqual(t, first.ReviewDigest, second.ReviewDigest)
	require.NotEqual(t, first.Items[0].SourceRevision, second.Items[0].SourceRevision)
	require.Equal(t, first.Items[0].CandidateGroup.ID, second.Items[0].CandidateGroup.ID)
}

func TestPreviewSingleMappingMigrationBlocksAmbiguousSourcesAndTargets(t *testing.T) {
	t.Run("multiple enabled source mappings", func(t *testing.T) {
		ctx := t.Context()
		db := openWorkspaceTestDB(t, ":memory:")
		defer db.Close()
		service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
		_, secondPointID := seedWriteGroupDevice(ctx, t, db)
		_, err := db.ExecContext(ctx, `INSERT INTO mappings (id,point_id,tag_id,transform_pipeline,status,enabled)
			VALUES (?,?,?,?,?,?)`, uuid.NewString(), secondPointID, fixture.tagID, `[]`, schema.MappingStatusActive, true)
		require.NoError(t, err)

		preview, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
		require.NoError(t, err)
		require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, preview.Items[0].Status)
		requireFindingCode(t, preview.Items[0].Issues, "multiple-enabled-source-mappings")
	})

	t.Run("multiple target mappings for tag", func(t *testing.T) {
		ctx := t.Context()
		db := openWorkspaceTestDB(t, ":memory:")
		defer db.Close()
		service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
		_, err := db.ExecContext(ctx, `INSERT INTO database_target_mappings
			(id,tag_id,connector_id,table_schema,table_name,column_name,write_mode,enabled)
			VALUES (?,?,?,?,?,?,?,?)`, "legacy-B", fixture.tagID, "another-connector", "main", "other_values", "temperature", "insert", true)
		require.NoError(t, err)

		preview, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
		require.NoError(t, err)
		require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, preview.Items[0].Status)
		requireFindingCode(t, preview.Items[0].Issues, "multiple-target-mappings")
	})
}

func TestPreviewSingleMappingMigrationReturnsSafeNotFoundForUnknownOrForeignSource(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)

	_, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"missing-source"})
	require.ErrorIs(t, err, ErrWriteGroupNotFound)
	require.ErrorIs(t, err, ErrNotFound)
	_, err = service.PreviewSingleMappingMigration(ctx, uuid.NewString(), []string{"legacy-A"})
	require.ErrorIs(t, err, ErrWriteGroupNotFound)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestPreviewSingleMappingMigrationBlocksDisabledMapping(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `UPDATE mappings SET enabled = 0 WHERE point_id = ? AND tag_id = ?`, fixture.pointID, fixture.tagID)
	require.NoError(t, err)

	preview, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
	require.NoError(t, err)
	require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, preview.Items[0].Status)
	requireFindingCode(t, preview.Items[0].Issues, "source-mapping-missing-or-disabled")
}

func TestPreviewSingleMappingMigrationPreservesPersistedIntentWithoutPersistentID(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `UPDATE database_target_mappings SET timestamp_column = ? WHERE id = ?`, " observed_at ", "legacy-A")
	require.NoError(t, err)

	preview, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
	require.NoError(t, err)
	item := preview.Items[0]
	require.Equal(t, WriteGroupMigrationPreviewStatusNeedsReview, item.Status)
	require.NotNil(t, item.BeforeIntent)
	require.NotNil(t, item.BeforeIntent.TimestampColumn)
	require.Equal(t, " observed_at ", *item.BeforeIntent.TimestampColumn)
	require.NotNil(t, item.CandidateGroup)
	require.Empty(t, item.CandidateGroup.ID)
}

func TestPreviewSingleMappingMigrationCountsOnlyEnabledTargetsAndAllEnabledSources(t *testing.T) {
	t.Run("disabled target does not create ambiguity", func(t *testing.T) {
		ctx := t.Context()
		db := openWorkspaceTestDB(t, ":memory:")
		defer db.Close()
		service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
		_, err := db.ExecContext(ctx, `INSERT INTO database_target_mappings
			(id,tag_id,connector_id,table_schema,table_name,column_name,write_mode,enabled)
			VALUES (?,?,?,?,?,?,?,?)`, "legacy-disabled", fixture.tagID, "other-connector", "main", "raw_values", "temperature", "insert", false)
		require.NoError(t, err)

		preview, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
		require.NoError(t, err)
		require.Equal(t, WriteGroupMigrationPreviewStatusNeedsReview, preview.Items[0].Status)
	})

	t.Run("enabled nonactive source is ambiguous", func(t *testing.T) {
		ctx := t.Context()
		db := openWorkspaceTestDB(t, ":memory:")
		defer db.Close()
		service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
		_, secondPointID := seedWriteGroupDevice(ctx, t, db)
		_, err := db.ExecContext(ctx, `INSERT INTO mappings (id,point_id,tag_id,transform_pipeline,status,enabled)
			VALUES (?,?,?,?,?,?)`, uuid.NewString(), secondPointID, fixture.tagID, `[]`, schema.MappingStatusDraft, true)
		require.NoError(t, err)

		preview, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
		require.NoError(t, err)
		require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, preview.Items[0].Status)
		requireFindingCode(t, preview.Items[0].Issues, "multiple-enabled-source-mappings")
	})
}

func newSingleMappingMigrationPreviewFixture(ctx context.Context, t *testing.T, db *sql.DB) (*WriteGroupService, writeGroupFixture) {
	t.Helper()
	fixture := seedWriteGroupFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `UPDATE database_connectors SET connection_config = ? WHERE id = ?`,
		`{"database":"persisted-db","schema":"main"}`, fixture.connectorID)
	require.NoError(t, err)
	workspaceSvc := NewService(NewSQLRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	_, err = workspaceSvc.UpdateDatabaseSetup(ctx, current.DatabaseSetupRevision, func(_ context.Context, _ *sql.Tx, record *Record) error {
		record.DatabaseConnectorID = fixture.connectorID
		return nil
	})
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO database_target_mappings
		(id,tag_id,connector_id,table_schema,table_name,column_name,write_mode,timestamp_column,write_interval_seconds,enabled)
		VALUES (?,?,?,?,?,?,?,?,?,?)`, "legacy-A", fixture.tagID, fixture.connectorID, "main", "raw_values", "temperature", "insert", "observed_at", 15, true)
	require.NoError(t, err)
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	return service, fixture
}

func requireFindingCode(t *testing.T, findings []WriteGroupMigrationFinding, code string) {
	t.Helper()
	for _, finding := range findings {
		if finding.Code == code {
			return
		}
	}
	require.Failf(t, "missing finding", "expected finding code %q in %#v", code, findings)
}
