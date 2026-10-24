package api

import (
	"net/http"
	"testing"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// profilePair is one circuit's two canonical profile groups on the same source:
// an instantaneous group and a cumulative group, both on the shared long table.
type profilePair struct {
	f          managedGroupFixture
	instant    map[string]any
	cumulative map[string]any
}

// insertCounterMeasurement saves a definition for point-A/tag-A. semantic and
// unit are written as given so a test can save a definition that does not
// confirm counter semantics.
func insertCounterMeasurement(t *testing.T, f managedGroupFixture, id, semantic, unit string) {
	t.Helper()
	_, err := f.db.ExecContext(t.Context(), `
		INSERT INTO measurement_definitions (
			id, workspace_id, device_id, point_id, tag_id, equipment_id, definition_revision, source_binding_revision,
			series_epoch, name, quantity, unit, semantic_kind
		) VALUES (?, ?, 'device-A', 'point-A', 'tag-A', 'device-A', 'def-1', 'bind-1', 'epoch-1', 'Energy', 'energy', ?, ?)`,
		id, f.record.ID, unit, semantic)
	require.NoError(t, err)
}

func ensureProfileGroup(t *testing.T, f managedGroupFixture, kind string, interval any, unit, measurementID string, storageMode ...string) map[string]any {
	t.Helper()
	request := profileRequest(t, f.writeGroupRouterFixture, kind, interval, unit, "main", "")
	group := request["group"].(map[string]any)
	group["layout_version"] = "managed_long_v1"
	group["storage_mode"] = "plain"
	if len(storageMode) > 0 {
		group["storage_mode"] = storageMode[0]
	}
	group["destination"].(map[string]any)["table_schema"] = f.namespace
	if measurementID != "" {
		group["members"].([]map[string]any)[0]["measurement_id"] = measurementID
	}
	response := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, response))
	return saved
}

// newProfilePair saves both profiles of one circuit and confirms the schema of
// each, exactly as Basic does: one explicit confirmation per canonical group.
func newProfilePair(t *testing.T, kind, semantic, unit string) profilePair {
	t.Helper()
	f := newManagedDestinationFixture(t, kind, "uint64")
	insertCounterMeasurement(t, f, "meas-energy", semantic, unit)
	pair := profilePair{f: f}
	pair.instant = ensureProfileGroup(t, f, "instantaneous", 15, "second", "")
	pair.cumulative = ensureProfileGroup(t, f, "cumulative", 10, "minute", "meas-energy")
	for _, group := range []map[string]any{pair.instant, pair.cumulative} {
		id := group["id"].(string)
		require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, id, previewLong(t, f, id)).Status)
	}
	return pair
}

func (p profilePair) request(t *testing.T, requestID string) workspace.RecordingStartRequest {
	t.Helper()
	current, err := p.f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	intents := make([]workspace.RecordingStartGroupIntent, 0, 2)
	for _, group := range []map[string]any{p.instant, p.cumulative} {
		reloaded := performJSONRequest(t, p.f.router, http.MethodGet, writeGroupsPath+"/"+group["id"].(string), nil)
		require.Equal(t, http.StatusOK, reloaded.Code, reloaded.Body.String())
		_, saved := groupSaveData(t, decodeJSONBody(t, reloaded))
		intents = append(intents, workspace.RecordingStartGroupIntent{
			GroupID: saved["id"].(string), ExpectedGroupRevision: saved["revision"].(string), ExpectedConnectorRevision: "connector-1",
		})
	}
	return workspace.RecordingStartRequest{
		RequestID: requestID, WorkspaceID: p.f.record.ID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		DeviceIDs: []string{"device-A"}, Groups: intents,
	}
}

func appliedRevisionOf(t *testing.T, f managedGroupFixture, groupID string) string {
	t.Helper()
	response := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+groupID, nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, response))
	applied, _ := saved["applied_revision"].(string)
	return applied
}

// TestRecordingStartAcceptsBothCanonicalProfileGroups pins the Basic start
// contract of 4.1: one recording-start intent names the instantaneous and the
// cumulative group of the same source, each at its own revision, and both are
// applied before the device activates. The first request is answered again for
// the same request id without a second activation.
func TestRecordingStartAcceptsBothCanonicalProfileGroups(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			pair := newProfilePair(t, kind, "counter", "kWh")
			svc := newRecordingStartServiceFixture(t, managedGroupFixture{writeGroupRouterFixture: pair.f.writeGroupRouterFixture, groups: pair.f.groups}, recordingStartSuccess())
			request := pair.request(t, "start-both-profiles")

			operation, err := svc.start.Start(t.Context(), request)
			require.NoError(t, err)
			require.Equal(t, recordingplan.SchemaOperationSucceeded, operation.Status, "%+v", operation)
			require.Len(t, operation.Groups, 2)
			for _, progress := range operation.Groups {
				require.True(t, progress.Saved && progress.Ready && progress.Applied, "%+v", progress)
			}
			require.NotEmpty(t, appliedRevisionOf(t, pair.f, pair.instant["id"].(string)))
			require.NotEmpty(t, appliedRevisionOf(t, pair.f, pair.cumulative["id"].(string)))
			require.Equal(t, 1, svc.activation.Calls())

			again, err := svc.start.Start(t.Context(), request)
			require.NoError(t, err)
			require.Equal(t, operation.OperationID, again.OperationID, "a resent request id reads the original operation")
			require.Equal(t, 1, svc.activation.Calls(), "a resend must not activate twice")
		})
	}
}

