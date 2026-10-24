package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func (r *SQLWriteGroupRepository) loadBasicManagedDeviceID(ctx context.Context, runner writeGroupSQLRunner, group *WriteGroup) error {
	group.BasicManagedDeviceID = ""
	if group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged {
		return nil
	}
	var deviceID string
	err := runner.QueryRowContext(ctx, r.query(`
		SELECT device_id
		FROM write_group_basic_keys
		WHERE workspace_id = $1 AND group_id = $2 AND canonical_role = $3
	`), group.WorkspaceID, group.ID, basicManagedCanonicalRole).Scan(&deviceID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read basic managed owner: %w", err)
	}
	group.BasicManagedDeviceID = strings.TrimSpace(deviceID)
	return nil
}

func (r *SQLWriteGroupRepository) loadMembers(ctx context.Context, runner writeGroupSQLRunner, group *WriteGroup) error {
	rows, err := runner.QueryContext(ctx, r.query(`
		SELECT device_id, point_id, tag_id, entity_key, source_revision, mapping_revision,
			measurement_id, target_column, required, max_age_seconds, circuit_id, circuit_name
		FROM write_group_members
		WHERE group_id = $1
		ORDER BY member_index ASC
	`), group.ID)
	if err != nil {
		return fmt.Errorf("list write-group members: %w", err)
	}
	group.Members = make([]WriteGroupMember, 0)
	for rows.Next() {
		var member WriteGroupMember
		var measurementID sql.NullString
		var maxAgeSeconds sql.NullInt64
		var entityKey, circuitID, circuitName sql.NullString
		if err := rows.Scan(&member.DeviceID, &member.PointID, &member.TagID, &entityKey, &member.SourceRevision,
			&member.MappingRevision, &measurementID, &member.TargetColumn, &member.Required, &maxAgeSeconds,
			&circuitID, &circuitName); err != nil {
			_ = rows.Close()
			return fmt.Errorf("scan write-group member: %w", err)
		}
		if entityKey.Valid {
			member.EntityKey = entityKey.String
		}
		member.CircuitID, member.CircuitName = circuitID.String, circuitName.String
		if measurementID.Valid && strings.TrimSpace(measurementID.String) != "" {
			value := measurementID.String
			member.MeasurementID = &value
		}
		if maxAgeSeconds.Valid {
			value := int(maxAgeSeconds.Int64)
			member.MaxAgeSeconds = &value
		}
		group.Members = append(group.Members, member)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return fmt.Errorf("iterate write-group members: %w", err)
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close write-group members: %w", err)
	}
	return nil
}
