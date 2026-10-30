package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/diagnostics"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/stretchr/testify/require"
)

// Faults live only at SQL/capture boundaries. The actual production composition,
// SQLite journal, native Tick gate and final proof reads remain in use. Injected
// classified errors test wiring, not the correctness of native classification.
func diagnosticFaultFixture(t *testing.T, brokers ...*diagnostics.Broker) (deliveryEventFixture, *lostcommit.Faults) {
	t.Helper()
	f := newDeliveryEventFixture(t)
	if len(brokers) > 0 {
		f.broker = brokers[0]
	}
	require.NoError(t, f.env.services.runtime.Stop(t.Context()))
	var seq int
	var name, path string
	require.NoError(t, f.env.db.QueryRow("PRAGMA database_list").Scan(&seq, &name, &path))
	fault := &lostcommit.Faults{}
	db := lostcommit.OpenFaults(path+"?_pragma=busy_timeout(5000)", fault)
	datalink.ApplySQLitePoolDefaults(db)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	services, err := wireGatewayServicesContext(t.Context(), db, manager, f.broker)
	require.NoError(t, err)
	f.env.db, f.env.services = db, services
	t.Cleanup(func() {
		fault.Reject.Store(nil)
		require.NoError(t, services.groupPipe.Stop(3*time.Second))
		require.NoError(t, services.runtime.Stop(context.Background()))
	})
	return f, fault
}

func TestProductionDiagnosticTickCaptureRemovesOnlyDuplicate(t *testing.T) {
	for _, kind := range []string{"counter", "owner", "joined"} {
		t.Run(kind, func(t *testing.T) {
			f, faults := diagnosticFaultFixture(t)
			require.NoError(t, f.env.services.groupPipe.Reconcile(t.Context()))
			injected := error(&counterproof.CounterBindingProofError{Reason: counterproof.ErrCounterBindingProofConflict})
			code := "domain.counter-replay.quarantined"
			if kind == "owner" || kind == "joined" {
				injected, code = groupdelivery.ErrOwnershipConflict, "domain.recording.refused"
			}
			if kind == "joined" {
				injected = errors.Join(injected, errors.New("password_fixture nested_secret private_host"))
			}
			var hits atomic.Int32
			reject := func(query string) error {
				if strings.Contains(query, "INSERT INTO wg_delivery_buckets") {
					hits.Add(1)
					return injected
				}
				return nil
			}
			faults.Reject.Store(&reject)
			f.env.services.groupPipe.TickAll(t.Context())
			faults.Reject.Store(nil)
			require.Equal(t, int32(1), hits.Load(), "one real native closure refusal")
			require.Len(t, deliveryEventsOfCode(t, f.broker, code), 1)
			generic := deliveryEventsOfCode(t, f.broker, "runtime.group_failed")
			if kind == "joined" {
				require.Len(t, generic, 1, "a covered leaf cannot hide a joined operational failure")
				raw, err := json.Marshal(collectorEventRecords(t, f.broker))
				require.NoError(t, err)
				require.NotContains(t, string(raw), "password_fixture")
			} else {
				require.Empty(t, generic, "accepted exact Tick observation must replace only its fixed generic duplicate")
			}
			require.Equal(t, f.item, f.state(t), "diagnostics cannot change the queued durable output")
		})
	}
}

func TestProductionDiagnosticOperationalFailuresStayGeneric(t *testing.T) {
	for _, mode := range []string{"receipt", "ready-heads"} {
		t.Run(mode, func(t *testing.T) {
			f, faults := diagnosticFaultFixture(t)
			var hits atomic.Int32
			reject := func(query string) error {
				matches := strings.Contains(query, "INSERT INTO wg_delivery_receipts")
				if mode == "ready-heads" {
					matches = strings.Contains(query, "SELECT o.partition_key, o.effect_key")
				}
				if matches {
					hits.Add(1)
					return errors.New("password_fixture nested_secret private_host")
				}
				return nil
			}
			faults.Reject.Store(&reject)
			f.start(t)
			require.Eventually(t, func() bool { return hits.Load() > 0 }, 5*time.Second, 10*time.Millisecond)
			require.NotEmpty(t, waitDeliveryEvents(t, f.broker, "runtime.group_failed"))
			require.NoError(t, f.env.services.groupPipe.Stop(3*time.Second))
			faults.Reject.Store(nil)
			raw, err := json.Marshal(collectorEventRecords(t, f.broker))
			require.NoError(t, err)
			for _, secret := range []string{"password_fixture", "nested_secret", "private_host"} {
				require.NotContains(t, string(raw), secret)
			}
			require.Empty(t, deliveryEventsOfCode(t, f.broker, "domain.delivery.first-confirmed"))
		})
	}
}

