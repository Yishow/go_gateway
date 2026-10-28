package workspace

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"math"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

// Private reserve evidence is not schema authority and never leaves the ledger.
type successorReserveSourceEvidence struct {
	Version     string                             `json:"version"`
	Source      *WriteGroup                        `json:"source"`
	SourceMySQL *dbtarget.MySQLRecordingCapability `json:"source_mysql,omitempty"`
	TargetMySQL *dbtarget.MySQLRecordingCapability `json:"target_mysql,omitempty"`
}

func decodeSuccessorReserveRecord(root *recordingplan.SuccessorReserveRecord) (*WriteGroupSuccessorEnvelope, *WriteGroupSuccessorPreparation, error) {
	if root == nil {
		return nil, nil, nil
	}
	var e WriteGroupSuccessorEnvelope
	var result WriteGroupSuccessorPreparation
	if decodeSuccessorJSON(root.Intent, &e) != nil || decodeSuccessorJSON(root.Result, &result) != nil {
		return nil, nil, ErrWriteGroupLineageCorrupt
	}
	normalized, _, digest, err := normalizeSuccessorEnvelope(e)
	if err != nil || !sameSuccessorJSON(e, normalized) || digest != root.IntentDigest || result.Group == nil {
		return nil, nil, ErrWriteGroupLineageCorrupt
	}
	if e.Generation == math.MaxInt64 || e.Request.ExpectedLineageRevision == math.MaxInt64 || result.Generation < 2 || result.LineageRevision < 2 {
		return nil, nil, ErrWriteGroupLineageCorrupt
	}
	scope := root.Scope
	if scope.WorkspaceID != e.Request.WorkspaceID || scope.RequestID != e.Request.OperationID || scope.LineageID != e.Request.LineageID || scope.PredecessorGroupID != e.PredecessorGroupID || scope.SuccessorGroupID != result.Group.ID ||
		result.Group.WorkspaceID != scope.WorkspaceID || result.OperationID != scope.RequestID || result.LineageID != scope.LineageID || result.PredecessorGroupID != scope.PredecessorGroupID || result.IntentDigest != digest || result.WorkspaceRevision == "" || result.Group.Revision == "" || result.Group.AppliedRevision != "" || result.Group.Status != WriteGroupStatusDraft || result.Generation != e.Generation+1 || result.LineageRevision != e.Request.ExpectedLineageRevision+1 {
		return nil, nil, ErrWriteGroupLineageCorrupt
	}
	// Server identity/revisions are the only differences permitted in the reserve result.
	actual := e
	actual.Request.Group = cloneWriteGroup(result.Group)
	if result.Group.isWideProfile() || result.Group.RecordingProfile == nil {
		actual.Request.Group.Destination.TableName = ""
	}
	_, _, actualDigest, err := normalizeSuccessorEnvelope(actual)
	if err != nil || actualDigest != digest {
		return nil, nil, ErrWriteGroupLineageCorrupt
	}
	if _, err := decodeSuccessorReserveSourceEvidence(root, &normalized, &result); err != nil {
		return nil, nil, err
	}
	return &normalized, &result, nil
}

func (s *WriteGroupService) recordedSuccessorReserveInTx(ctx context.Context, tx *sql.Tx, e WriteGroupSuccessorEnvelope, digest string) (*WriteGroupSuccessorPreparation, bool, error) {
	root, err := recordingplan.ReadSuccessorReserveInTx(ctx, tx, e.Request.WorkspaceID, e.Request.OperationID)
	if err != nil || root == nil {
		return nil, false, err
	}
	original, result, err := decodeSuccessorReserveRecord(root)
	if err != nil {
		return nil, false, err
	}
	if root.IntentDigest != digest || !sameSuccessorJSON(*original, e) {
		return nil, false, recreateReservationConflict("reservation_intent_changed")
	}
	return result, true, nil
}

func (s *WriteGroupService) checkpointSuccessorReserve(ctx context.Context, tx *sql.Tx, record *Record, e WriteGroupSuccessorEnvelope, digest string, result *WriteGroupSuccessorPreparation, evidence *successorReserveSourceEvidence) error {
	result.WorkspaceRevision = record.DatabaseSetupRevision
	intent, err := json.Marshal(e)
	if err != nil {
		return err
	}
	payload, err := json.Marshal(result)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(payload, result); err != nil {
		return err
	}
	source, err := json.Marshal(evidence)
	if err != nil {
		return err
	}
	binding, err := s.makeSuccessorReserveSchemaBinding(ctx, tx, result, evidence)
	if err != nil {
		return err
	}
	_, err = recordingplan.RecordSuccessorReserveInTx(ctx, tx, &recordingplan.SuccessorReserveRecord{Scope: recordingplan.SuccessorOperationScope{WorkspaceID: record.ID, RequestID: e.Request.OperationID, LineageID: result.LineageID, PredecessorGroupID: e.PredecessorGroupID, SuccessorGroupID: result.Group.ID}, IntentDigest: digest, Intent: intent, Result: payload, SchemaSourceBinding: binding, SourceEvidence: source})
	return err
}

func sameSuccessorJSON(a, b any) bool {
	left, e := json.Marshal(a)
	if e != nil {
		return false
	}
	right, e := json.Marshal(b)
	return e == nil && bytes.Equal(left, right)
}
func decodeSuccessorJSON(raw []byte, target any) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(target); err != nil {
		return err
	}
	if err := dec.Decode(new(any)); err != io.EOF {
		return ErrWriteGroupLineageCorrupt
	}
	return nil
}
