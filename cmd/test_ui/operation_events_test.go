package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/diagnostics"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestProductionOperationEventsCutoverUsesCommittedOriginalIdentity(t *testing.T) {
	f := newOperationEventFixture(t)
	cutover, err := f.env.services.writeGroups.ApplyCutover(t.Context(), f.request)
	require.NoError(t, err)
	require.NotEmpty(t, cutover.ID)
	persisted, err := f.env.services.writeGroups.GetCutover(t.Context(), f.group.ID, cutover.ID)
	require.NoError(t, err)
	require.Equal(t, cutover, persisted, "native durable operation precedes diagnostics assertion")
	record, err := f.env.services.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	require.NotEqual(t, f.request.ExpectedWorkspaceRevision, record.DatabaseSetupRevision)
	events := deliveryEventsOfCode(t, f.broker, "domain.cutover.applied")
	require.Len(t, events, 1, "actual committed cutover must reach production diagnostics")
	ref := events[0].Fields
	require.Equal(t, f.group.ID, ref["group_id"])
	require.Equal(t, f.group.AppliedRevision, ref["entity_revision"])
	require.Equal(t, record.DatabaseSetupRevision, ref["workspace_revision"])
	require.Equal(t, cutover.ID, ref["reference_id"])
	require.Equal(t, string(runtimeview.StageAccepted), ref["stage"])
}

func TestProductionOperationEventsCutoverReplayDoesNotEmitNewResult(t *testing.T) {
	f := newOperationEventFixture(t)
	original, err := f.env.services.writeGroups.ApplyCutover(t.Context(), f.request)
	require.NoError(t, err)
	require.Len(t, deliveryEventsOfCode(t, f.broker, "domain.cutover.applied"), 1)
	var before int
	require.NoError(t, f.env.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM write_group_cutovers`).Scan(&before))
	replay, err := f.env.services.writeGroups.ApplyCutover(t.Context(), f.request)
	require.NoError(t, err)
	require.Equal(t, original.ID, replay.ID)
	var after int
	require.NoError(t, f.env.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM write_group_cutovers`).Scan(&after))
	require.Equal(t, before, after)
	require.Len(t, deliveryEventsOfCode(t, f.broker, "domain.cutover.applied"), 1)
}

