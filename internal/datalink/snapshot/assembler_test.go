package snapshot

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func fixtureConfig() Config {
	return Config{
		WorkspaceID:      "workspace-A",
		GroupID:          "group-G",
		GroupRevision:    "rev-1",
		DestinationScope: "scope-X",
		Interval:         fixtureInterval,
		MaxFutureSkew:    2 * time.Second,
		FirstBucket:      fixtureStart,
		Members: []Member{
			{Key: "temperature", Required: true},
			{Key: "pressure", Required: true},
		},
	}
}

func newFixtureAssembler(t *testing.T, mutate func(*Config)) *Assembler {
	t.Helper()
	config := fixtureConfig()
	if mutate != nil {
		mutate(&config)
	}
	assembler, err := NewAssembler(config)
	require.NoError(t, err)
	return assembler
}

func offerGood(t *testing.T, a *Assembler, id, member string, seconds, now int) {
	t.Helper()
	result := a.Offer(goodSample(id, member, seconds, int64(seconds)), at(now))
	require.Contains(t, []OfferOutcome{OfferSelected, OfferNotSelected}, result.Outcome, id)
}

func memberResult(t *testing.T, outcome Outcome, key string) MemberResult {
	t.Helper()
	for _, member := range outcome.Members {
		if member.MemberKey == key {
			return member
		}
	}
	require.Failf(t, "member missing", "no result for %s", key)
	return MemberResult{}
}

func TestFreshnessMissingAndSilentBucketCompleteRow(t *testing.T) {
	a := newFixtureAssembler(t, nil)
	offerGood(t, a, "t3", "temperature", 3, 4)
	offerGood(t, a, "t8", "temperature", 8, 9)
	offerGood(t, a, "p6", "pressure", 6, 7)

	require.Empty(t, a.Tick(at(9)), "bucket is still open before its end")
	outcomes := a.Tick(at(10))
	require.Len(t, outcomes, 1)
	row := outcomes[0]
	require.Equal(t, OutcomeRow, row.Kind)
	require.Equal(t, fixtureStart, row.BucketStart)
	require.Equal(t, at(10), row.BucketEnd)
	require.Equal(t, "t8", memberResult(t, row, "temperature").Sample.SampleID)
	require.Equal(t, MemberOK, memberResult(t, row, "pressure").Status)
	require.Equal(t, at(10), a.NextClose())
	require.Empty(t, a.Tick(at(10)), "a closed bucket finalizes only once")
}

func TestFreshnessMissingAndSilentBucketMissingMemberSkipsRow(t *testing.T) {
	a := newFixtureAssembler(t, nil)
	offerGood(t, a, "t5", "temperature", 5, 6)

	outcomes := a.Tick(at(10))
	require.Len(t, outcomes, 1)
	require.Equal(t, OutcomeSkipped, outcomes[0].Kind)
	require.Equal(t, ReasonIncompleteRequired, outcomes[0].Reason)
	require.Equal(t, MemberOK, memberResult(t, outcomes[0], "temperature").Status)
	pressure := memberResult(t, outcomes[0], "pressure")
	require.Equal(t, MemberMissing, pressure.Status)
	require.Nil(t, pressure.Sample)
}

func TestFreshnessMissingAndSilentBucketBadMemberSkipsRow(t *testing.T) {
	a := newFixtureAssembler(t, nil)
	offerGood(t, a, "t5", "temperature", 5, 6)
	a.Offer(Sample{SampleID: "p9", MemberKey: "pressure", ObservedAt: at(9), Quality: schema.QualityBad, QualityReason: "read-failed"}, at(9))
	a.Offer(goodSample("p4", "pressure", 4, 4), at(9))

	outcomes := a.Tick(at(10))
	require.Len(t, outcomes, 1)
	require.Equal(t, OutcomeSkipped, outcomes[0].Kind)
	pressure := memberResult(t, outcomes[0], "pressure")
	require.Equal(t, MemberBad, pressure.Status)
	require.Equal(t, "read-failed", pressure.Reason, "the older good reading must not hide the failed read")
}

func TestFreshnessMissingAndSilentBucketGoodQualityWithoutValueIsInvalid(t *testing.T) {
	a := newFixtureAssembler(t, nil)
	offerGood(t, a, "t5", "temperature", 5, 6)
	a.Offer(Sample{SampleID: "p5", MemberKey: "pressure", ObservedAt: at(5), Quality: schema.QualityGood}, at(6))

	outcomes := a.Tick(at(10))
	require.Equal(t, OutcomeSkipped, outcomes[0].Kind)
	require.Equal(t, MemberInvalid, memberResult(t, outcomes[0], "pressure").Status)
}

func TestFreshnessMissingAndSilentBucketStaleMemberUsesMaxAgeFromBucketEnd(t *testing.T) {
	configure := func(c *Config) { c.Members[0].MaxAge = 2 * time.Second }

	stale := newFixtureAssembler(t, configure)
	offerGood(t, stale, "t5", "temperature", 5, 6) // end-observed = 5s > 2s
	offerGood(t, stale, "p9", "pressure", 9, 9)
	outcomes := stale.Tick(at(10))
	require.Equal(t, OutcomeSkipped, outcomes[0].Kind)
	require.Equal(t, MemberStale, memberResult(t, outcomes[0], "temperature").Status)

	fresh := newFixtureAssembler(t, configure)
	offerGood(t, fresh, "t8", "temperature", 8, 9) // end-observed = 2s <= 2s
	offerGood(t, fresh, "p9", "pressure", 9, 9)
	require.Equal(t, OutcomeRow, fresh.Tick(at(10))[0].Kind)
}

