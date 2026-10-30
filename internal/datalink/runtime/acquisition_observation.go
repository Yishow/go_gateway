package runtime

import (
	"context"
	"slices"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// acquisitionObservationInstall lives only during one projection application.
// Its scalar candidates come from that actual input, not a later editable row.
// The scheduler retains only a value copy bound to its exact point install lease.
type acquisitionObservationInstall struct {
	before     *workspace.Record
	version    string
	candidates map[string]collector.PointObservationBinding
	leases     map[string]uint64
}

func (s *Service) acquisitionWorkspaceSnapshot(ctx context.Context) (record *workspace.Record) {
	defer func() {
		if recover() != nil {
			record = nil
		}
	}()
	reader, ok := s.workspace.(sourceWorkspaceSnapshotReader)
	if !ok {
		return nil
	}
	record, err := reader.ReadSnapshot(ctx)
	if err != nil || record == nil {
		return nil
	}
	// Retain only detached observation metadata, never setup collections.
	return &workspace.Record{ID: record.ID, DatabaseSetupRevision: record.DatabaseSetupRevision, OrderedDeviceIDs: slices.Clone(record.OrderedDeviceIDs)}
}

func (s *Service) beginAcquisitionObservation(ctx context.Context, projection *workspace.RuntimeProjection) (capture *acquisitionObservationInstall) {
	if s == nil || s.scheduler == nil || !s.scheduler.HasAcquisitionObserver() {
		return nil
	}
	defer func() {
		if recover() != nil {
			capture = nil
		}
	}()
	before := s.acquisitionWorkspaceSnapshot(ctx)
	if projection == nil || before == nil || projection.WorkspaceID != before.ID || before.ID == "" || before.DatabaseSetupRevision == "" || projection.Version == "" {
		return nil
	}
	return &acquisitionObservationInstall{before: before, version: projection.Version, candidates: acquisitionObservationCandidates(projection, before), leases: make(map[string]uint64)}
}

func acquisitionObservationCandidates(projection *workspace.RuntimeProjection, before *workspace.Record) map[string]collector.PointObservationBinding {
	devices := make(map[string]*schema.Device)
	deviceCounts := make(map[string]int)
	rules := make(map[string]*schema.SourceRule)
	ruleCounts := make(map[string]int)
	links := make(map[string]*schema.SourceRuleLink)
	linkCounts := make(map[string]int)
	pointCounts := make(map[string]int)
	for _, device := range projection.Devices {
		if device != nil {
			devices[device.ID] = device
			deviceCounts[device.ID]++
		}
	}
	for _, rule := range projection.Rules {
		if rule != nil {
			rules[rule.ID] = rule
			ruleCounts[rule.ID]++
		}
	}
	for _, link := range projection.RuleLinks {
		if link != nil {
			links[link.PointID] = link
			linkCounts[link.PointID]++
		}
	}
	for _, point := range projection.Points {
		if point != nil {
			pointCounts[point.ID]++
		}
	}
	candidates := make(map[string]collector.PointObservationBinding)
	for _, point := range projection.Points {
		if point == nil || point.ID == "" || !point.Enabled || pointCounts[point.ID] != 1 || linkCounts[point.ID] != 1 {
			continue
		}
		device := devices[point.DeviceID]
		link := links[point.ID]
		rule := rules[link.RuleID]
		if device == nil || device.ID == "" || deviceCounts[device.ID] != 1 || device.Status != schema.DeviceStatusActive || !slices.Contains(before.OrderedDeviceIDs, device.ID) || !slices.Contains(projection.DeviceIDs, device.ID) {
			continue
		}
		if rule == nil || rule.ID == "" || ruleCounts[rule.ID] != 1 || rule.DeviceID != device.ID || !rule.Enabled || rule.RevisionID == "" {
			continue
		}
		fingerprint, err := measurement.AcquisitionConfigFingerprint(*device, *point)
		if err != nil || fingerprint == "" {
			continue
		}
		candidates[point.ID] = collector.PointObservationBinding{WorkspaceID: before.ID, WorkspaceRevision: before.DatabaseSetupRevision, DeviceID: device.ID, PointID: point.ID, RuleID: rule.ID, RuleRevision: rule.RevisionID, ConfigFingerprint: fingerprint}
	}
	return candidates
}

// installAcquisitionPoint preserves the native AddPoint path when observation is
// disabled or proof is missing. Both paths clear old proof before replacement.
func (s *Service) installAcquisitionPoint(capture *acquisitionObservationInstall, point *schema.Point) {
	if capture == nil {
		s.scheduler.AddPoint(point)
		return
	}
	capture.leases[point.ID] = s.scheduler.AddPointWithObservationLease(point)
}

func (s *Service) finishAcquisitionObservation(ctx context.Context, capture *acquisitionObservationInstall) {
	if capture == nil || len(capture.leases) == 0 {
		return
	}
	// Metadata providers and diagnostic binding cannot change the native result.
	defer func() { _ = recover() }()
	after := s.acquisitionWorkspaceSnapshot(ctx)
	if after == nil || after.ID != capture.before.ID || after.DatabaseSetupRevision != capture.before.DatabaseSetupRevision || !slices.Equal(after.OrderedDeviceIDs, capture.before.OrderedDeviceIDs) {
		return
	}
	for pointID, lease := range capture.leases {
		binding, ok := capture.candidates[pointID]
		if !ok || !s.sourceInputStillInstalled(binding.DeviceID, capture.version) {
			continue
		}
		// A rename-only reinstall may have the same acquisition fingerprint. Only
		// the exact scheduler-lifetime lease can bind; no result-time current lookup.
		s.scheduler.BindPointObservation(pointID, lease, binding)
	}
}
