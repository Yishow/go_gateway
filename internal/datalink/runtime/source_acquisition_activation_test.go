package runtime

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionActivationScopeRetainsOnlyProvenForeignSlots(t *testing.T) {
	for _, proven := range []bool{false, true} {
		t.Run(map[bool]string{false: "unknown-slot", true: "matching-installation"}[proven], func(t *testing.T) {
			live, _, projection := newAcquisitionInputSQLFixture(t)
			if proven {
				require.NoError(t, live.ApplyWorkspaceProjection(t.Context(), projection))
			} else {
				live.scheduler.AddDevice(projection.Devices[0])
				live.scheduler.AddPointWithObservationLease(projection.Points[1])
			}
			retained := projection.Points[1]
			before := live.scheduler.CapturePointLease(retained.ID, retained.RevisionID)
			require.NotZero(t, before)
			require.NoError(t, live.ApplyWorkspaceProjectionForDevicesAndPoints(t.Context(), projection, []string{"d"}, []string{projection.Points[0].ID}))
			require.Equal(t, before, live.scheduler.CapturePointLease(retained.ID, retained.RevisionID))
			token := live.scheduler.BeginAcquisitionGeometryMutation("d")
			require.Equal(t, proven, live.scheduler.CanRestoreAcquisitionGeometry(token), "a selected install must not bless an unproven retained slot")
			live.scheduler.FinishAcquisitionGeometryMutation(token)
		})
	}
}

func TestAcquisitionActivationScopeRejectsForeignPublicationReentry(t *testing.T) {
	live, reader, projection := newAcquisitionInputSQLFixture(t)
	foreign := runtimeProjectionDeviceState{RuntimeVersion: "foreign", Alignment: workspace.RuntimeProjectionAlignmentStale}
	changed := false
	reader.beforeRead = func() {
		if changed {
			return
		}
		changed = true
		live.projectionMu.Lock()
		live.projections = map[string]runtimeProjectionDeviceState{"d": foreign}
		live.projectionMu.Unlock()
	}
	err := live.ApplyWorkspaceProjectionForDevicesAndPoints(t.Context(), projection, []string{"d"}, []string{"a"})
	require.True(t, changed)
	require.Error(t, err, "original selected install cannot adopt a late foreign publication")
	require.Equal(t, foreign, live.deviceProjectionState("d"))
	token := live.scheduler.BeginAcquisitionGeometryMutation("d")
	require.False(t, live.scheduler.CanRestoreAcquisitionGeometry(token))
	live.scheduler.FinishAcquisitionGeometryMutation(token)
}

// The SQL acceptor is the real publication boundary, not an after-the-fact GET.
type activationAcceptanceProbe struct {
	*acquisitionInputReentryWorkspace
	beforeAccept func()
	reject       bool
}

func (p *activationAcceptanceProbe) AcceptAcquisitionSourceGraphs(ctx context.Context, graphs []workspace.AcquisitionSourceGraph, assign func() error) error {
	if p.beforeAccept != nil {
		p.beforeAccept()
	}
	if p.reject {
		return errors.New("acceptance refused")
	}
	return p.Service.AcceptAcquisitionSourceGraphs(ctx, graphs, assign)
}
func TestAcquisitionActivationScopeNewBindingsWaitForSQLAcceptance(t *testing.T) {
	for _, mode := range []string{"accepted", "rejected", "canceled"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			live, reader, projection := newAcquisitionInputSQLFixture(t)
			require.NoError(t, live.ApplyWorkspaceProjection(t.Context(), projection))
			old := live.scheduler.CapturePointLease("a", "pa")
			beforeState := live.deviceProjectionState("d")
			projection.Tags = []*schema.Tag{{ID: "tag", DataType: schema.DataTypeUint16}}
			projection.Mappings = []*schema.Mapping{{ID: "mapping", PointID: "a", TagID: "tag", Enabled: true}}
			retained := live.scheduler.CapturePointLease("b", "pb")
			probe := &activationAcceptanceProbe{acquisitionInputReentryWorkspace: reader, reject: mode == "rejected"}
			live.workspace = probe
			called := false
			probe.beforeAccept = func() {
				called = true
				if mode == "canceled" {
					cancel()
				}
				lease := live.scheduler.CapturePointLease("a", "pa")
				require.NotEqual(t, old, lease)
				bindings, rawAccepted := live.captureCollectedBindings(collector.CollectedValue{PointID: "a", InstallationLease: lease, InstallationScope: live.scheduler.InstallationScope()})
				require.True(t, rawAccepted, "raw quality observation keeps existing semantics")
				require.Empty(t, bindings, "new downstream binding cannot escape before SQL acceptance")
				require.Equal(t, beforeState, live.deviceProjectionState("d"), "new scope has not been accepted")
				require.True(t, live.scheduler.AcceptsInstallationSample("b", retained, live.scheduler.InstallationScope()), "unselected source remains readable")
			}
			err := live.ApplyWorkspaceProjectionForDevicesAndPoints(ctx, projection, []string{"d"}, []string{"a"})
			require.True(t, called)
			if mode != "accepted" {
				require.Error(t, err)
				require.Empty(t, live.mappingIndex["a"])
				require.Equal(t, beforeState, live.deviceProjectionState("d"))
			} else {
				require.NoError(t, err)
				require.Len(t, live.mappingIndex["a"], 1)
				require.True(t, live.scheduler.AcceptsInstallationSample("a", live.scheduler.CapturePointLease("a", "pa"), live.scheduler.InstallationScope()))
			}
			require.Equal(t, retained, live.scheduler.CapturePointLease("b", "pb"))
			token := live.scheduler.BeginAcquisitionGeometryMutation("d")
			require.Equal(t, mode == "accepted", live.scheduler.CanRestoreAcquisitionGeometry(token))
			live.scheduler.FinishAcquisitionGeometryMutation(token)
		})
	}
}

func TestAcquisitionActivationWholeDeviceCannotBlessUnownedExtraSlot(t *testing.T) {
	live, _, projection := newAcquisitionInputSQLFixture(t)
	live.scheduler.AddDevice(projection.Devices[0])
	extra := &schema.Point{ID: "unowned", DeviceID: "d", Address: "D8", DataType: schema.DataTypeUint16, Enabled: true}
	lease := live.scheduler.AddPointWithObservationLease(extra)
	require.NotZero(t, lease)
	require.NoError(t, live.ApplyWorkspaceProjectionForDevices(t.Context(), projection, []string{"d"}))
	require.Equal(t, lease, live.scheduler.CapturePointLease(extra.ID, extra.RevisionID), "unowned slot must not be silently removed")
	token := live.scheduler.BeginAcquisitionGeometryMutation("d")
	require.False(t, live.scheduler.CanRestoreAcquisitionGeometry(token), "complete SQL cannot bless an unjournaled extra slot")
	live.scheduler.FinishAcquisitionGeometryMutation(token)
}
