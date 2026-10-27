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
)

const (
	mysqlCapabilityBoundConnector   = "conn-mysql-bound"
	mysqlCapabilityBoundRevision    = "identity-rev-1"
	mysqlCapabilityForeignConnector = "conn-mysql-foreign"
	mysqlCapabilityCapabilitiesPath = "/capabilities"
	mysqlCapabilitySavedScopeSuffix = "&connector_id=" + mysqlCapabilityBoundConnector + "&expected_connector_revision=" + mysqlCapabilityBoundRevision
)

// fakeMySQLCapabilityResolver implements only the legacy saved-connector
// resolver. The capability probe is an optional narrow interface, so this fake
// also pins the safe "probe unavailable" behavior.
type fakeMySQLCapabilityResolver struct {
	connector        *schema.DatabaseConnector
	resolveErr       error
	resolvedID       string
	resolvedRevision string
	resolveCalls     int
	inspectCalls     int
}

func (f *fakeMySQLCapabilityResolver) ResolveSavedTarget(_ context.Context, connectorID, expectedRevision string) (*schema.DatabaseConnector, error) {
	f.resolveCalls++
	f.resolvedID, f.resolvedRevision = connectorID, expectedRevision
	if f.resolveErr != nil {
		return nil, f.resolveErr
	}
	return f.connector, nil
}

func (f *fakeMySQLCapabilityResolver) InspectTable(context.Context, string, string, string) (*dbtarget.TableInspection, error) {
	f.inspectCalls++
	return nil, nil
}

func (f *fakeMySQLCapabilityResolver) ExecuteSchemaStatements(context.Context, string, string, []string) (*dbtarget.SchemaStatementExecution, error) {
	return nil, nil
}

// fakeMySQLCapabilityProbeResolver adds the optional narrow probe without
// changing the legacy resolver interface every other fake already implements.
type fakeMySQLCapabilityProbeResolver struct {
	*fakeMySQLCapabilityResolver
	capability    *dbtarget.MySQLRecordingCapability
	probeErr      error
	probeID       string
	probeRevision string
	probeCalls    int
}

func (f *fakeMySQLCapabilityProbeResolver) ProbeMySQLRecordingAtRevision(_ context.Context, connectorID, expectedRevision string) (*dbtarget.MySQLRecordingCapability, error) {
	f.probeCalls++
	f.probeID, f.probeRevision = connectorID, expectedRevision
	if f.probeErr != nil {
		return nil, f.probeErr
	}
	return f.capability, nil
}

func newMySQLCapabilityProbeRouter(t *testing.T, capability *dbtarget.MySQLRecordingCapability, probeErr error) (*gin.Engine, *fakeMySQLCapabilityProbeResolver) {
	t.Helper()
	probe := &fakeMySQLCapabilityProbeResolver{
		fakeMySQLCapabilityResolver: &fakeMySQLCapabilityResolver{connector: &schema.DatabaseConnector{
			ID: mysqlCapabilityBoundConnector, Kind: schema.DatabaseConnectorKindMySQL,
			IdentityRevision: mysqlCapabilityBoundRevision, Enabled: true,
		}},
		capability: capability,
		probeErr:   probeErr,
	}
	router, _ := newMySQLCapabilityFixture(t, probe)
	return router, probe
}

type mysqlCapabilitySuccessEnvelope struct {
	Success bool `json:"success"`
	Data    struct {
		recordingplan.ConnectorCapability
		dbtarget.MySQLRecordingCapability
	} `json:"data"`
}

func decodeMySQLCapabilitySuccess(t *testing.T, response *httptest.ResponseRecorder) mysqlCapabilitySuccessEnvelope {
	t.Helper()
	var envelope mysqlCapabilitySuccessEnvelope
	if err := json.Unmarshal(response.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode capability success: %v; body=%s", err, response.Body.String())
	}
	if !envelope.Success {
		t.Fatalf("capability success envelope must succeed: %s", response.Body.String())
	}
	return envelope
}

type mysqlCapabilityErrorEnvelope struct {
	Success bool `json:"success"`
	Error   struct {
		Code      string `json:"code"`
		Message   string `json:"message"`
		Retryable bool   `json:"retryable"`
		Action    string `json:"action"`
	} `json:"error"`
}

