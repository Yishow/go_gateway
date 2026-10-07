package desktop

import "sync/atomic"

// browserTask permits one best-effort shell operation without blocking the UI.
// It is deliberately not part of database/service shutdown completion.
type browserTask struct{ running atomic.Bool }

func (task *browserTask) start(open func(string) error, url string, onError func()) bool {
	if !task.running.CompareAndSwap(false, true) {
		return false
	}
	go func() {
		defer func() {
			if recover() != nil && onError != nil {
				onError()
			}
			task.running.Store(false)
		}()
		if err := open(url); err != nil && onError != nil {
			onError()
		}
	}()
	return true
}
