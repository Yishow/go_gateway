package point

import (
	"context"
	"database/sql"
	"encoding/json"
	"go-gateway/internal/datalink/sqltx"
	"time"
)

// FormatInvalidation binds the old live installation and newly committed intent.
type FormatInvalidation struct {
	WorkspaceID, ExpectedWorkspaceRevision, WorkspaceRevision string
	PointID, ExpectedPointRevision, PointRevision             string
}

// FormatInvalidator must fail closed only for matching installations. It must
// not apply the new decoder or remove a newer installation.
type FormatInvalidator interface {
	InvalidatePointFormat(context.Context, FormatInvalidation) error
}

// SetFormatInvalidator configures the production fail-closed seam at startup.
func (s *Service) SetFormatInvalidator(invalidator FormatInvalidator) {
	s.formatInvalidator = invalidator
}

type formatInvalidationStore interface {
	CompleteFormatInvalidation(context.Context, FormatMutationRequest, *FormatMutationResult) error
	InspectFormatMutation(context.Context, FormatScope) (*FormatMutationResult, error)
}

// InspectFormatMutation reads the existing receipt; it never retries a callback.
func (s *Service) InspectFormatMutation(ctx context.Context, scope FormatScope) (*FormatMutationResult, error) {
	if !formatRequestID.MatchString(scope.RequestID) || scope.WorkspaceID == "" || scope.DeviceID == "" || scope.PointID == "" {
		return nil, ErrFormatInvalid
	}
	store, ok := s.repo.(formatInvalidationStore)
	if !ok {
		return nil, ErrFormatUnavailable
	}
	return store.InspectFormatMutation(ctx, scope)
}

func (r *SQLRepository) CompleteFormatInvalidation(ctx context.Context, req FormatMutationRequest, result *FormatMutationResult) (resultErr error) {
	ownedTx, err := sqltx.Begin(ctx, r.db, nil)
	if err != nil {
		return ErrFormatUnavailable
	}
	tx := ownedTx.Tx
	defer func() {
		if err := ownedTx.Rollback(); err != nil && err != sql.ErrTxDone {
			resultErr = formatTransactionFailure(err)
		}
	}()
	snapshot, err := r.formatSnapshot(ctx, tx, req)
	if err != nil {
		return ErrFormatUnavailable
	}
	rev, rawWorkspace, err := loadFormatWorkspaceRevision(ctx, tx, req.WorkspaceID)
	if err != nil || snapshot.point.RevisionID != result.Point.RevisionID || rev.Revision != result.WorkspaceRevision {
		return ErrFormatUnavailable
	}
	key := formatReceiptKey(req.WorkspaceID, req.RequestID)
	var old string
	if tx.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key=?`, key).Scan(&old) != nil {
		return ErrFormatUnavailable
	}
	var receipt formatReceipt
	if json.Unmarshal([]byte(old), &receipt) != nil || !validFormatReceipt(receipt, req.FormatScope) || receipt.Result.Point.RevisionID != result.Point.RevisionID || receipt.Result.InvalidationStatus != "pending" {
		return ErrFormatUnavailable
	}
	if err := releaseOwnedFormatDirty(ctx, tx, req, result, rawWorkspace); err != nil {
		return err
	}
	receipt.Result.InvalidationStatus = "invalidated"
	raw, err := json.Marshal(receipt)
	if err != nil {
		return ErrFormatUnavailable
	}
	updated, err := tx.ExecContext(ctx, `UPDATE system_settings SET value=?,updated_at=? WHERE key=? AND value=?`, string(raw), time.Now().UTC(), key, old)
	if err != nil {
		return ErrFormatUnavailable
	}
	n, err := updated.RowsAffected()
	if err != nil || n != 1 {
		return ErrFormatUnavailable
	}
	if err := ownedTx.Commit(); err != nil {
		return formatTransactionFailure(err)
	}
	result.InvalidationStatus = "invalidated"
	return nil
}

func (r *SQLRepository) InspectFormatMutation(ctx context.Context, scope FormatScope) (out *FormatMutationResult, resultErr error) {
	ownedTx, err := sqltx.Begin(ctx, r.db, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, ErrFormatUnavailable
	}
	tx := ownedTx.Tx
	defer func() {
		if err := ownedTx.Rollback(); err != nil && err != sql.ErrTxDone {
			out = nil
			resultErr = formatTransactionFailure(err)
		}
	}()
	if _, err = r.formatSnapshot(ctx, tx, FormatMutationRequest{FormatScope: scope}); err != nil {
		return nil, err
	}
	var raw string
	if err = tx.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key=?`, formatReceiptKey(scope.WorkspaceID, scope.RequestID)).Scan(&raw); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPointNotFound
		}
		return nil, ErrFormatUnavailable
	}
	var receipt formatReceipt
	if json.Unmarshal([]byte(raw), &receipt) != nil || receipt.Version != 1 || receipt.WorkspaceID != scope.WorkspaceID || receipt.RequestID != scope.RequestID {
		return nil, ErrFormatUnavailable
	}
	if receipt.Result.Point.ID != scope.PointID || receipt.Result.Point.DeviceID != scope.DeviceID {
		return nil, ErrFormatScope
	}
	if !validFormatReceipt(receipt, scope) {
		return nil, ErrFormatUnavailable
	}
	return &receipt.Result, nil
}