func TestProductionOperationEventsPreledgerRejectionIsUnavailable(t *testing.T) {
	f := newOperationEventFixture(t)
	f.request.ReviewDigest = "stale-digest"
	result, err := f.env.services.writeGroups.ApplyCutover(t.Context(), f.request)
	require.Error(t, err)
	require.Nil(t, result)
	var count int
	require.NoError(t, f.env.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM write_group_cutovers`).Scan(&count))
	require.Zero(t, count)
	require.Empty(t, deliveryEventsOfCode(t, f.broker, "domain.cutover.applied"))
	require.Empty(t, deliveryEventsOfCode(t, f.broker, "domain.cutover.rejected"), "no ledger ID must not become invented attributable operation")
	events := deliveryEventsOfCode(t, f.broker, "domain.unavailable")
	require.Len(t, events, 1)
	require.Empty(t, events[0].Fields)
}

func TestProductionOperationEventsMemoryFileAndSubscriberRedaction(t *testing.T) {
	f := newOperationEventFixture(t, true)
	sub, err := f.broker.Subscribe(diagnostics.Query{}, "")
	require.NoError(t, err)
	defer sub.Close()
	result, err := f.env.services.writeGroups.ApplyCutover(t.Context(), f.request)
	require.NoError(t, err)
	expected := deliveryEventsOfCode(t, f.broker, "domain.cutover.applied")
	require.Len(t, expected, 1)
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	live, err := sub.Next(ctx)
	require.NoError(t, err)
	sub.Ack(live.Cursor())
	require.Equal(t, expected[0], live)
	require.Equal(t, result.ID, live.Fields["reference_id"])
	raw, err := json.Marshal(collectorEventRecords(t, f.broker))
	require.NoError(t, err)
	liveRaw, err := json.Marshal(live)
	require.NoError(t, err)
	file, err := os.ReadFile(f.logPath)
	require.NoError(t, err)
	require.Contains(t, string(file), `"code":"domain.cutover.applied"`)
	for _, secret := range []string{"password_fixture", "private-host", "postgres://", "connection_config", "CREATE TABLE", "DROP TABLE", f.env.targets["connector-A"]} {
		require.NotContains(t, string(raw), secret)
		require.NotContains(t, string(liveRaw), secret)
		require.NotContains(t, string(file), secret)
	}
}

func TestProductionOperationEventsRepairSchemaIsNotHistoricalDataRecovery(t *testing.T) {
	f := newRepairOperationFixture(t)
	f.start(t)
	f.waitLongCommitted(t, 2)
	require.NoError(t, f.env.services.groupPipe.Stop(3*time.Second))
	var oldIDs []string
	rows, err := f.targetDB.QueryContext(t.Context(), `SELECT record_id FROM "`+f.long.Destination.TableName+`"`)
	require.NoError(t, err)
	for rows.Next() {
		var id string
		require.NoError(t, rows.Scan(&id))
		oldIDs = append(oldIDs, id)
	}
	require.NoError(t, rows.Close())
	require.GreaterOrEqual(t, len(oldIDs), 2, "retain every actual historical row, including legitimate native timer outputs")
	nextBucket := time.Now().UTC().Truncate(10 * time.Second).Add(20 * time.Second)
	f.clock.set(nextBucket.Add(time.Second))
	f.seed = f.env.pipeline(f.clock, "operation-next-window")
	require.NoError(t, f.seed.Reconcile(t.Context()))
	f.env.feed(t, f.seed, f.clock, []*workspace.WriteGroup{f.long}, nextBucket, 22.5, 200)
	pendingItems := f.pendingForBucket(t, nextBucket)
	pending := pendingItems[0]
	// A later saved draft is not the original frozen repair identity.
	saved, err := f.env.services.writeGroups.Get(t.Context(), f.long.ID)
	require.NoError(t, err)
	saved.Group.Name = "later password_fixture draft"
	updated, err := f.env.services.writeGroups.Update(t.Context(), f.long.ID, workspace.WriteGroupMutation{WorkspaceID: f.long.WorkspaceID, ExpectedWorkspaceRevision: saved.WorkspaceRevision, ExpectedGroupRevision: saved.Group.Revision, ExpectedConnectorRevision: saved.Group.Destination.ConnectorRevision, Group: saved.Group})
	require.NoError(t, err)
	require.NotEqual(t, pending.Key.GroupRevision, updated.Group.Revision)
	f.dropTable(t)
	f.start(t)
	for _, item := range pendingItems {
		require.Eventually(t, func() bool {
			current, err := groupdelivery.NewStore(f.env.db).GetOutbox(t.Context(), item.EffectKey)
			require.NoError(t, err)
			return current.State == groupdelivery.StateCommitted
		}, 8*time.Second, 20*time.Millisecond)
	}
	require.NoError(t, f.env.services.groupPipe.Stop(3*time.Second))
	events := waitDeliveryEvents(t, f.broker, "domain.repair.completed")
	require.Len(t, events, 1)
	fields := events[0].Fields
	require.Equal(t, f.long.ID, fields["group_id"])
	require.Equal(t, pending.Key.GroupRevision, fields["entity_revision"])
	require.Equal(t, pending.Destination.ConnectorRevision, fields["connector_revision"])
	id, ok := fields["reference_id"].(string)
	require.True(t, ok)
	op, err := f.env.services.recordingPlan.GetSchemaOperation(t.Context(), f.long.WorkspaceID, id)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, op.Status)
	require.NotEqual(t, pending.EffectKey, id, "real repair ledger ID is not a delivery effect or observation UUID")
	require.Contains(t, events[0].Message, "data recovery is not implied")
	var count int
	require.NoError(t, f.targetDB.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM "`+f.long.Destination.TableName+`"`).Scan(&count))
	require.GreaterOrEqual(t, count, 2, "both exact newly journaled effects reached the recreated table; timer outputs are independent")
	for _, old := range oldIDs {
		require.NoError(t, f.targetDB.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM "`+f.long.Destination.TableName+`" WHERE record_id=?`, old).Scan(&count))
		require.Zero(t, count)
	}
}

func TestProductionOperationEventsClosedOrFullCapturePreservesCutover(t *testing.T) {
	for _, mode := range []string{"full", "closed"} {
		t.Run(mode, func(t *testing.T) {
			f := newOperationEventFixture(t)
			var dropped uint64
			if mode == "closed" {
				require.NoError(t, f.broker.Close(t.Context()))
			} else {
				for i := range diagnostics.AdmissionRecords {
					key := groupdelivery.GroupKey{WorkspaceID: f.group.WorkspaceID, GroupID: fmt.Sprintf("retained-%d", i), GroupRevision: "R1"}
					ref := deliveryEventRef(key, f.record.DatabaseSetupRevision, "connector-A", "connector-1")
					ref.Code, ref.Stage, ref.Reason, ref.Action = "domain.repair.blocked", runtimeview.StageBlocked, runtimeview.ReasonUnavailable, runtimeview.ActionRefresh
					ref.Evidence = diagnostics.DomainEvidenceRef{Kind: "operation-result", ReferenceID: fmt.Sprintf("old-operation-%d", i)}
					emitDeliveryEvent(f.broker, ref)
				}
				require.NoError(t, f.broker.Flush(t.Context()))
				view, err := f.broker.Snapshot(diagnostics.Query{})
				require.NoError(t, err)
				dropped = view.CaptureDropped
			}
			result, err := f.env.services.writeGroups.ApplyCutover(t.Context(), f.request)
			require.NoError(t, err)
			persisted, err := f.env.services.writeGroups.GetCutover(t.Context(), f.group.ID, result.ID)
			require.NoError(t, err)
			require.Equal(t, result.ID, persisted.ID)
			if mode == "full" {
				require.NoError(t, f.broker.Flush(t.Context()))
				view, err := f.broker.Snapshot(diagnostics.Query{})
				require.NoError(t, err)
				require.Greater(t, view.CaptureDropped, dropped)
			}
		})
	}
}

func TestOperationEventObserverRejectsUnprovenOrContradictoryIdentity(t *testing.T) {
	for _, name := range []string{"missing-operation", "secret-operation", "unknown-outcome", "raw-reason", "completed-with-fault", "unavailable", "missing-workspace"} {
		t.Run(name, func(t *testing.T) {
			broker, err := diagnostics.New(diagnostics.Options{})
			require.NoError(t, err)
			defer broker.Close(context.Background())
			storage, cutover := newOperationEventObservers(broker, nil)
			o := grouppipeline.StorageOperationObservation{Available: true, Key: groupdelivery.GroupKey{WorkspaceID: "W", GroupID: "G", GroupRevision: "R1"}, WorkspaceRevision: "W1", ConnectorID: "C", ConnectorRevision: "C1", OperationID: "ledger-1", Outcome: grouppipeline.StorageObservationCompleted, Reason: "none"}
			switch name {
			case "missing-operation":
				o.OperationID = ""
			case "secret-operation":
				o.OperationID = "Server=private-host;Password=password_fixture"
			case "unknown-outcome":
				o.Outcome = "unknown"
			case "raw-reason":
				o.Outcome = grouppipeline.StorageObservationBlocked
				o.Reason = "driver password_fixture"
			case "completed-with-fault":
				o.Reason = "unavailable"
			case "unavailable":
				o.Available = false
			case "missing-workspace":
				o.WorkspaceRevision = ""
			}
			storage.Observe(o)
			cutover.Observe(t.Context(), workspace.CutoverObservation{Available: true, WorkspaceID: "W", WorkspaceRevision: "W1", GroupID: "G", GroupRevision: "R1", ConnectorID: "C", ConnectorRevision: "C1", OperationID: "", Outcome: "applied"})
			records := collectorEventRecords(t, broker)
			require.Len(t, records, 1, "unknown observations share existing bounded unavailable capture")
			require.Equal(t, "domain.unavailable", records[0].Code)
			require.Empty(t, records[0].Fields)
		})
	}
}

func TestOperationEventObserverNilFilteredAndPureScope(t *testing.T) {
	storage, cutover := newOperationEventObservers(nil, nil)
	require.Nil(t, storage.Observe)
	require.Nil(t, cutover.Observe)
	for _, level := range []string{"info", "warn", "error"} {
		t.Run(level, func(t *testing.T) {
			broker, err := diagnostics.New(diagnostics.Options{Level: level})
			require.NoError(t, err)
			defer broker.Close(context.Background())
			calls := 0
			cfg, c := newOperationEventObservers(broker, func(context.Context, groupdelivery.GroupKey) (string, bool) { calls++; return "W1", true })
			if level == "error" {
				require.Nil(t, cfg.Observe)
				require.Nil(t, c.Observe)
				require.Nil(t, cfg.ScopeRevision)
				require.Nil(t, c.ScopeRevision)
				require.Zero(t, cfg.RecordLimit)
			} else {
				require.Equal(t, level == "info", cfg.InfoEnabled)
				require.True(t, cfg.WarnEnabled)
				require.Equal(t, diagnostics.AdmissionRecords, cfg.RecordLimit)
				require.NotNil(t, cfg.OnDrop)
				require.NotNil(t, c.OnDrop)
			}
			require.Zero(t, calls, "config creation performs no metadata reads")
			view, err := broker.Snapshot(diagnostics.Query{})
			require.NoError(t, err)
			require.Zero(t, view.CaptureDropped)
		})
	}
	f := newOperationEventFixture(t)
	delivery, _ := newDeliveryEventObservers(f.broker, f.env.services.workspace)
	cfg, c := newOperationEventObservers(f.broker, delivery.ScopeRevision)
	var before, after int
	require.NoError(t, f.env.db.QueryRowContext(t.Context(), `SELECT total_changes()`).Scan(&before))
	current, err := f.env.services.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	rev, ok := cfg.ScopeRevision(t.Context(), f.item.Key)
	require.True(t, ok)
	require.Equal(t, current.DatabaseSetupRevision, rev)
	rev, ok = c.ScopeRevision(t.Context(), f.group.WorkspaceID, f.group.ID, f.group.AppliedRevision, "connector-A", "connector-1")
	require.True(t, ok)
	require.Equal(t, current.DatabaseSetupRevision, rev)
	_, ok = c.ScopeRevision(t.Context(), "foreign-workspace", f.group.ID, f.group.AppliedRevision, "connector-A", "connector-1")
	require.False(t, ok)
	require.NoError(t, f.env.db.QueryRowContext(t.Context(), `SELECT total_changes()`).Scan(&after))
	require.Equal(t, before, after)
}

func TestProductionOperationEventsRepairFailureAndUncertainty(t *testing.T) {
	for _, mode := range []string{"rollback", "unacknowledged"} {
		t.Run(mode, func(t *testing.T) {
			f := newRepairOperationFixture(t)
			delivery, _ := newDeliveryEventObservers(f.broker, f.env.services.workspace)
			cfg, _ := newOperationEventObservers(f.broker, delivery.ScopeRevision)
			decision, err := f.restoreWithFault(t, cfg, mode)
			if mode == "rollback" {
				require.NoError(t, err)
				require.Equal(t, groupdelivery.RestoreBlock, decision.Action)
				events := deliveryEventsOfCode(t, f.broker, "domain.repair.blocked")
				require.Len(t, events, 1)
				id, ok := events[0].Fields["reference_id"].(string)
				require.True(t, ok)
				op, err := f.env.services.recordingPlan.GetSchemaOperation(t.Context(), f.long.WorkspaceID, id)
				require.NoError(t, err)
				require.Equal(t, recordingplan.SchemaOperationFailed, op.Status)
			} else {
				require.Error(t, err)
				events := deliveryEventsOfCode(t, f.broker, "domain.unavailable")
				require.Len(t, events, 1)
				require.Empty(t, events[0].Fields)
			}
			require.Empty(t, deliveryEventsOfCode(t, f.broker, "domain.repair.completed"))
			raw, err := json.Marshal(collectorEventRecords(t, f.broker))
			require.NoError(t, err)
			require.NotContains(t, string(raw), "password_fixture")
			require.NotContains(t, string(raw), "private-host")
		})
	}
}

func TestProductionOperationEventsObserverPanicPreservesNativeRepair(t *testing.T) {
	f := newRepairOperationFixture(t)
	delivery, _ := newDeliveryEventObservers(f.broker, f.env.services.workspace)
	cfg, _ := newOperationEventObservers(f.broker, delivery.ScopeRevision)
	called := false
	cfg.Observe = func(grouppipeline.StorageOperationObservation) { called = true; panic("password_fixture") }
	decision, err := f.restoreWithFault(t, cfg, "")
	require.NoError(t, err)
	require.Equal(t, groupdelivery.RestoreRetryNow, decision.Action)
	require.True(t, called)
	var succeeded int
	require.NoError(t, f.env.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM managed_schema_operations WHERE action='storage_repair' AND status='succeeded'`).Scan(&succeeded))
	require.Equal(t, 1, succeeded)
}

