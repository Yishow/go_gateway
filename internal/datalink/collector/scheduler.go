// Package collector 提供資料收集排程器功能。
//
// 本套件實作輪詢群組調度、設備並發控制、重試/退避策略、
// 以及熔斷器機制等功能。
package collector

import (
	"context"
	"sync"
	"time"

	"go-gateway/internal/datalink/collector/health"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"

	"github.com/google/uuid"
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

	// AcquisitionID identifies one poll attempt across retries and mapped tags.
	AcquisitionID string

	// ObservedAt is the typed observation time selected from a trusted source
	// timestamp or the gateway completion time.
	ObservedAt time.Time

	// ReceivedAt is the gateway time immediately after the backend read ends.
	ReceivedAt time.Time

	// TimeOrigin identifies the source or gateway origin of ObservedAt.
	TimeOrigin string

	// ConfigFingerprint identifies the device and point configuration used by
	// this poll without exposing connection configuration contents.
	ConfigFingerprint string

	// Quality 資料品質標誌
	Quality schema.QualityFlag

	// Error 錯誤訊息
	Error string

	// QualityReason is a safe explanation for bad, missing or uncertain data.
	QualityReason string
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

	// BreakerConfig 熔斷器配置
	BreakerConfig health.BreakerConfig

	// Clock captures gateway acquisition completion time. Nil uses UTC now.
	Clock func() time.Time

	// AcquisitionIDFactory supplies one ID per device poll. Nil uses UUIDs.
	AcquisitionIDFactory func() string
}

// DefaultSchedulerConfig 預設排程器配置
func DefaultSchedulerConfig() SchedulerConfig {
	return SchedulerConfig{
		DefaultRetryCount:      3,
		DefaultRetryDelay:      1 * time.Second,
		MaxConcurrentPerDevice: 1, // 串列存取避免協議衝突
		ValueBufferSize:        1000,
		BreakerConfig:          health.DefaultBreakerConfig(),
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

	// 設備熔斷器
	deviceBreakers map[string]*health.CircuitBreaker

	// 狀態
	running bool
	mu      sync.RWMutex
	nowFunc func() time.Time

	acquisitionID func() string

	// 停止信號
	stopCh     chan struct{}
	stopDone   chan struct{}
	pollCtx    context.Context
	pollCancel context.CancelFunc
	wg         sync.WaitGroup
}

// groupTicker 輪詢群組 Ticker
type groupTicker struct {
	group  *schema.PollingGroup
	ticker *time.Ticker
	stopCh chan struct{}
	done   chan struct{}
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
	DataFormat     string
	Mode           schema.PointMode
	PollingGroupID string
}

// NewScheduler 建立新的排程器
func NewScheduler(config SchedulerConfig, connMgr *connector.ConnectionManager) *Scheduler {
	if connMgr == nil {
		connMgr = connector.GetConnectionManager()
	}

	return &Scheduler{
		config:         config,
		connMgr:        connMgr,
		groupTickers:   make(map[string]*groupTicker),
		valueChan:      make(chan CollectedValue, config.ValueBufferSize),
		deviceConfigs:  make(map[string]deviceConfig),
		pointInfos:     make(map[string]pointInfo),
		deviceLocks:    make(map[string]*sync.Mutex),
		deviceBreakers: make(map[string]*health.CircuitBreaker),
		stopCh:         make(chan struct{}),
		nowFunc:        config.Clock,
		acquisitionID:  config.AcquisitionIDFactory,
	}
}

func (s *Scheduler) now() time.Time {
	if s.nowFunc != nil {
		return s.nowFunc().UTC()
	}
	return time.Now().UTC()
}

func (s *Scheduler) nextAcquisitionID() string {
	if s.acquisitionID != nil {
		if id := s.acquisitionID(); id != "" {
			return id
		}
	}
	return uuid.NewString()
}
