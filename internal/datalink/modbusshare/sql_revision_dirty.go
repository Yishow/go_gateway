package modbusshare

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/sqltx"
)

// Mark the currently observed row unknown in one transaction. Serializing the
// typed record intentionally removes any point-format pending ownership marker.
// Exact raw CAS never writes an older revision over a concurrent operation.
func (s *SQLWorkspaceRevisionStore) markDirtyCurrentRevision(ctx context.Context, workspaceID string) (resultErr error) {
	if s == nil || s.db == nil {
		return newStorageError("mark workspace dirty", fmt.Errorf("database is not configured"))
	}
	owned, err := sqltx.Begin(ctx, s.db, nil)
	if err != nil {
		return newStorageError("begin dirty transaction", err)
	}
	defer func() {
		if err := owned.Rollback(); err != nil && err != sql.ErrTxDone {
			resultErr = errors.Join(resultErr, newStorageError("rollback dirty transaction", err))
		}
	}()
	key := workspaceRevisionKey(workspaceID)
	var raw string
	err = owned.Tx.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key=?`, key).Scan(&raw)
	record := workspaceRevisionRecord{Revision: initialRevision, Dirty: true}
	if errors.Is(err, sql.ErrNoRows) {
		if err = insertWorkspaceRevision(ctx, owned.Tx, key, record); err != nil {
			return err
		}
	} else {
		if err != nil {
			return newStorageError("read workspace revision", err)
		}
		if err = json.Unmarshal([]byte(raw), &record); err != nil {
			return newStorageError("decode workspace revision", err)
		}
		record.Dirty = true
		payload, err := json.Marshal(record)
		if err != nil {
			return newStorageError("encode workspace revision", err)
		}
		update, err := owned.Tx.ExecContext(ctx, `UPDATE system_settings SET value=?,updated_at=? WHERE key=? AND value=?`, string(payload), time.Now().UTC(), key, raw)
		if err != nil {
			return newStorageError("mark current workspace revision dirty", err)
		}
		count, err := update.RowsAffected()
		if err != nil {
			return newStorageError("mark workspace dirty rows", err)
		}
		if count != 1 {
			return NewError(ErrCodeRevisionConflict, "workspace changed while marking unknown", true)
		}
	}
	if err = owned.Commit(); err != nil {
		return newStorageError("commit dirty workspace revision", err)
	}
	return nil
}
