package connector

import "go-gateway/internal/datalink/optimizer"

// Snapshot returns detached request/binding storage, rejecting a policy attached
// to another logical point. FormatRevision and source proofs remain opaque:
// installation/ingestion, not a transport adapter, decides their eligibility.
func (r PlannedReadRequest) Snapshot() (optimizer.AcquisitionRequest, []PlannedPointBinding, error) {
	request, ok := r.Plan.Request(r.RequestIndex)
	if r.Plan.Fingerprint() == "" || !ok || r.AcquisitionID == "" || r.RequestID == "" {
		return optimizer.AcquisitionRequest{}, nil, ErrPlannedRead
	}
	if len(request.Points) == 0 || len(r.Bindings) != len(request.Points) {
		return request, nil, ErrPlannedRead
	}
	bindings := make(map[string]PlannedPointBinding, len(r.Bindings))
	for _, b := range r.Bindings {
		if b.Point.ID == "" {
			return request, nil, ErrPlannedRead
		}
		if _, ok := bindings[b.Point.ID]; ok {
			return request, nil, ErrPlannedRead
		}
		bindings[b.Point.ID] = b
	}
	ordered := make([]PlannedPointBinding, len(request.Points))
	for i, p := range request.Points {
		b, ok := bindings[p.Point.ID]
		if !ok || b.Point != p.Point || !b.Profile.Present() {
			return request, nil, ErrPlannedRead
		}
		profile := b.Profile.IdentityView()
		if profile.Resolution.PointRevision != b.PointRevision || profile.PointID != p.Point.ID || profile.DeviceID != request.Group.Identity.DeviceID ||
			ValidateReadProfile(ReadRequest{Address: p.Point.Address, Function: p.Point.Function, DataType: p.Point.DataType, Profile: b.Profile}) != nil ||
			profile.Resolution.EffectiveFormat != p.Point.Format || string(profile.Resolution.Source) != p.Point.FormatSource ||
			(profile.Resolution.ConnectionRevision != "" && profile.Resolution.ConnectionRevision != request.Group.Identity.ConnectionRevision) {
			return request, nil, ErrPlannedRead
		}
		ordered[i] = b
	}
	return request, ordered, nil
}
