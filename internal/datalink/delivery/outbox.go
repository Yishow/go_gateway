package delivery

import (
	"fmt"
	"sync"
	"time"
)

// Outbox 定義 Outbox 儲存庫的操作介面。
type Outbox interface {
	Enqueue(item *OutboxItem) error
	FetchPending(destinationID string, limit int) ([]*OutboxItem, error)
	MarkDelivered(itemID string, deliveredAt time.Time) error
	MarkFailed(itemID string, errStr string, maxRetries int) error
	GetMetrics(destinationID string) (DestinationMetrics, error)
}

// MemoryOutbox 記憶體測試用 Outbox 實作。
type MemoryOutbox struct {
	mu    sync.RWMutex
	items map[string]*OutboxItem // id -> item
	byKey map[string]*OutboxItem // destID + ":" + recordID -> item
}

// NewMemoryOutbox 建立新的記憶體 Outbox 實例。
func NewMemoryOutbox() *MemoryOutbox {
	return &MemoryOutbox{
		items: make(map[string]*OutboxItem),
		byKey: make(map[string]*OutboxItem),
	}
}

func (m *MemoryOutbox) Enqueue(item *OutboxItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if item == nil {
		return fmt.Errorf("outbox item cannot be nil")
	}

	key := item.DestinationID + ":" + item.RecordID
	if existing, exists := m.byKey[key]; exists {
		if item.CalculationRevision < existing.CalculationRevision {
			return fmt.Errorf("stale calculation revision %d, current is %d", item.CalculationRevision, existing.CalculationRevision)
		}
		// 更高或相同 revision，更新現有 item
		existing.CalculationRevision = item.CalculationRevision
		existing.Payload = item.Payload
		existing.Status = StatusPending
		existing.UpdatedAt = time.Now().UTC()
		return nil
	}

	if item.Status == "" {
		item.Status = StatusPending
	}
	item.CreatedAt = time.Now().UTC()
	item.UpdatedAt = time.Now().UTC()

	m.items[item.ID] = item
	m.byKey[key] = item
	return nil
}

func (m *MemoryOutbox) FetchPending(destinationID string, limit int) ([]*OutboxItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	now := time.Now().UTC()
	var pending []*OutboxItem

	for _, item := range m.items {
		if item.DestinationID == destinationID {
			if item.Status == StatusPending || (item.Status == StatusRetrying && !item.NextRetryAt.After(now)) {
				pending = append(pending, item)
				if limit > 0 && len(pending) >= limit {
					break
				}
			}
		}
	}
	return pending, nil
}

func (m *MemoryOutbox) MarkDelivered(itemID string, deliveredAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, exists := m.items[itemID]
	if !exists {
		return fmt.Errorf("item %s not found", itemID)
	}

	item.Status = StatusDelivered
	item.DeliveredAt = &deliveredAt
	item.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MemoryOutbox) MarkFailed(itemID string, errStr string, maxRetries int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, exists := m.items[itemID]
	if !exists {
		return fmt.Errorf("item %s not found", itemID)
	}

	item.RetryCount++
	item.LastError = errStr
	item.UpdatedAt = time.Now().UTC()

	item.Status, item.NextRetryAt = CalculateBackoff(item.RetryCount, maxRetries)
	return nil
}

func (m *MemoryOutbox) GetMetrics(destinationID string) (DestinationMetrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var metrics DestinationMetrics
	metrics.DestinationID = destinationID
	metrics.Status = StatusPending

	now := time.Now().UTC()
	var oldestTime *time.Time

	for _, item := range m.items {
		if item.DestinationID == destinationID {
			switch item.Status {
			case StatusDelivered:
				metrics.DeliveredCount++
			case StatusBlocked, StatusQuarantined:
				metrics.FailedCount++
			default:
				metrics.PendingCount++
				metrics.PendingBytes += int64(len(item.Payload))
				if oldestTime == nil || item.ObservedAt.Before(*oldestTime) {
					t := item.ObservedAt
					oldestTime = &t
				}
			}
		}
	}

	if oldestTime != nil {
		metrics.OldestPendingAgeSec = now.Sub(*oldestTime).Seconds()
	}

	return metrics, nil
}
