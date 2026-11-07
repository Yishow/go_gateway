package workspace

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermanentDeleteStorageScopePreservesSurvivorAndRemovesUnshared(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	h.db.SetMaxOpenConns(1)
	power := h.primary()
	energy := h.seedSource(power.deviceID, "energy")
	first := h.mustEnsure(h.longDraft(WriteGroupProfileInstantaneous, "main", power))
	second := h.mustEnsure(h.longDraft(WriteGroupProfileCumulative, "main", energy))
	scope := longStorageScopeOf(first)
	for _, table := range []string{"write_group_mysql_partition_maintenance", "write_group_sqlserver_partition_maintenance"} {
		_, err := h.db.ExecContext(t.Context(), `INSERT INTO `+table+` (workspace_id,connector_id,database_name,table_schema,storage_mode,layout_version,base_layout_digest,connector_revision,generation,current_evidence,revision,status,updated_at) VALUES(?,?,?,?,?,?,'digest','connector-revision-1',1,'{}',1,'ready',CURRENT_TIMESTAMP)`, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.Mode, scope.Layout)
		require.NoError(t, err)
	}
	current, err := h.service.Get(t.Context(), first.ID)
	require.NoError(t, err)
	_, err = h.service.Delete(t.Context(), first.ID, h.mutationFor(current))
	require.NoError(t, err)
	for _, table := range []string{"write_group_storage_scopes", "write_group_mysql_partition_maintenance", "write_group_sqlserver_partition_maintenance"} {
		require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM `+table), table)
	}
	current, err = h.service.Get(t.Context(), second.ID)
	require.NoError(t, err)
	_, err = h.service.Delete(t.Context(), second.ID, h.mutationFor(current))
	require.NoError(t, err)
	for _, table := range []string{"write_group_storage_scopes", "write_group_mysql_partition_maintenance", "write_group_sqlserver_partition_maintenance"} {
		require.Zero(t, h.countRows(`SELECT COUNT(*) FROM `+table), table)
	}
}

func TestPermanentDeleteRemovesUnsharedHistoricalStorageScopes(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	h.db.SetMaxOpenConns(1)
	group := h.mustEnsure(h.longDraft(WriteGroupProfileInstantaneous, "main", h.primary()))
	historical := cloneWriteGroup(group)
	historical.Destination.Database = "old-database"
	payload, err := json.Marshal(historical)
	require.NoError(t, err)
	_, err = h.db.ExecContext(t.Context(), `INSERT INTO write_group_versions(group_id,group_revision,effective_at,payload) VALUES(?,'old-revision','2099-01-01T00:00:00Z',?)`, group.ID, string(payload))
	require.NoError(t, err)
	_, err = h.db.ExecContext(t.Context(), `INSERT INTO write_group_storage_scopes(workspace_id,connector_id,database_name,table_schema,storage_mode,layout_version,table_name,owner_column,layout_digest)
		SELECT workspace_id,connector_id,'old-database',table_schema,storage_mode,layout_version,table_name,owner_column,layout_digest FROM write_group_storage_scopes`)
	require.NoError(t, err)
	current, err := h.service.Get(t.Context(), group.ID)
	require.NoError(t, err)
	_, err = h.service.Delete(t.Context(), group.ID, h.mutationFor(current))
	require.NoError(t, err)
	require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_group_storage_scopes`), "the historical scope also belongs solely to the deleted recording")
}

func TestPermanentDeleteRetainsScopeUsedBySurvivingHistory(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	h.db.SetMaxOpenConns(1)
	power := h.primary()
	first := h.mustEnsure(h.longDraft(WriteGroupProfileInstantaneous, "main", power))
	second := h.mustEnsure(h.longDraft(WriteGroupProfileCumulative, "main", h.seedSource(power.deviceID, "energy")))
	oldVersion, err := json.Marshal(second)
	require.NoError(t, err)
	_, err = h.db.ExecContext(t.Context(), `INSERT INTO write_group_versions(group_id,group_revision,effective_at,payload) VALUES(?,'old-revision','2099-01-01T00:00:00Z',?)`, second.ID, string(oldVersion))
	require.NoError(t, err)
	_, err = h.db.ExecContext(t.Context(), `UPDATE write_groups SET destination_database='new-database' WHERE id=?`, second.ID)
	require.NoError(t, err)
	current, err := h.service.Get(t.Context(), first.ID)
	require.NoError(t, err)
	_, err = h.service.Delete(t.Context(), first.ID, h.mutationFor(current))
	require.NoError(t, err)
	require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_group_storage_scopes`), "surviving historical versions retain the scope")
	current, err = h.service.Get(t.Context(), second.ID)
	require.NoError(t, err)
	_, err = h.service.Delete(t.Context(), second.ID, h.mutationFor(current))
	require.NoError(t, err)
	require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_group_storage_scopes`))
}
