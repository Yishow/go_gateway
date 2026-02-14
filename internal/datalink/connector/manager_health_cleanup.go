package connector

import (
	"context"
	"fmt"
	"time"
)

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
