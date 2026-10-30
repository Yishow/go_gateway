package runtime

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestSuspendPointRevisionPreservesNewerSchedulerInstall(t *testing.T) {
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
	svc := &Service{scheduler: scheduler}
	scheduler.AddPoint(&schema.Point{ID: "point", RevisionID: "new"})
	require.False(t, svc.SuspendPointRevision("point", "old"))
	require.True(t, svc.SuspendPointRevision("point", "new"))
}
