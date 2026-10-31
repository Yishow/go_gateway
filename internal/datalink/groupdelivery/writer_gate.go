package groupdelivery

import (
	"context"
	"database/sql"
	"errors"
	"runtime"
	"sync"
	"weak"
)

// ErrNestedStoreTransaction rejects recursive writer acquisition, including
// cross-database acquisition, which could otherwise create a lock-order cycle.
var ErrNestedStoreTransaction = errors.New("nested store transaction is not supported")

type storeWriterContextKey struct{}
type storeWriterGate struct{ token chan struct{} }

// Only inTx writers participate. Direct Exec writers, other database handles,
// and external SQLite writers still obey SQLite's existing busy/error contract.
// Weak keys avoid retaining a database (or its Store clones) for process life.
var storeWriterGates sync.Map // weak.Pointer[sql.DB] -> *storeWriterGate

func writerGateFor(db *sql.DB) *storeWriterGate {
	key := weak.Make(db)
	gate := &storeWriterGate{token: make(chan struct{}, 1)}
	actual, loaded := storeWriterGates.LoadOrStore(key, gate)
	if !loaded {
		runtime.AddCleanup(db, func(key weak.Pointer[sql.DB]) { storeWriterGates.Delete(key) }, key)
	}
	runtime.KeepAlive(db)
	return actual.(*storeWriterGate)
}

func acquireStoreWriter(ctx context.Context, db *sql.DB) (context.Context, func(), error) {
	if ctx.Value(storeWriterContextKey{}) != nil {
		return nil, nil, ErrNestedStoreTransaction
	}
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	gate := writerGateFor(db)
	select {
	case gate.token <- struct{}{}:
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	}
	var once sync.Once
	release := func() { once.Do(func() { <-gate.token }) }
	// Cancellation can race the select when both cases are ready. Never begin
	// a transaction with a permit granted to an already canceled caller.
	if err := ctx.Err(); err != nil {
		release()
		return nil, nil, err
	}
	return context.WithValue(ctx, storeWriterContextKey{}, true), release, nil
}
