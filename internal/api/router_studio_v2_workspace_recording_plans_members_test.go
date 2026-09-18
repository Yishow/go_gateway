package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

const recordingPlansPath = "/api/v1/datalink/studio-v2/workspace/recording-plans"

// recordingMembershipStub mirrors the workspace-scoped measurement reader and
// the device-scoped point reader used by recording-plan membership checks.
type recordingMembershipStub struct {
	definitions []measurement.MeasurementDefinition
	points      []*schema.Point
	listErr     error
}

func newRecordingMembershipStub(workspaceID string) *recordingMembershipStub {
	definition := func(id, workspace, device, pointID, equipment string) measurement.MeasurementDefinition {
		return measurement.MeasurementDefinition{ID: id, WorkspaceID: workspace, DeviceID: device, PointID: pointID, EquipmentID: equipment}
	}
	return &recordingMembershipStub{
		definitions: []measurement.MeasurementDefinition{
			definition("meas-a", workspaceID, "device-a", "point-a", "equipment-a"),
			definition("meas-b", workspaceID, "device-b", "point-b", "equipment-b"),
			definition("meas-disabled", workspaceID, "device-a", "point-disabled", "equipment-a"),
			definition("meas-foreign-device", workspaceID, "device-x", "point-x", "equipment-x"),
			definition("meas-other-workspace", "workspace-other", "device-a", "point-a", "equipment-a"),
			definition("meas-kw", workspaceID, "device-a", "point-kw", "equipment-meter"),
			definition("measurement-1", workspaceID, "device-a", "point-1", "equipment-line"),
			definition("measurement-2", workspaceID, "device-a", "point-2", "equipment-line"),
		},
		points: []*schema.Point{
			{ID: "point-a", DeviceID: "device-a", Enabled: true},
			{ID: "point-b", DeviceID: "device-b", Enabled: true},
			{ID: "point-disabled", DeviceID: "device-a", Enabled: false},
			{ID: "point-x", DeviceID: "device-x", Enabled: true},
			{ID: "point-kw", DeviceID: "device-a", Enabled: true},
			{ID: "point-1", DeviceID: "device-a", Enabled: true},
			{ID: "point-2", DeviceID: "device-a", Enabled: true},
		},
	}
}

func (s *recordingMembershipStub) ListByWorkspace(_ context.Context, workspaceID string) ([]measurement.MeasurementDefinition, error) {
	if s.listErr != nil {
		return nil, s.listErr
	}
	result := make([]measurement.MeasurementDefinition, 0, len(s.definitions))
	for _, definition := range s.definitions {
		if definition.WorkspaceID == workspaceID {
			result = append(result, definition)
		}
	}
	return result, nil
}

func (s *recordingMembershipStub) List(_ context.Context, filter point.ListFilter) ([]*schema.Point, error) {
	result := make([]*schema.Point, 0, len(s.points))
	for _, candidate := range s.points {
		if filter.DeviceID != nil && candidate.DeviceID == *filter.DeviceID {
			result = append(result, candidate)
		}
	}
	return result, nil
}

const twoDeviceDestination = `"destinations":[{"destination_id":"destination-1","connector_id":"conn-c1","connector_revision":"identity-c1"}]`

func TestStudioV2WorkspaceRecordingPlans_MultiDeviceMembershipRetainsEquipment(t *testing.T) {
	router, _, planRepo, _ := newStudioV2WorkspaceRecordingPlanRouter(t)
	body := `{"id":"plan-two-devices","name":"Two devices","members":[` +
		`{"member_id":"member-a","measurement_id":"meas-a","equipment_id":"equipment-a"},` +
		`{"member_id":"member-b","measurement_id":"meas-b","equipment_id":"equipment-b"}],` +
		`"streams":[{"stream_id":"stream-a","measurement_id":"meas-a","equipment_id":"equipment-a","mode":"raw_history"},` +
		`{"stream_id":"stream-b","measurement_id":"meas-b","equipment_id":"equipment-b","mode":"raw_history"}],` + twoDeviceDestination + `}`

	response := serveRecordingScopedRequest(router, http.MethodPost, recordingPlansPath, body)
	if response.Code != http.StatusCreated {
		t.Fatalf("two-device plan create = %d, want 201: %s", response.Code, response.Body.String())
	}
	saved, err := planRepo.GetPlanByID(context.Background(), "plan-two-devices")
	if err != nil {
		t.Fatalf("read saved plan: %v", err)
	}
	want := []recordingplan.PlanMember{
		{MemberID: "member-a", MeasurementID: "meas-a", EquipmentID: "equipment-a"},
		{MemberID: "member-b", MeasurementID: "meas-b", EquipmentID: "equipment-b"},
	}
	if !reflect.DeepEqual(saved.Members, want) {
		t.Fatalf("saved members = %+v, want %+v", saved.Members, want)
	}
}

