package audit

import (
	"context"
	"sync"
)

// MemoryRepository stores audit history in memory for tests and non-SQL wiring.
type MemoryRepository struct {
	mu      sync.RWMutex
	entries []Entry
}

// NewMemoryRepository creates an in-memory audit repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

// Create stores an audit history entry.
func (r *MemoryRepository) Create(_ context.Context, entry *Entry) error {
	if entry == nil {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, *entry)
	return nil
}

// List returns recent audit history entries.
func (r *MemoryRepository) List(_ context.Context, filter ListFilter) ([]Entry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	limit := normalizeLimit(filter.Limit)
	results := make([]Entry, 0, limit)
	for idx := len(r.entries) - 1; idx >= 0 && len(results) < limit; idx-- {
		entry := r.entries[idx]
		if entry.WorkspaceID != filter.WorkspaceID {
			continue
		}
		if filter.EventType != "" && entry.EventType != filter.EventType {
			continue
		}
		results = append(results, entry)
	}
	return results, nil
}
