// Package handlers 提供 Device Handler 擴展 API 的單元測試。
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/device"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

/**
 * setupDeviceRouterWithExtended 建立測試用的 Device Router（包含擴展端點）
 * @param t 測試實例
 * @returns *gin.Engine 測試路由器
 */
func setupDeviceRouterWithExtended() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := NewDeviceHandler()

	// 基礎端點
	r.GET("/datalink/devices", h.List)
	r.POST("/datalink/devices", h.Create)

	// 擴展端點
	r.POST("/datalink/devices/test-batch", h.TestConnectionBatch)
	r.POST("/datalink/devices/:id/activate", h.Activate)
	r.POST("/datalink/devices/:id/disable", h.Disable)

	return r
}

/**
 * TestDeviceHandler_Activate 測試啟用設備
 */
func TestDeviceHandler_Activate(t *testing.T) {
	r := setupDeviceRouterWithExtended()

	// 先建立一個設備
	newDevice := device.CreateDeviceRequest{
		Name:        "測試啟用設備",
		Protocol:    "modbus_tcp",
		ConnectionConfig: map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     502,
			"slave_id": 1,
		},
	}

	body, _ := json.Marshal(newDevice)
	req, _ := http.NewRequest("POST", "/datalink/devices", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	deviceID := createResp["data"].(map[string]interface{})["id"].(string)

	// 啟用設備
	req, _ = http.NewRequest("POST", "/datalink/devices/"+deviceID+"/activate", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var activateResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &activateResp)
	assert.True(t, activateResp["success"].(bool))
}

/**
 * TestDeviceHandler_Activate_NotFound 測試啟用不存在的設備
 */
