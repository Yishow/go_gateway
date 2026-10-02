package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type migrationPreviewRowScanner interface {
	Scan(...any) error
}

func readMigrationPreviewTarget(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	id string,
) (writeGroupMigrationTarget, error) {
	target, err := scanMigrationPreviewTarget(tx.QueryRowContext(ctx, repo.query(`
		SELECT id, tag_id, connector_id, table_schema, table_name, column_name,
			write_mode, timestamp_column, group_key, write_interval_seconds, enabled
		FROM database_target_mappings
		WHERE id = $1
	`), id))
	if errors.Is(err, sql.ErrNoRows) {
		return writeGroupMigrationTarget{}, writeGroupNotFound("preview migration source")
	}
	if err != nil {
		return writeGroupMigrationTarget{}, fmt.Errorf("read migration preview target mapping: %w", err)
	}
	return target, nil
}

func readMigrationPreviewTargetsForTag(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	tagID string,
) ([]writeGroupMigrationTarget, error) {
	return readMigrationPreviewTargets(ctx, tx, repo, `
		SELECT id, tag_id, connector_id, table_schema, table_name, column_name,
			write_mode, timestamp_column, group_key, write_interval_seconds, enabled
		FROM database_target_mappings
		WHERE tag_id = $1
		ORDER BY id ASC
	`, tagID)
}

func readMigrationPreviewTargetsForConnector(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	connectorID string,
) ([]writeGroupMigrationTarget, error) {
	return readMigrationPreviewTargets(ctx, tx, repo, `
		SELECT id, tag_id, connector_id, table_schema, table_name, column_name,
			write_mode, timestamp_column, group_key, write_interval_seconds, enabled
		FROM database_target_mappings
		WHERE connector_id = $1
		ORDER BY id ASC
	`, connectorID)
}

func readMigrationPreviewTargets(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	query string,
	arg string,
) ([]writeGroupMigrationTarget, error) {
	rows, err := tx.QueryContext(ctx, repo.query(query), arg)
	if err != nil {
		return nil, fmt.Errorf("read migration preview target mappings: %w", err)
	}
	targets := make([]writeGroupMigrationTarget, 0)
	for rows.Next() {
		target, scanErr := scanMigrationPreviewTarget(rows)
		if scanErr != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("scan migration preview target mapping: %w", scanErr)
		}
		targets = append(targets, target)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("iterate migration preview target mappings: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close migration preview target mappings: %w", err)
	}
	return targets, nil
}

func scanMigrationPreviewTarget(row migrationPreviewRowScanner) (writeGroupMigrationTarget, error) {
	var target writeGroupMigrationTarget
	var timestampColumn, groupKey sql.NullString
	var writeInterval sql.NullInt64
	if err := row.Scan(
		&target.id,
		&target.tagID,
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
		return writeGroupMigrationTarget{}, err
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
	return target, nil
}

func readMigrationPreviewSources(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	tagID string,
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
		WHERE m.tag_id = $1
		ORDER BY m.id ASC
	`), tagID)
	if err != nil {
		return nil, fmt.Errorf("read migration preview source mappings: %w", err)
	}
	sources := make([]writeGroupMigrationSource, 0)
	for rows.Next() {
		source, scanErr := scanMigrationPreviewSource(rows)
		if scanErr != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("scan migration preview source mapping: %w", scanErr)
		}
		sources = append(sources, source)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("iterate migration preview source mappings: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close migration preview source mappings: %w", err)
	}
	return sources, nil
}

func scanMigrationPreviewSource(row migrationPreviewRowScanner) (writeGroupMigrationSource, error) {
	var source writeGroupMigrationSource
	var (
		deviceID, deviceName, protocol, connectionConfig, deviceStatus sql.NullString
		pointID, pointDeviceID, pointName, pointAddress, pointFunction sql.NullString
		pointDataType, pointDataFormat, pointMode                      sql.NullString
		pointEnabled                                                   sql.NullBool
		tagID, tagKey, tagDisplayName, tagDataType, tagUnit            sql.NullString
		tagStatus, tagLabels                                           sql.NullString
		transformPipeline                                              sql.NullString
	)
	if err := row.Scan(
		&source.mapping.ID,
		&source.mapping.PointID,
		&source.mapping.TagID,
		&transformPipeline,
		&source.mapping.Status,
		&source.mapping.RuleCandidateID,
		&source.mapping.ProposedSignature,
		&source.mapping.LastAppliedSignature,
		&source.mapping.BlockingReason,
		&source.mapping.Enabled,
		&deviceID,
		&deviceName,
		&protocol,
		&connectionConfig,
		&deviceStatus,
		&pointID,
		&pointDeviceID,
		&pointName,
		&pointAddress,
		&pointFunction,
		&pointDataType,
		&pointDataFormat,
		&pointMode,
		&pointEnabled,
		&tagID,
		&tagKey,
		&tagDisplayName,
		&tagDataType,
		&tagUnit,
		&tagStatus,
		&tagLabels,
	); err != nil {
		return writeGroupMigrationSource{}, err
	}
	source.mapping.TransformPipeline = nullableString(transformPipeline)
	source.mappingState = strings.TrimSpace(source.mapping.Status)
	source.deviceStatus = nullableString(deviceStatus)
	source.tagStatus = nullableString(tagStatus)
	source.source = sourceRevisionSource{
		DeviceID:         nullableString(deviceID),
		DeviceName:       nullableString(deviceName),
		Protocol:         nullableString(protocol),
		ConnectionConfig: nullableString(connectionConfig),
		PointID:          nullableString(pointID),
		PointDeviceID:    nullableString(pointDeviceID),
		PointName:        nullableString(pointName),
		PointAddress:     nullableString(pointAddress),
		PointFunction:    nullableString(pointFunction),
		PointDataType:    nullableString(pointDataType),
		PointDataFormat:  nullableString(pointDataFormat),
		PointMode:        nullableString(pointMode),
		PointEnabled:     pointEnabled.Valid && pointEnabled.Bool,
		TagID:            nullableString(tagID),
		TagKey:           nullableString(tagKey),
		TagDisplayName:   nullableString(tagDisplayName),
		TagDataType:      nullableString(tagDataType),
		TagUnit:          nullableString(tagUnit),
		TagStatus:        nullableString(tagStatus),
		TagLabels:        nullableString(tagLabels),
	}
	return source, nil
}

func nullableString(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
