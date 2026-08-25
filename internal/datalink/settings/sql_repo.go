// Package settings 提供系統設定管理功能的 SQL Repository 實現。
package settings

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
)

// =============================================================================
// SQL Repository 實現
// =============================================================================

// SQLRepository SQL 設定儲存庫
type SQLRepository struct {
	db *sql.DB
}

// NewSQLRepository 建立新的 SQL 儲存庫
func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

// Get 取得設定值
func (r *SQLRepository) Get(ctx context.Context, key string) (*SettingItem, error) {
	if r == nil || r.db == nil {
		return nil, storageError("get", fmt.Errorf("database is not configured"))
	}
	query := `
		SELECT key, value, description, updated_at
		FROM system_settings WHERE key = ?
	`

	row := r.db.QueryRowContext(ctx, query, key)

	var item SettingItem
	var valueJSON string
	var description sql.NullString
	var updatedAt sql.NullString

	err := row.Scan(&item.Key, &valueJSON, &description, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: 設定不存在: %s", ErrSettingNotFound, key)
	}
	if err != nil {
		return nil, storageError("get", fmt.Errorf("read setting: %w", err))
	}

	// 解析 JSON 值
	if err := json.Unmarshal([]byte(valueJSON), &item.Value); err != nil {
		return nil, storageError("decode", fmt.Errorf("decode setting %s: %w", key, err))
	}

	if description.Valid {
		item.Description = description.String
	}

	if updatedAt.Valid {
		parsedUpdatedAt, err := common.ParseTimeString(updatedAt.String)
		if err != nil {
			return nil, fmt.Errorf("解析更新時間失敗: %w", err)
		}
		item.UpdatedAt = parsedUpdatedAt
	}

	return &item, nil
}

// Set 設定值
func (r *SQLRepository) Set(ctx context.Context, key string, value interface{}) error {
	if r == nil || r.db == nil {
		return storageError("set", fmt.Errorf("database is not configured"))
	}
	// 序列化值為 JSON
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return storageError("encode", err)
	}

	// 使用 UPSERT 語法
	query := `
		INSERT INTO system_settings (key, value, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at
	`

	result, err := r.db.ExecContext(ctx, query, key, string(valueJSON), time.Now())
	if err != nil {
		return storageError("set", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return storageError("set rows affected", err)
	}
	if rows != 1 {
		return storageError("set rows affected", fmt.Errorf("unexpected affected row count %d", rows))
	}

	return nil
}

// SetIfRevision atomically updates a SQL setting when its revision matches
// expectedRevision.
func (r *SQLRepository) SetIfRevision(ctx context.Context, key, expectedRevision string, value interface{}) error {
	if r == nil || r.db == nil {
		return storageError("set revision", fmt.Errorf("database is not configured"))
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return storageError("encode", err)
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return storageError("begin transaction", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			return
		}
	}()
	inserted, err := tx.ExecContext(ctx, `INSERT INTO system_settings (key, value, updated_at) VALUES (?, ?, ?) ON CONFLICT(key) DO NOTHING`, key, string(payload), time.Now())
	if err != nil {
		return storageError("insert revision", err)
	}
	insertedRows, rowsErr := inserted.RowsAffected()
	if rowsErr != nil {
		return storageError("insert revision rows affected", rowsErr)
	}
	if insertedRows == 1 {
		if expectedRevision != "" {
			return fmt.Errorf("%w", ErrRevisionConflict)
		}
		if err := tx.Commit(); err != nil {
			return storageError("commit revision", err)
		}
		return nil
	}
	if insertedRows != 0 {
		return storageError("insert revision rows affected", fmt.Errorf("unexpected affected row count %d", insertedRows))
	}
	result, err := tx.ExecContext(ctx, `UPDATE system_settings SET value = ?, updated_at = ? WHERE key = ? AND json_extract(value, '$.settings_revision') = ?`, string(payload), time.Now(), key, expectedRevision)
	if err != nil {
		return storageError("update revision", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return storageError("update revision rows affected", err)
	}
	if rows != 1 {
		return fmt.Errorf("%w", ErrRevisionConflict)
	}
	if err := tx.Commit(); err != nil {
		return storageError("commit revision", err)
	}
	return nil
}

// List 列出所有設定
func (r *SQLRepository) List(ctx context.Context) ([]*SettingItem, error) {
	if r == nil || r.db == nil {
		return nil, storageError("list", fmt.Errorf("database is not configured"))
	}
	query := `
		SELECT key, value, description, updated_at
		FROM system_settings
		ORDER BY key ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, storageError("list", err)
	}
	defer rows.Close()

	items := make([]*SettingItem, 0)
	for rows.Next() {
		var item SettingItem
		var valueJSON string
		var description sql.NullString
		var updatedAt sql.NullString

		err := rows.Scan(&item.Key, &valueJSON, &description, &updatedAt)
		if err != nil {
			return nil, storageError("scan", err)
		}

		// 解析 JSON 值
		if err := json.Unmarshal([]byte(valueJSON), &item.Value); err != nil {
			return nil, storageError("decode", fmt.Errorf("decode setting %s: %w", item.Key, err))
		}

		if description.Valid {
			item.Description = description.String
		}

		if updatedAt.Valid {
			parsedUpdatedAt, err := common.ParseTimeString(updatedAt.String)
			if err != nil {
				return nil, fmt.Errorf("解析更新時間失敗: %w", err)
			}
			item.UpdatedAt = parsedUpdatedAt
		}

		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, storageError("iterate", err)
	}

	return items, nil
}

// Delete 刪除設定
func (r *SQLRepository) Delete(ctx context.Context, key string) error {
	if r == nil || r.db == nil {
		return storageError("delete", fmt.Errorf("database is not configured"))
	}
	query := `DELETE FROM system_settings WHERE key = ?`

	result, err := r.db.ExecContext(ctx, query, key)
	if err != nil {
		return storageError("delete", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return storageError("delete rows affected", err)
	}
	if rows == 0 {
		return fmt.Errorf("設定不存在: %s", key)
	}

	return nil
}

// InitDefaults 初始化預設設定
func (r *SQLRepository) InitDefaults(ctx context.Context) error {
	defaults := map[string]interface{}{
		KeyWritePrecision:    "millisecond",
		KeyPartitionInterval: "monthly",
		KeyBatchSize:         1000,
		KeyDefaultRetryCount: 3,
		KeyDefaultRetryDelay: 1000,
	}

	for key, value := range defaults {
		// 檢查是否已存在
		_, err := r.Get(ctx, key)
		if errors.Is(err, ErrSettingNotFound) {
			// 不存在，建立預設值
			if err := r.Set(ctx, key, value); err != nil {
				return fmt.Errorf("初始化預設設定失敗 (%s): %w", key, err)
			}
		} else if err != nil {
			return fmt.Errorf("初始化預設設定讀取失敗 (%s): %w", key, err)
		}
	}

	return nil
}
