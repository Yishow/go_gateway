package collector

import (
	"context"
	"fmt"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
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
	s.mu.RUnlock()
	if gate == nil || s.admitScheduled(snapshot, snapshot.points) != nil || snapshot.ctx.Err() != nil {
		return
	}
	acquisitionID := s.nextAcquisitionID()
	trace := acquisitionTrace{planFingerprint: scan.plan.Fingerprint(), planReason: scan.coverageReason}
	if boundedAcquisitionProtocol(snapshot.device.Protocol) {
		if planErr != nil {
			s.emitScheduledFailure(snapshot, acquisitionID, planErr, "plan-invalid", trace)
			return
		}
		if len(scan.plan.Requests()) == 0 {
			for pointID, failure := range scan.failures {
				if snapshot.ctx.Err() != nil {
					return
				}
				s.reportScheduledHealth(snapshot, []pointInfo{scan.points[pointID]}, breaker, failure)
				s.emitValue(s.collectedValueWithReason(snapshot.device, scan.points[pointID], connector.ReadResult{Quality: schema.QualityBad}, failure, acquisitionID, s.now(), "plan-invalid", acquisitionTrace{planFingerprint: scan.plan.Fingerprint(), planReason: scan.coverageReason}))
			}
			return
		}
	}
	if breaker != nil && !breaker.AllowRequest() {
		s.emitScheduledFailure(snapshot, acquisitionID, fmt.Errorf("設備熔斷中 (狀態: %s)", breaker.State()), "circuit-open", trace)
		return
	}
	if err := gate.LockContext(ctx); err != nil {
		if snapshot.ctx.Err() == nil && s.admitScheduled(snapshot, snapshot.points) == nil {
			s.emitScheduledFailure(snapshot, acquisitionID, err, "scan-deadline", trace)
		}
		return
	}
	defer gate.Unlock()
	// Never let a stale queued caller replace the manager's current config.
	if ctx.Err() != nil {
		if snapshot.ctx.Err() == nil && s.admitScheduled(snapshot, snapshot.points) == nil {
			s.emitScheduledFailure(snapshot, acquisitionID, ctx.Err(), "scan-deadline", trace)
		}
		return
	}
	if s.admitScheduled(snapshot, snapshot.points) != nil {
		return
	}
	conn, err := s.connMgr.GetOrCreate(ctx, snapshot.device.ID, snapshot.device.Protocol, snapshot.device.Config)
	if err != nil {
		if snapshot.ctx.Err() != nil {
			return
		}
		s.reportScheduledHealth(snapshot, snapshot.points, breaker, err)
		s.emitScheduledFailure(snapshot, acquisitionID, fmt.Errorf("連線失敗: %w", err), "connection-failed", trace)
		return
	}
	if boundedAcquisitionProtocol(snapshot.device.Protocol) {
		if _, planned := conn.Protocol.(connector.PlannedReader); !planned {
			// A replaced/legacy adapter cannot escape declared restrictions by
			// silently running the unplanned scalar path on the entire Point set.
			for _, point := range snapshot.points {
				s.reportScheduledHealth(snapshot, []pointInfo{point}, breaker, connector.ErrPlannedRead)
			}
			s.emitScheduledFailure(snapshot, acquisitionID, connector.ErrPlannedRead, "planned-read-unavailable", trace)
			return
		}
		s.executeScheduledPlan(ctx, conn, scan, acquisitionID, breaker)
		return
	}
	// Other installed protocols retain their original scalar executor. This is
	// not a capability probe or a fallback for a failed supported PLC batch.
	for _, point := range snapshot.points {
		if ctx.Err() != nil || s.admitScheduled(snapshot, []pointInfo{point}) != nil {
			return
		}
		err := s.pollPointWithAcquisition(ctx, conn, point, acquisitionID)
		s.reportScheduledHealth(snapshot, []pointInfo{point}, breaker, err)
	}
}

func (s *Scheduler) emitScheduledFailure(snapshot scheduledSnapshot, acquisitionID string, err error, reason string, trace ...acquisitionTrace) {
	completed := s.now()
	for _, point := range snapshot.points {
		if snapshot.ctx.Err() != nil {
			return
		}
		s.emitValue(s.collectedValueWithReason(snapshot.device, point, connector.ReadResult{Quality: schema.QualityBad, Error: err.Error()}, nil, acquisitionID, completed, reason, trace...))
	}
}
