// Package pollinggroup 提供輪詢群組管理功能的 SQL Repository 實現。
package pollinggroup

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// SQL Repository 實現
// =============================================================================

// SQLRepository SQL 輪詢群組儲存庫
type SQLRepository struct {
	db *sql.DB
}

// NewSQLRepository 建立新的 SQL 儲存庫
func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

// Create 建立輪詢群組
func (r *SQLRepository) Create(ctx context.Context, group *schema.PollingGroup) error {
	if group.ID == "" {
		id, err := common.NewUUID()
		if err != nil {
			return fmt.Errorf("建立輪詢群組 ID 失敗: %w", err)
		}
		group.ID = id
	}

	query := `
		INSERT INTO polling_groups (id, name, description, interval_ms, priority, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		group.ID,
		group.Name,
		group.Description,
		group.IntervalMs,
		group.Priority,
		group.Enabled,
		group.CreatedAt,
		group.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("建立輪詢群組失敗: %w", err)
	}

	return nil
}

// Update 更新輪詢群組
func (r *SQLRepository) Update(ctx context.Context, group *schema.PollingGroup) error {
	query := `
		UPDATE polling_groups 
		SET name = ?, description = ?, interval_ms = ?, priority = ?, enabled = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		group.Name,
		group.Description,
		group.IntervalMs,
		group.Priority,
		group.Enabled,
		time.Now(),
		group.ID,
	)

	if err != nil {
		return fmt.Errorf("更新輪詢群組失敗: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("輪詢群組不存在: %s", group.ID)
	}

	return nil
}

// Delete 刪除輪詢群組
func (r *SQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM polling_groups WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("刪除輪詢群組失敗: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("輪詢群組不存在: %s", id)
	}

	return nil
}

// GetByID 根據 ID 取得輪詢群組
func (r *SQLRepository) GetByID(ctx context.Context, id string) (*schema.PollingGroup, error) {
	query := `
		SELECT id, name, description, interval_ms, priority, enabled, created_at, updated_at
		FROM polling_groups WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanGroup(row)
}

// List 列出輪詢群組
func (r *SQLRepository) List(ctx context.Context) ([]*schema.PollingGroup, error) {
	query := `
		SELECT id, name, description, interval_ms, priority, enabled, created_at, updated_at
		FROM polling_groups
		ORDER BY priority DESC, name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("查詢輪詢群組失敗: %w", err)
	}
	defer rows.Close()

	groups := make([]*schema.PollingGroup, 0)
	for rows.Next() {
		group, err := r.scanGroupFromRows(rows)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

// Count 計算輪詢群組數量
func (r *SQLRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM polling_groups`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("計算輪詢群組數量失敗: %w", err)
	}

	return count, nil
}

// Clear 清空所有輪詢群組
func (r *SQLRepository) Clear(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM polling_groups`)
	return err
}

// scanGroup 從單一行掃描輪詢群組
func (r *SQLRepository) scanGroup(row *sql.Row) (*schema.PollingGroup, error) {
	var group schema.PollingGroup
	var description sql.NullString
	var createdAt, updatedAt string

	err := row.Scan(
		&group.ID,
		&group.Name,
		&description,
		&group.IntervalMs,
		&group.Priority,
		&group.Enabled,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("輪詢群組不存在")
	}
	if err != nil {
		return nil, fmt.Errorf("掃描輪詢群組失敗: %w", err)
	}

	if description.Valid {
		group.Description = description.String
	}

	parsedCreatedAt, err := common.ParseTimeString(createdAt)
	if err != nil {
		return nil, fmt.Errorf("解析建立時間失敗: %w", err)
	}
	parsedUpdatedAt, err := common.ParseTimeString(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("解析更新時間失敗: %w", err)
	}
	group.CreatedAt = parsedCreatedAt
	group.UpdatedAt = parsedUpdatedAt

	return &group, nil
}

// scanGroupFromRows 從多行結果掃描輪詢群組
func (r *SQLRepository) scanGroupFromRows(rows *sql.Rows) (*schema.PollingGroup, error) {
	var group schema.PollingGroup
	var description sql.NullString
	var createdAt, updatedAt string

	err := rows.Scan(
		&group.ID,
		&group.Name,
		&description,
		&group.IntervalMs,
		&group.Priority,
		&group.Enabled,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("掃描輪詢群組失敗: %w", err)
	}

	if description.Valid {
		group.Description = description.String
	}

	parsedCreatedAt, err := common.ParseTimeString(createdAt)
	if err != nil {
		return nil, fmt.Errorf("解析建立時間失敗: %w", err)
	}
	parsedUpdatedAt, err := common.ParseTimeString(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("解析更新時間失敗: %w", err)
	}
	group.CreatedAt = parsedCreatedAt
	group.UpdatedAt = parsedUpdatedAt

	return &group, nil
}
