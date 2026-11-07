package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

// WriteGroupDeleteReference identifies another recording that still needs proof.
type WriteGroupDeleteReference struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// WriteGroupDeleteReferences describes surviving recordings requiring this evidence.
type WriteGroupDeleteReferences struct {
	ReferencedBy []WriteGroupDeleteReference `json:"referenced_by"`
	Total        int                         `json:"total"`
}

func (s *WriteGroupService) permanentDeleteGroupIDs(ctx context.Context, tx *sql.Tx, group *WriteGroup) (map[string]bool, error) {
	lineage, _, err := s.repo.loadWriteGroupLineage(ctx, tx, group.WorkspaceID, group.ID)
	if err != nil {
		return nil, err
	}
	ids, _, err := s.permanentDeleteIDsForLineage(ctx, tx, group.ID, lineage)
	return ids, err
}

func (s *WriteGroupService) checkPermanentDeleteSharedProof(ctx context.Context, tx *sql.Tx, group *WriteGroup, ids map[string]bool) error {
	type cutoverProof struct {
		id, boundary, source, detail string
	}
	rows, err := tx.QueryContext(ctx, s.repo.query(`SELECT id,boundary,source_group_id,detail FROM write_group_cutovers WHERE workspace_id=$1`), group.WorkspaceID)
	if err != nil {
		return fmt.Errorf("read recording cutovers: %w", err)
	}
	var cutovers []cutoverProof
	for rows.Next() {
		var proof cutoverProof
		if err := rows.Scan(&proof.id, &proof.boundary, &proof.source, &proof.detail); err != nil {
			rows.Close()
			return fmt.Errorf("scan recording cutover: %w", err)
		}
		cutovers = append(cutovers, proof)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("iterate recording cutovers: %w", err)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	rows, err = tx.QueryContext(ctx, s.repo.query(`SELECT f.fence_revision,f.group_id FROM write_group_ownership_fences f
		JOIN write_group_cutovers c ON c.id=f.fence_revision WHERE c.workspace_id=$1
		UNION SELECT e.cutover_id,e.group_id FROM write_group_revision_ends e
		JOIN write_group_cutovers c ON c.id=e.cutover_id WHERE c.workspace_id=$2`), group.WorkspaceID, group.WorkspaceID)
	if err != nil {
		return fmt.Errorf("read cutover proof owners: %w", err)
	}
	cutoverOwners := make(map[string]map[string]bool)
	for rows.Next() {
		var cutoverID, ownerID string
		if err := rows.Scan(&cutoverID, &ownerID); err != nil {
			rows.Close()
			return fmt.Errorf("scan cutover proof owner: %w", err)
		}
		if cutoverOwners[cutoverID] == nil {
			cutoverOwners[cutoverID] = make(map[string]bool)
		}
		cutoverOwners[cutoverID][ownerID] = true
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("iterate cutover proof owners: %w", err)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	refs := make(map[string]bool)
	for _, cutover := range cutovers {
		var detail cutoverDetail
		if json.Unmarshal([]byte(cutover.detail), &detail) != nil {
			return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_OPERATION_PENDING"}
		}
		owners := map[string]bool{cutover.source: true}
		for _, end := range detail.Ended {
			owners[end.GroupID] = true
		}
		for _, start := range detail.Started {
			owners[start.GroupID] = true
		}
		for id := range cutoverOwners[cutover.id] {
			owners[id] = true
		}
		relevant := false
		for id := range owners {
			if ids[id] {
				relevant = true
			}
		}
		if !relevant {
			continue
		}
		boundary, err := time.Parse(time.RFC3339Nano, cutover.boundary)
		if err != nil || boundary.After(s.clockNow()) {
			return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_OPERATION_PENDING"}
		}
		for id := range owners {
			if !ids[id] {
				refs[id] = true
			}
		}
	}
	for id := range ids {
		rows, err := tx.QueryContext(ctx, s.repo.query(`SELECT source_group_id,target_group_id,state FROM wg_delivery_sample_handoffs
			WHERE workspace_id=$1 AND (source_group_id=$2 OR target_group_id=$3)`), group.WorkspaceID, id, id)
		if err != nil {
			return fmt.Errorf("read recording handoffs: %w", err)
		}
		for rows.Next() {
			var source, target, state string
			if err := rows.Scan(&source, &target, &state); err != nil {
				rows.Close()
				return fmt.Errorf("scan recording handoff: %w", err)
			}
			if state != "finalized" {
				rows.Close()
				return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_OPERATION_PENDING"}
			}
			for _, owner := range []string{source, target} {
				if !ids[owner] {
					refs[owner] = true
				}
			}
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return fmt.Errorf("iterate recording handoffs: %w", err)
		}
		if err := rows.Close(); err != nil {
			return err
		}
	}
	if len(refs) == 0 {
		return nil
	}
	references := make([]WriteGroupDeleteReference, 0, len(refs))
	for id := range refs {
		var ref WriteGroupDeleteReference
		err := tx.QueryRowContext(ctx, s.repo.query(`SELECT id,name FROM write_groups WHERE workspace_id=$1 AND id=$2`), group.WorkspaceID, id).Scan(&ref.ID, &ref.Name)
		if errors.Is(err, sql.ErrNoRows) {
			return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_OPERATION_PENDING"}
		}
		if err != nil {
			return fmt.Errorf("read referenced recording: %w", err)
		}
		references = append(references, ref)
	}
	slices.SortFunc(references, func(a, b WriteGroupDeleteReference) int { return strings.Compare(a.ID, b.ID) })
	total := len(references)
	if len(references) > 20 {
		references = references[:20]
	}
	return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_HAS_SUCCESSOR", Details: WriteGroupDeleteReferences{ReferencedBy: references, Total: total}}
}
