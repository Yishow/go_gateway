package hsllogic

import (
	"testing"
	"time"
)

/**
 * TestPacketLog_LazyHexFormatting 測試延遲 Hex 格式化功能
 * @param t 測試實例
 */
func TestPacketLog_LazyHexFormatting(t *testing.T) {
	testData := []byte{0x01, 0x03, 0x00, 0x64}

	// 從物件池取得 PacketLog
	log := acquirePacketLog()
	defer releasePacketLog(log)

	log.Timestamp = time.Now()
	log.Direction = PacketSend
	log.Protocol = "modbus_tcp"
	log.DeviceID = "device_001"
	log.RawData = testData
	log.Description = "測試"

	// 在調用 GetHexString 之前，HexString 應該為空且未快取
	if log.hexCached {
		t.Error("hexCached 應該為 false")
	}
	if log.HexString != "" {
		t.Error("HexString 應該為空字串")
	}

	// 第一次調用 GetHexString，應該生成並快取
	hexStr := log.GetHexString()
	expectedHex := "01030064"
	if hexStr != expectedHex {
		t.Errorf("GetHexString() = %s, 期望 %s", hexStr, expectedHex)
	}
	if !log.hexCached {
		t.Error("hexCached 應該為 true")
	}

	// 第二次調用應該返回快取的值
	hexStr2 := log.GetHexString()
	if hexStr2 != expectedHex {
		t.Errorf("GetHexString() 第二次調用 = %s, 期望 %s", hexStr2, expectedHex)
	}

	// String() 方法應該包含 HexString
	logStr := log.String()
	if logStr == "" {
		t.Error("String() 不應該返回空字串")
	}
}

/**
 * TestPacketLog_String 測試 PacketLog 字串格式化
 * @param t 測試實例
 */
func TestPacketLog_String(t *testing.T) {
	log := acquirePacketLog()
	defer releasePacketLog(log)

	log.Timestamp = time.Date(2024, 1, 15, 10, 30, 45, 123000000, time.UTC)
	log.Direction = PacketSend
	log.Protocol = "modbus_tcp"
	log.DeviceID = "device_001"
	log.RawData = []byte{0x01, 0x03}
	log.Description = "讀取保持暫存器"
	log.Duration = 5 * time.Millisecond

	str := log.String()

	// 驗證字串包含預期內容
	expectedParts := []string{
		"2024-01-15 10:30:45",
		"[SEND]",
		"[modbus_tcp]",
		"[device_001]",
		"0103",
		"讀取保持暫存器",
		"5ms",
	}

	for _, part := range expectedParts {
		found := false
		for i := 0; i <= len(str)-len(part); i++ {
			if str[i:i+len(part)] == part {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("String() 應該包含 %q，實際值: %s", part, str)
		}
	}
}

/**
 * TestAcquireReleasePacketLog 測試物件池的 acquire/release 功能
 * @param t 測試實例
 */
func TestAcquireReleasePacketLog(t *testing.T) {
	// 取得物件
	log1 := acquirePacketLog()
	if log1 == nil {
		t.Fatal("acquirePacketLog() 應該返回非 nil 物件")
	}

	// 設定一些值
	log1.Protocol = "test"
	log1.DeviceID = "device_001"
	log1.hexCached = true

	// 歸還物件
	releasePacketLog(log1)

	// 再次取得物件（可能是同一個物件）
	log2 := acquirePacketLog()
	if log2 == nil {
		t.Fatal("acquirePacketLog() 應該返回非 nil 物件")
	}

	// 驗證物件已被重置（如果是同一個物件）
	if log2 == log1 {
		if log2.Protocol != "" {
			t.Error("歸還後的物件 Protocol 應該被重置為空")
		}
		if log2.DeviceID != "" {
			t.Error("歸還後的物件 DeviceID 應該被重置為空")
		}
		if log2.hexCached {
			t.Error("歸還後的物件 hexCached 應該被重置為 false")
		}
	}

	releasePacketLog(log2)
}
