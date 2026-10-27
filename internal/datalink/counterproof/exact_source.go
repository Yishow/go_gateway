package counterproof

import (
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// ExactRecordingSource freezes the confirmed measurement identity, scaling,
// and rollover policy of an integer counter member. It is backend-derived
// during write-group validation and immutable once applied.
type ExactRecordingSource struct {
	MeasurementID         string                     `json:"measurement_id"`
	DefinitionRevision    string                     `json:"definition_revision"`
	SourceBindingRevision string                     `json:"source_binding_revision"`
	SeriesEpoch           string                     `json:"series_epoch"`
	OriginalDataType      schema.DataType            `json:"original_data_type"`
	Scale                 string                     `json:"scale,omitempty"`
	Unit                  string                     `json:"unit"`
	SemanticKind          measurement.SemanticKind   `json:"semantic_kind"`
	CounterPolicy         *measurement.CounterPolicy `json:"counter_policy,omitempty"`
}

// Definition reconstructs the authoritative MeasurementDefinition fields that
// back this exact recording source proof.
func (s *ExactRecordingSource) Definition() measurement.MeasurementDefinition {
	if s == nil {
		return measurement.MeasurementDefinition{}
	}
	return measurement.MeasurementDefinition{
		ID:                    s.MeasurementID,
		DefinitionRevision:    s.DefinitionRevision,
		SourceBindingRevision: s.SourceBindingRevision,
		SeriesEpoch:           s.SeriesEpoch,
		Unit:                  s.Unit,
		SemanticKind:          s.SemanticKind,
		NumericEncoding:       string(s.OriginalDataType),
		CounterPolicy:         cloneOptionalCounterPolicy(s.CounterPolicy),
		ExactRecordingScale:   s.Scale,
	}
}

func cloneOptionalCounterPolicy(p *measurement.CounterPolicy) *measurement.CounterPolicy {
	if p == nil {
		return nil
	}
	return cloneExactCounterPolicy(p)
}
