package aggregation

import (
	"cmp"
	"fmt"
	"slices"

	"go-gateway/internal/datalink/measurement"
)

// CloseCounterWindow closes one window of one counter and returns its typed
// record together with the checkpoint to carry forward. readings are every
// eligible reading of the window before any decimation, in any order.
//
// The checkpoint is a value: prev is never modified, and the returned pointer is
// prev itself when the window does not advance the baseline, so a caller can
// tell whether anything changed and persist the new checkpoint together with the
// record. Only a window that follows the one the checkpoint was taken at can be
// folded; replaying a folded window is ErrCounterWindowOrder, never a second
// delta.
//
// A difference is produced only when it is proven: the same binding, strictly
// ordered time, and every adjacent pair of the baseline and the whole window
// either increasing or explained (a confirmed signed counter's decrease, or one
// wrap with rollover evidence). Anything else keeps the raw reading, reports a
// NULL delta with its reason and the events found, and moves the baseline to the
// window's last reliable reading without counting the uncertain change. A delta
// that crosses missing windows is a span total with its real interval.
func CloseCounterWindow(prev *CounterCheckpoint, binding CounterBinding, window CounterWindow, readings []CounterReading) (CounterWindowResult, *CounterCheckpoint, error) {
	key, err := binding.Key()
	if err != nil {
		return CounterWindowResult{}, prev, err
	}
	if window.Start.IsZero() || !window.End.After(window.Start) {
		return CounterWindowResult{}, prev, fmt.Errorf("%w: window must have an end after its start", ErrCounterWindowInvalid)
	}
	sequence, err := orderedReadings(window, readings)
	if err != nil {
		return CounterWindowResult{}, prev, err
	}
	var prevKey string
	if prev != nil {
		if prevKey, err = prev.Binding.Key(); err != nil {
			return CounterWindowResult{}, prev, fmt.Errorf("%w: %w", ErrCounterCheckpointInvalid, err)
		}
		if prev.WindowEnd.After(window.Start) {
			return CounterWindowResult{}, prev, fmt.Errorf("%w: baseline already covers up to %s", ErrCounterWindowOrder, prev.WindowEnd.UTC().Format("2006-01-02T15:04:05Z"))
		}
	}
	if len(sequence) == 0 {
		return CounterWindowResult{Status: CounterMissing, Reason: ReasonNoEligibleReading}, prev, nil
	}
	for _, reading := range sequence {
		if _, numeric := reading.Value.Sign(); !numeric {
			return CounterWindowResult{Status: CounterInvalid, Reason: ReasonNonNumericReading}, prev, nil
		}
	}

	latest := sequence[len(sequence)-1]
	raw := latest
	result := CounterWindowResult{Raw: &raw}
	advanced := &CounterCheckpoint{Binding: binding, WindowEnd: window.End.UTC(), Reading: latest}
	// A newest instant that disagrees with itself cannot be a baseline: whichever
	// value is chosen would silently shift every later difference.
	unreliableTail := disagreesAtLatest(sequence)

	if prev == nil || prevKey != key {
		result.Status, result.Reason = CounterBaseline, ReasonFirstReading
		if prev != nil {
			result.Reason = ReasonBindingChanged
		}
		if unreliableTail {
			result.Status, result.Reason = CounterUncertain, ReasonContradictory
			result.Events, result.EventsOmitted = tailConflictEvents(sequence)
			return result, prev, nil
		}
		return result, advanced, nil
	}

	walk := walkSequence(binding, prev.Reading, sequence)
	result.Events, result.EventsOmitted = walk.log.events, walk.log.omitted
	if walk.anomaly != "" {
		result.Status, result.Reason = CounterUncertain, walk.anomaly
		if unreliableTail {
			return result, prev, nil
		}
		return result, advanced, nil
	}
	diff, representable := exactDifference(latest.Value, prev.Reading.Value)
	for range walk.wraps {
		if representable {
			diff, representable = exactSum(diff, *binding.RolloverModulus)
		}
	}
	if !representable {
		result.Status, result.Reason = CounterUncertain, ReasonDeltaNotRepresented
		return result, advanced, nil
	}
	basis := DeltaBasisExact
	if latest.Value.Type() == measurement.ExactFloat64 {
		basis = DeltaBasisFloat64
	}
	result.Status = CounterOK
	result.Delta = &CounterDelta{
		Value: diff, Basis: basis, From: prev.Reading.ObservedAt.UTC(), To: latest.ObservedAt.UTC(),
		BaselineSampleID: prev.Reading.SampleID, EndSampleID: latest.SampleID,
	}
	if !prev.WindowEnd.Equal(window.Start) {
		result.Status, result.Delta.Span = CounterGap, true
	}
	return result, advanced, nil
}

// exactDifference is to - from, or false when the pair has no representable
// exact difference. That is a property of the data, not a failure of the caller.
func exactDifference(to, from measurement.ExactValue) (measurement.ExactValue, bool) {
	diff, err := measurement.Subtract(to, from)
	return diff, err == nil
}

func exactSum(a, b measurement.ExactValue) (measurement.ExactValue, bool) {
	sum, err := measurement.Add(a, b)
	return sum, err == nil
}

// orderedReadings validates and sorts a window's readings by observation time,
// with the sample id breaking a tie, on a copy.
func orderedReadings(window CounterWindow, readings []CounterReading) ([]CounterReading, error) {
	sequence := make([]CounterReading, 0, len(readings))
	for _, reading := range readings {
		observed := reading.ObservedAt.UTC()
		switch {
		case reading.SampleID == "" || observed.IsZero():
			return nil, fmt.Errorf("%w: a reading needs a sample id and an observation time", ErrCounterReadingInvalid)
		case observed.Before(window.Start) || !observed.Before(window.End):
			return nil, fmt.Errorf("%w: reading is outside the window", ErrCounterReadingInvalid)
		}
		reading.ObservedAt = observed
		if !reading.ReceivedAt.IsZero() {
			reading.ReceivedAt = reading.ReceivedAt.UTC()
		}
		sequence = append(sequence, reading)
	}
	slices.SortFunc(sequence, func(a, b CounterReading) int {
		return cmp.Or(a.ObservedAt.Compare(b.ObservedAt), cmp.Compare(a.SampleID, b.SampleID))
	})
	return sequence, nil
}
