package diagnostics

import (
	"context"
	"math"
	"strconv"
	"time"

	"go-gateway/internal/datalink/runtimeview"
)

// EmitDomain validates an untrusted reference before bounded coalescing. Its
// result is local capture acceptance, not admission or a domain outcome. Level
// filtering is intentional and consumes neither new state nor loss budget. A
// filtered confirmed transition can still finalize an observed exact warn fault.
func (b *Broker) EmitDomain(ctx context.Context, data []byte, scope runtimeview.Scope) bool {
	if b == nil || b.domain == nil {
		return false
	}
	visible := levelRank(ProjectDomainEvent(data, scope).Level) >= b.level
	if v, err := ParseDomainEventRef(data, scope); err == nil {
		if intakeObservation(v) {
			if !visible {
				return false
			}
			return b.domain.observeIntake(ctx, data, scope)
		}
		if fixedSourceObservation(v) {
			if !visible && v.Code != domainSourceApplied {
				return false
			}
			return b.domain.observeSource(ctx, data, scope, visible)
		}
		if fixedCollectorObservation(v) {
			if !visible && v.Code != domainCollectorRecovered {
				return false
			}
			return b.domain.observeCollector(ctx, data, scope, visible)
		}
		if v.Code == domainDeliveryFirstConfirmed {
			return b.domain.observeDeliveryCommit(ctx, data, scope, visible)
		}
		if fixedMappingObservation(v) {
			if !visible && v.Code != domainMappingRecovered {
				return false
			}
			return b.domain.observeMapping(ctx, data, scope, visible)
		}
	}
	if !visible {
		return false
	}
	return b.domain.Observe(ctx, data, scope)
}

// DomainInfoEnabled reports immutable informational capture eligibility, not health
// or capacity. Producers use it before allocating filtered intake descriptors.
func (b *Broker) DomainInfoEnabled() bool {
	return b != nil && b.domain != nil && levelRank(levelInfo) >= b.level
}

// DomainWarningsEnabled reports immutable warning-level capture eligibility.
// It is not an admission-capacity or runtime-health check, and never records loss.
func (b *Broker) DomainWarningsEnabled() bool {
	return b != nil && b.domain != nil && levelRank(levelWarn) >= b.level
}

// EmitDomainUnavailable captures a fixed, identity-free explanation. There are no
// candidate fields: missing producer proof cannot authorize partial attribution.
// Repeated unknown proof shares one bounded identity-free heartbeat summary.
func (b *Broker) EmitDomainUnavailable(ctx context.Context) bool {
	if b == nil || b.domain == nil {
		return false
	}
	e := project(Input{Code: domainUnavailableCode}, nil)
	if levelRank(e.Level) < b.level {
		return false
	}
	return b.domain.observeUnavailable(ctx)
}

// Only the validated projector/coalescer and fixed unavailable branch call this
// private admission path. It does not re-project safe fields through legacy Input,
// wait on broker contention, or count a rejection: OnDrop owns that accounting.
func (b *Broker) tryAdmitDomain(e Event) bool {
	if levelRank(e.Level) < b.level {
		return true
	}
	e.InstanceID = b.instance
	e.Timestamp = time.Now().UTC()
	e, size := boundEvent(e)
	if !b.mu.TryLock() {
		return false
	}
	defer b.mu.Unlock()
	b.foldDomainLossLocked()
	if b.closing || b.next == math.MaxUint64 || b.queue.count >= AdmissionRecords || b.queue.bytes+size > AdmissionBytes {
		return false
	}
	b.next++
	e.Sequence = strconv.FormatUint(b.next, 10)
	b.queue.push(record{event: e, sequence: b.next, size: size})
	b.signal()
	return true
}

// RecordDomainCaptureLoss accounts already-known bounded observation overflow.
// It accepts no candidate fields, emits no events and delegates to the single
// existing capture-loss owner. Callers must not also report the same loss.
func (b *Broker) RecordDomainCaptureLoss(n uint64) {
	if b != nil {
		b.recordDomainLoss(n)
	}
}

// The scalar accumulator is bounded and lock-free with respect to broker state.
// It cannot block a producer behind a reader, subscriber or file operation.
func (b *Broker) recordDomainLoss(n uint64) {
	for n != 0 {
		old := b.domainLoss.Load()
		next := old + n
		if math.MaxUint64-old < n {
			next = math.MaxUint64
		}
		if b.domainLoss.CompareAndSwap(old, next) {
			break
		}
	}
	b.signal()
}

// Loss joins the existing capture gap/cursor accounting under the broker lock.
// Rejected capture after Close cannot advance a drained cursor, but remains a
// counted loss at that boundary. No per-event auxiliary queue is introduced.
func (b *Broker) foldDomainLossLocked() {
	n := b.domainLoss.Swap(0)
	if n == 0 {
		return
	}
	first := b.next
	if !b.closing && b.next < math.MaxUint64 {
		first++
		b.next += min(n, math.MaxUint64-b.next)
	}
	if b.dropped == 0 {
		b.dropFirst = first
	}
	b.dropLast = b.next
	b.dropped += min(n, math.MaxUint64-b.dropped)
}
