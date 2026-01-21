package hsllogic

import (
	"sync"
	"testing"
)

/**
 * TestAddressCache_Get 測試位址快取的基本功能
 * @param t 測試實例
 */
func TestAddressCache_Get(t *testing.T) {
	cache := NewAddressCache()

	// 第一次取得（快取未命中）
	addr1, err := cache.Get(ProtocolModbus, "40001")
	if err != nil {
		t.Fatalf("Get() 錯誤: %v", err)
	}
	if addr1 == nil {
		t.Fatal("Get() 應該返回非 nil 位址")
	}

	// 驗證解析結果
	if addr1.Protocol != ProtocolModbus {
		t.Errorf("Protocol = %s, 期望 %s", addr1.Protocol, ProtocolModbus)
	}
	// Modbus 40001 解析後 Offset 為 1 (40001 - 40000)
	if addr1.Offset != 1 {
		t.Errorf("Offset = %d, 期望 1", addr1.Offset)
	}

	// 第二次取得相同位址（快取命中）
	addr2, err := cache.Get(ProtocolModbus, "40001")
	if err != nil {
		t.Fatalf("Get() 錯誤: %v", err)
	}

	// 驗證返回的是相同的物件（快取命中）
	if addr1 != addr2 {
		t.Error("快取命中應該返回相同的物件指標")
	}

	// 驗證統計資訊
	hits, misses, hitRate := cache.Stats()
	if hits != 1 {
		t.Errorf("hits = %d, 期望 1", hits)
	}
	if misses != 1 {
		t.Errorf("misses = %d, 期望 1", misses)
	}
	expectedHitRate := 0.5
	if hitRate != expectedHitRate {
		t.Errorf("hitRate = %f, 期望 %f", hitRate, expectedHitRate)
	}
}

/**
 * TestAddressCache_DifferentProtocols 測試不同協議的位址分別快取
 * @param t 測試實例
 */
func TestAddressCache_DifferentProtocols(t *testing.T) {
	cache := NewAddressCache()

	// Modbus 協議的 D100
	addr1, err1 := cache.Get(ProtocolModbus, "40101")
	if err1 != nil {
		t.Fatalf("Get(Modbus) 錯誤: %v", err1)
	}

	// Mitsubishi 協議的 D100
	addr2, err2 := cache.Get(ProtocolMitsubishi, "D100")
	if err2 != nil {
		t.Fatalf("Get(Mitsubishi) 錯誤: %v", err2)
	}

	// 應該是不同的物件
	if addr1 == addr2 {
		t.Error("不同協議的位址應該分別快取")
	}

	// 驗證快取大小
	size := cache.Size()
	if size != 2 {
		t.Errorf("Size() = %d, 期望 2", size)
	}
}

/**
 * TestAddressCache_Clear 測試清空快取
 * @param t 測試實例
 */
func TestAddressCache_Clear(t *testing.T) {
	cache := NewAddressCache()

	// 新增一些快取項目
	cache.Get(ProtocolModbus, "40001")
	cache.Get(ProtocolMitsubishi, "D100")

	// 驗證快取不為空
	if cache.Size() == 0 {
		t.Error("快取應該包含項目")
	}

	// 清空快取
	cache.Clear()

	// 驗證快取已清空
	if cache.Size() != 0 {
		t.Errorf("Clear() 後 Size() = %d, 期望 0", cache.Size())
	}

	// 驗證統計已重置
	hits, misses, _ := cache.Stats()
	if hits != 0 || misses != 0 {
		t.Errorf("Clear() 後統計未重置: hits=%d, misses=%d", hits, misses)
	}
}

/**
 * TestAddressCache_InvalidAddress 測試無效位址的處理
 * @param t 測試實例
 */
func TestAddressCache_InvalidAddress(t *testing.T) {
	cache := NewAddressCache()

	// 嘗試解析無效位址
	_, err := cache.Get(ProtocolModbus, "INVALID")
	if err == nil {
		t.Error("Get() 應該對無效位址返回錯誤")
	}

	// 驗證統計（應該記錄為 miss）
	hits, misses, _ := cache.Stats()
	if hits != 0 {
		t.Errorf("無效位址不應該產生 hit: hits=%d", hits)
	}
	if misses != 1 {
		t.Errorf("無效位址應該記錄 miss: misses=%d", misses)
	}
}

/**
 * TestAddressCache_Concurrent 測試並發存取
 * @param t 測試實例
 */
func TestAddressCache_Concurrent(t *testing.T) {
	cache := NewAddressCache()
	const goroutines = 100
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	// 啟動多個 goroutine 並發存取相同的位址
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_, err := cache.Get(ProtocolModbus, "40001")
				if err != nil {
					t.Errorf("並發 Get() 錯誤: %v", err)
				}
			}
		}()
	}

	wg.Wait()

	// 驗證快取只包含一個項目
	if cache.Size() != 1 {
		t.Errorf("並發情況下 Size() = %d, 期望 1", cache.Size())
	}

	// 驗證統計資訊
	hits, misses, _ := cache.Stats()
	total := hits + misses
	expectedTotal := uint64(goroutines * iterations)
	if total != expectedTotal {
		t.Errorf("總存取次數 = %d, 期望 %d", total, expectedTotal)
	}

	// 應該有很高的命中率（只有第一次是 miss）
	if hits < expectedTotal-10 {
		t.Errorf("命中次數過低: hits=%d, total=%d", hits, total)
	}
}

/**
 * TestGetGlobalAddressCache 測試全域快取實例
 * @param t 測試實例
 */
func TestGetGlobalAddressCache(t *testing.T) {
	cache1 := GetGlobalAddressCache()
	cache2 := GetGlobalAddressCache()

	// 應該返回相同的實例
	if cache1 != cache2 {
		t.Error("GetGlobalAddressCache() 應該返回相同的實例")
	}

	// 測試基本功能
	_, err := cache1.Get(ProtocolModbus, "40001")
	if err != nil {
		t.Fatalf("全域快取 Get() 錯誤: %v", err)
	}
}
