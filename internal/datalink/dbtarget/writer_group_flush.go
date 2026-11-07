package dbtarget

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// flushBuckets writes the due grouped buckets of one connector; an empty
// connectorID selects every connector.
func (w *Writer) flushBuckets(ctx context.Context, connectorID string, now time.Time, flushAll bool, retryDelay time.Duration) error {
	var failures []error
	for bucket := w.takeFlushableBucket(connectorID, now.UTC(), flushAll); bucket != nil; bucket = w.takeFlushableBucket(connectorID, now.UTC(), flushAll) {
		writeCtx, cancelWrite := context.WithTimeout(ctx, 5*time.Second)
		for column, mapping := range bucket.Mappings {
			current, err := w.mappingIsCurrent(writeCtx, mapping, bucket.Connector)
			switch {
			case err != nil:
				if writeCtx.Err() != nil && w.baseCtx.Err() != nil {
					w.restoreGrouped(bucket)
					cancelWrite()
					return errors.Join(append(failures, err)...)
				}
				failures = append(failures, fmt.Errorf("mapping %s 確認失敗: %w", mapping.ID, err))
				w.recordMappingSkip(bucket.Connector, mapping, bucket.Key.BucketStart, skipMappingCheckError)
				delete(bucket.Values, column)
			case !current:
				w.recordMappingSkip(bucket.Connector, mapping, bucket.Key.BucketStart, skipMappingNotCurrent)
				delete(bucket.Values, column)
			case w.suppress != nil && w.suppress(bucket.Key.ConnectorID, mapping.TagID):
				w.recordMappingSkip(bucket.Connector, mapping, bucket.Key.BucketStart, skipWriteGroupOwned)
				delete(bucket.Values, column)
			}
		}
		if len(bucket.Values) == 0 {
			cancelWrite()
			continue
		}
		writeErr := w.flushGroupedBucket(writeCtx, bucket)
		cancelWrite()
		attemptedAt := w.now().UTC()
		if writeErr != nil {
			outcomeCtx, cancel := writeOutcomeContext(ctx)
			summary := retryFailureSummary(writeErr, retryDelay)
			writeOutcomeErr := w.recordWriteOutcome(outcomeCtx, bucket.Connector, attemptedAt, deliveryOutcomeFailed, summary)
			outcomeErr := w.recordFlushOutcome(outcomeCtx, bucket.Connector, bucket.Key.BucketStart, deliveryOutcomeFailed, summary)
			cancel()
			if writeOutcomeErr != nil {
				failures = append(failures, fmt.Errorf("更新寫入狀態失敗: %w", writeOutcomeErr))
			}
			if outcomeErr != nil {
				failures = append(failures, fmt.Errorf("更新分組寫入狀態失敗: %w", outcomeErr))
			}
			recordWriteHistory(bucket.Key.ConnectorID, WriteHistoryRecord{
				ObservedAt:               bucket.Key.BucketStart,
				AttemptedAt:              &attemptedAt,
				Status:                   deliveryOutcomeFailed,
				RowCount:                 0,
				GroupKey:                 copyStringPointer(bucket.Key.GroupKey),
				TableName:                bucket.Key.TableName,
				EffectiveIntervalSeconds: bucket.Key.EffectiveIntervalSeconds,
				ErrorSummary:             writeErr.Error(),
				RetryAfterSeconds:        int(retryDelay.Seconds()),
				NextRetryAt:              optionalRetryAt(attemptedAt, retryDelay),
			})
			w.restoreGrouped(bucket)
			failures = append(failures, &snapshotRetryError{cause: fmt.Errorf("分組寫入失敗: %w", writeErr), retryAt: attemptedAt.Add(retryDelay)})
			return errors.Join(failures...)
		}
		if !flushAll {
			w.discardSupersededGrouped(bucket)
		}
		outcomeCtx, cancel := writeOutcomeContext(ctx)
		writeOutcomeErr := w.recordWriteOutcome(outcomeCtx, bucket.Connector, attemptedAt, deliveryOutcomeSuccess, "")
		outcomeErr := w.recordFlushOutcome(outcomeCtx, bucket.Connector, bucket.Key.BucketStart, deliveryOutcomeSuccess, "")
		cancel()
		recordWriteHistory(bucket.Key.ConnectorID, WriteHistoryRecord{
			ObservedAt:               bucket.Key.BucketStart,
			AttemptedAt:              &attemptedAt,
			Status:                   "success",
			RowCount:                 1,
			GroupKey:                 copyStringPointer(bucket.Key.GroupKey),
			TableName:                bucket.Key.TableName,
			EffectiveIntervalSeconds: bucket.Key.EffectiveIntervalSeconds,
		})
		if writeOutcomeErr != nil {
			failures = append(failures, fmt.Errorf("更新寫入狀態失敗: %w", writeOutcomeErr))
		}
		if outcomeErr != nil {
			failures = append(failures, fmt.Errorf("更新分組寫入狀態失敗: %w", outcomeErr))
		}
	}
	return errors.Join(failures...)
}

func (w *Writer) takeFlushableBucket(connectorID string, now time.Time, flushAll bool) *groupedWriteBucket {
	w.mu.Lock()
	defer w.mu.Unlock()

	var selected *groupedWriteBucket
	for key, bucket := range w.grouped {
		if connectorID != "" && key.ConnectorID != connectorID {
			continue
		}
		if !flushAll && !bucket.Retry && now.Before(bucket.DueAt) {
			continue
		}
		if selected != nil {
			if selected.Retry && !bucket.Retry {
				continue
			}
			if selected.Retry == bucket.Retry && !key.BucketStart.Before(selected.Key.BucketStart) {
				continue
			}
		}
		selected = bucket
	}
	if selected == nil {
		return nil
	}
	delete(w.grouped, selected.Key)
	interval := time.Duration(selected.Key.EffectiveIntervalSeconds) * time.Second
	intervalEnd := selected.Key.BucketStart.Add(interval)
	for column, mapping := range selected.Mappings {
		if latest, ok := w.groupedLocations[mapping.ID]; ok && latest.key == selected.Key {
			delete(w.groupedLocations, mapping.ID)
		}
		last, flushed := w.lastFlush[mapping.ID]
		if !selected.Retry && flushed && !intervalEnd.After(last) {
			w.recordMappingSkip(selected.Connector, mapping, selected.Key.BucketStart, skipIntervalFlushed)
			delete(selected.Values, column)
			delete(selected.Mappings, column)
			continue
		}
		if !flushed || intervalEnd.After(last) {
			w.lastFlush[mapping.ID] = intervalEnd
		}
	}
	return selected
}
