package collector

import (
	"context"
	"fmt"
	"slices"
)

// Design defaults, not measured throughput or optimal deployment settings.
const defaultMaxActiveScans = 32
const defaultPendingScanCapacity = 1024

func schedulerCapacities(config SchedulerConfig) (active, pending int, err error) {
	active, pending = defaultMaxActiveScans, defaultPendingScanCapacity
	if config.MaxActiveScans != nil {
		active = *config.MaxActiveScans
	}
	if config.PendingScanCapacity != nil {
		pending = *config.PendingScanCapacity
	}
	if active <= 0 || pending <= 0 {
		return 0, 0, fmt.Errorf("scheduler active and pending capacities must be positive")
	}
	return active, pending, nil
}

// DueDiagnostics is a detached run-scoped view; output-result loss is separate.
func (s *Scheduler) DueDiagnostics() DueQueueDiagnostics {
	s.mu.RLock()
	queue := s.due
	s.mu.RUnlock()
	if queue == nil {
		return DueQueueDiagnostics{}
	}
	return queue.diagnostics()
}

func (s *Scheduler) enqueueDueGroup(groupID string, expected *groupTicker) {
	snapshots := s.captureScheduledGroup(groupID, expected)
	// Stable order plus a per-ticker starting cursor prevents finite pending
	// capacity from permanently privileging the first sorted device each tick.
	slices.SortFunc(snapshots, func(a, b scheduledSnapshot) int { return comparePointID(a.device.ID, b.device.ID) })
	if len(snapshots) > 0 && expected != nil {
		offset := (expected.admissionCursor.Add(1) - 1) % uint64(len(snapshots))
		snapshots = append(snapshots[offset:len(snapshots):len(snapshots)], snapshots[:offset]...)
	}
	s.enqueueDueSnapshots(snapshots)
}

func (s *Scheduler) enqueueDueSnapshots(snapshots []scheduledSnapshot) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	queue := s.due
	if queue == nil || !s.running {
		return
	}
	for _, snapshot := range snapshots {
		// Short bookkeeping follows scheduler -> queue lock order; queue never
		// calls back into scheduler authority. A delayed old producer cannot replace
		// a newer ticker/config/installation's already-admitted pending snapshot.
		if snapshot.ctx.Err() == nil && s.scheduledCurrentLocked(snapshot, snapshot.points) == nil {
			queue.offer(snapshot)
		}
	}
}

func (s *Scheduler) runDueDispatcher(ctx context.Context, queue *dueQueue) {
	for {
		if ctx.Err() != nil {
			return
		}
		if snapshot, ok := queue.next(); ok {
			s.wg.Go(func() { //nolint:contextcheck // snapshot.ctx retains this exact run's cancellation and admission identity.
				defer queue.complete(snapshot.device.ID)
				s.pollScheduledSnapshot(snapshot)
			})
			continue
		}
		select {
		case <-ctx.Done():
			return
		case <-queue.wake:
		}
	}
}
