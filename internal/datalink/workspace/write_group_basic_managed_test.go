package workspace

import (
	"context"
	"database/sql"
	"path/filepath"
	"sync"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestWriteGroupServiceEnsureBasicManagedCreatesOnceAndReplaysAfterRestart(t *testing.T) {
	ctx := t.Context()
	dbPath := filepath.Join(t.TempDir(), "basic-managed.db")
	db := openWorkspaceTestDB(t, dbPath)
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)

	basicDraft := newBasicManagedDraft(fixture)
	basicDraft.BasicManagedDeviceID = "caller-spoof"
	mutation := WriteGroupMutation{
		WorkspaceID:               fixture.workspaceID,
		ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1",
		Group:                     basicDraft,
	}
	first, err := service.EnsureBasicManaged(ctx, fixture.deviceID, mutation)
	require.NoError(t, err)
	require.NotNil(t, first)
	require.NotNil(t, first.Group)
	require.Equal(t, fixture.deviceID, first.Group.BasicManagedDeviceID)
	require.Equal(t, WriteGroupStorageStrategyManaged, first.Group.Destination.StorageStrategy)
	require.Equal(t, managedGroupTableName(first.Group.ID), first.Group.Destination.TableName)
	require.Nil(t, first.Group.Members[0].MeasurementID)
	require.Equal(t, fixture.deviceID, first.Group.Members[0].DeviceID)
	require.Equal(t, fixture.pointID, first.Group.Members[0].PointID)
	require.Equal(t, fixture.tagID, first.Group.Members[0].TagID)

	// The retry deliberately retains the stale pre-create workspace revision.
	second, err := service.EnsureBasicManaged(ctx, fixture.deviceID, mutation)
	require.NoError(t, err)
	require.Equal(t, first.Group.ID, second.Group.ID)
	require.Equal(t, fixture.deviceID, second.Group.BasicManagedDeviceID)
	require.Equal(t, first.Group.Members, second.Group.Members)
	loaded, err := service.Get(ctx, first.Group.ID)
	require.NoError(t, err)
	require.Equal(t, fixture.deviceID, loaded.Group.BasicManagedDeviceID)

	require.NoError(t, db.Close())
	restartedDB := openWorkspaceTestDB(t, dbPath)
	defer restartedDB.Close()
	restartedService := NewWriteGroupService(
		NewService(NewSQLRepository(restartedDB)),
		NewSQLWriteGroupRepository(restartedDB),
	)
	replayed, err := restartedService.EnsureBasicManaged(ctx, fixture.deviceID, mutation)
	require.NoError(t, err)
	require.Equal(t, first.Group.ID, replayed.Group.ID)
	require.Equal(t, fixture.deviceID, replayed.Group.BasicManagedDeviceID)

	listed, err := restartedService.List(ctx)
	require.NoError(t, err)
	require.Len(t, listed.Groups, 1)
	require.Equal(t, fixture.deviceID, listed.Groups[0].BasicManagedDeviceID)
}

func TestWriteGroupServiceEnsureBasicManagedConcurrentRetryConverges(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, filepath.Join(t.TempDir(), "concurrent-basic-managed.db"))
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	basicDraft := newBasicManagedDraft(fixture)
	mutation := WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: basicDraft,
	}
	type outcome struct {
		result *WriteGroupSaveResult
		err    error
	}
	results := make(chan outcome, 2)
	for range 2 {
		go func() {
			result, err := service.EnsureBasicManaged(ctx, fixture.deviceID, mutation)
			results <- outcome{result: result, err: err}
		}()
	}
	first := <-results
	second := <-results
	require.NoError(t, first.err)
	require.NoError(t, second.err)
	require.Equal(t, first.result.Group.ID, second.result.Group.ID)

	listed, err := service.List(ctx)
	require.NoError(t, err)
	require.Len(t, listed.Groups, 1)
}

func TestWriteGroupServiceEnsureBasicManagedRejectsChangedIntent(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)

	basicDraft := newBasicManagedDraft(fixture)
	mutation := WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: basicDraft,
	}
	created, err := service.EnsureBasicManaged(ctx, fixture.deviceID, mutation)
	require.NoError(t, err)

	changed := newBasicManagedDraft(fixture)
	changed.Destination.TableName = "another_values"
	mutation.Group = changed
	_, err = service.EnsureBasicManaged(ctx, fixture.deviceID, mutation)
	require.ErrorIs(t, err, ErrWriteGroupBasicManagedConflict)

	listed, err := service.List(ctx)
	require.NoError(t, err)
	require.Len(t, listed.Groups, 1)
	require.Equal(t, created.Group.ID, listed.Groups[0].ID)

	workspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	updatedDraft := cloneWriteGroup(created.Group)
	updatedDraft.Destination.TableName = "another_values"
	updated, err := service.Update(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspace.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
		Group: updatedDraft,
	})
	require.NoError(t, err)
	require.Equal(t, created.Group.ID, updated.Group.ID)
	require.Equal(t, "another_values", updated.Group.Destination.TableName)

	mutation.ExpectedWorkspaceRevision = current.DatabaseSetupRevision
	mutation.Group = changed
	replayed, err := service.EnsureBasicManaged(ctx, fixture.deviceID, mutation)
	require.NoError(t, err)
	require.Equal(t, updated.Group.ID, replayed.Group.ID)
}

