package device

import (
	"context"
	"fmt"
	"time"
)

// =============================================================================
// 連線測試
// =============================================================================

// TestConnection 測試設備連線
func (s *Service) TestConnection(ctx context.Context, id string) error {
	// 檢查 ConnectionManager 是否已設定
	if s.connMgr == nil {
		baseErr := fmt.Errorf("連線管理器未初始化，無法測試連線")
		return s.recordTestResult(ctx, id, false, baseErr.Error(), baseErr)
	}

	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("取得設備失敗: %w", err)
	}

	// 嘗試建立連線
	conn, err := s.connMgr.GetOrCreate(ctx, id, device.Protocol, device.ConnectionConfig)
	if err != nil {
		return s.recordTestResult(ctx, id, false, err.Error(), err)
	}

	// 執行連線測試
	err = conn.Protocol.TestConnection(ctx)
	if err != nil {
		return s.recordTestResult(ctx, id, false, err.Error(), err)
	}

	// 記錄測試成功
	return s.recordTestResult(ctx, id, true, "", nil)
}

// TestConnectionResult 連線測試結果
type TestConnectionResult struct {
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	LatencyMs int64     `json:"latency_ms"`
}

// TestConnectionWithResult 測試設備連線並返回詳細結果
func (s *Service) TestConnectionWithResult(ctx context.Context, id string) (*TestConnectionResult, error) {
	if s.connMgr == nil {
		return nil, fmt.Errorf("連線管理器未初始化，無法測試連線")
	}

	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得設備失敗: %w", err)
	}

	start := time.Now()
	result := &TestConnectionResult{
		Timestamp: start,
	}

	// 嘗試建立連線
	conn, err := s.connMgr.GetOrCreate(ctx, id, device.Protocol, device.ConnectionConfig)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.LatencyMs = time.Since(start).Milliseconds()
		if updateErr := s.repo.UpdateTestResult(ctx, id, false, err.Error()); updateErr != nil {
			return result, fmt.Errorf("記錄測試結果失敗: %w", updateErr)
		}
		return result, nil
	}

	// 執行連線測試
	err = conn.Protocol.TestConnection(ctx)
	result.LatencyMs = time.Since(start).Milliseconds()

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		if updateErr := s.repo.UpdateTestResult(ctx, id, false, err.Error()); updateErr != nil {
			return result, fmt.Errorf("記錄測試結果失敗: %w", updateErr)
		}
	} else {
		result.Success = true
		if updateErr := s.repo.UpdateTestResult(ctx, id, true, ""); updateErr != nil {
			return result, fmt.Errorf("記錄測試結果失敗: %w", updateErr)
		}
	}

	return result, nil
}

func (s *Service) recordTestResult(ctx context.Context, id string, success bool, errMsg string, baseErr error) error {
	if err := s.repo.UpdateTestResult(ctx, id, success, errMsg); err != nil {
		if baseErr != nil {
			return fmt.Errorf("%w；記錄測試結果失敗: %w", baseErr, err)
		}
		return fmt.Errorf("記錄測試結果失敗: %w", err)
	}
	return baseErr
}
