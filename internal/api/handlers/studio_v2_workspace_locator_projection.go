package handlers

import (
	"encoding/json"
	"reflect"
	"time"

	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

type studioLocatorHydrationScope struct {
	CompleteWorkspace bool     `json:"complete_workspace"`
	DeviceIDs         []string `json:"device_ids"`
	RuleIDs           []string `json:"rule_ids"`
	MappingIDs        []string `json:"mapping_ids"`
	GroupIDs          []string `json:"group_ids"`
}

type studioLocatorEntry struct {
	*studioLocatorEditorProjection
	WorkspaceRevision     string                             `json:"workspace_revision"`
	HydrationScope        studioLocatorHydrationScope        `json:"hydration_scope"`
	Devices               []studioV2WorkspaceDeviceResponse  `json:"devices,omitempty"`
	SourceRules           []sourceRuleResponse               `json:"source_rules,omitempty"`
	Mappings              []studioV2WorkspaceMappingResponse `json:"mappings,omitempty"`
	WriteGroups           *workspace.WriteGroupListResult    `json:"write_groups,omitempty"`
	RuleHeads             map[string]string                  `json:"rule_heads"`
	GroupHeads            map[string]string                  `json:"group_heads"`
	RevalidatedMappingIDs []string                           `json:"revalidated_mapping_ids"`
}

func projectStudioLocator(record *workspace.Record, o *studioLocatorObservation) (*studioLocatorEntry, error) {
	entry := &studioLocatorEntry{WorkspaceRevision: record.DatabaseSetupRevision, HydrationScope: studioLocatorHydrationScope{DeviceIDs: []string{}, RuleIDs: []string{}, MappingIDs: []string{}, GroupIDs: []string{}}, RuleHeads: map[string]string{}, GroupHeads: map[string]string{}, RevalidatedMappingIDs: []string{}}
	if o.Device != nil {
		item := mapStudioV2WorkspaceDeviceResponse(o.Device, "", "")
		// Native editable configuration stays in this setup-only DTO. Diagnostics never do.
		item.LastTestError, item.AvailabilityReason, item.RuntimeApplyMessage = "", "", ""
		entry.Devices = []studioV2WorkspaceDeviceResponse{item}
		entry.HydrationScope.DeviceIDs = append(entry.HydrationScope.DeviceIDs, item.ID)
	}
	if o.Rule != nil {
		item := mapSourceRuleResponse(o.Rule)
		item.WorkspaceID = record.ID
		item.SaveState = workspaceSaveStateSaved
		if item.SkippedAddresses == nil {
			item.SkippedAddresses = []string{}
		}
		entry.SourceRules = []sourceRuleResponse{item}
		entry.HydrationScope.RuleIDs = append(entry.HydrationScope.RuleIDs, item.ID)
	}
	if o.Mapping != nil {
		item, err := projectStudioLocatorMapping(record.ID, o)
		if err != nil {
			return nil, err
		}
		entry.Mappings = []studioV2WorkspaceMappingResponse{item}
		entry.HydrationScope.MappingIDs = append(entry.HydrationScope.MappingIDs, item.ID)
	}
	if o.Group != nil {
		entry.WriteGroups = &workspace.WriteGroupListResult{WorkspaceID: record.ID, WorkspaceRevision: o.Group.WorkspaceRevision, Groups: []*workspace.WriteGroup{o.Group.Group}}
		entry.HydrationScope.GroupIDs = append(entry.HydrationScope.GroupIDs, o.Group.Group.ID)
	}
	if o.Editor != nil {
		if err := projectStudioLocatorEditor(record, o.Editor, entry); err != nil {
			return nil, err
		}
	}
	return entry, nil
}

func projectStudioLocatorMapping(workspaceID string, o *studioLocatorObservation) (studioV2WorkspaceMappingResponse, error) {
	empty := studioV2WorkspaceMappingResponse{}
	if o.Point == nil || o.Rule == nil || o.Link == nil || o.Tag == nil {
		return empty, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	// The native editor can represent only one scale followed by one cast. Do not
	// silently erase another operation or rewrite historical cast-before-scale.
	steps := runtimeSourceTransforms(o.Mapping.TransformPipeline)
	target, scale, offset := o.Point.DataType, 1.0, 0.0
	scaled, cast := false, false
	for _, step := range steps {
		switch step.Kind {
		case "scale":
			if scaled || cast || step.Multiplier == nil || step.Offset == nil {
				return empty, runtimeEvidenceError(runtimeview.ReasonUnavailable)
			}
			scale, offset, scaled = *step.Multiplier, *step.Offset, true
		case "cast":
			if cast || step.TargetType == nil {
				return empty, runtimeEvidenceError(runtimeview.ReasonUnavailable)
			}
			target, cast = schema.DataType(*step.TargetType), true
		default:
			return empty, runtimeEvidenceError(runtimeview.ReasonUnavailable)
		}
	}
	if !target.IsValid() || target != o.Tag.DataType {
		return empty, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	return studioV2WorkspaceMappingResponse{ID: o.Mapping.ID, WorkspaceID: workspaceID, PointID: o.Point.ID, RuleID: o.Rule.ID, DeviceID: o.Rule.DeviceID, Address: o.Link.Address, TagID: o.Tag.ID, TagKey: o.Tag.Key, DisplayName: o.Tag.DisplayName, Unit: o.Tag.Unit, TargetType: string(target), Scale: scale, Offset: offset, Enabled: o.Mapping.Enabled, SaveState: workspaceSaveStateSaved, CreatedAt: o.Mapping.CreatedAt.Format(time.RFC3339Nano), UpdatedAt: o.Mapping.UpdatedAt.Format(time.RFC3339Nano)}, nil
}

func cloneStudioLocatorValue[T any](value *T) (*T, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	var detached T
	if json.Unmarshal(raw, &detached) != nil {
		return nil, runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	return &detached, nil
}

func sameStudioLocatorRecord(a, b any) bool { return reflect.DeepEqual(a, b) }

func sameStudioLocatorObservation(a, b *studioLocatorObservation) bool {
	comparable := func(o *studioLocatorObservation) studioLocatorObservation {
		value := *o
		if o.Device != nil {
			// Existence/configuration is revalidated; live diagnostics are not a revision.
			d := *o.Device
			d.Status = ""
			d.ReadinessStatus = ""
			d.LastCollectedAt = nil
			d.CollectionCount = 0
			d.ErrorCount = 0
			d.LastTestAt = nil
			d.LastTestSuccess = nil
			d.LastTestError = ""
			d.UpdatedAt = time.Time{}
			value.Device = &d
		}
		if o.Point != nil {
			p := *o.Point
			p.LastReadAt = nil
			p.LastValue = nil
			p.LastError = ""
			p.UpdatedAt = time.Time{}
			value.Point = &p
		}
		return value
	}
	av, bv := comparable(a), comparable(b)
	av.Editor, bv.Editor = nil, nil
	return reflect.DeepEqual(av, bv) && sameStudioLocatorEditor(a.Editor, b.Editor)
}
