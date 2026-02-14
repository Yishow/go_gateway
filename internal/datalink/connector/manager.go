package connector

import (
	"sync"
	"time"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 連線管理器
// =============================================================================

// ConnectionManager 管理協議連線的生命週期和連線池
type ConnectionManager struct {
	mu          sync.RWMutex
	connections map[string]*ManagedConnection
	config      ConnectionManagerConfig
}

// ConnectionManagerConfig 連線管理器配置
type ConnectionManagerConfig struct {
	// MaxConnectionsPerDevice 每設備最大連線數
	MaxConnectionsPerDevice int

	// ConnectionTimeout 連線逾時時間
	ConnectionTimeout time.Duration

	// IdleTimeout 閒置連線逾時時間
	IdleTimeout time.Duration

	// HealthCheckInterval 健康檢查間隔
	HealthCheckInterval time.Duration

	// RetryCount 重連重試次數
	RetryCount int

	// RetryDelay 重連延遲
	RetryDelay time.Duration
}

// DefaultConnectionManagerConfig 預設連線管理器配置
func DefaultConnectionManagerConfig() ConnectionManagerConfig {
	return ConnectionManagerConfig{
		MaxConnectionsPerDevice: 1,
		ConnectionTimeout:       10 * time.Second,
		IdleTimeout:             5 * time.Minute,
		HealthCheckInterval:     30 * time.Second,
		RetryCount:              3,
		RetryDelay:              1 * time.Second,
	}
}

// ManagedConnection 受管理的連線
type ManagedConnection struct {
	// DeviceID 設備 ID
	DeviceID string

	// Protocol 協議連接器
	Protocol Protocol

	// Config 連線配置 (JSON)
	Config string

	// ProtocolType 協議類型
	ProtocolType schema.ProtocolType

	// LastUsed 最後使用時間
	LastUsed time.Time

	// LastError 最後錯誤
	LastError error

	// mu 連線鎖定 (用於序列化存取)
	mu sync.Mutex

	// inUse 是否正在使用中
	inUse bool
}

// NewConnectionManager 建立新的連線管理器
func NewConnectionManager(config ConnectionManagerConfig) *ConnectionManager {
	cm := &ConnectionManager{
		connections: make(map[string]*ManagedConnection),
		config:      config,
	}
	return cm
}

// 全域連線管理器
var globalConnectionManager *ConnectionManager
var cmOnce sync.Once

// GetConnectionManager 取得全域連線管理器
func GetConnectionManager() *ConnectionManager {
	cmOnce.Do(func() {
		globalConnectionManager = NewConnectionManager(DefaultConnectionManagerConfig())
	})
	return globalConnectionManager
}
