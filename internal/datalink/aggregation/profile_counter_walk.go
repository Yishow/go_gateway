package aggregation

import (
	"math/big"

	"go-gateway/internal/datalink/measurement"
)

// walked is what a pass over the baseline and one window's readings found.
type walked struct {
	log eventLog
	// anomaly is the reason of the first step that could not be proven; the
	// delta is unknown when it is set.
	anomaly string
	// wraps counts the single rollover wraps that were proven.
	wraps int
}

func (w *walked) flag(reason string) {
	if w.anomaly == "" {
		w.anomaly = reason
	}
}

// walkSequence examines every adjacent pair, baseline first, in time order. The
// endpoints alone would hide a reset in the middle of the window, so nothing is
// skipped: each unexplained step is an event, and the delta is only usable when
// none is left.
func walkSequence(binding CounterBinding, baseline CounterReading, sequence []CounterReading) walked {
	var w walked
	previous := baseline
	for _, reading := range sequence {
		switch order := reading.ObservedAt.Compare(previous.ObservedAt); {
		case order < 0:
			w.flag(ReasonNonMonotonicTime)
			w.log.add(EventNonMonotonicTime, ReasonNonMonotonicTime, previous, reading)
		case order == 0 && !reading.Value.Equal(previous.Value):
			w.flag(ReasonContradictory)
			w.log.add(EventContradictory, ReasonContradictory, previous, reading)
		default:
			w.step(binding, previous, reading)
		}
		previous = reading
	}
	return w
}

// step classifies the change between two readings at strictly increasing time
// (or the same value at one instant).
func (w *walked) step(binding CounterBinding, from, to CounterReading) {
	change, err := measurement.Subtract(to.Value, from.Value)
	if err != nil {
		w.flag(ReasonMixedBasis)
		w.log.add(EventMixedBasis, ReasonMixedBasis, from, to)
		return
	}
	if sign, _ := change.Sign(); sign >= 0 || binding.Kind == CounterKindSigned {
		return // a rise, or a drop of a counter confirmed to be signed
	}
	if reason := wrapEvidence(binding, from, to); reason != "" {
		w.flag(reason)
		w.log.add(EventDecrease, reason, from, to)
		return
	}
	w.wraps++
	w.log.add(EventRollover, "", from, to)
}

// wrapEvidence returns "" when the drop from -> to is proven to be exactly one
// wrap of the modulus, or the reason the evidence falls short. Proof needs the
// modulus, a maximum rate, readings inside [0, modulus) and a rate bound that
// cannot reach a second wrap over the real time between the readings. Every
// comparison is on exact rationals; a float source has none and proves nothing.
func wrapEvidence(binding CounterBinding, from, to CounterReading) string {
	if binding.RolloverModulus == nil {
		return ReasonDecreaseUnproven
	}
	if binding.MaxRatePerSecond == nil {
		return ReasonRolloverEvidence
	}
	modulus, modulusOK := binding.RolloverModulus.Rat()
	rate, rateOK := binding.MaxRatePerSecond.Rat()
	before, beforeOK := from.Value.Rat()
	after, afterOK := to.Value.Rat()
	if !modulusOK || !rateOK || !beforeOK || !afterOK {
		return ReasonRolloverEvidence
	}
	inRange := func(value *big.Rat) bool { return value.Sign() >= 0 && value.Cmp(modulus) < 0 }
	if !inRange(before) || !inRange(after) {
		return ReasonRolloverRange
	}
	seconds := big.NewRat(to.ObservedAt.Sub(from.ObservedAt).Nanoseconds(), 1_000_000_000)
	bound := new(big.Rat).Mul(rate, seconds)
	if bound.Cmp(modulus) >= 0 {
		return ReasonRolloverMultiWrap
	}
	wrapped := new(big.Rat).Sub(modulus, before)
	wrapped.Add(wrapped, after)
	if wrapped.Cmp(bound) > 0 {
		return ReasonRolloverRate
	}
	return ""
}

// disagreesAtLatest reports whether the readings at the newest observation
// instant do not all carry the same value.
func disagreesAtLatest(sequence []CounterReading) bool {
	latest := sequence[len(sequence)-1]
	for _, reading := range sequence {
		if reading.ObservedAt.Equal(latest.ObservedAt) && !reading.Value.Equal(latest.Value) {
			return true
		}
	}
	return false
}

// tailConflictEvents records the disagreement at the newest instant of a window
// that has no baseline to walk from.
func tailConflictEvents(sequence []CounterReading) (events []CounterEvent, omitted int) {
	var log eventLog
	latest := sequence[len(sequence)-1]
	for _, reading := range sequence {
		if reading.ObservedAt.Equal(latest.ObservedAt) && !reading.Value.Equal(latest.Value) {
			log.add(EventContradictory, ReasonContradictory, reading, latest)
		}
	}
	return log.events, log.omitted
}
