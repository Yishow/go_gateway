package sourcerule

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

type mcGuardRepositorySpy struct {
	*MemoryRepository
	calls []string
}

func (r *mcGuardRepositorySpy) Create(ctx context.Context, rule *schema.SourceRule) error {
	r.calls = append(r.calls, "Create")
	return r.MemoryRepository.Create(ctx, rule)
}

func (r *mcGuardRepositorySpy) Update(ctx context.Context, rule *schema.SourceRule) error {
	r.calls = append(r.calls, "Update")
	return r.MemoryRepository.Update(ctx, rule)
}

func (r *mcGuardRepositorySpy) Delete(ctx context.Context, id string) error {
	r.calls = append(r.calls, "Delete")
	return r.MemoryRepository.Delete(ctx, id)
}

func (r *mcGuardRepositorySpy) CreateLinks(ctx context.Context, links []*schema.SourceRuleLink) error {
	r.calls = append(r.calls, "CreateLinks")
	return r.MemoryRepository.CreateLinks(ctx, links)
}

func (r *mcGuardRepositorySpy) DeleteLinks(ctx context.Context, ruleID string) error {
	r.calls = append(r.calls, "DeleteLinks")
	return r.MemoryRepository.DeleteLinks(ctx, ruleID)
}

func (r *mcGuardRepositorySpy) ReplaceLinks(ctx context.Context, ruleID string, links []*schema.SourceRuleLink) error {
	r.calls = append(r.calls, "ReplaceLinks")
	return r.MemoryRepository.ReplaceLinks(ctx, ruleID, links)
}

func (r *mcGuardRepositorySpy) ReplaceCandidateSnapshots(ctx context.Context, snapshots []*schema.SourceRuleCandidateSnapshot) error {
	r.calls = append(r.calls, "ReplaceCandidateSnapshots")
	return r.MemoryRepository.ReplaceCandidateSnapshots(ctx, snapshots)
}

func (r *mcGuardRepositorySpy) DeleteCandidateSnapshots(ctx context.Context, ruleID, revisionID string) error {
	r.calls = append(r.calls, "DeleteCandidateSnapshots")
	return r.MemoryRepository.DeleteCandidateSnapshots(ctx, ruleID, revisionID)
}

type mcFixtureSnapshot struct {
	rule     *schema.SourceRule
	links    []*schema.SourceRuleLink
	points   []*schema.Point
	tags     []*schema.Tag
	mappings []*schema.Mapping
}

type mcFixture struct {
	svc        *Service
	repo       *MemoryRepository
	deviceID   string
	rule       *schema.SourceRule
	pointSvc   *point.Service
	tagSvc     *tag.Service
	mappingSvc *mapping.Service
}

func newMCFixture(t *testing.T, startAddress string, count int) mcFixture {
	t.Helper()
	ctx := t.Context()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-mc-guard")
	require.NoError(t, err)
	dev.Protocol = schema.ProtocolMC3E
	require.NoError(t, deviceRepo.Update(ctx, dev))

	rule := &schema.SourceRule{
		ID:               "rule-mc-guard",
		DeviceID:         dev.ID,
		StartAddress:     strings.ToUpper(startAddress),
		Count:            count,
		DataType:         schema.DataTypeInt16,
		NamingPrefix:     "LEGACY",
		Enabled:          true,
		SkippedAddresses: "[]",
		RevisionID:       "revision-mc-guard",
	}
	err = repo.Create(ctx, rule)
	require.NoError(t, err)
	legacyAddresses := legacyMCPlanForTest(startAddress, count)
	require.Len(t, legacyAddresses, count)
	links := make([]*schema.SourceRuleLink, 0, count)
	for index, address := range legacyAddresses {
		pointRecord, createErr := pointSvc.Create(ctx, point.CreatePointRequest{
			DeviceID: dev.ID,
			Name:     buildPointName(rule.NamingPrefix, address),
			Address:  address,
			DataType: schema.DataTypeInt16,
			Mode:     schema.PointModeReadOnly,
		})
		require.NoError(t, createErr)
		links = append(links, &schema.SourceRuleLink{
			ID:      fmt.Sprintf("link-mc-guard-%02d", index),
			RuleID:  rule.ID,
			Address: address,
			PointID: pointRecord.ID,
		})
	}
	require.NoError(t, repo.CreateLinks(ctx, links))
	applyRuleManagedLinks(ctx, t, repo, svc, rule.ID)

	links, err = svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, count)
	for _, link := range links {
		require.NotNil(t, link.TagID)
		require.NotNil(t, link.MappingID)
	}

	return mcFixture{
		svc:        svc,
		repo:       repo,
		deviceID:   dev.ID,
		rule:       rule,
		pointSvc:   pointSvc,
		tagSvc:     tagSvc,
		mappingSvc: mappingSvc,
	}
}

