package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type rowGroupMigrationFixture struct {
	workspaceID string
	connectorID string
	rowGroupID  string
	pointIDs    []string
	tagIDs      []string
	targetIDs   []string
}

func newRowGroupMigrationFixture(ctx context.Context, t *testing.T, db *sql.DB) (*WriteGroupService, rowGroupMigrationFixture) {
	t.Helper()
	base := seedWriteGroupFixture(ctx, t, db)
	_, secondPointID := seedWriteGroupDevice(ctx, t, db)
	secondTagID := uuid.NewString()
	_, err := db.ExecContext(ctx, `INSERT INTO tags (id,key,key_lower,display_name,data_type,status,labels)
		VALUES (?,?,?,?,?,?,?)`, secondTagID, "pressure", "pressure", "Pressure", schema.DataTypeFloat32, schema.TagStatusActive, `{}`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO mappings (id,point_id,tag_id,transform_pipeline,status,enabled)
		VALUES (?,?,?,?,?,?)`, uuid.NewString(), secondPointID, secondTagID, `[]`, schema.MappingStatusActive, true)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE database_connectors SET connection_config = ? WHERE id = ?`,
		`{"database":"persisted-db","schema":"main"}`, base.connectorID)
	require.NoError(t, err)
	workspaceSvc := NewService(NewSQLRepository(db))
	workspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	_, err = workspaceSvc.UpdateDatabaseSetup(ctx, workspace.DatabaseSetupRevision, func(_ context.Context, _ *sql.Tx, record *Record) error {
		record.DatabaseConnectorID = base.connectorID
		return nil
	})
	require.NoError(t, err)
	rowGroupID := "legacy-g"
	_, err = workspaceSvc.SaveDatabaseRowGroups(ctx, base.connectorID, "main", "raw_values", []DatabaseRowGroup{{
		ID:               rowGroupID,
		MemberPointIDs:   []string{base.pointID, secondPointID},
		GroupKeyColumns:  []string{"entity_id"},
		UniqueKeyColumns: []string{"entity_id", "observed_at"},
	}})
	require.NoError(t, err)
	for _, pointID := range []string{base.pointID, secondPointID} {
		_, err = workspaceSvc.SaveDatabaseTargetReference(ctx, pointID, rowGroupID)
		require.NoError(t, err)
	}
	_, err = db.ExecContext(ctx, `
		INSERT INTO database_target_mappings
		(id,tag_id,connector_id,table_schema,table_name,column_name,write_mode,timestamp_column,group_key,write_interval_seconds,enabled)
		VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		"target-a", base.tagID, base.connectorID, "main", "raw_values", "value", "insert", "observed_at", "legacy-g:p-a", 15, true)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `
		INSERT INTO database_target_mappings
		(id,tag_id,connector_id,table_schema,table_name,column_name,write_mode,timestamp_column,group_key,write_interval_seconds,enabled)
		VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		"target-b", secondTagID, base.connectorID, "main", "raw_values", "value", "insert", "observed_at", "legacy-g:p-b", 15, true)
	require.NoError(t, err)
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	return service, rowGroupMigrationFixture{
		workspaceID: base.workspaceID,
		connectorID: base.connectorID,
		rowGroupID:  rowGroupID,
		pointIDs:    []string{base.pointID, secondPointID},
		tagIDs:      []string{base.tagID, secondTagID},
		targetIDs:   []string{"target-a", "target-b"},
	}
}

func rowGroupReviewRequest(preview *WriteGroupMigrationPreview) WriteGroupMigrationReviewRequest {
	return WriteGroupMigrationReviewRequest{
		WorkspaceID:               preview.WorkspaceID,
		ExpectedWorkspaceRevision: preview.WorkspaceRevision,
		ExpectedConnectorRevision: preview.ConnectorRevision,
		ReviewDigest:              preview.ReviewDigest,
		SourceIDs:                 []string{preview.Items[0].SourceID},
		ConfirmSnapshotConversion: true,
	}
}

