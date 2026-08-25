package settings

import (
	"context"
	"database/sql"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestSQLRepository_SetIfRevision_ConcurrentConnectionsOnlyOneCommit(t *testing.T) {
	dsn := "file:settings-cas-" + t.Name() + "?mode=memory&cache=shared"
	db1, err := sql.Open("sqlite", dsn)
	require.NoError(t, err)
	defer db1.Close()
	db2, err := sql.Open("sqlite", dsn)
	require.NoError(t, err)
	defer db2.Close()
	migrationContent, err := os.ReadFile("../schema/migrations/001_initial_schema_sqlite.sql")
	require.NoError(t, err)
	_, err = db1.ExecContext(context.Background(), string(migrationContent))
	require.NoError(t, err)
	repo1, repo2 := NewSQLRepository(db1), NewSQLRepository(db2)
	ctx := context.Background()
	require.NoError(t, repo1.Set(ctx, "cas-key", map[string]any{"settings_revision": "rev-1", "value": 1}))

	var wg sync.WaitGroup
	errs := make([]error, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		errs[0] = repo1.SetIfRevision(ctx, "cas-key", "rev-1", map[string]any{"settings_revision": "rev-a", "value": 2})
	}()
	go func() {
		defer wg.Done()
		errs[1] = repo2.SetIfRevision(ctx, "cas-key", "rev-1", map[string]any{"settings_revision": "rev-b", "value": 3})
	}()
	wg.Wait()
	commits := 0
	for _, commitErr := range errs {
		if commitErr == nil {
			commits++
		} else {
			assert.Contains(t, commitErr.Error(), "conflict")
		}
	}
	assert.Equal(t, 1, commits)
}
