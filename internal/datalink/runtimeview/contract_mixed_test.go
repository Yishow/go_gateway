package runtimeview

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func mixedContractFixture(t *testing.T) RuleEvidenceView {
	v := committedFixture(t)
	g := &v.Groups[0]
	for _, f := range []*GroupFacts{g.Saved, g.Effective} {
		f.Layout = "managed_wide_profile_v1"
		f.RecordingProfile = nil
		f.IntervalSeconds = nil
		f.Reduction = "member_calculation"
		f.MixedRecording = &MixedRecordingFacts{RecordID: "03760df3-f1e1-4c52-9adf-43785e54be31", ResultShape: "member-calculation-v1"}
		for i := range f.Members {
			f.Members[i].SourceValueType = str("float64")
			f.Members[i].Calculation = &MemberCalculationFacts{Method: "mean", Interval: 15, Unit: "second", Alignment: "utc_epoch", SourceRevision: "source-1", MappingRevision: "mapping-1", ResultType: "float64"}
		}
	}
	current := g.Delivery.Current
	current.Sources = []SourceEvidence{{DeviceID: "D", PointID: "P", TagID: "T", Latest: &SourceResult{State: "not_due", Receipt: current.LastCommit, Calculation: &CalculationResult{ResultShape: "member-calculation-v1", Method: "mean", IntervalSeconds: 15, Alignment: "utc_epoch", DueAt: "2026-10-27T09:59:55Z", NextDue: str("2026-10-27T10:00:00Z"), OutcomeDigest: strings.Repeat("a", 64), OutcomeKind: "row", Status: "not_due", RowEmitted: true, SourceType: "float64", ResultType: "float64"}}}}
	return v
}

func TestMixedContractKeepsNotDueReceiptWithoutFakeWindow(t *testing.T) {
	v := mixedContractFixture(t)
	got, err := parseFixture(t, v)
	require.NoError(t, err)
	require.Equal(t, "effective", got.Groups[0].State)
	r := got.Groups[0].Delivery.Current.Sources[0].Latest
	require.NotNil(t, r.Receipt)
	require.Nil(t, r.WindowStart)
	require.False(t, r.Calculation.ResultPresent)
	require.False(t, r.Calculation.RawPresent)
	require.True(t, r.Calculation.RowEmitted)
	r.Calculation.NextDue = nil // An ended fast member has no future due at the slow member dispatch.
	_, err = parseFixture(t, *got)
	require.NoError(t, err)
	legacy := fixture(t)
	raw, err := json.Marshal(legacy)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "mixed_recording")
	require.NotContains(t, string(raw), "calculation")
	_, err = ParseRuleEvidence(raw, contractScope)
	require.NoError(t, err)
}

func TestMixedContractRejectsContradictoryOrForeignMeaning(t *testing.T) {
	cases := map[string]func(*RuleEvidenceView){
		"fake common interval": func(v *RuleEvidenceView) { x := int64(15); v.Groups[0].Effective.IntervalSeconds = &x },
		"missing intent":       func(v *RuleEvidenceView) { v.Groups[0].Effective.Members[0].Calculation = nil },
		"method changed": func(v *RuleEvidenceView) {
			v.Groups[0].Delivery.Current.Sources[0].Latest.Calculation.Method = "latest"
		},
		"cadence changed": func(v *RuleEvidenceView) {
			v.Groups[0].Delivery.Current.Sources[0].Latest.Calculation.IntervalSeconds = 60
		},
		"notdue raw": func(v *RuleEvidenceView) {
			v.Groups[0].Delivery.Current.Sources[0].Latest.Calculation.RawPresent = true
		},
		"notdue selected": func(v *RuleEvidenceView) {
			v.Groups[0].Delivery.Current.Sources[0].Latest.Calculation.ResultPresent = true
		},
		"skip with receipt": func(v *RuleEvidenceView) {
			c := v.Groups[0].Delivery.Current.Sources[0].Latest.Calculation
			c.RowEmitted = false
			c.OutcomeKind = "skipped"
		},
		"notdue aligned": func(v *RuleEvidenceView) {
			v.Groups[0].Delivery.Current.Sources[0].Latest.Calculation.AlignedStart = str("2026-10-27T10:00:00Z")
		},
		"notdue expired next": func(v *RuleEvidenceView) {
			c := v.Groups[0].Delivery.Current.Sources[0].Latest.Calculation
			c.NextDue = &c.DueAt
		},
		"long notdue": func(v *RuleEvidenceView) { v.Groups[0].Effective.Layout = "managed_long_v1" },
		"invalid digest": func(v *RuleEvidenceView) {
			v.Groups[0].Delivery.Current.Sources[0].Latest.Calculation.OutcomeDigest = "unverified"
		},
		"out of range time": func(v *RuleEvidenceView) {
			v.Groups[0].Delivery.Current.Sources[0].Latest.Calculation.DueAt = "2300-01-01T00:00:00Z"
		},
		"overflow": func(v *RuleEvidenceView) { v.Groups[0].Effective.Members[0].Calculation.Interval = 9223372037 },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			v := mixedContractFixture(t)
			mutate(&v)
			_, err := parseFixture(t, v)
			require.Error(t, err)
		})
	}
}

