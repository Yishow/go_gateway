package collector

import (
	"context"
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

	if s.stopDone != nil {
		select {
		case <-s.stopDone:
		default:
			return fmt.Errorf("scheduler is still stopping")
		}
	}
	s.stopDone = nil
	s.pollCtx, s.pollCancel = context.WithCancel(context.Background())
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
func (s *Scheduler) Stop() error {
	return s.StopContext(context.Background())
}

// StopContext initiates one shutdown and waits only until ctx expires. A
// deadline is a notification limit, not proof that protocol workers stopped.
func (s *Scheduler) StopContext(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.running = false
		close(s.stopCh)
		if s.pollCancel != nil {
			s.pollCancel()
		}
		for _, gt := range s.groupTickers {
			close(gt.stopCh)
			gt.ticker.Stop()
		}
		s.groupTickers = make(map[string]*groupTicker)
		done := make(chan struct{})
		s.stopDone = done
		go func() {
			s.wg.Wait()
			close(done)
		}()
	}
	done := s.stopDone
	s.mu.Unlock()
	return waitScheduler(ctx, done)
}

// WaitStopped observes the original shutdown without initiating it again.
func (s *Scheduler) WaitStopped(ctx context.Context) error {
	s.mu.RLock()
	done := s.stopDone
	running := s.running
	s.mu.RUnlock()
	if running {
		return fmt.Errorf("scheduler has not been stopped")
	}
	return waitScheduler(ctx, done)
}

func waitScheduler(ctx context.Context, done <-chan struct{}) error {
	if done == nil {
		return nil
	}
	select {
	case <-done:
		return nil
	default:
	}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
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
