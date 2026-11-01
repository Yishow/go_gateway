package collector

import (
	"sync"
	"time"
)

// DueQueueDiagnostics reports this scheduler run's due admission, not result loss.
// Overruns counts finished scans whose capture-to-completion exceeded cadence.
type DueQueueDiagnostics struct {
	Pending, Running              int
	Coalesced, Rejected, Overruns uint64
}

// dueQueue contains only admitted, undispatched tuples. Its one container per
// group never overwrites a different group's captured set. No scan runs here.
type dueQueue struct {
	mu           sync.Mutex
	limit        int
	groups       map[string]map[string]scheduledSnapshot
	deviceGroups map[string][]string
	ready        []string
	busy         map[string]scheduledSnapshot
	stats        DueQueueDiagnostics
	stopped      bool
	wake         chan struct{}
}

func newDueQueue(limit int) *dueQueue {
	return &dueQueue{limit: limit, groups: make(map[string]map[string]scheduledSnapshot), deviceGroups: make(map[string][]string), busy: make(map[string]scheduledSnapshot), wake: make(chan struct{}, 1)}
}

func (q *dueQueue) notify() {
	select {
	case q.wake <- struct{}{}:
	default:
	}
}

func (q *dueQueue) offer(snapshot scheduledSnapshot) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.stopped {
		return false
	}
	device, group := snapshot.device.ID, snapshot.groupID
	if _, exists := q.groups[group][device]; exists {
		q.groups[group][device] = snapshot
		q.stats.Coalesced++
		return true
	}
	if q.stats.Pending >= q.limit {
		q.stats.Rejected++
		return false
	}
	if q.groups[group] == nil {
		q.groups[group] = make(map[string]scheduledSnapshot)
	}
	q.groups[group][device] = snapshot
	if len(q.deviceGroups[device]) == 0 {
		q.ready = append(q.ready, device)
	}
	q.deviceGroups[device] = append(q.deviceGroups[device], group)
	q.stats.Pending++
	q.notify()
	return true
}

func (q *dueQueue) next() (scheduledSnapshot, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.stopped {
		return scheduledSnapshot{}, false
	}
	for range len(q.ready) {
		device := q.ready[0]
		q.ready = q.ready[1:]
		if _, running := q.busy[device]; running {
			q.ready = append(q.ready, device)
			continue
		}
		group := q.deviceGroups[device][0]
		q.deviceGroups[device] = q.deviceGroups[device][1:]
		if len(q.deviceGroups[device]) > 0 {
			q.ready = append(q.ready, device)
		} else {
			delete(q.deviceGroups, device)
		}
		snapshot := q.groups[group][device]
		delete(q.groups[group], device)
		if len(q.groups[group]) == 0 {
			delete(q.groups, group)
		}
		q.stats.Pending--
		q.busy[device] = snapshot
		return snapshot, true
	}
	return scheduledSnapshot{}, false
}

func (q *dueQueue) complete(device string) {
	q.mu.Lock()
	if snapshot, exists := q.busy[device]; exists && snapshot.cadence > 0 && time.Since(snapshot.capturedAt) > snapshot.cadence {
		q.stats.Overruns++
	}
	delete(q.busy, device)
	q.mu.Unlock()
	q.notify()
}

func (q *dueQueue) stop() {
	q.mu.Lock()
	q.stopped = true
	clear(q.groups)
	clear(q.deviceGroups)
	q.ready = nil
	q.stats.Pending = 0
	q.mu.Unlock()
	q.notify()
}

func (q *dueQueue) diagnostics() DueQueueDiagnostics {
	q.mu.Lock()
	defer q.mu.Unlock()
	stats := q.stats
	stats.Running = len(q.busy)
	return stats
}
