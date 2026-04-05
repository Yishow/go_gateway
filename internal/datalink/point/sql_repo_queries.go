package point

import (
	"context"
	"fmt"

	"go-gateway/internal/datalink/schema"
)

// GetByID 根據 ID 取得點位
func (r *SQLRepository) GetByID(ctx context.Context, id string) (*schema.Point, error) {
	query := `
		SELECT id, device_id, name, description, address, function, data_type, data_format, mode, 
		       polling_group_id, last_read_at, last_value, last_error, enabled, created_at, updated_at
		FROM points WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanPoint(row)
}

// ListByDevice 根據設備 ID 取得點位列表
func (r *SQLRepository) ListByDevice(ctx context.Context, deviceID string) ([]*schema.Point, error) {
	query := `
		SELECT id, device_id, name, description, address, function, data_type, data_format, mode,
		       polling_group_id, last_read_at, last_value, last_error, enabled, created_at, updated_at
		FROM points WHERE device_id = ?
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, deviceID)
	if err != nil {
		return nil, fmt.Errorf("查詢點位失敗: %w", err)
	}
	defer rows.Close()

	points := make([]*schema.Point, 0)
	for rows.Next() {
		point, err := r.scanPointFromRows(rows)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	return points, nil
}

// ListByPollingGroup 根據輪詢群組 ID 取得點位列表
func (r *SQLRepository) ListByPollingGroup(ctx context.Context, groupID string) ([]*schema.Point, error) {
	query := `
		SELECT id, device_id, name, description, address, function, data_type, data_format, mode,
		       polling_group_id, last_read_at, last_value, last_error, enabled, created_at, updated_at
		FROM points WHERE polling_group_id = ?
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("查詢點位失敗: %w", err)
	}
	defer rows.Close()

	points := make([]*schema.Point, 0)
	for rows.Next() {
		point, err := r.scanPointFromRows(rows)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	return points, nil
}

// List 列出點位
func (r *SQLRepository) List(ctx context.Context, filter ListFilter) ([]*schema.Point, error) {
	query := `
		SELECT id, device_id, name, description, address, function, data_type, data_format, mode,
		       polling_group_id, last_read_at, last_value, last_error, enabled, created_at, updated_at
		FROM points
		WHERE 1=1
	`
	args := []interface{}{}

	if filter.DeviceID != nil {
		query += ` AND device_id = ?`
		args = append(args, *filter.DeviceID)
	}

	if filter.PollingGroupID != nil {
		query += ` AND polling_group_id = ?`
		args = append(args, *filter.PollingGroupID)
	}

	query += ` ORDER BY name ASC`

	if filter.Limit > 0 {
		query += fmt.Sprintf(` LIMIT %d`, filter.Limit)
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(` OFFSET %d`, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("查詢點位失敗: %w", err)
	}
	defer rows.Close()

	points := make([]*schema.Point, 0)
	for rows.Next() {
		point, err := r.scanPointFromRows(rows)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	return points, nil
}
