package dbtarget

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func createDecoupleGroupMapping(t *testing.T, repo *SQLTargetMappingRepository, connectorID, tagID, column string) {
	t.Helper()
	group := "meter"
	require.NoError(t, repo.Create(t.Context(), &schema.DatabaseTargetMapping{
		ID: "mapping-" + tagID, TagID: tagID, ConnectorID: connectorID,
		TableName: "snapshot_values", ColumnName: column, GroupKey: &group,
		WriteMode: schema.DatabaseWriteModeInsert, Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}))
}

func TestSnapshotWorker_GroupedIntervalsDoNotMixAtBoundary(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	createDecoupleGroupMapping(t, mappings, connector.ID, "boundary-a", "d0")
	createDecoupleGroupMapping(t, mappings, connector.ID, "boundary-b", "d1")
	base := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick, Now: clock.Now})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })

	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "boundary-a", 10, base.Add(time.Second)), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "boundary-b", 20, base.Add(2*time.Second)), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "boundary-a", 11, base.Add(15*time.Second)), ErrSnapshotAccepted)
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()

	var firstA, firstB sql.NullFloat64
	require.Eventually(t, func() bool {
		return target.QueryRowContext(t.Context(), `SELECT d0, d1 FROM snapshot_values ORDER BY rowid LIMIT 1`).Scan(&firstA, &firstB) == nil
	}, 3*time.Second, 10*time.Millisecond, "the first grouped interval was not written")
	require.Equal(t, sql.NullFloat64{Float64: 10, Valid: true}, firstA)
	require.Equal(t, sql.NullFloat64{Float64: 20, Valid: true}, firstB)

	clock.Set(base.Add(30 * time.Second))
	tick <- clock.Now()
	var secondA, secondB sql.NullFloat64
	require.Eventually(t, func() bool {
		return target.QueryRowContext(t.Context(), `SELECT d0, d1 FROM snapshot_values ORDER BY rowid LIMIT 1 OFFSET 1`).Scan(&secondA, &secondB) == nil
	}, 3*time.Second, 10*time.Millisecond, "the next grouped interval was not written")
	require.Equal(t, sql.NullFloat64{Float64: 11, Valid: true}, secondA)
	require.False(t, secondB.Valid, "a value from the previous interval must not be copied into the next")
}

func TestSnapshotWorker_CloseFlushesNewGroupedInterval(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	createDecoupleGroupMapping(t, mappings, connector.ID, "close-boundary", "d0")
	base := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick, Now: clock.Now})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })

	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "close-boundary", 10, base.Add(time.Second)), ErrSnapshotAccepted)
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()
	require.Eventually(t, func() bool {
		var rows int
		return target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&rows) == nil && rows == 1
	}, 3*time.Second, 10*time.Millisecond, "the first interval was not written")

	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "close-boundary", 11, base.Add(16*time.Second)), ErrSnapshotAccepted)
	clock.Set(base.Add(17 * time.Second))
	closeCtx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	require.NoError(t, writer.Close(closeCtx))
	rows, err := target.QueryContext(t.Context(), `SELECT d0 FROM snapshot_values ORDER BY rowid`)
	require.NoError(t, err)
	defer rows.Close()
	var values []float64
	for rows.Next() {
		var value float64
		require.NoError(t, rows.Scan(&value))
		values = append(values, value)
	}
	require.NoError(t, rows.Err())
	require.Equal(t, []float64{10, 11}, values, fmt.Sprintf("Close returned success at %s", clock.Now()))
}

func TestSnapshotWorker_GroupedValuesStayInTheirIntervalsDuringSlowWrite(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	createDecoupleGroupMapping(t, mappings, connector.ID, "slow-a", "d0")
	createDecoupleGroupMapping(t, mappings, connector.ID, "slow-b", "d1")
	base := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	entered, release := make(chan struct{}), make(chan struct{})
	releaseOpen := sync.OnceFunc(func() { close(release) })
	var first sync.Once
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
		FlushTick: tick, Now: clock.Now,
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			first.Do(func() { close(entered); <-release })
			return openExternalDBManager(kind, config)
		},
	})
	t.Cleanup(func() { releaseOpen(); _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "slow-a", 10, base.Add(time.Second)), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "slow-b", 20, base.Add(time.Second)), ErrSnapshotAccepted)
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()
	select {
	case <-entered:
	case <-time.After(3 * time.Second):
		t.Fatal("the first grouped write did not start")
	}
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "slow-a", 11, base.Add(31*time.Second)), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "slow-b", 21, base.Add(32*time.Second)), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "slow-a", 12, base.Add(46*time.Second)), ErrSnapshotAccepted)
	clock.Set(base.Add(60 * time.Second))
	tick <- clock.Now()
	releaseOpen()
	require.Eventually(t, func() bool {
		var count int
		return target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&count) == nil && count == 3
	}, 3*time.Second, 10*time.Millisecond, "each observed interval must retain its own grouped row")
	rows, err := target.QueryContext(t.Context(), `SELECT d0, d1 FROM snapshot_values ORDER BY rowid`)
	require.NoError(t, err)
	defer rows.Close()
	type row struct{ a, b sql.NullFloat64 }
	var actual []row
	for rows.Next() {
		var value row
		require.NoError(t, rows.Scan(&value.a, &value.b))
		actual = append(actual, value)
	}
	require.NoError(t, rows.Err())
	require.Equal(t, []row{
		{sql.NullFloat64{Float64: 10, Valid: true}, sql.NullFloat64{Float64: 20, Valid: true}},
		{sql.NullFloat64{}, sql.NullFloat64{Float64: 21, Valid: true}},
		{sql.NullFloat64{Float64: 12, Valid: true}, sql.NullFloat64{}},
	}, actual)
}

