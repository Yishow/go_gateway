package api

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

func TestMixedAcceptanceFreshWideTargetPreservesLegacyRawHistory(t *testing.T) {
	f, request := mixedNativeRequest(t, true, []int64{15, 60, 60}, []bool{true, false, false})
	base := time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC)
	old := &mixedAcceptanceRig{f: f, clock: base.Add(7 * time.Second)}
	f.groups.WithClock(func() time.Time { return old.clock })
	raw := f.createRequest()
	rawGroup := raw["group"].(map[string]any)
	rawGroup["layout_version"] = "managed_wide_v1"
	rawGroup["members"] = []map[string]any{{"device_id": "device-A", "point_id": "point-a", "tag_id": "tag-a", "target_column": "raw_a", "required": true}}
	dest := rawGroup["destination"].(map[string]any)
	dest["storage_strategy"] = "managed"
	dest["table_schema"] = "main"
	dest["table_name"] = "legacy_raw"
	created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, raw)
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, created))
	id := saved["id"].(string)
	confirmLong(t, f, id, previewLong(t, f, id))
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	var err error
	old.frozen, err = f.groups.ResolveAppliedAt(t.Context(), id, base.Add(time.Minute))
	require.NoError(t, err)
	old.clock = old.frozen.EffectiveAt.Add(time.Second)
	old.store = newMixedAcceptanceStore(t, f.db, func() time.Time { return old.clock })
	old.reconcile(t)
	old.f.target, err = sql.Open("sqlite", f.file)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, old.f.target.Close()) })
	old.offer(t, 0, "legacy-raw", old.clock, float64(21.5))
	old.tick(t, old.frozen.EffectiveAt.Add(15*time.Second))
	require.Len(t, old.deliver(t), 1)
	table := old.frozen.Group.Destination.TableName
	before := mixedWideHistoricalRow(t, old.f.target, table)
	column := old.frozen.Group.Members[0].TargetColumn
	require.Equal(t, float64(21.5), before[column], "legacy stores its original raw value")
	original, err := f.groups.Get(t.Context(), id)
	require.NoError(t, err)
	request["expected_workspace_revision"] = original.WorkspaceRevision
	current := &mixedAcceptanceRig{f: f, clock: base.Add(67 * time.Second)}
	f.groups.WithClock(func() time.Time { return current.clock })
	created = performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	current.f.data, current.f.group = groupSaveData(t, decodeJSONBody(t, created))
	mixedID := current.f.group["id"].(string)
	confirmLong(t, f, mixedID, previewLong(t, f, mixedID))
	denied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+mixedID+"/apply", scopeBody(t, f, mixedID))
	require.Equal(t, http.StatusConflict, denied.Code, denied.Body.String())
	require.Contains(t, denied.Body.String(), "cutover_required")
	retained, err := f.groups.Get(t.Context(), id)
	require.NoError(t, err)
	require.Equal(t, original.Group, retained.Group, "a failed takeover does not automatically stop the old record")
	require.Equal(t, before, mixedWideHistoricalRow(t, old.f.target, table))
	var versions int
	require.NoError(t, f.db.QueryRow(`SELECT count(*) FROM wg_runtime_versions WHERE group_id=?`, mixedID).Scan(&versions))
	require.Zero(t, versions)
	// This is an explicit user lifecycle command, not an implicit side effect of
	// creating the calculated record. The original source-ownership guard stays.
	stopped := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/disable", scopeBody(t, f, id))
	require.Equal(t, http.StatusOK, stopped.Code, stopped.Body.String())
	afterStop, err := f.groups.Get(t.Context(), id)
	require.NoError(t, err)
	require.Equal(t, "disabled", string(afterStop.Group.Status))
	member := old.frozen.Group.Members[0]
	old.clock = current.clock
	rejected := old.pipe.AcceptSample(t.Context(), measurement.SampleEnvelope{SampleID: "post-explicit-stop", WorkspaceID: old.frozen.WorkspaceID, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID, SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision, ObservedAt: current.clock, ReceivedAt: current.clock, TimeOrigin: "gateway", Value: float64(99), Quality: schema.QualityGood})
	require.ErrorIs(t, rejected, groupdelivery.ErrIntakeStopped, "a cached old worker cannot bypass explicit stop")

	confirmLong(t, f, mixedID, previewLong(t, f, mixedID))
	applied = performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+mixedID+"/apply", scopeBody(t, f, mixedID))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	current.frozen, err = f.groups.ResolveAppliedAt(t.Context(), mixedID, current.clock)
	require.NoError(t, err)
	current.store = newMixedAcceptanceStore(t, f.db, func() time.Time { return current.clock })
	current.reconcile(t)
	current.f.target = old.f.target

	require.NotEqual(t, id, current.frozen.GroupID)
	require.NotEqual(t, table, current.frozen.Group.Destination.TableName, "raw and calculated wide columns must not share a physical target")
	current.offer(t, 0, "new-mean", base.Add(70*time.Second), float64(22))
	current.tick(t, base.Add(75*time.Second))
	require.Len(t, current.deliver(t), 1)
	after := mixedWideHistoricalRow(t, current.f.target, table)
	require.Equal(t, before, after, "every historical raw value, key, owner marker and provenance column remains unchanged")
	var result float64
	require.NoError(t, current.f.target.QueryRow(`SELECT "`+current.frozen.Group.Members[0].TargetColumn+`" FROM "`+current.frozen.Group.Destination.TableName+`"`).Scan(&result))
	require.Equal(t, float64(22), result)
	unchanged, err := f.groups.Get(t.Context(), id)
	require.NoError(t, err)
	require.Equal(t, afterStop.Group, unchanged.Group, "new record must not further mutate or repurpose the explicitly stopped record/target ownership")
}

func mixedWideHistoricalRow(t *testing.T, db *sql.DB, table string) map[string]any {
	t.Helper()
	rows, err := db.Query(`SELECT * FROM "` + table + `"`)
	require.NoError(t, err)
	defer rows.Close()
	columns, err := rows.Columns()
	require.NoError(t, err)
	require.True(t, rows.Next())
	values := make([]any, len(columns))
	dest := make([]any, len(columns))
	for i := range values {
		dest[i] = &values[i]
	}
	require.NoError(t, rows.Scan(dest...))
	result := map[string]any{}
	for i, c := range columns {
		result[c] = values[i]
	}
	require.False(t, rows.Next())
	require.NoError(t, rows.Err())
	return result
}
