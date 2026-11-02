package device

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type cleanupContractProtocol struct {
	connector.Protocol
	closed     int
	cleanupCtx context.Context
}

func (*cleanupContractProtocol) Connect(context.Context, string) error { return nil }
func (p *cleanupContractProtocol) IsConnected() bool                   { return p.closed == 0 }
func (p *cleanupContractProtocol) Close() error                        { p.closed++; return nil }
func (p *cleanupContractProtocol) CloseContext(ctx context.Context) error {
	p.cleanupCtx = ctx
	return p.Close()
}

type cleanupContractRepository struct {
	Repository
	cancel   context.CancelFunc
	observed context.Context
}

func (r *cleanupContractRepository) GetByID(ctx context.Context, id string) (*schema.Device, error) {
	value, err := r.Repository.GetByID(ctx, id)
	if r.cancel != nil {
		r.cancel()
	}
	return value, err
}
func (r *cleanupContractRepository) ClearTestResult(ctx context.Context, _ string) error {
	r.observed = ctx
	return ctx.Err()
}
func (r *cleanupContractRepository) Delete(ctx context.Context, _ string) error {
	r.observed = ctx
	return ctx.Err()
}
func (r *cleanupContractRepository) UpdateStatus(ctx context.Context, _ string, _ schema.DeviceStatus) error {
	r.observed = ctx
	return ctx.Err()
}

func TestDeviceMutationCleanupDoesNotInheritRequestCancellation(t *testing.T) {
	for _, operation := range []string{"update", "delete", "disable"} {
		t.Run(operation, func(t *testing.T) {
			typ := schema.ProtocolType(t.Name())
			var protocols []*cleanupContractProtocol
			connector.Register(typ, func() connector.Protocol {
				p := &cleanupContractProtocol{}
				protocols = append(protocols, p)
				return p
			})
			t.Cleanup(func() { connector.Unregister(typ) })
			manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
			t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
			_, err := manager.GetOrCreate(t.Context(), "target", typ, `{}`)
			require.NoError(t, err)
			_, err = manager.GetOrCreate(t.Context(), "other", typ, `{}`)
			require.NoError(t, err)
			require.Len(t, protocols, 2)
			memory := NewMemoryRepository()
			require.NoError(t, memory.Create(t.Context(), &schema.Device{
				ID: "target", Protocol: schema.ProtocolModbusTCP,
				ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":5}`,
			}))
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			repo := &cleanupContractRepository{Repository: memory}
			service := NewService(repo, manager)
			switch operation {
			case "update":
				repo.cancel = cancel // Request expires after its successful original read.
				_, err = service.Update(ctx, "target", UpdateDeviceRequest{ConnectionConfig: map[string]any{
					"host": "127.0.0.2", "port": 502, "slave_id": 1, "timeout": 5,
				}})
			case "delete":
				cancel()
				err = service.Delete(ctx, "target")
			case "disable":
				cancel()
				err = service.Disable(ctx, "target")
			}
			require.ErrorIs(t, err, context.Canceled)
			require.Same(t, ctx, repo.observed, "repository retains the original canceled request")
			require.Equal(t, 1, protocols[0].closed)
			require.NotNil(t, protocols[0].cleanupCtx)
			require.NoError(t, protocols[0].cleanupCtx.Err())
			require.Zero(t, protocols[1].closed, "unrelated device stays connected")
			_, err = manager.GetOrCreate(t.Context(), "target", typ, `{}`)
			require.NoError(t, err)
			require.Len(t, protocols, 3, "retired generation cannot be reused")
			require.Zero(t, protocols[2].closed, "fresh replacement is not retired by the prior call")
		})
	}
}
