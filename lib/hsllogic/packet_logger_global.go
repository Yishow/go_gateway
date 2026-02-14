package hsllogic

import (
	"sync"
	"time"
)

// =============================================================================
// 全域日誌器
// =============================================================================

var (
	globalLogger     *PacketLogger
	globalLoggerOnce sync.Once
)

// GetGlobalPacketLogger 取得全域報文日誌器
func GetGlobalPacketLogger() *PacketLogger {
	globalLoggerOnce.Do(func() {
		globalLogger = NewPacketLogger(DefaultPacketLoggerOptions)
	})
	return globalLogger
}

// SetGlobalPacketLogger 設定全域報文日誌器
func SetGlobalPacketLogger(logger *PacketLogger) {
	globalLogger = logger
}

// LogPacketSend 使用全域日誌器記錄發送報文
func LogPacketSend(protocol, deviceID string, data []byte, desc string) {
	GetGlobalPacketLogger().LogSend(protocol, deviceID, data, desc)
}

// LogPacketRecv 使用全域日誌器記錄接收報文
func LogPacketRecv(protocol, deviceID string, data []byte, desc string, duration time.Duration) {
	GetGlobalPacketLogger().LogRecv(protocol, deviceID, data, desc, duration)
}
