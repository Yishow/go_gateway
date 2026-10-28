package workspace

import "time"

// WriteGroupSuccessorActivationRequest selects a new boundary for a prepared reservation.
type WriteGroupSuccessorActivationRequest struct {
	Stage    WriteGroupSuccessorStageRequest `json:"stage"`
	Boundary time.Time                       `json:"boundary"`
}

// WriteGroupSuccessorActivationEnvelope binds local evidence, not permission to write.
type WriteGroupSuccessorActivationEnvelope struct {
	Version                string                               `json:"version"`
	ActivationKind         string                               `json:"activation_kind,omitempty"`
	Request                WriteGroupSuccessorActivationRequest `json:"request"`
	SourceAppliedRevision  string                               `json:"source_applied_revision"`
	SourceDefinitionDigest string                               `json:"source_definition_digest"`
	TargetDefinitionDigest string                               `json:"target_definition_digest"`
	CounterProofDigest     string                               `json:"counter_proof_digest"`
	SchemaRevision         string                               `json:"schema_revision"`
	SchemaDigest           string                               `json:"schema_digest"`
}

// WriteGroupSuccessorActivationReview is coherent local evidence without target I/O.
type WriteGroupSuccessorActivationReview struct {
	Envelope     WriteGroupSuccessorActivationEnvelope `json:"envelope"`
	ReviewDigest string                                `json:"review_digest"`
	History      *WriteGroupLifecycleHistory           `json:"history"`
	Blockers     []WriteGroupSuccessorBlocker          `json:"blockers"`
}

// WriteGroupSuccessorActivationConfirmation binds the reviewed local and external evidence.
type WriteGroupSuccessorActivationConfirmation struct {
	Envelope        WriteGroupSuccessorActivationEnvelope `json:"envelope"`
	ReviewDigest    string                                `json:"review_digest"`
	ReadinessDigest string                                `json:"readiness_digest"`
}

// WriteGroupSuccessorActivationReadiness is a fresh probe result, never a persisted grant.
type WriteGroupSuccessorActivationReadiness struct {
	Review          *WriteGroupSuccessorActivationReview `json:"review"`
	Readiness       *WriteGroupReadiness                 `json:"readiness"`
	ReadinessDigest string                               `json:"readiness_digest"`
	snapshot        *successorActivationSnapshot
}
