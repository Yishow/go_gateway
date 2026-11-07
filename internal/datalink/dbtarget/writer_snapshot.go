package dbtarget

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"time"

	"go-gateway/internal/datalink/schema"
)

// indexedMapping is one live mapping resolved with its enabled connector. Err
// is set when the connector could not be loaded, so ingest can surface the
// failure without touching any database.
type indexedMapping struct {
	Mapping   *schema.DatabaseTargetMapping
	Connector *schema.DatabaseConnector
	Err       error
}

// pendingWrite is the latest value of one ungrouped mapping that has not been
// written yet. A newer value replaces it; values are never queued.
type pendingWrite struct {
	Mapping         *schema.DatabaseTargetMapping
	Connector       *schema.DatabaseConnector
	IntervalSeconds int
	Value           any
	ObservedAt      time.Time
	DueAt           time.Time
	Retry           bool
}

func (p *pendingWrite) dueAt() time.Time {
	return p.DueAt
}

// RefreshMappings reloads the live mapping index used by ingest. The writer
// also refreshes it on every flush tick, so mapping changes take effect within
// one tick; callers may refresh explicitly after a change. On failure the
// previous index stays in use and the error is returned.
func (w *Writer) RefreshMappings(ctx context.Context) error {
	w.refreshActive.Add(1)
	w.refreshVersion.Add(1)
	defer func() {
		w.refreshVersion.Add(1)
		w.refreshActive.Add(-1)
	}()
	w.refreshMu.Lock()
	defer w.refreshMu.Unlock()
	if w.mappingRepo == nil || w.connectorRepo == nil {
		return fmt.Errorf("資料庫目標儲存庫未設定")
	}
	var mappings []*schema.DatabaseTargetMapping
	err := retryWriterConfigBusy(ctx, func() error {
		var err error
		mappings, err = w.mappingRepo.List(ctx, TargetMappingListFilter{})
		return err
	})
	var projection *liveTargetProjection
	if err == nil {
		err = retryWriterConfigBusy(ctx, func() error {
			reader, err := newRefreshTagReader(ctx, w.tagReader)
			if err != nil {
				return err
			}
			projection, err = filterLiveTargetProjection(ctx, mappings, reader)
			return err
		})
	}
	if err != nil {
		wrapped := fmt.Errorf("重新載入資料庫目標映射失敗: %w", err)
		w.indexMu.Lock()
		if w.baseCtx.Err() == nil {
			w.indexErr = wrapped
		}
		w.indexMu.Unlock()
		return wrapped
	}
	var connectorList []*schema.DatabaseConnector
	err = retryWriterConfigBusy(ctx, func() error {
		var err error
		connectorList, err = w.connectorRepo.List(ctx, ConnectorListFilter{})
		return err
	})
	if err != nil {
		wrapped := fmt.Errorf("重新載入啟用的資料庫連接器失敗: %w", err)
		w.indexMu.Lock()
		if w.baseCtx.Err() == nil {
			w.indexErr = wrapped
		}
		w.indexMu.Unlock()
		return wrapped
	}
	type connectorLookup struct {
		connector *schema.DatabaseConnector
		err       error
	}
	connectors := make(map[string]connectorLookup)
	connectorIDs := make(map[string]struct{}, len(connectorList))
	for _, connector := range connectorList {
		if connector != nil {
			connectors[connector.ID] = connectorLookup{connector: connector}
			if connector.Enabled {
				connectorIDs[connector.ID] = struct{}{}
			}
		}
	}
	index := make(map[string][]indexedMapping, len(projection.Mappings))
	liveMappingIDs := make(map[string]struct{}, len(mappings))
	for _, mapping := range mappings {
		if mapping != nil {
			liveMappingIDs[mapping.ID] = struct{}{}
		}
	}
	for _, mapping := range projection.Mappings {
		lookup, ok := connectors[mapping.ConnectorID]
		if !ok {
			lookup.err = retryWriterConfigBusy(ctx, func() error {
				var err error
				lookup.connector, err = w.connectorRepo.GetByID(ctx, mapping.ConnectorID)
				return err
			})
			if lookup.err == nil && lookup.connector == nil {
				lookup.err = fmt.Errorf("連接器 %s 不存在", mapping.ConnectorID)
			}
			connectors[mapping.ConnectorID] = lookup
		}
		if lookup.err == nil && !lookup.connector.Enabled {
			continue
		}
		index[mapping.TagID] = append(index[mapping.TagID], indexedMapping{
			Mapping: mapping, Connector: lookup.connector, Err: lookup.err,
		})
	}

	w.indexMu.Lock()
	w.index = index
	w.connectorIDs = connectorIDs
	w.indexErr = nil
	w.indexMu.Unlock()
	w.mu.Lock()
	for mappingID := range w.lastFlush {
		if _, exists := liveMappingIDs[mappingID]; !exists {
			delete(w.lastFlush, mappingID)
			delete(w.observed, mappingID)
		}
	}
	for mappingID := range w.observed {
		if _, exists := liveMappingIDs[mappingID]; !exists {
			delete(w.observed, mappingID)
		}
	}
	w.mu.Unlock()
	return nil
}

