package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteGroupLifecycleHistoryReservedDraftDoesNotEnableIntake(t *testing.T) {
	f := lifecycleAPIApplied(t)
	id := f.group["id"].(string)
	f.db.SetMaxOpenConns(1)
	// Seed the already-established lineage/reservation storage contract. This
	// read test does not call private recreation or invent a public mutation.
	tx, err := f.db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()
	_, err = tx.Exec(`UPDATE write_groups SET status='deleted' WHERE id=?`, id)
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_groups(id,workspace_id,revision,name,status,destination_connector_id,destination_connector_revision,destination_table_name) VALUES('history-reserved-G2',?,'G2-draft','Reserved','draft','connector-A','connector-1','new_owned_table')`, f.record.ID)
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_group_lineages(lineage_id,workspace_id,current_group_id,reserved_group_id,reservation_operation_id,reservation_intent_digest,revision) VALUES('history-lineage',?,?,'history-reserved-G2','history-operation','history-intent',2)`, f.record.ID, id)
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_group_generations(group_id,lineage_id,generation) VALUES(?,'history-lineage',1)`, id)
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_group_generations(group_id,lineage_id,generation,predecessor_group_id) VALUES('history-reserved-G2','history-lineage',2,?)`, id)
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_group_basic_keys(workspace_id,device_id,canonical_role,group_id,lineage_id) VALUES(?,'device-A','basic-managed',?,'history-lineage')`, f.record.ID, id)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())
	before := lifecycleAPIState(t, f.db)
	for _, requested := range []string{id, "history-reserved-G2"} {
		data := lifecycleAPIRead(t, f.writeGroupRouterFixture, requested)
		require.Contains(t, data, "history")
		history := data["history"].(map[string]any)
		require.Equal(t, true, history["canonical"])
		require.Equal(t, id, history["current_group_id"])
		require.Equal(t, "history-reserved-G2", history["reserved_group_id"])
		generations := history["generations"].([]any)
		require.Len(t, generations, 2)
		reserved := generations[1].(map[string]any)
		require.Equal(t, "history-reserved-G2", reserved["group_id"])
		require.Equal(t, false, reserved["current"])
		require.Equal(t, true, reserved["reserved"])
		require.Equal(t, "draft", reserved["status"])
		require.Equal(t, "", reserved["applied_revision"])
		require.Empty(t, reserved["revisions"])
		if requested == "history-reserved-G2" {
			for _, action := range lifecycleAPIActions(t, data) {
				require.Equal(t, false, action["allowed"])
				require.Equal(t, "successor_reserved", action["reason_code"])
			}
		}
	}
	require.Equal(t, before, lifecycleAPIState(t, f.db))
	var current, reserved, operation, digest string
	var revision int
	require.NoError(t, f.db.QueryRow(`SELECT current_group_id,reserved_group_id,reservation_operation_id,reservation_intent_digest,revision FROM write_group_lineages WHERE lineage_id='history-lineage'`).Scan(&current, &reserved, &operation, &digest, &revision))
	require.Equal(t, id, current)
	require.Equal(t, "history-reserved-G2", reserved)
	require.Equal(t, "history-operation", operation)
	require.Equal(t, "history-intent", digest)
	require.Equal(t, 2, revision)
}
