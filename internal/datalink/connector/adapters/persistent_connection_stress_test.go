package adapters

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"
)

/**
 * TestModbusTCPConnector_StressTest 壓力測試：高併發讀取與模式切換
 *
 * 測試場景：
 * 1. 建立一個使用 MockTransport 的 ModbusTCPConnector
 * 2. 啟動 N 個 goroutine 進行高頻讀取 (Simulates high concurrency reading)
 * 3. 啟動 1 個 goroutine 定期切換 Persistent/Short 模式 (Simulates mode toggling)
 * 4. 驗證在高併發與狀態變更下，程式不會崩潰且能持續運作
 * 5. 驗證資料讀取的正確性 (基於 MockTransport 的行為)
 */
func TestModbusTCPConnector_StressTest(t *testing.T) {
	// 1. 初始化 MockTransport
	mockTransport := NewMockTransport()

	// 預先填充一些資料到 MockTransport
	// 地址 0-99 填入對應的數值
	for i := uint16(0); i < 100; i++ {
		mockTransport.HoldingRegs[i] = i
	}

	// 建立 Connector
	conn := &ModbusTCPConnector{
		transport:      mockTransport,
		connected:      true, // 模擬已連線
		persistentMode: true, // 初始為長連接
	}

	// 初始化 config，避免讀取時出錯
	conn.config.SlaveID = 1
	conn.config.Timeout = 1

	// 重要：我們需要注入一個使用 mockTransport 的 client
	// 因為 ModbusTCPConnector 的 Read 方法會使用 c.client
	// 而我們無法輕易透過公開 API 注入自定義 transport 的 client
	// 這裡我們直接建立 client 並注入到 conn 中
	conn.client = modbus.NewClient(mockTransport, conn.config.SlaveID)

	// 定義併發參數
	concurrency := 20           // 併發讀取數
	duration := 2 * time.Second // 測試持續時間

	var successCount int64
	var errorCount int64

	// 用於通知測試結束
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var wg sync.WaitGroup

	// 2. 啟動並發讀取 Goroutines
	t.Logf("開始壓力測試：併發數 %d, 持續時間 %v", concurrency, duration)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 隨機延遲啟動，避免同時衝擊
			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)

			for {
				select {
				case <-ctx.Done():
					return
				default:
					// 隨機讀取一個地址
					addr := uint16(rand.Intn(50))

					req := connector.ReadRequest{
						Address:  fmt.Sprintf("4%04d", addr+1), // 40001 + addr
						Function: "03",                         // Holding Register
						DataType: schema.DataTypeUint16,
						Count:    1,
					}

					// 執行讀取
					result, err := conn.Read(ctx, req)

					if err != nil {
						// 錯誤處理
						if atomic.LoadInt64(&errorCount) < 5 {
							t.Logf("讀取錯誤: %v, Mode: %v, Connected: %v", err, conn.IsPersistentMode(), conn.IsConnected())
						}
						atomic.AddInt64(&errorCount, 1)

						// 如果是長連接模式且未連線，嘗試重連 (模擬調用者行為)
						if conn.IsPersistentMode() && !conn.IsConnected() {
							_ = conn.Reconnect(ctx)
						}
					} else {
						// 驗證數值正確性
						val, ok := result.Value.(uint16)
						if !ok || val != addr {
							t.Errorf("讀取數值錯誤: 預期 %d, 實際 %v", addr, result.Value)
							atomic.AddInt64(&errorCount, 1)
						} else {
							atomic.AddInt64(&successCount, 1)
						}
					}

					// 稍微休眠，模擬真實負載間隔
					time.Sleep(time.Millisecond * 10)
				}
			}
		}(i)
	}

	// 3. 啟動模式切換 Goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(200 * time.Millisecond) // 每 200ms 切換一次
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// 切換模式
				isPersistent := conn.IsPersistentMode()
				conn.SetPersistentConnection(!isPersistent)

				// 如果切換到長連接，需要手動確保連線（模擬 Connect）
				if !isPersistent { // 從 false -> true
					// 在真實場景，使用者會呼叫 Connect 或 Reconnect
					// 這裡我們簡單模擬
					conn.setConnected(true)
				}
			}
		}
	}()

	// 等待測試結束
	wg.Wait()

	t.Logf("壓力測試結束。成功: %d, 失敗: %d", successCount, errorCount)

	// 4. 驗證結果
	if successCount == 0 {
		t.Error("測試期間沒有成功的讀取")
	}

	// 允許少量錯誤（因為模式切換時可能會有短暫的連線不穩或狀態不一致）
	total := successCount + errorCount
	errorRate := float64(errorCount) / float64(total)
	if errorRate > 0.1 { // 允許 10% 的錯誤率
		t.Errorf("錯誤率過高: %.2f%%", errorRate*100)
	}
}
