package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

type SQLConnectorRepository struct {
	db *sql.DB
}

func NewSQLConnectorRepository(db *sql.DB) *SQLConnectorRepository {
	return &SQLConnectorRepository{db: db}
}

func (r *SQLConnectorRepository) Create(ctx context.Context, connector *schema.DatabaseConnector) error {
	query := `
		INSERT INTO database_connectors (
			id, name, kind, connection_config, status, last_check_at, last_check_error, enabled, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		connector.ID,
		connector.Name,
		connector.Kind,
		connector.ConnectionConfig,
		connector.Status,
		connector.LastCheckAt,
		connector.LastCheckError,
		connector.Enabled,
		connector.CreatedAt,
		connector.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("建立資料庫連接器失敗: %w", err)
	}
	return nil
}

func (r *SQLConnectorRepository) Update(ctx context.Context, connector *schema.DatabaseConnector) error {
	query := `
		UPDATE database_connectors
		SET name = ?, kind = ?, connection_config = ?, status = ?, last_check_at = ?, last_check_error = ?, enabled = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		connector.Name,
		connector.Kind,
		connector.ConnectionConfig,
		connector.Status,
		connector.LastCheckAt,
		connector.LastCheckError,
		connector.Enabled,
		connector.UpdatedAt,
		connector.ID,
	)
	if err != nil {
		return fmt.Errorf("更新資料庫連接器失敗: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得資料庫連接器更新筆數失敗: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("資料庫連接器不存在")
	}
	return nil
}

func (r *SQLConnectorRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM database_connectors WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("刪除資料庫連接器失敗: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得資料庫連接器刪除筆數失敗: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("資料庫連接器不存在")
	}
	return nil
}

func (r *SQLConnectorRepository) GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error) {
	query := `
		SELECT
			id,
			name,
			kind,
			connection_config,
			status,
			COALESCE(CAST(last_check_at AS TEXT), ''),
			last_check_error,
			enabled,
			CAST(created_at AS TEXT),
			CAST(updated_at AS TEXT)
		FROM database_connectors
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)
	return scanConnectorRow(row)
}

func (r *SQLConnectorRepository) List(ctx context.Context, filter ConnectorListFilter) ([]*schema.DatabaseConnector, error) {
	query := `
		SELECT
			id,
			name,
			kind,
			connection_config,
			status,
			COALESCE(CAST(last_check_at AS TEXT), ''),
			last_check_error,
			enabled,
			CAST(created_at AS TEXT),
			CAST(updated_at AS TEXT)
		FROM database_connectors
		WHERE 1=1
	`
	args := []any{}

	if filter.Enabled != nil {
		query += ` AND enabled = ?`
		args = append(args, *filter.Enabled)
	}

	query += ` ORDER BY created_at DESC, id DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("查詢資料庫連接器失敗: %w", err)
	}
	defer rows.Close()

	return scanConnectorRows(rows)
}

type SQLTargetMappingRepository struct {
	db *sql.DB
}

func NewSQLTargetMappingRepository(db *sql.DB) *SQLTargetMappingRepository {
	return &SQLTargetMappingRepository{db: db}
}

func (r *SQLTargetMappingRepository) Create(ctx context.Context, mapping *schema.DatabaseTargetMapping) error {
	query := `
		INSERT INTO database_target_mappings (
			id, tag_id, connector_id, table_schema, table_name, column_name, write_mode, timestamp_column, enabled, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		mapping.ID,
		mapping.TagID,
		mapping.ConnectorID,
		mapping.TableSchema,
		mapping.TableName,
		mapping.ColumnName,
		mapping.WriteMode,
		mapping.TimestampColumn,
		mapping.Enabled,
		mapping.CreatedAt,
		mapping.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("建立資料庫目標映射失敗: %w", err)
	}
	return nil
}

func (r *SQLTargetMappingRepository) Update(ctx context.Context, mapping *schema.DatabaseTargetMapping) error {
	query := `
		UPDATE database_target_mappings
		SET connector_id = ?, table_schema = ?, table_name = ?, column_name = ?, write_mode = ?, timestamp_column = ?, enabled = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		mapping.ConnectorID,
		mapping.TableSchema,
		mapping.TableName,
		mapping.ColumnName,
		mapping.WriteMode,
		mapping.TimestampColumn,
		mapping.Enabled,
		mapping.UpdatedAt,
		mapping.ID,
	)
	if err != nil {
		return fmt.Errorf("更新資料庫目標映射失敗: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得資料庫目標映射更新筆數失敗: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("資料庫目標映射不存在")
	}
	return nil
}

