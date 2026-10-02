package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"go-gateway/internal/datalink/measurement"
)

func (r *SQLWriteGroupRepository) validate(ctx context.Context, runner writeGroupSQLRunner, group *WriteGroup) error {
	if strings.TrimSpace(group.WorkspaceID) == "" {
		return fmt.Errorf("%w: workspace id is required", ErrWriteGroupValidation)
	}
	group.WorkspaceID = strings.TrimSpace(group.WorkspaceID)
	group.Name = strings.TrimSpace(group.Name)
	if group.Name == "" {
		return fmt.Errorf("%w: group name is required", ErrWriteGroupValidation)
	}
	if len(group.Members) == 0 {
		return fmt.Errorf("%w: at least one member is required", ErrWriteGroupValidation)
	}
	if strings.TrimSpace(group.Destination.ConnectorID) == "" || strings.TrimSpace(group.Destination.ConnectorRevision) == "" {
		return fmt.Errorf("%w: connector id and identity revision are required", ErrWriteGroupValidation)
	}
	group.Destination.ConnectorID = strings.TrimSpace(group.Destination.ConnectorID)
	group.Destination.ConnectorRevision = strings.TrimSpace(group.Destination.ConnectorRevision)
	group.Destination.TableSchema = strings.TrimSpace(group.Destination.TableSchema)
	group.Destination.TableName = strings.TrimSpace(group.Destination.TableName)
	group.RowPolicy.GroupKeyColumns = normalizeStringSet(group.RowPolicy.GroupKeyColumns)
	group.RowPolicy.UniqueKeyColumns = normalizeStringSet(group.RowPolicy.UniqueKeyColumns)
	if group.Destination.TableName == "" {
		return fmt.Errorf("%w: destination table name is required", ErrWriteGroupValidation)
	}
	if group.Destination.StorageStrategy == "" {
		group.Destination.StorageStrategy = WriteGroupStorageStrategyCustom
	}
	if group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged && group.Destination.StorageStrategy != WriteGroupStorageStrategyCustom {
		return fmt.Errorf("%w: unsupported storage strategy %q", ErrWriteGroupValidation, group.Destination.StorageStrategy)
	}

	workspaceRecord, err := readWorkspaceScope(ctx, runner, r.query, group.WorkspaceID)
	if err != nil {
		return err
	}
	var connectorKind, connectorRevision, connectorConfig string
	var connectorEnabled bool
	err = runner.QueryRowContext(ctx, r.query(`
		SELECT kind, identity_revision, enabled, connection_config
		FROM database_connectors
		WHERE id = $1
	`), group.Destination.ConnectorID).Scan(&connectorKind, &connectorRevision, &connectorEnabled, &connectorConfig)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %w", ErrWriteGroupConnectorNotFound, ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("read write-group connector: %w", err)
	}
	if !connectorEnabled {
		return fmt.Errorf("%w: connector is disabled", ErrWriteGroupValidation)
	}
	if connectorRevision != group.Destination.ConnectorRevision {
		return fmt.Errorf("%w: %w: connector identity revision is stale", ErrWriteGroupConnectorRevisionConflict, ErrSetupRevisionConflict)
	}
	connectorKind = strings.ToLower(strings.TrimSpace(connectorKind))
	if connectorKind != writeGroupConnectorKindSQLite && connectorKind != writeGroupConnectorKindPostgres {
		return fmt.Errorf("%w: %w: %s", ErrWriteGroupUnsupportedConnectorKind, ErrWriteGroupValidation, connectorKind)
	}
	if err := resolveWriteGroupDestinationScope(connectorKind, connectorConfig, &group.Destination); err != nil {
		return err
	}

	seenMembers := make(map[string]struct{}, len(group.Members))
	seenEntityColumns := make(map[string]struct{}, len(group.Members))
	for i := range group.Members {
		member := &group.Members[i]
		member.DeviceID = strings.TrimSpace(member.DeviceID)
		member.PointID = strings.TrimSpace(member.PointID)
		member.TagID = strings.TrimSpace(member.TagID)
		if strings.TrimSpace(member.EntityKey) == "" {
			member.EntityKey = ""
		}
		member.TargetColumn = strings.TrimSpace(member.TargetColumn)
		if member.DeviceID == "" || member.PointID == "" || member.TagID == "" || member.TargetColumn == "" {
			return fmt.Errorf("%w: member %d has missing identity or target column", ErrWriteGroupValidation, i)
		}
		memberKey := strings.Join([]string{member.DeviceID, member.PointID, member.TagID}, "\x00")
		if _, exists := seenMembers[memberKey]; exists {
			return fmt.Errorf("%w: duplicate member %d", ErrWriteGroupValidation, i)
		}
		seenMembers[memberKey] = struct{}{}
		entityColumnKey := strings.TrimSpace(member.EntityKey) + "\x00" + member.TargetColumn
		if _, exists := seenEntityColumns[entityColumnKey]; exists {
			return fmt.Errorf("%w: member %d collides on entity and target column", ErrWriteGroupValidation, i)
		}
		seenEntityColumns[entityColumnKey] = struct{}{}
		if !slices.Contains(workspaceRecord.OrderedDeviceIDs, member.DeviceID) {
			return fmt.Errorf("%w: %w", ErrWriteGroupNotFound, ErrNotFound)
		}

		var pointDeviceID string
		err = runner.QueryRowContext(ctx, r.query(`SELECT device_id FROM points WHERE id = $1`), member.PointID).Scan(&pointDeviceID)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: %w", ErrWriteGroupNotFound, ErrNotFound)
		}
		if err != nil {
			return fmt.Errorf("read write-group point: %w", err)
		}
		if pointDeviceID != member.DeviceID {
			return fmt.Errorf("%w: %w", ErrWriteGroupNotFound, ErrNotFound)
		}
		var tagExists int
		err = runner.QueryRowContext(ctx, r.query(`SELECT 1 FROM tags WHERE id = $1`), member.TagID).Scan(&tagExists)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: %w", ErrWriteGroupNotFound, ErrNotFound)
		}
		if err != nil {
			return fmt.Errorf("read write-group tag: %w", err)
		}

		var mapping mappingRevisionSource
		mappingQuery := `
			SELECT id, point_id, tag_id, transform_pipeline, status,
				COALESCE(rule_candidate_id, ''), COALESCE(proposed_signature, ''),
				COALESCE(last_applied_signature, ''), COALESCE(blocking_reason, ''), enabled
			FROM mappings
			WHERE point_id = $1 AND tag_id = $2 AND enabled = ` + r.trueLiteral() + `
			ORDER BY id ASC
			LIMIT 1
		`
		err = runner.QueryRowContext(ctx, r.query(mappingQuery), member.PointID, member.TagID).Scan(&mapping.ID, &mapping.PointID, &mapping.TagID,
			&mapping.TransformPipeline, &mapping.Status, &mapping.RuleCandidateID,
			&mapping.ProposedSignature, &mapping.LastAppliedSignature, &mapping.BlockingReason,
			&mapping.Enabled)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: %w", ErrWriteGroupNotFound, ErrNotFound)
		}
		if err != nil {
			return fmt.Errorf("read write-group mapping: %w", err)
		}
		mappingRevision, err := mapping.revision()
		if err != nil {
			return fmt.Errorf("%w: derive mapping revision: %w", ErrWriteGroupValidation, err)
		}
		if member.MappingRevision != "" && member.MappingRevision != mappingRevision {
			return fmt.Errorf("%w: %w: member %d mapping revision is stale", ErrWriteGroupSourceRevisionConflict, ErrSetupRevisionConflict, i)
		}
		member.MappingRevision = mappingRevision

		var source sourceRevisionSource
		err = runner.QueryRowContext(ctx, r.query(`
			SELECT d.id, d.name, d.protocol, d.connection_config,
				p.id, p.device_id, p.name, p.address, p.function, p.data_type,
				COALESCE(p.data_format, ''), p.mode, p.enabled,
				t.id, t.key, t.display_name, t.data_type, COALESCE(t.unit, ''),
				t.status, COALESCE(t.labels, '')
			FROM devices d
			JOIN points p ON p.device_id = d.id
			CROSS JOIN tags t
			WHERE d.id = $1 AND p.id = $2 AND t.id = $3
			`), member.DeviceID, member.PointID, member.TagID).Scan(
			&source.DeviceID, &source.DeviceName, &source.Protocol, &source.ConnectionConfig,
			&source.PointID, &source.PointDeviceID, &source.PointName, &source.PointAddress, &source.PointFunction,
			&source.PointDataType, &source.PointDataFormat, &source.PointMode, &source.PointEnabled,
			&source.TagID, &source.TagKey, &source.TagDisplayName, &source.TagDataType, &source.TagUnit,
			&source.TagStatus, &source.TagLabels)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: %w", ErrWriteGroupNotFound, ErrNotFound)
		}
		if err != nil {
			return fmt.Errorf("read write-group source: %w", err)
		}
		sourceRevision, err := source.revision()
		if err != nil {
			return fmt.Errorf("%w: derive source revision: %w", ErrWriteGroupValidation, err)
		}
		if member.SourceRevision != "" && member.SourceRevision != sourceRevision {
			return fmt.Errorf("%w: %w: member %d source revision is stale", ErrWriteGroupSourceRevisionConflict, ErrSetupRevisionConflict, i)
		}
		member.SourceRevision = sourceRevision

		if member.MeasurementID == nil {
			continue
		}
		measurementID := strings.TrimSpace(*member.MeasurementID)
		if measurementID == "" {
			member.MeasurementID = nil
			continue
		}
		var measurementDeviceID, measurementPointID string
		var measurementTagID sql.NullString
		err = runner.QueryRowContext(ctx, r.query(`
			SELECT device_id, point_id, tag_id
			FROM measurement_definitions
			WHERE id = $1 AND workspace_id = $2
		`), measurementID, group.WorkspaceID).Scan(&measurementDeviceID, &measurementPointID, &measurementTagID)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: %w", ErrWriteGroupNotFound, ErrNotFound)
		}
		if err != nil {
			return fmt.Errorf("read write-group measurement: %w", err)
		}
		if measurementDeviceID != member.DeviceID || measurementPointID != member.PointID || !measurementTagID.Valid || measurementTagID.String != member.TagID {
			return fmt.Errorf("%w: %w", ErrWriteGroupNotFound, ErrNotFound)
		}
		member.MeasurementID = &measurementID
	}
	return nil
}

