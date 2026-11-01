package collector

import (
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func geometryPublication(t *testing.T, s *Scheduler, scope string) (AcquisitionGeometryMutation, AcquisitionGeometryPublication) {
	t.Helper()
	token := s.BeginAcquisitionGeometryMutation("d")
	update := s.PrepareAcquisitionGeometry(token, schema.Device{ID: "d", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{}`}, AcquisitionGeometryDefinition{Revision: "next", Rules: map[string]string{"A": "a", "B": "b"}}, scope)
	require.NotNil(t, update)
	return token, update.Seal()
}

func TestAcquisitionGeometryPublicationIsOwnedAtomicAndNonblocking(t *testing.T) {
	s := seedAcquisitionGeometry(t)
	token, pub := geometryPublication(t, s, "")
	called := false
	s.mu.Lock()
	_, ok := s.PublishAcquisitionGeometry(pub, func() { called = true })
	s.mu.Unlock()
	require.False(t, ok, "SQL reservation must not wait for scheduler lock")
	require.False(t, called)
	receipt, ok := s.PublishAcquisitionGeometry(pub, func() { called = true })
	require.True(t, ok)
	require.True(t, called)
	require.Nil(t, s.currentAcquisitionGeometryLocked("d"), "accepted publication remains quarantined until its owner finishes")
	s.FinishAcquisitionGeometryMutation(token)
	require.Same(t, receipt.publication, s.currentAcquisitionGeometryLocked("d"))
	require.True(t, s.RevokeAcquisitionGeometry(receipt))
	require.Nil(t, s.currentAcquisitionGeometryLocked("d"))
	require.False(t, s.RevokeAcquisitionGeometry(receipt))
}

func TestAcquisitionGeometryPublicationRejectsForeignPointerOrSlotBeforeCallback(t *testing.T) {
	for _, change := range []string{"pending-operation", "remove-without-epoch", "foreign-add", "same-revision-reinstall"} {
		t.Run(change, func(t *testing.T) {
			s := seedAcquisitionGeometry(t)
			_, pub := geometryPublication(t, s, "")
			switch change {
			case "pending-operation":
				s.BeginAcquisitionGeometryMutation("d")
			case "remove-without-epoch":
				s.RemovePoint("p")
			case "foreign-add":
				s.AddPoint(&schema.Point{ID: "q", DeviceID: "d", Address: "D1", DataType: schema.DataTypeUint16})
			case "same-revision-reinstall":
				s.AddPoint(&schema.Point{ID: "p", DeviceID: "d", Address: "D0", DataType: schema.DataTypeUint16})
			}
			called := false
			_, ok := s.PublishAcquisitionGeometry(pub, func() { called = true })
			require.False(t, ok)
			require.False(t, called, "no mapping/scope writes before full receipt CAS")
		})
	}
}

func TestAcquisitionGeometryScopedPublicationPreservesForeignEmptyRules(t *testing.T) {
	for _, foreign := range []string{"a", "changed", ""} {
		t.Run("foreign="+foreign, func(t *testing.T) {
			s := seedAcquisitionGeometry(t)
			s.acquisitionGeometry["d"].publication.rules = map[string]string{"A": foreign, "B": "old-b"}
			_, pub := geometryPublication(t, s, "B")
			called := false
			_, ok := s.PublishAcquisitionGeometry(pub, func() { called = true })
			require.Equal(t, foreign == "a", ok)
			require.Equal(t, ok, called)
		})
	}
}

func TestAcquisitionGeometryCleanupCannotRevokeNewerPublication(t *testing.T) {
	s := seedAcquisitionGeometry(t)
	first, pub := geometryPublication(t, s, "")
	old, ok := s.PublishAcquisitionGeometry(pub, func() {})
	require.True(t, ok)
	s.FinishAcquisitionGeometryMutation(first)
	second, pub := geometryPublication(t, s, "")
	newer, ok := s.PublishAcquisitionGeometry(pub, func() {})
	require.True(t, ok)
	s.FinishAcquisitionGeometryMutation(second)
	require.False(t, s.RevokeAcquisitionGeometry(old))
	require.Same(t, newer.publication, s.currentAcquisitionGeometryLocked("d"))
}

func TestAcquisitionGeometrySealedPublicationCannotAdoptCopiedBuilderMutation(t *testing.T) {
	s := seedAcquisitionGeometry(t)
	token := s.BeginAcquisitionGeometryMutation("d")
	update := s.PrepareAcquisitionGeometry(token, schema.Device{ID: "d", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{}`}, AcquisitionGeometryDefinition{Revision: "source"}, "")
	shadow := *update
	publication := update.Seal()
	_, ok := s.PublishAcquisitionGeometry(publication, func() {})
	require.True(t, ok)
	s.FinishAcquisitionGeometryMutation(token)
	require.NotNil(t, s.currentAcquisitionGeometryLocked("d"))
	p := &schema.Point{ID: "p", DeviceID: "d", Address: "D0", DataType: schema.DataTypeUint16}
	lease := s.AddPointWithObservationLease(p)
	require.Nil(t, s.currentAcquisitionGeometryLocked("d"))
	require.True(t, shadow.Installed(p, connector.ReadProfile{}, lease))
	require.Nil(t, s.currentAcquisitionGeometryLocked("d"), "a copied mutable builder cannot rewrite the frozen published slot proof")
}

func TestAcquisitionGeometryFinalizedTokenCannotApplyDeviceConfiguration(t *testing.T) {
	s := seedAcquisitionGeometry(t)
	token := s.BeginAcquisitionGeometryMutation("d")
	update := s.PrepareAcquisitionGeometry(token, schema.Device{ID: "d", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{"timeout":2}`}, AcquisitionGeometryDefinition{}, "")
	s.FinishAcquisitionGeometryMutation(token)
	require.False(t, s.ApplyAcquisitionDevice(update), "a canceled/finalized owner cannot apply a delayed provider result")
	require.Equal(t, `{}`, s.deviceConfigs["d"].Config)
}
