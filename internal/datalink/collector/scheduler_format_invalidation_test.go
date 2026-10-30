package collector

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestRemovePointRevisionNeverRemovesNewInstallation(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	p := &schema.Point{ID: "point", RevisionID: "old"}
	s.AddPoint(p)
	require.False(t, s.RemovePointRevision(p.ID, "foreign"))
	require.Contains(t, s.pointInfos, p.ID)
	require.True(t, s.RemovePointRevision(p.ID, "old"))
	require.NotContains(t, s.pointInfos, p.ID)
	p.RevisionID = "new"
	s.AddPoint(p)
	require.False(t, s.RemovePointRevision(p.ID, "old"))
	require.Contains(t, s.pointInfos, p.ID)
	require.True(t, s.RemovePointRevision("absent", "old"))
}
