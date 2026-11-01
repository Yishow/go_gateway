package collector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
)

var errAcquisitionSuperseded = errors.New("acquisition installation superseded")

const acquisitionRequestRandom = "random"

const planReasonSourceCoverageUnverified = "source_coverage_unverified"

// scheduledSnapshot is a detached view of installed source slots, never a
// recording selection or a second configuration authority.
type scheduledSnapshot struct {
	device     deviceConfig
	geometry   *acquisitionGeometry
	points     []pointInfo
	groupID    string
	cadence    time.Duration
	ticker     *groupTicker
	ctx        context.Context
	running    bool
	capturedAt time.Time
}

type scheduledPlan struct {
	snapshot           scheduledSnapshot
	input              optimizer.AcquisitionInput
	plan               optimizer.AcquisitionPlan
	bindings           map[string]connector.PlannedPointBinding
	points             map[string]pointInfo
	failures           map[string]error
	deadline           time.Time
	capabilityRevision string
	coverageReason     string
}

func buildScheduledPlan(snapshot scheduledSnapshot) (scheduledPlan, error) {
	if snapshot.ctx == nil {
		snapshot.ctx = context.Background()
	}
	out := scheduledPlan{snapshot: snapshot, bindings: map[string]connector.PlannedPointBinding{}, points: map[string]pointInfo{}, failures: map[string]error{}}
	config, err := connector.ResolveAcquisitionConfig(snapshot.device.Protocol, snapshot.device.Config)
	if err != nil {
		return out, err
	}
	identity := config.Identity(snapshot.device.ID)
	out.capabilityRevision = config.Revision()
	coverage := config.Coverage()
	if snapshot.geometry != nil {
		// The immutable geometry has already passed canonical SQL acceptance and
		// a full installed-slot proof. It adds exclusions, never readability.
		coverage.Excluded = append(coverage.Excluded, snapshot.geometry.excluded...)
		coverage.Revision, err = measurement.HashStringTuple(coverage.Revision, snapshot.geometry.revision)
		if err != nil {
			return out, err
		}
	} else if len(coverage.Readable) > 0 {
		// Active slots alone cannot prove disabled/skipped source cells absent.
		// Absolute forbidden and declared excluded ranges always remain effective.
		coverage.Readable = nil
		out.coverageReason = planReasonSourceCoverageUnverified
	}
	identity.PollingGroup, identity.Cadence = snapshot.groupID, snapshot.cadence
	points := slices.Clone(snapshot.points)
	slices.SortFunc(points, func(a, b pointInfo) int { return comparePointID(a.ID, b.ID) })
	proof := []string{config.Revision(), snapshot.groupID, snapshot.cadence.String()}
	for _, installed := range points {
		out.points[installed.ID] = installed
		profileProof, err := json.Marshal(installed.readProfile.IdentityView())
		if err != nil {
			return out, err
		}
		// Include rejected points too: an invalid source is still membership, not
		// an omitted input. This closed identity has no live values or timestamps.
		proof = append(proof, installed.ID, installed.configRevision, installed.Address, installed.Function, string(installed.DataType), installed.DataFormat, string(installed.Mode), installed.PollingGroupID, string(profileProof))
		binding, bits, err := plannedBinding(snapshot.device, identity.Protocol, installed)
		if err != nil {
			out.failures[installed.ID] = err
			continue
		}
		out.bindings[installed.ID] = binding
		proof = append(proof, binding.ConfigFingerprint)
		out.input.Sources = append(out.input.Sources, optimizer.AcquisitionSource{Identity: identity, Capability: config.Capability(bits), Coverage: coverage, Points: []optimizer.EffectivePoint{binding.Point}})
	}
	out.input.ConfigurationRevision, err = measurement.HashStringTuple(proof...)
	if err != nil {
		return out, err
	}
	out.plan = optimizer.BuildAcquisitionPlan(out.input)
	for _, failure := range out.plan.Errors() {
		out.failures[failure.PointID] = errors.New(failure.Code)
	}
	requests := out.plan.Requests()
	maximumRequests := len(requests)
	if slices.ContainsFunc(requests, func(r optimizer.AcquisitionRequest) bool { return r.Kind == acquisitionRequestRandom }) {
		safeRequests := len(optimizer.BuildAcquisitionPlan(safeAcquisitionInput(out.input, nil)).Requests())
		if safeRequests > int(^uint(0)>>1)-maximumRequests {
			return out, errors.New("acquisition request bound overflows")
		}
		maximumRequests += safeRequests
	}
	timeout := config.ReadTimeout()
	if maximumRequests > 0 && (timeout <= 0 || int64(maximumRequests) > int64(^uint64(0)>>1)/int64(timeout)) {
		return out, errors.New("acquisition deadline exceeds supported duration")
	}
	out.deadline = snapshot.capturedAt.Add(time.Duration(maximumRequests) * timeout)
	return out, nil
}

