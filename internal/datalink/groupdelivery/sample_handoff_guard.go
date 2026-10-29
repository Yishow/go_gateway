package groupdelivery

import (
	"context"
	"database/sql"
	"encoding/json"
	"sort"
	"time"
)

// handoffView is computed only after receipt verification. State compaction by
// normal reclaim does not change ownership and therefore does not change it.
func handoffView(ctx context.Context, q ownershipQuery, key GroupKey) (HandoffView, []handoffReceipt, error) {
	receipts, err := handoffReceiptsForKey(ctx, q, key)
	if err != nil {
		return HandoffView{}, nil, err
	}
	if len(receipts) == 0 {
		return HandoffView{}, nil, nil
	}
	var parts []string
	for _, r := range receipts {
		if err = verifyHandoffReceipt(ctx, q, r, map[string]bool{}); err != nil {
			return HandoffView{}, nil, err
		}
		if r.State == "finalized" {
			r.State = "accepted"
		}
		data, e := json.Marshal(r)
		if e != nil {
			return HandoffView{}, nil, ErrOwnershipUnavailable
		}
		parts = append(parts, string(data))
	}
	sort.Strings(parts)
	raw, _ := json.Marshal(parts)
	return HandoffView{key: key, digest: handoffHash(string(raw))}, receipts, nil
}

func validateHandoffClosure(ctx context.Context, tx *sql.Tx, key GroupKey, closure Closure) error {
	if _, err := resolveCutoverOwnership(ctx, tx, OwnershipSource{WorkspaceID: key.WorkspaceID}, closure.NextClose, key); err != nil {
		return err
	}
	view, receipts, err := handoffView(ctx, tx, key)
	if err != nil {
		return err
	}
	if view != closure.HandoffView {
		return ErrOwnershipConflict
	}
	for _, r := range receipts {
		if r.State != "pending" && r.State != "quarantined" {
			continue
		}
		when := r.ObservedAt
		if r.Source == key {
			when = r.SourceBucket
		}
		at, e := time.Parse(time.RFC3339Nano, when)
		if e != nil {
			return ErrOwnershipUnavailable
		}
		if at.Before(closure.NextClose) {
			return ErrOwnershipConflict
		}
	}
	// Discover only this revision and its durable predecessor scopes. This
	// closes the gap before a transfer receipt is prepared without sweeping
	// unrelated workspace samples or relying on a cached boundary.
	rows, err := tx.QueryContext(ctx, `WITH RECURSIVE ancestors(group_id,revision) AS (
 SELECT ?,? UNION SELECT e.group_id,e.group_revision
 FROM ancestors a JOIN write_group_ownership_fences t ON t.workspace_id=? AND t.group_id=a.group_id AND t.group_revision=a.revision AND t.direction='started'
 JOIN write_group_ownership_fences e ON e.workspace_id=t.workspace_id AND e.fence_revision=t.fence_revision AND e.direction='ended'
 AND e.device_id=t.device_id AND e.point_id=t.point_id AND e.tag_id=t.tag_id)
 SELECT s.seq FROM wg_delivery_samples s JOIN ancestors a ON a.group_id=s.group_id AND a.revision=s.group_revision WHERE s.workspace_id=? AND s.consumed=0`, key.GroupID, key.GroupRevision, key.WorkspaceID, key.WorkspaceID)
	if err != nil {
		return ErrOwnershipUnavailable
	}
	var seqs []int64
	for rows.Next() {
		var seq int64
		if rows.Scan(&seq) != nil {
			rows.Close()
			return ErrOwnershipUnavailable
		}
		seqs = append(seqs, seq)
	}
	if rows.Err() != nil {
		rows.Close()
		return ErrOwnershipUnavailable
	}
	if rows.Close() != nil {
		return ErrOwnershipUnavailable
	}
	for _, seq := range seqs {
		row, e := readHandoffSample(ctx, tx, seq)
		if e != nil {
			return e
		}
		outgoing, e := outgoingHandoffs(ctx, tx, seq)
		if e != nil {
			return e
		}
		settled := len(outgoing) > 0
		for _, receipt := range outgoing {
			if e = verifyHandoffReceipt(ctx, tx, receipt, map[string]bool{}); e != nil {
				return e
			}
			if receipt.State == "pending" || receipt.State == "quarantined" {
				settled = false
			}
		}
		// Pending ancestry still owes the current descendant its original.
		// Only accepted copies permit this row to defer to downstream rows.
		if settled {
			continue
		}
		sample, e := decodeVerifiedSample(row.Payload, row.Digest)
		if e != nil {
			return ErrHandoffProofConflict
		}
		own, e := resolveCutoverOwnership(ctx, tx, sourceOf(row.Key, sample), sample.ObservedAt, GroupKey{})
		if e != nil {
			return e
		}
		if len(outgoing) > 0 {
			obligated := false
			for _, receipt := range outgoing {
				if receipt.State != "pending" && receipt.State != "quarantined" {
					continue
				}
				if key == receipt.Target || handoffDescendant(own, key, receipt.Target, map[GroupKey]bool{}) {
					obligated = true
				}
			}
			if !obligated {
				continue // An independent accepted fanout may still close.
			}
		}
		if !own.Affected[row.Key] || own.Owners[row.Key] {
			continue
		}
		if row.Key == key {
			bucket, e := time.Parse(time.RFC3339Nano, row.Bucket)
			if e != nil {
				return ErrOwnershipUnavailable
			}
			if bucket.Before(closure.NextClose) {
				return ErrOwnershipConflict
			}
		} else if own.Owners[key] && handoffDescendant(own, key, row.Key, map[GroupKey]bool{}) && sample.ObservedAt.Before(closure.NextClose) {
			return ErrOwnershipConflict
		}
	}
	return nil
}

