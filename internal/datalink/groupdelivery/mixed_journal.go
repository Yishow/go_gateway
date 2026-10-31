package groupdelivery

import (
	"context"
	"go-gateway/internal/datalink/snapshot"
	"time"
)

func validateMixedJournalWindow(ctx context.Context, q ownershipQuery, key GroupKey, bucket time.Time, sample snapshot.Sample) error {
	d, _, err := readMixedDescriptor(ctx, q, key)
	if err != nil {
		return err
	}
	if d == nil {
		return nil
	}
	intervals, err := d.intervals()
	if err != nil {
		return err
	}
	interval, ok := intervals[sample.MemberKey]
	if !ok {
		return ErrInvalidSample
	}
	start, _, err := snapshot.CheckedWindowBounds(sample.ObservedAt, interval)
	if err != nil {
		return err
	}
	if !bucket.Equal(start) {
		return ErrInvalidSample
	}
	return nil
}
func refuseMixedLegacyCursor(ctx context.Context, q ownershipQuery, key GroupKey) error {
	d, _, err := readMixedDescriptor(ctx, q, key)
	if err != nil {
		return err
	}
	if d != nil {
		return mixedInvalid("legacy closure for mixed descriptor")
	}
	var count int
	if err := q.QueryRowContext(ctx, `SELECT count(*) FROM wg_delivery_mixed_progress WHERE workspace_id=? AND group_id=? AND group_revision=?`, key.WorkspaceID, key.GroupID, key.GroupRevision).Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		return mixedInvalid("orphan mixed cursor")
	}
	return nil
}
