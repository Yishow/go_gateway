package collector

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/wirebudget"

	"github.com/stretchr/testify/require"
)

type boundedBackoffProtocol struct {
	scheduledExecutionProtocol
	adapterBackoff bool
	aSends, bSends *atomic.Int32
}

func (p *boundedBackoffProtocol) Connect(ctx context.Context, config string) error {
	if err := p.scheduledExecutionProtocol.Connect(ctx, config); err != nil {
		return err
	}
	isA := strings.Contains(config, "device-a")
	p.hook = func(ctx context.Context, req connector.PlannedReadRequest, _ int) ([]connector.PlannedPointResult, error) {
		if !isA {
			p.bSends.Add(1)
			return executionResults(req, time.Now(), nil), nil
		}
		p.aSends.Add(1)
		err := errors.New("retryable")
		if p.adapterBackoff {
			if err = wirebudget.Wait(ctx); err == nil {
				err = wirebudget.Take(ctx)
				if err == nil {
					p.aSends.Add(1)
				}
			}
		}
		return executionResults(req, time.Now(), err), err
	}
	return nil
}

func TestBoundedSchedulerBackoffYieldsButKeepsDeviceGate(t *testing.T) {
	for _, adapter := range []bool{false, true} {
		t.Run(map[bool]string{false: "scheduler", true: "adapter"}[adapter], func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var aSends, bSends atomic.Int32
				previous, existed := connector.GetFactory(schema.ProtocolMC3E)
				connector.Register(schema.ProtocolMC3E, func() connector.Protocol {
					return &boundedBackoffProtocol{scheduledExecutionProtocol: scheduledExecutionProtocol{timeout: time.Second}, adapterBackoff: adapter, aSends: &aSends, bSends: &bSends}
				})
				defer func() {
					if existed {
						connector.Register(schema.ProtocolMC3E, previous)
					} else {
						connector.Unregister(schema.ProtocolMC3E)
					}
				}()
				manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
				defer manager.CloseAll()
				config := DefaultSchedulerConfig()
				one := 1
				config.MaxActiveScans = &one
				config.DefaultRetryCount = 2
				config.DefaultRetryDelay = 100 * time.Millisecond
				s := NewScheduler(config, manager)
				for _, id := range []string{"a", "b"} {
					group := id
					s.AddDevice(&schema.Device{ID: id, Protocol: schema.ProtocolMC3E, ConnectionConfig: `{"host":"device-` + id + `"}`})
					s.AddPoint(&schema.Point{ID: id, DeviceID: id, Address: "D0", DataType: schema.DataTypeInt16, PollingGroupID: &group})
				}
				require.NoError(t, s.Start([]*schema.PollingGroup{{ID: "a", Enabled: true, IntervalMs: 10}, {ID: "b", Enabled: true, IntervalMs: 20}}))
				defer func() { require.NoError(t, s.Stop()) }()
				time.Sleep(35 * time.Millisecond)
				synctest.Wait()
				require.EqualValues(t, 1, aSends.Load())
				require.GreaterOrEqual(t, bSends.Load(), int32(1), "backoff must free capacity one for B")
				gateCtx, cancel := context.WithCancel(t.Context())
				gateEntered := make(chan error, 1)
				go func() {
					err := s.deviceLocks["a"].LockContext(gateCtx)
					if err == nil {
						s.deviceLocks["a"].Unlock()
					}
					gateEntered <- err
				}()
				synctest.Wait()
				select {
				case err := <-gateEntered:
					t.Fatalf("A lost scan ownership during backoff: %v", err)
				default:
				}
				cancel()
				require.ErrorIs(t, <-gateEntered, context.Canceled)
				require.NoError(t, s.Stop())
				before := aSends.Load() + bSends.Load()
				time.Sleep(time.Second)
				require.Equal(t, before, aSends.Load()+bSends.Load(), "stop cannot resume waiting retries or pending ticks")
			})
		})
	}
}
