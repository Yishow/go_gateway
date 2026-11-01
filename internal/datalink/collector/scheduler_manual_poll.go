package collector

import (
	"context"
	"sync"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/wirebudget"
)

// PollNow immediately polls only the requested installed points.
func (s *Scheduler) PollNow(pointIDs []string) []CollectedValue {
	return s.PollNowContext(context.Background(), pointIDs)
}

// PollNowContext uses the shared acquisition executor with caller-only delivery.
// Scheduled Stop does not own this caller's cancellation or returned results.
func (s *Scheduler) PollNowContext(ctx context.Context, pointIDs []string) []CollectedValue {
	if ctx == nil {
		ctx = context.Background()
	}
	snapshots := s.captureManualPoints(ctx, pointIDs)
	results := make([]CollectedValue, 0, len(pointIDs))
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, snapshot := range snapshots {
		wg.Go(func() {
			values := s.pollManualSnapshot(snapshot)
			mu.Lock()
			results = append(results, values...)
			mu.Unlock()
		})
	}
	wg.Wait()
	return results
}

// Capture all configuration authorities together; planning and I/O happen after
// releasing the lock. A manual intent retains original groups, not their tickers.
func (s *Scheduler) captureManualPoints(ctx context.Context, pointIDs []string) []scheduledSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	seen := make(map[string]bool, len(pointIDs))
	byDevice := map[string][]pointInfo{}
	for _, id := range pointIDs {
		if seen[id] {
			continue
		}
		seen[id] = true
		if point, ok := s.pointInfos[id]; ok && !point.recoveryPending {
			byDevice[point.DeviceID] = append(byDevice[point.DeviceID], point)
		}
	}
	capturedAt := time.Now()
	snapshots := make([]scheduledSnapshot, 0, len(byDevice))
	for id, points := range byDevice {
		device, ok := s.deviceConfigs[id]
		if !ok {
			continue
		}
		cadences := make(map[string]time.Duration)
		for _, point := range points {
			if ticker := s.groupTickers[point.PollingGroupID]; ticker != nil {
				cadences[point.PollingGroupID] = time.Duration(ticker.group.IntervalMs) * time.Millisecond
			}
		}
		snapshot := scheduledSnapshot{device: device, geometry: s.currentAcquisitionGeometryLocked(id), points: points, ctx: ctx, capturedAt: capturedAt, groupCadences: cadences}
		groupID := points[0].PollingGroupID
		sameGroup := true
		for _, point := range points {
			sameGroup = sameGroup && point.PollingGroupID == groupID
		}
		if sameGroup {
			// Identical installed sources have the same plan identity regardless
			// of the caller, while manual never borrows ticker/run ownership.
			snapshot.groupID, snapshot.cadence = groupID, cadences[groupID]
			snapshot.groupCadences = nil
		}
		snapshots = append(snapshots, snapshot)
	}
	return snapshots
}

func (s *Scheduler) pollManualSnapshot(snapshot scheduledSnapshot) []CollectedValue {
	acquisitionID := s.nextAcquisitionID()
	values := make([]CollectedValue, 0, len(snapshot.points))
	completed := make(map[string]bool, len(snapshot.points))
	deliver := func(value CollectedValue) {
		if completed[value.PointID] {
			return
		}
		completed[value.PointID] = true
		values = append(values, value)
	}
	s.pollAcquisitionSnapshot(snapshot, acquisitionID, true, deliver)
	// Scheduled cancellation suppresses publication. Manual callers instead need
	// one final Bad for each unresolved known Point, still bound to its old install.
	for _, point := range snapshot.points {
		if completed[point.ID] {
			continue
		}
		err := snapshot.ctx.Err()
		if err == nil {
			err = s.admitScheduled(snapshot, []pointInfo{point})
		}
		if err == nil {
			err = connector.ErrPlannedRead
		}
		deliver(s.collectedValueWithReason(snapshot.device, point, connector.ReadResult{Quality: schema.QualityBad}, err, acquisitionID, s.now(), "request-canceled"))
	}
	return values
}

func (s *Scheduler) ensureDeviceLocks() {
	for id := range s.deviceConfigs {
		if _, exists := s.deviceLocks[id]; !exists {
			s.deviceLocks[id] = &wirebudget.Mutex{}
		}
	}
}
