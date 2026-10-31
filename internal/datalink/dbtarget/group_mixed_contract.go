package dbtarget

import (
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"time"
)

const issueMixedContract = "mixed-contract-invalid"
const mixedTruncationStart = "start"

type mixedIntent struct {
	key      string
	method   snapshot.Calculation
	interval time.Duration
	source   schema.DataType
	result   measurement.ExactType
	required bool
	counter  *aggregation.CounterBinding
}

func (m GroupRowMember) mixedIntent() mixedIntent {
	return mixedIntent{m.MemberKey, m.Calculation, m.Interval, m.SourceValueType, m.Type, m.Required, m.Counter}
}
func (m LongRowMember) mixedIntent() mixedIntent {
	return mixedIntent{m.MemberKey, m.Calculation, m.Interval, m.SourceValueType, m.Type, m.Required, m.Counter}
}
func validMixedIntent(m mixedIntent, workspace, group, revision string) bool {
	if m.key == "" || m.interval < time.Second || m.interval%time.Second != 0 {
		return false
	}
	source, ok := measurement.ExactTypeForTag(m.source)
	if !ok {
		return false
	}
	expected := source
	switch m.method {
	case snapshot.CalculationLatest:
		if m.counter != nil {
			return false
		}
	case snapshot.CalculationMean, snapshot.CalculationDelta:
		switch source {
		case measurement.ExactInt64, measurement.ExactUint64:
			expected = measurement.ExactDecimal
		case measurement.ExactFloat64:
		default:
			return false
		}
		if m.method == snapshot.CalculationMean {
			if m.counter != nil {
				return false
			}
		} else {
			if m.counter == nil || m.counter.Validate() != nil || m.counter.WorkspaceID != workspace || m.counter.GroupID != group || m.counter.ProfileRevision != revision || m.counter.MemberKey != m.key {
				return false
			}
		}
	default:
		return false
	}
	return m.result == expected
}

func mixedWideIssues(s GroupRowSpec) []LayoutIssue {
	bad := func() []LayoutIssue { return []LayoutIssue{{Code: issueMixedContract}} }
	if s.ResultShape == "" {
		if s.DueAtColumn != "" {
			return bad()
		}
		for _, m := range s.Members {
			if m.Calculation != "" || m.Interval != 0 {
				return bad()
			}
		}
		return nil
	}
	if s.ResultShape != MixedResultShapeV1 || s.WorkspaceID == "" || s.GroupID == "" || s.GroupRevision == "" || !s.Partial || s.ProvenanceColumn == "" || s.RecordKeyColumn == "" || s.DueAtColumn == "" || s.BucketStartColumn != "" || s.EntityKeyed || s.EntityKeyColumn != "" || len(s.Members) == 0 {
		return bad()
	}
	seen := map[string]bool{}
	for _, m := range s.Members {
		if seen[m.MemberKey] || m.EntityKey != "" || m.Reduction != "" || !validMixedIntent(m.mixedIntent(), s.WorkspaceID, s.GroupID, s.GroupRevision) {
			return bad()
		}
		seen[m.MemberKey] = true
	}
	return nil
}
func mixedLongIssues(s LongRowSpec) []LayoutIssue {
	bad := func() []LayoutIssue { return []LayoutIssue{{Code: issueMixedContract}} }
	if s.ResultShape == "" {
		if s.ProfileKind == "mixed" {
			return bad()
		}
		for _, m := range s.Members {
			if m.Calculation != "" || m.Interval != 0 {
				return bad()
			}
		}
		return nil
	}
	if s.ResultShape != MixedResultShapeV1 || s.ProfileKind != "mixed" || s.ProfileInterval != 0 || s.ProfileUnit != "" || len(s.Members) == 0 {
		return bad()
	}
	if s.Dialect != SQLDialectSQLite && s.Dialect != SQLDialectPostgres {
		return []LayoutIssue{{Code: issueUnsupportedDialect}}
	}
	seen := map[string]bool{}
	for _, m := range s.Members {
		if seen[m.MemberKey] || m.Reduction != "" || !validMixedIntent(m.mixedIntent(), s.WorkspaceID, s.GroupID, s.GroupRevision) {
			return bad()
		}
		seen[m.MemberKey] = true
	}
	return nil
}

