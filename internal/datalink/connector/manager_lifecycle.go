package connector

import (
	"context"
	"fmt"
	"time"

	"go-gateway/internal/datalink/schema"
)

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
