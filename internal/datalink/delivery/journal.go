package delivery

import (
	"fmt"
	"sync"
	"time"
)

// JournalEntry 代表本機耐久日誌的一筆記錄。
type JournalEntry struct {
	Sequence   int64     `json:"sequence"`
	RecordID   string    `json:"record_id"`
	ObservedAt time.Time `json:"observed_at"`
	Payload    []byte    `json:"payload"`
}

// Journal 定義本機日誌的操作介面。
type Journal interface {
	Append(entry *JournalEntry) error
	ScanFrom(fromSequence int64, limit int) ([]*JournalEntry, error)
	TruncateBefore(sequence int64) error
	Close() error
}

// MemoryJournal 記憶體測試用 Journal 實作。
type MemoryJournal struct {
	mu      sync.RWMutex
	entries []*JournalEntry
}

// NewMemoryJournal 建立新的記憶體日誌。
func NewMemoryJournal() *MemoryJournal {
	return &MemoryJournal{
		entries: make([]*JournalEntry, 0),
	}
}

func (m *MemoryJournal) Append(entry *JournalEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if entry == nil {
		return fmt.Errorf("journal entry cannot be nil")
	}
	m.entries = append(m.entries, entry)
	return nil
}

func (m *MemoryJournal) ScanFrom(fromSequence int64, limit int) ([]*JournalEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*JournalEntry
	for _, e := range m.entries {
		if e.Sequence >= fromSequence {
			result = append(result, e)
			if limit > 0 && len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (m *MemoryJournal) TruncateBefore(sequence int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var remaining []*JournalEntry
	for _, e := range m.entries {
		if e.Sequence >= sequence {
			remaining = append(remaining, e)
		}
	}
	m.entries = remaining
	return nil
}

func (m *MemoryJournal) Close() error {
	return nil
}
