package recordingplan

import (
	"context"
	"errors"
	"time"
)

// runMySQLDDLWithHeartbeat keeps an active schema operation's lease alive
// while a single DDL statement is in flight. A failed heartbeat cancels the
// executor context; callers must still inspect the target catalog because the
// DDL may have committed before the failure was observed.
func runMySQLDDLWithHeartbeat(
	ctx context.Context,
	interval time.Duration,
	heartbeat func(context.Context) error,
	execute func(context.Context) error,
) (execErr, heartbeatErr error) {
	if interval <= 0 {
		return nil, errors.New("MySQL DDL heartbeat interval must be positive")
	}
	if heartbeat == nil {
		return nil, errors.New("MySQL DDL heartbeat callback is required")
	}
	if execute == nil {
		return nil, errors.New("MySQL DDL executor is required")
	}

	heartbeatCtx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	heartbeatDone := make(chan error, 1)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-heartbeatCtx.Done():
				heartbeatDone <- nil
				return
			case <-ticker.C:
				err := heartbeat(heartbeatCtx)
				if err == nil {
					continue
				}
				// A normal executor completion can race with a progress write
				// that observes the cancellation. That is not a heartbeat
				// failure; only report errors from an in-flight active lease.
				select {
				case <-heartbeatCtx.Done():
					heartbeatDone <- nil
					return
				default:
				}
				cancel(err)
				heartbeatDone <- err
				return
			}
		}
	}()

	execErr = execute(heartbeatCtx)
	cancel(nil)
	heartbeatErr = <-heartbeatDone
	return execErr, heartbeatErr
}
