package history

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Repository 定義歷史資料讀寫介面。
type Repository interface {
	SavePoints(workspaceID, planID string, points []HistoryPoint) error
	QueryPoints(ctx context.Context, query HistoryQuery) ([]HistoryPoint, error)
}

// MemoryHistoryRepository 記憶體歷史資料庫實作。
type MemoryHistoryRepository struct {
	mu     sync.RWMutex
	points map[string][]HistoryPoint // key: workspaceID + ":" + planID
}

// NewMemoryHistoryRepository 建立記憶體儲存庫。
func NewMemoryHistoryRepository() *MemoryHistoryRepository {
	return &MemoryHistoryRepository{
		points: make(map[string][]HistoryPoint),
	}
}

func (m *MemoryHistoryRepository) SavePoints(workspaceID, planID string, points []HistoryPoint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := workspaceID + ":" + planID
	m.points[key] = append(m.points[key], points...)
	return nil
}

func (m *MemoryHistoryRepository) QueryPoints(ctx context.Context, query HistoryQuery) ([]HistoryPoint, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := query.WorkspaceID + ":" + query.PlanID
	all := m.points[key]

	measFilter := make(map[string]bool)
	for _, id := range query.MeasurementIDs {
		measFilter[id] = true
	}

	var filtered []HistoryPoint
	for _, p := range all {
		if len(measFilter) > 0 && !measFilter[p.MeasurementID] {
			continue
		}
		if !query.StartTime.IsZero() && p.ObservedAt.Before(query.StartTime) {
			continue
		}
		if !query.EndTime.IsZero() && p.ObservedAt.After(query.EndTime) {
			continue
		}

		filtered = append(filtered, p)
		if query.Limit > 0 && len(filtered) >= query.Limit {
			break
		}
	}

	return filtered, nil
}

// Service 歷史資料與曲線查詢服務。
type Service struct {
	repo Repository
}

// NewService 建立歷史查詢服務。
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// QueryHistory 執行帶有 Workspace 與邊界驗證的歷史查詢。
func (s *Service) QueryHistory(ctx context.Context, query HistoryQuery) (*HistoryReport, error) {
	if query.WorkspaceID == "" {
		return nil, fmt.Errorf("workspace_id cannot be empty")
	}
	if query.PlanID == "" {
		return nil, fmt.Errorf("plan_id cannot be empty")
	}

	points, err := s.repo.QueryPoints(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query history points: %w", err)
	}

	return &HistoryReport{
		WorkspaceID: query.WorkspaceID,
		PlanID:      query.PlanID,
		Points:      points,
		TotalCount:  len(points),
		GeneratedAt: time.Now().UTC(),
	}, nil
}
