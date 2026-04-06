package sourcerule

import (
	"context"
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tagSnapshotPayload struct {
	Candidates []schema.SourceRuleTagCandidate `json:"candidates"`
}

func TestService_SQLCandidateSnapshotsPersistAndRestoreAfterRestart(t *testing.T) {
	ctx := context.Background()
	db := setupSQLRepoDB(t)
	defer db.Close()

	deviceRepo := device.NewSQLRepository(db)
	pointRepo := point.NewSQLRepository(db)
	groupRepo := pollinggroup.NewSQLRepository(db)
	tagRepo := tag.NewSQLRepository(db)
	mappingRepo := mapping.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, deviceRepo, "device-snapshot"))

	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, groupRepo)
	tagSvc := tag.NewService(tagRepo)
	mappingSvc := mapping.NewServiceWithTagResolver(mappingRepo, tagSvc.GetByID)
	svc := NewService(NewSQLRepository(db), deviceSvc, pointSvc, nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	created, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-snapshot",
		DeviceID:     "device-snapshot",
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	snapshots, err := svc.ListCandidateSnapshots(ctx, created.ID)
	require.NoError(t, err)
	assertSnapshotSet(t, snapshots, created.RevisionID, 2)

	restartedSvc := NewService(NewSQLRepository(db), device.NewService(device.NewSQLRepository(db), nil), point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db)), nil)

	restored, err := restartedSvc.ListCandidateSnapshots(ctx, created.ID)
	require.NoError(t, err)
	assertSnapshotSet(t, restored, created.RevisionID, 2)
}

func TestService_SQLCandidateSnapshotsPersistPerRevision(t *testing.T) {
	ctx := context.Background()
	db := setupSQLRepoDB(t)
	defer db.Close()

	deviceRepo := device.NewSQLRepository(db)
	pointRepo := point.NewSQLRepository(db)
	groupRepo := pollinggroup.NewSQLRepository(db)
	tagRepo := tag.NewSQLRepository(db)
	mappingRepo := mapping.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, deviceRepo, "device-history"))

	repo := NewSQLRepository(db)
	tagSvc := tag.NewService(tagRepo)
	mappingSvc := mapping.NewServiceWithTagResolver(mappingRepo, tagSvc.GetByID)
	svc := NewService(repo, device.NewService(deviceRepo, nil), point.NewService(pointRepo, groupRepo), nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	created, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-history",
		DeviceID:     "device-history",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	initialRevision := created.RevisionID
	count := 2
	updated, err := svc.Update(ctx, created.ID, UpdateRuleRequest{Count: &count})
	require.NoError(t, err)
	require.NotEqual(t, initialRevision, updated.RevisionID)

	initialSnapshots, err := repo.ListCandidateSnapshots(ctx, created.ID, initialRevision)
	require.NoError(t, err)
	assertSnapshotSet(t, initialSnapshots, initialRevision, 1)

	updatedSnapshots, err := repo.ListCandidateSnapshots(ctx, created.ID, updated.RevisionID)
	require.NoError(t, err)
	assertSnapshotSet(t, updatedSnapshots, updated.RevisionID, 2)
}

func TestService_SQLCandidateSnapshotsRestoreLatestRevisionAfterRestart(t *testing.T) {
	ctx := context.Background()
	db := setupSQLRepoDB(t)
	defer db.Close()

	deviceRepo := device.NewSQLRepository(db)
	pointRepo := point.NewSQLRepository(db)
	groupRepo := pollinggroup.NewSQLRepository(db)
	tagRepo := tag.NewSQLRepository(db)
	mappingRepo := mapping.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, deviceRepo, "device-restart-history"))

	repo := NewSQLRepository(db)
	tagSvc := tag.NewService(tagRepo)
	mappingSvc := mapping.NewServiceWithTagResolver(mappingRepo, tagSvc.GetByID)
	svc := NewService(repo, device.NewService(deviceRepo, nil), point.NewService(pointRepo, groupRepo), nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	created, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-restart-history",
		DeviceID:     "device-restart-history",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	count := 2
	updated, err := svc.Update(ctx, created.ID, UpdateRuleRequest{Count: &count})
	require.NoError(t, err)

	restartedSvc := NewService(NewSQLRepository(db), device.NewService(device.NewSQLRepository(db), nil), point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db)), nil)
	restored, err := restartedSvc.ListCandidateSnapshots(ctx, created.ID)
	require.NoError(t, err)
	assertSnapshotSet(t, restored, updated.RevisionID, 2)
}

func TestService_SQLCandidateSnapshotsKeepIdentityWhenPayloadChanges(t *testing.T) {
	ctx := context.Background()
	db := setupSQLRepoDB(t)
	defer db.Close()

	deviceRepo := device.NewSQLRepository(db)
	pointRepo := point.NewSQLRepository(db)
	groupRepo := pollinggroup.NewSQLRepository(db)
	tagRepo := tag.NewSQLRepository(db)
	mappingRepo := mapping.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, deviceRepo, "device-signature"))

	repo := NewSQLRepository(db)
	tagSvc := tag.NewService(tagRepo)
	mappingSvc := mapping.NewServiceWithTagResolver(mappingRepo, tagSvc.GetByID)
	svc := NewService(repo, device.NewService(deviceRepo, nil), point.NewService(pointRepo, groupRepo), nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	created, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-signature",
		DeviceID:     "device-signature",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	initialTag := firstTagCandidate(t, repo, created.ID, created.RevisionID)

	multiplier := 2.5
	updated, err := svc.Update(ctx, created.ID, UpdateRuleRequest{ScaleMultiplier: &multiplier})
	require.NoError(t, err)

	updatedTag := firstTagCandidate(t, repo, created.ID, updated.RevisionID)
	assert.Equal(t, initialTag.ID, updatedTag.ID)
	assert.Equal(t, initialTag.Identity, updatedTag.Identity)
	assert.NotEqual(t, initialTag.ProposedSignature, updatedTag.ProposedSignature)
	assert.Empty(t, initialTag.TransformPipeline)
	require.Len(t, updatedTag.TransformPipeline, 1)
	assert.Equal(t, schema.TransformScale, updatedTag.TransformPipeline[0].Type)
}

