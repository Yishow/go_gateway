// Package settings 提供系統設定管理功能的 SQL Repository 實現。
package settings

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
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
	query := `
		SELECT key, value, description, updated_at
		FROM settings WHERE key = ?
	`

	row := r.db.QueryRowContext(ctx, query, key)

	var item SettingItem
	var valueJSON string
	var description sql.NullString

	err := row.Scan(&item.Key, &valueJSON, &description, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("設定不存在: %s", key)
	}
	if err != nil {
		return nil, fmt.Errorf("讀取設定失敗: %w", err)
	}

	// 解析 JSON 值
	if err := json.Unmarshal([]byte(valueJSON), &item.Value); err != nil {
		item.Value = valueJSON // 如果不是 JSON，使用原始字串
	}

	if description.Valid {
		item.Description = description.String
	}

	return &item, nil
}

// Set 設定值
func (r *SQLRepository) Set(ctx context.Context, key string, value interface{}) error {
	// 序列化值為 JSON
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化設定值失敗: %w", err)
	}

	// 使用 UPSERT 語法
	query := `
		INSERT INTO settings (key, value, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at
	`

	_, err = r.db.ExecContext(ctx, query, key, string(valueJSON), time.Now())
	if err != nil {
		return fmt.Errorf("寫入設定失敗: %w", err)
	}

	return nil
}

// List 列出所有設定
func (r *SQLRepository) List(ctx context.Context) ([]*SettingItem, error) {
	query := `
		SELECT key, value, description, updated_at
		FROM settings
		ORDER BY key ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("查詢設定失敗: %w", err)
	}
	defer rows.Close()

	items := make([]*SettingItem, 0)
	for rows.Next() {
		var item SettingItem
		var valueJSON string
		var description sql.NullString

		err := rows.Scan(&item.Key, &valueJSON, &description, &item.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("掃描設定失敗: %w", err)
		}

		// 解析 JSON 值
		if err := json.Unmarshal([]byte(valueJSON), &item.Value); err != nil {
			item.Value = valueJSON
		}

		if description.Valid {
			item.Description = description.String
		}

		items = append(items, &item)
	}

	return items, nil
}

// Delete 刪除設定
func (r *SQLRepository) Delete(ctx context.Context, key string) error {
	query := `DELETE FROM settings WHERE key = ?`

	result, err := r.db.ExecContext(ctx, query, key)
	if err != nil {
		return fmt.Errorf("刪除設定失敗: %w", err)
	}

	rows, _ := result.RowsAffected()
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
		if err != nil {
			// 不存在，建立預設值
			if err := r.Set(ctx, key, value); err != nil {
				return fmt.Errorf("初始化預設設定失敗 (%s): %w", key, err)
			}
		}
	}

	return nil
}
