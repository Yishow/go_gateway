// Package tag 提供標籤管理功能的 SQL Repository 實現。
package tag

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// SQL Repository 實現
// =============================================================================

// SQLRepository SQL 標籤儲存庫
type SQLRepository struct {
	db *sql.DB
}

// NewSQLRepository 建立新的 SQL 儲存庫
func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

// Create 建立標籤
func (r *SQLRepository) Create(ctx context.Context, tag *schema.Tag) error {
	query := `
		INSERT INTO tags (id, key, key_lower, display_name, data_type, unit, description, status, labels, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		tag.ID,
		tag.Key,
		strings.ToLower(tag.Key),
		tag.DisplayName,
		tag.DataType,
		tag.Unit,
		tag.Description,
		tag.Status,
		tag.Labels,
		tag.CreatedAt,
		tag.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("建立標籤失敗: %w", err)
	}

	return nil
}

// BatchCreate 在交易內批量建立標籤
func (r *SQLRepository) BatchCreate(ctx context.Context, tags []*schema.Tag) error {
	if len(tags) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("開啟交易失敗: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			return
		}
	}()

	query := `
		INSERT INTO tags (id, key, key_lower, display_name, data_type, unit, description, status, labels, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("準備插入語句失敗: %w", err)
	}
	defer stmt.Close()

	for _, tag := range tags {
		_, err = stmt.ExecContext(ctx,
			tag.ID,
			tag.Key,
			strings.ToLower(tag.Key),
			tag.DisplayName,
			tag.DataType,
			tag.Unit,
			tag.Description,
			tag.Status,
			tag.Labels,
			tag.CreatedAt,
			tag.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("批量建立標籤 %s 失敗: %w", tag.Key, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交交易失敗: %w", err)
	}

	return nil
}

// Update 更新標籤
func (r *SQLRepository) Update(ctx context.Context, tag *schema.Tag) error {
	query := `
		UPDATE tags 
		SET display_name = ?, data_type = ?, unit = ?, description = ?, 
		    status = ?, labels = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		tag.DisplayName,
		tag.DataType,
		tag.Unit,
		tag.Description,
		tag.Status,
		tag.Labels,
		time.Now(),
		tag.ID,
	)

	if err != nil {
		return fmt.Errorf("更新標籤失敗: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得更新影響列數失敗: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("%w: %s", ErrTagNotFound, tag.ID)
	}

	return nil
}

// Delete 刪除標籤
func (r *SQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tags WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("刪除標籤失敗: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得刪除影響列數失敗: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("%w: %s", ErrTagNotFound, id)
	}

	return nil
}

// GetByID 根據 ID 取得標籤
func (r *SQLRepository) GetByID(ctx context.Context, id string) (*schema.Tag, error) {
	query := `
		SELECT id, key, display_name, data_type, unit, description, status, labels, created_at, updated_at
		FROM tags WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanTag(row)
}

// GetByKey 根據 Key 取得標籤
func (r *SQLRepository) GetByKey(ctx context.Context, key string) (*schema.Tag, error) {
	query := `
		SELECT id, key, display_name, data_type, unit, description, status, labels, created_at, updated_at
		FROM tags WHERE key_lower = ?
	`

	row := r.db.QueryRowContext(ctx, query, strings.ToLower(key))
	return r.scanTag(row)
}

// List 列出標籤
func (r *SQLRepository) List(ctx context.Context, filter ListFilter) ([]*schema.Tag, error) {
	query := `
		SELECT id, key, display_name, data_type, unit, description, status, labels, created_at, updated_at
		FROM tags
		WHERE 1=1
	`
	args := []interface{}{}

	if filter.Status != nil {
		query += ` AND status = ?`
		args = append(args, *filter.Status)
	}

	if filter.DataType != nil {
		query += ` AND data_type = ?`
		args = append(args, *filter.DataType)
	}

	query += ` ORDER BY key ASC`

	if filter.Limit > 0 {
		query += ` LIMIT ?`
		args = append(args, filter.Limit)
	}
	if filter.Offset > 0 {
		query += ` OFFSET ?`
		args = append(args, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("查詢標籤失敗: %w", err)
	}
	defer rows.Close()

	tags := make([]*schema.Tag, 0)
	for rows.Next() {
		tag, err := r.scanTagFromRows(rows)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// Count 計算標籤數量
func (r *SQLRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM tags`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("計算標籤數量失敗: %w", err)
	}

	return count, nil
}

// scanTag 從單一行掃描標籤
func (r *SQLRepository) scanTag(row *sql.Row) (*schema.Tag, error) {
	var tag schema.Tag
	var displayName, unit, description, labels sql.NullString
	var createdAt, updatedAt string

	err := row.Scan(
		&tag.ID,
		&tag.Key,
		&displayName,
		&tag.DataType,
		&unit,
		&description,
		&tag.Status,
		&labels,
		&createdAt,
		&updatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w", ErrTagNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("掃描標籤失敗: %w", err)
	}

	if displayName.Valid {
		tag.DisplayName = displayName.String
	}
	if unit.Valid {
		tag.Unit = unit.String
	}
	if description.Valid {
		tag.Description = description.String
	}
	if labels.Valid {
		tag.Labels = labels.String
	}

	parsedCreatedAt, err := common.ParseTimeString(createdAt)
	if err != nil {
		return nil, fmt.Errorf("解析建立時間失敗: %w", err)
	}
	parsedUpdatedAt, err := common.ParseTimeString(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("解析更新時間失敗: %w", err)
	}
	tag.CreatedAt = parsedCreatedAt
	tag.UpdatedAt = parsedUpdatedAt

	return &tag, nil
}

// scanTagFromRows 從多行結果掃描標籤
func (r *SQLRepository) scanTagFromRows(rows *sql.Rows) (*schema.Tag, error) {
	var tag schema.Tag
	var displayName, unit, description, labels sql.NullString
	var createdAt, updatedAt string

	err := rows.Scan(
		&tag.ID,
		&tag.Key,
		&displayName,
		&tag.DataType,
		&unit,
		&description,
		&tag.Status,
		&labels,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("掃描標籤失敗: %w", err)
	}

	if displayName.Valid {
		tag.DisplayName = displayName.String
	}
	if unit.Valid {
		tag.Unit = unit.String
	}
	if description.Valid {
		tag.Description = description.String
	}
	if labels.Valid {
		tag.Labels = labels.String
	}

	parsedCreatedAt, err := common.ParseTimeString(createdAt)
	if err != nil {
		return nil, fmt.Errorf("解析建立時間失敗: %w", err)
	}
	parsedUpdatedAt, err := common.ParseTimeString(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("解析更新時間失敗: %w", err)
	}
	tag.CreatedAt = parsedCreatedAt
	tag.UpdatedAt = parsedUpdatedAt

	return &tag, nil
}

// ExistsByKey 檢查 Key 是否存在
func (r *SQLRepository) ExistsByKey(ctx context.Context, key string) (bool, error) {
	query := `SELECT COUNT(*) FROM tags WHERE key_lower = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, strings.ToLower(key)).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("檢查標籤是否存在失敗: %w", err)
	}

	return count > 0, nil
}

// UpdateStatus 更新標籤狀態
func (r *SQLRepository) UpdateStatus(ctx context.Context, id string, status schema.TagStatus) error {
	query := `
		UPDATE tags 
		SET status = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("更新標籤狀態失敗: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得狀態更新影響列數失敗: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("%w: %s", ErrTagNotFound, id)
	}

	return nil
}
