package delivery

import (
	"context"
	"sync"
	"time"
)

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
type DeliveryWorker struct {
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

	for {
		select {
		case <-ctx.Done():
			w.flushBatch(context.Background())
			return
		case <-w.stopCh:
			w.flushBatch(context.Background())
			return
		case <-ticker.C:
			w.flushBatch(ctx)
		}
	}
}

// Stop 停止工作器並等待當前 batch 完成。
func (w *DeliveryWorker) Stop() {
	close(w.stopCh)
	w.wg.Wait()
}

func (w *DeliveryWorker) flushBatch(ctx context.Context) {
	items, err := w.outbox.FetchPending(w.destinationID, w.config.BatchSize)
	if err != nil || len(items) == 0 {
		return
	}

	for _, item := range items {
		// 檢查是否已有 receipt（防止重複送達）
		hasReceipt, _ := w.receiptLedger.HasReceipt(w.destinationID, item.RecordID, item.CalculationRevision)
		if hasReceipt {
			_ = w.outbox.MarkDelivered(item.ID, time.Now().UTC())
			continue
		}

		err := w.sender.Send(ctx, item)
		if err != nil {
			_ = w.outbox.MarkFailed(item.ID, err.Error(), w.config.MaxRetries)
		} else {
			now := time.Now().UTC()
			_ = w.receiptLedger.SaveReceipt(&Receipt{
				DestinationID:       w.destinationID,
				RecordID:            item.RecordID,
				CalculationRevision: item.CalculationRevision,
				Table:               item.Table,
				DeliveredAt:         now,
			})
			_ = w.outbox.MarkDelivered(item.ID, now)
		}
	}
}
