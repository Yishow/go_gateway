package connector

import "go-gateway/internal/datalink/schema"

// AcquisitionConnectionView is a read-only exact entry observation. Unavailable
// is deliberately distinct from absence and cannot authorize declared planning.
type AcquisitionConnectionView struct {
	State      string
	Token      AcquisitionConnectionToken
	Generation uint64
}

// PeekAcquisitionConnection never dials, waits for creation, or touches transport.
// Transport generation providers are required to read atomic state only.
func (cm *ConnectionManager) PeekAcquisitionConnection(device schema.Device) AcquisitionConnectionView {
	unavailable := AcquisitionConnectionView{State: "unavailable"}
	if cm == nil {
		return unavailable
	}
	identity, _, err := connectionIdentity(device.Protocol, device.ConnectionConfig)
	if err != nil {
		return unavailable
	}
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	for _, pending := range cm.cleanup {
		if pending.DeviceID == device.ID {
			return unavailable
		}
	}
	entry := cm.entries[device.ID]
	if entry == nil {
		if cm.connections[device.ID] != nil {
			return unavailable
		}
		return AcquisitionConnectionView{State: "absent", Generation: cm.versions[device.ID]}
	}
	if entry.identity != identity {
		return unavailable
	}
	select {
	case <-entry.done:
	default:
		return unavailable
	}
	if entry.err != nil || entry.conn == nil {
		return unavailable
	}
	token := entry.conn.acquisitionTokenLocked()
	if !token.Valid() || token.ManagerGeneration != entry.generation {
		return unavailable
	}
	return AcquisitionConnectionView{State: "current", Token: token, Generation: entry.generation}
}