func legacyMCPlanForTest(startAddress string, count int) []string {
	start := strings.ToUpper(strings.TrimSpace(startAddress))
	if len(start) < 2 || count < 0 {
		return nil
	}
	value, err := strconv.Atoi(start[1:])
	if err != nil {
		return nil
	}
	addresses := make([]string, 0, count)
	for offset := range count {
		addresses = append(addresses, fmt.Sprintf("%s%d", start[:1], value+offset))
	}
	return addresses
}

func snapshotMCFixture(t *testing.T, fixture mcFixture) mcFixtureSnapshot {
	t.Helper()
	ctx := t.Context()
	rule, err := fixture.svc.GetByID(ctx, fixture.rule.ID)
	require.NoError(t, err)
	links, err := fixture.svc.ListLinks(ctx, fixture.rule.ID)
	require.NoError(t, err)
	points, err := fixture.pointSvc.List(ctx, point.ListFilter{DeviceID: &fixture.deviceID})
	require.NoError(t, err)
	tags, err := fixture.tagSvc.List(ctx, tag.ListFilter{})
	require.NoError(t, err)
	mappings, err := fixture.mappingSvc.List(ctx, mapping.ListFilter{})
	require.NoError(t, err)
	return cloneMCFixtureSnapshot(mcFixtureSnapshot{rule: rule, links: links, points: points, tags: tags, mappings: mappings})
}

func cloneMCFixtureSnapshot(snapshot mcFixtureSnapshot) mcFixtureSnapshot {
	clone := mcFixtureSnapshot{
		rule:     cloneRule(snapshot.rule),
		links:    cloneSourceRuleLinks(snapshot.links),
		points:   make([]*schema.Point, len(snapshot.points)),
		tags:     make([]*schema.Tag, len(snapshot.tags)),
		mappings: make([]*schema.Mapping, len(snapshot.mappings)),
	}
	for index, item := range snapshot.points {
		if item == nil {
			continue
		}
		copyItem := *item
		if item.PollingGroupID != nil {
			value := *item.PollingGroupID
			copyItem.PollingGroupID = &value
		}
		if item.LastReadAt != nil {
			value := *item.LastReadAt
			copyItem.LastReadAt = &value
		}
		if item.LastValue != nil {
			value := *item.LastValue
			copyItem.LastValue = &value
		}
		clone.points[index] = &copyItem
	}
	for index, item := range snapshot.tags {
		if item == nil {
			continue
		}
		copyItem := *item
		clone.tags[index] = &copyItem
	}
	for index, item := range snapshot.mappings {
		if item == nil {
			continue
		}
		copyItem := *item
		clone.mappings[index] = &copyItem
	}
	return clone
}

func assertMCFixtureSnapshotEqual(t *testing.T, before, after mcFixtureSnapshot) {
	t.Helper()
	require.Equal(t, before.rule, after.rule)
	require.ElementsMatch(t, before.links, after.links)
	require.ElementsMatch(t, before.points, after.points)
	require.ElementsMatch(t, before.tags, after.tags)
	require.ElementsMatch(t, before.mappings, after.mappings)
}

type mcLinkIdentity struct {
	ID        string
	PointID   string
	Address   string
	TagID     string
	MappingID string
}

