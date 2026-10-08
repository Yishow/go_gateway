package storage

import (
	"context"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

// BatchWriter 批次寫入器
type BatchWriter struct {
	config      BatchWriterConfig
	underlying  Writer
	buffer      []TimeSeriesRecord
	mu          sync.Mutex
	flushTimer  *time.Timer
	stopCh      chan struct{}
	wg          sync.WaitGroup
	lifecycle   sync.Mutex
	closing     atomic.Bool
	closeDone   chan struct{}
	closeErr    error
	timerCancel context.CancelFunc
}

// NewBatchWriter 建立新的批次寫入器
func NewBatchWriter(underlying Writer, config BatchWriterConfig) *BatchWriter {
	bw := &BatchWriter{
		config:     config,
		underlying: underlying,
		buffer:     make([]TimeSeriesRecord, 0, config.BatchSize),
		stopCh:     make(chan struct{}),
	}

	// 啟動定時刷新
	bw.startFlushTimer()

	return bw
}

// Write 寫入單筆記錄
func (bw *BatchWriter) Write(ctx context.Context, record TimeSeriesRecord) error {
	if err := bw.checkOpen(); err != nil {
		return err
	}
	bw.mu.Lock()
	defer bw.mu.Unlock()
	if err := bw.checkOpen(); err != nil {
		return err
	}

	bw.buffer = append(bw.buffer, record)

	// 達到批次大小時刷新
	if len(bw.buffer) >= bw.config.BatchSize {
		return bw.flushLocked(ctx)
	}

	return nil
}

// WriteBatch 批次寫入記錄
func (bw *BatchWriter) WriteBatch(ctx context.Context, records []TimeSeriesRecord) error {
	if err := bw.checkOpen(); err != nil {
		return err
	}
	bw.mu.Lock()
	defer bw.mu.Unlock()
	if err := bw.checkOpen(); err != nil {
		return err
	}

	bw.buffer = append(bw.buffer, records...)

	// 達到批次大小時刷新
	if len(bw.buffer) >= bw.config.BatchSize {
		return bw.flushLocked(ctx)
	}

	return nil
}

// Flush 強制刷新緩衝區
func (bw *BatchWriter) Flush(ctx context.Context) error {
	if err := bw.checkOpen(); err != nil {
		return err
	}
	bw.mu.Lock()
	defer bw.mu.Unlock()
	if err := bw.checkOpen(); err != nil {
		return err
	}

	return bw.flushLocked(ctx)
}

// flushLocked 刷新緩衝區 (必須持有鎖)
func (bw *BatchWriter) flushLocked(ctx context.Context) error {
	if len(bw.buffer) == 0 {
		return nil
	}

	// 複製緩衝區
	toWrite := slices.Clone(bw.buffer)

	// 清空緩衝區
	bw.buffer = bw.buffer[:0]

	// 寫入底層
	err := bw.underlying.WriteBatch(ctx, toWrite)
	if err != nil && ctx.Err() != nil && bw.closing.Load() {
		// Close canceled this write; keep the batch for the final close flush.
		bw.buffer = append(bw.buffer, toWrite...)
	}
	return err
}

// startFlushTimer 啟動定時刷新
func (bw *BatchWriter) startFlushTimer() {
	timerCtx, cancel := context.WithCancel(context.Background())
	bw.timerCancel = cancel
	bw.flushTimer = time.NewTimer(bw.config.FlushInterval)

	bw.wg.Go(func() {
		for {
			select {
			case <-bw.stopCh:
				bw.flushTimer.Stop()
				return
			case <-bw.flushTimer.C:
				ctx, cancel := context.WithTimeout(timerCtx, bw.config.WriteTimeout)
				_ = bw.Flush(ctx)
				cancel()
				bw.flushTimer.Reset(bw.config.FlushInterval)
			}
		}
	})
}

// BufferSize 取得目前緩衝區大小
func (bw *BatchWriter) BufferSize() int {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	return len(bw.buffer)
}
