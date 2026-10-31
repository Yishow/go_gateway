package runtimeview

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
)

// Values/bounds/presence come from the existing assembler, not invented
// successful DTOs. The fast ended member remains not_due at slow nominal60.
func TestMixedContractAcceptsActualIndependentCutoffOutcomes(t *testing.T) {
	base := time.Date(2026, 10, 27, 9, 59, 0, 0, time.UTC)
	cfg := snapshot.Config{WorkspaceID: "W", GroupID: "G", GroupRevision: "A", DestinationScope: "db", Mode: snapshot.ModeMixedWide, StartsAt: base.Add(3 * time.Second), Storage: snapshot.StorageCapabilities{NullableValues: true, MemberQuality: true}, Members: []snapshot.Member{
		{Key: "fast", Required: true, Calculation: snapshot.CalculationLatest, Interval: 15 * time.Second, SourceValueType: schema.DataTypeFloat64},
		{Key: "slow", Required: true, Calculation: snapshot.CalculationLatest, Interval: 60 * time.Second, SourceValueType: schema.DataTypeFloat64},
	}}
	a, err := snapshot.NewAssembler(cfg)
	require.NoError(t, err)
	value, err := measurement.NewFloat64(12)
	require.NoError(t, err)
	for _, key := range []string{"fast", "slow"} {
		at := base.Add(20 * time.Second)
		got := a.Offer(snapshot.Sample{MemberKey: key, SampleID: key, ObservedAt: at, Quality: schema.QualityGood, Value: value}, at)
		require.NotEqual(t, snapshot.OfferRejected, got.Outcome)
	}
	require.NoError(t, a.SetMixedEndsAt(base.Add(23*time.Second)))
	plan := a.Plan(base.Add(time.Minute))
	require.NoError(t, plan.Err)
	require.Len(t, plan.Outcomes, 3)
	for _, out := range plan.Outcomes {
		for i, member := range out.Members {
			v := mixedContractFixture(t)
			r := v.Groups[0].Delivery.Current.Sources[0].Latest
			c := r.Calculation
			w := member.Window
			for _, facts := range []*GroupFacts{v.Groups[0].Saved, v.Groups[0].Effective} {
				facts.Members[0].Calculation.Method = "latest"
				facts.Members[0].Calculation.Interval = int64(cfg.Members[i].Interval / time.Second)
			}
			c.Method = "latest"
			c.IntervalSeconds = int64(w.Interval / time.Second)
			c.DueAt = out.DueAt.Format(time.RFC3339Nano)
			c.NextDue = producerTime(w.NextDue)
			c.AlignedStart = producerTime(w.AlignedStart)
			c.AlignedEnd = producerTime(w.AlignedEnd)
			c.ActualStart = producerTime(w.ActualStart)
			c.ActualEnd = producerTime(w.ActualEnd)
			c.Status = string(member.Status)
			c.ResultPresent = member.Usable()
			c.RawPresent = member.Sample != nil
			c.OutcomeKind = string(out.Kind)
			c.RowEmitted = out.Kind == snapshot.OutcomeRow
			r.WindowStart = c.ActualStart
			r.WindowEnd = c.ActualEnd
			r.Receipt = nil
			r.State = "unknown"
			if member.Usable() {
				r.State = "available"
			}
			if member.Status == snapshot.MemberNotDue {
				r.State = "not_due"
			}
			_, err = parseFixture(t, v)
			require.NoError(t, err, "due=%s member=%s", out.DueAt, member.MemberKey)
		}
	}
}
func producerTime(at time.Time) *string {
	if at.IsZero() {
		return nil
	}
	return str(at.UTC().Format(time.RFC3339Nano))
}
