package api

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMixedAcceptanceIndependentRecordsShareSamplesAndStopSeparately(t *testing.T) {
	for _, long := range []bool{false, true} {
		name := "wide_to_wide"
		if long {
			name = "wide_to_long"
		}
		t.Run(name, func(t *testing.T) {
			f, request := mixedNativeRequest(t, true, []int64{5, 60, 60}, []bool{true, false, false})
			group := request["group"].(map[string]any)
			first := group["members"].([]map[string]any)[0]
			first["calculation"] = "latest"
			group["members"] = []map[string]any{first}
			group["destination"].(map[string]any)["table_name"] = "record_a"
			base := time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC)
			a := applyMixedAcceptance(t, f, request, base.Add(2*time.Second))
			bytes, err := json.Marshal(request)
			require.NoError(t, err)
			var second map[string]any
			require.NoError(t, json.Unmarshal(bytes, &second))
			candidate := second["group"].(map[string]any)
			candidate["record_id"] = uuid.NewString()
			candidate["destination"].(map[string]any)["table_name"] = "record_b"
			member := candidate["members"].([]any)[0].(map[string]any)
			member["calculation"] = "mean"
			member["cadence"].(map[string]any)["interval"] = 60
			if long {
				candidate["layout_version"] = "managed_long_v1"
				candidate["storage_mode"] = "plain"
				candidate["destination"].(map[string]any)["table_name"] = ""
				delete(member, "target_column")
			}
			latest, err := f.groups.Get(t.Context(), a.frozen.GroupID)
			require.NoError(t, err)
			second["expected_workspace_revision"] = latest.WorkspaceRevision
			b := applyMixedAcceptance(t, f, second, base.Add(3*time.Second))
			require.NotEqual(t, a.frozen.Group.RecordID, b.frozen.Group.RecordID)
			a.store = b.store
			b.offer(t, 0, "shared-four", base.Add(4*time.Second), float64(10))
			b.tick(t, base.Add(5*time.Second))
			a.clock = b.clock
			require.Len(t, a.deliver(t), 1)
			require.Empty(t, b.deliver(t))
			b.offer(t, 0, "shared-six", base.Add(6*time.Second), float64(30))
			b.tick(t, base.Add(10*time.Second))
			a.clock = b.clock
			require.Len(t, a.deliver(t), 1)
			require.Empty(t, b.deliver(t))
			b.clock = base.Add(12 * time.Second)
			stopped := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+a.frozen.GroupID+"/disable", scopeBody(t, f, a.frozen.GroupID))
			require.Equal(t, http.StatusOK, stopped.Code, stopped.Body.String())
			b.reconcile(t)
			b.offer(t, 0, "only-b-after-a-stopped", base.Add(16*time.Second), float64(20))
			b.tick(t, base.Add(60*time.Second))
			require.Len(t, b.deliver(t), 1, "stopping record A cannot consume or discard record B's open 60s window")
			var n int
			require.NoError(t, a.f.target.QueryRowContext(t.Context(), `SELECT count(*) FROM "`+a.frozen.Group.Destination.TableName+`"`).Scan(&n))
			require.Equal(t, 2, n)
			column := a.frozen.Group.Members[0].TargetColumn
			rows, err := a.f.target.QueryContext(t.Context(), `SELECT "`+column+`" FROM "`+a.frozen.Group.Destination.TableName+`" ORDER BY due_at`)
			require.NoError(t, err)
			var values []float64
			for rows.Next() {
				var v float64
				require.NoError(t, rows.Scan(&v))
				values = append(values, v)
			}
			require.NoError(t, rows.Err())
			require.NoError(t, rows.Close())
			require.Equal(t, []float64{10, 30}, values)
			column = b.frozen.Group.Members[0].TargetColumn
			if long {
				column = "value_float"
			}
			var mean float64
			require.NoError(t, b.f.target.QueryRowContext(t.Context(), `SELECT "`+column+`" FROM "`+b.frozen.Group.Destination.TableName+`"`).Scan(&mean))
			require.Equal(t, float64(20), mean)
			retained, err := f.groups.Get(t.Context(), b.frozen.GroupID)
			require.NoError(t, err)
			require.NotEqual(t, "disabled", string(retained.Group.Status))
			require.Equal(t, b.frozen.AppliedRevision, retained.Group.AppliedRevision)
		})
	}
}
