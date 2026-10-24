package snapshot

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func counterBinding(member string) aggregation.CounterBinding {
	return aggregation.CounterBinding{
		WorkspaceID: "workspace-A", GroupID: "group-G", ProfileRevision: "rev-1", MemberKey: member,
		MeasurementID: "meas-" + member, DefinitionRevision: "def-1", SourceBindingRevision: "bind-1",
		SeriesEpoch: "epoch-1", Kind: aggregation.CounterKindMonotonic, Unit: "kWh",
	}
}

func counterConfig(mutate func(*Config)) Config {
	config := fixtureConfig()
	config.Mode = ModeCumulativeProfile
	config.Members = []Member{{Key: "alpha"}, {Key: "beta"}}
	config.Counters = map[string]aggregation.CounterBinding{"alpha": counterBinding("alpha"), "beta": counterBinding("beta")}
	if mutate != nil {
		mutate(&config)
	}
	return config
}

func newCounterAssembler(t *testing.T, mutate func(*Config)) *Assembler {
	t.Helper()
	assembler, err := NewAssembler(counterConfig(mutate))
	require.NoError(t, err)
	return assembler
}

// counterSample is a good uint64 reading observed seconds after the fixture start.
func counterSample(id, member string, seconds int, value uint64) Sample {
	sample := goodSample(id, member, seconds, 0)
	sample.Value = measurement.NewUint64(value)
	sample.ReceivedAt = sample.ObservedAt.Add(40 * time.Millisecond)
	return sample
}

func offerAll(t *testing.T, a *Assembler, samples ...Sample) {
	t.Helper()
	for _, sample := range samples {
		result := a.Offer(sample, sample.ObservedAt)
		require.Contains(t, []OfferOutcome{OfferSelected, OfferNotSelected}, result.Outcome, sample.SampleID)
	}
}

func windowEnd(n int) time.Time { return at(n * int(fixtureInterval/time.Second)) }

func deltaDigits(t *testing.T, outcome Outcome) string {
	t.Helper()
	counter := outcome.Members[0].Counter
	require.NotNil(t, counter)
	require.NotNil(t, counter.Delta, "status %s reason %s", counter.Status, counter.Reason)
	digits, ok := counter.Delta.Value.DecimalDigits()
	require.True(t, ok)
	return digits
}

func TestCumulativeModeKeepsBaselineInThePlanUntilCommit(t *testing.T) {
	a := newCounterAssembler(t, nil)
	offerAll(t, a, counterSample("a0", "alpha", 9, 9007199254740992), counterSample("b0", "beta", 9, 5))
	offerAll(t, a, counterSample("a1", "alpha", 19, 9007199254740993), counterSample("b1", "beta", 19, 8))

	plan := a.Plan(windowEnd(2))
	require.Len(t, plan.Outcomes, 4, "two windows of two members")
	require.Empty(t, a.CounterCheckpoints(), "planning does not advance any baseline")
	again := a.Plan(windowEnd(2))
	require.Equal(t, outcomesByMember(t, plan.Outcomes[2:])["alpha"].RecordID, outcomesByMember(t, again.Outcomes[2:])["alpha"].RecordID)
	require.Empty(t, a.CounterCheckpoints(), "a plan that is retried does not consume anything either")

	planned := plan.CounterCheckpoints()
	require.Len(t, planned, 2)
	require.Equal(t, "a1", planned["alpha"].Reading.SampleID, "the plan carries the baselines to persist with the closure")
	require.Equal(t, windowEnd(2), planned["alpha"].WindowEnd)

	byWindow := func(outcomes []Outcome) map[string]Outcome { return outcomesByMember(t, outcomes) }
	first, second := byWindow(plan.Outcomes[:2]), byWindow(plan.Outcomes[2:])
	require.Equal(t, aggregation.CounterBaseline, first["alpha"].Members[0].Counter.Status)
	require.Nil(t, first["alpha"].Members[0].Counter.Delta)
	require.Equal(t, "1", deltaDigits(t, second["alpha"]), "9007199254740992 to 9007199254740993 is exactly 1")
	require.Equal(t, "3", deltaDigits(t, second["beta"]), "members are compared independently")
	require.Equal(t, "9007199254740993", func() string {
		digits, _ := second["alpha"].Members[0].Counter.Raw.Value.DecimalDigits()
		return digits
	}())

	a.Commit(plan)
	require.Len(t, a.CounterCheckpoints(), 2, "a committed plan installs the baselines")
	require.Equal(t, windowEnd(2), a.NextClose())
}

