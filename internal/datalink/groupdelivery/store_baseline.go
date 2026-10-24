package groupdelivery

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/aggregation"
)

// decodedBaseline is a validated baseline: the member it belongs to, the
// identity of the binding it was taken under and the window it covers.
type decodedBaseline struct {
	CounterBaseline
	bindingKey string
	windowEnd  string
	digest     string
}

// validateBaselines decodes every baseline strictly before anything is written.
// A payload that does not decode, that belongs to another member than the one
// it is filed under, that names another workspace, group or profile revision
// than the closure's own key, or that covers a window the closure has not
// closed, is an invalid closure: foreign or damaged state is refused when it is
// written, not discovered after a restart. Only the measurement's own
// revisions (definition, source binding, epoch, unit, policy) may differ from an
// earlier baseline; those are legitimate and start a new one.
func validateBaselines(key GroupKey, nextClose time.Time, baselines []CounterBaseline) ([]decodedBaseline, error) {
	seen := make(map[string]struct{}, len(baselines))
	decoded := make([]decodedBaseline, 0, len(baselines))
	for _, baseline := range baselines {
		if baseline.MemberKey == "" {
			return nil, fmt.Errorf("%w: a counter baseline needs its member", ErrInvalidClosure)
		}
		if _, duplicate := seen[baseline.MemberKey]; duplicate {
			return nil, fmt.Errorf("%w: more than one baseline for a member", ErrInvalidClosure)
		}
		seen[baseline.MemberKey] = struct{}{}
		checkpoint, err := aggregation.DecodeCounterCheckpoint(baseline.Payload)
		if err != nil {
			return nil, fmt.Errorf("%w: counter baseline is not a valid checkpoint", ErrInvalidClosure)
		}
		if checkpoint.Binding.MemberKey != baseline.MemberKey {
			return nil, fmt.Errorf("%w: counter baseline belongs to another member", ErrInvalidClosure)
		}
		if checkpoint.Binding.WorkspaceID != key.WorkspaceID || checkpoint.Binding.GroupID != key.GroupID ||
			checkpoint.Binding.ProfileRevision != key.GroupRevision {
			return nil, fmt.Errorf("%w: counter baseline belongs to another group scope", ErrInvalidClosure)
		}
		if checkpoint.WindowEnd.After(nextClose) {
			return nil, fmt.Errorf("%w: counter baseline covers a window that is not closed", ErrInvalidClosure)
		}
		bindingKey, err := checkpoint.Binding.Key()
		if err != nil {
			return nil, fmt.Errorf("%w: counter baseline binding is incomplete", ErrInvalidClosure)
		}
		sum := sha256.Sum256(baseline.Payload)
		decoded = append(decoded, decodedBaseline{
			CounterBaseline: baseline, bindingKey: bindingKey,
			windowEnd: timestamp(checkpoint.WindowEnd), digest: hex.EncodeToString(sum[:]),
		})
	}
	return decoded, nil
}

// storeBaselines writes each baseline inside the closure transaction. A baseline
// moves forward only: the same window with the same bytes is a repeat, an older
// window is a late replay that changes nothing, and the same window with
// different bytes is a conflict that fails the whole closure.
func storeBaselines(ctx context.Context, tx *sql.Tx, key GroupKey, baselines []decodedBaseline, now string) error {
	for _, baseline := range baselines {
		var windowEnd, digest string
		err := tx.QueryRowContext(ctx, `
			SELECT window_end, payload_digest FROM wg_delivery_counter_baselines
			WHERE group_id = ? AND group_revision = ? AND member_key = ?`,
			key.GroupID, key.GroupRevision, baseline.MemberKey).Scan(&windowEnd, &digest)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO wg_delivery_counter_baselines (
					workspace_id, group_id, group_revision, member_key, binding_key, window_end,
					payload, payload_digest, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				key.WorkspaceID, key.GroupID, key.GroupRevision, baseline.MemberKey, baseline.bindingKey,
				baseline.windowEnd, string(baseline.Payload), baseline.digest, now, now); err != nil {
				return fmt.Errorf("store counter baseline: %w", err)
			}
		case err != nil:
			return fmt.Errorf("check counter baseline: %w", err)
		case windowEnd == baseline.windowEnd && digest == baseline.digest, baseline.windowEnd < windowEnd:
			continue
		case windowEnd == baseline.windowEnd:
			return ErrBaselineConflict
		default:
			if _, err := tx.ExecContext(ctx, `
				UPDATE wg_delivery_counter_baselines
				SET binding_key = ?, window_end = ?, payload = ?, payload_digest = ?, updated_at = ?
				WHERE group_id = ? AND group_revision = ? AND member_key = ?`,
				baseline.bindingKey, baseline.windowEnd, string(baseline.Payload), baseline.digest, now,
				key.GroupID, key.GroupRevision, baseline.MemberKey); err != nil {
				return fmt.Errorf("advance counter baseline: %w", err)
			}
		}
	}
	return nil
}

// restoreBaselines reads the committed baselines of one group revision and
// verifies each against the digest committed with it, the workspace it was
// filed under and the closed cursor: a baseline cannot be ahead of the windows
// the checkpoint says were closed.
func (s *Store) restoreBaselines(ctx context.Context, key GroupKey, nextClose time.Time) ([]CounterBaseline, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT member_key, payload, payload_digest, workspace_id, window_end FROM wg_delivery_counter_baselines
		WHERE group_id = ? AND group_revision = ? ORDER BY member_key`, key.GroupID, key.GroupRevision)
	if err != nil {
		return nil, fmt.Errorf("read counter baselines: %w", err)
	}
	defer rows.Close()
	var baselines []CounterBaseline
	for rows.Next() {
		var member, payload, digest, workspaceID, windowEnd string
		if err := rows.Scan(&member, &payload, &digest, &workspaceID, &windowEnd); err != nil {
			return nil, fmt.Errorf("scan counter baseline: %w", err)
		}
		sum := sha256.Sum256([]byte(payload))
		if hex.EncodeToString(sum[:]) != digest || workspaceID != key.WorkspaceID || windowEnd > timestamp(nextClose) {
			return nil, ErrBaselineCorrupt
		}
		baselines = append(baselines, CounterBaseline{MemberKey: member, Payload: []byte(payload)})
	}
	return baselines, rows.Err()
}
