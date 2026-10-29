package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/common"
)

// ErrWriteGroupLineageCorrupt identifies incomplete or contradictory canonical
// identity metadata. It deliberately does not unwrap either NotFound sentinel.
var ErrWriteGroupLineageCorrupt = errors.New("write-group lineage is corrupt")

type writeGroupLineageMetadata struct {
	LineageID, WorkspaceID, GenesisGroupID, CurrentGroupID string
	ReservedGroupID, ReservationOperationID,
	ReservationIntentDigest sql.NullString
	Revision, Generation int64
	PredecessorGroupID   sql.NullString
}

func lineageCorrupt(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrWriteGroupLineageCorrupt, fmt.Sprintf(format, args...))
}

// initializeCanonicalLineageInTx runs after a fresh group insert and before its
// key insert. Reads and replays must never call this initializer.
func (r *SQLWriteGroupRepository) initializeCanonicalLineageInTx(ctx context.Context, tx *sql.Tx, workspaceID, genesisGroupID string) (string, error) {
	if tx == nil || strings.TrimSpace(workspaceID) == "" || strings.TrimSpace(genesisGroupID) == "" {
		return "", fmt.Errorf("initialize canonical lineage: %w", ErrWriteGroupValidation)
	}
	metadata, found, err := r.loadWriteGroupLineage(ctx, tx, workspaceID, genesisGroupID)
	if err != nil {
		return "", err
	}
	if found {
		if metadata.Generation != 1 || metadata.GenesisGroupID != genesisGroupID ||
			metadata.CurrentGroupID != genesisGroupID || metadata.Revision != 1 || metadata.ReservedGroupID.Valid {
			return "", lineageCorrupt("canonical alias requires an unchanged first generation")
		}
		return metadata.LineageID, nil
	}
	lineageID, err := common.NewUUID()
	if err != nil {
		return "", fmt.Errorf("generate canonical lineage identity: %w", err)
	}
	if _, err = tx.ExecContext(ctx, r.query(`INSERT INTO write_group_lineages
		(lineage_id,workspace_id,current_group_id) VALUES ($1,$2,$3)`), lineageID, workspaceID, genesisGroupID); err != nil {
		return "", fmt.Errorf("insert canonical lineage: %w", err)
	}
	if _, err = tx.ExecContext(ctx, r.query(`INSERT INTO write_group_generations
		(group_id,lineage_id,generation) VALUES ($1,$2,1)`), genesisGroupID, lineageID); err != nil {
		return "", fmt.Errorf("insert canonical generation: %w", err)
	}
	return lineageID, nil
}

// loadWriteGroupLineage retains the requested generation while resolving its
// immutable genesis anchor. Every query uses the caller's runner/snapshot.
func (r *SQLWriteGroupRepository) loadWriteGroupLineage(ctx context.Context, runner writeGroupSQLRunner, workspaceID, groupID string) (*writeGroupLineageMetadata, bool, error) {
	metadata := &writeGroupLineageMetadata{}
	err := runner.QueryRowContext(ctx, r.query(`SELECT lineage_id,generation,predecessor_group_id
		FROM write_group_generations WHERE group_id=$1`), groupID).
		Scan(&metadata.LineageID, &metadata.Generation, &metadata.PredecessorGroupID)
	if errors.Is(err, sql.ErrNoRows) {
		if err := r.validateUnkeyedLegacyGroup(ctx, runner, workspaceID, groupID); err != nil {
			return nil, false, err
		}
		return nil, false, nil
	}
	if err != nil {
		return nil, false, lineageCorrupt("read generation: %v", err)
	}
	err = runner.QueryRowContext(ctx, r.query(`SELECT workspace_id,current_group_id,reserved_group_id,
		reservation_operation_id,reservation_intent_digest,revision FROM write_group_lineages WHERE lineage_id=$1`), metadata.LineageID).
		Scan(&metadata.WorkspaceID, &metadata.CurrentGroupID, &metadata.ReservedGroupID,
			&metadata.ReservationOperationID, &metadata.ReservationIntentDigest, &metadata.Revision)
	if err != nil {
		return nil, false, lineageCorrupt("read lineage: %v", err)
	}
	if strings.TrimSpace(metadata.LineageID) == "" || metadata.WorkspaceID != workspaceID ||
		strings.TrimSpace(workspaceID) == "" || metadata.Revision < 1 || metadata.Generation < 1 {
		return nil, false, lineageCorrupt("invalid lineage scope or revision")
	}
	if err := r.validateLineageGenerations(ctx, runner, metadata); err != nil {
		return nil, false, err
	}
	if err := r.validateLineageAliases(ctx, runner, metadata); err != nil {
		return nil, false, err
	}
	return metadata, true, nil
}

func (r *SQLWriteGroupRepository) validateUnkeyedLegacyGroup(ctx context.Context, runner writeGroupSQLRunner, workspaceID, groupID string) error {
	var references int
	err := runner.QueryRowContext(ctx, r.query(`SELECT
		(SELECT COUNT(*) FROM write_group_basic_keys WHERE group_id=$1) +
		(SELECT COUNT(*) FROM write_group_profile_keys WHERE group_id=$2) +
		(SELECT COUNT(*) FROM write_group_lineages WHERE current_group_id=$3 OR reserved_group_id=$4) +
		(SELECT COUNT(*) FROM write_group_generations WHERE predecessor_group_id=$5)`),
		groupID, groupID, groupID, groupID, groupID).Scan(&references)
	if err != nil {
		return lineageCorrupt("inspect generation references: %v", err)
	}
	if references != 0 {
		return lineageCorrupt("referenced group has no generation")
	}
	var actualWorkspace string
	err = runner.QueryRowContext(ctx, r.query(`SELECT workspace_id FROM write_groups WHERE id=$1`), groupID).Scan(&actualWorkspace)
	if err != nil || actualWorkspace != workspaceID {
		return lineageCorrupt("unkeyed group is absent or outside workspace")
	}
	return nil
}

// findCanonicalGroup reads the real key before any metadata joins. A key that
// exists cannot degrade to NotFound and accidentally enter ordinary creation.
func (r *SQLWriteGroupRepository) findCanonicalGroup(ctx context.Context, workspaceID, keyQuery string, args ...any) (*WriteGroup, error) {
	if r == nil || r.db == nil {
		return nil, lineageCorrupt("canonical repository unavailable")
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin canonical read: %w", err)
	}
	defer tx.Rollback()
	var genesisGroupID string
	var lineageID sql.NullString
	err = tx.QueryRowContext(ctx, r.query(keyQuery), args...).Scan(&genesisGroupID, &lineageID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("canonical key absent: %w", ErrWriteGroupNotFound)
	}
	if err != nil {
		return nil, lineageCorrupt("read canonical key: %v", err)
	}
	if !lineageID.Valid || strings.TrimSpace(lineageID.String) == "" {
		return nil, lineageCorrupt("canonical key has no lineage reference")
	}
	metadata, found, err := r.loadWriteGroupLineage(ctx, tx, workspaceID, genesisGroupID)
	if err != nil {
		return nil, err
	}
	if !found || metadata.GenesisGroupID != genesisGroupID || metadata.LineageID != lineageID.String {
		return nil, lineageCorrupt("canonical key does not name its genesis")
	}
	group, err := r.get(ctx, tx, workspaceID, metadata.CurrentGroupID)
	if err != nil {
		return nil, lineageCorrupt("read canonical head: %v", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("complete canonical read: %w", err)
	}
	return group, nil
}
