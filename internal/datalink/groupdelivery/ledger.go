package groupdelivery

import (
	"context"
	"time"

	"go-gateway/internal/datalink/snapshot"
)

// Ledger is a Store bound to one applied group revision, which is the shape a
// group boundary consumes.
type Ledger struct {
	store *Store
	key   GroupKey
}

// Ledger binds the store to a group key.
func (s *Store) Ledger(key GroupKey) *Ledger { return &Ledger{store: s, key: key} }

// AppendSample journals a sample; nil means it is committed.
func (l *Ledger) AppendSample(ctx context.Context, bucketStart time.Time, sample snapshot.Sample) error {
	return l.store.AppendSample(ctx, l.key, bucketStart, sample)
}

// Restore returns the checkpoint and unconsumed samples of the bound group.
func (l *Ledger) Restore(ctx context.Context) (Restored, error) {
	return l.store.Restore(ctx, l.key)
}

// CommitClosure persists a closure for the bound group.
func (l *Ledger) CommitClosure(ctx context.Context, closure Closure) error {
	return l.store.CommitClosure(ctx, l.key, closure)
}
