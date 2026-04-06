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

func TestService_ApplyDatabaseOutputCandidates_UsesCurrentRevisionEvenWhenCandidateIDIsStable(t *testing.T) {
	t.Parallel()

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

	overrideTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "factory.db.revision",
		DisplayName: "DB Revision",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)
	svc.SetDatabaseTargetMappingReader(DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
		return []*schema.DatabaseTargetMapping{
			{
				ID:          "db-map-revision",
				TagID:       overrideTag.ID,
				ConnectorID: "connector-revision",
				TableSchema: "public",
				TableName:   "measurements",
				ColumnName:  "line_a",
				WriteMode:   schema.DatabaseWriteModeInsert,
			},
		}, nil
	}))

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-db-apply-revision")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-db-apply-revision",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	initialRevision := rule.RevisionID
	tagCandidates := tagCandidatesFromMemoryRepo(t, repo, rule.ID, initialRevision)
	require.Len(t, tagCandidates, 1)
	_, err = svc.UpsertTagReviewDecision(ctx, rule.ID, UpsertTagReviewDecisionRequest{
		CandidateID:   tagCandidates[0].ID,
		Action:        schema.SourceRuleTagReviewDecisionActionOverride,
		OverrideTagID: overrideTag.ID,
	})
	require.NoError(t, err)

	initialCandidates := databaseOutputCandidatesFromMemoryRepo(t, repo, rule.ID, initialRevision)
	require.Len(t, initialCandidates, 1)
	initialCandidate := initialCandidates[0]
	require.NotNil(t, initialCandidate.MappingID)
	assert.Equal(t, "connector-revision", initialCandidate.ConnectorID)

	scaleMultiplier := 2.5
	updatedRule, err := svc.Update(ctx, rule.ID, UpdateRuleRequest{ScaleMultiplier: &scaleMultiplier})
	require.NoError(t, err)
	require.NotEqual(t, initialRevision, updatedRule.RevisionID)

	updatedCandidates := databaseOutputCandidatesFromMemoryRepo(t, repo, rule.ID, updatedRule.RevisionID)
	require.Len(t, updatedCandidates, 1)
	assert.Equal(t, initialCandidate.ID, updatedCandidates[0].ID)
	assert.Equal(t, initialCandidate.Identity, updatedCandidates[0].Identity)
	assert.Equal(t, initialCandidate.ProposedSignature, updatedCandidates[0].ProposedSignature)
	require.NotNil(t, updatedCandidates[0].MappingID)

	_, err = svc.ApplyDatabaseOutputCandidates(ctx, rule.ID, ApplyOutputCandidatesRequest{
		RevisionID:   initialRevision,
		CandidateIDs: []string{initialCandidate.ID},
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrOutputApplyRevisionConflict)

	response, err := svc.ApplyDatabaseOutputCandidates(ctx, rule.ID, ApplyOutputCandidatesRequest{
		RevisionID:   updatedRule.RevisionID,
		CandidateIDs: []string{initialCandidate.ID},
	})
	require.NoError(t, err)
	assert.Equal(t, rule.ID, response.SourceRuleID)
	assert.Equal(t, updatedRule.RevisionID, response.RevisionID)
	require.Len(t, response.Results, 1)
	assert.Equal(t, initialCandidate.ID, response.Results[0].CandidateID)
	assert.Equal(t, "success", response.Results[0].Status)
	assert.Equal(t, "db-map-revision", response.Results[0].MappingID)
	assert.Equal(t, "connector-revision", response.Results[0].ConnectorID)
}
