package workspace

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func lifecycleReadFixture(t *testing.T, applied bool) (*sql.DB, *WriteGroupService, *WriteGroup) {
	t.Helper()
	db := openWorkspaceTestDB(t, ":memory:")
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	db.SetMaxOpenConns(1)
	service, _, created := newReadinessGroup(t.Context(), t, db)
	service.WithClock(func() time.Time { return time.Date(2026, 10, 2, 12, 0, 7, 0, time.UTC) })
	if applied {
		service.WithTableInspector(writeGroupLifecycleInspector{})
		result, err := service.Apply(t.Context(), created.Group.ID, lifecycleReadMutation(t, service, created.Group))
		require.NoError(t, err)
		created = result
	}
	return db, service, created.Group
}

func lifecycleReadMutation(t *testing.T, service *WriteGroupService, group *WriteGroup) WriteGroupMutation {
	t.Helper()
	record, err := service.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	return WriteGroupMutation{WorkspaceID: group.WorkspaceID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, ExpectedGroupRevision: group.Revision, ExpectedConnectorRevision: group.Destination.ConnectorRevision}
}

func lifecycleReadGroup(t *testing.T, service *WriteGroupService, id string) *WriteGroup {
	t.Helper()
	result, err := service.Get(t.Context(), id)
	require.NoError(t, err)
	return result.Group
}

func lifecycleActionNames(view *WriteGroupLifecycle) []string {
	var names []string
	for _, action := range view.Actions {
		names = append(names, action.Action)
	}
	return names
}

func lifecycleReadState(t *testing.T, db *sql.DB) map[string]string {
	t.Helper()
	state := map[string]string{}
	for _, table := range []string{"system_settings", "write_groups", "write_group_versions", "wg_runtime_versions", "database_connectors", "wg_delivery_samples", "wg_delivery_outbox", "wg_delivery_receipts", "wg_delivery_checkpoints"} {
		rows, err := db.Query(`SELECT * FROM ` + table + ` ORDER BY rowid`)
		require.NoError(t, err)
		columns, err := rows.Columns()
		require.NoError(t, err)
		all := [][]any{}
		for rows.Next() {
			values, pointers := make([]any, len(columns)), make([]any, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}
			require.NoError(t, rows.Scan(pointers...))
			all = append(all, values)
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
		raw, err := json.Marshal(all)
		require.NoError(t, err)
		state[table] = string(raw)
	}
	return state
}
