package workspace

import (
	"context"
	"database/sql"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestCounterBindingApplyRuntimeInsertFailureIsAtomic(t *testing.T) {
	f := newCounterApplyFixture(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeInt64)
	_, err := f.db.ExecContext(t.Context(), `CREATE TRIGGER reject_counter_runtime BEFORE INSERT ON wg_runtime_versions BEGIN SELECT RAISE(ABORT, 'proof-save-failed'); END`)
	require.NoError(t, err)
	activated := false
	f.service.WithWriterOwnershipActivationBarrier(writerOwnershipActivationFunc(func(context.Context, *sql.Tx, WriterOwnershipActivationRequest) error { activated = true; return nil }))
	_, err = f.service.Apply(t.Context(), f.group.ID, f.mutation)
	require.Error(t, err)
	require.False(t, activated)
	f.assertUnapplied(t)
}

func TestCounterBindingApplyCannotBackfillHistoricalProofFromCurrentDefinition(t *testing.T) {
	f := newCounterApplyFixture(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeInt64)
	applied, err := f.service.Apply(t.Context(), f.group.ID, f.mutation)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `DELETE FROM wg_runtime_versions WHERE group_id = ?`, f.group.ID)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET definition_revision='def-2', series_epoch='epoch-2', unit='MWh' WHERE id='measurement-1'`)
	require.NoError(t, err)
	f.mutation.ExpectedWorkspaceRevision = applied.WorkspaceRevision
	_, err = f.service.Apply(t.Context(), f.group.ID, f.mutation)
	require.ErrorIs(t, err, ErrWriteGroupApplyNotReady)
	var count int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_runtime_versions WHERE group_id = ?`, f.group.ID).Scan(&count))
	require.Zero(t, count, "a missing old snapshot cannot be invented from today's definition")
	current, err := f.service.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, applied.WorkspaceRevision, current.DatabaseSetupRevision)
}

func TestCounterBindingApplyReapplyRefusesChangedMeaningAndPreservesOriginalRuntimeBytes(t *testing.T) {
	f := newCounterApplyFixture(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeInt64)
	applied, err := f.service.Apply(t.Context(), f.group.ID, f.mutation)
	require.NoError(t, err)
	_, before := readCounterApplySnapshot(t, f.db.QueryRowContext, f.group.ID)
	_, err = f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET definition_revision='def-2', series_epoch='epoch-2', unit='MWh' WHERE id='measurement-1'`)
	require.NoError(t, err)
	f.mutation.ExpectedWorkspaceRevision = applied.WorkspaceRevision
	// Task 4.1 strengthens the original immutability regression: changed meaning
	// cannot be acknowledged as a successful same-revision reapply.
	_, err = f.service.Apply(t.Context(), f.group.ID, f.mutation)
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	current, err := f.service.Get(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Equal(t, applied.Group.AppliedRevision, current.Group.AppliedRevision)
	require.Equal(t, applied.WorkspaceRevision, current.WorkspaceRevision)
	var versions int
	require.NoError(t, f.db.QueryRow(`SELECT COUNT(*) FROM wg_runtime_versions WHERE group_id=?`, f.group.ID).Scan(&versions))
	require.Equal(t, 1, versions)
	_, after := readCounterApplySnapshot(t, f.db.QueryRowContext, f.group.ID)
	require.Equal(t, before, after)
}

func TestCounterBindingApplyNilRuntimeLayoutKeepsNonCounterCompatibility(t *testing.T) {
	db := openWorkspaceTestDB(t, ":memory:")
	t.Cleanup(func() { _ = db.Close() })
	service, fixture, created := newReadinessGroup(t.Context(), t, db)
	mutation := WriteGroupMutation{WorkspaceID: fixture.workspaceID, ExpectedWorkspaceRevision: created.WorkspaceRevision,
		ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1"}
	require.NotPanics(t, func() {
		_, err := service.workspaceSvc.UpdateDatabaseSetup(t.Context(), mutation.ExpectedWorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
			_, err := service.applyGroupInTx(ctx, tx, record, created.Group.ID, mutation, &WriteGroupReadiness{}, nextAlignedCutover, false)
			return err
		})
		require.NoError(t, err, "legacy non-counter domain-only apply never required a runtime descriptor")
	})
}

func TestCounterBindingApplyNilRuntimeLayoutRefusesCounter(t *testing.T) {
	f := newCounterApplyFixture(t, WriteGroupLayoutManagedLongV1, WriteGroupProfileCumulative, schema.DataTypeInt64)
	ready, err := f.service.Readiness(t.Context(), f.group.ID)
	require.NoError(t, err)
	ready.runtimeLayout = nil
	require.NotPanics(t, func() {
		_, err := f.service.workspaceSvc.UpdateDatabaseSetup(t.Context(), f.mutation.ExpectedWorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
			_, err := f.service.applyGroupInTx(ctx, tx, record, f.group.ID, f.mutation, ready, nextAlignedCutover, false)
			return err
		})
		require.Error(t, err)
	})
	f.assertUnapplied(t)
}
