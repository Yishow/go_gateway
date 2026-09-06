package measurement

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MemoryRepository 實作內存式量測定義儲存庫。
type MemoryRepository struct {
	mu   sync.RWMutex
	defs map[string]MeasurementDefinition
}

// NewMemoryRepository 建立新的 MemoryRepository。
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		defs: make(map[string]MeasurementDefinition),
	}
}

// Create 新增量測定義。
func (r *MemoryRepository) Create(ctx context.Context, def *MeasurementDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.defs[def.ID]; exists {
		return fmt.Errorf("measurement_definition already exists: %s", def.ID)
	}

	if def.CreatedAt.IsZero() {
		def.CreatedAt = time.Now().UTC()
	}
	if def.UpdatedAt.IsZero() {
		def.UpdatedAt = time.Now().UTC()
	}

	r.defs[def.ID] = *def
	return nil
}

// Update 更新量測定義。
func (r *MemoryRepository) Update(ctx context.Context, def *MeasurementDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.defs[def.ID]; !exists {
		return fmt.Errorf("measurement_definition not found: %s", def.ID)
	}

	def.UpdatedAt = time.Now().UTC()
	r.defs[def.ID] = *def
	return nil
}

// GetByID 依 ID 取得量測定義。
func (r *MemoryRepository) GetByID(ctx context.Context, id string) (*MeasurementDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	def, exists := r.defs[id]
	if !exists {
		return nil, fmt.Errorf("measurement_definition not found: %s", id)
	}
	copyDef := def
	return &copyDef, nil
}

// ListByWorkspace 依 Workspace 取得量測定義列表。
func (r *MemoryRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]MeasurementDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []MeasurementDefinition
	for _, def := range r.defs {
		if def.WorkspaceID == workspaceID {
			list = append(list, def)
		}
	}
	return list, nil
}

// ListByDevice 依 Device 取得量測定義列表。
func (r *MemoryRepository) ListByDevice(ctx context.Context, workspaceID, deviceID string) ([]MeasurementDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []MeasurementDefinition
	for _, def := range r.defs {
		if def.WorkspaceID == workspaceID && def.DeviceID == deviceID {
			list = append(list, def)
		}
	}
	return list, nil
}

// Delete 刪除量測定義。
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.defs[id]; !exists {
		return fmt.Errorf("measurement_definition not found: %s", id)
	}
	delete(r.defs, id)
	return nil
}
