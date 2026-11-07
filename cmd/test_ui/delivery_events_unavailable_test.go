package main

import (
	"os"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"

	"github.com/stretchr/testify/require"
)

func TestProductionDeliveryEventReportsUnavailableDestination(t *testing.T) {
	f := newDeliveryEventFixture(t, true)
	path := f.env.targets["connector-A"]
	require.NotEmpty(t, path)
	offline := path + ".offline"
	require.NoError(t, os.Rename(path, offline))
	t.Cleanup(func() { require.NoError(t, os.Rename(offline, path)) })

	f.start(t)
	require.Eventually(t, func() bool { return f.state(t).State == groupdelivery.StateRetrying }, 8*time.Second, 20*time.Millisecond)
	require.Equal(t, "target-unavailable", f.state(t).LastErrorCode)
	events := waitDeliveryEvents(t, f.broker, "domain.delivery.retrying")
	require.Equal(t, "target-unavailable", events[0].Fields["reason"])
	require.Equal(t, "warn", events[0].Level)
	require.NoError(t, f.broker.Flush(t.Context()))
	data, err := os.ReadFile(f.logPath)
	require.NoError(t, err)
	require.Contains(t, string(data), `"reason":"target-unavailable"`)
	require.NotContains(t, string(data), path)
}
