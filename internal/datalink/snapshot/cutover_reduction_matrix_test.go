package snapshot

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func reductionCutoverConfig(mode Mode, reduction SampleReduction) Config {
	c := fixtureConfig()
	c.Mode, c.Interval = mode, 15*time.Second
	c.Storage = StorageCapabilities{NullableValues: true, MemberQuality: true}
	c.Members = []Member{{Key: "alpha", Required: true, SampleReduction: reduction, SourceValueType: schema.DataTypeInt64}}
	return c
}

func reductionCutoverOffer(t *testing.T, a *Assembler, samples ...Sample) {
	t.Helper()
	for _, sample := range samples {
		require.Contains(t, []OfferOutcome{OfferSelected, OfferNotSelected}, a.Offer(sample, at(12)).Outcome, sample.SampleID)
	}
}

func reductionCutoverOutcome(t *testing.T, a *Assembler) Outcome {
	t.Helper()
	plan := a.Plan(at(15))
	require.Len(t, plan.Outcomes, 1)
	require.Equal(t, at(0), a.NextClose(), "planning a reduction does not consume accepted inputs")
	return plan.Outcomes[0]
}

func TestCutoverReductionMeanFreshnessUsesActualEnd(t *testing.T) {
	for _, mode := range []Mode{ModeInstantaneousProfile, ModeWideProfile} {
		t.Run(string(mode), func(t *testing.T) {
			c := reductionCutoverConfig(mode, SampleReductionMean)
			c.Members[0].MaxAge = 2 * time.Second
			a, err := NewAssembler(c)
			require.NoError(t, err)
			reductionCutoverOffer(t, a, intSample("successor", 11, 100), intSample("latest", 9, 20), intSample("at-max-age", 8, 10), intSample("stale", 7, 1000))
			a.SetEndsAt(at(10))
			result := memberResult(t, reductionCutoverOutcome(t, a), "alpha")
			require.Equal(t, MemberOK, result.Status, "age must use actual end 10, not nominal end 15")
			require.Equal(t, "15", digitsOfExact(t, exactMean(t, result).Value))
			require.Equal(t, "30", digitsOfExact(t, result.Reduction.Sum))
			require.Equal(t, int64(2), result.Reduction.Count)
			require.Equal(t, "latest", result.Sample.SampleID)
			require.Equal(t, int64(1), result.Evidence.RejectedCount)
			require.Equal(t, []string{MeanReasonStale}, result.Evidence.RejectionReasons)
			require.Equal(t, CoveragePartial, result.Evidence.Coverage)
			require.Equal(t, ReasonBoundaryTruncated, result.Evidence.Reason)
		})
	}
}

