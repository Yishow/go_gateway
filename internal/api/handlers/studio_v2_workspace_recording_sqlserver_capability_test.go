package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

const (
	sqlServerCapabilityBoundConnector   = "conn-sqlserver-bound"
	sqlServerCapabilityBoundRevision    = "identity-rev-ss-1"
	sqlServerCapabilityForeignConnector = "conn-sqlserver-foreign"
	sqlServerCapabilityCapabilitiesPath = "/capabilities"
)

// fakeSQLServerCapabilityResolver implements only the legacy saved-connector
// resolver. The probe is an optional narrow interface, so this fake pins the
// safe "resolver unavailable" behavior.
type fakeSQLServerCapabilityResolver struct {
	connector        *schema.DatabaseConnector
	resolveErr       error
	resolvedID       string
	resolvedRevision string
	resolveCalls     int
}

func (f *fakeSQLServerCapabilityResolver) ResolveSavedTarget(_ context.Context, connectorID, expectedRevision string) (*schema.DatabaseConnector, error) {
	f.resolveCalls++
	f.resolvedID, f.resolvedRevision = connectorID, expectedRevision
	if f.resolveErr != nil {
		return nil, f.resolveErr
	}
	return f.connector, nil
}

func (f *fakeSQLServerCapabilityResolver) InspectTable(context.Context, string, string, string) (*dbtarget.TableInspection, error) {
	return nil, nil
}

func (f *fakeSQLServerCapabilityResolver) ExecuteSchemaStatements(context.Context, string, string, []string) (*dbtarget.SchemaStatementExecution, error) {
	return nil, nil
}

// fakeSQLServerCapabilityProbeResolver adds the optional narrow probe without
// changing the legacy resolver interface every other fake already implements.
type fakeSQLServerCapabilityProbeResolver struct {
	*fakeSQLServerCapabilityResolver
	capability    *dbtarget.SQLServerRecordingCapability
	probeErr      error
	probeID       string
	probeRevision string
	probeCalls    int
}

func (f *fakeSQLServerCapabilityProbeResolver) ProbeSQLServerRecordingAtRevision(
	_ context.Context, connectorID, expectedRevision string,
) (*dbtarget.SQLServerRecordingCapability, error) {
	f.probeCalls++
	f.probeID, f.probeRevision = connectorID, expectedRevision
	if f.probeErr != nil {
		return nil, f.probeErr
	}
	return f.capability, nil
}

func newSQLServerCapabilityFixture(t *testing.T, resolver recordingSavedConnectorResolver) (*gin.Engine, *fakeSQLServerCapabilityResolver) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	if _, err := workspaceSvc.BindDatabaseConnector(t.Context(), sqlServerCapabilityBoundConnector); err != nil {
		t.Fatalf("bind workspace connector: %v", err)
	}

	planSvc := recordingplan.NewService(recordingplan.NewMemoryRepository())
	handler := NewStudioV2WorkspaceRecordingPlansHandler(workspaceSvc, planSvc, resolver)
	router := gin.New()
	router.GET(sqlServerCapabilityCapabilitiesPath, handler.Capabilities)

	fake, _ := resolver.(*fakeSQLServerCapabilityResolver)
	return router, fake
}

func newSQLServerCapabilityProbeRouter(
	t *testing.T, capability *dbtarget.SQLServerRecordingCapability, probeErr error,
) (*gin.Engine, *fakeSQLServerCapabilityProbeResolver) {
	t.Helper()
	probe := &fakeSQLServerCapabilityProbeResolver{
		fakeSQLServerCapabilityResolver: &fakeSQLServerCapabilityResolver{connector: &schema.DatabaseConnector{
			ID: sqlServerCapabilityBoundConnector, Kind: schema.DatabaseConnectorKindSQLServer,
			IdentityRevision: sqlServerCapabilityBoundRevision, Enabled: true,
		}},
		capability: capability,
		probeErr:   probeErr,
	}
	router, _ := newSQLServerCapabilityFixture(t, probe)
	return router, probe
}

func getSQLServerCapability(t *testing.T, router *gin.Engine, query string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, sqlServerCapabilityCapabilitiesPath+query, http.NoBody)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

type sqlServerCapabilitySuccessEnvelope struct {
	Success bool `json:"success"`
	Data    struct {
		recordingplan.ConnectorCapability
		dbtarget.SQLServerRecordingCapability
	} `json:"data"`
}

