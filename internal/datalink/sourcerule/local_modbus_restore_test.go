package sourcerule

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_RestoreLocalModbusMappingState_RehydratesPersistedMappingsAfterRestart(t *testing.T) {
	ctx := context.Background()
	db := setupSQLRepoDB(t)
	defer db.Close()

	deviceRepo := device.NewSQLRepository(db)
	pointRepo := point.NewSQLRepository(db)
	groupRepo := pollinggroup.NewSQLRepository(db)
	tagRepo := tag.NewSQLRepository(db)
	mappingRepo := mapping.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, deviceRepo, "device-lm-restore"))

	repo := NewSQLRepository(db)
	tagSvc := tag.NewService(tagRepo)
	mappingSvc := mapping.NewServiceWithTagResolver(mappingRepo, tagSvc.GetByID)
	svc := NewService(repo, device.NewService(deviceRepo, nil), point.NewService(pointRepo, groupRepo), nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-lm-restore",
		DeviceID:     "device-lm-restore",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	links := applyRuleManagedLinksWithRepository(t, ctx, repo, svc, rule.ID)
	require.Len(t, links, 1)
	require.NotNil(t, links[0].TagID)

	svc.SetLocalModbusMappingReader(LocalModbusMappingListFunc(func(context.Context) ([]LocalModbusMappingRecord, error) {
		return []LocalModbusMappingRecord{
			{
				TagID:     *links[0].TagID,
				Register:  7,
				DataType:  schema.DataTypeInt16,
				UpdatedAt: time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC),
			},
		}, nil
	}))
	require.NoError(t, svc.persistCandidateSnapshots(ctx, rule, links))

	restartedTagSvc := tag.NewService(tag.NewSQLRepository(db))
	restartedModbusShareSvc := modbusshare.NewService(restartedTagSvc, 4096)
	restartedSvc := NewService(
		NewSQLRepository(db),
		device.NewService(device.NewSQLRepository(db), nil),
		point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db)),
		nil,
	)

	require.NoError(t, restartedSvc.RestoreLocalModbusMappingState(ctx, LocalModbusMappingUpsertFunc(func(ctx context.Context, mappingRecord LocalModbusMappingRecord) error {
		_, err := restartedModbusShareSvc.UpsertMapping(ctx, mappingRecord.TagID, mappingRecord.Register)
		return err
	})))

	restoredMappings := restartedModbusShareSvc.ListMappings()
	require.Len(t, restoredMappings, 1)
	assert.Equal(t, *links[0].TagID, restoredMappings[0].TagID)
	assert.Equal(t, uint16(7), restoredMappings[0].Register)
	assert.Equal(t, schema.DataTypeInt16, restoredMappings[0].DataType)
}

func TestService_SQLLocalModbusSnapshotsRestoreAfterRestart(t *testing.T) {
	ctx := context.Background()
	db := setupSQLRepoDB(t)
	defer db.Close()

	deviceRepo := device.NewSQLRepository(db)
	pointRepo := point.NewSQLRepository(db)
	groupRepo := pollinggroup.NewSQLRepository(db)
	tagRepo := tag.NewSQLRepository(db)
	mappingRepo := mapping.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, deviceRepo, "device-lm-snapshot-restart"))

	repo := NewSQLRepository(db)
	tagSvc := tag.NewService(tagRepo)
	mappingSvc := mapping.NewServiceWithTagResolver(mappingRepo, tagSvc.GetByID)
	svc := NewService(repo, device.NewService(deviceRepo, nil), point.NewService(pointRepo, groupRepo), nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-lm-snapshot-restart",
		DeviceID:     "device-lm-snapshot-restart",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	links := applyRuleManagedLinksWithRepository(t, ctx, repo, svc, rule.ID)
	require.Len(t, links, 1)
	require.NotNil(t, links[0].TagID)

	svc.SetLocalModbusMappingReader(LocalModbusMappingListFunc(func(context.Context) ([]LocalModbusMappingRecord, error) {
		return []LocalModbusMappingRecord{
			{
				TagID:     *links[0].TagID,
				Register:  7,
				DataType:  schema.DataTypeInt16,
				UpdatedAt: time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC),
			},
		}, nil
	}))
	require.NoError(t, svc.persistCandidateSnapshots(ctx, rule, links))

	restartedSvc := NewService(
		NewSQLRepository(db),
		device.NewService(device.NewSQLRepository(db), nil),
		point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db)),
		nil,
	)

	restoredSnapshots, err := restartedSvc.ListCandidateSnapshots(ctx, rule.ID)
	require.NoError(t, err)

	localModbusSnapshot := localModbusSnapshotFromSnapshots(t, restoredSnapshots)
	assert.Equal(t, rule.RevisionID, localModbusSnapshot.RevisionID)

	candidates, err := decodeCandidatePayload[schema.SourceRuleLocalModbusOutputCandidate](localModbusSnapshot.Payload)
	require.NoError(t, err)
	require.Len(t, candidates, 1)
	require.NotNil(t, candidates[0].TagID)
	require.NotNil(t, candidates[0].Register)
	assert.Equal(t, *links[0].TagID, *candidates[0].TagID)
	assert.Equal(t, uint16(7), *candidates[0].Register)
}

func applyRuleManagedLinksWithRepository(
	t *testing.T,
	ctx context.Context,
	repo Repository,
	svc *Service,
	ruleID string,
) []*schema.SourceRuleLink {
	t.Helper()

	rule, err := svc.GetByID(ctx, ruleID)
	require.NoError(t, err)

	links, err := svc.ListLinks(ctx, ruleID)
	require.NoError(t, err)
	require.NotEmpty(t, links)

	result := tagMappingSyncResult{
		updatedMappings: make(map[string]mappingRollbackState),
		updatedTags:     make(map[string]tagRollbackState),
	}
	for _, link := range links {
		pointRecord, err := svc.pointSvc.GetByID(ctx, link.PointID)
		require.NoError(t, err)

		tagRecord, mappingRecord, err := svc.ensureRuleTagMapping(ctx, nil, rule, pointRecord, link, rule.Enabled, &result)
		require.NoError(t, err)

		link.TagID = stringPtr(tagRecord.ID)
		link.MappingID = stringPtr(mappingRecord.ID)
		link.UpdatedAt = time.Now()
	}

	require.NoError(t, repo.DeleteLinks(ctx, ruleID))
	require.NoError(t, repo.CreateLinks(ctx, links))
	return links
}
