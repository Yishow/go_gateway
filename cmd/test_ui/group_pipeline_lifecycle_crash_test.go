package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func killLifecycleAfterAck(t *testing.T, env *outageEnv, group *workspace.WriteGroup, action string) {
	t.Helper()
	var seq int
	var name, path string
	require.NoError(t, env.db.QueryRowContext(t.Context(), `PRAGMA database_list`).Scan(&seq, &name, &path))
	ready := filepath.Join(filepath.Dir(path), "crash-ready")
	cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^TestProductionLifecycleCrashChild$", "-test.count=1")
	cmd.Env = append(os.Environ(), lifecycleCrashEnv+"="+path, "GW_LIFECYCLE_GROUP="+group.ID, "GW_LIFECYCLE_ACTION="+action)
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	require.NoError(t, cmd.Start())
	t.Cleanup(func() {
		if cmd.ProcessState == nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	})
	require.Eventually(t, func() bool { _, err := os.Stat(ready); return err == nil }, 15*time.Second, 20*time.Millisecond)
	require.NoError(t, cmd.Process.Kill())
	var exit *exec.ExitError
	require.ErrorAs(t, cmd.Wait(), &exit, output.String())
	require.False(t, exit.Exited(), "the process was killed before closure, not exited gracefully")
}

func lifecycleJournalCount(t *testing.T, env *outageEnv, revision string, openOnly bool) int {
	t.Helper()
	query := `SELECT COUNT(*) FROM wg_delivery_samples WHERE group_revision = ?`
	if openOnly {
		query += ` AND consumed = 0`
	}
	var count int
	require.NoError(t, env.db.QueryRowContext(t.Context(), query, revision).Scan(&count))
	return count
}

func TestProductionLifecycleKilledHistoricalJournalDrainsOriginalRevision(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		for _, action := range []string{"disable", "delete", "replace"} {
			t.Run(kind+"/"+action, func(t *testing.T) {
				target := newLifecycleTarget(t, kind)
				env, clock := target.env, &testClock{}
				clock.set(lifecycleBase.Add(-time.Second))
				env.services.writeGroups.WithClock(clock.now)
				group := env.createAndApply(t, "A")
				killLifecycleAfterAck(t, env, group, action)
				require.Equal(t, 2, lifecycleJournalCount(t, env, group.AppliedRevision, true))
				require.Zero(t, outboxCount(t, env.db, "connector-A", "1=1"), "kill preceded closure")
				target.down()
				clock.set(lifecycleBase.Add(21 * time.Second))
				reopenLifecycleConfiguration(t, env, clock)
				pipe := env.pipeline(clock, "node-1/recovered")
				require.NoError(t, pipe.Reconcile(t.Context()))
				// Replaying even a formerly matching historical sample cannot create
				// new intake in a drain-only boundary.
				require.NoError(t, pipe.AcceptSample(t.Context(), env.envelope(group, 0, "old-must-ignore", lifecycleBase.Add(3*time.Second), 999.0)))
				require.Equal(t, 2, lifecycleJournalCount(t, env, group.AppliedRevision, false))
				if action != "replace" {
					require.False(t, pipe.WantsSample("device-1", "point-t-A", "tag-t-A"))
				}
				pipe.TickAll(t.Context())
				require.Zero(t, lifecycleJournalCount(t, env, group.AppliedRevision, true))
				var effect, revision, connectorID, connectorRevision, tableSchema, table, payload string
				require.NoError(t, env.db.QueryRowContext(t.Context(), `SELECT effect_key, group_revision, connector_id, connector_revision, table_schema, table_name, payload
					FROM wg_delivery_outbox WHERE group_id = ?`, group.ID).Scan(&effect, &revision, &connectorID, &connectorRevision, &tableSchema, &table, &payload))
				require.Equal(t, group.AppliedRevision, revision)
				require.Equal(t, group.Destination.ConnectorID, connectorID)
				require.Equal(t, group.Destination.ConnectorRevision, connectorRevision)
				require.Equal(t, group.Destination.TableSchema, tableSchema)
				require.Equal(t, group.Destination.TableName, table)
				row, err := groupdelivery.DecodeRowPayload([]byte(payload))
				require.NoError(t, err)
				require.Equal(t, lifecycleBase, row.BucketStart)
				for range 2 {
					require.NoError(t, pipe.Reconcile(t.Context()))
					pipe.TickAll(t.Context())
				}
				require.Equal(t, 1, outboxCount(t, env.db, "connector-A", "1=1"))
				target.up()
				require.NoError(t, pipe.Start(t.Context()))
				t.Cleanup(func() { require.NoError(t, pipe.Stop(2*time.Second)) })
				require.Eventually(t, func() bool { return outboxCount(t, env.db, "connector-A", "state = 'sql_committed'") == 1 }, 10*time.Second, 20*time.Millisecond)
				require.Equal(t, [][2]float64{{2, 20}}, target.rows())
				var committedEffect string
				require.NoError(t, env.db.QueryRowContext(t.Context(), `SELECT effect_key FROM wg_delivery_receipts`).Scan(&committedEffect))
				require.Equal(t, effect, committedEffect, "delivery keeps the original effect identity")
			})
		}
	}
}