func (w *Writer) lookupMappings(tagID string) ([]indexedMapping, error) {
	w.indexMu.RLock()
	defer w.indexMu.RUnlock()
	return w.index[tagID], w.indexErr
}

// indexedConnectorIDs lists all enabled connectors from the latest refresh.
func (w *Writer) indexedConnectorIDs() map[string]struct{} {
	w.indexMu.RLock()
	defer w.indexMu.RUnlock()
	return maps.Clone(w.connectorIDs)
}

// storeLatest overwrites the pending value of an ungrouped mapping.
func (w *Writer) storeLatest(entry indexedMapping, value any, observedAt time.Time) error {
	interval := effectiveWriteIntervalSeconds(entry.Mapping, entry.Connector)
	if interval <= 0 {
		return fmt.Errorf("寫入間隔必須為正整數秒")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return fmt.Errorf("database target writer is closed")
	}
	byMapping := w.pending[entry.Mapping.ConnectorID]
	if byMapping == nil {
		byMapping = make(map[string]*pendingWrite)
		w.pending[entry.Mapping.ConnectorID] = byMapping
	}
	intervalDuration := time.Duration(interval) * time.Second
	due := observedAt.UTC().Truncate(intervalDuration).Add(intervalDuration)
	retry := false
	if previous := byMapping[entry.Mapping.ID]; previous != nil {
		if sameSnapshotBinding(previous.Mapping, previous.Connector, entry.Mapping, entry.Connector) {
			due, retry = previous.DueAt, previous.Retry
		} else {
			w.recordMappingSkip(previous.Connector, previous.Mapping, previous.ObservedAt, skipMappingNotCurrent)
		}
	}
	byMapping[entry.Mapping.ID] = &pendingWrite{
		Mapping: entry.Mapping, Connector: entry.Connector, IntervalSeconds: interval,
		Value: value, ObservedAt: observedAt, DueAt: due, Retry: retry,
	}
	bucket := observedAt.UTC().Truncate(intervalDuration)
	if last, ok := w.lastFlush[entry.Mapping.ID]; ok && !bucket.Add(intervalDuration).After(last) {
		bucket = last
	}
	w.markObserved(entry.Mapping.ID, interval, bucket)
	return nil
}

// connectorsWithPendingData lists connectors that still hold unwritten values.
func (w *Writer) connectorsWithPendingData() map[string]struct{} {
	w.mu.Lock()
	defer w.mu.Unlock()
	ids := make(map[string]struct{})
	for id, byMapping := range w.pending {
		if len(byMapping) > 0 {
			ids[id] = struct{}{}
		}
	}
	for key := range w.grouped {
		ids[key.ConnectorID] = struct{}{}
	}
	return ids
}

func (w *Writer) takeDuePending(connectorID string, now time.Time, all bool) *pendingWrite {
	w.mu.Lock()
	defer w.mu.Unlock()
	byMapping := w.pending[connectorID]
	ids := slices.Sorted(maps.Keys(byMapping))
	slices.SortStableFunc(ids, func(a, b string) int {
		if byMapping[a].Retry == byMapping[b].Retry {
			return 0
		}
		if byMapping[a].Retry {
			return -1
		}
		return 1
	})
	for _, mappingID := range ids {
		write := byMapping[mappingID]
		due := write.dueAt()
		interval := time.Duration(write.IntervalSeconds) * time.Second
		last, wasFlushed := w.lastFlush[mappingID]
		if wasFlushed && !due.After(last) {
			due = last.Add(interval)
		}
		if !write.Retry && now.Before(due) && (!all || wasFlushed) {
			continue
		}
		delete(byMapping, mappingID)
		if len(byMapping) == 0 {
			delete(w.pending, connectorID)
		}
		w.lastFlush[mappingID] = now.UTC().Truncate(interval)
		return write
	}
	return nil
}

