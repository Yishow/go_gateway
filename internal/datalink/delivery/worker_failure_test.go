package delivery

import (
	"context"
	"errors"
	"testing"
	"time"
)

type failingReceiptLedger struct {
	ReceiptLedger
	readErr, saveErr error
}

func (r failingReceiptLedger) HasReceipt(destinationID, recordID string, revision int64) (bool, error) {
	if r.readErr != nil {
		return false, r.readErr
	}
	return r.ReceiptLedger.HasReceipt(destinationID, recordID, revision)
}

func (r failingReceiptLedger) SaveReceipt(receipt *Receipt) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	return r.ReceiptLedger.SaveReceipt(receipt)
}

type deliverySenderFunc func(context.Context, *OutboxItem) error

func (f deliverySenderFunc) Send(ctx context.Context, item *OutboxItem) error {
	return f(ctx, item)
}

func pendingFailureItem(t *testing.T, outbox Outbox) {
	t.Helper()
	if err := outbox.Enqueue(&OutboxItem{
		ID:                  "out-failure",
		DestinationID:       "dest-failure",
		RecordID:            "record-failure",
		CalculationRevision: 1,
		Table:               "gw_record_samples",
		Payload:             []byte(`{"value":1}`),
		ObservedAt:          time.Now(),
	}); err != nil {
		t.Fatal(err)
	}
}

func TestDeliveryWorkerKeepsReceiptFailuresUnconfirmed(t *testing.T) {
	for _, failRead := range []bool{true, false} {
		name := "save"
		if failRead {
			name = "read"
		}
		t.Run(name, func(t *testing.T) {
			outbox := NewMemoryOutbox()
			pendingFailureItem(t, outbox)
			ledger := failingReceiptLedger{ReceiptLedger: NewMemoryReceiptLedger()}
			if failRead {
				ledger.readErr = errors.New("receipt read unavailable")
			} else {
				ledger.saveErr = errors.New("receipt write unavailable")
			}
			sender := &mockSender{}
			worker := NewDeliveryWorker("dest-failure", outbox, ledger, sender, WorkerConfig{})
			if err := worker.flushBatch(t.Context()); err == nil {
				t.Fatal("receipt failure must be observable")
			}
			metrics, err := outbox.GetMetrics("dest-failure")
			expectedPending, expectedFailed := int64(1), int64(0)
			if !failRead {
				expectedPending, expectedFailed = 0, 1
			}
			if err != nil || metrics.DeliveredCount != 0 || metrics.PendingCount != expectedPending || metrics.FailedCount != expectedFailed {
				t.Fatalf("unexpected receipt failure metrics: metrics=%+v err=%v", metrics, err)
			}
			if failRead && len(sender.sentItems) != 0 {
				t.Fatal("receipt lookup failure must not bypass duplicate detection")
			}
		})
	}
}

type failingDeliveryOutbox struct {
	Outbox
	markErr error
}

func (o *failingDeliveryOutbox) MarkDelivered(itemID string, deliveredAt time.Time) error {
	if o.markErr != nil {
		return o.markErr
	}
	return o.Outbox.MarkDelivered(itemID, deliveredAt)
}

func (o *failingDeliveryOutbox) MarkFailed(string, string, int) error {
	return o.markErr
}

func TestDeliveryWorkerRetainsReceiptWhenMarkDeliveredFails(t *testing.T) {
	storeErr := errors.New("outbox storage unavailable")
	outbox := &failingDeliveryOutbox{Outbox: NewMemoryOutbox(), markErr: storeErr}
	pendingFailureItem(t, outbox)
	sender := &mockSender{}
	worker := NewDeliveryWorker("dest-failure", outbox, NewMemoryReceiptLedger(), sender, WorkerConfig{})
	if err := worker.flushBatch(t.Context()); !errors.Is(err, storeErr) {
		t.Fatalf("outbox failure must be observable: %v", err)
	}
	outbox.markErr = nil
	if err := worker.flushBatch(t.Context()); err != nil {
		t.Fatal(err)
	}
	if len(sender.sentItems) != 1 {
		t.Fatalf("persisted receipt must prevent resend, sent %d times", len(sender.sentItems))
	}
}

func TestDeliveryWorkerReturnsSendAndOutboxErrors(t *testing.T) {
	storeErr := errors.New("outbox storage unavailable")
	sendErr := errors.New("destination unavailable")
	outbox := &failingDeliveryOutbox{Outbox: NewMemoryOutbox(), markErr: storeErr}
	pendingFailureItem(t, outbox)
	sender := deliverySenderFunc(func(context.Context, *OutboxItem) error { return sendErr })
	worker := NewDeliveryWorker("dest-failure", outbox, NewMemoryReceiptLedger(), sender, WorkerConfig{})
	if err := worker.flushBatch(t.Context()); !errors.Is(err, sendErr) || !errors.Is(err, storeErr) {
		t.Fatalf("both failures must be observable: %v", err)
	}
}

func TestDeliveryWorkerShutdownFlushRetainsContextValues(t *testing.T) {
	type contextKey struct{}
	ctx, cancel := context.WithCancel(context.WithValue(t.Context(), contextKey{}, "operation-identity"))
	cancel()
	outbox := NewMemoryOutbox()
	pendingFailureItem(t, outbox)
	called := false
	sender := deliverySenderFunc(func(ctx context.Context, _ *OutboxItem) error {
		called = true
		if ctx.Err() != nil || ctx.Value(contextKey{}) != "operation-identity" {
			t.Errorf("shutdown flush must detach cancellation while retaining values: err=%v value=%v", ctx.Err(), ctx.Value(contextKey{}))
		}
		return nil
	})
	worker := NewDeliveryWorker("dest-failure", outbox, NewMemoryReceiptLedger(), sender, WorkerConfig{})
	worker.Start(ctx)
	if !called {
		t.Fatal("shutdown did not flush pending item")
	}
}