func mixedRowError() error { return &GroupRowError{Code: issueMixedContract} }
func validateMixedResult(m mixedIntent, o snapshot.Outcome, r snapshot.MemberResult, wide bool) error {
	w := r.Window
	if o.DueAt.IsZero() || w == nil || w.Calculation != m.method || w.Interval != m.interval || r.MemberKey != m.key {
		return mixedRowError()
	}
	if _, err := snapshot.CheckedWindowAdd(o.DueAt, 0, "due at"); err != nil {
		return mixedRowError()
	}
	if r.Status == snapshot.MemberNotDue {
		if !wide || r.Reason != snapshot.ReasonNotDue || r.Sample != nil || r.Counter != nil || r.Reduction != nil || r.Evidence != nil || !w.AlignedStart.IsZero() || !w.AlignedEnd.IsZero() || !w.ActualStart.IsZero() || !w.ActualEnd.IsZero() || w.Truncated != "" {
			return mixedRowError()
		}
		if !w.NextDue.IsZero() {
			start, end, err := snapshot.CheckedWindowBounds(o.DueAt, m.interval)
			if err != nil || start.Equal(o.DueAt) || !end.Equal(w.NextDue) {
				return mixedRowError()
			}
		}
		return nil
	}
	start, end, err := snapshot.CheckedWindowBounds(w.AlignedStart, m.interval)
	if err != nil || !start.Equal(w.AlignedStart) || !end.Equal(w.AlignedEnd) || !end.Equal(o.DueAt) || w.ActualStart.Before(start) || w.ActualEnd.After(end) || !w.ActualEnd.After(w.ActualStart) {
		return mixedRowError()
	}
	cut := ""
	if w.ActualStart.After(start) {
		cut = mixedTruncationStart
	}
	if w.ActualEnd.Before(end) {
		if cut == "" {
			cut = "end"
		} else {
			cut = "both"
		}
	}
	if cut != w.Truncated {
		return mixedRowError()
	}
	if !wide && (!o.BucketStart.Equal(start) || !o.BucketEnd.Equal(end)) {
		return mixedRowError()
	}
	if !wide {
		actualStart, actualEnd := o.Window()
		if !actualStart.Equal(w.ActualStart) || !actualEnd.Equal(w.ActualEnd) || o.Truncated != w.Truncated {
			return mixedRowError()
		}
	}
	if m.method != snapshot.CalculationDelta && r.Counter != nil || m.method != snapshot.CalculationMean && (r.Reduction != nil || r.Evidence != nil) {
		return mixedRowError()
	}
	switch r.Status {
	case snapshot.MemberOK:
	case snapshot.MemberMissing, snapshot.MemberBad, snapshot.MemberStale, snapshot.MemberInvalid:
		if r.Reason == "" || m.required {
			return mixedRowError()
		}
	default:
		return mixedRowError()
	}
	if r.Sample != nil && r.Sample.MemberKey != m.key {
		return mixedRowError()
	}
	if r.Usable() {
		if r.Sample == nil || !snapshot.MixedInputMatches(m.source, m.method, m.counter, r.Sample.Value) || r.Sample.Quality != schema.QualityGood {
			return mixedRowError()
		}
		if r.Sample == nil || r.Sample.ObservedAt.Before(w.ActualStart) || !r.Sample.ObservedAt.Before(w.ActualEnd) {
			return mixedRowError()
		}
		if _, err := mixedMemberValue(m.method, r); err != nil {
			return err
		}
	}
	return nil
}
func (l *GroupRowLayout) validateMixedOutcome(o snapshot.Outcome) error {
	if l.spec.ResultShape == "" {
		if !o.DueAt.IsZero() {
			return mixedRowError()
		}
		for _, r := range o.Members {
			if r.Window != nil {
				return mixedRowError()
			}
		}
		return nil
	}
	if o.RecordID == "" || o.EffectKey == "" || !o.BucketStart.IsZero() || !o.BucketEnd.IsZero() || !o.WindowStart.IsZero() || !o.WindowEnd.IsZero() || o.EntityKey != "" || o.Truncated != "" || len(o.Members) != len(l.members) {
		return mixedRowError()
	}
	expected, err := snapshot.MixedDispatchRecordID(l.spec.WorkspaceID, l.spec.GroupID, l.spec.GroupRevision, o.DueAt)
	if err != nil || o.RecordID != expected {
		return mixedRowError()
	}
	seen := map[string]bool{}
	for _, r := range o.Members {
		m, ok := l.members[r.MemberKey]
		if !ok || seen[r.MemberKey] {
			return mixedRowError()
		}
		seen[r.MemberKey] = true
		if err := validateMixedResult(m.mixedIntent(), o, r, true); err != nil {
			return err
		}
	}
	return nil
}
func (l *LongRowLayout) validateMixedOutcome(o snapshot.Outcome, m LongRowMember, r snapshot.MemberResult) error {
	if l.spec.ResultShape == "" {
		if !o.DueAt.IsZero() || r.Window != nil {
			return mixedRowError()
		}
		return nil
	}
	if o.EntityKey != m.MemberKey {
		return mixedRowError()
	}
	expected, err := snapshot.ProfileRecordID(l.spec.WorkspaceID, l.spec.GroupID, l.spec.GroupRevision, m.MemberKey, "mixed", o.BucketStart)
	if err != nil || o.RecordID != expected {
		return mixedRowError()
	}
	return validateMixedResult(m.mixedIntent(), o, r, false)
}
