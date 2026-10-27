package workspace

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/testutil/lostcommit"
)

var errCounterBeforeCommit = errors.New("injected before local proof commit")

type counterApplyCommitRepository struct {
	*SQLRepository
	fail, after  bool
	commits      int
	pendingProof string
}

func (r *counterApplyCommitRepository) BeginDatabaseSetup(ctx context.Context) (DatabaseSetupTx, error) {
	tx, err := r.SQLRepository.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, err
	}
	return &counterApplyCommitTx{DatabaseSetupTx: tx, repo: r}, nil
}

type counterApplyCommitTx struct {
	DatabaseSetupTx
	repo *counterApplyCommitRepository
}

func (tx *counterApplyCommitTx) Commit() error {
	tx.repo.commits++
	if !tx.repo.fail {
		return tx.DatabaseSetupTx.Commit()
	}
	// Observe the real pending descriptor inside the same transaction; this is
	// not a stand-in for persistence or a fake successful commit.
	if err := tx.SQLTx().QueryRow(`SELECT payload FROM wg_runtime_versions`).Scan(&tx.repo.pendingProof); err != nil {
		return err
	}
	if !tx.repo.after {
		return errCounterBeforeCommit
	}
	if err := tx.DatabaseSetupTx.Commit(); err != nil {
		return err
	}
	return lostcommit.ErrResponseLost
}

func TestCounterBindingProofApplyCommitCrashBoundaries(t *testing.T) {
	for _, after := range []bool{false, true} {
		name := "before commit"
		if after {
			name = "committed response lost"
		}
		t.Run(name, func(t *testing.T) {
			f := newCounterApplyFixture(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeUint64)
			repo := &counterApplyCommitRepository{SQLRepository: NewSQLRepository(f.db), fail: true, after: after}
			f.service.workspaceSvc.repo = repo
			result, err := f.service.Apply(t.Context(), f.group.ID, f.mutation)
			require.Nil(t, result, "an unconfirmed commit is not a successful Apply response")
			wantErr := errCounterBeforeCommit
			if after {
				wantErr = lostcommit.ErrResponseLost
			}
			require.ErrorIs(t, err, wantErr)
			require.Equal(t, 1, repo.commits)
			require.Contains(t, repo.pendingProof, `"counter_bindings"`)
			if !after {
				f.assertUnapplied(t)
			} else {
				frozen, raw := counterEpochRuntime(t, f, f.group.Revision)
				require.Equal(t, repo.pendingProof, raw)
				require.Len(t, frozen.CounterBindings, 1)
				current, err := f.service.Get(t.Context(), f.group.ID)
				require.NoError(t, err)
				require.Equal(t, f.group.Revision, current.Group.AppliedRevision)
				require.NotEqual(t, f.mutation.ExpectedWorkspaceRevision, current.WorkspaceRevision)
			}
			repo.fail = false
			// Reload after ambiguous commit; a stale caller must not invent a second
			// revision or claim the old CAS token is still current.
			current, err := f.service.Get(t.Context(), f.group.ID)
			require.NoError(t, err)
			retry, err := f.service.Apply(t.Context(), f.group.ID, counterEpochMutation(t, f, current.Group))
			require.NoError(t, err)
			require.Equal(t, f.group.Revision, retry.Group.AppliedRevision)
			_, persisted := counterEpochRuntime(t, f, f.group.Revision)
			if after {
				require.Equal(t, repo.pendingProof, persisted)
			}
			var versions, samples int
			require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM wg_runtime_versions`).Scan(&versions))
			require.Equal(t, 1, versions)
			require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples`).Scan(&samples))
			require.Zero(t, samples)
		})
	}
}
