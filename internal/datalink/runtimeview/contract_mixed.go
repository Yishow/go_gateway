package runtimeview

import (
	"encoding/hex"
	"github.com/google/uuid"
	"go-gateway/internal/datalink/snapshot"
	"time"
)

// MixedRecordingFacts is a projection of the canonical mixed profile, never a
// legacy common interval or a permission to activate/edit the recording.
type MixedRecordingFacts struct {
	RecordID    string `json:"record_id"`
	ResultShape string `json:"result_shape" enum:"member-calculation-v1"`
}

type MemberCalculationFacts struct {
	Method          string  `json:"method" enum:"latest,mean,delta"`
	Interval        int64   `json:"interval" min:"1"`
	Unit            string  `json:"unit" enum:"second,minute"`
	Alignment       string  `json:"alignment" enum:"utc_epoch"`
	Required        bool    `json:"required"`
	SourceRevision  string  `json:"source_revision"`
	MappingRevision string  `json:"mapping_revision"`
	MeasurementID   *string `json:"measurement_id"`
	ResultType      string  `json:"result_type" enum:"bool,text,int64,uint64,decimal,float64"`
}

// CalculationResult keeps selected-result presence and historical raw presence
// distinct. RowEmitted may be true for a wide not_due NULL slot.
type CalculationResult struct {
	ResultShape     string  `json:"result_shape" enum:"member-calculation-v1"`
	Method          string  `json:"method" enum:"latest,mean,delta"`
	IntervalSeconds int64   `json:"interval_seconds" min:"1"`
	Alignment       string  `json:"alignment" enum:"utc_epoch"`
	DueAt           string  `json:"due_at" format:"time"`
	NextDue         *string `json:"next_due" format:"time"`
	OutcomeDigest   string  `json:"outcome_digest"`
	OutcomeKind     string  `json:"outcome_kind" enum:"row,skipped,no_data"`
	Status          string  `json:"status" enum:"ok,missing,bad,stale,invalid,not_due" enums:"ok,missing,bad,stale,invalid,not_due"`
	Reason          *string `json:"reason"`
	ResultPresent   bool    `json:"result_present"`
	RowEmitted      bool    `json:"row_emitted"`
	RawPresent      bool    `json:"raw_present"`
	SourceType      string  `json:"source_type" enum:"bool,int16,uint16,int32,uint32,int64,uint64,float32,float64,string"`
	ResultType      string  `json:"result_type" enum:"bool,text,int64,uint64,decimal,float64"`
	AlignedStart    *string `json:"aligned_start" format:"time"`
	AlignedEnd      *string `json:"aligned_end" format:"time"`
	ActualStart     *string `json:"actual_start" format:"time"`
	ActualEnd       *string `json:"actual_end" format:"time"`
	Rounded         *bool   `json:"rounded"`
}

func validateMixedFacts(f *GroupFacts) error {
	if f.MixedRecording == nil {
		if f.Reduction == "member_calculation" {
			return errUnavailable
		}
		for _, m := range f.Members {
			if m.Calculation != nil {
				return errUnavailable
			}
		}
		return nil
	}
	m := f.MixedRecording
	id, err := uuid.Parse(m.RecordID)
	if err != nil || id == uuid.Nil || id.String() != m.RecordID || m.ResultShape != "member-calculation-v1" || f.RecordingProfile != nil || f.IntervalSeconds != nil || f.Reduction != "member_calculation" {
		return errUnavailable
	}
	for _, member := range f.Members {
		c := member.Calculation
		if c == nil || member.SourceValueType == nil || c.SourceRevision == "" || c.MappingRevision == "" {
			return errUnavailable
		}
		if _, err := snapshot.CheckedCadenceDuration(c.Interval, c.Unit); err != nil {
			return errUnavailable
		}
		if c.Method == "delta" && c.MeasurementID == nil {
			return errUnavailable
		}
	}
	return nil
}

