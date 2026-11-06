package collector

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
)

// ErrAcquisitionPreviewStale means evidence changed during a read, not absence.
var ErrAcquisitionPreviewStale = errors.New("acquisition preview changed")

const acquisitionStatusUnavailable = "unavailable"

// AcquisitionPreviewProof is private-backed evidence for a final no-I/O recheck.
type AcquisitionPreviewProof struct {
	owner      *Scheduler
	receipt    AcquisitionGeometryReceipt
	snapshots  []scheduledSnapshot
	connection connector.AcquisitionConnectionView
	downgrade  acquisitionDowngrade
}

// PreviewAcquisition uses only installed capture and the execution planner.
// sourceRevision is the caller's verified persisted source graph, not an input plan.
func (s *Scheduler) PreviewAcquisition(ctx context.Context, workspaceID, deviceID, sourceRevision string) (AcquisitionPreview, AcquisitionPreviewProof, error) {
	out := AcquisitionPreview{DeviceID: deviceID, Status: acquisitionStatusUnavailable, Reason: "installation_unavailable"}
	proof := AcquisitionPreviewProof{}
	if s == nil {
		return out, proof, nil
	}
	if err := ctx.Err(); err != nil {
		return out, proof, err
	}
	s.mu.RLock()
	receipt := s.captureAcquisitionGeometryLocked(deviceID)
	if !receipt.authorized || workspaceID == "" || receipt.publication.workspaceID != workspaceID || sourceRevision == "" || receipt.publication.sourceRevision != sourceRevision {
		s.mu.RUnlock()
		return out, proof, nil
	}
	byGroup := map[string][]pointInfo{}
	for _, p := range receipt.slots {
		byGroup[p.PollingGroupID] = append(byGroup[p.PollingGroupID], p)
	}
	proof = AcquisitionPreviewProof{owner: s, receipt: receipt, downgrade: s.acquisitionDowngrades[deviceID]}
	for group, points := range byGroup {
		ticker := s.groupTickers[group]
		if ticker == nil || ticker.group == nil || !ticker.group.Enabled || ticker.group.IntervalMs <= 0 {
			s.mu.RUnlock()
			return out, AcquisitionPreviewProof{}, nil
		}
		proof.snapshots = append(proof.snapshots, scheduledSnapshot{device: receipt.device, geometry: receipt.publication, points: points, groupID: group, cadence: time.Duration(ticker.group.IntervalMs) * time.Millisecond, ticker: ticker, ctx: s.pollCtx, running: s.running, capturedAt: time.Now()})
	}
	s.mu.RUnlock()
	proof.connection = s.connMgr.PeekAcquisitionConnection(previewDevice(receipt.device))
	if proof.connection.State == acquisitionStatusUnavailable {
		out.Reason = "connection_context_unavailable"
		return out, AcquisitionPreviewProof{}, nil
	}
	slices.SortFunc(proof.snapshots, func(a, b scheduledSnapshot) int { return comparePointID(a.groupID, b.groupID) })
	out.Status, out.Reason, out.Estimate = "available", "", true
	out.Groups = make([]AcquisitionPreviewGroup, 0, len(proof.snapshots))
	heads := []string{deviceID, sourceRevision, receipt.device.connectionRevision, s.installationScope}
	for _, snapshot := range proof.snapshots {
		scan, err := buildScheduledPlan(snapshot)
		if err != nil {
			out.Status, out.Reason, out.Estimate = acquisitionStatusUnavailable, "plan_unavailable", false
			out.Groups = nil
			return out, AcquisitionPreviewProof{}, err
		}
		effective, basis := previewInitialPlan(scan, proof.connection, proof.downgrade)
		head, err := acquisitionGroupHead(scan, effective, proof.connection, proof.downgrade)
		if err != nil {
			return AcquisitionPreview{}, proof, err
		}
		heads = append(heads, head)
		group := AcquisitionPreviewGroup{PollingGroupID: snapshot.groupID, PollIntervalMS: snapshot.cadence.Milliseconds(), PlanRevision: effective.Fingerprint(), DeclaredPlanRevision: scan.plan.Fingerprint(), CapabilityRevision: scan.capabilityRevision, Basis: basis, AcquisitionPointCount: len(snapshot.points), PlannedRequests: len(effective.Requests()), UniqueWireUnits: []AcquisitionPreviewUnits{}, PointErrors: []AcquisitionPreviewPointError{}, FallbackReasons: []string{}, RequestLimits: []AcquisitionPreviewRequestLimit{}}
		for _, u := range effective.WireUnits() {
			group.UniqueWireUnits = append(group.UniqueWireUnits, AcquisitionPreviewUnits{Area: u.Group.Area, Function: u.Group.Function, AddressUnitBits: u.AddressUnitBits, UniqueAddressUnits: u.UniqueAddressUnits})
		}
		for id, failure := range scan.failures {
			code := "source_point_invalid"
			if errors.Is(failure, errSourceLayoutUnsupported) {
				code = "source-layout-unsupported"
			}
			group.PointErrors = append(group.PointErrors, AcquisitionPreviewPointError{PointID: id, Code: code})
		}
		for _, e := range effective.Errors() {
			if _, exists := scan.failures[e.PointID]; !exists {
				group.PointErrors = append(group.PointErrors, AcquisitionPreviewPointError{PointID: e.PointID, Code: e.Code})
			}
		}
		slices.SortFunc(group.PointErrors, func(a, b AcquisitionPreviewPointError) int { return comparePointID(a.PointID, b.PointID) })
		for _, r := range effective.Requests() {
			group.RequestLimits = append(group.RequestLimits, AcquisitionPreviewRequestLimit{Kind: r.Kind, Area: r.Group.Area, Function: r.Group.Function, NativeUnitLimit: r.Limit, RequestFrameBytes: r.FrameBytes.RequestBytes, ResponseFrameBytes: r.FrameBytes.ResponseBytes, FrameCostKnown: r.FrameCostKnown})
			if r.FallbackReason != "" && !slices.Contains(group.FallbackReasons, r.FallbackReason) {
				group.FallbackReasons = append(group.FallbackReasons, r.FallbackReason)
			}
		}
		if scan.coverageReason != "" {
			group.FallbackReasons = append(group.FallbackReasons, scan.coverageReason)
		}
		group.Actual = s.previewActual(snapshot, head)
		out.Groups = append(out.Groups, group)
	}
	revision, err := measurement.HashStringTuple(heads...)
	if err != nil {
		return AcquisitionPreview{}, proof, err
	}
	out.PreviewRevision = revision
	out.OutputLoss = s.previewOutputLoss(receipt)
	if !s.AcquisitionPreviewCurrent(proof) {
		return AcquisitionPreview{}, proof, ErrAcquisitionPreviewStale
	}
	return out, proof, nil
}

