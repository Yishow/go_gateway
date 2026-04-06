package sourcerule

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_CandidateSnapshots_PersistDatabaseOutputsPerRevision(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-db-revision-generation")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-db-revision-generation",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	initialRevision := rule.RevisionID
	initialCandidates := databaseOutputCandidatesFromMemoryRepo(t, repo, rule.ID, initialRevision)
	require.Len(t, initialCandidates, 1)
	assert.Equal(t, []string{"40001"}, databaseCandidateAddresses(initialCandidates))

	count := 2
	updatedRule, err := svc.Update(ctx, rule.ID, UpdateRuleRequest{Count: &count})
	require.NoError(t, err)
	require.NotEqual(t, initialRevision, updatedRule.RevisionID)

	historicalCandidates := databaseOutputCandidatesFromMemoryRepo(t, repo, rule.ID, initialRevision)
	require.Len(t, historicalCandidates, 1)
	assert.Equal(t, []string{"40001"}, databaseCandidateAddresses(historicalCandidates))

	updatedCandidates := databaseOutputCandidatesFromMemoryRepo(t, repo, rule.ID, updatedRule.RevisionID)
	require.Len(t, updatedCandidates, 2)
	assert.Equal(t, []string{"40001", "40002"}, databaseCandidateAddresses(updatedCandidates))

	latestSnapshots, err := svc.ListCandidateSnapshots(ctx, rule.ID)
	require.NoError(t, err)
	latestCandidates, latestSnapshot := databaseCandidatesFromSnapshots(t, latestSnapshots)
	assert.Equal(t, updatedRule.RevisionID, latestSnapshot.RevisionID)
	require.Len(t, latestCandidates, 2)
	assert.Equal(t, []string{"40001", "40002"}, databaseCandidateAddressesFromPayload(latestCandidates))
}

func databaseCandidateAddresses(candidates []schema.SourceRuleDatabaseOutputCandidate) []string {
	addresses := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		addresses = append(addresses, candidate.Address)
	}
	return addresses
}

func databaseCandidateAddressesFromPayload(candidates []map[string]any) []string {
	addresses := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		addresses = append(addresses, candidate["address"].(string))
	}
	return addresses
}