func newMySQLCapabilityFixture(t *testing.T, resolver recordingSavedConnectorResolver) (*gin.Engine, *fakeMySQLCapabilityResolver) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	record, err := workspaceSvc.GetOrCreate(t.Context())
	if err != nil {
		t.Fatalf("create workspace: %v", err)
	}
	if _, err := workspaceSvc.BindDatabaseConnector(t.Context(), mysqlCapabilityBoundConnector); err != nil {
		t.Fatalf("bind workspace connector: %v", err)
	}
	_ = record

	planSvc := recordingplan.NewService(recordingplan.NewMemoryRepository())
	handler := NewStudioV2WorkspaceRecordingPlansHandler(workspaceSvc, planSvc, resolver)
	router := gin.New()
	router.GET(mysqlCapabilityCapabilitiesPath, handler.Capabilities)

	fake, _ := resolver.(*fakeMySQLCapabilityResolver)
	return router, fake
}

func getMySQLCapability(t *testing.T, router *gin.Engine, query string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, mysqlCapabilityCapabilitiesPath+query, http.NoBody)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func decodeMySQLCapabilityError(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()
	var envelope mysqlCapabilityErrorEnvelope
	if err := json.Unmarshal(response.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode capability error: %v; body=%s", err, response.Body.String())
	}
	if envelope.Success {
		t.Fatalf("capability error response must not succeed: %s", response.Body.String())
	}
	if strings.TrimSpace(envelope.Error.Code) == "" || strings.TrimSpace(envelope.Error.Message) == "" {
		t.Fatalf("capability error must be typed and safe: %s", response.Body.String())
	}
}

func assertMySQLCapabilityBodySafe(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()
	body := response.Body.String()
	for _, leaked := range []string{mysqlCapabilityForeignConnector, "dsn", "password", "127.0.0.1"} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(leaked)) {
			t.Fatalf("capability response leaked %q: %s", leaked, body)
		}
	}
}

// A saved-scope query is all or nothing: any saved-scope parameter commits the
// caller to kind=mysql, connector_id and expected_connector_revision, and an
// incomplete query is a typed 422 before any saved connector is resolved.
func TestMySQLRecordingCapabilityRequiresCompleteSavedScope(t *testing.T) {
	fake := &fakeMySQLCapabilityResolver{}
	router, decodedFake := newMySQLCapabilityFixture(t, fake)
	if decodedFake == nil {
		t.Fatal("fixture must keep the fake resolver")
	}

	cases := []struct {
		name  string
		query string
	}{
		{name: "connector without revision", query: "?kind=mysql&connector_id=" + mysqlCapabilityBoundConnector},
		{name: "revision without connector", query: "?kind=mysql&expected_connector_revision=" + mysqlCapabilityBoundRevision},
		{name: "saved scope without kind", query: "?connector_id=" + mysqlCapabilityBoundConnector + "&expected_connector_revision=" + mysqlCapabilityBoundRevision},
		{name: "empty kind with saved scope", query: "?kind=&connector_id=" + mysqlCapabilityBoundConnector + "&expected_connector_revision=" + mysqlCapabilityBoundRevision},
		// The keys themselves are the saved-scope signal. An empty or whitespace
		// value must be rejected, not treated as an absent saved scope.
		{name: "empty saved scope values", query: "?kind=mysql&connector_id=&expected_connector_revision="},
		{name: "whitespace saved scope values", query: "?kind=mysql&connector_id=%20&expected_connector_revision=%20"},
		{name: "saved scope keys without values", query: "?kind=mysql&connector_id&expected_connector_revision"},
		{name: "saved scope keys without kind", query: "?connector_id=&expected_connector_revision="},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			response := getMySQLCapability(t, router, testCase.query)
			if response.Code != http.StatusUnprocessableEntity {
				t.Fatalf("incomplete saved scope must be 422, got %d: %s", response.Code, response.Body.String())
			}
			decodeMySQLCapabilityError(t, response)
			assertMySQLCapabilityBodySafe(t, response)
		})
	}
	if fake.resolveCalls != 0 || fake.inspectCalls != 0 {
		t.Fatalf("incomplete scope must not resolve or inspect a connector: %+v", fake)
	}
}