func readWorkspaceScope(ctx context.Context, runner writeGroupSQLRunner, query func(string) string, workspaceID string) (*Record, error) {
	var payload string
	err := runner.QueryRowContext(ctx, query(`SELECT value FROM system_settings WHERE key = $1`), storageKey).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %w", ErrWriteGroupNotFound, ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("read write-group workspace: %w", err)
	}
	var record Record
	if err := json.Unmarshal([]byte(payload), &record); err != nil {
		return nil, fmt.Errorf("decode write-group workspace: %w", err)
	}
	if record.ID != workspaceID {
		return nil, fmt.Errorf("%w: %w", ErrWriteGroupNotFound, ErrNotFound)
	}
	return &record, nil
}

type mappingRevisionSource struct {
	ID, PointID, TagID, TransformPipeline, Status                            string
	RuleCandidateID, ProposedSignature, LastAppliedSignature, BlockingReason string
	Enabled                                                                  bool
}

func (m mappingRevisionSource) revision() (string, error) {
	return hashWriteGroupRevision(m.ID, m.PointID, m.TagID, m.TransformPipeline, m.Status,
		m.RuleCandidateID, m.ProposedSignature, m.LastAppliedSignature, m.BlockingReason, fmt.Sprint(m.Enabled))
}

type sourceRevisionSource struct {
	DeviceID, DeviceName, Protocol, ConnectionConfig                                                          string
	PointID, PointDeviceID, PointName, PointAddress, PointFunction, PointDataType, PointDataFormat, PointMode string
	PointEnabled                                                                                              bool
	TagID, TagKey, TagDisplayName, TagDataType, TagUnit, TagStatus, TagLabels                                 string
}

func (s sourceRevisionSource) revision() (string, error) {
	return hashWriteGroupRevision(s.DeviceID, s.DeviceName, s.Protocol, s.ConnectionConfig,
		s.PointID, s.PointDeviceID, s.PointName, s.PointAddress, s.PointFunction, s.PointDataType,
		s.PointDataFormat, s.PointMode, fmt.Sprint(s.PointEnabled), s.TagID, s.TagKey,
		s.TagDisplayName, s.TagDataType, s.TagUnit, s.TagStatus, s.TagLabels)
}

func hashWriteGroupRevision(parts ...string) (string, error) {
	return measurement.HashStringTuple(parts...)
}
