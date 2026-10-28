package api

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
	"net/http"
	"testing"
	"time"
)

func TestPostgresSuccessorActivationPublicEligibility(t *testing.T) {
	r := newPostgresSuccessorRuntimeRig(t, false, true, false, false)
	f := r.activationAPI()
	review := f.postActivation(t, f.activationRequest(t))
	require.Empty(t, review["blockers"], "real PostgreSQL preparation must be eligible for activation")
}

func TestPostgresSuccessorActivationFourPlainTransitions(t *testing.T) {
	for _, c := range []struct {
		name                   string
		sourceLong, targetLong bool
	}{{"wide_long", false, true}, {"long_wide", true, false}, {"wide_wide", false, false}, {"long_long", true, true}} {
		t.Run(c.name, func(t *testing.T) {
			r := newPostgresSuccessorRuntimeRig(t, c.sourceLong, c.targetLong, false, false)
			f := r.activationAPI()
			before := r.targetState(t)
			local := successorAPIAllState(t, r.f.db)
			confirmation := f.activationConfirmation(t)
			require.Equal(t, local, successorAPIAllState(t, r.f.db), "review/readiness are read-only")
			event := f.postActivation(t, confirmation)
			require.NotEmpty(t, event["id"])
			next, err := r.f.groups.ResolveAppliedAt(t.Context(), r.prepared.Group.ID, r.boundary)
			require.NoError(t, err)
			require.Equal(t, "connector-B", next.Group.Destination.ConnectorID)
			require.Equal(t, "connector-2", next.Group.Destination.ConnectorRevision)
			require.Equal(t, r.targetNamespace, next.Group.Destination.TableSchema)
			require.NotEqual(t, r.groupID, next.Group.ID)
			_, err = r.f.groups.ResolveAppliedAt(t.Context(), r.prepared.Group.ID, r.boundary.Add(-time.Nanosecond))
			require.ErrorIs(t, err, workspace.ErrWriteGroupLifecycleBlocked)
			after := successorAPIAllState(t, r.f.db)
			require.Equal(t, local["write_group_storage_scopes"], after["write_group_storage_scopes"], "activation does not reset shared storage generation/owner")
			require.Equal(t, event["id"], f.postActivation(t, confirmation)["id"])
			require.Equal(t, after, successorAPIAllState(t, r.f.db))
			raw, err := json.Marshal(confirmation)
			require.NoError(t, err)
			var changed map[string]any
			require.NoError(t, json.Unmarshal(raw, &changed))
			changed["envelope"].(map[string]any)["request"].(map[string]any)["boundary"] = r.boundary.Add(time.Second)
			rejected := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", changed)
			require.GreaterOrEqual(t, rejected.Code, 400, rejected.Body.String())
			require.Less(t, rejected.Code, 500, rejected.Body.String())
			require.Equal(t, after, successorAPIAllState(t, r.f.db), "replay cannot replace frozen intent")

			require.Equal(t, before, r.targetState(t), "activation/replay never alter external history or markers")
		})
	}
}

func TestPostgresSuccessorActivationTimescaleDestination(t *testing.T) {
	r := newPostgresSuccessorRuntimeRig(t, true, true, false, true)
	var extension string
	require.NoError(t, r.target.QueryRow(`SELECT extversion FROM pg_extension WHERE extname='timescaledb'`).Scan(&extension))
	require.NotEmpty(t, extension)
	for _, scope := range []string{r.f.namespace, r.targetNamespace} {
		var n int
		require.NoError(t, r.target.QueryRow(`SELECT count(*) FROM timescaledb_information.hypertables WHERE hypertable_schema=$1`, scope).Scan(&n))
		require.Equal(t, 1, n, "both scopes are real hypertables")
		var column, interval string
		require.NoError(t, r.target.QueryRow(`SELECT column_name,time_interval::text FROM timescaledb_information.dimensions WHERE hypertable_schema=$1`, scope).Scan(&column, &interval))
		require.Equal(t, "bucket_start", column)
		require.Equal(t, "7 days", interval)
	}
	f := r.activationAPI()
	before := r.targetState(t)
	event := f.postActivation(t, f.activationConfirmation(t))
	require.NotEmpty(t, event["id"])
	next, err := r.f.groups.ResolveAppliedAt(t.Context(), r.prepared.Group.ID, r.boundary)
	require.NoError(t, err)
	require.Equal(t, workspace.WriteGroupStorageModeTimescale, next.Group.StorageMode)
	require.Equal(t, before, r.targetState(t), "no mode downgrade or marker mutation")
}
