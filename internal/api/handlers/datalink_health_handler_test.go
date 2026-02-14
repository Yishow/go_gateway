// Package handlers 提供 Datalink Health API Handler 的單元測試。
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/**
 * setupDatalinkHealthRouter 建立測試用的 Datalink Health Router
 * @param t 測試實例
 * @returns *gin.Engine 測試路由器
 */
func setupDatalinkHealthRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := NewDatalinkHealthHandler()
	r.GET("/datalink/health", h.Check)
	return r
}

/**
 * TestDatalinkHealthHandler_Check 測試健康檢查端點
 */
func TestDatalinkHealthHandler_Check(t *testing.T) {
	r := setupDatalinkHealthRouter()

	req, _ := http.NewRequest("GET", "/datalink/health", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

/**
 * TestDatalinkHealthHandler_Check_Structure 測試健康檢查回應結構
 */
func TestDatalinkHealthHandler_Check_Structure(t *testing.T) {
	r := setupDatalinkHealthRouter()

	req, _ := http.NewRequest("GET", "/datalink/health", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// 驗證頂層結構
	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["data"])

	// 驗證 data 結構
	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "status")
	assert.Contains(t, data, "service")
	assert.Contains(t, data, "timestamp")
	assert.Contains(t, data, "version")
}

/**
 * TestDatalinkHealthHandler_Check_Status 測試健康狀態
 */
func TestDatalinkHealthHandler_Check_Status(t *testing.T) {
	r := setupDatalinkHealthRouter()

	req, _ := http.NewRequest("GET", "/datalink/health", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})

	assert.Equal(t, "healthy", data["status"])
}

/**
 * TestDatalinkHealthHandler_Check_Service 測試服務名稱
 */
func TestDatalinkHealthHandler_Check_Service(t *testing.T) {
	r := setupDatalinkHealthRouter()

	req, _ := http.NewRequest("GET", "/datalink/health", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})

	assert.Equal(t, "datalink", data["service"])
}

/**
 * TestDatalinkHealthHandler_Check_Version 測試版本號
 */
func TestDatalinkHealthHandler_Check_Version(t *testing.T) {
	r := setupDatalinkHealthRouter()

	req, _ := http.NewRequest("GET", "/datalink/health", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})

	assert.Equal(t, "1.0.0", data["version"])
}

/**
 * TestDatalinkHealthHandler_Check_Timestamp 測試時間戳
 */
func TestDatalinkHealthHandler_Check_Timestamp(t *testing.T) {
	r := setupDatalinkHealthRouter()

	req, _ := http.NewRequest("GET", "/datalink/health", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})

	// 驗證 timestamp 存在且為字串
	timestamp, ok := data["timestamp"].(string)
	assert.True(t, ok, "timestamp should be a string")
	assert.NotEmpty(t, timestamp)

	// 驗證可以解析為 ISO 8601 格式
	_, err := time.Parse(time.RFC3339, timestamp)
	assert.NoError(t, err, "timestamp should be in RFC3339 format")
}

/**
 * TestDatalinkHealthHandler_Check_ContentType 測試回應 Content-Type
 */
