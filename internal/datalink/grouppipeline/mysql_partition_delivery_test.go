package grouppipeline

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestMySQLPartitionProofReadRetriesOnlyTransientFailures(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "proof.db")+"?_pragma=busy_timeout(1)")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	db.SetMaxOpenConns(2)
	_, err = db.ExecContext(t.Context(), `CREATE TABLE evidence (value TEXT)`)
	require.NoError(t, err)
	blocker, err := db.Conn(t.Context())
	require.NoError(t, err)
	defer blocker.Close()
	_, err = blocker.ExecContext(t.Context(), `BEGIN EXCLUSIVE`)
	require.NoError(t, err)
	var count int
	busy := db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM evidence`).Scan(&count)
	require.Error(t, busy)
	require.ErrorIs(t, mysqlPartitionProofReadError(busy), errMySQLPartitionProofReadTransient)
	_, err = blocker.ExecContext(t.Context(), `ROLLBACK`)
	require.NoError(t, err)
	invalid := db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM missing_evidence`).Scan(&count)
	require.Error(t, invalid)
	require.NotErrorIs(t, mysqlPartitionProofReadError(invalid), errMySQLPartitionProofReadTransient)
	require.NotErrorIs(t, mysqlPartitionProofReadError(recordingplan.ErrMySQLPartitionMaintenanceStateMismatch), errMySQLPartitionProofReadTransient)
	require.ErrorIs(t, mysqlPartitionProofReadError(context.Canceled), errMySQLPartitionProofReadTransient)
}

func TestMySQLPartitionResolverRetriesTransitionButBlocksScopeMismatch(t *testing.T) {
	for _, test := range []struct {
		name    string
		err     error
		blocked bool
	}{
		{name: "temporary local proof read", err: errMySQLPartitionProofReadTransient},
		{name: "partition changed before row transaction", err: dbtarget.ErrMySQLManagedPartitionChanged},
		{name: "confirmed scope mismatch", err: recordingplan.ErrMySQLPartitionMaintenanceStateMismatch, blocked: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			group, layout := mysqlNativeMonthlyFixture(t)
			snapshot := snapshotOf(group, at(0))
			snapshot.RuntimeLayout = layout
			resolver, item, opened := newMySQLResolverFixture(t, snapshot, mysqlRuntimeCapability(layout.MySQL))
			resolver.loadRuntimeSnapshot = func(context.Context, groupdelivery.GroupKey) (*workspace.WriteGroupAppliedSnapshot, error) {
				return snapshot, nil
			}
			resolver.loadMySQLProof = func(context.Context, *workspace.WriteGroupAppliedSnapshot) (*dbtarget.MySQLManagedRowProof, error) {
				return nil, test.err
			}
			target, err := resolver.Resolve(t.Context(), item)
			require.Error(t, err)
			require.Nil(t, target.DB, "a retry may never use an unverified target")
			if test.blocked {
				require.ErrorIs(t, err, groupdelivery.ErrTargetBlocked)
			} else {
				require.ErrorIs(t, err, test.err)
				require.NotErrorIs(t, err, groupdelivery.ErrTargetBlocked)
			}
			require.Equal(t, 1, opened.closeCalls)
		})
	}
}
