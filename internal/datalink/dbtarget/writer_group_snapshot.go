package dbtarget

import "go-gateway/internal/datalink/schema"

// A mapping revision change invalidates every retained interval for its former
// binding, including the oldest bucket that is not the latest location.
func (w *Writer) discardChangedGrouped(mapping *schema.DatabaseTargetMapping, connector *schema.DatabaseConnector) {
	for key, bucket := range w.grouped {
		for column, previous := range bucket.Mappings {
			if previous.ID != mapping.ID || sameSnapshotBinding(previous, bucket.Connector, mapping, connector) {
				continue
			}
			w.recordMappingSkip(bucket.Connector, previous, bucket.Key.BucketStart, skipMappingNotCurrent)
			delete(bucket.Values, column)
			delete(bucket.Mappings, column)
			if latest, ok := w.groupedLocations[mapping.ID]; ok && latest == (groupedLocation{key: key, column: column}) {
				delete(w.groupedLocations, mapping.ID)
			}
		}
		if len(bucket.Values) == 0 {
			delete(w.grouped, key)
		}
	}
}

// retainGroupedInterval keeps the oldest pending bucket and the newest value
// per mapping. Intermediate snapshots are superseded without moving members
// from one interval into another. The caller holds w.mu.
func (w *Writer) retainGroupedInterval(latest groupedLocation, next groupedWriteKey, mapping *schema.DatabaseTargetMapping, connector *schema.DatabaseConnector) bool {
	oldest := latest
	for key, bucket := range w.grouped {
		if !sameGroupedTarget(key, latest.key) || !sameConnectorIdentity(bucket.Connector, connector) || !key.BucketStart.Before(oldest.key.BucketStart) {
			continue
		}
		if candidate := bucket.Mappings[mapping.ColumnName]; candidate != nil && candidate.ID == mapping.ID {
			oldest = groupedLocation{key: key, column: mapping.ColumnName}
		}
	}

	if next.BucketStart.After(latest.key.BucketStart) {
		if oldest.key != latest.key {
			w.discardGroupedLocation(latest, mapping.ID)
		}
		return true
	}
	if oldest.key == latest.key || next == oldest.key {
		return true
	}
	if next.BucketStart.Before(oldest.key.BucketStart) {
		w.discardGroupedLocation(oldest, mapping.ID)
		return true
	}
	w.recordMappingSkip(connector, mapping, next.BucketStart, skipSnapshotReplaced)
	return false
}

// discardGroupedLocation is called under w.mu.
func (w *Writer) discardGroupedLocation(location groupedLocation, mappingID string) {
	bucket := w.grouped[location.key]
	if bucket == nil {
		return
	}
	mapping := bucket.Mappings[location.column]
	if mapping == nil || mapping.ID != mappingID {
		return
	}
	w.recordMappingSkip(bucket.Connector, mapping, bucket.Key.BucketStart, skipSnapshotReplaced)
	delete(bucket.Values, location.column)
	delete(bucket.Mappings, location.column)
	if current, ok := w.groupedLocations[mappingID]; ok && current == location {
		delete(w.groupedLocations, mappingID)
	}
	if len(bucket.Values) == 0 {
		delete(w.grouped, location.key)
	}
}

// Once a write finishes, keep only each mapping's latest pending value that
// arrived during it; the bucket of every remaining value stays unchanged.
func (w *Writer) discardSupersededGrouped(completed *groupedWriteBucket) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return
	}
	for key, bucket := range w.grouped {
		if !sameGroupedTarget(key, completed.Key) || !sameConnectorIdentity(bucket.Connector, completed.Connector) {
			continue
		}
		for column, mapping := range bucket.Mappings {
			latest, ok := w.groupedLocations[mapping.ID]
			if ok && latest.key != key {
				w.discardGroupedLocation(groupedLocation{key: key, column: column}, mapping.ID)
			}
		}
	}
}
