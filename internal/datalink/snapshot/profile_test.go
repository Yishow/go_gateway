package snapshot

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func profileConfig(mutate func(*Config)) Config {
	config := fixtureConfig()
	config.Mode = ModeInstantaneousProfile
	// A profile group has no required members and no entity partition.
	config.Members = []Member{{Key: "alpha"}, {Key: "beta", MaxAge: 3 * time.Second}, {Key: "gamma"}}
	if mutate != nil {
		mutate(&config)
	}
	return config
}

func newProfileAssembler(t *testing.T, mutate func(*Config)) *Assembler {
	t.Helper()
	assembler, err := NewAssembler(profileConfig(mutate))
	require.NoError(t, err)
	return assembler
}

func outcomesByMember(t *testing.T, outcomes []Outcome) map[string]Outcome {
	t.Helper()
	byMember := map[string]Outcome{}
	for _, outcome := range outcomes {
		require.Len(t, outcome.Members, 1, "a profile record carries exactly one member")
		require.Equal(t, outcome.Members[0].MemberKey, outcome.EntityKey)
		byMember[outcome.EntityKey] = outcome
	}
	return byMember
}

func TestProfileModeEmitsOneRecordPerMemberEvenForSilentWindows(t *testing.T) {
	a := newProfileAssembler(t, nil)
	silent := outcomesByMember(t, a.Tick(at(int(fixtureInterval/time.Second))))
	require.Len(t, silent, 3, "an empty window still has a record for every member")
	for key, outcome := range silent {
		require.Equal(t, OutcomeRow, outcome.Kind, key)
		require.Equal(t, MemberMissing, outcome.Members[0].Status, key)
		require.Equal(t, ReasonNoSamples, outcome.Members[0].Reason, key)
		require.False(t, outcome.Partial)
		require.NotEmpty(t, outcome.EffectKey)
	}
}

func TestProfileModeIdentityNamesTheMember(t *testing.T) {
	first := newProfileAssembler(t, nil)
	second := newProfileAssembler(t, nil)
	end := at(int(fixtureInterval / time.Second))
	firstIDs, secondIDs := map[string]string{}, map[string]string{}
	for key, outcome := range outcomesByMember(t, first.Tick(end)) {
		firstIDs[key] = outcome.RecordID
	}
	for key, outcome := range outcomesByMember(t, second.Tick(end)) {
		secondIDs[key] = outcome.RecordID
	}
	require.Equal(t, firstIDs, secondIDs, "identity is a pure function of the configuration and window")
	require.Len(t, map[string]struct{}{firstIDs["alpha"]: {}, firstIDs["beta"]: {}, firstIDs["gamma"]: {}}, 3)

	revised := newProfileAssembler(t, func(c *Config) { c.GroupRevision = "rev-2" })
	for key, outcome := range outcomesByMember(t, revised.Tick(end)) {
		require.NotEqual(t, firstIDs[key], outcome.RecordID, "a new applied revision never reuses a record identity")
	}
}

func TestProfileModeSelectsNewestEligibleSample(t *testing.T) {
	good := goodSample
	bad := func(id, member string, seconds int) Sample {
		sample := goodSample(id, member, seconds, 0)
		sample.Quality, sample.QualityReason, sample.Value = schema.QualityBad, "read-timeout", measurement.ExactValue{}
		return sample
	}
	run := func(t *testing.T, order []Sample) map[string]Outcome {
		a := newProfileAssembler(t, nil)
		for _, sample := range order {
			result := a.Offer(sample, at(10))
			require.Contains(t, []OfferOutcome{OfferSelected, OfferNotSelected}, result.Outcome, sample.SampleID)
		}
		return outcomesByMember(t, a.Tick(at(int(fixtureInterval/time.Second))))
	}
	samples := []Sample{
		good("alpha-old", "alpha", 2, 1), bad("alpha-new-bad", "alpha", 6), good("alpha-newest", "alpha", 5, 9),
		bad("beta-only-bad", "beta", 7),
		good("gamma-a", "gamma", 4, 10), good("gamma-b", "gamma", 4, 20), // same instant: greater ID wins
	}
	reversed := []Sample{samples[5], samples[4], samples[3], samples[2], samples[1], samples[0]}
	forward := run(t, samples)
	backward := run(t, reversed)
	require.Equal(t, forward, backward, "arrival order never changes a record")

	alpha := forward["alpha"].Members[0]
	require.Equal(t, MemberOK, alpha.Status)
	require.Equal(t, "alpha-newest", alpha.Sample.SampleID, "the newest eligible reading wins over a newer bad one")
	beta := forward["beta"].Members[0]
	require.Equal(t, MemberBad, beta.Status)
	require.Equal(t, "read-timeout", beta.Reason)
	require.Equal(t, "gamma-b", forward["gamma"].Members[0].Sample.SampleID)
}

