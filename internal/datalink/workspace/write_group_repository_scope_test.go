package workspace

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestWriteGroupRejectsStaleRevisionsWithoutPersisting(t *testing.T) {
	cases := []struct {
		name     string
		mutate   func(*WriteGroup)
		expected error
	}{
		{
			name: "connector identity revision",
			mutate: func(group *WriteGroup) {
				group.Destination.ConnectorRevision = "stale-connector-revision"
			},
			expected: ErrWriteGroupConnectorRevisionConflict,
		},
		{
			name: "source revision",
			mutate: func(group *WriteGroup) {
				group.Members[0].SourceRevision = "stale-source-revision"
			},
			expected: ErrWriteGroupSourceRevisionConflict,
		},
		{
			name: "mapping revision",
			mutate: func(group *WriteGroup) {
				group.Members[0].MappingRevision = "stale-mapping-revision"
			},
			expected: ErrWriteGroupSourceRevisionConflict,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			db := openWorkspaceTestDB(t, ":memory:")
			defer db.Close()
			fixture := seedWriteGroupFixture(ctx, t, db)
			group := newBasicWriteGroup(fixture)
			tc.mutate(group)

			err := NewSQLWriteGroupRepository(db).Create(ctx, group)
			require.ErrorIs(t, err, tc.expected)
			require.ErrorIs(t, err, ErrSetupRevisionConflict)
			requireWriteGroupCount(ctx, t, db, fixture.workspaceID, 0)
		})
	}
}

func TestWriteGroupForeignWorkspaceCreateGetListAreSafeNotFound(t *testing.T) {
	ctx := t.Context()
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()
	fixture := seedWriteGroupFixture(ctx, t, db)
	repo := NewSQLWriteGroupRepository(db)

	group := newBasicWriteGroup(fixture)
	require.NoError(t, repo.Create(ctx, group))

	foreignWorkspaceID := uuid.NewString()
	foreign := newBasicWriteGroup(fixture)
	foreign.WorkspaceID = foreignWorkspaceID
	err := repo.Create(ctx, foreign)
	require.ErrorIs(t, err, ErrWriteGroupNotFound)
	require.ErrorIs(t, err, ErrNotFound)

	_, err = repo.Get(ctx, foreignWorkspaceID, group.ID)
	require.ErrorIs(t, err, ErrWriteGroupNotFound)
	require.ErrorIs(t, err, ErrNotFound)
	_, err = repo.List(ctx, foreignWorkspaceID)
	require.ErrorIs(t, err, ErrWriteGroupNotFound)
	require.ErrorIs(t, err, ErrNotFound)
	requireWriteGroupCount(ctx, t, db, fixture.workspaceID, 1)
}

func TestWriteGroupRejectsCrossSourceEdgesWithoutPersisting(t *testing.T) {
	cases := []struct {
		name    string
		prepare func(*testing.T, context.Context, *sql.DB, writeGroupFixture) string
		mutate  func(*WriteGroup, writeGroupFixture, string)
	}{
		{
			name: "same workspace device with another device point",
			mutate: func(group *WriteGroup, _ writeGroupFixture, deviceID string) {
				group.Members[0].DeviceID = deviceID
			},
			prepare: func(t *testing.T, ctx context.Context, db *sql.DB, _ writeGroupFixture) string {
				deviceID, _ := seedWriteGroupDevice(ctx, t, db)
				// The point remains fixture.pointID, which belongs to the original device.
				return deviceID
			},
		},
		{
			name: "missing mapping",
			prepare: func(t *testing.T, ctx context.Context, db *sql.DB, fixture writeGroupFixture) string {
				_, err := db.ExecContext(ctx, `UPDATE mappings SET enabled = 0 WHERE point_id = ? AND tag_id = ?`, fixture.pointID, fixture.tagID)
				require.NoError(t, err)
				return ""
			},
		},
		{
			name: "missing tag",
			mutate: func(group *WriteGroup, _ writeGroupFixture, _ string) {
				group.Members[0].TagID = uuid.NewString()
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			db := openWorkspaceTestDB(t, ":memory:")
			defer db.Close()
			fixture := seedWriteGroupFixture(ctx, t, db)
			prepared := ""
			if tc.prepare != nil {
				prepared = tc.prepare(t, ctx, db, fixture)
			}
			group := newBasicWriteGroup(fixture)
			if tc.mutate != nil {
				tc.mutate(group, fixture, prepared)
			}

			err := NewSQLWriteGroupRepository(db).Create(ctx, group)
			require.Error(t, err)
			require.True(t, errors.Is(err, ErrWriteGroupNotFound) || errors.Is(err, ErrWriteGroupValidation))
			requireWriteGroupCount(ctx, t, db, fixture.workspaceID, 0)
		})
	}
}

func requireWriteGroupCount(ctx context.Context, t *testing.T, db *sql.DB, workspaceID string, expected int) {
	t.Helper()
	var count int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups WHERE workspace_id = ?`, workspaceID).Scan(&count))
	require.Equal(t, expected, count)
}

func seedWriteGroupDevice(ctx context.Context, t *testing.T, db *sql.DB) (deviceID, pointID string) {
	t.Helper()
	deviceID = uuid.NewString()
	pointID = uuid.NewString()
	_, err := db.ExecContext(ctx, `
		INSERT INTO devices (id, name, protocol, status, connection_config, readiness_status)
		VALUES (?, ?, ?, ?, ?, ?)
	`, deviceID, "Second PLC", schema.ProtocolModbusTCP, schema.DeviceStatusDraft, `{}`, `{}`)
	require.NoError(t, err)
	_, err = NewService(NewSQLRepository(db)).AttachDevice(ctx, deviceID)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `
		INSERT INTO points (id, device_id, name, address, function, data_type, mode, enabled)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, pointID, deviceID, "second-temperature", "40002", "FC03", schema.DataTypeFloat32, schema.PointModeReadOnly, true)
	require.NoError(t, err)
	return deviceID, pointID
}
