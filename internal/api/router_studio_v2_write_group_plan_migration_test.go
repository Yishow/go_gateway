package api

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

const planMigrationPreviewPath = writeGroupsPath + "/migrations/recording-plans/preview"
const planMigrationReviewPath = writeGroupsPath + "/migrations/recording-plans/review"

func newPlanMigrationRouterFixture(t *testing.T) (writeGroupRouterFixture, *recordingplan.RecordingPlan, string) {
	t.Helper()
	f, targetPath := newSingleMigrationRouterFixture(t)
	tagID := "tag-A"
	definition := &measurement.MeasurementDefinition{
		ID: "measurement-A", WorkspaceID: f.record.ID, DeviceID: "device-A", PointID: "point-A", TagID: &tagID,
		EquipmentID: "device-A", Name: "Temperature", Quantity: "temperature", Unit: "C",
		DefinitionRevision: "definition-1", SourceBindingRevision: "binding-1", SeriesEpoch: "epoch-1",
		SemanticKind: measurement.SemanticKindGauge,
	}
	require.NoError(t, measurement.NewSQLRepository(f.db).Create(t.Context(), definition))
	plan := &recordingplan.RecordingPlan{
		ID: "legacy-plan-A", WorkspaceID: f.record.ID, Revision: "plan-1", AppliedRevision: "legacy-applied-1",
		Name: "Original every-sample plan", Status: recordingplan.PlanStatusRunning, Timezone: "Asia/Taipei",
		Members:      []recordingplan.PlanMember{{MemberID: "member-A", MeasurementID: "measurement-A", EquipmentID: "device-A", Name: "Original member"}},
		Streams:      []recordingplan.PlanStream{{StreamID: "stream-A", MeasurementID: "measurement-A", Mode: recordingplan.StreamModeRawHistory, RawPolicy: recordingplan.RawPolicyEverySample, DestinationIDs: []string{"destination-A"}}},
		Destinations: []recordingplan.PlanDestination{{DestinationID: "destination-A", ConnectorID: "connector-A", ConnectorRevision: "connector-1", TablePrefix: "original_prefix", WriteIntervalSec: 15, BatchSize: 7}},
		Retention:    recordingplan.RetentionPolicy{RawDays: 12, SummaryDays: 30, EventsDays: 60, CorrectionHorizonHours: 8},
		Limits:       recordingplan.PlanLimits{MaxBatchSize: 7, MaxHoldSeconds: 9, MaxQueueBytes: 4096},
	}
	repo := recordingplan.NewSQLRepository(f.db)
	require.NoError(t, repo.CreatePlan(t.Context(), plan))
	saved, err := repo.GetPlanByWorkspace(t.Context(), plan.ID, f.record.ID)
	require.NoError(t, err)
	return f, saved, targetPath
}

func planMigrationPreviewData(t *testing.T, f writeGroupRouterFixture) map[string]any {
	t.Helper()
	response := performJSONRequest(t, f.router, http.MethodPost, planMigrationPreviewPath, map[string]any{
		"workspace_id": f.record.ID, "source_ids": []string{"legacy-plan-A"},
	})
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	return decodeJSONBody(t, response)["data"].(map[string]any)
}

func planMigrationReviewRequest(data map[string]any) map[string]any {
	return map[string]any{
		"workspace_id": data["workspace_id"], "expected_workspace_revision": data["workspace_revision"],
		"expected_connector_revision": data["connector_revision"], "review_digest": data["review_digest"],
		"source_ids": []string{"legacy-plan-A"}, "confirm_snapshot_conversion": true,
	}
}

