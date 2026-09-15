package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouterServeHTTPOptionsWithCORS(t *testing.T) {
	router := NewRouter(DefaultConfig(), Services{})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodOptions, "/api/v1/datalink/devices", http.NoBody)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
}

func TestRouterHealthEndpoint(t *testing.T) {
	router := NewRouter(DefaultConfig(), Services{})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/health", http.NoBody)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body APIResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.True(t, body.Success)

	data, ok := body.Data.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "datalink", data["service"])
}
