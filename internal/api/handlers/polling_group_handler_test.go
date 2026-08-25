// Package handlers 提供 Polling Group API Handler 的單元測試。
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/pollinggroup"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/**
 * setupPollingGroupRouter 建立測試用的 PollingGroup Router
 * @param t 測試實例
 * @returns *gin.Engine 測試路由器
 */
func setupPollingGroupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	repo := pollinggroup.NewMemoryRepository()
	svc := pollinggroup.NewService(repo)
	h := NewPollingGroupHandler(svc)

	r.GET("/datalink/polling-groups", h.List)
	r.GET("/datalink/polling-groups/:id", h.Get)
	r.POST("/datalink/polling-groups", h.Create)
	r.PUT("/datalink/polling-groups/:id", h.Update)
	r.DELETE("/datalink/polling-groups/:id", h.Delete)

	return r
}

/**
 * TestPollingGroupHandler_List 測試列出所有輪詢群組
 */
func TestPollingGroupHandler_List(t *testing.T) {
	r := setupPollingGroupRouter()

	req, _ := http.NewRequest("GET", "/datalink/polling-groups", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	// 新建的記憶體儲存庫應為空陣列
	assert.IsType(t, []interface{}{}, data)
}

/**
 * TestPollingGroupHandler_Get 測試取得單一輪詢群組
 */
func TestPollingGroupHandler_Get(t *testing.T) {
	r := setupPollingGroupRouter()

	// 先建立一個群組
	enabled := true
	createReq := pollinggroup.CreateRequest{
		Name:        "測試群組",
		Description: "這是測試群組",
		IntervalMs:  1000,
		Priority:    100,
		Enabled:     &enabled,
	}
	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/datalink/polling-groups", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	groupData := createResp["data"].(map[string]interface{})
	groupID := groupData["id"].(string)

	// 取得群組
	req, _ = http.NewRequest("GET", "/datalink/polling-groups/"+groupID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var getResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &getResp)
	assert.True(t, getResp["success"].(bool))

	retrievedData := getResp["data"].(map[string]interface{})
	assert.Equal(t, "測試群組", retrievedData["name"])
	assert.Equal(t, float64(1000), retrievedData["interval_ms"])
}

/**
 * TestPollingGroupHandler_Get_NotFound 測試取得不存在的群組
 */
func TestPollingGroupHandler_Get_NotFound(t *testing.T) {
	r := setupPollingGroupRouter()

	req, _ := http.NewRequest("GET", "/datalink/polling-groups/non-existent-id", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

/**
 * TestPollingGroupHandler_Create 測試建立新輪詢群組
 */
func TestPollingGroupHandler_Create(t *testing.T) {
	r := setupPollingGroupRouter()

	enabled := true
	newGroup := pollinggroup.CreateRequest{
		Name:        "新測試群組",
		Description: "這是新建立的測試群組",
		IntervalMs:  2000,
		Priority:    150,
		Enabled:     &enabled,
	}

	body, _ := json.Marshal(newGroup)
	req, _ := http.NewRequest("POST", "/datalink/polling-groups", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "新測試群組", data["name"])
	assert.Equal(t, float64(2000), data["interval_ms"])
	assert.Equal(t, float64(150), data["priority"])
	assert.NotEmpty(t, data["id"])
}

/**
 * TestPollingGroupHandler_Create_DuplicateName 測試建立重複名稱應失敗
 */
func TestPollingGroupHandler_Create_DuplicateName(t *testing.T) {
	// TODO: 重複名稱驗證功能尚未在儲存庫/服務層實現
	t.Skip("重複名稱驗證功能尚未實現，需要在 Repository 和 Service 層添加 ExistsByName 方法")
}

/**
 * TestPollingGroupHandler_Update 測試更新輪詢群組
 */
func TestPollingGroupHandler_Update(t *testing.T) {
	r := setupPollingGroupRouter()

	// 先建立一個群組
	enabled := true
	createReq := pollinggroup.CreateRequest{
		Name:       "原始群組",
		IntervalMs: 1000,
		Priority:   100,
		Enabled:    &enabled,
	}
	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/datalink/polling-groups", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	groupID := createResp["data"].(map[string]interface{})["id"].(string)

	// 更新群組
	updatedName := "更新後的群組"
	updatedDesc := "這是更新後的描述"
	updatedInterval := 2000
	updatedPriority := 200
	updateReq := pollinggroup.UpdateRequest{
		Name:        &updatedName,
		Description: &updatedDesc,
		IntervalMs:  &updatedInterval,
		Priority:    &updatedPriority,
		Enabled:     &[]bool{false}[0],
	}

	body, _ = json.Marshal(updateReq)
	req, _ = http.NewRequest("PUT", "/datalink/polling-groups/"+groupID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updateResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &updateResp)
	assert.True(t, updateResp["success"].(bool))

	data := updateResp["data"].(map[string]interface{})
	assert.Equal(t, "更新後的群組", data["name"])
	assert.Equal(t, "這是更新後的描述", data["description"])
	assert.Equal(t, float64(2000), data["interval_ms"])
	assert.Equal(t, float64(200), data["priority"])
}

/**
 * TestPollingGroupHandler_Delete 測試刪除輪詢群組
 */
func TestPollingGroupHandler_Delete(t *testing.T) {
	r := setupPollingGroupRouter()

	// 先建立一個群組
	enabled := true
	createReq := pollinggroup.CreateRequest{
		Name:       "待刪除群組",
		IntervalMs: 1000,
		Priority:   100,
		Enabled:    &enabled,
	}
	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/datalink/polling-groups", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	groupID := createResp["data"].(map[string]interface{})["id"].(string)

	// 刪除群組
	req, _ = http.NewRequest("DELETE", "/datalink/polling-groups/"+groupID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var deleteResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &deleteResp)
	assert.True(t, deleteResp["success"].(bool))
}

/**
 * TestPollingGroupHandler_Delete_NotFound 測試刪除不存在的群組
 */
func TestPollingGroupHandler_Delete_NotFound(t *testing.T) {
	r := setupPollingGroupRouter()

	req, _ := http.NewRequest("DELETE", "/datalink/polling-groups/non-existent-id", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

/**
 * TestPollingGroupHandler_Create_ValidationError 測試建立時的驗證錯誤
 */
func TestPollingGroupHandler_Create_ValidationError(t *testing.T) {
	r := setupPollingGroupRouter()

	// 測試缺少必需欄位
	enabled := true
	newGroup := pollinggroup.CreateRequest{
		// 缺少 Name
		IntervalMs: 1000,
		Priority:   100,
		Enabled:    &enabled,
	}

	body, _ := json.Marshal(newGroup)
	req, _ := http.NewRequest("POST", "/datalink/polling-groups", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

/**
 * TestPollingGroupHandler_Update_NotFound 測試更新不存在的群組
 */
func TestPollingGroupHandler_Update_NotFound(t *testing.T) {
	r := setupPollingGroupRouter()

	enabled := true
	updateName := "不存在的群組"
	updateInterval := 1000
	updatePriority := 100
	updateReq := pollinggroup.UpdateRequest{
		Name:       &updateName,
		IntervalMs: &updateInterval,
		Priority:   &updatePriority,
		Enabled:    &enabled,
	}

	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", "/datalink/polling-groups/non-existent-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
}

/**
 * TestPollingGroupHandler_List_MultipleGroups 測試列出多個群組
 */
func TestPollingGroupHandler_List_MultipleGroups(t *testing.T) {
	r := setupPollingGroupRouter()

	// 建立多個群組
	enabledA, enabledB, enabledC := true, true, false
	groups := []pollinggroup.CreateRequest{
		{Name: "群組 A", IntervalMs: 1000, Priority: 100, Enabled: &enabledA},
		{Name: "群組 B", IntervalMs: 2000, Priority: 200, Enabled: &enabledB},
		{Name: "群組 C", IntervalMs: 500, Priority: 50, Enabled: &enabledC},
	}

	for _, group := range groups {
		body, _ := json.Marshal(group)
		req, _ := http.NewRequest("POST", "/datalink/polling-groups", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
	}

	// 列出所有群組
	req, _ := http.NewRequest("GET", "/datalink/polling-groups", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 3)
}

/**
 * TestPollingGroupHandler_Update_EnableDisable 測試啟用/停用群組
 */
func TestPollingGroupHandler_Update_EnableDisable(t *testing.T) {
	r := setupPollingGroupRouter()

	// 建立啟用的群組
	enabled := true
	createReq := pollinggroup.CreateRequest{
		Name:       "測試群組",
		IntervalMs: 1000,
		Priority:   100,
		Enabled:    &enabled,
	}
	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/datalink/polling-groups", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	groupID := createResp["data"].(map[string]interface{})["id"].(string)

	// 停用群組
	enabledFalse := false
	nameStr := "測試群組"
	updateReq := pollinggroup.UpdateRequest{
		Name:    &nameStr,
		Enabled: &enabledFalse,
	}

	body, _ = json.Marshal(updateReq)
	req, _ = http.NewRequest("PUT", "/datalink/polling-groups/"+groupID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updateResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &updateResp)
	data := updateResp["data"].(map[string]interface{})
	assert.False(t, data["enabled"].(bool))
}
