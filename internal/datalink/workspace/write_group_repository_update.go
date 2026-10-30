package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// updateInTx is also used by the service's draft tombstone path. That path has
// already loaded and scoped the record, so it must retain the row even when a
// source member was removed after the original save.
func (r *SQLWriteGroupRepository) updateInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup, validate bool) error {
	if tx == nil {
		return fmt.Errorf("update write group in transaction: %w", ErrWriteGroupValidation)
	}
	if group == nil || group.ID == "" || group.WorkspaceID == "" {
		return fmt.Errorf("update write group: %w", ErrWriteGroupValidation)
	}
	existing, err := r.get(ctx, tx, group.WorkspaceID, group.ID)
	if err != nil {
		return err
	}
	if err := r.requireUnreservedWriteGroupInTx(ctx, tx, group.WorkspaceID, group.ID); err != nil {
		return err
	}
	candidate := cloneWriteGroup(group)
	// Basic ownership is read-only metadata; it has no write-group column and
	// must never be accepted from an update payload.
	candidate.BasicManagedDeviceID = ""
	if candidate.WritePolicy.RecreateMissing == nil {
		// An edit that does not mention the policy keeps what was saved; it
		// never creates one that was not confirmed before.
		candidate.WritePolicy.RecreateMissing = existing.WritePolicy.RecreateMissing
	}
	if validate {
		if candidate.Destination.StorageStrategy == WriteGroupStorageStrategyManaged && candidate.RowPolicy.RecordKeyColumn != "" && strings.TrimSpace(candidate.Destination.TableName) == "" && !candidate.isLongLayout() {
			candidate.Destination.TableName = managedGroupTableName(existing.ID)
		}
		if err := r.validate(ctx, tx, candidate); err != nil {
			return err
		}
		if err := r.checkProfileUpdateIdentity(ctx, tx, candidate, existing); err != nil {
			return err
		}
	}
	mergeWriteGroupServerFields(candidate, existing, r)
	if validate {
		if err := r.ensureStorageScopeInTx(ctx, tx, candidate); err != nil {
			return err
		}
	}
	if err := r.updateGroup(ctx, tx, candidate); err != nil {
		return err
	}
	*group = *cloneWriteGroup(candidate)
	return nil
}

func mergeWriteGroupServerFields(candidate, existing *WriteGroup, r *SQLWriteGroupRepository) {
	candidate.ID = existing.ID
	candidate.WorkspaceID = existing.WorkspaceID
	candidate.CreatedAt = existing.CreatedAt
	candidate.AppliedRevision = existing.AppliedRevision
	candidate.Migration = cloneWriteGroup(existing).Migration
	if sameWriteGroupSchema(candidate, existing) {
		candidate.Destination.SchemaRevision = existing.Destination.SchemaRevision
		candidate.Destination.SchemaDigest = existing.Destination.SchemaDigest
	} else {
		candidate.Destination.SchemaRevision = ""
		candidate.Destination.SchemaDigest = ""
	}
	candidate.Revision = r.newID()
	candidate.UpdatedAt = r.now()
}

func sameWriteGroupSchema(candidate, existing *WriteGroup) bool {
	if candidate == nil || existing == nil {
		return false
	}
	if candidate.Destination.ConnectorID != existing.Destination.ConnectorID ||
		candidate.Destination.ConnectorRevision != existing.Destination.ConnectorRevision ||
		candidate.Destination.Database != existing.Destination.Database ||
		candidate.Destination.TableSchema != existing.Destination.TableSchema ||
		candidate.Destination.TableName != existing.Destination.TableName ||
		candidate.Destination.StorageStrategy != existing.Destination.StorageStrategy ||
		candidate.RowPolicy.IntervalSeconds != existing.RowPolicy.IntervalSeconds ||
		candidate.RowPolicy.AllowedLatenessSeconds != existing.RowPolicy.AllowedLatenessSeconds ||
		candidate.RowPolicy.IncompletePolicy != existing.RowPolicy.IncompletePolicy ||
		candidate.RowPolicy.EntityKeyColumn != existing.RowPolicy.EntityKeyColumn ||
		!slices.Equal(candidate.RowPolicy.GroupKeyColumns, existing.RowPolicy.GroupKeyColumns) ||
		!slices.Equal(candidate.RowPolicy.UniqueKeyColumns, existing.RowPolicy.UniqueKeyColumns) ||
		candidate.RowPolicy.ValueColumn != existing.RowPolicy.ValueColumn ||
		candidate.RowPolicy.QualityColumn != existing.RowPolicy.QualityColumn ||
		candidate.RowPolicy.ProvenanceColumn != existing.RowPolicy.ProvenanceColumn {
		return false
	}
	// A profile edit changes the derived interval; stale schema evidence must go.
	if !sameRecordingProfile(candidate.RecordingProfile, existing.RecordingProfile) {
		return false
	}
	// The recovery policy is part of what a confirmation shows, so changing it
	// asks for a fresh confirmation.
	if !sameRecreateMissing(candidate.WritePolicy.RecreateMissing, existing.WritePolicy.RecreateMissing) {
		return false
	}
	if candidate.RowPolicy.RecordKeyColumn != existing.RowPolicy.RecordKeyColumn ||
		candidate.RowPolicy.BucketStartColumn != existing.RowPolicy.BucketStartColumn ||
		candidate.RowPolicy.GroupIDColumn != existing.RowPolicy.GroupIDColumn ||
		candidate.RowPolicy.DeviceIDColumn != existing.RowPolicy.DeviceIDColumn {
		return false
	}
	if len(candidate.Members) != len(existing.Members) {
		return false
	}
	for _, member := range candidate.Members {
		matched := slices.ContainsFunc(existing.Members, func(existingMember WriteGroupMember) bool {
			return member.DeviceID == existingMember.DeviceID &&
				member.PointID == existingMember.PointID &&
				member.TagID == existingMember.TagID &&
				member.TargetColumn == existingMember.TargetColumn &&
				member.EntityKey == existingMember.EntityKey &&
				member.ValueRole == existingMember.ValueRole &&
				sameStoredMetadata(string(member.SourceValueType), string(existingMember.SourceValueType)) &&
				sameStoredMetadata(string(member.ResultValueType), string(existingMember.ResultValueType)) &&
				sameExactRecordingSource(member.ExactRecordingSource, existingMember.ExactRecordingSource)
		})
		if !matched {
			return false
		}
	}
	return true
}

