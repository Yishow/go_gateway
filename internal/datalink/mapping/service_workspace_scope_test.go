package mapping

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Update_RebindsTagWhenTagIDChanges(t *testing.T) {
	repo := NewMemoryRepository()
	svc := newServiceWithDefaultTagResolver(repo)
	ctx := context.Background()

	created, err := svc.Create(ctx, CreateMappingRequest{
		PointID:           "point-rebind",
		TagID:             "tag-old",
		TransformPipeline: []schema.TransformStep{},
	})
	require.NoError(t, err)

	updated, err := svc.Update(ctx, created.ID, UpdateMappingRequest{
		TagID:             stringPtr("tag-new"),
		TransformPipeline: []schema.TransformStep{},
	})
	require.NoError(t, err)
	assert.Equal(t, "tag-new", updated.TagID)
}

func TestService_Update_ReturnsValidationErrorWhenTagIsAlreadyBound(t *testing.T) {
	repo := NewMemoryRepository()
	svc := newServiceWithDefaultTagResolver(repo)
	ctx := context.Background()

	_, err := svc.Create(ctx, CreateMappingRequest{
		PointID:           "point-a",
		TagID:             "tag-shared",
		TransformPipeline: []schema.TransformStep{},
	})
	require.NoError(t, err)

	other, err := svc.Create(ctx, CreateMappingRequest{
		PointID:           "point-b",
		TagID:             "tag-other",
		TransformPipeline: []schema.TransformStep{},
	})
	require.NoError(t, err)

	_, err = svc.Update(ctx, other.ID, UpdateMappingRequest{
		TagID:             stringPtr("tag-shared"),
		TransformPipeline: []schema.TransformStep{},
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrValidation)
	assert.Contains(t, err.Error(), "已綁定其他 point")
}

func stringPtr(value string) *string {
	return &value
}