func TestFreshnessMissingAndSilentBucketNoCarryForward(t *testing.T) {
	a := newFixtureAssembler(t, nil)
	offerGood(t, a, "t5", "temperature", 5, 6)
	offerGood(t, a, "p5", "pressure", 5, 6)
	offerGood(t, a, "t15", "temperature", 15, 16)

	outcomes := a.Tick(at(20))
	require.Len(t, outcomes, 2)
	require.Equal(t, OutcomeRow, outcomes[0].Kind)
	require.Equal(t, OutcomeSkipped, outcomes[1].Kind, "pressure from the previous bucket is never reused")
	require.Equal(t, MemberMissing, memberResult(t, outcomes[1], "pressure").Status)
}

func TestFreshnessMissingAndSilentBucketTimeDrivenNoData(t *testing.T) {
	a := newFixtureAssembler(t, nil)

	outcomes := a.Tick(at(30))
	require.Len(t, outcomes, 3, "three silent buckets each leave one scoped outcome")
	for i, outcome := range outcomes {
		require.Equal(t, OutcomeNoData, outcome.Kind)
		require.Equal(t, ReasonNoSamples, outcome.Reason)
		require.Equal(t, at(10*i), outcome.BucketStart)
		require.Equal(t, "", outcome.EntityKey)
	}
	require.Empty(t, a.Tick(at(30)))
	require.Empty(t, a.Tick(at(39)))
	require.Len(t, a.Tick(at(40)), 1)
}

func TestFreshnessMissingAndSilentBucketAllowedLateness(t *testing.T) {
	a := newFixtureAssembler(t, func(c *Config) { c.AllowedLateness = 3 * time.Second })
	offerGood(t, a, "t5", "temperature", 5, 6)

	require.Empty(t, a.Tick(at(12)), "still inside end+lateness")
	require.Equal(t, OfferSelected, a.Offer(goodSample("p9", "pressure", 9, 9), at(12)).Outcome, "late-but-allowed sample joins the open bucket")
	outcomes := a.Tick(at(13))
	require.Len(t, outcomes, 1)
	require.Equal(t, OutcomeRow, outcomes[0].Kind)
}

func TestFreshnessMissingAndSilentBucketCatchUpIsBoundedAndOrdered(t *testing.T) {
	a := newFixtureAssembler(t, func(c *Config) { c.MaxBucketsPerTick = 2 })

	first := a.Tick(at(50))
	require.Len(t, first, 2)
	second := a.Tick(at(50))
	require.Len(t, second, 2)
	third := a.Tick(at(50))
	require.Len(t, third, 1)
	require.Equal(t, fixtureStart, first[0].BucketStart)
	require.Equal(t, at(20), second[0].BucketStart)
	require.Equal(t, at(40), third[0].BucketStart)
	require.Empty(t, a.Tick(at(50)))
}

func TestFreshnessMissingAndSilentBucketRejectsInvalidConfig(t *testing.T) {
	for name, mutate := range map[string]func(*Config){
		"zero interval":        func(c *Config) { c.Interval = 0 },
		"fractional interval":  func(c *Config) { c.Interval = 1500 * time.Millisecond },
		"negative lateness":    func(c *Config) { c.AllowedLateness = -time.Second },
		"negative skew":        func(c *Config) { c.MaxFutureSkew = -time.Second },
		"unaligned first":      func(c *Config) { c.FirstBucket = at(3) },
		"zero first":           func(c *Config) { c.FirstBucket = time.Time{} },
		"no members":           func(c *Config) { c.Members = nil },
		"empty member key":     func(c *Config) { c.Members[0].Key = "" },
		"duplicate member":     func(c *Config) { c.Members[1].Key = "temperature" },
		"negative max age":     func(c *Config) { c.Members[0].MaxAge = -time.Second },
		"optional in skip_row": func(c *Config) { c.Members[0].Required = false },
		"unknown policy":       func(c *Config) { c.IncompletePolicy = "carry_forward" },
	} {
		config := fixtureConfig()
		config.Members = append([]Member(nil), config.Members...)
		mutate(&config)
		_, err := NewAssembler(config)
		require.ErrorIs(t, err, ErrInvalidConfig, name)
	}
}

func TestFreshnessMissingAndSilentBucketFutureSampleDoesNotHoldBucketOpen(t *testing.T) {
	a := newFixtureAssembler(t, nil)
	result := a.Offer(goodSample("far", "temperature", 3600, 1), at(5))
	require.Equal(t, OfferRejected, result.Outcome)
	require.Equal(t, ReasonFutureSkew, result.Reason)
	require.Zero(t, a.OpenBuckets(), "a rejected future sample creates no bucket")
	require.Len(t, a.Tick(at(10)), 1)
}

func TestFreshnessMissingAndSilentBucketUnknownMember(t *testing.T) {
	a := newFixtureAssembler(t, nil)
	result := a.Offer(Sample{SampleID: "x", MemberKey: "flow", ObservedAt: at(2), Quality: schema.QualityGood, Value: measurement.NewInt64(1)}, at(3))
	require.Equal(t, OfferRejected, result.Outcome)
	require.Equal(t, ReasonUnknownMember, result.Reason)
}
