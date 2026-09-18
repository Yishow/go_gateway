package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type recordingTargetResolverStub struct {
	connector   *schema.DatabaseConnector
	err         error
	calls       []recordingTargetResolverCall
	inspections map[string]dbtarget.TableInspection
	inspectErr  error
	inspected   []string
	executions  int
}

// ExecuteSchemaStatements counts attempts; this fixture has no real target database.
func (s *recordingTargetResolverStub) ExecuteSchemaStatements(context.Context, string, string, []string) (*dbtarget.SchemaStatementExecution, error) {
	s.executions++
	return nil, errors.New("fixture target cannot execute schema statements")
}

type recordingTargetResolverCall struct {
	connectorID      string
	expectedRevision string
}

func (s *recordingTargetResolverStub) ResolveSavedTarget(_ context.Context, connectorID, expectedRevision string) (*schema.DatabaseConnector, error) {
	s.calls = append(s.calls, recordingTargetResolverCall{connectorID: connectorID, expectedRevision: expectedRevision})
	if s.err != nil {
		return nil, s.err
	}
	return s.connector, nil
}

// InspectTable reports every managed table as missing unless a test seeds it.
func (s *recordingTargetResolverStub) InspectTable(_ context.Context, connectorID, schemaName, tableName string) (*dbtarget.TableInspection, error) {
	s.inspected = append(s.inspected, connectorID+"/"+schemaName+"/"+tableName)
	if s.inspectErr != nil {
		return nil, s.inspectErr
	}
	if inspection, ok := s.inspections[tableName]; ok {
		return &inspection, nil
	}
	return &dbtarget.TableInspection{Status: dbtarget.TableInspectionMissing, Schema: schemaName, Table: tableName}, nil
}

func recordingPostgresConnector(enabled bool) *schema.DatabaseConnector {
	return &schema.DatabaseConnector{
		ID:               "connector-pg",
		Kind:             schema.DatabaseConnectorKindPostgres,
		IdentityRevision: "identity-7",
		Enabled:          enabled,
	}
}

type recordingTargetResolverFixture struct {
	router            *gin.Engine
	plan              *recordingplan.RecordingPlan
	workspaceRevision string
}

// previewBody builds a valid preview request for the fixture plan. An override
// replaces a field; an empty override value removes that field.
func (f recordingTargetResolverFixture) previewBody(t *testing.T, overrides map[string]string) string {
	t.Helper()
	fields := map[string]string{
		"plan_id":                     f.plan.ID,
		"connector_id":                "connector-pg",
		"expected_connector_revision": "identity-7",
		"expected_workspace_revision": f.workspaceRevision,
		"expected_plan_revision":      f.plan.Revision,
	}
	for key, value := range overrides {
		if value == "" {
			delete(fields, key)
			continue
		}
		fields[key] = value
	}
	body, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("encode preview body: %v", err)
	}
	return string(body)
}

