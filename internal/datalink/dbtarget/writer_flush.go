package dbtarget

import (
	"context"
	"errors"
	"maps"
	"slices"
	"time"
)

// flushConnector reports the revisions that actually wrote a snapshot and any errors.
// Skipped or absent data must not be mistaken for a recovered destination.
func (w *Writer) flushConnector(ctx context.Context, connectorID string, now time.Time, all bool, retryDelay time.Duration) (map[string]struct{}, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var failures []error
	wrote := make(map[string]struct{})
	for write := w.takeDuePending(connectorID, now.UTC(), all); write != nil; write = w.takeDuePending(connectorID, now.UTC(), all) {
		writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		written, err := w.writePending(writeCtx, write, retryDelay)
		cancel()
		if written {
			wrote[write.Connector.IdentityRevision] = struct{}{}
		}
		if err != nil {
			failures = append(failures, err)
			if ctx.Err() != nil {
				return wrote, errors.Join(failures...)
			}
			var retryErr *snapshotRetryError
			if errors.As(err, &retryErr) {
				return wrote, errors.Join(failures...)
			}
		}
	}
	err := w.flushBuckets(ctx, connectorID, now, all, retryDelay, wrote)
	if err != nil {
		failures = append(failures, err)
	}
	return wrote, errors.Join(failures...)
}

// flushDue flushes every connector synchronously.
func (w *Writer) flushDue(ctx context.Context, now time.Time) error {
	ids := w.connectorsWithPendingData()
	var failures []error
	for _, id := range slices.Sorted(maps.Keys(ids)) {
		if _, err := w.flushConnector(ctx, id, now, false, 0); err != nil {
			failures = append(failures, err)
		}
	}
	return errors.Join(failures...)
}
