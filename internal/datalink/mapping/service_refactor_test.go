package mapping

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newServiceWithDefaultTagResolver(repo Repository) *Service {
	svc := NewService(repo)
	svc.SetTagResolver(func(ctx context.Context, tagID string) (*schema.Tag, error) {
		return &schema.Tag{ID: tagID, DataType: schema.DataTypeFloat64}, nil
	})
	return svc
}

func TestService_CreateAndList(t *testing.T) {
	repo := NewMemoryRepository()
	svc := newServiceWithDefaultTagResolver(repo)
	ctx := context.Background()

	req := CreateMappingRequest{
		PointID: "point-1",
		TagID:   "tag-1",
		TransformPipeline: []schema.TransformStep{
			{
				Type: schema.TransformScale,
				Params: map[string]interface{}{
					"multiplier": 2.0,
				},
			},
		},
	}

	mapping, err := svc.Create(ctx, req)
	require.NoError(t, err)
	assert.NotEmpty(t, mapping.ID)
	assert.True(t, mapping.Enabled)

	var steps []schema.TransformStep
	require.NoError(t, json.Unmarshal([]byte(mapping.TransformPipeline), &steps))
	assert.Len(t, steps, 1)
	assert.Equal(t, schema.TransformScale, steps[0].Type)

	list, err := svc.List(ctx, ListFilter{})
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, mapping.ID, list[0].ID)
}

func TestService_Update(t *testing.T) {
	repo := NewMemoryRepository()
	svc := newServiceWithDefaultTagResolver(repo)
	ctx := context.Background()

	created, err := svc.Create(ctx, CreateMappingRequest{
		PointID: "point-2",
		TagID:   "tag-2",
		TransformPipeline: []schema.TransformStep{
			{Type: schema.TransformScale, Params: map[string]interface{}{"multiplier": 1.0}},
		},
	})
	require.NoError(t, err)

	enabled := false
	updated, err := svc.Update(ctx, created.ID, UpdateMappingRequest{
		Enabled: &enabled,
		TransformPipeline: []schema.TransformStep{
			{Type: schema.TransformCast, Params: map[string]interface{}{"to_type": "int32"}},
		},
	})
	require.NoError(t, err)
	assert.False(t, updated.Enabled)

	var steps []schema.TransformStep
	require.NoError(t, json.Unmarshal([]byte(updated.TransformPipeline), &steps))
	assert.Len(t, steps, 1)
	assert.Equal(t, schema.TransformCast, steps[0].Type)
}

func TestParseConditionalParams(t *testing.T) {
	op, threshold, err := parseConditionalParams(map[string]interface{}{
		"condition":   "value >= 42",
		"true_value":  "ok",
		"false_value": "ng",
	})
	require.NoError(t, err)
	assert.Equal(t, ">=", op)
	assert.Equal(t, 42.0, threshold)

	op, threshold, err = parseConditionalParams(map[string]interface{}{
		"operator":  "lte",
		"threshold": 8,
	})
	require.NoError(t, err)
	assert.Equal(t, "<=", op)
	assert.Equal(t, 8.0, threshold)
}

