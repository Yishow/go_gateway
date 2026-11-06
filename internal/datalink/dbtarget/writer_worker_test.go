package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type failingWriteOutcomeRepository struct{ ConnectorRepository }

func (failingWriteOutcomeRepository) UpdateWriteOutcome(context.Context, string, time.Time, string, string) error {
	return errors.New("outcome storage unavailable")
}

type flushOutcomeSignalRepository struct {
	ConnectorRepository
	completed chan time.Time
}

func (r *flushOutcomeSignalRepository) UpdateFlushOutcome(ctx context.Context, id string, at time.Time, status, summary string) error {
	if err := r.ConnectorRepository.UpdateFlushOutcome(ctx, id, at, status, summary); err != nil {
		return err
	}
	if status == deliveryOutcomeSuccess {
		r.completed <- at
	}
	return nil
}

type blockingRefreshRepository struct {
	TargetMappingRepository
	blocked atomic.Bool
	entered chan struct{}
}

func (r *blockingRefreshRepository) List(ctx context.Context, filter TargetMappingListFilter) ([]*schema.DatabaseTargetMapping, error) {
	if r.blocked.Load() {
		close(r.entered)
		<-ctx.Done()
		return nil, ctx.Err()
	}
	return r.TargetMappingRepository.List(ctx, filter)
}

func TestWriter_CloseFlushesSnapshotWhenRefreshIsCanceled(t *testing.T) {
	db := openMigratedTestDB(t)
	connectors, mappings := NewSQLConnectorRepository(db), NewSQLTargetMappingRepository(db)
	path := filepath.Join(t.TempDir(), "close-refresh.db")
	target, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })
	_, err = target.ExecContext(t.Context(), `CREATE TABLE snapshot_values (d0 REAL)`)
	require.NoError(t, err)
	connector := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindSQLite, fmt.Sprintf(`{"dsn":%q}`, path), 15)
	createSnapshotMapping(t, mappings, connector.ID, "close-refresh", "d0")
	at := time.Date(2026, 11, 4, 10, 0, 1, 0, time.UTC)
	repo := &blockingRefreshRepository{TargetMappingRepository: mappings, entered: make(chan struct{})}
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, repo, WriterConfig{FlushTick: tick, Now: func() time.Time { return at.Add(time.Minute) }})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "close-refresh", 12, at), ErrSnapshotAccepted)
	repo.blocked.Store(true)
	tick <- at.Add(time.Minute)
	select {
	case <-repo.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("mapping refresh did not start")
	}
	require.NoError(t, writer.Close(context.Background()))
	var value float64
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT d0 FROM snapshot_values`).Scan(&value))
	require.Equal(t, 12.0, value)
}

func TestWriter_RecordsActualWriteEvenWhenOutcomeUpdateFails(t *testing.T) {
	db := openMigratedTestDB(t)
	connectors, mappings := NewSQLConnectorRepository(db), NewSQLTargetMappingRepository(db)
	path := filepath.Join(t.TempDir(), "outcome.db")
	target, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })
	_, err = target.ExecContext(t.Context(), `CREATE TABLE snapshot_values (d0 REAL)`)
	require.NoError(t, err)
	connector := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindSQLite, fmt.Sprintf(`{"dsn":%q}`, path), 15)
	createSnapshotMapping(t, mappings, connector.ID, "outcome", "d0")
	resetWriteHistory(connector.ID)
	writer := NewWriterWithConfig(failingWriteOutcomeRepository{connectors}, mappings, WriterConfig{FlushTick: make(chan time.Time)})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	at := time.Date(2026, 11, 4, 10, 0, 1, 0, time.UTC)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "outcome", 12, at), ErrSnapshotAccepted)
	require.ErrorContains(t, writer.flushDue(t.Context(), at.Add(time.Minute), false), "outcome storage unavailable")
	history := listWriteHistory(connector.ID, 10)
	require.Len(t, history, 1)
	require.Equal(t, deliveryOutcomeSuccess, history[0].Status)
	var value float64
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT d0 FROM snapshot_values`).Scan(&value))
	require.Equal(t, 12.0, value)
}

