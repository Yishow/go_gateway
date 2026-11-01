package collector

import (
	"fmt"
	"sync"
	"testing"

	"go-gateway/internal/datalink/schema"
)

func lossValue(scope, revision, point string) CollectedValue {
	return CollectedValue{InstallationScope: scope, PointRevision: revision, PointID: point, Quality: schema.QualityGood}
}

func assertOutputLossAccounting(t *testing.T, snapshot OutputLossSnapshot) {
	t.Helper()
	accounted := snapshot.Unattributed + snapshot.EvictedDetailLoss
	for _, series := range snapshot.Recent {
		accounted += series.Losses
	}
	if accounted != snapshot.Total {
		t.Fatalf("loss accounting: total=%d, accounted=%d: %+v", snapshot.Total, accounted, snapshot)
	}
	if len(snapshot.Recent) > 256 {
		t.Fatalf("unbounded recent loss series: %d", len(snapshot.Recent))
	}
}

func TestOutputLossCapturedIdentityAndUnattributed(t *testing.T) {
	var counter outputLossCounter
	old := lossValue("old-install", "old-revision", "point")
	counter.record(old)
	counter.record(old)
	counter.record(lossValue("new-install", "old-revision", "point"))
	counter.record(lossValue("old-install", "new-revision", "point"))
	counter.record(lossValue("", "revision", "point"))
	counter.record(lossValue("scope", "", "point"))
	counter.record(lossValue("scope", "revision", ""))
	snapshot := counter.snapshot()
	if snapshot.Total != 7 || snapshot.Unattributed != 3 || len(snapshot.Recent) != 3 {
		t.Fatalf("captured identities or missing identity attribution lost: %+v", snapshot)
	}
	if got := snapshot.Recent[0]; got.InstallationScope != "old-install" || got.PointRevision != "old-revision" || got.PointID != "point" || got.Losses != 2 {
		t.Fatalf("old identity changed: %+v", got)
	}
	assertOutputLossAccounting(t, snapshot)
	// Returned diagnostic values must not expose mutable internal state.
	snapshot.Recent[0].PointID = "changed-by-reader"
	snapshot.Recent[0].Losses = 999
	if got := counter.snapshot().Recent[0]; got.PointID != "point" || got.Losses != 2 {
		t.Fatalf("reader mutated counters: %+v", got)
	}
}

func TestOutputLossRecentSeriesBoundAndEvictionAccounting(t *testing.T) {
	var counter outputLossCounter
	for i := range 256 {
		counter.record(lossValue("scope", "revision", fmt.Sprint(i)))
	}
	// Touching an existing identity makes it recent without allocating a slot.
	counter.record(lossValue("scope", "revision", "0"))
	counter.record(lossValue("scope", "revision", "256"))
	first := counter.snapshot()
	if first.Total != 258 || first.EvictedDetailLoss != 1 || len(first.Recent) != 256 {
		t.Fatalf("first eviction: %+v", first)
	}
	for _, series := range first.Recent {
		if series.PointID == "1" {
			t.Fatal("least recently lost identity was not evicted")
		}
	}
	// A returning evicted identity starts new recent detail, not invented history.
	counter.record(lossValue("scope", "revision", "1"))
	second := counter.snapshot()
	if second.Total != 259 || second.EvictedDetailLoss != 2 || len(second.Recent) != 256 {
		t.Fatalf("reintroduced identity accounting: %+v", second)
	}
	assertOutputLossAccounting(t, first)
	assertOutputLossAccounting(t, second)
}

func TestOutputLossEvictsAllLossesOfOldSeries(t *testing.T) {
	var counter outputLossCounter
	for range 4 {
		counter.record(lossValue("scope", "old", "point"))
	}
	for i := range 256 {
		counter.record(lossValue("scope", "new", fmt.Sprint(i)))
	}
	snapshot := counter.snapshot()
	if snapshot.Total != 260 || snapshot.EvictedDetailLoss != 4 {
		t.Fatalf("evicted detail did not retain its entire count: %+v", snapshot)
	}
	assertOutputLossAccounting(t, snapshot)
}

func TestOutputLossConcurrentReadersAndWriters(t *testing.T) {
	var counter outputLossCounter
	var wg sync.WaitGroup
	for worker := range 8 {
		wg.Go(func() {
			for i := range 500 {
				counter.record(lossValue("scope", fmt.Sprint(worker), fmt.Sprint(i)))
				assertOutputLossAccounting(t, counter.snapshot())
			}
		})
	}
	wg.Wait()
	snapshot := counter.snapshot()
	if snapshot.Total != 4000 {
		t.Fatalf("concurrent loss count = %d, want 4000", snapshot.Total)
	}
	assertOutputLossAccounting(t, snapshot)
}

func TestOutputLossEmitEvictedValueKeepsCapturedIdentity(t *testing.T) {
	scheduler := NewScheduler(SchedulerConfig{ValueBufferSize: 1}, nil)
	old := lossValue("old-install", "old-revision", "point")
	latest := lossValue("new-install", "new-revision", "point")
	scheduler.emitValue(old)
	if got := scheduler.OutputLossDiagnostics(); got.Total != 0 {
		t.Fatalf("accepted value counted as loss: %+v", got)
	}
	// Deliberately install a different current identity: loss cannot look it up.
	scheduler.pointInfos["point"] = pointInfo{ID: "point", installationScope: "current-install", configRevision: "current-revision"}
	scheduler.emitValue(latest)
	snapshot := scheduler.OutputLossDiagnostics()
	if snapshot.Total != 1 || len(snapshot.Recent) != 1 {
		t.Fatalf("evicted value not counted: %+v", snapshot)
	}
	if got := snapshot.Recent[0]; got.InstallationScope != old.InstallationScope || got.PointRevision != old.PointRevision || got.PointID != old.PointID {
		t.Fatalf("loss borrowed current identity: %+v", got)
	}
	if got := <-scheduler.ValueChannel(); got.InstallationScope != latest.InstallationScope || got.PointRevision != latest.PointRevision {
		t.Fatalf("downstream received dropped/rewritten value: %+v", got)
	}
	select {
	case got := <-scheduler.ValueChannel():
		t.Fatalf("drop fabricated a downstream observation: %+v", got)
	default:
	}
}

func TestOutputLossEmitRejectedIncomingValue(t *testing.T) {
	scheduler := NewScheduler(SchedulerConfig{ValueBufferSize: 0}, nil)
	incoming := lossValue("captured", "revision", "point")
	scheduler.emitValue(incoming)
	snapshot := scheduler.OutputLossDiagnostics()
	if snapshot.Total != 1 || len(snapshot.Recent) != 1 || snapshot.Recent[0].InstallationScope != "captured" {
		t.Fatalf("unbuffered incoming loss not counted: %+v", snapshot)
	}
	select {
	case got := <-scheduler.ValueChannel():
		t.Fatalf("rejected incoming value fabricated delivery: %+v", got)
	default:
	}
	assertOutputLossAccounting(t, snapshot)
}
