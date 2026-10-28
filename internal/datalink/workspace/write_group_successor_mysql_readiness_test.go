package workspace

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccessorMySQLReadinessUsesMandatoryGuardedReload(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	connector := seedMySQLConnector(t, h.db, `{"dsn":"fixture:fixture-only@tcp(127.0.0.1:3306)/partition_test?tls=true"}`)
	_, err := h.workspace.BindDatabaseConnector(t.Context(), connector)
	require.NoError(t, err)
	group := h.longDraft(WriteGroupProfileInstantaneous, "guarded-reload", h.primary())
	group.Destination.ConnectorID, group.Destination.ConnectorRevision = connector, "connector-revision-mysql"
	saved := mustEnsureMySQLGroup(t, h, group)
	saved.Destination.SchemaRevision, saved.Destination.SchemaDigest = "confirmed", "proof"
	snapshot, err := h.service.readinessSnapshot(t.Context(), saved.ID, saved)
	require.NoError(t, err)
	h.service.WithManagedTableInspector(&fakeMySQLProbeInspector{err: errors.New("offline")})
	called := false
	stale := errors.New("reservation changed while probing")
	_, err = h.service.evaluateReadinessSnapshot(t.Context(), saved.ID, saved, snapshot, func(context.Context) (*writeGroupReadinessSnapshot, error) { called = true; return nil, stale })
	require.True(t, called, "MySQL must call the operation-bound reloader even when external probe fails")
	require.ErrorIs(t, err, stale)
}

func TestSuccessorMySQLActivationEngineGate(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		for _, mode := range []WriteGroupStorageMode{WriteGroupStorageModePlain, WriteGroupStorageModeNativeMonthly} {
			h := newCircuitProfileHarness(t, ":memory:")
			group := h.longDraft(WriteGroupProfileInstantaneous, "engine-gate", h.primary())
			group.StorageMode = mode
			group.Destination.StorageStrategy = WriteGroupStorageStrategyManaged
			require.True(t, successorActivationEngineSupported(group, family))
			group.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
			require.False(t, successorActivationEngineSupported(group, family))
		}
	}
}