func TestCumulativeModeBaselineSurvivesSerializedRestore(t *testing.T) {
	first := newCounterAssembler(t, nil)
	offerAll(t, first, counterSample("a0", "alpha", 9, 9007199254740992), counterSample("b0", "beta", 9, 40))
	first.Tick(windowEnd(1))
	stored := map[string][]byte{}
	for member, checkpoint := range first.CounterCheckpoints() {
		encoded, err := aggregation.EncodeCounterCheckpoint(checkpoint)
		require.NoError(t, err)
		stored[member] = encoded
	}

	// A new process starts at the next window and restores the baselines it saved.
	second := newCounterAssembler(t, func(c *Config) { c.FirstBucket = windowEnd(1) })
	restored := map[string]aggregation.CounterCheckpoint{}
	for member, encoded := range stored {
		checkpoint, err := aggregation.DecodeCounterCheckpoint(encoded)
		require.NoError(t, err)
		restored[member] = checkpoint
	}
	require.NoError(t, second.RestoreCounterCheckpoints(restored))
	offerAll(t, second, counterSample("a1", "alpha", 19, 9007199254740993), counterSample("b1", "beta", 19, 45))
	closed := outcomesByMember(t, second.Tick(windowEnd(2)))
	require.Equal(t, "1", deltaDigits(t, closed["alpha"]))
	require.Equal(t, "5", deltaDigits(t, closed["beta"]))
	require.Equal(t, aggregation.CounterOK, closed["alpha"].Members[0].Counter.Status)

	// Without the saved baseline the same reading is only a baseline: no delta is invented.
	fresh := newCounterAssembler(t, func(c *Config) { c.FirstBucket = windowEnd(1) })
	offerAll(t, fresh, counterSample("a1", "alpha", 19, 9007199254740993))
	lost := outcomesByMember(t, fresh.Tick(windowEnd(2)))["alpha"].Members[0].Counter
	require.Equal(t, aggregation.CounterBaseline, lost.Status)
	require.Nil(t, lost.Delta)

	require.Error(t, second.RestoreCounterCheckpoints(map[string]aggregation.CounterCheckpoint{"stranger": restored["alpha"]}))
}

func TestCumulativeModeJudgesTheWholeSequence(t *testing.T) {
	a := newCounterAssembler(t, nil)
	offerAll(t, a, counterSample("a0", "alpha", 9, 90))
	a.Tick(windowEnd(1))

	// 100, 10, 110 with the newest observation being a bad read: the bad read does
	// not hide the newest good value, and the reset in the middle is not hidden by
	// comparing only the endpoints.
	bad := counterSample("a-bad", "alpha", 18, 0)
	bad.Quality, bad.QualityReason, bad.Value = schema.QualityBad, "read-timeout", measurement.ExactValue{}
	offerAll(t, a,
		counterSample("a1", "alpha", 12, 100), counterSample("a2", "alpha", 14, 10), counterSample("a3", "alpha", 16, 110), bad,
	)
	closed := outcomesByMember(t, a.Tick(windowEnd(2)))["alpha"]
	counter := closed.Members[0].Counter
	require.Equal(t, aggregation.CounterUncertain, counter.Status)
	require.Equal(t, aggregation.ReasonDecreaseUnproven, counter.Reason)
	require.Nil(t, counter.Delta)
	digits, _ := counter.Raw.Value.DecimalDigits()
	require.Equal(t, "110", digits, "raw is the newest good reading")
	require.Equal(t, MemberOK, closed.Members[0].Status)
	require.Equal(t, "a3", closed.Members[0].Sample.SampleID)
}

