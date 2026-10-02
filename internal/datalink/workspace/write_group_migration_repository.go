package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type writeGroupMigrationMap struct {
	WorkspaceID    string
	SourceKind     string
	SourceID       string
	SourceRevision string
	GroupID        string
	BeforeIntent   string
	ReviewDigest   string
	AdapterVersion string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (r *SQLWriteGroupRepository) getMigrationMapInTx(
	ctx context.Context,
	tx *sql.Tx,
	workspaceID string,
	sourceKind string,
	sourceID string,
) (*writeGroupMigrationMap, bool, error) {
	if tx == nil {
		return nil, false, fmt.Errorf("read write-group migration map: %w", ErrWriteGroupValidation)
	}
	var migrationMap writeGroupMigrationMap
	err := tx.QueryRowContext(ctx, r.query(`
		SELECT workspace_id, source_kind, source_id, source_revision,
			group_id, before_intent, review_digest, adapter_version,
			created_at, updated_at
		FROM write_group_migration_maps
		WHERE workspace_id = $1 AND source_kind = $2 AND source_id = $3
	`), workspaceID, sourceKind, sourceID).Scan(
		&migrationMap.WorkspaceID,
		&migrationMap.SourceKind,
		&migrationMap.SourceID,
		&migrationMap.SourceRevision,
		&migrationMap.GroupID,
		&migrationMap.BeforeIntent,
		&migrationMap.ReviewDigest,
		&migrationMap.AdapterVersion,
		&migrationMap.CreatedAt,
		&migrationMap.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("read write-group migration map: %w", err)
	}
	return &migrationMap, true, nil
}

func (r *SQLWriteGroupRepository) insertMigrationMapInTx(
	ctx context.Context,
	tx *sql.Tx,
	migrationMap writeGroupMigrationMap,
) error {
	if tx == nil {
		return fmt.Errorf("insert write-group migration map: %w", ErrWriteGroupValidation)
	}
	if _, err := tx.ExecContext(ctx, r.query(`
		INSERT INTO write_group_migration_maps (
			workspace_id, source_kind, source_id, source_revision,
			group_id, before_intent, review_digest, adapter_version,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`), migrationMap.WorkspaceID, migrationMap.SourceKind, migrationMap.SourceID,
		migrationMap.SourceRevision, migrationMap.GroupID, migrationMap.BeforeIntent,
		migrationMap.ReviewDigest, migrationMap.AdapterVersion,
		migrationMap.CreatedAt, migrationMap.UpdatedAt); err != nil {
		return fmt.Errorf("insert write-group migration map: %w", err)
	}
	return nil
}
