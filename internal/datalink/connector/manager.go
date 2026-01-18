package connector

import (
	"context"
	"fmt"
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

// =============================================================================
// 連線管理方法
// =============================================================================

// GetOrCreate 取得或建立設備連線
//
// 參數:
//   - ctx: 上下文
//   - deviceID: 設備 ID
//   - protocolType: 協議類型
//   - config: JSON 格式的連線配置
//
// 返回:
//   - *ManagedConnection: 受管理的連線
//   - error: 錯誤 (若有)
func (cm *ConnectionManager) GetOrCreate(
	ctx context.Context,
	deviceID string,
	protocolType schema.ProtocolType,
	config string,
) (*ManagedConnection, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 檢查是否已存在連線
	if conn, ok := cm.connections[deviceID]; ok {
		// 驗證連線仍然有效
		if conn.Protocol.IsConnected() {
			conn.LastUsed = time.Now()
			return conn, nil
		}
		// 連線已斷開，移除舊連線
		_ = conn.Protocol.Close()
		delete(cm.connections, deviceID)
	}

	// 建立新連線
	conn, err := cm.createConnection(ctx, deviceID, protocolType, config)
	if err != nil {
		return nil, err
	}

	cm.connections[deviceID] = conn
	return conn, nil
}

// createConnection 建立新連線 (內部使用)
func (cm *ConnectionManager) createConnection(
	ctx context.Context,
	deviceID string,
	protocolType schema.ProtocolType,
	config string,
) (*ManagedConnection, error) {
	// 設定連線逾時
	connCtx, cancel := context.WithTimeout(ctx, cm.config.ConnectionTimeout)
	defer cancel()

	// 從註冊表取得連接器
	protocol, err := Get(connCtx, protocolType, config)
	if err != nil {
		return nil, fmt.Errorf("建立連線失敗 (設備 %s): %w", deviceID, err)
	}

	return &ManagedConnection{
		DeviceID:     deviceID,
		Protocol:     protocol,
		Config:       config,
		ProtocolType: protocolType,
		LastUsed:     time.Now(),
	}, nil
}

// Release 釋放連線 (標記為可用，不關閉)
func (cm *ConnectionManager) Release(deviceID string) {
	cm.mu.RLock()
	conn, ok := cm.connections[deviceID]
	cm.mu.RUnlock()

	if ok {
		conn.mu.Lock()
		conn.inUse = false
		conn.LastUsed = time.Now()
		conn.mu.Unlock()
	}
}

// Close 關閉指定設備的連線
func (cm *ConnectionManager) Close(deviceID string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	conn, ok := cm.connections[deviceID]
	if !ok {
		return nil
	}

	err := conn.Protocol.Close()
	delete(cm.connections, deviceID)
	return err
}

// CloseAll 關閉所有連線
func (cm *ConnectionManager) CloseAll() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	var lastErr error
	for deviceID, conn := range cm.connections {
		if err := conn.Protocol.Close(); err != nil {
			lastErr = err
		}
		delete(cm.connections, deviceID)
	}
	return lastErr
}

// Reconnect 重新連線指定設備
func (cm *ConnectionManager) Reconnect(ctx context.Context, deviceID string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	conn, ok := cm.connections[deviceID]
	if !ok {
		return fmt.Errorf("設備 %s 沒有現有連線", deviceID)
	}

	// 關閉舊連線
	_ = conn.Protocol.Close()

	// 重建連線 (帶重試)
	var lastErr error
	for i := 0; i <= cm.config.RetryCount; i++ {
		if i > 0 {
			time.Sleep(cm.config.RetryDelay)
		}

		newConn, err := cm.createConnection(ctx, deviceID, conn.ProtocolType, conn.Config)
		if err != nil {
			lastErr = err
			continue
		}

		cm.connections[deviceID] = newConn
		return nil
	}

	// 所有重試都失敗，移除連線
	delete(cm.connections, deviceID)
	return fmt.Errorf("重連失敗 (設備 %s): %w", deviceID, lastErr)
}

// =============================================================================
// 健康檢查和清理
// =============================================================================

