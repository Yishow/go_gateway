package workspace

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func lifecycleHistoryJSON(t *testing.T, service *WriteGroupService, id string) map[string]any {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()
	view, err := service.Lifecycle(ctx, id)
	require.NoError(t, err)
	raw, err := json.Marshal(view)
	require.NoError(t, err)
	var decoded map[string]any
	require.NoError(t, json.Unmarshal(raw, &decoded))
	history, ok := decoded["history"].(map[string]any)
	require.True(t, ok, "lifecycle must return separately scoped historical evidence")
	return history
}

func TestLifecycleHistoryPreservesImmutableDestinationsAndLegacyIdentity(t *testing.T) {
	db, service, original := lifecycleReadFixture(t, true)
	second := cloneWriteGroup(original)
	second.Revision, second.AppliedRevision = "R-historical-2", "R-historical-2"
	second.Destination.ConnectorID, second.Destination.ConnectorRevision = "historical-connector", "historical-C2"
	second.Destination.Database, second.Destination.TableName = "/private/secret.sqlite", "private_table"
	second.Members[0].TagID = "historical-tag"
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	require.NoError(t, service.repo.insertVersionInTx(t.Context(), tx, second, time.Date(2026, 10, 2, 13, 0, 0, 0, time.UTC)))
	require.NoError(t, service.repo.insertRevisionEndInTx(t.Context(), tx, original.WorkspaceID, original.ID, original.AppliedRevision, "C-end", time.Date(2026, 10, 2, 13, 0, 0, 0, time.UTC)))
	require.NoError(t, tx.Commit())
	before := lifecycleReadState(t, db)
	history := lifecycleHistoryJSON(t, service, original.ID)
	require.Equal(t, false, history["canonical"])
	require.NotContains(t, history, "lineage_id")
	generations := history["generations"].([]any)
	require.Len(t, generations, 1)
	generation := generations[0].(map[string]any)
	require.NotContains(t, generation, "generation")
	revisions := generation["revisions"].([]any)
	require.Len(t, revisions, 2)
	first := revisions[0].(map[string]any)
	last := revisions[1].(map[string]any)
	require.Equal(t, original.AppliedRevision, first["group_revision"])
	require.Equal(t, "2026-10-02T13:00:00Z", first["ends_at"])
	require.NotContains(t, last, "ends_at", "absence must not invent a boundary")
	require.Equal(t, "historical-connector", last["destination"].(map[string]any)["connector_id"])
	require.Equal(t, "historical-tag", last["sources"].([]any)[0].(map[string]any)["tag_id"])
	require.Equal(t, original.Destination.ConnectorID, generation["saved_destination"].(map[string]any)["connector_id"])
	raw, err := json.Marshal(history)
	require.NoError(t, err)
	for _, secret := range []string{"persisted-db", "persisted-schema", "/private", "private_table", "raw_values", "password", "connection_config"} {
		require.NotContains(t, string(raw), secret)
	}
	require.Equal(t, before, lifecycleReadState(t, db))
}

func TestLifecycleHistoryTombstoneAndReservedGenerationRemainReadOnly(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	h.db.SetMaxOpenConns(1)
	first, err := ensureLineageKind(h, "dual", h.mutation(t.Context(), lineageTestDraft(h, "dual")))
	require.NoError(t, err)
	lineage := seedLineageReadFixture(t, h, first.Group)
	successor := seedLineageGeneration(t, h, first.Group, lineage, "G-reserved", 2)
	_, err = h.db.Exec(`UPDATE write_groups SET status='deleted' WHERE id=?`, first.Group.ID)
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE write_group_lineages SET reserved_group_id=?,reservation_operation_id='O2',reservation_intent_digest='intent-2',revision=2 WHERE lineage_id=?`, successor.ID, lineage)
	require.NoError(t, err)
	before := lineageReadSnapshot(t, h)
	history := lifecycleHistoryJSON(t, h.service, first.Group.ID)
	require.Equal(t, true, history["canonical"])
	require.Equal(t, lineage, history["lineage_id"])
	require.Equal(t, first.Group.ID, history["current_group_id"])
	require.Equal(t, successor.ID, history["reserved_group_id"])
	generations := history["generations"].([]any)
	require.Len(t, generations, 2)
	old, next := generations[0].(map[string]any), generations[1].(map[string]any)
	require.Equal(t, "deleted", old["status"])
	require.Equal(t, true, old["current"])
	require.Equal(t, true, next["reserved"])
	require.Equal(t, float64(2), next["generation"])
	require.Equal(t, first.Group.ID, next["predecessor_group_id"])
	require.Empty(t, next["applied_revision"])
	require.Empty(t, next["revisions"])
	require.NotContains(t, next, "intake_enabled")
	require.NotContains(t, next, "schema_confirmed")
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

func TestLifecycleHistoryRejectsCorruptImmutableOrForeignBoundary(t *testing.T) {
	for _, mode := range []string{"foreign-payload", "foreign-end", "invalid-end", "missing-applied"} {
		t.Run(mode, func(t *testing.T) {
			db, service, group := lifecycleReadFixture(t, true)
			switch mode {
			case "foreign-payload":
				_, err := db.Exec(`UPDATE write_group_versions SET payload=json_set(payload,'$.workspace_id','foreign') WHERE group_id=?`, group.ID)
				require.NoError(t, err)
			case "foreign-end", "invalid-end":
				workspace, end := "foreign", "2026-10-02T13:00:00Z"
				if mode == "invalid-end" {
					workspace, end = group.WorkspaceID, "broken"
				}
				_, err := db.Exec(`INSERT INTO write_group_revision_ends(workspace_id,group_id,group_revision,ends_at,cutover_id) VALUES(?,?,?,?,?)`, workspace, group.ID, group.AppliedRevision, end, "C-end")
				require.NoError(t, err)
			case "missing-applied":
				_, err := db.Exec(`DELETE FROM write_group_versions WHERE group_id=?`, group.ID)
				require.NoError(t, err)
			}
			before := lifecycleReadState(t, db)
			view, err := service.Lifecycle(t.Context(), group.ID)
			require.Error(t, err, "unreadable ownership must not become empty successful history")
			require.Nil(t, view)
			require.Equal(t, before, lifecycleReadState(t, db))
		})
	}
}
