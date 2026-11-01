package workspace

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"slices"

	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
)

// AcquisitionSourceGraph is private-backed SQL evidence, not a public Complete
// flag or a grant to install points. It is meaningful only through acceptance.
type AcquisitionSourceGraph struct {
	data   *acquisitionSourceGraphData
	reason string
}
type acquisitionSourceGraphData struct {
	owner                                    *SQLRepository
	workspaceID, witness, semantic, original string
	input                                    sourcerule.AcquisitionSourceRows
	excluded                                 []optimizer.NativeRange
	ruleSemantics                            map[string]string
}

// Complete reports whether canonical SQL and source geometry were both proven.
func (g AcquisitionSourceGraph) Complete() bool { return g.data != nil && g.data.owner != nil }

// WorkspaceID returns the original verified ownership scope, never a new authority.
func (g AcquisitionSourceGraph) WorkspaceID() string {
	if !g.Complete() {
		return ""
	}
	return g.data.workspaceID
}

// IdentityRevision returns the existing frozen provenance digest, including
// Point and raw rule revisions. It is not a new persisted epoch or plan hash.
func (g AcquisitionSourceGraph) IdentityRevision() string {
	if !g.Complete() {
		return ""
	}
	return g.data.original
}

// Reason is a closed safety reason; it never contains stored configuration.
func (g AcquisitionSourceGraph) Reason() string {
	if g.Complete() {
		return ""
	}
	if g.reason != "" {
		return g.reason
	}
	return "source-snapshot-unavailable"
}

// Device returns the detached source-only device configuration.
func (g AcquisitionSourceGraph) Device() schema.Device {
	if !g.Complete() {
		return schema.Device{}
	}
	return g.data.input.Device
}

// Points returns detached source configurations, including disabled points.
func (g AcquisitionSourceGraph) Points() []*schema.Point {
	if !g.Complete() {
		return nil
	}
	return cloneAcquisitionRows(g.data.input).Points
}

// Rules returns detached source fields and original opaque rule revisions.
// A revision is not rewritten merely because SQL source contents still match.
func (g AcquisitionSourceGraph) Rules() []*schema.SourceRule {
	if !g.Complete() {
		return nil
	}
	return cloneAcquisitionRows(g.data.input).Rules
}

// Excluded returns physical source exclusions; it never grants readability.
func (g AcquisitionSourceGraph) Excluded() []optimizer.NativeRange {
	if !g.Complete() {
		return nil
	}
	return slices.Clone(g.data.excluded)
}

// SemanticRevision excludes output selectors, receipts and installation leases.
func (g AcquisitionSourceGraph) SemanticRevision() string {
	if !g.Complete() {
		return ""
	}
	return g.data.semantic
}

// RuleSemanticRevision allows a scoped publisher to preserve foreign sources.
func (g AcquisitionSourceGraph) RuleSemanticRevision(ruleID string) (string, bool) {
	if !g.Complete() {
		return "", false
	}
	revision, ok := g.data.ruleSemantics[ruleID]
	return revision, ok
}

// AcquisitionSourceGraph rejects hand-built/JSON projections and same-device
// partial copies, even if their public Version was copied from a valid source.
func (p *RuntimeProjection) AcquisitionSourceGraph(deviceID string) (AcquisitionSourceGraph, bool) {
	if p == nil {
		return AcquisitionSourceGraph{}, false
	}
	graph := p.acquisitionGraphs[deviceID]
	if !graph.Complete() {
		return graph, false
	}
	input, err := projectionAcquisitionRows(p, deviceID)
	if err != nil || p.WorkspaceID != graph.data.workspaceID {
		return AcquisitionSourceGraph{reason: "source-projection-changed"}, false
	}
	digest, err := acquisitionProjectionDigest(input)
	if err != nil || digest != graph.data.original {
		return AcquisitionSourceGraph{reason: "source-projection-changed"}, false
	}
	return graph, true
}

func (s *Service) attachAcquisitionSourceGraphs(ctx context.Context, projection *RuntimeProjection) {
	reader, ok := s.repo.(interface {
		ReadAcquisitionSourceGraph(context.Context, *RuntimeProjection, string) (AcquisitionSourceGraph, error)
	})
	if !ok {
		return
	}
	projection.acquisitionGraphs = map[string]AcquisitionSourceGraph{}
	for _, deviceID := range projection.DeviceIDs {
		graph, err := reader.ReadAcquisitionSourceGraph(ctx, projection, deviceID)
		if err != nil {
			graph = AcquisitionSourceGraph{reason: graph.Reason()}
		}
		projection.acquisitionGraphs[deviceID] = graph
	}
}

