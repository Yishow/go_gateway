package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/diagnostics"

	"github.com/stretchr/testify/require"
)

func TestProductionDeliveryEventsNewExactReceiptIsObserved(t *testing.T) {
	f := newDeliveryEventFixture(t)
	require.Empty(t, deliveryEventsOfCode(t, f.broker, "domain.delivery.first-confirmed"), "queued output is not a SQL commitment")
	f.start(t)
	item := f.waitCommitted(t)
	events := waitDeliveryEvents(t, f.broker, "domain.delivery.first-confirmed")
	require.Len(t, events, 1, "new actual destination/local receipt commitment must reach production diagnostics")
	fields := events[0].Fields
	require.Equal(t, f.record.ID, fields["workspace_id"])
	require.Equal(t, f.record.DatabaseSetupRevision, fields["workspace_revision"])
	require.Equal(t, item.Key.GroupID, fields["group_id"])
	require.Equal(t, item.Key.GroupRevision, fields["entity_revision"])
	require.Equal(t, item.Destination.ConnectorID, fields["connector_id"])
	require.Equal(t, item.Destination.ConnectorRevision, fields["connector_revision"])
	require.Equal(t, item.EffectKey, fields["reference_id"])
	require.Equal(t, item.EffectKey, fields["effect_key"])
	require.Equal(t, item.RecordID, fields["record_id"])
	require.Equal(t, item.PayloadDigest, fields["payload_digest"])
	require.Equal(t, "sql-committed", fields["stage"])
}

func TestProductionDeliveryEventsRetryThenExactReceiptRecovery(t *testing.T) {
	f := newDeliveryEventFixture(t)
	connection, err := f.destination.Conn(t.Context())
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = connection.ExecContext(context.Background(), "ROLLBACK")
		require.NoError(t, connection.Close())
	})
	_, err = connection.ExecContext(t.Context(), "BEGIN EXCLUSIVE")
	require.NoError(t, err)
	f.start(t)
	require.Eventually(t, func() bool { return f.state(t).State == groupdelivery.StateRetrying }, 8*time.Second, 20*time.Millisecond)
	retries := waitDeliveryEvents(t, f.broker, "domain.delivery.retrying")
	require.Len(t, retries, 1)
	require.Equal(t, f.item.EffectKey, retries[0].Fields["reference_id"])
	require.Empty(t, deliveryEventsOfCode(t, f.broker, "domain.delivery.first-confirmed"))
	_, err = connection.ExecContext(t.Context(), "ROLLBACK")
	require.NoError(t, err)
	f.waitCommitted(t)
	recovered := waitDeliveryEvents(t, f.broker, "domain.delivery.recovered")
	require.Len(t, recovered, 1)
	require.Equal(t, f.item.EffectKey, recovered[0].Fields["effect_key"])
	require.Len(t, waitDeliveryEvents(t, f.broker, "domain.delivery.first-confirmed"), 1)
}

