package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/pollinggroup"

	"github.com/stretchr/testify/require"
)

func TestPollingGroupHandler_Create_DuplicateName(t *testing.T) {
	router := setupPollingGroupRouter()
	payload, err := json.Marshal(pollinggroup.CreateRequest{Name: "unique-group", IntervalMs: 1000})
	require.NoError(t, err)
	for _, status := range []int{http.StatusCreated, http.StatusInternalServerError} {
		request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/datalink/polling-groups", bytes.NewReader(payload))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		require.Equal(t, status, response.Code, response.Body.String())
		var result struct {
			Success bool `json:"success"`
		}
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &result))
		require.Equal(t, status == http.StatusCreated, result.Success)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/datalink/polling-groups", http.NoBody))
	require.Equal(t, http.StatusOK, response.Code)
	var result struct {
		Data []struct {
			Name string `json:"name"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &result))
	require.Len(t, result.Data, 1)
	require.Equal(t, "unique-group", result.Data[0].Name)
}
