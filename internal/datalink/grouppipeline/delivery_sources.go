package grouppipeline

import (
	"context"
	snapshots "go-gateway/internal/datalink/snapshot"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"
)

// Source evidence states and the reasons it can be unavailable. Unavailable is never a statement about the
// destination: it only says this view cannot prove the per source result.
const (
	SourceEvidenceAvailable   = "available"
	SourceEvidenceUnavailable = "unavailable"

	SourceEvidenceNoSnapshot       = "applied_snapshot_unavailable"
	SourceEvidenceRevisionMismatch = "revision_mismatch"
)

// SourceEvidenceView is the per source result of the applied revision of a managed long or wide profile group.
type SourceEvidenceView struct {
	State             string `json:"state"`
	UnavailableReason string `json:"unavailable_reason,omitempty"`
	GroupRevision     string `json:"group_revision"`
	// ConnectorID and ConnectorRevision are the destination the applied revision froze, which is where every
	// record below was written. They are never the editable draft's destination.
	ConnectorID       string               `json:"connector_id,omitempty"`
	ConnectorRevision string               `json:"connector_revision,omitempty"`
	Sources           []SourceEvidenceItem `json:"sources"`
}

// SourceEvidenceItem is one applied source: its newest local record and the newest one the destination confirmed.
type SourceEvidenceItem struct {
	CircuitID        string                      `json:"circuit_id,omitempty"`
	DeviceID         string                      `json:"device_id"`
	PointID          string                      `json:"point_id"`
	TagID            string                      `json:"tag_id"`
	Latest           *groupdelivery.SourceRecord `json:"latest"`
	LastSQLCommitted *groupdelivery.SourceRecord `json:"last_sql_committed"`
}

// sourceEvidence reads the applied sources from the group's immutable applied snapshot, never from the editable
// draft, and asks the durable store for each source's newest record. Managed long and managed wide profile groups
// have evidence; legacy wide, custom and unapplied groups have none. A mean or a wide profile is trusted only when the
// frozen applied layout it was written under still proves out, and what each record says is checked against what that
// layout froze for the source.
func (p *Pipeline) sourceEvidence(ctx context.Context, groupID, appliedRevision string) (*SourceEvidenceView, error) {
	if appliedRevision == "" || p.deps.Groups == nil {
		return nil, nil
	}
	unavailable := func(reason string) *SourceEvidenceView {
		return &SourceEvidenceView{State: SourceEvidenceUnavailable, UnavailableReason: reason, GroupRevision: appliedRevision, Sources: []SourceEvidenceItem{}}
	}
	snapshot := p.appliedSnapshot(ctx, groupID)
	if snapshot == nil {
		return unavailable(SourceEvidenceNoSnapshot), nil
	}
	group := snapshot.Group
	layout := group.EffectiveLayoutVersion()
	wide := layout == workspace.WriteGroupLayoutManagedWideProfileV1
	if (layout != workspace.WriteGroupLayoutManagedLongV1 && !wide) || group.RecordingProfile == nil {
		return nil, nil
	}
	if snapshot.AppliedRevision != appliedRevision {
		return unavailable(SourceEvidenceRevisionMismatch), nil
	}
	if group.Destination.ConnectorID == "" || group.Destination.ConnectorRevision == "" {
		return unavailable(SourceEvidenceNoSnapshot), nil
	}
	reduction := ""
	if group.RecordingProfile.Kind == workspace.WriteGroupProfileInstantaneous {
		reduction = string(group.RecordingProfile.EffectiveSampleReduction())
	}
	// A wide profile and a mean read frozen recording facts, so the applied proof that froze them must still be valid.
	// A long group without a mean keeps the evidence it always had.
	if (group.RecordingProfile.Kind == workspace.WriteGroupProfileMixed || wide || reduction == string(workspace.WriteGroupSampleReductionMean)) && !appliedProofValid(snapshot, wide) {
		return unavailable(SourceEvidenceNoSnapshot), nil
	}
	scope := groupdelivery.SourceScope{
		WorkspaceID: snapshot.WorkspaceID, GroupID: groupID, GroupRevision: appliedRevision, ProfileKind: string(group.RecordingProfile.Kind),
		ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision,
	}
	if wide {
		policy := group.RowPolicy
		scope.Wide, scope.IntervalSeconds = true, group.RecordingProfile.IntervalSeconds()
		scope.RecordKeyColumn, scope.GroupIDColumn, scope.BucketStartColumn, scope.ProvenanceColumn =
			policy.RecordKeyColumn, policy.GroupIDColumn, policy.BucketStartColumn, policy.ProvenanceColumn
	}
	if group.RecordingProfile.Kind == workspace.WriteGroupProfileMixed {
		scope.ResultShape = string(group.RecordingProfile.ResultShape)
		scope.RecordIdentity = group.RecordID
		scope.DueAtColumn = group.RowPolicy.DueAtColumn
	}
	view := &SourceEvidenceView{
		State: SourceEvidenceAvailable, GroupRevision: appliedRevision, ConnectorID: scope.ConnectorID, ConnectorRevision: scope.ConnectorRevision,
		Sources: make([]SourceEvidenceItem, 0, len(group.Members)),
	}
	for _, member := range group.Members {
		ref := groupdelivery.SourceRef{
			CircuitID: member.CircuitID, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID,
			Column: member.TargetColumn, Reduction: reduction, Role: string(member.ValueRole),
			SourceType: string(member.SourceValueType), ResultType: string(member.ResultValueType),
		}
		if group.RecordingProfile.Kind == workspace.WriteGroupProfileMixed {
			if member.Cadence == nil {
				return unavailable(SourceEvidenceNoSnapshot), nil
			}
			interval, err := snapshots.CheckedCadenceDuration(member.Cadence.Interval, string(member.Cadence.Unit))
			if err != nil {
				return unavailable(SourceEvidenceNoSnapshot), nil //nolint:nilerr // Invalid frozen cadence is represented by typed unavailable evidence, not healthy data.
			}
			ref.Calculation = string(member.Calculation)
			ref.CadenceSeconds = int64(interval / time.Second)
			ref.Alignment = member.Cadence.Alignment
			ref.Required = member.Required
			if member.MeasurementID != nil {
				ref.MeasurementID = *member.MeasurementID
			}
		}
		records, err := p.deps.Store.SourceRecords(ctx, scope, ref)
		if err != nil {
			return nil, err
		}
		view.Sources = append(view.Sources, SourceEvidenceItem{
			CircuitID: member.CircuitID, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID,
			Latest: records.Latest, LastSQLCommitted: records.LastSQLCommitted,
		})
	}
	return view, nil
}

