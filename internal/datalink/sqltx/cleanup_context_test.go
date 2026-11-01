package sqltx

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type cleanupContextKey struct{}

func TestOwnedCleanupDetachesCancellationAndStartsFreshBudget(t *testing.T) {
	for _, commit := range []bool{true, false} {
		t.Run(map[bool]string{true: "commit", false: "rollback"}[commit], func(t *testing.T) {
			connector := &cleanupFaultConnector{}
			db := sql.OpenDB(connector)
			t.Cleanup(func() { require.NoError(t, db.Close()) })
			db.SetMaxOpenConns(1)
			_, err := db.ExecContext(t.Context(), "CREATE TABLE probe(value INTEGER); INSERT INTO probe VALUES(1)")
			require.NoError(t, err)
			parent, cancel := context.WithTimeout(context.WithValue(t.Context(), cleanupContextKey{}, "original-owner"), time.Hour)
			defer cancel()
			owned, err := Begin(parent, db, nil)
			require.NoError(t, err)
			require.NotNil(t, owned.cleanupParent)
			_, hasDeadline := owned.cleanupParent.Deadline()
			require.False(t, hasDeadline, "cleanup budget must not start at Begin")
			require.NoError(t, owned.Tx.Rollback(), "force terminal SQL Tx while the dedicated conn is still owned")
			cancel()
			require.NoError(t, owned.cleanupParent.Err())
			started := time.Now()
			observed := false
			connector.observeCleanup = func(ctx context.Context) {
				observed = true
				require.NoError(t, ctx.Err(), "canceled request must not cancel native cleanup")
				require.Equal(t, "original-owner", ctx.Value(cleanupContextKey{}))
				deadline, ok := ctx.Deadline()
				require.True(t, ok)
				require.False(t, deadline.Before(started.Add(5*time.Second)), "five-second budget begins with cleanup")
				require.False(t, deadline.After(time.Now().Add(5*time.Second)))
			}
			if commit {
				err = owned.Commit()
				require.ErrorIs(t, err, ErrCommitOutcomeUnknown)
			} else {
				err = owned.Rollback()
			}
			require.ErrorIs(t, err, sql.ErrTxDone)
			require.True(t, observed)
			var value int
			require.NoError(t, db.QueryRowContext(t.Context(), "SELECT value FROM probe").Scan(&value))
			require.Equal(t, 1, value, "clean dedicated connection preserves the private memory database")
		})
	}
}
