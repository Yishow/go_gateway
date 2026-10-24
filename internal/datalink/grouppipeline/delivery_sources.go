package grouppipeline

import (
	"context"

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

// SourceEvidenceView is the per source result of the applied revision of a managed long group.
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
// draft, and asks the durable store for each source's newest record. Wide, custom and unapplied groups have none.
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
	if group.EffectiveLayoutVersion() != workspace.WriteGroupLayoutManagedLongV1 || group.RecordingProfile == nil {
		return nil, nil
	}
	if snapshot.AppliedRevision != appliedRevision {
		return unavailable(SourceEvidenceRevisionMismatch), nil
	}
	if group.Destination.ConnectorID == "" || group.Destination.ConnectorRevision == "" {
		return unavailable(SourceEvidenceNoSnapshot), nil
	}
	scope := groupdelivery.SourceScope{
		WorkspaceID: snapshot.WorkspaceID, GroupID: groupID, GroupRevision: appliedRevision, ProfileKind: string(group.RecordingProfile.Kind),
		ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision,
	}
	view := &SourceEvidenceView{
		State: SourceEvidenceAvailable, GroupRevision: appliedRevision, ConnectorID: scope.ConnectorID, ConnectorRevision: scope.ConnectorRevision,
		Sources: make([]SourceEvidenceItem, 0, len(group.Members)),
	}
	for _, member := range group.Members {
		records, err := p.deps.Store.SourceRecords(ctx, scope, groupdelivery.SourceRef{
			CircuitID: member.CircuitID, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID,
		})
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

// appliedSnapshot is the immutable applied snapshot in force now, or nil when it cannot be resolved.
func (p *Pipeline) appliedSnapshot(ctx context.Context, groupID string) *workspace.WriteGroupAppliedSnapshot {
	snapshot, err := p.deps.Groups.ResolveAppliedAt(ctx, groupID, p.config.Now())
	if err != nil || snapshot == nil || snapshot.Group == nil {
		return nil
	}
	return snapshot
}
