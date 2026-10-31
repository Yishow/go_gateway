package snapshot

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"testing"
	"time"
)

func TestMixedSameSourceDifferentGroupsHaveIndependentBaselines(t *testing.T) {
	c := mixedCounterConfig(ModeMixedWide)
	a, err := NewAssembler(c)
	require.NoError(t, err)
	c.GroupID = "record-b"
	c.Members[2].Interval = 600 * time.Second
	bBinding := c.Counters["c"]
	bBinding.GroupID = c.GroupID
	c.Counters["c"] = bBinding
	b, err := NewAssembler(c)
	require.NoError(t, err)
	// Both retain the same source member and measurement identity; only record scope differs.
	require.Equal(t, a.config.Counters["c"].MeasurementID, b.config.Counters["c"].MeasurementID)
	mixedPolicyBaseline(t, a, 1000)
	mixedOffer(t, b, mixedIntegerSample("c", "baseline", 599, 2000))
	p := b.Plan(mixedTime(120))
	require.Empty(t, p.CounterCheckpoints())
	require.Empty(t, b.CounterCheckpoints())
	mixedOffer(t, a, mixedIntegerSample("c", "tail", 119, 1008))
	ap := a.Plan(mixedTime(120))
	require.Equal(t, "8", mixedPolicyDigits(t, mixedPolicyResult(t, ap)))
	a.Commit(ap)
	require.Empty(t, b.CounterCheckpoints())
	b.Commit(b.Plan(mixedTime(600)))
	require.Equal(t, "2000", func() string { d, _ := b.CounterCheckpoints()["c"].Reading.Value.DecimalDigits(); return d }())
	mixedOffer(t, b, mixedIntegerSample("c", "tail", 1199, 2003))
	bp := b.Plan(mixedTime(1200))
	require.Equal(t, "3", mixedPolicyDigits(t, mixedPolicyResult(t, bp)))
	require.Error(t, b.RestoreCounterCheckpoints(a.CounterCheckpoints()))
	require.Equal(t, "baseline", b.CounterCheckpoints()["c"].Reading.SampleID, "foreign restore cannot partially mutate current baseline")
	require.Equal(t, "tail", a.CounterCheckpoints()["c"].Reading.SampleID)
}

func TestMixedCounterBindingRevisionStartsNewBaseline(t *testing.T) {
	old := mixedPolicyAssembler(t, ModeMixedWide, nil)
	mixedPolicyBaseline(t, old, 1000)
	saved := old.CounterCheckpoints()
	progress, _ := old.MixedProgress()
	changes := map[string]func(*aggregation.CounterBinding){
		"measurement":    func(b *aggregation.CounterBinding) { b.MeasurementID = "other" },
		"definition":     func(b *aggregation.CounterBinding) { b.DefinitionRevision = "def-2" },
		"source binding": func(b *aggregation.CounterBinding) { b.SourceBindingRevision = "binding-2" },
		"epoch":          func(b *aggregation.CounterBinding) { b.SeriesEpoch = "epoch-2" },
		"unit":           func(b *aggregation.CounterBinding) { b.Unit = "MWh" },
		"kind":           func(b *aggregation.CounterBinding) { b.Kind = aggregation.CounterKindSigned },
		"rollover":       mixedRollover(true),
		"exact scale":    func(b *aggregation.CounterBinding) { b.ExactRecordingScale = "0.1" },
	}
	for name, mutate := range changes {
		t.Run(name, func(t *testing.T) {
			a := mixedPolicyAssembler(t, ModeMixedWide, mutate)
			require.NoError(t, a.RestoreMixedProgress(progress))
			require.NoError(t, a.RestoreCounterCheckpoints(saved))
			sample := mixedIntegerSample("c", "new-binding", 119, 1008)
			if scale := a.config.Counters["c"].ExactRecordingScale; scale != "" {
				var err error
				sample.Value, err = measurement.ScaleExactInteger(sample.Value, scale)
				require.NoError(t, err)
			}
			mixedOffer(t, a, sample)
			p := a.Plan(mixedTime(120))
			r := mixedPolicyResult(t, p)
			require.False(t, r.Usable())
			require.Equal(t, aggregation.CounterBaseline, r.Counter.Status)
			require.Equal(t, aggregation.ReasonBindingChanged, r.Counter.Reason)
			require.Nil(t, r.Counter.Delta)
			require.Equal(t, "baseline", a.CounterCheckpoints()["c"].Reading.SampleID)
			require.Equal(t, "new-binding", p.CounterCheckpoints()["c"].Reading.SampleID)
		})
	}
	// Rate alone changes the existing rollover identity too, not just adding a policy.
	rollover := mixedPolicyAssembler(t, ModeMixedWide, mixedRollover(true))
	mixedPolicyBaseline(t, rollover, 1000)
	changed := mixedPolicyAssembler(t, ModeMixedWide, func(b *aggregation.CounterBinding) {
		mixedRollover(true)(b)
		rate := measurement.NewUint64(50)
		b.MaxRatePerSecond = &rate
	})
	require.NoError(t, changed.RestoreMixedProgress(progress))
	require.NoError(t, changed.RestoreCounterCheckpoints(rollover.CounterCheckpoints()))
	mixedOffer(t, changed, mixedIntegerSample("c", "rate-change", 119, 1008))
	r := mixedPolicyResult(t, changed.Plan(mixedTime(120)))
	require.Equal(t, aggregation.ReasonBindingChanged, r.Counter.Reason)
	require.Nil(t, r.Counter.Delta)
}

