package main

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionAcceptanceScheduledPilot(t *testing.T) {
	runAcceptanceScheduled(t, 1, 3, time.Millisecond, 3)
}

func TestAcquisitionAcceptanceDeliveryPilot(t *testing.T) {
	runAcceptanceScheduled(t, 2, 128, time.Millisecond, 3)
}

// Join every required native ID to its immutable successful publication fact.
// Test-only first-N completion buffers are bounded; output delivery is separate.
func runAcceptanceScheduled(t *testing.T, devices, points int, rtt time.Duration, warmScans int) {
	t.Helper()
	goroutines := startAcceptanceGoroutineDiagnostic(t)
	setupStarted := time.Now()
	f, load, workspaceID := newAcceptanceScheduledLoad(t, devices, points, rtt)
	expectedGroups := make([]string, devices)
	for i, d := range load {
		pointsInSQL, err := f.services.point.ListByDevice(t.Context(), d.id)
		require.NoError(t, err)
		require.Len(t, pointsInSQL, points)
		for _, point := range pointsInSQL {
			require.NotNil(t, point.PollingGroupID)
			if expectedGroups[i] == "" {
				expectedGroups[i] = *point.PollingGroupID
			}
			require.Equal(t, expectedGroups[i], *point.PollingGroupID)
		}
	}
	t.Logf("fixture_setup_elapsed=%s (canonical persisted seed; excludes cycle; not UI hot authoring)", time.Since(setupStarted))
	native := newAcceptanceNativeCapture(load, warmScans+1)
	// Frozen geometry has one declared canonical plan per device. Retain only
	// that fingerprint, not another scan/value history or production registry.
	plans := make(map[string]string, devices)
	f.capture.acceptanceObserver = func(observation collector.AcquisitionObservation) {
		native.accept(observation)
		native.mu.Lock()
		defer native.mu.Unlock()
		if plan := plans[observation.DeviceID]; plan != "" && plan != observation.PlanFingerprint {
			native.bad = fmt.Errorf("declared native plan changed for device %s", observation.DeviceID)
		}
		plans[observation.DeviceID] = observation.PlanFingerprint
	}
	completions := newAcceptanceCompletionCapture(load, warmScans+1)
	f.capture.acceptanceCompletionObserver = completions.accept
	drain := startAcceptanceSampleDrain(t, f.capture, load, warmScans+1)
	observed := make([][]collector.AcquisitionActualMetrics, devices)
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	cpuBefore, cpuAvailable := acceptanceProcessCPU()
	started := time.Now()
	require.NoError(t, f.services.runtime.Start(t.Context()))
	t.Logf("observation phase=runtime-start-return elapsed=%s native_points=%d", time.Since(started), native.count())
	goroutines.count("after-runtime-start", "runtime/ingestor, scheduler, peer and two test drain workers active")
	t.Logf("observation_method=published-completion-v1; first %d publication facts/device bounded in test only, no latest polling/replanning; callback facts are not current source authorization; freshness is age at completion", warmScans+1)
	t.Cleanup(func() { drain.stopValidation(); require.NoError(t, f.services.scheduler.Stop()) })
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	// This is a test failure bound, not an early success/observation cutoff.
	deadline := time.NewTimer(10 * time.Minute)
	defer deadline.Stop()
	peakQueue, peakGoroutines := 0, runtime.NumGoroutine()
	for {
		all := true
		for i, d := range load {
			if len(observed[i]) >= warmScans+1 {
				continue
			}
			all = false
			expectedScan, complete, nativeErr := native.next(d.id, len(observed[i]))
			require.NoError(t, nativeErr)
			if !complete {
				continue
			}
			event, available, err := completions.next(d.id, expectedScan)
			require.NoError(t, err)
			if !available {
				continue
			}
			require.Equal(t, workspaceID, event.WorkspaceID)
			require.Equal(t, d.sourceHead, event.SourceRevision)
			require.Equal(t, d.id, event.DeviceID)
			require.Equal(t, expectedGroups[i], event.PollingGroupID)
			require.Equal(t, f.services.scheduler.InstallationScope(), event.InstallationScope)
			require.NotEmpty(t, event.IdentityHead)
			require.Equal(t, points, event.PointCount)
			actual := event.Metrics
			require.Equal(t, expectedScan, actual.AcquisitionID, "device=%s ordinal=%d: complete publication/native ID join required; never substitute a later sample", d.id, len(observed[i]))
			native.mu.Lock()
			declaredPlan := plans[d.id]
			native.mu.Unlock()
			require.Equal(t, declaredPlan, actual.DeclaredPlanRevision)
			// These exact MC geometries have one canonical random request; the
			// independent paired planner/adapter matrix verifies the same limits.
			require.EqualValues(t, 1, actual.WireAttempts)
			require.Zero(t, actual.RetryWireAttempts)
			require.Zero(t, actual.FallbackWireAttempts)
			require.Equal(t, "observed", actual.LogicalRequests.Status)
			require.EqualValues(t, 1, actual.LogicalRequests.Succeeded)
			require.Zero(t, actual.LogicalRequests.Failed)
			require.Zero(t, actual.LogicalRequests.Canceled)
			require.NotNil(t, actual.LastGoodAgeMS)
			observed[i] = append(observed[i], actual)
		}

		stats := f.services.scheduler.DueDiagnostics()
		peakQueue = max(peakQueue, stats.Pending)
		currentGoroutines := runtime.NumGoroutine()
		peakGoroutines = max(peakGoroutines, currentGoroutines)
		goroutines.observe(currentGoroutines)
		require.LessOrEqual(t, stats.Pending, 1024, "production default pending capacity")
		if all {
			break
		}
		select {
		case <-ticker.C:
		case <-deadline.C:
			t.Fatal("not every device completed the required measured scans")
		case <-t.Context().Done():
			t.Fatal(t.Context().Err())
		}
	}
	// Join only the producer; keep the original ingestor running and observe a
	// coherent fixed cutoff. Queued values are not completed output or drops.
	require.NoError(t, f.services.scheduler.Stop())
	goroutines.count("after-producer-stop", "scheduler joined; ingestor/runtime, peer, sample drain and ValueEvent reader still active")
	observerLoss, observerSurplus := completions.counts()
	t.Logf("completion observer first_required_loss=%d surplus_after_required=%d unknown_device=%d callback_calls=%d callback_body_duration=%s (includes bounded test channel/atomic recording, excludes producer detached copy; total process CPU/alloc includes both; not product output loss)", observerLoss, observerSurplus, completions.unknown.Load(), completions.callbackCalls.Load(), time.Duration(completions.callbackNanos.Load()))
	require.Zero(t, observerLoss, "missing required completion cannot be replaced")
	require.Zero(t, completions.unknown.Load())
	require.NoError(t, native.validationError())
	cut := acceptanceDeliverySnapshot(t, f, native, drain)
	stats := f.services.scheduler.DueDiagnostics()
	loss := f.services.scheduler.OutputLossDiagnostics()
	t.Logf("delivery cutoff=first-coherent-after-native-stop native_points=%d collected=%d queue_depth=%d sink_points=%d actual_output_loss=%d collected_not_sink=%d runtime_stats=%+v rejected=%d (error counters may overlap/postdate sink, not a causal partition); queued and ingestion-not-sink are NOT output drops or delivered values", cut.native, cut.ingestor.CollectedTotal, cut.queued, cut.sink, cut.dropped, cut.ingestor.CollectedTotal-cut.sink, cut.ingestor, cut.rejected)
	runtime.ReadMemStats(&after)
	cpuAfter, cpuStillAvailable := acceptanceProcessCPU()
	t.Logf("scheduled process_cpu_available=%t process_cpu=%s (includes full runtime bootstrap, scans, observer/drains, completion observation and stop; excludes persisted seed and compile)", cpuAvailable && cpuStillAvailable, cpuAfter-cpuBefore)
	for i, d := range load {
		cycles, ages := make([]time.Duration, warmScans), make([]time.Duration, warmScans)
		overruns := 0
		for j, actual := range observed[i][1:] {
			cycles[j] = time.Duration(actual.CycleLatencyMS * float64(time.Millisecond))
			ages[j] = time.Duration(*actual.LastGoodAgeMS * float64(time.Millisecond))
			if actual.Overrun {
				overruns++
			}
		}
		ids := make([]string, len(observed[i]))
		for j, actual := range observed[i] {
			ids[j] = actual.AcquisitionID
		}
		delivered := 0
		for _, id := range ids {
			delivered += cut.deliveries[d.id][id]
		}
		missing := len(ids)*points - delivered
		t.Logf("delivery device=%s native_selected_points=%d delivered_values_checked=%d missing_at_cutoff=%d (native observer exposes quality/identity, not dropped decoded values)", d.id, len(ids)*points, delivered, missing)
		wire := d.metrics.snapshot()
		require.EqualValues(t, 1, wire.Dials)
		t.Logf("scheduled device=%s consecutive_warm_n=%d cold_first_scan_ms=%f p50=%s p95=%s last_good_age_at_completion_p95=%s warm_overruns=%d frames_including_inflight=%+v", d.id, warmScans, observed[i][0].CycleLatencyMS, acceptancePercentile(cycles, 50), acceptancePercentile(cycles, 95), acceptancePercentile(ages, 95), overruns, wire)
	}
	t.Logf("scheduled elapsed=%s cadence=1s active_capacity=32 pending_capacity=1024 observed_peak_pending=%d observed_peak_goroutines=%d due=%+v output_loss_total=%d output_loss_unattributed=%d output_loss_evicted_detail=%d process_alloc_bytes=%d process_allocations=%d GOMAXPROCS=%d; setup/probe excluded; actual IDs joined to consecutive complete native quality/identity observations; delivered values checked separately; paired adapters prove all-point decoded equivalence; no SQL durability claim", time.Since(started), peakQueue, peakGoroutines, stats, loss.Total, loss.Unattributed, loss.EvictedDetailLoss, after.TotalAlloc-before.TotalAlloc, after.Mallocs-before.Mallocs, runtime.GOMAXPROCS(0))
	drain.stopValidation()
}
