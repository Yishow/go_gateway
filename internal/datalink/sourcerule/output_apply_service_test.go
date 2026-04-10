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

func TestService_ApplyDatabaseOutputCandidates_RejectsRevisionMismatch(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-db-apply-conflict")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-db-apply-conflict",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	candidates := databaseOutputCandidatesFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	require.Len(t, candidates, 1)

	_, err = svc.ApplyDatabaseOutputCandidates(ctx, rule.ID, ApplyOutputCandidatesRequest{
		RevisionID:   "stale-revision",
		CandidateIDs: []string{candidates[0].ID},
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrOutputApplyRevisionConflict)
}

func TestService_ApplyDatabaseOutputCandidates_RejectsMissingRevisionID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-db-apply-validation")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-db-apply-validation",
		DeviceID:     dev.ID,
		StartAddress: "40021",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	databaseCandidates := databaseOutputCandidatesFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	require.Len(t, databaseCandidates, 1)

	_, err = svc.ApplyDatabaseOutputCandidates(ctx, rule.ID, ApplyOutputCandidatesRequest{
		CandidateIDs: []string{databaseCandidates[0].ID},
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidOutputApplyRequest)
}

func TestService_ApplyDatabaseOutputCandidates_ReturnsPerItemResultsForPartialSuccess(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-db-apply-partial")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-db-apply-partial",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	tagCandidates := tagCandidatesFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	require.Len(t, tagCandidates, 2)
	overrideTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "factory.db.apply.partial",
		DisplayName: "DB Apply Partial",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)
	groupKey := "meter"
	writeIntervalSeconds := 15

	svc.SetDatabaseTargetMappingReader(DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
		return []*schema.DatabaseTargetMapping{
			{
				ID:                   "db-map-apply-partial",
				TagID:                overrideTag.ID,
				ConnectorID:          "connector-ready",
				TableSchema:          "public",
				TableName:            "measurements",
				ColumnName:           "line_a",
				GroupKey:             &groupKey,
				WriteMode:            schema.DatabaseWriteModeInsert,
				WriteIntervalSeconds: &writeIntervalSeconds,
			},
		}, nil
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

	databaseCandidates := databaseOutputCandidatesFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	require.Len(t, databaseCandidates, 2)

	var mappedCandidateID string
	var unmappedCandidateID string
	for _, candidate := range databaseCandidates {
		if candidate.MappingID == nil {
			unmappedCandidateID = candidate.ID
			continue
		}
		mappedCandidateID = candidate.ID
	}
	require.NotEmpty(t, mappedCandidateID)
	require.NotEmpty(t, unmappedCandidateID)

	response, err := svc.ApplyDatabaseOutputCandidates(ctx, rule.ID, ApplyOutputCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{mappedCandidateID, unmappedCandidateID},
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Len(t, response.Results, 2)

	assert.Equal(t, mappedCandidateID, response.Results[0].CandidateID)
	assert.Equal(t, "success", response.Results[0].Status)
	assert.Equal(t, "db-map-apply-partial", response.Results[0].MappingID)
	assert.Equal(t, "connector-ready", response.Results[0].ConnectorID)
	require.NotNil(t, response.Results[0].GroupKey)
	assert.Equal(t, "meter", *response.Results[0].GroupKey)
	assert.Equal(t, "public", response.Results[0].TableSchema)
	assert.Equal(t, "measurements", response.Results[0].TableName)
	assert.Equal(t, "line_a", response.Results[0].ColumnName)
	require.NotNil(t, response.Results[0].WriteIntervalSeconds)
	assert.Equal(t, 15, *response.Results[0].WriteIntervalSeconds)
	assert.Empty(t, response.Results[0].Code)
	assert.Empty(t, response.Results[0].Reason)

	assert.Equal(t, unmappedCandidateID, response.Results[1].CandidateID)
	assert.Equal(t, "failed", response.Results[1].Status)
	assert.Equal(t, "schema_missing", response.Results[1].Code)
	assert.NotEmpty(t, response.Results[1].Reason)
}

func TestService_ApplyLocalModbusOutputCandidates_DeferredSnapshotReturnsSkipped(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-local-modbus-apply")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-local-modbus-apply",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	response, err := svc.ApplyLocalModbusOutputCandidates(ctx, rule.ID, ApplyOutputCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{"lm-candidate-1"},
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Len(t, response.Results, 1)
	assert.Equal(t, "skipped", response.Results[0].Status)
	assert.Equal(t, "deferred", response.Results[0].Code)
	assert.Equal(t, "lm-candidate-1", response.Results[0].CandidateID)
}

func databaseOutputCandidatesFromMemoryRepo(
	t *testing.T,
	repo *MemoryRepository,
	ruleID string,
	revisionID string,
) []schema.SourceRuleDatabaseOutputCandidate {
	t.Helper()

	snapshots, err := repo.ListCandidateSnapshots(context.Background(), ruleID, revisionID)
	require.NoError(t, err)

	for _, snapshot := range snapshots {
		if snapshot.CandidateType != schema.SourceRuleCandidateTypeDatabaseOutputs {
			continue
		}
		var payload struct {
			Candidates []schema.SourceRuleDatabaseOutputCandidate `json:"candidates"`
		}
		require.NoError(t, json.Unmarshal([]byte(snapshot.Payload), &payload))
		return payload.Candidates
	}

	t.Fatalf("database output snapshot not found for rule %s revision %s", ruleID, revisionID)
	return nil
}
