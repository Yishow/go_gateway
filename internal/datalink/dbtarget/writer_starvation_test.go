package dbtarget

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestWriter_ContinuousUpdatesAtIntervalBoundaryStillWriteSnapshots(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	connector.DefaultWriteIntervalSeconds = 2
	require.NoError(t, connectors.Update(t.Context(), connector))
	createSnapshotMapping(t, mappings, connector.ID, "a", "d0")
	createSnapshotMapping(t, mappings, connector.ID, "z", "d1")
	base := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick, Now: clock.Now})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "a", 10, base.Add(time.Second)), ErrSnapshotAccepted)

	for i := 1; i <= 3; i++ {
		boundary := base.Add(time.Duration(2*i) * time.Second)
		require.ErrorIs(t, writer.WriteTagValue(t.Context(), "a", 9+2*i, boundary.Add(-time.Second)), ErrSnapshotAccepted)
		require.ErrorIs(t, writer.WriteTagValue(t.Context(), "z", 100+i, boundary.Add(-time.Second)), ErrSnapshotAccepted)
		clock.Set(boundary)
		require.ErrorIs(t, writer.WriteTagValue(t.Context(), "a", 10+2*i, boundary), ErrSnapshotAccepted)
		tick <- boundary
		// The second mapping is sorted after "a" and proves the worker finished
		// its due-selection pass, without relying on a sleep or timeout for absence.
		require.Eventually(t, func() bool {
			var count int
			return target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values WHERE d1 IS NOT NULL`).Scan(&count) == nil && count == i
		}, 2*time.Second, 10*time.Millisecond)
	}
	var count int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values WHERE d0 IS NOT NULL`).Scan(&count))
	require.Equal(t, 3, count, "a fresh value at every boundary must not postpone all writes")
}

func TestWriter_GroupedContinuousUpdatesAtBoundaryStillWriteSnapshots(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	connector.DefaultWriteIntervalSeconds = 2
	require.NoError(t, connectors.Update(t.Context(), connector))
	group := "meter"
	for _, entry := range []struct{ tag, column string }{{"a", "d0"}, {"z", "d1"}} {
		require.NoError(t, mappings.Create(t.Context(), &schema.DatabaseTargetMapping{
			ID: "mapping-" + entry.tag, TagID: entry.tag, ConnectorID: connector.ID,
			TableName: "snapshot_values", ColumnName: entry.column, GroupKey: &group,
			WriteMode: schema.DatabaseWriteModeInsert, Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}))
	}
	base := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick, Now: clock.Now})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "a", 10, base.Add(time.Second)), ErrSnapshotAccepted)
	for i := 1; i <= 3; i++ {
		boundary := base.Add(time.Duration(2*i) * time.Second)
		require.ErrorIs(t, writer.WriteTagValue(t.Context(), "a", 9+2*i, boundary.Add(-time.Second)), ErrSnapshotAccepted)
		require.ErrorIs(t, writer.WriteTagValue(t.Context(), "z", 100+i, boundary.Add(-time.Second)), ErrSnapshotAccepted)
		clock.Set(boundary)
		require.ErrorIs(t, writer.WriteTagValue(t.Context(), "a", 10+2*i, boundary), ErrSnapshotAccepted)
		tick <- boundary
		require.Eventually(t, func() bool {
			var count int
			return target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values WHERE d1 IS NOT NULL`).Scan(&count) == nil && count == i
		}, 2*time.Second, 10*time.Millisecond)
	}
	var count int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values WHERE d0 IS NOT NULL`).Scan(&count))
	require.Equal(t, 3, count, "new grouped values must not postpone all snapshots")
}
