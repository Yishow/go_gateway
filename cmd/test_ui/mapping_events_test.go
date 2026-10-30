package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/diagnostics"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestProductionMappingSimulatorCastFailureThenRecovery(t *testing.T) {
	f := newMappingEventFixture(t, 1)
	f.start(t)
	f.waitFailure(t)
	records := f.records(t)
	failed := findMappingEvent(records, "domain.mapping.failed", f.sources[0].mapping.ID)
	require.NotNil(t, failed, "actual failing cast must emit attributable mapping evidence")
	require.Equal(t, f.record.ID, failed.Fields["workspace_id"])
	require.Equal(t, f.record.DatabaseSetupRevision, failed.Fields["workspace_revision"])
	require.Equal(t, f.sources[0].rule.DeviceID, failed.Fields["device_id"])
	require.Equal(t, f.sources[0].mapping.PointID, failed.Fields["point_id"])
	require.Equal(t, f.sources[0].mapping.TagID, failed.Fields["tag_id"])
	require.Equal(t, f.sources[0].revision, failed.Fields["entity_revision"])
	require.Equal(t, "blocked", failed.Fields["stage"])
	require.Equal(t, "mapping-result", failed.Fields["evidence_kind"])
	require.NotEmpty(t, failed.Fields["reference_id"])
	before := f.services.runtime.Snapshot()
	f.setWord(t, 0, 123)
	require.Eventually(t, func() bool { return f.services.runtime.Snapshot().WriteSuccess > before.WriteSuccess }, 3*time.Second, 10*time.Millisecond, "native mapping and writer result must recover independently")
	recovered := findMappingEvent(f.records(t), "domain.mapping.recovered", f.sources[0].mapping.ID)
	require.NotNil(t, recovered)
	require.Equal(t, f.sources[0].revision, recovered.Fields["entity_revision"])
	require.Equal(t, "acquisition", recovered.Fields["stage"])
	require.NotEqual(t, failed.Fields["reference_id"], recovered.Fields["reference_id"])
}

func TestProductionMappingTwoDevicesAndRevisionChangeStaySeparate(t *testing.T) {
	f := newMappingEventFixture(t, 2)
	f.start(t)
	require.Eventually(t, func() bool {
		records := f.records(t)
		return findMappingEvent(records, "domain.mapping.failed", f.sources[0].mapping.ID) != nil && findMappingEvent(records, "domain.mapping.failed", f.sources[1].mapping.ID) != nil
	}, 3*time.Second, 10*time.Millisecond)
	oldRevision := f.sources[0].revision
	before := f.services.runtime.Snapshot()
	f.updatePipeline(t, 0, mappingCastPipeline("int32", nil))
	require.NotEqual(t, oldRevision, f.sources[0].revision)
	require.NoError(t, f.services.runtime.RefreshMappings(t.Context()))
	require.Eventually(t, func() bool { return f.services.runtime.Snapshot().WriteSuccess > before.WriteSuccess }, 3*time.Second, 10*time.Millisecond)
	require.Nil(t, findMappingEvent(f.records(t), "domain.mapping.recovered", f.sources[0].mapping.ID), "R2 success cannot recover the retained R1 failure")
	f.setWord(t, 1, 123)
	require.Eventually(t, func() bool {
		return findMappingEvent(f.records(t), "domain.mapping.recovered", f.sources[1].mapping.ID) != nil
	}, 3*time.Second, 10*time.Millisecond)
	recovered := findMappingEvent(f.records(t), "domain.mapping.recovered", f.sources[1].mapping.ID)
	require.Equal(t, f.sources[1].rule.DeviceID, recovered.Fields["device_id"])
	require.Equal(t, f.sources[1].revision, recovered.Fields["entity_revision"])
	require.Nil(t, findMappingEvent(f.records(t), "domain.mapping.recovered", f.sources[0].mapping.ID))
}

func TestProductionMappingSecretParamsNeverReachEvents(t *testing.T) {
	f := newMappingEventFixture(t, 1, true)
	sub, err := f.broker.Subscribe(diagnostics.Query{}, "")
	require.NoError(t, err)
	defer sub.Close()
	f.updatePipeline(t, 0, mappingCastPipeline("int16", map[string]any{"password": "password_fixture", "nested": map[string]any{"dsn": "Server=private;Password=secret"}}))
	f.start(t)
	f.waitFailure(t)
	records := f.records(t)
	failed := findMappingEvent(records, "domain.mapping.failed", f.sources[0].mapping.ID)
	require.NotNil(t, failed)
	require.Equal(t, f.sources[0].revision, failed.Fields["entity_revision"])
	raw, err := json.Marshal(records)
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	live, err := sub.Next(ctx)
	require.NoError(t, err)
	require.Equal(t, *failed, live)
	sub.Ack(live.Cursor())
	fileBytes, err := os.ReadFile(f.logPath)
	require.NoError(t, err)
	require.Contains(t, string(fileBytes), `"code":"domain.mapping.failed"`)
	for _, forbidden := range []string{"password_fixture", "private", "secret", "nested", "target_type", "40000", "127.0.0.1", "raw_error", "transform_pipeline"} {
		require.NotContains(t, string(raw), forbidden)
		require.NotContains(t, string(fileBytes), forbidden)
	}
}

