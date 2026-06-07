package workspace

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

var ErrRuntimeProjectionUnavailable = errors.New("workspace runtime projection unavailable")

type runtimeProjectionDeviceService interface {
	GetByID(ctx context.Context, id string) (*schema.Device, error)
}

type runtimeProjectionSourceRuleService interface {
	ListByDeviceIDs(ctx context.Context, deviceIDs []string) ([]*schema.SourceRule, error)
	ListLinks(ctx context.Context, ruleID string) ([]*schema.SourceRuleLink, error)
}

type runtimeProjectionPointService interface {
	List(ctx context.Context, filter point.ListFilter) ([]*schema.Point, error)
}

type runtimeProjectionMappingService interface {
	List(ctx context.Context, filter mapping.ListFilter) ([]*schema.Mapping, error)
}

type runtimeProjectionTagService interface {
	GetByID(ctx context.Context, id string) (*schema.Tag, error)
}

type runtimeProjectionConnectorService interface {
	GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error)
}

type runtimeProjectionTargetService interface {
	List(ctx context.Context, filter dbtarget.TargetMappingListFilter) ([]*schema.DatabaseTargetMapping, error)
}

type runtimeProjectionPollingGroupService interface {
	List(ctx context.Context) ([]*schema.PollingGroup, error)
}

// RuntimeProjectionAlignment identifies whether runtime is aligned to persisted workspace state.
type RuntimeProjectionAlignment string

const (
	RuntimeProjectionAlignmentAligned         RuntimeProjectionAlignment = "aligned"
	RuntimeProjectionAlignmentDeferred        RuntimeProjectionAlignment = "deferred"
	RuntimeProjectionAlignmentRestartRequired RuntimeProjectionAlignment = "restart-required"
	RuntimeProjectionAlignmentStale           RuntimeProjectionAlignment = "stale"
	RuntimeProjectionAlignmentDegraded        RuntimeProjectionAlignment = "degraded"
)

// RuntimeProjection is the persisted workspace scope consumed by runtime activation and restart.
type RuntimeProjection struct {
	WorkspaceID       string                          `json:"workspace_id"`
	Version           string                          `json:"version"`
	Alignment         RuntimeProjectionAlignment      `json:"alignment"`
	DeviceIDs         []string                        `json:"device_ids"`
	Devices           []*schema.Device                `json:"devices"`
	Rules             []*schema.SourceRule            `json:"rules"`
	RuleLinks         []*schema.SourceRuleLink        `json:"rule_links"`
	Points            []*schema.Point                 `json:"points"`
	PollingGroups     []*schema.PollingGroup          `json:"polling_groups"`
	Mappings          []*schema.Mapping               `json:"mappings"`
	Tags              []*schema.Tag                   `json:"tags"`
	DatabaseConnector *schema.DatabaseConnector       `json:"database_connector,omitempty"`
	DatabaseTargets   []*schema.DatabaseTargetMapping `json:"database_targets"`
}

// WithRuntimeProjectionServices configures persisted readers used to build runtime projection.
func (s *Service) WithRuntimeProjectionServices(
	deviceSvc runtimeProjectionDeviceService,
	ruleSvc runtimeProjectionSourceRuleService,
	pointSvc runtimeProjectionPointService,
	mappingSvc runtimeProjectionMappingService,
	tagSvc runtimeProjectionTagService,
	connectorSvc runtimeProjectionConnectorService,
	targetSvc runtimeProjectionTargetService,
	groupSvc runtimeProjectionPollingGroupService,
) *Service {
	s.projectionDevices = deviceSvc
	s.projectionRules = ruleSvc
	s.projectionPoints = pointSvc
	s.projectionMappings = mappingSvc
	s.projectionTags = tagSvc
	s.projectionConnectors = connectorSvc
	s.projectionTargets = targetSvc
	s.projectionGroups = groupSvc
	return s
}

// RuntimeProjection rebuilds the runtime scope from the persisted workspace.
func (s *Service) RuntimeProjection(ctx context.Context) (*RuntimeProjection, error) {
	if err := s.ensureRuntimeProjectionReaders(); err != nil {
		return nil, err
	}

	record, err := s.GetOrCreate(ctx)
	if err != nil {
		return nil, err
	}

	deviceOrder := indexStrings(record.OrderedDeviceIDs)
	devices, err := s.runtimeProjectionDevices(ctx, record.OrderedDeviceIDs)
	if err != nil {
		return nil, err
	}
	rules, links, err := s.runtimeProjectionRules(ctx, record.OrderedDeviceIDs, deviceOrder)
	if err != nil {
		return nil, err
	}
	points, err := s.runtimeProjectionPoints(ctx, record.OrderedDeviceIDs, links)
	if err != nil {
		return nil, err
	}
	mappings, tags, err := s.runtimeProjectionMappings(ctx, points, links)
	if err != nil {
		return nil, err
	}
	connector, targets, err := s.runtimeProjectionDatabaseTargets(ctx, record.DatabaseConnectorID, tags)
	if err != nil {
		return nil, err
	}
	groups, err := s.runtimeProjectionPollingGroups(ctx, points)
	if err != nil {
		return nil, err
	}

	projection := &RuntimeProjection{
		WorkspaceID:       record.ID,
		Alignment:         RuntimeProjectionAlignmentAligned,
		DeviceIDs:         append([]string{}, record.OrderedDeviceIDs...),
		Devices:           devices,
		Rules:             rules,
		RuleLinks:         links,
		Points:            points,
		PollingGroups:     groups,
		Mappings:          mappings,
		Tags:              tags,
		DatabaseConnector: connector,
		DatabaseTargets:   targets,
	}
	projection.Version, err = runtimeProjectionVersion(projection)
	if err != nil {
		return nil, err
	}
	return projection, nil
}

