package workspace

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/device"

	"github.com/stretchr/testify/require"
)

func TestService_ReadinessBlocksEmptyWorkspace(t *testing.T) {
	workspaceSvc := NewService(NewMemoryRepository()).WithReadinessServices(
		device.NewService(device.NewMemoryRepository(), nil),
		nil,
		nil,
		nil,
	)

	summary, err := workspaceSvc.Readiness(context.Background())
	require.NoError(t, err)
	require.False(t, summary.Ready)
	require.Equal(t, 1, summary.BlockingCount)
	requireReadinessIssue(t, summary, ReadinessIssue{
		Code:     "workspace-device-missing",
		Severity: ReadinessSeverityBlocking,
		Step:     ReadinessStep1,
		Scope:    "workspace",
	})
}
