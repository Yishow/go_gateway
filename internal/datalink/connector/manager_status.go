package connector

import (
	"time"

	"go-gateway/internal/datalink/schema"
)

// ConnectionStatus 連線狀態
type ConnectionStatus struct {
	DeviceID       string              `json:"device_id"`
	ProtocolType   schema.ProtocolType `json:"protocol_type"`
	Connected      bool                `json:"connected"`
	InUse          bool                `json:"in_use"`
	LastUsed       time.Time           `json:"last_used"`
	CleanupPending bool                `json:"cleanup_pending,omitempty"`
	LastError      string              `json:"last_error,omitempty"`
}

// GetStatus 取得設備連線狀態
func (cm *ConnectionManager) GetStatus(deviceID string) (ConnectionStatus, bool) {
	cm.mu.RLock()
	conn, ok := cm.connections[deviceID]
	pending := cm.pendingCleanup(deviceID, false)
	if !ok && len(pending) > 0 {
		conn = cm.cleanup[pending[0]]
		ok = true
	}
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

	if len(pending) > 0 {
		status.CleanupPending = true
		status.LastError = errConnectionCleanup.Error()
	}
	return status, true
}

// ListConnections 列出所有連線狀態
func (cm *ConnectionManager) ListConnections() []ConnectionStatus {
	cm.mu.RLock()
	ids := make(map[string]bool)
	for id := range cm.connections {
		ids[id] = true
	}
	for _, conn := range cm.cleanup {
		ids[conn.DeviceID] = true
	}
	cm.mu.RUnlock()
	statuses := make([]ConnectionStatus, 0, len(ids))
	for id := range ids {
		if status, ok := cm.GetStatus(id); ok {
			statuses = append(statuses, status)
		}
	}
	return statuses
}
