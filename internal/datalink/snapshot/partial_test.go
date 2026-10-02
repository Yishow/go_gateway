package snapshot

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func partialConfig(mutate func(*Config)) func(*Config) {
	return func(c *Config) {
		c.IncompletePolicy = IncompletePartial
		c.Storage = StorageCapabilities{NullableValues: true, MemberQuality: true}
		c.Members[0].Required = true
		c.Members[1].Required = false
		if mutate != nil {
			mutate(c)
		}
	}
}

func TestExplicitPartialPolicyRequiresVerifiedStorage(t *testing.T) {
	for name, storage := range map[string]StorageCapabilities{
		"no capabilities":     {},
		"nullable only":       {NullableValues: true},
		"member quality only": {MemberQuality: true},
	} {
		config := fixtureConfig()
		config.Members = append([]Member(nil), config.Members...)
		partialConfig(nil)(&config)
		config.Storage = storage
		_, err := NewAssembler(config)
		require.ErrorIs(t, err, ErrPartialUnsupported, name)
		require.ErrorIs(t, err, ErrInvalidConfig, name)
	}
	_, err := NewAssembler(func() Config {
		config := fixtureConfig()
		config.Members = append([]Member(nil), config.Members...)
		partialConfig(nil)(&config)
		return config
	}())
	require.NoError(t, err)
}

func TestExplicitPartialPolicyMissingMemberIsNullWithReason(t *testing.T) {
	a := newFixtureAssembler(t, partialConfig(nil))
	offerGood(t, a, "t5", "temperature", 5, 6)

	outcomes := a.Tick(at(10))
	require.Len(t, outcomes, 1)
	row := outcomes[0]
	require.Equal(t, OutcomeRow, row.Kind)
	require.True(t, row.Partial)
	require.True(t, memberResult(t, row, "temperature").Usable())
	pressure := memberResult(t, row, "pressure")
	require.False(t, pressure.Usable())
	require.Equal(t, MemberMissing, pressure.Status)
	require.Equal(t, ReasonNoSamples, pressure.Reason)
	require.Nil(t, pressure.Sample)
}

func TestExplicitPartialPolicyBadMemberKeepsProvenanceButNoValue(t *testing.T) {
	a := newFixtureAssembler(t, partialConfig(nil))
	offerGood(t, a, "t5", "temperature", 5, 6)
	a.Offer(Sample{
		SampleID: "p9", MemberKey: "pressure", ObservedAt: at(9), Quality: schema.QualityBad,
		QualityReason: "read-failed", Value: measurement.NewInt64(999),
	}, at(9))

	row := a.Tick(at(10))[0]
	require.Equal(t, OutcomeRow, row.Kind)
	require.True(t, row.Partial)
	pressure := memberResult(t, row, "pressure")
	require.False(t, pressure.Usable(), "a bad reading is never written as a value")
	require.Equal(t, "read-failed", pressure.Reason)
	require.Equal(t, "p9", pressure.Sample.SampleID, "observation metadata is kept for the local envelope")
}

func TestExplicitPartialPolicyStaleMemberIsNull(t *testing.T) {
	a := newFixtureAssembler(t, partialConfig(func(c *Config) { c.Members[1].MaxAge = 2 * time.Second }))
	offerGood(t, a, "t9", "temperature", 9, 9)
	offerGood(t, a, "p3", "pressure", 3, 4)

	row := a.Tick(at(10))[0]
	require.True(t, row.Partial)
	require.Equal(t, MemberStale, memberResult(t, row, "pressure").Status)
	require.False(t, memberResult(t, row, "pressure").Usable())
}

func TestExplicitPartialPolicyRequiredMemberStillGatesRow(t *testing.T) {
	a := newFixtureAssembler(t, partialConfig(nil))
	offerGood(t, a, "p5", "pressure", 5, 6)

	outcomes := a.Tick(at(10))
	require.Len(t, outcomes, 1)
	require.Equal(t, OutcomeSkipped, outcomes[0].Kind, "temperature is required, so no row is written")
	require.False(t, outcomes[0].Partial)
	require.Equal(t, ReasonIncompleteRequired, outcomes[0].Reason)
}

func TestExplicitPartialPolicyNoUsableMemberIsNotARow(t *testing.T) {
	a := newFixtureAssembler(t, partialConfig(func(c *Config) { c.Members[0].Required = false }))
	a.Offer(Sample{SampleID: "t5", MemberKey: "temperature", ObservedAt: at(5), Quality: schema.QualityBad, QualityReason: "read-failed"}, at(6))

	outcomes := a.Tick(at(10))
	require.Equal(t, OutcomeSkipped, outcomes[0].Kind)
	require.Equal(t, ReasonNoUsableMember, outcomes[0].Reason)

	silent := a.Tick(at(20))
	require.Equal(t, OutcomeNoData, silent[0].Kind)
}

func TestExplicitPartialPolicyCompleteRowIsNotMarkedPartial(t *testing.T) {
	a := newFixtureAssembler(t, partialConfig(nil))
	offerGood(t, a, "t5", "temperature", 5, 6)
	offerGood(t, a, "p6", "pressure", 6, 7)

	row := a.Tick(at(10))[0]
	require.Equal(t, OutcomeRow, row.Kind)
	require.False(t, row.Partial)
}

func TestExplicitPartialPolicyCompleteGoodRowNeedsNoExternalMetadata(t *testing.T) {
	// Default skip_row with zero storage capabilities: a custom table without
	// NULL or quality columns still gets complete rows.
	a := newFixtureAssembler(t, nil)
	offerGood(t, a, "t5", "temperature", 5, 6)
	offerGood(t, a, "p6", "pressure", 6, 7)

	row := a.Tick(at(10))[0]
	require.Equal(t, OutcomeRow, row.Kind)
	require.False(t, row.Partial)
	for _, member := range row.Members {
		require.True(t, member.Usable())
		require.NotNil(t, member.Sample, "member provenance stays available for the local durable envelope")
		require.NotEmpty(t, member.Sample.SampleID)
		require.Equal(t, schema.QualityGood, member.Sample.Quality)
	}
}
