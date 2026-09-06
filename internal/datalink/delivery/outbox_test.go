package delivery

import (
	"testing"
	"time"
)

func TestOutbox_EnqueueAndRevisionGuard(t *testing.T) {
	outbox := NewMemoryOutbox()

	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	item1 := &OutboxItem{
		ID:                  "out-1",
		DestinationID:       "dest-pg",
		RecordID:            "rec-summary-100",
		CalculationRevision: 2,
		Table:               "gw_record_intervals",
		Payload:             []byte(`{"val": 15.0}`),
		Status:              StatusPending,
		ObservedAt:          now,
	}

	err := outbox.Enqueue(item1)
	if err != nil {
		t.Fatalf("enqueue item1 failed: %v", err)
	}

	// 嘗試用較舊的 revision 1 覆寫 -> 應被阻擋 (P07)
	itemOld := &OutboxItem{
		ID:                  "out-2",
		DestinationID:       "dest-pg",
		RecordID:            "rec-summary-100",
		CalculationRevision: 1, // 舊 revision
		Table:               "gw_record_intervals",
		Payload:             []byte(`{"val": 10.0}`),
		Status:              StatusPending,
		ObservedAt:          now,
	}

	err = outbox.Enqueue(itemOld)
	if err == nil {
		t.Errorf("expected error when enqueuing stale calculation revision")
	}

	// 用更高的 revision 3 更新 -> 應成功
	itemNew := &OutboxItem{
		ID:                  "out-3",
		DestinationID:       "dest-pg",
		RecordID:            "rec-summary-100",
		CalculationRevision: 3,
		Table:               "gw_record_intervals",
		Payload:             []byte(`{"val": 18.0}`),
		Status:              StatusPending,
		ObservedAt:          now,
	}

	err = outbox.Enqueue(itemNew)
	if err != nil {
		t.Fatalf("enqueue newer revision failed: %v", err)
	}

	pending, err := outbox.FetchPending("dest-pg", 10)
	if err != nil {
		t.Fatalf("fetch pending failed: %v", err)
	}
	if len(pending) != 1 {
		t.Fatalf("expected 1 pending item, got %d", len(pending))
	}
	if pending[0].CalculationRevision != 3 {
		t.Errorf("expected revision 3, got %d", pending[0].CalculationRevision)
	}
}

func TestOutbox_DestinationIsolation(t *testing.T) {
	// P05: 一慢一快兩目的地 -> 獨立 backlog 隔離
	outbox := NewMemoryOutbox()
	now := time.Now().UTC()

	_ = outbox.Enqueue(&OutboxItem{
		ID:            "item-fast-1",
		DestinationID: "dest-fast",
		RecordID:      "rec-1",
		Status:        StatusPending,
		ObservedAt:    now,
	})
	_ = outbox.Enqueue(&OutboxItem{
		ID:            "item-slow-1",
		DestinationID: "dest-slow",
		RecordID:      "rec-1",
		Status:        StatusPending,
		ObservedAt:    now,
	})

	fastItems, _ := outbox.FetchPending("dest-fast", 10)
	slowItems, _ := outbox.FetchPending("dest-slow", 10)

	if len(fastItems) != 1 || len(slowItems) != 1 {
		t.Fatalf("expected 1 item per destination")
	}

	// 標記 fast 為 delivered，不影響 slow
	_ = outbox.MarkDelivered("item-fast-1", now)

	metricsFast, _ := outbox.GetMetrics("dest-fast")
	metricsSlow, _ := outbox.GetMetrics("dest-slow")

	if metricsFast.PendingCount != 0 || metricsFast.DeliveredCount != 1 {
		t.Errorf("unexpected fast metrics: %+v", metricsFast)
	}
	if metricsSlow.PendingCount != 1 || metricsSlow.DeliveredCount != 0 {
		t.Errorf("unexpected slow metrics: %+v", metricsSlow)
	}
}
