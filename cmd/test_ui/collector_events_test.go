package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/diagnostics"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestProductionCollectorEventsActualSimulatorDisconnectThenRecovery(t *testing.T) {
	f := newCollectorEventFixture(t, 1, true, false)
	f.start(t)
	healthy := f.poll(t, 0)
	require.Equal(t, schema.QualityGood, healthy.Quality)
	require.Empty(t, collectorEventRecords(t, f.broker), "healthy acquisition is not a log event")
	f.stopDevice(t, 0)
	failed := f.poll(t, 0)
	require.Equal(t, schema.QualityBad, failed.Quality)
	records := collectorEventRecords(t, f.broker)
	require.Len(t, records, 1)
	first := records[0]
	require.Equal(t, "domain.collector.failed", first.Code)
	require.Equal(t, f.record.ID, first.Fields["workspace_id"])
	require.Equal(t, f.record.DatabaseSetupRevision, first.Fields["workspace_revision"])
	require.Equal(t, f.devices[0].rule.DeviceID, first.Fields["device_id"])
	require.Equal(t, f.devices[0].point.ID, first.Fields["point_id"])
	require.Equal(t, f.devices[0].rule.RevisionID, first.Fields["entity_revision"])
	require.NotEqual(t, failed.ConfigFingerprint, first.Fields["entity_revision"])
	require.Equal(t, failed.AcquisitionID, first.Fields["reference_id"])
	f.recoverDevice(t, 0)
	recovered := f.poll(t, 0)
	require.Equal(t, schema.QualityGood, recovered.Quality)
	records = collectorEventRecords(t, f.broker)
	require.Len(t, records, 2)
	require.Equal(t, "domain.collector.recovered", records[1].Code)
	require.Equal(t, recovered.AcquisitionID, records[1].Fields["reference_id"])
	require.NotEqual(t, failed.AcquisitionID, recovered.AcquisitionID)
	require.Equal(t, first.Fields["entity_revision"], records[1].Fields["entity_revision"])
	raw, err := json.Marshal(records)
	require.NoError(t, err)
	for _, forbidden := range []string{"password_fixture", "127.0.0.1", "config_fingerprint", "quality_reason", "raw_error", "sample_value"} {
		require.NotContains(t, string(raw), forbidden)
	}
}

func TestProductionCollectorEventsMissingWorkspaceProofIsUnattributed(t *testing.T) {
	f := newCollectorEventFixture(t, 1, false, false)
	f.start(t)
	f.stopDevice(t, 0)
	failed := f.poll(t, 0)
	require.Equal(t, schema.QualityBad, failed.Quality)
	records := collectorEventRecords(t, f.broker)
	require.Len(t, records, 1)
	require.Equal(t, "domain.unavailable", records[0].Code)
	require.Empty(t, records[0].Fields)
}

func TestProductionCollectorEventsTwoDevicesKeepIndependentFailures(t *testing.T) {
	f := newCollectorEventFixture(t, 2, true, false)
	f.start(t)
	// Establish ordinary successful reads before the disconnect. This also
	// preserves the native breaker's initial rolling error-rate semantics.
	for i := range 2 {
		require.Equal(t, schema.QualityGood, f.poll(t, i).Quality)
	}
	require.Empty(t, collectorEventRecords(t, f.broker))
	f.stopDevice(t, 0)
	f.stopDevice(t, 1)
	failedA, failedB := f.poll(t, 0), f.poll(t, 1)
	require.Equal(t, schema.QualityBad, failedA.Quality)
	require.Equal(t, schema.QualityBad, failedB.Quality)
	records := collectorEventRecords(t, f.broker)
	require.Len(t, records, 2)
	for i := range 2 {
		require.Equal(t, f.devices[i].rule.DeviceID, records[i].Fields["device_id"])
		require.Equal(t, f.devices[i].point.ID, records[i].Fields["point_id"])
		require.Equal(t, f.devices[i].rule.RevisionID, records[i].Fields["entity_revision"])
	}
	f.recoverDevice(t, 1)
	require.Equal(t, schema.QualityGood, f.poll(t, 1).Quality)
	for range 3 {
		good := f.poll(t, 1)
		require.Equal(t, schema.QualityGood, good.Quality, "native result: %s / %s", good.QualityReason, good.Error)
	}
	records = collectorEventRecords(t, f.broker)
	require.Len(t, records, 3, "healthy B samples neither spam nor clear A")
	require.Equal(t, "domain.collector.recovered", records[2].Code)
	require.Equal(t, f.devices[1].rule.DeviceID, records[2].Fields["device_id"])
	f.recoverDevice(t, 0)
	require.Equal(t, schema.QualityGood, f.poll(t, 0).Quality)
	records = collectorEventRecords(t, f.broker)
	require.Len(t, records, 4)
	require.Equal(t, f.devices[0].rule.DeviceID, records[3].Fields["device_id"])
}

