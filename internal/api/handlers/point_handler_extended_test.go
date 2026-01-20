// Package handlers 提供 Point Handler 擴展 API 的單元測試。
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

/**
 * setupPointRouterWithExtended 建立測試用的 Point Router（包含擴展端點）
 * @param t 測試實例
 * @returns *gin.Engine 測試路由器
 */
func setupPointRouterWithExtended() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	
	// 建立設備 handler 來創建測試數據（暫時不使用，但保留以備將來需要）
	_ = NewDeviceHandler()

	// 基礎端點
	h := NewPointHandler()
	r.GET("/datalink/points", h.List)
	r.POST("/datalink/points", h.Create)

	// 擴展端點
	r.POST("/datalink/points/:id/poll", h.Poll)
	r.POST("/datalink/points/poll", h.PollBatch)

	return r
}

/**
 * TestPointHandler_Poll 測試單一點位輪詢
 */
func TestPointHandler_Poll(t *testing.T) {
	r := setupPointRouterWithExtended()

	// 先建立設備和點位
	newDevice := device.CreateDeviceRequest{
		Name:        "測試設備",
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

	var deviceResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &deviceResp)
	deviceID := deviceResp["data"].(map[string]interface{})["id"].(string)

	// 建立點位
	newPoint := point.CreatePointRequest{
		DeviceID: deviceID,
		Name:     "測試點位",
		Address:  "40001",
		DataType: "int16",
		Mode:     "read",
	}

	body, _ = json.Marshal(newPoint)
	req, _ = http.NewRequest("POST", "/datalink/points", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var pointResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &pointResp)
	pointID := pointResp["data"].(map[string]interface{})["id"].(string)

	// 輪詢點位
	req, _ = http.NewRequest("POST", "/datalink/points/"+pointID+"/poll", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, pointID, data["point_id"])
	assert.Contains(t, data, "timestamp")
	assert.Contains(t, data, "quality")
}

/**
 * TestPointHandler_Poll_NotFound 測試輪詢不存在的點位
 */