func TestPreviewRowGroupMigrationPreservesSharedColumnRowsAndProvenance(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newRowGroupMigrationFixture(ctx, t, db)

	preview, err := service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
	require.NoError(t, err)
	require.Equal(t, writeGroupRowGroupMigrationAdapterVersion, preview.AdapterVersion)
	require.Len(t, preview.Items, 1)
	item := preview.Items[0]
	require.Equal(t, WriteGroupMigrationPreviewStatusNeedsReview, item.Status)
	require.NotNil(t, item.BeforeRowGroupIntent)
	require.Equal(t, fixture.rowGroupID, item.BeforeRowGroupIntent.RowGroup.ID)
	require.Len(t, item.BeforeRowGroupIntent.Members, 2)
	require.NotNil(t, item.CandidateGroup)
	require.Equal(t, []string{"entity_id"}, item.CandidateGroup.RowPolicy.GroupKeyColumns)
	require.Equal(t, []string{"entity_id", "observed_at"}, item.CandidateGroup.RowPolicy.UniqueKeyColumns)
	require.Equal(t, fixture.targetIDs, item.CandidateGroup.Migration.SourceIDs)
	require.Equal(t, fixture.rowGroupID, item.CandidateGroup.Migration.LegacyRowGroupID)
	require.Equal(t, map[string]string{"target-a": fixture.pointIDs[0], "target-b": fixture.pointIDs[1]}, item.CandidateGroup.Migration.TargetMappingPoints)
	require.Equal(t, []string{"legacy-g:p-a", "legacy-g:p-b"}, []string{item.CandidateGroup.Members[0].EntityKey, item.CandidateGroup.Members[1].EntityKey})
	for _, code := range []string{
		"legacy-arrival-order-vs-max-observed",
		"snapshot-observed-time-tie-sample-id",
		"snapshot-utc-year-one-epoch-alignment",
		"snapshot-late-sample-does-not-reopen",
		"snapshot-missing-or-stale-skips-row",
		"legacy-timestamp-intent-not-bound",
	} {
		requireFindingCode(t, item.Differences, code)
	}
}

func TestCloneWriteGroupCopiesRowGroupProvenanceAndLayout(t *testing.T) {
	original := &WriteGroup{
		Members: []WriteGroupMember{{EntityKey: "legacy-g:p-a"}},
		RowPolicy: WriteGroupRowPolicy{
			GroupKeyColumns:  []string{"entity_id"},
			UniqueKeyColumns: []string{"entity_id", "observed_at"},
		},
		Migration: WriteGroupMigration{
			SourceIDs:           []string{"target-a"},
			TargetMappingPoints: map[string]string{"target-a": "point-a"},
		},
	}
	cloned := cloneWriteGroup(original)
	cloned.Members[0].EntityKey = "changed"
	cloned.RowPolicy.GroupKeyColumns[0] = "changed"
	cloned.RowPolicy.UniqueKeyColumns[0] = "changed"
	cloned.Migration.SourceIDs[0] = "changed"
	cloned.Migration.TargetMappingPoints["target-a"] = "changed"
	require.Equal(t, "legacy-g:p-a", original.Members[0].EntityKey)
	require.Equal(t, []string{"entity_id"}, original.RowPolicy.GroupKeyColumns)
	require.Equal(t, []string{"entity_id", "observed_at"}, original.RowPolicy.UniqueKeyColumns)
	require.Equal(t, []string{"target-a"}, original.Migration.SourceIDs)
	require.Equal(t, "point-a", original.Migration.TargetMappingPoints["target-a"])
}

