package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
)

type stubPointRuntimeSyncer struct {
	upserted *schema.Point
	removed  string
}

func (s *stubPointRuntimeSyncer) UpsertPoint(p *schema.Point) {
	s.upserted = p
}

func (s *stubPointRuntimeSyncer) RemovePoint(pointID string) {
	s.removed = pointID
}

func TestPointHandler_CreateSyncsRuntimePoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	syncer := &stubPointRuntimeSyncer{}
	handler := NewPointHandler(point.NewService(point.NewMemoryRepository(), nil), syncer)

	body, err := json.Marshal(point.CreatePointRequest{
		DeviceID: "device-1",
		Name:     "Flow",
		Address:  "40001",
		DataType: schema.DataTypeInt16,
		Mode:     schema.PointModeReadOnly,
	})
	if err != nil {
		t.Fatalf("marshal request failed: %v", err)
	}

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/points", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	handler.Create(c)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	if syncer.upserted == nil {
		t.Fatal("expected runtime point sync on create")
	}
	if syncer.upserted.Address != "40001" {
		t.Fatalf("unexpected synced point: %+v", syncer.upserted)
	}
}

func TestPointHandler_DeleteRemovesRuntimePoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := point.NewMemoryRepository()
	svc := point.NewService(repo, nil)
	created, err := svc.Create(context.Background(), point.CreatePointRequest{
		DeviceID: "device-1",
		Name:     "Flow",
		Address:  "40001",
		DataType: schema.DataTypeInt16,
		Mode:     schema.PointModeReadOnly,
	})
	if err != nil {
		t.Fatalf("create point failed: %v", err)
	}

	syncer := &stubPointRuntimeSyncer{}
	handler := NewPointHandler(svc, syncer)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/points/"+created.ID, http.NoBody)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Params = gin.Params{{Key: "id", Value: created.ID}}
	c.Request = req

	handler.Delete(c)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	if syncer.removed != created.ID {
		t.Fatalf("expected runtime remove %s, got %s", created.ID, syncer.removed)
	}
}
