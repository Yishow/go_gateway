package sourcerule

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLRepository_SourceRuleShare_RoundTripsCreateListGetUpdate(t *testing.T) {
	db := setupSQLRepoDB(t)
	defer db.Close()

	startRegister := 40001
	stride := 2
	now := time.Now().UTC()
	rule := &schema.SourceRule{
		ID:                 "rule-share-repository",
		DeviceID:           "device-1",
		StartAddress:       "40001",
		Count:              2,
		DataType:           schema.DataTypeInt16,
		NamingPrefix:       "SRC",
		Enabled:            true,
		ShareEnabled:       true,
		ShareStartRegister: &startRegister,
		ShareStride:        &stride,
		Origin:             "manual",
		SkippedAddresses:   `[]`,
		RevisionID:         "revision-1",
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	repo := NewSQLRepository(db)
	ctx := context.Background()
	require.NoError(t, repo.Create(ctx, rule))

	items, err := repo.List(ctx, ListFilter{})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, rule.ID, items[0].ID)
	assert.True(t, items[0].ShareEnabled)
	require.NotNil(t, items[0].ShareStartRegister)
	assert.Equal(t, startRegister, *items[0].ShareStartRegister)
	require.NotNil(t, items[0].ShareStride)
	assert.Equal(t, stride, *items[0].ShareStride)

	stored, err := repo.GetByID(ctx, rule.ID)
	require.NoError(t, err)
	assert.True(t, stored.ShareEnabled)
	require.NotNil(t, stored.ShareStartRegister)
	assert.Equal(t, startRegister, *stored.ShareStartRegister)
	require.NotNil(t, stored.ShareStride)
	assert.Equal(t, stride, *stored.ShareStride)

	stored.ShareEnabled = false
	stored.ShareStartRegister = nil
	stored.ShareStride = nil
	stored.UpdatedAt = time.Now().UTC()
	require.NoError(t, repo.Update(ctx, stored))

	updated, err := repo.GetByID(ctx, rule.ID)
	require.NoError(t, err)
	assert.False(t, updated.ShareEnabled)
	assert.Nil(t, updated.ShareStartRegister)
	assert.Nil(t, updated.ShareStride)
}

func TestService_Create_DisabledShareDefaultsToFalseAndNullable(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)
	_, err := seedActiveDevice(ctx, deviceRepo, "device-share-default")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-share-default",
		DeviceID:     "device-share-default",
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      false,
	})
	require.NoError(t, err)
	assert.False(t, rule.ShareEnabled)
	assert.Nil(t, rule.ShareStartRegister)
	assert.Nil(t, rule.ShareStride)

	stored, err := svc.GetByID(ctx, rule.ID)
	require.NoError(t, err)
	assert.False(t, stored.ShareEnabled)
	assert.Nil(t, stored.ShareStartRegister)
	assert.Nil(t, stored.ShareStride)
}

func TestService_Create_RejectsInvalidEnabledShareConfig(t *testing.T) {
	startBelowMinimum := 40000
	validStride := 1
	validStart := 40001
	zeroStride := 0
	cases := []struct {
		name   string
		start  *int
		stride *int
	}{
		{name: "start below minimum", start: &startBelowMinimum, stride: &validStride},
		{name: "stride below minimum", start: &validStart, stride: &zeroStride},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewService(NewMemoryRepository(), nil, nil, nil)
			_, err := svc.Create(context.Background(), CreateRuleRequest{
				DeviceID:           "device-share-invalid",
				StartAddress:       "40001",
				Count:              1,
				DataType:           schema.DataTypeInt16,
				NamingPrefix:       "SRC",
				Enabled:            true,
				ShareEnabled:       true,
				ShareStartRegister: tc.start,
				ShareStride:        tc.stride,
			})
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrValidation)
		})
	}
}
