package modbusshare

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

func workspaceRevisionKey(workspaceID string) string {
	return workspaceRevisionKeyPrefix + workspaceID
}

func desiredMappingsKey(workspaceID string) string { return desiredMappingsKeyPrefix + workspaceID }

type workspaceRevisionExecer interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func insertWorkspaceRevision(ctx context.Context, execer workspaceRevisionExecer, key string, record workspaceRevisionRecord) error {
	payload, err := json.Marshal(record)
	if err != nil {
		return newStorageError("encode workspace revision", err)
	}
	result, err := execer.ExecContext(ctx, `INSERT INTO system_settings (key, value, updated_at) VALUES (?, ?, ?)`, key, string(payload), time.Now().UTC())
	if err != nil {
		return newStorageError("insert workspace revision", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return newStorageError("insert workspace revision rows affected", err)
	}
	if rows != 1 {
		return newStorageError("insert workspace revision rows affected", fmt.Errorf("unexpected affected row count %d", rows))
	}
	return nil
}

func upsertSystemSetting(ctx context.Context, execer workspaceRevisionExecer, key, value string) error {
	result, err := execer.ExecContext(ctx, `UPDATE system_settings SET value = ?, updated_at = ? WHERE key = ?`, value, time.Now().UTC(), key)
	if err != nil {
		return newStorageError("update desired mappings", err)
	}
	rows, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		return newStorageError("update desired mappings rows affected", rowsErr)
	}
	if rows == 0 {
		insertResult, insertErr := execer.ExecContext(ctx, `INSERT INTO system_settings (key, value, updated_at) VALUES (?, ?, ?)`, key, value, time.Now().UTC())
		if insertErr != nil {
			return newStorageError("insert desired mappings", insertErr)
		}
		insertRows, insertRowsErr := insertResult.RowsAffected()
		if insertRowsErr != nil {
			return newStorageError("insert desired mappings rows affected", insertRowsErr)
		}
		if insertRows != 1 {
			return newStorageError("insert desired mappings rows affected", fmt.Errorf("unexpected affected row count %d", insertRows))
		}
		return nil
	}
	if rows != 1 {
		return newStorageError("update desired mappings rows affected", fmt.Errorf("unexpected affected row count %d", rows))
	}
	return nil
}