// The spec's 10@8, 20@9, 100@11 fixture is shared by every numeric profile.
// State uses the same observation identities and timestamps with frozen bools.
func TestCutoverReductionOwnerPartitionMatrix(t *testing.T) {
	for _, tc := range []struct {
		name      string
		mode      Mode
		reduction SampleReduction
		source    schema.DataType
	}{
		{"legacy-wide", ModeWide, SampleReductionLatest, schema.DataTypeInt64},
		{"custom-wide-partial", ModeWide, SampleReductionLatest, schema.DataTypeInt64},
		{"long-latest", ModeInstantaneousProfile, SampleReductionLatest, schema.DataTypeInt64},
		{"wide-latest", ModeWideProfile, SampleReductionLatest, schema.DataTypeInt64},
		{"long-mean", ModeInstantaneousProfile, SampleReductionMean, schema.DataTypeInt64},
		{"wide-mean", ModeWideProfile, SampleReductionMean, schema.DataTypeInt64},
		{"long-state", ModeInstantaneousProfile, SampleReductionLatest, schema.DataTypeBool},
		{"wide-state", ModeWideProfile, SampleReductionLatest, schema.DataTypeBool},
		{"long-counter-raw", ModeInstantaneousProfile, SampleReductionLatest, schema.DataTypeUint64},
		{"wide-counter-raw", ModeWideProfile, SampleReductionLatest, schema.DataTypeUint64},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, exactBoundary := range []bool{false, true} {
				for _, reverse := range []bool{false, true} {
					t.Run(fmt.Sprintf("exact=%t/reverse=%t", exactBoundary, reverse), func(t *testing.T) {
						c := reductionCutoverConfig(tc.mode, tc.reduction)
						c.Members[0].SourceValueType = tc.source
						if tc.name == "custom-wide-partial" {
							c.IncompletePolicy, c.Members[0].Required = IncompletePartial, false
						}
						old, err := NewAssembler(c)
						require.NoError(t, err)
						c.GroupRevision, c.StartsAt, c.Until = "rev-successor", at(10), at(15)
						next, err := NewAssembler(c)
						require.NoError(t, err)
						samples := []Sample{intSample("pre8", 8, 10), intSample("pre9", 9, 20), intSample("post11", 11, 100)}
						if exactBoundary {
							samples = append(samples, intSample("exact10", 10, 80))
						}
						for i := range samples {
							switch tc.source {
							case schema.DataTypeBool:
								samples[i].Value = measurement.NewBool(samples[i].SampleID == "pre9")
							case schema.DataTypeUint64:
								samples[i].Value = measurement.NewUint64(uint64(samples[i].Value.Value().(int64)))
							}
							samples[i].ReceivedAt, samples[i].SourceRevision, samples[i].MappingRevision = at(12), "source-frozen", "mapping-frozen"
						}
						if reverse {
							slices.Reverse(samples)
						}
						for _, a := range []*Assembler{old, next} {
							reductionCutoverOffer(t, a, samples...)
							require.Equal(t, OfferDuplicate, a.Offer(samples[0], at(12)).Outcome)
						}
						old.SetEndsAt(at(10))
						before, after := reductionCutoverOutcome(t, old), reductionCutoverOutcome(t, next)
						require.Equal(t, []time.Time{at(0), at(10), at(10), at(15)}, []time.Time{before.WindowStart, before.WindowEnd, after.WindowStart, after.WindowEnd})
						require.Equal(t, []string{"end", "start"}, []string{before.Truncated, after.Truncated})
						require.NotEqual(t, before.EffectKey, after.EffectKey)
						oldMember, newMember := memberResult(t, before, "alpha"), memberResult(t, after, "alpha")
						for _, expectedID := range []string{"pre9", "post11"} {
							result := oldMember
							if expectedID == "post11" {
								result = newMember
							}
							require.True(t, result.Usable())
							for _, sample := range samples {
								if sample.SampleID == expectedID {
									require.Equal(t, sample, *result.Sample, "source value and frozen provenance must not be rewritten")
								}
							}
						}
						if tc.reduction == SampleReductionMean {
							require.Equal(t, "15", digitsOfExact(t, exactMean(t, oldMember).Value))
							require.Equal(t, "30", digitsOfExact(t, oldMember.Reduction.Sum))
							require.Equal(t, int64(2), oldMember.Reduction.Count)
							wantValue, wantSum, wantCount := "100", "100", int64(1)
							if exactBoundary {
								wantValue, wantSum, wantCount = "90", "180", 2
							}
							require.Equal(t, wantValue, digitsOfExact(t, exactMean(t, newMember).Value))
							require.Equal(t, wantSum, digitsOfExact(t, newMember.Reduction.Sum))
							require.Equal(t, wantCount, newMember.Reduction.Count)
							for _, result := range []MemberResult{oldMember, newMember} {
								require.Equal(t, CoveragePartial, result.Evidence.Coverage)
								require.Equal(t, ReasonBoundaryTruncated, result.Evidence.Reason)
								require.Zero(t, result.Evidence.RejectedCount)
							}
						} else {
							require.Nil(t, oldMember.Reduction)
							require.Nil(t, newMember.Reduction)
						}
						entity := ""
						if tc.mode.perMember() {
							entity = "alpha"
						}
						oldSet := old.buckets[at(0)].selectors[entity].Accepted("alpha")
						newSet := next.selectorWithinActualBounds(at(0), entity, next.buckets[at(0)].selectors[entity]).Accepted("alpha")
						require.Len(t, oldSet, 2)
						require.Equal(t, []string{"pre8", "pre9"}, []string{oldSet[0].SampleID, oldSet[1].SampleID})
						union := append(append([]Sample(nil), oldSet...), newSet...)
						require.ElementsMatch(t, samples, union, "the owner partition neither drops nor double-counts accepted inputs")
					})
				}
			}
		})
	}
}

