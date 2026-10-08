package workspace

import (
	"context"
	"slices"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

type runtimeProjectionSummary struct {
	DeviceIDs      []string
	RuleIDs        []string
	RuleLinkIDs    []string
	PointIDs       []string
	MappingIDs     []string
	TagIDs         []string
	DatabaseTarget []string
	PollingGroups  []string
}

func summarizeRuntimeProjection(projection *RuntimeProjection) runtimeProjectionSummary {
	return runtimeProjectionSummary{
		DeviceIDs:      slices.Clone(projection.DeviceIDs),
		RuleIDs:        sourceRuleIDs(projection.Rules),
		RuleLinkIDs:    sourceRuleLinkIDs(projection.RuleLinks),
		PointIDs:       pointIDs(projection.Points),
		MappingIDs:     mappingIDs(projection.Mappings),
		TagIDs:         tagIDs(projection.Tags),
		DatabaseTarget: databaseTargetIDs(projection.DatabaseTargets),
		PollingGroups:  pollingGroupIDs(projection.PollingGroups),
	}
}

func idSlice[T any](records []*T, idFn func(*T) string) []string {
	if records == nil {
		return []string{}
	}
	ids := make([]string, len(records))
	for i, r := range records {
		ids[i] = idFn(r)
	}
	return ids
}

func sourceRuleIDs(records []*schema.SourceRule) []string {
	return idSlice(records, func(r *schema.SourceRule) string { return r.ID })
}

func sourceRuleLinkIDs(records []*schema.SourceRuleLink) []string {
	return idSlice(records, func(r *schema.SourceRuleLink) string { return r.ID })
}

func pointIDs(records []*schema.Point) []string {
	return idSlice(records, func(r *schema.Point) string { return r.ID })
}

func mappingIDs(records []*schema.Mapping) []string {
	return idSlice(records, func(r *schema.Mapping) string { return r.ID })
}

func tagIDs(records []*schema.Tag) []string {
	return idSlice(records, func(r *schema.Tag) string { return r.ID })
}

func databaseTargetIDs(records []*schema.DatabaseTargetMapping) []string {
	return idSlice(records, func(r *schema.DatabaseTargetMapping) string { return r.ID })
}

func pollingGroupIDs(records []*schema.PollingGroup) []string {
	return idSlice(records, func(r *schema.PollingGroup) string { return r.ID })
}

type runtimeProjectionDeviceStub struct {
	records map[string]*schema.Device
}

func (s *runtimeProjectionDeviceStub) GetByID(_ context.Context, id string) (*schema.Device, error) {
	record := *s.records[id]
	return &record, nil
}

type runtimeProjectionRuleStub struct {
	rules []*schema.SourceRule
	links map[string][]*schema.SourceRuleLink
}

func (s *runtimeProjectionRuleStub) ListByDeviceIDs(_ context.Context, deviceIDs []string) ([]*schema.SourceRule, error) {
	deviceSet := make(map[string]struct{}, len(deviceIDs))
	for _, deviceID := range deviceIDs {
		deviceSet[deviceID] = struct{}{}
	}
	result := make([]*schema.SourceRule, 0)
	for _, rule := range s.rules {
		if _, ok := deviceSet[rule.DeviceID]; ok {
			record := *rule
			result = append(result, &record)
		}
	}
	return result, nil
}

func (s *runtimeProjectionRuleStub) ListLinks(_ context.Context, ruleID string) ([]*schema.SourceRuleLink, error) {
	links := s.links[ruleID]
	result := make([]*schema.SourceRuleLink, 0, len(links))
	for _, link := range links {
		record := *link
		result = append(result, &record)
	}
	return result, nil
}

type runtimeProjectionPointStub struct {
	records []*schema.Point
}

func (s *runtimeProjectionPointStub) List(_ context.Context, filter point.ListFilter) ([]*schema.Point, error) {
	result := make([]*schema.Point, 0)
	for _, record := range s.records {
		if filter.DeviceID != nil && record.DeviceID != *filter.DeviceID {
			continue
		}
		copied := *record
		result = append(result, &copied)
	}
	return result, nil
}

type runtimeProjectionMappingStub struct {
	records []*schema.Mapping
}

func (s *runtimeProjectionMappingStub) List(_ context.Context, filter mapping.ListFilter) ([]*schema.Mapping, error) {
	result := make([]*schema.Mapping, 0)
	for _, record := range s.records {
		if filter.PointID != nil && record.PointID != *filter.PointID {
			continue
		}
		if filter.Enabled != nil && record.Enabled != *filter.Enabled {
			continue
		}
		copied := *record
		result = append(result, &copied)
	}
	return result, nil
}

type runtimeProjectionTagStub struct {
	records map[string]*schema.Tag
}

func (s *runtimeProjectionTagStub) GetByID(_ context.Context, id string) (*schema.Tag, error) {
	record := *s.records[id]
	return &record, nil
}

type runtimeProjectionTargetStub struct {
	records []*schema.DatabaseTargetMapping
}

func (s *runtimeProjectionTargetStub) List(_ context.Context, filter dbtarget.TargetMappingListFilter) ([]*schema.DatabaseTargetMapping, error) {
	result := make([]*schema.DatabaseTargetMapping, 0)
	for _, record := range s.records {
		if filter.ConnectorID != nil && record.ConnectorID != *filter.ConnectorID {
			continue
		}
		if filter.TagID != nil && record.TagID != *filter.TagID {
			continue
		}
		if filter.Enabled != nil && record.Enabled != *filter.Enabled {
			continue
		}
		copied := *record
		result = append(result, &copied)
	}
	return result, nil
}

type runtimeProjectionGroupStub struct {
	records []*schema.PollingGroup
}

func (s *runtimeProjectionGroupStub) List(context.Context) ([]*schema.PollingGroup, error) {
	result := make([]*schema.PollingGroup, 0, len(s.records))
	for _, record := range s.records {
		copied := *record
		result = append(result, &copied)
	}
	return result, nil
}