func TestProductionCollectorEventsPeriodicSimulatorDisconnectAndRecovery(t *testing.T) {
	f := newCollectorEventFixture(t, 1, true, true)
	f.stopDevice(t, 0)
	f.start(t)
	require.Eventually(t, func() bool {
		records := collectorEventRecords(t, f.broker)
		return len(records) > 0 && records[0].Code == "domain.collector.failed"
	}, 3*time.Second, 10*time.Millisecond)
	f.recoverDevice(t, 0)
	require.Eventually(t, func() bool {
		for _, event := range collectorEventRecords(t, f.broker) {
			if event.Code == "domain.collector.recovered" {
				return event.Fields["device_id"] == f.devices[0].rule.DeviceID && event.Fields["entity_revision"] == f.devices[0].rule.RevisionID && event.Fields["reference_id"] != ""
			}
		}
		return false
	}, 3*time.Second, 10*time.Millisecond)
}

func TestProductionCollectorEventsUnknownProofBurstIsBounded(t *testing.T) {
	f := newCollectorEventFixture(t, 1, false, false)
	f.start(t)
	for range 50 {
		require.Equal(t, schema.QualityGood, f.poll(t, 0).Quality)
	}
	records := collectorEventRecords(t, f.broker)
	require.Len(t, records, 1)
	require.Equal(t, "domain.unavailable", records[0].Code)
	require.Empty(t, records[0].Fields)
}

func TestProductionCollectorEventsClosedCapturePreservesNativeResult(t *testing.T) {
	f := newCollectorEventFixture(t, 1, true, false)
	f.start(t)
	require.NoError(t, f.broker.Close(t.Context()))
	good := f.poll(t, 0)
	require.Equal(t, schema.QualityGood, good.Quality)
	require.EqualValues(t, 12345, good.Value)
	f.stopDevice(t, 0)
	bad := f.poll(t, 0)
	require.Equal(t, schema.QualityBad, bad.Quality)
	require.NotEmpty(t, bad.Error, "native acquisition error remains available to its existing owner")
	require.NotEmpty(t, bad.AcquisitionID)
}

func TestProductionCollectorEventsFullCapturePreservesNativeResult(t *testing.T) {
	f := newCollectorEventFixture(t, 1, true, false)
	f.start(t)
	scope := runtimeview.Scope{WorkspaceID: f.record.ID, WorkspaceRevision: f.record.DatabaseSetupRevision, DeviceID: f.devices[0].rule.DeviceID}
	for i := range diagnostics.AdmissionRecords {
		id := fmt.Sprintf("retained-fault-%d", i)
		ref := diagnostics.DomainEventRef{ContractVersion: 1, Code: "domain.source.rejected", Binding: runtimeview.Binding{WorkspaceID: scope.WorkspaceID, WorkspaceRevision: scope.WorkspaceRevision, EntityRevision: "R1"}, Locator: runtimeview.EntityLocator{Step: 2, WorkspaceID: scope.WorkspaceID, DeviceID: &scope.DeviceID, RuleID: &id}, Stage: "blocked", Reason: "unavailable", Action: "refresh", Evidence: diagnostics.DomainEvidenceRef{Kind: "source-runtime-transition", ReferenceID: id}}
		raw, err := json.Marshal(ref)
		require.NoError(t, err)
		f.broker.EmitDomain(t.Context(), raw, scope)
	}
	require.NoError(t, f.broker.Flush(t.Context()))
	before, err := f.broker.Snapshot(diagnostics.Query{})
	require.NoError(t, err)
	f.stopDevice(t, 0)
	failed := f.poll(t, 0)
	require.Equal(t, schema.QualityBad, failed.Quality)
	require.NotEmpty(t, failed.Error)
	require.NoError(t, f.broker.Flush(t.Context()))
	after, err := f.broker.Snapshot(diagnostics.Query{})
	require.NoError(t, err)
	require.Greater(t, after.CaptureDropped, before.CaptureDropped, "full capture has explicit loss, not a changed acquisition")
	f.recoverDevice(t, 0)
	good := f.poll(t, 0)
	require.Equal(t, schema.QualityGood, good.Quality)
	require.EqualValues(t, 12345, good.Value)
}