func mcLinkIdentities(links []*schema.SourceRuleLink) map[string]mcLinkIdentity {
	result := make(map[string]mcLinkIdentity, len(links))
	for _, link := range links {
		if link == nil {
			continue
		}
		result[link.ID] = mcLinkIdentity{
			ID:        link.ID,
			PointID:   link.PointID,
			Address:   link.Address,
			TagID:     optionalMCID(link.TagID),
			MappingID: optionalMCID(link.MappingID),
		}
	}
	return result
}

type mcPointIdentity struct {
	ID       string
	DeviceID string
	Name     string
	Address  string
	Function string
	DataType schema.DataType
	Enabled  bool
}

func mcPointIdentities(points []*schema.Point) map[string]mcPointIdentity {
	result := make(map[string]mcPointIdentity, len(points))
	for _, item := range points {
		if item == nil {
			continue
		}
		result[item.ID] = mcPointIdentity{ID: item.ID, DeviceID: item.DeviceID, Name: item.Name, Address: item.Address, Function: item.Function, DataType: item.DataType, Enabled: item.Enabled}
	}
	return result
}

type mcTagIdentity struct {
	ID       string
	Key      string
	DataType schema.DataType
	Labels   string
}

func mcTagIdentities(tags []*schema.Tag) map[string]mcTagIdentity {
	result := make(map[string]mcTagIdentity, len(tags))
	for _, item := range tags {
		if item == nil {
			continue
		}
		result[item.ID] = mcTagIdentity{ID: item.ID, Key: item.Key, DataType: item.DataType, Labels: item.Labels}
	}
	return result
}

type mcMappingIdentity struct {
	ID                   string
	PointID              string
	TagID                string
	TransformPipeline    string
	Status               schema.MappingStatus
	RuleCandidateID      string
	ProposedSignature    string
	LastAppliedSignature string
	BlockingReason       string
	Enabled              bool
}

func mcMappingIdentities(mappings []*schema.Mapping) map[string]mcMappingIdentity {
	result := make(map[string]mcMappingIdentity, len(mappings))
	for _, item := range mappings {
		if item == nil {
			continue
		}
		result[item.ID] = mcMappingIdentity{ID: item.ID, PointID: item.PointID, TagID: item.TagID, TransformPipeline: item.TransformPipeline, Status: item.Status, RuleCandidateID: item.RuleCandidateID, ProposedSignature: item.ProposedSignature, LastAppliedSignature: item.LastAppliedSignature, BlockingReason: item.BlockingReason, Enabled: item.Enabled}
	}
	return result
}

func optionalMCID(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func TestMCAddressPlan_CorrectedRadixVectors(t *testing.T) {
	cases := []struct {
		name     string
		start    string
		count    int
		expected []string
		wantErr  bool
	}{
		{
			name:  "W hexadecimal",
			start: "W0",
			count: 17,
			expected: []string{
				"W0", "W1", "W2", "W3", "W4", "W5", "W6", "W7", "W8", "W9",
				"WA", "WB", "WC", "WD", "WE", "WF", "W10",
			},
		},
		{
			name:  "X octal",
			start: "X0",
			count: 17,
			expected: []string{
				"X0", "X1", "X2", "X3", "X4", "X5", "X6", "X7", "X10",
				"X11", "X12", "X13", "X14", "X15", "X16", "X17", "X20",
			},
		},
		{
			name:     "D decimal",
			start:    "D0",
			count:    2,
			expected: []string{"D0", "D1"},
		},
		{name: "X8 is not octal", start: "X8", count: 1, wantErr: true},
		{name: "YA is not octal", start: "YA", count: 1, wantErr: true},
		{name: "WG is not hexadecimal", start: "WG", count: 1, wantErr: true},
		{name: "Z0 is not an MC area", start: "Z0", count: 1, wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rule := &schema.SourceRule{StartAddress: tc.start, Count: tc.count, DataType: schema.DataTypeInt16, SkippedAddresses: "[]"}
			got, err := correctedMCAddressPlan(rule)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected, got)
		})
	}
}

