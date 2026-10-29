package workspace

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccessorActivationConcurrentExactConfirmation(t *testing.T) {
	h, id, req := atomicActivationFixture(t)
	var seq int
	var name, path string
	require.NoError(t, h.db.QueryRow("PRAGMA database_list").Scan(&seq, &name, &path))
	peerDB := openWorkspaceTestDB(t, path)
	t.Cleanup(func() { peerDB.Close() })
	peerDB.SetMaxOpenConns(1)
	for _, db := range []*sql.DB{h.db, peerDB} {
		_, err := db.Exec("PRAGMA busy_timeout=1000")
		require.NoError(t, err)
	}
	peer := activationTestService(h.service)
	peer.repo = NewSQLWriteGroupRepository(peerDB)
	peer.workspaceSvc = NewService(NewSQLRepository(peerDB))
	services := []*WriteGroupService{h.service, peer}
	results := make([]*WriteGroupCutover, 2)
	errs := make([]error, 2)
	gate := make(chan struct{})
	var wg sync.WaitGroup
	for i := range services {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-gate
			results[i], errs[i] = services[i].ConfirmSuccessorActivation(t.Context(), id, req)
		}(i)
	}
	close(gate)
	wg.Wait()
	for _, err := range errs {
		require.NoError(t, err, "identical competing confirmation must resolve durable winner")
	}
	require.Equal(t, results[0], results[1])
	require.Equal(t, 1, h.countRows("SELECT COUNT(*) FROM write_group_cutovers"))
	require.Equal(t, 2, h.countRows("SELECT COUNT(*) FROM write_group_generations"))
}