func (w *Writer) writePending(ctx context.Context, write *pendingWrite, retryDelay time.Duration) (bool, error) {
	mapping, connector := write.Mapping, write.Connector
	current, err := w.mappingIsCurrent(ctx, mapping, connector)
	if err != nil {
		if ctx.Err() != nil && w.baseCtx.Err() != nil {
			w.restorePending(write)
			return false, err
		}
		w.recordMappingSkip(connector, mapping, write.ObservedAt, skipMappingCheckError)
		return false, err
	}
	if !current {
		w.recordMappingSkip(connector, mapping, write.ObservedAt, skipMappingNotCurrent)
		return false, nil
	}
	if w.suppress != nil && w.suppress(mapping.ConnectorID, mapping.TagID) {
		// Ownership moved to a canonical write group after ingest.
		w.recordMappingSkip(connector, mapping, write.ObservedAt, skipWriteGroupOwned)
		return false, nil
	}
	writeErr := w.writeMapping(ctx, connector, mapping, write.Value, write.ObservedAt)
	attemptedAt := w.now().UTC()
	outcomeCtx, cancel := writeOutcomeContext(ctx)
	defer cancel()
	if writeErr != nil {
		var failures []error
		failure := fmt.Errorf("mapping %s 寫入失敗: %w", mapping.ID, writeErr)
		failures = append(failures, &snapshotRetryError{cause: failure, retryAt: attemptedAt.Add(retryDelay), connectorRevision: connector.IdentityRevision})
		if outcomeErr := w.recordWriteOutcome(outcomeCtx, connector, attemptedAt, deliveryOutcomeFailed, retryFailureSummary(writeErr, retryDelay)); outcomeErr != nil {
			failures = append(failures, fmt.Errorf("mapping %s 更新寫入狀態失敗: %w", mapping.ID, outcomeErr))
		}
		recordWriteHistory(connector.ID, WriteHistoryRecord{
			ObservedAt:               write.ObservedAt,
			AttemptedAt:              &attemptedAt,
			Status:                   deliveryOutcomeFailed,
			RowCount:                 0,
			TableName:                mapping.TableName,
			GroupKey:                 cloneOptionalString(mapping.GroupKey),
			EffectiveIntervalSeconds: write.IntervalSeconds,
			ErrorSummary:             writeErr.Error(),
			RetryAfterSeconds:        int(retryDelay.Seconds()),
			NextRetryAt:              optionalRetryAt(attemptedAt, retryDelay),
		})
		w.restorePending(write)
		return false, errors.Join(failures...)
	}
	outcomeErr := w.recordWriteOutcome(outcomeCtx, connector, attemptedAt, deliveryOutcomeSuccess, "")
	recordWriteHistory(connector.ID, WriteHistoryRecord{
		ObservedAt:               write.ObservedAt.UTC(),
		AttemptedAt:              &attemptedAt,
		Status:                   deliveryOutcomeSuccess,
		RowCount:                 1,
		TableName:                mapping.TableName,
		EffectiveIntervalSeconds: write.IntervalSeconds,
	})
	if outcomeErr != nil {
		return true, fmt.Errorf("mapping %s 更新寫入狀態失敗: %w", mapping.ID, outcomeErr)
	}
	return true, nil
}

// A buffered value belongs to the mapping and connector revision observed at
// ingest; disabling or redirecting them must not write to the former target.
func (w *Writer) mappingIsCurrent(ctx context.Context, mapping *schema.DatabaseTargetMapping, connector *schema.DatabaseConnector) (bool, error) {
	version := w.refreshVersion.Load()
	entries, err := w.lookupMappings(mapping.TagID)
	if err != nil {
		return false, fmt.Errorf("確認資料庫目標映射失敗: %w", err)
	}
	for _, entry := range entries {
		if entry.Mapping.ID != mapping.ID {
			continue
		}
		if entry.Err != nil {
			return false, fmt.Errorf("確認資料庫連接器失敗: %w", entry.Err)
		}
		if !sameSnapshotBinding(mapping, connector, entry.Mapping, entry.Connector) {
			return false, nil
		}
		if w.refreshing.Load() || w.refreshActive.Load() > 0 || w.refreshVersion.Load() != version {
			return w.mappingIsLive(ctx, mapping, connector)
		}
		return true, nil
	}
	return false, nil
}