func TestServiceUpdate_MCRadixLegacyRangeConflictPreservesCompleteFixture(t *testing.T) {
	corrected := []string{"W0", "W1", "W2", "W3", "W4", "W5", "W6", "W7", "W8", "W9", "WA", "WB", "WC", "WD", "WE", "WF", "W10"}

	cases := map[string]func() UpdateRuleRequest{
		"name patch": func() UpdateRuleRequest {
			name := "RENAMED"
			return UpdateRuleRequest{NamingPrefix: &name}
		},
		"scale patch": func() UpdateRuleRequest {
			scale := 2.0
			return UpdateRuleRequest{ScaleMultiplier: &scale}
		},
		"range patch": func() UpdateRuleRequest {
			start := "W1"
			count := 17
			return UpdateRuleRequest{StartAddress: &start, Count: &count}
		},
	}
	for name, request := range cases {
		t.Run(name, func(t *testing.T) {
			fixture := newMCFixture(t, "W0", 17)
			links, err := fixture.svc.ListLinks(t.Context(), fixture.rule.ID)
			require.NoError(t, err)
			stored := make([]string, 0, len(links))
			for _, link := range links {
				stored = append(stored, link.Address)
			}
			require.NotEqual(t, corrected, stored, "the fixture must represent the legacy decimal W range")

			before := snapshotMCFixture(t, fixture)
			spy := &mcGuardRepositorySpy{MemoryRepository: fixture.repo}
			fixture.svc.repo = spy
			_, err = fixture.svc.Update(t.Context(), fixture.rule.ID, request())
			require.Error(t, err)
			require.ErrorIs(t, err, ErrValidation)
			require.Contains(t, err.Error(), "mc_address_plan_conflict")
			require.Contains(t, err.Error(), fixture.rule.ID)
			require.Contains(t, err.Error(), "W0")
			require.Empty(t, spy.calls)
			assertMCFixtureSnapshotEqual(t, before, snapshotMCFixture(t, fixture))
		})
	}
}

func TestServiceUpdate_MCRadixLegacyInvalidX8ConflictPreservesFixture(t *testing.T) {
	fixture := newMCFixture(t, "X8", 1)
	before := snapshotMCFixture(t, fixture)
	spy := &mcGuardRepositorySpy{MemoryRepository: fixture.repo}
	fixture.svc.repo = spy
	name := "RENAMED"
	_, err := fixture.svc.Update(t.Context(), fixture.rule.ID, UpdateRuleRequest{NamingPrefix: &name})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrValidation)
	require.Contains(t, err.Error(), "mc_address_plan_conflict")
	require.Empty(t, spy.calls)
	assertMCFixtureSnapshotEqual(t, before, snapshotMCFixture(t, fixture))
}

func TestServiceUpdate_MCAddressPlanCompatibleRuleStillUpdatesNormally(t *testing.T) {
	fixture := newMCFixture(t, "W0", 1)
	name := "CURRENT"
	updated, err := fixture.svc.Update(t.Context(), fixture.rule.ID, UpdateRuleRequest{NamingPrefix: &name})
	require.NoError(t, err)
	require.Equal(t, "CURRENT", updated.NamingPrefix)
}

func TestServiceSyncDerivedPointState_MCRadixLegacyLinksRemainAuthoritative(t *testing.T) {
	fixture := newMCFixture(t, "W0", 17)
	before := snapshotMCFixture(t, fixture)

	require.NoError(t, fixture.svc.SyncDerivedPointState(t.Context()))
	after := snapshotMCFixture(t, fixture)
	require.Equal(t, before.rule.ID, after.rule.ID)
	require.Equal(t, before.rule.StartAddress, after.rule.StartAddress)
	require.Equal(t, before.rule.Count, after.rule.Count)
	require.Equal(t, before.rule.DataType, after.rule.DataType)
	require.Equal(t, before.rule.RevisionID, after.rule.RevisionID)
	require.Equal(t, mcLinkIdentities(before.links), mcLinkIdentities(after.links))
	require.Equal(t, mcPointIdentities(before.points), mcPointIdentities(after.points))
	require.Equal(t, mcTagIdentities(before.tags), mcTagIdentities(after.tags))
	require.Equal(t, mcMappingIdentities(before.mappings), mcMappingIdentities(after.mappings))
}
