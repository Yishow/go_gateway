package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPermanentDeleteDrainedTwoGenerationRecording(t *testing.T) {
	ctx := t.Context()
	db, service, predecessor := lifecyclePreservationFixture(t)
	service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
	_, err := db.ExecContext(ctx, `PRAGMA foreign_keys=ON`)
	require.NoError(t, err)
	destinationPath := filepath.Join(t.TempDir(), "destination.sqlite")
	destination, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, destination.Close()) })
	_, err = destination.ExecContext(ctx, `ATTACH DATABASE ? AS public`, destinationPath)
	require.NoError(t, err)
	_, err = destination.ExecContext(ctx, `CREATE TABLE public.gw_long_a (value INTEGER)`)
	require.NoError(t, err)
	_, err = destination.ExecContext(ctx, `WITH RECURSIVE n(v) AS (VALUES(1) UNION ALL SELECT v+1 FROM n WHERE v<1200) INSERT INTO public.gw_long_a SELECT v FROM n`)
	require.NoError(t, err)
	config, err := json.Marshal(map[string]string{"path": destinationPath})
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE database_connectors SET kind='sqlite',connection_config=? WHERE id=?`, string(config), predecessor.Destination.ConnectorID)
	require.NoError(t, err)
	predecessor.Destination.Database, predecessor.Destination.TableSchema, predecessor.Destination.TableName = destinationPath, "public", "gw_long_a"
	_, err = db.ExecContext(ctx, `UPDATE write_groups SET destination_database=?,destination_table_schema='public',destination_table_name='gw_long_a' WHERE id='G1'`, destinationPath)
	require.NoError(t, err)
	version, err := json.Marshal(predecessor)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE write_group_versions SET payload=? WHERE group_id='G1' AND group_revision=?`, string(version), predecessor.AppliedRevision)
	require.NoError(t, err)
	seedLifecycleAcceptedData(t, db, predecessor)
	_, err = db.ExecContext(ctx, `UPDATE wg_delivery_samples SET consumed=1`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE wg_delivery_outbox SET state='sql_committed'`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE wg_delivery_outbox SET state='operator_skipped' WHERE effect_key='Q1'`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `DELETE FROM managed_schema_operations`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE write_groups SET status='deleted' WHERE id='G1'`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO write_group_lineages(lineage_id,workspace_id,current_group_id) VALUES('L1',?,'G1')`, predecessor.WorkspaceID)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO write_group_generations(group_id,lineage_id,generation) VALUES('G1','L1',1)`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO write_group_basic_keys(workspace_id,device_id,canonical_role,group_id,lineage_id) VALUES(?,?,'basic','G1','L1')`, predecessor.WorkspaceID, predecessor.Members[0].DeviceID)
	require.NoError(t, err)
	head := cloneWriteGroup(predecessor)
	head.ID, head.Revision, head.AppliedRevision, head.Status = "G2", "R2", "", WriteGroupStatusDisabled
	require.NoError(t, service.repo.insertGroup(ctx, db, head))
	_, err = db.ExecContext(ctx, `INSERT INTO write_group_generations(group_id,lineage_id,generation,predecessor_group_id) VALUES('G2','L1',2,'G1')`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE write_group_lineages SET current_group_id='G2',revision=2 WHERE lineage_id='L1'`)
	require.NoError(t, err)
	seedOwnedAdditionalDeliveryRows(t, db, predecessor)
	operationDetail := `{"scope":{"workspace_id":"` + predecessor.WorkspaceID + `","lineage_id":"L1","predecessor_group_id":"G1","successor_group_id":"G2"}}`
	_, err = db.ExecContext(ctx, `INSERT INTO managed_schema_operations(operation_id,token,workspace_id,scope_key,owner,action,status,created_at,updated_at,detail) VALUES('Odone','Tdone',?,'scope-done','owner','successor_reserve','succeeded',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,?)`, predecessor.WorkspaceID, operationDetail)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO managed_schema_operations(operation_id,token,workspace_id,scope_key,owner,action,status,created_at,updated_at,detail) VALUES('Ofailed','Tfailed',?,'scope-failed','owner','successor_first_start','failed',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,?)`, predecessor.WorkspaceID, operationDetail)
	require.NoError(t, err)

	current, err := service.Get(ctx, "G2")
	require.NoError(t, err)
	deleted, err := service.Delete(ctx, "G2", WriteGroupMutation{
		WorkspaceID: current.Group.WorkspaceID, ExpectedWorkspaceRevision: current.WorkspaceRevision,
		ExpectedGroupRevision: current.Group.Revision, ExpectedConnectorRevision: current.Group.Destination.ConnectorRevision,
	})
	require.NoError(t, err)
	require.True(t, deleted.Deleted)
	require.Equal(t, "G2", deleted.GroupID)
	require.Equal(t, "L1", deleted.LineageID)
	require.Equal(t, 2, deleted.DeletedGenerations)
	var groups int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM write_groups WHERE id IN ('G1','G2')`).Scan(&groups))
	require.Zero(t, groups, "permanent deletion must remove both generations")
	for _, table := range []string{"write_group_versions", "write_group_members", "write_group_revision_ends", "write_group_basic_keys", "write_group_profile_keys", "write_group_storage_events", "write_group_ownership_fences", "write_group_migration_maps", "write_group_generations", "wg_runtime_versions", "wg_delivery_samples", "wg_delivery_checkpoints", "wg_delivery_buckets", "wg_delivery_outbox", "wg_delivery_counter_baselines", "wg_counter_proof_associations", "wg_delivery_mixed_progress", "wg_delivery_mixed_outcomes"} {
		var count int
		require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+table+` WHERE group_id IN ('G1','G2')`).Scan(&count), table)
		require.Zero(t, count, table)
	}
	for _, query := range []string{
		`SELECT COUNT(*) FROM write_group_lineages WHERE lineage_id='L1'`,
		`SELECT COUNT(*) FROM write_group_cutovers WHERE source_group_id IN ('G1','G2')`,
		`SELECT COUNT(*) FROM wg_delivery_sample_handoffs WHERE source_group_id IN ('G1','G2') OR target_group_id IN ('G1','G2')`,
		`SELECT COUNT(*) FROM wg_delivery_receipts WHERE effect_key='Qcommitted'`,
		`SELECT COUNT(*) FROM managed_schema_operations WHERE operation_id='Odone'`,
		`SELECT COUNT(*) FROM managed_schema_operations WHERE operation_id='Ofailed'`,
	} {
		var count int
		require.NoError(t, db.QueryRowContext(ctx, query).Scan(&count))
		require.Zero(t, count, query)
	}
	var destinationRows int
	require.NoError(t, destination.QueryRowContext(ctx, `SELECT COUNT(*) FROM public.gw_long_a`).Scan(&destinationRows))
	require.Equal(t, 1200, destinationRows)
	fixture := writeGroupFixture{workspaceID: predecessor.WorkspaceID, deviceID: predecessor.Members[0].DeviceID, pointID: predecessor.Members[0].PointID, tagID: predecessor.Members[0].TagID, connectorID: predecessor.Destination.ConnectorID}
	workspace, err := service.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	recreated, err := service.EnsureBasicManaged(ctx, fixture.deviceID, WriteGroupMutation{WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: workspace.DatabaseSetupRevision, ExpectedConnectorRevision: "connector-revision-1", Group: newBasicManagedDraft(fixture)})
	require.NoError(t, err)
	require.NotEmpty(t, recreated.Group.ID)
}

func TestPermanentDeleteRefusesPendingDelivery(t *testing.T) {
	for _, pending := range []string{"samples", "pending", "retrying", "unknown", "retrying_and_unknown", "blocked", "quarantined", "sending"} {
		t.Run(pending, func(t *testing.T) {
			db, service, group := lifecyclePreservationFixture(t)
			seedLifecycleAcceptedData(t, db, group)
			_, err := db.ExecContext(t.Context(), `DELETE FROM managed_schema_operations`)
			require.NoError(t, err)
			_, err = db.ExecContext(t.Context(), `UPDATE wg_delivery_samples SET consumed=1`)
			require.NoError(t, err)
			_, err = db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET state='sql_committed'`)
			require.NoError(t, err)
			if pending == "samples" {
				_, err = db.ExecContext(t.Context(), `UPDATE wg_delivery_samples SET consumed=0`)
				require.NoError(t, err)
				for _, id := range []string{"S2", "S3"} {
					_, err = db.ExecContext(t.Context(), `INSERT INTO wg_delivery_samples(workspace_id,group_id,group_revision,sample_id,member_key,observed_at,bucket_start,payload,payload_digest,consumed,created_at)
						SELECT workspace_id,group_id,group_revision,?,member_key,observed_at,bucket_start,payload,payload_digest,0,created_at FROM wg_delivery_samples WHERE sample_id='S1'`, id)
					require.NoError(t, err)
				}
			} else if pending == "retrying_and_unknown" {
				_, err = db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET state='retrying' WHERE effect_key='Qretry'`)
				require.NoError(t, err)
				_, err = db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET state='unknown' WHERE effect_key='Qunknown'`)
				require.NoError(t, err)
			} else {
				_, err = db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET state=? WHERE effect_key='Q1'`, pending)
				require.NoError(t, err)
			}
			service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
			disabled, err := service.Disable(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
			require.NoError(t, err)
			before := lifecycleReadState(t, db)
			_, err = service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, disabled.Group))
			require.ErrorContains(t, err, "WRITE_GROUP_DELETE_BACKLOG")
			var refusal *WriteGroupDeleteError
			require.ErrorAs(t, err, &refusal)
			require.Equal(t, "WRITE_GROUP_DELETE_BACKLOG", refusal.Code)
			backlog, ok := refusal.Details.(WriteGroupDeleteBacklog)
			require.True(t, ok)
			if pending == "samples" {
				require.Equal(t, 3, backlog.SamplesUnconsumed)
				require.Empty(t, backlog.Outbox)
			} else if pending == "retrying_and_unknown" {
				require.Zero(t, backlog.SamplesUnconsumed)
				require.Equal(t, map[string]int{"retrying": 1, "unknown": 1}, backlog.Outbox)
			} else {
				require.Zero(t, backlog.SamplesUnconsumed)
				require.Equal(t, map[string]int{pending: 1}, backlog.Outbox)
			}
			require.Equal(t, before, lifecycleReadState(t, db))
		})
	}
}

func TestPermanentDeleteRefusesActiveIntake(t *testing.T) {
	db, service, group := lifecyclePreservationFixture(t)
	service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
	before := lifecycleReadState(t, db)
	_, err := service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
	require.ErrorContains(t, err, "WRITE_GROUP_DELETE_ACTIVE")
	require.Equal(t, before, lifecycleReadState(t, db))
}

func TestPermanentDeleteAllowsNotYetActiveIntake(t *testing.T) {
	db, service, group := lifecyclePreservationFixture(t)
	service.WithClock(func() time.Time { return time.Date(2026, 10, 2, 12, 0, 7, 0, time.UTC) })
	service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
	intake, err := service.IntakeEligibility(t.Context(), group.ID)
	require.NoError(t, err)
	require.False(t, intake.Accepting)
	_, err = service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
	require.NoError(t, err)
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM write_groups WHERE id=?`, group.ID).Scan(&count))
	require.Zero(t, count)
}

