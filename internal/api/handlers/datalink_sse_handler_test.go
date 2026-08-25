package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type previewRuntimeSource struct {
	subscribed chan struct{}
	values     chan runtime.ValueEvent
	deviceID   string
	pointIDs   []string
}

func TestDatalinkSSEHandler_MissingMappingIDIsInvalidBeforeDependencyCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewDatalinkSSEHandler(nil, nil)
	router := gin.New()
	router.GET("/preview", handler.PreviewStream)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, newHandlerTestRequest("GET", "/preview", nil))

	require.Equal(t, http.StatusUnprocessableEntity, response.Code)
	require.Contains(t, response.Body.String(), `"code":"preview_invalid_request"`)
}

func TestDatalinkSSEHandler_UnavailableDependenciesReturn503(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewDatalinkSSEHandler(nil, &previewRuntimeSource{})
	router := gin.New()
	router.GET("/preview", handler.PreviewStream)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, newHandlerTestRequest("GET", "/preview?mapping_id=map-1", nil))

	require.Equal(t, http.StatusServiceUnavailable, response.Code)
	require.Contains(t, response.Body.String(), `"code":"preview_unavailable"`)
}

func (s *previewRuntimeSource) SubscribeValueEvents(deviceID string, pointIDs []string) (values <-chan runtime.ValueEvent, cancel func()) {
	s.deviceID = deviceID
	s.pointIDs = append([]string(nil), pointIDs...)
	if s.subscribed != nil {
		close(s.subscribed)
	}
	return s.values, func() {}
}

func TestDatalinkSSEHandler_PreviewStreamsRealRuntimeValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	tagRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "preview.tag", DataType: schema.DataTypeFloat64})
	require.NoError(t, err)
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	mappingRecord, err := mappingSvc.Create(ctx, mapping.CreateMappingRequest{PointID: "point-1", TagID: tagRecord.ID, TransformPipeline: []schema.TransformStep{}})
	require.NoError(t, err)
	source := &previewRuntimeSource{subscribed: make(chan struct{}), values: make(chan runtime.ValueEvent, 1)}
	handler := NewDatalinkSSEHandler(mappingSvc, source).WithPreviewScopeValidator(func(_ context.Context, workspaceID string, _ *schema.Mapping) (string, error) {
		if workspaceID != "ws-a" {
			return "", errors.New("mapping is outside workspace")
		}
		return "device-a", nil
	})
	router := gin.New()
	router.GET("/preview", handler.PreviewStream)
	requestCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	req := newHandlerTestRequest("GET", "/preview?mapping_id="+mappingRecord.ID+"&workspace_id=ws-a", nil).WithContext(requestCtx)
	response := httptest.NewRecorder()
	done := make(chan struct{})
	go func() {
		router.ServeHTTP(response, req)
		close(done)
	}()
	select {
	case <-source.subscribed:
	case <-time.After(time.Second):
		t.Fatal("preview did not subscribe to runtime source")
	}
	source.values <- runtime.ValueEvent{PointID: "point-1", RawValue: 42.5, TransformedValue: 43.5, Quality: schema.QualityGood, Timestamp: time.Date(2026, 8, 25, 1, 2, 3, 0, time.UTC)}
	time.Sleep(20 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("preview stream did not close")
	}
	body := response.Body.String()
	require.Contains(t, body, `"raw_value":42.5`)
	require.Contains(t, body, `"final_value":43.5`)
	require.NotContains(t, body, "1010")
	require.True(t, strings.Contains(body, `"type":"connected"`))
	require.Equal(t, "device-a", source.deviceID)
}

func TestDatalinkSSEHandler_ClosedRuntimeStreamEmitsTypedSafeError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	tagRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "preview.closed", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	mappingRecord, err := mappingSvc.Create(ctx, mapping.CreateMappingRequest{PointID: "point-closed", TagID: tagRecord.ID})
	require.NoError(t, err)
	values := make(chan runtime.ValueEvent)
	close(values)
	handler := NewDatalinkSSEHandler(mappingSvc, &previewRuntimeSource{values: values}).WithPreviewScopeValidator(func(_ context.Context, workspaceID string, _ *schema.Mapping) (string, error) {
		if workspaceID != "ws-a" {
			return "", errors.New("mapping is outside workspace")
		}
		return "device-a", nil
	})
	router := gin.New()
	router.GET("/preview", handler.PreviewStream)
	req := newHandlerTestRequest("GET", "/preview?mapping_id="+mappingRecord.ID+"&workspace_id=ws-a", nil)
	req.Header.Set("X-Request-ID", "preview-test-request")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	body := response.Body.String()
	require.Contains(t, body, `"type":"error"`)
	require.Contains(t, body, `"code":"preview_stream_closed"`)
	require.Contains(t, body, `"request_id":"preview-test-request"`)
	require.NotContains(t, body, `"error":"runtime value stream closed"`)
	var event struct {
		Error *TypedAPIErrorEnvelope `json:"error"`
	}
	for _, line := range strings.Split(body, "\n") {
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		if json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &event) == nil && event.Error != nil {
			require.Equal(t, ErrCodePreviewStreamClosed, event.Error.Code)
			require.Equal(t, "preview-test-request", event.Error.RequestID)
			return
		}
	}
	t.Fatal("typed SSE error event not found")
}

func TestDatalinkSSEHandler_CrossWorkspacePreviewIsRejectedBeforeSubscription(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	tagRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "preview.foreign", DataType: schema.DataTypeInt16})
	require.NoError(t, err)
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	mappingRecord, err := mappingSvc.Create(ctx, mapping.CreateMappingRequest{PointID: "point-foreign", TagID: tagRecord.ID})
	require.NoError(t, err)
	source := &previewRuntimeSource{values: make(chan runtime.ValueEvent)}
	handler := NewDatalinkSSEHandler(mappingSvc, source).WithPreviewScopeValidator(func(_ context.Context, workspaceID string, _ *schema.Mapping) (string, error) {
		if workspaceID != "ws-a" {
			return "", errors.New("mapping is outside workspace")
		}
		return "device-a", nil
	})
	router := gin.New()
	router.GET("/preview", handler.PreviewStream)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, newHandlerTestRequest("GET", "/preview?mapping_id="+mappingRecord.ID+"&workspace_id=ws-b", nil))

	require.Equal(t, http.StatusForbidden, response.Code)
	require.Contains(t, response.Body.String(), `"code":"modbus_share_workspace_scope"`)
	require.Empty(t, source.deviceID)
}
