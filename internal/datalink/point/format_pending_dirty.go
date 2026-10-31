package point

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

const formatPendingDirtyKey = "point_format_pending_v1"

type formatPendingDirty struct {
	Version           int    `json:"version"`
	ReceiptKey        string `json:"receipt_key"`
	WorkspaceID       string `json:"workspace_id"`
	PointID           string `json:"point_id"`
	PointRevision     string `json:"point_revision"`
	WorkspaceRevision string `json:"workspace_revision"`
}

func formatDirtyOwner(req FormatMutationRequest, result *FormatMutationResult) formatPendingDirty {
	return formatPendingDirty{Version: 1, ReceiptKey: formatReceiptKey(req.WorkspaceID, req.RequestID), WorkspaceID: req.WorkspaceID, PointID: req.PointID, PointRevision: result.Point.RevisionID, WorkspaceRevision: result.WorkspaceRevision}
}

// A receipt alone cannot own dirty state: another operation can mark the SAME
// revision unknown. Its authoritative writer revokes this same-row witness.
func releaseOwnedFormatDirty(ctx context.Context, tx *sql.Tx, req FormatMutationRequest, result *FormatMutationResult, raw string) error {
	var fields map[string]json.RawMessage
	if json.Unmarshal([]byte(raw), &fields) != nil {
		return ErrFormatUnavailable
	}
	var owner formatPendingDirty
	if json.Unmarshal(fields[formatPendingDirtyKey], &owner) != nil || owner != formatDirtyOwner(req, result) {
		return ErrFormatUnavailable
	}
	var dirty bool
	if json.Unmarshal(fields["dirty"], &dirty) != nil || !dirty {
		return ErrFormatUnavailable
	}
	delete(fields, formatPendingDirtyKey)
	fields["dirty"] = json.RawMessage("false")
	next, err := json.Marshal(fields)
	if err != nil {
		return ErrFormatUnavailable
	}
	updated, err := tx.ExecContext(ctx, `UPDATE system_settings SET value=?,updated_at=? WHERE key=? AND value=?`, string(next), time.Now().UTC(), "modbus_share_workspace_revision:"+req.WorkspaceID, raw)
	if err != nil {
		return ErrFormatUnavailable
	}
	count, err := updated.RowsAffected()
	if err != nil || count != 1 {
		return ErrFormatUnavailable
	}
	return nil
}