func TestDeviceHandler_Activate_NotFound(t *testing.T) {
	r := setupDeviceRouterWithExtended()

	req, _ := http.NewRequest("POST", "/datalink/devices/non-existent-id/activate", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

/**
 * TestDeviceHandler_Disable 測試停用設備
 */
func TestDeviceHandler_Disable(t *testing.T) {
	r := setupDeviceRouterWithExtended()

	// 先建立一個設備
	newDevice := device.CreateDeviceRequest{
		Name:        "測試停用設備",
		Protocol:    "modbus_tcp",
		ConnectionConfig: map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     502,
			"slave_id": 1,
		},
	}

	body, _ := json.Marshal(newDevice)
	req, _ := http.NewRequest("POST", "/datalink/devices", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	deviceID := createResp["data"].(map[string]interface{})["id"].(string)

	// 停用設備
	req, _ = http.NewRequest("POST", "/datalink/devices/"+deviceID+"/disable", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var disableResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &disableResp)
	assert.True(t, disableResp["success"].(bool))
}

/**
 * TestDeviceHandler_Disable_NotFound 測試停用不存在的設備
 */
func TestDeviceHandler_Disable_NotFound(t *testing.T) {
	r := setupDeviceRouterWithExtended()

	req, _ := http.NewRequest("POST", "/datalink/devices/non-existent-id/disable", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

/**
 * TestDeviceHandler_TestConnectionBatch 測試批量測試連線
 */
func TestDeviceHandler_TestConnectionBatch(t *testing.T) {
	r := setupDeviceRouterWithExtended()

	// 先建立多個設備
	devices := []device.CreateDeviceRequest{
		{
			Name:        "測試設備 1",
			Protocol:    "modbus_tcp",
			ConnectionConfig: map[string]interface{}{
				"host":     "127.0.0.1",
				"port":     502,
				"slave_id": 1,
			},
		},
		{
			Name:        "測試設備 2",
			Protocol:    "modbus_tcp",
			ConnectionConfig: map[string]interface{}{
				"host":     "127.0.0.2",
				"port":     502,
				"slave_id": 2,
			},
		},
	}

	deviceIDs := make([]string, 0, len(devices))
	for _, dev := range devices {
		body, _ := json.Marshal(dev)
		req, _ := http.NewRequest("POST", "/datalink/devices", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var createResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &createResp)
		deviceID := createResp["data"].(map[string]interface{})["id"].(string)
		deviceIDs = append(deviceIDs, deviceID)
	}

	// 批量測試連線
	batchReq := map[string]interface{}{
		"device_ids": deviceIDs,
	}

	body, _ := json.Marshal(batchReq)
	req, _ := http.NewRequest("POST", "/datalink/devices/test-batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Equal(t, 2, len(data))

	// 驗證每個結果都有必要的欄位
	for _, result := range data {
		resultMap := result.(map[string]interface{})
		assert.Contains(t, resultMap, "latency_ms")
		assert.Contains(t, resultMap, "success")
		assert.Contains(t, resultMap, "timestamp")
	}
}

func TestDeviceHandler_TestConnectionBatch_ResultStructure(t *testing.T) {
	r := setupDeviceRouterWithExtended()

	batchReq := map[string]interface{}{
		"device_ids": []string{"non-existent-id"},
	}

	body, _ := json.Marshal(batchReq)
	req, _ := http.NewRequest("POST", "/datalink/devices/test-batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Len(t, data, 1)

	resultMap := data[0].(map[string]interface{})
	assert.Contains(t, resultMap, "latency_ms")
	assert.Contains(t, resultMap, "success")
	assert.Contains(t, resultMap, "timestamp")
	assert.Contains(t, resultMap, "error")
	assert.False(t, resultMap["success"].(bool))
}

/**
 * TestDeviceHandler_TestConnectionBatch_Empty 測試批量測試空陣列
 */
func TestDeviceHandler_TestConnectionBatch_Empty(t *testing.T) {
	r := setupDeviceRouterWithExtended()

	batchReq := map[string]interface{}{
		"device_ids": []string{},
	}

	body, _ := json.Marshal(batchReq)
	req, _ := http.NewRequest("POST", "/datalink/devices/test-batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Equal(t, 0, len(data))
}

/**
 * TestDeviceHandler_TestConnectionBatch_MissingField 測試批量測試缺少欄位
 */
func TestDeviceHandler_TestConnectionBatch_MissingField(t *testing.T) {
	r := setupDeviceRouterWithExtended()

	// 缺少 device_ids 欄位
	batchReq := map[string]interface{}{}

	body, _ := json.Marshal(batchReq)
	req, _ := http.NewRequest("POST", "/datalink/devices/test-batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

/**
 * TestDeviceHandler_TestConnectionBatch_InvalidDeviceIDs 測試批量測試無效的設備 ID
 */
func TestDeviceHandler_TestConnectionBatch_InvalidDeviceIDs(t *testing.T) {
	r := setupDeviceRouterWithExtended()

	batchReq := map[string]interface{}{
		"device_ids": []string{"non-existent-id-1", "non-existent-id-2"},
	}

	body, _ := json.Marshal(batchReq)
	req, _ := http.NewRequest("POST", "/datalink/devices/test-batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Equal(t, 2, len(data))

	// 驗證結果都是失敗的
	for _, result := range data {
		resultMap := result.(map[string]interface{})
		assert.False(t, resultMap["success"].(bool))
		assert.NotEmpty(t, resultMap["error"])
	}
}

/**
 * TestDeviceHandler_ActivateThenDisable 測試啟用後停用設備
 */
func TestDeviceHandler_ActivateThenDisable(t *testing.T) {
	r := setupDeviceRouterWithExtended()

	// 建立設備
	newDevice := device.CreateDeviceRequest{
		Name:        "測試循環設備",
		Protocol:    "modbus_tcp",
		ConnectionConfig: map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     502,
			"slave_id": 1,
		},
	}

	body, _ := json.Marshal(newDevice)
	req, _ := http.NewRequest("POST", "/datalink/devices", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	deviceID := createResp["data"].(map[string]interface{})["id"].(string)

	// 啟用設備
	req, _ = http.NewRequest("POST", "/datalink/devices/"+deviceID+"/activate", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 停用設備
	req, _ = http.NewRequest("POST", "/datalink/devices/"+deviceID+"/disable", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
}

/**
 * TestDeviceHandler_TestConnectionBatch_LargeList 測試大批量設備連線測試
 */
func TestDeviceHandler_TestConnectionBatch_LargeList(t *testing.T) {
	r := setupDeviceRouterWithExtended()

	// 建立 10 個設備
	deviceIDs := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		newDevice := device.CreateDeviceRequest{
			Name:        "批量測試設備 " + string(rune('0'+i)),
			Protocol:    "modbus_tcp",
			ConnectionConfig: map[string]interface{}{
				"host":     "127.0.0.1",
				"port":     502 + i,
				"slave_id": 1,
			},
		}

		body, _ := json.Marshal(newDevice)
		req, _ := http.NewRequest("POST", "/datalink/devices", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var createResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &createResp)
		deviceID := createResp["data"].(map[string]interface{})["id"].(string)
		deviceIDs = append(deviceIDs, deviceID)
	}

	// 批量測試
	batchReq := map[string]interface{}{
		"device_ids": deviceIDs,
	}

	body, _ := json.Marshal(batchReq)
	req, _ := http.NewRequest("POST", "/datalink/devices/test-batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Equal(t, 10, len(data))
}
