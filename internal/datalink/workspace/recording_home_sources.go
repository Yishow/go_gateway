package workspace

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"slices"
	"strings"
)

// ReadRecordingHomeSources observes complete readonly source content twice.
// Agreement is not an atomic snapshot and grants no command authority.
func ReadRecordingHomeSources(ctx context.Context, workspaceID string, read RecordingHomeSourceReader) (RecordingHomeSourcesObservation, error) {
	unknown := RecordingHomeSourcesObservation{WorkspaceID: workspaceID, State: "unknown", ReasonCode: "source-inconsistent"}
	if strings.TrimSpace(workspaceID) == "" {
		return unknown, ErrRecordingHomeInvalidScope
	}
	if ctx == nil || read == nil {
		return unknown, ErrRecordingHomeSourcesUnavailable
	}
	readPass := func() (RecordingHomeSourceInput, error) {
		if err := ctx.Err(); err != nil {
			return RecordingHomeSourceInput{}, err
		}
		value, err := read(ctx, workspaceID)
		if err != nil {
			return RecordingHomeSourceInput{}, err
		}
		if err := ctx.Err(); err != nil {
			return RecordingHomeSourceInput{}, err
		}
		// Detach from repository-owned slices before invoking the reader again.
		b, err := json.Marshal(value)
		if err != nil {
			return RecordingHomeSourceInput{}, err
		}
		var detached RecordingHomeSourceInput
		if err := json.Unmarshal(b, &detached); err != nil {
			return RecordingHomeSourceInput{}, err
		}
		return normalizeHomeSources(detached), nil
	}
	before, err := readPass()
	if err != nil {
		return unknown, homeReadError(err)
	}
	beforeKey, err := homeFingerprint(before)
	if err != nil {
		return unknown, ErrRecordingHomeSourcesUnavailable
	}
	_, beforeValid := projectHomeSources(before, workspaceID)
	after, err := readPass()
	if err != nil {
		return unknown, homeReadError(err)
	}
	summary, afterValid := projectHomeSources(after, workspaceID)
	if !beforeValid || !afterValid {
		return unknown, nil
	}
	afterKey, err := homeFingerprint(after)
	if err != nil {
		return unknown, ErrRecordingHomeSourcesUnavailable
	}
	if beforeKey != afterKey {
		unknown.ReasonCode = "source-changed"
		return unknown, nil
	}
	// The digest describes presentation only. Full revision/content evidence above
	// fences the two reads but is not a command token or exposed digest input.
	digest, err := homePresentationDigest(workspaceID, summary)
	if err != nil {
		return unknown, ErrRecordingHomeSourcesUnavailable
	}
	return RecordingHomeSourcesObservation{WorkspaceID: workspaceID, State: "available", SourceDigest: hex.EncodeToString(digest[:]), RecordingHomeSourceSummary: &summary}, nil
}
func homeReadError(err error) error {
	if errors.Is(err, ErrRecordingHomeSourcesInconsistent) {
		return nil
	}
	if errors.Is(err, ErrRecordingHomeWorkspaceNotFound) {
		return ErrRecordingHomeWorkspaceNotFound
	}
	if errors.Is(err, context.Canceled) {
		return context.Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return context.DeadlineExceeded
	}
	return ErrRecordingHomeSourcesUnavailable
}
func normalizeHomeSources(v RecordingHomeSourceInput) RecordingHomeSourceInput {
	if v.Links == nil {
		v.Links = []RecordingHomeSourceLink{}
	}
	if v.Workspace.OrderedDeviceIDs == nil {
		v.Workspace.OrderedDeviceIDs = []string{}
	}
	if v.Devices == nil {
		v.Devices = []RecordingHomeSourceDevice{}
	}
	if v.Rules == nil {
		v.Rules = []RecordingHomeSourceRule{}
	}
	if v.Mappings == nil {
		v.Mappings = []RecordingHomeSourceMapping{}
	}
	points := make([]RecordingHomeSourcePoint, 0, len(v.Points))
	for _, p := range v.Points {
		if slices.Contains(v.Workspace.OrderedDeviceIDs, p.DeviceID) {
			points = append(points, p)
		}
	}
	v.Points = points
	for i := range v.Devices {
		if v.Devices[i].Protocol == nil {
			s := "unknown"
			v.Devices[i].Protocol = &s
		}
		if v.Devices[i].AvailabilityStatus == nil {
			s := "unknown"
			v.Devices[i].AvailabilityStatus = &s
		}
	}
	for i := range v.Rules {
		if v.Rules[i].SkippedAddresses == nil {
			v.Rules[i].SkippedAddresses = []string{}
		}
	}
	return v
}
func homeFingerprint(v RecordingHomeSourceInput) (string, error) {
	// Sorting copies prevents changing the order used by the legacy presentation.
	v.Devices = slices.Clone(v.Devices)
	v.Rules = slices.Clone(v.Rules)
	v.Mappings = slices.Clone(v.Mappings)
	v.Points = slices.Clone(v.Points)
	v.Links = slices.Clone(v.Links)
	slices.SortFunc(v.Devices, func(a, b RecordingHomeSourceDevice) int { return strings.Compare(a.ID, b.ID) })
	slices.SortFunc(v.Rules, func(a, b RecordingHomeSourceRule) int { return strings.Compare(a.ID, b.ID) })
	slices.SortFunc(v.Mappings, func(a, b RecordingHomeSourceMapping) int { return strings.Compare(a.ID, b.ID) })
	slices.SortFunc(v.Points, func(a, b RecordingHomeSourcePoint) int { return strings.Compare(a.ID, b.ID) })
	slices.SortFunc(v.Links, func(a, b RecordingHomeSourceLink) int { return strings.Compare(a.ID, b.ID) })
	b, err := json.Marshal(v)
	return string(b), err
}
func homeUnique[T any](items []T, id func(T) string) (map[string]T, bool) {
	m := make(map[string]T, len(items))
	for _, v := range items {
		key := id(v)
		if strings.TrimSpace(key) == "" {
			return nil, false
		}
		if _, ok := m[key]; ok {
			return nil, false
		}
		m[key] = v
	}
	return m, true
}
func projectHomeSources(v RecordingHomeSourceInput, scope string) (RecordingHomeSourceSummary, bool) {
	out := RecordingHomeSourceSummary{Devices: []RecordingHomeDeviceSummary{}, Candidates: []RecordingHomeCandidate{}, Excluded: []RecordingHomeExcluded{}}
	if v.Workspace.ID != scope || len(v.Workspace.OrderedDeviceIDs) != len(v.Devices) {
		return out, false
	}
	members, ok := homeUnique(v.Workspace.OrderedDeviceIDs, func(s string) string { return s })
	if !ok {
		return out, false
	}
	devices, ok := homeUnique(v.Devices, func(d RecordingHomeSourceDevice) string { return d.ID })
	if !ok {
		return out, false
	}
	rules, ok := homeUnique(v.Rules, func(r RecordingHomeSourceRule) string { return r.ID })
	if !ok {
		return out, false
	}
	points, ok := homeUnique(v.Points, func(p RecordingHomeSourcePoint) string { return p.ID })
	if !ok {
		return out, false
	}
	if _, ok = homeUnique(v.Mappings, func(m RecordingHomeSourceMapping) string { return m.ID }); !ok {
		return out, false
	}
	for _, d := range v.Devices {
		if _, ok = members[d.ID]; !ok {
			return out, false
		}
		out.Devices = append(out.Devices, RecordingHomeDeviceSummary{ID: d.ID, Name: d.Name, Protocol: *d.Protocol, Availability: *d.AvailabilityStatus})
	}
	for _, r := range v.Rules {
		if r.WorkspaceID != scope {
			return out, false
		}
		if _, ok = devices[r.DeviceID]; !ok {
			return out, false
		}
	}
	claimed := map[string]bool{}
	for _, m := range v.Mappings {
		r, hasRule := rules[m.RuleID]
		p, hasPoint := points[m.PointID]
		if m.WorkspaceID != scope || !hasRule || !hasPoint || p.DeviceID != m.DeviceID || r.DeviceID != m.DeviceID || p.Address != m.Address || claimed[p.ID] {
			return out, false
		}
		claimed[p.ID] = true
		if !p.Enabled || !r.Enabled || !m.Enabled || slices.Contains(r.SkippedAddresses, p.Address) {
			continue
		}
		label := homeTrimLabel(m.DisplayName)
		if label == "" {
			label = p.Address
		}
		if m.TagID == "" {
			out.Excluded = append(out.Excluded, RecordingHomeExcluded{Label: label, Reason: "tag-missing"})
			continue
		}
		out.Candidates = append(out.Candidates, RecordingHomeCandidate{Key: p.ID + "|" + m.TagID, DeviceID: p.DeviceID, PointID: p.ID, TagID: m.TagID, TagKey: m.TagKey, Label: label, DeviceName: devices[p.DeviceID].Name, Address: p.Address, TargetType: m.TargetType})
	}
	return out, true
}

func homePresentationDigest(workspaceID string, summary RecordingHomeSourceSummary) ([32]byte, error) {
	summary.Devices = slices.Clone(summary.Devices)
	summary.Candidates = slices.Clone(summary.Candidates)
	summary.Excluded = slices.Clone(summary.Excluded)
	slices.SortFunc(summary.Devices, func(a, b RecordingHomeDeviceSummary) int { return strings.Compare(a.ID, b.ID) })
	slices.SortFunc(summary.Candidates, func(a, b RecordingHomeCandidate) int { return strings.Compare(a.Key, b.Key) })
	slices.SortFunc(summary.Excluded, func(a, b RecordingHomeExcluded) int {
		if n := strings.Compare(a.Label, b.Label); n != 0 {
			return n
		}
		return strings.Compare(a.Reason, b.Reason)
	})
	value := struct {
		WorkspaceID string `json:"workspace_id"`
		RecordingHomeSourceSummary
	}{workspaceID, summary}
	b, err := json.Marshal(value)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(b), nil
}
