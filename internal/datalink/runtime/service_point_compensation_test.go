package runtime

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
	"sync"
	"testing"
)

type compensationProjectionReader struct {
	projection *workspace.RuntimeProjection
	hook       func()
}

func (r *compensationProjectionReader) RuntimeProjection(context.Context) (*workspace.RuntimeProjection, error) {
	if r.hook != nil {
		f := r.hook
		r.hook = nil
		f()
	}
	return r.projection, nil
}

func TestSourceProjectionDoesNotOverwriteConcurrentSameRevisionLease(t *testing.T) {
	for _, operation := range []string{"install", "disable", "remove"} {
		t.Run(operation, func(t *testing.T) {
			scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
			d := &schema.Device{ID: "d", Status: schema.DeviceStatusActive, Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{"host":"127.0.0.1","port":502}`}
			p := &schema.Point{ID: "p", DeviceID: d.ID, Address: "40001", DataType: schema.DataTypeInt32, RevisionID: "same", Enabled: true}
			old := scheduler.AddPointWithObservationLease(p)
			projected := *p
			if operation == "disable" {
				projected.Enabled = false
			}
			projection := &workspace.RuntimeProjection{WorkspaceID: "ws", DeviceIDs: []string{d.ID}, Devices: []*schema.Device{d}, Points: []*schema.Point{&projected}, Alignment: workspace.RuntimeProjectionAlignmentAligned}
			if operation == "remove" {
				projection.Points = nil
			}
			var newer uint64
			reader := &compensationProjectionReader{projection: projection, hook: func() { newer = scheduler.AddPointWithObservationLease(p) }}
			svc := &Service{scheduler: scheduler, workspace: reader, mappingIndex: map[string][]mappingBinding{}, pointMetaIndex: map[string]pointMeta{p.ID: {DeviceID: d.ID, Address: p.Address}}}
			svc.running.Store(true)
			outcome := svc.ReconcileSourceRule(t.Context(), sourcerule.RuntimeReconcileRequest{Operation: sourcerule.RuntimeReconcileOperationUpdate, Scope: sourcerule.RuntimeReconcileScope{RuleID: "rule", DeviceID: d.ID}})
			require.NotEqual(t, old, newer)
			require.Equal(t, newer, scheduler.CapturePointLease(p.ID, p.RevisionID))
			require.Equal(t, sourcerule.RuntimeReconcileStatusStale, outcome.Status)
		})
	}
}

func TestOwnedPointRemovalPreservesForeignMetadata(t *testing.T) {
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
	svc := &Service{scheduler: scheduler}
	p := &schema.Point{ID: "p", DeviceID: "d", Address: "40001", RevisionID: "same", Enabled: true, DataType: schema.DataTypeInt32}
	old := svc.InstallPointForMutation(p)
	foreign := *p
	foreign.Address = "40003"
	svc.UpsertPoint(&foreign)
	_, ok := svc.RemovePointOwned(p.ID, old)
	require.False(t, ok)
	require.Equal(t, foreign.Address, svc.lookupPointMeta(p.ID).Address)
	current := svc.CapturePointLease(p.ID, p.RevisionID)
	_, ok = svc.RemovePointOwned(p.ID, current)
	require.True(t, ok)
	require.Empty(t, svc.lookupPointMeta(p.ID).Address)
}

func TestPointMetadataInstallationInterleavings(t *testing.T) {
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
	svc := &Service{scheduler: scheduler}
	p := &schema.Point{ID: "p", DeviceID: "d", Address: "40001", RevisionID: "same", Enabled: true, DataType: schema.DataTypeInt32}
	var wg sync.WaitGroup
	for worker := 0; worker < 3; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				switch worker {
				case 0:
					svc.UpsertPoint(p)
				case 1:
					svc.RemovePoint(p.ID)
				case 2:
					ctx := svc.beginSourceProjectionLeaseGuard(context.Background(), p.DeviceID)
					_, _ = svc.installProjectionPoint(ctx, p)
					_ = svc.removeProjectionPoint(ctx, p.ID)
				}
			}
		}(worker)
	}
	wg.Wait()
	svc.UpsertPoint(p)
	require.Equal(t, p.Address, svc.lookupPointMeta(p.ID).Address)
	require.NotZero(t, svc.CapturePointLease(p.ID, p.RevisionID))
	svc.RemovePoint(p.ID)
	require.Empty(t, svc.lookupPointMeta(p.ID))
	require.Zero(t, svc.CapturePointLease(p.ID, p.RevisionID))
}
