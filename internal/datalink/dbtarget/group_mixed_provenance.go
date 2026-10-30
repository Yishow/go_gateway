package dbtarget

import (
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"
	"time"
)

type mixedWindowEvidence struct {
	AlignedStart string `json:"aligned_start"`
	AlignedEnd   string `json:"aligned_end"`
	ActualStart  string `json:"actual_start"`
	ActualEnd    string `json:"actual_end"`
	Truncated    string `json:"truncated,omitempty"`
}

// mixedMeanEvidence adds the projection flag without changing legacy JSON.
type mixedMeanEvidence struct {
	*longMeanProvenance
	Rounded *bool `json:"rounded,omitempty"`
}
type mixedResultEvidence struct {
	ResultShape     string                  `json:"result_shape"`
	Member          string                  `json:"member"`
	Calculation     snapshot.Calculation    `json:"calculation"`
	IntervalSeconds int64                   `json:"interval_seconds"`
	Alignment       string                  `json:"alignment"`
	DueAt           string                  `json:"due_at"`
	NextDue         string                  `json:"next_due,omitempty"`
	Status          snapshot.MemberStatus   `json:"status"`
	Reason          string                  `json:"reason,omitempty"`
	SourceType      string                  `json:"source_type"`
	ResultType      measurement.ExactType   `json:"result_type"`
	Window          *mixedWindowEvidence    `json:"window,omitempty"`
	SampleID        string                  `json:"sample_id,omitempty"`
	ObservedAt      string                  `json:"observed_at,omitempty"`
	ReceivedAt      string                  `json:"received_at,omitempty"`
	SourceRevision  string                  `json:"source_revision,omitempty"`
	MappingRevision string                  `json:"mapping_revision,omitempty"`
	Quality         string                  `json:"quality,omitempty"`
	Raw             *measurement.ExactValue `json:"raw,omitempty"`
	Mean            *mixedMeanEvidence      `json:"mean,omitempty"`
	Counter         *longCounterProvenance  `json:"counter,omitempty"`
}

func mixedProvenance(m mixedIntent, r snapshot.MemberResult, due time.Time) mixedResultEvidence {
	e := mixedResultEvidence{ResultShape: MixedResultShapeV1, Member: m.key, Calculation: m.method, IntervalSeconds: int64(m.interval / time.Second), Alignment: "utc_epoch", DueAt: longObservedText(due), Status: r.Status, Reason: r.Reason, SourceType: string(m.source), ResultType: m.result}
	if w := r.Window; w != nil {
		e.NextDue = longObservedText(w.NextDue)
		if r.Status != snapshot.MemberNotDue {
			e.Window = &mixedWindowEvidence{longObservedText(w.AlignedStart), longObservedText(w.AlignedEnd), longObservedText(w.ActualStart), longObservedText(w.ActualEnd), w.Truncated}
		}
	}
	if s := r.Sample; s != nil {
		e.SampleID = s.SampleID
		e.ObservedAt = longObservedText(s.ObservedAt)
		e.ReceivedAt = longObservedText(s.ReceivedAt)
		e.SourceRevision = s.SourceRevision
		e.MappingRevision = s.MappingRevision
		e.Quality = string(s.Quality)
		if s.Value.Type() != "" {
			value := s.Value
			e.Raw = &value
		}
	}
	if m.method == snapshot.CalculationMean && r.Status != snapshot.MemberNotDue {
		e.Mean = &mixedMeanEvidence{longMeanProvenance: meanProvenance("", m.source, m.result, r)}
		if r.Reduction != nil {
			rounded := r.Reduction.Rounded
			e.Mean.Rounded = &rounded
		}
	}
	if m.counter != nil && r.Counter != nil {
		e.Counter = wideCounterProvenance(m.counter, r.Counter)
		if raw := r.Counter.Raw; raw != nil {
			value := raw.Value
			e.Counter.Raw = &value
			e.Counter.RawSampleID = raw.SampleID
			e.Counter.RawObservedAt = longObservedText(raw.ObservedAt)
		}
	}
	return e
}
func (l *GroupRowLayout) mixedProvenanceEntries(o snapshot.Outcome) []mixedResultEvidence {
	entries := make([]mixedResultEvidence, 0, len(o.Members))
	for _, r := range o.Members {
		entries = append(entries, mixedProvenance(l.members[r.MemberKey].mixedIntent(), r, o.DueAt))
	}
	return entries
}
