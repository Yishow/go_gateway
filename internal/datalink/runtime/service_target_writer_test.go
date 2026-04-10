package runtime

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"

	"github.com/stretchr/testify/require"
)

type closeAwareTargetWriter struct {
	closeCalls int
}

func (m *closeAwareTargetWriter) WriteTagValue(context.Context, string, any, time.Time) error {
	return nil
}

func (m *closeAwareTargetWriter) Close(context.Context) error {
	m.closeCalls++
	return nil
}

func TestService_Stop_ClosesTargetWriterWhenSupported(t *testing.T) {
	targetWriter := &closeAwareTargetWriter{}
	service := &Service{
		scheduler: collector.NewScheduler(collector.DefaultSchedulerConfig(), nil),
		writer:    &mockWriter{},
		target:    targetWriter,
		stopCh:    make(chan struct{}),
	}
	service.running.Store(true)

	require.NoError(t, service.Stop(context.Background()))
	require.Equal(t, 1, targetWriter.closeCalls)
}
