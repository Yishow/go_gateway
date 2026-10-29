package groupdelivery

import (
	"context"
	"time"

	"go-gateway/internal/datalink/counterproof"
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

// PersistsCounterBaselines states that CommitClosure stores a closure's counter
// baselines in the closure transaction and Restore returns them, which a
// cumulative profile requires of any durable ledger.
func (l *Ledger) PersistsCounterBaselines() bool { return true }

// AppendSample journals a sample; nil means it is committed.
func (l *Ledger) AppendSample(ctx context.Context, bucketStart time.Time, sample snapshot.Sample) error {
	return l.store.AppendSample(ctx, l.key, bucketStart, sample)
}

// VerifyCounterDuplicate rechecks a previously accepted sample and its proof.
func (l *Ledger) VerifyCounterDuplicate(ctx context.Context, sample snapshot.Sample) error {
	return l.store.VerifyCounterDuplicate(ctx, l.key, sample)
}

// Restore returns the checkpoint and unconsumed samples of the bound group.
func (l *Ledger) Restore(ctx context.Context) (Restored, error) {
	return l.store.Restore(ctx, l.key)
}

// CommitClosure persists a closure for the bound group.
func (l *Ledger) CommitClosure(ctx context.Context, closure Closure) error {
	return l.store.CommitClosure(ctx, l.key, closure)
}

// PrepareCounterReplay audits legacy counter associations before rebuilding memory.
func (l *Ledger) PrepareCounterReplay(ctx context.Context, expected []counterproof.CounterBindingReference) error {
	return l.store.AssociateLegacyCounterProofs(ctx, l.key, expected)
}

// RecognizeTransferredSample reads durable transfer authority, never appending.
func (l *Ledger) RecognizeTransferredSample(ctx context.Context, incoming snapshot.Sample) (HandoffRecognition, error) {
	return l.store.RecognizeTransferredSample(ctx, l.key, incoming)
}
