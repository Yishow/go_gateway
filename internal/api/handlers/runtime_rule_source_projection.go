package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"sort"
	"strings"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/schema"
)

func projectRuntimeSource(records *runtimeSourceRecords, scope runtimeview.Scope) (*runtimeview.SourceView, error) {
	rule, p := &records.Rule, &records.Point
	facts := &runtimeview.SourceFacts{
		Binding: runtimeview.Binding{WorkspaceID: scope.WorkspaceID, WorkspaceRevision: scope.WorkspaceRevision, EntityRevision: rule.RevisionID},
		Enabled: rule.Enabled && p.Enabled, Address: runtimeEvidenceString(runtimeSafeAddress(p.Address)),
		ReadType: runtimeEvidenceString(string(runtimeSafeDataType(p.DataType))), Transforms: []runtimeview.TransformStep{},
	}
	if runtimePointFunctionPattern.MatchString(p.Function) {
		facts.ProtocolFunction = runtimeEvidenceString(p.Function)
	}
	if format := runtimeSafeEnum(p.DataFormat, "ABCD", "BADC", "CDAB", "DCBA"); format != "" {
		facts.Transforms = append(facts.Transforms, runtimeview.TransformStep{Kind: "byte-order", Order: &format})
	}
	if group := records.Polling; group != nil && runtimeSafeID(group.ID) != "" && group.IntervalMs > 0 {
		facts.Polling = &runtimeview.PollingEvidence{GroupID: group.ID, IntervalMS: int64(group.IntervalMs), Source: "polling-group"}
	}
	result := &runtimeview.SourceView{RuleID: rule.ID, PointID: p.ID, SavedRevision: runtimeEvidenceString(rule.RevisionID), State: "saved-only", Proof: "unavailable", Saved: facts}
	if m := records.Mapping; m != nil {
		result.MappingID, result.TagID = runtimeEvidenceString(m.ID), runtimeEvidenceString(records.Tag.ID)
		revision, err := measurement.MappingRevision(*m)
		if err != nil {
			return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
		facts.MappingRevision = &revision
		if m.Status == schema.MappingStatusActive || m.Status == schema.MappingStatusValidated {
			facts.ResultType = runtimeEvidenceString(string(runtimeSafeDataType(records.Tag.DataType)))
			facts.Transforms = append(facts.Transforms, runtimeSourceTransforms(m.TransformPipeline)...)
		} else {
			facts.Transforms = append(facts.Transforms, runtimeview.TransformStep{Kind: "unsupported"})
		}
	}
	// No public authority provides a revision-bound immutable source snapshot.
	// A running collector or matching projection version must not fill Effective.
	return result, nil
}

type runtimeSourceTransform struct {
	Type   string          `json:"type"`
	Order  int             `json:"order"`
	Params json.RawMessage `json:"params"`
}

func runtimeSourceTransforms(raw string) []runtimeview.TransformStep {
	// Runtime ingestion explicitly treats a blank stored pipeline as no steps.
	if strings.TrimSpace(raw) == "" {
		return []runtimeview.TransformStep{}
	}
	unsupported := []runtimeview.TransformStep{{Kind: "unsupported"}}
	var steps []runtimeSourceTransform
	if !runtimeSourceDecode(raw, &steps) {
		return unsupported
	}
	for _, step := range steps {
		if step.Order < 0 {
			return unsupported
		}
	}
	// The existing mapping engine retains input order for ties (including all
	// zero orders). This is presentation of that ordering, not a new policy.
	sort.SliceStable(steps, func(i, j int) bool { return steps[i].Order < steps[j].Order })
	result := make([]runtimeview.TransformStep, 0, len(steps))
	for _, step := range steps {
		projected := runtimeview.TransformStep{Kind: "unsupported"}
		switch step.Type {
		case "scale":
			var params struct {
				Scale  *float64 `json:"scale"`
				Offset *float64 `json:"offset"`
			}
			if runtimeSourceDecode(string(step.Params), &params) && params.Scale != nil && params.Offset != nil && !math.IsNaN(*params.Scale) && !math.IsInf(*params.Scale, 0) && !math.IsNaN(*params.Offset) && !math.IsInf(*params.Offset, 0) {
				projected = runtimeview.TransformStep{Kind: "scale", Multiplier: params.Scale, Offset: params.Offset}
			}
		case "cast":
			var params struct {
				TargetType string `json:"target_type"`
			}
			if runtimeSourceDecode(string(step.Params), &params) && runtimeSafeDataType(schema.DataType(params.TargetType)) != "" {
				projected = runtimeview.TransformStep{Kind: "cast", TargetType: &params.TargetType}
			}
		}
		result = append(result, projected)
	}
	return result
}

func runtimeSourceDecode(raw string, output any) bool {
	decoder := json.NewDecoder(bytes.NewBufferString(raw))
	decoder.DisallowUnknownFields()
	if decoder.Decode(output) != nil {
		return false
	}
	return decoder.Decode(new(any)) == io.EOF
}
