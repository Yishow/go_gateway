package api

import (
	"database/sql"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/testutil/lostcommit"
	"modernc.org/sqlite"
)

func TestSuccessorFirstStartHTTPRecoversCommittedTransactionLostACK(t *testing.T) {
	lose := &atomic.Bool{}
	r := newFirstStartAPIRig(t, "sqlite", false, true, false, workspace.WriteGroupStorageModePlain, func(path string) *sql.DB { return lostcommit.Wrap(&sqlite.Driver{}, path, lose) })
	request := r.confirmation(t)
	lose.Store(true)
	result := r.activationAPI().postActivation(t, request)
	lose.Store(false)
	require.Equal(t, "succeeded", result["status"])
	require.Equal(t, r.prepared.Group.ID, result["group_id"])
	after := successorAPIAllState(t, r.f.db)
	require.Equal(t, result, r.activationAPI().postActivation(t, request))
	require.Equal(t, result, successorStatus(t, r.successorAPIFixture)["first_start"].(map[string]any)["result"])
	require.Equal(t, after, successorAPIAllState(t, r.f.db), "an acknowledgement lost after local COMMIT must not create another generation or runtime version")
	var versions int
	require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM write_group_versions WHERE group_id=?`, r.prepared.Group.ID).Scan(&versions))
	require.Equal(t, 1, versions)
}
