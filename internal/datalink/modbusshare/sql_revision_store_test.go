package modbusshare

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestSQLWorkspaceRevisionStore_StorageFailuresAreTyped(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	store := NewSQLWorkspaceRevisionStore(db)
	require.NoError(t, db.Close())

	_, _, err = store.GetRevision(context.Background(), "ws-closed")
	var storageErr *StorageError
	require.ErrorAs(t, err, &storageErr)
	require.Equal(t, "read workspace revision", storageErr.Operation)
}

func TestSQLWorkspaceRevisionStore_PersistsAndCASChecksAcrossInstances(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	require.NoError(t, datalink.NewMigrator().Migrate(db))

	store := NewSQLWorkspaceRevisionStore(db)
	ctx := context.Background()
	rev, dirty, err := store.GetRevision(ctx, "ws-restart")
	require.NoError(t, err)
	require.Equal(t, "rev-1", rev)
	require.False(t, dirty)
	require.NoError(t, store.UpdateRevision(ctx, "ws-restart", "rev-1", "rev-2"))

	restarted := NewSQLWorkspaceRevisionStore(db)
	rev, dirty, err = restarted.GetRevision(ctx, "ws-restart")
	require.NoError(t, err)
	require.Equal(t, "rev-2", rev)
	require.False(t, dirty)
	err = restarted.UpdateRevision(ctx, "ws-restart", "rev-1", "rev-3")
	require.Error(t, err)
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	require.Equal(t, ErrCodeRevisionConflict, shareErr.Code)
}

func TestSQLWorkspaceRevisionStore_CommitsDesiredMappingsWithRevision(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	store := NewSQLWorkspaceRevisionStore(db)
	desired := []DesiredMapping{{WorkspaceID: "ws-a", SourceRuleID: "rule-a", SourceRuleRevision: "r1", TagID: "tag-a", ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}}
	require.NoError(t, store.CommitDesiredMappings(context.Background(), "ws-a", "rev-1", "rev-2", desired))
	got, err := store.GetDesiredMappings(context.Background(), "ws-a")
	require.NoError(t, err)
	require.Equal(t, desired, got)
	rev, dirty, err := store.GetRevision(context.Background(), "ws-a")
	require.NoError(t, err)
	require.Equal(t, "rev-2", rev)
	require.False(t, dirty)
}

func TestReconciler_SQLProjectionFailureRestoresOldDesiredAndDirtyUnknownOnRollbackFailure(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	require.NoError(t, datalink.NewMigrator().Migrate(db))

	store := NewSQLWorkspaceRevisionStore(db)
	svc := NewService(nil, 65536)
	svc.SetHydrationState(HydrationState{State: HydrationStateReady, WorkspaceID: "ws-sql", Readiness: true})
	t.Cleanup(func() { _ = svc.Stop() })
	reconciler := NewReconciler(svc, store).WithOwnershipValidator(func(context.Context, DesiredMapping) error {
		return nil
	})
	old := DesiredMapping{WorkspaceID: "ws-sql", SourceRuleID: "rule", SourceRuleRevision: "rev", TagID: "tag", DataType: schema.DataTypeInt16, ZeroBasedRegister: 0, ShareStartRegister: 40001, SpanRegisters: 1, StrideRegisters: 1}
	first, err := reconciler.Reconcile(context.Background(), ReconcileRequest{WorkspaceID: "ws-sql", ExpectedWorkspaceRevision: "rev-1", ExpectedSettingsRevision: "set-1", DesiredMappings: []DesiredMapping{old}})
	require.NoError(t, err)

	svc.SetProjectionSwapHook(func(map[string]TagMirrorMapping) error { return errors.New("injected runtime projection failure") })
	_, err = reconciler.Reconcile(context.Background(), ReconcileRequest{WorkspaceID: "ws-sql", ExpectedWorkspaceRevision: first.NewWorkspaceRevision, ExpectedSettingsRevision: "set-1", DesiredMappings: nil})
	var shareErr *Error
	require.ErrorAs(t, err, &shareErr)
	require.Equal(t, ErrCodeReconcileFailed, shareErr.Code)
	got, err := store.GetDesiredMappings(context.Background(), "ws-sql")
	require.NoError(t, err)
	require.Equal(t, []DesiredMapping{old}, got)

	svc.SetRestoreMemoryHook(func([]byte) error { return errors.New("injected memory rollback failure") })
	require.NoError(t, svc.ApplySettings(context.Background(), Settings{
		Enabled: true, BindAddress: "127.0.0.1", Port: reservePort(t), SlaveID: 1,
		CapacityRegisters: 64, SettingsRevision: "set-1",
	}))
	_, err = reconciler.Reconcile(context.Background(), ReconcileRequest{WorkspaceID: "ws-sql", ExpectedWorkspaceRevision: first.NewWorkspaceRevision, ExpectedSettingsRevision: "set-1", DesiredMappings: nil})
	require.ErrorAs(t, err, &shareErr)
	require.Equal(t, ErrCodeDirtyUnknown, shareErr.Code)
	got, err = store.GetDesiredMappings(context.Background(), "ws-sql")
	require.NoError(t, err)
	require.Equal(t, []DesiredMapping{old}, got)

	err = svc.StartConfiguredListener(context.Background())
	require.ErrorAs(t, err, &shareErr)
	require.Equal(t, ErrCodeHydrationRequired, shareErr.Code)

	svc.SetHydrationState(HydrationState{State: HydrationStateReady, WorkspaceID: "ws-sql", Readiness: true})
	_, err = reconciler.Reconcile(context.Background(), ReconcileRequest{
		WorkspaceID:               "ws-sql",
		ExpectedWorkspaceRevision: first.NewWorkspaceRevision,
		ExpectedSettingsRevision:  "set-1",
		DesiredMappings:           nil,
	})
	require.ErrorAs(t, err, &shareErr)
	require.Equal(t, ErrCodeDirtyUnknown, shareErr.Code)
}
