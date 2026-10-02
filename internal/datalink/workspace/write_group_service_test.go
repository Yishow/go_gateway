package workspace

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteGroupServiceCreatePersistsDraftAndProjection(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)

	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	result, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID:               fixture.workspaceID,
		ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1",
		Group:                     newBasicWriteGroup(fixture),
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Group)
	require.NotEqual(t, current.DatabaseSetupRevision, result.WorkspaceRevision)
	require.NotEmpty(t, result.Group.ID)
	require.NotEmpty(t, result.Group.Revision)
	require.Equal(t, WriteGroupStatusDraft, result.Group.Status)
	require.Empty(t, result.Group.AppliedRevision)
	require.Len(t, result.Group.Members, 1)

	reloaded, err := service.Get(ctx, result.Group.ID)
	require.NoError(t, err)
	require.Equal(t, result.WorkspaceRevision, reloaded.WorkspaceRevision)
	require.Equal(t, result.Group, reloaded.Group)

	workspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Len(t, workspace.DatabaseRowGroups, 1)
	require.Equal(t, result.Group.ID, workspace.DatabaseRowGroups[0].ID)
	require.Equal(t, []string{fixture.pointID}, workspace.DatabaseRowGroups[0].MemberPointIDs)
	require.Equal(t, []DatabaseTargetRef{{PointID: fixture.pointID, RowGroupID: result.Group.ID}}, workspace.DatabaseTargetRefs)
}

func TestWriteGroupServiceAtomicGroupSaveAndCAS(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)

	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: newBasicWriteGroup(fixture),
	})
	require.NoError(t, err)

	oldWorkspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	oldGroup := cloneWriteGroup(created.Group)
	candidateA := cloneWriteGroup(oldGroup)
	candidateA.Name = "Client A"
	candidateB := cloneWriteGroup(oldGroup)
	candidateB.Name = "Client B"
	first, err := service.Update(ctx, oldGroup.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: oldWorkspace.DatabaseSetupRevision,
		ExpectedGroupRevision: oldGroup.Revision, ExpectedConnectorRevision: "connector-revision-1", Group: candidateA,
	})
	require.NoError(t, err)
	_, err = service.Update(ctx, oldGroup.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: oldWorkspace.DatabaseSetupRevision,
		ExpectedGroupRevision: oldGroup.Revision, ExpectedConnectorRevision: "connector-revision-1", Group: candidateB,
	})
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	require.ErrorIs(t, err, ErrSetupRevisionConflict)
	require.NotEqual(t, oldWorkspace.DatabaseSetupRevision, first.WorkspaceRevision)
	require.NotEqual(t, oldGroup.Revision, first.Group.Revision)

	final, err := service.Get(ctx, oldGroup.ID)
	require.NoError(t, err)
	require.Equal(t, "Client A", final.Group.Name)
	require.Equal(t, first.Group.Revision, final.Group.Revision)
	require.Equal(t, first.WorkspaceRevision, final.WorkspaceRevision)
}

func TestWriteGroupServiceRollsBackCanonicalGroupOnWorkspaceSaveFailure(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: newBasicWriteGroup(fixture),
	})
	require.NoError(t, err)
	beforeWorkspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	beforeGroup, err := service.Get(ctx, created.Group.ID)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `CREATE TRIGGER fail_write_group_workspace_save
		BEFORE UPDATE ON system_settings
		WHEN NEW.key = 'studio_v2_workspace'
		BEGIN SELECT RAISE(ABORT, 'private-storage-diagnostic'); END`)
	require.NoError(t, err)
	updatedPayload := cloneWriteGroup(created.Group)
	updatedPayload.Name = "must roll back"
	_, err = service.Update(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeWorkspace.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
		Group: updatedPayload,
	})
	require.Error(t, err)

	afterWorkspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	afterGroup, err := service.Get(ctx, created.Group.ID)
	require.NoError(t, err)
	require.Equal(t, beforeWorkspace, afterWorkspace)
	require.Equal(t, beforeGroup, afterGroup)
}

