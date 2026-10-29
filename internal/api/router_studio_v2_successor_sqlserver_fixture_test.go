package api

import (
	"database/sql"
	"encoding/json"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

const sqlServerSuccessorFixtureBEnv = "GW_SQLSERVER_TEST_FIXTURE_B"

// Both opt-in manifests reuse the circuit fixture contract. No database/schema,
// server, container, TLS bypass or fabricated capability is created here.
func sqlServerSuccessorFixtures(t *testing.T) [2]*SQLServerCircuitFixture {
	t.Helper()
	a := SQLServerCircuitFixtureOrSkip(t)
	pathB := strings.TrimSpace(os.Getenv(sqlServerSuccessorFixtureBEnv))
	if pathB == "" {
		t.Skip("NOT_RUN: GW_SQLSERVER_TEST_FIXTURE_B is not set; real A/B acceptance remains pending")
	}
	pathA := os.Getenv(sqlServerCircuitFixtureEnv)
	t.Setenv(sqlServerCircuitFixtureEnv, pathB)
	b := SQLServerCircuitFixtureOrSkip(t)
	t.Setenv(sqlServerCircuitFixtureEnv, pathA)
	require.NotEqual(t, strings.ToLower(a.Manifest.Database), strings.ToLower(b.Manifest.Database), "A/B must be distinct disposable databases")
	var endpoint string
	for _, fixture := range []*SQLServerCircuitFixture{a, b} {
		var config map[string]any
		require.NoError(t, json.Unmarshal([]byte(fixture.ConnectionConfigJSON()), &config))
		host := strings.ToLower(strings.Trim(strings.TrimSpace(stringConfigValue(config, "host")), "[]"))
		ip := net.ParseIP(host)
		require.True(t, host == "localhost" || ip != nil && ip.IsLoopback(), "fixture must stay on the authorized existing computer")
		port := strings.TrimSpace(stringConfigValue(config, "port"))
		if port == "" {
			port = "1433"
		}
		current := net.JoinHostPort(host, port)
		if endpoint != "" {
			require.Equal(t, endpoint, current, "both databases must be on the same existing server")
		}
		endpoint = current
		// Parsing is offline; this existing helper rejects disabled verification.
		validatedTLS, _ := sqlServerCircuitFaultDB(t, fixture)
		require.NoError(t, validatedTLS.Close())
	}
	return [2]*SQLServerCircuitFixture{a, b}
}

func newSQLServerSuccessorDestination(t *testing.T, mode workspace.WriteGroupStorageMode) (managedGroupFixture, *sql.DB, [2]*SQLServerCircuitFixture) {
	t.Helper()
	fixtures := sqlServerSuccessorFixtures(t)
	f := managedGroupFixture{writeGroupRouterFixture: newWriteGroupRouterFixture(t)}
	for i, fixture := range fixtures {
		var err error
		if i == 0 {
			_, err = f.db.Exec(`UPDATE database_connectors SET kind='sqlserver',connection_config=? WHERE id='connector-A'`, fixture.ConnectionConfigJSON())
		} else {
			_, err = f.db.Exec(`INSERT INTO database_connectors(id,name,kind,connection_config,identity_revision,status,enabled) VALUES('connector-B','Other','sqlserver',?,'connector-2','ready',1)`, fixture.ConnectionConfigJSON())
		}
		require.NoError(t, err)
	}
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	var destinations [2]*sql.DB
	for i, fixture := range fixtures {
		id, revision := []string{"connector-A", "connector-B"}[i], []string{"connector-1", "connector-2"}[i]
		capability, err := targets.ProbeSQLServerRecordingAtRevision(t.Context(), id, revision)
		require.NoError(t, err)
		require.NoError(t, fixture.ValidateCapability(capability))
		checkSQLServerCircuitGate(t, mode, capability)
		opened, err := targets.OpenDestination(t.Context(), id, revision)
		require.NoError(t, err)
		if opened.Close != nil {
			t.Cleanup(func() { require.NoError(t, opened.Close()) })
		}
		destinations[i] = opened.DB
		var otherTables int
		require.NoError(t, opened.DB.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM sys.tables WHERE schema_id=SCHEMA_ID(@p1) AND name<>@p2`, fixture.Manifest.Schema, dbtarget.EffectReceiptTable).Scan(&otherTables))
		require.Zero(t, otherTables, "use an otherwise empty disposable scope; shared receipts are preserved")
	}
	f.target, f.namespace = destinations[0], fixtures[0].Manifest.Schema
	var err error
	f.record, err = f.workspace.BindDatabaseConnector(t.Context(), "connector-A")
	require.NoError(t, err)
	_, err = f.db.Exec(`UPDATE tags SET data_type='uint64' WHERE id='tag-A'`)
	require.NoError(t, err)
	f.groups = workspace.NewWriteGroupService(f.workspace, workspace.NewSQLWriteGroupRepository(f.db)).WithTableInspector(dbtarget.NewReadOnlyTableInspector(targets)).WithBacklogOwnershipGuard(grouppipeline.BacklogGuard{})
	f.groups.WithManagedTableInspector(f.groups.ManagedTableInspector(targets))
	f.router = NewRouter(&DatalinkServices{Workspace: f.workspace, WriteGroups: f.groups, DBTarget: targets, RecordingPlan: recordingplan.NewService(recordingplan.NewSQLRepository(f.db))})
	return f, destinations[1], fixtures
}
