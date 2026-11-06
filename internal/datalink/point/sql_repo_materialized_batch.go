package point

import (
	"context"
	"database/sql"

	"go-gateway/internal/datalink/sqltx"
)

func (r *SQLRepository) updateMaterializedRevisions(ctx context.Context, updates []RevisionRestore) (resultErr error) {
	if len(updates) == 0 {
		return ErrFormatInvalid
	}
	owned, err := sqltx.Begin(ctx, r.db, nil)
	if err != nil {
		return ErrFormatUnavailable
	}
	defer func() {
		if err := owned.Rollback(); err != nil && err != sql.ErrTxDone { //nolint:errorlint // Only bare ErrTxDone is benign.
			resultErr = formatTransactionFailure(err)
		}
	}()
	tx := owned.Tx
	seen := make(map[string]bool, len(updates))
	deviceID := ""
	for _, update := range updates {
		p := update.Point
		if p == nil || p.ID == "" || update.ExpectedRevision == "" || seen[p.ID] {
			return ErrFormatInvalid
		}
		seen[p.ID] = true
		if deviceID == "" {
			deviceID = p.DeviceID
		} else if deviceID != p.DeviceID {
			return ErrFormatInvalid
		}
		if err := validateSQLOverrideUpdate(ctx, tx, p, update.ExpectedRevision); err != nil {
			return err
		}
	}
	// Point addresses cannot contain an underscore through the public API. The
	// temporary names are visible only within this transaction and are never
	// published to runtime or persisted on failure.
	for _, update := range updates {
		p := update.Point
		result, err := tx.ExecContext(ctx, `UPDATE points SET address=? WHERE id=? AND device_id=? AND revision_id=?`, "__source_move_"+p.ID, p.ID, p.DeviceID, update.ExpectedRevision)
		if err != nil {
			return ErrFormatUnavailable
		}
		count, err := result.RowsAffected()
		if err != nil {
			return ErrFormatUnavailable
		}
		if count != 1 {
			return ErrFormatConflict
		}
	}
	for _, update := range updates {
		if err := writePointConfigRevision(ctx, tx, update.Point, update.ExpectedRevision); err != nil {
			return err
		}
	}
	if err := owned.Commit(); err != nil {
		return formatTransactionFailure(err)
	}
	return nil
}