func TestWriter_RecordsFailureAfterDestinationTimeout(t *testing.T) {
	db := openMigratedTestDB(t)
	connectors, mappings := NewSQLConnectorRepository(db), NewSQLTargetMappingRepository(db)
	connector := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindPostgres, `{"host":"192.0.2.1","database":"lost"}`, 15)
	createSnapshotMapping(t, mappings, connector.ID, "timeout", "d0")
	resetWriteHistory(connector.ID)
	entered, release := make(chan struct{}), make(chan struct{})
	markEntered := sync.OnceFunc(func() { close(entered) })
	releaseOpen := sync.OnceFunc(func() { close(release) })
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: make(chan time.Time),
		OpenDestination: func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
			markEntered()
			<-release
			return nil, context.DeadlineExceeded
		},
	})
	t.Cleanup(func() { releaseOpen(); _ = writer.Close(context.Background()) })
	at := time.Date(2026, 11, 4, 10, 0, 1, 0, time.UTC)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "timeout", 12, at), ErrSnapshotAccepted)
	ctx, cancel := context.WithTimeout(t.Context(), 300*time.Millisecond)
	defer cancel()
	err := writer.flushDue(ctx, at.Add(time.Minute), false)
	releaseOpen()
	require.ErrorIs(t, err, context.DeadlineExceeded)
	select {
	case <-entered:
	default:
		t.Fatal("the connector was not attempted before the deadline")
	}
	updated, err := connectors.GetByID(t.Context(), connector.ID)
	require.NoError(t, err)
	require.Equal(t, deliveryOutcomeFailed, updated.LastWriteStatus)
	history := listWriteHistory(connector.ID, 10)
	require.Len(t, history, 1)
	require.Equal(t, deliveryOutcomeFailed, history[0].Status)
}

func TestWriter_RecordsGroupedFailureAfterDestinationTimeout(t *testing.T) {
	db := openMigratedTestDB(t)
	connectors, mappings := NewSQLConnectorRepository(db), NewSQLTargetMappingRepository(db)
	connector := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindPostgres, `{"host":"192.0.2.1","database":"lost"}`, 15)
	group := "meter"
	require.NoError(t, mappings.Create(t.Context(), &schema.DatabaseTargetMapping{
		ID: "group-timeout", TagID: "group-timeout", ConnectorID: connector.ID, TableName: "snapshot_values",
		ColumnName: "d0", GroupKey: &group, WriteMode: schema.DatabaseWriteModeInsert,
		Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}))
	resetWriteHistory(connector.ID)
	entered, release := make(chan struct{}), make(chan struct{})
	markEntered := sync.OnceFunc(func() { close(entered) })
	releaseOpen := sync.OnceFunc(func() { close(release) })
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: make(chan time.Time),
		OpenDestination: func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
			markEntered()
			<-release
			return nil, context.DeadlineExceeded
		},
	})
	t.Cleanup(func() { releaseOpen(); _ = writer.Close(context.Background()) })
	at := time.Date(2026, 11, 4, 10, 0, 1, 0, time.UTC)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "group-timeout", 12, at), ErrSnapshotAccepted)
	ctx, cancel := context.WithTimeout(t.Context(), 300*time.Millisecond)
	defer cancel()
	err := writer.flushDue(ctx, at.Add(time.Minute), false)
	releaseOpen()
	require.ErrorIs(t, err, context.DeadlineExceeded)
	select {
	case <-entered:
	default:
		t.Fatal("the connector was not attempted before the deadline")
	}
	updated, err := connectors.GetByID(t.Context(), connector.ID)
	require.NoError(t, err)
	require.Equal(t, deliveryOutcomeFailed, updated.LastFlushStatus)
	history := listWriteHistory(connector.ID, 10)
	require.Len(t, history, 1)
	require.Equal(t, deliveryOutcomeFailed, history[0].Status)
}

