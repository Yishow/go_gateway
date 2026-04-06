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

func TestService_TagReviewDecisionCarriesForwardAcrossCompatibleRevisions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		seedDecision func(t *testing.T, ctx context.Context, svc *Service, tagSvc *tag.Service, candidate schema.SourceRuleTagCandidate)
		assertResult func(t *testing.T, ctx context.Context, response *ApplyTagCandidatesResponse, tagSvc *tag.Service, mappingSvc *mapping.Service)
	}{
		{
			name: "rename",
			seedDecision: func(t *testing.T, ctx context.Context, svc *Service, tagSvc *tag.Service, candidate schema.SourceRuleTagCandidate) {
				t.Helper()
				_, err := svc.UpsertTagReviewDecision(ctx, "rule-carry-forward", UpsertTagReviewDecisionRequest{
					CandidateID: candidate.ID,
					Action:      schema.SourceRuleTagReviewDecisionActionRename,
					TagKey:      "factory.compat.rename",
				})
				require.NoError(t, err)
			},
			assertResult: func(t *testing.T, ctx context.Context, response *ApplyTagCandidatesResponse, tagSvc *tag.Service, mappingSvc *mapping.Service) {
				t.Helper()
				require.Len(t, response.Results, 1)
				assert.Equal(t, "applied", response.Results[0].Status)
				tagRecord, err := tagSvc.GetByID(ctx, response.Results[0].TagID)
				require.NoError(t, err)
				assert.Equal(t, "factory.compat.rename", tagRecord.Key)
			},
		},
		{
			name: "skip",
			seedDecision: func(t *testing.T, ctx context.Context, svc *Service, tagSvc *tag.Service, candidate schema.SourceRuleTagCandidate) {
				t.Helper()
				_, err := svc.UpsertTagReviewDecision(ctx, "rule-carry-forward", UpsertTagReviewDecisionRequest{
					CandidateID: candidate.ID,
					Action:      schema.SourceRuleTagReviewDecisionActionSkip,
				})
				require.NoError(t, err)
			},
			assertResult: func(t *testing.T, ctx context.Context, response *ApplyTagCandidatesResponse, tagSvc *tag.Service, mappingSvc *mapping.Service) {
				t.Helper()
				require.Len(t, response.Results, 1)
				assert.Equal(t, "failed", response.Results[0].Status)
				assert.Contains(t, response.Results[0].Error, "marked skip")

				tags, err := tagSvc.List(ctx, tag.ListFilter{})
				require.NoError(t, err)
				assert.Empty(t, tags)

				mappings, err := mappingSvc.List(ctx, mapping.ListFilter{})
				require.NoError(t, err)
				assert.Empty(t, mappings)
			},
		},
		{
			name: "override",
			seedDecision: func(t *testing.T, ctx context.Context, svc *Service, tagSvc *tag.Service, candidate schema.SourceRuleTagCandidate) {
				t.Helper()
				existingTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{
					Key:         "factory.compat.override",
					DisplayName: "Compat Override",
					DataType:    schema.DataTypeInt16,
				})
				require.NoError(t, err)

				_, err = svc.UpsertTagReviewDecision(ctx, "rule-carry-forward", UpsertTagReviewDecisionRequest{
					CandidateID:   candidate.ID,
					Action:        schema.SourceRuleTagReviewDecisionActionOverride,
					OverrideTagID: existingTag.ID,
				})
				require.NoError(t, err)
			},
			assertResult: func(t *testing.T, ctx context.Context, response *ApplyTagCandidatesResponse, tagSvc *tag.Service, mappingSvc *mapping.Service) {
				t.Helper()
				require.Len(t, response.Results, 1)
				assert.Equal(t, "applied", response.Results[0].Status)

				mappingRecord, err := mappingSvc.GetByID(ctx, response.Results[0].MappingID)
				require.NoError(t, err)
				tagRecord, err := tagSvc.GetByID(ctx, mappingRecord.TagID)
				require.NoError(t, err)
				assert.Equal(t, "factory.compat.override", tagRecord.Key)
			},
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
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

			dev, err := seedActiveDevice(ctx, deviceRepo, "device-carry-forward-"+testCase.name)
			require.NoError(t, err)

			rule, err := svc.Create(ctx, CreateRuleRequest{
				ID:           "rule-carry-forward",
				DeviceID:     dev.ID,
				StartAddress: "40001",
				Count:        1,
				DataType:     schema.DataTypeInt16,
				NamingPrefix: "SRC",
				Enabled:      true,
			})
			require.NoError(t, err)
			initialRevision := rule.RevisionID

			initialCandidate := firstTagCandidateFromMemoryRepo(t, repo, rule.ID, initialRevision)
			testCase.seedDecision(t, ctx, svc, tagSvc, initialCandidate)

			updatedRule, err := svc.Update(ctx, rule.ID, UpdateRuleRequest{
				NamingPrefix: stringPtr("ALT"),
			})
			require.NoError(t, err)
			assert.NotEqual(t, initialRevision, updatedRule.RevisionID)

			updatedCandidate := firstTagCandidateFromMemoryRepo(t, repo, rule.ID, updatedRule.RevisionID)
			assert.Equal(t, initialCandidate.ID, updatedCandidate.ID)

			decision, err := svc.GetTagReviewDecision(ctx, rule.ID, updatedCandidate.ID)
			require.NoError(t, err)
			require.NotNil(t, decision)
			assert.False(t, decision.Stale)
			assert.Empty(t, decision.StaleRevisionID)
			assert.Nil(t, decision.StaleAt)

			response, err := svc.ApplyTagCandidates(ctx, rule.ID, ApplyTagCandidatesRequest{
				RevisionID:   updatedRule.RevisionID,
				CandidateIDs: []string{updatedCandidate.ID},
			})
			require.NoError(t, err)

			testCase.assertResult(t, ctx, response, tagSvc, mappingSvc)
		})
	}
}
