package dbtarget

import (
	"sort"
	"time"

	"go-gateway/internal/datalink/schema"
)

func sameGroupedTarget(a, b groupedWriteKey) bool {
	a.BucketStart, b.BucketStart = time.Time{}, time.Time{}
	return a == b
}

func sameConnectorIdentity(a, b *schema.DatabaseConnector) bool {
	return a != nil && b != nil && a.ID == b.ID && a.Kind == b.Kind &&
		a.IdentityRevision == b.IdentityRevision && a.ConnectionConfig == b.ConnectionConfig
}

// movePendingGroup carries uncommitted members and their original deadline to
// the newest bucket. The caller holds w.mu.
func (w *Writer) movePendingGroup(previous, next groupedWriteKey, connector *schema.DatabaseConnector) groupedWriteKey {
	old := w.grouped[previous]
	if previous == next || old == nil || !sameGroupedTarget(previous, next) || !sameConnectorIdentity(old.Connector, connector) {
		return next
	}
	if next.BucketStart.Before(previous.BucketStart) {
		return previous
	}
	target := w.grouped[next]
	if target == nil {
		target = &groupedWriteBucket{Key: next, Connector: connector, Values: make(map[string]any), Mappings: make(map[string]*schema.DatabaseTargetMapping), DueAt: next.BucketStart.Add(time.Duration(next.EffectiveIntervalSeconds) * time.Second)}
		w.grouped[next] = target
	}
	if !sameConnectorIdentity(old.Connector, target.Connector) {
		return next
	}
	for column, mapping := range old.Mappings {
		if _, occupied := target.Mappings[column]; occupied {
			w.recordMappingSkip(old.Connector, mapping, old.Key.BucketStart, skipColumnReplaced)
			if latest, ok := w.groupedLocations[mapping.ID]; ok && latest == (groupedLocation{key: previous, column: column}) {
				delete(w.groupedLocations, mapping.ID)
			}
			continue
		}
		target.Values[column] = old.Values[column]
		target.Mappings[column] = mapping
		w.groupedLocations[mapping.ID] = groupedLocation{key: next, column: column}
	}
	if old.DueAt.Before(target.DueAt) {
		target.DueAt = old.DueAt
	}
	target.Retry = target.Retry || old.Retry
	delete(w.grouped, previous)
	return next
}

// restoreGrouped keeps at most the latest uncommitted value per mapping. A
// value accepted while the destination was blocked takes precedence.
func (w *Writer) restoreGrouped(bucket *groupedWriteBucket) {
	w.mu.Lock()
	defer w.mu.Unlock()
	key := bucket.Key
	for candidate, pending := range w.grouped {
		if sameGroupedTarget(candidate, key) && sameConnectorIdentity(bucket.Connector, pending.Connector) && candidate.BucketStart.After(key.BucketStart) {
			key = candidate
		}
	}
	target := w.grouped[key]
	if target == nil {
		target = &groupedWriteBucket{Key: key, Connector: bucket.Connector, Values: make(map[string]any), Mappings: make(map[string]*schema.DatabaseTargetMapping), DueAt: bucket.DueAt}
		w.grouped[key] = target
	}
	if !sameConnectorIdentity(bucket.Connector, target.Connector) {
		for column, mapping := range bucket.Mappings {
			if _, exists := bucket.Values[column]; exists {
				w.recordMappingSkip(bucket.Connector, mapping, bucket.Key.BucketStart, skipMappingNotCurrent)
			}
		}
		return
	}
	if bucket.DueAt.Before(target.DueAt) {
		target.DueAt = bucket.DueAt
	}
	// Pending updates can span several intervals while the destination is
	// blocked. Merge them first, newest bucket first, before restoring values
	// from the failed attempt.
	var candidates []groupedWriteKey
	for candidate, pending := range w.grouped {
		if candidate != key && sameGroupedTarget(candidate, key) && sameConnectorIdentity(bucket.Connector, pending.Connector) {
			candidates = append(candidates, candidate)
		}
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].BucketStart.After(candidates[j].BucketStart) })
	for _, candidate := range candidates {
		pending := w.grouped[candidate]
		if pending.DueAt.Before(target.DueAt) {
			target.DueAt = pending.DueAt
		}
		for column, mapping := range pending.Mappings {
			if _, occupied := target.Mappings[column]; occupied {
				w.recordMappingSkip(pending.Connector, mapping, pending.Key.BucketStart, skipColumnReplaced)
				if latest, ok := w.groupedLocations[mapping.ID]; ok && latest == (groupedLocation{key: candidate, column: column}) {
					delete(w.groupedLocations, mapping.ID)
				}
				continue
			}
			target.Values[column] = pending.Values[column]
			target.Mappings[column] = mapping
			w.groupedLocations[mapping.ID] = groupedLocation{key: key, column: column}
		}
		delete(w.grouped, candidate)
	}
	for column, mapping := range bucket.Mappings {
		value, exists := bucket.Values[column]
		if !exists {
			continue
		}
		if newer, exists := w.groupedLocations[mapping.ID]; exists {
			if pending := w.grouped[newer.key]; pending != nil && !sameSnapshotBinding(mapping, bucket.Connector, pending.Mappings[newer.column], pending.Connector) {
				w.recordMappingSkip(bucket.Connector, mapping, bucket.Key.BucketStart, skipMappingNotCurrent)
			}
			continue
		}
		if _, occupied := target.Mappings[column]; occupied {
			w.recordMappingSkip(bucket.Connector, mapping, bucket.Key.BucketStart, skipColumnReplaced)
			continue
		}
		target.Values[column] = value
		target.Mappings[column] = mapping
		w.groupedLocations[mapping.ID] = groupedLocation{key: key, column: column}
	}
	if len(target.Values) == 0 {
		delete(w.grouped, key)
		return
	}
	target.Retry = true
}
