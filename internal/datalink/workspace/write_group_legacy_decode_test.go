package workspace

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// legacyWriteGroupPayload is a snapshot payload written before circuit
// profiles, layout versions and member circuit identity existed.
const legacyWriteGroupPayload = `{
	"id": "legacy-group", "workspace_id": "workspace-1", "revision": "revision-1",
	"applied_revision": "revision-1", "name": "legacy wide", "status": "running",
	"members": [{"device_id": "device-1", "point_id": "point-1", "tag_id": "tag-1",
		"target_column": "v_legacy", "required": true}],
	"destination": {"connector_id": "connector-1", "connector_revision": "connector-revision-1",
		"database": "line-a", "table_schema": "main", "table_name": "gw_group_legacy",
		"storage_strategy": "managed"},
	"row_policy": {"interval_seconds": 60, "allowed_lateness_seconds": 5},
	"write_policy": {"mode": "append", "dedupe_capability": "receipt"},
	"migration": {}, "created_at": "2026-01-01T00:00:00Z", "updated_at": "2026-01-01T00:00:00Z"
}`

func TestLegacyGroupDecode(t *testing.T) {
	t.Run("legacy snapshot payload decodes as wide without a profile", func(t *testing.T) {
		var group WriteGroup
		require.NoError(t, json.Unmarshal([]byte(legacyWriteGroupPayload), &group))

		require.Equal(t, WriteGroupLayoutManagedWideV1, group.EffectiveLayoutVersion())
		require.Empty(t, group.LayoutVersion, "absent layout must stay absent, not be rewritten")
		require.Nil(t, group.RecordingProfile)
		require.Empty(t, group.Members[0].CircuitID)
		require.Empty(t, group.Members[0].CircuitName)
		require.Equal(t, 60, group.RowPolicy.IntervalSeconds)

		encoded, err := json.Marshal(&group)
		require.NoError(t, err)
		require.NotContains(t, string(encoded), "layout_version")
		require.NotContains(t, string(encoded), "recording_profile")
		require.NotContains(t, string(encoded), "circuit_id")
	})

	t.Run("a stored layout version is returned unchanged", func(t *testing.T) {
		group := WriteGroup{LayoutVersion: WriteGroupLayoutManagedWideV1}
		require.Equal(t, WriteGroupLayoutManagedWideV1, group.EffectiveLayoutVersion())
	})

	t.Run("pre-migration rows and applied snapshots stay readable as wide", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		ctx := t.Context()
		now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		// Only the columns that existed before circuit profiles are written, so
		// every new column keeps its NULL default.
		_, err := h.db.ExecContext(ctx, `
			INSERT INTO write_groups (
				id, workspace_id, revision, applied_revision, name, status,
				destination_connector_id, destination_connector_revision, destination_database,
				destination_table_schema, destination_table_name, destination_storage_strategy,
				row_policy, write_policy, migration, created_at, updated_at
			) VALUES ('legacy-group', ?, 'revision-1', 'revision-1', 'legacy wide', 'running',
				?, 'connector-revision-1', 'line-a', 'main', 'gw_group_legacy', 'managed',
				'{"interval_seconds":60}', '{"mode":"append"}', '{}', ?, ?)
		`, h.fixture.workspaceID, h.fixture.connectorID, now, now)
		require.NoError(t, err)
		primary := h.primary()
		_, err = h.db.ExecContext(ctx, `
			INSERT INTO write_group_members (group_id, member_index, device_id, point_id, tag_id, target_column)
			VALUES ('legacy-group', 0, ?, ?, ?, 'v_legacy')
		`, primary.deviceID, primary.pointID, primary.tagID)
		require.NoError(t, err)
		snapshotPayload := strings.ReplaceAll(legacyWriteGroupPayload, `"workspace-1"`, `"`+h.fixture.workspaceID+`"`)
		_, err = h.db.ExecContext(ctx, `
			INSERT INTO write_group_versions (group_id, group_revision, effective_at, payload)
			VALUES ('legacy-group', 'revision-1', ?, ?)
		`, now, snapshotPayload)
		require.NoError(t, err)

		loaded, err := h.service.Get(ctx, "legacy-group")
		require.NoError(t, err)
		require.Equal(t, WriteGroupLayoutManagedWideV1, loaded.Group.EffectiveLayoutVersion())
		require.Nil(t, loaded.Group.RecordingProfile)
		require.Empty(t, loaded.Group.Members[0].CircuitID)
		require.Equal(t, 60, loaded.Group.RowPolicy.IntervalSeconds)

		listed, err := h.service.List(ctx)
		require.NoError(t, err)
		require.Len(t, listed.Groups, 1)
		require.Nil(t, listed.Groups[0].RecordingProfile)

		snapshot, err := h.service.repo.versionByRevisionInTx(ctx, h.db, h.fixture.workspaceID, "legacy-group", "revision-1")
		require.NoError(t, err)
		require.Equal(t, WriteGroupLayoutManagedWideV1, snapshot.Group.EffectiveLayoutVersion())
		require.Nil(t, snapshot.Group.RecordingProfile)
	})

	t.Run("an unrecognized stored layout fails closed instead of reading as wide", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		ctx := t.Context()
		created := h.mustEnsure(h.draft(WriteGroupProfileInstantaneous, "circuit-a", "A", h.primary()))
		_, err := h.db.ExecContext(ctx, `UPDATE write_groups SET layout_version = 'managed_future_v9' WHERE id = ?`, created.ID)
		require.NoError(t, err)

		_, err = h.service.Get(ctx, created.ID)
		require.ErrorIs(t, err, ErrWriteGroupValidation)
		_, err = h.service.List(ctx)
		require.ErrorIs(t, err, ErrWriteGroupValidation)

		_, err = h.db.ExecContext(ctx, `UPDATE write_groups SET profile_unit = NULL, layout_version = NULL WHERE id = ?`, created.ID)
		require.NoError(t, err)
		_, err = h.service.Get(ctx, created.ID)
		require.ErrorIs(t, err, ErrWriteGroupValidation, "a half-written profile is corrupt, not wide")
	})

	t.Run("saving a legacy group does not invent a layout, profile or interval", func(t *testing.T) {
		h := newCircuitProfileHarness(t, ":memory:")
		ctx := t.Context()
		legacy := newBasicManagedDraft(h.fixture)
		legacy.Destination.StorageStrategy = WriteGroupStorageStrategyManaged
		legacy.RowPolicy.IntervalSeconds = 60
		created, err := h.service.Create(ctx, h.mutation(ctx, legacy))
		require.NoError(t, err)

		edit := cloneWriteGroup(created.Group)
		edit.Name = "legacy renamed"
		mutation := h.mutation(ctx, edit)
		mutation.ExpectedGroupRevision = created.Group.Revision
		updated, err := h.service.Update(ctx, created.Group.ID, mutation)
		require.NoError(t, err)

		require.Empty(t, updated.Group.LayoutVersion)
		require.Nil(t, updated.Group.RecordingProfile)
		require.Equal(t, 60, updated.Group.RowPolicy.IntervalSeconds)
		require.Equal(t, 1, h.countRows(
			`SELECT COUNT(*) FROM write_groups WHERE id = ? AND layout_version IS NULL AND profile_kind IS NULL`, created.Group.ID))
		require.Equal(t, 0, h.countRows(`SELECT COUNT(*) FROM write_group_profile_keys`))
	})
}
