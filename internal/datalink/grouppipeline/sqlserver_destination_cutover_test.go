package grouppipeline

import (
	"database/sql"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// TestSQLServerDestinationCutover proves that accepted effects stay bound to
// the SQL Server destination captured before a boundary, while post-boundary
// effects resolve against the newly applied destination. The catalog and SQL
// Server connections are scripted; the durable gateway journal is local SQLite.
func TestSQLServerDestinationCutover(t *testing.T) {
	fixture := newFixture(t)
	pipeline := fixture.pipeline()
	base := sqlServerRuntimeSnapshot(t)
	old := sqlServerCutoverSnapshot(t, base, "sqlserver-rev-a", "connector-1", "crev-1", at(0))

	// Close and freeze the old accepted effect before the new descriptor exists.
	// Ending the old boundary is a separate runtime fact and never rewrites its
	// immutable descriptor or pending payload.
	require.NoError(t, pipeline.saveRuntimeSnapshot(t.Context(), old))
	oldBoundarySnapshot := *old
	oldBoundarySnapshot.EndsAt = at(20)
	oldBoundary, reason := pipeline.buildFrozen(t.Context(), &oldBoundarySnapshot, time.Time{})
	require.Empty(t, reason)
	require.NotNil(t, oldBoundary)
	fixture.clock = at(25)
	require.NoError(t, oldBoundary.boundary.AcceptSample(t.Context(), envelope(old.Group, "before-cutover", at(19), 19.5)))
	require.NoError(t, oldBoundary.boundary.Tick(t.Context(), at(30)))
	frozenOld := sqlServerFrozenOutboxForRevision(t, fixture.store.DB(), old.Group.Revision)
	require.NotEmpty(t, frozenOld)

	next := sqlServerCutoverSnapshot(t, base, "sqlserver-rev-b", "connector-b", "crev-b", at(20))
	require.NoError(t, pipeline.saveRuntimeSnapshot(t.Context(), next))
	nextBoundary, reason := pipeline.buildFrozen(t.Context(), next, time.Time{})
	require.Empty(t, reason)
	require.NotNil(t, nextBoundary)
	require.NoError(t, nextBoundary.boundary.AcceptSample(t.Context(), envelope(next.Group, "at-cutover", at(20), 20.5)))
	require.NoError(t, nextBoundary.boundary.Tick(t.Context(), at(30)))
	require.Equal(t, frozenOld, sqlServerFrozenOutboxForRevision(t, fixture.store.DB(), old.Group.Revision),
		"adding a new destination descriptor does not rewrite accepted old payload or destination bytes")

	rows, err := fixture.store.DB().QueryContext(t.Context(), `SELECT effect_key FROM wg_delivery_outbox ORDER BY group_revision`)
	require.NoError(t, err)
	var keys []string
	for rows.Next() {
		var key string
		require.NoError(t, rows.Scan(&key))
		keys = append(keys, key)
	}
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
	require.Len(t, keys, 3, "the old partial boundary can close its preceding bucket as well")
	var oldItems, nextItems []groupdelivery.OutboxItem
	for _, key := range keys {
		item, err := fixture.store.GetOutbox(t.Context(), key)
		require.NoError(t, err)
		switch item.Key.GroupRevision {
		case old.Group.Revision:
			require.Equal(t, old.Group.Destination.ConnectorID, item.Destination.ConnectorID)
			require.Equal(t, old.Group.Destination.ConnectorRevision, item.Destination.ConnectorRevision)
			require.Equal(t, old.Group.Destination.TableName, item.Destination.TableName)
			oldItems = append(oldItems, item)
		case next.Group.Revision:
			require.Equal(t, next.Group.Destination.ConnectorID, item.Destination.ConnectorID)
			require.Equal(t, next.Group.Destination.ConnectorRevision, item.Destination.ConnectorRevision)
			require.Equal(t, next.Group.Destination.TableName, item.Destination.TableName)
			nextItems = append(nextItems, item)
		default:
			require.FailNow(t, "unexpected outbox revision", item.Key.GroupRevision)
		}
	}
	require.NotEmpty(t, oldItems)
	require.Len(t, nextItems, 1)
	require.NotEqual(t, oldItems[0].Destination.Scope, nextItems[0].Destination.Scope)
	require.Len(t, frozenOld, len(oldItems))

	// Each outbox row resolves against its own immutable runtime version even
	// though a newer destination revision now exists for the same group.
	for _, item := range oldItems {
		oldResolver, _, opened := sqlServerResolverFixture(t, old)
		oldTarget, err := oldResolver.Resolve(t.Context(), item)
		require.NoError(t, err)
		require.Equal(t, schema.DatabaseConnectorKindSQLServer, oldTarget.Kind)
		require.Equal(t, old.Group.Destination.TableName, oldTarget.SQLServer.Group.Inspection.Table)
		require.NoError(t, oldTarget.Close())
		require.Equal(t, 1, opened.closeCalls)
	}
	nextResolver, _, opened := sqlServerResolverFixture(t, next)
	nextTarget, err := nextResolver.Resolve(t.Context(), nextItems[0])
	require.NoError(t, err)
	require.Equal(t, schema.DatabaseConnectorKindSQLServer, nextTarget.Kind)
	require.Equal(t, next.Group.Destination.TableName, nextTarget.SQLServer.Group.Inspection.Table)
	require.NoError(t, nextTarget.Close())
	require.Equal(t, 1, opened.closeCalls)
}

type sqlServerFrozenOutbox struct {
	Key, Payload, Digest, ConnectorID, ConnectorRevision, TableSchema, TableName string
}

func sqlServerFrozenOutboxForRevision(t *testing.T, db *sql.DB, revision string) []sqlServerFrozenOutbox {
	t.Helper()
	rows, err := db.QueryContext(t.Context(), `SELECT effect_key, payload, payload_digest, connector_id,
		connector_revision, table_schema, table_name FROM wg_delivery_outbox WHERE group_revision = ? ORDER BY effect_key`, revision)
	require.NoError(t, err)
	defer rows.Close()
	var out []sqlServerFrozenOutbox
	for rows.Next() {
		var item sqlServerFrozenOutbox
		require.NoError(t, rows.Scan(&item.Key, &item.Payload, &item.Digest, &item.ConnectorID,
			&item.ConnectorRevision, &item.TableSchema, &item.TableName))
		out = append(out, item)
	}
	require.NoError(t, rows.Err())
	return out
}

// sqlServerCutoverSnapshot builds one new immutable scope using the existing
// scripted catalog fixture, not a live SQL Server or a copied catalog builder.
func sqlServerCutoverSnapshot(
	t *testing.T,
	base *workspace.WriteGroupAppliedSnapshot,
	revision, connectorID, connectorRevision string,
	effectiveAt time.Time,
) *workspace.WriteGroupAppliedSnapshot {
	t.Helper()
	group := *base.Group
	group.Revision, group.AppliedRevision = revision, revision
	group.Destination.ConnectorID, group.Destination.ConnectorRevision = connectorID, connectorRevision
	group.Destination.SchemaRevision = "confirmed-" + revision
	table, owner, ok := workspace.LongStorageIdentity(&group)
	require.True(t, ok)
	group.Destination.TableName = table

	prior := base.RuntimeLayout.SQLServer.Group.Evidence
	options := recordingplan.SQLServerGroupSchemaOptions{
		Database: group.Destination.Database, Schema: group.Destination.TableSchema,
		ServerVersion: prior.ServerVersion, ServerEdition: prior.ServerEdition,
		EngineEdition: prior.EngineEdition, CompatibilityLevel: prior.CompatibilityLevel,
		CapabilityRevision: "cap-" + connectorID, Collation: prior.ColumnCollations["record_id"],
	}
	layout, err := recordingplan.ManagedSQLServerLongGroupLayout(table, owner, options)
	require.NoError(t, err)
	groupProof, groupState := sqlServerRuntimeCatalog(layout.Columns, layout.Indexes, table, owner, 201, options)
	receiptColumns := []recordingplan.GroupSchemaColumn{
		{Name: "effect_key", SQLType: "CHAR(64)", PrimaryKey: true},
		{Name: "payload_digest", SQLType: "CHAR(64)"},
		{Name: "committed_at", SQLType: "DATETIME2(6)"},
		{Name: "_gw_effect_receipts_v1", SQLType: "NVARCHAR(16)"},
	}
	receiptProof, receiptState := sqlServerRuntimeCatalog(receiptColumns, nil, dbtarget.EffectReceiptTable, "_gw_effect_receipts_v1", 202, options)
	digest, err := recordingplan.VerifyManagedSQLServerCatalog(layout, group.Destination.Database, groupState, receiptState)
	require.NoError(t, err)
	group.Destination.SchemaDigest = digest

	frozen := &workspace.WriteGroupRuntimeLayout{
		Dialect: dbtarget.SQLDialectSQLServer, Columns: groupProof.Inspection.Columns,
		TagTypes: base.RuntimeLayout.TagTypes, ReceiptTableReady: true, SchemaDigest: digest,
		SQLServer: &dbtarget.SQLServerManagedRowProof{
			ConnectorID: connectorID, ConnectorRevision: connectorRevision, Group: groupProof, Receipt: receiptProof,
		},
	}
	return &workspace.WriteGroupAppliedSnapshot{
		WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: revision, AppliedRevision: revision,
		EffectiveAt: effectiveAt, Group: &group, RuntimeLayout: frozen,
	}
}
