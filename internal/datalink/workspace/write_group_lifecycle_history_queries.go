package workspace

import (
	"context"
	"database/sql"
	"sort"
	"strings"

	"go-gateway/internal/datalink/groupdelivery"
)

func (r *SQLWriteGroupRepository) lifecycleGenerationIDs(ctx context.Context, tx *sql.Tx, lineage string) ([]string, error) {
	rows, err := tx.QueryContext(ctx, r.query(`SELECT group_id FROM write_group_generations WHERE lineage_id=$1 ORDER BY generation,group_id`), lineage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *SQLWriteGroupRepository) lifecycleRevisionsInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup) ([]string, error) {
	// Include persisted delivery references as well as versions: a broken or
	// missing immutable owner must fail, not disappear from impact review.
	revisions, err := groupdelivery.LifecyclePendingRevisions(ctx, tx, group.WorkspaceID, group.ID)
	if err != nil {
		return nil, err
	}
	wanted := map[string]bool{}
	for _, revision := range revisions {
		wanted[revision] = true
	}
	if group.AppliedRevision != "" {
		wanted[group.AppliedRevision] = true
	}
	rows, err := tx.QueryContext(ctx, r.query(`SELECT group_revision FROM write_group_versions WHERE group_id=$1
 UNION SELECT group_revision FROM write_group_revision_ends WHERE group_id=$2`), group.ID, group.ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var revision string
		if err := rows.Scan(&revision); err != nil {
			_ = rows.Close()
			return nil, err
		}
		wanted[revision] = true
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	revisions = []string{}
	for revision := range wanted {
		if strings.TrimSpace(revision) == "" {
			return nil, ErrWriteGroupAppliedRevisionUnavailable
		}
		revisions = append(revisions, revision)
	}
	sort.Strings(revisions)
	return revisions, nil
}
