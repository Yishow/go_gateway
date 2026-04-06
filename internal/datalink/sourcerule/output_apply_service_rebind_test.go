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

func TestService_ApplyDatabaseOutputCandidates_DoesNotSilentlyRebindToNewScope(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-db-apply-rebind")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-db-apply-rebind",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	tagCandidates := tagCandidatesFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	require.Len(t, tagCandidates, 1)
	overrideTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "factory.db.rebind",
		DisplayName: "DB Rebind",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)

	currentMappings := []*schema.DatabaseTargetMapping{
		{
			ID:          "db-map-rebind",
			TagID:       overrideTag.ID,
			ConnectorID: "connector-a",
			TableSchema: "public",
			TableName:   "measurements",
			ColumnName:  "line_a",
			WriteMode:   schema.DatabaseWriteModeInsert,
		},
	}
	svc.SetDatabaseTargetMappingReader(DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
		return currentMappings, nil
	}))
	svc.SetDatabaseTargetConnectorValidator(DatabaseTargetConnectorValidatorFunc(func(context.Context, string) (*DatabaseTargetConnectorValidation, error) {
		return &DatabaseTargetConnectorValidation{Ready: true}, nil
	}))

	_, err = svc.UpsertTagReviewDecision(ctx, rule.ID, UpsertTagReviewDecisionRequest{
		CandidateID:   tagCandidates[0].ID,
		Action:        schema.SourceRuleTagReviewDecisionActionOverride,
		OverrideTagID: overrideTag.ID,
	})
	require.NoError(t, err)

	originalCandidates := databaseOutputCandidatesFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	require.Len(t, originalCandidates, 1)
	originalCandidateID := originalCandidates[0].ID
	require.Equal(t, "connector-a", originalCandidates[0].ConnectorID)
	require.Equal(t, "measurements", originalCandidates[0].TableName)

	currentMappings = []*schema.DatabaseTargetMapping{
		{
			ID:          "db-map-rebind",
			TagID:       overrideTag.ID,
			ConnectorID: "connector-b",
			TableSchema: "analytics",
			TableName:   "measurements_v2",
			ColumnName:  "line_b",
			WriteMode:   schema.DatabaseWriteModeInsert,
		},
	}
	_, err = svc.RecomputeCandidateView(ctx, rule.ID)
	require.NoError(t, err)

	reboundCandidates := databaseOutputCandidatesFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	require.Len(t, reboundCandidates, 1)
	require.NotEqual(t, originalCandidateID, reboundCandidates[0].ID)
	require.Equal(t, "connector-b", reboundCandidates[0].ConnectorID)
	require.Equal(t, "measurements_v2", reboundCandidates[0].TableName)

	response, err := svc.ApplyDatabaseOutputCandidates(ctx, rule.ID, ApplyOutputCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{originalCandidateID},
	})
	require.NoError(t, err)
	require.Len(t, response.Results, 1)

	assert.Equal(t, originalCandidateID, response.Results[0].CandidateID)
	assert.Equal(t, "failed", response.Results[0].Status)
	assert.Equal(t, "validation", response.Results[0].Code)
	assert.Contains(t, response.Results[0].Reason, "not found in revision")
}
