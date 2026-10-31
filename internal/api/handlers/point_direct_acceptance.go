package handlers

import (
	"context"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

type pointDirectProofAuthority interface {
	CaptureDirectReadProof(string) (collector.PointMutationLease, string, bool)
	AcceptsDirectReadProof(collector.PointMutationLease, string) bool
}

func (r *defaultPointDirectReader) captureDirectEvidence(ctx context.Context, p *schema.Point) (*collector.DirectReadEvidence, error) {
	evidence := &collector.DirectReadEvidence{Point: *p}
	if authority, ok := r.profiles.(pointDirectProofAuthority); ok {
		evidence.Mutation, evidence.Scope, evidence.Guarded = authority.CaptureDirectReadProof(p.ID)
	}
	if p.FormatMode != "" && p.FormatMode != schema.PointFormatLegacy && !evidence.Guarded {
		return evidence, connector.ErrReadProfile
	}
	d, err := r.deviceSvc.GetByID(ctx, p.DeviceID)
	if err != nil {
		return evidence, err
	}
	evidence.ConnectionRevision = connector.ReadConnectionRevision(d.ConnectionConfig)
	evidence.Protocol = d.Protocol
	return evidence, nil
}

// AcceptsDirectRead compares current authority against the captured execution
// decision. It never changes the decoded value or adopts a newer profile.
func (r *defaultPointDirectReader) AcceptsDirectRead(ctx context.Context, cv collector.CollectedValue) bool {
	evidence := cv.DirectEvidence
	if r == nil || evidence == nil {
		return false
	}
	current, err := r.pointSvc.GetByID(ctx, cv.PointID)
	if err != nil || current == nil {
		return false
	}
	old := evidence.Point
	if current.ID != old.ID || current.DeviceID != old.DeviceID || current.RevisionID != old.RevisionID || current.Address != old.Address || current.Function != old.Function || current.DataType != old.DataType || current.DataFormat != old.DataFormat || current.FormatMode != old.FormatMode || current.FormatOverride != old.FormatOverride || current.FormatVersion != old.FormatVersion {
		return false
	}
	d, err := r.deviceSvc.GetByID(ctx, current.DeviceID)
	if err != nil || d.Protocol != evidence.Protocol || connector.ReadConnectionRevision(d.ConnectionConfig) != evidence.ConnectionRevision {
		return false
	}
	if evidence.Guarded {
		authority, ok := r.profiles.(pointDirectProofAuthority)
		if !ok || !authority.AcceptsDirectReadProof(evidence.Mutation, evidence.Scope) {
			return false
		}
	}
	if current.FormatMode != "" && current.FormatMode != schema.PointFormatLegacy {
		if !evidence.Guarded || r.profiles == nil {
			return false
		}
		profile, err := r.profiles.ResolvePointReadProfile(ctx, current)
		if err != nil || profile != evidence.Profile {
			return false
		}
		// Recheck after provider I/O; absence may have been invalidated while reading.
		authority := r.profiles.(pointDirectProofAuthority)
		return authority.AcceptsDirectReadProof(evidence.Mutation, evidence.Scope)
	}
	return !evidence.Profile.Present() || evidence.Profile.Resolution().Mode == schema.PointFormatLegacy
}
