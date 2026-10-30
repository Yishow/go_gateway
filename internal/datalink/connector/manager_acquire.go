package connector

import (
	"context"
	"errors"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/wirebudget"
)

type connectionEntry struct {
	identity   [32]byte
	generation uint64
	serial     string
	done       chan struct{}
	cancel     context.CancelFunc
	conn       *ManagedConnection
	err        error
	cleanupMu  wirebudget.Mutex
	cleaned    bool
}

// GetOrCreate shares one creation for an immutable protocol/configuration generation.
// No protocol operation, dial, or per-connection lock runs under the manager lock.
func (cm *ConnectionManager) GetOrCreate(ctx context.Context, deviceID string, typ schema.ProtocolType, config string) (*ManagedConnection, error) {
	return cm.getOrCreate(ctx, deviceID, typ, config, nil, false)
}

// expected fences reconnect retries; callers of the public API must order snapshots upstream.
func (cm *ConnectionManager) getOrCreate(ctx context.Context, deviceID string, typ schema.ProtocolType, config string, expected *uint64, force bool) (*ManagedConnection, error) {
	identity, serial, err := connectionIdentity(typ, config)
	if err != nil {
		return nil, err
	}
	if err = ctx.Err(); err != nil {
		return nil, err
	}
	cm.mu.Lock()
	if expected != nil && cm.versions[deviceID] != *expected {
		cm.mu.Unlock()
		return nil, errConnectionSuperseded
	}
	for _, pending := range cm.cleanup {
		if pending.DeviceID == deviceID {
			cm.mu.Unlock()
			return nil, errConnectionCleanup
		}
	}
	old := cm.entries[deviceID]
	if !force && old != nil && old.identity == identity {
		cm.mu.Unlock()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-old.done:
		}
		cm.mu.RLock()
		current := cm.entries[deviceID] == old
		cm.mu.RUnlock()
		if !current {
			return nil, errConnectionSuperseded
		}
		if old.err != nil {
			return nil, old.err
		}
		if err := old.conn.mu.LockContext(ctx); err != nil {
			return nil, err
		}
		live, statusErr := protocolConnectedContext(ctx, old.conn.Protocol)
		old.conn.LastUsed = time.Now()
		old.conn.mu.Unlock()
		if statusErr != nil {
			return nil, statusErr
		}
		cm.mu.RLock()
		current = cm.entries[deviceID] == old
		cm.mu.RUnlock()
		if !current {
			return nil, errConnectionSuperseded
		}
		if live {
			return old.conn, nil
		}
		version, retired, retireErr := cm.retireIfCurrentContext(ctx, deviceID, old)
		if retireErr != nil {
			return nil, retireErr
		}
		if !retired {
			return nil, errConnectionSuperseded
		}
		return cm.getOrCreate(ctx, deviceID, typ, config, &version, false)
	}
	if owner := cm.serialOwners[serial]; serial != "" && owner != nil {
		ready := false
		if owner == old {
			select {
			case <-old.done:
				ready = true
			default:
			}
		}
		if !ready {
			cm.mu.Unlock()
			return nil, errSerialAlias
		}
	}
	cm.nextGeneration++
	cm.versions[deviceID] = cm.nextGeneration
	if expected != nil {
		*expected = cm.nextGeneration
	}
	dialCtx, cancel := context.WithCancel(ctx)
	entry := &connectionEntry{identity: identity, generation: cm.nextGeneration, serial: serial, done: make(chan struct{}), cancel: cancel}
	cm.entries[deviceID] = entry
	delete(cm.connections, deviceID)
	if serial != "" {
		cm.serialOwners[serial] = entry
	}
	cm.mu.Unlock()
	var conn *ManagedConnection
	if old != nil {
		err = cm.retireContext(ctx, old)
		if err != nil {
			cm.mu.Lock()
			if old.serial != "" && cm.serialOwners[old.serial] == entry {
				cm.serialOwners[old.serial] = old
			}
			cm.mu.Unlock()
		}
	}
	if err == nil {
		conn, err = cm.createConnection(dialCtx, deviceID, typ, config)
	}
	entry.conn = conn

	if err == nil {
		err = dialCtx.Err()
	}
	cm.mu.Lock()
	if cm.entries[deviceID] != entry {
		err = errConnectionSuperseded
	}
	if err == nil {
		entry.conn = conn
		cm.connections[deviceID] = conn
	} else if cm.entries[deviceID] == entry {
		delete(cm.entries, deviceID)
	}
	entry.err = err
	if err == nil {
		close(entry.done)
	}
	cm.mu.Unlock()
	if err != nil {
		if conn != nil {
			if cm.cleanupEntry(entry) != nil {
				err = errConnectionCleanup
			}
		} else {
			cm.releaseSerial(entry)
		}
	}

	cancel()
	if err != nil {
		entry.err = err
		close(entry.done)
	}
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (cm *ConnectionManager) createConnection(ctx context.Context, deviceID string, typ schema.ProtocolType, config string) (*ManagedConnection, error) {
	connCtx, cancel := context.WithTimeout(ctx, cm.config.ConnectionTimeout)
	defer cancel()
	factory, ok := GetFactory(typ)
	if !ok {
		return nil, errors.New("建立連線失敗: connection protocol not registered")
	}
	protocol := factory()
	conn := &ManagedConnection{DeviceID: deviceID, Protocol: protocol, Config: config, ProtocolType: typ, LastUsed: time.Now()}
	if err := protocol.Connect(connCtx, config); err != nil {
		if connCtx.Err() != nil {
			return conn, connCtx.Err()
		}
		return conn, errors.New("建立連線失敗: connection establishment failed")
	}
	if err := connCtx.Err(); err != nil {
		return conn, err
	}
	return conn, nil
}

func closeManaged(conn *ManagedConnection) error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.Protocol.Close() != nil {
		return errConnectionCleanup
	}
	return nil
}
func (cm *ConnectionManager) releaseSerial(entry *connectionEntry) {
	cm.mu.Lock()
	if cm.serialOwners[entry.serial] == entry {
		delete(cm.serialOwners, entry.serial)
	}
	cm.mu.Unlock()
}
func (cm *ConnectionManager) retire(entry *connectionEntry) error {
	return cm.retireContext(context.Background(), entry)
}
func (cm *ConnectionManager) retireIfCurrentContext(ctx context.Context, id string, entry *connectionEntry) (uint64, bool, error) {
	cm.mu.Lock()
	if cm.entries[id] != entry {
		cm.mu.Unlock()
		return 0, false, nil
	}
	delete(cm.entries, id)
	cm.nextGeneration++
	cm.versions[id] = cm.nextGeneration
	version := cm.nextGeneration
	delete(cm.connections, id)
	cm.mu.Unlock()
	err := cm.retireContext(ctx, entry)
	return version, true, err
}
