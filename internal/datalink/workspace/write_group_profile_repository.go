package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// nullableText stores an empty circuit field as NULL, the legacy value.
func nullableText(value string) any {
	if value == "" {
		return nil
	}
	return value
}

// profileColumns holds the nullable profile and storage columns of one group,
// each bound as a SQL parameter exactly as stored (nil is NULL).
type profileColumns struct {
	layout, kind, interval, unit, mode, reduction, shape any
}

// profileColumnValues maps a group to its nullable profile and storage columns.
func profileColumnValues(group *WriteGroup) profileColumns {
	columns := profileColumns{
		layout: nullableText(string(group.LayoutVersion)),
		mode:   nullableText(string(group.StorageMode)),
	}
	if group.RecordingProfile != nil {
		columns.kind, columns.interval, columns.unit = string(group.RecordingProfile.Kind), group.RecordingProfile.Interval, string(group.RecordingProfile.Unit)
		columns.reduction = nullableText(string(group.RecordingProfile.SampleReduction))
		columns.shape = nullableText(string(group.RecordingProfile.ResultShape))
		if group.RecordingProfile.Kind == WriteGroupProfileMixed {
			columns.interval, columns.unit = nil, nil
		}
	}
	return columns
}

// applyProfileColumns decodes the nullable profile columns. All-NULL columns
// are a pre-profile group; a partially written profile is corrupt and fails
// loudly instead of being read as wide. An omitted reduction stays empty and
// reads as latest through EffectiveSampleReduction; it is never backfilled, and
// a stored unknown value fails instead of degrading to latest.
func applyProfileColumns(group *WriteGroup, layout, kind, unit, mode sql.NullString, interval sql.NullInt64, reduction, shape sql.NullString) error {
	group.LayoutVersion = WriteGroupLayoutVersion(layout.String)
	group.StorageMode = WriteGroupStorageMode(mode.String)
	if !supportedLayoutVersion(group.LayoutVersion) {
		return fmt.Errorf("decode write-group layout version %q: %w", layout.String, ErrWriteGroupValidation)
	}
	if group.StorageMode != "" && group.StorageMode != WriteGroupStorageModePlain &&
		group.StorageMode != WriteGroupStorageModeTimescale && group.StorageMode != WriteGroupStorageModeNativeMonthly {
		return fmt.Errorf("decode write-group storage mode %q: %w", mode.String, ErrWriteGroupValidation)
	}
	if kind.String == string(WriteGroupProfileMixed) {
		if !shape.Valid || shape.String != string(WriteGroupResultShapeMemberCalculationV1) || interval.Valid || unit.Valid || reduction.Valid {
			return fmt.Errorf("decode mixed recording profile: %w", ErrWriteGroupValidation)
		}
		group.RecordingProfile = &WriteGroupRecordingProfile{Kind: WriteGroupProfileMixed, ResultShape: WriteGroupResultShape(shape.String)}
		return nil
	}
	if shape.Valid {
		return fmt.Errorf("decode legacy result shape: %w", ErrWriteGroupValidation)
	}
	if !kind.Valid && !interval.Valid && !unit.Valid {
		if reduction.Valid && reduction.String != "" {
			return fmt.Errorf("decode write-group sample reduction %q without a recording profile: %w", reduction.String, ErrWriteGroupValidation)
		}
		group.RecordingProfile = nil
		return nil
	}
	if !kind.Valid || !interval.Valid || !unit.Valid {
		return fmt.Errorf("decode write-group recording profile: %w", ErrWriteGroupValidation)
	}
	group.RecordingProfile = &WriteGroupRecordingProfile{
		Kind: WriteGroupProfileKind(kind.String), Interval: int(interval.Int64), Unit: WriteGroupIntervalUnit(unit.String),
	}
	if reduction.Valid {
		if reduction.String == "" {
			return fmt.Errorf("decode write-group sample reduction: %w", ErrWriteGroupValidation)
		}
		switch parsed := WriteGroupSampleReduction(reduction.String); parsed {
		case WriteGroupSampleReductionLatest, WriteGroupSampleReductionMean:
			group.RecordingProfile.SampleReduction = parsed
		default:
			return fmt.Errorf("decode write-group sample reduction %q: %w", reduction.String, ErrWriteGroupValidation)
		}
	}
	return nil
}

