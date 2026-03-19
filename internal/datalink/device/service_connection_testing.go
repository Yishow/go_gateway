package device

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
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

// TestConnectionStageStatus 連線測試階段狀態。
type TestConnectionStageStatus string

const (
	// TestConnectionStageSuccess 代表階段成功。
	TestConnectionStageSuccess TestConnectionStageStatus = "success"
	// TestConnectionStageFailed 代表階段失敗。
	TestConnectionStageFailed TestConnectionStageStatus = "failed"
	// TestConnectionStageSkipped 代表階段略過。
	TestConnectionStageSkipped TestConnectionStageStatus = "skipped"
)

// TestConnectionStageResult 連線測試單一階段結果。
type TestConnectionStageResult struct {
	Status    TestConnectionStageStatus `json:"status"`
	Message   string                    `json:"message,omitempty"`
	Error     string                    `json:"error,omitempty"`
	LatencyMs int64                     `json:"latency_ms"`
}

// TestConnectionResult 連線測試結果
type TestConnectionResult struct {
	Success     bool                      `json:"success"`
	Error       string                    `json:"error,omitempty"`
	Timestamp   time.Time                 `json:"timestamp"`
	LatencyMs   int64                     `json:"latency_ms"`
	Connect     TestConnectionStageResult `json:"connect"`
	Probe       TestConnectionStageResult `json:"probe"`
	CanActivate bool                      `json:"can_activate"`
	CanCollect  bool                      `json:"can_collect"`
}

// DraftTestConnectionRequest 代表尚未儲存的設備草稿連線測試請求。
type DraftTestConnectionRequest struct {
	Protocol         schema.ProtocolType    `json:"protocol"`
	ConnectionConfig map[string]interface{} `json:"connection_config"`
}

// testConnectionTimeout 連線測試的整體逾時上限
const testConnectionTimeout = 15 * time.Second

// TestConnectionWithResult 測試設備連線並返回詳細結果
func (s *Service) TestConnectionWithResult(ctx context.Context, id string) (*TestConnectionResult, error) {
	if s.connMgr == nil {
		return nil, fmt.Errorf("連線管理器未初始化，無法測試連線")
	}

	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得設備失敗: %w", err)
	}

	return s.runConnectionTest(ctx, device, true, true)
}

// TestDraftConnectionWithResult 測試尚未儲存的設備草稿連線並返回詳細結果。
func (s *Service) TestDraftConnectionWithResult(
	ctx context.Context,
	req DraftTestConnectionRequest,
) (*TestConnectionResult, error) {
	if err := validateConnectionConfig(req.Protocol, req.ConnectionConfig); err != nil {
		return nil, fmt.Errorf("連線配置無效: %w", err)
	}

	configJSON, err := json.Marshal(req.ConnectionConfig)
	if err != nil {
		return nil, fmt.Errorf("序列化連線配置失敗: %w", err)
	}

	device := &schema.Device{
		ID:               "__draft__",
		Protocol:         req.Protocol,
		ConnectionConfig: string(configJSON),
	}

	return s.runConnectionTest(ctx, device, false, false)
}

func (s *Service) runConnectionTest(
	ctx context.Context,
	device *schema.Device,
	persistResult bool,
	useConnectionManager bool,
) (*TestConnectionResult, error) {
	// 為整體連線測試流程設定逾時上限
	testCtx, cancel := context.WithTimeout(ctx, testConnectionTimeout)
	defer cancel()

	start := time.Now()
	result := &TestConnectionResult{
		Timestamp: start,
		Connect: TestConnectionStageResult{
			Status: TestConnectionStageSkipped,
		},
		Probe: TestConnectionStageResult{
			Status: TestConnectionStageSkipped,
		},
	}

	connectStartedAt := time.Now()
	var protocol connector.Protocol
	var err error
	if useConnectionManager {
		if s.connMgr == nil {
			return nil, fmt.Errorf("連線管理器未初始化，無法測試連線")
		}
		conn, connErr := s.connMgr.GetOrCreate(testCtx, device.ID, device.Protocol, device.ConnectionConfig)
		err = connErr
		if conn != nil {
			protocol = conn.Protocol
		}
	} else {
		protocol, err = connector.Get(testCtx, device.Protocol, device.ConnectionConfig)
	}
	result.Connect.LatencyMs = time.Since(connectStartedAt).Milliseconds()
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.LatencyMs = time.Since(start).Milliseconds()
		result.Connect.Status = TestConnectionStageFailed
		result.Connect.Error = err.Error()
		if persistResult {
			if updateErr := s.repo.UpdateTestResult(ctx, device.ID, false, result.Error); updateErr != nil {
				return result, fmt.Errorf("記錄測試結果失敗: %w", updateErr)
			}
		}
		return result, nil
	}
	if !useConnectionManager {
		defer protocol.Close()
	}

	result.Connect.Status = TestConnectionStageSuccess
	result.Connect.Message = "connect ok"

	probeStartedAt := time.Now()
	probeEnabled, probeErr := probeReadWithProtocol(testCtx, protocol, device)
	result.Probe.LatencyMs = time.Since(probeStartedAt).Milliseconds()
	result.LatencyMs = time.Since(start).Milliseconds()

	switch {
	case probeErr != nil:
		result.Success = false
		result.Error = fmt.Sprintf("讀取探測失敗: %s", probeErr.Error())
		result.Probe.Status = TestConnectionStageFailed
		result.Probe.Error = probeErr.Error()
		if persistResult {
			if updateErr := s.repo.UpdateTestResult(ctx, device.ID, false, result.Error); updateErr != nil {
				return result, fmt.Errorf("記錄測試結果失敗: %w", updateErr)
			}
		}
	case !probeEnabled:
		result.Success = true
		result.Probe.Status = TestConnectionStageSkipped
		result.Probe.Message = "probe skipped"
		result.CanActivate = true
		result.CanCollect = true
		if persistResult {
			if updateErr := s.repo.UpdateTestResult(ctx, device.ID, true, ""); updateErr != nil {
				return result, fmt.Errorf("記錄測試結果失敗: %w", updateErr)
			}
		}
	default:
		result.Success = true
		result.Probe.Status = TestConnectionStageSuccess
		result.Probe.Message = "probe ok"
		result.CanActivate = true
		result.CanCollect = true
		if persistResult {
			if updateErr := s.repo.UpdateTestResult(ctx, device.ID, true, ""); updateErr != nil {
				return result, fmt.Errorf("記錄測試結果失敗: %w", updateErr)
			}
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
