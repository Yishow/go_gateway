package workspace

import (
	"context"
	"errors"
	"time"

	"go-gateway/internal/datalink/recordingplan"
)

// ErrRecordingStartInvalid identifies an incomplete or conflicting start intent.
var ErrRecordingStartInvalid = errors.New("recording start intent is invalid")

// ErrRecordingStartBusy reports another live start owner without disabling output.
var ErrRecordingStartBusy = errors.New("recording start scope is busy")

// ErrRecordingStartShareNotReady preserves Share's fail-closed alignment gate.
var ErrRecordingStartShareNotReady = errors.New("recording start share projection is not ready")

// RecordingStartGroupIntent carries one saved group and an optional CAS draft.
type RecordingStartGroupIntent struct {
	GroupID                   string      `json:"group_id"`
	ExpectedGroupRevision     string      `json:"expected_group_revision"`
	ExpectedConnectorRevision string      `json:"expected_connector_revision"`
	Draft                     *WriteGroup `json:"draft,omitempty"`
}

// RecordingStartRequest identifies the operator's exact selected intent.
type RecordingStartRequest struct {
	RequestID                 string                      `json:"request_id"`
	WorkspaceID               string                      `json:"workspace_id"`
	ExpectedWorkspaceRevision string                      `json:"expected_workspace_revision"`
	DeviceIDs                 []string                    `json:"device_ids"`
	Groups                    []RecordingStartGroupIntent `json:"groups"`
	ReadinessToken            string                      `json:"readiness_token,omitempty"`
	SettingsRevision          string                      `json:"settings_revision,omitempty"`
	WorkspaceRevision         string                      `json:"workspace_revision,omitempty"`
}

// RecordingStartGroupProgress distinguishes saving, readiness and Apply facts.
type RecordingStartGroupProgress struct {
	GroupID         string `json:"group_id"`
	GroupRevision   string `json:"group_revision"`
	AppliedRevision string `json:"applied_revision,omitempty"`
	Saved           bool   `json:"saved"`
	Ready           bool   `json:"ready"`
	Applied         bool   `json:"applied"`
}

// RecordingStartDeviceProgress records activation, independently from SQL delivery.
type RecordingStartDeviceProgress struct {
	DeviceID  string `json:"device_id"`
	Activated bool   `json:"activated"`
	Reason    string `json:"reason,omitempty"`
}

// RecordingStartOperation is the safe view of the existing ledger's start action.
type RecordingStartOperation struct {
	OperationID   string                              `json:"operation_id"`
	Action        string                              `json:"action"`
	Status        recordingplan.SchemaOperationStatus `json:"status"`
	IntentDigest  string                              `json:"intent_digest"`
	Reason        string                              `json:"reason,omitempty"`
	NextAction    string                              `json:"next_action,omitempty"`
	WorkspaceID   string                              `json:"workspace_id"`
	SetupRevision string                              `json:"setup_revision"`
	DeviceIDs     []string                            `json:"device_ids"`
	Groups        []RecordingStartGroupProgress       `json:"groups"`
	Devices       []RecordingStartDeviceProgress      `json:"devices"`
	Stage         string                              `json:"stage"`
	CreatedAt     time.Time                           `json:"created_at"`
	UpdatedAt     time.Time                           `json:"updated_at"`
}

type recordingStartActivator interface {
	ActivateScope(context.Context, []string, []string) (*ActivationResponse, error)
}

type recordingStartGroupScope struct {
	SourceDigest   string `json:"source_digest"`
	SchemaRevision string `json:"schema_revision"`
	SchemaDigest   string `json:"schema_digest"`
}

type recordingStartProgress struct {
	Version           int                        `json:"version"`
	BarrierVerified   bool                       `json:"barrier_verified"`
	Request           RecordingStartRequest      `json:"request"`
	WorkspaceRevision string                     `json:"workspace_revision"`
	Scopes            []recordingStartGroupScope `json:"scopes"`
	DeviceDigests     []string                   `json:"device_digests"`
	View              RecordingStartOperation    `json:"view"`
}
