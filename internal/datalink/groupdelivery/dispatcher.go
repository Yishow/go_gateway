package groupdelivery

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Defaults applied when a bound is unset or not positive.
const (
	DefaultBatchSize     = 50
	DefaultMaxPartitions = 4
)

// DispatcherConfig bounds one delivery cycle.
type DispatcherConfig struct {
	// BatchSize is the most items delivered per partition per cycle.
	BatchSize int
	// MaxPartitions is the most partitions delivered concurrently.
	MaxPartitions int
	Now           func() time.Time
}

// CycleReport counts what one cycle did.
type CycleReport struct {
	Partitions  int
	Delivered   int
	Retrying    int
	Blocked     int
	Quarantined int
	Unknown     int
}

// Dispatcher delivers outbox rows partition by partition, in order. A row that
// is not finished holds back only its own partition.
type Dispatcher struct {
	store  *Store
	sender *Sender
	config DispatcherConfig
}

// NewDispatcher builds a dispatcher with defaults for unset bounds.
func NewDispatcher(store *Store, sender *Sender, config DispatcherConfig) *Dispatcher {
	if config.BatchSize <= 0 {
		config.BatchSize = DefaultBatchSize
	}
	if config.MaxPartitions <= 0 {
		config.MaxPartitions = DefaultMaxPartitions
	}
	if config.Now == nil {
		config.Now = func() time.Time { return time.Now().UTC() }
	}
	return &Dispatcher{store: store, sender: sender, config: config}
}

// RunOnce performs one bounded delivery cycle: every partition with a due head
// is worked by at most MaxPartitions workers, each delivering up to BatchSize
// consecutive rows in order and stopping at the first row that does not commit.
func (d *Dispatcher) RunOnce(ctx context.Context) (CycleReport, error) {
	return d.runCycle(ctx, ctx)
}

// runCycle is RunOnce with a separate grace context: once grace ends no new row
// is started, while work (in-flight attempts) only stops when work ends.
func (d *Dispatcher) runCycle(ctx, grace context.Context) (CycleReport, error) {
	heads, err := d.store.ReadyHeads(ctx, d.config.Now())
	if err != nil {
		return CycleReport{}, err
	}
	report := CycleReport{Partitions: len(heads)}
	var (
		mu   sync.Mutex
		errs []error
		wg   sync.WaitGroup
	)
	slots := make(chan struct{}, d.config.MaxPartitions)
	for _, head := range heads {
		slots <- struct{}{}
		wg.Add(1)
		go func(head Head) {
			defer wg.Done()
			defer func() { <-slots }()
			partial, partErr := d.deliverPartition(ctx, grace, head)
			mu.Lock()
			defer mu.Unlock()
			report.Delivered += partial.Delivered
			report.Retrying += partial.Retrying
			report.Blocked += partial.Blocked
			report.Quarantined += partial.Quarantined
			report.Unknown += partial.Unknown
			if partErr != nil {
				errs = append(errs, partErr)
			}
		}(head)
	}
	wg.Wait()
	return report, errors.Join(errs...)
}

func (d *Dispatcher) deliverPartition(ctx, grace context.Context, head Head) (CycleReport, error) {
	var report CycleReport
	effectKey := head.EffectKey
	for i := 0; i < d.config.BatchSize; i++ {
		if graceEnded(grace) {
			return report, nil
		}
		result, err := d.sender.Deliver(ctx, effectKey)
		if err != nil {
			// Losing the claim to a concurrent cycle is normal, not a failure.
			if errors.Is(err, ErrNotDeliverable) || errors.Is(err, ErrFenced) {
				return report, nil
			}
			return report, err
		}
		switch result.State {
		case StateCommitted:
			if !result.Already {
				report.Delivered++
			}
		case StateRetrying:
			report.Retrying++
		case StateBlocked:
			report.Blocked++
		case StateQuarantined:
			report.Quarantined++
		case StateUnknown:
			report.Unknown++
		}
		if result.State != StateCommitted {
			return report, nil
		}
		next, ok, err := d.store.PartitionHead(ctx, head.PartitionKey, d.config.Now())
		if err != nil {
			return report, err
		}
		if !ok {
			return report, nil
		}
		effectKey = next
	}
	return report, nil
}

// graceEnded reports whether no new row may be started.
func graceEnded(grace context.Context) bool {
	select {
	case <-grace.Done():
		return true
	default:
		return false
	}
}
