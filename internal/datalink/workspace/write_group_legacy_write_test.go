package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestPreflightLegacyTargetWriteRejectsMigratedMapping(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
	_, err := service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(preview))
	require.NoError(t, err)

	err = service.PreflightLegacyTargetWrite(ctx, "legacy-A", fixture.tagID, fixture.connectorID)
	require.ErrorIs(t, err, ErrWriteGroupLegacyWriteConflict)
}

func TestCheckLegacyTargetWriteInTxRejectsMigratedMapping(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
	_, err := service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(preview))
	require.NoError(t, err)

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()
	err = service.CheckLegacyTargetWriteInTx(ctx, tx, "legacy-A", fixture.tagID, fixture.connectorID)
	require.ErrorIs(t, err, ErrWriteGroupLegacyWriteConflict)
}

func TestPreflightLegacyRowGroupReplacementRejectsOwnedGroup(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newRowGroupMigrationFixture(ctx, t, db)
	preview, err := service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
	require.NoError(t, err)
	_, err = service.ReviewRowGroupMigration(ctx, rowGroupReviewRequest(preview))
	require.NoError(t, err)

	err = service.PreflightLegacyRowGroupReplacement(ctx, fixture.connectorID, []DatabaseRowGroup{})
	require.ErrorIs(t, err, ErrWriteGroupLegacyWriteConflict)
}

// The first save of a destination has no connector ID yet. An empty row-group
// replacement then changes nothing and must not be refused, otherwise a fresh
// workspace can never save its destination; a non-empty one still needs a connector.
func TestPreflightLegacyRowGroupReplacementAllowsEmptySetForANewDestination(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, _ := newRowGroupMigrationFixture(ctx, t, db)

	require.NoError(t, service.PreflightLegacyRowGroupReplacement(ctx, "", []DatabaseRowGroup{}))
	require.ErrorIs(t, service.PreflightLegacyRowGroupReplacement(ctx, "", []DatabaseRowGroup{{ID: "rg"}}), ErrWriteGroupValidation)

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()
	record, err := readLegacyWriteWorkspace(ctx, tx)
	require.NoError(t, err)
	require.NoError(t, service.CheckLegacyRowGroupReplacementInTx(ctx, tx, record, " ", []DatabaseRowGroup{}))
}

func TestLegacyTargetWriteUsesDurableBeforeIntentAfterSourceRemovalAndCanonicalEdit(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
	review, err := service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(preview))
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `DELETE FROM database_target_mappings WHERE id = ?`, "legacy-A")
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO database_connectors
		(id, name, kind, connection_config, identity_revision, status, enabled)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		"connector-2", "SQLite target 2", "sqlite", `{"database":"persisted-db","schema":"main"}`,
		"connector-revision-2", "ready", true)
	require.NoError(t, err)
	groupResult, err := service.Get(ctx, review.Groups[0].ID)
	require.NoError(t, err)
	group := groupResult.Group
	group.Name = "canonical-edit"
	group.Destination.ConnectorID = "connector-2"
	group.Destination.ConnectorRevision = "connector-revision-2"
	_, err = service.Update(ctx, group.ID, WriteGroupMutation{
		WorkspaceID:               group.WorkspaceID,
		ExpectedWorkspaceRevision: groupResult.WorkspaceRevision,
		ExpectedGroupRevision:     group.Revision,
		ExpectedConnectorRevision: "connector-revision-2",
		Group:                     group,
	})
	require.NoError(t, err)

	err = service.PreflightLegacyTargetWrite(ctx, "", fixture.tagID, fixture.connectorID)
	require.ErrorIs(t, err, ErrWriteGroupLegacyWriteConflict)
	err = service.PreflightLegacyTargetWrite(ctx, "legacy-A", "", "")
	require.ErrorIs(t, err, ErrWriteGroupLegacyWriteConflict)
}

func TestLegacyTargetWriteProtectsDisabledAndDeletedCanonicalGroups(t *testing.T) {
	for _, status := range []WriteGroupStatus{WriteGroupStatusDisabled, WriteGroupStatusDeleted} {
		t.Run(string(status), func(t *testing.T) {
			ctx := t.Context()
			db := openWorkspaceTestDB(t, ":memory:")
			defer db.Close()
			service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
			preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
			review, err := service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(preview))
			require.NoError(t, err)
			_, err = db.ExecContext(ctx, `UPDATE write_groups SET status = ? WHERE id = ?`, status, review.Groups[0].ID)
			require.NoError(t, err)

			err = service.PreflightLegacyTargetWrite(ctx, "legacy-A", fixture.tagID, fixture.connectorID)
			require.ErrorIs(t, err, ErrWriteGroupLegacyWriteConflict)
		})
	}
}