func (s *Service) ensureRuntimeProjectionReaders() error {
	if s.projectionDevices == nil ||
		s.projectionRules == nil ||
		s.projectionPoints == nil ||
		s.projectionMappings == nil ||
		s.projectionTags == nil ||
		s.projectionGroups == nil {
		return ErrRuntimeProjectionUnavailable
	}
	return nil
}

func (s *Service) runtimeProjectionDevices(ctx context.Context, deviceIDs []string) ([]*schema.Device, error) {
	devices := make([]*schema.Device, 0, len(deviceIDs))
	for _, deviceID := range deviceIDs {
		deviceRecord, err := s.projectionDevices.GetByID(ctx, deviceID)
		if err != nil {
			return nil, fmt.Errorf("read workspace projection device %s: %w", deviceID, err)
		}
		devices = append(devices, deviceRecord)
	}
	return devices, nil
}

func (s *Service) runtimeProjectionRules(
	ctx context.Context,
	deviceIDs []string,
	deviceOrder map[string]int,
) ([]*schema.SourceRule, []*schema.SourceRuleLink, error) {
	rules, err := s.projectionRules.ListByDeviceIDs(ctx, deviceIDs)
	if err != nil {
		return nil, nil, fmt.Errorf("read workspace projection rules: %w", err)
	}
	sort.SliceStable(rules, func(i, j int) bool {
		leftDevice := deviceOrder[rules[i].DeviceID]
		rightDevice := deviceOrder[rules[j].DeviceID]
		if leftDevice != rightDevice {
			return leftDevice < rightDevice
		}
		return rules[i].ID < rules[j].ID
	})

	links := make([]*schema.SourceRuleLink, 0)
	for _, rule := range rules {
		ruleLinks, err := s.projectionRules.ListLinks(ctx, rule.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("read workspace projection rule links %s: %w", rule.ID, err)
		}
		sort.SliceStable(ruleLinks, func(i, j int) bool {
			if ruleLinks[i].Address != ruleLinks[j].Address {
				return ruleLinks[i].Address < ruleLinks[j].Address
			}
			return ruleLinks[i].ID < ruleLinks[j].ID
		})
		links = append(links, ruleLinks...)
	}
	return rules, links, nil
}

func (s *Service) runtimeProjectionPoints(ctx context.Context, deviceIDs []string, links []*schema.SourceRuleLink) ([]*schema.Point, error) {
	livePointIDs := runtimeProjectionLinkPointIDs(links)
	points := make([]*schema.Point, 0)
	for _, deviceID := range deviceIDs {
		devicePoints, err := s.projectionPoints.List(ctx, point.ListFilter{DeviceID: &deviceID, Limit: 100000})
		if err != nil {
			return nil, fmt.Errorf("read workspace projection points %s: %w", deviceID, err)
		}
		sort.SliceStable(devicePoints, func(i, j int) bool {
			return devicePoints[i].ID < devicePoints[j].ID
		})
		for _, pointRecord := range devicePoints {
			if _, live := livePointIDs[pointRecord.ID]; !live {
				continue
			}
			points = append(points, pointRecord)
		}
	}
	return points, nil
}

func (s *Service) runtimeProjectionMappings(ctx context.Context, points []*schema.Point, links []*schema.SourceRuleLink) ([]*schema.Mapping, []*schema.Tag, error) {
	enabled := true
	liveMappingIDs := runtimeProjectionLinkMappingIDs(links)
	mappings := make([]*schema.Mapping, 0)
	tagIDs := make([]string, 0)
	seenTags := make(map[string]struct{})

	for _, pointRecord := range points {
		pointMappings, err := s.projectionMappings.List(ctx, mapping.ListFilter{
			PointID: &pointRecord.ID,
			Enabled: &enabled,
			Limit:   100000,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("read workspace projection mappings %s: %w", pointRecord.ID, err)
		}
		sort.SliceStable(pointMappings, func(i, j int) bool {
			return pointMappings[i].ID < pointMappings[j].ID
		})
		for _, mappingRecord := range pointMappings {
			if _, live := liveMappingIDs[mappingRecord.ID]; !live {
				continue
			}
			mappings = append(mappings, mappingRecord)
			if _, ok := seenTags[mappingRecord.TagID]; ok {
				continue
			}
			seenTags[mappingRecord.TagID] = struct{}{}
			tagIDs = append(tagIDs, mappingRecord.TagID)
		}
	}

	tags := make([]*schema.Tag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		tagRecord, err := s.projectionTags.GetByID(ctx, tagID)
		if err != nil {
			return nil, nil, fmt.Errorf("read workspace projection tag %s: %w", tagID, err)
		}
		tags = append(tags, tagRecord)
	}
	return mappings, tags, nil
}

func indexStrings(values []string) map[string]int {
	index := make(map[string]int, len(values))
	for i, value := range values {
		index[value] = i
	}
	return index
}
