package collector

import (
	"context"
	"errors"
	"slices"
	"sync/atomic"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/wirebudget"
)

type acquisitionEvidenceKey struct{}
type acquisitionEvidenceTuple struct{ device, group string }
type acquisitionScanEvidence struct {
	snapshot scheduledSnapshot
	head     string
	metrics  AcquisitionActualMetrics
	lastGood time.Time
}

// AcquisitionCompletionObservation is a historical fact from the original
// successful publication, not callback-time source/SQL authorization. Consumers
// must not infer ordering across devices or retain mutable scheduler internals.
type AcquisitionCompletionObservation struct {
	WorkspaceID, SourceRevision, DeviceID, PollingGroupID string
	InstallationScope, IdentityHead                       string
	PointCount                                            int
	Metrics                                               AcquisitionActualMetrics
}

func (s *Scheduler) beginAcquisitionEvidence(ctx context.Context, scan scheduledPlan, id string) (context.Context, *acquisitionScanEvidence) {
	e := &acquisitionScanEvidence{snapshot: scan.snapshot, metrics: AcquisitionActualMetrics{AcquisitionID: id, DeclaredPlanRevision: scan.plan.Fingerprint(), LogicalRequests: AcquisitionLogicalOutcomes{Status: "not_started"}, Phases: []AcquisitionActualPhase{}}}
	s.setAcquisitionEvidenceInitial(e, scan, scan.plan, "initial")
	return context.WithValue(ctx, acquisitionEvidenceKey{}, e), e
}
func evidenceFrom(ctx context.Context) *acquisitionScanEvidence {
	e, ok := ctx.Value(acquisitionEvidenceKey{}).(*acquisitionScanEvidence)
	if !ok {
		return nil
	}
	return e
}

func (s *Scheduler) setAcquisitionEvidenceInitial(e *acquisitionScanEvidence, scan scheduledPlan, plan optimizer.AcquisitionPlan, phase string) {
	if e == nil {
		return
	}
	connection := s.connMgr.PeekAcquisitionConnection(previewDevice(scan.snapshot.device))
	s.mu.RLock()
	downgrade := s.acquisitionDowngrades[scan.snapshot.device.ID]
	s.mu.RUnlock()
	head, err := acquisitionGroupHead(scan, plan, connection, downgrade)
	if err != nil {
		e.head = ""
		return
	}
	e.head = head
	e.metrics.InitialPlanRevision = plan.Fingerprint()
	e.metrics.Phases = []AcquisitionActualPhase{{Phase: phase, PlanRevision: plan.Fingerprint()}}
}

// beginEvidenceRequest counts actual admitted sends, including adapter retries.
// Outcomes are recorded before executor cancellation exits or fallback breaks.
func beginEvidenceRequest(ctx context.Context, phase string) (observedContext context.Context, recordOutcome func([]connector.PlannedPointResult, error)) {
	e := evidenceFrom(ctx)
	if e == nil {
		return ctx, func([]connector.PlannedPointResult, error) {}
	}
	var sends atomic.Uint64
	ctx = wirebudget.WithSendObserver(ctx, func() { sends.Add(1) })
	return ctx, func(results []connector.PlannedPointResult, err error) {
		n := sends.Load()
		e.metrics.WireAttempts += n
		if n > 1 {
			e.metrics.RetryWireAttempts += n - 1
		}
		if phase == "fallback" {
			e.metrics.FallbackWireAttempts += n
		}
		if e.metrics.LogicalRequests.AcquisitionLogicalCounts == nil {
			e.metrics.LogicalRequests = AcquisitionLogicalOutcomes{Status: "observed", AcquisitionLogicalCounts: &AcquisitionLogicalCounts{}}
		}
		switch {
		case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
			e.metrics.LogicalRequests.Canceled++
		case err == nil && plannedResultsGood(results):
			e.metrics.LogicalRequests.Succeeded++
		default:
			e.metrics.LogicalRequests.Failed++
		}
		for _, result := range results {
			if err == nil && result.Quality == schema.QualityGood && result.Error == "" && result.ReceivedAt.After(e.lastGood) {
				e.lastGood = result.ReceivedAt
			}
		}
	}
}

