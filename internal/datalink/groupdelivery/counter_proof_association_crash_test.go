package groupdelivery

import (
	"sync/atomic"
	"testing"

	"go-gateway/internal/testutil/lostcommit"

	"github.com/stretchr/testify/require"
)

func TestLegacyCounterProofAssociationLostCommitResponseRetriesWithoutRewriting(t *testing.T) {
	store, path, descriptor, refs := legacyAssociationFixture(t)
	before := associationJournal(t, store)
	lose := &atomic.Bool{}
	lose.Store(true)
	interrupted := lostcommit.Open(path+"?_pragma=busy_timeout(5000)", lose)
	err := NewStore(interrupted).AssociateLegacyCounterProofs(t.Context(), testKey, refs)
	require.ErrorIs(t, err, lostcommit.ErrResponseLost, "a lost response is not reported as a successful upgrade")
	require.Equal(t, 2, countRows(t, store.DB(), "wg_counter_proof_associations"), "the real SQLite transaction committed both or neither")
	audit := associationAudit(t, store)
	require.NoError(t, interrupted.Close())
	require.NoError(t, store.DB().Close())
	recovered := NewStore(openStoreDB(t, path))
	require.NoError(t, recovered.AssociateLegacyCounterProofs(t.Context(), testKey, refs))
	require.Equal(t, audit, associationAudit(t, recovered))
	require.Equal(t, before, associationJournal(t, recovered))
	persisted, err := recovered.RuntimeVersion(t.Context(), testKey)
	require.NoError(t, err)
	require.Equal(t, descriptor, persisted)
}
