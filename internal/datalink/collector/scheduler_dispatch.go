package collector

import (
	"context"
	"slices"
	"sync"
	"time"

	"go-gateway/internal/datalink/schema"
)

// startGroupTicker 啟動群組 Ticker (必須持有 mu 鎖)
func (s *Scheduler) startGroupTicker(group *schema.PollingGroup) {
	captured := *group
	group = &captured
	interval := time.Duration(group.IntervalMs) * time.Millisecond
	gt := &groupTicker{group: group, ticker: time.NewTicker(interval), stopCh: make(chan struct{}), done: make(chan struct{})}
	s.groupTickers[group.ID] = gt
	s.wg.Go(func() { s.runGroupTicker(gt) })
}

// runGroupTicker 運行群組 Ticker
func (s *Scheduler) runGroupTicker(gt *groupTicker) {
	defer close(gt.done)
	for {
		select {
		case <-gt.stopCh:
			return
		case <-gt.ticker.C:
			s.pollGroup(gt.group.ID, gt)
		}
	}
}

// pollGroup retains the existing per-device parallelism and group wait. Fair
// due admission/coalescing is a separate migration, not part of planned reads.
func (s *Scheduler) pollGroup(groupID string, expected ...*groupTicker) {
	snapshots := s.captureScheduledGroup(groupID, expected...)
	var wg sync.WaitGroup
	for _, snapshot := range snapshots {
		wg.Go(func() { s.pollScheduledSnapshot(snapshot) })
	}
	wg.Wait()
}

// pollDevicePoints retains the existing internal direct-call seam. Production
// tickers capture all group slots and their device configs together above.
func (s *Scheduler) pollDevicePoints(deviceID string, points []pointInfo) {
	s.mu.RLock()
	device, exists := s.deviceConfigs[deviceID]
	ctx := s.pollCtx
	if ctx == nil {
		ctx = context.Background()
	}
	snapshot := scheduledSnapshot{device: device, geometry: s.currentAcquisitionGeometryLocked(deviceID), points: slices.Clone(points), ctx: ctx, running: s.running, capturedAt: time.Now()}
	if len(points) > 0 {
		snapshot.groupID = points[0].PollingGroupID
		snapshot.ticker = s.groupTickers[snapshot.groupID]
		if snapshot.ticker != nil {
			snapshot.cadence = time.Duration(snapshot.ticker.group.IntervalMs) * time.Millisecond
		}
	}
	s.mu.RUnlock()
	if exists {
		s.pollScheduledSnapshot(snapshot)
	}
}
