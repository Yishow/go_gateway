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

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/datalink/devices", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
}

func TestRouterHealthEndpoint(t *testing.T) {
	router := NewRouter(DefaultConfig(), Services{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/health", nil)
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
