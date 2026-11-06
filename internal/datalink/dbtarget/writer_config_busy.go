package dbtarget

import (
	"context"
	"errors"
	"time"

	"modernc.org/sqlite"
)

// Configuration refresh and outcome updates may briefly overlap on SQLite.
// Retry only SQLITE_BUSY; persistent locks and other failures remain visible.
func retryWriterConfigBusy(ctx context.Context, operation func() error) error {
	var err error
	for attempt := range 5 {
		err = operation()
		var sqliteErr *sqlite.Error
		if err == nil || !errors.As(err, &sqliteErr) || sqliteErr.Code()&0xff != 5 || attempt == 4 {
			return err
		}
		timer := time.NewTimer(time.Duration(10<<attempt) * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return errors.Join(err, ctx.Err())
		case <-timer.C:
		}
	}
	return err
}
