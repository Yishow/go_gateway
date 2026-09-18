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

// sqlRunner is the query surface shared by *sql.DB and *sql.Tx.
type sqlRunner interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type SQLConnectorRepository struct {
	db sqlRunner
}

func NewSQLConnectorRepository(db *sql.DB) *SQLConnectorRepository {
	return &SQLConnectorRepository{db: db}
}

// WithTx returns the repository bound to a caller-owned local transaction.
func (r *SQLConnectorRepository) WithTx(tx *sql.Tx) ConnectorRepository {
	return &SQLConnectorRepository{db: tx}
}

func (r *SQLConnectorRepository) Create(ctx context.Context, connector *schema.DatabaseConnector) error {
	if strings.TrimSpace(connector.IdentityRevision) == "" {
		revision, err := common.NewUUID()
		if err != nil {
			return fmt.Errorf("建立資料庫連接器 identity revision 失敗: %w", err)
		}
		connector.IdentityRevision = revision
	}
	query := `INSERT INTO database_connectors (
		id, name, kind, connection_config, identity_revision, status, last_check_at, last_check_error, enabled, default_write_interval_seconds,
		last_schema_ensure_at, last_schema_ensure_status, last_schema_ensure_error,
		last_write_at, last_write_status, last_write_error,
		last_flush_at, last_flush_status, last_flush_error,
		created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(
		ctx,
		query,
		connector.ID,
		connector.Name,
		connector.Kind,
		connector.ConnectionConfig,
		connector.IdentityRevision,
		connector.Status,
		connector.LastCheckAt,
		connector.LastCheckError,
		connector.Enabled,
		connector.DefaultWriteIntervalSeconds,
		connector.LastSchemaEnsureAt,
		connector.LastSchemaEnsureStatus,
		connector.LastSchemaEnsureError,
		connector.LastWriteAt,
		connector.LastWriteStatus,
		connector.LastWriteError,
		connector.LastFlushAt,
		connector.LastFlushStatus,
		connector.LastFlushError,
		connector.CreatedAt,
		connector.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("建立資料庫連接器失敗: %w", err)
	}
	return nil
}

func (r *SQLConnectorRepository) Update(ctx context.Context, connector *schema.DatabaseConnector) error {
	return r.UpdateWithExpectedIdentityRevision(ctx, connector, connector.IdentityRevision)
}

func (r *SQLConnectorRepository) UpdateWithExpectedIdentityRevision(
	ctx context.Context,
	connector *schema.DatabaseConnector,
	expectedRevision string,
) error {
	expectedRevision = strings.TrimSpace(expectedRevision)
	if expectedRevision == "" || strings.TrimSpace(connector.IdentityRevision) == "" {
		return fmt.Errorf("%w: 資料庫連接器 identity revision 不可為空", ErrConnectorRevisionConflict)
	}
	query := `UPDATE database_connectors
		SET name = ?, kind = ?, connection_config = ?, identity_revision = ?, status = ?, last_check_at = ?, last_check_error = ?,
			enabled = ?, default_write_interval_seconds = ?,
			last_schema_ensure_at = ?, last_schema_ensure_status = ?, last_schema_ensure_error = ?,
			last_write_at = ?, last_write_status = ?, last_write_error = ?,
			last_flush_at = ?, last_flush_status = ?, last_flush_error = ?,
		updated_at = ?
		WHERE id = ? AND identity_revision = ?`

	result, err := r.db.ExecContext(
		ctx,
		query,
		connector.Name,
		connector.Kind,
		connector.ConnectionConfig,
		connector.IdentityRevision,
		connector.Status,
		connector.LastCheckAt,
		connector.LastCheckError,
		connector.Enabled,
		connector.DefaultWriteIntervalSeconds,
		connector.LastSchemaEnsureAt,
		connector.LastSchemaEnsureStatus,
		connector.LastSchemaEnsureError,
		connector.LastWriteAt,
		connector.LastWriteStatus,
		connector.LastWriteError,
		connector.LastFlushAt,
		connector.LastFlushStatus,
		connector.LastFlushError,
		connector.UpdatedAt,
		connector.ID,
		expectedRevision,
	)
	if err != nil {
		return fmt.Errorf("更新資料庫連接器失敗: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得資料庫連接器更新筆數失敗: %w", err)
	}
	if rows == 0 {
		var exists int
		err := r.db.QueryRowContext(ctx, `SELECT 1 FROM database_connectors WHERE id = ?`, connector.ID).Scan(&exists)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: 資料庫連接器不存在", ErrConnectorNotFound)
		}
		if err != nil {
			return fmt.Errorf("檢查資料庫連接器更新衝突失敗: %w", err)
		}
		return fmt.Errorf("%w: 資料庫連接器 identity revision 已變更", ErrConnectorRevisionConflict)
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
		return fmt.Errorf("%w: 資料庫連接器不存在", ErrConnectorNotFound)
	}
	return nil
}

func (r *SQLConnectorRepository) GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error) {
	query := `SELECT ` + connectorSelectColumns + `
	FROM database_connectors
	WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)
	return scanConnectorRow(row)
}

