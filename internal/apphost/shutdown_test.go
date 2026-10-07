package apphost

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSharedDeadlineRetainsStorageUntilWorkerActuallyEnds(t *testing.T) {
	release := make(chan struct{})
	states := make(chan ShutdownState, 16)
	var closed atomic.Bool
	stop := NewShutdown(15*time.Millisecond, []Phase{
		{Name: "runtime", Stop: func(ctx context.Context) error { <-ctx.Done(); return ctx.Err() }, Wait: func(context.Context) error { <-release; return nil }},
		{Name: "database", Stop: func(context.Context) error { closed.Store(true); return nil }},
	}, func(state ShutdownState) { states <- state })
	stop.Begin()
	for {
		select {
		case state := <-states:
			if state.Process == "stop-timeout" {
				if closed.Load() {
					t.Fatal("DB closed beneath worker")
				}
				close(release)
				if err := stop.Wait(t.Context()); !errors.Is(err, ErrShutdownTimeout) {
					t.Fatal(err)
				}
				if !closed.Load() {
					t.Fatal("original cleanup not completed")
				}
				return
			}
		case <-time.After(time.Second):
			close(release)
			t.Fatal("no timeout notification")
		}
	}
}

func TestShutdownIsOnceAndUsesOneDeadline(t *testing.T) {
	var calls atomic.Int32
	var deadlines []time.Time
	phases := make([]Phase, 3)
	for i := range phases {
		phases[i] = Phase{Name: "runtime", Stop: func(ctx context.Context) error {
			calls.Add(1)
			d, ok := ctx.Deadline()
			if !ok {
				t.Error("deadline missing")
			}
			deadlines = append(deadlines, d)
			return nil
		}}
	}
	stop := NewShutdown(time.Second, phases, nil)
	var wg sync.WaitGroup
	for range 32 {
		wg.Go(stop.Begin)
	}
	wg.Wait()
	if err := stop.Wait(t.Context()); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 3 {
		t.Fatal("phase ran twice")
	}
	for _, d := range deadlines {
		if !d.Equal(deadlines[0]) {
			t.Fatal("phase reset deadline")
		}
	}
}

func TestTimeoutCannotBeOverwrittenByDelayedNotification(t *testing.T) {
	entered, allowCallback, allowWorker := make(chan struct{}), make(chan struct{}), make(chan struct{})
	states := make(chan string, 16)
	stop := NewShutdown(10*time.Millisecond, []Phase{{Name: "runtime", Stop: func(context.Context) error { <-allowWorker; return nil }}}, func(state ShutdownState) {
		if state.Process == "stopping" {
			close(entered)
			<-allowCallback
		}
		states <- state.Process
	})
	stop.Begin()
	<-entered
	time.Sleep(25 * time.Millisecond)
	close(allowCallback)
	for {
		select {
		case state := <-states:
			if state == "stop-timeout" {
				goto timedOut
			}
		case <-time.After(time.Second):
			t.Fatal("no timeout")
		}
	}
timedOut:
	select {
	case state := <-states:
		if state == "stopping" {
			t.Fatal("late stale callback")
		}
	default:
	}
	close(allowWorker)
	_ = stop.Wait(t.Context())
}