func mixedDueContractFixture(t *testing.T) RuleEvidenceView {
	v := mixedContractFixture(t)
	r := v.Groups[0].Delivery.Current.Sources[0].Latest
	c := r.Calculation
	r.Receipt = nil
	r.State = "available"
	c.Status = "ok"
	c.ResultPresent = true
	c.RawPresent = true
	c.DueAt = "2026-10-27T10:00:00Z"
	c.NextDue = str("2026-10-27T10:00:15Z")
	c.AlignedStart = str("2026-10-27T09:59:45Z")
	c.AlignedEnd = &c.DueAt
	c.ActualStart = c.AlignedStart
	c.ActualEnd = c.AlignedEnd
	r.WindowStart = c.ActualStart
	r.WindowEnd = c.ActualEnd
	return v
}

func TestMixedContractRejectsProducerContradictions(t *testing.T) {
	for name, mutate := range map[string]func(*SourceResult){
		"missing as available": func(r *SourceResult) {
			r.Calculation.Status = "missing"
			r.Calculation.ResultPresent = false
			r.Calculation.Reason = str("no-samples")
		},
		"unsupported member partial": func(r *SourceResult) {
			r.State = "partial"
			r.Coverage = str("partial")
			r.Calculation.Status = "partial"
			r.Calculation.ResultPresent = false
		},
		"ok without selected result": func(r *SourceResult) { r.State = "unknown"; r.Calculation.ResultPresent = false },
		"due labelled not due":       func(r *SourceResult) { r.State = "not_due" },
		"usable without raw":         func(r *SourceResult) { r.Calculation.RawPresent = false },
		"59 minute window": func(r *SourceResult) {
			c := r.Calculation
			c.AlignedStart = str("2026-10-27T09:00:00Z")
			c.AlignedEnd = str("2026-10-27T09:59:00Z")
			c.ActualStart = c.AlignedStart
			c.ActualEnd = c.AlignedEnd
			r.WindowStart = c.ActualStart
			r.WindowEnd = c.ActualEnd
		},
		"wrong nominal due":   func(r *SourceResult) { r.Calculation.DueAt = "2026-10-27T10:00:15Z"; r.Calculation.NextDue = nil },
		"misaligned next":     func(r *SourceResult) { r.Calculation.NextDue = str("2026-10-27T10:00:01Z") },
		"skipped next window": func(r *SourceResult) { r.Calculation.NextDue = str("2026-10-27T10:00:30Z") },
	} {
		t.Run(name, func(t *testing.T) {
			v := mixedDueContractFixture(t)
			_, err := parseFixture(t, v)
			require.NoError(t, err)
			mutate(v.Groups[0].Delivery.Current.Sources[0].Latest)
			_, err = parseFixture(t, v)
			require.Error(t, err)
		})
	}
}

func TestMixedContractPartialMeanAndUnknownDeltaRemainDistinct(t *testing.T) {
	v := mixedDueContractFixture(t)
	r := v.Groups[0].Delivery.Current.Sources[0].Latest
	r.State = "partial"
	r.Coverage = str("partial")
	r.Calculation.ActualStart = str("2026-10-27T09:59:48Z")
	r.Calculation.ActualEnd = str("2026-10-27T09:59:53Z")
	r.Calculation.NextDue = nil
	r.WindowStart = r.Calculation.ActualStart
	r.WindowEnd = r.Calculation.ActualEnd
	_, err := parseFixture(t, v)
	require.NoError(t, err)
	for _, f := range []*GroupFacts{v.Groups[0].Saved, v.Groups[0].Effective} {
		f.Members[0].Calculation.Method = "delta"
		f.Members[0].Calculation.MeasurementID = str("meter")
	}
	r.State = "unknown"
	r.Coverage = nil
	r.DeltaKnown = new(bool)
	r.Calculation.Method = "delta"
	r.Calculation.Status = "missing"
	r.Calculation.ResultPresent = false
	r.Calculation.Reason = str("gap")
	_, err = parseFixture(t, v)
	require.NoError(t, err, "reliable raw baseline/gap is not a selected delta")
}
