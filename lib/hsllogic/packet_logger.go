package hsllogic

import (
	"os"
	"sync"
	"time"
)

// =============================================================================
// PacketLogger 報文日誌器
// =============================================================================

// =============================================================================
// PacketLogger 報文日誌器
// =============================================================================

// PacketLogger 報文日誌器，用於記錄通訊報文軌跡
type PacketLogger struct {
	options       PacketLoggerOptions
	mu            sync.RWMutex
	buffer        chan *PacketLog
	currentFile   *os.File
	currentSize   int64
	stopChan      chan struct{}
	wg            sync.WaitGroup
	fileIndex     int
	lastFlushTime time.Time
}

// NewPacketLogger 建立新的報文日誌器
func NewPacketLogger(options PacketLoggerOptions) *PacketLogger {
	if options.BufferSize <= 0 {
		options.BufferSize = DefaultPacketLoggerOptions.BufferSize
	}
	if options.FlushInterval <= 0 {
		options.FlushInterval = DefaultPacketLoggerOptions.FlushInterval
	}
	if options.MaxFileSize <= 0 {
		options.MaxFileSize = DefaultPacketLoggerOptions.MaxFileSize
	}
	if options.MaxFiles <= 0 {
		options.MaxFiles = DefaultPacketLoggerOptions.MaxFiles
	}

	pl := &PacketLogger{
		options:       options,
		buffer:        make(chan *PacketLog, options.BufferSize),
		stopChan:      make(chan struct{}),
		lastFlushTime: time.Now(),
	}

	if options.Enabled && options.LogToFile {
		pl.wg.Add(1)
		go pl.writeLoop()
	}

	return pl
}

// =============================================================================
// 公開方法
// =============================================================================

// LogSend 記錄發送報文
func (pl *PacketLogger) LogSend(protocol, deviceID string, data []byte, desc string) {
	pl.log(PacketSend, protocol, deviceID, data, desc, 0)
}

// LogRecv 記錄接收報文
func (pl *PacketLogger) LogRecv(protocol, deviceID string, data []byte, desc string, duration time.Duration) {
	pl.log(PacketRecv, protocol, deviceID, data, desc, duration)
}

// LogTransaction 記錄完整的通訊交易 (發送 + 接收)
func (pl *PacketLogger) LogTransaction(
	protocol, deviceID string,
	sendData, recvData []byte,
	sendDesc, recvDesc string,
	duration time.Duration,
) {
	pl.LogSend(protocol, deviceID, sendData, sendDesc)
	pl.LogRecv(protocol, deviceID, recvData, recvDesc, duration)
}

// Close 關閉日誌器
func (pl *PacketLogger) Close() error {
	close(pl.stopChan)
	pl.wg.Wait()

	pl.mu.Lock()
	defer pl.mu.Unlock()

	if pl.currentFile != nil {
		return pl.currentFile.Close()
	}
	return nil
}

// Flush 手動刷新緩衝區
func (pl *PacketLogger) Flush() {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	if pl.currentFile != nil {
		if err := pl.currentFile.Sync(); err != nil {
			return
		}
	}
}
