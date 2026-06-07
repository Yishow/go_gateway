package audit

import "time"

// EventType identifies one operator-facing workspace audit event type.
type EventType string

const (
	// EventTypeWorkspaceActivation records a Studio V2 workspace activation attempt.
	EventTypeWorkspaceActivation EventType = "workspace_activation"
	// EventTypeDatabaseConfigSaved records a persisted workspace database connector change.
	EventTypeDatabaseConfigSaved EventType = "database_config_saved"
	// EventTypeDatabaseTargetSaved records a persisted workspace database target change.
	EventTypeDatabaseTargetSaved EventType = "database_target_saved"
)

// Result identifies the outcome of an audited event.
type Result string

const (
	// ResultSuccess means the audited event succeeded.
	ResultSuccess Result = "success"
	// ResultPartialSuccess means the audited event completed with mixed per-scope results.
	ResultPartialSuccess Result = "partial_success"
	// ResultFailure means the audited event failed.
	ResultFailure Result = "failure"
)

// Entry is one persisted workspace audit history record.
type Entry struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	EventType   EventType `json:"event_type"`
	Result      Result    `json:"result"`
	Scope       string    `json:"scope"`
	ReferenceID string    `json:"reference_id,omitempty"`
	Details     string    `json:"details,omitempty"`
	OccurredAt  time.Time `json:"occurred_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// RecordEvent is the validated input for creating an audit history entry.
type RecordEvent struct {
	WorkspaceID string
	EventType   EventType
	Result      Result
	Scope       string
	ReferenceID string
	Details     any
	OccurredAt  time.Time
}

// ListFilter scopes recent audit history queries.
type ListFilter struct {
	WorkspaceID string
	EventType   EventType
	Limit       int
}
