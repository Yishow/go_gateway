package main

import (
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestProductionMixedScopedStartIndependentNativeSources(t *testing.T) {
	f := newMixedScopedProductionFixture(t, true)
	status, response := httpJSON(t, f.router, http.MethodPost, testWriteBase+"/recording-start", f.request)
	require.Equal(t, http.StatusOK, status, "%+v", response)
	operation := response["data"].(map[string]any)

	require.Equal(t, "succeeded", operation["status"], "%+v", operation)
	require.Equal(t, []any{"mixed-selected", "mixed-selected-second"}, operation["device_ids"])
	groups := operation["groups"].([]any)
	require.Len(t, groups, 1)
	require.Equal(t, f.group.ID, groups[0].(map[string]any)["group_id"])
	require.Equal(t, true, groups[0].(map[string]any)["applied"])
	var groupCount int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM write_groups`).Scan(&groupCount))
	require.Equal(t, 1, groupCount, "mixed start must never create a hidden Basic recording")
	saved, err := f.services.writeGroups.Get(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Equal(t, f.group.RecordID, saved.Group.RecordID)
	require.Equal(t, f.group.Members, saved.Group.Members)
	require.Equal(t, f.group.RecordingProfile, saved.Group.RecordingProfile)
	statuses := f.services.scheduler.CollectionStatus()
	require.Len(t, statuses, 3)
	require.NotContains(t, statuses, f.excludedPoint, "unselected same-device source must not be installed")
	for _, member := range f.group.Members {
		require.Contains(t, statuses, member.PointID)
		require.Equal(t, member.DeviceID, statuses[member.PointID].DeviceID)
		require.True(t, statuses[member.PointID].Running)
	}
	other, err := f.services.device.GetByID(t.Context(), "mixed-unselected")
	require.NoError(t, err)
	require.Equal(t, schema.DeviceStatusDraft, other.Status, "not a whole-workspace activation")
	selected := map[string]string{}
	for _, member := range f.group.Members {
		selected[member.PointID] = member.DeviceID
	}
	observed := map[string]bool{}
	deadline := time.NewTimer(7 * time.Second)
	defer deadline.Stop()
	for len(observed) < len(selected) {
		select {
		case sample := <-f.capture.samples:
			require.Contains(t, selected, sample.PointID)
			require.Equal(t, selected[sample.PointID], sample.DeviceID)
			require.Equal(t, schema.QualityGood, sample.Quality)
			require.NotEmpty(t, sample.AcquisitionID)
			require.NotEmpty(t, sample.SourceRevision)
			observed[sample.PointID] = true
		case <-deadline.C:
			t.Fatalf("native observations missing selected sources: observed=%v selected=%v", observed, selected)
		}
	}

	// Lost/late client responses are recovered from the same durable operation;
	// replay does not need fresh revisions or reapply the saved record.
	replayStatus, replay := httpJSON(t, f.router, http.MethodPost, testWriteBase+"/recording-start", f.request)
	require.Equal(t, http.StatusOK, replayStatus, "%+v", replay)
	require.Equal(t, operation, replay["data"])
	changed := f.request
	changed.DeviceIDs = []string{"mixed-selected", "mixed-selected-second", "mixed-unselected"}
	changedStatus, changedResponse := httpJSON(t, f.router, http.MethodPost, testWriteBase+"/recording-start", changed)
	require.Equal(t, http.StatusConflict, changedStatus, "%+v", changedResponse)
	id := operation["operation_id"].(string)
	lookupStatus, lookup := httpJSON(t, f.router, http.MethodGet, testWriteBase+"/recording-start/operations/"+id, nil)
	require.Equal(t, http.StatusOK, lookupStatus)
	require.Equal(t, operation, lookup["data"])
	unknownStatus, _ := httpJSON(t, f.router, http.MethodGet, testWriteBase+"/recording-start/operations/unknown-operation", nil)
	require.Equal(t, http.StatusNotFound, unknownStatus)
	after, err := f.services.writeGroups.Get(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Equal(t, saved.Group.Revision, after.Group.Revision)
	require.Equal(t, saved.Group.AppliedRevision, after.Group.AppliedRevision)
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM managed_schema_operations WHERE action='recording_start'`).Scan(&groupCount))
	require.Equal(t, 1, groupCount)
	t.Log("real scoped HTTP start → group Apply → native Modbus scheduled samples; no full-window or SQL-result assertion")
}

func TestProductionMixedScopedStartRejectsUnreviewedScope(t *testing.T) {
	for _, mode := range []string{"workspace", "group revision", "connector revision", "foreign member", "missing preparation"} {
		t.Run(mode, func(t *testing.T) {
			f := newMixedScopedProductionFixture(t, mode != "missing preparation")
			request := f.request
			// Clone the only mutable slice before changing the submitted intent.
			request.Groups = append(request.Groups[:0:0], request.Groups...)
			switch mode {
			case "workspace":
				request.WorkspaceID = "foreign-workspace"
			case "group revision":
				request.Groups[0].ExpectedGroupRevision = "stale"
			case "connector revision":
				request.Groups[0].ExpectedConnectorRevision = "stale"
			case "foreign member":
				request.DeviceIDs = []string{"mixed-unselected"}
			}
			status, response := httpJSON(t, f.router, http.MethodPost, testWriteBase+"/recording-start", request)
			if mode == "workspace" || mode == "foreign member" {
				require.Equal(t, http.StatusBadRequest, status, "%+v", response)
			} else {
				require.Equal(t, http.StatusOK, status, "%+v", response)
				op := response["data"].(map[string]any)
				require.Equal(t, "failed", op["status"])
				reason := "stale_intent"
				if mode == "missing preparation" {
					reason = "preparation_required"
				}
				require.Equal(t, reason, op["reason"])
			}
			saved, err := f.services.writeGroups.Get(t.Context(), f.group.ID)
			require.NoError(t, err)
			require.Equal(t, f.group.Revision, saved.Group.Revision)
			require.Empty(t, saved.Group.AppliedRevision)
			require.Empty(t, f.services.scheduler.CollectionStatus())
			require.Empty(t, f.capture.samples)
			for _, id := range []string{"mixed-selected", "mixed-selected-second", "mixed-unselected"} {
				device, err := f.services.device.GetByID(t.Context(), id)
				require.NoError(t, err)
				require.Equal(t, schema.DeviceStatusDraft, device.Status)
			}
		})
	}
}
