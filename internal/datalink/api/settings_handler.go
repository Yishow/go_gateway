package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// SettingsHandler 系統設定 API 處理器
// =============================================================================

// SettingsHandler 系統設定 API 處理器
type SettingsHandler struct {
	mu       sync.RWMutex
	settings map[string]*SettingItem
}

// SettingItem 設定項目
type SettingItem struct {
	Key         string      `json:"key"`
	Value       interface{} `json:"value"`
	Description string      `json:"description"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// NewSettingsHandler 建立新的設定處理器
func NewSettingsHandler() *SettingsHandler {
	h := &SettingsHandler{
		settings: make(map[string]*SettingItem),
	}

	// 初始化預設設定
	h.initDefaults()

	return h
}

// initDefaults 初始化預設設定
func (h *SettingsHandler) initDefaults() {
	now := time.Now()

	h.settings[schema.SettingWritePrecision] = &SettingItem{
		Key:         schema.SettingWritePrecision,
		Value:       string(schema.TimePrecisionSecond),
		Description: "時序資料寫入精度 (second 或 millisecond)",
		UpdatedAt:   now,
	}

	h.settings[schema.SettingPartitionInterval] = &SettingItem{
		Key:         schema.SettingPartitionInterval,
		Value:       string(schema.PartitionMonthly),
		Description: "Postgres 時序表分區間隔 (daily, weekly, monthly)",
		UpdatedAt:   now,
	}

	h.settings[schema.SettingBatchSize] = &SettingItem{
		Key:         schema.SettingBatchSize,
		Value:       1000,
		Description: "批次寫入大小",
		UpdatedAt:   now,
	}

	h.settings[schema.SettingFlushInterval] = &SettingItem{
		Key:         schema.SettingFlushInterval,
		Value:       5000,
		Description: "批次刷新間隔 (毫秒)",
		UpdatedAt:   now,
	}

	h.settings[schema.SettingDefaultRetryCount] = &SettingItem{
		Key:         schema.SettingDefaultRetryCount,
		Value:       3,
		Description: "預設重試次數",
		UpdatedAt:   now,
	}

	h.settings[schema.SettingDefaultRetryDelay] = &SettingItem{
		Key:         schema.SettingDefaultRetryDelay,
		Value:       1000,
		Description: "預設重試延遲 (毫秒)",
		UpdatedAt:   now,
	}
}

// =============================================================================
// API 操作
// =============================================================================

// List 列出所有設定
// GET /settings
func (h *SettingsHandler) List(w http.ResponseWriter, _ *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make([]*SettingItem, 0, len(h.settings))
	for _, item := range h.settings {
		result = append(result, item)
	}

	writeJSON(w, http.StatusOK, result)
}

// Get 取得單一設定
// GET /settings/{key}
func (h *SettingsHandler) Get(w http.ResponseWriter, _ *http.Request, key string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	item, exists := h.settings[key]
	if !exists {
		writeError(w, http.StatusNotFound, "設定不存在: "+key)
		return
	}

	writeJSON(w, http.StatusOK, item)
}

// UpdateSettingRequest 更新設定請求
type UpdateSettingRequest struct {
	Value interface{} `json:"value"`
}

// Update 更新設定
// PUT/PATCH /settings/{key}
func (h *SettingsHandler) Update(w http.ResponseWriter, r *http.Request, key string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	item, exists := h.settings[key]
	if !exists {
		writeError(w, http.StatusNotFound, "設定不存在: "+key)
		return
	}

	var req UpdateSettingRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式: "+err.Error())
		return
	}

	// 驗證設定值
	if err := h.validateSetting(key, req.Value); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	item.Value = req.Value
	item.UpdatedAt = time.Now()

	writeJSON(w, http.StatusOK, item)
}

// =============================================================================
// 設定驗證
// =============================================================================

// validateSetting 驗證設定值
func (h *SettingsHandler) validateSetting(key string, value interface{}) error {
	switch key {
	case schema.SettingWritePrecision:
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("write_precision 必須為字串")
		}
		if str != string(schema.TimePrecisionSecond) && str != string(schema.TimePrecisionMillisecond) {
			return fmt.Errorf("write_precision 必須為 'second' 或 'millisecond'")
		}

	case schema.SettingPartitionInterval:
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("partition_interval 必須為字串")
		}
		if str != string(schema.PartitionDaily) &&
			str != string(schema.PartitionWeekly) &&
			str != string(schema.PartitionMonthly) {
			return fmt.Errorf("partition_interval 必須為 'daily', 'weekly' 或 'monthly'")
		}

	case schema.SettingBatchSize:
		n, ok := toInt(value)
		if !ok || n < 1 || n > 100000 {
			return fmt.Errorf("batch_size 必須為 1-100000 之間的整數")
		}

	case schema.SettingFlushInterval:
		n, ok := toInt(value)
		if !ok || n < 100 || n > 60000 {
			return fmt.Errorf("flush_interval 必須為 100-60000 之間的整數 (毫秒)")
		}

	case schema.SettingDefaultRetryCount:
		n, ok := toInt(value)
		if !ok || n < 0 || n > 10 {
			return fmt.Errorf("default_retry_count 必須為 0-10 之間的整數")
		}

	case schema.SettingDefaultRetryDelay:
		n, ok := toInt(value)
		if !ok || n < 100 || n > 30000 {
			return fmt.Errorf("default_retry_delay 必須為 100-30000 之間的整數 (毫秒)")
		}
	}

	return nil
}

// toInt 轉換為整數
func toInt(v interface{}) (int, bool) {
	switch val := v.(type) {
	case int:
		return val, true
	case int64:
		return int(val), true
	case float64:
		return int(val), true
	case float32:
		return int(val), true
	default:
		return 0, false
	}
}

// =============================================================================
// 設定存取器 (供其他模組使用)
// =============================================================================

// GetWritePrecision 取得寫入精度設定
func (h *SettingsHandler) GetWritePrecision() schema.TimePrecision {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if item, exists := h.settings[schema.SettingWritePrecision]; exists {
		if str, ok := item.Value.(string); ok {
			return schema.TimePrecision(str)
		}
	}
	return schema.TimePrecisionSecond
}

// GetPartitionInterval 取得分區間隔設定
func (h *SettingsHandler) GetPartitionInterval() schema.PartitionInterval {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if item, exists := h.settings[schema.SettingPartitionInterval]; exists {
		if str, ok := item.Value.(string); ok {
			return schema.PartitionInterval(str)
		}
	}
	return schema.PartitionMonthly
}

// GetBatchSize 取得批次大小設定
func (h *SettingsHandler) GetBatchSize() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if item, exists := h.settings[schema.SettingBatchSize]; exists {
		if n, ok := toInt(item.Value); ok {
			return n
		}
	}
	return 1000
}

// GetFlushInterval 取得刷新間隔設定 (毫秒)
func (h *SettingsHandler) GetFlushInterval() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if item, exists := h.settings[schema.SettingFlushInterval]; exists {
		if n, ok := toInt(item.Value); ok {
			return n
		}
	}
	return 5000
}

// GetDefaultRetryCount 取得預設重試次數
func (h *SettingsHandler) GetDefaultRetryCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if item, exists := h.settings[schema.SettingDefaultRetryCount]; exists {
		if n, ok := toInt(item.Value); ok {
			return n
		}
	}
	return 3
}

// GetDefaultRetryDelay 取得預設重試延遲 (毫秒)
func (h *SettingsHandler) GetDefaultRetryDelay() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if item, exists := h.settings[schema.SettingDefaultRetryDelay]; exists {
		if n, ok := toInt(item.Value); ok {
			return n
		}
	}
	return 1000
}