func TestCumulativeModeSilentWindowKeepsTheBaseline(t *testing.T) {
	a := newCounterAssembler(t, nil)
	offerAll(t, a, counterSample("a0", "alpha", 9, 1000))
	a.Tick(windowEnd(1))
	before := a.CounterCheckpoints()["alpha"]

	silent := outcomesByMember(t, a.Tick(windowEnd(2)))["alpha"]
	require.Equal(t, OutcomeRow, silent.Kind)
	require.Equal(t, MemberMissing, silent.Members[0].Status)
	require.Nil(t, silent.Members[0].Counter, "an absent reading is NULL, not a zero delta")
	require.Equal(t, before, a.CounterCheckpoints()["alpha"], "a missing window does not move the baseline")

	offerAll(t, a, counterSample("a2", "alpha", 29, 1002))
	gap := outcomesByMember(t, a.Tick(windowEnd(3)))["alpha"].Members[0].Counter
	require.Equal(t, aggregation.CounterGap, gap.Status)
	require.True(t, gap.Delta.Span, "the difference across the silent window is a span total")
	require.Equal(t, "2", deltaDigitsOf(t, *gap))
	require.Equal(t, at(9), gap.Delta.From)
	require.Equal(t, at(29), gap.Delta.To)
}

func TestCumulativeModeIdentityDiffersFromInstantaneous(t *testing.T) {
	cumulative := newCounterAssembler(t, nil)
	instantaneous := newProfileAssembler(t, func(c *Config) { c.Members = []Member{{Key: "alpha"}, {Key: "beta"}} })
	end := windowEnd(1)
	cumulativeID := outcomesByMember(t, cumulative.Tick(end))["alpha"].RecordID
	instantaneousID := outcomesByMember(t, instantaneous.Tick(end))["alpha"].RecordID
	require.NotEmpty(t, cumulativeID)
	require.NotEqual(t, cumulativeID, instantaneousID,
		"the two profiles of one source and window never share a record identity")
}

func TestCumulativeModeConfigValidation(t *testing.T) {
	for name, mutate := range map[string]func(*Config){
		"a member without a binding": func(c *Config) { delete(c.Counters, "beta") },
		"an unconfirmed counter kind": func(c *Config) {
			binding := c.Counters["beta"]
			binding.Kind = ""
			c.Counters["beta"] = binding
		},
		"a binding of another group revision": func(c *Config) {
			binding := c.Counters["beta"]
			binding.ProfileRevision = "rev-0"
			c.Counters["beta"] = binding
		},
		"a binding of another member": func(c *Config) {
			binding := c.Counters["beta"]
			binding.MemberKey = "alpha"
			c.Counters["beta"] = binding
		},
		"no bindings at all": func(c *Config) { c.Counters = nil },
	} {
		t.Run(name, func(t *testing.T) {
			_, err := NewAssembler(counterConfig(mutate))
			require.ErrorIs(t, err, ErrInvalidConfig)
		})
	}
	_, err := NewAssembler(profileConfig(func(c *Config) {
		c.Counters = map[string]aggregation.CounterBinding{"alpha": counterBinding("alpha")}
	}))
	require.ErrorIs(t, err, ErrInvalidConfig, "counter bindings are not an instantaneous concept")
	require.Error(t, newProfileAssembler(t, nil).RestoreCounterCheckpoints(nil),
		"baselines can only be restored into a cumulative assembler")
}

// TestCumulativeModeRestoredBaselineOfAnotherBindingIsNotSubtracted restores
// baselines taken under one binding into an assembler configured with a changed
// series epoch or rollover policy: the first window rebuilds the baseline.
func TestCumulativeModeRestoredBaselineOfAnotherBindingIsNotSubtracted(t *testing.T) {
	old := newCounterAssembler(t, nil)
	offerAll(t, old, counterSample("a0", "alpha", 9, 1000))
	old.Tick(windowEnd(1))
	saved := old.CounterCheckpoints()

	modulus := measurement.NewUint64(65536)
	rate := measurement.NewUint64(100)
	for name, mutate := range map[string]func(*aggregation.CounterBinding){
		"series epoch":    func(b *aggregation.CounterBinding) { b.SeriesEpoch = "epoch-2" },
		"definition":      func(b *aggregation.CounterBinding) { b.DefinitionRevision = "def-2" },
		"source binding":  func(b *aggregation.CounterBinding) { b.SourceBindingRevision = "bind-2" },
		"rollover policy": func(b *aggregation.CounterBinding) { b.RolloverModulus, b.MaxRatePerSecond = &modulus, &rate },
	} {
		t.Run(name, func(t *testing.T) {
			changed := newCounterAssembler(t, func(c *Config) {
				c.FirstBucket = windowEnd(1)
				binding := c.Counters["alpha"]
				mutate(&binding)
				c.Counters["alpha"] = binding
			})
			require.NoError(t, changed.RestoreCounterCheckpoints(saved))
			offerAll(t, changed, counterSample("a1", "alpha", 19, 1005))
			counter := outcomesByMember(t, changed.Tick(windowEnd(2)))["alpha"].Members[0].Counter
			require.Equal(t, aggregation.CounterBaseline, counter.Status)
			require.Equal(t, aggregation.ReasonBindingChanged, counter.Reason)
			require.Nil(t, counter.Delta, "the old baseline is never reused under a changed meaning")
		})
	}
}

