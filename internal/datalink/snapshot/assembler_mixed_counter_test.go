package snapshot

import (
	"fmt"
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func mixedCounterConfig(mode Mode) Config {
	c := mixedConfig(mode)
	c.Members = []Member{
		{Key: "a", Required: true, Calculation: CalculationMean, Interval: 15 * time.Second, SourceValueType: schema.DataTypeInt64},
		{Key: "b", Required: true, Calculation: CalculationMean, Interval: 15 * time.Second, SourceValueType: schema.DataTypeInt64},
		{Key: "c", Required: true, Calculation: CalculationDelta, Interval: 60 * time.Second, SourceValueType: schema.DataTypeUint64},
	}
	binding := counterBinding("c")
	binding.WorkspaceID = c.WorkspaceID
	binding.GroupID = c.GroupID
	binding.ProfileRevision = c.GroupRevision
	c.Counters = map[string]aggregation.CounterBinding{"c": binding}
	return c
}
func mixedIntegerSample(member, id string, sec int, value int64) Sample {
	s := mixedSample(member, id, sec, 0)
	s.Value = measurement.NewInt64(value)
	if member == "c" {
		s.Value = measurement.NewUint64(uint64(value))
	}
	return s
}

func TestMixedMeanMeanDeltaReuseExistingReducers(t *testing.T) {
	c := mixedCounterConfig(ModeMixedWide)
	a, err := NewAssembler(c)
	require.NoError(t, err)
	mixedOffer(t, a, mixedIntegerSample("c", "baseline", 59, 1000))
	p := a.Plan(mixedTime(60))
	require.NoError(t, p.Err)
	require.Len(t, p.Outcomes, 4)
	require.Empty(t, a.CounterCheckpoints())
	require.Len(t, p.CounterCheckpoints(), 1)
	first := p.Outcomes[3]
	require.Equal(t, OutcomeSkipped, first.Kind)
	require.Equal(t, aggregation.ReasonFirstReading, first.Members[2].Reason)
	require.False(t, first.Members[2].Usable())
	require.Equal(t, aggregation.CounterBaseline, first.Members[2].Counter.Status)
	a.Commit(p)
	for _, s := range []Sample{
		mixedIntegerSample("a", "a1", 106, 20), mixedIntegerSample("a", "a2", 110, 22), mixedIntegerSample("a", "a3", 114, 24),
		mixedIntegerSample("b", "b1", 106, 10), mixedIntegerSample("b", "b2", 114, 14), mixedIntegerSample("c", "c1", 119, 1008),
	} {
		mixedOffer(t, a, s)
	}
	p = a.Plan(mixedTime(120))
	require.NoError(t, p.Err)
	require.Len(t, p.CounterCheckpoints(), 1)
	row := p.Outcomes[3]
	require.Equal(t, OutcomeRow, row.Kind)
	require.Len(t, row.Members, 3)
	for i, want := range []string{"22", "12"} {
		require.NotNil(t, row.Members[i].Reduction)
		digits, ok := row.Members[i].Reduction.Value.DecimalDigits()
		require.True(t, ok)
		require.Equal(t, want, digits)
		require.Nil(t, row.Members[i].Counter)
	}
	counter := row.Members[2]
	require.True(t, counter.Usable())
	digits, ok := counter.Counter.Delta.Value.DecimalDigits()
	require.True(t, ok)
	require.Equal(t, "8", digits)
	require.Nil(t, counter.Reduction)
	require.Equal(t, mixedTime(105), row.Members[0].Window.AlignedStart)
	require.Equal(t, mixedTime(60), counter.Window.AlignedStart)
	require.Equal(t, "baseline", a.CounterCheckpoints()["c"].Reading.SampleID, "planning must not commit baseline")
	require.Equal(t, p.Outcomes, a.Plan(mixedTime(120)).Outcomes)
	a.Commit(p)
	require.Len(t, a.CounterCheckpoints(), 1)
	require.Equal(t, "c1", a.CounterCheckpoints()["c"].Reading.SampleID)
}

func TestMixedNotDueNeverAdvancesCounterAndRestartKeepsPhase(t *testing.T) {
	c := mixedCounterConfig(ModeMixedWide)
	a, err := NewAssembler(c)
	require.NoError(t, err)
	mixedOffer(t, a, mixedIntegerSample("c", "baseline", 10, 1000))
	p := a.Plan(mixedTime(15))
	require.Empty(t, p.CounterCheckpoints())
	a.Commit(p)
	require.Empty(t, a.CounterCheckpoints())
	mixedOffer(t, a, mixedIntegerSample("c", "end", 59, 1008))
	p = a.Plan(mixedTime(60))
	a.Commit(p)
	require.Equal(t, "end", a.CounterCheckpoints()["c"].Reading.SampleID)
	progress, _ := a.MixedProgress()
	b, err := NewAssembler(c)
	require.NoError(t, err)
	require.NoError(t, b.RestoreMixedProgress(progress))
	require.NoError(t, b.RestoreCounterCheckpoints(a.CounterCheckpoints()))
	mixedOffer(t, b, mixedIntegerSample("c", "next", 119, 1016))
	p = b.Plan(mixedTime(120))
	require.Equal(t, "8", func() string { d, _ := p.Outcomes[3].Members[2].Counter.Delta.Value.DecimalDigits(); return d }())
	bad := a.CounterCheckpoints()
	checkpoint := bad["c"]
	checkpoint.Binding.MemberKey = "a"
	require.Error(t, b.RestoreCounterCheckpoints(map[string]aggregation.CounterCheckpoint{"a": checkpoint}))
}

func TestMixedLongBaselineFailureIsDueNotNotDue(t *testing.T) {
	c := mixedCounterConfig(ModeMixedLong)
	a, err := NewAssembler(c)
	require.NoError(t, err)
	mixedOffer(t, a, mixedIntegerSample("c", "tail", 59, 1000))
	mixedOffer(t, a, mixedIntegerSample("c", "exact-boundary", 60, 1008))
	p := a.Plan(mixedTime(60))
	require.NoError(t, p.Err)
	require.Len(t, p.Outcomes, 9)
	row := p.Outcomes[8]
	require.Equal(t, "c", row.EntityKey)
	require.Equal(t, OutcomeSkipped, row.Kind)
	require.Equal(t, aggregation.ReasonFirstReading, row.Members[0].Reason)
	require.Equal(t, "tail", p.CounterCheckpoints()["c"].Reading.SampleID)
	require.Equal(t, mixedTime(0), row.Members[0].Window.AlignedStart)
	require.Equal(t, mixedTime(60), row.DueAt)
	a.Commit(p)
	p = a.Plan(mixedTime(120))
	row = p.Outcomes[8]
	require.Equal(t, OutcomeRow, row.Kind)
	require.Equal(t, "exact-boundary", row.Members[0].Sample.SampleID)
	require.Equal(t, "8", func() string { d, _ := row.Members[0].Counter.Delta.Value.DecimalDigits(); return d }())
}

func TestMixedFrozenCounterTypeCannotPoisonBaselineOrHideWrongMiddle(t *testing.T) {
	for _, mode := range []Mode{ModeMixedWide, ModeMixedLong} {
		for _, haveBaseline := range []bool{false, true} {
			for _, wrongTail := range []bool{false, true} {
				name := fmt.Sprintf("%s/baseline=%t/wrong-tail=%t", mode, haveBaseline, wrongTail)
				t.Run(name, func(t *testing.T) {
					c := mixedCounterConfig(mode)
					a, err := NewAssembler(c)
					require.NoError(t, err)
					offset := 0
					if haveBaseline {
						mixedOffer(t, a, mixedIntegerSample("c", "previous", 59, 1000))
						a.Commit(a.Plan(mixedTime(60)))
						offset = 60
					}
					good := mixedIntegerSample("c", "good", offset+58, 1008)
					wrong := mixedIntegerSample("c", "wrong", offset+59, 1009)
					wrong.Value = measurement.NewInt64(1009)
					if !wrongTail {
						wrong.ObservedAt = mixedTime(offset + 58)
						good.ObservedAt = mixedTime(offset + 59)
					}
					mixedOffer(t, a, good)
					mixedOffer(t, a, wrong)
					p := a.Plan(mixedTime(offset + 60))
					require.NoError(t, p.Err)
					row := p.Outcomes[len(p.Outcomes)-1]
					r := row.Members[len(row.Members)-1]
					require.False(t, r.Usable(), "wrong-type sequence must never become normal usage by dropping the bad sample")
					require.Equal(t, ReasonTypeMismatch, r.Reason)
					require.NotNil(t, r.Counter)
					require.Nil(t, r.Counter.Delta)
					require.Equal(t, aggregation.CounterUncertain, r.Counter.Status)
					next, advanced := p.CounterCheckpoints()["c"]
					if wrongTail {
						require.False(t, advanced, "wrong frozen-type tail must never advance baseline")
					} else {
						require.True(t, advanced, "existing engine can retain its reliable correct-type tail")
						require.Equal(t, "good", next.Reading.SampleID)
						require.Equal(t, measurement.ExactUint64, next.Reading.Value.Type())
					}
					a.Commit(p)
					if wrongTail && haveBaseline {
						require.Equal(t, "previous", a.CounterCheckpoints()["c"].Reading.SampleID)
					}
					if wrongTail && !haveBaseline {
						require.Empty(t, a.CounterCheckpoints())
					}
				})
			}
		}
	}
}

func TestMixedCounterRestoreRejectsWrongFrozenTypeAndSourceRevision(t *testing.T) {
	c := mixedCounterConfig(ModeMixedWide)
	c.Members[2].SourceRevision = "source-1"
	c.Members[2].MappingRevision = "mapping-1"
	a, err := NewAssembler(c)
	require.NoError(t, err)
	s := mixedIntegerSample("c", "foreign", 59, 1000)
	require.Equal(t, ReasonRevisionMismatch, a.Offer(s, s.ObservedAt).Reason)
	require.Zero(t, a.OpenBuckets())
	s.SampleID = "local"
	s.SourceRevision = "source-1"
	s.MappingRevision = "mapping-1"
	mixedOffer(t, a, s)
	a.Commit(a.Plan(mixedTime(60)))
	checkpoints := a.CounterCheckpoints()
	bad := checkpoints["c"]
	bad.Reading.Value = measurement.NewInt64(1000)
	checkpoints["c"] = bad
	fresh, err := NewAssembler(c)
	require.NoError(t, err)
	require.Error(t, fresh.RestoreCounterCheckpoints(checkpoints))
	require.Empty(t, fresh.CounterCheckpoints())
	require.NoError(t, fresh.RestoreCounterCheckpoints(a.CounterCheckpoints()))
	require.Equal(t, measurement.ExactUint64, fresh.CounterCheckpoints()["c"].Reading.Value.Type())
}

func TestMixedWrongNonnumericCounterKeepsEngineInvalidity(t *testing.T) {
	a, err := NewAssembler(mixedCounterConfig(ModeMixedWide))
	require.NoError(t, err)
	wrong := mixedIntegerSample("c", "wrong-text", 58, 0)
	wrong.Value = measurement.NewText("bad")
	mixedOffer(t, a, wrong)
	mixedOffer(t, a, mixedIntegerSample("c", "correct-tail", 59, 1000))
	p := a.Plan(mixedTime(60))
	r := p.Outcomes[3].Members[2]
	require.Equal(t, MemberInvalid, r.Status)
	require.Equal(t, ReasonTypeMismatch, r.Reason)
	require.Equal(t, aggregation.CounterInvalid, r.Counter.Status)
	require.Equal(t, aggregation.ReasonNonNumericReading, r.Counter.Reason)
	require.Nil(t, r.Counter.Delta)
	require.Empty(t, p.CounterCheckpoints(), "adapter must not override engine refusal to advance nonnumeric sequence")
}
