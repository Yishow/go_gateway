package sourcerule

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
)

// 建立規則時尚無 database target mapping（候選 ready 但沒有 mapping_id）；
// 之後使用者建立 target，GET candidates 應自動重算並補上 mapping_id，
// 而不是停在過期快照讓 apply 一次就 schema_missing。
func TestService_GetCandidateView_RecomputesDatabaseOutputsWhenMappingAppears(t *testing.T) {
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
		Key:         "factory.db.late",
		DisplayName: "DB Late",
		DataType:    schema.DataTypeInt16,
	})
	require.NoError(t, err)

	var mappings []*schema.DatabaseTargetMapping
	svc.SetDatabaseTargetMappingReader(DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
		return mappings, nil
	}))

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-db-late-mapping")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-db-late-mapping",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "LATE",
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

	// 第一輪：還沒有 mapping，候選 ready 且沒有 mapping_id。
	view, err := svc.GetCandidateView(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, view.DatabaseOutputs.Candidates, 1)
	first, ok := view.DatabaseOutputs.Candidates[0].(schema.SourceRuleDatabaseOutputCandidate)
	require.True(t, ok)
	assert.Nil(t, first.MappingID)

	// 使用者此刻建立了 database target mapping。
	mappings = []*schema.DatabaseTargetMapping{
		{
			ID:          "db-map-late",
			TagID:       overrideTag.ID,
			ConnectorID: "connector-late",
			TableSchema: "public",
			TableName:   "measurements",
			ColumnName:  "late_a",
			WriteMode:   schema.DatabaseWriteModeInsert,
		},
	}

	// 第二輪 GET：應自動重算，補上 mapping_id 與欄位資訊。
	view, err = svc.GetCandidateView(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, view.DatabaseOutputs.Candidates, 1)
	second, ok := view.DatabaseOutputs.Candidates[0].(schema.SourceRuleDatabaseOutputCandidate)
	require.True(t, ok)
	require.NotNil(t, second.MappingID, "GET candidates 應在 mapping 出現後自動重算 mapping_id")
	assert.Equal(t, "db-map-late", *second.MappingID)
	assert.Equal(t, "public", second.TableSchema)
	assert.Equal(t, "late_a", second.ColumnName)
}
