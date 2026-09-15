package history

import (
	"time"
)

// Resolution 表示歷史查詢的聚合粒度。
type Resolution string

const (
	ResolutionRaw    Resolution = "raw"
	ResolutionMinute Resolution = "1m"
	ResolutionHour   Resolution = "1h"
	ResolutionDay    Resolution = "1d"
)

// HistoryQuery 歷史資料與曲線查詢請求。
type HistoryQuery struct { //nolint:revive // Preserve the exported Go name and its existing callers during lint maintenance.
	WorkspaceID    string     `json:"workspace_id"`
	PlanID         string     `json:"plan_id"`
	MeasurementIDs []string   `json:"measurement_ids"`
	StartTime      time.Time  `json:"start_time"`
	EndTime        time.Time  `json:"end_time"`
	Resolution     Resolution `json:"resolution"`
	Limit          int        `json:"limit"`
	Cursor         string     `json:"cursor,omitempty"`
	RevisionCutoff int64      `json:"revision_cutoff,omitempty"`
}

// HistoryPoint 單一時間點的歷史與曲線資料。
type HistoryPoint struct { //nolint:revive // Preserve the exported Go name and its existing callers during lint maintenance.
	MeasurementID    string    `json:"measurement_id"`
	ObservedAt       time.Time `json:"observed_at"`
	ValueNumeric     *float64  `json:"value_numeric,omitempty"`
	ValueString      *string   `json:"value_string,omitempty"`
	TimeWeightedMean *float64  `json:"time_weighted_mean,omitempty"`
	SampledMin       *float64  `json:"sampled_min,omitempty"`
	SampledMax       *float64  `json:"sampled_max,omitempty"`
	UsageDelta       *float64  `json:"usage_delta,omitempty"`
	Quality          string    `json:"quality"`
	CoverageRatio    float64   `json:"coverage_ratio"`
	IsEstimated      bool      `json:"is_estimated"`
	IsProvisional    bool      `json:"is_provisional"`
}

// HistoryReport 歷史查詢結果回傳模型。
type HistoryReport struct { //nolint:revive // Preserve the exported Go name and its existing callers during lint maintenance.
	WorkspaceID string         `json:"workspace_id"`
	PlanID      string         `json:"plan_id"`
	Points      []HistoryPoint `json:"points"`
	TotalCount  int            `json:"total_count"`
	NextCursor  string         `json:"next_cursor,omitempty"`
	GeneratedAt time.Time      `json:"generated_at"`
}