// TestCumulativeProfileRefusesAnUnsafeSavedMeasurement is the 2.4 hand-off: a
// cumulative source is ready only when its saved definition really belongs to
// the member and confirms counter semantics. A non-empty measurement id is not
// enough, and every refusal keeps the group unapplied and the device inactive.
func TestCumulativeProfileRefusesAnUnsafeSavedMeasurement(t *testing.T) {
	cases := []struct {
		name   string
		code   string
		mutate func(t *testing.T, p profilePair)
	}{
		{"a gauge definition is not a counter", "counter-measurement-not-counter", func(t *testing.T, p profilePair) {
			_, err := p.f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET semantic_kind='gauge' WHERE id='meas-energy'`)
			require.NoError(t, err)
		}},
		{"a definition moved to another point is foreign", "counter-measurement-foreign", func(t *testing.T, p profilePair) {
			_, err := p.f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET point_id='point-elsewhere' WHERE id='meas-energy'`)
			require.NoError(t, err)
		}},
		{"a definition of another workspace is foreign", "counter-measurement-foreign", func(t *testing.T, p profilePair) {
			_, err := p.f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET workspace_id='other-workspace' WHERE id='meas-energy'`)
			require.NoError(t, err)
		}},
		{"a deleted definition is unavailable", "counter-measurement-unavailable", func(t *testing.T, p profilePair) {
			_, err := p.f.db.ExecContext(t.Context(), `DELETE FROM measurement_definitions WHERE id='meas-energy'`)
			require.NoError(t, err)
		}},
		{"a counter without a unit is unconfirmed", "counter-measurement-unconfirmed", func(t *testing.T, p profilePair) {
			_, err := p.f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET unit='' WHERE id='meas-energy'`)
			require.NoError(t, err)
		}},
		{"a counter that allows negative steps is unconfirmed", "counter-measurement-unconfirmed", func(t *testing.T, p profilePair) {
			_, err := p.f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET counter_policy='{"allow_negative":true}' WHERE id='meas-energy'`)
			require.NoError(t, err)
		}},
	}
	for _, kind := range []string{"sqlite", "postgres"} {
		for _, tc := range cases {
			t.Run(kind+"/"+tc.name, func(t *testing.T) {
				pair := newProfilePair(t, kind, "counter", "kWh")
				cumulativeID := pair.cumulative["id"].(string)
				ready, err := pair.f.groups.Readiness(t.Context(), cumulativeID)
				require.NoError(t, err)
				require.True(t, ready.Ready, "a confirmed counter is ready before the mutation: %+v", ready.Issues)

				tc.mutate(t, pair)
				ready, err = pair.f.groups.Readiness(t.Context(), cumulativeID)
				require.NoError(t, err)
				require.False(t, ready.Ready)
				codes := make([]string, 0, len(ready.Issues))
				for _, issue := range ready.Issues {
					codes = append(codes, issue.Code)
				}
				require.Contains(t, codes, tc.code)

				svc := newRecordingStartServiceFixture(t, managedGroupFixture{writeGroupRouterFixture: pair.f.writeGroupRouterFixture, groups: pair.f.groups}, recordingStartSuccess())
				operation, err := svc.start.Start(t.Context(), pair.request(t, "start-unsafe"))
				require.NoError(t, err)
				require.NotEqual(t, recordingplan.SchemaOperationSucceeded, operation.Status, "%+v", operation)
				require.Zero(t, svc.activation.Calls(), "an unsafe counter must never activate the device")
				require.Empty(t, appliedRevisionOf(t, pair.f, cumulativeID), "the unsafe cumulative group stays unapplied")

				applied := performJSONRequest(t, pair.f.router, http.MethodPost, writeGroupsPath+"/"+cumulativeID+"/apply", scopeBody(t, pair.f, cumulativeID))
				require.GreaterOrEqual(t, applied.Code, http.StatusBadRequest, applied.Body.String())
				require.Empty(t, appliedRevisionOf(t, pair.f, cumulativeID))
			})
		}
	}
}

// TestInstantaneousProfileNeedsNoMeasurement keeps the other half of the rule:
// a raw source recorded instantaneously never needs a measurement definition.
func TestInstantaneousProfileNeedsNoMeasurement(t *testing.T) {
	f := newManagedDestinationFixture(t, "sqlite", "float32")
	group := ensureProfileGroup(t, f, "instantaneous", 15, "second", "")
	id := group["id"].(string)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, id, previewLong(t, f, id)).Status)
	ready, err := f.groups.Readiness(t.Context(), id)
	require.NoError(t, err)
	require.True(t, ready.Ready, "%+v", ready.Issues)
	var measurements int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM measurement_definitions`).Scan(&measurements))
	require.Zero(t, measurements)
}