func (r *SQLTargetMappingRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM database_target_mappings WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("刪除資料庫目標映射失敗: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得資料庫目標映射刪除筆數失敗: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("資料庫目標映射不存在")
	}
	return nil
}

func (r *SQLTargetMappingRepository) DeleteByConnectorID(ctx context.Context, connectorID string) error {
	if _, err := r.db.ExecContext(
		ctx,
		`DELETE FROM database_target_mappings WHERE connector_id = ?`,
		connectorID,
	); err != nil {
		return fmt.Errorf("刪除資料庫連接器對應映射失敗: %w", err)
	}
	return nil
}

func (r *SQLTargetMappingRepository) GetByID(ctx context.Context, id string) (*schema.DatabaseTargetMapping, error) {
	query := `
		SELECT
			id,
			tag_id,
			connector_id,
			table_schema,
			table_name,
			column_name,
			write_mode,
			timestamp_column,
			enabled,
			CAST(created_at AS TEXT),
			CAST(updated_at AS TEXT)
		FROM database_target_mappings
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)
	return scanTargetMappingRow(row)
}

func (r *SQLTargetMappingRepository) List(ctx context.Context, filter TargetMappingListFilter) ([]*schema.DatabaseTargetMapping, error) {
	query := `
		SELECT
			id,
			tag_id,
			connector_id,
			table_schema,
			table_name,
			column_name,
			write_mode,
			timestamp_column,
			enabled,
			CAST(created_at AS TEXT),
			CAST(updated_at AS TEXT)
		FROM database_target_mappings
		WHERE 1=1
	`
	args := []any{}

	if filter.ConnectorID != nil {
		query += ` AND connector_id = ?`
		args = append(args, *filter.ConnectorID)
	}
	if filter.TagID != nil {
		query += ` AND tag_id = ?`
		args = append(args, *filter.TagID)
	}
	if filter.Enabled != nil {
		query += ` AND enabled = ?`
		args = append(args, *filter.Enabled)
	}

	query += ` ORDER BY created_at DESC, id DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("查詢資料庫目標映射失敗: %w", err)
	}
	defer rows.Close()

	return scanTargetMappingRows(rows)
}

func scanConnectorRow(row *sql.Row) (*schema.DatabaseConnector, error) {
	var connector schema.DatabaseConnector
	var lastCheckAt string
	var createdAt string
	var updatedAt string

	err := row.Scan(
		&connector.ID,
		&connector.Name,
		&connector.Kind,
		&connector.ConnectionConfig,
		&connector.Status,
		&lastCheckAt,
		&connector.LastCheckError,
		&connector.Enabled,
		&createdAt,
		&updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("資料庫連接器不存在")
	}
	if err != nil {
		return nil, fmt.Errorf("掃描資料庫連接器失敗: %w", err)
	}

	parsedCreatedAt, err := common.ParseTimeString(createdAt)
	if err != nil {
		return nil, fmt.Errorf("解析資料庫連接器建立時間失敗: %w", err)
	}
	parsedUpdatedAt, err := common.ParseTimeString(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("解析資料庫連接器更新時間失敗: %w", err)
	}
	connector.CreatedAt = parsedCreatedAt
	connector.UpdatedAt = parsedUpdatedAt

	if strings.TrimSpace(lastCheckAt) != "" {
		parsedLastCheckAt, err := common.ParseTimeString(lastCheckAt)
		if err != nil {
			return nil, fmt.Errorf("解析資料庫連接器檢查時間失敗: %w", err)
		}
		connector.LastCheckAt = &parsedLastCheckAt
	}

	return &connector, nil
}

