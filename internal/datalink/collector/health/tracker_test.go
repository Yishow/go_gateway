// Package health 提供設備健康度追蹤功能
package health

import (
	"sync"
	"testing"
	"time"
)

// =============================================================================
// 任務 1.1: [RED] 建立滑動視窗測試
// =============================================================================

func TestTracker_SlidingWindow(t *testing.T) {
	// 建立 10 次請求的滑動視窗追蹤器
	tracker := NewTracker(10)

	// 模擬輸入 10 次結果 (8 成功 2 失敗)
	for i := 0; i < 8; i++ {
		tracker.RecordSuccess(10 * time.Millisecond)
	}
	for i := 0; i < 2; i++ {
		tracker.RecordFailure(nil)
	}

	// 斷言錯誤率為 20%
	errorRate := tracker.ErrorRate()
	if errorRate < 0.19 || errorRate > 0.21 {
		t.Errorf("預期錯誤率約為 0.20，實際為 %f", errorRate)
	}
}

func TestTracker_SlidingWindow_OverwriteOldValues(t *testing.T) {
	// 建立 5 次請求的滑動視窗追蹤器
	tracker := NewTracker(5)

	// 先填滿 5 次失敗
	for i := 0; i < 5; i++ {
		tracker.RecordFailure(nil)
	}

	// 錯誤率應為 100%
	if tracker.ErrorRate() != 1.0 {
		t.Errorf("預期錯誤率為 1.0，實際為 %f", tracker.ErrorRate())
	}

	// 再輸入 5 次成功，舊的失敗應該被覆蓋
	for i := 0; i < 5; i++ {
		tracker.RecordSuccess(5 * time.Millisecond)
	}

	// 錯誤率應為 0%
	if tracker.ErrorRate() != 0.0 {
		t.Errorf("預期錯誤率為 0.0，實際為 %f", tracker.ErrorRate())
	}
}

func TestTracker_Stats(t *testing.T) {
	tracker := NewTracker(10)

	// 記錄不同延遲的成功請求
	durations := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
		40 * time.Millisecond,
		50 * time.Millisecond,
	}

	for _, d := range durations {
		tracker.RecordSuccess(d)
	}

	stats := tracker.GetStats()

	// 驗證總請求數
	if stats.TotalRequests != 5 {
		t.Errorf("預期總請求數為 5，實際為 %d", stats.TotalRequests)
	}

	// 驗證成功數
	if stats.SuccessCount != 5 {
		t.Errorf("預期成功數為 5，實際為 %d", stats.SuccessCount)
	}

	// 驗證失敗數
	if stats.FailureCount != 0 {
		t.Errorf("預期失敗數為 0，實際為 %d", stats.FailureCount)
	}

	// 驗證平均延遲 (應約為 30ms)
	avgLatency := stats.AvgLatency
	if avgLatency < 29*time.Millisecond || avgLatency > 31*time.Millisecond {
		t.Errorf("預期平均延遲約為 30ms，實際為 %v", avgLatency)
	}
}

// =============================================================================
// 任務 1.3: [REFACTOR] 併發安全性測試
// =============================================================================

func TestTracker_Concurrency(t *testing.T) {
	tracker := NewTracker(100)

	var wg sync.WaitGroup
	goroutines := 10
	recordsPerGoroutine := 100

	// 多個 Goroutine 同時寫入
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < recordsPerGoroutine; j++ {
				if j%2 == 0 {
					tracker.RecordSuccess(time.Duration(j) * time.Millisecond)
				} else {
					tracker.RecordFailure(nil)
				}
			}
		}()
	}

	wg.Wait()

	// 不應該 panic，且統計數據應該合理
	stats := tracker.GetStats()
	if stats.TotalRequests == 0 {
		t.Error("併發寫入後，總請求數不應為 0")
	}

	// 錯誤率應該在 0 到 1 之間
	errorRate := tracker.ErrorRate()
	if errorRate < 0 || errorRate > 1 {
		t.Errorf("錯誤率應在 0-1 之間，實際為 %f", errorRate)
	}
}

func TestTracker_ConcurrentReadWrite(t *testing.T) {
	tracker := NewTracker(50)

	var wg sync.WaitGroup

	// Writer goroutines
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				tracker.RecordSuccess(10 * time.Millisecond)
				tracker.RecordFailure(nil)
			}
		}()
	}

	// Reader goroutines
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				_ = tracker.ErrorRate()
				_ = tracker.GetStats()
			}
		}()
	}

	wg.Wait()

	// 只要不 panic 就算通過
}
