package hsllogic

import (
	"io"
	"testing"
	"time"
)

/**
 * BenchmarkPacketLogger_LogWithPool 測試使用物件池的 PacketLogger 效能
 * @param b 基準測試實例
 */
func BenchmarkPacketLogger_LogWithPool(b *testing.B) {
	// 建立不寫入檔案的日誌器（僅測試物件分配）
	logger := NewPacketLogger(PacketLoggerOptions{
		Enabled:       true,
		LogToFile:     false,
		LogToConsole:  false,
		ConsoleWriter: io.Discard,
	})
	defer logger.Close()

	testData := []byte{0x01, 0x03, 0x00, 0x64, 0x00, 0x01, 0x84, 0x0A}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.LogSend("modbus_tcp", "device_001", testData, "讀取保持暫存器")
	}
}

/**
 * BenchmarkPacketLogger_LogTransaction 測試完整交易記錄的效能
 * @param b 基準測試實例
 */
func BenchmarkPacketLogger_LogTransaction(b *testing.B) {
	logger := NewPacketLogger(PacketLoggerOptions{
		Enabled:       true,
		LogToFile:     false,
		LogToConsole:  false,
		ConsoleWriter: io.Discard,
	})
	defer logger.Close()

	sendData := []byte{0x01, 0x03, 0x00, 0x64, 0x00, 0x01, 0x84, 0x0A}
	recvData := []byte{0x01, 0x03, 0x02, 0x00, 0x2A, 0xB8, 0x44}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.LogTransaction(
			"modbus_tcp",
			"device_001",
			sendData,
			recvData,
			"讀取保持暫存器",
			"回應資料",
			5*time.Millisecond,
		)
	}
}

/**
 * BenchmarkPacketLogger_Concurrent 測試並發情境下的效能
 * @param b 基準測試實例
 */
func BenchmarkPacketLogger_Concurrent(b *testing.B) {
	logger := NewPacketLogger(PacketLoggerOptions{
		Enabled:       true,
		LogToFile:     false,
		LogToConsole:  false,
		ConsoleWriter: io.Discard,
		BufferSize:    10000,
	})
	defer logger.Close()

	testData := []byte{0x01, 0x03, 0x00, 0x64, 0x00, 0x01, 0x84, 0x0A}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.LogSend("modbus_tcp", "device_001", testData, "讀取保持暫存器")
		}
	})
}

/**
 * BenchmarkFormatHexBytes 測試十六進位格式化的效能
 * @param b 基準測試實例
 */
func BenchmarkFormatHexBytes(b *testing.B) {
	testData := []byte{0x01, 0x03, 0x00, 0x64, 0x00, 0x01, 0x84, 0x0A}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = formatHexBytes(testData)
	}
}

/**
 * BenchmarkAcquireReleasePacketLog 測試物件池的 acquire/release 效能
 * @param b 基準測試實例
 */
func BenchmarkAcquireReleasePacketLog(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log := acquirePacketLog()
		log.Timestamp = time.Now()
		log.Direction = PacketSend
		log.Protocol = "modbus_tcp"
		log.DeviceID = "device_001"
		log.RawData = []byte{0x01, 0x03}
		log.HexString = "0103"
		releasePacketLog(log)
	}
}
