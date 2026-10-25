package api

import (
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

// savedFacts are the saved fields of a group that reading, reloading or a
// resent Ensure must never change.
var savedFacts = []string{"id", "revision", "applied_revision", "status", "layout_version", "row_policy", "write_policy", "recording_profile", "destination", "members"}

func pick(group map[string]any, keys []string) map[string]any {
	picked := make(map[string]any, len(keys))
	for _, key := range keys {
		if value, ok := group[key]; ok {
			picked[key] = value
		}
	}
	return picked
}

// withCurrentRevision refreshes the workspace revision of a prepared request, as a
// client that has just reloaded would.
func withCurrentRevision(t *testing.T, f managedGroupFixture, request map[string]any) map[string]any {
	t.Helper()
	current, err := f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	request["expected_workspace_revision"] = current.DatabaseSetupRevision
	return request
}

// latestLongRequest is the Ensure request of an instantaneous long group that
// states no reduction: the saved shape every earlier group has.
func latestLongRequest(t *testing.T, f managedGroupFixture) map[string]any {
	t.Helper()
	request := profileRequest(t, f.writeGroupRouterFixture, "instantaneous", 15, "second", "main", "")
	group := request["group"].(map[string]any)
	group["layout_version"] = "managed_long_v1"
	group["storage_mode"] = "plain"
	group["destination"].(map[string]any)["table_schema"] = f.namespace
	return request
}

func saveGroup(t *testing.T, f managedGroupFixture, path string, request map[string]any, status int) map[string]any {
	t.Helper()
	response := performJSONRequest(t, f.router, http.MethodPost, path, withCurrentRevision(t, f, request))
	require.Equal(t, status, response.Code, response.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, response))
	return saved
}

