// Package handlers 提供 Device Handler 擴展 API 的單元測試。
package handlers

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "go-gateway/internal/datalink/connector/adapters"
)

/**
 * setupDeviceRouterWithExtended 建立測試用的 Device Router（包含擴展端點）
 * @param t 測試實例
 * @returns *gin.Engine 測試路由器
 */
func setupDeviceRouterWithExtended() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Setup dependency injection
	repo := device.NewMemoryRepository()
	svc := device.NewService(repo, nil)
	h := NewDeviceHandler(svc)

	// 基礎端點
	r.GET("/datalink/devices", h.List)
	r.POST("/datalink/devices", h.Create)

	// 靜態路徑必須在參數路徑之前註冊
	r.POST("/datalink/devices/test-batch", h.TestConnectionBatch)

	// 參數路徑 (:id)
	r.GET("/datalink/devices/:id", h.Get)
	r.POST("/datalink/devices/:id/test", h.TestConnection)
	r.POST("/datalink/devices/:id/activate", h.Activate)
	r.POST("/datalink/devices/:id/disable", h.Disable)

	return r
}

/**
 * TestDeviceHandler_Activate 測試啟用設備
 */
func TestDeviceHandler_Activate(t *testing.T) {
	// TODO: 此測試需要實際的設備連線才能通過
	// Activate 方法會先測試連線，在沒有實際設備的測試環境中會失敗
	t.Skip("需要 Mock ConnectionManager 以避免實際連線測試")

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

func TestDeviceHandler_TestConnection_SurfacesConnectAndProbeSuccess(t *testing.T) {
	bank := memory.NewMemoryBank(64)
	require.NoError(t, bank.WriteWord(0, 123))

	server := virtualmodbus.NewServer(bank)
	require.NoError(t, server.Start(0))
	defer server.Stop()

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	repo := device.NewMemoryRepository()
	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	svc := device.NewService(repo, connMgr)
	h := NewDeviceHandler(svc)
	r.POST("/datalink/devices", h.Create)
	r.POST("/datalink/devices/:id/test", h.TestConnection)

	createReq := device.CreateDeviceRequest{
		Name:     "probe-ok",
		Protocol: "modbus_tcp",
		ConnectionConfig: map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     server.Port(),
			"slave_id": 1,
			"timeout":  2,
		},
	}

	body, err := json.Marshal(createReq)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/datalink/devices", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &createResp))
	deviceID := createResp["data"].(map[string]interface{})["id"].(string)

	req, err = http.NewRequest("POST", "/datalink/devices/"+deviceID+"/test", nil)
	require.NoError(t, err)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	require.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.True(t, data["success"].(bool))
	assert.Equal(t, true, data["can_activate"])
	assert.Equal(t, true, data["can_collect"])
	require.Contains(t, data, "connect")
	require.Contains(t, data, "probe")
	assert.Equal(t, "success", data["connect"].(map[string]interface{})["status"])
	assert.Equal(t, "success", data["probe"].(map[string]interface{})["status"])
}

func TestDeviceHandler_TestConnection_SurfacesProbeFailureWhileConnectSucceeds(t *testing.T) {
	bank := memory.NewMemoryBank(64)
	require.NoError(t, bank.WriteWord(0, 123))

	server := virtualmodbus.NewServer(bank)
	require.NoError(t, server.Start(0))
	defer server.Stop()

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	repo := device.NewMemoryRepository()
	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	svc := device.NewService(repo, connMgr)
	h := NewDeviceHandler(svc)
	r.POST("/datalink/devices", h.Create)
	r.POST("/datalink/devices/:id/test", h.TestConnection)

	createReq := device.CreateDeviceRequest{
		Name:     "probe-failed",
		Protocol: "modbus_tcp",
		ConnectionConfig: map[string]interface{}{
			"host":           "127.0.0.1",
			"port":           server.Port(),
			"slave_id":       1,
			"timeout":        2,
			"probe_address":  "49999",
			"probe_function": "03",
		},
	}

	body, err := json.Marshal(createReq)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/datalink/devices", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &createResp))
	deviceID := createResp["data"].(map[string]interface{})["id"].(string)

	req, err = http.NewRequest("POST", "/datalink/devices/"+deviceID+"/test", nil)
	require.NoError(t, err)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	require.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.False(t, data["success"].(bool))
	assert.Equal(t, false, data["can_activate"])
	assert.Equal(t, false, data["can_collect"])
	require.Contains(t, data, "error")
	assert.Contains(t, data["error"].(string), "讀取探測失敗")
	require.Contains(t, data, "connect")
	require.Contains(t, data, "probe")
	assert.Equal(t, "success", data["connect"].(map[string]interface{})["status"])
	assert.Equal(t, "failed", data["probe"].(map[string]interface{})["status"])
	assert.NotEmpty(t, data["probe"].(map[string]interface{})["error"])
}

