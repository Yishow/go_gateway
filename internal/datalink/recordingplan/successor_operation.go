package recordingplan

import (
	"context"
	"database/sql"
	"encoding/json"
)

// SuccessorOperationScope is the exact local lineage and generation identity.
// It is private service evidence, not an operator response.
type SuccessorOperationScope struct {
	WorkspaceID        string `json:"-"`
	RequestID          string `json:"-"`
	LineageID          string `json:"-"`
	PredecessorGroupID string `json:"-"`
	SuccessorGroupID   string `json:"-"`
}

// SuccessorReserveRecord preserves the original normalized reservation and
// its immutable result. Owning workspace code validates the inner payloads.
type SuccessorReserveRecord struct {
	Scope               SuccessorOperationScope `json:"-"`
	IntentDigest        string                  `json:"-"`
	Intent              json.RawMessage         `json:"-"`
	Result              json.RawMessage         `json:"-"`
	SchemaSourceBinding json.RawMessage         `json:"-"`
	SourceEvidence      json.RawMessage         `json:"-"`
}

// SuccessorFirstStartRecord preserves one first activation of a reserved root.
type SuccessorFirstStartRecord struct {
	Scope              SuccessorOperationScope `json:"-"`
	ReserveOperationID string                  `json:"-"`
	IntentDigest       string                  `json:"-"`
	Intent             json.RawMessage         `json:"-"`
	Result             json.RawMessage         `json:"-"`
}

// SuccessorReserveOperationID is independent of the frozen intent digest.
func SuccessorReserveOperationID(workspaceID, requestID string) string {
	return successorOperationID("successor-reserve/v1", workspaceID, requestID)
}

// SuccessorFirstStartOperationID is independent of the confirmed boundary.
func SuccessorFirstStartOperationID(workspaceID, reserveOperationID string) string {
	return successorOperationID("successor-first-start/v1", workspaceID, reserveOperationID)
}

// RecordSuccessorReserveInTx records reservation success in the caller's commit.
// Invoke it only after the new workspace revision is known. An existing result
// is immutable: complete intent, result and evidence must all match.
func RecordSuccessorReserveInTx(ctx context.Context, tx *sql.Tx, record *SuccessorReserveRecord) (*SuccessorReserveRecord, error) {
	if record == nil {
		return nil, ErrSchemaOperationResult
	}
	detail := successorOperationDetail{
		Version: successorDetailVersion, Action: SuccessorReserveAction,
		Scope: storedSuccessorScope(record.Scope), IntentDigest: record.IntentDigest,
		Intent: record.Intent, Result: record.Result, SchemaSourceBinding: record.SchemaSourceBinding, SourceEvidence: record.SourceEvidence,
	}
	if err := writeSuccessorOperationInTx(ctx, tx, detail); err != nil {
		return nil, err
	}
	return ReadSuccessorReserveInTx(ctx, tx, record.Scope.WorkspaceID, record.Scope.RequestID)
}

// ReadSuccessorReserveInTx reads an original request; an absent record is nil.
// It grants no mutation, schema execution or intake authority.
func ReadSuccessorReserveInTx(ctx context.Context, tx *sql.Tx, workspaceID, requestID string) (*SuccessorReserveRecord, error) {
	detail, err := readSuccessorOperationInTx(ctx, tx, SuccessorReserveAction, workspaceID, requestID)
	if err != nil || detail == nil {
		return nil, err
	}
	return &SuccessorReserveRecord{
		Scope: detail.Scope.public(), IntentDigest: detail.IntentDigest, Intent: detail.Intent, Result: detail.Result,
		SchemaSourceBinding: detail.SchemaSourceBinding, SourceEvidence: detail.SourceEvidence,
	}, nil
}

// RecordSuccessorFirstStartInTx records first activation in the caller's commit.
// Its reserve root must already exist with the identical generation scope.
func RecordSuccessorFirstStartInTx(ctx context.Context, tx *sql.Tx, record *SuccessorFirstStartRecord) (*SuccessorFirstStartRecord, error) {
	if record == nil || record.ReserveOperationID != SuccessorReserveOperationID(record.Scope.WorkspaceID, record.Scope.RequestID) {
		return nil, ErrSchemaOperationMismatch
	}
	root, err := ReadSuccessorReserveInTx(ctx, tx, record.Scope.WorkspaceID, record.Scope.RequestID)
	if err != nil {
		return nil, err
	}
	if root == nil || root.Scope != record.Scope {
		return nil, ErrSchemaOperationMismatch
	}
	detail := successorOperationDetail{
		Version: successorDetailVersion, Action: SuccessorFirstStartAction,
		Scope: storedSuccessorScope(record.Scope), ReserveOperationID: record.ReserveOperationID,
		IntentDigest: record.IntentDigest, Intent: record.Intent, Result: record.Result,
	}
	if err := writeSuccessorOperationInTx(ctx, tx, detail); err != nil {
		return nil, err
	}
	return ReadSuccessorFirstStartInTx(ctx, tx, record.Scope.WorkspaceID, record.ReserveOperationID)
}

// ReadSuccessorFirstStartInTx reads first activation for the exact reserve root.
// The returned evidence is private and must be projected into an operator DTO.
func ReadSuccessorFirstStartInTx(ctx context.Context, tx *sql.Tx, workspaceID, reserveOperationID string) (*SuccessorFirstStartRecord, error) {
	detail, err := readSuccessorOperationInTx(ctx, tx, SuccessorFirstStartAction, workspaceID, reserveOperationID)
	if err != nil || detail == nil {
		return nil, err
	}
	root, err := ReadSuccessorReserveInTx(ctx, tx, workspaceID, detail.Scope.RequestID)
	if err != nil {
		return nil, err
	}
	if root == nil || root.Scope != detail.Scope.public() {
		return nil, ErrSchemaOperationMismatch
	}
	return &SuccessorFirstStartRecord{
		Scope: detail.Scope.public(), ReserveOperationID: detail.ReserveOperationID,
		IntentDigest: detail.IntentDigest, Intent: detail.Intent, Result: detail.Result,
	}, nil
}
