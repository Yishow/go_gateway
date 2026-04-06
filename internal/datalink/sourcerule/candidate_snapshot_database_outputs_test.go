package sourcerule

import (
	"context"
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type databaseOutputSnapshotPayload struct {
	Candidates []map[string]any `json:"candidates"`
}

func TestService_CandidateSnapshots_BuildDatabaseOutputsFromEffectiveTagReviewState(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-db-effective")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-db-effective",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        3,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	tagCandidates := tagCandidatesFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	require.Len(t, tagCandidates, 3)

	overrideTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "factory.db.override",
		DisplayName: "DB Override",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)

	_, err = svc.UpsertTagReviewDecision(ctx, rule.ID, UpsertTagReviewDecisionRequest{
		CandidateID: tagCandidates[0].ID,
		Action:      schema.SourceRuleTagReviewDecisionActionRename,
		TagKey:      "factory.db.renamed",
	})
	require.NoError(t, err)
	_, err = svc.UpsertTagReviewDecision(ctx, rule.ID, UpsertTagReviewDecisionRequest{
		CandidateID: tagCandidates[1].ID,
		Action:      schema.SourceRuleTagReviewDecisionActionSkip,
	})
	require.NoError(t, err)
	_, err = svc.UpsertTagReviewDecision(ctx, rule.ID, UpsertTagReviewDecisionRequest{
		CandidateID:   tagCandidates[2].ID,
		Action:        schema.SourceRuleTagReviewDecisionActionOverride,
		OverrideTagID: overrideTag.ID,
	})
	require.NoError(t, err)

	snapshots, err := svc.ListCandidateSnapshots(ctx, rule.ID)
	require.NoError(t, err)
	databaseCandidates, databaseSnapshot := databaseCandidatesFromSnapshots(t, snapshots)
	assert.Equal(t, schema.SourceRuleCandidateStatusReady, databaseSnapshot.Status)
	assert.Empty(t, databaseSnapshot.Reason)
	require.Len(t, databaseCandidates, 2)

	candidateByAddress := make(map[string]map[string]any, len(databaseCandidates))
	for _, candidate := range databaseCandidates {
		candidateByAddress[candidate["address"].(string)] = candidate
	}

	require.Contains(t, candidateByAddress, "40001")
	assert.Equal(t, "factory.db.renamed", candidateByAddress["40001"]["tag_key"])
	assert.Nil(t, candidateByAddress["40001"]["tag_id"])

	require.Contains(t, candidateByAddress, "40003")
	assert.Equal(t, "factory.db.override", candidateByAddress["40003"]["tag_key"])
	assert.Equal(t, overrideTag.ID, candidateByAddress["40003"]["tag_id"])
}

func databaseCandidatesFromSnapshots(
	t *testing.T,
	snapshots []*schema.SourceRuleCandidateSnapshot,
) ([]map[string]any, *schema.SourceRuleCandidateSnapshot) {
	t.Helper()

	for _, snapshot := range snapshots {
		if snapshot == nil || snapshot.CandidateType != schema.SourceRuleCandidateTypeDatabaseOutputs {
			continue
		}
		var payload databaseOutputSnapshotPayload
		require.NoError(t, json.Unmarshal([]byte(snapshot.Payload), &payload))
		return payload.Candidates, snapshot
	}

	t.Fatalf("database output snapshot not found")
	return nil, nil
}
