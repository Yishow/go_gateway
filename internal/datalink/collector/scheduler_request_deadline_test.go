package collector

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/wirebudget"
	"io"
	"sync/atomic"
	"testing"
	"time"
)

type requestDeadlineProtocol struct {
	connector.Protocol
	mode         string
	sends        atomic.Int32
	parentCancel context.CancelFunc
}

func (p *requestDeadlineProtocol) AcquisitionReadTimeout() time.Duration {
	return 15 * time.Millisecond
}
func (p *requestDeadlineProtocol) Read(ctx context.Context, _ connector.ReadRequest) (connector.ReadResult, error) {
	if err := wirebudget.Take(ctx); err != nil {
		return connector.ReadResult{}, err
	}
	p.sends.Add(1)
	if p.mode == "backoff" {
		return connector.ReadResult{Quality: schema.QualityBad}, io.EOF
	}
	<-ctx.Done()
	if p.parentCancel != nil {
		p.parentCancel()
	}
	if p.mode == "late-good" {
		return connector.ReadResult{Quality: schema.QualityGood, Value: 99}, nil
	}
	return connector.ReadResult{Quality: schema.QualityBad}, ctx.Err()
}
func TestConfiguredRequestDeadlineEmitsExactlyOneBad(t *testing.T) {
	for _, mode := range []string{"read", "backoff", "late-good"} {
		for _, retries := range []int{0, 3} {
			t.Run(mode+string(rune('0'+retries)), func(t *testing.T) {
				protocol := &requestDeadlineProtocol{mode: mode}
				settings := DefaultSchedulerConfig()
				settings.DefaultRetryCount = retries
				settings.DefaultRetryDelay = 30 * time.Millisecond
				scheduler := NewScheduler(settings, connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()))
				conn := &connector.ManagedConnection{DeviceID: "d", ProtocolType: schema.ProtocolModbusTCP, Protocol: protocol}
				point := pointInfo{ID: "p", DeviceID: "d", Address: "0", DataType: schema.DataTypeUint16}
				err := scheduler.pollPointWithAcquisition(t.Context(), conn, point, "scan")
				require.Error(t, err)
				if mode != "backoff" || retries > 0 {
					require.ErrorIs(t, err, context.DeadlineExceeded)
				}
				require.EqualValues(t, 1, protocol.sends.Load())
				var value CollectedValue
				select {
				case value = <-scheduler.ValueChannel():
				case <-time.After(150 * time.Millisecond):
					t.Fatal("missing final Bad observation")
				}
				require.Equal(t, schema.QualityBad, value.Quality)
				require.NotEmpty(t, value.Error)
				if mode == "late-good" {
					require.Nil(t, value.Value)
				}
				select {
				case <-scheduler.ValueChannel():
					t.Fatal("duplicate final value")
				default:
				}
			})
		}
	}
}
func TestParentCancellationSuppressesLateDeadlineObservation(t *testing.T) {
	for _, retries := range []int{0, 3} {
		t.Run(string(rune('0'+retries)), func(t *testing.T) {
			parent, cancel := context.WithCancel(t.Context())
			defer cancel()
			protocol := &requestDeadlineProtocol{mode: "late-good", parentCancel: cancel}
			settings := DefaultSchedulerConfig()
			settings.DefaultRetryCount = retries
			scheduler := NewScheduler(settings, connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()))
			conn := &connector.ManagedConnection{DeviceID: "d", ProtocolType: schema.ProtocolModbusTCP, Protocol: protocol}
			err := scheduler.pollPointWithAcquisition(parent, conn, pointInfo{ID: "p", Address: "0", DataType: schema.DataTypeUint16}, "scan")
			require.ErrorIs(t, err, context.Canceled)
			require.EqualValues(t, 1, protocol.sends.Load())
			select {
			case <-scheduler.ValueChannel():
				t.Fatal("parent cancellation published a late observation")
			default:
			}
		})
	}
}
