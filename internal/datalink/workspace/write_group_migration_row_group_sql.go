package workspace

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func readRowGroupTargetsForPoint(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	connectorID string,
	pointID string,
) ([]writeGroupRowGroupTarget, error) {
	rows, err := tx.QueryContext(ctx, repo.query(`
		SELECT dt.id, dt.tag_id, m.point_id, dt.connector_id,
			dt.table_schema, dt.table_name, dt.column_name, dt.write_mode,
			dt.timestamp_column, dt.group_key, dt.write_interval_seconds, dt.enabled
		FROM database_target_mappings dt
		JOIN mappings m ON m.tag_id = dt.tag_id
		WHERE dt.connector_id = $1 AND m.point_id = $2
		ORDER BY dt.id ASC
	`), connectorID, pointID)
	if err != nil {
		return nil, fmt.Errorf("read row-group target mappings: %w", err)
	}
	defer rows.Close()
	targets := make([]writeGroupRowGroupTarget, 0)
	for rows.Next() {
		var target writeGroupRowGroupTarget
		var timestampColumn, groupKey sql.NullString
		var writeInterval sql.NullInt64
		if err := rows.Scan(
			&target.id,
			&target.tagID,
			&target.pointID,
			&target.connectorID,
			&target.tableSchema,
			&target.tableName,
			&target.columnName,
			&target.writeMode,
			&timestampColumn,
			&groupKey,
			&writeInterval,
			&target.enabled,
		); err != nil {
			return nil, fmt.Errorf("scan row-group target mapping: %w", err)
		}
		if timestampColumn.Valid {
			value := timestampColumn.String
			target.timestampColumn = &value
		}
		if groupKey.Valid {
			value := groupKey.String
			target.groupKey = &value
		}
		if writeInterval.Valid {
			value := int(writeInterval.Int64)
			target.writeIntervalSeconds = &value
		}
		target.connectorID = strings.TrimSpace(target.connectorID)
		target.tableSchema = strings.TrimSpace(target.tableSchema)
		target.tableName = strings.TrimSpace(target.tableName)
		target.columnName = strings.TrimSpace(target.columnName)
		target.writeMode = strings.TrimSpace(target.writeMode)
		target.pointID = strings.TrimSpace(target.pointID)
		target.id = strings.TrimSpace(target.id)
		target.tagID = strings.TrimSpace(target.tagID)
		targets = append(targets, target)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate row-group target mappings: %w", err)
	}
	return targets, nil
}

func readMigrationPreviewSourcesForPoint(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	pointID string,
) ([]writeGroupMigrationSource, error) {
	rows, err := tx.QueryContext(ctx, repo.query(`
		SELECT m.id, m.point_id, m.tag_id, m.transform_pipeline, m.status,
			COALESCE(m.rule_candidate_id, ''), COALESCE(m.proposed_signature, ''),
			COALESCE(m.last_applied_signature, ''), COALESCE(m.blocking_reason, ''), m.enabled,
			d.id, d.name, d.protocol, d.connection_config, d.status,
			p.id, p.device_id, p.name, p.address, p.function, p.data_type,
			COALESCE(p.data_format, ''), p.mode, p.enabled,
			t.id, t.key, t.display_name, t.data_type, COALESCE(t.unit, ''),
			t.status, COALESCE(t.labels, '')
		FROM mappings m
		LEFT JOIN points p ON p.id = m.point_id
		LEFT JOIN devices d ON d.id = p.device_id
		LEFT JOIN tags t ON t.id = m.tag_id
		WHERE m.point_id = $1
		ORDER BY m.id ASC
	`), pointID)
	if err != nil {
		return nil, fmt.Errorf("read row-group source mappings: %w", err)
	}
	defer rows.Close()
	sources := make([]writeGroupMigrationSource, 0)
	for rows.Next() {
		source, scanErr := scanMigrationPreviewSource(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan row-group source mapping: %w", scanErr)
		}
		sources = append(sources, source)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate row-group source mappings: %w", err)
	}
	return sources, nil
}
