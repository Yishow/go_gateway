package modbusshare

import (
	"context"
	"time"

	"go-gateway/internal/datalink/schema"
)

// DesiredMapping describes a desired tag-to-register mapping in holding register space.
type DesiredMapping struct {
	WorkspaceID        string          `json:"workspace_id"`
	SourceRuleID       string          `json:"source_rule_id"`
	SourceRuleRevision string          `json:"source_rule_revision"`
	TagID              string          `json:"tag_id"`
	MappingID          string          `json:"mapping_id"`
	DataType           schema.DataType `json:"data_type"`
	ShareStartRegister uint32          `json:"share_start_register"` // Human 40001 domain
	ZeroBasedRegister  uint16          `json:"zero_based_register"`  // Zero-based index (40001 -> 0)
	SpanRegisters      int             `json:"span_registers"`
	StrideRegisters    int             `json:"stride_registers"`
	CapacityRegisters  int             `json:"capacity_registers"`
	TagKey             string          `json:"tag_key,omitempty"`
	DisplayName        string          `json:"display_name,omitempty"`
	OwnershipProof     *OwnershipProof `json:"ownership_proof,omitempty"`
}

// OwnershipProof records the persisted relationship used by the backend to
// authorize a canonical Share mapping. A tag id alone is never sufficient.
type OwnershipProof struct {
	Verified           bool   `json:"verified"`
	WorkspaceID        string `json:"workspace_id"`
	SourceRuleID       string `json:"source_rule_id"`
	SourceRuleRevision string `json:"source_rule_revision"`
	Basis              string `json:"basis"`
}

// SpanRange defines a contiguous register range [Start, Start+Count).
type SpanRange struct {
	TagID    string          `json:"tag_id"`
	Start    uint16          `json:"start"`
	Count    int             `json:"count"`
	DataType schema.DataType `json:"data_type"`
}

// Diagnostic represents an operator-safe diagnostic message.
type Diagnostic struct {
	Code      string `json:"code"`
	Severity  string `json:"severity"` // "info", "warning", "error"
	Message   string `json:"message"`
	TagID     string `json:"tag_id,omitempty"`
	RuleID    string `json:"rule_id,omitempty"`
	Retryable bool   `json:"retryable"`
	Action    string `json:"action,omitempty"`
}

func cloneDiagnostic(diagnostic *Diagnostic) *Diagnostic {
	if diagnostic == nil {
		return nil
	}
	clone := *diagnostic
	return &clone
}

// RecoveryStatus is safe, actionable state exposed when durable Share
// projection or readiness cannot be trusted.
type RecoveryStatus struct {
	Code      string `json:"code"`
	Retryable bool   `json:"retryable"`
	Action    string `json:"action"`
	RequestID string `json:"request_id,omitempty"`
}

// ReconcileRequest describes the input payload to reconcile a workspace's Share state.
type ReconcileRequest struct {
	WorkspaceID               string           `json:"workspace_id"`
	ExpectedWorkspaceRevision string           `json:"expected_workspace_revision"`
	ExpectedSettingsRevision  string           `json:"expected_settings_revision"`
	ReadinessToken            string           `json:"readiness_token"`
	DesiredMappings           []DesiredMapping `json:"desired_mappings"`
	// CanonicalPlanSignature identifies a server-produced candidate snapshot.
	// When present, reconcile uses the server snapshot instead of browser data.
	CanonicalPlanSignature string `json:"canonical_plan_signature,omitempty"`
	// Restore preserves the durable workspace revision while rebuilding the
	// process runtime projection after restart.
	Restore bool `json:"-"`
}

// ReconcileOutcome describes the result of a Share reconciliation.
type ReconcileOutcome struct {
	Outcome              string             `json:"outcome"` // "aligned", "applied", "disabled", "failed", "dirty_unknown", "invalidated_unknown"
	NewWorkspaceRevision string             `json:"new_workspace_revision"`
	NewReadinessToken    string             `json:"new_readiness_token,omitempty"`
	SettingsRevision     string             `json:"settings_revision"`
	AppliedCount         int                `json:"applied_count"`
	RemovedCount         int                `json:"removed_count"`
	InvalidatedCount     int                `json:"invalidated_count"`
	RemovedSpans         []SpanRange        `json:"removed_spans"`
	InvalidatedSpans     []SpanRange        `json:"invalidated_spans"`
	Mappings             []TagMirrorMapping `json:"mappings"`
	Diagnostics          []Diagnostic       `json:"diagnostics,omitempty"`
}

// Settings represents the persisted global Modbus Share configuration.
type Settings struct {
	Enabled           bool      `json:"enabled"`
	BindAddress       string    `json:"bind_address"`
	Port              int       `json:"port"`
	SlaveID           uint8     `json:"slave_id"`
	CapacityRegisters int       `json:"capacity_registers"`
	SettingsRevision  string    `json:"settings_revision"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// OwnershipChecker proves that a tag is durably related to a workspace.
// Share runtime state must never infer ownership from a tag id alone.
type OwnershipChecker func(context.Context, string, string) bool

// DesiredOwnershipChecker proves the complete persisted mapping identity.
type DesiredOwnershipChecker func(context.Context, DesiredMapping) error

// DefaultSettings returns safe default settings (disabled by default).
func DefaultSettings() Settings {
	return Settings{
		Enabled:           false,
		BindAddress:       defaultBindAddress,
		Port:              0,
		SlaveID:           1,
		CapacityRegisters: 32768,
		SettingsRevision:  "default",
		UpdatedAt:         time.Now().UTC(),
	}
}

// HydrationState constants.
const (
	HydrationStateReady              = "ready"
	HydrationStatePending            = "pending"
	HydrationStateFailed             = "failed"
	shareStateDisabled               = "disabled"
	shareStateRunning                = "running"
	shareStateStopped                = "stopped"
	defaultBindAddress               = "127.0.0.1"
	initialRevision                  = "rev-1"
	settingsRevisionConflictMessage  = "settings revision conflict"
	workspaceRevisionConflictMessage = "workspace revision conflict"
	dirtyUnknownState                = "dirty_unknown"
	retryReconcileAction             = "retry reconcile"
	failedMarkWorkspaceDirtyMessage  = "failed to mark workspace dirty"
	runRecoveryReconcileAction       = "run recovery reconcile"
)

// HydrationState represents the current workspace hydration status.
type HydrationState struct {
	State             string `json:"state"` // "ready", "pending", "failed"
	WorkspaceID       string `json:"workspace_id"`
	WorkspaceRevision string `json:"workspace_revision"`
	SettingsRevision  string `json:"settings_revision"`
	Readiness         bool   `json:"readiness"`
	ReadinessToken    string `json:"readiness_token,omitempty"`
}
