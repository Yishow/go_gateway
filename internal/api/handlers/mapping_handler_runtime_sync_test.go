package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
)

type stubMappingRuntimeRefresher struct {
	refreshCalls int
	refreshErr   error
}

func (s *stubMappingRuntimeRefresher) RefreshMappings(ctx context.Context) error {
	s.refreshCalls++
	return s.refreshErr
}

func activeTagResolver(_ context.Context, id string) (*schema.Tag, error) {
	if id == "" {
		return nil, errors.New("missing tag id")
	}
	return &schema.Tag{
		ID:       id,
		Key:      "test.tag",
		DataType: schema.DataTypeInt16,
		Status:   schema.TagStatusActive,
	}, nil
}

func TestMappingHandler_CreateRefreshesRuntimeMappings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	refresher := &stubMappingRuntimeRefresher{}
	handler := NewMappingHandler(
		mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), activeTagResolver),
		refresher,
	)

	body, err := json.Marshal(mapping.CreateMappingRequest{
		PointID:           "point-1",
		TagID:             "tag-1",
		TransformPipeline: []schema.TransformStep{},
	})
	if err != nil {
		t.Fatalf("marshal request failed: %v", err)
	}

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/mappings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	handler.Create(c)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	if refresher.refreshCalls != 1 {
		t.Fatalf("expected runtime mappings refresh on create, got %d", refresher.refreshCalls)
	}
}

func TestMappingHandler_DeleteRefreshesRuntimeMappings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := mapping.NewMemoryRepository()
	svc := mapping.NewServiceWithTagResolver(repo, activeTagResolver)
	created, err := svc.Create(context.Background(), mapping.CreateMappingRequest{
		PointID:           "point-1",
		TagID:             "tag-1",
		TransformPipeline: []schema.TransformStep{},
	})
	if err != nil {
		t.Fatalf("create mapping failed: %v", err)
	}

	refresher := &stubMappingRuntimeRefresher{}
	handler := NewMappingHandler(svc, refresher)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/mappings/"+created.ID, http.NoBody)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Params = gin.Params{{Key: "id", Value: created.ID}}
	c.Request = req

	handler.Delete(c)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	if refresher.refreshCalls != 1 {
		t.Fatalf("expected runtime mappings refresh on delete, got %d", refresher.refreshCalls)
	}
}
