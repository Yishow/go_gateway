package snapshot

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func wideProfileConfig(mutate func(*Config)) Config {
	config := fixtureConfig()
	config.Mode = ModeWideProfile
	config.Storage = StorageCapabilities{NullableValues: true, MemberQuality: true}
	config.Members = []Member{{Key: "alpha", SampleReduction: SampleReductionMean}, {Key: "beta"}}
	if mutate != nil {
		mutate(&config)
	}
	return config
}

func newWideProfileAssembler(t *testing.T, mutate func(*Config)) *Assembler {
	t.Helper()
	assembler, err := NewAssembler(wideProfileConfig(mutate))
	require.NoError(t, err)
	return assembler
}

func closeWideWindow(t *testing.T, a *Assembler, now int) Outcome {
	t.Helper()
	outcomes := a.Tick(at(now))
	require.Len(t, outcomes, 1, "a wide profile closes one record per window, not one per member")
	return outcomes[0]
}

// TestWideProfileWindowKeepsEachMembersReduction pins that one record carries
// every member under its own rule: the numeric gauge averages, the other member
// keeps its newest eligible reading, and neither changes the other.
func TestWideProfileWindowKeepsEachMembersReduction(t *testing.T) {
	a := newWideProfileAssembler(t, nil)
	accept(t, a, intSample("a1", 1, 1), intSample("a2", 2, 2), goodSample("b1", "beta", 1, 10), goodSample("b2", "beta", 5, 30))
	outcome := closeWideWindow(t, a, 10)

	require.Equal(t, OutcomeRow, outcome.Kind)
	require.False(t, outcome.Partial)
	require.Len(t, outcome.Members, 2)
	alpha, beta := outcome.Members[0], outcome.Members[1]
	require.Equal(t, "alpha", alpha.MemberKey)
	require.Equal(t, "1.5", digitsOfExact(t, exactMean(t, alpha).Value))
	require.NotNil(t, alpha.Evidence)
	require.Equal(t, CoverageUnknown, alpha.Evidence.Coverage, "a usable mean without whole-window proof is unknown, never complete")
	require.Equal(t, int64(30), beta.Sample.Value.Value(), "a latest member keeps the newest reading, not the mean 20")
	require.Nil(t, beta.Reduction)
	require.NotEmpty(t, outcome.RecordID)
	require.NotEmpty(t, outcome.EffectKey)
}

// TestWideProfileAllMissingWindowIsOneNULLRow pins that a window where nothing
// was observed is still exactly one record of missing members, never no_data, a
// zero or an older reading.
func TestWideProfileAllMissingWindowIsOneNULLRow(t *testing.T) {
	a := newWideProfileAssembler(t, nil)
	outcome := closeWideWindow(t, a, 10)
	require.Equal(t, OutcomeRow, outcome.Kind)
	require.True(t, outcome.Partial)
	for _, member := range outcome.Members {
		require.Equal(t, [2]string{string(MemberMissing), ReasonNoSamples}, [2]string{string(member.Status), member.Reason}, member.MemberKey)
		require.Nil(t, member.Sample)
		require.Nil(t, member.Reduction)
	}
	require.Equal(t, CoverageMissing, outcome.Members[0].Evidence.Coverage)
	require.NotEmpty(t, outcome.EffectKey, "the record still has its identity")

	// The next window is a new record, and nothing is carried over from this one.
	accept(t, a, intSample("a1", 12, 5))
	next := closeWideWindow(t, a, 20)
	require.Equal(t, MemberMissing, next.Members[1].Status)
}

// TestWideProfileRequiredMemberGatesTheRow pins that only an explicit required
// member that has no usable value skips the whole row, and says so; an optional
// member in the same state is NULL.
func TestWideProfileRequiredMemberGatesTheRow(t *testing.T) {
	required := func(c *Config) { c.Members[0].Required = true }

	skipped := closeWideWindow(t, newWideProfileAssembler(t, required), 10)
	require.Equal(t, [2]string{string(OutcomeSkipped), ReasonIncompleteRequired}, [2]string{string(skipped.Kind), skipped.Reason})

	bad := newWideProfileAssembler(t, required)
	reading := intSample("a1", 1, 1)
	reading.Quality = schema.QualityBad
	accept(t, bad, reading)
	require.Equal(t, OutcomeSkipped, closeWideWindow(t, bad, 10).Kind, "a bad reading is not a usable required value")

	served := newWideProfileAssembler(t, required)
	accept(t, served, intSample("a1", 1, 4), intSample("a2", 2, 5))
	row := closeWideWindow(t, served, 10)
	require.Equal(t, OutcomeRow, row.Kind)
	require.True(t, row.Partial, "the optional member is still NULL")
	require.Equal(t, MemberMissing, row.Members[1].Status)
}

// TestWideProfileWindowBoundsAndCutover pins the fixed half-open UTC window and
// the real bounds of a window a forward cutover cut short: an observation before
// the start belongs to the previous owner, and the cut window is named, not
// passed off as complete.
func TestWideProfileWindowBoundsAndCutover(t *testing.T) {
	a := newWideProfileAssembler(t, func(c *Config) { c.StartsAt = at(7) })
	accept(t, a, intSample("before", 5, 100), intSample("a1", 8, 2), intSample("a2", 9, 4))
	outcome := closeWideWindow(t, a, 10)
	require.Equal(t, [3]any{"start", at(7), at(10)}, [3]any{outcome.Truncated, outcome.WindowStart, outcome.WindowEnd})
	alpha := outcome.Members[0]
	require.Equal(t, "3", digitsOfExact(t, exactMean(t, alpha).Value), "the observation before the cut is not part of the mean")
	require.Equal(t, [2]any{CoveragePartial, ReasonBoundaryTruncated}, [2]any{alpha.Evidence.Coverage, alpha.Evidence.Reason})
}

// TestWideProfileConfigFailsClosed pins what this mode needs: storage that holds
// NULL and member provenance, and members that are not entity-keyed.
func TestWideProfileConfigFailsClosed(t *testing.T) {
	_, err := NewAssembler(wideProfileConfig(func(c *Config) { c.Storage = StorageCapabilities{} }))
	require.ErrorIs(t, err, ErrPartialUnsupported)
	_, err = NewAssembler(wideProfileConfig(func(c *Config) { c.Storage = StorageCapabilities{NullableValues: true} }))
	require.ErrorIs(t, err, ErrPartialUnsupported, "nullable alone is not enough without member provenance")
	_, err = NewAssembler(wideProfileConfig(func(c *Config) { c.Members[1].EntityKey = "plant-A" }))
	require.ErrorIs(t, err, ErrInvalidConfig)
	_, err = NewAssembler(wideProfileConfig(func(c *Config) { c.Mode = ModeInstantaneousProfile; c.Members[1].EntityKey = "" }))
	require.NoError(t, err, "the per-member profile mode is unchanged")
}