func TestPermanentDeleteClearsExistingTombstone(t *testing.T) {
	db, service, group := lifecyclePreservationFixture(t)
	_, err := db.ExecContext(t.Context(), `UPDATE write_groups SET status='deleted' WHERE id=?`, group.ID)
	require.NoError(t, err)
	service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
	_, err = service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
	require.NoError(t, err)
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM write_groups WHERE id=?`, group.ID).Scan(&count))
	require.Zero(t, count)
}

func TestPermanentDeleteRefusesReservedSuccessor(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	h.db.SetMaxOpenConns(1)
	first, err := ensureLineageKind(h, "basic", h.mutation(t.Context(), lineageTestDraft(h, "basic")))
	require.NoError(t, err)
	lineage := requireLineageGenesis(t, h.db, first.Group)
	seedLineageGeneration(t, h, first.Group, lineage, "G2", 2)
	_, err = h.db.ExecContext(t.Context(), `UPDATE write_group_lineages SET reserved_group_id='G2',reservation_operation_id='O2',reservation_intent_digest='intent-2' WHERE lineage_id=?`, lineage)
	require.NoError(t, err)
	before := lineageReadSnapshot(t, h)
	_, err = h.service.Delete(t.Context(), first.Group.ID, h.mutationFor(first))
	require.ErrorContains(t, err, "WRITE_GROUP_DELETE_OPERATION_PENDING")
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

func TestPermanentDeleteRefusesPredecessorWithSuccessor(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	h.db.SetMaxOpenConns(1)
	first, err := ensureLineageKind(h, "basic", h.mutation(t.Context(), lineageTestDraft(h, "basic")))
	require.NoError(t, err)
	lineage := requireLineageGenesis(t, h.db, first.Group)
	seedLineageGeneration(t, h, first.Group, lineage, "G2", 2)
	_, err = h.db.ExecContext(t.Context(), `UPDATE write_groups SET status='deleted' WHERE id=?`, first.Group.ID)
	require.NoError(t, err)
	_, err = h.db.ExecContext(t.Context(), `UPDATE write_groups SET status='disabled' WHERE id='G2'`)
	require.NoError(t, err)
	_, err = h.db.ExecContext(t.Context(), `UPDATE write_group_lineages SET current_group_id='G2',revision=2 WHERE lineage_id=?`, lineage)
	require.NoError(t, err)
	before := lineageReadSnapshot(t, h)
	_, err = h.service.Delete(t.Context(), first.Group.ID, h.mutationFor(first))
	require.ErrorContains(t, err, "WRITE_GROUP_DELETE_HAS_SUCCESSOR")
	require.Equal(t, before, lineageReadSnapshot(t, h))
}
