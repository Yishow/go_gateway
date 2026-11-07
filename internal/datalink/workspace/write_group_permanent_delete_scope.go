package workspace

import (
	"context"
	"database/sql"
	"fmt"
)

// permanentDeleteScope is read once from the same transaction as the guards.
// Groups are ordered newest first to respect generation foreign keys at purge.
type permanentDeleteScope struct {
	lineage *writeGroupLineageMetadata
	ids     map[string]bool
	groups  []*WriteGroup
}

func (s *WriteGroupService) loadPermanentDeleteScope(ctx context.Context, tx *sql.Tx, group *WriteGroup) (*permanentDeleteScope, error) {
	lineage, _, err := s.repo.loadWriteGroupLineage(ctx, tx, group.WorkspaceID, group.ID)
	if err != nil {
		return nil, err
	}
	if err := s.checkPermanentDeleteLineage(ctx, tx, lineage); err != nil {
		return nil, err
	}
	ids, ordered, err := s.permanentDeleteIDsForLineage(ctx, tx, group.ID, lineage)
	if err != nil {
		return nil, err
	}
	scope := &permanentDeleteScope{lineage: lineage, ids: ids, groups: make([]*WriteGroup, 0, len(ordered))}
	for _, id := range ordered {
		generation := group
		if id != group.ID {
			generation, err = s.repo.GetInTx(ctx, tx, group.WorkspaceID, id)
			if err != nil {
				return nil, err
			}
		}
		scope.groups = append(scope.groups, generation)
	}
	return scope, nil
}

func (s *WriteGroupService) permanentDeleteIDsForLineage(ctx context.Context, tx *sql.Tx, groupID string, lineage *writeGroupLineageMetadata) (map[string]bool, []string, error) {
	ids := map[string]bool{groupID: true}
	if lineage == nil {
		return ids, []string{groupID}, nil
	}
	rows, err := tx.QueryContext(ctx, s.repo.query(`SELECT group_id FROM write_group_generations WHERE lineage_id=$1 ORDER BY generation DESC`), lineage.LineageID)
	if err != nil {
		return nil, nil, fmt.Errorf("read recording generations: %w", err)
	}
	defer rows.Close()
	var ordered []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, nil, fmt.Errorf("scan recording generation: %w", err)
		}
		ids[id] = true
		ordered = append(ordered, id)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("iterate recording generations: %w", err)
	}
	return ids, ordered, nil
}
