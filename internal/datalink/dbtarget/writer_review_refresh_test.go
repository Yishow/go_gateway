package dbtarget

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type reviewCountingConnectors struct {
	ConnectorRepository
	gets atomic.Int64
}

type reviewHeldRefresh struct {
	TargetMappingRepository
	hold    atomic.Bool
	once    sync.Once
	entered chan struct{}
	release chan struct{}
}

func (r *reviewHeldRefresh) List(ctx context.Context, filter TargetMappingListFilter) ([]*schema.DatabaseTargetMapping, error) {
	if r.hold.Load() {
		r.once.Do(func() { close(r.entered) })
		select {
		case <-r.release:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return r.TargetMappingRepository.List(ctx, filter)
}

func (r *reviewCountingConnectors) GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error) {
	r.gets.Add(1)
	return r.ConnectorRepository.GetByID(ctx, id)
}

type reviewBatchTags struct {
	tags    map[string]*schema.Tag
	gets    atomic.Int64
	lists   atomic.Int64
	listErr error
}

func (r *reviewBatchTags) GetByID(_ context.Context, id string) (*schema.Tag, error) {
	r.gets.Add(1)
	if entity := r.tags[id]; entity != nil {
		return entity, nil
	}
	return nil, tag.ErrTagNotFound
}

func (r *reviewBatchTags) List(context.Context, tag.ListFilter) ([]*schema.Tag, error) {
	r.lists.Add(1)
	if r.listErr != nil {
		return nil, r.listErr
	}
	entities := make([]*schema.Tag, 0, len(r.tags))
	for _, entity := range r.tags {
		entities = append(entities, entity)
	}
	return entities, nil
}

func TestWriter_TagBatchRefreshFailureRejectsIngestUntilRecovery(t *testing.T) {
	connectors, mappings, connector, _ := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "batch-failure", "d0")
	reader := &reviewBatchTags{tags: map[string]*schema.Tag{"batch-failure": {ID: "batch-failure"}}}
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{TagReader: reader, FlushTick: make(chan time.Time)})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	reader.listErr = fmt.Errorf("tag store unavailable")
	require.ErrorContains(t, writer.RefreshMappings(t.Context()), "tag store unavailable")
	accepted, err := writer.AcceptTagValue(t.Context(), "batch-failure", 12, time.Now())
	require.ErrorContains(t, err, "tag store unavailable")
	require.False(t, accepted)
	reader.listErr = nil
	require.NoError(t, writer.RefreshMappings(t.Context()))
	accepted, err = writer.AcceptTagValue(t.Context(), "batch-failure", 12, time.Now())
	require.NoError(t, err)
	require.True(t, accepted)
}

func TestWriter_RefreshReusesConnectorListAndBatchesTags(t *testing.T) {
	connectors, mappings, connector, _ := backoffTestTarget(t)
	reader := &reviewBatchTags{tags: make(map[string]*schema.Tag)}
	for i := range 20 {
		id := fmt.Sprintf("tag-%02d", i)
		reader.tags[id] = &schema.Tag{ID: id}
		createSnapshotMapping(t, mappings, connector.ID, id, "d0")
	}
	repo := &reviewCountingConnectors{ConnectorRepository: connectors}
	writer := NewWriterWithConfig(repo, mappings, WriterConfig{TagReader: reader, FlushTick: make(chan time.Time)})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	repo.gets.Store(0)
	reader.gets.Store(0)
	reader.lists.Store(0)
	require.NoError(t, writer.RefreshMappings(t.Context()))
	assert.Zero(t, repo.gets.Load(), "listed connectors must not be reloaded individually")
	assert.EqualValues(t, 1, reader.lists.Load(), "tags should be loaded in one batch")
	assert.Zero(t, reader.gets.Load(), "batch-capable tag readers must not be called per tag")
	accepted, err := writer.AcceptTagValue(t.Context(), "tag-00", 12, time.Now())
	require.NoError(t, err)
	require.True(t, accepted)
}

func TestWriter_SlowRefreshDoesNotBlockDueSnapshot(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "pending", "d0")
	repo := &reviewHeldRefresh{TargetMappingRepository: mappings, entered: make(chan struct{}), release: make(chan struct{})}
	var releaseOnce sync.Once
	releaseRefresh := func() { releaseOnce.Do(func() { close(repo.release) }) }
	base := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	opened := make(chan struct{}, 1)
	writer := NewWriterWithConfig(connectors, repo, WriterConfig{FlushTick: tick, Now: clock.Now,
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			opened <- struct{}{}
			return openExternalDBManagerFunc(kind, config)
		},
	})
	t.Cleanup(func() { releaseRefresh(); _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "pending", 12, base.Add(time.Second)), ErrSnapshotAccepted)
	repo.hold.Store(true)
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()
	select {
	case <-repo.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("refresh did not start")
	}
	select {
	case <-opened:
	case <-time.After(2 * time.Second):
		t.Fatal("a blocked refresh prevented the due connector worker from waking")
	}
	require.Eventually(t, func() bool {
		var value float64
		return target.QueryRowContext(t.Context(), `SELECT d0 FROM snapshot_values`).Scan(&value) == nil && value == 12
	}, 2*time.Second, 10*time.Millisecond)
}