func decodeSQLServerCapabilitySuccess(t *testing.T, response *httptest.ResponseRecorder) sqlServerCapabilitySuccessEnvelope {
	t.Helper()
	var envelope sqlServerCapabilitySuccessEnvelope
	if err := json.Unmarshal(response.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode capability success: %v; body=%s", err, response.Body.String())
	}
	if !envelope.Success {
		t.Fatalf("capability success envelope must succeed: %s", response.Body.String())
	}
	return envelope
}

type sqlServerErrorEnvelope struct {
	Success bool `json:"success"`
	Error   struct {
		Code      string `json:"code"`
		Message   string `json:"message"`
		Retryable bool   `json:"retryable"`
		Action    string `json:"action"`
	} `json:"error"`
}

func decodeSQLServerCapabilityError(t *testing.T, response *httptest.ResponseRecorder) sqlServerErrorEnvelope {
	t.Helper()
	var envelope sqlServerErrorEnvelope
	if err := json.Unmarshal(response.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode capability error: %v; body=%s", err, response.Body.String())
	}
	if envelope.Success {
		t.Fatalf("capability error response must not succeed: %s", response.Body.String())
	}
	if strings.TrimSpace(envelope.Error.Code) == "" || strings.TrimSpace(envelope.Error.Message) == "" {
		t.Fatalf("capability error must be typed and safe: %s", response.Body.String())
	}
	return envelope
}

// TestSQLServerRecordingCapabilityRequiresCompleteSavedScope pins the strict
// saved-scope branch: any scope key commits the caller, and an incomplete or
// empty scope is a typed 422 before any connector is resolved. The bare kind
// stays on the static unsupported matrix.
func TestSQLServerRecordingCapabilityRequiresCompleteSavedScope(t *testing.T) {
	fake := &fakeSQLServerCapabilityResolver{}
	router, decoded := newSQLServerCapabilityFixture(t, fake)
	require.NotNil(t, decoded, "fixture must keep the fake resolver")

	cases := []struct{ name, query string }{
		{name: "connector without revision", query: "?kind=sqlserver&connector_id=" + sqlServerCapabilityBoundConnector},
		{name: "revision without connector", query: "?kind=sqlserver&expected_connector_revision=" + sqlServerCapabilityBoundRevision},
		{name: "empty saved scope values", query: "?kind=sqlserver&connector_id=&expected_connector_revision="},
		{name: "whitespace saved scope values", query: "?kind=sqlserver&connector_id=%20&expected_connector_revision=%20"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			response := getSQLServerCapability(t, router, tc.query)
			require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
			envelope := decodeSQLServerCapabilityError(t, response)
			require.Equal(t, "RECORDING_CONNECTOR_SCOPE_INVALID", envelope.Error.Code)
			require.Zero(t, decoded.resolveCalls, "an incomplete scope must not resolve any connector")
		})
	}

	bare := getSQLServerCapability(t, router, "?kind=sqlserver")
	require.Equal(t, http.StatusOK, bare.Code, bare.Body.String())
	envelope := decodeSQLServerCapabilitySuccess(t, bare)
	require.False(t, envelope.Data.SupportsManagedSchema,
		"the bare kind stays on the static unsupported matrix")
	require.Zero(t, decoded.resolveCalls)
}

// TestSQLServerRecordingCapabilityRejectsForeignScope pins the workspace
// binding and the sqlserver-only dialect inside the strict branch.
func TestSQLServerRecordingCapabilityRejectsForeignScope(t *testing.T) {
	probe := &fakeSQLServerCapabilityProbeResolver{
		fakeSQLServerCapabilityResolver: &fakeSQLServerCapabilityResolver{connector: &schema.DatabaseConnector{
			ID: sqlServerCapabilityBoundConnector, Kind: schema.DatabaseConnectorKindSQLServer,
			IdentityRevision: sqlServerCapabilityBoundRevision, Enabled: true,
		}},
	}
	router, _ := newSQLServerCapabilityFixture(t, probe)

	t.Run("foreign workspace connector", func(t *testing.T) {
		response := getSQLServerCapability(t, router,
			"?kind=sqlserver&connector_id="+sqlServerCapabilityForeignConnector+"&expected_connector_revision="+sqlServerCapabilityBoundRevision)
		require.Equal(t, http.StatusNotFound, response.Code, response.Body.String())
		require.Zero(t, probe.resolveCalls, "a foreign connector is refused before any resolve")
		require.Zero(t, probe.probeCalls, "a foreign connector is refused before any probe")
	})
	t.Run("non sqlserver kind with scope", func(t *testing.T) {
		response := getSQLServerCapability(t, router,
			"?kind=postgres&connector_id="+sqlServerCapabilityBoundConnector+"&expected_connector_revision="+sqlServerCapabilityBoundRevision)
		require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
		decodeSQLServerCapabilityError(t, response)
		require.Zero(t, probe.probeCalls)
	})
}

