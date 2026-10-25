package snapshot

import (
	"testing"

	"go-gateway/internal/datalink/aggregation"

	"github.com/stretchr/testify/require"
)

func wideCounterConfig(mutate func(*Config)) Config {
	config := counterConfig(nil)
	config.Mode = ModeWideCumulative
	config.Storage = StorageCapabilities{NullableValues: true, MemberQuality: true}
	if mutate != nil {
		mutate(&config)
	}
	return config
}

func newWideCounterAssembler(t *testing.T, mutate func(*Config)) *Assembler {
	t.Helper()
	assembler, err := NewAssembler(wideCounterConfig(mutate))
	require.NoError(t, err)
	return assembler
}

func closeWideCounterWindow(t *testing.T, a *Assembler, n int) Outcome {
	t.Helper()
	outcomes := a.Tick(windowEnd(n))
	require.Len(t, outcomes, 1, "a wide cumulative group closes one record per window, not one per member")
	return outcomes[0]
}

// TestWideCumulativeRecordCarriesEveryMembersCounter pins that one record holds
// each member's raw reading beside its own exact difference, from the same engine
// as the per-member profile: a baseline first, then the exact difference.
func TestWideCumulativeRecordCarriesEveryMembersCounter(t *testing.T) {
	a := newWideCounterAssembler(t, nil)
	offerAll(t, a, counterSample("a0", "alpha", 9, 9007199254740992), counterSample("b0", "beta", 9, 5))
	first := closeWideCounterWindow(t, a, 1)
	require.Equal(t, OutcomeRow, first.Kind)
	require.False(t, first.Partial)
	for _, member := range first.Members {
		require.Equal(t, aggregation.CounterBaseline, member.Counter.Status, member.MemberKey)
		require.Nil(t, member.Counter.Delta, "the first window only sets the baseline")
	}
	require.Equal(t, "9007199254740992", digitsOfExact(t, first.Members[0].Sample.Value))

	offerAll(t, a, counterSample("a1", "alpha", 19, 9007199254740993), counterSample("b1", "beta", 19, 8))
	second := closeWideCounterWindow(t, a, 2)
	require.Equal(t, []string{"1", "3"}, []string{deltaDigitsOf(t, *second.Members[0].Counter), deltaDigitsOf(t, *second.Members[1].Counter)})
	require.Equal(t, aggregation.CounterOK, second.Members[0].Counter.Status)
	require.Equal(t, "9007199254740993", digitsOfExact(t, second.Members[0].Sample.Value), "the raw reading is always kept")
	require.NotEqual(t, first.RecordID, second.RecordID)
}

// TestWideCumulativeBaselineMovesOnlyWithTheCommittedPlan pins the durable
// contract the profile mode already has: planning never advances a baseline, a
// retried plan consumes nothing, and the plan carries what the closure persists.
func TestWideCumulativeBaselineMovesOnlyWithTheCommittedPlan(t *testing.T) {
	a := newWideCounterAssembler(t, nil)
	offerAll(t, a, counterSample("a0", "alpha", 9, 100), counterSample("b0", "beta", 9, 7))
	plan := a.Plan(windowEnd(1))
	require.Empty(t, a.CounterCheckpoints(), "planning does not advance any baseline")
	planned := plan.CounterCheckpoints()
	require.Equal(t, planned, a.Plan(windowEnd(1)).CounterCheckpoints(), "a retried plan carries the same baselines and consumes nothing")
	require.Empty(t, a.CounterCheckpoints())
	require.Len(t, planned, 2)
	require.Equal(t, "a0", planned["alpha"].Reading.SampleID)

	a.Commit(plan)
	require.Equal(t, planned, a.CounterCheckpoints())
}

func TestWideCumulativeResetIsNeverAbsoluteOrZero(t *testing.T) {
	a := newWideCounterAssembler(t, nil)
	offerAll(t, a, counterSample("a0", "alpha", 9, 100), counterSample("b0", "beta", 9, 7))
	closeWideCounterWindow(t, a, 1)
	offerAll(t, a, counterSample("a1", "alpha", 19, 40), counterSample("b1", "beta", 19, 9))
	row := closeWideCounterWindow(t, a, 2)
	require.Equal(t, aggregation.CounterUncertain, row.Members[0].Counter.Status)
	require.Nil(t, row.Members[0].Counter.Delta, "an unproven decrease has no difference: not abs(), not zero")
	require.Equal(t, "40", digitsOfExact(t, row.Members[0].Sample.Value))
	require.Equal(t, "2", deltaDigitsOf(t, *row.Members[1].Counter), "another member's result is independent")
}