func TestProductionDiagnosticSettledDeliveryAlreadyHasNoGeneric(t *testing.T) {
	f := newDeliveryEventFixture(t)
	f.start(t)
	f.waitCommitted(t)
	require.Len(t, waitDeliveryEvents(t, f.broker, "domain.delivery.first-confirmed"), 1)
	require.Empty(t, deliveryEventsOfCode(t, f.broker, "runtime.group_failed"))
}

func TestDiagnosticCaptureAdaptersAcknowledgeOnlyValidatedLocalCapture(t *testing.T) {
	for _, mode := range []string{"valid", "unavailable", "unsafe", "unknown", "reason", "filtered", "closed", "full"} {
		t.Run(mode, func(t *testing.T) {
			options := diagnostics.Options{}
			if mode == "filtered" {
				options.Level = "error"
			}
			broker, err := diagnostics.New(options)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
			_, owner := newDeliveryEventObservers(broker, nil)
			if mode == "filtered" {
				require.Nil(t, owner.ObserveCaptured)
				return
			}
			require.NotNil(t, owner.ObserveCaptured, "production owner factory must expose local capture acknowledgment")
			o := diagnosticOwnerObservation()
			switch mode {
			case "unavailable":
				o.Available = false
			case "unsafe":
				o.Key.GroupID = "password=secret"
			case "unknown":
				o.Kind = "future"
			case "reason":
				o.Reason = "password=secret"
			case "closed":
				require.NoError(t, broker.Close(t.Context()))
			case "full":
				for i := range diagnostics.AdmissionRecords {
					fill := o
					fill.Key.GroupID = fmt.Sprintf("retained-%d", i)
					owner.ObserveCaptured(fill)
				}
				require.NoError(t, broker.Flush(t.Context()))
			}
			require.Equal(t, mode == "valid", owner.ObserveCaptured(o))
			if mode == "valid" {
				require.True(t, owner.ObserveCaptured(o), "same local episode coalesces without another record")
				owner.Observe(o) // legacy adapter shares conversion, not parallel authority
				require.Len(t, deliveryEventsOfCode(t, broker, "domain.recording.refused"), 1)
			}
		})
	}
	_, owner := newDeliveryEventObservers(nil, nil)
	require.Nil(t, owner.ObserveCaptured)
}