func (s *Scheduler) finishAcquisitionEvidence(e *acquisitionScanEvidence) {
	if e == nil || e.head == "" || e.snapshot.cadence <= 0 {
		return
	}
	e.metrics.ObservedAt = s.now()
	elapsed := time.Since(e.snapshot.capturedAt)
	e.metrics.CycleLatencyMS = float64(max(elapsed, 0)) / float64(time.Millisecond)
	e.metrics.Overrun = elapsed > e.snapshot.cadence
	s.mu.Lock()
	if s.scheduledCurrentLocked(e.snapshot, e.snapshot.points) != nil || s.groupTickers[e.snapshot.groupID] == nil {
		s.mu.Unlock()
		return
	}
	key := acquisitionEvidenceTuple{e.snapshot.device.ID, e.snapshot.groupID}
	old := s.acquisitionEvidence[key]
	if old != nil && !old.snapshot.capturedAt.Before(e.snapshot.capturedAt) {
		s.mu.Unlock()
		return
	}
	if s.acquisitionEvidence == nil {
		s.acquisitionEvidence = map[acquisitionEvidenceTuple]*acquisitionScanEvidence{}
	}
	s.acquisitionEvidence[key] = e
	s.mu.Unlock()
	s.observeAcquisitionCompletion(e)
}

func (s *Scheduler) previewActual(snapshot scheduledSnapshot, head string) AcquisitionActual {
	out := AcquisitionActual{Status: "unavailable", Reason: "no_matching_completed_scan"}
	s.mu.RLock()
	e := s.acquisitionEvidence[acquisitionEvidenceTuple{snapshot.device.ID, snapshot.groupID}]
	if e == nil || e.head != head {
		s.mu.RUnlock()
		return out
	}
	metrics := e.metrics
	metrics.Phases = slices.Clone(metrics.Phases)
	if metrics.LogicalRequests.AcquisitionLogicalCounts != nil {
		counts := *metrics.LogicalRequests.AcquisitionLogicalCounts
		metrics.LogicalRequests.AcquisitionLogicalCounts = &counts
	}
	good := e.lastGood
	s.mu.RUnlock()
	if !good.IsZero() {
		age := float64(max(s.now().Sub(good), 0)) / float64(time.Millisecond)
		metrics.LastGoodAgeMS = &age
	}
	return AcquisitionActual{Status: "available", AcquisitionActualMetrics: &metrics}
}

func (s *Scheduler) pruneAcquisitionEvidenceLocked() {
	for key, e := range s.acquisitionEvidence {
		if s.groupTickers[key.group] == nil || s.scheduledCurrentLocked(e.snapshot, e.snapshot.points) != nil {
			delete(s.acquisitionEvidence, key)
		}
	}
}

// Called only after successful publication and after releasing scheduler.mu.
// The event remains a publication-time fact if Stop/invalidation wins before
// delivery; consumers must not treat it as callback-time source authorization.
func (s *Scheduler) observeAcquisitionCompletion(e *acquisitionScanEvidence) {
	observer := s.config.AcquisitionCompletionObserver
	if observer == nil || e == nil {
		return
	}
	metrics := e.metrics
	metrics.Phases = slices.Clone(metrics.Phases)
	if metrics.LogicalRequests.AcquisitionLogicalCounts != nil {
		counts := *metrics.LogicalRequests.AcquisitionLogicalCounts
		metrics.LogicalRequests.AcquisitionLogicalCounts = &counts
	}
	if metrics.LastGoodAgeMS != nil {
		age := *metrics.LastGoodAgeMS
		metrics.LastGoodAgeMS = &age
	}
	if !e.lastGood.IsZero() {
		age := float64(max(metrics.ObservedAt.Sub(e.lastGood), 0)) / float64(time.Millisecond)
		metrics.LastGoodAgeMS = &age
	}
	event := AcquisitionCompletionObservation{DeviceID: e.snapshot.device.ID, PollingGroupID: e.snapshot.groupID, IdentityHead: e.head, PointCount: len(e.snapshot.points), Metrics: metrics}
	if e.snapshot.geometry != nil {
		event.WorkspaceID = e.snapshot.geometry.workspaceID
		event.SourceRevision = e.snapshot.geometry.sourceRevision
	}
	if len(e.snapshot.points) > 0 {
		event.InstallationScope = e.snapshot.points[0].installationScope
	}
	// Same optional nonblocking/panic-isolated contract as the Point observer.
	// No goroutine, retained history or callback-owned authority is introduced.
	defer func() {
		if recover() != nil {
			return // Diagnostic panic cannot change published acquisition facts.
		}
	}()
	observer(event)
}
