package collector

import (
	"fmt"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 運行控制
// =============================================================================

// Start 啟動排程器
func (s *Scheduler) Start(groups []*schema.PollingGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("排程器已在運行")
	}

	s.ensureDeviceLocks()

	s.running = true
	s.stopCh = make(chan struct{})

	// 啟動所有已啟用的輪詢群組
	for _, group := range groups {
		if group.Enabled {
			s.startGroupTicker(group)
		}
	}

	return nil
}

// Stop 停止排程器
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	// 關閉停止信號
	close(s.stopCh)

	// 停止所有群組 Ticker
	for _, gt := range s.groupTickers {
		close(gt.stopCh)
		gt.ticker.Stop()
	}
	s.groupTickers = make(map[string]*groupTicker)

	// 等待所有 goroutine 結束
	s.wg.Wait()

	s.running = false
}

// IsRunning 檢查是否正在運行
func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// ValueChannel 取得值輸出通道
func (s *Scheduler) ValueChannel() <-chan CollectedValue {
	return s.valueChan
}