func TestCutoverReductionKeepsMissingAndBadQualityPolicies(t *testing.T) {
	for _, mode := range []Mode{ModeWide, ModeInstantaneousProfile, ModeWideProfile} {
		for _, reduction := range []SampleReduction{SampleReductionLatest, SampleReductionMean} {
			if mode == ModeWide && reduction == SampleReductionMean {
				continue // Legacy/custom has no saved mean policy; do not open a new matrix.
			}
			for _, required := range []bool{false, true} {
				for _, fault := range []string{"missing", "bad", "null", "bad-after-good", "null-after-good"} {
					t.Run(fmt.Sprintf("%s/%s/required=%t/%s", mode, reduction, required, fault), func(t *testing.T) {
						c := reductionCutoverConfig(mode, reduction)
						c.Members[0].Required = required
						if mode == ModeWide && !required {
							c.IncompletePolicy = IncompletePartial
						}
						if !mode.perMember() {
							c.Members = append(c.Members, Member{Key: "anchor", Required: true})
						}
						a, err := NewAssembler(c)
						require.NoError(t, err)
						reductionCutoverOffer(t, a, intSample("successor-good", 11, 100))
						if !mode.perMember() {
							reductionCutoverOffer(t, a, goodSample("anchor", "anchor", 9, 1))
						}
						if strings.HasSuffix(fault, "after-good") {
							reductionCutoverOffer(t, a, intSample("earlier-good", 8, 10))
						}
						if fault != "missing" {
							bad := intSample("unusable", 9, 999)
							if strings.HasPrefix(fault, "bad") {
								bad.Quality, bad.QualityReason = schema.QualityBad, "private-read-detail"
							} else {
								bad.Value = measurement.ExactValue{}
							}
							reductionCutoverOffer(t, a, bad)
						}
						a.SetEndsAt(at(10))
						outcome := reductionCutoverOutcome(t, a)
						result := memberResult(t, outcome, "alpha")
						wantStatus := MemberMissing
						if strings.HasPrefix(fault, "bad") {
							wantStatus = MemberBad
						} else if strings.HasPrefix(fault, "null") {
							wantStatus = MemberInvalid
						}
						if mode != ModeWide && strings.HasSuffix(fault, "after-good") {
							wantStatus = MemberOK
						}
						require.Equal(t, wantStatus, result.Status)
						require.Equal(t, wantStatus == MemberOK, result.Usable(), "unusable values remain NULL, never zero or successor values")
						wantKind := OutcomeRow
						if !mode.perMember() && required && wantStatus != MemberOK {
							wantKind = OutcomeSkipped
						}
						require.Equal(t, wantKind, outcome.Kind)
						require.Equal(t, wantKind == OutcomeRow && !mode.perMember() && wantStatus != MemberOK, outcome.Partial)
						if wantKind == OutcomeSkipped {
							require.Equal(t, ReasonIncompleteRequired, outcome.Reason)
						}
						if fault == "missing" {
							require.Nil(t, result.Sample)
							require.Equal(t, ReasonNoSamples, result.Reason)
						} else if wantStatus == MemberOK {
							require.Equal(t, "earlier-good", result.Sample.SampleID)
						} else {
							require.Equal(t, "unusable", result.Sample.SampleID)
						}
						if reduction == SampleReductionMean {
							if wantStatus == MemberOK {
								require.Equal(t, "10", digitsOfExact(t, exactMean(t, result).Value))
								require.Equal(t, int64(1), result.Reduction.Count)
								require.Equal(t, int64(1), result.Evidence.RejectedCount)
								require.Equal(t, CoveragePartial, result.Evidence.Coverage)
							} else {
								require.Nil(t, result.Reduction)
								require.Equal(t, CoverageMissing, result.Evidence.Coverage)
							}
							require.NotContains(t, result.Evidence.RejectionReasons, "private-read-detail")
							require.NotEqual(t, "private-read-detail", result.Evidence.Reason)
						}
					})
				}
			}
		}
	}
}

func TestCutoverReductionMeanRetainsDistinctAcceptanceOrigins(t *testing.T) {
	for _, mode := range []Mode{ModeInstantaneousProfile, ModeWideProfile} {
		t.Run(string(mode), func(t *testing.T) {
			a, err := NewAssembler(reductionCutoverConfig(mode, SampleReductionMean))
			require.NoError(t, err)
			origin, native := intSample("same-id", 8, 10), intSample("same-id", 9, 20)
			origin.AcceptanceOrigin = AcceptanceOrigin{"workspace-A", "original-owner", "rev-original"}
			reductionCutoverOffer(t, a, intSample("successor", 11, 100), native, origin)
			require.Equal(t, OfferDuplicate, a.Offer(origin, at(12)).Outcome)
			conflict := native
			conflict.Value = measurement.NewInt64(999)
			require.Equal(t, OfferConflict, a.Offer(conflict, at(12)).Outcome)
			a.SetEndsAt(at(10))
			result := memberResult(t, reductionCutoverOutcome(t, a), "alpha")
			require.Equal(t, "15", digitsOfExact(t, exactMean(t, result).Value))
			require.Equal(t, "30", digitsOfExact(t, result.Reduction.Sum))
			require.Equal(t, int64(2), result.Reduction.Count)
			require.Equal(t, int64(2), result.Evidence.ValidCount)
			require.Equal(t, native, *result.Sample)
		})
	}
}
