package pollinggroup

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionGroupDeleteFailureCleansSQLBeforeFinishingTokens(t *testing.T) {
	for _, failure := range []string{"cancel", "delete-trigger", "commit-busy"} {
		t.Run(failure, func(t *testing.T) {
			db, repo, dsn := acquisitionDeleteFixture(t)
			foreign, err := sql.Open("sqlite", dsn)
			require.NoError(t, err)
			defer foreign.Close()
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			lifecycle := &acquisitionDeleteLifecycle{t: t}
			var reader *sql.Tx
			switch failure {
			case "cancel":
				lifecycle.begin = func() error { cancel(); return nil }
			case "delete-trigger":
				_, err = db.ExecContext(t.Context(), `CREATE TRIGGER reject_group_delete BEFORE DELETE ON polling_groups BEGIN SELECT RAISE(ABORT,'delete rejected'); END`)
				require.NoError(t, err)
			case "commit-busy":
				reader, err = foreign.BeginTx(t.Context(), nil)
				require.NoError(t, err)
				defer reader.Rollback()
				var count int
				require.NoError(t, reader.QueryRowContext(t.Context(), `SELECT count(*) FROM polling_groups`).Scan(&count))
			}
			lifecycle.onFinish = func() {
				// Test-only writer probe: rollback/cleanup must precede the runtime
				// finish callback, which may take scheduler.Lock after SQL ends.
				if reader != nil {
					require.NoError(t, reader.Rollback())
				}
				_, err := foreign.ExecContext(t.Context(), `UPDATE polling_groups SET priority=priority WHERE id='other'`)
				require.NoError(t, err, "finish callback ran while the writer reservation was held")
			}
			repo.SetAcquisitionMutationLifecycle(lifecycle)
			require.Error(t, repo.Delete(ctx, "g"))
			require.Equal(t, 1, lifecycle.begun)
			require.Equal(t, 1, lifecycle.finished)
			require.False(t, lifecycle.active)
			requireAcquisitionGroupMembers(t, db, 3)
			_, err = repo.GetByID(t.Context(), "g")
			require.NoError(t, err)
		})
	}
}

func TestAcquisitionGroupDeleteWithoutHookPreservesSeedBehavior(t *testing.T) {
	db, repo, _ := acquisitionDeleteFixture(t)
	require.NoError(t, repo.Delete(t.Context(), "g"))
	requireAcquisitionGroupMembers(t, db, 0)
	require.ErrorContains(t, repo.Delete(t.Context(), "g"), "輪詢群組不存在")
}