func TestWriter_BlockedConnectorDoesNotDelayOtherConnector(t *testing.T) {
	db := openMigratedTestDB(t)
	connectors, mappings := NewSQLConnectorRepository(db), NewSQLTargetMappingRepository(db)
	path := filepath.Join(t.TempDir(), "fast.db")
	target, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })
	_, err = target.ExecContext(t.Context(), `CREATE TABLE snapshot_values (d0 REAL)`)
	require.NoError(t, err)
	slow := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindPostgres, `{"host":"192.0.2.1","database":"lost"}`, 15)
	fast := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindSQLite, fmt.Sprintf(`{"dsn":%q}`, path), 15)
	createSnapshotMapping(t, mappings, slow.ID, "slow", "d0")
	createSnapshotMapping(t, mappings, fast.ID, "fast", "d0")
	blocked := newBlockingDestination(t)
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
		FlushTick: tick,
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			if kind == schema.DatabaseConnectorKindPostgres {
				return blocked.Open(kind, config)
			}
			return openExternalDBManagerFunc(kind, config)
		},
	})
	t.Cleanup(func() { blocked.Release(); _ = writer.Close(context.Background()) })
	at := time.Date(2026, 11, 4, 10, 0, 1, 0, time.UTC)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "slow", 1, at), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "fast", 12, at), ErrSnapshotAccepted)
	tick <- at.Add(time.Minute)
	select {
	case <-blocked.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("slow connector was not attempted")
	}
	require.Eventually(t, func() bool {
		var value float64
		return target.QueryRowContext(t.Context(), `SELECT d0 FROM snapshot_values`).Scan(&value) == nil && value == 12
	}, 2*time.Second, 10*time.Millisecond, "the other connector must complete while the first is blocked")
}

func TestWriter_CloseHonorsDeadlineWhenConnectorIsBlocked(t *testing.T) {
	db := openMigratedTestDB(t)
	connectors, mappings := NewSQLConnectorRepository(db), NewSQLTargetMappingRepository(db)
	connector := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindPostgres, `{"host":"192.0.2.1","database":"lost"}`, 15)
	createSnapshotMapping(t, mappings, connector.ID, "blocked", "d0")
	blocked := newBlockingDestination(t)
	tick := make(chan time.Time, 1)
	at := time.Date(2026, 11, 4, 10, 0, 1, 0, time.UTC)
	var attempts atomic.Int64
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
		FlushTick: tick, Now: func() time.Time { return at.Add(time.Minute) },
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			attempts.Add(1)
			return blocked.Open(kind, config)
		},
	})
	t.Cleanup(func() { blocked.Release(); _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "blocked", 1, at), ErrSnapshotAccepted)
	tick <- at.Add(time.Minute)
	select {
	case <-blocked.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("connector was not attempted")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Millisecond)
	defer cancel()
	require.ErrorIs(t, writer.Close(ctx), context.DeadlineExceeded)
	blocked.Release()
	require.ErrorContains(t, writer.Close(context.Background()), "backoff")
	require.EqualValues(t, 1, attempts.Load(), "shutdown must not retry before the backoff expires")
}

func TestWriter_GroupedOwnershipIsCheckedAtFlush(t *testing.T) {
	db := openMigratedTestDB(t)
	connectors, mappings := NewSQLConnectorRepository(db), NewSQLTargetMappingRepository(db)
	path := filepath.Join(t.TempDir(), "group.db")
	target, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })
	_, err = target.ExecContext(t.Context(), `CREATE TABLE snapshot_values (d0 REAL, d1 REAL)`)
	require.NoError(t, err)
	connector := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindSQLite, fmt.Sprintf(`{"dsn":%q}`, path), 15)
	group := "meter"
	for i, tagID := range []string{"owned", "free"} {
		require.NoError(t, mappings.Create(t.Context(), &schema.DatabaseTargetMapping{
			ID: "map-" + tagID, TagID: tagID, ConnectorID: connector.ID, TableName: "snapshot_values",
			ColumnName: fmt.Sprintf("d%d", i), GroupKey: &group, WriteMode: schema.DatabaseWriteModeInsert,
			Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}))
	}
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick,
		SuppressMapping: func(_, tagID string) bool { return tagID == "owned" },
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	at := time.Date(2026, 11, 4, 10, 0, 1, 0, time.UTC)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "owned", 10, at), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "free", 12, at), ErrSnapshotAccepted)
	tick <- at.Add(15 * time.Second)
	require.Eventually(t, func() bool {
		var owned, free sql.NullFloat64
		err := target.QueryRowContext(t.Context(), `SELECT d0, d1 FROM snapshot_values`).Scan(&owned, &free)
		return err == nil && !owned.Valid && free.Valid && free.Float64 == 12
	}, 2*time.Second, 10*time.Millisecond)
}

