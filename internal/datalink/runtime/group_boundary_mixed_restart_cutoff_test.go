package runtime

import (
	"database/sql"
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
)

func installMixedRestartFence(t *testing.T, db *sql.DB, source *longFixture, end time.Time) {
	t.Helper()
	target := mixedBoundaryFixture(t)
	target.config.Group.ID = "mixed-cutover-target"
	target.config.Group.Destination.TableName = "mixed_target_table"
	target.config.StartsAt = end
	mixedDurableFixture(t, target, db)
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	runtimeHandoffVersion(t, tx, source, source.config.StartsAt)
	runtimeHandoffVersion(t, tx, target, end)
	g := source.config.Group
	moved := []workspace.CutoverMemberRef{}
	for _, m := range g.Members {
		moved = append(moved, workspace.CutoverMemberRef{DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID})
	}
	detail, err := json.Marshal(map[string]any{"ownership_fence_revision": "mixed-restart", "ended": []workspace.CutoverEnd{{GroupID: g.ID, Revision: g.AppliedRevision}}, "started": []workspace.CutoverStart{{GroupID: target.config.Group.ID, Revision: target.config.Group.AppliedRevision}}, "moved": moved})
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_group_cutovers(id,workspace_id,kind,source_group_id,boundary,review_digest,detail,intent) VALUES('mixed-restart',?,'forward',?,?,'review',?,'{}')`, g.WorkspaceID, g.ID, end.Format(time.RFC3339Nano), string(detail))
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_group_revision_ends(workspace_id,group_id,group_revision,ends_at,cutover_id) VALUES(?,?,?,?,'mixed-restart')`, g.WorkspaceID, g.ID, g.AppliedRevision, end.Format(time.RFC3339Nano))
	require.NoError(t, err)
	for direction, f := range map[string]*longFixture{"ended": source, "started": target} {
		owner := f.config.Group
		for _, m := range owner.Members {
			_, err = tx.Exec(`INSERT INTO write_group_ownership_fences(fence_revision,workspace_id,direction,group_id,group_revision,device_id,point_id,tag_id,profile_kind,connector_id,connector_revision,database_name,table_schema,table_name) VALUES('mixed-restart',?,?,?,?,?,?,?,?,?,?,?,?,?)`, owner.WorkspaceID, direction, owner.ID, owner.AppliedRevision, m.DeviceID, m.PointID, m.TagID, string(owner.RecordingProfile.Kind), owner.Destination.ConnectorID, owner.Destination.ConnectorRevision, owner.Destination.Database, owner.Destination.TableSchema, owner.Destination.TableName)
			require.NoError(t, err)
		}
	}
	require.NoError(t, tx.Commit())
}

func TestMixedDurableRestartAcceptsCanonicalNewCutoff(t *testing.T) {
	db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
	f := mixedBoundaryFixture(t)
	store := mixedDurableFixture(t, f, db)
	b, err := NewGroupBoundaryContext(t.Context(), f.config)
	require.NoError(t, err)
	f.accept(t, b, f.envelope(0, "first", boundaryAt(10), int64(20)), f.envelope(1, "slow", boundaryAt(10), "ON"), f.envelope(1, "future", boundaryAt(29), "OFF"))
	require.NoError(t, b.Tick(t.Context(), boundaryAt(15)))
	installMixedRestartFence(t, db, f, boundaryAt(23))
	g := f.config.Group
	batch, err := store.PrepareCutoverHandoffs(t.Context(), groupdelivery.GroupKey{WorkspaceID: g.WorkspaceID, GroupID: g.ID, GroupRevision: g.AppliedRevision}, 0, 128)
	require.NoError(t, err)
	require.Equal(t, 1, batch.Accepted)
	for _, wrong := range []time.Time{boundaryAt(22), boundaryAt(24), time.Time{}} {
		f.config.EndsAt = wrong
		_, err = NewGroupBoundaryContext(t.Context(), f.config)
		require.ErrorIs(t, err, groupdelivery.ErrOwnershipConflict)
	}
	f.config.EndsAt = boundaryAt(23)
	restored, err := NewGroupBoundaryContext(t.Context(), f.config)
	require.NoError(t, err, "cause=%v", errors.Unwrap(err))
	require.NoError(t, restored.Tick(t.Context(), boundaryAt(60)))
	require.True(t, restored.Retired())
	var raw string
	require.NoError(t, db.QueryRow(`SELECT outcome FROM wg_delivery_mixed_outcomes WHERE due_at=?`, boundaryAt(60).Format("2006-01-02T15:04:05.000000000Z")).Scan(&raw))
	var out snapshot.Outcome
	require.NoError(t, json.Unmarshal([]byte(raw), &out))
	require.Equal(t, "slow", out.Members[1].Sample.SampleID)
	require.True(t, out.Members[1].Window.ActualEnd.Equal(boundaryAt(23)))
	again, err := NewGroupBoundaryContext(t.Context(), f.config)
	require.NoError(t, err)
	require.True(t, again.Retired())
}

func TestMixedDurableRestartRejectsUnprovenNewCutoff(t *testing.T) {
	db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
	f := mixedBoundaryFixture(t)
	mixedDurableFixture(t, f, db)
	b, err := NewGroupBoundaryContext(t.Context(), f.config)
	require.NoError(t, err)
	require.NoError(t, b.Tick(t.Context(), boundaryAt(15)))
	var before string
	require.NoError(t, db.QueryRow(`SELECT checkpoint FROM wg_delivery_mixed_progress`).Scan(&before))
	f.config.EndsAt = boundaryAt(23)
	_, err = NewGroupBoundaryContext(t.Context(), f.config)
	require.ErrorIs(t, err, groupdelivery.ErrOwnershipConflict)
	var after string
	require.NoError(t, db.QueryRow(`SELECT checkpoint FROM wg_delivery_mixed_progress`).Scan(&after))
	require.Equal(t, before, after)
}

func TestMixedDurableRestartCannotReopenClosedHistory(t *testing.T) {
	db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
	f := mixedBoundaryFixture(t)
	mixedDurableFixture(t, f, db)
	b, err := NewGroupBoundaryContext(t.Context(), f.config)
	require.NoError(t, err)
	require.NoError(t, b.Tick(t.Context(), boundaryAt(15)))
	installMixedRestartFence(t, db, f, boundaryAt(10))
	var before string
	require.NoError(t, db.QueryRow(`SELECT checkpoint FROM wg_delivery_mixed_progress`).Scan(&before))
	f.config.EndsAt = boundaryAt(10)
	_, err = NewGroupBoundaryContext(t.Context(), f.config)
	require.Error(t, err)
	var after string
	require.NoError(t, db.QueryRow(`SELECT checkpoint FROM wg_delivery_mixed_progress`).Scan(&after))
	require.Equal(t, before, after)
}
