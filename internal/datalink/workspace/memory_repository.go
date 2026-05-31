package workspace

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	mu     sync.RWMutex
	record *Record
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (r *MemoryRepository) Get(context.Context) (*Record, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.record == nil {
		return nil, ErrNotFound
	}

	return cloneRecord(r.record), nil
}

func (r *MemoryRepository) Save(_ context.Context, record *Record) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.record = cloneRecord(record)
	return nil
}
