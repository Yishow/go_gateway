package dbtarget

import (
	"slices"
	"strings"
	"sync"
	"time"
)

type writeHistoryStore struct {
	mu          sync.RWMutex
	byConnector map[string][]WriteHistoryRecord
}

var globalWriteHistoryStore = &writeHistoryStore{
	byConnector: map[string][]WriteHistoryRecord{},
}

func recordWriteHistory(connectorID string, record WriteHistoryRecord) {
	key := strings.TrimSpace(connectorID)
	if key == "" {
		return
	}
	record.Status = strings.TrimSpace(record.Status)
	record.ErrorSummary = strings.TrimSpace(record.ErrorSummary)
	record.TableName = strings.TrimSpace(record.TableName)
	if record.ObservedAt.IsZero() {
		record.ObservedAt = time.Now().UTC()
	} else {
		record.ObservedAt = record.ObservedAt.UTC()
	}
	if record.GroupKey != nil {
		trimmed := strings.TrimSpace(*record.GroupKey)
		record.GroupKey = &trimmed
	}

	globalWriteHistoryStore.mu.Lock()
	defer globalWriteHistoryStore.mu.Unlock()

	records := append(globalWriteHistoryStore.byConnector[key], record)
	if len(records) > 1000 {
		records = slices.Clone(records[len(records)-1000:])
	}
	globalWriteHistoryStore.byConnector[key] = records
}

func listWriteHistory(connectorID string, limit int) []WriteHistoryRecord {
	key := strings.TrimSpace(connectorID)
	if key == "" {
		return []WriteHistoryRecord{}
	}
	if limit <= 0 {
		limit = 20
	}
	limit = min(limit, 200)

	globalWriteHistoryStore.mu.RLock()
	defer globalWriteHistoryStore.mu.RUnlock()

	records := globalWriteHistoryStore.byConnector[key]
	if len(records) == 0 {
		return []WriteHistoryRecord{}
	}
	copied := slices.Clone(records)
	slices.SortFunc(copied, func(a, b WriteHistoryRecord) int {
		return b.ObservedAt.Compare(a.ObservedAt)
	})

	take := min(limit, len(copied))
	return slices.Clone(copied[:take])
}

func resetWriteHistory(connectorID string) {
	key := strings.TrimSpace(connectorID)
	if key == "" {
		return
	}
	globalWriteHistoryStore.mu.Lock()
	defer globalWriteHistoryStore.mu.Unlock()
	delete(globalWriteHistoryStore.byConnector, key)
}
