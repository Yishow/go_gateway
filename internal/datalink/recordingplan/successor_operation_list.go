package recordingplan

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

// ReadSuccessorFirstStartsInTx reads every first-start in a workspace snapshot.
// Malformed, nonterminal or orphaned evidence fails closed instead of vanishing.
func ReadSuccessorFirstStartsInTx(ctx context.Context, tx *sql.Tx, workspaceID string) ([]*SuccessorFirstStartRecord, error) {
	if tx == nil || !successorIdentity(workspaceID) {
		return nil, ErrSchemaOperationResult
	}
	rows, err := tx.QueryContext(ctx, adaptPlaceholders(`SELECT operation_id,detail FROM managed_schema_operations
  WHERE workspace_id=$1 AND action=$2 ORDER BY operation_id`), workspaceID, SuccessorFirstStartAction)
	if err != nil {
		return nil, fmt.Errorf("list successor first starts: %w", err)
	}
	type candidate struct{ id, detail string }
	candidates := []candidate{}
	for rows.Next() {
		var c candidate
		if err = rows.Scan(&c.id, &c.detail); err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("read successor first start candidate: %w", err)
		}
		candidates = append(candidates, c)
	}
	readErr := rows.Err()
	closeErr := rows.Close()
	if readErr != nil {
		return nil, fmt.Errorf("read successor first starts: %w", readErr)
	}
	if closeErr != nil {
		return nil, fmt.Errorf("close successor first starts: %w", closeErr)
	}
	records := make([]*SuccessorFirstStartRecord, 0, len(candidates))
	for _, c := range candidates {
		canonical, err := successorJSONObject(json.RawMessage(c.detail))
		if err != nil {
			return nil, err
		}
		var detail successorOperationDetail
		decoder := json.NewDecoder(bytes.NewReader(canonical))
		decoder.DisallowUnknownFields()
		if err = decoder.Decode(&detail); err != nil {
			return nil, ErrSchemaOperationResult
		}
		if err = detail.normalize(); err != nil {
			return nil, err
		}
		if detail.Action != SuccessorFirstStartAction || detail.Scope.WorkspaceID != workspaceID || detail.id() != c.id {
			return nil, ErrSchemaOperationMismatch
		}
		record, err := ReadSuccessorFirstStartInTx(ctx, tx, workspaceID, detail.ReserveOperationID)
		if err != nil {
			return nil, err
		}
		if record == nil {
			return nil, ErrSchemaOperationMismatch
		}
		records = append(records, record)
	}
	return records, nil
}