func TestProfileModeClassifiesStaleAndInvalid(t *testing.T) {
	a := newProfileAssembler(t, nil)
	stale := goodSample("beta-1", "beta", 2, 5) // 8s before the end with a 3s limit
	typeMismatch := goodSample("gamma-1", "gamma", 4, 0)
	typeMismatch.Quality, typeMismatch.QualityReason, typeMismatch.Value = schema.QualityBad, ReasonTypeMismatch, measurement.ExactValue{}
	missingValue := goodSample("alpha-1", "alpha", 3, 0)
	missingValue.Value = measurement.ExactValue{}
	for _, sample := range []Sample{stale, typeMismatch, missingValue} {
		require.Contains(t, []OfferOutcome{OfferSelected, OfferNotSelected}, a.Offer(sample, at(10)).Outcome)
	}
	outcomes := outcomesByMember(t, a.Tick(at(int(fixtureInterval/time.Second))))
	require.Equal(t, MemberStale, outcomes["beta"].Members[0].Status)
	require.Equal(t, MemberInvalid, outcomes["gamma"].Members[0].Status)
	require.Equal(t, ReasonTypeMismatch, outcomes["gamma"].Members[0].Reason)
	require.Equal(t, MemberInvalid, outcomes["alpha"].Members[0].Status)
}

func TestProfileModeKeepsClosedWindowsClosed(t *testing.T) {
	a := newProfileAssembler(t, nil)
	end := at(int(fixtureInterval / time.Second))
	a.Tick(end)
	late := goodSample("alpha-late", "alpha", 3, 1)
	require.Equal(t, OfferLate, a.Offer(late, end).Outcome)
	require.Equal(t, uint64(1), a.Late())
	require.Empty(t, a.Tick(end))
}

func TestProfileModeConfigValidation(t *testing.T) {
	_, err := NewAssembler(profileConfig(func(c *Config) { c.Members[0].EntityKey = "line-1" }))
	require.ErrorIs(t, err, ErrInvalidConfig, "a profile member is its own partition")
	_, err = NewAssembler(profileConfig(func(c *Config) { c.Mode = "bogus" }))
	require.ErrorIs(t, err, ErrInvalidConfig)

	// Wide rules are untouched: skip_row still needs every member required.
	wide := fixtureConfig()
	wide.Members[0].Required = false
	_, err = NewAssembler(wide)
	require.ErrorIs(t, err, ErrInvalidConfig)
}

func TestProfileRecordIDRequiresEveryPart(t *testing.T) {
	start := at(0)
	base := []string{"workspace-A", "group-G", "rev-1", "member-1", "instantaneous"}
	id, err := ProfileRecordID(base[0], base[1], base[2], base[3], base[4], start)
	require.NoError(t, err)
	for i := range base {
		parts := append([]string(nil), base...)
		parts[i] = ""
		_, err := ProfileRecordID(parts[0], parts[1], parts[2], parts[3], parts[4], start)
		require.ErrorIs(t, err, ErrInvalidConfig, "part %d", i)
		parts[i] = base[i] + "-changed"
		other, err := ProfileRecordID(parts[0], parts[1], parts[2], parts[3], parts[4], start)
		require.NoError(t, err)
		require.NotEqual(t, id, other, "part %d must take part in the identity", i)
	}
	_, err = ProfileRecordID(base[0], base[1], base[2], base[3], base[4], time.Time{})
	require.ErrorIs(t, err, ErrInvalidConfig)
	wide, err := RecordID(base[0], base[1], base[2], "", start)
	require.NoError(t, err)
	require.NotEqual(t, id, wide)
	next, err := ProfileRecordID(base[0], base[1], base[2], base[3], base[4], start.Add(time.Nanosecond))
	require.NoError(t, err)
	require.NotEqual(t, id, next)
}

func TestProfileModeKeepsCatchUpWorkBounded(t *testing.T) {
	// 3 members and a limit of 7 records per tick allow 2 windows per tick, so a
	// long pause never turns into one unbounded batch, and nothing is lost.
	a := newProfileAssembler(t, func(c *Config) { c.MaxBucketsPerTick = 7 })
	intervalSeconds := int(fixtureInterval / time.Second)
	now := at(10 * intervalSeconds)
	seen := map[string]int{}
	ticks := 0
	var starts []time.Time
	for {
		outcomes := a.Tick(now)
		if len(outcomes) == 0 {
			break
		}
		ticks++
		require.LessOrEqual(t, len(outcomes), 7, "a tick never exceeds the configured work bound")
		for _, outcome := range outcomes {
			seen[outcome.RecordID]++
			if outcome.EntityKey == "alpha" {
				starts = append(starts, outcome.BucketStart)
			}
		}
	}
	require.Equal(t, 5, ticks)
	require.Len(t, seen, 10*3, "every window of every member closes exactly once")
	for _, count := range seen {
		require.Equal(t, 1, count)
	}
	require.IsIncreasing(t, unixNanos(starts), "windows close in order")
}

func unixNanos(times []time.Time) []int64 {
	out := make([]int64, len(times))
	for i, value := range times {
		out[i] = value.UnixNano()
	}
	return out
}