func TestWriteGroupServiceRejectsStaleWorkspaceAndConnector(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: newBasicWriteGroup(fixture),
	})
	require.NoError(t, err)

	workspaceBeforeStale, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	_, err = workspaceSvc.UpdateDatabaseSetup(ctx, workspaceBeforeStale.DatabaseSetupRevision, func(context.Context, *sql.Tx, *Record) error {
		return nil
	})
	require.NoError(t, err)
	stalePayload := cloneWriteGroup(created.Group)
	stalePayload.Name = "stale workspace"
	_, err = service.Update(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforeStale.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
		Group: stalePayload,
	})
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	require.ErrorIs(t, err, ErrSetupRevisionConflict)

	currentWorkspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE database_connectors SET identity_revision = ? WHERE id = ?`, "connector-revision-2", fixture.connectorID)
	require.NoError(t, err)
	staleConnectorPayload := cloneWriteGroup(created.Group)
	staleConnectorPayload.Name = "stale connector"
	_, err = service.Update(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: currentWorkspace.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
		Group: staleConnectorPayload,
	})
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	require.ErrorIs(t, err, ErrWriteGroupConnectorRevisionConflict)

	unchanged, err := service.Get(ctx, created.Group.ID)
	require.NoError(t, err)
	require.Equal(t, created.Group, unchanged.Group)
}

func TestWriteGroupServiceForeignAndUnknownAreSafeNotFound(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)

	foreignWorkspaceID := "foreign-workspace"
	foreignPayload := newBasicWriteGroup(fixture)
	foreignPayload.WorkspaceID = foreignWorkspaceID
	_, err = service.Create(ctx, WriteGroupMutation{
		WorkspaceID: foreignWorkspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: foreignPayload,
	})
	require.ErrorIs(t, err, ErrWriteGroupNotFound)
	require.ErrorIs(t, err, ErrNotFound)
	_, err = service.Get(ctx, "missing-group")
	require.ErrorIs(t, err, ErrWriteGroupNotFound)
	require.ErrorIs(t, err, ErrNotFound)

	_, err = service.Update(ctx, "missing-group", WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedGroupRevision: "missing-revision", ExpectedConnectorRevision: "connector-revision-1",
		Group: newBasicWriteGroup(fixture),
	})
	require.ErrorIs(t, err, ErrWriteGroupNotFound)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestWriteGroupServiceDeleteTombstonesDraftAndBlocksApplied(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: newBasicWriteGroup(fixture),
	})
	require.NoError(t, err)
	workspaceBeforeDelete, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	deleted, err := service.Delete(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforeDelete.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	require.Equal(t, WriteGroupStatusDeleted, deleted.Group.Status)
	reloaded, err := service.Get(ctx, created.Group.ID)
	require.NoError(t, err)
	require.Equal(t, WriteGroupStatusDeleted, reloaded.Group.Status)
	workspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Len(t, workspace.DatabaseRowGroups, 1)
	require.Equal(t, created.Group.ID, workspace.DatabaseRowGroups[0].ID)

	_, err = db.ExecContext(ctx, `UPDATE write_groups SET applied_revision = ? WHERE id = ?`, "applied-1", created.Group.ID)
	require.NoError(t, err)
	beforeBlocked, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	_, err = service.Delete(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeBlocked.DatabaseSetupRevision,
		ExpectedGroupRevision: deleted.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
	require.True(t, errors.Is(err, ErrWriteGroupLifecycleBlocked))
}

func TestWriteGroupServiceUpdateCannotReviveDeletedGroup(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: newBasicWriteGroup(fixture),
	})
	require.NoError(t, err)
	beforeDelete, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	deleted, err := service.Delete(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: beforeDelete.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	afterDelete, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	revive := cloneWriteGroup(deleted.Group)
	revive.Name = "should stay deleted"
	_, err = service.Update(ctx, deleted.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: afterDelete.DatabaseSetupRevision,
		ExpectedGroupRevision: deleted.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
		Group: revive,
	})
	require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
	final, err := service.Get(ctx, deleted.Group.ID)
	require.NoError(t, err)
	require.Equal(t, WriteGroupStatusDeleted, final.Group.Status)
	require.Equal(t, deleted.Group.Name, final.Group.Name)
}

func TestWriteGroupServiceProjectionPreservesUnrelatedRefs(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	legacy := DatabaseRowGroup{
		ID: "legacy-group", ConnectorID: "legacy-connector", TableSchema: "legacy", TableName: "legacy_values",
		MemberPointIDs: []string{fixture.pointID}, GroupKeyColumns: []string{"site_id"}, UniqueKeyColumns: []string{},
	}
	legacyRef := DatabaseTargetRef{PointID: fixture.pointID, RowGroupID: legacy.ID}
	seeded, err := workspaceSvc.UpdateDatabaseSetup(ctx, current.DatabaseSetupRevision, func(_ context.Context, _ *sql.Tx, record *Record) error {
		record.DatabaseRowGroups = []DatabaseRowGroup{legacy}
		record.DatabaseTargetRefs = []DatabaseTargetRef{legacyRef}
		return nil
	})
	require.NoError(t, err)
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: seeded.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: newBasicWriteGroup(fixture),
	})
	require.NoError(t, err)
	workspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Len(t, workspace.DatabaseRowGroups, 2)
	require.Equal(t, legacy, workspace.DatabaseRowGroups[0])
	require.Equal(t, legacyRef, workspace.DatabaseTargetRefs[0])
	require.Len(t, workspace.DatabaseTargetRefs, 1)
	require.Equal(t, created.Group.ID, workspace.DatabaseRowGroups[1].ID)
}

func TestWriteGroupServiceSchemaProofClearsOnColumnChange(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	created, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: newBasicWriteGroup(fixture),
	})
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE write_groups SET destination_schema_revision = ?, destination_schema_digest = ? WHERE id = ?`,
		"schema-revision-1", "schema-digest-1", created.Group.ID)
	require.NoError(t, err)
	withProof, err := service.Get(ctx, created.Group.ID)
	require.NoError(t, err)
	require.Equal(t, "schema-revision-1", withProof.Group.Destination.SchemaRevision)
	workspaceBeforeRename, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	rename := cloneWriteGroup(withProof.Group)
	rename.Name = "renamed"
	renamed, err := service.Update(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforeRename.DatabaseSetupRevision,
		ExpectedGroupRevision: withProof.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
		Group: rename,
	})
	require.NoError(t, err)
	require.Equal(t, "schema-revision-1", renamed.Group.Destination.SchemaRevision)

	workspaceBeforeColumnChange, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	columnChange := cloneWriteGroup(renamed.Group)
	columnChange.Members[0].TargetColumn = "temperature_c"
	changed, err := service.Update(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspaceBeforeColumnChange.DatabaseSetupRevision,
		ExpectedGroupRevision: renamed.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
		Group: columnChange,
	})
	require.NoError(t, err)
	require.Empty(t, changed.Group.Destination.SchemaRevision)
	require.Empty(t, changed.Group.Destination.SchemaDigest)
}
