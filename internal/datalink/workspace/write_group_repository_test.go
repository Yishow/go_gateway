package workspace

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBasicGroupWithoutMeasurement(t *testing.T) {
	ctx := t.Context()
	dbPath := filepath.Join(t.TempDir(), "write-groups.db")
	db := openWorkspaceTestDB(t, dbPath)
	fixture := seedWriteGroupFixture(ctx, t, db)
	repo := NewSQLWriteGroupRepository(db)

	group := &WriteGroup{
		ID:              "client-supplied-id",
		WorkspaceID:     fixture.workspaceID,
		Revision:        "client-supplied-revision",
		AppliedRevision: "client-supplied-applied-revision",
		Name:            "raw tags",
		Status:          WriteGroupStatusReady,
		Members: []WriteGroupMember{{
			DeviceID:     fixture.deviceID,
			PointID:      fixture.pointID,
			TagID:        fixture.tagID,
			TargetColumn: "temperature",
			Required:     true,
		}},
		Destination: WriteGroupDestination{
			ConnectorID:       fixture.connectorID,
			ConnectorRevision: "connector-revision-1",
			Database:          "line-a",
			TableSchema:       "main",
			TableName:         "raw_values",
			StorageStrategy:   WriteGroupStorageStrategyCustom,
		},
	}

	require.NoError(t, repo.Create(ctx, group))
	require.NotEqual(t, "client-supplied-id", group.ID)
	require.NotEqual(t, "client-supplied-revision", group.Revision)
	require.NotEmpty(t, group.ID)
	require.NotEmpty(t, group.Revision)
	require.Empty(t, group.AppliedRevision)
	require.Equal(t, WriteGroupStatusDraft, group.Status)
	require.Nil(t, group.Members[0].MeasurementID)

	got, err := repo.Get(ctx, fixture.workspaceID, group.ID)
	require.NoError(t, err)
	require.Equal(t, group, got)
	require.Equal(t, fixture.deviceID, got.Members[0].DeviceID)
	require.Equal(t, fixture.pointID, got.Members[0].PointID)
	require.Equal(t, fixture.tagID, got.Members[0].TagID)
	require.Equal(t, "persisted-db", got.Destination.Database)
	require.Equal(t, "persisted-schema", got.Destination.TableSchema)
	require.NotEmpty(t, got.Members[0].MappingRevision)
	require.NotEmpty(t, got.Members[0].SourceRevision)
	sourceRevision := got.Members[0].SourceRevision

	got.Members[0].DeviceID = "mutated-after-read"
	got.Destination.TableName = "mutated-after-read"
	listed, err := repo.List(ctx, fixture.workspaceID)
	require.NoError(t, err)
	require.Len(t, listed, 1)
	require.Equal(t, fixture.deviceID, listed[0].Members[0].DeviceID)
	require.Equal(t, "raw_values", listed[0].Destination.TableName)

	_, err = db.ExecContext(ctx, `UPDATE devices SET status = ? WHERE id = ?`, schema.DeviceStatusActive, fixture.deviceID)
	require.NoError(t, err)
	stableGroup := newBasicWriteGroup(fixture)
	require.NoError(t, repo.Create(ctx, stableGroup))
	stableReload, err := repo.Get(ctx, fixture.workspaceID, stableGroup.ID)
	require.NoError(t, err)
	require.Equal(t, sourceRevision, stableReload.Members[0].SourceRevision)

	require.NoError(t, db.Close())
	restartedDB := openWorkspaceTestDB(t, dbPath)
	defer restartedDB.Close()
	reloaded, err := NewSQLWriteGroupRepository(restartedDB).Get(ctx, fixture.workspaceID, group.ID)
	require.NoError(t, err)
	require.Equal(t, group, reloaded)
	_, err = NewSQLWriteGroupRepository(restartedDB).Get(ctx, fixture.workspaceID, "missing-group")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestWriteGroupCreateAcceptsMatchingOptionalMeasurement(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	measurementID := uuid.NewString()
	_, err := db.ExecContext(ctx, `
		INSERT INTO measurement_definitions (
			id, workspace_id, device_id, point_id, tag_id, equipment_id,
			definition_revision, source_binding_revision, name, quantity, semantic_kind
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, measurementID, fixture.workspaceID, fixture.deviceID, fixture.pointID, fixture.tagID,
		fixture.deviceID, "definition-1", "source-1", "Temperature", "temperature", "gauge")
	require.NoError(t, err)

	group := newBasicWriteGroup(fixture)
	group.Members[0].MeasurementID = &measurementID
	repo := NewSQLWriteGroupRepository(db)
	require.NoError(t, repo.Create(ctx, group))
	got, err := repo.Get(ctx, fixture.workspaceID, group.ID)
	require.NoError(t, err)
	require.NotNil(t, got.Members[0].MeasurementID)
	require.Equal(t, measurementID, *got.Members[0].MeasurementID)
}

func TestWriteGroupCreateRejectsMismatchedOptionalMeasurement(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	measurementID := uuid.NewString()
	_, err := db.ExecContext(ctx, `
		INSERT INTO measurement_definitions (
			id, workspace_id, device_id, point_id, tag_id, equipment_id,
			definition_revision, source_binding_revision, name, quantity, semantic_kind
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, measurementID, fixture.workspaceID, fixture.deviceID, "point-for-another-member", fixture.tagID,
		fixture.deviceID, "definition-1", "source-1", "Temperature", "temperature", "gauge")
	require.NoError(t, err)

	group := newBasicWriteGroup(fixture)
	group.Members[0].MeasurementID = &measurementID
	repo := NewSQLWriteGroupRepository(db)
	require.Error(t, repo.Create(ctx, group))

	var count int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups WHERE workspace_id = ?`, fixture.workspaceID).Scan(&count))
	require.Zero(t, count)
}

func TestWriteGroupCreateInTxSharesCallerTransaction(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	repo := NewSQLWriteGroupRepository(db)
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	group := newBasicWriteGroup(fixture)
	require.NoError(t, repo.CreateInTx(ctx, tx, group))
	_, err = repo.GetInTx(ctx, tx, fixture.workspaceID, group.ID)
	require.NoError(t, err)
	require.NoError(t, tx.Rollback())

	_, err = repo.Get(ctx, fixture.workspaceID, group.ID)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestWriteGroupCreateRejectsForeignOrMissingMemberWithoutPersisting(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	repo := NewSQLWriteGroupRepository(db)

	foreign := newBasicWriteGroup(fixture)
	foreign.Members[0].DeviceID = "device-not-in-workspace"
	err := repo.Create(ctx, foreign)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrWriteGroupValidation) || errors.Is(err, ErrNotFound))

	missing := newBasicWriteGroup(fixture)
	missing.Members[0].PointID = "point-missing"
	err = repo.Create(ctx, missing)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrWriteGroupValidation) || errors.Is(err, ErrNotFound))

	var count int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups WHERE workspace_id = ?`, fixture.workspaceID).Scan(&count))
	require.Zero(t, count)
}

