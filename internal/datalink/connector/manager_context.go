package connector

import (
	"context"
	"errors"
)

func protocolConnectedContext(ctx context.Context, protocol Protocol) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if checker, ok := protocol.(interface {
		IsConnectedContext(context.Context) (bool, error)
	}); ok {
		return checker.IsConnectedContext(ctx)
	}
	return protocol.IsConnected(), ctx.Err()
}

func (cm *ConnectionManager) retireContext(ctx context.Context, entry *connectionEntry) error {
	entry.cancel()
	select {
	case <-entry.done:
		if entry.conn != nil {
			return cm.cleanupEntryContext(ctx, entry)
		}
	default:
	}
	return nil
}

func (cm *ConnectionManager) cleanupEntryContext(ctx context.Context, entry *connectionEntry) error {
	if err := entry.cleanupMu.LockContext(ctx); err != nil {
		// Do not orphan a retired generation if cancellation happened before
		// cleanup could take its own lock. An active cleaner already owns it.
		if entry.cleanupMu.TryLock() {
			if !entry.cleaned {
				cm.mu.Lock()
				cm.cleanup[entry] = entry.conn
				cm.mu.Unlock()
			}
			entry.cleanupMu.Unlock()
		}
		return errors.Join(errConnectionCleanup, err)
	}
	defer entry.cleanupMu.Unlock()
	if entry.cleaned {
		return nil
	}
	err := closeManagedContext(ctx, entry.conn)
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if err != nil {
		cm.cleanup[entry] = entry.conn
		if ctx.Err() != nil {
			return errors.Join(errConnectionCleanup, ctx.Err())
		}
		return errConnectionCleanup
	}
	entry.cleaned = true
	delete(cm.cleanup, entry)
	if cm.serialOwners[entry.serial] == entry {
		delete(cm.serialOwners, entry.serial)
	}
	return nil
}

func closeManagedContext(ctx context.Context, conn *ManagedConnection) error {
	if err := conn.mu.LockContext(ctx); err != nil {
		return err
	}
	defer conn.mu.Unlock()
	if closer, ok := conn.Protocol.(interface{ CloseContext(context.Context) error }); ok {
		return closer.CloseContext(ctx)
	}
	return conn.Protocol.Close()
}