func TestNewRouter_BasicPlanMigrationPreservesOriginalAndBlocksReview(t *testing.T) {
	f, original, targetPath := newPlanMigrationRouterFixture(t)
	data := planMigrationPreviewData(t, f)
	require.Equal(t, f.record.ID, data["workspace_id"])
	require.Equal(t, f.record.DatabaseSetupRevision, data["workspace_revision"])
	require.Equal(t, "connector-1", data["connector_revision"])
	require.Equal(t, "recording-plan-v1", data["adapter_version"])
	require.NotEmpty(t, data["review_digest"])
	items := data["items"].([]any)
	require.Len(t, items, 1)
	item := items[0].(map[string]any)
	require.Equal(t, original.ID, item["source_id"])
	require.NotEmpty(t, item["source_revision"])
	require.Equal(t, "blocked", item["status"])
	require.Equal(t, "open_write_groups", item["repair_action"])
	require.NotEmpty(t, item["issues"])
	require.NotContains(t, item, "candidate_group")
	intent := item["before_recording_plan_intent"].(map[string]any)
	require.Equal(t, original.ID, intent["source_id"])
	encoded, err := json.Marshal(original)
	require.NoError(t, err)
	var expected map[string]any
	require.NoError(t, json.Unmarshal(encoded, &expected))
	require.Equal(t, expected, intent["plan"])
	sources := intent["sources"].([]any)
	require.Len(t, sources, 1)
	require.Equal(t, "measurement-A", sources[0].(map[string]any)["measurement_id"])
	require.Equal(t, data, planMigrationPreviewData(t, f))

	request := planMigrationReviewRequest(data)
	request["apply"] = true
	request["group"] = f.createRequest()["group"]
	review := performJSONRequest(t, f.router, http.MethodPost, planMigrationReviewPath, request)
	require.Equal(t, http.StatusUnprocessableEntity, review.Code, review.Body.String())
	body := decodeJSONBody(t, review)
	assertWriteGroupError(t, body, "WRITE_GROUP_PLAN_MIGRATION_BLOCKED")
	require.Equal(t, "open_write_groups", body["error"].(map[string]any)["action"])
	saved, err := recordingplan.NewSQLRepository(f.db).GetPlanByWorkspace(t.Context(), original.ID, f.record.ID)
	require.NoError(t, err)
	require.Equal(t, original, saved)
	current, err := f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, f.record, current)
	for _, table := range []string{"write_groups", "write_group_migration_maps", "write_group_versions"} {
		var count int
		require.NoError(t, f.db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM "+table).Scan(&count))
		require.Zero(t, count)
	}
	var enabled int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT enabled FROM database_target_mappings WHERE id='legacy-A'`).Scan(&enabled))
	require.Equal(t, 1, enabled)
	_, err = os.Stat(targetPath)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestNewRouter_BasicPlanMigrationRejectsStaleSourceAndDestination(t *testing.T) {
	for name, change := range map[string]string{
		"plan":        `UPDATE recording_plans SET name='changed original intent' WHERE id='legacy-plan-A'`,
		"measurement": `UPDATE measurement_definitions SET unit='F' WHERE id='measurement-A'`,
		"source":      `UPDATE points SET name='changed source metadata' WHERE id='point-A'`,
		"destination": `UPDATE database_connectors SET default_write_interval_seconds=30 WHERE id='connector-A'`,
	} {
		t.Run(name, func(t *testing.T) {
			f, _, _ := newPlanMigrationRouterFixture(t)
			before := planMigrationPreviewData(t, f)
			_, err := f.db.ExecContext(t.Context(), change)
			require.NoError(t, err)
			after := planMigrationPreviewData(t, f)
			require.NotEqual(t, before["review_digest"], after["review_digest"])
			review := performJSONRequest(t, f.router, http.MethodPost, planMigrationReviewPath, planMigrationReviewRequest(before))
			require.Equal(t, http.StatusConflict, review.Code, review.Body.String())
			assertWriteGroupError(t, decodeJSONBody(t, review), "revision_mismatch")
		})
	}
}

func TestNewRouter_BasicPlanMigrationRejectsInvalidAndForeignRequests(t *testing.T) {
	f, _, _ := newPlanMigrationRouterFixture(t)
	for _, request := range []map[string]any{{}, {"workspace_id": f.record.ID}, {"workspace_id": f.record.ID, "source_ids": []string{" "}}} {
		response := performJSONRequest(t, f.router, http.MethodPost, planMigrationPreviewPath, request)
		require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
	}
	for _, request := range []map[string]any{
		{"workspace_id": "unknown-workspace", "source_ids": []string{"legacy-plan-A"}},
		{"workspace_id": f.record.ID, "source_ids": []string{"unknown-plan"}},
	} {
		response := performJSONRequest(t, f.router, http.MethodPost, planMigrationPreviewPath, request)
		require.Equal(t, http.StatusNotFound, response.Code, response.Body.String())
		assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_NOT_FOUND")
	}
	data := planMigrationPreviewData(t, f)
	for _, field := range []string{"expected_workspace_revision", "expected_connector_revision", "review_digest", "source_ids", "confirm_snapshot_conversion"} {
		request := planMigrationReviewRequest(data)
		delete(request, field)
		response := performJSONRequest(t, f.router, http.MethodPost, planMigrationReviewPath, request)
		require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
	}
	request := planMigrationReviewRequest(data)
	request["confirm_snapshot_conversion"] = false
	response := performJSONRequest(t, f.router, http.MethodPost, planMigrationReviewPath, request)
	require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
	request = planMigrationReviewRequest(data)
	request["expected_workspace_revision"] = "old-workspace"
	response = performJSONRequest(t, f.router, http.MethodPost, planMigrationReviewPath, request)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	request = planMigrationReviewRequest(data)
	request["expected_connector_revision"] = "old-connector"
	response = performJSONRequest(t, f.router, http.MethodPost, planMigrationReviewPath, request)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	_, err := f.db.ExecContext(t.Context(), `UPDATE recording_plans SET workspace_id='foreign-workspace' WHERE id='legacy-plan-A'`)
	require.NoError(t, err)
	response = performJSONRequest(t, f.router, http.MethodPost, planMigrationPreviewPath, map[string]any{"workspace_id": f.record.ID, "source_ids": []string{"legacy-plan-A"}})
	require.Equal(t, http.StatusNotFound, response.Code, response.Body.String())
	response = performJSONRequest(t, f.router, http.MethodPost, planMigrationReviewPath, planMigrationReviewRequest(data))
	require.Equal(t, http.StatusNotFound, response.Code, response.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_NOT_FOUND")
}

func TestNewRouter_BasicPlanMigrationHidesForeignMeasurement(t *testing.T) {
	f, _, _ := newPlanMigrationRouterFixture(t)
	_, err := f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET workspace_id='foreign-workspace' WHERE id='measurement-A'`)
	require.NoError(t, err)
	response := performJSONRequest(t, f.router, http.MethodPost, planMigrationPreviewPath, map[string]any{"workspace_id": f.record.ID, "source_ids": []string{"legacy-plan-A"}})
	require.Equal(t, http.StatusNotFound, response.Code, response.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_NOT_FOUND")
	require.NotContains(t, response.Body.String(), "foreign-workspace")
}
