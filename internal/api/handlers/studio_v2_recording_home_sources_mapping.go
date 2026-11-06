package handlers

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
)

func (h *StudioV2RecordingHomeSourcesHandler) readMapping(ctx context.Context, workspaceID string, r *schema.SourceRule, link *schema.SourceRuleLink, p *schema.Point) (*workspace.RecordingHomeSourceMapping, error) {
	// Even a direct ID cannot prove uniqueness. Inspect the complete point scope
	// without adopting a recovered mapping or changing the persisted link.
	records, err := h.mappings.List(ctx, mapping.ListFilter{PointID: &link.PointID})
	if err != nil {
		return nil, err
	}
	if link.MappingID == nil && link.TagID == nil && len(records) == 0 {
		return nil, nil
	}
	if link.MappingID == nil || link.TagID == nil || *link.MappingID == "" || *link.TagID == "" || len(records) != 1 {
		return nil, workspace.ErrRecordingHomeSourcesInconsistent
	}
	listed := records[0]
	if listed == nil || listed.ID != *link.MappingID || listed.PointID != p.ID || listed.TagID != *link.TagID {
		return nil, workspace.ErrRecordingHomeSourcesInconsistent
	}
	m, err := h.mappings.GetByID(ctx, *link.MappingID)
	if errors.Is(err, mapping.ErrMappingNotFound) {
		return nil, workspace.ErrRecordingHomeSourcesInconsistent
	}
	if err != nil {
		return nil, err
	}
	if m == nil || m.ID != listed.ID || m.PointID != listed.PointID || m.TagID != listed.TagID || m.Enabled != listed.Enabled || m.TransformPipeline != listed.TransformPipeline {
		return nil, workspace.ErrRecordingHomeSourcesInconsistent
	}
	t, err := h.tags.GetByID(ctx, *link.TagID)
	if errors.Is(err, tag.ErrTagNotFound) {
		return nil, workspace.ErrRecordingHomeSourcesInconsistent
	}
	if err != nil {
		return nil, err
	}
	if t == nil || t.ID != *link.TagID || (!isRuleManagedWorkspaceTagOwnedBy(t, r.ID, link.Address) && !isLegacyDirectWorkspaceBinding(link, m, t)) {
		return nil, workspace.ErrRecordingHomeSourcesInconsistent
	}
	target, _, _, err := decodeWorkspaceMappingPipeline(p.DataType, m.TransformPipeline)
	if err != nil {
		return nil, workspace.ErrRecordingHomeSourcesInconsistent
	}
	return &workspace.RecordingHomeSourceMapping{ID: m.ID, WorkspaceID: workspaceID, PointID: p.ID, RuleID: r.ID, DeviceID: r.DeviceID, Address: link.Address, TagID: t.ID, TagKey: t.Key, DisplayName: t.DisplayName, TargetType: string(target), Enabled: m.Enabled}, nil
}
