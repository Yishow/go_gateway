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

func TestService_UpsertTagReviewDecision_PersistsAcrossCompatibleRevision(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		buildReq  func(t *testing.T, tagSvc *tag.Service, candidate schema.SourceRuleTagCandidate) UpsertTagReviewDecisionRequest
		assertion func(t *testing.T, decision *schema.SourceRuleTagReviewDecision)
	}{
		{
			name: "rename",
			buildReq: func(_ *testing.T, _ *tag.Service, candidate schema.SourceRuleTagCandidate) UpsertTagReviewDecisionRequest {
				return UpsertTagReviewDecisionRequest{
					CandidateID: candidate.ID,
					Action:      schema.SourceRuleTagReviewDecisionActionRename,
					TagKey:      "factory.temperature",
				}
			},
			assertion: func(t *testing.T, decision *schema.SourceRuleTagReviewDecision) {
				t.Helper()
				assert.Equal(t, schema.SourceRuleTagReviewDecisionActionRename, decision.Action)
				assert.Equal(t, "factory.temperature", decision.TagKey)
				assert.Nil(t, decision.OverrideTagID)
			},
		},
		{
			name: "skip",
			buildReq: func(_ *testing.T, _ *tag.Service, candidate schema.SourceRuleTagCandidate) UpsertTagReviewDecisionRequest {
				return UpsertTagReviewDecisionRequest{
					CandidateID: candidate.ID,
					Action:      schema.SourceRuleTagReviewDecisionActionSkip,
				}
			},
			assertion: func(t *testing.T, decision *schema.SourceRuleTagReviewDecision) {
				t.Helper()
				assert.Equal(t, schema.SourceRuleTagReviewDecisionActionSkip, decision.Action)
				assert.Empty(t, decision.TagKey)
				assert.Nil(t, decision.OverrideTagID)
			},
		},
		{
			name: "override",
			buildReq: func(t *testing.T, tagSvc *tag.Service, candidate schema.SourceRuleTagCandidate) UpsertTagReviewDecisionRequest {
				t.Helper()
				existing, err := tagSvc.Create(context.Background(), tag.CreateTagRequest{
					Key:      "existing.tag.key",
					DataType: candidate.DataType,
				})
				require.NoError(t, err)
				return UpsertTagReviewDecisionRequest{
					CandidateID:   candidate.ID,
					Action:        schema.SourceRuleTagReviewDecisionActionOverride,
					OverrideTagID: existing.ID,
				}
			},
			assertion: func(t *testing.T, decision *schema.SourceRuleTagReviewDecision) {
				t.Helper()
				assert.Equal(t, schema.SourceRuleTagReviewDecisionActionOverride, decision.Action)
				assert.Equal(t, "existing.tag.key", decision.TagKey)
				require.NotNil(t, decision.OverrideTagID)
				assert.NotEmpty(t, *decision.OverrideTagID)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

			dev, err := seedActiveDevice(ctx, deviceRepo, "device-review-"+tc.name)
			require.NoError(t, err)

			rule, err := svc.Create(ctx, CreateRuleRequest{
				ID:           "rule-review-" + tc.name,
				DeviceID:     dev.ID,
				StartAddress: "40001",
				Count:        1,
				DataType:     schema.DataTypeInt16,
				NamingPrefix: "SRC",
				Enabled:      true,
			})
			require.NoError(t, err)

			candidate := firstTagCandidateFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
			req := tc.buildReq(t, tagSvc, candidate)

			decision, err := svc.UpsertTagReviewDecision(ctx, rule.ID, req)
			require.NoError(t, err)
			require.NotNil(t, decision)
			tc.assertion(t, decision)

			updatedRule, err := svc.Update(ctx, rule.ID, UpdateRuleRequest{
				ScaleMultiplier:    float64Ptr(2),
				ScaleMultiplierSet: true,
			})
			require.NoError(t, err)

			updatedCandidate := firstTagCandidateFromMemoryRepo(t, repo, rule.ID, updatedRule.RevisionID)
			assert.Equal(t, candidate.ID, updatedCandidate.ID)

			restored, err := svc.GetTagReviewDecision(ctx, rule.ID, updatedCandidate.ID)
			require.NoError(t, err)
			require.NotNil(t, restored)
			tc.assertion(t, restored)
		})
	}
}