func projectionAcquisitionRows(projection *RuntimeProjection, deviceID string) (sourcerule.AcquisitionSourceRows, error) {
	var result sourcerule.AcquisitionSourceRows
	if projection == nil || projection.WorkspaceID == "" || deviceID == "" || !slices.Contains(projection.DeviceIDs, deviceID) {
		return result, ErrAcquisitionSourceGraph
	}
	deviceFound := false
	for _, device := range projection.Devices {
		if device != nil && device.ID == deviceID {
			if deviceFound {
				return result, ErrAcquisitionSourceGraph
			}
			result.Device, deviceFound = *device, true
		}
	}
	if !deviceFound {
		return result, ErrAcquisitionSourceGraph
	}
	rules, points := map[string]bool{}, map[string]bool{}
	for _, rule := range projection.Rules {
		if rule != nil && rule.DeviceID == deviceID {
			result.Rules = append(result.Rules, rule)
			rules[rule.ID] = true
		}
	}
	for _, point := range projection.Points {
		if point != nil && point.DeviceID == deviceID {
			result.Points = append(result.Points, point)
			points[point.ID] = true
		}
	}
	for _, link := range projection.RuleLinks {
		if link == nil {
			return result, ErrAcquisitionSourceGraph
		}
		if rules[link.RuleID] {
			result.Links = append(result.Links, link)
		} else if points[link.PointID] {
			return result, ErrAcquisitionSourceGraph
		}
	}
	return cloneAcquisitionRows(result), nil
}

func cloneAcquisitionRows(input sourcerule.AcquisitionSourceRows) sourcerule.AcquisitionSourceRows {
	d := input.Device
	result := sourcerule.AcquisitionSourceRows{Device: schema.Device{ID: d.ID, Protocol: d.Protocol, Status: d.Status, ConnectionConfig: d.ConnectionConfig}}
	for _, r := range input.Rules {
		result.Rules = append(result.Rules, &schema.SourceRule{ID: r.ID, DeviceID: r.DeviceID, StartAddress: r.StartAddress, Count: r.Count, DataType: r.DataType, Enabled: r.Enabled, SkippedAddresses: r.SkippedAddresses, DataFormat: r.DataFormat, RevisionID: r.RevisionID, LayoutMode: r.LayoutMode, MixedItems: r.MixedItems})
	}
	for _, p := range input.Points {
		point := &schema.Point{ID: p.ID, DeviceID: p.DeviceID, Address: p.Address, Function: p.Function, DataType: p.DataType, DataFormat: p.DataFormat, FormatMode: p.FormatMode, FormatOverride: p.FormatOverride, FormatVersion: p.FormatVersion, RevisionID: p.RevisionID, Mode: p.Mode, Enabled: p.Enabled}
		if p.PollingGroupID != nil {
			group := *p.PollingGroupID
			point.PollingGroupID = &group
		}
		result.Points = append(result.Points, point)
	}
	for _, l := range input.Links {
		result.Links = append(result.Links, &schema.SourceRuleLink{ID: l.ID, RuleID: l.RuleID, Address: l.Address, PointID: l.PointID})
	}
	return result
}

// This private provenance digest additionally binds the original raw rule
// revision supplied to installed ReadProfile; it is not plan semantic identity.
func acquisitionProjectionDigest(input sourcerule.AcquisitionSourceRows) (string, error) {
	witness, err := sourcerule.AcquisitionSourceDigest(input, true)
	if err != nil {
		return "", err
	}
	revisions := make([][2]string, 0, len(input.Rules))
	for _, rule := range input.Rules {
		revisions = append(revisions, [2]string{rule.ID, rule.RevisionID})
	}
	slices.SortFunc(revisions, func(a, b [2]string) int {
		if a[0] < b[0] {
			return -1
		}
		if a[0] > b[0] {
			return 1
		}
		return 0
	})
	raw, err := json.Marshal(struct {
		Witness   string
		Revisions [][2]string
	}{witness, revisions})
	if err != nil {
		return "", ErrAcquisitionSourceGraph
	}
	digest := sha256.Sum256(raw)
	return hex.EncodeToString(digest[:]), nil
}