func TestService_Update_EnableBlockedWhenPreviewExecutionFails(t *testing.T) {
	repo := NewMemoryRepository()
	svc := newServiceWithDefaultTagResolver(repo)
	ctx := context.Background()

	created, err := svc.Create(ctx, CreateMappingRequest{
		PointID:           "point-exec-fail",
		TagID:             "tag-exec-fail",
		TransformPipeline: []schema.TransformStep{},
	})
	require.NoError(t, err)

	disabled := false
	_, err = svc.Update(ctx, created.ID, UpdateMappingRequest{Enabled: &disabled})
	require.NoError(t, err)

	enabled := true
	_, err = svc.Update(ctx, created.ID, UpdateMappingRequest{
		Enabled: &enabled,
		TransformPipeline: []schema.TransformStep{
			{Type: schema.TransformScale, Params: map[string]interface{}{"multiplier": 2.0}},
		},
		PreviewRawValue: "abc",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "映射預覽驗證失敗")

	current, getErr := svc.GetByID(ctx, created.ID)
	require.NoError(t, getErr)
	assert.False(t, current.Enabled)
}

func TestService_Update_EnableBlockedWhenOutputNotCastable(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	svc.SetTagResolver(func(ctx context.Context, tagID string) (*schema.Tag, error) {
		return &schema.Tag{ID: tagID, DataType: schema.DataTypeBool}, nil
	})
	ctx := context.Background()

	created, err := svc.Create(ctx, CreateMappingRequest{
		PointID:           "point-cast-fail",
		TagID:             "tag-cast-fail",
		TransformPipeline: []schema.TransformStep{},
		PreviewRawValue:   false,
	})
	require.NoError(t, err)

	disabled := false
	_, err = svc.Update(ctx, created.ID, UpdateMappingRequest{Enabled: &disabled})
	require.NoError(t, err)

	enabled := true
	_, err = svc.Update(ctx, created.ID, UpdateMappingRequest{
		Enabled:         &enabled,
		PreviewRawValue: "not-bool",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "輸出無法轉換為 tag data_type")

	current, getErr := svc.GetByID(ctx, created.ID)
	require.NoError(t, getErr)
	assert.False(t, current.Enabled)
}

func TestService_Create_BlockedWhenPreviewOutputNotCastable(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	svc.SetTagResolver(func(ctx context.Context, tagID string) (*schema.Tag, error) {
		return &schema.Tag{ID: tagID, DataType: schema.DataTypeBool}, nil
	})

	_, err := svc.Create(context.Background(), CreateMappingRequest{
		PointID: "point-create-cast-fail",
		TagID:   "tag-create-cast-fail",
		TransformPipeline: []schema.TransformStep{
			{Type: schema.TransformScale, Params: map[string]interface{}{"multiplier": 2.0}},
		},
		PreviewRawValue: 2.0,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "映射預覽驗證失敗")
}

func TestService_Create_ReturnsErrorWhenResolverMissing(t *testing.T) {
	svc := NewService(NewMemoryRepository())

	_, err := svc.Create(context.Background(), CreateMappingRequest{
		PointID: "point-missing-resolver",
		TagID:   "tag-missing-resolver",
		TransformPipeline: []schema.TransformStep{
			{Type: schema.TransformScale, Params: map[string]interface{}{"multiplier": 1.0}},
		},
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrTagResolverNotConfigured)
}

func TestService_Create_RejectsSecondMappingForSamePoint(t *testing.T) {
	repo := NewMemoryRepository()
	svc := newServiceWithDefaultTagResolver(repo)
	ctx := context.Background()

	_, err := svc.Create(ctx, CreateMappingRequest{
		PointID:           "point-unique",
		TagID:             "tag-unique-a",
		TransformPipeline: []schema.TransformStep{},
	})
	require.NoError(t, err)

	_, err = svc.Create(ctx, CreateMappingRequest{
		PointID:           "point-unique",
		TagID:             "tag-unique-b",
		TransformPipeline: []schema.TransformStep{},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "point point-unique 已綁定其他 tag")
}

func TestService_Create_RejectsSecondMappingForSameTag(t *testing.T) {
	repo := NewMemoryRepository()
	svc := newServiceWithDefaultTagResolver(repo)
	ctx := context.Background()

	_, err := svc.Create(ctx, CreateMappingRequest{
		PointID:           "point-unique-a",
		TagID:             "tag-unique",
		TransformPipeline: []schema.TransformStep{},
	})
	require.NoError(t, err)

	_, err = svc.Create(ctx, CreateMappingRequest{
		PointID:           "point-unique-b",
		TagID:             "tag-unique",
		TransformPipeline: []schema.TransformStep{},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "tag tag-unique 已綁定其他 point")
}

func TestService_Create_RespectsDisabledFlag(t *testing.T) {
	repo := NewMemoryRepository()
	svc := newServiceWithDefaultTagResolver(repo)
	ctx := context.Background()

	enabled := false
	created, err := svc.Create(ctx, CreateMappingRequest{
		PointID:           "point-disabled",
		TagID:             "tag-disabled",
		Enabled:           &enabled,
		TransformPipeline: []schema.TransformStep{},
	})
	require.NoError(t, err)
	assert.False(t, created.Enabled)
}

func TestService_UpdateEnable_ReturnsErrorWhenResolverMissing(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	now := time.Now()
	ctx := context.Background()

	err := repo.Create(ctx, &schema.Mapping{
		ID:                "mapping-no-resolver",
		PointID:           "point-no-resolver",
		TagID:             "tag-no-resolver",
		TransformPipeline: "[]",
		Enabled:           false,
		CreatedAt:         now,
		UpdatedAt:         now,
	})
	require.NoError(t, err)

	enabled := true
	_, err = svc.Update(ctx, "mapping-no-resolver", UpdateMappingRequest{
		Enabled:         &enabled,
		PreviewRawValue: 1.0,
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrTagResolverNotConfigured)
}