func TestSQLServerRecordingCapabilityRejectsStaleSavedRevision(t *testing.T) {
	probe := &fakeSQLServerCapabilityProbeResolver{
		fakeSQLServerCapabilityResolver: &fakeSQLServerCapabilityResolver{connector: &schema.DatabaseConnector{
			ID: sqlServerCapabilityBoundConnector, Kind: schema.DatabaseConnectorKindSQLServer,
			IdentityRevision: sqlServerCapabilityBoundRevision, Enabled: true,
		}},
	}
	router, _ := newSQLServerCapabilityFixture(t, probe)

	response := getSQLServerCapability(t, router,
		"?kind=sqlserver&connector_id="+sqlServerCapabilityBoundConnector+"&expected_connector_revision=stale-identity-rev")
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	decodeSQLServerCapabilityError(t, response)
	require.Equal(t, 1, probe.resolveCalls, "the stale revision is checked against the saved connector")
	require.Zero(t, probe.probeCalls, "a stale revision never reaches the probe")
}

func TestSQLServerRecordingCapabilityWithoutProbeIsUnavailable(t *testing.T) {
	router, _ := newSQLServerCapabilityFixture(t, &fakeSQLServerCapabilityResolver{connector: &schema.DatabaseConnector{
		ID: sqlServerCapabilityBoundConnector, Kind: schema.DatabaseConnectorKindSQLServer,
		IdentityRevision: sqlServerCapabilityBoundRevision, Enabled: true,
	}})
	response := getSQLServerCapability(t, router,
		"?kind=sqlserver&connector_id="+sqlServerCapabilityBoundConnector+"&expected_connector_revision="+sqlServerCapabilityBoundRevision)
	require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
	decodeSQLServerCapabilityError(t, response)
}

func TestSQLServerRecordingCapabilityRejectsDisabledConnector(t *testing.T) {
	probe := &fakeSQLServerCapabilityProbeResolver{
		fakeSQLServerCapabilityResolver: &fakeSQLServerCapabilityResolver{connector: &schema.DatabaseConnector{
			ID: sqlServerCapabilityBoundConnector, Kind: schema.DatabaseConnectorKindSQLServer,
			IdentityRevision: sqlServerCapabilityBoundRevision, Enabled: false,
		}},
	}
	router, _ := newSQLServerCapabilityFixture(t, probe)
	response := getSQLServerCapability(t, router,
		"?kind=sqlserver&connector_id="+sqlServerCapabilityBoundConnector+"&expected_connector_revision="+sqlServerCapabilityBoundRevision)
	require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
	decodeSQLServerCapabilityError(t, response)
	require.Zero(t, probe.probeCalls, "a disabled connector is refused before any probe")
}

// TestSQLServerRecordingCapabilityRejectsForeignProbeAnswer pins that a probe
// answer binding another connector is an adapter inconsistency, never a
// client-visible capability.
func TestSQLServerRecordingCapabilityRejectsForeignProbeAnswer(t *testing.T) {
	router, probe := newSQLServerCapabilityProbeRouter(t, &dbtarget.SQLServerRecordingCapability{
		ConnectorID: sqlServerCapabilityForeignConnector, ConnectorRevision: sqlServerCapabilityBoundRevision,
		Connected: true, MetadataStatus: "verified",
	}, nil)

	response := getSQLServerCapability(t, router,
		"?kind=sqlserver&connector_id="+sqlServerCapabilityBoundConnector+"&expected_connector_revision="+sqlServerCapabilityBoundRevision)
	require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
	decodeSQLServerCapabilityError(t, response)
	require.Equal(t, 1, probe.probeCalls)
}

// TestSQLServerRecordingCapabilityProbeErrorStaysTypedAndSafe keeps raw
// database or TLS failures out of the operator response.
func TestSQLServerRecordingCapabilityProbeErrorStaysTypedAndSafe(t *testing.T) {
	router, probe := newSQLServerCapabilityProbeRouter(t, nil, errors.New("tls: invalid certificate presented by 10.1.2.3"))
	response := getSQLServerCapability(t, router,
		"?kind=sqlserver&connector_id="+sqlServerCapabilityBoundConnector+"&expected_connector_revision="+sqlServerCapabilityBoundRevision)
	require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
	envelope := decodeSQLServerCapabilityError(t, response)
	require.NotContains(t, strings.ToLower(envelope.Error.Message), "certificate")
	require.NotContains(t, strings.ToLower(response.Body.String()), "10.1.2.3")
	require.Equal(t, 1, probe.probeCalls)
}