func (r *SQLWriteGroupRepository) updateGroup(ctx context.Context, tx *sql.Tx, group *WriteGroup) error {
	if err := ValidateMemberCalculationIntent(group); err != nil {
		return err
	}
	rowPolicy, err := json.Marshal(group.RowPolicy)
	if err != nil {
		return fmt.Errorf("encode write-group row policy: %w", err)
	}
	writePolicy, err := json.Marshal(group.WritePolicy)
	if err != nil {
		return fmt.Errorf("encode write-group write policy: %w", err)
	}
	migration, err := json.Marshal(group.Migration)
	if err != nil {
		return fmt.Errorf("encode write-group migration: %w", err)
	}
	profile := profileColumnValues(group)
	result, err := tx.ExecContext(ctx, r.query(`
		UPDATE write_groups SET
			revision = $1, applied_revision = $2, name = $3, status = $4,
			destination_connector_id = $5, destination_connector_revision = $6,
			destination_database = $7, destination_table_schema = $8, destination_table_name = $9,
			destination_storage_strategy = $10, destination_schema_revision = $11,
			destination_schema_digest = $12, row_policy = $13, write_policy = $14,
			migration = $15, updated_at = $16,
			layout_version = $17, profile_kind = $18, profile_interval = $19, profile_unit = $20, storage_mode = $21,
			sample_reduction = $22, result_shape = $23
		WHERE workspace_id = $24 AND id = $25
	`), group.Revision, group.AppliedRevision, group.Name, group.Status,
		group.Destination.ConnectorID, group.Destination.ConnectorRevision,
		group.Destination.Database, group.Destination.TableSchema, group.Destination.TableName,
		group.Destination.StorageStrategy, group.Destination.SchemaRevision, group.Destination.SchemaDigest,
		string(rowPolicy), string(writePolicy), string(migration), group.UpdatedAt,
		profile.layout, profile.kind, profile.interval, profile.unit, profile.mode, profile.reduction, profile.shape,
		group.WorkspaceID, group.ID)
	if err != nil {
		return fmt.Errorf("update write group: %w", err)
	}
	if affected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("confirm write-group update: %w", err)
	} else if affected != 1 {
		return fmt.Errorf("update write group: %w", ErrWriteGroupNotFound)
	}
	if _, err := tx.ExecContext(ctx, r.query(`DELETE FROM write_group_members WHERE group_id = $1`), group.ID); err != nil {
		return fmt.Errorf("replace write-group members: %w", err)
	}
	for index, member := range group.Members {
		var measurementID any
		if member.MeasurementID != nil {
			measurementID = *member.MeasurementID
		}
		var exactRecordingSource any
		if member.ExactRecordingSource != nil {
			payload, err := json.Marshal(member.ExactRecordingSource)
			if err != nil {
				return fmt.Errorf("encode exact recording source member %d: %w", index, err)
			}
			exactRecordingSource = string(payload)
		}
		if _, err := tx.ExecContext(ctx, r.query(`
			INSERT INTO write_group_members (
				group_id, member_index, device_id, point_id, tag_id,
				entity_key, source_revision, mapping_revision, measurement_id, target_column,
				required, max_age_seconds, circuit_id, circuit_name,
				value_role, source_value_type, result_value_type, exact_recording_source, calculation
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		`), group.ID, index, member.DeviceID, member.PointID, member.TagID,
			member.EntityKey, member.SourceRevision, member.MappingRevision, measurementID, member.TargetColumn,
			member.Required, member.MaxAgeSeconds, nullableText(member.CircuitID), nullableText(member.CircuitName),
			nullableText(string(member.ValueRole)), nullableText(string(member.SourceValueType)), nullableText(string(member.ResultValueType)),
			exactRecordingSource, nullableText(string(member.Calculation))); err != nil {
			return fmt.Errorf("replace write-group member %d: %w", index, err)
		}
	}
	return nil
}
