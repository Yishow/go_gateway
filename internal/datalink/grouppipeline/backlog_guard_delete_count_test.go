package grouppipeline

import (
	"context"
	"database/sql"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestBacklogGuardCountsPermanentDeletionEvidence(t *testing.T) {
	for _, kind := range []string{"sample", "outbox"} {
		t.Run(kind, func(t *testing.T) {
			db, group := backlogOwnershipFixture(t, kind)
			tx, err := db.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			defer tx.Rollback()
			counter, ok := any(BacklogGuard{}).(interface {
				CountWriteGroupDeleteBacklog(context.Context, *sql.Tx, *workspace.WriteGroup, *workspace.SQLWriteGroupRepository) (workspace.WriteGroupDeleteBacklog, error)
			})
			require.True(t, ok, "the production backlog guard must expose deletion counts")
			counts, err := counter.CountWriteGroupDeleteBacklog(t.Context(), tx, group, workspace.NewSQLWriteGroupRepository(db))
			require.NoError(t, err)
			if kind == "sample" {
				require.Equal(t, 1, counts.SamplesUnconsumed)
				require.Empty(t, counts.Outbox)
			} else {
				require.Zero(t, counts.SamplesUnconsumed)
				require.Equal(t, map[string]int{"pending": 1}, counts.Outbox)
			}
		})
	}
}