// TestCumulativeModeCarriesEventsAndSpanThroughThePlan shows the discontinuity
// evidence and the span total travel in the closed outcome a durable caller
// freezes, and that a missing window leaves the baseline in the plan untouched.
func TestCumulativeModeCarriesEventsAndSpanThroughThePlan(t *testing.T) {
	a := newCounterAssembler(t, nil)
	offerAll(t, a, counterSample("a0", "alpha", 9, 90))
	offerAll(t, a, counterSample("a1", "alpha", 12, 100), counterSample("a2", "alpha", 15, 10), counterSample("a3", "alpha", 18, 110))
	offerAll(t, a, counterSample("a5", "alpha", 49, 120))

	plan := a.Plan(windowEnd(5))
	byWindow := map[time.Time]aggregation.CounterWindowResult{}
	for _, outcome := range plan.Outcomes {
		if outcome.EntityKey != "alpha" {
			continue
		}
		if counter := outcome.Members[0].Counter; counter != nil {
			byWindow[outcome.BucketStart] = *counter
		}
	}
	reset := byWindow[windowEnd(1)]
	require.Equal(t, aggregation.CounterUncertain, reset.Status)
	require.Len(t, reset.Events, 1)
	require.Equal(t, aggregation.EventDecrease, reset.Events[0].Kind)

	span := byWindow[windowEnd(4)]
	require.Equal(t, aggregation.CounterGap, span.Status, "windows 2 to 4 were silent")
	require.True(t, span.Delta.Span)
	require.Equal(t, "10", deltaDigitsOf(t, span), "120 against the reset window's last reliable 110")
	require.Equal(t, at(18), span.Delta.From)
	require.Equal(t, at(49), span.Delta.To)
	require.Equal(t, "a5", plan.CounterCheckpoints()["alpha"].Reading.SampleID)
}

func deltaDigitsOf(t *testing.T, counter aggregation.CounterWindowResult) string {
	t.Helper()
	require.NotNil(t, counter.Delta)
	digits, ok := counter.Delta.Value.DecimalDigits()
	require.True(t, ok)
	return digits
}

// TestCumulativeModeRestoreRefusesACheckpointOfAnotherScope keeps the fixed
// scope (workspace, group, profile revision) out of the "binding changed" path.
func TestCumulativeModeRestoreRefusesACheckpointOfAnotherScope(t *testing.T) {
	source := newCounterAssembler(t, nil)
	offerAll(t, source, counterSample("a0", "alpha", 9, 1000))
	source.Tick(windowEnd(1))
	saved := source.CounterCheckpoints()["alpha"]

	for name, mutate := range map[string]func(*aggregation.CounterCheckpoint){
		"workspace":        func(c *aggregation.CounterCheckpoint) { c.Binding.WorkspaceID = "workspace-B" },
		"group":            func(c *aggregation.CounterCheckpoint) { c.Binding.GroupID = "group-other" },
		"profile revision": func(c *aggregation.CounterCheckpoint) { c.Binding.ProfileRevision = "rev-0" },
		"member":           func(c *aggregation.CounterCheckpoint) { c.Binding.MemberKey = "beta" },
	} {
		t.Run(name, func(t *testing.T) {
			foreign := saved
			mutate(&foreign)
			target := newCounterAssembler(t, func(c *Config) { c.FirstBucket = windowEnd(1) })
			require.ErrorIs(t, target.RestoreCounterCheckpoints(map[string]aggregation.CounterCheckpoint{"alpha": foreign}), ErrInvalidConfig)
			require.Empty(t, target.CounterCheckpoints(), "nothing of a refused restore is installed")
		})
	}
}
