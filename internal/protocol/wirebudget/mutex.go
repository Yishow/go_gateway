package wirebudget

import (
	"context"
	"sync"
)

// Mutex is a zero-value-compatible operation gate. A canceled waiter does
// not need to wait for another request's network deadline or create a goroutine.
type Mutex struct {
	once  sync.Once
	token chan struct{}
}

func (m *Mutex) init() {
	m.once.Do(func() { m.token = make(chan struct{}, 1); m.token <- struct{}{} })
}
func (m *Mutex) Lock()   { _ = m.LockContext(context.Background()) }
func (m *Mutex) Unlock() { m.token <- struct{}{} }
func (m *Mutex) LockContext(ctx context.Context) error {
	m.init()
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-m.token:
	}
	if err := ctx.Err(); err != nil {
		m.Unlock()
		return err
	}
	return nil
}

// TryLock acquires only an immediately available gate, for nonblocking cleanup bookkeeping.
func (m *Mutex) TryLock() bool {
	m.init()
	select {
	case <-m.token:
		return true
	default:
		return false
	}
}
