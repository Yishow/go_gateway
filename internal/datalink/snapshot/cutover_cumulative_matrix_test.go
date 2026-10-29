package snapshot

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"

	"github.com/stretchr/testify/require"
)

// These are the existing saved-profile intervals: long cumulative ten minutes
// and wide cumulative fifteen seconds. Both partition the very same readings.
func cutoverCounterAssembler(t *testing.T, mode Mode, successor bool, missingRequired *bool) *Assembler {
	t.Helper()
	c := counterConfig(nil)
	c.Mode, c.Interval = mode, 10*time.Minute
	if mode == ModeWideCumulative {
		c.Interval = 15 * time.Second
	}
	c.FirstBucket = at(600).Add(-c.Interval)
	c.Members = c.Members[:1]
	delete(c.Counters, "beta")
	c.Storage = StorageCapabilities{NullableValues: true, MemberQuality: true}
	if missingRequired != nil {
		c.Members = append(c.Members, Member{Key: "beta", Required: *missingRequired})
		c.Counters["beta"] = counterBinding("beta")
	}
	if successor {
		c.FirstBucket, c.StartsAt = at(600), at(612).Add(500*time.Millisecond)
		c.GroupID = "successor"
		for key, binding := range c.Counters {
			binding.GroupID = c.GroupID
			c.Counters[key] = binding
		}
	}
	a, err := NewAssembler(c)
	require.NoError(t, err)
	return a
}

func cutoverCounterAccepted() []Sample {
	exact := counterSample("exact", "alpha", 612, 1)
	exact.ObservedAt = at(612).Add(500 * time.Millisecond)
	exact.ReceivedAt = exact.ObservedAt.Add(40 * time.Millisecond)
	return []Sample{
		counterSample("after", "alpha", 613, 5), exact,
		counterSample("latest-before", "alpha", 609, 110),
		counterSample("hidden-reset", "alpha", 605, 10),
		counterSample("rise-before", "alpha", 603, 100),
	}
}

func cutoverCounterOutcome(t *testing.T, plan Plan, member string) Outcome {
	t.Helper()
	for _, outcome := range plan.Outcomes {
		for _, result := range outcome.Members {
			if result.MemberKey == member {
				return outcome
			}
		}
	}
	t.Fatalf("missing outcome for %s", member)
	return Outcome{}
}

func TestCutoverCumulativeMatrixActualCheckpointBounds(t *testing.T) {
	for _, mode := range []Mode{ModeCumulativeProfile, ModeWideCumulative} {
		t.Run(string(mode), func(t *testing.T) {
			old := cutoverCounterAssembler(t, mode, false, nil)
			offerAll(t, old, counterSample("baseline", "alpha", 599, 90))
			old.Tick(at(600))
			offerAll(t, old, cutoverCounterAccepted()...)
			boundary := at(612).Add(500 * time.Millisecond)
			old.SetEndsAt(boundary)
			before := old.CounterCheckpoints()
			plan := old.Plan(at(1200))
			row := cutoverCounterOutcome(t, plan, "alpha")
			require.Equal(t, boundary, row.WindowEnd)
			require.Equal(t, "end", row.Truncated)
			require.Equal(t, "latest-before", plan.CounterCheckpoints()["alpha"].Reading.SampleID)
			require.Equal(t, boundary, plan.CounterCheckpoints()["alpha"].WindowEnd,
				"a counter checkpoint must stop at its actual owner end, not the nominal bucket end")
			require.Equal(t, before, old.CounterCheckpoints(), "planning never advances the durable baseline")
			old.Commit(plan)
			require.Equal(t, plan.CounterCheckpoints(), old.CounterCheckpoints())
		})
	}
}