func TestWriter_RejectsIngestWhenMappingRefreshFails(t *testing.T) {
	db := openMigratedTestDB(t)
	connectors, mappings := NewSQLConnectorRepository(db), NewSQLTargetMappingRepository(db)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: make(chan time.Time)})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.NoError(t, db.Close())
	require.Error(t, writer.RefreshMappings(t.Context()))
	accepted, err := writer.AcceptTagValue(t.Context(), "tag-1", 12, time.Now())
	require.False(t, accepted)
	require.Error(t, err)
	require.False(t, errors.Is(err, ErrSnapshotAccepted))
}

func TestWriter_DoesNotFlushDisabledConnector(t *testing.T) {
	db := openMigratedTestDB(t)
	connectors, mappings := NewSQLConnectorRepository(db), NewSQLTargetMappingRepository(db)
	path := filepath.Join(t.TempDir(), "disabled.db")
	target, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })
	_, err = target.ExecContext(t.Context(), `CREATE TABLE snapshot_values (d0 REAL)`)
	require.NoError(t, err)
	connector := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindSQLite, fmt.Sprintf(`{"dsn":%q}`, path), 15)
	createSnapshotMapping(t, mappings, connector.ID, "disabled", "d0")
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: make(chan time.Time)})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	at := time.Date(2026, 11, 4, 10, 0, 1, 0, time.UTC)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "disabled", 12, at), ErrSnapshotAccepted)
	connector.Enabled = false
	require.NoError(t, connectors.Update(t.Context(), connector))
	require.NoError(t, writer.RefreshMappings(t.Context()))
	require.NoError(t, writer.flushDue(t.Context(), at.Add(time.Minute), false))
	var rows int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&rows))
	require.Zero(t, rows, "a disabled connector must not receive a buffered snapshot")
}

func TestWriter_DoesNotFlushDisabledGroupedConnector(t *testing.T) {
	db := openMigratedTestDB(t)
	connectors, mappings := NewSQLConnectorRepository(db), NewSQLTargetMappingRepository(db)
	path := filepath.Join(t.TempDir(), "disabled-group.db")
	target, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })
	_, err = target.ExecContext(t.Context(), `CREATE TABLE snapshot_values (d0 REAL)`)
	require.NoError(t, err)
	connector := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindSQLite, fmt.Sprintf(`{"dsn":%q}`, path), 15)
	group := "meter"
	require.NoError(t, mappings.Create(t.Context(), &schema.DatabaseTargetMapping{
		ID: "disabled-group", TagID: "disabled-group", ConnectorID: connector.ID, TableName: "snapshot_values",
		ColumnName: "d0", GroupKey: &group, WriteMode: schema.DatabaseWriteModeInsert,
		Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}))
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: make(chan time.Time)})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	at := time.Date(2026, 11, 4, 10, 0, 1, 0, time.UTC)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "disabled-group", 12, at), ErrSnapshotAccepted)
	connector.Enabled = false
	require.NoError(t, connectors.Update(t.Context(), connector))
	require.NoError(t, writer.RefreshMappings(t.Context()))
	require.NoError(t, writer.flushDue(t.Context(), at.Add(time.Minute), false))
	var rows int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&rows))
	require.Zero(t, rows, "a disabled connector must not receive a buffered grouped snapshot")
}

