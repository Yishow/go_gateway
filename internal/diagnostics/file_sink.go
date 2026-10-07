package diagnostics

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
)

// FileOptions identifies the rotation set and protected industrial data.
type FileOptions struct {
	Path, DatabasePath, DatabaseID string
	ProtectedPaths                 []string
	write                          func([]byte) error
}

// FileStatus separates a synchronized event from subsequent degradation. Path is
// native-only and is populated only by a successful current identity verification.
type FileStatus struct {
	Saved, Degraded                    bool
	Path                               string
	Dropped, DiskErrors, RecoveredGaps uint64
	LostFrom, LostTo                   string
}

// FileSink writes from its own bounded queue without blocking broker producers.
type FileSink struct {
	mu                  sync.Mutex
	owned               *ownedFiles
	queue               recordQueue
	status              FileStatus
	accepted, completed uint64
	lostLow, lostHigh   uint64
	closing             bool
	wake, done, changed chan struct{}
	write               func([]byte) error
}

// OpenFileSink validates and locks the complete set before starting the writer.
func OpenFileSink(o FileOptions) (*FileSink, error) {
	owned, err := openOwnedFiles(o)
	if err != nil {
		return nil, err
	}
	f := &FileSink{owned: owned, queue: newQueue(AdmissionRecords), wake: make(chan struct{}, 1), done: make(chan struct{}), changed: make(chan struct{}), write: owned.append}
	if o.write != nil {
		f.write = o.write
	}
	go f.run()
	return f, nil
}
func (f *FileSink) enqueue(r record) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closing || f.queue.count >= AdmissionRecords || f.queue.bytes+r.size > AdmissionBytes {
		f.lossLocked(r)
		f.status.Degraded = true
		return
	}
	f.queue.push(r)
	increment(&f.accepted)
	select {
	case f.wake <- struct{}{}:
	default:
	}
}
func (f *FileSink) run() {
	defer close(f.done)
	defer f.owned.close()
	for range f.wake {
		for {
			f.mu.Lock()
			if f.queue.count == 0 {
				closing := f.closing
				f.mu.Unlock()
				if closing {
					return
				}
				break
			}
			r := f.queue.pop()
			f.mu.Unlock()
			data, err := json.Marshal(r.event)
			if err == nil {
				err = f.write(append(data, '\n'))
			}
			f.mu.Lock()
			increment(&f.completed)
			if err != nil {
				increment(&f.status.DiskErrors)
				f.lossLocked(r)
				f.status.Degraded = true
			} else {
				f.status.Saved = true
				if f.status.Degraded {
					increment(&f.status.RecoveredGaps)
				}
				f.status.Degraded = false
			}
			close(f.changed)
			f.changed = make(chan struct{})
			f.mu.Unlock()
		}
	}
}
func (f *FileSink) lossLocked(r record) {
	increment(&f.status.Dropped)
	if f.status.LostFrom == "" || r.sequence < f.lostLow {
		f.lostLow = r.sequence
		f.status.LostFrom = r.event.Cursor()
	}
	if f.status.LostTo == "" || r.sequence > f.lostHigh {
		f.lostHigh = r.sequence
		f.status.LostTo = r.event.Cursor()
	}
}

func (f *FileSink) health() FileStatus { f.mu.Lock(); defer f.mu.Unlock(); return f.status }

// Health returns cached status without filesystem IO; it never provides a path.
func (f *FileSink) Health() FileStatus { return f.health() }

// Status performs read-only verification before reporting a saved native path.
// Callers needing a responsive disk-fault dialog must bound their wait externally.
func (f *FileSink) Status() FileStatus {
	s := f.health()
	if s.Saved {
		s.Path = f.owned.verifiedPath()
		if s.Path == "" {
			s.Saved = false
			s.Degraded = true
		}
	}
	return s
}

// Flush waits for accepted file records, honoring the caller's wait deadline.
func (f *FileSink) Flush(ctx context.Context) error {
	f.mu.Lock()
	target := f.accepted
	for f.completed < target {
		changed := f.changed
		f.mu.Unlock()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-changed:
		}
		f.mu.Lock()
	}
	degraded := f.status.Degraded
	f.mu.Unlock()
	if degraded {
		return errors.New("managed diagnostic storage is degraded")
	}
	return nil
}

// Close keeps ownership while draining. A timed-out OS disk stall remains pending
// rather than falsely reporting a flush or releasing handles beneath the writer.
func (f *FileSink) Close(ctx context.Context) error {
	f.mu.Lock()
	f.closing = true
	select {
	case f.wake <- struct{}{}:
	default:
	}
	f.mu.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-f.done:
	}
	if f.health().Degraded {
		return errors.New("managed diagnostic storage is degraded")
	}
	return nil
}
