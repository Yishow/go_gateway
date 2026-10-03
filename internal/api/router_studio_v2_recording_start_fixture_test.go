package api

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type recordingStartActivationOutcome struct {
	response *workspace.ActivationResponse
	err      error
}

type recordingStartActivationStub struct {
	mu        sync.Mutex
	workspace string
	outcomes  []recordingStartActivationOutcome
	calls     int
}

func (s *recordingStartActivationStub) ActivateScope(_ context.Context, deviceIDs, _ []string) (*workspace.ActivationResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	index := s.calls - 1
	if len(s.outcomes) > 0 {
		if index >= len(s.outcomes) {
			index = len(s.outcomes) - 1
		}
		outcome := s.outcomes[index]
		if outcome.err != nil {
			return nil, outcome.err
		}
		if outcome.response != nil {
			response := *outcome.response
			if response.WorkspaceID == "" {
				response.WorkspaceID = s.workspace
			}
			response.Results = append([]workspace.ActivationResult(nil), outcome.response.Results...)
			return &response, nil
		}
	}
	response := &workspace.ActivationResponse{WorkspaceID: s.workspace, Results: make([]workspace.ActivationResult, 0, len(deviceIDs))}
	for _, deviceID := range deviceIDs {
		response.Results = append(response.Results, workspace.ActivationResult{
			DeviceID: deviceID, Status: workspace.ActivationResultStatusSuccess, Message: "stub activated",
		})
	}
	return response, nil
}

func (s *recordingStartActivationStub) Calls() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls
}

type recordingStartReadinessDeviceStub struct{}

func (recordingStartReadinessDeviceStub) GetByID(_ context.Context, id string) (*schema.Device, error) {
	if id == "" {
		return nil, errors.New("device id is empty")
	}
	return &schema.Device{ID: id}, nil
}

func (recordingStartReadinessDeviceStub) CheckReadiness(_ context.Context, id string) (*schema.DeviceReadiness, error) {
	if id == "" {
		return nil, errors.New("device id is empty")
	}
	return &schema.DeviceReadiness{
		DeviceID: id, ActivationAllowed: true, ApplyAllowed: true, PlanningAllowed: true,
		ConnectStatus: schema.ReadinessStageStatusSuccess, ProbeStatus: schema.ReadinessStageStatusSuccess,
	}, nil
}

type recordingStartFixture struct {
	managedGroupFixture
	start      *workspace.RecordingStartService
	activation *recordingStartActivationStub
	ledger     *recordingplan.Service
}

func newRecordingStartFixture(t *testing.T, outcomes ...recordingStartActivationOutcome) recordingStartFixture {
	t.Helper()
	managed := newManagedGroupFixture(t, "sqlite")
	token := managed.preview(t)
	confirm := managed.schemaRequest()
	confirm["token"], confirm["operation_id"] = token.Token, token.OperationID
	response := performJSONRequest(t, managed.router, http.MethodPost, managed.confirmPath(), confirm)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	reloaded := performJSONRequest(t, managed.router, http.MethodGet, writeGroupsPath+"/"+managed.group["id"].(string), nil)
	require.Equal(t, http.StatusOK, reloaded.Code, reloaded.Body.String())
	managed.data, managed.group = groupSaveData(t, decodeJSONBody(t, reloaded))
	return newRecordingStartServiceFixture(t, managed, outcomes...)
}

func newRecordingStartServiceFixture(t *testing.T, managed managedGroupFixture, outcomes ...recordingStartActivationOutcome) recordingStartFixture {
	t.Helper()
	managed.workspace.WithReadinessServices(recordingStartReadinessDeviceStub{}, nil, nil, nil).WithWriteGroupReadiness(managed.groups)
	activation := &recordingStartActivationStub{workspace: managed.record.ID, outcomes: outcomes}
	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(managed.db))
	start := workspace.NewRecordingStartService(managed.workspace, managed.groups, activation, ledger)
	return recordingStartFixture{managedGroupFixture: managed, start: start, activation: activation, ledger: ledger}
}

func (f recordingStartFixture) request(requestID string) workspace.RecordingStartRequest {
	return workspace.RecordingStartRequest{
		RequestID: requestID, WorkspaceID: f.record.ID,
		ExpectedWorkspaceRevision: f.data["workspace_revision"].(string), DeviceIDs: []string{"device-A"},
		Groups: []workspace.RecordingStartGroupIntent{{
			GroupID: f.group["id"].(string), ExpectedGroupRevision: f.group["revision"].(string),
			ExpectedConnectorRevision: "connector-1",
		}},
	}
}

func recordingStartSuccess() recordingStartActivationOutcome {
	return recordingStartActivationOutcome{response: &workspace.ActivationResponse{
		Results: []workspace.ActivationResult{{DeviceID: "device-A", Status: workspace.ActivationResultStatusSuccess}},
	}}
}

func recordingStartFailure() recordingStartActivationOutcome {
	return recordingStartActivationOutcome{response: &workspace.ActivationResponse{
		Results: []workspace.ActivationResult{{DeviceID: "device-A", Status: workspace.ActivationResultStatusFailed, Message: "stub failure"}},
	}}
}
