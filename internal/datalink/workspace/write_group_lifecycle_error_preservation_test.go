package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

func TestLifecycleStopAndRemovePreserveAcceptedData(t *testing.T) {
	for _, operation := range []string{"stop", "remove"} {
		connectors := []string{"unchanged"}
		if operation == "stop" {
			connectors = append(connectors, "changed", "removed")
		}
		for _, connector := range connectors {
			t.Run(operation+"/"+connector, func(t *testing.T) {
				db, service, group := lifecyclePreservationFixture(t)
				store := seedLifecycleAcceptedData(t, db, group)
				// This service regression isolates stop/remove mutation from the guard;
				// grouppipeline's integration tests exercise the actual ownership proof.
				service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
				if connector == "changed" {
					_, err := db.ExecContext(t.Context(), `UPDATE database_connectors SET identity_revision='C2', connection_config='{"database":"different"}' WHERE id=?`, group.Destination.ConnectorID)
					require.NoError(t, err)
				}
				if connector == "removed" {
					_, err := db.ExecContext(t.Context(), `DELETE FROM database_connectors WHERE id=?`, group.Destination.ConnectorID)
					require.NoError(t, err)
				}
				require.True(t, mustIntakeEligibility(t, service, group.ID).Accepting)
				before := lifecyclePreservedData(t, db)
				q1, err := store.GetOutbox(t.Context(), "Q1")
				require.NoError(t, err)
				require.Equal(t, "G1", q1.Key.GroupID)
				require.Equal(t, "R1", q1.Key.GroupRevision)
				mutation := lifecycleReadMutation(t, service, group)
				var saved *WriteGroupSaveResult
				if operation == "stop" {
					saved, err = service.Disable(t.Context(), group.ID, mutation)
				} else {
					saved, err = service.Disable(t.Context(), group.ID, mutation)
					require.NoError(t, err)
					beforeDelete := lifecycleReadState(t, db)
					_, deleteErr := service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, saved.Group))
					require.ErrorContains(t, deleteErr, "WRITE_GROUP_DELETE_BACKLOG")
					var refusal *WriteGroupDeleteError
					require.ErrorAs(t, deleteErr, &refusal)
					counts := refusal.Details.(WriteGroupDeleteBacklog)
					require.Equal(t, 1, counts.SamplesUnconsumed)
					require.Equal(t, map[string]int{"pending": 1, "retrying": 1, "blocked": 1, "unknown": 1, "sending": 1}, counts.Outbox)
					require.Equal(t, beforeDelete, lifecycleReadState(t, db))
				}
				require.NoError(t, err)
				require.Equal(t, WriteGroupStatusDisabled, saved.Group.Status)
				require.Equal(t, "R1", saved.Group.AppliedRevision)
				require.Equal(t, group.Destination, saved.Group.Destination)
				require.False(t, mustIntakeEligibility(t, service, group.ID).Accepting)
				require.Equal(t, before, lifecyclePreservedData(t, db), "accepted samples, pending states, versions, receipts, checkpoints and operations remain byte-for-byte unchanged")
				after, err := store.GetOutbox(t.Context(), "Q1")
				require.NoError(t, err)
				require.Equal(t, q1, after, "Q1 keeps its exact R1 / D1 binding and payload")
				historical, err := service.ResolveAppliedAt(t.Context(), group.ID, time.Date(2026, 10, 2, 12, 0, 16, 0, time.UTC))
				require.NoError(t, err)
				require.Equal(t, "R1", historical.AppliedRevision)
				require.Equal(t, group.Destination, historical.Group.Destination)
				reloaded, err := service.Get(t.Context(), group.ID)
				require.NoError(t, err)
				require.Equal(t, saved.Group, reloaded.Group, "the stopped group remains readable")
				for _, effect := range []string{"Q1", "Qretry"} {
					_, err := store.BeginDelivery(t.Context(), effect, "existing-worker", time.Minute)
					require.NoError(t, err, "stop/remove must not cancel accepted delivery")
					claimed, err := store.GetOutbox(t.Context(), effect)
					require.NoError(t, err)
					require.Equal(t, q1.Destination, claimed.Destination)
					require.Equal(t, "R1", claimed.Key.GroupRevision)
				}
			})
		}
	}
}

func TestLifecycleBacklogRemovalPreservesAllAcceptedData(t *testing.T) {
	db, service, group := lifecyclePreservationFixture(t)
	seedLifecycleAcceptedData(t, db, group)
	service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
	stopped, err := service.Disable(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
	require.NoError(t, err)
	before := lifecyclePreservedData(t, db)
	config := lifecycleReadState(t, db)
	_, err = service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, stopped.Group))
	require.ErrorContains(t, err, "WRITE_GROUP_DELETE_BACKLOG")
	var refusal *WriteGroupDeleteError
	require.ErrorAs(t, err, &refusal)
	require.Equal(t, 1, refusal.Details.(WriteGroupDeleteBacklog).SamplesUnconsumed)
	require.Equal(t, before, lifecyclePreservedData(t, db))
	require.Equal(t, config, lifecycleReadState(t, db))
	require.False(t, mustIntakeEligibility(t, service, group.ID).Accepting, "failed remove retains the explicit stop")
}