func comparePointID(a, b string) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func plannedBinding(cfg deviceConfig, protocol string, pt pointInfo) (connector.PlannedPointBinding, int, error) {
	point := pt.identityPoint()
	profile := pt.readProfile
	if !profile.Present() {
		var connection struct {
			DataFormat string `json:"data_format"`
		}
		if json.Unmarshal([]byte(cfg.Config), &connection) != nil {
			return connector.PlannedPointBinding{}, 0, connector.ErrReadProfile
		}
		var err error
		profile, err = connector.ResolveReadProfile(point, decoding.FormatContext{Protocol: cfg.Protocol, ConnectionFormat: connection.DataFormat, ConnectionRevision: connector.ReadConnectionRevision(cfg.Config)})
		if err != nil {
			return connector.PlannedPointBinding{}, 0, err
		}
	}
	if !profile.MatchesConnection(cfg.Protocol, cfg.Config) || connector.ValidatePointReadProfile(&point, profile) != nil {
		return connector.PlannedPointBinding{}, 0, connector.ErrReadProfile
	}
	resolution := profile.Resolution()
	sourceID, sourceRevision := profile.IdentityView().RuleID, resolution.RuleRevision
	if sourceID == "" {
		sourceID, sourceRevision = pt.ID, pt.configRevision
	}
	effective := optimizer.EffectivePoint{ID: pt.ID, SourceIdentity: sourceID, SourceRevision: sourceRevision, Address: pt.Address, Function: pt.Function, DataType: pt.DataType, Format: resolution.EffectiveFormat, FormatSource: string(resolution.Source), FormatRevision: pt.configRevision}
	bits, width, err := optimizer.AcquisitionPointShape(protocol, effective)
	if err != nil {
		return connector.PlannedPointBinding{}, 0, err
	}
	effective.NativeWidth = width
	fingerprint, err := measurement.AcquisitionConfigFingerprint(schema.Device{ID: cfg.ID, Protocol: cfg.Protocol, ConnectionConfig: cfg.Config}, point, profile)
	if err != nil {
		return connector.PlannedPointBinding{}, 0, err
	}
	return connector.PlannedPointBinding{Point: effective, Profile: profile, PointRevision: pt.configRevision, ConfigFingerprint: fingerprint, InstallationScope: pt.installationScope, InstallationLease: pt.observationLease}, bits, nil
}

// safeAcquisitionInput preserves all frozen metadata and explicit coverage. A
// completed set excludes already-final points when one fallback is authorized.
func safeAcquisitionInput(input optimizer.AcquisitionInput, completed map[string]bool) optimizer.AcquisitionInput {
	out := optimizer.AcquisitionInput{ConfigurationRevision: input.ConfigurationRevision}
	for _, source := range input.Sources {
		safeSource := source
		safeSource.Capability.RandomSupport = "unsupported"
		safeSource.Points = nil
		for _, point := range source.Points {
			if !completed[point.ID] {
				safeSource.Points = append(safeSource.Points, point)
			}
		}
		if len(safeSource.Points) > 0 {
			out.Sources = append(out.Sources, safeSource)
		}
	}
	return out
}

func (s *Scheduler) admitScheduled(snapshot scheduledSnapshot, points []pointInfo) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.scheduledCurrentLocked(snapshot, points)
}

func (s *Scheduler) scheduledCurrentLocked(snapshot scheduledSnapshot, points []pointInfo) error {
	if s.deviceConfigs[snapshot.device.ID] != snapshot.device || (snapshot.running && (!s.running || s.pollCtx != snapshot.ctx)) {
		return errAcquisitionSuperseded
	}
	if snapshot.ticker != nil && s.groupTickers[snapshot.groupID] != snapshot.ticker {
		return errAcquisitionSuperseded
	}
	for _, captured := range points {
		current, ok := s.pointInfos[captured.ID]
		if !ok || current.recoveryPending || !s.installationConnectionCurrentLocked(current) || current.observationLease != captured.observationLease || current.installationScope != captured.installationScope || current.configRevision != captured.configRevision {
			return fmt.Errorf("%w: point installation changed", errAcquisitionSuperseded)
		}
	}
	return nil
}
