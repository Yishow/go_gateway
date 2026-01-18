package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/connector"
	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	if !connector.IsRegistered(schema.ProtocolModbusTCP) {
		t.Fatal("modbus_tcp not registered - imports might be missing or init failed")
	}
}

func setupDeviceRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := NewDeviceHandler() // This uses MemoryRepo and seeds default data
	r.GET("/devices", h.List)
	r.POST("/devices", h.Create)
	r.GET("/devices/:id", h.Get)
	return r
}

func setupTagRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := NewTagHandler()
	r.GET("/tags", h.List)
	r.POST("/tags", h.Create)
	return r
}

func TestDeviceHandler_List(t *testing.T) {
	r := setupDeviceRouter()

	req, _ := http.NewRequest("GET", "/devices", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	// By default, NewDeviceHandler seeds one device
	data := response["data"].([]interface{})
	assert.NotEmpty(t, data)
	first := data[0].(map[string]interface{})
	assert.Equal(t, "Demo Modbus Device", first["name"])
}

func TestDeviceHandler_Create(t *testing.T) {
	r := setupDeviceRouter()

	newDevice := device.CreateDeviceRequest{
		Name:        "Test Device",
		Protocol:    schema.ProtocolModbusTCP,
		ConnectionConfig: map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     502,
			"slave_id": 1,
		},
	}
	body, _ := json.Marshal(newDevice)
	req, _ := http.NewRequest("POST", "/devices", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Test Device", data["name"])
	assert.NotEmpty(t, data["id"])
}

func setupMappingRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := NewMappingHandler()
	r.POST("/mappings", h.Create)
	r.POST("/mappings/preview", h.Preview)
	return r
}

func TestMappingHandler_Create(t *testing.T) {
	r := setupMappingRouter()

	newMapping := map[string]interface{}{
		"point_id": "test-point-id",
		"tag_id":   "test-tag-id",
		"transform_pipeline": []map[string]interface{}{
			{
				"type": "scale",
				"params": map[string]interface{}{
					"multiplier": 10,
				},
			},
		},
	}
	body, _ := json.Marshal(newMapping)
	req, _ := http.NewRequest("POST", "/mappings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestMappingHandler_Preview(t *testing.T) {
	r := setupMappingRouter()

	previewReq := map[string]interface{}{
		"raw_value": 5,
		"transform_pipeline": []map[string]interface{}{
			{
				"type": "scale",
				"params": map[string]interface{}{
					"multiplier": 2,
					"offset": 1,
				},
			},
		},
	}
	body, _ := json.Marshal(previewReq)
	req, _ := http.NewRequest("POST", "/mappings/preview", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	t.Logf("Response: %+v", response)
	data := response["data"].(map[string]interface{})
	
	// Check step results for debug
	if stepResults, ok := data["step_results"].([]interface{}); ok {
		t.Logf("Step Results: %+v", stepResults)
	}

	// 5 * 2 + 1 = 11
	// json unmarshals numbers to float64
	assert.Equal(t, float64(11), data["final_value"])
}



func TestTagHandler_List(t *testing.T) {
	r := setupTagRouter()

	req, _ := http.NewRequest("GET", "/tags", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 2) // Seeded 2 tags
}

func TestTagHandler_Create(t *testing.T) {
	r := setupTagRouter()

	newTag := tag.CreateTagRequest{
		Key: "test_key",
		DisplayName: "Test Key",
		DataType: "int16",
	}
	body, _ := json.Marshal(newTag)
	req, _ := http.NewRequest("POST", "/tags", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "test_key", data["key"])
}
