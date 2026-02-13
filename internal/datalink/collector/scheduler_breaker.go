package collector

import "go-gateway/internal/datalink/collector/health"

// =============================================================================
// 熔斷器相關方法
// =============================================================================

// GetDeviceBreakerState 取得設備熔斷器狀態
func (s *Scheduler) GetDeviceBreakerState(deviceID string) (health.State, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	breaker, exists := s.deviceBreakers[deviceID]
	if !exists {
		return "", false
	}
	return breaker.State(), true
}

// GetDeviceBreakerStats 取得設備熔斷器統計
func (s *Scheduler) GetDeviceBreakerStats(deviceID string) (health.BreakerStats, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	breaker, exists := s.deviceBreakers[deviceID]
	if !exists {
		return health.BreakerStats{}, false
	}
	return breaker.GetStats(), true
}

// GetAllDeviceBreakerStates 取得所有設備熔斷器狀態
func (s *Scheduler) GetAllDeviceBreakerStates() map[string]health.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	states := make(map[string]health.State)
	for id, breaker := range s.deviceBreakers {
		states[id] = breaker.State()
	}
	return states
}

// ResetDeviceBreaker 重置設備熔斷器
func (s *Scheduler) ResetDeviceBreaker(deviceID string) bool {
	s.mu.RLock()
	breaker, exists := s.deviceBreakers[deviceID]
	s.mu.RUnlock()

	if !exists {
		return false
	}

	breaker.Reset()
	return true
}
