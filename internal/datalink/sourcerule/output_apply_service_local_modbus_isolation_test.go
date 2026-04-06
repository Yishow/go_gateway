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

func TestService_ApplyDatabaseOutputCandidates_PreservesLocalModbusSnapshot(t *testing.T) {
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
		Key:         "factory.db.apply.isolated",
		DisplayName: "DB Apply Isolated",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)
	svc.SetDatabaseTargetMappingReader(DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
		return []*schema.DatabaseTargetMapping{
			{
				ID:          "db-map-apply-isolated",
				TagID:       overrideTag.ID,
				ConnectorID: "connector-isolated",
				TableSchema: "public",
				TableName:   "measurements",
				ColumnName:  "line_a",
				WriteMode:   schema.DatabaseWriteModeInsert,
			},
		}, nil
	}))
	svc.SetDatabaseTargetConnectorValidator(DatabaseTargetConnectorValidatorFunc(func(context.Context, string) (*DatabaseTargetConnectorValidation, error) {
		return &DatabaseTargetConnectorValidation{Ready: true}, nil
	}))

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-db-apply-isolated")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-db-apply-isolated",
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
	_, err = svc.UpsertTagReviewDecision(ctx, rule.ID, UpsertTagReviewDecisionRequest{
		CandidateID:   tagCandidates[0].ID,
		Action:        schema.SourceRuleTagReviewDecisionActionOverride,
		OverrideTagID: overrideTag.ID,
	})
	require.NoError(t, err)

	databaseCandidates := databaseOutputCandidatesFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	require.Len(t, databaseCandidates, 1)
	require.NotNil(t, databaseCandidates[0].MappingID)

	localModbusPayload := `{"candidates":[{"tag_id":"tag-local","tag_key":"LM_APPLY","register":9}]}`
	setLocalModbusSnapshotForRule(t, repo, rule.ID, rule.RevisionID, localModbusPayload, schema.SourceRuleCandidateStatusBlocked, "preserved local modbus apply state")

	response, err := svc.ApplyDatabaseOutputCandidates(ctx, rule.ID, ApplyOutputCandidatesRequest{
		RevisionID:   rule.RevisionID,
		CandidateIDs: []string{databaseCandidates[0].ID},
	})
	require.NoError(t, err)
	require.Len(t, response.Results, 1)
	assert.Equal(t, "success", response.Results[0].Status)
	assert.Equal(t, "db-map-apply-isolated", response.Results[0].MappingID)
	assert.Equal(t, "connector-isolated", response.Results[0].ConnectorID)

	assertLocalModbusSnapshotForRule(t, repo, rule.ID, rule.RevisionID, localModbusPayload, schema.SourceRuleCandidateStatusBlocked, "preserved local modbus apply state")
}

func setLocalModbusSnapshotForRule(
	t *testing.T,
	repo *MemoryRepository,
	ruleID string,
	revisionID string,
	payload string,
	status schema.SourceRuleCandidateStatus,
	reason string,
) {
	t.Helper()

	snapshots, err := repo.ListCandidateSnapshots(context.Background(), ruleID, revisionID)
	require.NoError(t, err)

	nextSnapshots := make([]*schema.SourceRuleCandidateSnapshot, 0, len(snapshots))
	for _, snapshot := range snapshots {
		candidateSnapshot := *snapshot
		if candidateSnapshot.CandidateType == schema.SourceRuleCandidateTypeLocalModbusOutputs {
			candidateSnapshot.Payload = payload
			candidateSnapshot.Status = status
			candidateSnapshot.Reason = reason
		}
		nextSnapshots = append(nextSnapshots, &candidateSnapshot)
	}
	require.NoError(t, repo.ReplaceCandidateSnapshots(context.Background(), nextSnapshots))
}

func assertLocalModbusSnapshotForRule(
	t *testing.T,
	repo *MemoryRepository,
	ruleID string,
	revisionID string,
	payload string,
	status schema.SourceRuleCandidateStatus,
	reason string,
) {
	t.Helper()

	snapshots, err := repo.ListCandidateSnapshots(context.Background(), ruleID, revisionID)
	require.NoError(t, err)

	for _, snapshot := range snapshots {
		if snapshot.CandidateType != schema.SourceRuleCandidateTypeLocalModbusOutputs {
			continue
		}
		assert.Equal(t, status, snapshot.Status)
		assert.Equal(t, reason, snapshot.Reason)
		assert.JSONEq(t, payload, snapshot.Payload)
		return
	}

	t.Fatalf("local modbus snapshot not found for rule %s revision %s", ruleID, revisionID)
}