func TestDeviceHandler_TestDraftConnection_AllowsUnsavedPayload(t *testing.T) {
	bank := memory.NewMemoryBank(64)
	require.NoError(t, bank.WriteWord(0, 123))

	server := virtualmodbus.NewServer(bank)
	require.NoError(t, server.Start(0))
	defer server.Stop()

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	repo := device.NewMemoryRepository()
	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	svc := device.NewService(repo, connMgr)
	h := NewDeviceHandler(svc)
	r.POST("/datalink/devices/test-draft", h.TestDraftConnection)

	reqBody := map[string]interface{}{
		"protocol": "modbus_tcp",
		"connection_config": map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     server.Port(),
			"slave_id": 1,
			"timeout":  2,
		},
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/datalink/devices/test-draft", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	require.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.True(t, data["success"].(bool))
	assert.Equal(t, true, data["can_activate"])
	assert.Equal(t, true, data["can_collect"])
	require.Contains(t, data, "connect")
	require.Contains(t, data, "probe")
	assert.Equal(t, "success", data["connect"].(map[string]interface{})["status"])
	assert.Equal(t, "success", data["probe"].(map[string]interface{})["status"])
}

func TestDeviceHandler_TestDraftConnection_SurfacesConnectFailure(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	require.NoError(t, listener.Close())

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	repo := device.NewMemoryRepository()
	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	svc := device.NewService(repo, connMgr)
	h := NewDeviceHandler(svc)
	r.POST("/datalink/devices/test-draft", h.TestDraftConnection)

	reqBody := map[string]interface{}{
		"protocol": "modbus_tcp",
		"connection_config": map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     port,
			"slave_id": 1,
			"timeout":  1,
		},
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/datalink/devices/test-draft", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	require.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.False(t, data["success"].(bool))
	assert.Equal(t, false, data["can_activate"])
	assert.Equal(t, false, data["can_collect"])
	assert.Contains(t, data["error"].(string), "connection refused")
	assert.Equal(t, "failed", data["connect"].(map[string]interface{})["status"])
	assert.Equal(t, "skipped", data["probe"].(map[string]interface{})["status"])
}

func TestDeviceHandler_TestDraftConnection_SurfacesProbeFailure(t *testing.T) {
	bank := memory.NewMemoryBank(64)
	require.NoError(t, bank.WriteWord(0, 123))

	server := virtualmodbus.NewServer(bank)
	require.NoError(t, server.Start(0))
	defer server.Stop()

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	repo := device.NewMemoryRepository()
	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	svc := device.NewService(repo, connMgr)
	h := NewDeviceHandler(svc)
	r.POST("/datalink/devices/test-draft", h.TestDraftConnection)

	reqBody := map[string]interface{}{
		"protocol": "modbus_tcp",
		"connection_config": map[string]interface{}{
			"host":           "127.0.0.1",
			"port":           server.Port(),
			"slave_id":       1,
			"timeout":        2,
			"probe_address":  "49999",
			"probe_function": "03",
		},
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/datalink/devices/test-draft", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	require.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.False(t, data["success"].(bool))
	assert.Equal(t, false, data["can_activate"])
	assert.Equal(t, false, data["can_collect"])
	assert.Contains(t, data["error"].(string), "讀取探測失敗")
	assert.Equal(t, "success", data["connect"].(map[string]interface{})["status"])
	assert.Equal(t, "failed", data["probe"].(map[string]interface{})["status"])
	assert.NotEmpty(t, data["probe"].(map[string]interface{})["error"])
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
	// TODO: 此測試需要實際的設備連線才能通過
	t.Skip("需要 Mock ConnectionManager 以避免實際連線測試")

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
