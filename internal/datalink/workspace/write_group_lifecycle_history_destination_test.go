package workspace

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/snapshot"
)

func lifecycleHistoryRecordKeyFixture(t *testing.T, dialect, key string, withOutbox bool) (*sql.DB, *WriteGroupService, *WriteGroup) {
	t.Helper()
	db, service, group := lifecycleReadFixture(t, true)
	group.RowPolicy.RecordKeyColumn = "record_id"
	group.WritePolicy.DedupeCapability = "none"
	raw, err := json.Marshal(group)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE write_group_versions SET payload=? WHERE group_id=?`, string(raw), group.ID)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE wg_runtime_versions SET payload=json_set(payload,'$.group.row_policy.record_key_column','record_id','$.runtime_layout.dialect',?) WHERE group_id=?`, dialect, group.ID)
	require.NoError(t, err)
	if withOutbox {
		d := group.Destination
		scope, err := snapshot.DestinationScope(d.ConnectorID, d.ConnectorRevision, d.Database, d.TableSchema, d.TableName)
		require.NoError(t, err)
		_, err = db.Exec(`INSERT INTO wg_delivery_outbox(effect_key,record_id,workspace_id,group_id,group_revision,bucket_start,partition_key,destination_scope,connector_id,connector_revision,database_name,table_schema,table_name,dedupe_capability,record_key_column,payload,payload_digest,next_retry_at,created_at,updated_at) VALUES('effect','record',?,?,?,'2026-10-02T12:00:00Z','partition',?,?,?,?,?,?,?,?,'{}','digest','2026-10-02T12:00:00Z','2026-10-02T12:00:00Z','2026-10-02T12:00:00Z')`, group.WorkspaceID, group.ID, group.AppliedRevision, scope, d.ConnectorID, d.ConnectorRevision, d.Database, d.TableSchema, d.TableName, group.WritePolicy.DedupeCapability, key)
		require.NoError(t, err)
	}
	return db, service, group
}

func TestLifecycleHistoryRecordKeyUsesFrozenRuntimeDialect(t *testing.T) {
	for _, dialect := range []string{"sqlite", "postgres", "mysql", "sqlserver"} {
		t.Run(dialect, func(t *testing.T) {
			key := ""
			if dialect == "mysql" || dialect == "sqlserver" {
				key = "record_id"
			}
			db, service, group := lifecycleHistoryRecordKeyFixture(t, dialect, key, true)
			_, err := db.Exec(`UPDATE database_connectors SET kind='unknown',identity_revision='changed' WHERE id=?`, group.Destination.ConnectorID)
			require.NoError(t, err)
			before := lifecycleReadState(t, db)
			view, err := service.Lifecycle(t.Context(), group.ID)
			require.NoError(t, err, "immutable runtime destination, not live connector or raw row policy, defines outbox key")
			require.NotNil(t, view)
			require.Equal(t, before, lifecycleReadState(t, db))
		})
	}
}

func TestLifecycleHistoryRecordKeyRejectsMissingOrContradictoryEvidence(t *testing.T) {
	for _, mode := range []string{"missing", "malformed", "foreign-row", "duplicate-row", "foreign-root", "foreign-nested", "wrong-destination", "wrong-key", "unknown-dialect"} {
		t.Run(mode, func(t *testing.T) {
			db, service, group := lifecycleHistoryRecordKeyFixture(t, "mysql", "record_id", true)
			query := map[string]string{
				"missing":           `DELETE FROM wg_runtime_versions`,
				"malformed":         `UPDATE wg_runtime_versions SET payload='broken'`,
				"foreign-row":       `UPDATE wg_runtime_versions SET workspace_id='foreign'`,
				"duplicate-row":     `INSERT INTO wg_runtime_versions(workspace_id,group_id,group_revision,payload) SELECT 'foreign',group_id,group_revision,payload FROM wg_runtime_versions`,
				"foreign-root":      `UPDATE wg_runtime_versions SET payload=json_set(payload,'$.workspace_id','foreign')`,
				"foreign-nested":    `UPDATE wg_runtime_versions SET payload=json_set(payload,'$.group.workspace_id','foreign')`,
				"wrong-destination": `UPDATE wg_runtime_versions SET payload=json_set(payload,'$.group.destination.table_name','other')`,
				"wrong-key":         `UPDATE wg_runtime_versions SET payload=json_set(payload,'$.group.row_policy.record_key_column','other')`,
				"unknown-dialect":   `UPDATE wg_runtime_versions SET payload=json_set(payload,'$.runtime_layout.dialect','unknown')`,
			}[mode]
			_, err := db.Exec(query)
			require.NoError(t, err)
			before := lifecycleReadState(t, db)
			view, err := service.Lifecycle(t.Context(), group.ID)
			require.Error(t, err)
			require.Nil(t, view)
			require.Equal(t, before, lifecycleReadState(t, db))
		})
	}
}

func TestLifecycleHistoryWithoutOutboxDoesNotRequireUnusedRecordKeyEvidence(t *testing.T) {
	db, service, group := lifecycleHistoryRecordKeyFixture(t, "postgres", "", false)
	_, err := db.Exec(`DELETE FROM wg_runtime_versions`)
	require.NoError(t, err)
	view, err := service.Lifecycle(t.Context(), group.ID)
	require.NoError(t, err)
	require.NotNil(t, view)
}
