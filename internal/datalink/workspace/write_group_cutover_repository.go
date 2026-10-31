package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const cutoverTimeLayout = time.RFC3339Nano

type cutoverDetail struct {
	Ended   []CutoverEnd       `json:"ended"`
	Started []CutoverStart     `json:"started"`
	Moved   []CutoverMemberRef `json:"moved"`

	OwnershipFenceRevision string `json:"ownership_fence_revision,omitempty"`
}

func (r *SQLWriteGroupRepository) insertCutoverInTx(ctx context.Context, tx *sql.Tx, c *WriteGroupCutover) error {
	detail, err := json.Marshal(cutoverDetail{Ended: c.Ended, Started: c.Started, Moved: c.Moved, OwnershipFenceRevision: c.OwnershipFenceRevision})
	if err != nil {
		return fmt.Errorf("encode cutover detail: %w", err)
	}
	if _, err := tx.ExecContext(ctx, r.query(`
		INSERT INTO write_group_cutovers (id, workspace_id, kind, source_group_id, boundary, prior_cutover_id, review_digest, detail, created_at, intent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`), c.ID, c.WorkspaceID, c.Kind, c.SourceGroupID, c.Boundary.UTC().Format(cutoverTimeLayout), c.PriorCutoverID,
		c.ReviewDigest, string(detail), c.CreatedAt.UTC(), c.intent); err != nil {
		return fmt.Errorf("record cutover: %w", err)
	}
	return nil
}

// insertRevisionEndInTx fixes the instant a revision stops owning its sources.
// A revision already ended cannot be ended again: the first end is immutable.
func (r *SQLWriteGroupRepository) insertRevisionEndInTx(ctx context.Context, tx *sql.Tx, workspaceID, groupID, revision, cutoverID string, ends time.Time) error {
	result, err := tx.ExecContext(ctx, r.query(`
		INSERT INTO write_group_revision_ends (workspace_id, group_id, group_revision, ends_at, cutover_id)
		VALUES ($1, $2, $3, $4, $5) ON CONFLICT (group_id, group_revision) DO NOTHING
	`), workspaceID, groupID, revision, ends.UTC().Format(cutoverTimeLayout), cutoverID)
	if err != nil {
		return fmt.Errorf("end write-group revision: %w", err)
	}
	if affected, err := result.RowsAffected(); err != nil {
		return fmt.Errorf("confirm write-group revision end: %w", err)
	} else if affected != 1 {
		return writeGroupRevisionConflict("write group revision already ended")
	}
	return nil
}

// revisionEndInTx returns the immutable end of a revision, if it has one.
func (r *SQLWriteGroupRepository) revisionEndInTx(ctx context.Context, runner writeGroupSQLRunner, groupID, revision string) (time.Time, error) {
	var text string
	err := runner.QueryRowContext(ctx, r.query(`SELECT ends_at FROM write_group_revision_ends WHERE group_id = $1 AND group_revision = $2`),
		groupID, revision).Scan(&text)
	if errors.Is(err, sql.ErrNoRows) {
		return time.Time{}, nil
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("read write-group revision end: %w", err)
	}
	ends, err := time.Parse(cutoverTimeLayout, text)
	if err != nil {
		return time.Time{}, fmt.Errorf("decode write-group revision end: %w", err)
	}
	return ends.UTC(), nil
}

const cutoverColumns = `id, workspace_id, kind, source_group_id, boundary, prior_cutover_id, review_digest, detail, created_at, intent`

type cutoverScanner interface{ Scan(dest ...any) error }

func scanCutover(row cutoverScanner) (*WriteGroupCutover, error) {
	var c WriteGroupCutover
	var boundary, detail string
	if err := row.Scan(&c.ID, &c.WorkspaceID, &c.Kind, &c.SourceGroupID, &boundary, &c.PriorCutoverID, &c.ReviewDigest, &detail, &c.CreatedAt, &c.intent); err != nil {
		return nil, err
	}
	var err error
	if c.Boundary, err = time.Parse(cutoverTimeLayout, boundary); err != nil {
		return nil, fmt.Errorf("decode cutover boundary: %w", err)
	}
	c.Boundary, c.CreatedAt = c.Boundary.UTC(), c.CreatedAt.UTC()
	var decoded cutoverDetail
	if err := json.Unmarshal([]byte(detail), &decoded); err != nil {
		return nil, fmt.Errorf("decode cutover detail: %w", err)
	}
	c.Ended, c.Started, c.Moved = decoded.Ended, decoded.Started, decoded.Moved
	c.OwnershipFenceRevision = decoded.OwnershipFenceRevision
	return &c, nil
}

func (r *SQLWriteGroupRepository) cutoverInTx(ctx context.Context, runner writeGroupSQLRunner, workspaceID, id string) (*WriteGroupCutover, error) {
	c, err := scanCutover(runner.QueryRowContext(ctx, r.query(`SELECT `+cutoverColumns+` FROM write_group_cutovers WHERE workspace_id = $1 AND id = $2`), workspaceID, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, writeGroupNotFound("read cutover")
	}
	if err != nil {
		return nil, fmt.Errorf("read cutover: %w", err)
	}
	return c, nil
}

// cutoverByDigestInTx finds the cutover a review digest was applied as, with the
// intent it recorded. A nil cutover without error means none was applied.
func (r *SQLWriteGroupRepository) cutoverByDigestInTx(ctx context.Context, runner writeGroupSQLRunner, workspaceID, digest string) (*WriteGroupCutover, string, error) {
	c, err := scanCutover(runner.QueryRowContext(ctx, r.query(`SELECT `+cutoverColumns+` FROM write_group_cutovers WHERE workspace_id = $1 AND review_digest = $2`), workspaceID, digest))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, "", nil
	}
	if err != nil {
		return nil, "", fmt.Errorf("read cutover by review: %w", err)
	}
	return c, c.intent, nil
}

