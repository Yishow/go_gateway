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
	"go-gateway/internal/datalink/schema"
)

func (r *SQLWriteGroupRepository) validate(ctx context.Context, runner writeGroupSQLRunner, group *WriteGroup) error {
	return r.validateWithMeasurements(ctx, runner, group, make(map[string]*measurement.MeasurementDefinition))
}

func (r *SQLWriteGroupRepository) validateWithMeasurements(ctx context.Context, runner writeGroupSQLRunner, group *WriteGroup, definitions map[string]*measurement.MeasurementDefinition) error {
	if err := r.validateDefinitionWithMeasurements(ctx, runner, group, definitions); err != nil {
		return err
	}
	return r.checkProfileSourceOwnership(ctx, runner, group)
}

// validateDefinitionWithMeasurements validates every non-ownership property.
// Ordinary callers always use validateWithMeasurements above. The successor
// preparation path separately proves its exact predecessor-only overlap.
func (r *SQLWriteGroupRepository) validateDefinitionWithMeasurements(ctx context.Context, runner writeGroupSQLRunner, group *WriteGroup, definitions map[string]*measurement.MeasurementDefinition) error {
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
	// A long group's table is derived from its storage scope below, so it may
	// arrive without a name.
	if group.Destination.TableName == "" && !group.isLongLayout() {
		return fmt.Errorf("%w: destination table name is required", ErrWriteGroupValidation)
	}
	if group.Destination.StorageStrategy == "" {
		group.Destination.StorageStrategy = WriteGroupStorageStrategyCustom
	}
	if group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged && group.Destination.StorageStrategy != WriteGroupStorageStrategyCustom {
		return fmt.Errorf("%w: unsupported storage strategy %q", ErrWriteGroupValidation, group.Destination.StorageStrategy)
	}
	if err := normalizeWriteGroupProfile(group); err != nil {
		return err
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
	if group.StorageMode == WriteGroupStorageModeNativeMonthly && connectorKind != writeGroupConnectorKindMySQL && connectorKind != writeGroupConnectorKindSQLServer {
		return fmt.Errorf("%w: native monthly partitioning requires a verified MySQL-family or SQL Server destination", ErrWriteGroupValidation)
	}
	if connectorKind != writeGroupConnectorKindSQLite && connectorKind != writeGroupConnectorKindPostgres {
		if connectorKind != writeGroupConnectorKindMySQL && connectorKind != writeGroupConnectorKindSQLServer || !group.isLongLayout() ||
			group.StorageMode != WriteGroupStorageModePlain && group.StorageMode != WriteGroupStorageModeNativeMonthly {
			return fmt.Errorf("%w: %w: %s", ErrWriteGroupUnsupportedConnectorKind, ErrWriteGroupValidation, connectorKind)
		}
	}
	if err := resolveWriteGroupDestinationScope(connectorKind, connectorConfig, &group.Destination); err != nil {
		return err
	}
	if group.isLongLayout() {
		// The shared table is named by its storage scope. A caller may echo that
		// name back, but cannot choose another one.
		derived := longStorageScopeOf(group).tableName()
		if group.Destination.TableName != "" && group.Destination.TableName != derived {
			return fmt.Errorf("%w: a long group's table is chosen by its storage scope", ErrWriteGroupValidation)
		}
		group.Destination.TableName = derived
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

		if member.MeasurementID == nil || strings.TrimSpace(*member.MeasurementID) == "" {
			if member.ExactRecordingSource != nil {
				return fmt.Errorf("%w: member %d carries exact recording source without measurement_id", ErrWriteGroupValidation, i)
			}
			member.MeasurementID = nil
			if err := applyRecordingMetadata(group, member, i, schema.DataType(source.TagDataType)); err != nil {
				return err
			}
			continue
		}
		measurementID := strings.TrimSpace(*member.MeasurementID)
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
		if err := r.deriveAndValidateExactRecordingSourceInTx(ctx, runner, group, member, i, schema.DataType(source.TagDataType), definitions); err != nil {
			return err
		}
		if err := applyRecordingMetadata(group, member, i, schema.DataType(source.TagDataType)); err != nil {
			return err
		}
		if member.ValueRole == WriteGroupValueRoleCounterRaw {
			if err := confirmCounterBinding(definitions[*member.MeasurementID]); err != nil {
				return fmt.Errorf("member %d: %w", i, err)
			}
		}
	}
	return nil
}

// applyRecordingMetadata verifies and derives the member recording metadata of
// a new recording shape: a mean reduction, the new wide profile layout, or a
// member that explicitly carries metadata. Legacy members are never touched
// and never backfilled. The source type is the tag's own data type; a caller
// type that differs from it is refused instead of silently saved.
func applyRecordingMetadata(group *WriteGroup, member *WriteGroupMember, index int, tagType schema.DataType) error {
	newShape := group.RecordingProfile.EffectiveSampleReduction() == WriteGroupSampleReductionMean ||
		group.EffectiveLayoutVersion() == WriteGroupLayoutManagedWideProfileV1 ||
		memberMetadataUnsupported(member)
	if !newShape {
		return nil
	}
	role := member.ValueRole
	if role == "" {
		// Only an instantaneous shape declares roles; a cumulative profile
		// counts and never averages, so its members keep the raw mapping
		// without one.
		if group.RecordingProfile.Kind == WriteGroupProfileInstantaneous {
			return fmt.Errorf("%w: member %d needs an explicit value role for this recording shape", ErrWriteGroupValidation, index)
		}
	}
	if member.SourceValueType != "" && member.SourceValueType != tagType {
		return fmt.Errorf("%w: member %d source type %q does not match the tag's %q",
			ErrWriteGroupValidation, index, member.SourceValueType, tagType)
	}
	result, err := deriveMemberResultType(group.RecordingProfile.EffectiveSampleReduction(), role, tagType)
	if err != nil {
		return fmt.Errorf("member %d: %w", index, err)
	}
	expectedResult := result
	if member.ExactRecordingSource != nil {
		expectedResult = measurement.ExactDecimal
	}
	if member.ResultValueType != "" && member.ResultValueType != expectedResult {
		return fmt.Errorf("%w: member %d result type %q does not match the derived %q",
			ErrWriteGroupValidation, index, member.ResultValueType, expectedResult)
	}
	member.SourceValueType = tagType
	member.ResultValueType = expectedResult
	if role == WriteGroupValueRoleCounterRaw {
		if member.MeasurementID == nil || strings.TrimSpace(*member.MeasurementID) == "" {
			return fmt.Errorf("%w: member %d needs the confirmed counter binding for counter_raw", ErrWriteGroupValidation, index)
		}
		// The saved counter evidence never turns a bool or text tag into a
		// counter: only a numeric source can feed one, as the cumulative gate
		// already requires.
		if exact, ok := measurement.ExactTypeForTag(tagType); !ok || exact == measurement.ExactBool || exact == measurement.ExactText {
			return fmt.Errorf("%w: member %d counter_raw requires a numeric source tag, not %q",
				ErrWriteGroupValidation, index, tagType)
		}
	}
	return nil
}

// confirmCounterBinding reuses the existing confirmed counter evidence: the
// saved definition must itself confirm counter semantics on the same local
// transaction, never a guess from a name or a unit.
func confirmCounterBinding(def *measurement.MeasurementDefinition) error {
	if def == nil {
		return fmt.Errorf("read write-group measurement: the counter binding is unavailable")
	}
	if confirmation := measurement.ConfirmCounter(*def); confirmation != measurement.CounterConfirmed {
		return fmt.Errorf("%w: the counter binding is not confirmed (%s)", ErrWriteGroupValidation, confirmation)
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