func decodeRecordingPreview(t *testing.T, w *httptest.ResponseRecorder) recordingplan.SchemaPreviewToken {
	t.Helper()
	var response struct {
		Success bool                             `json:"success"`
		Data    recordingplan.SchemaPreviewToken `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || !response.Success {
		t.Fatalf("decode schema preview: success=%v err=%v body=%s", response.Success, err, w.Body.String())
	}
	return response.Data
}

func TestStudioV2WorkspaceRecordingPlans_SchemaPreviewUsesSavedConnectorRevisionAndServerDialect(t *testing.T) {
	resolver := &recordingTargetResolverStub{connector: recordingPostgresConnector(true)}
	fixture := newRecordingTargetResolverRouter(t, resolver)

	w := serveRecordingTargetResolverRequest(fixture.router, fixture.previewBody(t, nil))

	if w.Code != http.StatusOK {
		t.Fatalf("schema preview failed: %d %s", w.Code, w.Body.String())
	}
	data := decodeRecordingPreview(t, w)
	if data.ConnectorID != "connector-pg" || data.ConnectorRevision != "identity-7" || data.WorkspaceRevision != fixture.workspaceRevision ||
		data.PlanRevision != fixture.plan.Revision || data.Dialect != "postgres" || data.Schema != "public" {
		t.Fatalf("schema preview lost its server-resolved scope: %+v", data)
	}
	if data.Action != recordingplan.SchemaApplyAction || !strings.HasPrefix(data.OperationID, "op-") || len(data.Digest) != 64 {
		t.Fatalf("schema preview must carry its operation identity and digest: %+v", data)
	}
	if len(data.Statements) == 0 || !strings.Contains(data.Statements[0], "DOUBLE PRECISION") {
		t.Fatalf("schema preview did not use server postgres dialect: %v", data.Statements)
	}
	if len(resolver.calls) != 1 || resolver.calls[0].connectorID != "connector-pg" || resolver.calls[0].expectedRevision != "identity-7" {
		t.Fatalf("saved connector resolver call mismatch: %+v", resolver.calls)
	}
	if len(resolver.inspected) != 6 || resolver.inspected[0] != "connector-pg/public/gw_record_samples" {
		t.Fatalf("preview must inspect each managed table in the saved schema: %v", resolver.inspected)
	}
}

func TestStudioV2WorkspaceRecordingPlans_SchemaPreviewRejectsDialectSpoof(t *testing.T) {
	resolver := &recordingTargetResolverStub{connector: recordingPostgresConnector(true)}
	fixture := newRecordingTargetResolverRouter(t, resolver)

	w := serveRecordingTargetResolverRequest(fixture.router, fixture.previewBody(t, map[string]string{"dialect": "sqlite"}))

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("dialect spoof must be rejected with 422, got %d: %s", w.Code, w.Body.String())
	}
	if strings.Contains(w.Body.String(), "DOUBLE PRECISION") {
		t.Fatalf("dialect spoof response must not expose generated DDL: %s", w.Body.String())
	}
	if len(resolver.calls) != 1 || len(resolver.inspected) != 0 {
		t.Fatalf("dialect validation must resolve once and never inspect: calls=%+v inspected=%v", resolver.calls, resolver.inspected)
	}
}

func TestStudioV2WorkspaceRecordingPlans_ForeignConnectorIsRejectedBeforeResolverRead(t *testing.T) {
	resolver := &recordingTargetResolverStub{connector: &schema.DatabaseConnector{
		ID:               "connector-foreign",
		Kind:             schema.DatabaseConnectorKindSQLite,
		IdentityRevision: "identity-7",
		Enabled:          true,
	}}
	fixture := newRecordingTargetResolverRouter(t, resolver)

	w := serveRecordingTargetResolverRequest(fixture.router, fixture.previewBody(t, map[string]string{"connector_id": "connector-foreign"}))

	if w.Code != http.StatusNotFound {
		t.Fatalf("foreign connector must return safe 404, got %d: %s", w.Code, w.Body.String())
	}
	if len(resolver.calls) != 0 || len(resolver.inspected) != 0 {
		t.Fatalf("foreign connector must be rejected before resolver read: %+v", resolver.calls)
	}
}

func TestStudioV2WorkspaceRecordingPlans_MissingRevisionAndDisabledConnectorAreRejected(t *testing.T) {
	for _, field := range []string{"expected_connector_revision", "expected_workspace_revision", "expected_plan_revision"} {
		t.Run("missing "+field, func(t *testing.T) {
			resolver := &recordingTargetResolverStub{connector: recordingPostgresConnector(true)}
			fixture := newRecordingTargetResolverRouter(t, resolver)
			w := serveRecordingTargetResolverRequest(fixture.router, fixture.previewBody(t, map[string]string{field: ""}))
			if w.Code != http.StatusBadRequest {
				t.Fatalf("missing %s must return 400, got %d: %s", field, w.Code, w.Body.String())
			}
			if len(resolver.calls) != 0 || len(resolver.inspected) != 0 {
				t.Fatalf("missing %s must not resolve target: %+v", field, resolver.calls)
			}
		})
	}

	resolver := &recordingTargetResolverStub{connector: recordingPostgresConnector(false)}
	fixture := newRecordingTargetResolverRouter(t, resolver)
	w := serveRecordingTargetResolverRequest(fixture.router, fixture.previewBody(t, nil))
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("disabled connector must return 422, got %d: %s", w.Code, w.Body.String())
	}
	if len(resolver.calls) != 1 || len(resolver.inspected) != 0 {
		t.Fatalf("disabled connector should be checked by resolver only: calls=%+v inspected=%v", resolver.calls, resolver.inspected)
	}
}

func TestStudioV2WorkspaceRecordingPlans_SchemaPreviewRejectsStaleRevisionsBeforeTargetRead(t *testing.T) {
	for name, override := range map[string]map[string]string{
		"workspace revision": {"expected_workspace_revision": "workspace-stale"},
		"plan revision":      {"expected_plan_revision": "rev-0"},
	} {
		t.Run(name, func(t *testing.T) {
			resolver := &recordingTargetResolverStub{connector: recordingPostgresConnector(true)}
			fixture := newRecordingTargetResolverRouter(t, resolver)

			w := serveRecordingTargetResolverRequest(fixture.router, fixture.previewBody(t, override))

			assertRecordingPlanError(t, w, http.StatusConflict, "RECORDING_SCHEMA_PREVIEW_STALE", false)
			if len(resolver.calls) != 0 || len(resolver.inspected) != 0 {
				t.Fatalf("stale %s must be rejected before target access: calls=%+v inspected=%v", name, resolver.calls, resolver.inspected)
			}
		})
	}
}

func TestStudioV2WorkspaceRecordingPlans_SchemaPreviewBlocksUnsafeTargets(t *testing.T) {
	tests := []struct {
		name        string
		inspections map[string]dbtarget.TableInspection
		inspectErr  error
		overrides   map[string]string
		status      int
		code        string
		retryable   bool
		inspects    bool
	}{
		{
			name: "incompatible existing table",
			inspections: map[string]dbtarget.TableInspection{"gw_record_samples": {
				Status: dbtarget.TableInspectionExists, Table: "gw_record_samples", Columns: []dbtarget.ColumnInfo{{Name: "record_id"}},
			}},
			status: http.StatusUnprocessableEntity, code: "RECORDING_SCHEMA_INCOMPATIBLE", inspects: true,
		},
		{
			name:        "forbidden table",
			inspections: map[string]dbtarget.TableInspection{"gw_record_events": {Status: dbtarget.TableInspectionForbidden, Reason: "permission_denied"}},
			status:      http.StatusUnprocessableEntity, code: "RECORDING_SCHEMA_PERMISSION_DENIED", inspects: true,
		},
		{
			name:        "failed table inspection",
			inspections: map[string]dbtarget.TableInspection{"gw_record_samples": {Status: dbtarget.TableInspectionFailed, Reason: "connection_failed"}},
			status:      http.StatusServiceUnavailable, code: "RECORDING_SCHEMA_TARGET_UNAVAILABLE", retryable: true, inspects: true,
		},
		{
			name:       "inspection error",
			inspectErr: errors.New("private target diagnostic"),
			status:     http.StatusServiceUnavailable, code: "RECORDING_SCHEMA_TARGET_UNAVAILABLE", retryable: true, inspects: true,
		},
		{
			name:      "unsafe table prefix",
			overrides: map[string]string{"table_prefix": "gw;drop table x;"},
			status:    http.StatusUnprocessableEntity, code: "RECORDING_SCHEMA_PREFIX_INVALID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &recordingTargetResolverStub{connector: recordingPostgresConnector(true), inspections: tt.inspections, inspectErr: tt.inspectErr}
			fixture := newRecordingTargetResolverRouter(t, resolver)

			w := serveRecordingTargetResolverRequest(fixture.router, fixture.previewBody(t, tt.overrides))

			assertRecordingPlanError(t, w, tt.status, tt.code, tt.retryable)
			if body := w.Body.String(); strings.Contains(body, "CREATE TABLE") || strings.Contains(body, "private target diagnostic") || strings.Contains(body, "tok-") {
				t.Fatalf("blocked preview must not return statements, tokens or diagnostics: %s", body)
			}
			if tt.inspects == (len(resolver.inspected) == 0) {
				t.Fatalf("unexpected target inspection for %s: %v", tt.name, resolver.inspected)
			}
		})
	}
}

func TestStudioV2WorkspaceRecordingPlans_ForeignPlanIsRejectedBeforeTargetResolve(t *testing.T) {
	resolver := &recordingTargetResolverStub{connector: recordingPostgresConnector(true)}
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	_, err := workspaceSvc.GetOrCreate(context.Background())
	if err != nil {
		t.Fatalf("create active workspace: %v", err)
	}
	record, err := workspaceSvc.BindDatabaseConnector(context.Background(), "connector-pg")
	if err != nil {
		t.Fatalf("bind active connector: %v", err)
	}
	planRepo := recordingplan.NewMemoryRepository()
	planSvc := recordingplan.NewService(planRepo)
	foreign := &recordingplan.RecordingPlan{
		ID:           "plan-foreign",
		WorkspaceID:  "workspace-foreign",
		Name:         "Foreign plan",
		Streams:      []recordingplan.PlanStream{{StreamID: "stream-1", MeasurementID: "measurement-1", Mode: recordingplan.StreamModeRawHistory}},
		Destinations: []recordingplan.PlanDestination{{ConnectorID: "connector-pg", ConnectorRevision: "identity-7"}},
	}
	if err := planSvc.CreatePlan(context.Background(), foreign); err != nil {
		t.Fatalf("seed foreign plan: %v", err)
	}
	handler := handlers.NewStudioV2WorkspaceRecordingPlansHandler(workspaceSvc, planSvc, resolver)
	router := gin.New()
	router.POST("/schema-preview", handler.SchemaPreview)
	fixture := recordingTargetResolverFixture{router: router, plan: foreign, workspaceRevision: record.DatabaseSetupRevision}

	w := serveRecordingTargetResolverRequest(router, fixture.previewBody(t, nil))

	if w.Code != http.StatusNotFound {
		t.Fatalf("foreign plan must return safe 404, got %d: %s", w.Code, w.Body.String())
	}
	if len(resolver.calls) != 0 || len(resolver.inspected) != 0 {
		t.Fatalf("foreign plan must be rejected before target resolve: %+v", resolver.calls)
	}
}

func TestStudioV2WorkspaceRecordingPlans_CreateAndUpdateValidateSavedDestinations(t *testing.T) {
	resolver := &recordingTargetResolverStub{connector: recordingPostgresConnector(true)}
	router := newRecordingTargetResolverRouter(t, resolver).router

	createBody := `{"id":"plan-created","name":"Created plan","members":[{"member_id":"member-2","measurement_id":"measurement-2","equipment_id":"equipment-line"}],"streams":[{"stream_id":"stream-2","measurement_id":"measurement-2","mode":"raw_history"}],"destinations":[{"destination_id":"destination-2","connector_id":"connector-pg","connector_revision":"identity-7"}]}`
	w := serveRecordingTargetResolverRequestAt(router, http.MethodPost, "/recording-plans", createBody)
	if w.Code != http.StatusCreated {
		t.Fatalf("valid destination create failed: %d %s", w.Code, w.Body.String())
	}
	if len(resolver.calls) != 1 || resolver.calls[0].expectedRevision != "identity-7" {
		t.Fatalf("create must resolve the saved connector revision once: %+v", resolver.calls)
	}

	resolver.calls = nil
	updateBody := `{"name":"Updated plan","members":[{"member_id":"member-1","measurement_id":"measurement-1","equipment_id":"equipment-line"}],"streams":[{"stream_id":"stream-1","measurement_id":"measurement-1","mode":"raw_history"}],"destinations":[{"destination_id":"destination-1","connector_id":"connector-pg","connector_revision":"identity-7"}]}`
	w = serveRecordingTargetResolverRequestAt(router, http.MethodPut, "/recording-plans/plan-target", updateBody)
	if w.Code != http.StatusOK {
		t.Fatalf("valid destination update failed: %d %s", w.Code, w.Body.String())
	}
	if len(resolver.calls) != 1 || resolver.calls[0].expectedRevision != "identity-7" {
		t.Fatalf("update must resolve the saved connector revision once: %+v", resolver.calls)
	}

	resolver.calls = nil
	foreignBody := `{"id":"plan-foreign-destination","name":"Foreign destination","streams":[{"stream_id":"stream-3","measurement_id":"measurement-3","mode":"raw_history"}],"destinations":[{"destination_id":"destination-3","connector_id":"connector-foreign","connector_revision":"identity-7"}]}`
	w = serveRecordingTargetResolverRequestAt(router, http.MethodPost, "/recording-plans", foreignBody)
	if w.Code != http.StatusNotFound {
		t.Fatalf("foreign create destination must return safe 404, got %d: %s", w.Code, w.Body.String())
	}
	if len(resolver.calls) != 0 {
		t.Fatalf("foreign create destination must not resolve a target: %+v", resolver.calls)
	}

	workspaceBody := `{"workspace_id":"workspace-foreign","id":"plan-foreign-workspace","name":"Foreign workspace","streams":[{"stream_id":"stream-4","measurement_id":"measurement-4","mode":"raw_history"}],"destinations":[{"destination_id":"destination-4","connector_id":"connector-pg","connector_revision":"identity-7"}]}`
	w = serveRecordingTargetResolverRequestAt(router, http.MethodPost, "/recording-plans", workspaceBody)
	if w.Code != http.StatusNotFound {
		t.Fatalf("foreign workspace create must return safe 404, got %d: %s", w.Code, w.Body.String())
	}
	if len(resolver.calls) != 0 {
		t.Fatalf("foreign workspace create must not resolve a target: %+v", resolver.calls)
	}

	emptyDestinationsBody := `{"id":"plan-empty-destinations","name":"Empty destinations","streams":[{"stream_id":"stream-5","measurement_id":"measurement-5","mode":"raw_history"}],"destinations":[]}`
	w = serveRecordingTargetResolverRequestAt(router, http.MethodPost, "/recording-plans", emptyDestinationsBody)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("create without destinations must return 422, got %d: %s", w.Code, w.Body.String())
	}
	if len(resolver.calls) != 0 {
		t.Fatalf("create without destinations must not resolve a target: %+v", resolver.calls)
	}
}

func newRecordingTargetResolverRouter(t *testing.T, resolver *recordingTargetResolverStub) recordingTargetResolverFixture {
	t.Helper()
	connectorID := "connector-pg"
	connectorRevision := "identity-7"

	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	if _, err := workspaceSvc.GetOrCreate(context.Background()); err != nil {
		t.Fatalf("create active workspace: %v", err)
	}
	if _, err := workspaceSvc.BindDatabaseConnector(context.Background(), connectorID); err != nil {
		t.Fatalf("bind active connector: %v", err)
	}
	record, err := workspaceSvc.GetOrCreate(context.Background())
	if err != nil {
		t.Fatalf("read active workspace: %v", err)
	}
	planRepo := recordingplan.NewMemoryRepository()
	planSvc := recordingplan.NewService(planRepo)
	plan := &recordingplan.RecordingPlan{
		ID:           "plan-target",
		WorkspaceID:  record.ID,
		Name:         "Target plan",
		Streams:      []recordingplan.PlanStream{{StreamID: "stream-1", MeasurementID: "measurement-1", Mode: recordingplan.StreamModeRawHistory}},
		Destinations: []recordingplan.PlanDestination{{ConnectorID: connectorID, ConnectorRevision: connectorRevision}},
	}
	if err := planSvc.CreatePlan(context.Background(), plan); err != nil {
		t.Fatalf("seed target plan: %v", err)
	}
	if _, err := workspaceSvc.AttachDevice(context.Background(), "device-a"); err != nil {
		t.Fatalf("attach device fixture: %v", err)
	}
	current, err := workspaceSvc.GetOrCreate(context.Background())
	if err != nil {
		t.Fatalf("read workspace revision: %v", err)
	}
	handler := handlers.NewStudioV2WorkspaceRecordingPlansHandler(workspaceSvc, planSvc, resolver)
	membership := newRecordingMembershipStub(record.ID)
	handler.SetRecordingMembershipServices(membership, membership)
	router := gin.New()
	router.POST("/recording-plans", handler.Create)
	router.PUT("/recording-plans/:id", handler.Update)
	router.POST("/schema-preview", handler.SchemaPreview)
	router.POST("/schema-apply", handler.SchemaApply)
	return recordingTargetResolverFixture{router: router, plan: plan, workspaceRevision: current.DatabaseSetupRevision}
}

// MySQL has no verified managed schema behavior, so a preview is refused
// before any statement is built and the target is never read.
func TestStudioV2WorkspaceRecordingPlans_MySQLManagedSchemaStaysUnavailable(t *testing.T) {
	resolver := &recordingTargetResolverStub{connector: &schema.DatabaseConnector{
		ID: "connector-pg", Kind: schema.DatabaseConnectorKindMySQL, IdentityRevision: "identity-7", Enabled: true,
	}}
	fixture := newRecordingTargetResolverRouter(t, resolver)

	w := serveRecordingTargetResolverRequest(fixture.router, fixture.previewBody(t, nil))

	assertRecordingPlanError(t, w, http.StatusUnprocessableEntity, "SCHEMA_PREVIEW_FAILED", false)
	if len(resolver.inspected) != 0 || resolver.executions != 0 {
		t.Fatalf("an unverified adapter must not touch the target: inspected=%v executions=%d", resolver.inspected, resolver.executions)
	}
	if strings.Contains(w.Body.String(), "CREATE TABLE") {
		t.Fatalf("a refused preview must not expose generated DDL: %s", w.Body.String())
	}
}

// A target that refuses the confirmed statements ends as a recorded failure:
// nothing is reported as created and no verified digest is claimed.
func TestStudioV2WorkspaceRecordingPlans_SchemaApplyRecordsARefusedExecutionAsFailed(t *testing.T) {
	resolver := &recordingTargetResolverStub{connector: recordingPostgresConnector(true)}
	fixture := newRecordingTargetResolverRouter(t, resolver)
	preview := decodeRecordingPreview(t, serveRecordingTargetResolverRequest(fixture.router, fixture.previewBody(t, nil)))
	body, err := json.Marshal(map[string]string{
		"token": preview.Token, "operation_id": preview.OperationID, "expected_workspace_revision": fixture.workspaceRevision,
		"expected_plan_revision": fixture.plan.Revision, "expected_connector_revision": "identity-7",
	})
	if err != nil {
		t.Fatalf("encode apply body: %v", err)
	}

	w := serveRecordingTargetResolverRequestAt(fixture.router, http.MethodPost, "/schema-apply", string(body))

	refused := decodeRecordingOperation(t, w).Data
	if w.Code != http.StatusOK || refused.Status != recordingplan.SchemaOperationFailed || refused.ExecutedStatements != 0 || refused.VerifiedDigest != "" {
		t.Fatalf("a refused execution must be recorded as a failure: %d %s", w.Code, w.Body.String())
	}
	if refused.NextAction == "" {
		t.Fatalf("a failed operation must name the next step: %+v", refused)
	}
	if resolver.executions != 1 {
		t.Fatalf("the confirmed batch must be attempted exactly once, got %d executions", resolver.executions)
	}
}

func serveRecordingTargetResolverRequest(router *gin.Engine, body string) *httptest.ResponseRecorder {
	return serveRecordingTargetResolverRequestAt(router, http.MethodPost, "/schema-preview", body)
}

func serveRecordingTargetResolverRequestAt(router *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequestWithContext(context.Background(), method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