// TestRecordingWorkflowCoexistence pins what the new recording workflow leaves
// alone and what it may only change through an explicit, reviewed step: a saved
// 60 second wide group, a custom group and a long group that never chose a
// reduction keep their saved behavior; a latest long group becomes a mean only
// as a new revision with fresh confirmation and Apply; and the canonical owner
// of a source cannot be bypassed by another group ID or another layout.
func TestRecordingWorkflowCoexistence(t *testing.T) {
	t.Run("saved 60 second wide, custom and latest long groups stay as saved when read, listed and ensured again", func(t *testing.T) {
		f := newManagedDestinationFixture(t, "sqlite")
		seedProfileSource(t, f.writeGroupRouterFixture, "second")

		wideRequest := f.createRequest()
		wide := wideRequest["group"].(map[string]any)
		wide["name"] = "Legacy wide"
		wide["destination"].(map[string]any)["storage_strategy"] = "managed"
		wide["destination"].(map[string]any)["table_schema"] = f.namespace
		wide["row_policy"] = map[string]any{"interval_seconds": 60, "allowed_lateness_seconds": 0}
		wideSaved := saveGroup(t, f, writeGroupsPath, wideRequest, http.StatusCreated)

		customRequest := f.createRequest()
		custom := customRequest["group"].(map[string]any)
		custom["name"] = "Custom"
		custom["members"] = []map[string]any{{"device_id": "device-A", "point_id": "point-second", "tag_id": "tag-second", "target_column": "second_col", "required": true}}
		customSaved := saveGroup(t, f, writeGroupsPath, customRequest, http.StatusCreated)

		longRequest := latestLongRequest(t, f)
		longSaved := saveGroup(t, f, ensureCircuitProfilePath, longRequest, http.StatusOK)
		require.NotContains(t, longSaved["recording_profile"].(map[string]any), "sample_reduction", "an omitted reduction stays omitted")

		saved := map[string]map[string]any{}
		for _, group := range []map[string]any{wideSaved, customSaved, longSaved} {
			saved[group["id"].(string)] = pick(group, savedFacts)
		}
		require.Equal(t, float64(60), wideSaved["row_policy"].(map[string]any)["interval_seconds"])
		require.Equal(t, "custom", customSaved["destination"].(map[string]any)["storage_strategy"])

		for id, want := range saved {
			require.Equal(t, want, pick(getGroup(t, f, id), savedFacts), "reading %s changes nothing", id)
		}
		listed := decodeJSONBody(t, performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil))["data"].(map[string]any)["groups"].([]any)
		require.Len(t, listed, 3)
		for _, entry := range listed {
			group := entry.(map[string]any)
			require.Equal(t, saved[group["id"].(string)], pick(group, savedFacts), "listing %s changes nothing", group["id"])
		}
		again := saveGroup(t, f, ensureCircuitProfilePath, latestLongRequest(t, f), http.StatusOK)
		require.Equal(t, saved[longSaved["id"].(string)], pick(again, savedFacts), "a resent Ensure returns the saved group and never converts it")
		require.Len(t, decodeJSONBody(t, performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil))["data"].(map[string]any)["groups"].([]any), 3)
	})

	t.Run("an explicit mean revision of a latest long group takes effect only after fresh confirmation and Apply", func(t *testing.T) {
		f := newManagedDestinationFixture(t, "sqlite")
		longSaved := saveGroup(t, f, ensureCircuitProfilePath, latestLongRequest(t, f), http.StatusOK)
		id := longSaved["id"].(string)
		require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, id, previewLong(t, f, id)).Status)
		applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
		require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
		oldGroup := getGroup(t, f, id)
		oldRevision := oldGroup["revision"].(string)
		require.Equal(t, oldRevision, oldGroup["applied_revision"])
		oldVersion := (&cutoverRig{t: t, f: f}).versionPayload(id, oldRevision)
		oldSnapshot, err := f.groups.ResolveAppliedAt(t.Context(), id, time.Now().Add(time.Hour))
		require.NoError(t, err)
		require.Equal(t, "latest", string(oldSnapshot.Group.RecordingProfile.EffectiveSampleReduction()))

		edit := getGroupData(t, f, id)
		group := edit["group"].(map[string]any)
		group["recording_profile"].(map[string]any)["sample_reduction"] = "mean"
		for _, member := range group["members"].([]any) {
			member.(map[string]any)["value_role"] = "gauge"
		}
		updated := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+id, groupEditRequest(f.writeGroupRouterFixture, edit, group))
		require.Equal(t, http.StatusOK, updated.Code, updated.Body.String())
		_, revised := groupSaveData(t, decodeJSONBody(t, updated))
		require.Equal(t, id, revised["id"], "the same group and canonical owner")
		require.NotEqual(t, oldRevision, revised["revision"], "an explicit mean is a new revision")
		require.Equal(t, oldRevision, revised["applied_revision"], "the prior revision stays in effect until a successful Apply")
		stillOld, err := f.groups.ResolveAppliedAt(t.Context(), id, time.Now().Add(time.Hour))
		require.NoError(t, err)
		require.Equal(t, oldRevision, stillOld.Group.Revision)
		require.Equal(t, "latest", string(stillOld.Group.RecordingProfile.EffectiveSampleReduction()))

		refused := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
		require.NotEqual(t, http.StatusOK, refused.Code, "without a fresh confirmation of the new revision: %s", refused.Body.String())
		require.Equal(t, oldRevision, getGroup(t, f, id)["applied_revision"])

		require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, id, previewLong(t, f, id)).Status)
		// The first Apply is scheduled for the next aligned window; the existing guard keeps one pending
		// cutover, so the second Apply is made once that window has passed on the service's own clock.
		afterFirstWindow := time.Now().Add(time.Hour)
		f.groups.WithClock(func() time.Time { return afterFirstWindow })
		t.Cleanup(func() { f.groups.WithClock(nil) })
		reapplied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
		require.Equal(t, http.StatusOK, reapplied.Code, reapplied.Body.String())
		now := getGroup(t, f, id)
		require.Equal(t, revised["revision"], now["applied_revision"])
		meanSnapshot, err := f.groups.ResolveAppliedAt(t.Context(), id, afterFirstWindow.Add(time.Hour))
		require.NoError(t, err)
		require.Equal(t, "mean", string(meanSnapshot.Group.RecordingProfile.EffectiveSampleReduction()))
		require.Equal(t, oldVersion, (&cutoverRig{t: t, f: f}).versionPayload(id, oldRevision), "the old version and its payload destination are not rewritten")
		require.Equal(t, 2, (&cutoverRig{t: t, f: f}).localCount(`SELECT COUNT(*) FROM write_group_versions WHERE group_id = ?`, id))
		var tagType string
		require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT data_type FROM tags WHERE id = 'tag-A'`).Scan(&tagType))
		require.Equal(t, "float32", tagType, "the Tag's own type is never rewritten by the choice")
	})

	t.Run("another layout is not a way around the canonical owner of a circuit profile", func(t *testing.T) {
		f := newManagedDestinationFixture(t, "sqlite")
		longSaved := saveGroup(t, f, ensureCircuitProfilePath, latestLongRequest(t, f), http.StatusOK)

		wideMean := wideProfileRequest(t, f, []meanSource{{name: "A", tagType: "float32", role: "gauge"}})
		wideMean["group"].(map[string]any)["recording_profile"].(map[string]any)["sample_reduction"] = "mean"
		response := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, withCurrentRevision(t, f, wideMean))
		require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
		assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_CIRCUIT_PROFILE_CONFLICT")
		listed := decodeJSONBody(t, performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil))["data"].(map[string]any)["groups"].([]any)
		require.Len(t, listed, 1, "no second owner was created")
		require.Equal(t, longSaved["id"], listed[0].(map[string]any)["id"])
	})

	t.Run("the generic create route with a client-chosen ID is not a way around the canonical owner either", func(t *testing.T) {
		f := newManagedDestinationFixture(t, "sqlite")
		longSaved := saveGroup(t, f, ensureCircuitProfilePath, latestLongRequest(t, f), http.StatusOK)

		other := withCurrentRevision(t, f, latestLongRequest(t, f))
		other["group"].(map[string]any)["id"] = "client-chosen-id"
		other["group"].(map[string]any)["name"] = "Another ID for the same circuit"
		response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, other)
		require.GreaterOrEqual(t, response.Code, 400, "a second group for the same circuit profile is refused: %s", response.Body.String())
		require.Less(t, response.Code, 500, response.Body.String())
		listed := decodeJSONBody(t, performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil))["data"].(map[string]any)["groups"].([]any)
		require.Len(t, listed, 1, "still one owner")
		require.Equal(t, longSaved["id"], listed[0].(map[string]any)["id"])
	})

	t.Run("a second wide group of the same kind over the same source cannot even be saved beside the first", func(t *testing.T) {
		f := newManagedDestinationFixture(t, "sqlite")
		sources := []meanSource{{name: "A", tagType: "float32", role: "gauge"}}
		first := saveGroup(t, f, ensureCircuitProfilePath, wideProfileRequest(t, f, sources), http.StatusOK)
		firstID := first["id"].(string)
		require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, firstID, previewLong(t, f, firstID)).Status)
		applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+firstID+"/apply", scopeBody(t, f, firstID))
		require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())

		otherCircuit := wideProfileRequest(t, f, sources)
		renameCircuit(otherCircuit, "other-circuit")
		refused := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, withCurrentRevision(t, f, otherCircuit))
		require.Equal(t, http.StatusConflict, refused.Code, refused.Body.String())
		assertWriteGroupError(t, decodeJSONBody(t, refused), "WRITE_GROUP_CIRCUIT_PROFILE_CONFLICT")
		listed := decodeJSONBody(t, performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil))["data"].(map[string]any)["groups"].([]any)
		require.Len(t, listed, 1, "the source has one owner of that kind")
	})

	t.Run("the instantaneous and the cumulative wide group of one circuit apply side by side", func(t *testing.T) {
		f := newManagedDestinationFixture(t, "sqlite", "uint64")
		insertMeasurement(t, f, "meas-energy", "A", "epoch-3", "counter")
		instantaneous := wideProfileRequest(t, f, []meanSource{{name: "A", tagType: "uint64", role: "gauge"}})
		cumulative := profileRequest(t, f.writeGroupRouterFixture, "cumulative", 1, "minute", "main", "")
		cumulativeGroup := cumulative["group"].(map[string]any)
		cumulativeGroup["layout_version"] = wideProfileLayoutVersion
		cumulativeGroup["destination"].(map[string]any)["table_schema"] = f.namespace
		member := cumulativeGroup["members"].([]map[string]any)[0]
		member["required"], member["measurement_id"] = false, "meas-energy"

		ids := make([]string, 0, 2)
		for _, request := range []map[string]any{instantaneous, cumulative} {
			saved := saveGroup(t, f, ensureCircuitProfilePath, request, http.StatusOK)
			id := saved["id"].(string)
			require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, id, previewLong(t, f, id)).Status)
			ids = append(ids, id)
		}
		for _, id := range ids {
			applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
			require.Equal(t, http.StatusOK, applied.Code, "two profiles of one source are not two writers of one kind: %s", applied.Body.String())
			require.Equal(t, getGroup(t, f, id)["revision"], getGroup(t, f, id)["applied_revision"])
		}
	})

	t.Run("an ordinary Apply of a new wide mean group cannot record beside the wide group that owns the source", func(t *testing.T) {
		r := newCutoverRigWith(t, "sqlite", cutoverOptions{wideSeconds: 60, wideTarget: true})
		wideID, targetID := r.wide["id"].(string), r.long["id"].(string)
		code, body := r.post(writeGroupsPath+"/"+targetID+"/apply", scopeBody(t, r.f, targetID))
		require.Equal(t, http.StatusConflict, code, "%v", body)
		require.Equal(t, "cutover_required", errorOf(t, body)["code"], "only a reviewed cutover may hand the source over")
		require.Equal(t, "", getGroup(t, r.f, targetID)["applied_revision"], "the new group was not applied")
		snapshot, err := r.f.groups.ResolveAppliedAt(t.Context(), wideID, time.Now().Add(time.Hour))
		require.NoError(t, err)
		require.True(t, snapshot.EndsAt.IsZero(), "the legacy group still owns its source")
	})
}
