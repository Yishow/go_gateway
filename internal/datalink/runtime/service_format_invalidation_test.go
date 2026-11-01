package runtime

import (
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuspendPointRevisionPreservesNewerSchedulerInstall(t *testing.T) {
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
	svc := &Service{scheduler: scheduler}
	scheduler.AddPoint(&schema.Point{ID: "point", RevisionID: "new"})
	require.False(t, svc.SuspendPointRevision("point", "old"))
	require.True(t, svc.SuspendPointRevision("point", "new"))
}