func TestWriteGroupServiceEnsureBasicManagedRejectsDestinationScopeChange(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, filepath.Join(t.TempDir(), "destination-change.db"))
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)

	mutation := WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: newBasicManagedDraft(fixture),
	}
	created, err := service.EnsureBasicManaged(ctx, fixture.deviceID, mutation)
	require.NoError(t, err)

	for _, change := range []func(*WriteGroup){
		func(group *WriteGroup) { group.Destination.Database = "/tmp/another-target.db" },
		func(group *WriteGroup) { group.Destination.TableSchema = "another_schema" },
	} {
		candidate := newBasicManagedDraft(fixture)
		change(candidate)
		mutation.Group = candidate
		_, err = service.EnsureBasicManaged(ctx, fixture.deviceID, mutation)
		require.ErrorIs(t, err, ErrWriteGroupBasicManagedConflict)
	}

	listed, err := service.List(ctx)
	require.NoError(t, err)
	require.Len(t, listed.Groups, 1)
	require.Equal(t, created.Group.ID, listed.Groups[0].ID)
}

func TestWriteGroupServiceEnsureBasicManagedConcurrentAcrossServicesConverges(t *testing.T) {
	ctx := t.Context()
	dbPath := filepath.Join(t.TempDir(), "concurrent-services.db")
	seedDB := openWorkspaceTestDB(t, dbPath)
	fixture := seedWriteGroupFixture(ctx, t, seedDB)
	require.NoError(t, seedDB.Close())
	dbA := openWorkspaceTestDB(t, dbPath)
	dbB := openWorkspaceTestDB(t, dbPath)
	t.Cleanup(func() {
		require.NoError(t, dbA.Close())
		require.NoError(t, dbB.Close())
	})
	workspaceA := NewService(NewSQLRepository(dbA))
	workspaceB := NewService(NewSQLRepository(dbB))
	serviceA := NewWriteGroupService(workspaceA, NewSQLWriteGroupRepository(dbA))
	serviceB := NewWriteGroupService(workspaceB, NewSQLWriteGroupRepository(dbB))
	current, err := workspaceA.GetOrCreate(ctx)
	require.NoError(t, err)
	basicDraft := newBasicManagedDraft(fixture)
	mutation := WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: basicDraft,
	}

	type outcome struct {
		result *WriteGroupSaveResult
		err    error
	}
	results := make(chan outcome, 2)
	var wg sync.WaitGroup
	for _, service := range []*WriteGroupService{serviceA, serviceB} {
		service := service
		wg.Go(func() {
			result, err := service.EnsureBasicManaged(ctx, fixture.deviceID, mutation)
			results <- outcome{result: result, err: err}
		})
	}
	wg.Wait()
	close(results)

	var outcomes []outcome
	var outcomeErrors []error
	for result := range results {
		outcomes = append(outcomes, result)
		outcomeErrors = append(outcomeErrors, result.err)
	}
	require.Len(t, outcomes, 2)
	for _, err := range outcomeErrors {
		require.NoError(t, err)
	}
	for _, result := range outcomes {
		require.NotNil(t, result.result)
	}
	require.Equal(t, outcomes[0].result.Group.ID, outcomes[1].result.Group.ID)
	listed, err := serviceA.List(ctx)
	require.NoError(t, err)
	require.Len(t, listed.Groups, 1)
}

func TestWriteGroupServiceEnsureBasicManagedKeepsTombstone(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)

	mutation := WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: newBasicManagedDraft(fixture),
	}
	created, err := service.EnsureBasicManaged(ctx, fixture.deviceID, mutation)
	require.NoError(t, err)

	workspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	deleted, err := service.Delete(ctx, created.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspace.DatabaseSetupRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
	})
	require.NoError(t, err)
	require.Equal(t, WriteGroupStatusDeleted, deleted.Group.Status)

	_, err = service.EnsureBasicManaged(ctx, fixture.deviceID, mutation)
	require.ErrorIs(t, err, ErrWriteGroupLifecycleBlocked)
	var keyGroupID string
	require.NoError(t, db.QueryRowContext(ctx, `
		SELECT group_id FROM write_group_basic_keys
		WHERE workspace_id = ? AND device_id = ? AND canonical_role = ?
	`, fixture.workspaceID, fixture.deviceID, basicManagedCanonicalRole).Scan(&keyGroupID))
	require.Equal(t, created.Group.ID, keyGroupID)
}

