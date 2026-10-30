package runtime

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
	"testing"
)

type compensationWorkspaceFunc func(context.Context) (*workspace.RuntimeProjection, error)

func (f compensationWorkspaceFunc) RuntimeProjection(ctx context.Context) (*workspace.RuntimeProjection, error) {
	return f(ctx)
}

type compensationShareFailure struct{ closed int }

func (f *compensationShareFailure) ReconcileSourceRuleShare(context.Context, sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
	return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusFailed}, errors.New("share failed")
}
func (f *compensationShareFailure) FailClosed(context.Context) { f.closed++ }

func TestSourceRuleCompensationRestoresOwnedRemovalsAcrossDeviceProjection(t *testing.T) {
	for _, operation := range []string{"disable", "delete", "foreign-gap"} {
		t.Run(operation, func(t *testing.T) {
			ctx := t.Context()
			dr := device.NewMemoryRepository()
			ready := true
			d := &schema.Device{ID: "d", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: `{"host":"127.0.0.1","port":502}`, LastTestSuccess: &ready}
			require.NoError(t, dr.Create(ctx, d))
			ps := point.NewService(point.NewMemoryRepository(), nil)
			scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
			runtime := &Service{scheduler: scheduler, pointSvc: ps, mappingIndex: map[string][]mappingBinding{}, pointMetaIndex: map[string]pointMeta{}}
			rules := sourcerule.NewService(sourcerule.NewMemoryRepository(), device.NewService(dr, nil), ps, runtime)
			main, err := rules.Create(ctx, sourcerule.CreateRuleRequest{ID: "main", DeviceID: d.ID, StartAddress: "40001", Count: 2, DataType: schema.DataTypeInt32, NamingPrefix: "MAIN", Enabled: true})
			require.NoError(t, err)
			other, err := rules.Create(ctx, sourcerule.CreateRuleRequest{ID: "other", DeviceID: d.ID, StartAddress: "40011", Count: 1, DataType: schema.DataTypeInt32, NamingPrefix: "OTHER", Enabled: true})
			require.NoError(t, err)
			links, err := rules.ListLinks(ctx, main.ID)
			require.NoError(t, err)
			original := map[string]*schema.Point{}
			for _, link := range links {
				p, e := ps.GetByID(ctx, link.PointID)
				require.NoError(t, e)
				original[p.ID] = p
				require.NotZero(t, scheduler.CapturePointLease(p.ID, p.RevisionID))
			}
			otherLinks, err := rules.ListLinks(ctx, other.ID)
			require.NoError(t, err)
			otherPoint, err := ps.GetByID(ctx, otherLinks[0].PointID)
			require.NoError(t, err)
			var foreignLease uint64
			runtime.workspace = compensationWorkspaceFunc(func(ctx context.Context) (*workspace.RuntimeProjection, error) {
				allRules, e := rules.List(ctx, sourcerule.ListFilter{})
				if e != nil {
					return nil, e
				}
				points, e := ps.ListByDevice(ctx, d.ID)
				if e != nil {
					return nil, e
				}
				var allLinks []*schema.SourceRuleLink
				for _, r := range allRules {
					ls, e := rules.ListLinks(ctx, r.ID)
					if e != nil {
						return nil, e
					}
					allLinks = append(allLinks, ls...)
				}
				if operation == "foreign-gap" && foreignLease == 0 {
					foreignLease = scheduler.AddPointWithObservationLease(otherPoint)
				}
				return &workspace.RuntimeProjection{WorkspaceID: "ws", Version: "projection", Alignment: workspace.RuntimeProjectionAlignmentAligned, DeviceIDs: []string{d.ID}, Devices: []*schema.Device{d}, Points: points, Rules: allRules, RuleLinks: allLinks}, nil
			})
			runtime.running.Store(true)
			share := &compensationShareFailure{}
			rules.SetShareRuntimeReconciler(share)
			if operation == "delete" {
				_, err = rules.DeleteWithRuntimeReconcile(ctx, main.ID)
			} else {
				_, err = rules.DisableWithRuntimeReconcile(ctx, main.ID)
			}
			var typed *modbusshare.Error
			require.ErrorAs(t, err, &typed)
			if operation == "foreign-gap" {
				require.Equal(t, modbusshare.ErrCodeDirtyUnknown, typed.Code)
				require.Positive(t, share.closed)
				require.Equal(t, foreignLease, scheduler.CapturePointLease(otherPoint.ID, otherPoint.RevisionID))
			} else {
				require.Equal(t, modbusshare.ErrCodeReconcileFailed, typed.Code)
				for id, p := range original {
					got, e := ps.GetByID(ctx, id)
					require.NoError(t, e)
					require.Equal(t, p.RevisionID, got.RevisionID)
					require.True(t, got.Enabled)
					require.NotZero(t, scheduler.CapturePointLease(id, p.RevisionID))
				}
				require.NotZero(t, scheduler.CapturePointLease(otherPoint.ID, otherPoint.RevisionID))
			}
		})
	}
}