func (r *SQLWriteGroupRepository) listCutoversInTx(ctx context.Context, runner writeGroupSQLRunner, workspaceID, sourceGroupID string) ([]*WriteGroupCutover, error) {
	rows, err := runner.QueryContext(ctx, r.query(`SELECT `+cutoverColumns+` FROM write_group_cutovers
		WHERE workspace_id = $1 AND source_group_id = $2 ORDER BY boundary, id`), workspaceID, sourceGroupID)
	if err != nil {
		return nil, fmt.Errorf("list cutovers: %w", err)
	}
	defer rows.Close()
	out := []*WriteGroupCutover{}
	for rows.Next() {
		c, err := scanCutover(rows)
		if err != nil {
			return nil, fmt.Errorf("scan cutover: %w", err)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// rolledBackInTx reports whether a forward cutover already has its rollback.
func (r *SQLWriteGroupRepository) rolledBackInTx(ctx context.Context, runner writeGroupSQLRunner, id string) (bool, error) {
	var n int
	if err := runner.QueryRowContext(ctx, r.query(`SELECT COUNT(*) FROM write_group_cutovers WHERE prior_cutover_id = $1`), id).Scan(&n); err != nil {
		return false, fmt.Errorf("read cutover rollback: %w", err)
	}
	return n > 0, nil
}

// activeOverlapInTx lists the sources of candidate that another applied,
// not-ended wide or custom group of the workspace still owns.
func (r *SQLWriteGroupRepository) activeOverlapInTx(ctx context.Context, runner writeGroupSQLRunner, candidate *WriteGroup) ([]CutoverMemberRef, error) {
	candidateRecord, err := r.mixedRecordIdentityInTx(ctx, runner, candidate)
	if err != nil {
		return nil, err
	}
	if candidateRecord != nil {
		// One scan, not one scan per active owner. Independent records never gain
		// permission to write another record's wide table through source reuse.
		if err := r.checkMixedTargetOwnership(ctx, runner, candidate); err != nil {
			return nil, err
		}
	}
	rows, err := runner.QueryContext(ctx, r.query(`
		SELECT id, applied_revision FROM write_groups
		WHERE workspace_id = $1 AND id <> $2 AND applied_revision <> '' AND status IN ('draft', 'ready', 'running')
	`), candidate.WorkspaceID, candidate.ID)
	if err != nil {
		return nil, fmt.Errorf("list active write groups: %w", err)
	}
	type applied struct{ id, revision string }
	var active []applied
	for rows.Next() {
		var a applied
		if err := rows.Scan(&a.id, &a.revision); err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("scan active write group: %w", err)
		}
		active = append(active, a)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("list active write groups: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("list active write groups: %w", err)
	}
	wanted := make(map[CutoverMemberRef]struct{}, len(candidate.Members))
	for _, m := range candidate.Members {
		wanted[memberRefOf(m)] = struct{}{}
	}
	var overlap []CutoverMemberRef
	for _, a := range active {
		ended, err := r.revisionEndInTx(ctx, runner, a.id, a.revision)
		if err != nil {
			return nil, err
		}
		if !ended.IsZero() {
			continue
		}
		version, err := r.versionByRevisionInTx(ctx, runner, candidate.WorkspaceID, a.id, a.revision)
		if errors.Is(err, ErrWriteGroupAppliedRevisionUnavailable) {
			if candidateRecord != nil {
				// An unknown active writer cannot prove an independent record.
				// Keep the historical legacy-candidate behavior unchanged.
				return nil, err
			}
			continue
		}
		if err != nil {
			return nil, err
		}
		if version.Group.isLongLayout() {
			continue
		}
		if candidateRecord != nil && isMixedWriteGroup(version.Group) {
			if version.Group.WorkspaceID != candidate.WorkspaceID || version.Group.ID != a.id || version.Group.Revision != a.revision || version.Group.AppliedRevision != a.revision {
				return nil, lineageCorrupt("active mixed source reference has a foreign applied identity")
			}
			ownerRecord, err := r.mixedRecordIdentityInTx(ctx, runner, version.Group)
			if err != nil {
				return nil, err
			}
			if independentMixedRecordIdentities(candidateRecord, ownerRecord) {
				continue
			}
		}
		if candidate.isWideProfile() && version.Group.isWideProfile() && profileKindOf(version.Group) != profileKindOf(candidate) {
			// The instantaneous and the cumulative profile of one source are two groups of different kinds,
			// not two writers; a group of the same kind, like a legacy wide or custom group, still owns it.
			continue
		}
		for _, m := range version.Group.Members {
			if _, ok := wanted[memberRefOf(m)]; ok {
				overlap = append(overlap, memberRefOf(m))
			}
		}
	}
	sortMemberRefs(overlap)
	return overlap, nil
}