func TestPointHandler_Poll_NotFound(t *testing.T) {
	r := setupPointRouterWithExtended()

	req, _ := http.NewRequest("POST", "/datalink/points/non-existent-id/poll", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

/**
 * TestPointHandler_PollBatch 測試批量輪詢
 */
func TestPointHandler_PollBatch(t *testing.T) {
	r := setupPointRouterWithExtended()

	// 建立設備
	newDevice := device.CreateDeviceRequest{
		Name:        "批量測試設備",
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

	var deviceResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &deviceResp)
	deviceID := deviceResp["data"].(map[string]interface{})["id"].(string)

	// 建立多個點位
	pointIDs := make([]string, 0, 3)
	for i := 0; i < 3; i++ {
		newPoint := point.CreatePointRequest{
			DeviceID: deviceID,
			Name:     "批量測試點位 " + string(rune('0'+i)),
			Address:  "4000" + string(rune('1'+i)),
			DataType: "int16",
			Mode:     "read",
		}

		body, _ = json.Marshal(newPoint)
		req, _ = http.NewRequest("POST", "/datalink/points", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var pointResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &pointResp)
		pointID := pointResp["data"].(map[string]interface{})["id"].(string)
		pointIDs = append(pointIDs, pointID)
	}

	// 批量輪詢
	batchReq := map[string]interface{}{
		"point_ids": pointIDs,
	}

	body, _ = json.Marshal(batchReq)
	req, _ = http.NewRequest("POST", "/datalink/points/poll", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Equal(t, 3, len(data))
}

/**
 * TestPointHandler_PollBatch_Empty 測試批量輪詢空陣列
 */
func TestPointHandler_PollBatch_Empty(t *testing.T) {
	r := setupPointRouterWithExtended()

	batchReq := map[string]interface{}{
		"point_ids": []string{},
	}

	body, _ := json.Marshal(batchReq)
	req, _ := http.NewRequest("POST", "/datalink/points/poll", bytes.NewBuffer(body))
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
 * TestPointHandler_PollBatch_MissingField 測試批量輪詢缺少欄位
 */
func TestPointHandler_PollBatch_MissingField(t *testing.T) {
	r := setupPointRouterWithExtended()

	// 缺少 point_ids 欄位
	batchReq := map[string]interface{}{}

	body, _ := json.Marshal(batchReq)
	req, _ := http.NewRequest("POST", "/datalink/points/poll", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

/**
 * TestPointHandler_PollBatch_PartialSuccess 測試批量輪詢部分成功
 */
func TestPointHandler_PollBatch_PartialSuccess(t *testing.T) {
	r := setupPointRouterWithExtended()

	// 建立一個點位
	newDevice := device.CreateDeviceRequest{
		Name:        "部分測試設備",
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

	var deviceResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &deviceResp)
	deviceID := deviceResp["data"].(map[string]interface{})["id"].(string)

	newPoint := point.CreatePointRequest{
		DeviceID: deviceID,
		Name:     "部分測試點位",
		Address:  "40001",
		DataType: "int16",
		Mode:     "read",
	}

	body, _ = json.Marshal(newPoint)
	req, _ = http.NewRequest("POST", "/datalink/points", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var pointResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &pointResp)
	pointID := pointResp["data"].(map[string]interface{})["id"].(string)

	// 批量輪詢（包含不存在的點位）
	batchReq := map[string]interface{}{
		"point_ids": []string{pointID, "non-existent-id-1", "non-existent-id-2"},
	}

	body, _ = json.Marshal(batchReq)
	req, _ = http.NewRequest("POST", "/datalink/points/poll", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Equal(t, 3, len(data))
}

/**
 * TestPointHandler_PollBatch_LargeList 測試大批量點位輪詢
 */
func TestPointHandler_PollBatch_LargeList(t *testing.T) {
	r := setupPointRouterWithExtended()

	// 建立設備
	newDevice := device.CreateDeviceRequest{
		Name:        "大批量測試設備",
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

	var deviceResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &deviceResp)
	deviceID := deviceResp["data"].(map[string]interface{})["id"].(string)

	// 建立 10 個點位
	pointIDs := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		newPoint := point.CreatePointRequest{
			DeviceID: deviceID,
			Name:     "大批量點位 " + string(rune('0'+i)),
			Address:  "4000" + string(rune('1'+i%10)),
			DataType: "int16",
			Mode:     "read",
		}

		body, _ = json.Marshal(newPoint)
		req, _ = http.NewRequest("POST", "/datalink/points", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var pointResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &pointResp)
		pointID := pointResp["data"].(map[string]interface{})["id"].(string)
		pointIDs = append(pointIDs, pointID)
	}

	// 批量輪詢
	batchReq := map[string]interface{}{
		"point_ids": pointIDs,
	}

	body, _ = json.Marshal(batchReq)
	req, _ = http.NewRequest("POST", "/datalink/points/poll", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Equal(t, 10, len(data))
}

/**
 * TestPointHandler_Poll_ResponseStructure 測試輪詢回應結構
 */
func TestPointHandler_Poll_ResponseStructure(t *testing.T) {
	r := setupPointRouterWithExtended()

	// 建立設備
	newDevice := device.CreateDeviceRequest{
		Name:        "結構測試設備",
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

	var deviceResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &deviceResp)
	deviceID := deviceResp["data"].(map[string]interface{})["id"].(string)

	// 建立點位
	newPoint := point.CreatePointRequest{
		DeviceID: deviceID,
		Name:     "結構測試點位",
		Address:  "40001",
		DataType: "int16",
		Mode:     "read",
	}

	body, _ = json.Marshal(newPoint)
	req, _ = http.NewRequest("POST", "/datalink/points", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var pointResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &pointResp)
	pointID := pointResp["data"].(map[string]interface{})["id"].(string)

	// 輪詢點位
	req, _ = http.NewRequest("POST", "/datalink/points/"+pointID+"/poll", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})

	// 驗證欄位存在
	assert.Contains(t, data, "point_id")
	assert.Contains(t, data, "value")
	assert.Contains(t, data, "timestamp")
	assert.Contains(t, data, "quality")
	assert.Contains(t, data, "error")

	// 驗證數據類型
	assert.IsType(t, "", data["point_id"])
	assert.IsType(t, "", data["timestamp"])
	assert.IsType(t, float64(0), data["quality"])
}

/**
 * TestPointHandler_PollBatch_ResponseStructure 測試批量輪詢回應結構
 */
func TestPointHandler_PollBatch_ResponseStructure(t *testing.T) {
	r := setupPointRouterWithExtended()

	batchReq := map[string]interface{}{
		"point_ids": []string{},
	}

	body, _ := json.Marshal(batchReq)
	req, _ := http.NewRequest("POST", "/datalink/points/poll", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.IsType(t, []interface{}{}, data)

	// 驗證每個結果的結構
	for _, item := range data {
		result := item.(map[string]interface{})
		assert.Contains(t, result, "point_id")
		assert.Contains(t, result, "value")
		assert.Contains(t, result, "timestamp")
		assert.Contains(t, result, "quality")
		assert.Contains(t, result, "error")
	}
}
