package collector

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type deviceGateProtocol struct {
	connector.Protocol
	connects, reads atomic.Int32
}

func (p *deviceGateProtocol) Connect(context.Context, string) error { p.connects.Add(1); return nil }
func (p *deviceGateProtocol) Read(context.Context, connector.ReadRequest) (connector.ReadResult, error) {
	p.reads.Add(1)
	return connector.ReadResult{Quality: schema.QualityGood}, nil
}
func (p *deviceGateProtocol) Close() error      { return nil }
func (p *deviceGateProtocol) IsConnected() bool { return true }
func TestSchedulerDeviceGateHonorsCanceledWaiters(t *testing.T) {
	for _, manual := range []bool{false, true} {
		t.Run(map[bool]string{false: "scheduled", true: "manual"}[manual], func(t *testing.T) {
			typ := schema.ProtocolType("device-gate-cancel")
			protocol := &deviceGateProtocol{}
			connector.Register(typ, func() connector.Protocol { return protocol })
			defer connector.Unregister(typ)
			manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
			defer manager.CloseAll()
			scheduler := NewScheduler(DefaultSchedulerConfig(), manager)
			scheduler.AddDevice(&schema.Device{ID: "d", Protocol: typ, ConnectionConfig: "{}"})
			scheduler.AddPoint(&schema.Point{ID: "p", DeviceID: "d", Address: "0", DataType: schema.DataTypeUint16})
			gate := scheduler.deviceLocks["d"]
			gate.Lock()
			release := sync.OnceFunc(gate.Unlock)
			defer release()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
			defer cancel()
			scheduler.pollCtx = ctx
			done := make(chan []CollectedValue, 1)
			go func() {
				if manual {
					done <- scheduler.PollNowContext(ctx, []string{"p"})
					return
				}
				scheduler.pollDevicePoints("d", []pointInfo{scheduler.pointInfos["p"]})
				done <- nil
			}()
			var results []CollectedValue
			select {
			case results = <-done:
			case <-time.After(150 * time.Millisecond):
				t.Error("device gate ignored caller cancellation")
				release()
				results = <-done
			}
			require.Zero(t, protocol.connects.Load())
			require.Zero(t, protocol.reads.Load())
			if manual {
				require.Len(t, results, 1)
				require.Equal(t, schema.QualityBad, results[0].Quality)
				require.Contains(t, results[0].Error, context.DeadlineExceeded.Error())
			}
		})
	}
}