func (r *SQLConnectorRepository) List(ctx context.Context, filter ConnectorListFilter) ([]*schema.DatabaseConnector, error) {
	query := `SELECT ` + connectorSelectColumns + `
	FROM database_connectors
	WHERE 1=1`
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
	db sqlRunner
}

func NewSQLTargetMappingRepository(db *sql.DB) *SQLTargetMappingRepository {
	return &SQLTargetMappingRepository{db: db}
}

// WithTx returns the repository bound to a caller-owned local transaction.
func (r *SQLTargetMappingRepository) WithTx(tx *sql.Tx) TargetMappingRepository {
	return &SQLTargetMappingRepository{db: tx}
}

func (r *SQLTargetMappingRepository) Create(ctx context.Context, mapping *schema.DatabaseTargetMapping) error {
	query := `INSERT INTO database_target_mappings (
		id, tag_id, connector_id, table_schema, table_name, column_name, write_mode, timestamp_column, group_key, write_interval_seconds, enabled, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

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
		mapping.GroupKey,
		mapping.WriteIntervalSeconds,
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
	query := `UPDATE database_target_mappings
		SET connector_id = ?, table_schema = ?, table_name = ?, column_name = ?, write_mode = ?, timestamp_column = ?, group_key = ?, write_interval_seconds = ?, enabled = ?, updated_at = ?
		WHERE id = ?`

	result, err := r.db.ExecContext(
		ctx,
		query,
		mapping.ConnectorID,
		mapping.TableSchema,
		mapping.TableName,
		mapping.ColumnName,
		mapping.WriteMode,
		mapping.TimestampColumn,
		mapping.GroupKey,
		mapping.WriteIntervalSeconds,
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
	query := `SELECT
		id, tag_id, connector_id, table_schema, table_name, column_name, write_mode, timestamp_column,
		group_key, write_interval_seconds, enabled, CAST(created_at AS TEXT), CAST(updated_at AS TEXT)
	FROM database_target_mappings
	WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)
	return scanTargetMappingRow(row)
}

func (r *SQLTargetMappingRepository) List(ctx context.Context, filter TargetMappingListFilter) ([]*schema.DatabaseTargetMapping, error) {
	query := `SELECT
		id, tag_id, connector_id, table_schema, table_name, column_name, write_mode, timestamp_column,
		group_key, write_interval_seconds, enabled, CAST(created_at AS TEXT), CAST(updated_at AS TEXT)
	FROM database_target_mappings
	WHERE 1=1`
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

func scanTargetMappingRow(row *sql.Row) (*schema.DatabaseTargetMapping, error) {
	var mapping schema.DatabaseTargetMapping
	var timestampColumn, groupKey sql.NullString
	var writeIntervalSeconds sql.NullInt64
	var createdAt, updatedAt string

	err := row.Scan(
		&mapping.ID,
		&mapping.TagID,
		&mapping.ConnectorID,
		&mapping.TableSchema,
		&mapping.TableName,
		&mapping.ColumnName,
		&mapping.WriteMode,
		&timestampColumn,
		&groupKey,
		&writeIntervalSeconds,
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
	if groupKey.Valid {
		value := groupKey.String
		mapping.GroupKey = &value
	}
	if writeIntervalSeconds.Valid {
		value := int(writeIntervalSeconds.Int64)
		mapping.WriteIntervalSeconds = &value
	}

	return &mapping, nil
}

func scanTargetMappingRows(rows *sql.Rows) ([]*schema.DatabaseTargetMapping, error) {
	var mappings []*schema.DatabaseTargetMapping

	for rows.Next() {
		var mapping schema.DatabaseTargetMapping
		var timestampColumn, groupKey sql.NullString
		var writeIntervalSeconds sql.NullInt64
		var createdAt, updatedAt string

		if err := rows.Scan(
			&mapping.ID,
			&mapping.TagID,
			&mapping.ConnectorID,
			&mapping.TableSchema,
			&mapping.TableName,
			&mapping.ColumnName,
			&mapping.WriteMode,
			&timestampColumn,
			&groupKey,
			&writeIntervalSeconds,
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
		if groupKey.Valid {
			value := groupKey.String
			mapping.GroupKey = &value
		}
		if writeIntervalSeconds.Valid {
			value := int(writeIntervalSeconds.Int64)
			mapping.WriteIntervalSeconds = &value
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