func TestProductionMappingHealthySamplesDoNotSpam(t *testing.T) {
	f := newMappingEventFixture(t, 1)
	f.setWord(t, 0, 123)
	f.start(t)
	require.Eventually(t, func() bool { return f.services.runtime.Snapshot().WriteSuccess >= 4 }, 3*time.Second, 10*time.Millisecond)
	require.Empty(t, f.records(t))
}

func TestProductionMappingMissingScopeIsBoundedUnavailable(t *testing.T) {
	f := newMappingEventFixture(t, 1)
	f.record.DatabaseSetupRevision = ""
	require.NoError(t, workspace.NewSQLRepository(f.db).Save(t.Context(), f.record))
	f.start(t)
	require.Eventually(t, func() bool { return f.services.runtime.Snapshot().MappingError >= 4 }, 3*time.Second, 10*time.Millisecond)
	records := f.records(t)
	require.Len(t, records, 1)
	require.Equal(t, "domain.unavailable", records[0].Code)
	require.Empty(t, records[0].Fields)
}

func TestProductionMappingMalformedPersistedTransformNeverExposesRawError(t *testing.T) {
	f := newMappingEventFixture(t, 1)
	_, err := f.db.Exec(`UPDATE mappings SET transform_pipeline=? WHERE id=?`, `[{"type":"password_fixture@private-host","params":{"password":"secret"}}]`, f.sources[0].mapping.ID)
	require.NoError(t, err)
	f.start(t)
	f.waitFailure(t)
	records := f.records(t)
	require.NotNil(t, findMappingEvent(records, "domain.mapping.failed", f.sources[0].mapping.ID))
	raw, err := json.Marshal(records)
	require.NoError(t, err)
	for _, forbidden := range []string{"password_fixture", "private-host", "secret", "transform_pipeline", "40000"} {
		require.NotContains(t, string(raw), forbidden)
	}
}

func TestProductionMappingFullOrClosedCapturePreservesNativeProcessing(t *testing.T) {
	for _, kind := range []string{"full", "closed"} {
		t.Run(kind, func(t *testing.T) {
			f := newMappingEventFixture(t, 1)
			var dropped uint64
			if kind == "closed" {
				require.NoError(t, f.broker.Close(t.Context()))
			} else {
				scope := runtimeview.Scope{WorkspaceID: f.record.ID, WorkspaceRevision: f.record.DatabaseSetupRevision, DeviceID: f.sources[0].rule.DeviceID}
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
			f.waitFailure(t)
			if kind == "full" {
				require.NoError(t, f.broker.Flush(t.Context()))
				snapshot, err := f.broker.Snapshot(diagnostics.Query{})
				require.NoError(t, err)
				require.Greater(t, snapshot.CaptureDropped, dropped)
			}
			before := f.services.runtime.Snapshot()
			f.setWord(t, 0, 123)
			require.Eventually(t, func() bool { return f.services.runtime.Snapshot().WriteSuccess > before.WriteSuccess }, 3*time.Second, 10*time.Millisecond)
		})
	}
}

func TestMappingEventObserverRejectsIncompleteUnknownOrUnsafeMetadata(t *testing.T) {
	cases := map[string]func(*datalinkruntime.MappingObservation){
		"workspace":          func(o *datalinkruntime.MappingObservation) { o.WorkspaceID = "" },
		"workspace-revision": func(o *datalinkruntime.MappingObservation) { o.WorkspaceRevision = "" },
		"device":             func(o *datalinkruntime.MappingObservation) { o.DeviceID = "" },
		"point":              func(o *datalinkruntime.MappingObservation) { o.PointID = "" },
		"tag":                func(o *datalinkruntime.MappingObservation) { o.TagID = "" },
		"mapping":            func(o *datalinkruntime.MappingObservation) { o.MappingID = "" },
		"mapping-revision":   func(o *datalinkruntime.MappingObservation) { o.MappingRevision = "" },
		"acquisition":        func(o *datalinkruntime.MappingObservation) { o.AcquisitionID = "" },
		"completion":         func(o *datalinkruntime.MappingObservation) { o.ReceivedAt = time.Time{} },
		"unsafe-id":          func(o *datalinkruntime.MappingObservation) { o.MappingID = "password=secret" },
		"unknown-result":     func(o *datalinkruntime.MappingObservation) { o.Result = "future" },
		"unavailable":        func(o *datalinkruntime.MappingObservation) { o.Result = datalinkruntime.MappingObservationUnavailable },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			broker, err := diagnostics.New(diagnostics.Options{})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
			observation := datalinkruntime.MappingObservation{WorkspaceID: "W", WorkspaceRevision: "W1", DeviceID: "D", PointID: "P", TagID: "T", MappingID: "M", MappingRevision: "M1", AcquisitionID: "A", ReceivedAt: time.Now(), Result: datalinkruntime.MappingObservationFailed}
			mutate(&observation)
			newMappingEventObserver(broker)(observation)
			records := collectorEventRecords(t, broker)
			require.Len(t, records, 1)
			require.Equal(t, "domain.unavailable", records[0].Code)
			require.Empty(t, records[0].Fields)
		})
	}
}

func TestMappingEventObserverNilBrokerIsAbsent(t *testing.T) {
	require.Nil(t, newMappingEventObserver(nil))
}