func TestWriter_RemovedMappingReleasesIntervalState(t *testing.T) {
	db := openMigratedTestDB(t)
	connectors, mappings := NewSQLConnectorRepository(db), NewSQLTargetMappingRepository(db)
	path := filepath.Join(t.TempDir(), "removed.db")
	target, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })
	_, err = target.ExecContext(t.Context(), `CREATE TABLE snapshot_values (d0 REAL)`)
	require.NoError(t, err)
	connector := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindSQLite, fmt.Sprintf(`{"dsn":%q}`, path), 15)
	createSnapshotMapping(t, mappings, connector.ID, "removed", "d0")
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: make(chan time.Time)})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	at := time.Date(2026, 11, 4, 10, 0, 1, 0, time.UTC)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "removed", 12, at), ErrSnapshotAccepted)
	require.NoError(t, writer.flushDue(t.Context(), at.Add(time.Minute), false))
	require.NoError(t, mappings.Delete(t.Context(), "mapping-removed"))
	require.NoError(t, writer.RefreshMappings(t.Context()))
	writer.mu.Lock()
	_, retained := writer.lastFlush["mapping-removed"]
	writer.mu.Unlock()
	require.False(t, retained, "removed mappings must not retain interval state indefinitely")
}

func TestWriter_GroupedValuesDoNotQueueWhileConnectorIsBlocked(t *testing.T) {
	db := openMigratedTestDB(t)
	db.SetMaxOpenConns(1)
	connectors, mappings := NewSQLConnectorRepository(db), NewSQLTargetMappingRepository(db)
	path := filepath.Join(t.TempDir(), "group-latest.db")
	target, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })
	_, err = target.ExecContext(t.Context(), `CREATE TABLE snapshot_values (d0 REAL)`)
	require.NoError(t, err)
	connector := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindSQLite, fmt.Sprintf(`{"dsn":%q}`, path), 15)
	group := "meter"
	require.NoError(t, mappings.Create(t.Context(), &schema.DatabaseTargetMapping{
		ID: "group-latest", TagID: "group-latest", ConnectorID: connector.ID, TableName: "snapshot_values",
		ColumnName: "d0", GroupKey: &group, WriteMode: schema.DatabaseWriteModeInsert,
		Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}))
	entered, release := make(chan struct{}), make(chan struct{})
	var once sync.Once
	releaseOpen := func() { once.Do(func() { close(release) }) }
	var first sync.Once
	tick := make(chan time.Time, 2)
	at := time.Date(2026, 11, 4, 10, 0, 1, 0, time.UTC)
	clock := &backoffTestClock{now: at}
	outcomes := &flushOutcomeSignalRepository{ConnectorRepository: connectors, completed: make(chan time.Time, 2)}
	writer := NewWriterWithConfig(outcomes, mappings, WriterConfig{FlushTick: tick, Now: clock.Now,
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			first.Do(func() { close(entered); <-release })
			return openExternalDBManagerFunc(kind, config)
		},
	})
	t.Cleanup(func() { releaseOpen(); _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "group-latest", 1, at), ErrSnapshotAccepted)
	clock.Set(at.Add(15 * time.Second))
	tick <- clock.Now()
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("connector did not start its first write")
	}
	for i := 1; i <= 10; i++ {
		require.ErrorIs(t, writer.WriteTagValue(t.Context(), "group-latest", i, at.Add(time.Duration(i)*15*time.Second)), ErrSnapshotAccepted)
	}
	clock.Set(at.Add(11 * 15 * time.Second))
	tick <- clock.Now()
	releaseOpen()
	waitCtx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()
	for i, want := range []time.Time{at.Truncate(15 * time.Second), at.Add(150 * time.Second).Truncate(15 * time.Second)} {
		select {
		case completed := <-outcomes.completed:
			require.Equal(t, want, completed, "flush %d must complete for its snapshot bucket", i+1)
		case <-waitCtx.Done():
			t.Fatalf("flush %d did not finish: %v", i+1, waitCtx.Err())
		}
	}
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
	require.Equal(t, []float64{1, 10}, values, "only the in-flight row and the latest grouped snapshot may be written")
}