func TestStudioV2WorkspaceRecordingPlans_InvalidMembershipBlocksCreate(t *testing.T) {
	memberA := `{"member_id":"member-a","measurement_id":"meas-a","equipment_id":"equipment-a"}`
	streamA := `{"stream_id":"stream-a","measurement_id":"meas-a","mode":"raw_history"}`
	tests := []struct {
		name    string
		members string
		streams string
	}{
		{name: "missing measurement", members: `{"member_id":"member-m","measurement_id":"meas-missing","equipment_id":"equipment-a"}`, streams: `{"stream_id":"stream-m","measurement_id":"meas-missing","mode":"raw_history"}`},
		{name: "disabled point", members: `{"member_id":"member-d","measurement_id":"meas-disabled","equipment_id":"equipment-a"}`, streams: `{"stream_id":"stream-d","measurement_id":"meas-disabled","mode":"raw_history"}`},
		{name: "other workspace measurement", members: `{"member_id":"member-o","measurement_id":"meas-other-workspace","equipment_id":"equipment-a"}`, streams: `{"stream_id":"stream-o","measurement_id":"meas-other-workspace","mode":"raw_history"}`},
		{name: "device outside workspace", members: `{"member_id":"member-x","measurement_id":"meas-foreign-device","equipment_id":"equipment-x"}`, streams: `{"stream_id":"stream-x","measurement_id":"meas-foreign-device","mode":"raw_history"}`},
		{name: "equipment mismatch", members: `{"member_id":"member-a","measurement_id":"meas-a","equipment_id":"equipment-b"}`, streams: streamA},
		{name: "stream without member", members: memberA, streams: streamA + `,{"stream_id":"stream-b","measurement_id":"meas-b","mode":"raw_history"}`},
		{name: "no members", members: ``, streams: streamA},
		{name: "duplicate measurement", members: memberA + `,{"member_id":"member-a2","measurement_id":"meas-a","equipment_id":"equipment-a"}`, streams: streamA},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _, planRepo, _ := newStudioV2WorkspaceRecordingPlanRouter(t)
			body := `{"id":"plan-invalid-members","name":"Invalid members","members":[` + tt.members + `],"streams":[` + tt.streams + `],` + twoDeviceDestination + `}`

			response := serveRecordingScopedRequest(router, http.MethodPost, recordingPlansPath, body)
			assertRecordingPlanError(t, response, http.StatusUnprocessableEntity, "RECORDING_PLAN_MEMBERS_INVALID", false)
			if _, err := planRepo.GetPlanByID(context.Background(), "plan-invalid-members"); err == nil {
				t.Fatal("invalid membership must not persist a plan")
			}
		})
	}
}

func TestStudioV2WorkspaceRecordingPlans_InvalidMembershipBlocksUpdate(t *testing.T) {
	router, _, planRepo, _ := newStudioV2WorkspaceRecordingPlanRouter(t)
	createBody := `{"id":"plan-update-members","name":"Before","members":[{"member_id":"member-a","measurement_id":"meas-a","equipment_id":"equipment-a"}],` +
		`"streams":[{"stream_id":"stream-a","measurement_id":"meas-a","mode":"raw_history"}],` + twoDeviceDestination + `}`
	if response := serveRecordingScopedRequest(router, http.MethodPost, recordingPlansPath, createBody); response.Code != http.StatusCreated {
		t.Fatalf("seed plan = %d: %s", response.Code, response.Body.String())
	}

	updateBody := `{"name":"After","members":[{"member_id":"member-m","measurement_id":"meas-missing","equipment_id":"equipment-a"}],` +
		`"streams":[{"stream_id":"stream-m","measurement_id":"meas-missing","mode":"raw_history"}],` + twoDeviceDestination + `}`
	response := serveRecordingScopedRequest(router, http.MethodPut, recordingPlansPath+"/plan-update-members", updateBody)
	assertRecordingPlanError(t, response, http.StatusUnprocessableEntity, "RECORDING_PLAN_MEMBERS_INVALID", false)

	saved, err := planRepo.GetPlanByID(context.Background(), "plan-update-members")
	if err != nil {
		t.Fatalf("read unchanged plan: %v", err)
	}
	if saved.Name != "Before" || saved.Revision != "rev-1" || len(saved.Members) != 1 || saved.Members[0].MeasurementID != "meas-a" {
		t.Fatalf("rejected update changed the plan: %+v", saved)
	}
}

