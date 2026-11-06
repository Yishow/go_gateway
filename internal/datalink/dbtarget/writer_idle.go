package dbtarget

import (
	"time"

	"go-gateway/internal/datalink/schema"
)

const noNewValueReason = "no_new_value"

type idleIntervalState struct {
	boundary time.Time
	seconds  int
}

type idleMapping struct {
	mapping   *schema.DatabaseTargetMapping
	connector *schema.DatabaseConnector
}

// Only interval occupancy is retained, never a queue of values. The bound
// matches the existing per-connector in-memory history limit.
type observedIntervalSet struct {
	seconds int
	buckets map[time.Time]struct{}
}

// markObserved is called under w.mu.
func (w *Writer) markObserved(mappingID string, seconds int, bucket time.Time) {
	observed := w.observed[mappingID]
	if observed == nil || observed.seconds != seconds {
		observed = &observedIntervalSet{seconds: seconds, buckets: make(map[time.Time]struct{})}
		w.observed[mappingID] = observed
	}
	observed.buckets[bucket] = struct{}{}
	if len(observed.buckets) > 1000 {
		var oldest time.Time
		for candidate := range observed.buckets {
			if oldest.IsZero() || candidate.Before(oldest) {
				oldest = candidate
			}
		}
		delete(observed.buckets, oldest)
	}
}

// takeObserved returns and removes completed occupied intervals. Future ones
// remain available if a worker receives a delayed tick.
func (w *Writer) takeObserved(mappingID string, seconds int, before time.Time) map[time.Time]struct{} {
	w.mu.Lock()
	defer w.mu.Unlock()
	observed := w.observed[mappingID]
	if observed == nil || observed.seconds != seconds {
		return nil
	}
	completed := make(map[time.Time]struct{})
	for bucket := range observed.buckets {
		if bucket.Before(before) {
			completed[bucket] = struct{}{}
			delete(observed.buckets, bucket)
		}
	}
	return completed
}

func (w *Writer) idleMappings(connectorID string) []idleMapping {
	w.indexMu.RLock()
	var mappings []idleMapping
	for _, entries := range w.index {
		for _, entry := range entries {
			if entry.Err == nil && entry.Mapping.ConnectorID == connectorID {
				mappings = append(mappings, idleMapping{mapping: entry.Mapping, connector: entry.Connector})
			}
		}
	}
	w.indexMu.RUnlock()

	return mappings
}

func (w *Writer) initialIdleIntervals(connectorID string, now time.Time) map[string]idleIntervalState {
	states := make(map[string]idleIntervalState)
	for _, entry := range w.idleMappings(connectorID) {
		seconds := effectiveWriteIntervalSeconds(entry.mapping, entry.connector)
		if seconds > 0 {
			states[entry.mapping.ID] = idleIntervalState{boundary: now.UTC().Truncate(time.Duration(seconds) * time.Second), seconds: seconds}
		}
	}
	return states
}

// recordIdleIntervals logs complete intervals with no newly accepted value.
// A retry of an earlier interval remains a separate history event.
func (w *Writer) recordIdleIntervals(connectorID string, worker *connectorWorker, now time.Time) {
	active := make(map[string]struct{})
	for _, entry := range w.idleMappings(connectorID) {
		mapping := entry.mapping
		seconds := effectiveWriteIntervalSeconds(mapping, entry.connector)
		if seconds <= 0 {
			continue
		}
		active[mapping.ID] = struct{}{}
		interval := time.Duration(seconds) * time.Second
		boundary := now.UTC().Truncate(interval)
		state, exists := worker.idle[mapping.ID]
		if !exists || state.seconds != seconds {
			worker.idle[mapping.ID] = idleIntervalState{boundary: boundary, seconds: seconds}
			continue
		}
		if !boundary.After(state.boundary) {
			continue
		}
		occupied := w.takeObserved(mapping.ID, seconds, boundary)
		if w.suppress == nil || !w.suppress(connectorID, mapping.TagID) {
			// The in-memory history keeps at most 1000 entries per connector.
			count := min(int(boundary.Sub(state.boundary)/interval), 1000)
			for i := count - 1; i >= 0; i-- {
				at := boundary.Add(-time.Duration(i) * interval)
				if _, received := occupied[at.Add(-interval)]; received {
					continue
				}
				attemptedAt := now.UTC()
				recordWriteHistory(connectorID, WriteHistoryRecord{
					ObservedAt: at, AttemptedAt: &attemptedAt, Status: "skipped", MappingID: mapping.ID,
					GroupKey: cloneOptionalString(mapping.GroupKey), TableName: mapping.TableName,
					EffectiveIntervalSeconds: seconds, SkipReason: noNewValueReason,
				})
			}
		}
		worker.idle[mapping.ID] = idleIntervalState{boundary: boundary, seconds: seconds}
	}
	for mappingID := range worker.idle {
		if _, exists := active[mappingID]; !exists {
			delete(worker.idle, mappingID)
		}
	}
}
