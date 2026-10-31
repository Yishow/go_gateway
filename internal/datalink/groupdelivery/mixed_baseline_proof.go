package groupdelivery

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/counterproof"
)

// verifyMixedBaselineProof checks state against the very same sealed immutable
// facts that admitted its observations. It never looks up an editable definition.
func verifyMixedBaselineProof(ctx context.Context, q ownershipQuery, key GroupKey, baselines []CounterBaseline) error {
	if len(baselines) == 0 {
		return nil
	}
	var raw string
	if err := q.QueryRowContext(ctx, `SELECT payload FROM wg_runtime_versions WHERE workspace_id=? AND group_id=? AND group_revision=?`, key.WorkspaceID, key.GroupID, key.GroupRevision).Scan(&raw); err != nil {
		return err
	}
	proofs, err := counterproof.ResolveSnapshot([]byte(raw), counterproof.SnapshotScope{WorkspaceID: key.WorkspaceID, GroupID: key.GroupID, AppliedRevision: key.GroupRevision})
	if err != nil {
		return err
	}
	byMember := map[string]counterproof.CounterBindingProof{}
	for _, p := range proofs {
		byMember[p.Member.SampleMemberKey()] = p
	}
	for _, baseline := range baselines {
		cp, err := aggregation.DecodeCounterCheckpoint(baseline.Payload)
		if err != nil {
			return err
		}
		p, ok := byMember[baseline.MemberKey]
		if !ok {
			return ErrBaselineCorrupt
		}
		b := cp.Binding
		if b.WorkspaceID != key.WorkspaceID || b.GroupID != key.GroupID || b.ProfileRevision != key.GroupRevision || b.MemberKey != baseline.MemberKey {
			return ErrBaselineCorrupt
		}
		if b.MeasurementID != p.MeasurementID || b.DefinitionRevision != p.DefinitionRevision || b.SourceBindingRevision != p.SourceBindingRevision || b.SeriesEpoch != p.SeriesEpoch || string(b.Kind) != string(p.SemanticKind) || b.Unit != strings.TrimSpace(p.Unit) || cp.Reading.Value.Type() != p.ResultValueType {
			return ErrBaselineCorrupt
		}
		scale := ""
		if p.ScaleMode == counterproof.CounterScaleExact {
			scale = p.ExactScale
		}
		if b.ExactRecordingScale != scale {
			return ErrBaselineCorrupt
		}
		modulus, rate, err := p.CounterPolicy.ExactRollover()
		if err != nil {
			return err
		}
		// ExactValue contains immutable private representation; JSON is its canonical
		// public exact form, preserving decimal evidence without float conversion.
		got, err := json.Marshal([]any{b.RolloverModulus, b.MaxRatePerSecond})
		if err != nil {
			return ErrBaselineCorrupt
		}
		want, err := json.Marshal([]any{modulus, rate})
		if err != nil {
			return ErrBaselineCorrupt
		}
		if !reflect.DeepEqual(got, want) {
			return ErrBaselineCorrupt
		}
	}
	return nil
}
