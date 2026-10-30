package collector

import (
	"context"
	"fmt"
	"sync"

	"go-gateway/internal/datalink/collector/health"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/wirebudget"
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
	if ctx == nil {
		ctx = context.Background()
	}
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
		wg.Go(func() {
			s.mu.RLock()
			lock, exists := s.deviceLocks[deviceID]
			deviceCfg, cfgExists := s.deviceConfigs[deviceID]
			breaker := s.deviceBreakers[deviceID]
			s.mu.RUnlock()

			if !exists || !cfgExists || breaker == nil {
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

			if breaker != nil && !breaker.AllowRequest() {
				acquisitionID := s.nextAcquisitionID()
				resultsMu.Lock()
				for _, pt := range points {
					results = append(results, s.collectedValueWithReason(deviceCfg, pt, connector.ReadResult{
						Quality: schema.QualityBad,
						Error:   fmt.Sprintf("設備熔斷中 (狀態: %s)", breaker.State()),
					}, nil, acquisitionID, s.now(), "circuit-open"))
				}
				resultsMu.Unlock()
				return
			}

			if err := lock.LockContext(ctx); err != nil {
				acquisitionID := s.nextAcquisitionID()
				resultsMu.Lock()
				for _, pt := range points {
					results = append(results, s.collectedValueWithReason(deviceCfg, pt, connector.ReadResult{
						Quality: schema.QualityBad,
						Error:   err.Error(),
					}, err, acquisitionID, s.now(), "request-canceled"))
				}
				resultsMu.Unlock()
				return
			}
			defer lock.Unlock()

			conn, err := s.connMgr.GetOrCreate(ctx, deviceID, deviceCfg.Protocol, deviceCfg.Config)
			if err != nil {
				if breaker != nil {
					breaker.ReportResult(err)
				}

				acquisitionID := s.nextAcquisitionID()
				resultsMu.Lock()
				for _, pt := range points {
					results = append(results, s.collectedValueWithReason(deviceCfg, pt, connector.ReadResult{
						Quality: schema.QualityBad,
						Error:   fmt.Sprintf("連線失敗: %v", err),
					}, nil, acquisitionID, s.now(), "connection-failed"))
				}
				resultsMu.Unlock()
				return
			}

			acquisitionID := s.nextAcquisitionID()
			for _, pt := range points {
				req := connector.ReadRequest{
					Address:    pt.Address,
					Function:   pt.Function,
					DataType:   pt.DataType,
					DataFormat: pt.DataFormat,
					Profile:    pt.readProfile,
				}
				count, countErr := connector.PointReadCountForAddress(conn.ProtocolType, pt.Address, pt.DataType)
				if countErr != nil {
					cv := s.collectedValue(managedConfig(conn), pt, connector.ReadResult{
						Quality: schema.QualityBad,
						Error:   countErr.Error(),
					}, countErr, acquisitionID, s.now())
					resultsMu.Lock()
					results = append(results, cv)
					resultsMu.Unlock()
					continue
				}
				req.Count = count

				readCtx, cancelRead := s.acquisitionReadContext(ctx, conn)
				result, readErr := conn.Read(readCtx, req)
				cancelRead()
				if breaker != nil {
					breaker.ReportResult(readErr)
				}
				cv := s.collectedValue(managedConfig(conn), pt, result, readErr, acquisitionID, s.now())

				resultsMu.Lock()
				results = append(results, cv)
				resultsMu.Unlock()
			}
		})
	}

	wg.Wait()
	return results
}

func (s *Scheduler) ensureDeviceLocks() {
	for id := range s.deviceConfigs {
		if _, exists := s.deviceLocks[id]; !exists {
			s.deviceLocks[id] = &wirebudget.Mutex{}
		}
	}
}

func (s *Scheduler) ensureDeviceLock(deviceID string) {
	if _, exists := s.deviceLocks[deviceID]; !exists {
		s.deviceLocks[deviceID] = &wirebudget.Mutex{}
	}
}
