package dbtarget

import (
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

func recordWriteHistory(connectorID string, status string, rowCount int, errorSummary string) {
	key := strings.TrimSpace(connectorID)
	if key == "" {
		return
	}
	record := WriteHistoryRecord{
		Timestamp:    time.Now().UTC(),
		Status:       strings.TrimSpace(status),
		RowCount:     rowCount,
		ErrorSummary: strings.TrimSpace(errorSummary),
	}

	globalWriteHistoryStore.mu.Lock()
	defer globalWriteHistoryStore.mu.Unlock()

	records := append(globalWriteHistoryStore.byConnector[key], record)
	if len(records) > 1000 {
		records = append([]WriteHistoryRecord(nil), records[len(records)-1000:]...)
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
	if limit > 200 {
		limit = 200
	}

	globalWriteHistoryStore.mu.RLock()
	defer globalWriteHistoryStore.mu.RUnlock()

	records := globalWriteHistoryStore.byConnector[key]
	if len(records) == 0 {
		return []WriteHistoryRecord{}
	}

	result := make([]WriteHistoryRecord, 0, limit)
	for i := len(records) - 1; i >= 0 && len(result) < limit; i-- {
		result = append(result, records[i])
	}
	return result
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
