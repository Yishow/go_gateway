package collector

import (
	"context"
	"fmt"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/wirebudget"
)

// captureScheduledGroup takes one coherent installed view. The global lock is
// released before pure planning, waiting for a device, dialing or wire I/O.
func (s *Scheduler) captureScheduledGroup(groupID string, expected ...*groupTicker) []scheduledSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ctx := s.pollCtx
	if ctx == nil {
		ctx = context.Background()
	}
	ticker := s.groupTickers[groupID]
	if len(expected) > 0 && ticker != expected[0] {
		return nil
	}
	cadence := time.Duration(0)
	if ticker != nil {
		cadence = time.Duration(ticker.group.IntervalMs) * time.Millisecond
	}
	byDevice := map[string][]pointInfo{}
	for _, point := range s.pointInfos {
		if point.PollingGroupID == groupID && !point.recoveryPending {
			byDevice[point.DeviceID] = append(byDevice[point.DeviceID], point)
		}
	}
	result := make([]scheduledSnapshot, 0, len(byDevice))
	capturedAt := time.Now()
	for id, points := range byDevice {
		if device, exists := s.deviceConfigs[id]; exists {
			result = append(result, scheduledSnapshot{device: device, geometry: s.currentAcquisitionGeometryLocked(id), points: points, groupID: groupID, cadence: cadence, ticker: ticker, ctx: ctx, running: s.running, capturedAt: capturedAt})
		}
	}
	return result
}

func (s *Scheduler) pollScheduledSnapshot(snapshot scheduledSnapshot) {
	s.pollAcquisitionSnapshot(snapshot, "", false, s.emitValue)
}

