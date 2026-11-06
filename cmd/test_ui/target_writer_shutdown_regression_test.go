package main

import (
	"context"
	"database/sql"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/apphost"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

func TestProductionShutdownWaitsForDatabaseTargetBeforeClosingSharedStorage(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "configuration.db"))
	require.NoError(t, err)
	datalink.ApplySQLitePoolDefaults(db)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	connectors, mappings := dbtarget.NewSQLConnectorRepository(db), dbtarget.NewSQLTargetMappingRepository(db)
	now := time.Now().UTC()
	require.NoError(t, connectors.Create(t.Context(), &schema.DatabaseConnector{
		ID: "shutdown-target", Name: "shutdown-target", Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: `{"dsn":"unused.db"}`, Enabled: true, DefaultWriteIntervalSeconds: 15,
		CreatedAt: now, UpdatedAt: now,
	}))
	require.NoError(t, mappings.Create(t.Context(), &schema.DatabaseTargetMapping{
		ID: "shutdown-mapping", TagID: "shutdown-tag", ConnectorID: "shutdown-target",
		TableName: "snapshot_values", ColumnName: "value", WriteMode: schema.DatabaseWriteModeInsert,
		Enabled: true, CreatedAt: now, UpdatedAt: now,
	}))
	entered, release, checked := make(chan struct{}), make(chan struct{}), make(chan error, 1)
	unblock := sync.OnceFunc(func() { close(release) })
	var checkOnce sync.Once
	tick := make(chan time.Time, 1)
	target := dbtarget.NewWriterWithConfig(connectors, mappings, dbtarget.WriterConfig{
		FlushTick: tick,
		SuppressMapping: func(string, string) bool {
			checkOnce.Do(func() {
				close(entered)
				<-release
				checked <- db.PingContext(context.Background())
			})
			return true
		},
	})
	t.Cleanup(func() { unblock(); _ = target.Close(context.Background()) })
	productionTarget := newProductionTargetWriter(target, nil)
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	groups := pollinggroup.NewMemoryRepository()
	service, err := datalinkruntime.NewService(datalinkruntime.DefaultConfig(), datalinkruntime.Dependencies{
		Scheduler: collector.NewScheduler(collector.DefaultSchedulerConfig(), manager),
		Writer:    storage.NewSQLiteWriter(db), TargetWriter: productionTarget,
		DeviceService:  device.NewService(device.NewMemoryRepository(), manager),
		PointService:   point.NewService(point.NewMemoryRepository(), groups),
		MappingService: mapping.NewService(mapping.NewMemoryRepository()),
		TagService:     tag.NewService(tag.NewMemoryRepository()), PollingGroupService: pollinggroup.NewService(groups),
	})
	require.NoError(t, err)
	require.ErrorIs(t, productionTarget.WriteTagValue(t.Context(), "shutdown-tag", 12, now), dbtarget.ErrSnapshotAccepted)
	tick <- now.Add(time.Minute)
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("database target did not enter its ownership check")
	}
	stopResult := make(chan error, 1)
	storageClosed := make(chan struct{})
	shutdown := apphost.NewShutdown(50*time.Millisecond, []apphost.Phase{
		{Name: "runtime", Stop: func(ctx context.Context) error {
			err := service.Stop(ctx)
			stopResult <- err
			return err
		}, Wait: service.WaitStopped},
		{Name: "database", Stop: func(context.Context) error {
			err := db.Close()
			close(storageClosed)
			return err
		}},
	}, nil)
	t.Cleanup(func() { unblock(); _ = shutdown.Wait(context.Background()) })
	shutdown.Begin()
	select {
	case err := <-stopResult:
		require.ErrorIs(t, err, context.DeadlineExceeded)
	case <-time.After(time.Second):
		t.Fatal("runtime Stop did not honor the caller deadline")
	}
	select {
	case <-storageClosed:
		t.Fatal("shared database closed while the target worker was still active")
	case <-time.After(100 * time.Millisecond):
	}
	unblock()
	require.NoError(t, <-checked, "the target worker must finish while the shared database is open")
	require.ErrorIs(t, shutdown.Wait(t.Context()), context.DeadlineExceeded)
	select {
	case <-storageClosed:
	default:
		t.Fatal("shared database was not closed after target completion")
	}
}
