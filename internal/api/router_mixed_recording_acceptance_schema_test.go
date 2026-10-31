package api

import (
	"encoding/json"
	"go-gateway/internal/datalink/workspace"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMixedAcceptanceScaledSchemaAndApply(t *testing.T) {
	f, request := mixedScaledSaveFixture(t, "delta")
	now := time.Date(2026, 10, 31, 12, 0, 7, 0, time.UTC)
	f.groups.WithClock(func() time.Time { return now })
	saved := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)
	require.Equal(t, http.StatusCreated, saved.Code, saved.Body.String())
	f.data, f.group = groupSaveData(t, decodeJSONBody(t, saved))
	id := f.group["id"].(string)
	before, err := f.groups.Readiness(t.Context(), id)
	require.NoError(t, err)
	require.False(t, before.Ready)
	confirmLong(t, f, id, previewLong(t, f, id))
	ready, err := f.groups.Readiness(t.Context(), id)
	require.NoError(t, err)
	require.True(t, ready.ConfigReady, "%+v", ready.Issues)
	require.True(t, ready.SchemaReady, "%+v", ready.Issues)
	require.True(t, ready.Ready, "complete real schema and counter proof must serve supported mixed delta: %+v", ready.Issues)
	applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	frozen, err := f.groups.ResolveAppliedAt(t.Context(), id, now)
	require.NoError(t, err)
	require.Equal(t, now, frozen.EffectiveAt, "new mixed starts at Apply clock without rephasing member windows")
	require.Zero(t, frozen.Group.RowPolicy.IntervalSeconds)
	var payload string
	require.NoError(t, f.db.QueryRow(`SELECT payload FROM wg_runtime_versions WHERE workspace_id=? AND group_id=? AND group_revision=?`, f.record.ID, id, frozen.AppliedRevision).Scan(&payload))
	var durable workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal([]byte(payload), &durable))
	require.NotNil(t, durable.RuntimeLayout)
	require.NotEmpty(t, durable.CounterBindings)
}