func TestService_SQLCandidateSnapshotsChangeIdentityWhenDerivedTargetDataTypeChanges(t *testing.T) {
	ctx := context.Background()
	db := setupSQLRepoDB(t)
	defer db.Close()

	deviceRepo := device.NewSQLRepository(db)
	pointRepo := point.NewSQLRepository(db)
	groupRepo := pollinggroup.NewSQLRepository(db)
	tagRepo := tag.NewSQLRepository(db)
	mappingRepo := mapping.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, deviceRepo, "device-identity"))

	repo := NewSQLRepository(db)
	tagSvc := tag.NewService(tagRepo)
	mappingSvc := mapping.NewServiceWithTagResolver(mappingRepo, tagSvc.GetByID)
	svc := NewService(repo, device.NewService(deviceRepo, nil), point.NewService(pointRepo, groupRepo), nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	created, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-identity",
		DeviceID:     "device-identity",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	initialTag := firstTagCandidate(t, repo, created.ID, created.RevisionID)
	require.Len(t, initialTag.Identity.TargetBindingScope, 1)
	assert.Equal(t, "int16", initialTag.Identity.TargetBindingScope[0].Value)

	targetType := schema.DataTypeFloat64
	updated, err := svc.Update(ctx, created.ID, UpdateRuleRequest{TargetDataType: &targetType})
	require.NoError(t, err)

	updatedTag := firstTagCandidate(t, repo, created.ID, updated.RevisionID)
	require.Len(t, updatedTag.Identity.TargetBindingScope, 1)
	assert.Equal(t, "float64", updatedTag.Identity.TargetBindingScope[0].Value)
	assert.NotEqual(t, initialTag.ID, updatedTag.ID)
	assert.NotEqual(t, initialTag.Identity, updatedTag.Identity)
	assert.NotEqual(t, initialTag.ProposedSignature, updatedTag.ProposedSignature)
}

func assertSnapshotSet(t *testing.T, snapshots []*schema.SourceRuleCandidateSnapshot, revisionID string, expectedTagCount int) {
	t.Helper()

	require.Len(t, snapshots, 3)
	byType := make(map[schema.SourceRuleCandidateType]*schema.SourceRuleCandidateSnapshot, len(snapshots))
	for _, snapshot := range snapshots {
		byType[snapshot.CandidateType] = snapshot
		assert.Equal(t, revisionID, snapshot.RevisionID)
	}

	tagSnapshot := byType[schema.SourceRuleCandidateTypeTags]
	require.NotNil(t, tagSnapshot)
	assert.Equal(t, schema.SourceRuleCandidateStatusReady, tagSnapshot.Status)
	assert.Empty(t, tagSnapshot.Reason)

	var tags tagSnapshotPayload
	require.NoError(t, json.Unmarshal([]byte(tagSnapshot.Payload), &tags))
	require.Len(t, tags.Candidates, expectedTagCount)
	if expectedTagCount > 0 {
		assert.NotEmpty(t, tags.Candidates[0].ID)
		assert.Equal(t, schema.SourceRuleCandidateKindTag, tags.Candidates[0].Identity.CandidateKind)
		assert.NotEmpty(t, tags.Candidates[0].ProposedSignature)
		assert.Nil(t, tags.Candidates[0].TagID)
		assert.Nil(t, tags.Candidates[0].MappingID)
	}

	databaseSnapshot := byType[schema.SourceRuleCandidateTypeDatabaseOutputs]
	require.NotNil(t, databaseSnapshot)
	assert.Equal(t, schema.SourceRuleCandidateStatusDeferred, databaseSnapshot.Status)
	assert.NotEmpty(t, databaseSnapshot.Reason)

	localModbusSnapshot := byType[schema.SourceRuleCandidateTypeLocalModbusOutputs]
	require.NotNil(t, localModbusSnapshot)
	assert.Equal(t, schema.SourceRuleCandidateStatusDeferred, localModbusSnapshot.Status)
	assert.NotEmpty(t, localModbusSnapshot.Reason)
}

func firstTagCandidate(t *testing.T, repo *SQLRepository, ruleID, revisionID string) schema.SourceRuleTagCandidate {
	t.Helper()

	snapshots, err := repo.ListCandidateSnapshots(context.Background(), ruleID, revisionID)
	require.NoError(t, err)

	byType := make(map[schema.SourceRuleCandidateType]*schema.SourceRuleCandidateSnapshot, len(snapshots))
	for _, snapshot := range snapshots {
		byType[snapshot.CandidateType] = snapshot
	}

	tagSnapshot := byType[schema.SourceRuleCandidateTypeTags]
	require.NotNil(t, tagSnapshot)

	var payload tagSnapshotPayload
	require.NoError(t, json.Unmarshal([]byte(tagSnapshot.Payload), &payload))
	require.NotEmpty(t, payload.Candidates)
	return payload.Candidates[0]
}