func TestMySQLRecordingCapabilityRejectsNonMySQLSavedScope(t *testing.T) {
	fake := &fakeMySQLCapabilityResolver{connector: &schema.DatabaseConnector{
		ID: mysqlCapabilityBoundConnector, Kind: schema.DatabaseConnectorKindPostgres,
		IdentityRevision: mysqlCapabilityBoundRevision, Enabled: true,
	}}
	router, _ := newMySQLCapabilityFixture(t, fake)

	response := getMySQLCapability(t, router, "?kind=postgres&connector_id="+mysqlCapabilityBoundConnector+"&expected_connector_revision="+mysqlCapabilityBoundRevision)
	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("non-mysql saved scope must be 422, got %d: %s", response.Code, response.Body.String())
	}
	decodeMySQLCapabilityError(t, response)
	assertMySQLCapabilityBodySafe(t, response)
	if fake.resolveCalls != 0 {
		t.Fatalf("kind mismatch must be refused before resolving a connector: %+v", fake)
	}
}

func TestMySQLRecordingCapabilityRejectsUnboundConnector(t *testing.T) {
	fake := &fakeMySQLCapabilityResolver{connector: &schema.DatabaseConnector{
		ID: mysqlCapabilityForeignConnector, Kind: schema.DatabaseConnectorKindMySQL,
		IdentityRevision: mysqlCapabilityBoundRevision, Enabled: true,
	}}
	router, _ := newMySQLCapabilityFixture(t, fake)

	response := getMySQLCapability(t, router, "?kind=mysql&connector_id="+mysqlCapabilityForeignConnector+"&expected_connector_revision="+mysqlCapabilityBoundRevision)
	if response.Code != http.StatusNotFound {
		t.Fatalf("a connector outside the workspace binding must be a safe 404, got %d: %s", response.Code, response.Body.String())
	}
	decodeMySQLCapabilityError(t, response)
	assertMySQLCapabilityBodySafe(t, response)
	if fake.resolveCalls != 0 {
		t.Fatalf("an unbound connector must be refused before resolution: %+v", fake)
	}
}

func TestMySQLRecordingCapabilityRejectsStaleSavedRevision(t *testing.T) {
	fake := &fakeMySQLCapabilityResolver{connector: &schema.DatabaseConnector{
		ID: mysqlCapabilityBoundConnector, Kind: schema.DatabaseConnectorKindMySQL,
		IdentityRevision: "identity-rev-2", Enabled: true,
	}}
	router, _ := newMySQLCapabilityFixture(t, fake)

	response := getMySQLCapability(t, router, "?kind=mysql"+mysqlCapabilitySavedScopeSuffix)
	if response.Code != http.StatusConflict {
		t.Fatalf("a stale saved revision must be 409, got %d: %s", response.Code, response.Body.String())
	}
	decodeMySQLCapabilityError(t, response)
	assertMySQLCapabilityBodySafe(t, response)
	if fake.inspectCalls != 0 {
		t.Fatalf("stale revision must not inspect the target: %+v", fake)
	}
}

func TestMySQLRecordingCapabilityWithoutProbeIsUnavailable(t *testing.T) {
	fake := &fakeMySQLCapabilityResolver{connector: &schema.DatabaseConnector{
		ID: mysqlCapabilityBoundConnector, Kind: schema.DatabaseConnectorKindMySQL,
		IdentityRevision: mysqlCapabilityBoundRevision, Enabled: true,
	}}
	router, _ := newMySQLCapabilityFixture(t, fake)

	response := getMySQLCapability(t, router, "?kind=mysql"+mysqlCapabilitySavedScopeSuffix)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("a resolver without the readonly probe must be a safe 503, got %d: %s", response.Code, response.Body.String())
	}
	decodeMySQLCapabilityError(t, response)
	assertMySQLCapabilityBodySafe(t, response)
}

func TestMySQLRecordingCapabilityWithoutResolverIsUnavailable(t *testing.T) {
	router, _ := newMySQLCapabilityFixture(t, nil)

	response := getMySQLCapability(t, router, "?kind=mysql"+mysqlCapabilitySavedScopeSuffix)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("a missing resolver must be a safe 503, got %d: %s", response.Code, response.Body.String())
	}
	decodeMySQLCapabilityError(t, response)
	assertMySQLCapabilityBodySafe(t, response)
}

