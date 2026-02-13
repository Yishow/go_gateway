package collector

import (
	"sync"

	"go-gateway/internal/datalink/collector/health"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 配置方法
// =============================================================================

// AddDevice 新增設備配置
func (s *Scheduler) AddDevice(device *schema.Device) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.deviceConfigs[device.ID] = deviceConfig{
		ID:       device.ID,
		Protocol: device.Protocol,
		Config:   device.ConnectionConfig,
	}

	// 建立設備鎖
	if _, exists := s.deviceLocks[device.ID]; !exists {
		s.deviceLocks[device.ID] = &sync.Mutex{}
	}

	// 建立設備熔斷器
	if _, exists := s.deviceBreakers[device.ID]; !exists {
		s.deviceBreakers[device.ID] = health.NewCircuitBreaker(s.config.BreakerConfig)
	}
}

// RemoveDevice 移除設備配置
func (s *Scheduler) RemoveDevice(deviceID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.deviceConfigs, deviceID)
	delete(s.deviceBreakers, deviceID)

	// 移除相關點位
	for id, info := range s.pointInfos {
		if info.DeviceID == deviceID {
			delete(s.pointInfos, id)
		}
	}
}

// AddPoint 新增點位
func (s *Scheduler) AddPoint(point *schema.Point) {
	s.mu.Lock()
	defer s.mu.Unlock()

	groupID := ""
	if point.PollingGroupID != nil {
		groupID = *point.PollingGroupID
	}

	s.pointInfos[point.ID] = pointInfo{
		ID:             point.ID,
		DeviceID:       point.DeviceID,
		Address:        point.Address,
		Function:       point.Function,
		DataType:       point.DataType,
		PollingGroupID: groupID,
	}
}

// RemovePoint 移除點位
func (s *Scheduler) RemovePoint(pointID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.pointInfos, pointID)
}

// AddPollingGroup 新增輪詢群組
func (s *Scheduler) AddPollingGroup(group *schema.PollingGroup) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ensureDeviceLocks()

	// 如果已在運行，啟動新群組的 Ticker
	if s.running && group.Enabled {
		s.startGroupTicker(group)
	}
}

// RemovePollingGroup 移除輪詢群組
func (s *Scheduler) RemovePollingGroup(groupID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if gt, exists := s.groupTickers[groupID]; exists {
		close(gt.stopCh)
		gt.ticker.Stop()
		delete(s.groupTickers, groupID)
	}
}