// insertProfileKeyInTx records the create-once key for a profile group in the
// caller's transaction. A key held by another group is a conflict.
func (r *SQLWriteGroupRepository) insertProfileKeyInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup) error {
	key, err := circuitProfileKeyOf(group)
	if err != nil {
		return err
	}
	lineageID, err := r.initializeCanonicalLineageInTx(ctx, tx, key.WorkspaceID, group.ID)
	if err != nil {
		return err
	}
	result, err := tx.ExecContext(ctx, r.query(`
		INSERT INTO write_group_profile_keys (
			workspace_id, circuit_id, selection_scope, profile_kind, canonical_role, group_id, lineage_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (workspace_id, circuit_id, selection_scope, profile_kind, canonical_role) DO NOTHING
	`), key.WorkspaceID, key.CircuitID, key.SelectionScope, string(key.Kind), key.Role, group.ID, lineageID)
	if err != nil {
		return fmt.Errorf("insert circuit profile key: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("confirm circuit profile key: %w", err)
	}
	if affected != 1 {
		return fmt.Errorf("circuit profile key already belongs to another group: %w", ErrWriteGroupCircuitProfileConflict)
	}
	return nil
}

func (r *SQLWriteGroupRepository) findProfileGroup(ctx context.Context, key circuitProfileKey) (*WriteGroup, error) {
	return r.findCanonicalGroup(ctx, key.WorkspaceID, `
		SELECT group_id, lineage_id
		FROM write_group_profile_keys
		WHERE workspace_id = $1 AND circuit_id = $2 AND selection_scope = $3
			AND profile_kind = $4 AND canonical_role = $5
	`, key.WorkspaceID, key.CircuitID, key.SelectionScope, string(key.Kind), key.Role)
}

// loadProfileKey returns the key held by a group; found is false for a group
// created without a profile.
func (r *SQLWriteGroupRepository) loadProfileKey(ctx context.Context, runner writeGroupSQLRunner, workspaceID, groupID string) (key circuitProfileKey, found bool, err error) {
	metadata, found, err := r.loadWriteGroupLineage(ctx, runner, workspaceID, groupID)
	if err != nil || !found {
		return circuitProfileKey{}, false, err
	}
	key = circuitProfileKey{WorkspaceID: workspaceID}
	var kind string
	err = runner.QueryRowContext(ctx, r.query(`
		SELECT circuit_id, selection_scope, profile_kind, canonical_role
		FROM write_group_profile_keys
		WHERE workspace_id = $1 AND group_id = $2
	`), workspaceID, metadata.GenesisGroupID).Scan(&key.CircuitID, &key.SelectionScope, &kind, &key.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return circuitProfileKey{}, false, nil
	}
	if err != nil {
		return circuitProfileKey{}, false, fmt.Errorf("read circuit profile owner: %w", err)
	}
	key.Kind = WriteGroupProfileKind(kind)
	return key, true, nil
}

// checkProfileUpdateIdentity keeps a group's key and layout fixed after
// create: the key is its identity, and a layout change is an explicit
// cutover rather than an edit.
func (r *SQLWriteGroupRepository) checkProfileUpdateIdentity(ctx context.Context, runner writeGroupSQLRunner, candidate, existing *WriteGroup) error {
	if candidate.EffectiveLayoutVersion() != existing.EffectiveLayoutVersion() {
		return fmt.Errorf("%w: layout version cannot change on update", ErrWriteGroupValidation)
	}
	candidate.LayoutVersion = existing.LayoutVersion
	if candidate.StorageMode != existing.StorageMode {
		return fmt.Errorf("%w: storage mode cannot change on update", ErrWriteGroupValidation)
	}
	stored, found, err := r.loadProfileKey(ctx, runner, existing.WorkspaceID, existing.ID)
	if err != nil {
		return err
	}
	if !found {
		if candidate.RecordingProfile != nil {
			return fmt.Errorf("%w: a recording profile can only be set when the group is created", ErrWriteGroupValidation)
		}
		return nil
	}
	if candidate.RecordingProfile == nil {
		return fmt.Errorf("%w: a profile group cannot drop its recording profile", ErrWriteGroupValidation)
	}
	derived, err := circuitProfileKeyOf(candidate)
	if err != nil {
		return err
	}
	if derived != stored {
		return fmt.Errorf("%w: circuit and profile kind are fixed after create", ErrWriteGroupValidation)
	}
	return nil
}

// checkProfileSourceOwnership rejects a source already owned by another live
// group of the same profile kind and destination. The same source in the
// other profile kind is a distinct output and stays allowed.
func (r *SQLWriteGroupRepository) checkProfileSourceOwnership(ctx context.Context, runner writeGroupSQLRunner, group *WriteGroup) error {
	if err := r.checkMixedTargetOwnership(ctx, runner, group); err != nil {
		return err
	}
	if isMixedWriteGroup(group) {
		_, err := mixedRecordProfileKey(group)
		return err
	}
	if group.RecordingProfile == nil {
		return nil
	}
	rows, err := runner.QueryContext(ctx, r.query(`
		SELECT g.id, m.device_id, m.point_id, m.tag_id
		FROM write_groups g
		JOIN write_group_members m ON m.group_id = g.id
		WHERE g.workspace_id = $1 AND g.id <> $2 AND g.status <> $3 AND g.profile_kind = $4
			AND g.destination_connector_id = $5 AND g.destination_database = $6 AND g.destination_table_schema = $7
	`), group.WorkspaceID, group.ID, WriteGroupStatusDeleted, string(group.RecordingProfile.Kind),
		group.Destination.ConnectorID, group.Destination.Database, group.Destination.TableSchema)
	if err != nil {
		return fmt.Errorf("read circuit profile source ownership: %w", err)
	}
	defer rows.Close()
	wanted := make(map[[3]string]struct{}, len(group.Members))
	for _, member := range group.Members {
		wanted[[3]string{member.DeviceID, member.PointID, member.TagID}] = struct{}{}
	}
	for rows.Next() {
		var ownerID string
		var source [3]string
		if err := rows.Scan(&ownerID, &source[0], &source[1], &source[2]); err != nil {
			return fmt.Errorf("scan circuit profile source ownership: %w", err)
		}
		if _, owned := wanted[source]; owned {
			return fmt.Errorf("%w: source %s/%s/%s is already owned by group %s for the %s profile",
				ErrWriteGroupCircuitProfileConflict, source[0], source[1], source[2], ownerID, group.RecordingProfile.Kind)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate circuit profile source ownership: %w", err)
	}
	return nil
}
