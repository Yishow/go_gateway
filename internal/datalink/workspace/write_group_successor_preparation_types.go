package workspace

import "time"

// WriteGroupSuccessorReviewRequest is an intent, not a saved group or DDL authority.
type WriteGroupSuccessorReviewRequest struct {
	WorkspaceID                          string      `json:"workspace_id"`
	ExpectedWorkspaceRevision            string      `json:"expected_workspace_revision"`
	LineageID                            string      `json:"lineage_id"`
	ExpectedLineageRevision              int64       `json:"expected_lineage_revision"`
	ExpectedPredecessorRevision          string      `json:"expected_predecessor_revision"`
	ExpectedPredecessorConnectorRevision string      `json:"expected_predecessor_connector_revision"`
	ExpectedTargetConnectorRevision      string      `json:"expected_target_connector_revision"`
	OperationID                          string      `json:"operation_id"`
	Boundary                             time.Time   `json:"boundary"`
	Group                                *WriteGroup `json:"group"`
}

// WriteGroupSuccessorEnvelope is the original frozen review identity. Current
// stage revisions are supplied separately and never replace these original CAS.
type WriteGroupSuccessorEnvelope struct {
	Version                string                           `json:"version"`
	PredecessorGroupID     string                           `json:"predecessor_group_id"`
	Generation             int64                            `json:"generation"`
	SourceAppliedRevision  string                           `json:"source_applied_revision"`
	SourceDefinitionDigest string                           `json:"source_definition_digest"`
	Request                WriteGroupSuccessorReviewRequest `json:"request"`
}

// WriteGroupSuccessorBlocker contains closed, operator-safe classifications.
type WriteGroupSuccessorBlocker struct {
	ReasonCode string `json:"reason_code"`
	NextAction string `json:"next_action"`
}

// WriteGroupSuccessorReview is advisory; reserve must revalidate in its transaction.
type WriteGroupSuccessorReview struct {
	Envelope       WriteGroupSuccessorEnvelope  `json:"envelope"`
	ReviewDigest   string                       `json:"review_digest"`
	History        *WriteGroupLifecycleHistory  `json:"history"`
	Blockers       []WriteGroupSuccessorBlocker `json:"blockers"`
	RequiredChecks []WriteGroupSuccessorBlocker `json:"required_checks"`
}

// WriteGroupSuccessorCurrentCAS authorizes a new preparation stage, not replay.
type WriteGroupSuccessorCurrentCAS struct {
	WorkspaceRevision            string `json:"workspace_revision"`
	LineageRevision              int64  `json:"lineage_revision"`
	PredecessorRevision          string `json:"predecessor_revision"`
	PredecessorConnectorRevision string `json:"predecessor_connector_revision"`
	GroupID                      string `json:"group_id"`
	GroupRevision                string `json:"group_revision"`
	TargetConnectorRevision      string `json:"target_connector_revision"`
}

// WriteGroupSuccessorReservationRequest carries the exact reviewed envelope.
type WriteGroupSuccessorReservationRequest struct {
	Envelope     WriteGroupSuccessorEnvelope `json:"envelope"`
	ReviewDigest string                      `json:"review_digest"`
}

// WriteGroupSuccessorStageRequest binds a schema stage to an existing reservation.
type WriteGroupSuccessorStageRequest struct {
	Envelope     WriteGroupSuccessorEnvelope   `json:"envelope"`
	ReviewDigest string                        `json:"review_digest"`
	Current      WriteGroupSuccessorCurrentCAS `json:"current"`
}

// WriteGroupSuccessorPreparation is a reserved draft, never active recording.
type WriteGroupSuccessorPreparation struct {
	WorkspaceRevision  string      `json:"workspace_revision"`
	LineageID          string      `json:"lineage_id"`
	LineageRevision    int64       `json:"lineage_revision"`
	Generation         int64       `json:"generation"`
	PredecessorGroupID string      `json:"predecessor_group_id"`
	OperationID        string      `json:"operation_id"`
	IntentDigest       string      `json:"intent_digest"`
	Group              *WriteGroup `json:"group"`
}
