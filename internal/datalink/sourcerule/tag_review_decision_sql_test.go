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

func TestService_SQLTagReviewDecisionPersistsAfterRestart(t *testing.T) {
	ctx := context.Background()
	db := setupSQLRepoDB(t)
	defer db.Close()

	deviceRepo := device.NewSQLRepository(db)
	pointRepo := point.NewSQLRepository(db)
	groupRepo := pollinggroup.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, deviceRepo, "device-review-sql"))

	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, groupRepo)
	svc := NewService(NewSQLRepository(db), deviceSvc, pointSvc, nil)

	created, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-review-sql",
		DeviceID:     "device-review-sql",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	candidate := firstTagCandidate(t, NewSQLRepository(db), created.ID, created.RevisionID)
	decision, err := svc.UpsertTagReviewDecision(ctx, created.ID, UpsertTagReviewDecisionRequest{
		CandidateID: candidate.ID,
		Action:      schema.SourceRuleTagReviewDecisionActionRename,
		TagKey:      "factory.sql.temperature",
	})
	require.NoError(t, err)
	require.NotNil(t, decision)

	restartedSvc := NewService(NewSQLRepository(db), device.NewService(device.NewSQLRepository(db), nil), point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db)), nil)
	restored, err := restartedSvc.GetTagReviewDecision(ctx, created.ID, candidate.ID)
	require.NoError(t, err)
	require.NotNil(t, restored)
	assert.Equal(t, schema.SourceRuleTagReviewDecisionActionRename, restored.Action)
	assert.Equal(t, "factory.sql.temperature", restored.TagKey)
	assert.Nil(t, restored.OverrideTagID)
}