// TestMySQLSavedScopeStaysOnMySQLBranch protects the dispatch order: a MySQL
// saved-scope query is resolved by the MySQL branch (it reads the saved
// connector, then refuses the sqlserver-kind connector), and the sqlserver
// probe never serves it.
func TestMySQLSavedScopeStaysOnMySQLBranch(t *testing.T) {
	router, probe := newSQLServerCapabilityProbeRouter(t, &dbtarget.SQLServerRecordingCapability{
		ConnectorID: sqlServerCapabilityBoundConnector, ConnectorRevision: sqlServerCapabilityBoundRevision,
		Connected: true, MetadataStatus: "verified",
	}, nil)

	response := getSQLServerCapability(t, router,
		"?kind=mysql&connector_id="+sqlServerCapabilityBoundConnector+"&expected_connector_revision="+sqlServerCapabilityBoundRevision)
	require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
	decodeSQLServerCapabilityError(t, response)
	require.Equal(t, 1, probe.resolveCalls, "the MySQL branch resolves the saved connector")
	require.Zero(t, probe.probeCalls, "the sqlserver probe must not serve MySQL queries")
}

// TestSQLServerRecordingScopeRequestedOnlyForSQLServerKind pins the predicate
// gate: only kind=sqlserver with scope keys commits to the sqlserver branch.
func TestSQLServerRecordingScopeRequestedOnlyForSQLServerKind(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name, kind string
		want       bool
	}{
		{name: "mysql scope stays off", kind: "mysql", want: false},
		{name: "postgres scope stays off", kind: "postgres", want: false},
		{name: "sqlserver scope commits", kind: "sqlserver", want: true},
		{name: "case-insensitive sqlserver commits", kind: "SQLSERVER", want: true},
		{name: "bare kind without scope stays off", kind: "", want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequestWithContext(t.Context(), http.MethodGet,
				"/capabilities?kind="+tc.kind+"&connector_id=x&expected_connector_revision=y", http.NoBody)
			require.Equal(t, tc.want, sqlserverRecordingScopeRequested(c))
		})
	}
}

// TestSQLServerRecordingCapabilityFlattensStaticMatrixAndServerFacts proves a// TestSQLServerRecordingCapabilityFlattensStaticMatrixAndServerFacts proves a
// connected, metadata-verified probe still reports every operation flag false
// and never copies client-claimed server facts into the answer.
func TestSQLServerRecordingCapabilityFlattensStaticMatrixAndServerFacts(t *testing.T) {
	router, probe := newSQLServerCapabilityProbeRouter(t, &dbtarget.SQLServerRecordingCapability{
		ConnectorID: sqlServerCapabilityBoundConnector, ConnectorRevision: sqlServerCapabilityBoundRevision,
		Database: "inventory", Schema: "dbo", ServerVersion: "16.0.1000",
		ServerEdition: "Standard Edition", EngineEdition: 3, CompatibilityLevel: 160,
		Connected: true, MetadataStatus: "verified",
	}, nil)

	response := getSQLServerCapability(t, router,
		"?kind=sqlserver&connector_id="+sqlServerCapabilityBoundConnector+"&expected_connector_revision="+
			sqlServerCapabilityBoundRevision+"&server_version=9.9.9&supports_managed_long=true")
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	envelope := decodeSQLServerCapabilitySuccess(t, response)

	require.Equal(t, sqlServerCapabilityBoundConnector, envelope.Data.ConnectorID)
	require.Equal(t, sqlServerCapabilityBoundRevision, envelope.Data.ConnectorRevision)
	require.Equal(t, "16.0.1000", envelope.Data.ServerVersion, "client-claimed versions are ignored")
	require.True(t, envelope.Data.Connected)
	require.Equal(t, "verified", envelope.Data.MetadataStatus)
	require.False(t, envelope.Data.SupportsManagedLong)
	require.False(t, envelope.Data.SupportsNativePartition)
	require.False(t, envelope.Data.SupportsGroupTestWrite)
	require.False(t, envelope.Data.SupportsMissingRepair)
	require.False(t, envelope.Data.SupportsManagedSchema,
		"the static sqlserver kind matrix stays unsupported")
	require.Equal(t, 1, probe.probeCalls, "one probe per saved-scope query")
}
