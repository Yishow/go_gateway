package snapshot

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/schema"
	"math"
	"testing"
	"time"
)

func TestMixedConfigRejectsAmbiguousOrForeignIntent(t *testing.T) {
	tests := map[string]func(*Config){
		"group interval":           func(c *Config) { c.Interval = time.Minute },
		"legacy first bucket":      func(c *Config) { c.FirstBucket = mixedTime(0) },
		"legacy until":             func(c *Config) { c.Until = mixedTime(60) },
		"no activation":            func(c *Config) { c.StartsAt = time.Time{} },
		"fractional cadence":       func(c *Config) { c.Members[0].Interval = 1500 * time.Millisecond },
		"unknown calculation":      func(c *Config) { c.Members[0].Calculation = "average" },
		"omitted calculation":      func(c *Config) { c.Members[0].Calculation = "" },
		"bool mean":                func(c *Config) { c.Members[0].SourceValueType = schema.DataTypeBool },
		"text mean":                func(c *Config) { c.Members[0].SourceValueType = schema.DataTypeString },
		"missing frozen type":      func(c *Config) { c.Members[0].SourceValueType = "" },
		"legacy reduction":         func(c *Config) { c.Members[0].SampleReduction = SampleReductionMean },
		"unbound delta":            func(c *Config) { c.Members[0].Calculation = CalculationDelta },
		"latest baseline":          func(c *Config) { c.Counters = map[string]aggregation.CounterBinding{"b": counterBinding("b")} },
		"mean baseline":            func(c *Config) { c.Counters = map[string]aggregation.CounterBinding{"a": counterBinding("a")} },
		"non nullable":             func(c *Config) { c.Storage.NullableValues = false },
		"no quality":               func(c *Config) { c.Storage.MemberQuality = false },
		"timestamp outside domain": func(c *Config) { c.StartsAt = time.Date(2300, 1, 1, 0, 0, 0, 0, time.UTC) },
		"minimum floor overflow":   func(c *Config) { c.StartsAt = time.Unix(0, math.MinInt64) },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			c := mixedConfig(ModeMixedWide)
			mutate(&c)
			_, err := NewAssembler(c)
			require.Error(t, err)
		})
	}
	c := mixedCounterConfig(ModeMixedWide)
	binding := c.Counters["c"]
	binding.GroupID = "foreign"
	c.Counters["c"] = binding
	_, err := NewAssembler(c)
	require.Error(t, err)
}

func TestMixedOverflowDoesNotAdvanceOrConsumeState(t *testing.T) {
	c := mixedConfig(ModeMixedWide)
	a, err := NewAssembler(c)
	require.NoError(t, err)
	mixedOffer(t, a, mixedSample("a", "keep", 10, 22))
	p := a.Plan(time.Date(2300, 1, 1, 0, 0, 0, 0, time.UTC))
	require.Error(t, p.Err)
	require.True(t, p.Empty())
	a.Commit(p)
	require.Equal(t, 1, a.OpenBuckets())
	progress, _ := a.MixedProgress()
	require.True(t, progress.ClosedThrough.IsZero())
	c.MaxFutureSkew = time.Duration(math.MaxInt64)
	_, err = NewAssembler(c)
	require.Error(t, err)
	sample := mixedSample("a", "outside", 10, 22)
	sample.ObservedAt = time.Date(2300, 1, 1, 0, 0, 0, 0, time.UTC)
	require.Equal(t, ReasonWindowArithmetic, a.Check(sample, sample.ObservedAt).Reason)
	require.Equal(t, ReasonWindowArithmetic, a.Check(mixedSample("a", "clock", 10, 22), time.Date(2300, 1, 1, 0, 0, 0, 0, time.UTC)).Reason)
}

func TestMixedLaterClosureArithmeticFailsWithoutPartialAdvance(t *testing.T) {
	c := mixedConfig(ModeMixedWide)
	c.Members = c.Members[:1]
	c.Members[0].Interval = time.Second
	// First end fits, but computing another due boundary cannot be represented.
	c.StartsAt = time.Unix(0, math.MaxInt64-2*int64(time.Second))
	a, err := NewAssembler(c)
	require.NoError(t, err)
	p := a.Plan(time.Unix(0, math.MaxInt64))
	require.Error(t, p.Err)
	var arithmetic *WindowArithmeticError
	require.ErrorAs(t, p.Err, &arithmetic)
	require.True(t, p.Empty())
	require.Empty(t, p.Outcomes)
	a.Commit(p)
	progress, _ := a.MixedProgress()
	require.True(t, progress.ClosedThrough.IsZero())
	c.AllowedLateness = 3 * time.Second
	_, err = NewAssembler(c)
	require.ErrorAs(t, err, &arithmetic)
}
