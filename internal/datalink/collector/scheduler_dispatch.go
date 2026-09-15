package collector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-gateway/internal/datalink/collector/health"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 內部方法
// =============================================================================

// startGroupTicker 啟動群組 Ticker (必須持有 mu 鎖)
func (s *Scheduler) startGroupTicker(group *schema.PollingGroup) {
	interval := time.Duration(group.IntervalMs) * time.Millisecond

	gt := &groupTicker{
		group:  group,
		ticker: time.NewTicker(interval),
		stopCh: make(chan struct{}),
		done:   make(chan struct{}),
	}

	s.groupTickers[group.ID] = gt

	s.wg.Add(1)
	go s.runGroupTicker(gt)
}

// runGroupTicker 運行群組 Ticker
func (s *Scheduler) runGroupTicker(gt *groupTicker) {
	defer s.wg.Done()
	defer close(gt.done)

	for {
		select {
		case <-gt.stopCh:
			return
		case <-gt.ticker.C:
			s.pollGroup(gt.group.ID)
		}
	}
}

// pollGroup 輪詢群組
func (s *Scheduler) pollGroup(groupID string) {
	s.mu.RLock()
	// 取得群組中的所有點位
	pointsToPoll := make([]pointInfo, 0)
	for _, info := range s.pointInfos {
		if info.PollingGroupID == groupID {
			pointsToPoll = append(pointsToPoll, info)
		}
	}
	s.mu.RUnlock()

	// 依設備分組
	devicePoints := make(map[string][]pointInfo)
	for _, info := range pointsToPoll {
		devicePoints[info.DeviceID] = append(devicePoints[info.DeviceID], info)
	}

	// 各設備並行處理，同一設備內依序執行。
	var wg sync.WaitGroup
	for deviceID, points := range devicePoints {
		wg.Add(1)
		go func(devID string, pts []pointInfo) {
			defer wg.Done()
			s.pollDevicePoints(devID, pts)
		}(deviceID, points)
	}
	wg.Wait()
}

// pollDevicePoints 輪詢設備的點位
func (s *Scheduler) pollDevicePoints(deviceID string, points []pointInfo) {
	// 取得設備鎖與熔斷器
	s.mu.RLock()
	lock, exists := s.deviceLocks[deviceID]
	deviceCfg, cfgExists := s.deviceConfigs[deviceID]
	breaker := s.deviceBreakers[deviceID]
	s.mu.RUnlock()

	if !exists || !cfgExists {
		s.mu.Lock()
		s.ensureDeviceLock(deviceID)
		lock = s.deviceLocks[deviceID]
		deviceCfg, cfgExists = s.deviceConfigs[deviceID]
		if s.deviceBreakers[deviceID] == nil {
			s.deviceBreakers[deviceID] = health.NewCircuitBreaker(s.config.BreakerConfig)
		}
		breaker = s.deviceBreakers[deviceID]
		s.mu.Unlock()
		if !cfgExists {
			return
		}
	}

	// 檢查熔斷器是否允許請求
	if breaker != nil && !breaker.AllowRequest() {
		// 熔斷中，跳過採集並產生熔斷錯誤
		now := time.Now()
		for _, pt := range points {
			s.emitValue(CollectedValue{
				PointID:   pt.ID,
				DeviceID:  deviceID,
				Timestamp: now,
				Quality:   schema.QualityBad,
				Error:     fmt.Sprintf("設備熔斷中 (狀態: %s)", breaker.State()),
			})
		}
		return
	}

	// 鎖定設備確保串列存取
	lock.Lock()
	defer lock.Unlock()

	// 取得或建立連線
	ctx := context.Background()
	conn, err := s.connMgr.GetOrCreate(ctx, deviceID, deviceCfg.Protocol, deviceCfg.Config)
	if err != nil {
		// 連線失敗，報告給熔斷器
		if breaker != nil {
			breaker.ReportResult(err)
		}

		// 為所有點位產生錯誤結果
		now := time.Now()
		for _, pt := range points {
			s.emitValue(CollectedValue{
				PointID:   pt.ID,
				DeviceID:  deviceID,
				Timestamp: now,
				Quality:   schema.QualityBad,
				Error:     fmt.Sprintf("連線失敗: %v", err),
			})
		}
		return
	}

	// 逐一讀取點位，並報告結果給熔斷器
	for _, pt := range points {
		err := s.pollPoint(conn, pt)

		// 報告結果給熔斷器
		if breaker != nil {
			breaker.ReportResult(err)
		}
	}
}