// appliedProofValid is the very proof the pipeline builds a writer from: the applied snapshot carries the frozen
// runtime layout and it still proves the recording shape. Anything less cannot support per source evidence.
func appliedProofValid(snapshot *workspace.WriteGroupAppliedSnapshot, wide bool) bool {
	if snapshot.RuntimeLayout == nil {
		return false
	}
	if mixedRecordingGroup(snapshot.Group) {
		return validMixedRuntimeStorageProof(snapshot.Group, snapshot.RuntimeLayout)
	}
	if wide {
		return validCanonicalManagedRuntimeProof(snapshot.Group, snapshot.RuntimeLayout)
	}
	return validLongRuntimeProof(snapshot.Group, snapshot.RuntimeLayout)
}

// appliedSnapshot is the immutable applied snapshot in force now, or nil when it cannot be resolved.
func (p *Pipeline) appliedSnapshot(ctx context.Context, groupID string) *workspace.WriteGroupAppliedSnapshot {
	snapshot, err := p.deps.Groups.ResolveAppliedAt(ctx, groupID, p.config.Now())
	if err != nil || snapshot == nil || snapshot.Group == nil {
		return nil
	}
	// The saved revision never carries the runtime layout. Prefer the verified
	// runtime version Apply froze for exactly this workspace, group and applied
	// revision; an invalid or missing one is never trusted and the caller then
	// fails closed on the layout-less resolved snapshot.
	key := groupdelivery.GroupKey{
		WorkspaceID: snapshot.WorkspaceID, GroupID: snapshot.GroupID, GroupRevision: snapshot.AppliedRevision,
	}
	frozen, loadErr := p.loadRuntimeSnapshot(ctx, key)
	if loadErr != nil {
		// A stored version that cannot be verified (malformed, mismatched key or a
		// failed proof) must never be replaced by the layout-less saved revision.
		return nil
	}
	if frozen != nil {
		return frozen
	}
	// No frozen runtime version at all: a legacy revision keeps the resolved
	// snapshot, while a new shape has no layout there and fails closed.
	return snapshot
}
