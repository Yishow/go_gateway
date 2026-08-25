// Package handlers 提供 Mapping Handler 擴展 API 的單元測試。
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

/**
 * setupMappingRouterWithExtended 建立測試用的 Mapping Router（包含擴展端點）
 * @param t 測試實例
 * @returns *gin.Engine 測試路由器
 */
func setupMappingRouterWithExtended() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	repo := mapping.NewMemoryRepository()
	svc := mapping.NewService(repo)
	svc.SetTagResolver(func(ctx context.Context, tagID string) (*schema.Tag, error) {
		return &schema.Tag{ID: tagID, DataType: schema.DataTypeFloat64}, nil
	})
	h := NewMappingHandler(svc)

	// 基礎端點
	r.GET("/datalink/mappings", h.List)
	r.POST("/datalink/mappings", h.Create)

	// 擴展端點
	r.POST("/datalink/mappings/preview", h.Preview)
	r.POST("/datalink/mappings/validate-pipeline", h.ValidatePipeline)

	return r
}

/**
 * TestMappingHandler_ValidatePipeline 測試驗證轉換管線
 */
func TestMappingHandler_ValidatePipeline(t *testing.T) {
	r := setupMappingRouterWithExtended()

	// 測試有效的管線
	validPipeline := []schema.TransformStep{
		{Type: schema.TransformScale, Params: map[string]interface{}{"multiplier": 10, "offset": 0}},
		{Type: schema.TransformScale, Params: map[string]interface{}{"multiplier": 1, "offset": 5}},
	}

	reqBody := map[string]interface{}{
		"pipeline": validPipeline,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/datalink/mappings/validate-pipeline", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.True(t, data["valid"].(bool))
}

/**
 * TestMappingHandler_ValidatePipeline_Empty 測試驗證空管線
 */
func TestMappingHandler_ValidatePipeline_Empty(t *testing.T) {
	r := setupMappingRouterWithExtended()

	// 空管線是有效的
	reqBody := map[string]interface{}{
		"pipeline": []schema.TransformStep{},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/datalink/mappings/validate-pipeline", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.True(t, data["valid"].(bool))
}

/**
 * TestMappingHandler_ValidatePipeline_InvalidType 測試驗證無效的轉換類型
 */
func TestMappingHandler_ValidatePipeline_InvalidType(t *testing.T) {
	r := setupMappingRouterWithExtended()

	// 測試無效的轉換類型
	invalidPipeline := []schema.TransformStep{
		{Type: "invalid_type", Params: map[string]interface{}{}},
	}

	reqBody := map[string]interface{}{
		"pipeline": invalidPipeline,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/datalink/mappings/validate-pipeline", bytes.NewBuffer(body))
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
 * TestMappingHandler_ValidatePipeline_ScaleMissingParams 測試 Scale 缺少參數
 */
func TestMappingHandler_ValidatePipeline_ScaleMissingParams(t *testing.T) {
	r := setupMappingRouterWithExtended()

	// Scale 類型必須有 params
	invalidPipeline := []schema.TransformStep{
		{Type: schema.TransformScale, Params: nil},
	}

	reqBody := map[string]interface{}{
		"pipeline": invalidPipeline,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/datalink/mappings/validate-pipeline", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.False(t, data["valid"].(bool))
	assert.Contains(t, data["error"].(string), "缺少 params")
}

/**
 * TestMappingHandler_ValidatePipeline_MultipleValidTypes 測試多個有效轉換類型
 */
func TestMappingHandler_ValidatePipeline_MultipleValidTypes(t *testing.T) {
	r := setupMappingRouterWithExtended()

	// 測試所有支持的轉換類型
	validPipeline := []schema.TransformStep{
		{Type: schema.TransformDecode, Params: map[string]interface{}{"type": "int16"}},
		{Type: schema.TransformCast, Params: map[string]interface{}{"type": "float64"}},
		{Type: schema.TransformScale, Params: map[string]interface{}{"multiplier": 10, "offset": 0}},
		{Type: schema.TransformScale, Params: map[string]interface{}{"multiplier": 1, "offset": 5}},
		{Type: schema.TransformLookup, Params: map[string]interface{}{"table": "test"}},
		{Type: schema.TransformConditional, Params: map[string]interface{}{"condition": "x > 0"}},
		{Type: schema.TransformFormula, Params: map[string]interface{}{"formula": "x * 2"}},
	}

	reqBody := map[string]interface{}{
		"pipeline": validPipeline,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/datalink/mappings/validate-pipeline", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.True(t, data["valid"].(bool))
}

/**
 * TestMappingHandler_Preview 測試預覽轉換管線
 */
func TestMappingHandler_Preview(t *testing.T) {
	r := setupMappingRouterWithExtended()

	// 測試預覽
	pipeline := []schema.TransformStep{
		{Type: schema.TransformScale, Params: map[string]interface{}{"multiplier": 10, "offset": 0}},
	}

	reqBody := map[string]interface{}{
		"raw_value":          5,
		"transform_pipeline": pipeline,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/datalink/mappings/preview", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(5), data["raw_value"])
	// 5 * 10 = 50
	assert.Equal(t, float64(50), data["final_value"])
	assert.Contains(t, data, "step_results")
}

/**
 * TestMappingHandler_Preview_ComplexPipeline 測試複雜管線預覽
 */
func TestMappingHandler_Preview_ComplexPipeline(t *testing.T) {
	r := setupMappingRouterWithExtended()

	// 複雜管線
	pipeline := []schema.TransformStep{
		{Type: schema.TransformScale, Params: map[string]interface{}{"multiplier": 10, "offset": 0}},
		{Type: schema.TransformScale, Params: map[string]interface{}{"multiplier": 1, "offset": 5}},
		{Type: schema.TransformFormula, Params: map[string]interface{}{"formula": "x > 0 ? x : 0"}},
		{Type: schema.TransformFormula, Params: map[string]interface{}{"formula": "x < 100 ? x : 100"}},
	}

	reqBody := map[string]interface{}{
		"raw_value":          5,
		"transform_pipeline": pipeline,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/datalink/mappings/preview", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})

	// 5 * 10 + 5 = 55, clamp to [0, 100] = 55, round to 2 decimals = 55.00
	assert.Equal(t, float64(55), data["final_value"])
	assert.Contains(t, data, "step_results")
}

/**
 * TestMappingHandler_Preview_InvalidPipeline 測試無效管線預覽
 */
func TestMappingHandler_Preview_InvalidPipeline(t *testing.T) {
	r := setupMappingRouterWithExtended()

	// 無效的管線格式
	invalidPipeline := []schema.TransformStep{
		{Type: "invalid_type", Params: map[string]interface{}{}},
	}

	reqBody := map[string]interface{}{
		"raw_value":          5,
		"transform_pipeline": invalidPipeline,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/datalink/mappings/preview", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
	apiError := response["error"].(map[string]interface{})
	assert.Equal(t, ErrCodePreviewInvalidRequest, apiError["code"])
	assert.NotContains(t, w.Body.String(), "invalid_type")
}

/**
 * TestMappingHandler_Preview_MissingField 測試缺少必需欄位
 */
func TestMappingHandler_Preview_MissingField(t *testing.T) {
	r := setupMappingRouterWithExtended()

	// 缺少 transform_pipeline
	reqBody := map[string]interface{}{
		"raw_value": 5,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/datalink/mappings/preview", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

/**
 * TestMappingHandler_ValidatePipeline_MissingField 測試驗證管線缺少欄位
 */
func TestMappingHandler_ValidatePipeline_MissingField(t *testing.T) {
	r := setupMappingRouterWithExtended()

	// 缺少 pipeline 欄位
	reqBody := map[string]interface{}{}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/datalink/mappings/validate-pipeline", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

/**
 * TestMappingHandler_Preview_DifferentValueTypes 測試不同數值類型
 */
func TestMappingHandler_Preview_DifferentValueTypes(t *testing.T) {
	r := setupMappingRouterWithExtended()

	testCases := []struct {
		rawValue interface{}
		expected interface{}
	}{
		{5, float64(5)},
		{3.14, float64(3.14)},
		{"100", "100"},
		{true, true},
	}

	for _, tc := range testCases {
		pipeline := []schema.TransformStep{}
		if _, ok := tc.rawValue.(float64); ok {
			pipeline = []schema.TransformStep{
				{Type: schema.TransformScale, Params: map[string]interface{}{"multiplier": 2, "offset": 0}},
			}
		}

		reqBody := map[string]interface{}{
			"raw_value":          tc.rawValue,
			"transform_pipeline": pipeline,
		}

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/datalink/mappings/preview", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))
	}
}

/**
 * TestMappingHandler_ValidatePipeline_AllTransformTypes 測試所有轉換類型
 */
func TestMappingHandler_ValidatePipeline_AllTransformTypes(t *testing.T) {
	r := setupMappingRouterWithExtended()

	// 測試所有支持的轉換類型
	transformTypes := []schema.TransformType{
		schema.TransformDecode,
		schema.TransformCast,
		schema.TransformScale,
		schema.TransformLookup,
		schema.TransformConditional,
		schema.TransformFormula,
	}

	for _, transformType := range transformTypes {
		pipeline := []schema.TransformStep{
			{Type: transformType, Params: map[string]interface{}{"test": true}},
		}

		// Scale 類型需要特定的參數
		if transformType == schema.TransformScale {
			pipeline = []schema.TransformStep{
				{Type: transformType, Params: map[string]interface{}{"multiplier": 1, "offset": 0}},
			}
		}

		reqBody := map[string]interface{}{
			"pipeline": pipeline,
		}

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/datalink/mappings/validate-pipeline", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Transform type %s should be valid", transformType)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool), "Transform type %s should return success", transformType)

		data := response["data"].(map[string]interface{})
		assert.True(t, data["valid"].(bool), "Transform type %s should be valid", transformType)
	}
}

/**
 * TestMappingHandler_Preview_ResponseStructure 測試預覽回應結構
 */
func TestMappingHandler_Preview_ResponseStructure(t *testing.T) {
	r := setupMappingRouterWithExtended()

	pipeline := []schema.TransformStep{
		{Type: schema.TransformScale, Params: map[string]interface{}{"multiplier": 10, "offset": 0}},
	}

	reqBody := map[string]interface{}{
		"raw_value":          5,
		"transform_pipeline": pipeline,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/datalink/mappings/preview", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// 驗證頂層結構
	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["data"])

	// 驗證 data 結構
	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "raw_value")
	assert.Contains(t, data, "final_value")
	assert.Contains(t, data, "step_results")

	// 驗證 step_results 結構
	stepResults := data["step_results"].([]interface{})
	assert.IsType(t, []interface{}{}, stepResults)
}

/**
 * TestMappingHandler_ValidatePipeline_ResponseStructure 測試驗證回應結構
 */
func TestMappingHandler_ValidatePipeline_ResponseStructure(t *testing.T) {
	r := setupMappingRouterWithExtended()

	pipeline := []schema.TransformStep{
		{Type: schema.TransformScale, Params: map[string]interface{}{"multiplier": 10, "offset": 0}},
	}

	reqBody := map[string]interface{}{
		"pipeline": pipeline,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/datalink/mappings/validate-pipeline", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// 驗證頂層結構
	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["data"])

	// 驗證 data 結構
	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "valid")
	assert.Contains(t, data, "error")
}

/**
 * TestMappingHandler_Preview_WithCast 測試包含 Cast 的管線
 */
func TestMappingHandler_Preview_WithCast(t *testing.T) {
	r := setupMappingRouterWithExtended()

	pipeline := []schema.TransformStep{
		{Type: schema.TransformCast, Params: map[string]interface{}{"type": "float64"}},
		{Type: schema.TransformScale, Params: map[string]interface{}{"multiplier": 2, "offset": 0}},
	}

	reqBody := map[string]interface{}{
		"raw_value":          5,
		"transform_pipeline": pipeline,
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/datalink/mappings/preview", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	// cast(5, float64) * 2 = 10.0
	assert.Equal(t, float64(10), data["final_value"])
}
