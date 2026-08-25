package modbusshare

import (
	"context"
	"database/sql"
	"testing"

	"go-gateway/internal/datalink"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestService_StatusForWorkspaceScopesMappingsAndSharesLifecycleState(t *testing.T) {
	svc := NewService(nil, 128)
	svc.SetHydrationState(HydrationState{
		State:             HydrationStateReady,
		WorkspaceID:       "ws-a",
		WorkspaceRevision: "rev-a",
		SettingsRevision:  "settings-a",
		Readiness:         true,
	})
	svc.ReplaceMappings(map[string]TagMirrorMapping{
		"tag-a": {WorkspaceID: "ws-a", TagID: "tag-a"},
		"tag-b": {WorkspaceID: "ws-b", TagID: "tag-b"},
	})

	status := svc.StatusForWorkspace("ws-a")
	require.Equal(t, "ws-a", status.WorkspaceID)
	require.Equal(t, 1, status.MappingCount)
	require.Equal(t, status.LifecycleState, status.ListenerState)
}

func TestReconciler_RevisionStatusReadsDurableDirtyMarker(t *testing.T) {
	store := newMockWorkspaceRevisionStore()
	store.revisions["ws-a"] = "rev-a"
	store.dirty["ws-a"] = true
	reconciler := NewReconciler(NewService(nil, 128), store)

	status, err := reconciler.RevisionStatus(context.Background(), "ws-a")
	require.NoError(t, err)
	require.Equal(t, "rev-a", status.Revision)
	require.True(t, status.Dirty)
}

func TestReconciler_RevisionStatusReadsDirtyMarkerAfterStoreRestart(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	store := NewSQLWorkspaceRevisionStore(db)
	require.NoError(t, store.MarkDirty(context.Background(), "ws-restart"))

	restarted := NewReconciler(nil, NewSQLWorkspaceRevisionStore(db))
	status, err := restarted.RevisionStatus(context.Background(), "ws-restart")
	require.NoError(t, err)
	require.Equal(t, initialRevision, status.Revision)
	require.True(t, status.Dirty)
}
