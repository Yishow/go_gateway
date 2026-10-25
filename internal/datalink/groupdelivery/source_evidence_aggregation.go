package groupdelivery

import (
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"
)

const (
	sourceReductionMean = "mean"
	sourceRoleGauge     = "gauge"
)

// SourceAggregation is the value-free evidence of one mean window, verified against what the applied layout froze for
// the source. It names the method and the actual counts and coverage; it never carries a value, a sum or its encoding.
// Unknown coverage stays unknown: complete would need a saved whole-window proof, and this change wires none.
type SourceAggregation struct {
	Method        string `json:"method"`
	Version       string `json:"version"`
	Basis         string `json:"basis"`
	SourceType    string `json:"source_type"`
	ResultType    string `json:"result_type"`
	ValidCount    int64  `json:"valid_count"`
	RejectedCount int64  `json:"rejected_count"`
	Coverage      string `json:"coverage"`
	Reason        string `json:"reason,omitempty"`
}

// sourceMean is the part of a mean block this projection reads. The stored sum, scale and rounding are not read.
type sourceMean struct {
	Reduction     string `json:"reduction"`
	Role          string `json:"role"`
	SourceType    string `json:"source_type"`
	ResultType    string `json:"result_type"`
	Version       string `json:"version"`
	Method        string `json:"method"`
	Basis         string `json:"basis"`
	ValidCount    int64  `json:"valid_count"`
	RejectedCount int64  `json:"rejected_count"`
	Coverage      string `json:"coverage"`
	Reason        string `json:"reason"`
}

// A stored coverage word is legal evidence only when something can back it. The three below describe what was
// observed; complete is excluded because no stored proof can back it.
var sourceMeanCoverages = map[string]bool{
	string(snapshot.CoveragePartial): true, string(snapshot.CoverageUnknown): true, string(snapshot.CoverageMissing): true,
}

func (r SourceRef) expectsMean() bool {
	return r.Reduction == sourceReductionMean && r.Role == sourceRoleGauge
}

// verifySourceAggregation checks a record's mean block against the source's frozen recording facts. It returns the
// evidence to report, or false when the block is missing where one is required, present where none belongs, or says
// anything the applied layout and the record's own status do not both support.
func verifySourceAggregation(ref SourceRef, status SourceStatus, mean *sourceMean) (*SourceAggregation, bool) {
	if mean == nil {
		return nil, !ref.expectsMean()
	}
	if !ref.expectsMean() || ref.SourceType == "" || mean.Reduction != sourceReductionMean || mean.Role != ref.Role ||
		mean.Version != measurement.MeanAggregationVersion || mean.Method != measurement.MeanAggregationMethod {
		return nil, false
	}
	exactDecimal := mean.Basis == string(measurement.MeanBasisExact) && mean.ResultType == string(measurement.ExactDecimal)
	floatMean := mean.Basis == string(measurement.MeanBasisFloat64) && mean.ResultType == string(measurement.ExactFloat64)
	if (!exactDecimal && !floatMean) || mean.SourceType != ref.SourceType || mean.ResultType != ref.ResultType {
		return nil, false
	}
	if mean.ValidCount < 0 || mean.RejectedCount < 0 || !sourceMeanCoverages[mean.Coverage] || (mean.ValidCount == 0) != (mean.Coverage == string(snapshot.CoverageMissing)) {
		return nil, false
	}
	if reading := status == SourceStatusOK || status == SourceStatusPartial; (mean.ValidCount > 0) != reading {
		return nil, false
	}
	return &SourceAggregation{
		Method: mean.Method, Version: mean.Version, Basis: mean.Basis, SourceType: mean.SourceType, ResultType: mean.ResultType,
		ValidCount: mean.ValidCount, RejectedCount: mean.RejectedCount, Coverage: mean.Coverage, Reason: sourceReason(mean.Reason),
	}, true
}
