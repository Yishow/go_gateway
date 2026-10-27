package grouppipeline

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestSQLServerDestinationResolverNeverFallsBackWithoutProof(t *testing.T) {
	snap := mysqlRuntimeSnapshot(t)
	resolver, item, opened := newMySQLResolverFixture(t, snap, nil)
	opened.opened.Kind = schema.DatabaseConnectorKindSQLServer

	target, err := resolver.Resolve(t.Context(), item)

	require.ErrorIs(t, err, groupdelivery.ErrTargetBlocked)
	require.Nil(t, target.DB)
	require.Equal(t, 1, opened.closeCalls)
}

func TestSQLServerPipelineDialect(t *testing.T) {
	dialect, supported := dialectFor(schema.DatabaseConnectorKindSQLServer)
	require.True(t, supported)
	require.Equal(t, dbtarget.SQLDialectSQLServer, dialect)
}

func TestSQLServerPipelineFrozenProofAndBoundary(t *testing.T) {
	fixture := newFixture(t)
	fixture.dest.connector.Kind = schema.DatabaseConnectorKindSQLServer
	pipeline := fixture.pipeline()
	snap := sqlServerRuntimeSnapshot(t)
	require.NoError(t, pipeline.saveRuntimeSnapshot(t.Context(), snap))
	frozen, err := pipeline.loadRuntimeSnapshot(t.Context(), groupKey(snap))
	require.NoError(t, err)
	require.True(t, validLongRuntimeProof(frozen.Group, frozen.RuntimeLayout))
	built, reason := pipeline.build(t.Context(), snap)
	require.Empty(t, reason)
	require.NotNil(t, built)
	require.NoError(t, built.boundary.AcceptSample(t.Context(), envelope(snap.Group, "sample", at(1), 21.5)))
	require.NoError(t, built.boundary.Tick(t.Context(), at(15)))
	var effectKey string
	require.NoError(t, fixture.store.DB().QueryRowContext(t.Context(), `SELECT effect_key FROM wg_delivery_outbox`).Scan(&effectKey))
	item, err := fixture.store.GetOutbox(t.Context(), effectKey)
	require.NoError(t, err)
	require.True(t, sqlServerFrozenScopeMatches(item, frozen))
	// This is the real local journal/outbox boundary, not SQL Server readback.
	require.Equal(t, "record_id", item.Destination.RecordKeyColumn)
}

func TestSQLServerDestinationResolverUsesFrozenRevisionAndIndependentProof(t *testing.T) {
	snap := sqlServerRuntimeSnapshot(t)
	encoded, err := json.Marshal(snap)
	require.NoError(t, err)
	var restored workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(encoded, &restored))
	resolver, item, opened := sqlServerResolverFixture(t, &restored)
	target, err := resolver.Resolve(t.Context(), item)
	require.NoError(t, err)
	require.Equal(t, schema.DatabaseConnectorKindSQLServer, target.Kind)
	require.Equal(t, restored.RuntimeLayout.SQLServer, target.SQLServer)
	require.NotSame(t, restored.RuntimeLayout.SQLServer, target.SQLServer)
	target.SQLServer.Group.Evidence.PrimaryKey[0] = "edited"
	require.Equal(t, "bucket_start", restored.RuntimeLayout.SQLServer.Group.Evidence.PrimaryKey[0])
	require.Zero(t, opened.closeCalls)
	require.NoError(t, target.Close())
	require.Equal(t, 1, opened.closeCalls)
}

