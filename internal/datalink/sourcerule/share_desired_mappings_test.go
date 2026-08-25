package sourcerule

import (
	"context"
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildDesiredShareMappingsForDevices_UsesEachRuleCandidateRegister(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)
	dev, err := seedActiveDevice(ctx, deviceRepo, "device-share-plan")
	require.NoError(t, err)

	ruleA, err := svc.Create(ctx, CreateRuleRequest{ID: "rule-share-plan-a", DeviceID: dev.ID, StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt16, NamingPrefix: "A", Enabled: true, ShareEnabled: true, ShareStartRegister: intPtr(40001), ShareStride: intPtr(2)})
	require.NoError(t, err)
	ruleB, err := svc.Create(ctx, CreateRuleRequest{ID: "rule-share-plan-b", DeviceID: dev.ID, StartAddress: "40011", Count: 1, DataType: schema.DataTypeInt32, NamingPrefix: "B", Enabled: true, ShareEnabled: true, ShareStartRegister: intPtr(40011), ShareStride: intPtr(3)})
	require.NoError(t, err)
	require.Len(t, applyRuleManagedLinks(t, ctx, repo, svc, ruleA.ID), 1)
	require.Len(t, applyRuleManagedLinks(t, ctx, repo, svc, ruleB.ID), 1)

	linksA, err := svc.ListLinks(ctx, ruleA.ID)
	require.NoError(t, err)
	linksB, err := svc.ListLinks(ctx, ruleB.ID)
	require.NoError(t, err)
	require.NotNil(t, linksA[0].TagID)
	require.NotNil(t, linksB[0].TagID)
	require.NoError(t, svc.persistCandidateSnapshots(ctx, ruleA, linksA))
	require.NoError(t, svc.persistCandidateSnapshots(ctx, ruleB, linksB))

	candidateA, _ := localModbusSingleCandidateFromRule(t, ctx, svc, ruleA.ID)
	candidateB, _ := localModbusSingleCandidateFromRule(t, ctx, svc, ruleB.ID)
	// Candidate registers are canonical zero-based coordinates. Rule B's
	// durable 40011 start therefore produces zero-based 10, not 20.
	registerA, registerB := uint16(1), uint16(10)
	candidateA.Register, candidateB.Register = &registerA, &registerB
	candidateA.Status = schema.SourceRuleLocalModbusOutputStatusReady
	candidateB.Status = schema.SourceRuleLocalModbusOutputStatusReady
	setReadyCandidateSnapshot(t, repo, ruleA, candidateA)
	setReadyCandidateSnapshot(t, repo, ruleB, candidateB)

	desired, err := svc.BuildDesiredShareMappingsForDevices(ctx, "workspace-share-plan", modbusshare.Settings{Enabled: true, CapacityRegisters: 100}, []string{dev.ID})
	require.NoError(t, err)
	require.Len(t, desired, 2)
	byRule := map[string]modbusshare.DesiredMapping{desired[0].SourceRuleID: desired[0], desired[1].SourceRuleID: desired[1]}
	assert.Equal(t, uint16(1), byRule[ruleA.ID].ZeroBasedRegister)
	assert.Equal(t, uint16(10), byRule[ruleB.ID].ZeroBasedRegister)
	assert.Equal(t, 1, byRule[ruleA.ID].SpanRegisters)
	assert.Equal(t, 2, byRule[ruleB.ID].SpanRegisters)
	assert.Equal(t, 100, byRule[ruleA.ID].CapacityRegisters)
	assert.NotEmpty(t, byRule[ruleA.ID].MappingID)
	assert.NotEmpty(t, byRule[ruleA.ID].TagKey)
	assert.NotEmpty(t, byRule[ruleA.ID].DisplayName)
	assert.NotEqual(t, byRule[ruleA.ID].SourceRuleRevision, byRule[ruleB.ID].SourceRuleRevision)
}

func TestShareDesiredMappingOwnershipCheckerIgnoresCanonicalProofMetadata(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)
	dev, err := seedActiveDevice(ctx, deviceRepo, "device-ownership-proof")
	require.NoError(t, err)
	rule, err := svc.Create(ctx, CreateRuleRequest{ID: "rule-ownership-proof", DeviceID: dev.ID, StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt16, NamingPrefix: "PROOF", Enabled: true, ShareEnabled: true, ShareStartRegister: intPtr(40001), ShareStride: intPtr(1)})
	require.NoError(t, err)
	require.Len(t, applyRuleManagedLinks(t, ctx, repo, svc, rule.ID), 1)
	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.NoError(t, svc.persistCandidateSnapshots(ctx, rule, links))
	candidate, _ := localModbusSingleCandidateFromRule(t, ctx, svc, rule.ID)
	register := uint16(0)
	candidate.Register = &register
	candidate.Status = schema.SourceRuleLocalModbusOutputStatusReady
	setReadyCandidateSnapshot(t, repo, rule, candidate)

	desired, err := svc.BuildDesiredShareMappingsForDevices(ctx, "workspace-proof", modbusshare.Settings{Enabled: true, CapacityRegisters: 64}, []string{dev.ID})
	require.NoError(t, err)
	require.Len(t, desired, 1)
	desired[0].OwnershipProof = &modbusshare.OwnershipProof{Verified: true, WorkspaceID: "workspace-proof", SourceRuleID: rule.ID, SourceRuleRevision: rule.RevisionID, Basis: "persisted workspace/source-rule/tag relationship"}

	checker := NewShareDesiredMappingOwnershipChecker(candidateScopeWorkspaceStub{id: "workspace-proof", devices: []string{dev.ID}}, svc, tagSvc, func() modbusshare.Settings {
		return modbusshare.Settings{Enabled: true, CapacityRegisters: 64}
	}, mappingSvc)
	require.NoError(t, checker(ctx, desired[0]))
}

func setReadyCandidateSnapshot(t *testing.T, repo *MemoryRepository, rule *schema.SourceRule, candidate schema.SourceRuleLocalModbusOutputCandidate) {
	t.Helper()
	snapshots, err := repo.ListCandidateSnapshots(context.Background(), rule.ID, rule.RevisionID)
	require.NoError(t, err)
	payload, err := json.Marshal(struct {
		Candidates []schema.SourceRuleLocalModbusOutputCandidate `json:"candidates"`
	}{Candidates: []schema.SourceRuleLocalModbusOutputCandidate{candidate}})
	require.NoError(t, err)
	for _, snapshot := range snapshots {
		if snapshot.CandidateType != schema.SourceRuleCandidateTypeLocalModbusOutputs {
			continue
		}
		copySnapshot := *snapshot
		copySnapshot.Payload = string(payload)
		copySnapshot.Status = schema.SourceRuleCandidateStatusReady
		copySnapshot.Reason = ""
		next := make([]*schema.SourceRuleCandidateSnapshot, 0, len(snapshots))
		for _, item := range snapshots {
			if item.CandidateType == schema.SourceRuleCandidateTypeLocalModbusOutputs {
				next = append(next, &copySnapshot)
			} else {
				next = append(next, item)
			}
		}
		require.NoError(t, repo.ReplaceCandidateSnapshots(context.Background(), next))
		return
	}
	t.Fatalf("local modbus snapshot not found for rule %s", rule.ID)
}

func intPtr(value int) *int { return &value }
