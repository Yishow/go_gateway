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

func TestService_CandidateSnapshots_InfersGroupedDatabaseMetadataFromSlashTagKey(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-db-grouped")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-db-grouped",
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
		CandidateID: tagCandidates[0].ID,
		Action:      schema.SourceRuleTagReviewDecisionActionRename,
		TagKey:      "meter/A1",
	})
	require.NoError(t, err)

	databaseCandidates := databaseOutputCandidatesFromMemoryRepo(t, repo, rule.ID, rule.RevisionID)
	require.Len(t, databaseCandidates, 1)

	candidate := databaseCandidates[0]
	require.NotNil(t, candidate.GroupKey)
	assert.Equal(t, "meter", *candidate.GroupKey)
	assert.Equal(t, "a1", candidate.ColumnName)
	require.NotNil(t, candidate.WriteIntervalSeconds)
	assert.Equal(t, 15, *candidate.WriteIntervalSeconds)
	assert.Contains(t, candidate.Identity.TargetBindingScope, schema.SourceRuleCandidateScopeField{
		Key:   "database_group_key",
		Value: "meter",
	})
	assert.Contains(t, candidate.Identity.TargetBindingScope, schema.SourceRuleCandidateScopeField{
		Key:   "database_write_interval_seconds",
		Value: "15",
	})
}

func TestService_CandidateSnapshots_UsesPersistedDatabaseMappingScope(t *testing.T) {
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
		Key:         "factory.db.bound",
		DisplayName: "DB Bound",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)
	svc.SetDatabaseTargetMappingReader(DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
		return []*schema.DatabaseTargetMapping{
			{
				ID:              "db-map-1",
				TagID:           overrideTag.ID,
				ConnectorID:     "connector-1",
				TableSchema:     "public",
				TableName:       "measurements",
				ColumnName:      "line_a",
				WriteMode:       schema.DatabaseWriteModeUpsert,
				TimestampColumn: stringPtr("ts"),
			},
		}, nil
	}))

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-db-scope")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-db-scope",
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

	snapshots, err := svc.ListCandidateSnapshots(ctx, rule.ID)
	require.NoError(t, err)
	databaseCandidates, databaseSnapshot := databaseCandidatesFromSnapshots(t, snapshots)
	assert.Equal(t, schema.SourceRuleCandidateStatusReady, databaseSnapshot.Status)
	require.Len(t, databaseCandidates, 1)

	candidate := databaseCandidates[0]
	assert.Equal(t, "db-map-1", candidate["mapping_id"])
	assert.Equal(t, "connector-1", candidate["connector_id"])
	assert.Equal(t, "public", candidate["table_schema"])
	assert.Equal(t, "measurements", candidate["table_name"])
	assert.Equal(t, "line_a", candidate["column_name"])
	assert.Equal(t, "upsert", candidate["write_mode"])
	assert.Equal(t, "ts", candidate["timestamp_column"])

	identity := candidate["identity"].(map[string]any)
	scope := identity["target_binding_scope"].([]any)
	assert.Contains(t, scope, map[string]any{"key": "database_connector_id", "value": "connector-1"})
	assert.Contains(t, scope, map[string]any{"key": "database_table_schema", "value": "public"})
	assert.Contains(t, scope, map[string]any{"key": "database_table_name", "value": "measurements"})
	assert.Contains(t, scope, map[string]any{"key": "database_column_name", "value": "line_a"})
}

func TestService_CandidateSnapshots_UsesConnectorDefaultWriteIntervalWhenMappingOverrideIsEmpty(t *testing.T) {
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
		Key:         "factory.db.connector-default",
		DisplayName: "DB Connector Default",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)
	svc.SetDatabaseTargetMappingReader(DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
		return []*schema.DatabaseTargetMapping{
			{
				ID:          "db-map-connector-default",
				TagID:       overrideTag.ID,
				ConnectorID: "connector-45",
				TableSchema: "public",
				TableName:   "measurements",
				ColumnName:  "line_a",
				WriteMode:   schema.DatabaseWriteModeInsert,
			},
		}, nil
	}))
	svc.SetDatabaseTargetConnectorReader(DatabaseTargetConnectorGetFunc(func(context.Context, string) (*schema.DatabaseConnector, error) {
		return &schema.DatabaseConnector{
			ID:                          "connector-45",
			DefaultWriteIntervalSeconds: 45,
		}, nil
	}))

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-db-connector-default")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-db-connector-default",
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
	require.NotNil(t, databaseCandidates[0].WriteIntervalSeconds)
	assert.Equal(t, 45, *databaseCandidates[0].WriteIntervalSeconds)
}

func TestService_CandidateSnapshots_MarksInvalidDatabaseScopeOutOfSync(t *testing.T) {
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
		Key:         "factory.db.broken",
		DisplayName: "DB Broken",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)
	svc.SetDatabaseTargetMappingReader(DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
		return []*schema.DatabaseTargetMapping{
			{
				ID:          "db-map-broken",
				TagID:       overrideTag.ID,
				ConnectorID: "connector-broken",
				TableSchema: "public",
				TableName:   "measurements",
				ColumnName:  "line_a",
				WriteMode:   schema.DatabaseWriteModeInsert,
			},
		}, nil
	}))
	svc.SetDatabaseTargetConnectorValidator(DatabaseTargetConnectorValidatorFunc(func(context.Context, string) (*DatabaseTargetConnectorValidation, error) {
		return &DatabaseTargetConnectorValidation{
			Ready: false,
			Issues: []DatabaseTargetValidationIssue{
				{
					Severity:  "error",
					MappingID: "db-map-broken",
					Code:      "table_missing",
					Message:   "找不到資料表: public.measurements",
				},
			},
		}, nil
	}))

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-db-invalid")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-db-invalid",
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

	snapshots, err := svc.ListCandidateSnapshots(ctx, rule.ID)
	require.NoError(t, err)
	databaseCandidates, databaseSnapshot := databaseCandidatesFromSnapshots(t, snapshots)
	assert.Equal(t, schema.SourceRuleCandidateStatusBlocked, databaseSnapshot.Status)
	assert.Contains(t, databaseSnapshot.Reason, "找不到資料表")
	require.Len(t, databaseCandidates, 1)
	assert.Equal(t, "out_of_sync", databaseCandidates[0]["status"])
	assert.Contains(t, databaseCandidates[0]["blocking_reason"], "找不到資料表")
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