func TestProductionDeliveryEventsReceiptConflictDoesNotClaimCommit(t *testing.T) {
	f := newDeliveryEventFixture(t)
	_, err := f.env.db.ExecContext(t.Context(), `INSERT INTO wg_delivery_receipts(effect_key,payload_digest,committed_at) VALUES(?,?,?)`, f.item.EffectKey, "wrong-digest", time.Now().UTC().Format(time.RFC3339Nano))
	require.NoError(t, err)
	f.start(t)
	require.Eventually(t, func() bool { return f.state(t).State == groupdelivery.StateCommitted }, 8*time.Second, 20*time.Millisecond, "native result remains committed despite diagnostic proof conflict")
	var count int
	require.NoError(t, f.destination.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM readings").Scan(&count))
	require.Equal(t, 1, count)
	require.Empty(t, deliveryEventsOfCode(t, f.broker, "domain.delivery.first-confirmed"))
	require.Empty(t, deliveryEventsOfCode(t, f.broker, "domain.delivery.recovered"))
	unknown := waitDeliveryEvents(t, f.broker, "domain.unavailable")
	require.Len(t, unknown, 1)
	require.Empty(t, unknown[0].Fields)
}

func TestProductionDeliveryEventsRawTriggerErrorIsOnlyTypedQuarantine(t *testing.T) {
	f := newDeliveryEventFixture(t, true)
	sub, err := f.broker.Subscribe(diagnostics.Query{}, "")
	require.NoError(t, err)
	defer sub.Close()
	_, err = f.destination.ExecContext(t.Context(), `CREATE TRIGGER refuse_raw BEFORE INSERT ON readings BEGIN SELECT RAISE(FAIL, 'password_fixture nested_secret private_host'); END`)
	require.NoError(t, err)
	f.start(t)
	require.Eventually(t, func() bool { return f.state(t).State == groupdelivery.StateQuarantined }, 8*time.Second, 20*time.Millisecond)
	events := waitDeliveryEvents(t, f.broker, "domain.delivery.quarantined")
	require.Len(t, events, 1)
	require.Equal(t, f.item.EffectKey, events[0].Fields["reference_id"])
	require.Empty(t, deliveryEventsOfCode(t, f.broker, "domain.delivery.first-confirmed"))
	raw, err := json.Marshal(collectorEventRecords(t, f.broker))
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	var live diagnostics.Event
	received := make([]diagnostics.Event, 0, 2)
	found := false
	for len(received) < diagnostics.LiveRecords {
		live, err = sub.Next(ctx)
		require.NoError(t, err)
		received = append(received, live)
		sub.Ack(live.Cursor())
		if live.Code == "domain.delivery.quarantined" {
			found = true
			break
		}
	}
	require.True(t, found, "typed quarantine must arrive within the bounded live sequence")
	require.Equal(t, events[0], live)
	liveBytes, err := json.Marshal(received)
	require.NoError(t, err)
	fileBytes, err := os.ReadFile(f.logPath)
	require.NoError(t, err)
	require.Contains(t, string(fileBytes), `"code":"domain.delivery.quarantined"`)
	for _, forbidden := range []string{"password_fixture", "nested_secret", "private_host", "connection_config", "raw_error"} {
		require.NotContains(t, string(raw), forbidden)
		require.NotContains(t, string(fileBytes), forbidden)
		require.NotContains(t, string(liveBytes), forbidden)
	}
}

func TestProductionDeliveryEventsFrozenOutputDoesNotBorrowLatestDraft(t *testing.T) {
	f := newDeliveryEventFixture(t)
	current, err := f.env.services.writeGroups.Get(t.Context(), f.group.ID)
	require.NoError(t, err)
	current.Group.Name = "New editable draft"
	current.Group.Destination.ConnectorID = "connector-B"
	saved, err := f.env.services.writeGroups.Update(t.Context(), f.group.ID, workspace.WriteGroupMutation{WorkspaceID: f.record.ID, ExpectedWorkspaceRevision: f.record.DatabaseSetupRevision, ExpectedGroupRevision: f.group.Revision, ExpectedConnectorRevision: "connector-1", Group: current.Group})
	require.NoError(t, err)
	require.NotEqual(t, f.item.Key.GroupRevision, saved.Group.Revision)
	f.start(t)
	item := f.waitCommitted(t)
	events := waitDeliveryEvents(t, f.broker, "domain.delivery.first-confirmed")
	require.Len(t, events, 1)
	require.Equal(t, item.Key.GroupRevision, events[0].Fields["entity_revision"])
	require.Equal(t, "connector-A", events[0].Fields["connector_id"])
	require.Equal(t, saved.WorkspaceRevision, events[0].Fields["workspace_revision"], "observed transition revision, not historical output revision")
}

func TestProductionDeliveryEventsFullOrClosedCapturePreservesNativeCommit(t *testing.T) {
	for _, kind := range []string{"full", "closed"} {
		t.Run(kind, func(t *testing.T) {
			f := newDeliveryEventFixture(t)
			var dropped uint64
			if kind == "closed" {
				require.NoError(t, f.broker.Close(t.Context()))
			} else {
				scope := runtimeview.Scope{WorkspaceID: f.record.ID, WorkspaceRevision: f.record.DatabaseSetupRevision, DeviceID: "device-1"}
				for i := range diagnostics.AdmissionRecords {
					id := fmt.Sprintf("retained-%d", i)
					ref := diagnostics.DomainEventRef{ContractVersion: 1, Code: "domain.source.rejected", Binding: runtimeview.Binding{WorkspaceID: scope.WorkspaceID, WorkspaceRevision: scope.WorkspaceRevision, EntityRevision: "R1"}, Locator: runtimeview.EntityLocator{Step: 2, WorkspaceID: scope.WorkspaceID, DeviceID: &scope.DeviceID, RuleID: &id}, Stage: "blocked", Reason: "unavailable", Action: "refresh", Evidence: diagnostics.DomainEvidenceRef{Kind: "source-runtime-transition", ReferenceID: id}}
					raw, err := json.Marshal(ref)
					require.NoError(t, err)
					f.broker.EmitDomain(t.Context(), raw, scope)
				}
				require.NoError(t, f.broker.Flush(t.Context()))
				snapshot, err := f.broker.Snapshot(diagnostics.Query{})
				require.NoError(t, err)
				dropped = snapshot.CaptureDropped
			}
			f.start(t)
			f.waitCommitted(t)
			if kind == "full" {
				require.Eventually(t, func() bool {
					require.NoError(t, f.broker.Flush(t.Context()))
					snapshot, err := f.broker.Snapshot(diagnostics.Query{})
					require.NoError(t, err)
					return snapshot.CaptureDropped > dropped
				}, time.Second, 10*time.Millisecond)
			}
		})
	}
}

func TestProductionDeliveryEventsRestartDoesNotReplayOldReceipt(t *testing.T) {
	f := newDeliveryEventFixture(t)
	f.start(t)
	f.waitCommitted(t)
	require.Len(t, waitDeliveryEvents(t, f.broker, "domain.delivery.first-confirmed"), 1)
	require.NoError(t, f.env.services.groupPipe.Stop(3*time.Second))
	require.NoError(t, f.env.services.runtime.Stop(t.Context()))
	require.NoError(t, f.broker.Close(t.Context()))
	broker, err := diagnostics.New(diagnostics.Options{})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	f.env.services, err = wireGatewayServicesContext(t.Context(), f.env.db, manager, broker)
	require.NoError(t, err)
	f.start(t)
	require.NoError(t, f.env.services.groupPipe.Stop(3*time.Second))
	require.Empty(t, deliveryEventsOfCode(t, broker, "domain.delivery.first-confirmed"))
	require.Empty(t, deliveryEventsOfCode(t, broker, "domain.delivery.recovered"))
	var count int
	require.NoError(t, f.destination.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM readings").Scan(&count))
	require.Equal(t, 1, count)
}

func TestDeliveryEventObserverRejectsMissingOrContradictoryProof(t *testing.T) {
	cases := map[string]func(*groupdelivery.DeliveryObservation){
		"unavailable":        func(o *groupdelivery.DeliveryObservation) { o.Available = false },
		"workspace":          func(o *groupdelivery.DeliveryObservation) { o.Key.WorkspaceID = "" },
		"revision":           func(o *groupdelivery.DeliveryObservation) { o.Key.GroupRevision = "" },
		"observed-scope":     func(o *groupdelivery.DeliveryObservation) { o.WorkspaceRevision = "" },
		"connector":          func(o *groupdelivery.DeliveryObservation) { o.ConnectorID = "" },
		"connector-revision": func(o *groupdelivery.DeliveryObservation) { o.ConnectorRevision = "" },
		"effect":             func(o *groupdelivery.DeliveryObservation) { o.EffectKey = "" },
		"record":             func(o *groupdelivery.DeliveryObservation) { o.RecordID = "" },
		"digest":             func(o *groupdelivery.DeliveryObservation) { o.PayloadDigest = "" },
		"receipt":            func(o *groupdelivery.DeliveryObservation) { o.ReceiptMatched = false },
		"timestamp":          func(o *groupdelivery.DeliveryObservation) { o.CommittedAt = time.Time{} },
		"reason":             func(o *groupdelivery.DeliveryObservation) { o.Reason = runtimeview.ReasonUnavailable },
		"queued":             func(o *groupdelivery.DeliveryObservation) { o.State = groupdelivery.StatePending },
		"unknown-state":      func(o *groupdelivery.DeliveryObservation) { o.State = "password_fixture" },
		"unsafe":             func(o *groupdelivery.DeliveryObservation) { o.Key.GroupID = "Server=private;Password=secret" },
		"false-retry-receipt": func(o *groupdelivery.DeliveryObservation) {
			o.State = groupdelivery.StateRetrying
			o.Reason = runtimeview.ReasonUnavailable
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			broker, err := diagnostics.New(diagnostics.Options{})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
			o := groupdelivery.DeliveryObservation{Available: true, Key: groupdelivery.GroupKey{WorkspaceID: "W", GroupID: "G", GroupRevision: "R1"}, WorkspaceRevision: "W1", ConnectorID: "C", ConnectorRevision: "C1", EffectKey: "E", RecordID: "R", PayloadDigest: "digest", State: groupdelivery.StateCommitted, Reason: runtimeview.ReasonNone, CommittedAt: time.Now(), ReceiptMatched: true}
			mutate(&o)
			observer := newDeliveryEventObserver(broker)
			for range 5 {
				observer(o)
			}
			records := collectorEventRecords(t, broker)
			require.Len(t, records, 1)
			require.Equal(t, "domain.unavailable", records[0].Code)
			require.Empty(t, records[0].Fields)
		})
	}
}

func TestOwnerEventObserverSeparatesInputsFromDurableDelivery(t *testing.T) {
	for _, kind := range []grouppipeline.OwnerObservationKind{grouppipeline.OwnerCounterQuarantine, grouppipeline.OwnerRecordingRefusal} {
		t.Run(string(kind), func(t *testing.T) {
			broker, err := diagnostics.New(diagnostics.Options{})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
			reason, code, evidence := "counter-replay-proof-missing", "domain.counter-replay.quarantined", "counter-replay-inputs"
			if kind == grouppipeline.OwnerRecordingRefusal {
				reason, code, evidence = "recording-handoff-proof-conflict", "domain.recording.refused", "recording-refusal"
			}
			newOwnerEventObserver(broker)(grouppipeline.OwnerObservation{Available: true, Key: groupdelivery.GroupKey{WorkspaceID: "W", GroupID: "G", GroupRevision: "R1"}, WorkspaceRevision: "W1", ConnectorID: "C", ConnectorRevision: "C1", ObservationID: "observation-1", Kind: kind, Reason: reason, PendingSamples: 0})
			records := collectorEventRecords(t, broker)
			require.Len(t, records, 1)
			require.Equal(t, code, records[0].Code)
			require.Equal(t, evidence, records[0].Fields["evidence_kind"])
			require.Equal(t, "0", records[0].Fields["pending_samples"])
			require.Equal(t, "observation-1", records[0].Fields["reference_id"])
			require.NotContains(t, records[0].Fields, "effect_key")
			require.NotContains(t, records[0].Fields, "record_id")
		})
	}
}

func TestOwnerEventObserverUnknownUnsafeOrMissingProofIsUnavailable(t *testing.T) {
	for _, mode := range []string{"unavailable", "foreign", "reason", "kind", "observation"} {
		t.Run(mode, func(t *testing.T) {
			broker, err := diagnostics.New(diagnostics.Options{})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
			o := grouppipeline.OwnerObservation{Available: true, Key: groupdelivery.GroupKey{WorkspaceID: "W", GroupID: "G", GroupRevision: "R1"}, WorkspaceRevision: "W1", ConnectorID: "C", ConnectorRevision: "C1", ObservationID: "O", Kind: grouppipeline.OwnerRecordingRefusal, Reason: "recording-ownership-conflict", PendingSamples: 9}
			switch mode {
			case "unavailable":
				o.Available = false
			case "foreign":
				o.Key.WorkspaceID = "password=secret"
			case "reason":
				o.Reason = "Server=private;Password=secret"
			case "kind":
				o.Kind = "future"
			case "observation":
				o.ObservationID = ""
			}
			newOwnerEventObserver(broker)(o)
			records := collectorEventRecords(t, broker)
			require.Len(t, records, 1)
			require.Equal(t, "domain.unavailable", records[0].Code)
			require.Empty(t, records[0].Fields)
		})
	}
}

func TestDeliveryEventObserverNilBrokerHasNoProviderOrCapture(t *testing.T) {
	delivery, owner := newDeliveryEventObservers(nil, nil)
	require.Nil(t, delivery.ScopeRevision)
	require.Nil(t, delivery.Observe)
	require.Nil(t, owner.ScopeRevision)
	require.Nil(t, owner.Observe)
	require.Nil(t, owner.OnDrop)
	require.Zero(t, owner.RecordLimit)
	require.Nil(t, newDeliveryEventObserver(nil))
	require.Nil(t, newOwnerEventObserver(nil))
}

func TestDeliveryEventScopeReaderIsPureAndDoesNotBootstrap(t *testing.T) {
	f := newDeliveryEventFixture(t)
	delivery, owner := newDeliveryEventObservers(f.broker, f.env.services.workspace)
	require.Equal(t, diagnostics.AdmissionRecords, owner.RecordLimit)
	require.NotNil(t, owner.OnDrop)
	var before, after int64
	changes := func() *int64 {
		var n int64
		require.NoError(t, f.env.db.QueryRowContext(t.Context(), "SELECT total_changes()").Scan(&n))
		return &n
	}
	before = *changes()
	for _, read := range []func(context.Context, groupdelivery.GroupKey) (string, bool){delivery.ScopeRevision, owner.ScopeRevision} {
		revision, ok := read(t.Context(), f.item.Key)
		require.True(t, ok)
		require.Equal(t, f.record.DatabaseSetupRevision, revision)
		foreign := f.item.Key
		foreign.WorkspaceID = "another-workspace"
		revision, ok = read(t.Context(), foreign)
		require.False(t, ok)
		require.Empty(t, revision)
	}
	after = *changes()
	require.Equal(t, before, after)
	_, err := f.env.db.ExecContext(t.Context(), "DELETE FROM system_settings WHERE key='studio_v2_workspace'")
	require.NoError(t, err)
	before = *changes()
	for _, read := range []func(context.Context, groupdelivery.GroupKey) (string, bool){delivery.ScopeRevision, owner.ScopeRevision} {
		revision, ok := read(t.Context(), f.item.Key)
		require.False(t, ok)
		require.Empty(t, revision)
	}
	after = *changes()
	require.Equal(t, before, after)
	var count int
	require.NoError(t, f.env.db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM system_settings WHERE key='studio_v2_workspace'").Scan(&count))
	require.Zero(t, count)
}

func TestDeliveryEventObserverFaultStatesNeverClaimSQLCommit(t *testing.T) {
	cases := map[string]string{groupdelivery.StateRetrying: "domain.delivery.retrying", groupdelivery.StateUnknown: "domain.delivery.unknown", groupdelivery.StateBlocked: "domain.delivery.blocked", groupdelivery.StateQuarantined: "domain.delivery.quarantined"}
	for state, code := range cases {
		t.Run(state, func(t *testing.T) {
			broker, err := diagnostics.New(diagnostics.Options{})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
			newDeliveryEventObserver(broker)(groupdelivery.DeliveryObservation{Available: true, Key: groupdelivery.GroupKey{WorkspaceID: "W", GroupID: "G", GroupRevision: "R1"}, WorkspaceRevision: "W1", ConnectorID: "C", ConnectorRevision: "C1", EffectKey: "E", RecordID: "R", PayloadDigest: "digest", State: state, Reason: runtimeview.ReasonUnavailable})
			records := collectorEventRecords(t, broker)
			require.Len(t, records, 1)
			require.Equal(t, code, records[0].Code)
			require.Equal(t, "durable-transition", records[0].Fields["evidence_kind"])
			require.Equal(t, "E", records[0].Fields["reference_id"])
			require.NotContains(t, records[0].Fields, "committed_at")
			require.NotContains(t, records[0].Fields, "payload_digest")
		})
	}
}

func TestDeliveryEventOwnerFilteredBurstDoesNotCountCaptureLoss(t *testing.T) {
	for _, level := range []string{"error", "warn", "info"} {
		t.Run(level, func(t *testing.T) {
			f := newDeliveryEventFixture(t)
			keys := seedDeliveryOwnerProofFaults(t, f)
			broker, err := diagnostics.New(diagnostics.Options{Level: level})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
			delivery, owner := newDeliveryEventObservers(broker, f.env.services.workspace)
			require.NotNil(t, delivery.Observe, "delivery processing stays independent of owner warn filtering")
			require.NotNil(t, delivery.ScopeRevision)
			reads, dropCalls := 0, 0
			var lost uint64
			if owner.Observe != nil {
				require.Equal(t, diagnostics.AdmissionRecords, owner.RecordLimit, "production cap remains unchanged")
				// Two actual native descriptors exceed this smaller test cap. The production
				// cap itself is verified above and is never changed outside this fixture.
				owner.RecordLimit = 1
				scope := owner.ScopeRevision
				owner.ScopeRevision = func(ctx context.Context, key groupdelivery.GroupKey) (string, bool) { reads++; return scope(ctx, key) }
				drop := owner.OnDrop
				owner.OnDrop = func(n uint64) { dropCalls++; lost += n; drop(n) }
			}
			pipe := grouppipeline.New(grouppipeline.Dependencies{Groups: f.env.services.writeGroups, Tags: f.env.services.tag, Destinations: f.env.services.dbTarget, Inspector: dbtarget.NewReadOnlyTableInspector(f.env.services.dbTarget), Store: groupdelivery.NewStore(f.env.db), OwnerObservation: owner}, grouppipeline.Config{})
			require.NoError(t, pipe.Reconcile(t.Context()))
			quarantines, err := pipe.CounterReplayQuarantines(t.Context(), f.group.ID)
			require.NoError(t, err)
			require.Len(t, quarantines, 2)
			for _, key := range keys {
				count, err := groupdelivery.NewStore(f.env.db).PendingSampleCount(t.Context(), key)
				require.NoError(t, err)
				require.Equal(t, 1, count)
			}
			require.Equal(t, groupdelivery.StatePending, f.state(t).State, "diagnostic overflow cannot finalize or quarantine an output")
			require.NoError(t, broker.Flush(t.Context()))
			observed, err := broker.Snapshot(diagnostics.Query{})
			require.NoError(t, err)
			if level == "error" {
				require.Zero(t, observed.CaptureDropped, "intentionally filtered owner faults are not capture loss")
				require.Zero(t, reads)
				require.Zero(t, dropCalls)
				require.Zero(t, lost)
				require.Nil(t, owner.Observe)
				require.Nil(t, owner.ScopeRevision)
				require.Nil(t, owner.OnDrop)
				require.Zero(t, owner.RecordLimit)
			} else {
				require.Equal(t, uint64(1), observed.CaptureDropped)
				require.Equal(t, 1, dropCalls)
				require.Equal(t, uint64(1), lost)
				require.Equal(t, 2, reads)
				require.Len(t, deliveryEventsOfCode(t, broker, "domain.counter-replay.quarantined"), 1)
			}
		})
	}
}
