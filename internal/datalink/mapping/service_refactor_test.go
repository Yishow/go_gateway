package mapping

import (
	"context"
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_CreateAndList(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
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
	svc := NewService(repo)
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
