package runtime

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

type provenanceCounter struct {
	Reason  string `json:"reason"`
	Counter struct {
		Kind             string `json:"kind"`
		Basis            string `json:"basis"`
		SpanTotal        bool   `json:"span_total"`
		RolloverModulus  string `json:"rollover_modulus"`
		MaxRatePerSecond string `json:"max_rate_per_second"`
		BaselineSampleID string `json:"baseline_sample_id"`
		EndSampleID      string `json:"end_sample_id"`
		EventsOmitted    int    `json:"events_omitted"`
		Events           []struct {
			Kind      string `json:"kind"`
			Reason    string `json:"reason"`
			FromValue string `json:"from_value"`
			ToValue   string `json:"to_value"`
			From      string `json:"from"`
			To        string `json:"to"`
		} `json:"events"`
	} `json:"counter"`
}

func readProvenance(t *testing.T, f *longFixture, bucket time.Duration) provenanceCounter {
	t.Helper()
	var decoded provenanceCounter
	require.NoError(t, json.Unmarshal([]byte(readCell(t, f, "tag-uint", bucket, "provenance").text.String), &decoded))
	return decoded
}

func exactOf(t *testing.T, digits string) *measurement.ExactValue {
	t.Helper()
	value, err := measurement.NewDecimal(digits)
	require.NoError(t, err)
	return &value
}

// closeWindows accepts the readings of each ten minute window in order, closes
// it, and finally commits the captured rows to the real SQLite table.
func closeWindows(t *testing.T, f *longFixture, tagIndex int, windows ...map[time.Duration]any) {
	t.Helper()
	boundary := f.build(t)
	for window, readings := range windows {
		f.clock = at(time.Duration(window+1) * counterWindow)
		for offset, value := range readings {
			start := time.Duration(window) * counterWindow
			id := fmt.Sprintf("w%d-%d", window, offset/time.Second)
			f.accept(t, boundary, f.envelope(tagIndex, id, at(start+offset), value))
		}
		require.NoError(t, boundary.Tick(t.Context(), f.clock))
	}
	f.insertRows(t)
}

// TestCumulativeProfileResetWithinWindowThroughBoundary drives 100, 10, 110 in
// one window and reads the committed record back from SQLite.
func TestCumulativeProfileResetWithinWindowThroughBoundary(t *testing.T) {
	f := newCounterFixture(t)
	closeWindows(t, f, 1,
		map[time.Duration]any{9 * time.Minute: uint64(90)},
		map[time.Duration]any{2 * time.Minute: uint64(100), 5 * time.Minute: uint64(10), 8 * time.Minute: uint64(110)},
		map[time.Duration]any{5 * time.Minute: uint64(125)},
	)
	require.Equal(t, "text:reset_or_uncertain", readCell(t, f, "tag-uint", counterWindow, "status").String())
	require.Equal(t, "text:110", readCell(t, f, "tag-uint", counterWindow, "value_decimal").String(), "the raw reading is kept")
	for _, column := range []string{"delta_value", "delta_from", "delta_to"} {
		require.Equal(t, "null:", readCell(t, f, "tag-uint", counterWindow, column).String(), column)
	}
	provenance := readProvenance(t, f, counterWindow)
	require.Equal(t, "decrease-unproven", provenance.Reason)
	require.Len(t, provenance.Counter.Events, 1)
	event := provenance.Counter.Events[0]
	require.Equal(t, "decrease", event.Kind)
	require.Equal(t, "100", event.FromValue)
	require.Equal(t, "10", event.ToValue)
	require.Equal(t, at(counterWindow+2*time.Minute).Format(longTimeLayout), event.From)
	require.Equal(t, at(counterWindow+5*time.Minute).Format(longTimeLayout), event.To)

	// The uncertain change is not counted: the next window starts from the last reliable 110.
	require.Equal(t, "text:ok", readCell(t, f, "tag-uint", 2*counterWindow, "status").String())
	require.Equal(t, "text:15", readCell(t, f, "tag-uint", 2*counterWindow, "delta_value").String())
}