func TestSQLServerDestinationResolverRejectsDriftAndCloses(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(*workspace.WriteGroupAppliedSnapshot, *dbtarget.OpenedDestination, *groupdelivery.OutboxItem)
	}{
		{"gate closed", func(_ *workspace.WriteGroupAppliedSnapshot, o *dbtarget.OpenedDestination, _ *groupdelivery.OutboxItem) {
			o.SQLServer.SupportsManagedLong = false
		}},
		{"metadata hidden", func(_ *workspace.WriteGroupAppliedSnapshot, o *dbtarget.OpenedDestination, _ *groupdelivery.OutboxItem) {
			o.SQLServer.MetadataStatus = "unconfirmed"
		}},
		{"connector edited", func(_ *workspace.WriteGroupAppliedSnapshot, o *dbtarget.OpenedDestination, _ *groupdelivery.OutboxItem) {
			o.SQLServer.ConnectorRevision = "new"
		}},
		{"database changed", func(_ *workspace.WriteGroupAppliedSnapshot, o *dbtarget.OpenedDestination, _ *groupdelivery.OutboxItem) {
			o.SQLServer.Database = "other"
		}},
		{"schema changed", func(_ *workspace.WriteGroupAppliedSnapshot, o *dbtarget.OpenedDestination, _ *groupdelivery.OutboxItem) {
			o.SQLServer.Schema = "other"
		}},
		{"permission denied", func(_ *workspace.WriteGroupAppliedSnapshot, o *dbtarget.OpenedDestination, _ *groupdelivery.OutboxItem) {
			o.SQLServer.Permissions.InsertSchema = false
		}},
		{"capability changed", func(_ *workspace.WriteGroupAppliedSnapshot, o *dbtarget.OpenedDestination, _ *groupdelivery.OutboxItem) {
			o.SQLServer.CapabilityRevision = "other"
		}},
		{"version changed", func(_ *workspace.WriteGroupAppliedSnapshot, o *dbtarget.OpenedDestination, _ *groupdelivery.OutboxItem) {
			o.SQLServer.ServerVersion = "17.0.1"
		}},
		{"receipt proof missing", func(s *workspace.WriteGroupAppliedSnapshot, _ *dbtarget.OpenedDestination, _ *groupdelivery.OutboxItem) {
			s.RuntimeLayout.SQLServer.Receipt = nil
		}},
		{"proof absent", func(s *workspace.WriteGroupAppliedSnapshot, _ *dbtarget.OpenedDestination, _ *groupdelivery.OutboxItem) {
			s.RuntimeLayout.SQLServer = nil
		}},
		{"unverified indexes", func(s *workspace.WriteGroupAppliedSnapshot, _ *dbtarget.OpenedDestination, _ *groupdelivery.OutboxItem) {
			s.RuntimeLayout.SQLServer.Receipt.Indexes = nil
		}},
		{"table recreated", func(s *workspace.WriteGroupAppliedSnapshot, _ *dbtarget.OpenedDestination, _ *groupdelivery.OutboxItem) {
			s.RuntimeLayout.SQLServer.Group.Evidence.ObjectID++
		}},
		{"applied revision changed", func(s *workspace.WriteGroupAppliedSnapshot, _ *dbtarget.OpenedDestination, _ *groupdelivery.OutboxItem) {
			s.AppliedRevision = "other"
		}},
		{"outbox revision changed", func(_ *workspace.WriteGroupAppliedSnapshot, _ *dbtarget.OpenedDestination, i *groupdelivery.OutboxItem) {
			i.Key.GroupRevision = "other"
		}},
		{"outbox scope changed", func(_ *workspace.WriteGroupAppliedSnapshot, _ *dbtarget.OpenedDestination, i *groupdelivery.OutboxItem) {
			i.Destination.Scope = "other"
		}},
		{"native monthly", func(s *workspace.WriteGroupAppliedSnapshot, _ *dbtarget.OpenedDestination, _ *groupdelivery.OutboxItem) {
			s.Group.StorageMode = workspace.WriteGroupStorageModeNativeMonthly
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			snap := sqlServerRuntimeSnapshot(t)
			resolver, item, opened := sqlServerResolverFixture(t, snap)
			tc.edit(snap, opened.opened, &item)
			target, err := resolver.Resolve(t.Context(), item)
			require.ErrorIs(t, err, groupdelivery.ErrTargetBlocked)
			require.Nil(t, target.DB)
			require.Equal(t, 1, opened.closeCalls)
		})
	}
}

func TestSQLServerPipelineRejectsLegacyAndUnconfirmedSnapshots(t *testing.T) {
	for _, legacy := range []bool{true, false} {
		fixture := newFixture(t)
		fixture.dest.connector.Kind = schema.DatabaseConnectorKindSQLServer
		pipeline := fixture.pipeline()
		if legacy {
			built, reason := pipeline.build(t.Context(), snapshotOf(groupAt("rev-1"), at(0)))
			require.Nil(t, built)
			require.Equal(t, reasonDestinationUnsupported, reason)
			continue
		}
		snap := sqlServerRuntimeSnapshot(t)
		snap.Group.Destination.StorageStrategy = workspace.WriteGroupStorageStrategyCustom
		require.NoError(t, pipeline.saveRuntimeSnapshot(t.Context(), snap))
		built, reason := pipeline.build(t.Context(), snap)
		require.Nil(t, built)
		require.Equal(t, reasonSnapshotUnavailable, reason)
	}
}

func TestSQLServerRuntimeProofCannotRecoverWithoutCatalog(t *testing.T) {
	snap := sqlServerRuntimeSnapshot(t)
	snap.RuntimeLayout.SQLServer = nil
	require.False(t, validLongRuntimeProof(snap.Group, snap.RuntimeLayout))
	snap.RuntimeLayout.SQLServer = sqlServerRuntimeSnapshot(t).RuntimeLayout.SQLServer
	snap.EffectiveAt = time.Time{}
	item := mysqlOutboxItem(t, snap)
	require.False(t, sqlServerFrozenScopeMatches(item, snap))
}

func TestSQLServerDestinationResolverUnavailableSnapshotCloses(t *testing.T) {
	for _, scenario := range []string{"missing", "read failure", "no loader", "no connection"} {
		t.Run(scenario, func(t *testing.T) {
			snap := sqlServerRuntimeSnapshot(t)
			resolver, item, opened := sqlServerResolverFixture(t, snap)
			switch scenario {
			case "missing":
				resolver.loadRuntimeSnapshot = func(context.Context, groupdelivery.GroupKey) (*workspace.WriteGroupAppliedSnapshot, error) {
					return nil, nil
				}
			case "read failure":
				resolver.loadRuntimeSnapshot = func(context.Context, groupdelivery.GroupKey) (*workspace.WriteGroupAppliedSnapshot, error) {
					return nil, errors.New("private store error")
				}
			case "no loader":
				resolver.loadRuntimeSnapshot = nil
			case "no connection":
				opened.opened.DB = nil
			}
			_, err := resolver.Resolve(t.Context(), item)
			require.ErrorIs(t, err, groupdelivery.ErrTargetBlocked)
			require.NotContains(t, err.Error(), "private store error")
			require.Equal(t, 1, opened.closeCalls)
		})
	}
}
