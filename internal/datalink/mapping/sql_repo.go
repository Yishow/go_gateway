// Package mapping 提供映射管理功能的 SQL Repository 實現。
package mapping

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

// SQLRepository SQL 映射儲存庫
type SQLRepository struct {
	db *sql.DB
}

// NewSQLRepository 建立新的 SQL 儲存庫
func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

// Create 建立映射
func (r *SQLRepository) Create(ctx context.Context, mapping *schema.Mapping) error {
	if mapping.ID == "" {
		id, err := common.NewUUID()
		if err != nil {
			return fmt.Errorf("建立映射 ID 失敗: %w", err)
		}
		mapping.ID = id
	}

	query := `
		INSERT INTO mappings (id, point_id, tag_id, transform_pipeline, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		mapping.ID,
		mapping.PointID,
		mapping.TagID,
		mapping.TransformPipeline,
		mapping.Enabled,
		mapping.CreatedAt,
		mapping.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("建立映射失敗: %w", err)
	}

	return nil
}

// Update 更新映射
func (r *SQLRepository) Update(ctx context.Context, mapping *schema.Mapping) error {
	query := `
		UPDATE mappings 
		SET transform_pipeline = ?, enabled = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		mapping.TransformPipeline,
		mapping.Enabled,
		time.Now(),
		mapping.ID,
	)

	if err != nil {
		return fmt.Errorf("更新映射失敗: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("映射不存在: %s", mapping.ID)
	}

	return nil
}

// Delete 刪除映射
func (r *SQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM mappings WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("刪除映射失敗: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("映射不存在: %s", id)
	}

	return nil
}

// GetByID 根據 ID 取得映射
func (r *SQLRepository) GetByID(ctx context.Context, id string) (*schema.Mapping, error) {
	query := `
		SELECT id, point_id, tag_id, transform_pipeline, enabled, created_at, updated_at
		FROM mappings WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanMapping(row)
}

// GetByPointID 根據點位 ID 取得映射列表
func (r *SQLRepository) GetByPointID(ctx context.Context, pointID string) ([]*schema.Mapping, error) {
	query := `
		SELECT id, point_id, tag_id, transform_pipeline, enabled, created_at, updated_at
		FROM mappings WHERE point_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, pointID)
	if err != nil {
		return nil, fmt.Errorf("查詢映射失敗: %w", err)
	}
	defer rows.Close()

	return r.scanMappings(rows)
}

// GetByTagID 根據標籤 ID 取得映射列表
func (r *SQLRepository) GetByTagID(ctx context.Context, tagID string) ([]*schema.Mapping, error) {
	query := `
		SELECT id, point_id, tag_id, transform_pipeline, enabled, created_at, updated_at
		FROM mappings WHERE tag_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, tagID)
	if err != nil {
		return nil, fmt.Errorf("查詢映射失敗: %w", err)
	}
	defer rows.Close()

	return r.scanMappings(rows)
}

// List 列出映射
func (r *SQLRepository) List(ctx context.Context, filter ListFilter) ([]*schema.Mapping, error) {
	query := `
		SELECT id, point_id, tag_id, transform_pipeline, enabled, created_at, updated_at
		FROM mappings
		WHERE 1=1
	`
	args := []interface{}{}

	if filter.PointID != nil {
		query += ` AND point_id = ?`
		args = append(args, *filter.PointID)
	}

	if filter.TagID != nil {
		query += ` AND tag_id = ?`
		args = append(args, *filter.TagID)
	}

	if filter.Enabled != nil {
		query += ` AND enabled = ?`
		args = append(args, *filter.Enabled)
	}

	query += ` ORDER BY created_at DESC`

	if filter.Limit > 0 {
		query += fmt.Sprintf(` LIMIT %d`, filter.Limit)
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(` OFFSET %d`, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("查詢映射失敗: %w", err)
	}
	defer rows.Close()

	return r.scanMappings(rows)
}

// scanMapping 從單一行掃描映射
func (r *SQLRepository) scanMapping(row *sql.Row) (*schema.Mapping, error) {
	var mapping schema.Mapping
	var createdAt, updatedAt string

	err := row.Scan(
		&mapping.ID,
		&mapping.PointID,
		&mapping.TagID,
		&mapping.TransformPipeline,
		&mapping.Enabled,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("映射不存在")
	}
	if err != nil {
		return nil, fmt.Errorf("掃描映射失敗: %w", err)
	}

	parsedCreatedAt, err := common.ParseTimeString(createdAt)
	if err != nil {
		return nil, fmt.Errorf("解析建立時間失敗: %w", err)
	}
	parsedUpdatedAt, err := common.ParseTimeString(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("解析更新時間失敗: %w", err)
	}
	mapping.CreatedAt = parsedCreatedAt
	mapping.UpdatedAt = parsedUpdatedAt

	return &mapping, nil
}

// scanMappings 從多行結果掃描映射列表
func (r *SQLRepository) scanMappings(rows *sql.Rows) ([]*schema.Mapping, error) {
	mappings := make([]*schema.Mapping, 0)

	for rows.Next() {
		var mapping schema.Mapping
		var createdAt, updatedAt string

		err := rows.Scan(
			&mapping.ID,
			&mapping.PointID,
			&mapping.TagID,
			&mapping.TransformPipeline,
			&mapping.Enabled,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("掃描映射失敗: %w", err)
		}

		parsedCreatedAt, err := common.ParseTimeString(createdAt)
		if err != nil {
			return nil, fmt.Errorf("解析建立時間失敗: %w", err)
		}
		parsedUpdatedAt, err := common.ParseTimeString(updatedAt)
		if err != nil {
			return nil, fmt.Errorf("解析更新時間失敗: %w", err)
		}
		mapping.CreatedAt = parsedCreatedAt
		mapping.UpdatedAt = parsedUpdatedAt

		mappings = append(mappings, &mapping)
	}

	return mappings, nil
}