func TestSnapshotWorker_LateGroupedValueDoesNotMoveToNextInterval(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	createDecoupleGroupMapping(t, mappings, connector.ID, "late-a", "d0")
	createDecoupleGroupMapping(t, mappings, connector.ID, "late-b", "d1")
	base := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick, Now: clock.Now})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "late-a", 10, base.Add(time.Second)), ErrSnapshotAccepted)
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()
	require.Eventually(t, func() bool {
		var count int
		return target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&count) == nil && count == 1
	}, 3*time.Second, 10*time.Millisecond)

	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "late-a", 99, base.Add(10*time.Second)), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "late-b", 21, base.Add(16*time.Second)), ErrSnapshotAccepted)
	clock.Set(base.Add(30 * time.Second))
	tick <- clock.Now()
	var lateA, currentB sql.NullFloat64
	require.Eventually(t, func() bool {
		return target.QueryRowContext(t.Context(), `SELECT d0, d1 FROM snapshot_values ORDER BY rowid LIMIT 1 OFFSET 1`).Scan(&lateA, &currentB) == nil
	}, 3*time.Second, 10*time.Millisecond)
	require.False(t, lateA.Valid, "a late value from the previous interval must not be moved into this row")
	require.Equal(t, sql.NullFloat64{Float64: 21, Valid: true}, currentB)
	var count int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&count))
	require.Equal(t, 2, count, "an already written interval must not be written again")
	history, err := NewConnectorService(connectors, mappings).ListWriteHistory(t.Context(), connector.ID, 20)
	require.NoError(t, err)
	var skipped bool
	for _, record := range history {
		if record.MappingID == "mapping-late-a" && record.SkipReason == "interval_already_flushed" {
			skipped = true
			break
		}
	}
	require.True(t, skipped, "the skipped late value must be visible in write history")
}

func TestSnapshotWorker_GroupedSnapshotsDoNotReturnAfterMappingChangesBack(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	createDecoupleGroupMapping(t, mappings, connector.ID, "changed", "d0")
	base := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base.Add(31 * time.Second)}
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: make(chan time.Time), Now: clock.Now})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "changed", 10, base.Add(time.Second)), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "changed", 11, base.Add(16*time.Second)), ErrSnapshotAccepted)

	mapping, err := mappings.GetByID(t.Context(), "mapping-changed")
	require.NoError(t, err)
	mapping.ColumnName = "d1"
	require.NoError(t, mappings.Update(t.Context(), mapping))
	require.NoError(t, writer.RefreshMappings(t.Context()))
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "changed", 12, base.Add(17*time.Second)), ErrSnapshotAccepted)
	mapping.ColumnName = "d0"
	require.NoError(t, mappings.Update(t.Context(), mapping))
	require.NoError(t, writer.RefreshMappings(t.Context()))

	closeCtx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	require.NoError(t, writer.Close(closeCtx))
	var rows int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&rows))
	require.Zero(t, rows, "a superseded mapping must not regain its earlier pending snapshots")
}

func TestSnapshotWorker_CloseFlushesAllRetainedGroupedIntervals(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	createDecoupleGroupMapping(t, mappings, connector.ID, "close-a", "d0")
	createDecoupleGroupMapping(t, mappings, connector.ID, "close-b", "d1")
	base := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base.Add(46 * time.Second)}
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: make(chan time.Time), Now: clock.Now})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "close-a", 10, base.Add(time.Second)), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "close-b", 20, base.Add(16*time.Second)), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "close-b", 21, base.Add(31*time.Second)), ErrSnapshotAccepted)
	closeCtx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	require.NoError(t, writer.Close(closeCtx))

	rows, err := target.QueryContext(t.Context(), `SELECT d0, d1 FROM snapshot_values ORDER BY rowid`)
	require.NoError(t, err)
	defer rows.Close()
	type row struct{ a, b sql.NullFloat64 }
	var actual []row
	for rows.Next() {
		var value row
		require.NoError(t, rows.Scan(&value.a, &value.b))
		actual = append(actual, value)
	}
	require.NoError(t, rows.Err())
	require.Equal(t, []row{
		{sql.NullFloat64{Float64: 10, Valid: true}, sql.NullFloat64{}},
		{sql.NullFloat64{}, sql.NullFloat64{Float64: 20, Valid: true}},
		{sql.NullFloat64{}, sql.NullFloat64{Float64: 21, Valid: true}},
	}, actual, "Close must not drop a retained bucket after reporting success")
}
