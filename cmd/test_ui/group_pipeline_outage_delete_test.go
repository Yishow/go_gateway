package main

import (
	"os"
	"testing"
	"time"

	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/testutil/legacyrecording"

	"github.com/stretchr/testify/require"
)

func recordingDeleteState(t *testing.T, env *outageEnv, id string) map[string][][]any {
	t.Helper()
	state := make(map[string][][]any)
	for _, query := range []string{
		`SELECT * FROM write_groups WHERE id=? ORDER BY rowid`,
		`SELECT * FROM write_group_versions WHERE group_id=? ORDER BY rowid`,
		`SELECT * FROM write_group_members WHERE group_id=? ORDER BY rowid`,
		`SELECT * FROM wg_runtime_versions WHERE group_id=? ORDER BY rowid`,
		`SELECT * FROM wg_delivery_samples WHERE group_id=? ORDER BY rowid`,
		`SELECT * FROM wg_delivery_outbox WHERE group_id=? ORDER BY rowid`,
	} {
		rows, err := env.db.QueryContext(t.Context(), query, id)
		require.NoError(t, err)
		columns, err := rows.Columns()
		require.NoError(t, err)
		for rows.Next() {
			values, destinations := make([]any, len(columns)), make([]any, len(columns))
			for i := range values {
				destinations[i] = &values[i]
			}
			require.NoError(t, rows.Scan(destinations...))
			state[query] = append(state[query], values)
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
	}
	return state
}

func TestProductionGroupOutageRecoveryLegacyTombstoneBacklogIsCompletedByWorkers(t *testing.T) {
	env := newOutageEnv(t)
	group := env.createAndApply(t, "A")
	ctx := t.Context()
	resolved, err := env.services.writeGroups.ResolveAppliedAt(ctx, group.ID, time.Now().UTC().Add(time.Minute))
	require.NoError(t, err)
	base := resolved.EffectiveAt.Truncate(10 * time.Second).Add(10 * time.Second)
	clock := &testClock{}
	clock.set(base.Add(time.Second))
	pipe := env.pipeline(clock, "node-1/first")
	require.NoError(t, pipe.Start(ctx))
	defer func() { _ = pipe.Stop(2 * time.Second) }()

	// Destination offline: two accepted buckets wait in the durable outbox.
	offline := env.targets["connector-A"] + ".offline"
	require.NoError(t, os.Rename(env.targets["connector-A"], offline))
	env.feed(t, pipe, clock, []*workspace.WriteGroup{group}, base, 21.5, 100)
	env.feed(t, pipe, clock, []*workspace.WriteGroup{group}, base.Add(10*time.Second), 22.5, 200)
	require.Eventually(t, func() bool { return outboxCount(t, env.db, "connector-A", "state != 'sql_committed'") == 2 }, 10*time.Second, 25*time.Millisecond)
	require.NoError(t, pipe.Stop(2*time.Second)) // freeze worker updates while comparing the refused delete

	before := recordingDeleteState(t, env, group.ID)
	mutation := lifecycleMutation(t, env, group.ID)
	_, err = env.services.writeGroups.Delete(ctx, group.ID, mutation)
	var refusal *workspace.WriteGroupDeleteError
	require.ErrorAs(t, err, &refusal)
	require.Equal(t, "WRITE_GROUP_DELETE_BACKLOG", refusal.Code)
	require.Equal(t, mutation, lifecycleMutation(t, env, group.ID))
	require.Equal(t, before, recordingDeleteState(t, env, group.ID))
	_, err = env.services.writeGroups.Disable(ctx, group.ID, lifecycleMutation(t, env, group.ID))
	require.NoError(t, err)
	before = recordingDeleteState(t, env, group.ID)
	mutation = lifecycleMutation(t, env, group.ID)
	_, err = env.services.writeGroups.Delete(ctx, group.ID, mutation)
	require.ErrorAs(t, err, &refusal)
	require.Equal(t, "WRITE_GROUP_DELETE_BACKLOG", refusal.Code)
	require.Equal(t, mutation, lifecycleMutation(t, env, group.ID))
	require.Equal(t, before, recordingDeleteState(t, env, group.ID))

	// Formerly removed groups remain drainable, but cannot be purged with backlog.
	require.NoError(t, legacyrecording.Tombstone(ctx, env.db, group.ID))
	tombstone, err := env.services.writeGroups.Get(ctx, group.ID)
	require.NoError(t, err)
	require.Equal(t, workspace.WriteGroupStatusDeleted, tombstone.Group.Status)
	before = recordingDeleteState(t, env, group.ID)
	mutation = lifecycleMutation(t, env, group.ID)
	_, err = env.services.writeGroups.Delete(ctx, group.ID, mutation)
	require.ErrorAs(t, err, &refusal)
	require.Equal(t, "WRITE_GROUP_DELETE_BACKLOG", refusal.Code)
	require.Equal(t, mutation, lifecycleMutation(t, env, group.ID))
	require.Equal(t, before, recordingDeleteState(t, env, group.ID))
	require.NoError(t, pipe.Reconcile(ctx))
	before = recordingDeleteState(t, env, group.ID)
	err = pipe.AcceptSample(ctx, env.envelope(group, 0, "after-stop", clock.now(), 999.0))
	require.Error(t, err, "a stale transport interest must not accept new intake")
	require.Equal(t, before, recordingDeleteState(t, env, group.ID))
	require.Equal(t, 2, outboxCount(t, env.db, "connector-A", "state != 'sql_committed'"))
	view, err := pipe.Delivery(ctx, group.ID)
	require.NoError(t, err)
	require.Len(t, view.Backlog, 1)
	require.Equal(t, 2, view.Backlog[0].Pending)
	require.Equal(t, group.AppliedRevision, view.Backlog[0].GroupRevision)
	require.Equal(t, group.Destination.ConnectorID, view.Backlog[0].ConnectorID)
	require.Equal(t, group.Destination.ConnectorRevision, view.Backlog[0].ConnectorRevision)
	require.Equal(t, group.Destination.TableSchema, view.Backlog[0].TableSchema)
	require.Equal(t, group.Destination.TableName, view.Backlog[0].TableName)

	// Workers deliver the accepted rows to their frozen destination, once and in order.
	require.NoError(t, pipe.Start(ctx))
	require.NoError(t, os.Rename(offline, env.targets["connector-A"]))
	require.Eventually(t, func() bool { return len(destinationRows(t, env.targets["connector-A"])) == 2 }, 15*time.Second, 25*time.Millisecond)
	require.Equal(t, [][2]float64{{21.5, 100}, {22.5, 200}}, destinationRows(t, env.targets["connector-A"]))
	require.Eventually(t, func() bool { return outboxCount(t, env.db, "connector-A", "state = 'sql_committed'") == 2 }, 5*time.Second, 25*time.Millisecond)
	view, err = pipe.Delivery(ctx, group.ID)
	require.NoError(t, err)
	require.Equal(t, 2, view.Stages.SQLCommitted)
	require.Empty(t, view.Backlog)
	require.NoError(t, pipe.Stop(2*time.Second))
	deleted, err := env.services.writeGroups.Delete(ctx, group.ID, lifecycleMutation(t, env, group.ID))
	require.NoError(t, err)
	require.True(t, deleted.Deleted)
	require.Equal(t, [][2]float64{{21.5, 100}, {22.5, 200}}, destinationRows(t, env.targets["connector-A"]), "permanent deletion preserves destination rows")
}