func lifecyclePreservationFixture(t *testing.T) (*sql.DB, *WriteGroupService, *WriteGroup) {
	t.Helper()
	db := openWorkspaceTestDB(t, ":memory:")
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	db.SetMaxOpenConns(1)
	fixture := seedWriteGroupFixture(t.Context(), t, db)
	_, err := db.ExecContext(t.Context(), `UPDATE database_connectors SET kind='postgres',connection_config='{"database":"persisted-db","schema":"persisted-schema"}' WHERE id=?`, fixture.connectorID)
	require.NoError(t, err)
	repo := NewSQLWriteGroupRepository(db)
	sequence := 0
	repo.newID = func() string {
		sequence++
		if sequence == 1 {
			return "G1"
		}
		return fmt.Sprintf("R%d", sequence-1)
	}
	service := NewWriteGroupService(NewService(NewSQLRepository(db)), repo).WithTableInspector(writeGroupLifecycleInspector{})
	now := time.Date(2026, 10, 2, 12, 0, 7, 0, time.UTC)
	service.WithClock(func() time.Time { return now })
	record, err := service.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	candidate := newBasicWriteGroup(fixture)
	candidate.RowPolicy.IntervalSeconds = 15
	created, err := service.Create(t.Context(), WriteGroupMutation{WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, ExpectedConnectorRevision: "connector-revision-1", Group: candidate})
	require.NoError(t, err)
	applied, err := service.Apply(t.Context(), created.Group.ID, lifecycleReadMutation(t, service, created.Group))
	require.NoError(t, err)
	now = time.Date(2026, 10, 2, 12, 0, 16, 0, time.UTC)
	require.Equal(t, "G1", applied.Group.ID)
	require.Equal(t, "R1", applied.Group.AppliedRevision)
	return db, service, applied.Group
}

func seedLifecycleAcceptedData(t *testing.T, db *sql.DB, group *WriteGroup) *groupdelivery.Store {
	t.Helper()
	store := groupdelivery.NewStore(db)
	key := groupdelivery.GroupKey{WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.AppliedRevision}
	d := group.Destination
	scope, err := snapshot.DestinationScope(d.ConnectorID, d.ConnectorRevision, d.Database, d.TableSchema, d.TableName)
	require.NoError(t, err)
	// D1 is the exact frozen destination represented by these persisted fields.
	destination := groupdelivery.FrozenDestination{Scope: scope, ConnectorID: d.ConnectorID, ConnectorRevision: d.ConnectorRevision, Database: d.Database, TableSchema: d.TableSchema, TableName: d.TableName, DedupeCapability: groupdelivery.DedupeNone}
	start := time.Date(2026, 10, 2, 12, 0, 15, 0, time.UTC)
	buckets := make([]groupdelivery.ClosedBucket, 0, 6)
	for i, effect := range []string{"Q1", "Qretry", "Qblocked", "Qunknown", "Qsending", "Qcommitted"} {
		bucket := start.Add(time.Duration(i) * 15 * time.Second)
		row := &dbtarget.EncodedRow{RecordID: effect, EffectKey: effect, BucketStart: bucket, Cells: []dbtarget.EncodedCell{{Column: "temperature", Value: 21.5}}}
		buckets = append(buckets, groupdelivery.ClosedBucket{Outcome: snapshot.Outcome{Kind: snapshot.OutcomeRow, RecordID: effect, EffectKey: effect, BucketStart: bucket}, Row: row})
	}
	next := start.Add(90 * time.Second)
	require.NoError(t, store.CommitClosure(t.Context(), key, groupdelivery.Closure{Destination: destination, NextClose: next, Buckets: buckets}))
	for effect, state := range map[string]string{"Qretry": "retrying", "Qblocked": "blocked", "Qunknown": "unknown", "Qsending": "sending", "Qcommitted": "sql_committed"} {
		_, err := db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET state=?,retry_count=2,last_error_code='preserved-code' WHERE effect_key=?`, state, effect)
		require.NoError(t, err)
	}
	_, err = db.ExecContext(t.Context(), `INSERT INTO wg_delivery_receipts(effect_key,payload_digest,committed_at) SELECT effect_key,payload_digest,'2026-10-02T12:02:00Z' FROM wg_delivery_outbox WHERE effect_key='Qcommitted'`)
	require.NoError(t, err)
	member := group.Members[0]
	memberJSON, err := json.Marshal([]string{member.DeviceID, member.PointID, member.TagID})
	require.NoError(t, err)
	require.NoError(t, store.AppendSample(t.Context(), key, next, snapshot.Sample{SampleID: "S1", MemberKey: string(memberJSON), ObservedAt: next.Add(time.Second), Quality: schema.QualityGood, Value: measurement.NewInt64(42), SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision}))
	for _, state := range []string{"pending", "succeeded"} {
		_, err := db.ExecContext(t.Context(), `INSERT INTO managed_schema_operations(operation_id,token,workspace_id,scope_key,owner,action,status,created_at,updated_at,detail) VALUES(?,?,?,?,?,'create',?,'2026-10-02T12:00:00Z','2026-10-02T12:00:00Z','preserved history')`, "O-"+state, "T-"+state, group.WorkspaceID, "scope-"+state, group.ID, state)
		require.NoError(t, err)
	}
	return store
}

func lifecyclePreservedData(t *testing.T, db *sql.DB) map[string]string {
	t.Helper()
	state := lifecycleReadState(t, db)
	delete(state, "system_settings")
	delete(state, "write_groups")
	delete(state, "database_connectors")
	for _, table := range []string{"managed_schema_operations", "wg_delivery_buckets", "wg_delivery_counter_baselines", "wg_counter_proof_associations", "wg_delivery_sample_handoffs", "write_group_storage_scopes", "write_group_storage_events"} {
		rows, err := db.QueryContext(t.Context(), `SELECT * FROM `+table+` ORDER BY rowid`)
		require.NoError(t, err)
		columns, err := rows.Columns()
		require.NoError(t, err)
		all := [][]any{}
		for rows.Next() {
			values, pointers := make([]any, len(columns)), make([]any, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}
			require.NoError(t, rows.Scan(pointers...))
			all = append(all, values)
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
		raw, err := json.Marshal(all)
		require.NoError(t, err)
		state[table] = string(raw)
	}
	return state
}
