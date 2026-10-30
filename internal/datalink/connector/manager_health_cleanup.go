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

	if err := conn.mu.LockContext(ctx); err != nil {
		return err
	}
	defer conn.mu.Unlock()

	err := conn.Protocol.TestConnection(ctx)
	if err != nil {
		conn.LastError = err
	}
	return err
}

// CleanupIdle 清理閒置連線
func (cm *ConnectionManager) CleanupIdle() int {
	cm.mu.RLock()
	snapshot := make(map[string]*ManagedConnection, len(cm.connections))
	for id, conn := range cm.connections {
		snapshot[id] = conn
	}
	cm.mu.RUnlock()
	now := time.Now()
	cleaned := 0
	for id, conn := range snapshot {
		conn.mu.Lock()
		if conn.inUse || now.Sub(conn.LastUsed) <= cm.config.IdleTimeout {
			conn.mu.Unlock()
			continue
		}
		cm.mu.Lock()
		if cm.connections[id] != conn {
			cm.mu.Unlock()
			conn.mu.Unlock()
			continue
		}
		entry := cm.entries[id]
		delete(cm.connections, id)
		delete(cm.entries, id)
		cm.nextGeneration++
		cm.versions[id] = cm.nextGeneration
		cm.mu.Unlock()
		conn.mu.Unlock()
		if entry != nil {
			entry.cancel()
			if cm.cleanupEntry(entry) == nil {
				cleaned++
			}
		} else if closeManaged(conn) == nil {
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