func TestReviewRowGroupMigrationPersistsCanonicalRowsAndReusesStableMap(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newRowGroupMigrationFixture(ctx, t, db)

	preview, err := service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
	require.NoError(t, err)
	first, err := service.ReviewRowGroupMigration(ctx, rowGroupReviewRequest(preview))
	require.NoError(t, err)
	require.Len(t, first.Groups, 1)
	group := first.Groups[0]
	require.Equal(t, fixture.rowGroupID, group.Migration.LegacyRowGroupID)
	require.Equal(t, fixture.targetIDs, group.Migration.SourceIDs)
	require.Equal(t, []string{"legacy-g:p-a", "legacy-g:p-b"}, []string{group.Members[0].EntityKey, group.Members[1].EntityKey})
	require.Empty(t, group.AppliedRevision)

	var beforeIntent string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT before_intent FROM write_group_migration_maps WHERE source_kind = ? AND source_id = ?`, "legacy-row-group", fixture.rowGroupID).Scan(&beforeIntent))
	var decoded WriteGroupRowGroupMigrationIntent
	require.NoError(t, json.Unmarshal([]byte(beforeIntent), &decoded))
	require.Equal(t, fixture.rowGroupID, decoded.RowGroup.ID)

	workspaceSvc := service.workspaceSvc
	reloaded, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Len(t, reloaded.DatabaseRowGroups, 1)
	require.Equal(t, fixture.rowGroupID, reloaded.DatabaseRowGroups[0].ID)
	require.Equal(t, []string{"entity_id"}, reloaded.DatabaseRowGroups[0].GroupKeyColumns)
	require.Equal(t, []DatabaseTargetRef{{PointID: fixture.pointIDs[0], RowGroupID: fixture.rowGroupID}, {PointID: fixture.pointIDs[1], RowGroupID: fixture.rowGroupID}}, reloaded.DatabaseTargetRefs)

	edited := cloneWriteGroup(group)
	edited.Name = "canonical row-group edit"
	updated, err := service.Update(ctx, group.ID, WriteGroupMutation{
		WorkspaceID:               fixture.workspaceID,
		ExpectedWorkspaceRevision: first.WorkspaceRevision,
		ExpectedGroupRevision:     group.Revision,
		ExpectedConnectorRevision: preview.ConnectorRevision,
		Group:                     edited,
	})
	require.NoError(t, err)
	afterEditWorkspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Len(t, afterEditWorkspace.DatabaseRowGroups, 1)
	require.Equal(t, fixture.rowGroupID, afterEditWorkspace.DatabaseRowGroups[0].ID)
	require.Equal(t, []string{"entity_id"}, afterEditWorkspace.DatabaseRowGroups[0].GroupKeyColumns)
	secondPreview, err := service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
	require.NoError(t, err)
	second, err := service.ReviewRowGroupMigration(ctx, rowGroupReviewRequest(secondPreview))
	require.NoError(t, err)
	require.Equal(t, updated.WorkspaceRevision, second.WorkspaceRevision)
	require.Equal(t, group.ID, second.Groups[0].ID)
	require.Equal(t, "canonical row-group edit", second.Groups[0].Name)
	require.Equal(t, updated.Group.Revision, second.Groups[0].Revision)
	var mapCount, groupCount int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_group_migration_maps`).Scan(&mapCount))
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups`).Scan(&groupCount))
	require.Equal(t, 1, mapCount)
	require.Equal(t, 1, groupCount)
	var targetEnabled int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT enabled FROM database_target_mappings WHERE id = ?`, fixture.targetIDs[0]).Scan(&targetEnabled))
	require.Equal(t, 1, targetEnabled)

}

func TestReviewRowGroupMigrationReusesStableMapAfterRestart(t *testing.T) {
	ctx := t.Context()
	dbPath := filepath.Join(t.TempDir(), "row-group.db")
	db := openWorkspaceTestDB(t, dbPath)
	service, fixture := newRowGroupMigrationFixture(ctx, t, db)

	preview, err := service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
	require.NoError(t, err)
	first, err := service.ReviewRowGroupMigration(ctx, rowGroupReviewRequest(preview))
	require.NoError(t, err)
	require.Len(t, first.Groups, 1)
	firstGroupID := first.Groups[0].ID
	require.NoError(t, db.Close())

	restartedDB := openWorkspaceTestDB(t, dbPath)
	defer restartedDB.Close()
	restartedWorkspace := NewService(NewSQLRepository(restartedDB))
	restarted := NewWriteGroupService(restartedWorkspace, NewSQLWriteGroupRepository(restartedDB))

	restartedPreview, err := restarted.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
	require.NoError(t, err)
	second, err := restarted.ReviewRowGroupMigration(ctx, rowGroupReviewRequest(restartedPreview))
	require.NoError(t, err)
	require.Len(t, second.Groups, 1)
	require.Equal(t, firstGroupID, second.Groups[0].ID)
	var mapCount, groupCount int
	require.NoError(t, restartedDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_group_migration_maps`).Scan(&mapCount))
	require.NoError(t, restartedDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups`).Scan(&groupCount))
	require.Equal(t, 1, mapCount)
	require.Equal(t, 1, groupCount)
}

