package groupdelivery

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/testutil/lostcommit"
)

var handoffMetadataReaders = []struct {
	name string
	read func(context.Context, *Store, GroupKey) (int, error)
}{
	{"GroupStatus", func(ctx context.Context, store *Store, key GroupKey) (int, error) {
		status, err := store.GroupStatus(ctx, key.GroupID)
		return status.Stages.Collecting, err
	}},
	{"PendingSampleCount", func(ctx context.Context, store *Store, key GroupKey) (int, error) {
		return store.PendingSampleCount(ctx, key)
	}},
	{"RevisionStageCounts", func(ctx context.Context, store *Store, key GroupKey) (int, error) {
		stages, err := store.RevisionStageCounts(ctx, key.GroupID, key.GroupRevision)
		return stages.Collecting, err
	}},
}

func acceptedMetadataHandoff(t *testing.T) *Store {
	t.Helper()
	store, _ := seedPreparedHandoff(t)
	batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	require.Equal(t, 1, batch.Accepted)
	return store
}

func TestHandoffMetadataReadersUseOneConnection(t *testing.T) {
	for _, reader := range handoffMetadataReaders {
		t.Run(reader.name, func(t *testing.T) {
			store := acceptedMetadataHandoff(t)
			count, err := reader.read(t.Context(), store, otherKey)
			require.NoError(t, err, "healthy handoff passes before limiting the pool")
			require.Equal(t, 1, count)
			store.DB().SetMaxOpenConns(1)
			ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
			defer cancel()
			for range 2 {
				for _, key := range []GroupKey{otherKey, testKey} {
					count, err = reader.read(ctx, store, key)
					require.NoError(t, err, "nested receipt/ownership reads must use the pinned connection")
					want := 0
					if key == otherKey {
						want = 1
					}
					require.Equal(t, want, count, "only the target owns the transferred input")
					require.Zero(t, store.DB().Stats().InUse, "each validation releases its connection")
				}
			}
			require.Equal(t, 2, countRows(t, store.DB(), "wg_delivery_samples"))
			require.Zero(t, openCount(t, store.DB(), "wg_delivery_samples", "consumed=1"))
			require.Equal(t, 1, countRows(t, store.DB(), "wg_delivery_sample_handoffs"))
		})
	}
}

func TestHandoffMetadataReadersIsolateDamagedScope(t *testing.T) {
	store := acceptedMetadataHandoff(t)
	healthy := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "healthy-metadata", GroupRevision: "rev-healthy"}
	require.NoError(t, store.AppendSample(t.Context(), healthy, at(20), sampleAt("healthy", "temp", 23, 6)))
	_, err := store.DB().Exec(`UPDATE wg_delivery_sample_handoffs SET record='{}'`)
	require.NoError(t, err)
	store.DB().SetMaxOpenConns(1)
	for _, reader := range handoffMetadataReaders {
		t.Run(reader.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
			defer cancel()
			count, err := reader.read(ctx, store, otherKey)
			require.ErrorIs(t, err, ErrOwnershipUnavailable)
			require.Zero(t, count, "a damaged scope must not report successful counts")
			require.Zero(t, store.DB().Stats().InUse, "failed verification rolls back")
			count, err = reader.read(ctx, store, healthy)
			require.NoError(t, err, "an unrelated healthy scope still progresses")
			require.Equal(t, 1, count)
			require.Zero(t, store.DB().Stats().InUse)
		})
	}
	// Revision-specific readers do not validate another revision's damaged audit.
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	unrelatedRevision := otherKey
	unrelatedRevision.GroupRevision = "not-applied"
	count, err := store.PendingSampleCount(ctx, unrelatedRevision)
	require.NoError(t, err)
	require.Zero(t, count)
	stages, err := store.RevisionStageCounts(ctx, otherKey.GroupID, unrelatedRevision.GroupRevision)
	require.NoError(t, err)
	require.Equal(t, StageCounts{}, stages)
}

