package point

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