func TestDatalinkHealthHandler_Check_ContentType(t *testing.T) {
	r := setupDatalinkHealthRouter()

	req, _ := http.NewRequest("GET", "/datalink/health", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

/**
 * TestDatalinkHealthHandler_Check_ResponseTime 測試回應時間
 */
func TestDatalinkHealthHandler_Check_ResponseTime(t *testing.T) {
	r := setupDatalinkHealthRouter()

	start := time.Now()
	req, _ := http.NewRequest("GET", "/datalink/health", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	duration := time.Since(start)

	assert.Equal(t, http.StatusOK, w.Code)

	// 驗證回應時間合理（應小於 100ms）
	assert.Less(t, duration.Milliseconds(), int64(100), "Health check should respond quickly")
}

/**
 * TestDatalinkHealthHandler_Check_MultipleRequests 測試多次請求一致性
 */
func TestDatalinkHealthHandler_Check_MultipleRequests(t *testing.T) {
	r := setupDatalinkHealthRouter()

	var responses []map[string]interface{}

	// 發送 5 次請求
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", "/datalink/health", http.NoBody)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		responses = append(responses, response)

		// 短暫等待以確保時間戳不同
		time.Sleep(10 * time.Millisecond)
	}

	// 驗證所有回應結構一致
	for i, response := range responses {
		assert.True(t, response["success"].(bool), "Request %d should be successful", i)
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "healthy", data["status"], "Request %d should have healthy status", i)
		assert.Equal(t, "datalink", data["service"], "Request %d should have correct service name", i)
		assert.Equal(t, "1.0.0", data["version"], "Request %d should have correct version", i)
	}

	// 驗證時間戳不同
	timestamps := make(map[string]bool)
	for _, response := range responses {
		data := response["data"].(map[string]interface{})
		timestamp := data["timestamp"].(string)
		timestamps[timestamp] = true
	}
	assert.Greater(t, len(timestamps), 1, "Timestamps should be different for each request")
}

/**
 * TestDatalinkHealthHandler_Check_InvalidMethod 測試不支援的 HTTP 方法
 */
func TestDatalinkHealthHandler_Check_InvalidMethod(t *testing.T) {
	r := setupDatalinkHealthRouter()

	// 測試 POST 方法（不支援）
	req, _ := http.NewRequest("POST", "/datalink/health", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Gin 默認返回 404 或 405，取決於配置
	assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusMethodNotAllowed)
}

/**
 * TestDatalinkHealthHandler_Check_InvalidPath 測試不存在的路徑
 */
func TestDatalinkHealthHandler_Check_InvalidPath(t *testing.T) {
	r := setupDatalinkHealthRouter()

	req, _ := http.NewRequest("GET", "/datalink/invalid", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

/**
 * TestDatalinkHealthHandler_Check_ResponseBody 測試回應內容完整性
 */
func TestDatalinkHealthHandler_Check_ResponseBody(t *testing.T) {
	r := setupDatalinkHealthRouter()

	req, _ := http.NewRequest("GET", "/datalink/health", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	body := w.Body.String()
	assert.NotEmpty(t, body, "Response body should not be empty")

	// 驗證包含預期欄位
	assert.Contains(t, body, "success")
	assert.Contains(t, body, "data")
	assert.Contains(t, body, "status")
	assert.Contains(t, body, "service")
	assert.Contains(t, body, "timestamp")
	assert.Contains(t, body, "version")
}

/**
 * TestDatalinkHealthHandler_Check_Headers 測試回應標頭
 */
func TestDatalinkHealthHandler_Check_Headers(t *testing.T) {
	r := setupDatalinkHealthRouter()

	req, _ := http.NewRequest("GET", "/datalink/health", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 驗證 Content-Type 標頭
	assert.NotEmpty(t, w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

/**
 * TestDatalinkHealthHandler_Check_JSONFormat 測試 JSON 格式正確性
 */
func TestDatalinkHealthHandler_Check_JSONFormat(t *testing.T) {
	r := setupDatalinkHealthRouter()

	req, _ := http.NewRequest("GET", "/datalink/health", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 驗證可以解析為 JSON
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")

	// 驗證 JSON 結構
	assert.IsType(t, map[string]interface{}{}, response)
	assert.IsType(t, true, response["success"])
	assert.IsType(t, map[string]interface{}{}, response["data"])

	// 驗證 data 內部結構
	data := response["data"].(map[string]interface{})
	assert.IsType(t, "", data["status"])
	assert.IsType(t, "", data["service"])
	assert.IsType(t, "", data["timestamp"])
	assert.IsType(t, "", data["version"])
}

/**
 * TestDatalinkHealthHandler_Check_ConcurrentRequests 測試並發請求
 */
func TestDatalinkHealthHandler_Check_ConcurrentRequests(t *testing.T) {
	r := setupDatalinkHealthRouter()

	numRequests := 10
	results := make(chan bool, numRequests)

	// 發送並發請求
	for i := 0; i < numRequests; i++ {
		go func() {
			req, _ := http.NewRequest("GET", "/datalink/health", http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			success := response["success"].(bool)
			results <- success
		}()
	}

	// 等待所有請求完成
	successCount := 0
	for i := 0; i < numRequests; i++ {
		if <-results {
			successCount++
		}
	}

	// 驗證所有請求都成功
	assert.Equal(t, numRequests, successCount, "All concurrent requests should succeed")
}

/**
 * TestDatalinkHealthHandler_Check_TimestampPrecision 測試時間戳精度
 */
func TestDatalinkHealthHandler_Check_TimestampPrecision(t *testing.T) {
	r := setupDatalinkHealthRouter()

	req, _ := http.NewRequest("GET", "/datalink/health", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	timestampStr := data["timestamp"].(string)

	// 驗證可以解析為 RFC3339（包含秒和毫秒）
	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	assert.NoError(t, err)

	// 驗證時間戳是最近的（1 秒內）
	now := time.Now()
	diff := now.Sub(timestamp)
	assert.Less(t, diff.Abs(), time.Second, "Timestamp should be within 1 second")
}