// validateHandoffScope prevents metadata readers from silently treating a
// malformed audit row as a successful transfer.
func validateHandoffScope(ctx context.Context, q ownershipQuery, groupID, revision string) error {
	if db, ok := q.(*sql.DB); ok {
		// Pin the entire verification walk: nested ownership reads must not
		// wait for another pool connection or observe a different snapshot.
		// Existing transactional callers retain their own transaction boundary.
		tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
		if err != nil {
			return &handoffStorageError{cause: err}
		}
		defer tx.Rollback()
		observed := &handoffReadObserver{ownershipQuery: tx}
		err = validateHandoffScope(ctx, observed, groupID, revision)
		if observed.failure != nil {
			return &handoffStorageError{cause: observed.failure}
		}
		if err != nil {
			return err
		}
		if err = tx.Commit(); err != nil {
			return &handoffStorageError{cause: err}
		}
		return nil
	}
	query := `SELECT DISTINCT workspace_id,source_group_id,source_revision FROM wg_delivery_sample_handoffs WHERE (source_group_id=? OR target_group_id=?)`
	args := []any{groupID, groupID}
	if revision != "" {
		query += ` AND (source_revision=? OR target_revision=?)`
		args = append(args, revision, revision)
	}
	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return &handoffStorageError{cause: err}
	}
	var keys []GroupKey
	for rows.Next() {
		var key GroupKey
		if err := rows.Scan(&key.WorkspaceID, &key.GroupID, &key.GroupRevision); err != nil {
			rows.Close()
			return &handoffStorageError{cause: err}
		}
		keys = append(keys, key)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return &handoffStorageError{cause: err}
	}
	if err := rows.Close(); err != nil {
		return &handoffStorageError{cause: err}
	}
	for _, key := range keys {
		if _, _, err := handoffView(ctx, q, key); err != nil {
			return err
		}
	}
	return nil
}

func validateAllHandoffReceipts(ctx context.Context, q ownershipQuery) error {
	rs, err := queryHandoffs(ctx, q, "1=1")
	if err != nil {
		return err
	}
	for _, r := range rs {
		if err = verifyHandoffReceipt(ctx, q, r, map[string]bool{}); err != nil {
			return err
		}
	}
	return nil
}

// Discovery may retain a damaged scope, but must not swallow actual SQL read
// failures as mere damaged metadata. No raw cause is exposed by Error().
type handoffReadObserver struct {
	ownershipQuery
	failure error
}

func (q *handoffReadObserver) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	rows, err := q.ownershipQuery.QueryContext(ctx, query, args...)
	if err != nil && q.failure == nil {
		q.failure = err
	}
	return rows, err
}
func (q *handoffReadObserver) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	row := q.ownershipQuery.QueryRowContext(ctx, query, args...)
	if err := row.Err(); err != nil && q.failure == nil {
		q.failure = err
	}
	return row
}
