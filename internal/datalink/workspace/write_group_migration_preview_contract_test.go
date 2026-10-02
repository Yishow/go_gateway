package workspace

import (
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestPreviewSingleMappingMigrationWireFindingsAreArrays(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	for _, disabled := range []bool{false, true} {
		if disabled {
			_, err := db.ExecContext(ctx, `UPDATE mappings SET enabled = 0 WHERE tag_id = ?`, fixture.tagID)
			require.NoError(t, err)
		}
		preview, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
		require.NoError(t, err)
		encoded, err := json.Marshal(preview.Items[0])
		require.NoError(t, err)
		var wire map[string]any
		require.NoError(t, json.Unmarshal(encoded, &wire))
		_, differencesAreArray := wire["differences"].([]any)
		_, issuesAreArray := wire["issues"].([]any)
		require.True(t, differencesAreArray, string(encoded))
		require.True(t, issuesAreArray, string(encoded))
	}
}

func TestPreviewSingleMappingMigrationDefaultIntervalChangesSourceRevision(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `UPDATE database_target_mappings SET write_interval_seconds = NULL WHERE id = 'legacy-A'`)
	require.NoError(t, err)
	first, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE database_connectors SET default_write_interval_seconds = 30 WHERE id = ?`, fixture.connectorID)
	require.NoError(t, err)
	second, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
	require.NoError(t, err)
	require.NotEqual(t, first.Items[0].SourceRevision, second.Items[0].SourceRevision)
	require.NotEqual(t, first.ReviewDigest, second.ReviewDigest)
	require.NotNil(t, second.Items[0].CandidateGroup)
	require.Equal(t, 30, second.Items[0].CandidateGroup.RowPolicy.IntervalSeconds)
}

func TestPreviewSingleMappingMigrationChecksEveryPotentialColumnCollision(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `INSERT INTO database_target_mappings
		(id,tag_id,connector_id,table_schema,table_name,column_name,write_mode,group_key,enabled)
		VALUES ('legacy-B','tag-B',?,'main','raw_values','other','insert','legacy-group',1),
		('legacy-C','tag-C',?,'main','raw_values','temperature','insert',NULL,1)`, fixture.connectorID, fixture.connectorID)
	require.NoError(t, err)
	preview, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
	require.NoError(t, err)
	require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, preview.Items[0].Status)
	requireFindingCode(t, preview.Items[0].Issues, "shared-target-column")
	require.Nil(t, preview.Items[0].CandidateGroup)
}

func TestPreviewSingleMappingMigrationRejectsUnsupportedLegacyWriteMode(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `UPDATE database_target_mappings SET write_mode = 'unknown' WHERE id = 'legacy-A'`)
	require.NoError(t, err)
	preview, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
	require.NoError(t, err)
	require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, preview.Items[0].Status)
	requireFindingCode(t, preview.Items[0].Issues, "write-mode-unsupported")
	require.Equal(t, "unknown", preview.Items[0].BeforeIntent.WriteMode)
	require.Nil(t, preview.Items[0].CandidateGroup)
}

func TestPreviewSingleMappingMigrationPreservesMissingTagBindingForRepair(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `DELETE FROM tags WHERE id = ?`, fixture.tagID)
	require.NoError(t, err)
	preview, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
	require.NoError(t, err)
	item := preview.Items[0]
	require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, item.Status)
	requireFindingCode(t, item.Issues, "source-tag-missing-or-disabled")
	require.Equal(t, fixture.tagID, item.BeforeIntent.TagID)
	require.Nil(t, item.CandidateGroup)
}

func TestPreviewSingleMappingMigrationCannotMoveLegacySchemaScope(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `UPDATE database_connectors SET kind = 'postgres', connection_config = ? WHERE id = ?`,
		`{"database":"line-a","schema":"authorized"}`, fixture.connectorID)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE database_target_mappings SET table_schema = 'legacy' WHERE id = 'legacy-A'`)
	require.NoError(t, err)
	preview, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
	require.NoError(t, err)
	item := preview.Items[0]
	require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, item.Status)
	requireFindingCode(t, item.Issues, "connector-scope-mismatch")
	require.Equal(t, "legacy", item.BeforeIntent.TableSchema)
	require.Nil(t, item.CandidateGroup)
}

func TestPreviewSingleMappingMigrationForeignSourcesStayNotFoundWhenAmbiguous(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `DELETE FROM mappings WHERE tag_id = ?`, fixture.tagID)
	require.NoError(t, err)
	for _, id := range []string{"foreign-source-A", "foreign-source-B"} {
		pointID := id + "-point"
		_, err := db.ExecContext(ctx, `INSERT INTO devices (id,name,protocol,status,connection_config,readiness_status)
			VALUES (?,?,?,'draft','{}','{}')`, id, id, schema.ProtocolModbusTCP)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, `INSERT INTO points (id,device_id,name,address,function,data_type,mode,enabled)
			VALUES (?,?,?,'40001','FC03',?,'read',1)`, pointID, id, pointID, schema.DataTypeFloat32)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, `INSERT INTO mappings (id,point_id,tag_id,transform_pipeline,status,enabled)
			VALUES (?,?,?,?,?,?)`, id, pointID, fixture.tagID, `[]`, schema.MappingStatusActive, true)
		require.NoError(t, err)
	}
	_, err = service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{"legacy-A"})
	require.ErrorIs(t, err, ErrWriteGroupNotFound)
}
