package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/settings"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupSettingsRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	repo := settings.NewMemoryRepository()
	// Seed 預設設定
	_ = repo.Set(nil, "test_key", "test_value")

	svc := settings.NewService(repo)
	h := NewSettingsHandler(svc)

	r.GET("/datalink/settings", h.List)
	r.PUT("/datalink/settings/:key", h.Update)

	return r
}

func TestSettingsHandler_List(t *testing.T) {
	r := setupSettingsRouter()

	req, _ := http.NewRequest("GET", "/datalink/settings", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

func TestSettingsHandler_Update(t *testing.T) {
	r := setupSettingsRouter()

	reqBody := map[string]interface{}{
		"value": "updated_value",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/datalink/settings/test_key", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "updated_value", data["value"])
}