func TestHandoffMetadataReadersRetainSafeStorageCause(t *testing.T) {
	store := acceptedMetadataHandoff(t)
	path := handoffDatabasePath(t, store)
	for _, reader := range handoffMetadataReaders {
		for _, phase := range []string{"begin", "scope-read", "ownership-roots", "ownership-version", "commit"} {
			t.Run(reader.name+"/"+phase, func(t *testing.T) {
				cause := errors.New("private-db-path-and-raw-SQL")
				faults := &lostcommit.Faults{Lose: &atomic.Bool{}}
				var commits atomic.Int32
				afterCommit := func() { commits.Add(1) }
				faults.AfterCommit.Store(&afterCommit)
				if phase == "commit" {
					faults.Lose.Store(true)
					cause = lostcommit.ErrResponseLost
				}
				queries := map[string]string{
					"scope-read":        "SELECT DISTINCT workspace_id,source_group_id,source_revision",
					"ownership-roots":   "SELECT c.id,c.boundary,c.detail",
					"ownership-version": "SELECT effective_at,payload FROM write_group_versions",
				}
				reject := func(query string) error {
					if match := queries[phase]; match != "" && strings.Contains(query, match) {
						return cause
					}
					return nil
				}
				faults.Reject.Store(&reject)
				db := lostcommit.OpenFaults(path, faults)
				defer db.Close()
				db.SetMaxOpenConns(1)
				ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
				defer cancel()
				if phase == "begin" {
					cancel()
					cause = context.Canceled
				}
				count, err := reader.read(ctx, NewStore(db), otherKey)
				require.Zero(t, count)
				require.ErrorIs(t, err, cause)
				require.ErrorIs(t, err, ErrOwnershipUnavailable)
				var storage *handoffStorageError
				require.ErrorAs(t, err, &storage)
				require.Equal(t, ErrOwnershipUnavailable.Error(), err.Error())
				require.NotContains(t, err.Error(), cause.Error())
				require.Zero(t, db.Stats().InUse)
				wantCommits := int32(0)
				if phase == "commit" {
					wantCommits = 1
				}
				require.Equal(t, wantCommits, commits.Load(), "failed validation must not commit")
				faults.Reject.Store(nil)
				faults.Lose.Store(false)
				retryCtx, retryCancel := context.WithTimeout(t.Context(), 2*time.Second)
				defer retryCancel()
				count, err = reader.read(retryCtx, NewStore(db), otherKey)
				require.NoError(t, err, "cleared storage fault permits retry on the same pool")
				require.Equal(t, 1, count)
				require.Equal(t, wantCommits+1, commits.Load())
				require.Zero(t, db.Stats().InUse)
			})
		}
	}
}

func TestHandoffMetadataValidationReusesCallerTransaction(t *testing.T) {
	store := acceptedMetadataHandoff(t)
	faults := &lostcommit.Faults{}
	var commits atomic.Int32
	afterCommit := func() { commits.Add(1) }
	faults.AfterCommit.Store(&afterCommit)
	db := lostcommit.OpenFaults(handoffDatabasePath(t, store), faults)
	defer db.Close()
	db.SetMaxOpenConns(1)
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	require.NoError(t, err)
	defer tx.Rollback()
	for range 2 {
		accepted, err := readLegacyAcceptedSet(ctx, tx, otherKey)
		require.NoError(t, err, "the existing transactional caller must not acquire a second connection")
		require.Empty(t, accepted, "transferred input is not a native legacy association candidate")
		require.Equal(t, 1, db.Stats().InUse, "the caller still owns its transaction")
		require.Zero(t, commits.Load(), "validation must not commit the caller's transaction")
	}
	require.NoError(t, tx.Commit())
	require.Equal(t, int32(1), commits.Load())
	require.Zero(t, db.Stats().InUse)
}

func TestHandoffMetadataReaderEmptyInputsAvoidStorage(t *testing.T) {
	store, db, _ := newTestStore(t)
	require.NoError(t, db.Close())
	count, err := store.PendingSampleCount(t.Context(), GroupKey{GroupID: testKey.GroupID})
	require.ErrorIs(t, err, ErrInvalidGroupKey)
	require.Zero(t, count)
	stages, err := store.RevisionStageCounts(t.Context(), testKey.GroupID, "")
	require.NoError(t, err)
	require.Equal(t, StageCounts{}, stages)
}
