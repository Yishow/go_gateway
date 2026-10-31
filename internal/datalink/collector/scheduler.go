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
	"go-gateway/internal/protocol/wirebudget"

	"github.com/google/uuid"
)

// =============================================================================
// 收集結果
// =============================================================================

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

	// AcquisitionObserver must be nonblocking; it receives final metadata only.
	// It is immutable after construction and cannot change acquisition results.
	AcquisitionObserver func(AcquisitionObservation)
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

	// Monotonic for this scheduler lifetime, including removed/re-added points.
	pointObservationLease uint64
	installationScope     string

	// 設備鎖 (確保同一設備串列存取)
	deviceLocks map[string]*wirebudget.Mutex

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
	readProfile        connector.ReadProfile
	configRevision     string
	ID                 string
	DeviceID           string
	Address            string
	Function           string
	DataType           schema.DataType
	DataFormat         string
	Mode               schema.PointMode
	PollingGroupID     string
	observationLease   uint64
	installationScope  string
	observationBinding PointObservationBinding
}

// NewScheduler 建立新的排程器
func NewScheduler(config SchedulerConfig, connMgr *connector.ConnectionManager) *Scheduler {
	if connMgr == nil {
		connMgr = connector.GetConnectionManager()
	}

	return &Scheduler{
		installationScope: uuid.NewString(),
		config:            config,
		connMgr:           connMgr,
		groupTickers:      make(map[string]*groupTicker),
		valueChan:         make(chan CollectedValue, config.ValueBufferSize),
		deviceConfigs:     make(map[string]deviceConfig),
		pointInfos:        make(map[string]pointInfo),
		deviceLocks:       make(map[string]*wirebudget.Mutex),
		deviceBreakers:    make(map[string]*health.CircuitBreaker),
		stopCh:            make(chan struct{}),
		nowFunc:           config.Clock,
		acquisitionID:     config.AcquisitionIDFactory,
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
