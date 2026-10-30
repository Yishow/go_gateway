package connector

import "context"

// cleanupEntry retains the exact generation and protocol handle until a Close
// succeeds. Only explicit Close/CloseAll retry failures; no background retry loop.
func (cm *ConnectionManager) cleanupEntry(entry *connectionEntry) error {
	return cm.cleanupEntryContext(context.Background(), entry)
}

func (cm *ConnectionManager) pendingCleanup(deviceID string, all bool) []*connectionEntry {
	var pending []*connectionEntry
	for entry, conn := range cm.cleanup {
		if all || conn.DeviceID == deviceID {
			pending = append(pending, entry)
		}
	}
	return pending
}
