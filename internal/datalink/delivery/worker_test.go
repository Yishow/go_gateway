package delivery

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

type mockSender struct {
	shouldFail int32
	sentItems  []*OutboxItem
}

func (s *mockSender) Send(ctx context.Context, item *OutboxItem) error {
	if atomic.LoadInt32(&s.shouldFail) == 1 {
		return fmt.Errorf("connection refused")
	}
	s.sentItems = append(s.sentItems, item)
	return nil
}

func TestDeliveryWorker_SuccessDeliveryAndReceipt(t *testing.T) {
	outbox := NewMemoryOutbox()
	receiptLedger := NewMemoryReceiptLedger()
	sender := &mockSender{}

	worker := NewDeliveryWorker("dest-1", outbox, receiptLedger, sender, WorkerConfig{
		BatchSize:     5,
		FlushInterval: 10 * time.Millisecond,
		MaxRetries:    3,
	})

	now := time.Now().UTC()
	_ = outbox.Enqueue(&OutboxItem{
		ID:                  "out-1",
		DestinationID:       "dest-1",
		RecordID:            "rec-1",
		CalculationRevision: 1,
		Status:              StatusPending,
		ObservedAt:          now,
	})

	ctx, cancel := context.WithCancel(context.Background())
	go worker.Start(ctx)

	time.Sleep(50 * time.Millisecond)
	cancel()
	worker.Stop()

	// 檢查已成功送達
	metrics, _ := outbox.GetMetrics("dest-1")
	if metrics.DeliveredCount != 1 || metrics.PendingCount != 0 {
		t.Errorf("expected 1 delivered item, got metrics: %+v", metrics)
	}

	// 檢查 Receipt 是否存在
	hasReceipt, _ := receiptLedger.HasReceipt("dest-1", "rec-1", 1)
	if !hasReceipt {
		t.Errorf("expected receipt to be recorded")
	}
}

func TestDeliveryWorker_FailureAndRetry(t *testing.T) {
	outbox := NewMemoryOutbox()
	receiptLedger := NewMemoryReceiptLedger()
	sender := &mockSender{shouldFail: 1} // 模擬連線失敗

	worker := NewDeliveryWorker("dest-fail", outbox, receiptLedger, sender, WorkerConfig{
		BatchSize:     5,
		FlushInterval: 10 * time.Millisecond,
		MaxRetries:    2,
	})

	now := time.Now().UTC()
	_ = outbox.Enqueue(&OutboxItem{
		ID:                  "out-fail-1",
		DestinationID:       "dest-fail",
		RecordID:            "rec-fail-1",
		CalculationRevision: 1,
		Status:              StatusPending,
		ObservedAt:          now,
	})

	ctx, cancel := context.WithCancel(context.Background())
	go worker.Start(ctx)

	time.Sleep(50 * time.Millisecond)
	cancel()
	worker.Stop()

	metrics, _ := outbox.GetMetrics("dest-fail")
	if metrics.DeliveredCount != 0 {
		t.Errorf("expected 0 delivered items on failure")
	}
}