func TestPreviewRowGroupMigrationBlocksMissingIdentityAndSameEntityCollision(t *testing.T) {
	for _, test := range []struct {
		name    string
		mutate  func(context.Context, *testing.T, *sql.DB, rowGroupMigrationFixture)
		finding string
	}{
		{name: "missing group key", mutate: func(ctx context.Context, t *testing.T, db *sql.DB, fixture rowGroupMigrationFixture) {
			_, err := db.ExecContext(ctx, `UPDATE database_target_mappings SET group_key = NULL WHERE id = ?`, fixture.targetIDs[1])
			require.NoError(t, err)
		}, finding: "row-group-group-key-missing"},
		{name: "missing unique metadata", mutate: func(ctx context.Context, t *testing.T, db *sql.DB, fixture rowGroupMigrationFixture) {
			workspaceSvc := NewService(NewSQLRepository(db))
			workspace, err := workspaceSvc.GetOrCreate(ctx)
			require.NoError(t, err)
			workspace.DatabaseRowGroups[0].UniqueKeyColumns = nil
			_, err = workspaceSvc.SaveDatabaseRowGroups(ctx, fixture.connectorID, "main", "raw_values", workspace.DatabaseRowGroups)
			require.NoError(t, err)
		}, finding: "row-group-unique-key-columns-missing"},
		{name: "same entity and column", mutate: func(ctx context.Context, t *testing.T, db *sql.DB, fixture rowGroupMigrationFixture) {
			_, err := db.ExecContext(ctx, `UPDATE database_target_mappings SET group_key = ? WHERE id = ?`, "same-entity", fixture.targetIDs[0])
			require.NoError(t, err)
			_, err = db.ExecContext(ctx, `UPDATE database_target_mappings SET group_key = ? WHERE id = ?`, "same-entity", fixture.targetIDs[1])
			require.NoError(t, err)
		}, finding: "row-group-entity-column-collision"},
		{name: "trimmed entity collision", mutate: func(ctx context.Context, t *testing.T, db *sql.DB, fixture rowGroupMigrationFixture) {
			_, err := db.ExecContext(ctx, `UPDATE database_target_mappings SET group_key = ? WHERE id = ?`, " legacy-g:p-a ", fixture.targetIDs[1])
			require.NoError(t, err)
		}, finding: "row-group-entity-column-collision"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()
			db := openWorkspaceTestDB(t, ":memory:")
			defer db.Close()
			service, fixture := newRowGroupMigrationFixture(ctx, t, db)
			test.mutate(ctx, t, db, fixture)
			preview, err := service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
			require.NoError(t, err)
			require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, preview.Items[0].Status)
			requireFindingCode(t, preview.Items[0].Issues, test.finding)
			_, err = service.ReviewRowGroupMigration(ctx, rowGroupReviewRequest(preview))
			require.ErrorIs(t, err, ErrWriteGroupValidation)
			var count int
			require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups`).Scan(&count))
			require.Zero(t, count)
		})
	}
}

func TestPreviewRowGroupMigrationBlocksScopeIntervalTimestampAndOtherConnector(t *testing.T) {
	for _, test := range []struct {
		name    string
		mutate  func(context.Context, *testing.T, *sql.DB, rowGroupMigrationFixture)
		finding string
	}{
		{name: "saved row group scope", mutate: func(ctx context.Context, t *testing.T, db *sql.DB, fixture rowGroupMigrationFixture) {
			workspaceSvc := NewService(NewSQLRepository(db))
			workspace, err := workspaceSvc.GetOrCreate(ctx)
			require.NoError(t, err)
			workspace.DatabaseRowGroups[0].TableName = "saved-a"
			_, err = workspaceSvc.SaveDatabaseRowGroups(ctx, fixture.connectorID, "main", "saved-a", workspace.DatabaseRowGroups)
			require.NoError(t, err)
		}, finding: "row-group-target-scope-mismatch"},
		{name: "mixed interval", mutate: func(ctx context.Context, t *testing.T, db *sql.DB, fixture rowGroupMigrationFixture) {
			_, err := db.ExecContext(ctx, `UPDATE database_target_mappings SET write_interval_seconds = ? WHERE id = ?`, 30, fixture.targetIDs[1])
			require.NoError(t, err)
		}, finding: "row-group-interval-mismatch"},
		{name: "timestamp mismatch", mutate: func(ctx context.Context, t *testing.T, db *sql.DB, fixture rowGroupMigrationFixture) {
			_, err := db.ExecContext(ctx, `UPDATE database_target_mappings SET timestamp_column = ? WHERE id = ?`, "received_at", fixture.targetIDs[1])
			require.NoError(t, err)
		}, finding: "row-group-timestamp-mismatch"},
		{name: "other connector output", mutate: func(ctx context.Context, t *testing.T, db *sql.DB, fixture rowGroupMigrationFixture) {
			otherConnector := uuid.NewString()
			_, err := db.ExecContext(ctx, `INSERT INTO database_connectors (id,name,kind,connection_config,identity_revision,status,enabled)
				VALUES (?,?,?,?,?,?,?)`, otherConnector, "Other target", schema.DatabaseConnectorKindSQLite, `{"database":"other","schema":"main"}`, "other-revision", schema.DatabaseConnectorStatusReady, true)
			require.NoError(t, err)
			_, err = db.ExecContext(ctx, `INSERT INTO database_target_mappings
				(id,tag_id,connector_id,table_schema,table_name,column_name,write_mode,timestamp_column,group_key,write_interval_seconds,enabled)
				VALUES (?,?,?,?,?,?,?,?,?,?,?)`, "other-target", fixture.tagIDs[0], otherConnector, "main", "raw_values", "value", "insert", "observed_at", "other-entity", 15, true)
			require.NoError(t, err)
		}, finding: "row-group-multiple-output-connectors"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()
			db := openWorkspaceTestDB(t, ":memory:")
			defer db.Close()
			service, fixture := newRowGroupMigrationFixture(ctx, t, db)
			test.mutate(ctx, t, db, fixture)
			preview, err := service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
			require.NoError(t, err)
			require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, preview.Items[0].Status)
			requireFindingCode(t, preview.Items[0].Issues, test.finding)
		})
	}
}

func TestPreviewRowGroupMigrationSourceRevisionIncludesPersistedRowGroupLayout(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newRowGroupMigrationFixture(ctx, t, db)
	first, err := service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
	require.NoError(t, err)
	workspaceSvc := NewService(NewSQLRepository(db))
	workspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	workspace.DatabaseRowGroups[0].UniqueKeyColumns = []string{"entity_id", "received_at"}
	_, err = workspaceSvc.SaveDatabaseRowGroups(ctx, fixture.connectorID, "main", "raw_values", workspace.DatabaseRowGroups)
	require.NoError(t, err)
	second, err := service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
	require.NoError(t, err)
	require.NotEqual(t, first.Items[0].SourceRevision, second.Items[0].SourceRevision)
	require.NotEqual(t, first.ReviewDigest, second.ReviewDigest)
	_, err = service.ReviewRowGroupMigration(ctx, rowGroupReviewRequest(first))
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	var groups, maps int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups`).Scan(&groups))
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_group_migration_maps`).Scan(&maps))
	require.Zero(t, groups)
	require.Zero(t, maps)
}

func TestPreviewRowGroupMigrationAllowsMatchedAppliedSourceSignatureAsProvenance(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newRowGroupMigrationFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `UPDATE mappings SET rule_candidate_id = ?, proposed_signature = ?, last_applied_signature = ? WHERE point_id = ?`,
		"candidate-1", "signature-1", "signature-1", fixture.pointIDs[0])
	require.NoError(t, err)
	preview, err := service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
	require.NoError(t, err)
	require.Equal(t, WriteGroupMigrationPreviewStatusNeedsReview, preview.Items[0].Status)
}

func TestPreviewRowGroupMigrationForeignSourceIsSafeNotFoundWithoutSelectedTarget(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newRowGroupMigrationFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `DELETE FROM database_target_mappings WHERE id = ?`, fixture.targetIDs[1])
	require.NoError(t, err)
	foreignDeviceID, foreignPointID := seedWriteGroupDevice(ctx, t, db)
	foreignTagID := uuid.NewString()
	_, err = db.ExecContext(ctx, `INSERT INTO tags (id,key,key_lower,display_name,data_type,status,labels)
		VALUES (?,?,?,?,?,?,?)`, foreignTagID, "foreign", "foreign", "Foreign", schema.DataTypeFloat32, schema.TagStatusActive, `{}`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO mappings (id,point_id,tag_id,transform_pipeline,status,enabled)
		VALUES (?,?,?,?,?,?)`, uuid.NewString(), foreignPointID, foreignTagID, `[]`, schema.MappingStatusActive, true)
	require.NoError(t, err)
	workspaceSvc := NewService(NewSQLRepository(db))
	workspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	remainingDevices := make([]string, 0, len(workspace.OrderedDeviceIDs))
	for _, deviceID := range workspace.OrderedDeviceIDs {
		if deviceID != foreignDeviceID {
			remainingDevices = append(remainingDevices, deviceID)
		}
	}
	_, err = workspaceSvc.UpdateDatabaseSetup(ctx, workspace.DatabaseSetupRevision, func(_ context.Context, _ *sql.Tx, record *Record) error {
		record.OrderedDeviceIDs = remainingDevices
		return nil
	})
	require.NoError(t, err)
	workspace, err = workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	workspace.DatabaseRowGroups[0].MemberPointIDs = []string{fixture.pointIDs[0], foreignPointID}
	_, err = workspaceSvc.SaveDatabaseRowGroups(ctx, fixture.connectorID, "main", "raw_values", workspace.DatabaseRowGroups)
	require.NoError(t, err)
	_, err = service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
	require.ErrorIs(t, err, ErrWriteGroupNotFound)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestRowGroupMigrationCannotClaimTargetOwnedBySingleMappingMigration(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newRowGroupMigrationFixture(ctx, t, db)
	_, err := db.ExecContext(ctx, `DELETE FROM database_target_mappings WHERE id = ?`, fixture.targetIDs[1])
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE database_target_mappings SET group_key = NULL WHERE id = ?`, fixture.targetIDs[0])
	require.NoError(t, err)
	singlePreview, err := service.PreviewSingleMappingMigration(ctx, fixture.workspaceID, []string{fixture.targetIDs[0]})
	require.NoError(t, err)
	require.Equal(t, WriteGroupMigrationPreviewStatusNeedsReview, singlePreview.Items[0].Status)
	_, err = service.ReviewSingleMappingMigration(ctx, WriteGroupMigrationReviewRequest{
		WorkspaceID:               singlePreview.WorkspaceID,
		ExpectedWorkspaceRevision: singlePreview.WorkspaceRevision,
		ExpectedConnectorRevision: singlePreview.ConnectorRevision,
		ReviewDigest:              singlePreview.ReviewDigest,
		SourceIDs:                 []string{fixture.targetIDs[0]},
		ConfirmSnapshotConversion: true,
	})
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO database_target_mappings
		(id,tag_id,connector_id,table_schema,table_name,column_name,write_mode,timestamp_column,group_key,write_interval_seconds,enabled)
		VALUES (?,?,?,?,?,?,?,?,?,?,?)`, fixture.targetIDs[1], fixture.tagIDs[1], fixture.connectorID, "main", "raw_values", "value", "insert", "observed_at", "legacy-g:p-b", 15, true)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE database_target_mappings SET group_key = ? WHERE id = ?`, "legacy-g:p-a", fixture.targetIDs[0])
	require.NoError(t, err)
	rowGroupPreview, err := service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
	require.NoError(t, err)
	require.Equal(t, WriteGroupMigrationPreviewStatusBlocked, rowGroupPreview.Items[0].Status)
	requireFindingCode(t, rowGroupPreview.Items[0].Issues, "migration-source-already-owned")
	var groupCount int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups`).Scan(&groupCount))
	require.Equal(t, 1, groupCount)
}

func TestReviewRowGroupMigrationRollsBackCanonicalGroupAndProjection(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newRowGroupMigrationFixture(ctx, t, db)
	before, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `CREATE TRIGGER fail_row_group_migration_map
		AFTER INSERT ON write_group_migration_maps
		BEGIN SELECT RAISE(FAIL, 'forced row-group migration-map failure'); END`)
	require.NoError(t, err)
	preview, err := service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
	require.NoError(t, err)
	_, err = service.ReviewRowGroupMigration(ctx, rowGroupReviewRequest(preview))
	require.Error(t, err)
	var groupCount, mapCount int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups`).Scan(&groupCount))
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_group_migration_maps`).Scan(&mapCount))
	require.Zero(t, groupCount)
	require.Zero(t, mapCount)
	after, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Equal(t, before, after)
	require.Equal(t, fixture.rowGroupID, after.DatabaseRowGroups[0].ID)
}
