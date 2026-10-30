package snapshot

import (
	"fmt"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/schema"
)

// ReasonCounterUnavailable marks a window whose cumulative record could not be
// computed because the baseline or the window broke the calculation contract.
const ReasonCounterUnavailable = "counter-unavailable"

// countersInvalid reports why the counter bindings do not fit the configured
// mode and members, or the empty string.
func (c Config) countersInvalid() string {
	if !c.Mode.Cumulative() {
		if len(c.Counters) > 0 {
			return "counter bindings belong to the cumulative profile only"
		}
		return ""
	}
	if len(c.Counters) != len(c.Members) {
		return "every cumulative member needs a confirmed counter binding"
	}
	for _, member := range c.Members {
		binding, ok := c.Counters[member.Key]
		switch {
		case !ok || binding.Validate() != nil:
			return "every cumulative member needs a confirmed counter binding"
		case binding.MemberKey != member.Key || binding.WorkspaceID != c.WorkspaceID ||
			binding.GroupID != c.GroupID || binding.ProfileRevision != c.GroupRevision:
			return "counter binding does not belong to this group revision"
		}
	}
	return ""
}

// counterReadings returns every good, typed reading the bucket accepted for a
// member, before any decimation. It is the sequence a cumulative window is
// judged on, in no particular order.
func (s *Selector) counterReadings(member string) []aggregation.CounterReading {
	var readings []aggregation.CounterReading
	for _, sample := range s.seen {
		if sample.MemberKey != member || sample.Quality != schema.QualityGood || sample.Value.Type() == "" {
			continue
		}
		readings = append(readings, aggregation.CounterReading{
			SampleID: sample.SampleID, ObservedAt: sample.ObservedAt, ReceivedAt: sample.ReceivedAt, Value: sample.Value,
		})
	}
	return readings
}

// applyCounter closes a member's window against its baseline and attaches the
// typed cumulative record. The baseline it advances lives in the plan, so
// nothing changes until the caller commits the plan.
func (a *Assembler) applyCounter(plan *Plan, member Member, start, end time.Time, selector *Selector, result *MemberResult) {
	prev, ok := plan.counters[member.Key]
	if !ok {
		prev = a.counters[member.Key]
	}
	start, end, _ = a.windowBounds(start, end)
	readings := selector.counterReadings(member.Key)
	counter, next, err := aggregation.CloseCounterWindow(
		prev, a.config.Counters[member.Key], aggregation.CounterWindow{Start: start, End: end}, readings,
	)
	if err != nil {
		// A broken contract is never turned into a value or a delta.
		result.Status, result.Reason = MemberInvalid, ReasonCounterUnavailable
		return
	}
	if counter.Status == aggregation.CounterInvalid {
		result.Status, result.Reason = MemberInvalid, counter.Reason
	}
	if a.config.Mode.mixed() {
		next = a.enforceMixedCounterType(member, readings, &counter, next, prev, result)
	}
	result.Counter = &counter
	if next != prev {
		plan.counters[member.Key] = next
	}
}

// CounterCheckpoints returns the baselines this plan advances, by member key.
// A durable caller persists them in the same transaction as the closure and
// only then commits the plan.
func (p Plan) CounterCheckpoints() map[string]aggregation.CounterCheckpoint {
	return copyCheckpoints(p.counters)
}

// CounterCheckpoints returns the committed baselines by member key.
func (a *Assembler) CounterCheckpoints() map[string]aggregation.CounterCheckpoint {
	return copyCheckpoints(a.counters)
}

// RestoreCounterCheckpoints installs durable baselines before the first Plan.
// The workspace, group and profile revision a checkpoint names are the fixed
// scope of this assembler: a checkpoint of another scope is foreign state and is
// refused, never mistaken for a binding that changed. A baseline whose
// measurement revision (definition, source binding, epoch, unit, policy) no
// longer matches the member is legitimate: the calculation starts a new baseline
// from it instead of subtracting across bindings.
func (a *Assembler) RestoreCounterCheckpoints(checkpoints map[string]aggregation.CounterCheckpoint) error {
	if !a.config.Mode.Cumulative() && !a.config.Mode.mixed() {
		return fmt.Errorf("%w: counter baselines belong to the cumulative profile only", ErrInvalidConfig)
	}
	restored := make(map[string]*aggregation.CounterCheckpoint, len(checkpoints))
	for member, checkpoint := range checkpoints {
		if m, known := a.members[member]; !known || (a.config.Mode.mixed() && m.Calculation != CalculationDelta) {
			return fmt.Errorf("%w: baseline for unknown member", ErrInvalidConfig)
		}
		if a.config.Mode.mixed() && !meanValueMatchesSource(a.members[member].SourceValueType, checkpoint.Reading.Value) {
			return fmt.Errorf("%w: baseline does not match frozen source type", aggregation.ErrCounterCheckpointInvalid)
		}
		if _, err := checkpoint.Binding.Key(); err != nil {
			return fmt.Errorf("%w: %w", aggregation.ErrCounterCheckpointInvalid, err)
		}
		binding := checkpoint.Binding
		if binding.MemberKey != member || binding.WorkspaceID != a.config.WorkspaceID ||
			binding.GroupID != a.config.GroupID || binding.ProfileRevision != a.config.GroupRevision {
			return fmt.Errorf("%w: baseline belongs to another scope", ErrInvalidConfig)
		}
		owned := cloneCounterCheckpoint(checkpoint)
		restored[member] = &owned
	}
	a.counters = restored
	return nil
}

func copyCheckpoints(source map[string]*aggregation.CounterCheckpoint) map[string]aggregation.CounterCheckpoint {
	copied := make(map[string]aggregation.CounterCheckpoint, len(source))
	for member, checkpoint := range source {
		copied[member] = cloneCounterCheckpoint(*checkpoint)
	}
	return copied
}
