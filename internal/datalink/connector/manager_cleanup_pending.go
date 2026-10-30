package connector

// cleanupEntry retains the exact generation and protocol handle until a Close
// succeeds. Only explicit Close/CloseAll retry failures; no background retry loop.
func (cm *ConnectionManager) cleanupEntry(entry *connectionEntry) error {
	entry.cleanupMu.Lock()
	defer entry.cleanupMu.Unlock()
	if entry.cleaned {
		return nil
	}
	err := closeManaged(entry.conn)
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if err != nil {
		cm.cleanup[entry] = entry.conn
		return errConnectionCleanup
	}
	entry.cleaned = true
	delete(cm.cleanup, entry)
	if cm.serialOwners[entry.serial] == entry {
		delete(cm.serialOwners, entry.serial)
	}
	return nil
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
