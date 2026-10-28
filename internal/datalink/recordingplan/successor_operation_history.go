package recordingplan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ReadSuccessorFirstStartsForLineageInTx validates all local first-starts and
// rejects foreign evidence naming the exact lineage, a generation, or its saved
// reserve root. Foreign evidence is only used to reject, never to grant authority.
func ReadSuccessorFirstStartsForLineageInTx(ctx context.Context, tx *sql.Tx, workspaceID, lineageID string, groupIDs []string) ([]*SuccessorFirstStartRecord, error) {
	if !successorIdentity(lineageID) || len(groupIDs) == 0 {
		return nil, ErrSchemaOperationResult
	}
	groups := make(map[string]bool, len(groupIDs))
	for _, id := range groupIDs {
		if !successorIdentity(id) {
			return nil, ErrSchemaOperationResult
		}
		groups[id] = true
	}
	local, err := ReadSuccessorFirstStartsInTx(ctx, tx, workspaceID)
	if err != nil {
		return nil, err
	}
	rows, err := tx.QueryContext(ctx, adaptPlaceholders(`SELECT detail FROM managed_schema_operations
  WHERE workspace_id<>$1 AND action=$2`), workspaceID, SuccessorFirstStartAction)
	if err != nil {
		return nil, fmt.Errorf("read foreign successor first starts: %w", err)
	}
	roots := map[string]bool{}
	for rows.Next() {
		var raw string
		if err = rows.Scan(&raw); err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("read foreign successor evidence: %w", err)
		}
		facts := successorHistoryAttribution(raw, lineageID, groups)
		if facts.attributed {
			_ = rows.Close()
			return nil, ErrSchemaOperationMismatch
		}
		for id := range facts.roots {
			roots[id] = true
		}
	}
	readErr := rows.Err()
	closeErr := rows.Close()
	if readErr != nil {
		return nil, fmt.Errorf("read foreign successor histories: %w", readErr)
	}
	if closeErr != nil {
		return nil, fmt.Errorf("close foreign successor histories: %w", closeErr)
	}
	// A fully namespaced reserve ID, unlike the user's reusable request ID, can
	// identify an existing local root. Its scope must also name this lineage.
	for rootID := range roots {
		var raw string
		err = tx.QueryRowContext(ctx, adaptPlaceholders(`SELECT detail FROM managed_schema_operations
   WHERE operation_id=$1 AND workspace_id=$2 AND action=$3`), rootID, workspaceID, SuccessorReserveAction).Scan(&raw)
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read successor history root: %w", err)
		}
		if successorHistoryAttribution(raw, lineageID, groups).attributed {
			return nil, ErrSchemaOperationMismatch
		}
	}
	return local, nil
}
