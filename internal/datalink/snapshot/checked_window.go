package snapshot

import (
	"fmt"
	"math"
	"time"
)

// MaxCadenceSeconds is a representation bound, not an advertised product or
// runtime capability. It matches the existing lifecycle time.Duration bound.
const MaxCadenceSeconds int64 = math.MaxInt64 / int64(time.Second)

// MaxCadenceMinutes is the whole-minute equivalent of MaxCadenceSeconds.
const MaxCadenceMinutes int64 = MaxCadenceSeconds / 60

// WindowArithmeticError rejects an unrepresentable input or boundary. Callers
// must propagate it; the accompanying zero time is never a usable window.
type WindowArithmeticError struct {
	Operation string
	Reason    string
}

func (e *WindowArithmeticError) Error() string {
	return fmt.Sprintf("snapshot %s: %s", e.Operation, e.Reason)
}

func windowError(operation, reason string) error {
	return &WindowArithmeticError{Operation: operation, Reason: reason}
}

// CheckedCadenceDuration converts integer member intent without first doing a
// potentially overflowing multiplication or platform-dependent int conversion.
// It does not change the legacy profile interval limits.
func CheckedCadenceDuration(value int64, unit string) (time.Duration, error) {
	var scale time.Duration
	switch unit {
	case "second":
		scale = time.Second
	case "minute":
		scale = time.Minute
	default:
		return 0, windowError("cadence", "unsupported unit")
	}
	if value <= 0 {
		return 0, windowError("cadence", "interval must be positive")
	}
	if value > math.MaxInt64/int64(scale) {
		return 0, windowError("cadence", "interval exceeds time.Duration")
	}
	return time.Duration(value) * scale, nil
}

func checkedWindowNano(at time.Time, operation string) (int64, error) {
	n := at.UnixNano()
	// UnixNano is undefined outside its representable range; a round trip
	// prevents a far-future or far-past time being silently folded into it.
	if !time.Unix(0, n).Equal(at) {
		return 0, windowError(operation, "timestamp exceeds UnixNano range")
	}
	return n, nil
}

// CheckedWindowBounds returns the UTC epoch-aligned half-open window containing
// at. Both boundaries must be representable. Legacy BucketStart is unchanged.
func CheckedWindowBounds(at time.Time, interval time.Duration) (time.Time, time.Time, error) {
	if interval < time.Second || interval%time.Second != 0 {
		return time.Time{}, time.Time{}, windowError("window", "interval must be positive whole seconds")
	}
	n, err := checkedWindowNano(at, "window")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	width := int64(interval)
	quotient := n / width
	if n%width < 0 {
		quotient--
	}
	// Integer division truncates toward zero. For negative products this bound
	// precisely excludes a mathematical floor falling below MinInt64.
	if quotient < math.MinInt64/width || quotient > math.MaxInt64/width {
		return time.Time{}, time.Time{}, windowError("start", "aligned start exceeds UnixNano range")
	}
	start := time.Unix(0, quotient*width).UTC()
	end, err := CheckedWindowAdd(start, interval, "end")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return start, end, nil
}

// CheckedWindowAdd checks one nonnegative boundary offset independently. Use
// separate calls for interval, lateness and future skew, never add durations
// together before checking. This also validates the original timestamp at zero.
func CheckedWindowAdd(at time.Time, offset time.Duration, operation string) (time.Time, error) {
	n, err := checkedWindowNano(at, operation)
	if err != nil {
		return time.Time{}, err
	}
	if offset < 0 {
		return time.Time{}, windowError(operation, "offset must not be negative")
	}
	if n > math.MaxInt64-int64(offset) {
		return time.Time{}, windowError(operation, "boundary exceeds UnixNano range")
	}
	return time.Unix(0, n+int64(offset)).UTC(), nil
}
