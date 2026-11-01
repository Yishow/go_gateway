package collector

import (
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionGeometryRollbackUsesOwnedPublicationAndFreshCompensationLease(t *testing.T) {
	for _, change := range []string{"current", "foreign-slot", "foreign-pending", "epoch", "capture-not-owned"} {
		t.Run(change, func(t *testing.T) {
			s := seedAcquisitionGeometry(t)
			token, publication := geometryPublication(t, s, "")
			owned, ok := s.PublishAcquisitionGeometry(publication, func() {})
			require.True(t, ok)
			point := &schema.Point{ID: "p", DeviceID: "d", Address: "D0", DataType: schema.DataTypeUint16}
			freshLease := s.AddPointWithObservationLease(point)
			epoch := s.pointObservationLease
			if change == "capture-not-owned" {
				owned = s.CaptureAcquisitionGeometry("d")
			}
			rollback := s.PrepareAcquisitionGeometryRollback(token, owned, schema.Device{ID: "d", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{}`}, AcquisitionGeometryDefinition{Revision: "source-semantic"}, &epoch)
			if change == "capture-not-owned" {
				require.Nil(t, rollback)
				return
			}
			require.NotNil(t, rollback)
			require.True(t, rollback.Restored("p", "", freshLease))
			switch change {
			case "foreign-slot":
				s.AddPoint(&schema.Point{ID: "foreign", DeviceID: "d", Address: "D1", DataType: schema.DataTypeUint16})
			case "foreign-pending":
				s.BeginAcquisitionGeometryMutation("d")
			case "epoch":
				s.AddDevice(&schema.Device{ID: "other"})
				s.AddPoint(&schema.Point{ID: "other", DeviceID: "other"})
			}
			called := false
			_, ok = s.PublishAcquisitionGeometry(rollback.Seal(), func() { called = true })
			require.Equal(t, change == "current", ok)
			require.Equal(t, ok, called)
			if ok {
				s.FinishAcquisitionGeometryMutation(token)
				require.Equal(t, "source-semantic", s.currentAcquisitionGeometryLocked("d").revision)
				require.Equal(t, freshLease, s.CapturePointLease("p", ""))
			}
		})
	}
}

func TestAcquisitionGeometryRollbackCannotReviveUnknownOrChangeFrozenProfile(t *testing.T) {
	s := seedAcquisitionGeometry(t)
	invalidated := s.BeginAcquisitionGeometryMutation("d")
	s.FinishAcquisitionGeometryMutation(invalidated)
	token := s.BeginAcquisitionGeometryMutation("d")
	require.False(t, s.CanRestoreAcquisitionGeometry(token))
	require.Nil(t, s.PrepareAcquisitionGeometryRollback(token, AcquisitionGeometryReceipt{}, schema.Device{ID: "d"}, AcquisitionGeometryDefinition{Revision: "source-semantic"}, nil))
	s = seedAcquisitionGeometry(t)
	token, publication := geometryPublication(t, s, "")
	owned, ok := s.PublishAcquisitionGeometry(publication, func() {})
	require.True(t, ok)
	rollback := s.PrepareAcquisitionGeometryRollback(token, owned, schema.Device{ID: "d", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{}`}, AcquisitionGeometryDefinition{Revision: "source-semantic"}, nil)
	require.NotNil(t, rollback)
	wrong := &schema.Point{ID: "p", DeviceID: "d", Address: "D20", DataType: schema.DataTypeUint16}
	lease := s.AddPointWithObservationLease(wrong, connector.ReadProfile{})
	require.True(t, rollback.Restored("p", "", lease))
	_, ok = s.PublishAcquisitionGeometry(rollback.Seal(), func() { t.Fatal("new lease alone cannot replace original source fields") })
	require.False(t, ok)
}
