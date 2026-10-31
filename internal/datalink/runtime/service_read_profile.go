package runtime

import (
	"context"
	"encoding/json"
	"slices"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

type readProfilesKey struct{}

func modernPoint(p *schema.Point) bool {
	return p != nil && p.FormatMode != "" && p.FormatMode != schema.PointFormatLegacy
}

// Prepare mandatory decoding independently of optional observation. All provider
// reads and resolution happen before metadata/scheduler locks are acquired.
func prepareProjectionProfiles(ctx context.Context, projection *workspace.RuntimeProjection, deviceIDs ...string) (context.Context, error) {
	profiles := map[string]connector.ReadProfile{}
	for _, p := range projection.Points {
		if p == nil || !p.Enabled || !modernPoint(p) || (len(deviceIDs) > 0 && !slices.Contains(deviceIDs, p.DeviceID)) {
			continue
		}
		active := false
		for _, d := range projection.Devices {
			if d != nil && d.ID == p.DeviceID && d.Status == schema.DeviceStatusActive {
				active = true
			}
		}
		if !active {
			continue
		}
		profile, err := projectionReadProfile(projection, p)
		if err != nil {
			return ctx, err
		}
		profiles[p.ID] = profile
	}
	return context.WithValue(ctx, readProfilesKey{}, profiles), nil
}
func projectionReadProfile(projection *workspace.RuntimeProjection, p *schema.Point, allowInactive ...bool) (connector.ReadProfile, error) {
	invalid := func() (connector.ReadProfile, error) { return connector.ReadProfile{}, connector.ErrReadProfile }
	if projection == nil || p == nil || !slices.Contains(projection.DeviceIDs, p.DeviceID) {
		return invalid()
	}
	var device *schema.Device
	for _, d := range projection.Devices {
		if d != nil && d.ID == p.DeviceID {
			if device != nil {
				return invalid()
			}
			device = d
		}
	}
	if device == nil || (device.Status != schema.DeviceStatusActive && (len(allowInactive) == 0 || !allowInactive[0])) {
		return invalid()
	}
	var persisted *schema.Point
	for _, candidate := range projection.Points {
		if candidate != nil && candidate.ID == p.ID {
			if persisted != nil {
				return invalid()
			}
			persisted = candidate
		}
	}
	if persisted == nil || persisted.DeviceID != p.DeviceID || persisted.RevisionID != p.RevisionID || persisted.Address != p.Address || persisted.Function != p.Function || persisted.DataType != p.DataType || persisted.FormatMode != p.FormatMode || persisted.FormatOverride != p.FormatOverride || persisted.FormatVersion != p.FormatVersion || persisted.DataFormat != p.DataFormat {
		return invalid()
	}
	defaults := decoding.FormatContext{Protocol: device.Protocol, ConnectionRevision: connector.ReadConnectionRevision(device.ConnectionConfig)}
	var config struct {
		DataFormat string `json:"data_format"`
	}
	if json.Unmarshal([]byte(device.ConnectionConfig), &config) != nil {
		return invalid()
	}
	defaults.ConnectionFormat = config.DataFormat
	var link *schema.SourceRuleLink
	for _, candidate := range projection.RuleLinks {
		if candidate != nil && candidate.PointID == p.ID {
			if link != nil {
				return invalid()
			}
			link = candidate
		}
	}
	if link != nil {
		if link.Address != p.Address {
			return invalid()
		}
		var rule *schema.SourceRule
		for _, candidate := range projection.Rules {
			if candidate != nil && candidate.ID == link.RuleID {
				if rule != nil {
					return invalid()
				}
				rule = candidate
			}
		}
		if rule == nil || rule.DeviceID != p.DeviceID || rule.RevisionID == "" {
			return invalid()
		}
		defaults.RuleID = rule.ID
		defaults.RuleOwned = true
		defaults.RuleRevision = rule.RevisionID
		defaults.RuleFormat = rule.DataFormat
	}
	return connector.ResolveReadProfile(*p, defaults)
}
func (s *Service) pointReadProfile(ctx context.Context, p *schema.Point, allowInactive ...bool) (connector.ReadProfile, error) {
	if !modernPoint(p) {
		return connector.ReadProfile{}, nil
	}
	if profiles, ok := ctx.Value(readProfilesKey{}).(map[string]connector.ReadProfile); ok {
		profile, found := profiles[p.ID]
		if !found {
			return connector.ReadProfile{}, connector.ErrReadProfile
		}
		return profile, connector.ValidatePointReadProfile(p, profile)
	}
	if s.workspace == nil {
		return connector.ReadProfile{}, connector.ErrReadProfile
	}
	projection, err := s.workspace.RuntimeProjection(ctx)
	if err != nil {
		return connector.ReadProfile{}, connector.ErrReadProfile
	}
	return projectionReadProfile(projection, p, allowInactive...)
}

// ResolvePointReadProfile supplies the same persisted authority to direct manual
// reads; it does not install a point or create acquisition/readiness proof.
func (s *Service) ResolvePointReadProfile(ctx context.Context, p *schema.Point) (connector.ReadProfile, error) {
	if p != nil && !modernPoint(p) {
		return connector.ReadProfile{}, nil
	}
	if s == nil || p == nil {
		return connector.ReadProfile{}, connector.ErrReadProfile
	}
	return s.pointReadProfile(ctx, p, true)
}

func (s *Service) pointInstallationProfile(ctx context.Context, p *schema.Point) (connector.ReadProfile, error) {
	if p != nil && !p.Enabled {
		return connector.ReadProfile{}, nil
	}
	return s.pointReadProfile(ctx, p)
}

func (s *Service) prepareProjectionReadProfiles(ctx context.Context, projection *workspace.RuntimeProjection, deviceIDs ...string) (context.Context, error) {
	// Capture stop proofs before resolution. A failed candidate cannot leave an
	// older matching acquisition alive, nor may cleanup remove a foreign install.
	proofs := map[string]collector.PointMutationLease{}
	for _, p := range projection.Points {
		if p == nil || !modernPoint(p) || !p.Enabled || (len(deviceIDs) > 0 && !slices.Contains(deviceIDs, p.DeviceID)) {
			continue
		}
		proof, ok := s.CaptureOrdinaryPoint(p.ID)
		if !ok {
			continue
		}
		if guard := sourceProjectionGuard(ctx); guard != nil {
			proof.Lease = guard.expected[p.ID]
		}
		proofs[p.ID] = proof
	}
	prepared, err := prepareProjectionProfiles(ctx, projection, deviceIDs...)
	if err != nil {
		for _, proof := range proofs {
			s.InvalidateOrdinaryPoint(proof)
		}
	}
	return prepared, err
}