func previewDevice(cfg deviceConfig) schema.Device {
	return schema.Device{ID: cfg.ID, Protocol: cfg.Protocol, ConnectionConfig: cfg.Config}
}
func previewInitialPlan(scan scheduledPlan, connection connector.AcquisitionConnectionView, downgrade acquisitionDowngrade) (plan optimizer.AcquisitionPlan, basis string) {
	if connection.State == "absent" {
		return scan.plan, "declared_no_connection"
	}
	if downgrade.token == connection.Token && downgrade.capabilityRevision == scan.capabilityRevision {
		return initialSafeAcquisitionPlan(scan), "safe_current_connection"
	}
	return scan.plan, "declared_current_connection"
}
func acquisitionGroupHead(scan scheduledPlan, initial optimizer.AcquisitionPlan, connection connector.AcquisitionConnectionView, downgrade acquisitionDowngrade) (string, error) {
	scope := ""
	if scan.snapshot.geometry != nil {
		scope = scan.snapshot.geometry.workspaceID + "/" + scan.snapshot.geometry.sourceRevision
	}
	fields := make([]string, 0, 7+4*len(scan.snapshot.points))
	fields = append(fields, scope, scan.plan.Fingerprint(), initial.Fingerprint(), fmt.Sprint(connection), fmt.Sprint(downgrade), scan.snapshot.groupID, scan.snapshot.cadence.String())
	points := slices.Clone(scan.snapshot.points)
	slices.SortFunc(points, func(a, b pointInfo) int { return comparePointID(a.ID, b.ID) })
	for _, p := range points {
		fields = append(fields, p.ID, p.configRevision, p.installationScope, fmt.Sprint(p.observationLease))
	}
	return measurement.HashStringTuple(fields...)
}

// AcquisitionPreviewCurrent revalidates the original capture; never recaptures a replacement.
func (s *Scheduler) AcquisitionPreviewCurrent(proof AcquisitionPreviewProof) bool {
	if s == nil || proof.owner != s {
		return false
	}
	if s.connMgr.PeekAcquisitionConnection(previewDevice(proof.receipt.device)) != proof.connection {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.currentAcquisitionGeometryLocked(proof.receipt.deviceID) != proof.receipt.publication || !s.acquisitionSlotsCurrentLocked(proof.receipt.deviceID, proof.receipt.slots) || s.acquisitionDowngrades[proof.receipt.deviceID] != proof.downgrade {
		return false
	}
	for _, snapshot := range proof.snapshots {
		if s.scheduledCurrentLocked(snapshot, snapshot.points) != nil {
			return false
		}
	}
	return true
}
func (s *Scheduler) previewOutputLoss(receipt AcquisitionGeometryReceipt) *AcquisitionPreviewLoss {
	out := &AcquisitionPreviewLoss{Status: acquisitionStatusUnavailable, Reason: "no_matching_retained_detail"}
	for _, loss := range s.OutputLossDiagnostics().Recent {
		p, ok := receipt.slots[loss.PointID]
		if ok && p.installationScope == loss.InstallationScope && p.configRevision == loss.PointRevision {
			out.RetainedDetail = append(out.RetainedDetail, AcquisitionPreviewLossDetail{PointID: loss.PointID, Losses: loss.Losses})
		}
	}
	if len(out.RetainedDetail) > 0 {
		out.Status, out.Reason = "partial", ""
	}
	return out
}

// Initial failures are finalized before warm-safe planning in the executor.
// Preview excludes exactly that same set, preserving the canonical fingerprint.
func initialSafeAcquisitionPlan(scan scheduledPlan) optimizer.AcquisitionPlan {
	completed := make(map[string]bool, len(scan.failures))
	for id := range scan.failures {
		completed[id] = true
	}
	return optimizer.BuildAcquisitionPlan(safeAcquisitionInput(scan.input, completed))
}
