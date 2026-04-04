package point

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubPollingGroupSyncer struct {
	groups []*schema.PollingGroup
}

func (s *stubPollingGroupSyncer) AddPollingGroup(group *schema.PollingGroup) {
	if group == nil {
		return
	}
	s.groups = append(s.groups, group)
}

func TestServiceCreateSetsDefaultMode(t *testing.T) {
	repo := NewMemoryRepository()
	groupRepo := NewMemoryPollingGroupRepository()
	svc := NewService(repo, groupRepo)
	ctx := context.Background()

	point, err := svc.Create(ctx, CreatePointRequest{
		DeviceID: "dev-1",
		Name:     "P1",
		Address:  "D100",
		DataType: schema.DataTypeInt16,
	})
	require.NoError(t, err)
	assert.Equal(t, schema.PointModeReadOnly, point.Mode)
	assert.True(t, point.Enabled)
}

func TestCreatePollingGroupValidatesInterval(t *testing.T) {
	repo := NewMemoryRepository()
	groupRepo := NewMemoryPollingGroupRepository()
	svc := NewService(repo, groupRepo)
	ctx := context.Background()

	_, err := svc.CreatePollingGroup(ctx, CreatePollingGroupRequest{
		Name:       "fast",
		IntervalMs: 50,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "輪詢間隔不能小於 100ms")
}

func TestEnsureDefaultPollingGroupSyncsAutoCreatedGroup(t *testing.T) {
	repo := NewMemoryRepository()
	groupRepo := NewMemoryPollingGroupRepository()
	svc := NewService(repo, groupRepo)
	syncer := &stubPollingGroupSyncer{}
	svc.SetPollingGroupSyncer(syncer)
	ctx := context.Background()

	group, err := svc.EnsureDefaultPollingGroup(ctx)
	require.NoError(t, err)
	require.NotNil(t, group)
	require.Len(t, syncer.groups, 1)
	assert.Equal(t, group.ID, syncer.groups[0].ID)
}

func TestEnsureDefaultPollingGroupSyncsReenabledGroup(t *testing.T) {
	repo := NewMemoryRepository()
	groupRepo := NewMemoryPollingGroupRepository()
	svc := NewService(repo, groupRepo)
	ctx := context.Background()

	group, err := svc.CreatePollingGroup(ctx, CreatePollingGroupRequest{
		Name:       "legacy",
		IntervalMs: 1000,
	})
	require.NoError(t, err)

	disabled := false
	_, err = svc.UpdatePollingGroup(ctx, group.ID, UpdatePollingGroupRequest{Enabled: &disabled})
	require.NoError(t, err)

	syncer := &stubPollingGroupSyncer{}
	svc.SetPollingGroupSyncer(syncer)

	ensured, err := svc.EnsureDefaultPollingGroup(ctx)
	require.NoError(t, err)
	require.NotNil(t, ensured)
	require.Len(t, syncer.groups, 1)
	assert.Equal(t, group.ID, syncer.groups[0].ID)
	assert.True(t, ensured.Enabled)
}

func TestGetLastValueReturnsRawStringOnInvalidJSON(t *testing.T) {
	repo := NewMemoryRepository()
	groupRepo := NewMemoryPollingGroupRepository()
	svc := NewService(repo, groupRepo)
	ctx := context.Background()

	point, err := svc.Create(ctx, CreatePointRequest{
		DeviceID: "dev-1",
		Name:     "P2",
		Address:  "D101",
		DataType: schema.DataTypeString,
	})
	require.NoError(t, err)

	current, err := repo.GetByID(ctx, point.ID)
	require.NoError(t, err)
	raw := "not-json"
	current.LastValue = &raw
	require.NoError(t, repo.Update(ctx, current))

	value, err := svc.GetLastValue(ctx, point.ID)
	require.NoError(t, err)
	assert.Equal(t, "not-json", value)
}

func TestBatchCreate_DryRunDoesNotWriteDB(t *testing.T) {
	repo := NewMemoryRepository()
	groupRepo := NewMemoryPollingGroupRepository()
	svc := NewService(repo, groupRepo)
	ctx := context.Background()

	result, err := svc.BatchCreate(ctx, BatchCreatePointsRequest{
		DeviceID: "dev-1",
		DataType: schema.DataTypeUint16,
		DryRun:   true,
		Points: []BatchPointItem{
			{Name: "ok", Address: "40001", Function: "03"},
			{Name: "", Address: "bad@addr", Function: "03"},
		},
	})
	require.NoError(t, err)
	assert.True(t, result.DryRun)
	assert.False(t, result.Applied)
	assert.Equal(t, 0, result.CreatedCount)
	assert.Len(t, result.Points, 0)
	assert.NotEmpty(t, result.Errors)
	assert.Equal(t, 0, repo.Count(ctx))
}

func TestBatchCreate_ApplyIfClean(t *testing.T) {
	repo := NewMemoryRepository()
	groupRepo := NewMemoryPollingGroupRepository()
	svc := NewService(repo, groupRepo)
	ctx := context.Background()

	result, err := svc.BatchCreate(ctx, BatchCreatePointsRequest{
		DeviceID:     "dev-1",
		DataType:     schema.DataTypeUint16,
		DryRun:       true,
		ApplyIfClean: true,
		Points: []BatchPointItem{
			{Name: "p1", Address: "40001", Function: "03"},
			{Name: "p2", Address: "40002", Function: "03"},
		},
	})
	require.NoError(t, err)
	assert.True(t, result.DryRun)
	assert.True(t, result.Applied)
	assert.Equal(t, 2, result.CreatedCount)
	assert.Len(t, result.Errors, 0)
	assert.Equal(t, 2, repo.Count(ctx))
}

func TestBatchCreate_BlockWhenDuplicateAddressFunctionExists(t *testing.T) {
	repo := NewMemoryRepository()
	groupRepo := NewMemoryPollingGroupRepository()
	svc := NewService(repo, groupRepo)
	ctx := context.Background()

	_, err := svc.Create(ctx, CreatePointRequest{
		DeviceID: "dev-1",
		Name:     "existing",
		Address:  "40001",
		DataType: schema.DataTypeUint16,
	})
	require.NoError(t, err)

	result, err := svc.BatchCreate(ctx, BatchCreatePointsRequest{
		DeviceID: "dev-1",
		DataType: schema.DataTypeUint16,
		Points: []BatchPointItem{
			{Name: "dup", Address: "40001", Function: "03"},
		},
	})
	require.NoError(t, err)
	assert.False(t, result.Applied)
	assert.Equal(t, 0, result.CreatedCount)
	assert.NotEmpty(t, result.Errors)
	assert.Equal(t, 1, repo.Count(ctx))
}

func TestBatchCreate_AllowSameAddressWithDifferentFunction(t *testing.T) {
	repo := NewMemoryRepository()
	groupRepo := NewMemoryPollingGroupRepository()
	svc := NewService(repo, groupRepo)
	ctx := context.Background()

	_, err := svc.Create(ctx, CreatePointRequest{
		DeviceID: "dev-1",
		Name:     "existing",
		Address:  "40001",
		Function: "03",
		DataType: schema.DataTypeUint16,
	})
	require.NoError(t, err)

	result, err := svc.BatchCreate(ctx, BatchCreatePointsRequest{
		DeviceID: "dev-1",
		DataType: schema.DataTypeUint16,
		Points: []BatchPointItem{
			{Name: "allow-diff-fn", Address: "40001", Function: "04"},
		},
	})
	require.NoError(t, err)
	assert.True(t, result.Applied)
	assert.Equal(t, 1, result.CreatedCount)
	assert.Len(t, result.Errors, 0)
	assert.Equal(t, 2, repo.Count(ctx))
}

func TestPointValidation_CreateAndBatchDryRunUseSameFormat(t *testing.T) {
	repo := NewMemoryRepository()
	groupRepo := NewMemoryPollingGroupRepository()
	svc := NewService(repo, groupRepo)
	ctx := context.Background()

	_, createErr := svc.Create(ctx, CreatePointRequest{
		DeviceID: "dev-1",
		Name:     "bad-address",
		Address:  "bad@addr",
		Function: "03",
		DataType: schema.DataTypeUint16,
	})
	require.Error(t, createErr)
	assert.Contains(t, createErr.Error(), "point 驗證失敗: field=address")
	assert.Contains(t, createErr.Error(), "reason=格式無效")

	result, err := svc.BatchCreate(ctx, BatchCreatePointsRequest{
		DeviceID: "dev-1",
		DataType: schema.DataTypeUint16,
		DryRun:   true,
		Points: []BatchPointItem{
			{Name: "bad-address", Address: "bad@addr", Function: "03"},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0], "point 驗證失敗: field=address")
	assert.Contains(t, result.Errors[0], "reason=格式無效")
}

func TestPointValidation_NormalizeFunctionAcrossCreateUpdateAndBatch(t *testing.T) {
	repo := NewMemoryRepository()
	groupRepo := NewMemoryPollingGroupRepository()
	svc := NewService(repo, groupRepo)
	ctx := context.Background()

	created, err := svc.Create(ctx, CreatePointRequest{
		DeviceID: "dev-1",
		Name:     "norm",
		Address:  "40001",
		Function: "read_holding",
		DataType: schema.DataTypeUint16,
	})
	require.NoError(t, err)
	assert.Equal(t, "03", created.Function)

	function := "fc04"
	updated, err := svc.Update(ctx, created.ID, UpdatePointRequest{
		Function: &function,
	})
	require.NoError(t, err)
	assert.Equal(t, "04", updated.Function)

	result, err := svc.BatchCreate(ctx, BatchCreatePointsRequest{
		DeviceID: "dev-1",
		DataType: schema.DataTypeUint16,
		DryRun:   true,
		Points: []BatchPointItem{
			{Name: "alias-dup", Address: "40001", Function: "input"},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0], "field=address+function")
}