func TestStudioV2WorkspaceRecordingPlans_MembershipReadersFailClosed(t *testing.T) {
	body := `{"id":"plan-unavailable","name":"Unavailable","members":[{"member_id":"member-a","measurement_id":"meas-a","equipment_id":"equipment-a"}],` +
		`"streams":[{"stream_id":"stream-a","measurement_id":"meas-a","mode":"raw_history"}],` + twoDeviceDestination + `}`
	tests := []struct {
		name string
		wire func(*handlers.StudioV2WorkspaceRecordingPlansHandler, string)
	}{
		{name: "readers not wired", wire: func(*handlers.StudioV2WorkspaceRecordingPlansHandler, string) {}},
		{name: "measurement read fails", wire: func(handler *handlers.StudioV2WorkspaceRecordingPlansHandler, workspaceID string) {
			membership := newRecordingMembershipStub(workspaceID)
			membership.listErr = errors.New("private measurement store diagnostic")
			handler.SetRecordingMembershipServices(membership, membership)
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
			record, err := workspaceSvc.BindDatabaseConnector(context.Background(), "conn-c1")
			if err != nil {
				t.Fatalf("bind connector fixture: %v", err)
			}
			if _, err := workspaceSvc.AttachDevice(context.Background(), "device-a"); err != nil {
				t.Fatalf("attach device fixture: %v", err)
			}
			planRepo := recordingplan.NewMemoryRepository()
			handler := handlers.NewStudioV2WorkspaceRecordingPlansHandler(workspaceSvc, recordingplan.NewService(planRepo), recordingPlanTestResolver{connector: &schema.DatabaseConnector{
				ID: "conn-c1", Kind: schema.DatabaseConnectorKindSQLite, IdentityRevision: "identity-c1", Enabled: true,
			}})
			tt.wire(handler, record.ID)
			router := gin.New()
			router.POST(recordingPlansPath, handler.Create)

			response := serveRecordingScopedRequest(router, http.MethodPost, recordingPlansPath, body)
			assertRecordingPlanError(t, response, http.StatusServiceUnavailable, "RECORDING_PLAN_UNAVAILABLE", true)
			if _, err := planRepo.GetPlanByID(context.Background(), "plan-unavailable"); err == nil {
				t.Fatal("unavailable membership readers must not persist a plan")
			}
		})
	}
}

func TestStudioV2WorkspaceRecordingPlans_MissingPlanUsesSameSafeNotFound(t *testing.T) {
	router, _, _, _ := newStudioV2WorkspaceRecordingPlanRouter(t)
	updateBody := `{"name":"Missing","members":[{"member_id":"member-a","measurement_id":"meas-a","equipment_id":"equipment-a"}],` +
		`"streams":[{"stream_id":"stream-a","measurement_id":"meas-a","mode":"raw_history"}],` + twoDeviceDestination + `}`

	assertScopedNotFound(t, serveRecordingScopedRequest(router, http.MethodGet, recordingPlansPath+"/plan-missing", ""))
	assertScopedNotFound(t, serveRecordingScopedRequest(router, http.MethodPut, recordingPlansPath+"/plan-missing", updateBody))
	assertScopedNotFound(t, serveRecordingScopedRequest(router, http.MethodDelete, recordingPlansPath+"/plan-missing", ""))
}

func assertRecordingPlanError(t *testing.T, response *httptest.ResponseRecorder, status int, code string, retryable bool) {
	t.Helper()
	var payload struct {
		Success bool                           `json:"success"`
		Error   handlers.TypedAPIErrorEnvelope `json:"error"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error response: %v; body=%s", err, response.Body.String())
	}
	if response.Code != status || payload.Success || payload.Error.Code != code || payload.Error.Retryable != retryable || payload.Error.RequestID == "" {
		t.Fatalf("response = %d %s, want %d code=%s retryable=%v", response.Code, response.Body.String(), status, code, retryable)
	}
	if strings.Contains(response.Body.String(), "diagnostic") {
		t.Fatalf("private diagnostic leaked into response: %s", response.Body.String())
	}
}
