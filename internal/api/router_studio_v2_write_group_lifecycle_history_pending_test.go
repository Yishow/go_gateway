package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func lifecycleHistorySeedPending(t *testing.T, f managedGroupFixture) {
	t.Helper()
	raw, err := json.Marshal(f.group)
	require.NoError(t, err)
	var group workspace.WriteGroup
	require.NoError(t, json.Unmarshal(raw, &group))
	d := group.Destination
	scope, err := snapshot.DestinationScope(d.ConnectorID, d.ConnectorRevision, d.Database, d.TableSchema, d.TableName)
	require.NoError(t, err)
	key := groupdelivery.GroupKey{WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.AppliedRevision}
	store := groupdelivery.NewStore(f.db)
	at := time.Date(2026, 10, 28, 2, 0, 0, 0, time.UTC)
	for i := range 3 {
		sample := snapshot.Sample{SampleID: fmt.Sprintf("history-%d", i), MemberKey: "tag-A", ObservedAt: at.Add(time.Duration(i) * time.Second), Quality: schema.QualityGood, Value: measurement.NewInt64(int64(i + 1)), SourceRevision: "source-1", MappingRevision: "mapping-1"}
		require.NoError(t, store.AppendSample(t.Context(), key, at, sample))
	}
	dedupe := strings.ToLower(strings.TrimSpace(group.WritePolicy.DedupeCapability))
	if dedupe == "" {
		dedupe = "none"
	}
	for i, state := range []string{groupdelivery.StatePending, groupdelivery.StateUnknown} {
		_, err := f.db.Exec(`INSERT INTO wg_delivery_outbox(effect_key,record_id,workspace_id,group_id,group_revision,bucket_start,partition_key,destination_scope,connector_id,connector_revision,database_name,table_schema,table_name,dedupe_capability,record_key_column,payload,payload_digest,state,next_retry_at,last_error_code,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			fmt.Sprintf("history-effect-%d", i), fmt.Sprintf("history-record-%d", i), key.WorkspaceID, key.GroupID, key.GroupRevision,
			at.Format(time.RFC3339Nano), "history-partition", scope, d.ConnectorID, d.ConnectorRevision, d.Database, d.TableSchema, d.TableName,
			dedupe, "", `{"private-sample":"never-display"}`, "frozen-digest", state, at.Format(time.RFC3339Nano), "private-delivery-error", at.Format(time.RFC3339Nano), at.Format(time.RFC3339Nano))
		require.NoError(t, err)
	}
}

func TestWriteGroupLifecycleHistoryMixedPendingUnitsRemainSeparate(t *testing.T) {
	f := lifecycleAPIApplied(t)
	lifecycleHistorySeedPending(t, f)
	before := lifecycleAPIState(t, f.db)
	data := lifecycleAPIRead(t, f.writeGroupRouterFixture, f.group["id"].(string))
	require.Contains(t, data, "history")
	generation := data["history"].(map[string]any)["generations"].([]any)[0].(map[string]any)
	revision := generation["revisions"].([]any)[0].(map[string]any)
	pending := revision["pending"].(map[string]any)
	require.EqualValues(t, 3, pending["samples"])
	outbox := pending["outbox"].(map[string]any)
	require.EqualValues(t, 2, outbox["total"])
	require.EqualValues(t, 2, outbox["outstanding"])
	require.EqualValues(t, 1, outbox["queued"])
	require.EqualValues(t, 1, outbox["unknown"])
	require.EqualValues(t, 0, outbox["sql_committed"])
	for _, field := range []string{"total", "records", "record_count", "intake_enabled", "readback_verified", "cleanup_completed"} {
		require.NotContains(t, pending, field)
	}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	for _, value := range []string{f.file, "private-delivery-error", "private-sample", "never-display", "frozen-digest"} {
		require.NotContains(t, string(raw), value)
	}
	require.Equal(t, before, lifecycleAPIState(t, f.db))
}

func TestWriteGroupLifecycleHistoryOrphanPendingRevisionFailsClosed(t *testing.T) {
	f := lifecycleAPIApplied(t)
	lifecycleHistorySeedPending(t, f)
	_, err := f.db.Exec(`UPDATE wg_delivery_outbox SET group_revision='private-orphan-revision' WHERE effect_key='history-effect-1'`)
	require.NoError(t, err)
	before := lifecycleAPIState(t, f.db)
	response := performJSONRequest(t, f.router, "GET", writeGroupsPath+"/"+f.group["id"].(string)+"/lifecycle", nil)
	require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
	require.NotContains(t, response.Body.String(), "private-")
	require.NotContains(t, response.Body.String(), `"history"`)
	require.Equal(t, before, lifecycleAPIState(t, f.db))
}