func TestWriteGroupServiceEnsureBasicManagedPreservesAdvancedAndSeparatesDevices(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	workspaceSvc := NewService(NewSQLRepository(db))
	service := NewWriteGroupService(workspaceSvc, NewSQLWriteGroupRepository(db))
	current, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)

	advanced := newBasicWriteGroup(fixture)
	advanced.Name = "existing advanced"
	advanced.BasicManagedDeviceID = "caller-spoof"
	createdAdvanced, err := service.Create(ctx, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: advanced,
	})
	require.NoError(t, err)

	workspace, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	first, err := service.EnsureBasicManaged(ctx, fixture.deviceID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspace.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: newBasicManagedDraft(fixture),
	})
	require.NoError(t, err)
	require.NotEqual(t, createdAdvanced.Group.ID, first.Group.ID)
	require.Empty(t, createdAdvanced.Group.BasicManagedDeviceID)
	require.Equal(t, WriteGroupStorageStrategyCustom, createdAdvanced.Group.Destination.StorageStrategy)

	workspace, err = workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	advancedUpdate := cloneWriteGroup(createdAdvanced.Group)
	advancedUpdate.BasicManagedDeviceID = "caller-spoof-update"
	updatedAdvanced, err := service.Update(ctx, createdAdvanced.Group.ID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspace.DatabaseSetupRevision,
		ExpectedGroupRevision: createdAdvanced.Group.Revision, ExpectedConnectorRevision: "connector-revision-1",
		Group: advancedUpdate,
	})
	require.NoError(t, err)
	require.Empty(t, updatedAdvanced.Group.BasicManagedDeviceID)
	reloadedAdvanced, err := service.Get(ctx, createdAdvanced.Group.ID)
	require.NoError(t, err)
	require.Empty(t, reloadedAdvanced.Group.BasicManagedDeviceID)

	secondFixture := seedAdditionalWriteGroupDevice(ctx, t, db, fixture.workspaceID, fixture.connectorID)
	workspace, err = workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	second, err := service.EnsureBasicManaged(ctx, secondFixture.deviceID, WriteGroupMutation{
		WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspace.DatabaseSetupRevision,
		ExpectedConnectorRevision: "connector-revision-1", Group: newBasicManagedDraft(secondFixture),
	})
	require.NoError(t, err)
	require.NotEqual(t, first.Group.ID, second.Group.ID)
	require.Equal(t, "raw tags", second.Group.Name)

	listed, err := service.List(ctx)
	require.NoError(t, err)
	require.Len(t, listed.Groups, 3)
	require.Empty(t, listed.Groups[0].BasicManagedDeviceID)
	require.Equal(t, fixture.deviceID, listed.Groups[1].BasicManagedDeviceID)
	require.Equal(t, secondFixture.deviceID, listed.Groups[2].BasicManagedDeviceID)
}

func newBasicManagedDraft(fixture writeGroupFixture) *WriteGroup {
	draft := newBasicWriteGroup(fixture)
	draft.Destination.Database = ""
	draft.Destination.TableSchema = ""
	draft.Destination.TableName = ""
	return draft
}

func seedAdditionalWriteGroupDevice(ctx context.Context, t *testing.T, db *sql.DB, workspaceID, connectorID string) writeGroupFixture {
	t.Helper()
	fixture := writeGroupFixture{
		workspaceID: workspaceID, deviceID: uuid.NewString(), pointID: uuid.NewString(),
		tagID: uuid.NewString(), connectorID: connectorID,
	}
	_, err := db.ExecContext(ctx, `
		INSERT INTO devices (id, name, protocol, status, connection_config, readiness_status)
		VALUES (?, ?, ?, ?, ?, ?)
	`, fixture.deviceID, "PLC two", schema.ProtocolModbusTCP, schema.DeviceStatusDraft, `{}`, `{}`)
	require.NoError(t, err)
	_, err = NewService(NewSQLRepository(db)).AttachDevice(ctx, fixture.deviceID)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `
		INSERT INTO points (id, device_id, name, address, function, data_type, mode, enabled)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, fixture.pointID, fixture.deviceID, "pressure", "40002", "FC03", schema.DataTypeFloat32, schema.PointModeReadOnly, true)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `
		INSERT INTO tags (id, key, key_lower, display_name, data_type, status, labels)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, fixture.tagID, "pressure", "pressure", "Pressure", schema.DataTypeFloat32, schema.TagStatusActive, `{}`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `
		INSERT INTO mappings (id, point_id, tag_id, transform_pipeline, status, enabled)
		VALUES (?, ?, ?, ?, ?, ?)
	`, uuid.NewString(), fixture.pointID, fixture.tagID, `[]`, schema.MappingStatusActive, true)
	require.NoError(t, err)
	return fixture
}
