package history

import (
	"context"
	"testing"
	"time"
)

func floatPtr(v float64) *float64 {
	return &v
}

func TestHistoryService_QueryPoints(t *testing.T) {
	repo := NewMemoryHistoryRepository()
	service := NewService(repo)

	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)

	// 注入測試資料
	p1 := HistoryPoint{
		MeasurementID: "meas-kw",
		ObservedAt:    baseTime,
		ValueNumeric:  floatPtr(100.5),
		Quality:       "good",
		CoverageRatio: 1.0,
	}
	p2 := HistoryPoint{
		MeasurementID: "meas-kw",
		ObservedAt:    baseTime.Add(10 * time.Minute),
		ValueNumeric:  floatPtr(105.0),
		Quality:       "good",
		CoverageRatio: 1.0,
	}

	_ = repo.SavePoints("ws-1", "plan-1", []HistoryPoint{p1, p2})

	ctx := context.Background()
	report, err := service.QueryHistory(ctx, HistoryQuery{
		WorkspaceID:    "ws-1",
		PlanID:         "plan-1",
		MeasurementIDs: []string{"meas-kw"},
		StartTime:      baseTime.Add(-1 * time.Minute),
		EndTime:        baseTime.Add(20 * time.Minute),
		Resolution:     ResolutionRaw,
		Limit:          10,
	})

	if err != nil {
		t.Fatalf("query history failed: %v", err)
	}
	if len(report.Points) != 2 {
		t.Fatalf("expected 2 points, got %d", len(report.Points))
	}
	if *report.Points[0].ValueNumeric != 100.5 || *report.Points[1].ValueNumeric != 105.0 {
		t.Errorf("unexpected query points content")
	}
}

func TestHistoryService_WorkspaceIsolation(t *testing.T) {
	repo := NewMemoryHistoryRepository()
	service := NewService(repo)

	baseTime := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	p1 := HistoryPoint{
		MeasurementID: "meas-1",
		ObservedAt:    baseTime,
		ValueNumeric:  floatPtr(50.0),
	}
	_ = repo.SavePoints("ws-1", "plan-1", []HistoryPoint{p1})

	// 嘗試用 ws-2 查詢 -> 應為空
	ctx := context.Background()
	report, err := service.QueryHistory(ctx, HistoryQuery{
		WorkspaceID:    "ws-2",
		PlanID:         "plan-1",
		MeasurementIDs: []string{"meas-1"},
		StartTime:      baseTime.Add(-1 * time.Minute),
		EndTime:        baseTime.Add(10 * time.Minute),
	})

	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if len(report.Points) != 0 {
		t.Errorf("expected 0 points for isolated workspace, got %d", len(report.Points))
	}
}