func TestLegacyRowGroupWriteUsesDurableBeforeIntentAfterRawTargetsAndDestinationChange(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newRowGroupMigrationFixture(ctx, t, db)
	preview, err := service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
	require.NoError(t, err)
	review, err := service.ReviewRowGroupMigration(ctx, rowGroupReviewRequest(preview))
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `DELETE FROM database_target_mappings WHERE id IN (?, ?)`, fixture.targetIDs[0], fixture.targetIDs[1])
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO database_connectors
		(id, name, kind, connection_config, identity_revision, status, enabled)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		"connector-2", "SQLite row target 2", "sqlite", `{"database":"persisted-db","schema":"main"}`,
		"connector-revision-2", "ready", true)
	require.NoError(t, err)
	groupResult, err := service.Get(ctx, review.Groups[0].ID)
	require.NoError(t, err)
	group := groupResult.Group
	group.Destination.ConnectorID = "connector-2"
	group.Destination.ConnectorRevision = "connector-revision-2"
	_, err = service.Update(ctx, group.ID, WriteGroupMutation{
		WorkspaceID:               group.WorkspaceID,
		ExpectedWorkspaceRevision: groupResult.WorkspaceRevision,
		ExpectedGroupRevision:     group.Revision,
		ExpectedConnectorRevision: "connector-revision-2",
		Group:                     group,
	})
	require.NoError(t, err)

	err = service.PreflightLegacyTargetWrite(ctx, "", fixture.tagIDs[0], fixture.connectorID)
	require.ErrorIs(t, err, ErrWriteGroupLegacyWriteConflict)
	err = service.PreflightLegacyTargetWrite(ctx, fixture.targetIDs[0], "", "")
	require.ErrorIs(t, err, ErrWriteGroupLegacyWriteConflict)
}

func TestLegacyTargetWriteFailsClosedForMissingOrMalformedMigrationMap(t *testing.T) {
	for _, testCase := range []struct {
		name   string
		mutate string
	}{
		{name: "malformed before intent", mutate: `UPDATE write_group_migration_maps SET before_intent = ?`},
		{name: "missing migration map", mutate: `DELETE FROM write_group_migration_maps`},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := t.Context()
			db := openWorkspaceTestDB(t, ":memory:")
			defer db.Close()
			service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
			preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
			_, err := service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(preview))
			require.NoError(t, err)
			if strings.HasPrefix(testCase.mutate, "UPDATE") {
				_, err = db.ExecContext(ctx, testCase.mutate+` WHERE source_kind = ? AND source_id = ?`, `{}`, writeGroupMigrationSourceKindLegacySingleMapping, "legacy-A")
			} else {
				_, err = db.ExecContext(ctx, testCase.mutate+` WHERE source_kind = ? AND source_id = ?`, writeGroupMigrationSourceKindLegacySingleMapping, "legacy-A")
			}
			require.NoError(t, err)

			err = service.PreflightLegacyTargetWrite(ctx, "legacy-A", "", "")
			require.ErrorIs(t, err, ErrWriteGroupLegacyWriteConflict)
			require.ErrorIs(t, err, ErrWriteGroupValidation)
		})
	}
}

func TestLegacyTargetWriteFailsClosedForUnknownMigrationKind(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
	review, err := service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(preview))
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE write_groups SET migration = ? WHERE id = ?`,
		`{"source_kind":"future-kind","source_ids":["legacy-A"]}`, review.Groups[0].ID)
	require.NoError(t, err)

	err = service.PreflightLegacyTargetWrite(ctx, "legacy-A", "", "")
	require.ErrorIs(t, err, ErrWriteGroupLegacyWriteConflict)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
}

func TestLegacyTargetWriteFailsClosedForDuplicateRowBeforeIntentSourceIDs(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newRowGroupMigrationFixture(ctx, t, db)
	preview, err := service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
	require.NoError(t, err)
	_, err = service.ReviewRowGroupMigration(ctx, rowGroupReviewRequest(preview))
	require.NoError(t, err)
	var rawIntent string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT before_intent FROM write_group_migration_maps WHERE source_kind = ? AND source_id = ?`, writeGroupMigrationSourceKindLegacyRowGroup, fixture.rowGroupID).Scan(&rawIntent))
	var intent WriteGroupRowGroupMigrationIntent
	require.NoError(t, json.Unmarshal([]byte(rawIntent), &intent))
	require.Len(t, intent.Members, 2)
	intent.Members[1].SourceID = intent.Members[0].SourceID
	updatedIntent, err := json.Marshal(intent)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE write_group_migration_maps SET before_intent = ? WHERE source_kind = ? AND source_id = ?`, string(updatedIntent), writeGroupMigrationSourceKindLegacyRowGroup, fixture.rowGroupID)
	require.NoError(t, err)

	err = service.PreflightLegacyTargetWrite(ctx, fixture.targetIDs[1], "", "")
	require.ErrorIs(t, err, ErrWriteGroupLegacyWriteConflict)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
}

func TestLegacyTargetWriteOperationalReadErrorsRemainUnavailableAndDoNotPersist(t *testing.T) {
	for _, testCase := range []struct {
		name  string
		table string
	}{
		{name: "list groups or members", table: "write_group_members"},
		{name: "read migration map", table: "write_group_migration_maps"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := t.Context()
			db := openWorkspaceTestDB(t, ":memory:")
			defer db.Close()
			service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
			preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
			_, err := service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(preview))
			require.NoError(t, err)
			_, err = db.ExecContext(ctx, "DROP TABLE "+testCase.table)
			require.NoError(t, err)

			persistCalled := false
			err = service.RunLegacyTargetWrite(ctx, "legacy-A", &schema.DatabaseTargetMapping{
				ID:          "legacy-A",
				TagID:       fixture.tagID,
				ConnectorID: fixture.connectorID,
			}, func(context.Context, *sql.Tx) error {
				persistCalled = true
				return nil
			})
			require.ErrorIs(t, err, ErrWriteGroupServiceUnavailable)
			require.NotErrorIs(t, err, ErrWriteGroupLegacyWriteConflict)
			require.NotErrorIs(t, err, ErrWriteGroupValidation)
			require.Contains(t, err.Error(), testCase.table)
			require.False(t, persistCalled)
		})
	}
}

func TestRunLegacyTargetWriteWithoutWorkspaceDoesNotCreateOne(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	_, err := db.ExecContext(ctx, `DELETE FROM system_settings WHERE key = ?`, storageKey)
	require.NoError(t, err)

	candidate := &schema.DatabaseTargetMapping{ID: "unowned", TagID: "tag-unowned", ConnectorID: "connector-unowned"}
	require.NoError(t, service.RunLegacyTargetWrite(ctx, candidate.ID, candidate, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO database_target_mappings
			(id, tag_id, connector_id, table_schema, table_name, column_name, write_mode, enabled)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, candidate.ID, candidate.TagID, candidate.ConnectorID, "main", "raw_values", "value", "insert", true)
		return err
	}))
	var workspaceCount int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM system_settings WHERE key = ?`, storageKey).Scan(&workspaceCount))
	require.Zero(t, workspaceCount)
}

