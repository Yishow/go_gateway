package dbtarget

import (
	"context"
	"fmt"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestWriter_HealthySlowMappingsDoNotShareOneFlushDeadline(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	for i := range 40 {
		createSnapshotMapping(t, mappings, connector.ID, fmt.Sprintf("slow-%02d", i), "d0")
	}
	base := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick, Now: clock.Now,
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			time.Sleep(150 * time.Millisecond)
			return openExternalDBManagerFunc(kind, config)
		},
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	for i := range 40 {
		require.ErrorIs(t, writer.WriteTagValue(t.Context(), fmt.Sprintf("slow-%02d", i), i, base.Add(time.Second)), ErrSnapshotAccepted)
	}
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()
	require.Eventually(t, func() bool {
		return countHistoryStatus(connector.ID, "success") == 40 || countHistoryStatus(connector.ID, "failed") > 0
	}, 12*time.Second, 20*time.Millisecond)
	require.Zero(t, countHistoryStatus(connector.ID, "failed"), "healthy writes must not inherit an earlier mapping's timeout")
	require.Equal(t, 40, countHistoryStatus(connector.ID, "success"))
	var rows int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&rows))
	require.Equal(t, 40, rows)
	updated, err := connectors.GetByID(t.Context(), connector.ID)
	require.NoError(t, err)
	require.Equal(t, deliveryOutcomeSuccess, updated.LastWriteStatus)
	require.Empty(t, updated.LastWriteError)
}