type writeGroupFixture struct {
	workspaceID, deviceID, pointID, tagID, connectorID string
}

func seedWriteGroupFixture(ctx context.Context, t *testing.T, db *sql.DB) writeGroupFixture {
	t.Helper()
	workspace, err := NewService(NewSQLRepository(db)).GetOrCreate(ctx)
	require.NoError(t, err)
	fixture := writeGroupFixture{
		workspaceID: workspace.ID,
		deviceID:    uuid.NewString(),
		pointID:     uuid.NewString(),
		tagID:       uuid.NewString(),
		connectorID: uuid.NewString(),
	}
	_, err = db.ExecContext(ctx, `
		INSERT INTO devices (id, name, protocol, status, connection_config, readiness_status)
		VALUES (?, ?, ?, ?, ?, ?)
	`, fixture.deviceID, "PLC", schema.ProtocolModbusTCP, schema.DeviceStatusDraft, `{}`, `{}`)
	require.NoError(t, err)
	_, err = NewService(NewSQLRepository(db)).AttachDevice(ctx, fixture.deviceID)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `
		INSERT INTO points (id, device_id, name, address, function, data_type, mode, enabled)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, fixture.pointID, fixture.deviceID, "temperature", "40001", "FC03", schema.DataTypeFloat32, schema.PointModeReadOnly, true)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `
		INSERT INTO tags (id, key, key_lower, display_name, data_type, status, labels)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, fixture.tagID, "temperature", "temperature", "Temperature", schema.DataTypeFloat32, schema.TagStatusActive, `{}`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `
		INSERT INTO mappings (id, point_id, tag_id, transform_pipeline, status, enabled)
		VALUES (?, ?, ?, ?, ?, ?)
	`, uuid.NewString(), fixture.pointID, fixture.tagID, `[]`, schema.MappingStatusActive, true)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `
		INSERT INTO database_connectors (id, name, kind, connection_config, identity_revision, status, enabled)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, fixture.connectorID, "SQLite target", schema.DatabaseConnectorKindSQLite, `{"database":"persisted-db","schema":"persisted-schema"}`, "connector-revision-1", schema.DatabaseConnectorStatusReady, true)
	require.NoError(t, err)
	return fixture
}

func newBasicWriteGroup(fixture writeGroupFixture) *WriteGroup {
	return &WriteGroup{
		WorkspaceID: fixture.workspaceID,
		Name:        "raw tags",
		Members: []WriteGroupMember{{
			DeviceID: fixture.deviceID, PointID: fixture.pointID, TagID: fixture.tagID,
			TargetColumn: "temperature", Required: true,
		}},
		Destination: WriteGroupDestination{
			ConnectorID: fixture.connectorID, ConnectorRevision: "connector-revision-1",
			Database:    "line-a",
			TableSchema: "main", TableName: "raw_values", StorageStrategy: WriteGroupStorageStrategyCustom,
		},
	}
}
