package dbtarget

import (
	"context"
	"errors"
	"testing"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type managedSQLServerProbeForTest interface {
	ProbeSQLServerRecordingAtRevision(context.Context, string, string) (*SQLServerRecordingCapability, error)
}

func TestSQLServerManagedInspectorForwardsSavedRevisionProbe(t *testing.T) {
	internalDB := openMigratedTestDB(t)
	repo := NewSQLConnectorRepository(internalDB)
	service := NewConnectorService(repo)
	config := sqlServerConnectorConfigForTest()
	connector := createSQLServerRecordingSavedConnector(t, repo, "managed-sqlserver-probe", schema.DatabaseConnectorKindSQLServer, config, true)
	inspector := NewManagedTableInspector(service, internalDB)
	prober, ok := any(inspector).(managedSQLServerProbeForTest)
	require.True(t, ok, "production managed inspector must implement the workspace SQL Server probe contract")

	original := openExternalDBManagerFunc
	openCalls := 0
	openExternalDBManagerFunc = func(kind schema.DatabaseConnectorKind, saved ConnectionConfig) (*datalinkbase.DBManager, error) {
		openCalls++
		require.Equal(t, schema.DatabaseConnectorKindSQLServer, kind)
		require.Equal(t, config["password"], saved["password"])
		return nil, errors.New("disposable fixture unavailable")
	}
	t.Cleanup(func() { openExternalDBManagerFunc = original })
	capability, err := prober.ProbeSQLServerRecordingAtRevision(t.Context(), connector.ID, connector.IdentityRevision)
	require.NoError(t, err)
	require.Equal(t, connector.ID, capability.ConnectorID)
	require.Equal(t, connector.IdentityRevision, capability.ConnectorRevision)
	require.False(t, capability.Connected)
	require.False(t, capability.SupportsNativePartition)
	require.Equal(t, 1, openCalls)

	_, err = prober.ProbeSQLServerRecordingAtRevision(t.Context(), connector.ID, "stale")
	require.ErrorIs(t, err, ErrConnectorRevisionConflict)
	require.Equal(t, 1, openCalls, "stale probe must not open the destination")
}

func TestSQLServerManagedInspectorProbeRequiresSavedService(t *testing.T) {
	for _, inspector := range []*ManagedTableInspector{nil, {}} {
		prober, ok := any(inspector).(managedSQLServerProbeForTest)
		require.True(t, ok)
		_, err := prober.ProbeSQLServerRecordingAtRevision(t.Context(), "connector", "revision")
		require.ErrorIs(t, err, ErrManagedDestinationInvalid)
	}
}
