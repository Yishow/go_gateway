package recordingplan

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MemoryRepository 實作內存式記錄方案儲存庫。
type MemoryRepository struct {
	mu     sync.RWMutex
	plans  map[string]RecordingPlan
	tokens map[string]SchemaPreviewToken
}

// NewMemoryRepository 建立新的 MemoryRepository。
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		plans:  make(map[string]RecordingPlan),
		tokens: make(map[string]SchemaPreviewToken),
	}
}

// CreatePlan 新增記錄方案。
func (r *MemoryRepository) CreatePlan(ctx context.Context, plan *RecordingPlan) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plans[plan.ID]; exists {
		return fmt.Errorf("recording_plan already exists: %s", plan.ID)
	}

	if plan.CreatedAt.IsZero() {
		plan.CreatedAt = time.Now().UTC()
	}
	if plan.UpdatedAt.IsZero() {
		plan.UpdatedAt = time.Now().UTC()
	}

	r.plans[plan.ID] = *plan
	return nil
}

// UpdatePlan 更新記錄方案。
func (r *MemoryRepository) UpdatePlan(ctx context.Context, plan *RecordingPlan) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plans[plan.ID]; !exists {
		return fmt.Errorf("recording_plan not found: %s", plan.ID)
	}

	plan.UpdatedAt = time.Now().UTC()
	r.plans[plan.ID] = *plan
	return nil
}

// GetPlanByID 依 ID 取得記錄方案。
func (r *MemoryRepository) GetPlanByID(ctx context.Context, id string) (*RecordingPlan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.plans[id]
	if !exists {
		return nil, fmt.Errorf("recording_plan not found: %s", id)
	}
	copyPlan := p
	return &copyPlan, nil
}

// ListPlansByWorkspace 依 Workspace 取得記錄方案列表。
func (r *MemoryRepository) ListPlansByWorkspace(ctx context.Context, workspaceID string) ([]RecordingPlan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []RecordingPlan
	for _, p := range r.plans {
		if p.WorkspaceID == workspaceID {
			list = append(list, p)
		}
	}
	return list, nil
}

// DeletePlan 刪除記錄方案。
func (r *MemoryRepository) DeletePlan(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plans[id]; !exists {
		return fmt.Errorf("recording_plan not found: %s", id)
	}
	delete(r.plans, id)
	return nil
}

// SavePreviewToken 儲存預覽 Token。
func (r *MemoryRepository) SavePreviewToken(ctx context.Context, token *SchemaPreviewToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if token.CreatedAt.IsZero() {
		token.CreatedAt = time.Now().UTC()
	}
	r.tokens[token.Token] = *token
	return nil
}

// GetPreviewToken 取得預覽 Token。
func (r *MemoryRepository) GetPreviewToken(ctx context.Context, token string) (*SchemaPreviewToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, exists := r.tokens[token]
	if !exists {
		return nil, fmt.Errorf("preview token not found: %s", token)
	}
	copyToken := t
	return &copyToken, nil
}

// DeletePreviewToken 刪除預覽 Token。
func (r *MemoryRepository) DeletePreviewToken(ctx context.Context, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tokens, token)
	return nil
}
