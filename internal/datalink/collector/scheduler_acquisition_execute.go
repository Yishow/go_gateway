package collector

import (
	"context"
	"errors"

	"go-gateway/internal/datalink/collector/health"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
)

type acquisitionDowngrade struct {
	token              connector.AcquisitionConnectionToken
	capabilityRevision string
}

func (s *Scheduler) executeScheduledPlan(ctx context.Context, conn *connector.ManagedConnection, scan scheduledPlan, acquisitionID string, breaker *health.CircuitBreaker) {
	s.executeAcquisitionPlan(ctx, conn, scan, acquisitionID, breaker, s.emitValue)
}

func (s *Scheduler) executeAcquisitionPlan(ctx context.Context, conn *connector.ManagedConnection, scan scheduledPlan, acquisitionID string, breaker *health.CircuitBreaker, deliver func(CollectedValue)) {
	completed := make(map[string]bool, len(scan.points))
	for pointID, err := range scan.failures {
		if scan.snapshot.ctx.Err() != nil {
			return
		}
		pt := scan.points[pointID]
		s.reportScheduledHealth(scan.snapshot, []pointInfo{pt}, breaker, err)
		deliver(s.collectedValueWithReason(scan.snapshot.device, pt, connector.ReadResult{Quality: schema.QualityBad}, err, acquisitionID, s.now(), "plan-invalid", acquisitionTrace{planFingerprint: scan.plan.Fingerprint(), planReason: scan.coverageReason}))
		completed[pointID] = true
	}
	plan, phase := scan.plan, "initial"
	if s.scheduledDowngraded(ctx, conn, scan) {
		plan = optimizer.BuildAcquisitionPlan(safeAcquisitionInput(scan.input, completed))
		phase = "safe"
	}
	s.setAcquisitionEvidenceInitial(evidenceFrom(ctx), scan, plan, phase)
	fallbackUsed := false
	for {
		fallback := false
		for index, request := range plan.Requests() {
			if scan.snapshot.ctx.Err() != nil {
				return
			}
			requestCtx, recordOutcome := beginEvidenceRequest(ctx, phase)
			results, err := s.readScheduledRequest(requestCtx, conn, scan, plan, index, acquisitionID, phase)
			recordOutcome(results, err)
			requestPoints := make([]pointInfo, 0, len(request.Points))
			for _, point := range request.Points {
				requestPoints = append(requestPoints, scan.points[point.Point.ID])
			}
			if scan.snapshot.ctx.Err() != nil {
				return
			}
			if request.Kind == acquisitionRequestRandom && !fallbackUsed && explicitCommandRejection(err) && len(results) > 0 {
				// A parsed unsupported-command response authorizes at most one safe
				// replan. An arbitrary timeout, bad frame or address error does not.
				s.rememberScheduledDowngrade(ctx, conn, scan, results[0].ConnectionToken)
				fallbackUsed, fallback = true, true
				break
			}
			for _, result := range results {
				id := result.Binding.Point.ID
				if completed[id] {
					continue
				}
				pt := scan.points[id]
				value := s.collectedValueWithReason(scan.snapshot.device, pt, result.ReadResult, err, acquisitionID, result.ReceivedAt, "", acquisitionTrace{requestID: result.RequestID, planFingerprint: result.PlanFingerprint, planReason: scan.coverageReason})
				if breaker != nil {
					pointErr := err
					if pointErr == nil && result.Error != "" {
						pointErr = errors.New(result.Error)
					}
					s.reportScheduledHealth(scan.snapshot, requestPoints, breaker, pointErr)
				}
				deliver(value)
				completed[id] = true
			}
		}
		if !fallback {
			return
		}
		plan = optimizer.BuildAcquisitionPlan(safeAcquisitionInput(scan.input, completed))
		phase = "fallback"
		if evidence := evidenceFrom(ctx); evidence != nil {
			evidence.metrics.Phases = append(evidence.metrics.Phases, AcquisitionActualPhase{Phase: phase, PlanRevision: plan.Fingerprint()})
		}
		// Safe planning may reject a scalar impossible under its contiguous cap.
		// It receives one final bad result, never a fabricated value or a retry plan.
		for _, failure := range plan.Errors() {
			if completed[failure.PointID] || scan.snapshot.ctx.Err() != nil {
				continue
			}
			pt := scan.points[failure.PointID]
			s.reportScheduledHealth(scan.snapshot, []pointInfo{pt}, breaker, errors.New(failure.Code))
			deliver(s.collectedValueWithReason(scan.snapshot.device, pt, connector.ReadResult{Quality: schema.QualityBad}, errors.New(failure.Code), acquisitionID, s.now(), "fallback-invalid", acquisitionTrace{planFingerprint: plan.Fingerprint(), planReason: scan.coverageReason}))
			completed[failure.PointID] = true
		}
	}
}

// Health keeps its historical Point-result unit, not a new wire-request unit.
// Eligibility covers the whole request: rejecting a batch after one member was
// removed is not a device failure for its unchanged siblings. The established
// scheduler->breaker lock order contains no I/O or callback into the scheduler.
func (s *Scheduler) reportScheduledHealth(snapshot scheduledSnapshot, points []pointInfo, breaker *health.CircuitBreaker, err error) {
	if breaker == nil || errors.Is(err, errAcquisitionSuperseded) || errors.Is(err, errActiveCapacityWait) {
		return
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if (snapshot.ctx != nil && snapshot.ctx.Err() != nil) || s.scheduledCurrentLocked(snapshot, points) != nil {
		return
	}
	breaker.ReportResult(err)
}

func (s *Scheduler) scheduledDowngraded(ctx context.Context, conn *connector.ManagedConnection, scan scheduledPlan) bool {
	token, err := conn.CurrentAcquisitionToken(ctx)
	if err != nil || !token.Valid() {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	downgrade, exists := s.acquisitionDowngrades[scan.snapshot.device.ID]
	return exists && downgrade.token == token && downgrade.capabilityRevision == scan.capabilityRevision
}

func (s *Scheduler) rememberScheduledDowngrade(ctx context.Context, conn *connector.ManagedConnection, scan scheduledPlan, finalToken connector.AcquisitionConnectionToken) {
	current, err := conn.CurrentAcquisitionToken(ctx)
	if err != nil || !finalToken.Valid() || current != finalToken {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.deviceConfigs[scan.snapshot.device.ID] != scan.snapshot.device {
		return
	}
	old := s.acquisitionDowngrades[scan.snapshot.device.ID]
	// A late result must not replace a downgrade from a newer successful connect.
	if old.token.ManagerGeneration > finalToken.ManagerGeneration || (old.token.ManagerGeneration == finalToken.ManagerGeneration && old.token.TransportGeneration > finalToken.TransportGeneration) {
		return
	}
	if s.acquisitionDowngrades == nil {
		s.acquisitionDowngrades = make(map[string]acquisitionDowngrade)
	}
	s.acquisitionDowngrades[scan.snapshot.device.ID] = acquisitionDowngrade{finalToken, scan.capabilityRevision}
}
