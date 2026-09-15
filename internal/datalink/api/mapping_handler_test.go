package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMappingHandler_Create_ReturnsValidationErrorWhenResolverMissing(t *testing.T) {
	handler := NewMappingHandler(mapping.NewService(mapping.NewMemoryRepository()), nil)

	reqBody := map[string]interface{}{
		"point_id": "point-1",
		"tag_id":   "tag-1",
		"transform_pipeline": []map[string]interface{}{
			{
				"type": "scale",
				"params": map[string]interface{}{
					"multiplier": 2,
				},
			},
		},
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/v1/datalink/mappings", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.Create(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var resp APIResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.False(t, resp.Success)
	require.NotNil(t, resp.Error)
	assert.Contains(t, resp.Error.Message, mapping.ErrTagResolverNotConfigured.Error())
}

func TestMappingHandler_Create_SucceedsWhenResolverConfigured(t *testing.T) {
	svc := mapping.NewService(mapping.NewMemoryRepository())
	svc.SetTagResolver(func(ctx context.Context, tagID string) (*schema.Tag, error) {
		return &schema.Tag{ID: tagID, DataType: schema.DataTypeFloat64}, nil
	})
	handler := NewMappingHandler(svc, nil)

	reqBody := map[string]interface{}{
		"point_id": "point-1",
		"tag_id":   "tag-1",
		"transform_pipeline": []map[string]interface{}{
			{
				"type": "scale",
				"params": map[string]interface{}{
					"multiplier": 2,
				},
			},
		},
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/v1/datalink/mappings", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.Create(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}