func TestMySQLRecordingCapabilityProbeProjectsServerFacts(t *testing.T) {
	router, probe := newMySQLCapabilityProbeRouter(t, &dbtarget.MySQLRecordingCapability{
		ConnectorID: mysqlCapabilityBoundConnector, ConnectorRevision: mysqlCapabilityBoundRevision, Database: "gw",
		ServerFamily: "mysql", ServerVersion: "8.4.5", CapabilityRevision: "server-probe-1",
		Connected: true, MetadataStatus: "readable", TransactionalEngine: true,
		Reasons: []string{"unverified"},
	}, nil)

	response := getMySQLCapability(t, router, "?kind=mysql"+mysqlCapabilitySavedScopeSuffix+
		"&server_family=mariadb&server_version=11.4.0&capability_revision=client-claim"+
		"&connected=false&supports_managed_long=true&supports_native_partition=true"+
		"&supports_group_test_write=true&supports_missing_repair=true")
	if response.Code != http.StatusOK {
		t.Fatalf("saved-scope capability must succeed, got %d: %s", response.Code, response.Body.String())
	}
	data := decodeMySQLCapabilitySuccess(t, response).Data
	if data.Kind != "mysql" || !data.Supported || data.SupportsManagedSchema || data.SupportsGroupTestWrites {
		t.Fatalf("static mysql capability must stay masked: %+v", data.ConnectorCapability)
	}
	if data.ConnectorID != mysqlCapabilityBoundConnector || data.ConnectorRevision != mysqlCapabilityBoundRevision {
		t.Fatalf("probe result must bind the saved scope: %+v", data.MySQLRecordingCapability)
	}
	if data.ServerFamily != "mysql" || data.ServerVersion != "8.4.5" || data.CapabilityRevision != "server-probe-1" {
		t.Fatalf("server facts must come from the probe, not the client: %+v", data.MySQLRecordingCapability)
	}
	if !data.Connected || data.MetadataStatus != "readable" || !data.TransactionalEngine {
		t.Fatalf("probe connection/metadata facts were lost: %+v", data.MySQLRecordingCapability)
	}
	if data.SupportsManagedLong || data.SupportsNativePartition || data.SupportsGroupTestWrite || data.SupportsMissingRepair {
		t.Fatalf("unverified operations must stay false: %+v", data.MySQLRecordingCapability)
	}
	if probe.probeCalls != 1 || probe.probeID != mysqlCapabilityBoundConnector || probe.probeRevision != mysqlCapabilityBoundRevision {
		t.Fatalf("probe must run once on the saved revision: %+v", probe)
	}
	if probe.resolveCalls != 1 || probe.inspectCalls != 0 {
		t.Fatalf("probe path must resolve once and never inspect/DDL: %+v", probe)
	}
}

func TestMySQLRecordingCapabilityForbiddenMetadataStaysConnected(t *testing.T) {
	router, _ := newMySQLCapabilityProbeRouter(t, &dbtarget.MySQLRecordingCapability{
		ConnectorID: mysqlCapabilityBoundConnector, ConnectorRevision: mysqlCapabilityBoundRevision, Database: "gw",
		ServerFamily: "mysql", ServerVersion: "8.4.5", CapabilityRevision: "server-probe-1",
		Connected: true, MetadataStatus: "forbidden", TransactionalEngine: true,
		Reasons: []string{"metadata-forbidden"},
	}, nil)

	response := getMySQLCapability(t, router, "?kind=mysql"+mysqlCapabilitySavedScopeSuffix)
	if response.Code != http.StatusOK {
		t.Fatalf("a metadata refusal is not a connection failure: %d %s", response.Code, response.Body.String())
	}
	data := decodeMySQLCapabilitySuccess(t, response).Data
	if !data.Connected || data.MetadataStatus != "forbidden" {
		t.Fatalf("connection success and metadata refusal must both be visible: %+v", data.MySQLRecordingCapability)
	}
	if data.SupportsManagedLong || data.SupportsNativePartition || data.SupportsGroupTestWrite || data.SupportsMissingRepair {
		t.Fatalf("a metadata refusal must block every dependent operation: %+v", data.MySQLRecordingCapability)
	}
	if data.SupportsManagedSchema || data.SupportsGroupTestWrites {
		t.Fatalf("static mutation gates must stay masked: %+v", data.ConnectorCapability)
	}
	assertMySQLCapabilityBodySafe(t, response)
}

