package groupdelivery

import (
	"context"
	"database/sql"
	"encoding/json"
	"sort"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/snapshot"
)

// OwnershipView attests the existing cutover authority and actual owner bounds
// observed with restored inputs. Its zero value proves no cutover only when
// CommitClosure independently observes the same absence in its write transaction.
// Private, comparable fields prevent a caller from relabelling a stale view.
type OwnershipView struct {
	key        GroupKey
	digest     string
	start, end time.Time
	interval   time.Duration
}

// OwnershipView reads a consistent routing hint without replaying the journal.
// It may trigger a rebuild, but must never be attached to an existing assembler:
// only Restore captures ownership and the complete replay state together.
func (s *Store) OwnershipView(ctx context.Context, key GroupKey) (OwnershipView, error) {
	if !key.Valid() {
		return OwnershipView{}, ErrInvalidGroupKey
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return OwnershipView{}, ErrOwnershipUnavailable
	}
	defer tx.Rollback()
	view, err := captureOwnershipView(ctx, tx, key)
	if err != nil {
		return OwnershipView{}, err
	}
	if tx.Commit() != nil {
		return OwnershipView{}, ErrOwnershipUnavailable
	}
	return view, nil
}

func captureOwnershipView(ctx context.Context, q ownershipQuery, key GroupKey) (OwnershipView, error) {
	ownership, err := resolveCutoverOwnership(ctx, q, OwnershipSource{WorkspaceID: key.WorkspaceID}, time.Time{}, key)
	if err != nil {
		return OwnershipView{}, err
	}
	if !ownership.Affected[key] {
		return OwnershipView{}, nil
	}
	bounds, ok := ownership.bounds[key]
	if !ok || bounds.Start.IsZero() {
		return OwnershipView{}, ErrOwnershipUnavailable
	}
	var seconds int64
	if err := q.QueryRowContext(ctx, `SELECT json_extract(payload,'$.row_policy.interval_seconds') FROM write_group_versions WHERE group_id=? AND group_revision=?`, key.GroupID, key.GroupRevision).Scan(&seconds); err != nil || seconds <= 0 || seconds > int64((1<<63-1)/time.Second) {
		return OwnershipView{}, ErrOwnershipUnavailable
	}
	ids := make([]string, 0, len(ownership.roots))
	for id := range ownership.roots {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	raw, _ := json.Marshal(ids)
	return OwnershipView{key: key, digest: handoffHash(string(raw)), start: bounds.Start.UTC(), end: bounds.End.UTC(), interval: time.Duration(seconds) * time.Second}, nil
}

func ownershipCeil(at time.Time, interval time.Duration) time.Time {
	floor := snapshot.BucketStart(at, interval)
	if floor.Equal(at) {
		return floor
	}
	return floor.Add(interval)
}

// ValidateBounds refuses an old assembler configuration paired with a newly
// restored view. A checkpoint may advance the first bucket; it cannot remove
// an actual cutover end or widen the first, partial owner window.
func (v OwnershipView) ValidateBounds(key GroupKey, firstBucket, startsAt, until, endsAt time.Time, interval time.Duration) error {
	if v == (OwnershipView{}) {
		return nil
	}
	if v.key != key || interval != v.interval || firstBucket.IsZero() || !snapshot.BucketStart(firstBucket, interval).Equal(firstBucket) {
		return ErrOwnershipConflict
	}
	floor := snapshot.BucketStart(v.start, interval)
	if firstBucket.Before(floor) || (!startsAt.IsZero() && !startsAt.Equal(v.start)) {
		return ErrOwnershipConflict
	}
	if firstBucket.Equal(floor) && !v.start.Equal(floor) && !startsAt.Equal(v.start) {
		return ErrOwnershipConflict
	}
	if !v.end.IsZero() {
		ceiling := ownershipCeil(v.end, interval)
		if until.IsZero() || !snapshot.BucketStart(until, interval).Equal(until) || until.After(ceiling) {
			return ErrOwnershipConflict
		}
		// An already disabled/draining boundary may stop earlier. Do not
		// extend it merely to reach a later, scheduled cutover boundary.
		if until.Equal(ceiling) {
			if !endsAt.Equal(v.end) {
				return ErrOwnershipConflict
			}
		} else if !endsAt.IsZero() && !endsAt.Equal(until) {
			return ErrOwnershipConflict
		}
	}
	return nil
}

func validateOwnershipClosure(ctx context.Context, tx *sql.Tx, key GroupKey, closure Closure) error {
	current, err := captureOwnershipView(ctx, tx, key)
	if err != nil {
		return err
	}
	if current != closure.OwnershipView {
		return ErrOwnershipConflict
	}
	if current == (OwnershipView{}) {
		return nil
	} // legacy callers without a relevant cutover
	if !snapshot.BucketStart(closure.NextClose, current.interval).Equal(closure.NextClose) || closure.NextClose.Before(snapshot.BucketStart(current.start, current.interval)) || (!current.end.IsZero() && closure.NextClose.After(ownershipCeil(current.end, current.interval))) {
		return ErrOwnershipConflict
	}
	for _, bucket := range closure.Buckets {
		outcome := bucket.Outcome
		start, end := outcome.BucketStart, outcome.BucketEnd
		if !snapshot.BucketStart(start, current.interval).Equal(start) || !end.Equal(start.Add(current.interval)) || end.After(closure.NextClose) {
			return ErrOwnershipConflict
		}
		if current.start.After(start) {
			start = current.start
		}
		if !current.end.IsZero() && current.end.Before(end) {
			end = current.end
		}
		actualStart, actualEnd := outcome.Window()
		if !start.Before(end) || !actualStart.Equal(start) || !actualEnd.Equal(end) {
			return ErrOwnershipConflict
		}
		for _, member := range outcome.Members {
			if member.Sample != nil && (member.Sample.ObservedAt.Before(start) || !member.Sample.ObservedAt.Before(end)) {
				return ErrOwnershipConflict
			}
		}
	}
	for _, baseline := range closure.Baselines {
		checkpoint, err := aggregation.DecodeCounterCheckpoint(baseline.Payload)
		if err != nil || checkpoint.WindowEnd.Before(current.start) || (!current.end.IsZero() && checkpoint.WindowEnd.After(current.end)) || checkpoint.Reading.ObservedAt.Before(current.start) || (!current.end.IsZero() && !checkpoint.Reading.ObservedAt.Before(current.end)) {
			return ErrOwnershipConflict
		}
	}
	return nil
}