func diagnosticOwnerObservation() grouppipeline.OwnerObservation {
	return grouppipeline.OwnerObservation{Available: true, Key: groupdelivery.GroupKey{WorkspaceID: "W", GroupID: "G", GroupRevision: "R1"}, WorkspaceRevision: "W1", ConnectorID: "C", ConnectorRevision: "C1", ObservationID: "O", Kind: grouppipeline.OwnerRecordingRefusal, Reason: "recording-ownership-conflict", PendingSamples: 2}
}
func TestDiagnosticCaptureAcknowledgmentDoesNotPromiseFileRetention(t *testing.T) {
	dir := canonicalTempDir(t)
	sink, err := diagnostics.OpenFileSink(diagnostics.FileOptions{Path: filepath.Join(dir, "logs", "events.jsonl"), DatabasePath: filepath.Join(dir, "configuration.db"), DatabaseID: strings.Repeat("a", 64)})
	require.NoError(t, err)
	require.NoError(t, sink.Close(t.Context()))
	broker, err := diagnostics.New(diagnostics.Options{File: sink})
	require.NoError(t, err)
	t.Cleanup(func() { require.ErrorContains(t, broker.Close(context.Background()), "storage is degraded") })
	_, owner := newDeliveryEventObservers(broker, nil)
	require.NotNil(t, owner.ObserveCaptured)
	require.True(t, owner.ObserveCaptured(diagnosticOwnerObservation()), "acceptance ends at existing coalescer, not downstream file")
	f, faults := diagnosticFaultFixture(t, broker)
	require.NoError(t, f.env.services.groupPipe.Reconcile(t.Context()))
	reject := func(query string) error {
		if strings.Contains(query, "INSERT INTO wg_delivery_buckets") {
			return groupdelivery.ErrOwnershipConflict
		}
		return nil
	}
	faults.Reject.Store(&reject)
	f.env.services.groupPipe.TickAll(t.Context())
	faults.Reject.Store(nil)
	require.Equal(t, f.item, f.state(t), "downstream loss cannot change native refusal/output")
	require.ErrorContains(t, broker.Flush(t.Context()), "storage is degraded")
	view, err := broker.Snapshot(diagnostics.Query{})
	require.NoError(t, err)
	require.Positive(t, view.FileDropped, "existing sink loss accounting stays honest")
	for _, event := range view.Records {
		require.NotEqual(t, "runtime.group_failed", event.Code, "downstream loss after local capture does not promise another generic record")
	}
}

// Hard-limit quota originates in AcceptSample, not Tick; preserve its existing
// native refusal and single typed event without introducing quota Tick coverage.
func TestProductionDiagnosticNativeQuotaRefusalRemainsUnchanged(t *testing.T) {
	f := newIntakeEventFixture(t)
	clock := intakeEventClock()
	pipe, _ := newObservedQuotaPipeline(t, f, f.broker, groupdelivery.QuotaConfig{GlobalMaxBytes: 1}, clock)
	refusal := quotaEventRefusal(t, f, pipe, clock)
	require.Equal(t, groupdelivery.QuotaScopeGlobal, refusal.Scope)
	require.Equal(t, int64(1), refusal.MaxBytes)
	require.Len(t, deliveryEventsOfCode(t, f.broker, "domain.quota.refused"), 1)
	require.Empty(t, deliveryEventsOfCode(t, f.broker, "runtime.group_failed"))
	require.Equal(t, f.item, f.state(t))
}

func TestProductionDiagnosticFilteredOrCoalescerRejectedCaptureFallsBack(t *testing.T) {
	for _, mode := range []string{"filtered", "full"} {
		t.Run(mode, func(t *testing.T) {
			options := diagnostics.Options{}
			if mode == "filtered" {
				options.Level = "error"
			}
			broker, err := diagnostics.New(options)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
			f, faults := diagnosticFaultFixture(t, broker)
			require.NoError(t, f.env.services.groupPipe.Reconcile(t.Context()))
			if mode == "full" {
				_, observer := newDeliveryEventObservers(broker, nil)
				for i := range diagnostics.AdmissionRecords {
					o := diagnosticOwnerObservation()
					o.Key.GroupID = fmt.Sprintf("occupied-%d", i)
					observer.ObserveCaptured(o)
					require.NoError(t, broker.Flush(t.Context()))
				}
			}
			before, err := broker.Snapshot(diagnostics.Query{})
			require.NoError(t, err)
			reject := func(query string) error {
				if strings.Contains(query, "INSERT INTO wg_delivery_buckets") {
					return groupdelivery.ErrOwnershipConflict
				}
				return nil
			}
			faults.Reject.Store(&reject)
			f.env.services.groupPipe.TickAll(t.Context())
			faults.Reject.Store(nil)
			require.Len(t, deliveryEventsOfCode(t, broker, "runtime.group_failed"), 1)
			require.Equal(t, f.item, f.state(t))
			after, err := broker.Snapshot(diagnostics.Query{})
			require.NoError(t, err)
			if mode == "full" {
				require.Greater(t, after.CaptureDropped, before.CaptureDropped)
			} else {
				require.Equal(t, before.CaptureDropped, after.CaptureDropped)
			}
		})
	}
}
