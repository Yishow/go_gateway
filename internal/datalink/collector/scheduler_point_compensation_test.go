package collector

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestPointCompensationRejectsSameRevisionNewLease(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	p := &schema.Point{ID: "p", RevisionID: "same-revision", Enabled: true}
	owned := s.AddPointWithObservationLease(p)
	require.Equal(t, owned, s.CapturePointLease(p.ID, p.RevisionID))
	newer := s.AddPointWithObservationLease(p)
	require.NotEqual(t, owned, newer)
	require.False(t, s.SuspendPointLease(p.ID, owned))
	restored := *p
	restored.RevisionID = "old-revision"
	require.False(t, s.RestorePointLease(&restored, owned))
	require.Equal(t, newer, s.CapturePointLease(p.ID, p.RevisionID))
	require.True(t, s.RestorePointLease(&restored, newer))
	require.NotZero(t, s.CapturePointLease(p.ID, restored.RevisionID))
	require.False(t, s.RestorePointLease(&restored, newer))
}
func TestPointCompensationCannotReviveAbsentInstallation(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	p := &schema.Point{ID: "p", RevisionID: "r", Enabled: true}
	lease := s.AddPointWithObservationLease(p)
	require.True(t, s.SuspendPointLease(p.ID, lease))
	require.False(t, s.RestorePointLease(p, lease))
	require.Zero(t, s.CapturePointLease(p.ID, p.RevisionID))
}

func TestOwnedRemovalReceiptRestoresNormalBatchAndRejectsABA(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	a := &schema.Point{ID: "a", DeviceID: "d", RevisionID: "r", Enabled: true}
	b := &schema.Point{ID: "b", DeviceID: "d", RevisionID: "r", Enabled: true}
	la := s.AddPointWithObservationLease(a)
	lb := s.AddPointWithObservationLease(b)
	ea, ok := s.RemovePointOwned(a.ID, la)
	require.True(t, ok)
	eb, ok := s.RemovePointOwned(b.ID, lb)
	require.True(t, ok)
	require.Equal(t, ea, eb)
	restoredA, ok := s.RestoreRemovedPointOwned(a, la, ea)
	require.True(t, ok)
	restoredB, ok := s.RestoreRemovedPointOwned(b, lb, restoredA)
	require.True(t, ok)
	require.Greater(t, restoredB, restoredA)
	epoch, ok := s.RemovePointOwned(a.ID, restoredA)
	require.True(t, ok)
	foreign := s.AddPointWithObservationLease(a)
	require.True(t, s.SuspendPointLease(a.ID, foreign))
	_, ok = s.RestoreRemovedPointOwned(a, restoredA, epoch)
	require.False(t, ok)
}
func TestOwnedRemovalBatchRejectsForeignInstallBetweenRestorations(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	a := &schema.Point{ID: "a", DeviceID: "d", RevisionID: "r", Enabled: true}
	b := &schema.Point{ID: "b", DeviceID: "d", RevisionID: "r", Enabled: true}
	foreign := &schema.Point{ID: "foreign", DeviceID: "d", RevisionID: "r", Enabled: true}
	la := s.AddPointWithObservationLease(a)
	lb := s.AddPointWithObservationLease(b)
	epoch, ok := s.RemovePointOwned(a.ID, la)
	require.True(t, ok)
	_, ok = s.RemovePointOwned(b.ID, lb)
	require.True(t, ok)
	restored, ok := s.RestoreRemovedPointOwned(a, la, epoch)
	require.True(t, ok)
	foreignLease := s.AddPointWithObservationLease(foreign)
	_, ok = s.RestoreRemovedPointOwned(b, lb, restored)
	require.False(t, ok)
	require.True(t, s.SuspendPointLease(a.ID, restored))
	require.Equal(t, foreignLease, s.CapturePointLease(foreign.ID, foreign.RevisionID))
}

func TestOwnedPointInstallationRefusesExhaustedLeaseProof(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	p := &schema.Point{ID: "p", DeviceID: "d", RevisionID: "r", Enabled: true, DataType: schema.DataTypeInt32}
	old := s.AddPointWithObservationLease(p)
	s.pointObservationLease = ^uint64(0)
	_, ok := s.InstallPointAtLease(p, old)
	require.False(t, ok)
	require.Equal(t, old, s.CapturePointLease(p.ID, p.RevisionID))
	_, ok = s.RestorePointLeaseReceipt(p, old)
	require.False(t, ok)
	require.Equal(t, old, s.CapturePointLease(p.ID, p.RevisionID))
}
