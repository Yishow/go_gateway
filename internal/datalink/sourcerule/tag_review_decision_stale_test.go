package sourcerule

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Update_MarksTagReviewDecisionStaleWhenCandidateIdentityChanges(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-stale-decision")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-stale-decision",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	initialCandidate := firstTagCandidateFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	_, err = svc.UpsertTagReviewDecision(ctx, rule.ID, UpsertTagReviewDecisionRequest{
		CandidateID: initialCandidate.ID,
		Action:      schema.SourceRuleTagReviewDecisionActionRename,
		TagKey:      "factory.temperature",
	})
	require.NoError(t, err)

	targetType := schema.DataTypeFloat64
	updatedRule, err := svc.Update(ctx, rule.ID, UpdateRuleRequest{
		TargetDataType: &targetType,
	})
	require.NoError(t, err)

	updatedCandidate := firstTagCandidateFromMemoryRepo(t, repo, rule.ID, updatedRule.RevisionID)
	assert.NotEqual(t, initialCandidate.ID, updatedCandidate.ID)

	staleDecision, err := svc.GetTagReviewDecision(ctx, rule.ID, initialCandidate.ID)
	require.NoError(t, err)
	require.NotNil(t, staleDecision)
	assert.True(t, staleDecision.Stale)
	assert.Equal(t, updatedRule.RevisionID, staleDecision.StaleRevisionID)
	require.NotNil(t, staleDecision.StaleAt)

	currentDecision, err := svc.GetTagReviewDecision(ctx, rule.ID, updatedCandidate.ID)
	require.NoError(t, err)
	assert.Nil(t, currentDecision)
}

func TestService_SQLTagReviewDecisionStaleStatePersistsAfterRestart(t *testing.T) {
	ctx := context.Background()
	db := setupSQLRepoDB(t)
	defer db.Close()

	deviceRepo := device.NewSQLRepository(db)
	pointRepo := point.NewSQLRepository(db)
	groupRepo := pollinggroup.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, deviceRepo, "device-stale-sql"))

	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, groupRepo)
	svc := NewService(NewSQLRepository(db), deviceSvc, pointSvc, nil)

	created, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-stale-sql",
		DeviceID:     "device-stale-sql",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	initialCandidate := firstTagCandidate(t, NewSQLRepository(db), created.ID, created.RevisionID)
	_, err = svc.UpsertTagReviewDecision(ctx, created.ID, UpsertTagReviewDecisionRequest{
		CandidateID: initialCandidate.ID,
		Action:      schema.SourceRuleTagReviewDecisionActionSkip,
	})
	require.NoError(t, err)

	targetType := schema.DataTypeFloat64
	updatedRule, err := svc.Update(ctx, created.ID, UpdateRuleRequest{
		TargetDataType: &targetType,
	})
	require.NoError(t, err)

	restartedSvc := NewService(NewSQLRepository(db), device.NewService(device.NewSQLRepository(db), nil), point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db)), nil)
	staleDecision, err := restartedSvc.GetTagReviewDecision(ctx, created.ID, initialCandidate.ID)
	require.NoError(t, err)
	require.NotNil(t, staleDecision)
	assert.True(t, staleDecision.Stale)
	assert.Equal(t, updatedRule.RevisionID, staleDecision.StaleRevisionID)
	require.NotNil(t, staleDecision.StaleAt)
}