func TestProductionOperationEventsFilteredOrPanickingCutoverCaptureIsInert(t *testing.T) {
	for _, mode := range []string{"warn-filter", "error-filter", "panic"} {
		t.Run(mode, func(t *testing.T) {
			f := newOperationEventFixture(t)
			level := "info"
			if mode == "warn-filter" {
				level = "warn"
			}
			if mode == "error-filter" {
				level = "error"
			}
			broker, err := diagnostics.New(diagnostics.Options{Level: level})
			require.NoError(t, err)
			defer broker.Close(context.Background())
			delivery, _ := newDeliveryEventObservers(broker, f.env.services.workspace)
			_, cfg := newOperationEventObservers(broker, delivery.ScopeRevision)
			reads := 0
			if cfg.ScopeRevision != nil {
				original := cfg.ScopeRevision
				cfg.ScopeRevision = func(ctx context.Context, w, g, r, c, cr string) (string, bool) {
					reads++
					return original(ctx, w, g, r, c, cr)
				}
			}
			called := false
			if mode == "panic" {
				cfg.Observe = func(context.Context, workspace.CutoverObservation) { called = true; panic("password_fixture") }
			} else {
				cfg.RecordLimit = 0
			}
			svc := workspace.NewWriteGroupService(f.env.services.workspace, workspace.NewSQLWriteGroupRepository(f.env.db), cfg).
				WithClock(f.clock.now).WithBacklogOwnershipGuard(grouppipeline.BacklogGuard{}).
				WithTableInspector(dbtarget.NewReadOnlyTableInspector(f.env.services.dbTarget)).
				WithManagedTableInspector(dbtarget.NewManagedTableInspector(f.env.services.dbTarget, f.env.db))
			result, err := svc.ApplyCutover(t.Context(), f.request)
			require.NoError(t, err)
			require.NotEmpty(t, result.ID)
			if mode == "panic" {
				require.True(t, called)
				require.Equal(t, 2, reads)
			} else {
				require.Zero(t, reads)
			}
			view, err := broker.Snapshot(diagnostics.Query{})
			require.NoError(t, err)
			if mode != "panic" {
				require.Zero(t, view.CaptureDropped, "level filtering before cap is not capture loss")
			}
		})
	}
}
