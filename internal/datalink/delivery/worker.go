package delivery

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

var errReceiptBlockFailed = errors.New("receipt failure could not be durably blocked")

// DestinationSender 定義向特定目的地派送資料的執行介面。
type DestinationSender interface {
	Send(ctx context.Context, item *OutboxItem) error
}

// WorkerConfig 交付工作器配置。
type WorkerConfig struct {
	BatchSize     int           `json:"batch_size"`
	FlushInterval time.Duration `json:"flush_interval"`
	MaxRetries    int           `json:"max_retries"`
}

// DeliveryWorker 負責單一目的地的背景可靠交付與回執管理。
type DeliveryWorker struct { //nolint:revive // Preserve the exported Go name and its existing callers during lint maintenance.
	destinationID string
	outbox        Outbox
	receiptLedger ReceiptLedger
	sender        DestinationSender
	config        WorkerConfig
	stopCh        chan struct{}
	wg            sync.WaitGroup
}

// NewDeliveryWorker 建立新的交付工作器。
func NewDeliveryWorker(
	destinationID string,
	outbox Outbox,
	receiptLedger ReceiptLedger,
	sender DestinationSender,
	config WorkerConfig,
) *DeliveryWorker {
	if config.BatchSize <= 0 {
		config.BatchSize = 100
	}
	if config.FlushInterval <= 0 {
		config.FlushInterval = 1 * time.Second
	}
	if config.MaxRetries <= 0 {
		config.MaxRetries = 5
	}

	return &DeliveryWorker{
		destinationID: destinationID,
		outbox:        outbox,
		receiptLedger: receiptLedger,
		sender:        sender,
		config:        config,
		stopCh:        make(chan struct{}),
	}
}

// Start 啟動背景派送循環。
func (w *DeliveryWorker) Start(ctx context.Context) {
	w.wg.Add(1)
	defer w.wg.Done()

	ticker := time.NewTicker(w.config.FlushInterval)
	defer ticker.Stop()
	flush := func(flushCtx context.Context) error {
		if err := w.flushBatch(flushCtx); err != nil {
			if errors.Is(err, errReceiptBlockFailed) {
				slog.ErrorContext(flushCtx, "delivery destination stopped after receipt and block storage failures; reconcile delivery status before restart", "destination_id", w.destinationID, "error", err)
			} else {
				slog.ErrorContext(flushCtx, "delivery batch failed", "destination_id", w.destinationID, "error", err)
			}
			return err
		}
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			if err := flush(context.WithoutCancel(ctx)); err != nil {
				return
			}
			return
		case <-w.stopCh:
			if err := flush(context.WithoutCancel(ctx)); err != nil {
				return
			}
			return
		case <-ticker.C:
			if err := flush(ctx); errors.Is(err, errReceiptBlockFailed) {
				return
			}
		}
	}
}

// Stop 停止工作器並等待當前 batch 完成。
func (w *DeliveryWorker) Stop() {
	close(w.stopCh)
	w.wg.Wait()
}

func (w *DeliveryWorker) flushBatch(ctx context.Context) error {
	items, err := w.outbox.FetchPending(w.destinationID, w.config.BatchSize)
	if err != nil {
		return fmt.Errorf("fetch pending deliveries: %w", err)
	}
	var failures []error
	for _, item := range items {
		if err := w.deliverItem(ctx, item); err != nil {
			wrappedErr := fmt.Errorf("deliver record %s: %w", item.RecordID, err)
			failures = append(failures, wrappedErr)
			if errors.Is(err, errReceiptBlockFailed) {
				return errors.Join(failures...)
			}
		}
	}
	return errors.Join(failures...)
}

func (w *DeliveryWorker) deliverItem(ctx context.Context, item *OutboxItem) error {
	hasReceipt, err := w.receiptLedger.HasReceipt(w.destinationID, item.RecordID, item.CalculationRevision)
	if err != nil {
		return fmt.Errorf("read receipt: %w", err)
	}
	if hasReceipt {
		return w.outbox.MarkDelivered(item.ID, time.Now().UTC())
	}
	if err := w.sender.Send(ctx, item); err != nil {
		return errors.Join(err, w.outbox.MarkFailed(item.ID, err.Error(), w.config.MaxRetries))
	}
	now := time.Now().UTC()
	if err := w.receiptLedger.SaveReceipt(&Receipt{
		DestinationID:       w.destinationID,
		RecordID:            item.RecordID,
		CalculationRevision: item.CalculationRevision,
		Table:               item.Table,
		DeliveredAt:         now,
	}); err != nil {
		receiptErr := fmt.Errorf("save receipt: %w", err)
		if markErr := w.outbox.MarkFailed(item.ID, receiptErr.Error(), 1); markErr != nil {
			return errors.Join(errReceiptBlockFailed, receiptErr, markErr)
		}
		return receiptErr
	}
	return w.outbox.MarkDelivered(item.ID, now)
}
