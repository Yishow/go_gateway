// Package handlers 提供 Tag Handler 擴展 API 的單元測試。
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/tag"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

/**
 * setupTagRouter 建立測試用的 Tag Router（包含擴展端點）
 * @param t 測試實例
 * @returns *gin.Engine 測試路由器
 */
func setupTagRouterWithExtended() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := NewTagHandler()

	// 基礎端點
	r.GET("/datalink/tags", h.List)
	r.POST("/datalink/tags", h.Create)

	// 擴展端點
	r.POST("/datalink/tags/:id/activate", h.Activate)
	r.POST("/datalink/tags/:id/retire", h.Retire)
	r.POST("/datalink/tags/batch", h.BatchCreate)
	r.POST("/datalink/tags/validate-key", h.ValidateKey)

	return r
}

/**
 * TestTagHandler_Activate 測試啟用標籤
 */
func TestTagHandler_Activate(t *testing.T) {
	r := setupTagRouterWithExtended()

	// 先建立一個標籤
	newTag := tag.CreateTagRequest{
		Key:         "test.activate.tag",
		DisplayName: "測試啟用標籤",
		DataType:    "float64",
		Unit:        "°C",
	}

	body, _ := json.Marshal(newTag)
	req, _ := http.NewRequest("POST", "/datalink/tags", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	tagID := createResp["data"].(map[string]interface{})["id"].(string)

	// 啟用標籤
	req, _ = http.NewRequest("POST", "/datalink/tags/"+tagID+"/activate", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var activateResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &activateResp)
	assert.True(t, activateResp["success"].(bool))

	data := activateResp["data"].(map[string]interface{})
	assert.Equal(t, "active", data["status"])
}

/**
 * TestTagHandler_Activate_NotFound 測試啟用不存在的標籤
 */
func TestTagHandler_Activate_NotFound(t *testing.T) {
	r := setupTagRouterWithExtended()

	req, _ := http.NewRequest("POST", "/datalink/tags/non-existent-id/activate", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

/**
 * TestTagHandler_Retire 測試退役標籤
 */
func TestTagHandler_Retire(t *testing.T) {
	r := setupTagRouterWithExtended()

	// 先建立一個標籤
	newTag := tag.CreateTagRequest{
		Key:         "test.retire.tag",
		DisplayName: "測試退役標籤",
		DataType:    "float64",
		Unit:        "°C",
	}

	body, _ := json.Marshal(newTag)
	req, _ := http.NewRequest("POST", "/datalink/tags", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	tagID := createResp["data"].(map[string]interface{})["id"].(string)

	// 退役標籤
	req, _ = http.NewRequest("POST", "/datalink/tags/"+tagID+"/retire", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var retireResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &retireResp)
	assert.True(t, retireResp["success"].(bool))

	data := retireResp["data"].(map[string]interface{})
	assert.Equal(t, "retired", data["status"])
}

/**
 * TestTagHandler_Retire_NotFound 測試退役不存在的標籤
 */
func TestTagHandler_Retire_NotFound(t *testing.T) {
	r := setupTagRouterWithExtended()

	req, _ := http.NewRequest("POST", "/datalink/tags/non-existent-id/retire", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

/**
 * TestTagHandler_BatchCreate 測試批量建立標籤
 */
func TestTagHandler_BatchCreate(t *testing.T) {
	r := setupTagRouterWithExtended()

	// 準備批量建立請求
	batchReq := map[string]interface{}{
		"tags": []tag.CreateTagRequest{
			{
				Key:         "batch.tag.1",
				DisplayName: "批量標籤 1",
				DataType:    "float64",
				Unit:        "°C",
			},
			{
				Key:         "batch.tag.2",
				DisplayName: "批量標籤 2",
				DataType:    "int32",
				Unit:        "rpm",
			},
			{
				Key:         "batch.tag.3",
				DisplayName: "批量標籤 3",
				DataType:    "bool",
				Unit:        "",
			},
		},
	}

	body, _ := json.Marshal(batchReq)
	req, _ := http.NewRequest("POST", "/datalink/tags/batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	created := data["created"].([]interface{})
	assert.Equal(t, 3, len(created))

	errors := data["errors"].([]interface{})
	assert.Equal(t, 0, len(errors))
}

/**
 * TestTagHandler_BatchCreate_PartialSuccess 測試批量建立部分成功
 */
func TestTagHandler_BatchCreate_PartialSuccess(t *testing.T) {
	r := setupTagRouterWithExtended()

	// 先建立一個標籤（造成重複）
	newTag := tag.CreateTagRequest{
		Key:         "duplicate.tag",
		DisplayName: "重複標籤",
		DataType:    "float64",
		Unit:        "°C",
	}

	body, _ := json.Marshal(newTag)
	req, _ := http.NewRequest("POST", "/datalink/tags", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 準備批量建立請求（包含重複）
	batchReq := map[string]interface{}{
		"tags": []tag.CreateTagRequest{
			{
				Key:         "batch.unique.tag",
				DisplayName: "唯一標籤",
				DataType:    "float64",
				Unit:        "°C",
			},
			{
				Key:         "duplicate.tag", // 會重複
				DisplayName: "重複標籤",
				DataType:    "float64",
				Unit:        "°C",
			},
		},
	}

	body, _ = json.Marshal(batchReq)
	req, _ = http.NewRequest("POST", "/datalink/tags/batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	created := data["created"].([]interface{})
	assert.Equal(t, 1, len(created))

	errors := data["errors"].([]interface{})
	assert.Equal(t, 1, len(errors))
}

/**
 * TestTagHandler_BatchCreate_Empty 測試批量建立空陣列
 */
func TestTagHandler_BatchCreate_Empty(t *testing.T) {
	r := setupTagRouterWithExtended()

	batchReq := map[string]interface{}{
		"tags": []tag.CreateTagRequest{},
	}

	body, _ := json.Marshal(batchReq)
	req, _ := http.NewRequest("POST", "/datalink/tags/batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	created := data["created"].([]interface{})
	assert.Equal(t, 0, len(created))

	errors := data["errors"].([]interface{})
	assert.Equal(t, 0, len(errors))
}

/**
 * TestTagHandler_ValidateKey 測試驗證標籤鍵
 */
func TestTagHandler_ValidateKey(t *testing.T) {
	r := setupTagRouterWithExtended()

	// 測試有效的標籤鍵
	validateReq := map[string]interface{}{
		"key": "test.valid.tag",
	}

	body, _ := json.Marshal(validateReq)
	req, _ := http.NewRequest("POST", "/datalink/tags/validate-key", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.True(t, data["valid"].(bool))
	assert.Equal(t, "test.valid.tag", data["normalized"])
	assert.False(t, data["exists"].(bool))
}

/**
 * TestTagHandler_ValidateKey_Invalid 測試驗證無效的標籤鍵
 */
func TestTagHandler_ValidateKey_Invalid(t *testing.T) {
	r := setupTagRouterWithExtended()

	// 測試無效的標籤鍵（包含非法字元）
	validateReq := map[string]interface{}{
		"key": "invalid@tag#key",
	}

	body, _ := json.Marshal(validateReq)
	req, _ := http.NewRequest("POST", "/datalink/tags/validate-key", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.False(t, data["valid"].(bool))
	assert.NotEmpty(t, data["error"])
}

/**
 * TestTagHandler_ValidateKey_Exists 測試驗證已存在的標籤鍵
 */
func TestTagHandler_ValidateKey_Exists(t *testing.T) {
	r := setupTagRouterWithExtended()

	// 先建立一個標籤
	newTag := tag.CreateTagRequest{
		Key:         "existing.tag.key",
		DisplayName: "已存在的標籤",
		DataType:    "float64",
		Unit:        "°C",
	}

	body, _ := json.Marshal(newTag)
	req, _ := http.NewRequest("POST", "/datalink/tags", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// 驗證已存在的標籤鍵
	validateReq := map[string]interface{}{
		"key": "existing.tag.key",
	}

	body, _ = json.Marshal(validateReq)
	req, _ = http.NewRequest("POST", "/datalink/tags/validate-key", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.True(t, data["valid"].(bool))
	assert.Equal(t, "existing.tag.key", data["normalized"])
	assert.True(t, data["exists"].(bool))
}

/**
 * TestTagHandler_ValidateKey_Normalization 測試標籤鍵正規化
 */
func TestTagHandler_ValidateKey_Normalization(t *testing.T) {
	r := setupTagRouterWithExtended()

	testCases := []struct {
		input      string
		normalized string
	}{
		{"Test.Tag.Key", "test.tag.key"},
		{"  spaced.tag  ", "spaced.tag"},
		{"Mixed.Case.Tag", "mixed.case.tag"},
		{"UPPER.CASE.TAG", "upper.case.tag"},
	}

	for _, tc := range testCases {
		validateReq := map[string]interface{}{
			"key": tc.input,
		}

		body, _ := json.Marshal(validateReq)
		req, _ := http.NewRequest("POST", "/datalink/tags/validate-key", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})
		assert.Equal(t, tc.normalized, data["normalized"], "Input: %s", tc.input)
	}
}

/**
 * TestTagHandler_ActivateThenRetire 測試啟用後退役標籤
 */
func TestTagHandler_ActivateThenRetire(t *testing.T) {
	r := setupTagRouterWithExtended()

	// 建立標籤
	newTag := tag.CreateTagRequest{
		Key:         "test.cycle.tag",
		DisplayName: "測試循環標籤",
		DataType:    "float64",
		Unit:        "°C",
	}

	body, _ := json.Marshal(newTag)
	req, _ := http.NewRequest("POST", "/datalink/tags", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	tagID := createResp["data"].(map[string]interface{})["id"].(string)

	// 啟用標籤
	req, _ = http.NewRequest("POST", "/datalink/tags/"+tagID+"/activate", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var activateResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &activateResp)
	data := activateResp["data"].(map[string]interface{})
	assert.Equal(t, "active", data["status"])

	// 退役標籤
	req, _ = http.NewRequest("POST", "/datalink/tags/"+tagID+"/retire", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var retireResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &retireResp)
	data = retireResp["data"].(map[string]interface{})
	assert.Equal(t, "retired", data["status"])
}

/**
 * TestTagHandler_BatchCreate_MissingField 測試批量建立缺少必需欄位
 */
func TestTagHandler_BatchCreate_MissingField(t *testing.T) {
	r := setupTagRouterWithExtended()

	// 缺少 tags 欄位
	batchReq := map[string]interface{}{}

	body, _ := json.Marshal(batchReq)
	req, _ := http.NewRequest("POST", "/datalink/tags/batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

/**
 * TestTagHandler_ValidateKey_MissingField 測試驗證標籤鍵缺少欄位
 */
func TestTagHandler_ValidateKey_MissingField(t *testing.T) {
	r := setupTagRouterWithExtended()

	// 缺少 key 欄位
	validateReq := map[string]interface{}{}

	body, _ := json.Marshal(validateReq)
	req, _ := http.NewRequest("POST", "/datalink/tags/validate-key", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}