func TestCutoverCumulativeMatrixFullSequenceAndNewBaseline(t *testing.T) {
	for _, mode := range []Mode{ModeCumulativeProfile, ModeWideCumulative} {
		t.Run(string(mode), func(t *testing.T) {
			old, next := cutoverCounterAssembler(t, mode, false, nil), cutoverCounterAssembler(t, mode, true, nil)
			offerAll(t, old, counterSample("baseline", "alpha", 599, 90))
			old.Tick(at(600))
			accepted := cutoverCounterAccepted()
			offerAll(t, old, accepted...)
			offerAll(t, next, accepted...)
			old.SetEndsAt(next.config.StartsAt)
			for _, tc := range []struct {
				owner *Assembler
				ids   []string
			}{{old, []string{"rise-before", "hidden-reset", "latest-before"}}, {next, []string{"exact", "after"}}} {
				entity := tc.owner.config.entityOf(tc.owner.config.Members[0])
				selector := tc.owner.selectorWithinActualBounds(at(600), entity, tc.owner.buckets[at(600)].selectors[entity])
				var ids []string
				for _, reading := range selector.counterReadings("alpha") {
					ids = append(ids, reading.SampleID)
				}
				require.ElementsMatch(t, tc.ids, ids, "every eligible pre-decimation input survives; equality belongs only to the successor")
			}
			oldRow := cutoverCounterOutcome(t, old.Plan(at(1200)), "alpha")
			result := memberResult(t, oldRow, "alpha")
			require.Equal(t, "latest-before", result.Sample.SampleID)
			require.Equal(t, "110", digitsOfExact(t, result.Counter.Raw.Value))
			require.Equal(t, aggregation.CounterUncertain, result.Counter.Status)
			require.Equal(t, aggregation.ReasonDecreaseUnproven, result.Counter.Reason)
			require.Nil(t, result.Counter.Delta, "an eligible hidden reset cannot be erased by latest-only decimation")
			require.Len(t, result.Counter.Events, 1, "successor-owned decreases cannot contaminate old evidence")
			require.Equal(t, "rise-before", result.Counter.Events[0].FromSampleID)
			require.Equal(t, "hidden-reset", result.Counter.Events[0].ToSampleID)
			require.Equal(t, "end", oldRow.Truncated)
			nextRow := cutoverCounterOutcome(t, next.Plan(at(600).Add(next.config.Interval)), "alpha")
			fresh := memberResult(t, nextRow, "alpha")
			require.Equal(t, next.config.StartsAt, nextRow.WindowStart)
			require.Equal(t, "start", nextRow.Truncated)
			require.Equal(t, aggregation.CounterBaseline, fresh.Counter.Status)
			require.Equal(t, aggregation.ReasonFirstReading, fresh.Counter.Reason)
			require.Nil(t, fresh.Counter.Delta, "the new owner never subtracts the old owner's 110")
			require.Equal(t, "after", fresh.Counter.Raw.SampleID)
			require.Equal(t, "5", digitsOfExact(t, fresh.Counter.Raw.Value))
			require.Empty(t, fresh.Counter.Events)
			require.Equal(t, cutoverCounterAccepted(), accepted, "projection cannot rewrite original accepted input")
		})
	}
}

func TestCutoverCumulativeMatrixMissingPolicyPreserved(t *testing.T) {
	for _, mode := range []Mode{ModeCumulativeProfile, ModeWideCumulative} {
		for _, required := range []bool{false, true} {
			name := string(mode) + "/optional"
			if required {
				name = string(mode) + "/required"
			}
			t.Run(name, func(t *testing.T) {
				a := cutoverCounterAssembler(t, mode, false, &required)
				offerAll(t, a, counterSample("a-baseline", "alpha", 599, 90), counterSample("b-baseline", "beta", 599, 7))
				a.Tick(at(600))
				before := a.CounterCheckpoints()["beta"]
				offerAll(t, a, counterSample("a-before", "alpha", 609, 100), counterSample("b-outside", "beta", 613, 200))
				a.SetEndsAt(at(612).Add(500 * time.Millisecond))
				plan := a.Plan(at(1200))
				row := cutoverCounterOutcome(t, plan, "beta")
				missing := memberResult(t, row, "beta")
				require.Equal(t, MemberMissing, missing.Status)
				require.Nil(t, missing.Sample)
				require.Nil(t, missing.Counter, "NULL, not successor data, zero delta, or an old carried value")
				if mode == ModeWideCumulative && required {
					require.Equal(t, OutcomeSkipped, row.Kind)
					require.Equal(t, ReasonIncompleteRequired, row.Reason)
				} else {
					require.Equal(t, OutcomeRow, row.Kind)
				}
				require.Equal(t, "end", row.Truncated)
				a.Commit(plan)
				require.Equal(t, before, a.CounterCheckpoints()["beta"], "a missing member does not advance its baseline")
			})
		}
	}
}
