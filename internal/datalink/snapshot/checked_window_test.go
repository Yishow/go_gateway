package snapshot

import (
	"errors"
	"math"
	"testing"
	"time"
)

func TestCheckedCadenceDuration(t *testing.T) {
	for _, tc := range []struct {
		value int64
		unit  string
		want  time.Duration
		bad   bool
	}{
		{5, "second", 5 * time.Second, false}, {15, "second", 15 * time.Second, false},
		{60, "second", time.Minute, false}, {10, "minute", 10 * time.Minute, false},
		{9223372036, "second", 9223372036 * time.Second, false},
		{153722867, "minute", 153722867 * time.Minute, false},
		{9223372037, "second", 0, true}, {153722868, "minute", 0, true},
		{math.MaxInt64, "minute", 0, true}, {0, "second", 0, true}, {-1, "second", 0, true}, {1, "hour", 0, true},
	} {
		got, err := CheckedCadenceDuration(tc.value, tc.unit)
		if tc.bad {
			var typed *WindowArithmeticError
			if !errors.As(err, &typed) || got != 0 {
				t.Fatalf("%+v: %v %v", tc, got, err)
			}
		} else if err != nil || got != tc.want {
			t.Fatalf("%+v: %v %v", tc, got, err)
		}
	}
}

func TestCheckedWindowEpochAlignment(t *testing.T) {
	for _, seconds := range []int64{5, 15, 60, 600} {
		interval, _ := CheckedCadenceDuration(seconds, "second")
		for _, sec := range []int64{-601, -16, -1, 0, 1, 16, 601, 1793380001} {
			now := time.Unix(sec, 123).UTC()
			start, end, err := CheckedWindowBounds(now, interval)
			if err != nil || start.After(now) || !end.After(now) || end.Sub(start) != interval || start.Unix()%seconds != 0 {
				t.Fatalf("%d %d: %v %v %v", seconds, sec, start, end, err)
			}
			restartStart, restartEnd, err := CheckedWindowBounds(now.In(time.FixedZone("offset", 8*3600)), interval)
			if err != nil || !start.Equal(restartStart) || !end.Equal(restartEnd) {
				t.Fatal("restart/location phase changed")
			}
			next, _, err := CheckedWindowBounds(end, interval)
			if err != nil || !next.Equal(end) {
				t.Fatal("exact boundary did not advance")
			}
		}
	}
}

func TestCheckedWindowRejectsOverflow(t *testing.T) {
	for _, tc := range []struct {
		now      time.Time
		interval time.Duration
	}{
		{time.Unix(0, math.MinInt64), 15 * time.Second},
		{time.Unix(0, math.MaxInt64), 15 * time.Second},
		{time.Date(2300, 1, 1, 0, 0, 0, 0, time.UTC), time.Second},
		{time.Unix(0, 0), 0}, {time.Unix(0, 0), time.Nanosecond},
	} {
		start, end, err := CheckedWindowBounds(tc.now, tc.interval)
		var typed *WindowArithmeticError
		if !errors.As(err, &typed) || !start.IsZero() || !end.IsZero() {
			t.Errorf("%+v: %v %v %v", tc, start, end, err)
		}
	}
}

func TestCheckedWindowDeadlineAndSkew(t *testing.T) {
	for _, op := range []string{"lateness", "future_skew"} {
		got, err := CheckedWindowAdd(time.Unix(0, math.MaxInt64-1), time.Nanosecond, op)
		if err != nil || got.UnixNano() != math.MaxInt64 {
			t.Fatalf("edge %v %v", got, err)
		}
		for _, tc := range []struct {
			now   time.Time
			delay time.Duration
		}{
			{time.Unix(0, math.MaxInt64), time.Nanosecond},
			{time.Unix(0, 0), -time.Nanosecond},
			{time.Date(2300, 1, 1, 0, 0, 0, 0, time.UTC), 0},
		} {
			got, err := CheckedWindowAdd(tc.now, tc.delay, op)
			var typed *WindowArithmeticError
			if !errors.As(err, &typed) || typed.Operation != op || !got.IsZero() {
				t.Errorf("%s %+v %v %v", op, tc, got, err)
			}
		}
	}
}

func TestCheckedWindowRepresentableEdges(t *testing.T) {
	// The first whole-second boundary above MinInt64 is valid, while the
	// preceding fractional timestamp cannot have a representable aligned start.
	for _, nanos := range []int64{-9223372036000000000, -1000000000, 0, 9223372035000000000} {
		start, end, err := CheckedWindowBounds(time.Unix(0, nanos), time.Second)
		if err != nil || start.UnixNano() != nanos || end.UnixNano() != nanos+int64(time.Second) {
			t.Fatalf("boundary %d: %v %v %v", nanos, start, end, err)
		}
	}
	// Large independent offsets may be legal from a negative timestamp; do not
	// reject solely because summing the durations first would overflow.
	current := time.Unix(0, math.MinInt64)
	got, err := CheckedWindowAdd(current, time.Duration(math.MaxInt64), "end")
	if err != nil || got.UnixNano() != -1 {
		t.Fatalf("large first offset: %v %v", got, err)
	}
	got, err = CheckedWindowAdd(got, time.Duration(math.MaxInt64), "lateness")
	if err != nil || got.UnixNano() != math.MaxInt64-1 {
		t.Fatalf("large second offset: %v %v", got, err)
	}
}