// HealthCheck 對指定設備執行健康檢查
func (cm *ConnectionManager) HealthCheck(ctx context.Context, deviceID string) error {
	cm.mu.RLock()
	conn, ok := cm.connections[deviceID]
	cm.mu.RUnlock()

	if !ok {
		return fmt.Errorf("設備 %s 沒有連線", deviceID)
	}

	conn.mu.Lock()
	defer conn.mu.Unlock()

	err := conn.Protocol.TestConnection(ctx)
	if err != nil {
		conn.LastError = err
	}
	return err
}

// CleanupIdle 清理閒置連線
func (cm *ConnectionManager) CleanupIdle() int {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	now := time.Now()
	cleaned := 0

	for deviceID, conn := range cm.connections {
		conn.mu.Lock()
		idle := now.Sub(conn.LastUsed)
		inUse := conn.inUse
		conn.mu.Unlock()

		if !inUse && idle > cm.config.IdleTimeout {
			_ = conn.Protocol.Close()
			delete(cm.connections, deviceID)
			cleaned++
		}
	}

	return cleaned
}

// StartCleanupRoutine 啟動定期清理 goroutine
func (cm *ConnectionManager) StartCleanupRoutine(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(cm.config.IdleTimeout / 2)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				cm.CleanupIdle()
			}
		}
	}()
}

// =============================================================================
// 連線狀態查詢
// =============================================================================

// ConnectionStatus 連線狀態
type ConnectionStatus struct {
	DeviceID     string             `json:"device_id"`
	ProtocolType schema.ProtocolType `json:"protocol_type"`
	Connected    bool               `json:"connected"`
	InUse        bool               `json:"in_use"`
	LastUsed     time.Time          `json:"last_used"`
	LastError    string             `json:"last_error,omitempty"`
}

// GetStatus 取得設備連線狀態
func (cm *ConnectionManager) GetStatus(deviceID string) (ConnectionStatus, bool) {
	cm.mu.RLock()
	conn, ok := cm.connections[deviceID]
	cm.mu.RUnlock()

	if !ok {
		return ConnectionStatus{}, false
	}

	conn.mu.Lock()
	defer conn.mu.Unlock()

	status := ConnectionStatus{
		DeviceID:     conn.DeviceID,
		ProtocolType: conn.ProtocolType,
		Connected:    conn.Protocol.IsConnected(),
		InUse:        conn.inUse,
		LastUsed:     conn.LastUsed,
	}

	if conn.LastError != nil {
		status.LastError = conn.LastError.Error()
	}

	return status, true
}

// ListConnections 列出所有連線狀態
func (cm *ConnectionManager) ListConnections() []ConnectionStatus {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	statuses := make([]ConnectionStatus, 0, len(cm.connections))
	for _, conn := range cm.connections {
		conn.mu.Lock()
		status := ConnectionStatus{
			DeviceID:     conn.DeviceID,
			ProtocolType: conn.ProtocolType,
			Connected:    conn.Protocol.IsConnected(),
			InUse:        conn.inUse,
			LastUsed:     conn.LastUsed,
		}
		if conn.LastError != nil {
			status.LastError = conn.LastError.Error()
		}
		conn.mu.Unlock()
		statuses = append(statuses, status)
	}

	return statuses
}

// =============================================================================
// ManagedConnection 方法
// =============================================================================

// Lock 鎖定連線以進行操作
func (mc *ManagedConnection) Lock() {
	mc.mu.Lock()
	mc.inUse = true
}

// Unlock 解鎖連線
func (mc *ManagedConnection) Unlock() {
	mc.inUse = false
	mc.LastUsed = time.Now()
	mc.mu.Unlock()
}

// Read 透過連線讀取資料 (自動鎖定)
func (mc *ManagedConnection) Read(ctx context.Context, req ReadRequest) (ReadResult, error) {
	mc.Lock()
	defer mc.Unlock()

	result, err := mc.Protocol.Read(ctx, req)
	if err != nil {
		mc.LastError = err
	}
	return result, err
}

// Write 透過連線寫入資料 (自動鎖定)
func (mc *ManagedConnection) Write(ctx context.Context, req WriteRequest) error {
	mc.Lock()
	defer mc.Unlock()

	err := mc.Protocol.Write(ctx, req)
	if err != nil {
		mc.LastError = err
	}
	return err
}
