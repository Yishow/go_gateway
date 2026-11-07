package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// WriteGroupDeleteError describes an operator-actionable refusal to remove a recording.
type WriteGroupDeleteError struct {
	Code    string
	Details any
}

// Error returns the stable refusal code without exposing private diagnostics.
func (e *WriteGroupDeleteError) Error() string { return e.Code }

// Unwrap preserves the existing lifecycle-blocked classification.
func (e *WriteGroupDeleteError) Unwrap() error { return ErrWriteGroupLifecycleBlocked }

// WriteGroupDeleteBacklog separates accepted samples from each outbox state.
type WriteGroupDeleteBacklog struct {
	SamplesUnconsumed int            `json:"samples_unconsumed"`
	Outbox            map[string]int `json:"outbox"`
}

// WriteGroupDeletionBacklogGuard counts pending deletion evidence in the caller's transaction.
type WriteGroupDeletionBacklogGuard interface {
	CountWriteGroupDeleteBacklog(context.Context, *sql.Tx, *WriteGroup, *SQLWriteGroupRepository) (WriteGroupDeleteBacklog, error)
}

// CountWriteGroupDeleteBacklog returns the transaction-scoped pending evidence.
func (r *SQLWriteGroupRepository) CountWriteGroupDeleteBacklog(ctx context.Context, tx *sql.Tx, groupID string) (WriteGroupDeleteBacklog, error) {
	result := WriteGroupDeleteBacklog{Outbox: make(map[string]int)}
	if r == nil || tx == nil || groupID == "" {
		return result, ErrWriteGroupServiceUnavailable
	}
	err := tx.QueryRowContext(ctx, r.query(`SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id=$1 AND consumed=0`), groupID).Scan(&result.SamplesUnconsumed)
	if err != nil {
		return result, fmt.Errorf("count unconsumed recording samples: %w", err)
	}
	rows, err := tx.QueryContext(ctx, r.query(`SELECT state,COUNT(*) FROM wg_delivery_outbox WHERE group_id=$1 AND state NOT IN ('sql_committed','operator_skipped') GROUP BY state`), groupID)
	if err != nil {
		return result, fmt.Errorf("count pending recording deliveries: %w", err)
	}
	for rows.Next() {
		var state string
		var count int
		if err := rows.Scan(&state, &count); err != nil {
			rows.Close()
			return result, fmt.Errorf("read pending recording deliveries: %w", err)
		}
		result.Outbox[state] = count
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return result, fmt.Errorf("iterate pending recording deliveries: %w", err)
	}
	if err := rows.Close(); err != nil {
		return result, err
	}
	rows, err = tx.QueryContext(ctx, r.query(`SELECT status,COUNT(*) FROM gw_delivery_outbox WHERE destination_id=$1 AND status<>'delivered' GROUP BY status`), groupID)
	if err != nil {
		return result, fmt.Errorf("count pending legacy recording deliveries: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var state string
		var count int
		if err := rows.Scan(&state, &count); err != nil {
			return result, fmt.Errorf("read pending legacy recording deliveries: %w", err)
		}
		result.Outbox[state] += count
	}
	if err := rows.Err(); err != nil {
		return result, fmt.Errorf("iterate pending legacy recording deliveries: %w", err)
	}
	return result, nil
}

func (b WriteGroupDeleteBacklog) pending() bool {
	return b.SamplesUnconsumed > 0 || len(b.Outbox) > 0
}

func (s *WriteGroupService) checkPermanentDeleteIntakeAndBacklog(ctx context.Context, tx *sql.Tx, groups []*WriteGroup) error {
	combined := WriteGroupDeleteBacklog{Outbox: make(map[string]int)}
	for _, generation := range groups {
		id := generation.ID
		if generation.Status != WriteGroupStatusDisabled && generation.Status != WriteGroupStatusDeleted {
			snapshot, resolveErr := s.repo.resolveVersionInTx(ctx, tx, generation.WorkspaceID, id, s.clockNow())
			if resolveErr != nil && !errors.Is(resolveErr, ErrWriteGroupAppliedRevisionUnavailable) {
				return resolveErr
			}
			if errors.Is(resolveErr, ErrWriteGroupAppliedRevisionUnavailable) && generation.AppliedRevision != "" {
				if _, err := s.repo.VersionByRevisionInTx(ctx, tx, generation.WorkspaceID, id, generation.AppliedRevision); err != nil {
					return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_OPERATION_PENDING"}
				}
			}
			if resolveErr == nil {
				end, err := s.repo.revisionEndInTx(ctx, tx, id, snapshot.AppliedRevision)
				if err != nil {
					return err
				}
				if end.IsZero() || s.clockNow().Before(end) {
					return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_ACTIVE"}
				}
			}
		}
		var backlog WriteGroupDeleteBacklog
		var err error
		if guard, ok := s.backlogGuard.(WriteGroupDeletionBacklogGuard); ok {
			backlog, err = guard.CountWriteGroupDeleteBacklog(ctx, tx, generation, s.repo)
		} else {
			backlog, err = s.repo.CountWriteGroupDeleteBacklog(ctx, tx, id)
		}
		if err != nil {
			return err
		}
		combined.SamplesUnconsumed += backlog.SamplesUnconsumed
		for state, count := range backlog.Outbox {
			combined.Outbox[state] += count
		}
	}
	if combined.pending() {
		return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_BACKLOG", Details: combined}
	}
	return nil
}

func (s *WriteGroupService) checkPermanentDeleteLineage(ctx context.Context, tx *sql.Tx, lineage *writeGroupLineageMetadata) error {
	if lineage == nil {
		return nil
	}
	if lineage.ReservedGroupID.Valid {
		return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_OPERATION_PENDING"}
	}
	var later int
	err := tx.QueryRowContext(ctx, s.repo.query(`SELECT COUNT(*) FROM write_group_generations g
		JOIN write_groups w ON w.id=g.group_id
		WHERE g.lineage_id=$1 AND g.generation>$2 AND w.status<>'deleted'`), lineage.LineageID, lineage.Generation).Scan(&later)
	if err != nil {
		return fmt.Errorf("check later recording generations: %w", err)
	}
	if later > 0 {
		return &WriteGroupDeleteError{Code: "WRITE_GROUP_DELETE_HAS_SUCCESSOR"}
	}
	return nil
}