func scanConnectorRows(rows *sql.Rows) ([]*schema.DatabaseConnector, error) {
	var connectors []*schema.DatabaseConnector

	for rows.Next() {
		var connector schema.DatabaseConnector
		var lastCheckAt string
		var createdAt string
		var updatedAt string

		if err := rows.Scan(
			&connector.ID,
			&connector.Name,
			&connector.Kind,
			&connector.ConnectionConfig,
			&connector.Status,
			&lastCheckAt,
			&connector.LastCheckError,
			&connector.Enabled,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("掃描資料庫連接器失敗: %w", err)
		}

		parsedCreatedAt, err := common.ParseTimeString(createdAt)
		if err != nil {
			return nil, fmt.Errorf("解析資料庫連接器建立時間失敗: %w", err)
		}
		parsedUpdatedAt, err := common.ParseTimeString(updatedAt)
		if err != nil {
			return nil, fmt.Errorf("解析資料庫連接器更新時間失敗: %w", err)
		}
		connector.CreatedAt = parsedCreatedAt
		connector.UpdatedAt = parsedUpdatedAt

		if strings.TrimSpace(lastCheckAt) != "" {
			parsedLastCheckAt, err := common.ParseTimeString(lastCheckAt)
			if err != nil {
				return nil, fmt.Errorf("解析資料庫連接器檢查時間失敗: %w", err)
			}
			connector.LastCheckAt = &parsedLastCheckAt
		}

		connectors = append(connectors, &connector)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍歷資料庫連接器失敗: %w", err)
	}

	return connectors, nil
}

func scanTargetMappingRow(row *sql.Row) (*schema.DatabaseTargetMapping, error) {
	var mapping schema.DatabaseTargetMapping
	var timestampColumn sql.NullString
	var createdAt string
	var updatedAt string

	err := row.Scan(
		&mapping.ID,
		&mapping.TagID,
		&mapping.ConnectorID,
		&mapping.TableSchema,
		&mapping.TableName,
		&mapping.ColumnName,
		&mapping.WriteMode,
		&timestampColumn,
		&mapping.Enabled,
		&createdAt,
		&updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("資料庫目標映射不存在")
	}
	if err != nil {
		return nil, fmt.Errorf("掃描資料庫目標映射失敗: %w", err)
	}

	parsedCreatedAt, err := common.ParseTimeString(createdAt)
	if err != nil {
		return nil, fmt.Errorf("解析資料庫目標映射建立時間失敗: %w", err)
	}
	parsedUpdatedAt, err := common.ParseTimeString(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("解析資料庫目標映射更新時間失敗: %w", err)
	}
	mapping.CreatedAt = parsedCreatedAt
	mapping.UpdatedAt = parsedUpdatedAt

	if timestampColumn.Valid {
		value := timestampColumn.String
		mapping.TimestampColumn = &value
	}

	return &mapping, nil
}

func scanTargetMappingRows(rows *sql.Rows) ([]*schema.DatabaseTargetMapping, error) {
	var mappings []*schema.DatabaseTargetMapping

	for rows.Next() {
		var mapping schema.DatabaseTargetMapping
		var timestampColumn sql.NullString
		var createdAt string
		var updatedAt string

		if err := rows.Scan(
			&mapping.ID,
			&mapping.TagID,
			&mapping.ConnectorID,
			&mapping.TableSchema,
			&mapping.TableName,
			&mapping.ColumnName,
			&mapping.WriteMode,
			&timestampColumn,
			&mapping.Enabled,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("掃描資料庫目標映射失敗: %w", err)
		}

		parsedCreatedAt, err := common.ParseTimeString(createdAt)
		if err != nil {
			return nil, fmt.Errorf("解析資料庫目標映射建立時間失敗: %w", err)
		}
		parsedUpdatedAt, err := common.ParseTimeString(updatedAt)
		if err != nil {
			return nil, fmt.Errorf("解析資料庫目標映射更新時間失敗: %w", err)
		}
		mapping.CreatedAt = parsedCreatedAt
		mapping.UpdatedAt = parsedUpdatedAt

		if timestampColumn.Valid {
			value := timestampColumn.String
			mapping.TimestampColumn = &value
		}

		mappings = append(mappings, &mapping)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍歷資料庫目標映射失敗: %w", err)
	}

	return mappings, nil
}

var _ ConnectorRepository = (*SQLConnectorRepository)(nil)
var _ TargetMappingRepository = (*SQLTargetMappingRepository)(nil)