func (s *Scheduler) pollAcquisitionSnapshot(snapshot scheduledSnapshot, acquisitionID string, manual bool, deliver func(CollectedValue)) {
	if snapshot.ctx == nil {
		snapshot.ctx = context.Background()
	}
	ctx := snapshot.ctx
	var scan scheduledPlan
	var planErr error
	if boundedAcquisitionProtocol(snapshot.device.Protocol) {
		scan, planErr = buildScheduledPlan(snapshot)
		if planErr == nil && len(scan.plan.Requests()) > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithDeadline(ctx, scan.deadline)
			defer cancel()
		}
	}
	s.mu.RLock()
	gate, breaker := s.deviceLocks[snapshot.device.ID], s.deviceBreakers[snapshot.device.ID]
	capacity := s.active
	s.mu.RUnlock()
	if gate == nil || s.admitScheduled(snapshot, snapshot.points) != nil || snapshot.ctx.Err() != nil {
		return
	}
	if acquisitionID == "" {
		acquisitionID = s.nextAcquisitionID()
	}
	if !manual && planErr == nil && boundedAcquisitionProtocol(snapshot.device.Protocol) {
		var evidence *acquisitionScanEvidence
		ctx, evidence = s.beginAcquisitionEvidence(ctx, scan, acquisitionID)
		defer s.finishAcquisitionEvidence(evidence)
	}
	if s.capacityErr != nil {
		s.deliverAcquisitionFailure(snapshot, acquisitionID, s.capacityErr, "capacity-invalid", deliver)
		return
	}
	trace := acquisitionTrace{planFingerprint: scan.plan.Fingerprint(), planReason: scan.coverageReason}
	if boundedAcquisitionProtocol(snapshot.device.Protocol) {
		if planErr != nil {
			s.deliverAcquisitionFailure(snapshot, acquisitionID, planErr, "plan-invalid", deliver, trace)
			return
		}
		if len(scan.plan.Requests()) == 0 {
			for pointID, failure := range scan.failures {
				if snapshot.ctx.Err() != nil {
					return
				}
				s.reportScheduledHealth(snapshot, []pointInfo{scan.points[pointID]}, breaker, failure)
				deliver(s.collectedValueWithReason(snapshot.device, scan.points[pointID], connector.ReadResult{Quality: schema.QualityBad}, failure, acquisitionID, s.now(), "plan-invalid", acquisitionTrace{planFingerprint: scan.plan.Fingerprint(), planReason: scan.coverageReason}))
			}
			return
		}
	}
	if breaker != nil && !breaker.AllowRequest() {
		s.deliverAcquisitionFailure(snapshot, acquisitionID, fmt.Errorf("設備熔斷中 (狀態: %s)", breaker.State()), "circuit-open", deliver, trace)
		return
	}
	if err := gate.LockContext(ctx); err != nil {
		if snapshot.ctx.Err() == nil && s.admitScheduled(snapshot, snapshot.points) == nil {
			s.deliverAcquisitionFailure(snapshot, acquisitionID, err, "scan-deadline", deliver, trace)
		}
		return
	}
	defer gate.Unlock()
	if capacity != nil {
		lease, err := capacity.acquire(ctx)
		if err != nil {
			if snapshot.ctx.Err() == nil && s.admitScheduled(snapshot, snapshot.points) == nil {
				s.deliverAcquisitionFailure(snapshot, acquisitionID, err, "scan-deadline", deliver, trace)
			}
			return
		}
		defer lease.release()
		ctx = wirebudget.WithWaitCooperator(ctx, lease)
	}
	// Never let a stale queued caller replace the manager's current config.
	if ctx.Err() != nil {
		if snapshot.ctx.Err() == nil && s.admitScheduled(snapshot, snapshot.points) == nil {
			s.deliverAcquisitionFailure(snapshot, acquisitionID, ctx.Err(), "scan-deadline", deliver, trace)
		}
		return
	}
	if s.admitScheduled(snapshot, snapshot.points) != nil {
		return
	}
	conn, err := s.connMgr.GetOrCreate(ctx, snapshot.device.ID, snapshot.device.Protocol, snapshot.device.Config)
	if err != nil {
		if evidence := evidenceFrom(ctx); evidence != nil {
			s.setAcquisitionEvidenceInitial(evidence, scan, scan.plan, "initial")
		}
		if snapshot.ctx.Err() != nil {
			return
		}
		s.reportScheduledHealth(snapshot, snapshot.points, breaker, err)
		s.deliverAcquisitionFailure(snapshot, acquisitionID, fmt.Errorf("連線失敗: %w", err), "connection-failed", deliver, trace)
		return
	}
	if boundedAcquisitionProtocol(snapshot.device.Protocol) {
		if _, planned := conn.Protocol.(connector.PlannedReader); !planned {
			// A replaced/legacy adapter cannot escape declared restrictions by
			// silently running the unplanned scalar path on the entire Point set.
			for _, point := range snapshot.points {
				s.reportScheduledHealth(snapshot, []pointInfo{point}, breaker, connector.ErrPlannedRead)
			}
			s.deliverAcquisitionFailure(snapshot, acquisitionID, connector.ErrPlannedRead, "planned-read-unavailable", deliver, trace)
			return
		}
		s.executeAcquisitionPlan(ctx, conn, scan, acquisitionID, breaker, deliver)
		return
	}
	// Other installed protocols retain their original scalar executor. This is
	// not a capability probe or a fallback for a failed supported PLC batch.
	for _, point := range snapshot.points {
		if ctx.Err() != nil || s.admitScheduled(snapshot, []pointInfo{point}) != nil {
			return
		}
		attempts := s.config.DefaultRetryCount
		if manual {
			attempts = 0
		}
		err := s.pollScalarPoint(ctx, conn, point, acquisitionID, attempts, deliver)
		s.reportScheduledHealth(snapshot, []pointInfo{point}, breaker, err)
	}
}

func (s *Scheduler) deliverAcquisitionFailure(snapshot scheduledSnapshot, acquisitionID string, err error, reason string, deliver func(CollectedValue), trace ...acquisitionTrace) {
	completed := s.now()
	for _, point := range snapshot.points {
		if snapshot.ctx.Err() != nil {
			return
		}
		deliver(s.collectedValueWithReason(snapshot.device, point, connector.ReadResult{Quality: schema.QualityBad, Error: err.Error()}, nil, acquisitionID, completed, reason, trace...))
	}
}
