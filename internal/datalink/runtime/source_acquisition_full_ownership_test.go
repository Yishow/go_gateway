package runtime

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionFullProjectionNeverAdoptsForeignRemovalTargets(t *testing.T) {
	for _, kind := range []string{"new-device", "removed-device-reinstalled"} {
		t.Run(kind, func(t *testing.T) {
			f := newAcquisitionInstallFixture(t, true)
			require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
			id := "foreign-device"
			if kind == "removed-device-reinstalled" {
				id = "D2"
				f.projection.Devices = f.projection.Devices[:1]
				f.projection.Points = f.projection.Points[:1]
				f.projection.DeviceIDs = f.projection.DeviceIDs[:1]
			}
			p := &schema.Point{ID: "foreign-point", DeviceID: id, Address: "40009", DataType: schema.DataTypeInt16, Enabled: true}
			var foreignLease uint64
			f.reader.beforeRead = func(int) {
				if foreignLease != 0 {
					return
				}
				f.service.scheduler.AddDevice(&schema.Device{ID: id, Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{}`})
				foreignLease = f.service.scheduler.AddPointWithObservationLease(p)
				f.service.pointMetaMu.Lock()
				f.service.setPointMetaLocked(p)
				f.service.pointMetaMu.Unlock()
			}
			_ = f.service.ApplyWorkspaceProjection(t.Context(), f.projection)
			require.NotZero(t, foreignLease)
			require.Equal(t, foreignLease, f.service.scheduler.CapturePointLease(p.ID, p.RevisionID), "old full projection may reject, but must not remove a newly observed foreign target")
			f.service.pointMetaMu.RLock()
			_, exists := f.service.pointMetaIndex[p.ID]
			f.service.pointMetaMu.RUnlock()
			require.True(t, exists)
		})
	}
}
