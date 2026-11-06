package dbtarget

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

// snapshotRetryError identifies a destination failure, not a failure to
// persist the outcome after a successful write.
type snapshotRetryError struct {
	cause   error
	retryAt time.Time
}

func (e *snapshotRetryError) Error() string { return e.cause.Error() }
func (e *snapshotRetryError) Unwrap() error { return e.cause }

func retryDelay(failures int) time.Duration {
	if failures >= 6 {
		return 60 * time.Second
	}
	return time.Duration(1<<failures) * time.Second
}

func optionalRetryAt(at time.Time, delay time.Duration) *time.Time {
	if delay <= 0 {
		return nil
	}
	next := at.Add(delay)
	return &next
}

func retryFailureSummary(err error, delay time.Duration) string {
	summary := "資料庫寫入失敗，請檢查目的地狀態"
	var networkError net.Error
	var sqlState interface{ SQLState() string }
	switch {
	case isSQLServerPermissionDeniedNumber(err) || isPermissionDeniedError(err):
		summary = "資料庫寫入失敗：權限不足，請檢查資料庫帳號權限"
	case IsSchemaShapeError(err):
		summary = "資料庫寫入失敗：目標資料表結構不符，請檢查資料表與欄位設定"
	case errors.Is(err, context.DeadlineExceeded), errors.As(err, &networkError),
		errors.As(err, &sqlState) && strings.HasPrefix(sqlState.SQLState(), "08"):
		summary = "資料庫寫入失敗：連線異常，請檢查連線設定與服務狀態"
	}
	if delay <= 0 {
		return summary
	}
	return fmt.Sprintf("%s，%d 秒後重試", summary, int(delay.Seconds()))
}

// restorePending preserves a newer value if ingest updated it during the
// failed attempt. The retry bit allows the snapshot to be retried before the
// next ordinary write interval.
func (w *Writer) restorePending(write *pendingWrite) {
	w.mu.Lock()
	defer w.mu.Unlock()
	byMapping := w.pending[write.Mapping.ConnectorID]
	if byMapping == nil {
		byMapping = make(map[string]*pendingWrite)
		w.pending[write.Mapping.ConnectorID] = byMapping
	}
	if latest := byMapping[write.Mapping.ID]; latest != nil {
		if sameSnapshotBinding(write.Mapping, write.Connector, latest.Mapping, latest.Connector) {
			latest.Retry = true
		} else {
			w.recordMappingSkip(write.Connector, write.Mapping, write.ObservedAt, skipMappingNotCurrent)
		}
		return
	}
	write.Retry = true
	byMapping[write.Mapping.ID] = write
}
