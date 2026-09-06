package device

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 記憶體 Repository (用於測試和原型)
// =============================================================================

// MemoryRepository 記憶體內設備儲存庫
// 主要用於單元測試和原型開發，生產環境應使用資料庫實作
type MemoryRepository struct {
	mu      sync.RWMutex
	devices map[string]*schema.Device
}

// NewMemoryRepository 建立新的記憶體儲存庫
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		devices: make(map[string]*schema.Device),
	}
}

// Create 建立新設備
func (r *MemoryRepository) Create(ctx context.Context, device *schema.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[device.ID]; exists {
		return fmt.Errorf("設備 ID 已存在: %s", device.ID)
	}

	// 複製設備以避免外部修改
	deviceCopy := *device
	r.devices[device.ID] = &deviceCopy

	return nil
}

// Update 更新設備
func (r *MemoryRepository) Update(ctx context.Context, device *schema.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[device.ID]; !exists {
		return fmt.Errorf("設備不存在: %s", device.ID)
	}

	deviceCopy := *device
	r.devices[device.ID] = &deviceCopy

	return nil
}

// Delete 刪除設備
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[id]; !exists {
		return fmt.Errorf("設備不存在: %s", id)
	}

	delete(r.devices, id)
	return nil
}

// GetByID 根據 ID 取得設備
func (r *MemoryRepository) GetByID(ctx context.Context, id string) (*schema.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	device, exists := r.devices[id]
	if !exists {
		return nil, fmt.Errorf("設備不存在: %s", id)
	}

	// 返回複製以避免外部修改
	deviceCopy := *device
	return &deviceCopy, nil
}

// List 列出所有設備
func (r *MemoryRepository) List(ctx context.Context, filter ListFilter) ([]*schema.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*schema.Device, 0, len(r.devices))

	for _, device := range r.devices {
		// 套用篩選條件
		if filter.Protocol != nil && device.Protocol != *filter.Protocol {
			continue
		}
		if filter.Status != nil && device.Status != *filter.Status {
			continue
		}

		deviceCopy := *device
		result = append(result, &deviceCopy)
	}

	// 套用分頁
	if filter.Offset > 0 {
		if filter.Offset >= len(result) {
			return []*schema.Device{}, nil
		}
		result = result[filter.Offset:]
	}
	if filter.Limit > 0 && len(result) > filter.Limit {
		result = result[:filter.Limit]
	}

	return result, nil
}

// UpdateTestResult 更新連線測試結果
func (r *MemoryRepository) UpdateTestResult(ctx context.Context, id string, success bool, errMsg string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	device, exists := r.devices[id]
	if !exists {
		return fmt.Errorf("設備不存在: %s", id)
	}

	now := time.Now()
	device.LastTestAt = &now
	device.LastTestSuccess = &success
	device.LastTestError = errMsg
	device.UpdatedAt = now

	return nil
}

// ClearTestResult 清除已失效的連線測試結果
func (r *MemoryRepository) ClearTestResult(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	device, exists := r.devices[id]
	if !exists {
		return fmt.Errorf("設備不存在: %s", id)
	}

	device.LastTestAt = nil
	device.LastTestSuccess = nil
	device.LastTestError = ""
	device.UpdatedAt = time.Now()
	return nil
}

// UpdateStatus 更新設備狀態
func (r *MemoryRepository) UpdateStatus(ctx context.Context, id string, status schema.DeviceStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	device, exists := r.devices[id]
	if !exists {
		return fmt.Errorf("設備不存在: %s", id)
	}

	device.Status = status
	device.UpdatedAt = time.Now()

	return nil
}

// Count 計算設備數量
func (r *MemoryRepository) Count(ctx context.Context) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.devices)
}

// Clear 清空所有設備 (用於測試)
func (r *MemoryRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.devices = make(map[string]*schema.Device)
}