func TestProductionLifecycleRapidReenableJournalsAfterOldCutoff(t *testing.T) {
	env, clock := newOutageEnv(t), &testClock{}
	clock.set(lifecycleBase.Add(-time.Second))
	env.services.writeGroups.WithClock(clock.now)
	old := env.createAndApply(t, "A")
	clock.set(lifecycleBase.Add(time.Second))
	pipe := env.pipeline(clock, "node-1/reenable")
	require.NoError(t, pipe.Reconcile(t.Context()))
	lifecycleAccept(t, env, pipe, old, clock, 11)
	clock.set(lifecycleBase.Add(12 * time.Second))
	_, err := env.services.writeGroups.Disable(t.Context(), old.ID, lifecycleMutation(t, env, old.ID))
	require.NoError(t, err)
	require.NoError(t, pipe.Reconcile(t.Context()))
	clock.set(lifecycleBase.Add(15 * time.Second))
	next, err := env.services.writeGroups.Apply(t.Context(), old.ID, lifecycleMutation(t, env, old.ID))
	require.NoError(t, err)
	require.NotEqual(t, old.AppliedRevision, next.Group.AppliedRevision)
	require.NoError(t, pipe.Reconcile(t.Context()))
	lifecycleAccept(t, env, pipe, next.Group, clock, 21)
	require.Equal(t, 2, lifecycleJournalCount(t, env, old.AppliedRevision, true))
	require.Equal(t, 2, lifecycleJournalCount(t, env, next.Group.AppliedRevision, true))
	reopenLifecycleConfiguration(t, env, clock)
	recovered := env.pipeline(clock, "node-1/reenable-restart")
	require.NoError(t, recovered.Reconcile(t.Context()))
	clock.set(lifecycleBase.Add(31 * time.Second))
	recovered.TickAll(t.Context())
	require.Equal(t, 2, outboxCount(t, env.db, "connector-A", "1=1"))
	require.Zero(t, lifecycleJournalCount(t, env, old.AppliedRevision, true))
	require.Zero(t, lifecycleJournalCount(t, env, next.Group.AppliedRevision, true))
}

func TestProductionLifecycleUnprovableRecoveryKeepsJournalAndBlocksLegacy(t *testing.T) {
	for _, failure := range []string{"missing-offline", "revision-mismatch", "changed-source"} {
		t.Run(failure, func(t *testing.T) {
			target := newLifecycleTarget(t, "sqlite")
			env, clock := target.env, &testClock{}
			clock.set(lifecycleBase.Add(-time.Second))
			env.services.writeGroups.WithClock(clock.now)
			group := env.createAndApply(t, "A")
			clock.set(lifecycleBase.Add(time.Second))
			first := env.pipeline(clock, "node-1/unprovable")
			require.NoError(t, first.Reconcile(t.Context()))
			lifecycleAccept(t, env, first, group, clock, 2)
			if failure == "revision-mismatch" {
				_, err := env.db.ExecContext(t.Context(), `UPDATE wg_runtime_versions SET payload = json_set(payload, '$.applied_revision', 'foreign-version')`)
				require.NoError(t, err)
			} else {
				_, err := env.db.ExecContext(t.Context(), `DELETE FROM wg_runtime_versions`)
				require.NoError(t, err)
				if failure == "changed-source" {
					_, err = env.db.ExecContext(t.Context(), `UPDATE points SET data_type = 'uint64' WHERE id = 'point-t-A'`)
					require.NoError(t, err)
				} else {
					target.down()
				}
			}
			clock.set(lifecycleBase.Add(21 * time.Second))
			reopenLifecycleConfiguration(t, env, clock)
			restarted := env.pipeline(clock, "node-1/unprovable-restart")
			require.NoError(t, restarted.Reconcile(t.Context()))
			require.False(t, restarted.WantsSample("device-1", "point-t-A", "tag-t-A"))
			require.NotEmpty(t, restarted.Status())
			for _, status := range restarted.Status() {
				require.Equal(t, "blocked", status.State)
			}
			restarted.TickAll(t.Context())
			require.Equal(t, 2, lifecycleJournalCount(t, env, group.AppliedRevision, true))
			require.Zero(t, outboxCount(t, env.db, "connector-A", "1=1"))
			_, err := env.services.dbMapping.Create(t.Context(), dbtarget.CreateTargetMappingRequest{
				TagID: "tag-t-A", ConnectorID: "connector-A", TableSchema: "main", TableName: "readings", ColumnName: "temperature", WriteMode: "insert",
			})
			require.ErrorIs(t, err, workspace.ErrWriteGroupLegacyWriteConflict, "blocked recovery cannot authorize a legacy writer to take over")
		})
	}
}
