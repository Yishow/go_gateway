package delivery

import (
	"sync"
)

// QuotaMonitor 監控本機佇列與容量水位。
type QuotaMonitor struct {
	mu           sync.RWMutex
	maxBytes     int64
	usedBytes    int64
	warningRatio float64 // 預設 0.80 (80%)
	critRatio    float64 // 預設 0.95 (95%)
}

// NewQuotaMonitor 建立配額監控器。
func NewQuotaMonitor(maxBytes int64) *QuotaMonitor {
	if maxBytes <= 0 {
		maxBytes = 500 * 1024 * 1024 // 預設 500MB
	}
	return &QuotaMonitor{
		maxBytes:     maxBytes,
		warningRatio: 0.80,
		critRatio:    0.95,
	}
}

// UpdateUsage 更新目前使用空間並計算配額警戒狀態。
func (q *QuotaMonitor) UpdateUsage(usedBytes int64) QuotaStatus {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.usedBytes = usedBytes
	ratio := float64(q.usedBytes) / float64(q.maxBytes)

	isWarn := ratio >= q.warningRatio
	isCrit := ratio >= q.critRatio

	return QuotaStatus{
		MaxQueueBytes:     q.maxBytes,
		UsedQueueBytes:    q.usedBytes,
		UsageRatio:        ratio,
		IsWarning:         isWarn,
		IsCritical:        isCrit,
		IsIntakeSuspended: isCrit,
	}
}

// GetStatus 取得當前配額狀態。
func (q *QuotaMonitor) GetStatus() QuotaStatus {
	q.mu.RLock()
	defer q.mu.RUnlock()

	ratio := float64(q.usedBytes) / float64(q.maxBytes)
	return QuotaStatus{
		MaxQueueBytes:     q.maxBytes,
		UsedQueueBytes:    q.usedBytes,
		UsageRatio:        ratio,
		IsWarning:         ratio >= q.warningRatio,
		IsCritical:        ratio >= q.critRatio,
		IsIntakeSuspended: ratio >= q.critRatio,
	}
}
