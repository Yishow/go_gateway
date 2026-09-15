package collector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-gateway/internal/datalink/collector/health"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 手動觸發
// =============================================================================

// PollNow 立即輪詢指定的點位
func (s *Scheduler) PollNow(pointIDs []string) []CollectedValue {
	return s.PollNowContext(context.Background(), pointIDs)
}

// PollNowContext immediately polls points using the caller's context.
func (s *Scheduler) PollNowContext(ctx context.Context, pointIDs []string) []CollectedValue {
	results := make([]CollectedValue, 0, len(pointIDs))

	s.mu.RLock()
	// 取得點位資訊
	pointsToPoll := make([]pointInfo, 0, len(pointIDs))
	for _, id := range pointIDs {
		if info, exists := s.pointInfos[id]; exists {
			pointsToPoll = append(pointsToPoll, info)
		}
	}
	s.mu.RUnlock()

	// 依設備分組
	devicePoints := make(map[string][]pointInfo)
	for _, info := range pointsToPoll {
		devicePoints[info.DeviceID] = append(devicePoints[info.DeviceID], info)
	}

	// 收集結果
	var resultsMu sync.Mutex
	var wg sync.WaitGroup

	for deviceID, points := range devicePoints {
		wg.Add(1)
		go func(devID string, pts []pointInfo) {
			defer wg.Done()

			s.mu.RLock()
			lock, exists := s.deviceLocks[devID]
			deviceCfg, cfgExists := s.deviceConfigs[devID]
			breaker := s.deviceBreakers[devID]
			s.mu.RUnlock()

			if !exists || !cfgExists || breaker == nil {
				s.mu.Lock()
				s.ensureDeviceLock(devID)
				lock = s.deviceLocks[devID]
				deviceCfg, cfgExists = s.deviceConfigs[devID]
				if s.deviceBreakers[devID] == nil {
					s.deviceBreakers[devID] = health.NewCircuitBreaker(s.config.BreakerConfig)
				}
				breaker = s.deviceBreakers[devID]
				s.mu.Unlock()
				if !cfgExists {
					return
				}
			}

			if breaker != nil && !breaker.AllowRequest() {
				now := time.Now()
				resultsMu.Lock()
				for _, pt := range pts {
					results = append(results, CollectedValue{
						PointID:   pt.ID,
						DeviceID:  devID,
						Timestamp: now,
						Quality:   schema.QualityBad,
						Error:     fmt.Sprintf("設備熔斷中 (狀態: %s)", breaker.State()),
					})
				}
				resultsMu.Unlock()
				return
			}

			lock.Lock()
			defer lock.Unlock()

			conn, err := s.connMgr.GetOrCreate(ctx, devID, deviceCfg.Protocol, deviceCfg.Config)
			if err != nil {
				if breaker != nil {
					breaker.ReportResult(err)
				}

				now := time.Now()
				resultsMu.Lock()
				for _, pt := range pts {
					results = append(results, CollectedValue{
						PointID:   pt.ID,
						DeviceID:  devID,
						Timestamp: now,
						Quality:   schema.QualityBad,
						Error:     fmt.Sprintf("連線失敗: %v", err),
					})
				}
				resultsMu.Unlock()
				return
			}

			for _, pt := range pts {
				req := connector.ReadRequest{
					Address:    pt.Address,
					Function:   pt.Function,
					DataType:   pt.DataType,
					DataFormat: pt.DataFormat,
				}
				req.Count = schema.RegisterCountForDataType(pt.DataType)

				result, readErr := conn.Read(ctx, req)
				if breaker != nil {
					breaker.ReportResult(readErr)
				}
				cv := CollectedValue{
					PointID:   pt.ID,
					DeviceID:  devID,
					Value:     result.Value,
					RawBytes:  result.RawBytes,
					Timestamp: result.Timestamp,
					Quality:   result.Quality,
					Error:     result.Error,
				}
				if readErr != nil {
					cv.Quality = schema.QualityBad
					cv.Error = readErr.Error()
				}

				resultsMu.Lock()
				results = append(results, cv)
				resultsMu.Unlock()
			}
		}(deviceID, points)
	}

	wg.Wait()
	return results
}

func (s *Scheduler) ensureDeviceLocks() {
	for id := range s.deviceConfigs {
		if _, exists := s.deviceLocks[id]; !exists {
			s.deviceLocks[id] = &sync.Mutex{}
		}
	}
}

func (s *Scheduler) ensureDeviceLock(deviceID string) {
	if _, exists := s.deviceLocks[deviceID]; !exists {
		s.deviceLocks[deviceID] = &sync.Mutex{}
	}
}