func TestWriter_SlowRefreshDoesNotDeliverDisabledMapping(t *testing.T) {
	checkSlowRefreshDoesNotDeliverDisabledMapping(t, false)
}

func TestWriter_SlowRefreshDoesNotDeliverDisabledGroupedMapping(t *testing.T) {
	checkSlowRefreshDoesNotDeliverDisabledMapping(t, true)
}

func TestWriter_SlowExplicitRefreshDoesNotDeliverDisabledMapping(t *testing.T) {
	connectors, mappings, connector, target := backoffTestTarget(t)
	createSnapshotMapping(t, mappings, connector.ID, "first", "d0")
	createSnapshotMapping(t, mappings, connector.ID, "second", "d1")
	repo := &reviewHeldRefresh{TargetMappingRepository: mappings, entered: make(chan struct{}), release: make(chan struct{})}
	var refreshOnce, openOnce sync.Once
	releaseRefresh := func() { refreshOnce.Do(func() { close(repo.release) }) }
	releaseOpen := make(chan struct{})
	stopOpen := func() { openOnce.Do(func() { close(releaseOpen) }) }
	opened := make(chan struct{})
	var attempts atomic.Int64
	base := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, repo, WriterConfig{FlushTick: tick, Now: clock.Now,
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			if attempts.Add(1) == 1 {
				close(opened)
				<-releaseOpen
			}
			return openExternalDBManagerFunc(kind, config)
		},
	})
	t.Cleanup(func() { releaseRefresh(); stopOpen(); _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "first", 1, base.Add(time.Second)), ErrSnapshotAccepted)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "second", 2, base.Add(time.Second)), ErrSnapshotAccepted)
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()
	select {
	case <-opened:
	case <-time.After(2 * time.Second):
		t.Fatal("first snapshot did not start")
	}
	require.Eventually(t, func() bool { return !writer.refreshing.Load() }, 2*time.Second, 10*time.Millisecond)
	entry, err := mappings.GetByID(t.Context(), "mapping-second")
	require.NoError(t, err)
	entry.Enabled = false
	require.NoError(t, mappings.Update(t.Context(), entry))
	repo.hold.Store(true)
	done := make(chan error, 1)
	go func() { done <- writer.RefreshMappings(context.Background()) }()
	select {
	case <-repo.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("explicit refresh did not start")
	}
	stopOpen()
	require.Eventually(t, func() bool {
		for _, record := range listWriteHistory(connector.ID, 20) {
			if record.MappingID == entry.ID && record.SkipReason == "mapping_not_current" {
				return true
			}
		}
		return false
	}, 2*time.Second, 10*time.Millisecond, "explicit refresh must not let the worker deliver a disabled mapping")
	var rows int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values WHERE d1 IS NOT NULL`).Scan(&rows))
	require.Zero(t, rows)
	releaseRefresh()
	require.NoError(t, <-done)
}

func checkSlowRefreshDoesNotDeliverDisabledMapping(t *testing.T, grouped bool) {
	t.Helper()
	connectors, mappings, connector, target := backoffTestTarget(t)
	if grouped {
		group := "meter"
		require.NoError(t, mappings.Create(t.Context(), &schema.DatabaseTargetMapping{
			ID: "mapping-disabled-during-refresh", TagID: "disabled-during-refresh", ConnectorID: connector.ID,
			TableName: "snapshot_values", ColumnName: "d0", GroupKey: &group,
			WriteMode: schema.DatabaseWriteModeInsert, Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}))
	} else {
		createSnapshotMapping(t, mappings, connector.ID, "disabled-during-refresh", "d0")
	}
	repo := &reviewHeldRefresh{TargetMappingRepository: mappings, entered: make(chan struct{}), release: make(chan struct{})}
	var releaseOnce sync.Once
	releaseRefresh := func() { releaseOnce.Do(func() { close(repo.release) }) }
	base := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, repo, WriterConfig{FlushTick: tick, Now: clock.Now})
	t.Cleanup(func() { releaseRefresh(); _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "disabled-during-refresh", 12, base.Add(time.Second)), ErrSnapshotAccepted)
	entry, err := mappings.GetByID(t.Context(), "mapping-disabled-during-refresh")
	require.NoError(t, err)
	entry.Enabled = false
	require.NoError(t, mappings.Update(t.Context(), entry))
	repo.hold.Store(true)
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()
	select {
	case <-repo.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("refresh did not start")
	}
	require.Eventually(t, func() bool {
		for _, record := range listWriteHistory(connector.ID, 20) {
			if record.MappingID == entry.ID && record.SkipReason == "mapping_not_current" {
				return true
			}
		}
		return false
	}, 2*time.Second, 10*time.Millisecond, "a concurrently disabled mapping must be rejected even while refresh is blocked")
	var rows int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&rows))
	require.Zero(t, rows)
}