func TestCollectorEventObserverRejectsMissingForeignAndContradictoryProof(t *testing.T) {
	cases := map[string]func(*collector.AcquisitionObservation){
		"missing-workspace":           func(o *collector.AcquisitionObservation) { o.Binding.WorkspaceID = "" },
		"missing-workspace-revision":  func(o *collector.AcquisitionObservation) { o.Binding.WorkspaceRevision = "" },
		"foreign-device":              func(o *collector.AcquisitionObservation) { o.DeviceID = "different-device" },
		"foreign-point":               func(o *collector.AcquisitionObservation) { o.PointID = "different-point" },
		"missing-rule":                func(o *collector.AcquisitionObservation) { o.Binding.RuleID = "" },
		"unsafe-rule":                 func(o *collector.AcquisitionObservation) { o.Binding.RuleID = "Server=private;Password=secret" },
		"missing-source-revision":     func(o *collector.AcquisitionObservation) { o.Binding.RuleRevision = "" },
		"unsafe-source-revision":      func(o *collector.AcquisitionObservation) { o.Binding.RuleRevision = "host:502" },
		"empty-actual-fingerprint":    func(o *collector.AcquisitionObservation) { o.ConfigFingerprint = "" },
		"changed-fingerprint":         func(o *collector.AcquisitionObservation) { o.ConfigFingerprint = "different-fingerprint" },
		"missing-binding-fingerprint": func(o *collector.AcquisitionObservation) { o.Binding.ConfigFingerprint = "" },
		"missing-acquisition":         func(o *collector.AcquisitionObservation) { o.AcquisitionID = "" },
		"unsafe-acquisition":          func(o *collector.AcquisitionObservation) { o.AcquisitionID = "password=secret" },
		"missing-completion":          func(o *collector.AcquisitionObservation) { o.ReceivedAt = time.Time{} },
		"unknown-quality":             func(o *collector.AcquisitionObservation) { o.Quality = "future-quality" },
		"non-normalized-bad":          func(o *collector.AcquisitionObservation) { o.QualityReason = "" },
		"contradictory-good":          func(o *collector.AcquisitionObservation) { o.Quality = schema.QualityGood },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			broker, err := diagnostics.New(diagnostics.Options{})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
			observation := collectorEventObservation()
			mutate(&observation)
			newCollectorEventObserver(broker)(observation)
			records := collectorEventRecords(t, broker)
			require.Len(t, records, 1)
			require.Equal(t, "domain.unavailable", records[0].Code)
			require.Empty(t, records[0].Fields, "invalid proof must not retain partial attribution")
		})
	}
}

func TestCollectorEventObserverClosedProjectionOmitsFingerprintAndQualityReason(t *testing.T) {
	broker, err := diagnostics.New(diagnostics.Options{})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
	observer := newCollectorEventObserver(broker)
	observation := collectorEventObservation()
	observation.QualityReason = "Server=private;Password=secret;sample_value=12345"
	observer(observation)
	records := collectorEventRecords(t, broker)
	require.Len(t, records, 1)
	require.Equal(t, "domain.collector.failed", records[0].Code)
	require.Equal(t, "original-rule-revision", records[0].Fields["entity_revision"])
	require.Equal(t, "acquisition-id", records[0].Fields["reference_id"])
	require.EqualValues(t, 1, records[0].Fields["locator_step"])
	require.NotContains(t, records[0].Fields, "rule_id", "collector locator remains the existing Step 1 device locator")
	raw, err := json.Marshal(records)
	require.NoError(t, err)
	for _, forbidden := range []string{"private", "secret", "sample_value", "12345", "only-internal-fingerprint", "quality_reason"} {
		require.NotContains(t, string(raw), forbidden)
	}
}

func TestCollectorEventObserverNilBrokerIsAbsent(t *testing.T) {
	require.Nil(t, newCollectorEventObserver(nil))
}

func collectorEventObservation() collector.AcquisitionObservation {
	return collector.AcquisitionObservation{
		Binding:  collector.PointObservationBinding{WorkspaceID: "workspace", WorkspaceRevision: "observed-W1", DeviceID: "device", PointID: "point", RuleID: "rule", RuleRevision: "original-rule-revision", ConfigFingerprint: "only-internal-fingerprint"},
		DeviceID: "device", PointID: "point", AcquisitionID: "acquisition-id", ConfigFingerprint: "only-internal-fingerprint",
		Quality: schema.QualityBad, QualityReason: "read-failed", ReceivedAt: time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC),
	}
}
