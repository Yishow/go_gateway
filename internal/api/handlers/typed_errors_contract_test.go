package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypedAPIError_PreviewInvalidRequestContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &MappingHandler{}
	router := gin.New()
	router.POST("/mappings/preview", handler.Preview)

	// Missing transform_pipeline field
	badPayload := []byte(`{"raw_value": 123}`)
	req := newHandlerTestRequest(http.MethodPost, "/mappings/preview", bytes.NewReader(badPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	var resp struct {
		Success bool                  `json:"success"`
		Error   TypedAPIErrorEnvelope `json:"error"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.False(t, resp.Success)
	assert.Equal(t, ErrCodePreviewInvalidRequest, resp.Error.Code)
	assert.NotEmpty(t, resp.Error.RequestID)
	assert.False(t, resp.Error.Retryable)
}

func TestTypedAPIError_PreviewMalformedJSONDoesNotExposeBackendDetail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &MappingHandler{}
	router := gin.New()
	router.POST("/mappings/preview", handler.Preview)

	const backendDetail = "json: cannot unmarshal string into Go struct field MappingPreviewRequest.raw_value"
	req := newHandlerTestRequest(http.MethodPost, "/mappings/preview", bytes.NewBufferString(`{"raw_value":`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.NotContains(t, w.Body.String(), backendDetail)
	var resp struct {
		Success bool                  `json:"success"`
		Error   TypedAPIErrorEnvelope `json:"error"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, ErrCodePreviewInvalidRequest, resp.Error.Code)
	assert.Equal(t, "preview request is invalid", resp.Error.Message)
}

func TestTypedAPIError_RequestIDRejectsWhitespace(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &MappingHandler{}
	router := gin.New()
	router.POST("/mappings/preview", handler.Preview)

	req := newHandlerTestRequest(http.MethodPost, "/mappings/preview", bytes.NewBufferString(`{"raw_value":`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", "operator supplied id")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp struct {
		Error TypedAPIErrorEnvelope `json:"error"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotEqual(t, "operator supplied id", resp.Error.RequestID)
	assert.NotContains(t, resp.Error.RequestID, " ")
}

func TestTypedAPIError_RuntimeStatusUnavailableUsesStableEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewRuntimeHandler(nil, nil, nil, nil, nil, nil)
	router := gin.New()
	router.GET("/runtime/status", handler.Status)

	req := newHandlerTestRequest(http.MethodGet, "/runtime/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusServiceUnavailable, w.Code)
	var resp struct {
		Success bool                  `json:"success"`
		Error   TypedAPIErrorEnvelope `json:"error"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.False(t, resp.Success)
	assert.Equal(t, ErrCodeRuntimeSnapshotUnavailable, resp.Error.Code)
	assert.Equal(t, "runtime snapshot is unavailable", resp.Error.Message)
	assert.NotEmpty(t, resp.Error.RequestID)
}

func TestTypedAPIError_RuntimeStreamMissingDeviceContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewRuntimeStreamHandler(nil)
	router := gin.New()
	router.GET("/runtime/stream", handler.Stream)

	req := newHandlerTestRequest(http.MethodGet, "/runtime/stream", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// When source is nil, returns runtime_stream_unavailable
	require.Equal(t, http.StatusServiceUnavailable, w.Code)
	var resp struct {
		Success bool                  `json:"success"`
		Error   TypedAPIErrorEnvelope `json:"error"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.False(t, resp.Success)
	assert.Equal(t, ErrCodeRuntimeStreamUnavailable, resp.Error.Code)
	assert.NotEmpty(t, resp.Error.RequestID)
	assert.True(t, resp.Error.Retryable)
}
