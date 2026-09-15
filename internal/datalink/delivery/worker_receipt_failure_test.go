package delivery

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestDeliveryWorkerBlocksAfterReceiptFailureMemory(t *testing.T) {
	outbox := NewMemoryOutbox()
	pendingFailureItem(t, outbox)
	saveErr := errors.New("receipt write unavailable")
	ledger := failingReceiptLedger{ReceiptLedger: NewMemoryReceiptLedger(), saveErr: saveErr}
	var sent int
	sender := deliverySenderFunc(func(context.Context, *OutboxItem) error {
		sent++
		return nil
	})
	worker := NewDeliveryWorker("dest-failure", outbox, ledger, sender, WorkerConfig{})

	if err := worker.flushBatch(t.Context()); !errors.Is(err, saveErr) {
		t.Fatalf("receipt failure must be returned: %v", err)
	}
	assertBlockedMetrics(t, outbox)
	if sent != 1 {
		t.Fatalf("expected one send before receipt failure, got %d", sent)
	}
	if err := worker.flushBatch(t.Context()); err != nil {
		t.Fatalf("blocked item must not be retried: %v", err)
	}
	if sent != 1 {
		t.Fatalf("blocked item was resent, sent %d times", sent)
	}
}

func TestDeliveryWorkerBlocksAfterReceiptFailureSQLite(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	outbox := NewSQLOutbox(db)
	pendingFailureItem(t, outbox)
	saveErr := errors.New("receipt write unavailable")
	ledger := failingReceiptLedger{ReceiptLedger: NewSQLReceiptLedger(db), saveErr: saveErr}
	var sent int
	sender := deliverySenderFunc(func(context.Context, *OutboxItem) error {
		sent++
		return nil
	})
	worker := NewDeliveryWorker("dest-failure", outbox, ledger, sender, WorkerConfig{})

	if err := worker.flushBatch(t.Context()); !errors.Is(err, saveErr) {
		t.Fatalf("receipt failure must be returned: %v", err)
	}
	blockedOutbox := NewSQLOutbox(db)
	assertBlockedMetrics(t, blockedOutbox)
	if sent != 1 {
		t.Fatalf("expected one send before receipt failure, got %d", sent)
	}
	retryWorker := NewDeliveryWorker("dest-failure", blockedOutbox, ledger, sender, WorkerConfig{})
	if err := retryWorker.flushBatch(t.Context()); err != nil {
		t.Fatalf("blocked item must not be retried: %v", err)
	}
	if sent != 1 {
		t.Fatalf("blocked item was resent, sent %d times", sent)
	}
}

func TestDeliveryWorkerReturnsReceiptAndBlockErrors(t *testing.T) {
	saveErr := errors.New("receipt write unavailable")
	markErr := errors.New("outbox block unavailable")
	outbox := &failingDeliveryOutbox{Outbox: NewMemoryOutbox(), markErr: markErr}
	pendingFailureItem(t, outbox)
	ledger := failingReceiptLedger{ReceiptLedger: NewMemoryReceiptLedger(), saveErr: saveErr}
	worker := NewDeliveryWorker("dest-failure", outbox, ledger, deliverySenderFunc(func(context.Context, *OutboxItem) error {
		return nil
	}), WorkerConfig{})

	err := worker.flushBatch(t.Context())
	if !errors.Is(err, errReceiptBlockFailed) || !errors.Is(err, saveErr) || !errors.Is(err, markErr) {
		t.Fatalf("receipt and block failures must be observable: %v", err)
	}
}

func TestDeliveryWorkerStopsBatchAfterReceiptBlockFailure(t *testing.T) {
	saveErr := errors.New("receipt write unavailable")
	markErr := errors.New("outbox block unavailable")
	outbox := &orderedPendingOutbox{
		items: []*OutboxItem{
			{ID: "first", DestinationID: "dest-failure", RecordID: "first-record"},
			{ID: "second", DestinationID: "dest-failure", RecordID: "second-record"},
		},
		markErr: markErr,
	}
	ledger := failingReceiptLedger{ReceiptLedger: NewMemoryReceiptLedger(), saveErr: saveErr}
	var sent atomic.Int32
	worker := NewDeliveryWorker("dest-failure", outbox, ledger, deliverySenderFunc(func(context.Context, *OutboxItem) error {
		sent.Add(1)
		return nil
	}), WorkerConfig{})

	err := worker.flushBatch(t.Context())
	if !errors.Is(err, errReceiptBlockFailed) || !errors.Is(err, saveErr) || !errors.Is(err, markErr) {
		t.Fatalf("batch must return both storage failures: %v", err)
	}
	if got := sent.Load(); got != 1 {
		t.Fatalf("batch must stop after the first unblocked item, sent %d times", got)
	}
}

func TestDeliveryWorkerStopsAfterReceiptBlockFailure(t *testing.T) {
	saveErr := errors.New("receipt write unavailable")
	markErr := errors.New("outbox block unavailable")
	outbox := &failingDeliveryOutbox{Outbox: NewMemoryOutbox(), markErr: markErr}
	pendingFailureItem(t, outbox)
	ledger := failingReceiptLedger{ReceiptLedger: NewMemoryReceiptLedger(), saveErr: saveErr}
	var sent atomic.Int32
	worker := NewDeliveryWorker("dest-failure", outbox, ledger, deliverySenderFunc(func(context.Context, *OutboxItem) error {
		sent.Add(1)
		return nil
	}), WorkerConfig{FlushInterval: 10 * time.Millisecond})
	done := make(chan struct{})
	go func() {
		worker.Start(t.Context())
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		worker.Stop()
		t.Fatal("worker did not stop after receipt and block failures")
	}
	if got := sent.Load(); got != 1 {
		t.Fatalf("worker must not resend after block failure, sent %d times", got)
	}
}

func assertBlockedMetrics(t *testing.T, outbox Outbox) {
	t.Helper()
	metrics, err := outbox.GetMetrics("dest-failure")
	if err != nil {
		t.Fatalf("get delivery metrics: %v", err)
	}
	if metrics.PendingCount != 0 || metrics.FailedCount != 1 || metrics.DeliveredCount != 0 {
		t.Fatalf("receipt failure must block exactly one item: metrics=%+v", metrics)
	}
}

type orderedPendingOutbox struct {
	items   []*OutboxItem
	markErr error
}

func (o *orderedPendingOutbox) Enqueue(*OutboxItem) error {
	return nil
}

func (o *orderedPendingOutbox) FetchPending(string, int) ([]*OutboxItem, error) {
	return o.items, nil
}

func (o *orderedPendingOutbox) MarkDelivered(string, time.Time) error {
	return nil
}

func (o *orderedPendingOutbox) MarkFailed(string, string, int) error {
	return o.markErr
}

func (o *orderedPendingOutbox) GetMetrics(string) (DestinationMetrics, error) {
	return DestinationMetrics{}, nil
}