func TestLegacyRowGroupReplacementAllowsEqualOwnedProjectionAndRejectsChange(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newRowGroupMigrationFixture(ctx, t, db)
	preview, err := service.PreviewRowGroupMigration(ctx, fixture.workspaceID, []string{fixture.rowGroupID})
	require.NoError(t, err)
	review, err := service.ReviewRowGroupMigration(ctx, rowGroupReviewRequest(preview))
	require.NoError(t, err)

	current, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.NoError(t, service.PreflightLegacyRowGroupReplacement(ctx, fixture.connectorID, current.DatabaseRowGroups))
	changed := append([]DatabaseRowGroup(nil), current.DatabaseRowGroups...)
	changed[0].TableName = "other_values"
	err = service.PreflightLegacyRowGroupReplacement(ctx, fixture.connectorID, changed)
	require.ErrorIs(t, err, ErrWriteGroupLegacyWriteConflict)

	// The review result is intentionally retained in the test so the equal
	// projection assertion is tied to the migrated row-group owner.
	require.Equal(t, fixture.rowGroupID, review.Groups[0].Migration.LegacyRowGroupID)
}

func TestRunLegacyTargetWriteRechecksOwnershipAfterPreflight(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)

	require.NoError(t, service.PreflightLegacyTargetWrite(ctx, "legacy-A", fixture.tagID, fixture.connectorID))
	preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
	_, err := service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(preview))
	require.NoError(t, err)
	err = service.RunLegacyTargetWrite(ctx, "legacy-A", &schema.DatabaseTargetMapping{
		ID:          "legacy-A",
		TagID:       fixture.tagID,
		ConnectorID: fixture.connectorID,
	}, func(context.Context, *sql.Tx) error { return nil })
	require.ErrorIs(t, err, ErrWriteGroupLegacyWriteConflict)
}

func TestRunLegacyTargetWriteRollsBackFailedPersistence(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	var before string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT column_name FROM database_target_mappings WHERE id = ?`, "legacy-A").Scan(&before))
	wantErr := errors.New("forced legacy write failure")
	err := service.RunLegacyTargetWrite(ctx, "legacy-A", &schema.DatabaseTargetMapping{
		ID:          "legacy-A",
		TagID:       fixture.tagID,
		ConnectorID: fixture.connectorID,
	}, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE database_target_mappings SET column_name = ? WHERE id = ?`, "rolled-back", "legacy-A")
		if err != nil {
			return err
		}
		return wantErr
	})
	require.ErrorIs(t, err, wantErr)
	var after string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT column_name FROM database_target_mappings WHERE id = ?`, "legacy-A").Scan(&after))
	require.Equal(t, before, after)
}

func TestPreflightConnectorDeleteRejectsConnectorOwnedByCanonicalGroup(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	service, fixture := newSingleMappingMigrationPreviewFixture(ctx, t, db)
	preview := mustSingleMappingMigrationPreview(t, service, fixture.workspaceID)
	_, err := service.ReviewSingleMappingMigration(ctx, reviewRequestFromPreview(preview))
	require.NoError(t, err)

	require.ErrorIs(t, service.PreflightConnectorDelete(ctx, fixture.connectorID), ErrWriteGroupLegacyWriteConflict)
	require.NoError(t, service.PreflightConnectorDelete(ctx, "connector-without-groups"), "an unrelated connector may be deleted")
	require.ErrorIs(t, service.PreflightConnectorDelete(ctx, " "), ErrWriteGroupValidation)
}
