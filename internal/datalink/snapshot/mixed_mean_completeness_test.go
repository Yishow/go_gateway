package snapshot

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"testing"
	"time"
)

func TestMixedUnavailableMeanRetainsCoverageAndReasons(t *testing.T) {
	for _, required := range []bool{false, true} {
		for _, kind := range []string{"missing", "bad", "stale", "wrong type"} {
			t.Run(fmt.Sprintf("required=%t/%s", required, kind), func(t *testing.T) {
				c := mixedConfig(ModeMixedWide)
				c.Members[0].Required = required
				c.Members[0].MaxAge = 2 * time.Second
				a, err := NewAssembler(c)
				require.NoError(t, err)
				reason := ReasonNoSamples
				if kind != "missing" {
					s := mixedSample("a", "observation", 14, 999)
					switch kind {
					case "bad":
						s.Quality = schema.QualityBad
						s.QualityReason = "internal-secret-probe"
						reason = MeanReasonBad
					case "stale":
						s.ObservedAt = mixedTime(10)
						reason = MeanReasonStale
					case "wrong type":
						s.Value = measurement.NewText("999")
						reason = ReasonTypeMismatch
					}
					mixedOffer(t, a, s)
				}
				row, r := mixedDueMember(t, a.Plan(mixedTime(15)), 15, "a")
				require.False(t, r.Usable())
				require.Nil(t, r.Reduction)
				require.NotNil(t, r.Evidence)
				require.Equal(t, CoverageMissing, r.Evidence.Coverage)
				require.Zero(t, r.Evidence.ValidCount)
				require.Equal(t, reason, r.Reason)
				require.Equal(t, reason, r.Evidence.Reason)
				require.NotContains(t, r.Evidence.Reason, "internal-secret-probe")
				if kind == "missing" {
					require.Zero(t, r.Evidence.RejectedCount)
					require.Nil(t, r.Sample)
					require.True(t, r.Evidence.FirstObservedAt.IsZero())
				} else {
					require.Equal(t, int64(1), r.Evidence.RejectedCount)
					require.NotNil(t, r.Sample)
					require.Equal(t, []string{reason}, r.Evidence.RejectionReasons)
				}
				if required {
					require.Equal(t, OutcomeSkipped, row.Kind)
				} else {
					require.Equal(t, OutcomeRow, row.Kind)
					require.True(t, row.Partial)
				}
			})
		}
	}
}

func TestMixedPartialMeanKeepsValidResultWithoutClaimingFullCoverage(t *testing.T) {
	c := mixedConfig(ModeMixedWide)
	c.StartsAt = mixedTime(7)
	a, err := NewAssembler(c)
	require.NoError(t, err)
	mixedOffer(t, a, mixedSample("a", "first", 8, 20))
	bad := mixedSample("a", "bad", 10, 999)
	bad.Quality = schema.QualityBad
	mixedOffer(t, a, bad)
	mixedOffer(t, a, mixedSample("a", "last", 14, 24))
	row, r := mixedDueMember(t, a.Plan(mixedTime(15)), 15, "a")
	require.Equal(t, OutcomeRow, row.Kind)
	require.True(t, r.Usable())
	require.Equal(t, 22.0, r.Reduction.Value.Value())
	require.Equal(t, CoveragePartial, r.Evidence.Coverage)
	require.Equal(t, ReasonBoundaryTruncated, r.Evidence.Reason)
	require.Equal(t, int64(2), r.Evidence.ValidCount)
	require.Equal(t, int64(1), r.Evidence.RejectedCount)
	require.Equal(t, []string{MeanReasonBad}, r.Evidence.RejectionReasons)
	require.Equal(t, mixedTime(8), r.Evidence.FirstObservedAt)
	require.Equal(t, mixedTime(14), r.Evidence.LastObservedAt)
	require.Equal(t, mixedTime(7), r.Window.ActualStart)
	require.Equal(t, mixedTime(0), r.Window.AlignedStart)
}

func TestMixedRealZeroIsNotConfusedWithUnavailable(t *testing.T) {
	c := mixedCounterConfig(ModeMixedWide)
	a, err := NewAssembler(c)
	require.NoError(t, err)
	mixedPolicyBaseline(t, a, 1000)
	mixedOffer(t, a, mixedIntegerSample("a", "negative", 106, -1))
	mixedOffer(t, a, mixedIntegerSample("a", "positive", 114, 1))
	mixedOffer(t, a, mixedIntegerSample("b", "zero", 114, 0))
	mixedOffer(t, a, mixedIntegerSample("c", "unchanged", 119, 1000))
	p := a.Plan(mixedTime(120))
	row, r := mixedDueMember(t, p, 120, "c")
	require.Equal(t, OutcomeRow, row.Kind)
	require.True(t, r.Usable())
	require.Equal(t, "0", mixedPolicyDigits(t, r))
	for _, key := range []string{"a", "b"} {
		_, mean := mixedDueMember(t, p, 120, key)
		require.True(t, mean.Usable())
		digits, ok := mean.Reduction.Value.DecimalDigits()
		require.True(t, ok)
		require.Equal(t, "0", digits)
	}
}

func TestMixedOptionalNullAndSparseWideRequireStorageProof(t *testing.T) {
	for _, mode := range []Mode{ModeMixedWide, ModeMixedLong} {
		for _, quality := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/quality=%t", mode, quality), func(t *testing.T) {
				c := mixedCounterConfig(mode)
				c.Members[2].Required = false
				c.Storage = StorageCapabilities{NullableValues: quality, MemberQuality: !quality}
				_, err := NewAssembler(c)
				require.ErrorIs(t, err, ErrPartialUnsupported)
			})
		}
	}
	c := mixedCounterConfig(ModeMixedWide)
	c.Storage = StorageCapabilities{}
	_, err := NewAssembler(c)
	require.ErrorIs(t, err, ErrPartialUnsupported, "even required sparse columns need not_due NULL support")
}