func TestMySQLRecordingCapabilityProbeFailuresStaySafe(t *testing.T) {
	cases := []struct {
		name       string
		probeErr   error
		wantStatus int
		leak       string
	}{
		{name: "stale revision", probeErr: dbtarget.ErrConnectorRevisionConflict, wantStatus: http.StatusConflict},
		{name: "raw transport", probeErr: errors.New("dial tcp 10.10.10.10:3306: connect: connection refused"), wantStatus: http.StatusServiceUnavailable, leak: "10.10.10.10"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			router, _ := newMySQLCapabilityProbeRouter(t, nil, testCase.probeErr)
			response := getMySQLCapability(t, router, "?kind=mysql"+mysqlCapabilitySavedScopeSuffix)
			if response.Code != testCase.wantStatus {
				t.Fatalf("probe failure must map to %d, got %d: %s", testCase.wantStatus, response.Code, response.Body.String())
			}
			decodeMySQLCapabilityError(t, response)
			assertMySQLCapabilityBodySafe(t, response)
			if testCase.leak != "" && strings.Contains(response.Body.String(), testCase.leak) {
				t.Fatalf("probe failure leaked raw detail: %s", response.Body.String())
			}
		})
	}
}

func TestMySQLRecordingCapabilityProbeMustBindSavedScope(t *testing.T) {
	cases := []struct {
		name       string
		capability *dbtarget.MySQLRecordingCapability
		wantStatus int
	}{
		{name: "nil result", capability: nil, wantStatus: http.StatusServiceUnavailable},
		{name: "wrong connector", capability: &dbtarget.MySQLRecordingCapability{ConnectorID: mysqlCapabilityForeignConnector, ConnectorRevision: mysqlCapabilityBoundRevision}, wantStatus: http.StatusServiceUnavailable},
		{name: "wrong revision", capability: &dbtarget.MySQLRecordingCapability{ConnectorID: mysqlCapabilityBoundConnector, ConnectorRevision: "identity-rev-2"}, wantStatus: http.StatusServiceUnavailable},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			router, _ := newMySQLCapabilityProbeRouter(t, testCase.capability, nil)
			response := getMySQLCapability(t, router, "?kind=mysql"+mysqlCapabilitySavedScopeSuffix)
			if response.Code != testCase.wantStatus {
				t.Fatalf("an unbound probe result must not succeed, got %d: %s", response.Code, response.Body.String())
			}
			decodeMySQLCapabilityError(t, response)
			assertMySQLCapabilityBodySafe(t, response)
		})
	}
}

// Legacy queries stay byte-compatible: no saved-scope parameters means the
// static kind matrix is returned exactly as before.
func TestMySQLRecordingCapabilityLegacyQueriesRemainStatic(t *testing.T) {
	router, _ := newMySQLCapabilityFixture(t, &fakeMySQLCapabilityResolver{})

	single := getMySQLCapability(t, router, "?kind=mysql")
	if single.Code != http.StatusOK {
		t.Fatalf("legacy kind=mysql query failed: %d %s", single.Code, single.Body.String())
	}
	var singleBody struct {
		Success bool                              `json:"success"`
		Data    recordingplan.ConnectorCapability `json:"data"`
	}
	if err := json.Unmarshal(single.Body.Bytes(), &singleBody); err != nil {
		t.Fatalf("decode legacy kind=mysql: %v; body=%s", err, single.Body.String())
	}
	if !singleBody.Success || singleBody.Data.Kind != "mysql" || singleBody.Data.SupportsManagedSchema || singleBody.Data.SupportsGroupTestWrites {
		t.Fatalf("legacy kind=mysql must stay masked: %+v", singleBody)
	}

	list := getMySQLCapability(t, router, "")
	if list.Code != http.StatusOK {
		t.Fatalf("legacy capability list failed: %d %s", list.Code, list.Body.String())
	}
	var listBody struct {
		Success bool                                `json:"success"`
		Data    []recordingplan.ConnectorCapability `json:"data"`
	}
	if err := json.Unmarshal(list.Body.Bytes(), &listBody); err != nil {
		t.Fatalf("decode legacy capability list: %v; body=%s", err, list.Body.String())
	}
	if !listBody.Success || len(listBody.Data) != 5 {
		t.Fatalf("legacy capability list changed: %s", list.Body.String())
	}
}
