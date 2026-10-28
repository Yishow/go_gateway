package recordingplan

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

func writeSuccessorOperationInTx(ctx context.Context, tx *sql.Tx, detail successorOperationDetail) error {
	if tx == nil {
		return ErrSchemaOperationResult
	}
	if err := detail.normalize(); err != nil {
		return err
	}
	encoded, err := json.Marshal(detail)
	if err != nil {
		return fmt.Errorf("encode successor result: %w", err)
	}
	now := time.Now().UTC()
	id := detail.id()
	_, err = tx.ExecContext(ctx, adaptPlaceholders(claimSchemaOperationQuery),
		id, detail.Action+":"+id, detail.Scope.WorkspaceID, detail.scopeKey(), "", detail.Action, string(SchemaOperationSucceeded),
		0, "", "", "", now, now, now, detail.IntentDigest, "", "", "", string(encoded))
	if err != nil {
		return fmt.Errorf("record successor result: %w", err)
	}
	requestID := detail.Scope.RequestID
	if detail.Action == SuccessorFirstStartAction {
		requestID = detail.ReserveOperationID
	}
	existing, err := readSuccessorOperationInTx(ctx, tx, detail.Action, detail.Scope.WorkspaceID, requestID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrSchemaOperationMismatch
	}
	prior, err := json.Marshal(existing)
	if err != nil {
		return fmt.Errorf("encode existing successor result: %w", err)
	}
	if !bytes.Equal(encoded, prior) {
		return ErrSchemaOperationMismatch
	}
	return nil
}

func readSuccessorOperationInTx(ctx context.Context, tx *sql.Tx, action, workspaceID, requestID string) (*successorOperationDetail, error) {
	if tx == nil || !successorIdentity(workspaceID) || !successorIdentity(requestID) {
		return nil, ErrSchemaOperationResult
	}
	var id string
	switch action {
	case SuccessorReserveAction:
		id = SuccessorReserveOperationID(workspaceID, requestID)
	case SuccessorFirstStartAction:
		id = SuccessorFirstStartOperationID(workspaceID, requestID)
	default:
		return nil, ErrSchemaOperationMismatch
	}
	var rowToken, rowWorkspace, rowScope, rowAction, status, digest, encoded, owner, verified string
	var executed int
	var completed sql.NullTime
	err := tx.QueryRowContext(ctx, adaptPlaceholders(`SELECT token,workspace_id,scope_key,action,status,payload_digest,detail,owner,verified_digest,executed_statements,completed_at
  FROM managed_schema_operations WHERE operation_id=$1`), id).Scan(
		&rowToken, &rowWorkspace, &rowScope, &rowAction, &status, &digest, &encoded, &owner, &verified, &executed, &completed)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read successor result: %w", err)
	}
	if rowWorkspace != workspaceID || rowAction != action || rowToken != action+":"+id {
		return nil, ErrSchemaOperationMismatch
	}
	if status != string(SchemaOperationSucceeded) || owner != "" || verified != "" || executed != 0 || !completed.Valid {
		return nil, ErrSchemaOperationResult
	}
	// Detail is not a public DTO and may never fall back to a partial result.
	canonical, err := successorJSONObject(json.RawMessage(encoded))
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(canonical))
	dec.DisallowUnknownFields()
	var detail successorOperationDetail
	if err = dec.Decode(&detail); err != nil {
		return nil, ErrSchemaOperationResult
	}
	if err = dec.Decode(new(any)); err != io.EOF {
		return nil, ErrSchemaOperationResult
	}
	if err = detail.normalize(); err != nil {
		return nil, err
	}
	if detail.Action != action || detail.Scope.WorkspaceID != workspaceID || detail.id() != id || detail.scopeKey() != rowScope || detail.IntentDigest != digest {
		return nil, ErrSchemaOperationMismatch
	}
	return &detail, nil
}