func TestMixedMethodOrRecordRevisionCannotBorrowCounterCheckpoint(t *testing.T) {
	old := mixedPolicyAssembler(t, ModeMixedWide, nil)
	mixedPolicyBaseline(t, old, 1000)
	saved := old.CounterCheckpoints()
	for _, method := range []Calculation{CalculationLatest, CalculationMean} {
		t.Run(string(method), func(t *testing.T) {
			c := mixedCounterConfig(ModeMixedWide)
			c.Members[2].Calculation = method
			delete(c.Counters, "c")
			a, err := NewAssembler(c)
			require.NoError(t, err)
			require.Error(t, a.RestoreCounterCheckpoints(saved))
			mixedOffer(t, a, mixedIntegerSample("c", "raw", 59, 1000))
			p := a.Plan(mixedTime(60))
			require.Empty(t, p.CounterCheckpoints())
			a.Commit(p)
			require.Empty(t, a.CounterCheckpoints())
		})
	}
	for _, field := range []string{"workspace", "group", "profile revision", "member"} {
		t.Run(field, func(t *testing.T) {
			bad := saved["c"]
			switch field {
			case "workspace":
				bad.Binding.WorkspaceID = "other"
			case "group":
				bad.Binding.GroupID = "copy"
			case "profile revision":
				bad.Binding.ProfileRevision = "rev-old"
			case "member":
				bad.Binding.MemberKey = "another"
			}
			fresh := mixedPolicyAssembler(t, ModeMixedWide, nil)
			require.Error(t, fresh.RestoreCounterCheckpoints(map[string]aggregation.CounterCheckpoint{"c": bad}))
			require.Empty(t, fresh.CounterCheckpoints())
		})
	}
	c := mixedCounterConfig(ModeMixedWide)
	c.GroupRevision = "new-method-revision"
	binding := c.Counters["c"]
	binding.ProfileRevision = c.GroupRevision
	c.Counters["c"] = binding
	fresh, err := NewAssembler(c)
	require.NoError(t, err)
	require.Error(t, fresh.RestoreCounterCheckpoints(saved))
	mixedOffer(t, fresh, mixedIntegerSample("c", "new-revision-first", 59, 1008))
	r := mixedPolicyResult(t, fresh.Plan(mixedTime(60)))
	require.Equal(t, aggregation.CounterBaseline, r.Counter.Status)
	require.Nil(t, r.Counter.Delta)
}