// TestWideCumulativeMissingMemberIsNULLAndKeepsItsBaseline pins that a member
// without an eligible reading is an optional NULL beside members that have one,
// moves no baseline and invents no difference; only a required one skips the row.
func TestWideCumulativeMissingMemberIsNULLAndKeepsItsBaseline(t *testing.T) {
	a := newWideCounterAssembler(t, nil)
	offerAll(t, a, counterSample("a0", "alpha", 9, 100), counterSample("b0", "beta", 9, 7))
	closeWideCounterWindow(t, a, 1)
	before := a.CounterCheckpoints()["beta"]

	offerAll(t, a, counterSample("a1", "alpha", 19, 130))
	row := closeWideCounterWindow(t, a, 2)
	require.Equal(t, OutcomeRow, row.Kind)
	require.True(t, row.Partial)
	require.Equal(t, [2]any{MemberMissing, (*aggregation.CounterWindowResult)(nil)}, [2]any{row.Members[1].Status, row.Members[1].Counter})
	require.Equal(t, before, a.CounterCheckpoints()["beta"], "a missing window does not move the baseline")
	require.Equal(t, "30", deltaDigitsOf(t, *row.Members[0].Counter))

	required := newWideCounterAssembler(t, func(c *Config) { c.Members[1].Required = true })
	offerAll(t, required, counterSample("a0", "alpha", 9, 100))
	skipped := closeWideCounterWindow(t, required, 1)
	require.Equal(t, [2]string{string(OutcomeSkipped), ReasonIncompleteRequired}, [2]string{string(skipped.Kind), skipped.Reason})
}

func TestWideCumulativeRestoreAndConfigFailClosed(t *testing.T) {
	first := newWideCounterAssembler(t, nil)
	offerAll(t, first, counterSample("a0", "alpha", 9, 100), counterSample("b0", "beta", 9, 7))
	closeWideCounterWindow(t, first, 1)

	second := newWideCounterAssembler(t, func(c *Config) { c.FirstBucket = windowEnd(1) })
	require.NoError(t, second.RestoreCounterCheckpoints(first.CounterCheckpoints()))
	offerAll(t, second, counterSample("a1", "alpha", 19, 101), counterSample("b1", "beta", 19, 9))
	require.Equal(t, "1", deltaDigitsOf(t, *closeWideCounterWindow(t, second, 2).Members[0].Counter), "the restored baseline is subtracted from")

	instantaneous := newWideProfileAssembler(t, nil)
	require.ErrorIs(t, instantaneous.RestoreCounterCheckpoints(first.CounterCheckpoints()), ErrInvalidConfig,
		"baselines belong to the cumulative modes only")

	for name, mutate := range map[string]func(*Config){
		"a member without a binding":    func(c *Config) { delete(c.Counters, "beta") },
		"a mean as a counter baseline":  func(c *Config) { c.Members[0].SampleReduction = SampleReductionMean },
		"storage that cannot hold NULL": func(c *Config) { c.Storage = StorageCapabilities{} },
		"a member that is entity-keyed": func(c *Config) { c.Members[1].EntityKey = "plant-A" },
		"a binding of another member": func(c *Config) {
			binding := c.Counters["beta"]
			binding.MemberKey = "alpha"
			c.Counters["beta"] = binding
		},
		"a binding of another group rev.": func(c *Config) {
			binding := c.Counters["beta"]
			binding.ProfileRevision = "rev-0"
			c.Counters["beta"] = binding
		},
	} {
		_, err := NewAssembler(wideCounterConfig(mutate))
		require.Error(t, err, name)
	}
	_, err := NewAssembler(wideProfileConfig(func(c *Config) { c.Counters = counterConfig(nil).Counters }))
	require.ErrorIs(t, err, ErrInvalidConfig, "the instantaneous wide mode refuses counter bindings")
}
