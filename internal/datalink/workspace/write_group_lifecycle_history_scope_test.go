package workspace

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLifecycleHistoryUsesCallerSnapshotWithoutAnotherConnection(t *testing.T) {
	db, service, group := lifecycleReadFixture(t, true)
	before := lifecycleReadState(t, db)
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()
	historical := cloneWriteGroup(group)
	historical.Destination.ConnectorID = "transaction-only-connector"
	raw, err := json.Marshal(historical)
	require.NoError(t, err)
	_, err = tx.ExecContext(ctx, `UPDATE write_group_versions SET payload=? WHERE group_id=?`, string(raw), group.ID)
	require.NoError(t, err)
	view, err := service.repo.lifecycleHistoryInTx(ctx, tx, group.WorkspaceID, group.ID)
	require.NoError(t, err)
	require.Len(t, view.Generations, 1)
	require.Len(t, view.Generations[0].Revisions, 1)
	require.Equal(t, "transaction-only-connector", view.Generations[0].Revisions[0].Destination.ConnectorID)
	require.Equal(t, group.Destination.ConnectorID, view.Generations[0].SavedDestination.ConnectorID)
	require.NoError(t, tx.Rollback())
	require.Equal(t, before, lifecycleReadState(t, db))
}

func TestLifecycleHistoryCannotHideOrphanDeliveryRevision(t *testing.T) {
	for _, mode := range []string{"unknown-version", "foreign-workspace"} {
		t.Run(mode, func(t *testing.T) {
			db, service, group := lifecycleReadFixture(t, true)
			workspace, revision := group.WorkspaceID, "orphan-revision"
			if mode == "foreign-workspace" {
				workspace, revision = "foreign", group.AppliedRevision
			}
			_, err := db.Exec(`INSERT INTO wg_delivery_samples(workspace_id,group_id,group_revision,sample_id,member_key,observed_at,bucket_start,payload,payload_digest,created_at) VALUES(?,?,?,?,?,'2026-10-02T12:00:10Z','2026-10-02T12:00:00Z','private-payload','digest','2026-10-02T12:00:10Z')`, workspace, group.ID, revision, "S-orphan", "member")
			require.NoError(t, err)
			before := lifecycleReadState(t, db)
			result, err := service.Lifecycle(t.Context(), group.ID)
			require.Error(t, err)
			require.Nil(t, result)
			require.Equal(t, before, lifecycleReadState(t, db))
		})
	}
}
