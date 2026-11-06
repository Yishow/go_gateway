package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

func (h *StudioV2RecordingHomeSourcesHandler) readSources(ctx context.Context, id string) (workspace.RecordingHomeSourceInput, error) {
	var out workspace.RecordingHomeSourceInput
	record, err := h.workspace.ReadSnapshot(ctx)
	if errors.Is(err, workspace.ErrNotFound) {
		return out, workspace.ErrRecordingHomeWorkspaceNotFound
	}
	if err != nil {
		return out, err
	}
	if record == nil || record.ID != id {
		return out, workspace.ErrRecordingHomeWorkspaceNotFound
	}
	out.Workspace = workspace.RecordingHomeSourceWorkspace{ID: record.ID, OrderedDeviceIDs: append([]string{}, record.OrderedDeviceIDs...)}
	members := map[string]bool{}
	for _, deviceID := range record.OrderedDeviceIDs {
		if deviceID == "" || members[deviceID] {
			return out, workspace.ErrRecordingHomeSourcesInconsistent
		}
		members[deviceID] = true
		d, readErr := h.devices.GetByID(ctx, deviceID)
		if errors.Is(readErr, sql.ErrNoRows) {
			return out, workspace.ErrRecordingHomeSourcesInconsistent
		}
		if readErr != nil {
			return out, readErr
		}
		if d == nil || d.ID != deviceID {
			return out, workspace.ErrRecordingHomeSourcesInconsistent
		}
		protocol, availability := string(d.Protocol), availabilityStatusOf(d)
		out.Devices = append(out.Devices, workspace.RecordingHomeSourceDevice{ID: d.ID, Name: d.Name, Protocol: &protocol, AvailabilityStatus: &availability})
	}
	rules, err := h.rules.ListByDeviceIDs(ctx, record.OrderedDeviceIDs)
	if err != nil {
		return out, err
	}
	seenRules := map[string]bool{}
	for _, r := range rules {
		if r == nil || r.ID == "" || seenRules[r.ID] || !members[r.DeviceID] {
			return out, workspace.ErrRecordingHomeSourcesInconsistent
		}
		seenRules[r.ID] = true
		skipped := []string{}
		// Unlike the legacy rules GET, malformed persisted skip data cannot become
		// an apparently complete observation by silently dropping exclusions.
		if r.SkippedAddresses != "" {
			if json.Unmarshal([]byte(r.SkippedAddresses), &skipped) != nil || skipped == nil {
				return out, workspace.ErrRecordingHomeSourcesInconsistent
			}
		}
		out.Rules = append(out.Rules, workspace.RecordingHomeSourceRule{ID: r.ID, WorkspaceID: id, DeviceID: r.DeviceID, RevisionID: r.RevisionID, Enabled: r.Enabled, SkippedAddresses: skipped})
	}
	// No pagination or enabled filter: disabled and unmapped points also affect
	// completeness. Foreign-device points are discarded exactly as in legacy.
	points, err := h.points.List(ctx, point.ListFilter{})
	if err != nil {
		return out, err
	}
	byPoint := map[string]*schema.Point{}
	for _, p := range points {
		if p == nil {
			return out, workspace.ErrRecordingHomeSourcesInconsistent
		}
		if !members[p.DeviceID] {
			continue
		}
		if p.ID == "" || byPoint[p.ID] != nil {
			return out, workspace.ErrRecordingHomeSourcesInconsistent
		}
		detached := *p
		byPoint[p.ID] = &detached
		out.Points = append(out.Points, workspace.RecordingHomeSourcePoint{ID: p.ID, DeviceID: p.DeviceID, Address: p.Address, Enabled: p.Enabled})
	}
	seenLinks, linkedPoints := map[string]bool{}, map[string]bool{}
	for _, r := range rules {
		links, readErr := h.rules.ListLinks(ctx, r.ID)
		if readErr != nil {
			return out, readErr
		}
		for _, link := range links {
			if link == nil || link.ID == "" || seenLinks[link.ID] || link.RuleID != r.ID || linkedPoints[link.PointID] {
				return out, workspace.ErrRecordingHomeSourcesInconsistent
			}
			seenLinks[link.ID], linkedPoints[link.PointID] = true, true
			out.Links = append(out.Links, workspace.RecordingHomeSourceLink{ID: link.ID, RuleID: link.RuleID, PointID: link.PointID, Address: link.Address, MappingID: link.MappingID, TagID: link.TagID})
			p := byPoint[link.PointID]
			if p == nil || p.DeviceID != r.DeviceID || p.Address != link.Address {
				return out, workspace.ErrRecordingHomeSourcesInconsistent
			}
			projected, readErr := h.readMapping(ctx, id, r, link, p)
			if readErr != nil {
				return out, readErr
			}
			if projected != nil {
				// Legacy rules transport exposes an empty persisted string as null,
				// so it cannot supply the array required by a mapped candidate.
				// Unconfigured links may keep the canonical parser's empty default.
				if r.SkippedAddresses == "" {
					return out, workspace.ErrRecordingHomeSourcesInconsistent
				}
				out.Mappings = append(out.Mappings, *projected)
			}
		}
	}
	return out, nil
}