func TestCumulativeProfileSignedAndRolloverThroughBoundary(t *testing.T) {
	t.Run("a confirmed signed counter records 10 to 7 as -3", func(t *testing.T) {
		f := newCounterFixture(t)
		source := f.config.Counters["tag-int"]
		source.Kind = aggregation.CounterKindSigned
		f.config.Counters["tag-int"] = source
		closeWindows(t, f, 0,
			map[time.Duration]any{9 * time.Minute: int64(10)},
			map[time.Duration]any{9 * time.Minute: int64(7)},
		)
		require.Equal(t, "text:ok", readCell(t, f, "tag-int", counterWindow, "status").String())
		require.Equal(t, "text:-3", readCell(t, f, "tag-int", counterWindow, "delta_value").String())
		require.Equal(t, "integer:7", readCell(t, f, "tag-int", counterWindow, "value_int").String())
	})

	t.Run("an unconfirmed counter records the same drop as unknown", func(t *testing.T) {
		f := newCounterFixture(t)
		closeWindows(t, f, 0,
			map[time.Duration]any{9 * time.Minute: int64(10)},
			map[time.Duration]any{9 * time.Minute: int64(7)},
		)
		require.Equal(t, "text:reset_or_uncertain", readCell(t, f, "tag-int", counterWindow, "status").String())
		require.Equal(t, "null:", readCell(t, f, "tag-int", counterWindow, "delta_value").String())
	})

	rollover := func(t *testing.T, modulus, rate string) *longFixture {
		f := newCounterFixture(t)
		source := f.config.Counters["tag-uint"]
		source.RolloverModulus = exactOf(t, modulus)
		if rate != "" {
			source.MaxRatePerSecond = exactOf(t, rate)
		}
		f.config.Counters["tag-uint"] = source
		closeWindows(t, f, 1,
			map[time.Duration]any{9 * time.Minute: uint64(65530)},
			map[time.Duration]any{9 * time.Minute: uint64(4)},
		)
		return f
	}
	t.Run("65530 to 4 is 10 with a single-wrap bound", func(t *testing.T) {
		f := rollover(t, "65536", "100")
		require.Equal(t, "text:ok", readCell(t, f, "tag-uint", counterWindow, "status").String())
		require.Equal(t, "text:10", readCell(t, f, "tag-uint", counterWindow, "delta_value").String())
		require.Equal(t, "text:4", readCell(t, f, "tag-uint", counterWindow, "value_decimal").String())
		provenance := readProvenance(t, f, counterWindow)
		require.Equal(t, "65536", provenance.Counter.RolloverModulus)
		require.Equal(t, "100", provenance.Counter.MaxRatePerSecond)
		require.Len(t, provenance.Counter.Events, 1)
		require.Equal(t, "rollover", provenance.Counter.Events[0].Kind)
	})
	t.Run("a modulus alone is not evidence", func(t *testing.T) {
		f := rollover(t, "65536", "")
		require.Equal(t, "text:reset_or_uncertain", readCell(t, f, "tag-uint", counterWindow, "status").String())
		require.Equal(t, "null:", readCell(t, f, "tag-uint", counterWindow, "delta_value").String())
		require.Equal(t, "rollover-evidence-insufficient", readProvenance(t, f, counterWindow).Reason)
	})
}

func TestCumulativeProfileGapSpanThroughBoundary(t *testing.T) {
	f := newCounterFixture(t)
	closeWindows(t, f, 1,
		map[time.Duration]any{0: uint64(1000)},
		nil,
		map[time.Duration]any{0: uint64(1002)},
		map[time.Duration]any{0: uint64(1005)},
	)
	missing := func(column string) string { return readCell(t, f, "tag-uint", counterWindow, column).String() }
	require.Equal(t, "text:missing", missing("status"))
	require.Equal(t, "null:", missing("value_decimal"))
	require.Equal(t, "null:", missing("delta_value"), "the missing window gets none of the span")

	span := func(column string) string { return readCell(t, f, "tag-uint", 2*counterWindow, column).String() }
	require.Equal(t, "text:gap", span("status"))
	require.Equal(t, "text:1002", span("value_decimal"))
	require.Equal(t, "text:2", span("delta_value"))
	require.Equal(t, "text:"+at(0).Format(longTimeLayout), span("delta_from"))
	require.Equal(t, "text:"+at(2*counterWindow).Format(longTimeLayout), span("delta_to"), "the real interval of the span")
	provenance := readProvenance(t, f, 2*counterWindow)
	require.True(t, provenance.Counter.SpanTotal)
	require.Equal(t, "w0-0", provenance.Counter.BaselineSampleID)
	require.Equal(t, "w2-0", provenance.Counter.EndSampleID)

	next := func(column string) string { return readCell(t, f, "tag-uint", 3*counterWindow, column).String() }
	require.Equal(t, "text:ok", next("status"))
	require.Equal(t, "text:3", next("delta_value"))
	require.False(t, readProvenance(t, f, 3*counterWindow).Counter.SpanTotal, "an ordinary window is not a span")
}

func TestCumulativeProfileEventsStayBounded(t *testing.T) {
	f := newCounterFixture(t)
	readings := map[time.Duration]any{}
	for i := range 40 {
		value := uint64(100)
		if i%2 == 1 {
			value = 10
		}
		readings[time.Duration(i)*10*time.Second] = value
	}
	closeWindows(t, f, 1, map[time.Duration]any{9 * time.Minute: uint64(5)}, readings)
	provenance := readProvenance(t, f, counterWindow)
	require.Len(t, provenance.Counter.Events, 16, "a flapping counter cannot grow a record without limit")
	require.Equal(t, 4, provenance.Counter.EventsOmitted)
	require.Equal(t, "text:reset_or_uncertain", readCell(t, f, "tag-uint", counterWindow, "status").String())
}

func TestCumulativeProfileRefusesAnInvalidRolloverPolicy(t *testing.T) {
	for name, mutate := range map[string]func(*longFixture){
		"a float modulus": func(f *longFixture) {
			float, _ := measurement.NewFloat64(65536)
			source := f.config.Counters["tag-uint"]
			source.RolloverModulus = &float
			f.config.Counters["tag-uint"] = source
		},
		"rollover on a signed counter": func(f *longFixture) {
			source := f.config.Counters["tag-uint"]
			source.Kind, source.RolloverModulus = aggregation.CounterKindSigned, exactOf(t, "65536")
			f.config.Counters["tag-uint"] = source
		},
		"a rate without a modulus": func(f *longFixture) {
			source := f.config.Counters["tag-uint"]
			source.MaxRatePerSecond = exactOf(t, "100")
			f.config.Counters["tag-uint"] = source
		},
	} {
		t.Run(name, func(t *testing.T) {
			f := newCounterFixture(t)
			mutate(f)
			_, err := NewGroupBoundary(f.config)
			var blocked *GroupBoundaryError
			require.ErrorAs(t, err, &blocked)
			require.Equal(t, "counter-policy-invalid", blocked.Code)
		})
	}
}
