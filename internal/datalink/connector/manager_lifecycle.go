package connector

import (
	"context"
	"errors"
	"fmt"
	"time"
)

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

// Close invalidates the current generation before closing its resources.
func (cm *ConnectionManager) Close(deviceID string) error {
	cm.mu.Lock()
	entry := cm.entries[deviceID]
	pending := cm.pendingCleanup(deviceID, false)
	cm.nextGeneration++
	cm.versions[deviceID] = cm.nextGeneration
	delete(cm.entries, deviceID)
	delete(cm.connections, deviceID)
	cm.mu.Unlock()
	var lastErr error
	if entry != nil {
		lastErr = cm.retire(entry)
	}
	for _, old := range pending {
		if err := cm.cleanupEntry(old); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// CloseAll invalidates all current generations, including pending dials.
func (cm *ConnectionManager) CloseAll() error {
	cm.mu.Lock()
	entries := cm.entries
	pending := cm.pendingCleanup("", true)
	for id := range cm.versions {
		cm.nextGeneration++
		cm.versions[id] = cm.nextGeneration
	}
	cm.entries = make(map[string]*connectionEntry)
	cm.connections = make(map[string]*ManagedConnection)
	cm.mu.Unlock()
	var lastErr error
	for _, entry := range entries {
		if err := cm.retire(entry); err != nil {
			lastErr = err
		}
	}
	for _, old := range pending {
		if err := cm.cleanupEntry(old); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Reconnect replaces only the observed generation; concurrent configuration changes win.
func (cm *ConnectionManager) Reconnect(ctx context.Context, deviceID string) error {
	cm.mu.RLock()
	entry := cm.entries[deviceID]
	conn := cm.connections[deviceID]
	cm.mu.RUnlock()
	if conn == nil {
		return fmt.Errorf("device has no existing connection")
	}
	expected := entry.generation
	var lastErr error
	for i := 0; i <= cm.config.RetryCount; i++ {
		if i > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(cm.config.RetryDelay):
			}
		}
		_, lastErr = cm.getOrCreate(ctx, deviceID, conn.ProtocolType, conn.Config, &expected, true)
		if errors.Is(lastErr, errConnectionSuperseded) {
			return lastErr
		}
		if lastErr == nil {
			return nil
		}
	}
	return fmt.Errorf("reconnection failed: %w", lastErr)
}
