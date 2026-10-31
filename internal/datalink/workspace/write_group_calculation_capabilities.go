package workspace

import (
	"encoding/json"
	"go-gateway/internal/datalink/snapshot"
)

// WriteGroupCalculationCapabilities describes this backend's preparation
// contract. It is not a saved candidate, readiness token or activation proof.
type WriteGroupCalculationCapabilities struct {
	Version             int                                `json:"version"`
	ResultShape         WriteGroupResultShape              `json:"result_shape"`
	DraftSupported      bool                               `json:"draft_supported"`
	ActivationSupported bool                               `json:"activation_supported"`
	EditorSupported     bool                               `json:"editor_supported"`
	Reason              string                             `json:"reason"`
	Methods             []WriteGroupMemberCalculation      `json:"methods"`
	Cadence             WriteGroupCadenceCapability        `json:"cadence"`
	SchemaDialects      []string                           `json:"schema_dialects"`
	SchemaLayouts       []WriteGroupLayoutVersion          `json:"schema_layouts"`
	Unsupported         []WriteGroupCalculationRestriction `json:"unsupported_combinations"`
}

// WriteGroupCadenceCapability exposes the existing checked duration bound, not
// an arbitrary product maximum. Timestamps/window arithmetic are still checked.
type WriteGroupCadenceCapability struct {
	Alignment           string                   `json:"alignment"`
	Units               []WriteGroupIntervalUnit `json:"units"`
	Minimum             int64                    `json:"minimum"`
	MaxSeconds          int64                    `json:"max_seconds"`
	MaxMinutes          int64                    `json:"max_minutes"`
	WindowBoundsChecked bool                     `json:"window_bounds_checked"`
}

type WriteGroupCalculationRestriction struct {
	Method          WriteGroupMemberCalculation `json:"method"`
	Phase           string                      `json:"phase"`
	SourceTransform string                      `json:"source_transform"`
	Reason          string                      `json:"reason"`
}

// RecordingCalculationCapabilities returns fresh slices so consumers cannot
// mutate the backend contract. Public activation remains closed until the real
// scoped save/schema/runtime/SQL acceptance gate has passed.
func RecordingCalculationCapabilities() WriteGroupCalculationCapabilities {
	return WriteGroupCalculationCapabilities{
		Version: 1, ResultShape: WriteGroupResultShapeMemberCalculationV1,
		DraftSupported: true, ActivationSupported: false, EditorSupported: false,
		Reason:         "production_acceptance_pending",
		Methods:        []WriteGroupMemberCalculation{WriteGroupMemberCalculationLatest, WriteGroupMemberCalculationMean, WriteGroupMemberCalculationDelta},
		Cadence:        WriteGroupCadenceCapability{Alignment: WriteGroupAlignmentUTCEpoch, Units: []WriteGroupIntervalUnit{WriteGroupIntervalUnitSecond, WriteGroupIntervalUnitMinute}, Minimum: 1, MaxSeconds: snapshot.MaxCadenceSeconds, MaxMinutes: snapshot.MaxCadenceMinutes, WindowBoundsChecked: true},
		SchemaDialects: []string{"sqlite", "postgres"},
		SchemaLayouts:  []WriteGroupLayoutVersion{WriteGroupLayoutManagedWideProfileV1, WriteGroupLayoutManagedLongV1},
		Unsupported: []WriteGroupCalculationRestriction{
			{Method: WriteGroupMemberCalculationLatest, SourceTransform: "exact_counter_scale", Phase: "draft", Reason: "mixed_exact_source_draft_unsupported"},
			{Method: WriteGroupMemberCalculationMean, SourceTransform: "exact_counter_scale", Phase: "draft", Reason: "mixed_exact_source_draft_unsupported"},
			{Method: WriteGroupMemberCalculationDelta, SourceTransform: "exact_counter_scale", Phase: "draft", Reason: "mixed_exact_source_draft_unsupported"},
			{Method: WriteGroupMemberCalculationLatest, SourceTransform: "exact_counter_scale", Phase: "runtime", Reason: "mixed_exact_source_latest_unsupported"},
		},
	}
}

// MarshalJSON adds a read-only consumer capability alongside (not inside) the
// persisted group. Legacy profile/member JSON and revision material stay intact.
func (r WriteGroupSaveResult) MarshalJSON() ([]byte, error) {
	type plain WriteGroupSaveResult
	return json.Marshal(struct {
		plain
		Capabilities WriteGroupCalculationCapabilities `json:"calculation_capabilities"`
	}{plain(r), RecordingCalculationCapabilities()})
}

func (r WriteGroupListResult) MarshalJSON() ([]byte, error) {
	type plain WriteGroupListResult
	return json.Marshal(struct {
		plain
		Capabilities WriteGroupCalculationCapabilities `json:"calculation_capabilities"`
	}{plain(r), RecordingCalculationCapabilities()})
}