func validateMixedResult(r *SourceResult, facts *GroupFacts, source SourceEvidence) error {
	if r.Calculation == nil {
		if facts.MixedRecording != nil && r.State != "unavailable" {
			return errUnavailable
		}
		return nil
	}
	c := r.Calculation
	digest, digestErr := hex.DecodeString(c.OutcomeDigest)
	if digestErr != nil || len(digest) != 32 || (c.Rounded != nil && c.Method != "mean") {
		return errUnavailable
	}
	for _, raw := range []*string{&c.DueAt, c.NextDue, c.AlignedStart, c.AlignedEnd, c.ActualStart, c.ActualEnd} {
		if raw == nil {
			continue
		}
		at, err := time.Parse(time.RFC3339Nano, *raw)
		if err != nil {
			return errUnavailable
		}
		if _, err = snapshot.CheckedWindowAdd(at, 0, "calculation evidence"); err != nil {
			return errUnavailable
		}
	}
	if c.NextDue != nil && !afterTime(*c.NextDue, c.DueAt) {
		return errUnavailable
	}
	if facts.MixedRecording == nil || c.ResultShape != facts.MixedRecording.ResultShape {
		return errUnavailable
	}
	var member *Member
	for i := range facts.Members {
		m := &facts.Members[i]
		if m.DeviceID == source.DeviceID && m.PointID == source.PointID && m.TagID == source.TagID {
			member = m
			break
		}
	}
	if member == nil || member.Calculation == nil || member.SourceValueType == nil {
		return errForeign
	}
	intent := member.Calculation
	duration, err := snapshot.CheckedCadenceDuration(intent.Interval, intent.Unit)
	if err != nil || c.IntervalSeconds != int64(duration/time.Second) || c.Method != intent.Method || c.Alignment != intent.Alignment || c.SourceType != *member.SourceValueType || c.ResultType != intent.ResultType {
		return errStale
	}
	if c.RowEmitted != (c.OutcomeKind == "row") || c.ResultPresent != (c.Status == "ok") || (c.ResultPresent && !c.RawPresent) || (r.State == "available" && !c.ResultPresent) || ((r.State == "not_due") != (c.Status == "not_due")) || (r.Receipt != nil && !c.RowEmitted) {
		return errUnavailable
	}
	if c.Method == "delta" {
		if r.DeltaKnown == nil || *r.DeltaKnown != c.ResultPresent {
			return errUnavailable
		}
	} else if r.DeltaKnown != nil {
		return errUnavailable
	}
	due, _ := time.Parse(time.RFC3339Nano, c.DueAt)
	if c.Status == "not_due" {
		if c.NextDue != nil {
			start, end, err := snapshot.CheckedWindowBounds(due, duration)
			next, _ := time.Parse(time.RFC3339Nano, *c.NextDue)
			if err != nil || start.Equal(due) || !end.Equal(next) {
				return errUnavailable
			}
		}
		if facts.Layout != "managed_wide_profile_v1" || c.AlignedStart != nil || c.AlignedEnd != nil || r.State != "not_due" || c.ResultPresent || c.RawPresent || c.ActualStart != nil || c.ActualEnd != nil || r.WindowStart != nil || r.WindowEnd != nil {
			return errUnavailable
		}
	} else {
		if c.ActualStart == nil || c.ActualEnd == nil || c.AlignedStart == nil || c.AlignedEnd == nil || r.WindowStart == nil || r.WindowEnd == nil {
			return errUnavailable
		}
		alignedStart, _ := time.Parse(time.RFC3339Nano, *c.AlignedStart)
		alignedEnd, _ := time.Parse(time.RFC3339Nano, *c.AlignedEnd)
		start, end, err := snapshot.CheckedWindowBounds(alignedStart, duration)
		if err != nil || !start.Equal(alignedStart) || !end.Equal(alignedEnd) || !end.Equal(due) {
			return errUnavailable
		}
		if c.NextDue != nil {
			expected, err := snapshot.CheckedWindowAdd(end, duration, "next member due")
			next, _ := time.Parse(time.RFC3339Nano, *c.NextDue)
			if err != nil || !expected.Equal(next) {
				return errUnavailable
			}
		}
		if !sameTime(*c.ActualStart, *r.WindowStart) || !sameTime(*c.ActualEnd, *r.WindowEnd) || !afterTime(*c.ActualEnd, *c.ActualStart) || afterTime(*c.AlignedStart, *c.ActualStart) || afterTime(*c.ActualEnd, *c.AlignedEnd) {
			return errUnavailable
		}
	}
	return nil
}
