package collector

import (
	"log"
	"time"

	"go-gateway/internal/datalink/schema"
)

const (
	autoIntervalWindowSize  = 5
	autoIntervalStepUnitMs  = 100
	autoIntervalMinShrinkMs = 500
	defaultIntervalMs       = 1000
)

// DeviceIntervalStatus 描述設備的設定輪詢間隔與目前有效輪詢間隔（毫秒）。
type DeviceIntervalStatus struct {
	ConfiguredIntervalMs int
	EffectiveIntervalMs  int
	AutoInterval         bool
}

type deviceIntervalState struct {
	autoInterval bool
	configuredMs int
	effectiveMs  int
	recentScans  []time.Duration
	onTimeStreak int
}

// DeviceIntervalStatus 回傳指定設備的設定輪詢間隔與目前有效輪詢間隔。
func (s *Scheduler) DeviceIntervalStatus(deviceID string) (DeviceIntervalStatus, bool) {
	if s == nil || deviceID == "" {
		return DeviceIntervalStatus{}, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, hasDevice := s.deviceConfigs[deviceID]
	state, hasState := s.deviceIntervals[deviceID]
	if !hasDevice && !hasState {
		return DeviceIntervalStatus{}, false
	}

	configuredMs := s.configuredIntervalMsLocked(deviceID)
	effectiveMs := configuredMs
	autoInterval := false
	if state != nil {
		autoInterval = state.autoInterval
		if state.autoInterval && state.effectiveMs > 0 {
			effectiveMs = max(configuredMs, state.effectiveMs)
		}
	}
	return DeviceIntervalStatus{
		ConfiguredIntervalMs: configuredMs,
		EffectiveIntervalMs:  effectiveMs,
		AutoInterval:         autoInterval,
	}, true
}

// RecordScanDuration 記錄單輪設備掃描耗時並依自動調整規則更新有效間隔（毫秒）。
func (s *Scheduler) RecordScanDuration(deviceID string, duration time.Duration) int {
	if s == nil || deviceID == "" {
		return 0
	}
	s.mu.Lock()
	configuredMs, effectiveMs, stretchedFrom, maxScan := s.recordScanDurationLocked(deviceID, duration)
	s.mu.Unlock()

	if stretchedFrom > 0 && effectiveMs > stretchedFrom {
		log.Printf(
			"[collector] device %s auto-interval stretched from %dms to %dms (configured=%dms, max_recent_scan=%s)",
			deviceID,
			stretchedFrom,
			effectiveMs,
			configuredMs,
			maxScan,
		)
	}
	return effectiveMs
}

func (s *Scheduler) recordScanDurationLocked(deviceID string, duration time.Duration) (configuredMs, effectiveMs, stretchedFrom int, maxScan time.Duration) {
	configuredMs = s.configuredIntervalMsLocked(deviceID)
	state := s.ensureDeviceIntervalStateLocked(deviceID)
	if state.configuredMs <= 0 {
		state.configuredMs = configuredMs
	}
	if state.effectiveMs <= 0 {
		state.effectiveMs = configuredMs
	}
	if !state.autoInterval {
		state.configuredMs = configuredMs
		state.effectiveMs = configuredMs
		state.recentScans = nil
		state.onTimeStreak = 0
		return configuredMs, configuredMs, 0, 0
	}

	if duration < 0 {
		duration = 0
	}
	state.recentScans = append(state.recentScans, duration)
	if len(state.recentScans) > autoIntervalWindowSize {
		state.recentScans = state.recentScans[len(state.recentScans)-autoIntervalWindowSize:]
	}

	configuredDur := time.Duration(configuredMs) * time.Millisecond
	if duration > configuredDur {
		state.onTimeStreak = 0
	} else if state.effectiveMs > configuredMs {
		state.onTimeStreak++
	}
	for _, d := range state.recentScans {
		if d > maxScan {
			maxScan = d
		}
	}
	targetMs := max(configuredMs, ceilDurationTimes1Point2To100Ms(maxScan))
	nextMs := state.effectiveMs
	if targetMs > nextMs {
		stretchedFrom = nextMs
		nextMs = targetMs
	} else if targetMs < nextMs && maxScan > configuredDur {
		// An old overrun left the five-scan window; retain the current
		// window's required margin without waiting for on-time scans.
		nextMs = targetMs
	} else if nextMs > configuredMs && state.onTimeStreak >= autoIntervalWindowSize {
		state.onTimeStreak = 0
		step := max(autoIntervalMinShrinkMs, ceilIntTo100Ms((nextMs-configuredMs)/2))
		nextMs = max(targetMs, nextMs-step)
	}
	if nextMs != state.effectiveMs {
		state.effectiveMs = nextMs
		s.resetDeviceGroupTickerLocked(deviceID, nextMs)
	}
	return configuredMs, state.effectiveMs, stretchedFrom, maxScan
}

func (s *Scheduler) setDeviceAutoIntervalLocked(deviceID string, autoInterval bool) {
	if deviceID == "" {
		return
	}
	state := s.ensureDeviceIntervalStateLocked(deviceID)
	configuredMs := s.configuredIntervalMsLocked(deviceID)
	state.autoInterval = autoInterval
	state.configuredMs = configuredMs
	if !autoInterval || state.effectiveMs < configuredMs {
		state.effectiveMs = configuredMs
		if !autoInterval {
			state.recentScans = nil
			state.onTimeStreak = 0
			s.resetDeviceGroupTickerLocked(deviceID, configuredMs)
		}
	}
}

func (s *Scheduler) storePollingGroupLocked(group *schema.PollingGroup) {
	if group == nil || group.ID == "" {
		return
	}
	if s.configuredGroups == nil {
		s.configuredGroups = make(map[string]*schema.PollingGroup)
	}
	copied := *group
	s.configuredGroups[group.ID] = &copied

	if group.IntervalMs <= 0 {
		return
	}
	if group.OwnerDeviceID != nil && *group.OwnerDeviceID != "" {
		devID := *group.OwnerDeviceID
		state := s.ensureDeviceIntervalStateLocked(devID)
		if state.configuredMs != group.IntervalMs {
			state.configuredMs = group.IntervalMs
			state.effectiveMs = group.IntervalMs
			state.recentScans = nil
			state.onTimeStreak = 0
		}
	}
}

func (s *Scheduler) removePollingGroupConfigLocked(groupID string) {
	if s.configuredGroups != nil {
		delete(s.configuredGroups, groupID)
	}
}

func (s *Scheduler) ensureDeviceIntervalStateLocked(deviceID string) *deviceIntervalState {
	if s.deviceIntervals == nil {
		s.deviceIntervals = make(map[string]*deviceIntervalState)
	}
	state, ok := s.deviceIntervals[deviceID]
	if !ok || state == nil {
		state = &deviceIntervalState{}
		s.deviceIntervals[deviceID] = state
	}
	return state
}

func (s *Scheduler) configuredIntervalMsLocked(deviceID string) int {
	for _, gt := range s.groupTickers {
		if gt != nil && gt.group != nil && gt.group.OwnerDeviceID != nil && *gt.group.OwnerDeviceID == deviceID && gt.group.IntervalMs > 0 {
			return gt.group.IntervalMs
		}
	}
	for _, g := range s.configuredGroups {
		if g != nil && g.OwnerDeviceID != nil && *g.OwnerDeviceID == deviceID && g.IntervalMs > 0 {
			return g.IntervalMs
		}
	}
	for _, pt := range s.pointInfos {
		if pt.DeviceID != deviceID || pt.PollingGroupID == "" {
			continue
		}
		if gt := s.groupTickers[pt.PollingGroupID]; gt != nil && gt.group != nil && gt.group.IntervalMs > 0 {
			return gt.group.IntervalMs
		}
		if g := s.configuredGroups[pt.PollingGroupID]; g != nil && g.IntervalMs > 0 {
			return g.IntervalMs
		}
	}
	if state := s.deviceIntervals[deviceID]; state != nil && state.configuredMs > 0 {
		return state.configuredMs
	}
	return defaultIntervalMs
}

func (s *Scheduler) effectiveGroupIntervalLocked(group *schema.PollingGroup) time.Duration {
	if group == nil || group.IntervalMs <= 0 {
		return time.Duration(defaultIntervalMs) * time.Millisecond
	}
	intervalMs := group.IntervalMs
	if group.OwnerDeviceID != nil && *group.OwnerDeviceID != "" {
		if state := s.deviceIntervals[*group.OwnerDeviceID]; state != nil && state.autoInterval && state.effectiveMs > intervalMs {
			intervalMs = state.effectiveMs
		}
	}
	return time.Duration(intervalMs) * time.Millisecond
}

func (s *Scheduler) resetDeviceGroupTickerLocked(deviceID string, intervalMs int) {
	if intervalMs <= 0 {
		return
	}
	dur := time.Duration(intervalMs) * time.Millisecond
	for _, gt := range s.groupTickers {
		if gt == nil || gt.group == nil || gt.ticker == nil {
			continue
		}
		if gt.group.OwnerDeviceID != nil && *gt.group.OwnerDeviceID == deviceID {
			gt.ticker.Reset(dur)
		}
	}
}

func ceilDurationTimes1Point2To100Ms(d time.Duration) int {
	if d <= 0 {
		return 0
	}
	scaledNs := (d.Nanoseconds()*6 + 4) / 5
	unitNs := int64(autoIntervalStepUnitMs) * int64(time.Millisecond)
	steps := (scaledNs + unitNs - 1) / unitNs
	return int(steps * int64(autoIntervalStepUnitMs))
}

func ceilIntTo100Ms(ms int) int {
	if ms <= 0 {
		return 0
	}
	return ((ms + autoIntervalStepUnitMs - 1) / autoIntervalStepUnitMs) * autoIntervalStepUnitMs
}
