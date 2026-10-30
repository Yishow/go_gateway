package collector

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestOrdinaryLeaseProofRejectsUnrelatedAdvanceWithoutReplacingPoint(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	p := &schema.Point{ID: "p", RevisionID: "old", Enabled: true}
	lease := s.AddPointWithObservationLease(p)
	proof, ok := s.CaptureOrdinaryPoint(p.ID)
	require.True(t, ok)
	s.AddPoint(&schema.Point{ID: "unrelated", RevisionID: "new", Enabled: true})
	next := *p
	next.RevisionID = "attempt"
	require.False(t, s.ApplyOrdinaryPoint(&next, proof))
	require.Equal(t, lease, s.CapturePointLease(p.ID, p.RevisionID))
}
func TestOrdinaryAbsentProofRejectsFormatInvalidationWithoutInstall(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	proof, ok := s.CaptureOrdinaryPoint("p")
	require.True(t, ok)
	require.True(t, s.InvalidatePointFormat("p", "old"))
	require.False(t, s.ApplyOrdinaryPoint(&schema.Point{ID: "p", RevisionID: "attempt", Enabled: true}, proof))
	require.Zero(t, s.CapturePointLease("p", "attempt"))
}

func TestOrdinaryUnknownInvalidationStopsOwnedLeaseDespiteUnrelatedAdvance(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	p := &schema.Point{ID: "p", RevisionID: "old", Enabled: true}
	s.AddPoint(p)
	proof, ok := s.CaptureOrdinaryPoint(p.ID)
	require.True(t, ok)
	s.AddPoint(&schema.Point{ID: "other", RevisionID: "foreign", Enabled: true})
	require.True(t, s.InvalidateOrdinaryPoint(proof))
	require.Zero(t, s.CapturePointLease(p.ID, p.RevisionID))
	require.NotZero(t, s.CapturePointLease("other", "foreign"))
}
