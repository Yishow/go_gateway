package groupdelivery

import "time"

// SourceDeliveryWindowSkipped is a local calculation closure without an SQL
// row. It is distinct from an operator choosing to skip an outbox effect.
const SourceDeliveryWindowSkipped SourceDelivery = "window_skipped"

// SourceCalculation is value-free evidence of the selected calculation. Raw
// presence is historical provenance, never current source health. A usable
// result can exist even when another required member prevents row emission.
type SourceCalculation struct {
	ResultShape     string     `json:"result_shape"`
	Method          string     `json:"method"`
	IntervalSeconds int64      `json:"interval_seconds"`
	Alignment       string     `json:"alignment"`
	DueAt           time.Time  `json:"due_at"`
	NextDue         *time.Time `json:"next_due,omitempty"`
	OutcomeDigest   string     `json:"outcome_digest"`
	OutcomeKind     string     `json:"outcome_kind"`
	Status          string     `json:"status"`
	Reason          string     `json:"reason,omitempty"`
	ResultPresent   bool       `json:"result_present"`
	RowEmitted      bool       `json:"row_emitted"`
	RawPresent      bool       `json:"raw_present"`
	SourceType      string     `json:"source_type"`
	ResultType      string     `json:"result_type"`
	AlignedStart    *time.Time `json:"aligned_start,omitempty"`
	AlignedEnd      *time.Time `json:"aligned_end,omitempty"`
	ActualStart     *time.Time `json:"actual_start,omitempty"`
	ActualEnd       *time.Time `json:"actual_end,omitempty"`
	Rounded         *bool      `json:"rounded,omitempty"`
}
