package recordingplan

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMySQLDDLHeartbeatCancelsExecutorOnHeartbeatFailure(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "lost ownership", err: ErrSchemaOperationNotOwned},
		{name: "heartbeat error", err: errors.New("progress write failed")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executorStarted := make(chan struct{})
			executorCanceled := make(chan struct{})
			heartbeat := func(context.Context) error { return tt.err }
			executor := func(ctx context.Context) error {
				close(executorStarted)
				<-ctx.Done()
				close(executorCanceled)
				return ctx.Err()
			}

			result := make(chan struct {
				execErr      error
				heartbeatErr error
			}, 1)
			go func() {
				execErr, heartbeatErr := runMySQLDDLWithHeartbeat(t.Context(), 5*time.Millisecond, heartbeat, executor)
				result <- struct {
					execErr      error
					heartbeatErr error
				}{execErr, heartbeatErr}
			}()
			<-executorStarted

			select {
			case <-executorCanceled:
			case <-time.After(time.Second):
				t.Fatal("heartbeat failure did not cancel the blocked DDL executor")
			}
			got := <-result
			require.ErrorIs(t, got.execErr, context.Canceled)
			require.ErrorIs(t, got.heartbeatErr, tt.err)
		})
	}
}
