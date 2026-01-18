// Package collector 提供資料收集排程器功能。
//
// 本套件實作輪詢群組調度、設備並發控制、重試/退避策略等功能。
package collector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 收集結果
// =============================================================================

// CollectedValue 收集到的值
type CollectedValue struct {
	// PointID 點位 ID
	PointID string

	// DeviceID 設備 ID
	DeviceID string

	// TagID 標籤 ID (若有映射)
	TagID string

	// Value 收集到的值
	Value interface{}

	// RawBytes 原始位元組
	RawBytes []byte

	// Timestamp 收集時間戳記
	Timestamp time.Time

	// Quality 資料品質標誌
	Quality schema.QualityFlag

	// Error 錯誤訊息
	Error string
}

// =============================================================================
// 排程器配置
// =============================================================================

// SchedulerConfig 排程器配置
type SchedulerConfig struct {
	// DefaultRetryCount 預設重試次數
	DefaultRetryCount int

	// DefaultRetryDelay 預設重試延遲
	DefaultRetryDelay time.Duration

	// MaxConcurrentPerDevice 每設備最大並發數
	MaxConcurrentPerDevice int

	// ValueBufferSize 值緩衝區大小
	ValueBufferSize int
}

// DefaultSchedulerConfig 預設排程器配置
func DefaultSchedulerConfig() SchedulerConfig {
	return SchedulerConfig{
		DefaultRetryCount:      3,
		DefaultRetryDelay:      1 * time.Second,
		MaxConcurrentPerDevice: 1, // 串列存取避免協議衝突
		ValueBufferSize:        1000,
	}
}

// =============================================================================
// 排程器
// =============================================================================

// Scheduler 資料收集排程器
type Scheduler struct {
	config  SchedulerConfig
	connMgr *connector.ConnectionManager

	// 輪詢群組 Ticker
	groupTickers map[string]*groupTicker

	// 值輸出通道
	valueChan chan CollectedValue

	// 設備連接器快取
	deviceConfigs map[string]deviceConfig

	// 點位資訊快取
	pointInfos map[string]pointInfo

	// 設備鎖 (確保同一設備串列存取)
	deviceLocks map[string]*sync.Mutex

	// 狀態
	running bool
	mu      sync.RWMutex

	// 停止信號
	stopCh chan struct{}
	wg     sync.WaitGroup
}

// groupTicker 輪詢群組 Ticker
type groupTicker struct {
	group  *schema.PollingGroup
	ticker *time.Ticker
	stopCh chan struct{}
}

// deviceConfig 設備配置快取
type deviceConfig struct {
	ID       string
	Protocol schema.ProtocolType
	Config   string
}

// pointInfo 點位資訊快取
type pointInfo struct {
	ID             string
	DeviceID       string
	Address        string
	Function       string
	DataType       schema.DataType
	PollingGroupID string
}

// NewScheduler 建立新的排程器
func NewScheduler(config SchedulerConfig, connMgr *connector.ConnectionManager) *Scheduler {
	if connMgr == nil {
		connMgr = connector.GetConnectionManager()
	}

	return &Scheduler{
		config:        config,
		connMgr:       connMgr,
		groupTickers:  make(map[string]*groupTicker),
		valueChan:     make(chan CollectedValue, config.ValueBufferSize),
		deviceConfigs: make(map[string]deviceConfig),
		pointInfos:    make(map[string]pointInfo),
		deviceLocks:   make(map[string]*sync.Mutex),
		stopCh:        make(chan struct{}),
	}
}

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
}

// RemoveDevice 移除設備配置
func (s *Scheduler) RemoveDevice(deviceID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.deviceConfigs, deviceID)

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
	}

	s.groupTickers[group.ID] = gt

	s.wg.Add(1)
	go s.runGroupTicker(gt)
}

// runGroupTicker 運行群組 Ticker
func (s *Scheduler) runGroupTicker(gt *groupTicker) {
	defer s.wg.Done()

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

	// 並行處理各設備 (但同一設備內串列)
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
	// 取得設備鎖
	s.mu.RLock()
	lock, exists := s.deviceLocks[deviceID]
	deviceCfg, cfgExists := s.deviceConfigs[deviceID]
	s.mu.RUnlock()

	if !exists || !cfgExists {
		return
	}

	// 鎖定設備確保串列存取
	lock.Lock()
	defer lock.Unlock()

	// 取得或建立連線
	ctx := context.Background()
	conn, err := s.connMgr.GetOrCreate(ctx, deviceID, deviceCfg.Protocol, deviceCfg.Config)
	if err != nil {
		// 連線失敗，為所有點位產生錯誤結果
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

	// 逐一讀取點位
	for _, pt := range points {
		s.pollPoint(conn, pt)
	}
}

// pollPoint 輪詢單一點位
func (s *Scheduler) pollPoint(conn *connector.ManagedConnection, pt pointInfo) {
	req := connector.ReadRequest{
		Address:  pt.Address,
		Function: pt.Function,
		DataType: pt.DataType,
		Count:    1,
	}

	var result connector.ReadResult
	var err error

	// 重試邏輯
	for attempt := 0; attempt <= s.config.DefaultRetryCount; attempt++ {
		if attempt > 0 {
			time.Sleep(s.config.DefaultRetryDelay)
		}

		result, err = conn.Read(context.Background(), req)
		if err == nil && result.Quality == schema.QualityGood {
			break
		}
	}

	// 發送結果
	cv := CollectedValue{
		PointID:   pt.ID,
		DeviceID:  pt.DeviceID,
		Value:     result.Value,
		RawBytes:  result.RawBytes,
		Timestamp: result.Timestamp,
		Quality:   result.Quality,
		Error:     result.Error,
	}

	if err != nil {
		cv.Quality = schema.QualityBad
		cv.Error = err.Error()
	}

	s.emitValue(cv)
}

// emitValue 發送收集到的值
func (s *Scheduler) emitValue(cv CollectedValue) {
	select {
	case s.valueChan <- cv:
	default:
		// 緩衝區已滿，丟棄舊值
		select {
		case <-s.valueChan:
			s.valueChan <- cv
		default:
		}
	}
}

// =============================================================================
// 手動觸發
// =============================================================================

// PollNow 立即輪詢指定的點位
func (s *Scheduler) PollNow(pointIDs []string) []CollectedValue {
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
			s.mu.RUnlock()

			if !exists || !cfgExists {
				return
			}

			lock.Lock()
			defer lock.Unlock()

			ctx := context.Background()
			conn, err := s.connMgr.GetOrCreate(ctx, devID, deviceCfg.Protocol, deviceCfg.Config)
			if err != nil {
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
					Address:  pt.Address,
					Function: pt.Function,
					DataType: pt.DataType,
					Count:    1,
				}

				result, readErr := conn.Read(ctx, req)
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
