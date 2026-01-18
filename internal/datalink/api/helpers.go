package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

// =============================================================================
// JSON 回應輔助函數
// =============================================================================

// APIResponse 標準 API 回應結構
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *APIMeta    `json:"meta,omitempty"`
}

// APIError 錯誤回應結構
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// APIMeta 分頁/元資料結構
type APIMeta struct {
	Total  int `json:"total,omitempty"`
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// writeJSON 寫入 JSON 回應
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	resp := APIResponse{
		Success: status >= 200 && status < 300,
		Data:    data,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

// writeJSONWithMeta 寫入帶分頁資訊的 JSON 回應
func writeJSONWithMeta(w http.ResponseWriter, status int, data interface{}, meta APIMeta) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	resp := APIResponse{
		Success: status >= 200 && status < 300,
		Data:    data,
		Meta:    &meta,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

// writeError 寫入錯誤回應
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	code := httpStatusToCode(status)
	resp := APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
}

// writeErrorWithDetails 寫入帶詳細資訊的錯誤回應
func writeErrorWithDetails(w http.ResponseWriter, status int, message, details string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	code := httpStatusToCode(status)
	resp := APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
}

// httpStatusToCode 將 HTTP 狀態碼轉換為錯誤代碼
func httpStatusToCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusMethodNotAllowed:
		return "METHOD_NOT_ALLOWED"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusUnprocessableEntity:
		return "VALIDATION_ERROR"
	case http.StatusInternalServerError:
		return "INTERNAL_ERROR"
	case http.StatusServiceUnavailable:
		return "SERVICE_UNAVAILABLE"
	default:
		return "ERROR"
	}
}

// =============================================================================
// 路徑解析輔助函數
// =============================================================================

// parseIDAndAction 從路徑解析 ID 和操作
// 例如: "/api/v1/datalink/devices/abc123/test" -> id="abc123", action="test"
func parseIDAndAction(path, prefix string) (id, action string) {
	// 移除前綴
	path = strings.TrimPrefix(path, prefix)
	path = strings.Trim(path, "/")

	// 分割路徑
	parts := strings.Split(path, "/")
	if len(parts) >= 1 {
		id = parts[0]
	}
	if len(parts) >= 2 {
		action = parts[1]
	}

	return id, action
}

// =============================================================================
// 請求解析輔助函數
// =============================================================================

// decodeJSON 解碼 JSON 請求體
func decodeJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// getQueryParam 取得查詢參數
func getQueryParam(r *http.Request, key, defaultValue string) string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getQueryParamInt 取得整數查詢參數
func getQueryParamInt(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}

	var result int
	if _, err := parseIntValue(value, &result); err != nil {
		return defaultValue
	}
	return result
}

// getQueryParamBool 取得布林查詢參數
func getQueryParamBool(r *http.Request, key string, defaultValue bool) bool {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}

	switch strings.ToLower(value) {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return defaultValue
	}
}

// parseIntValue 解析整數值
func parseIntValue(s string, result *int) (bool, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return false, nil
		}
		n = n*10 + int(c-'0')
	}
	*result = n
	return true, nil
}
