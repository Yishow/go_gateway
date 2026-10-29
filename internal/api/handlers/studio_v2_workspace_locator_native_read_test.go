package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func locatorNativeQuery(t *testing.T, f workspaceDatabaseFixture) url.Values {
	t.Helper()
	r, err := f.workspaceSvc.ReadSnapshot(t.Context())
	require.NoError(t, err)
	return url.Values{"locator_workspace_id": {r.ID}, "locator_workspace_revision": {r.DatabaseSetupRevision}}
}
func locatorNativeServe(handler gin.HandlerFunc, q url.Values) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/native?"+q.Encode(), nil)
	handler(c)
	return w
}
func TestStudioLocatorNativeMetadataNeverBootstrapsMissingWorkspace(t *testing.T) {
	f := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, f)
	q := locatorNativeQuery(t, f)
	_, err := f.mainDB.Exec(`DELETE FROM system_settings WHERE key='studio_v2_workspace'`)
	require.NoError(t, err)
	response := locatorNativeServe(f.handler.GetMetadata, q)
	assert.Equal(t, 409, response.Code, response.Body.String())
	var n int
	require.NoError(t, f.mainDB.QueryRow(`SELECT COUNT(*) FROM system_settings WHERE key='studio_v2_workspace'`).Scan(&n))
	assert.Zero(t, n)
}

type locatorNativeOperations struct {
	schemaOperationReader
	change func()
}

func (r locatorNativeOperations) GetSchemaOperation(ctx context.Context, w, id string) (*recordingplan.SchemaOperation, error) {
	r.change()
	return r.schemaOperationReader.GetSchemaOperation(ctx, w, id)
}
func (r locatorNativeOperations) ListUnresolvedTestWrites(ctx context.Context, w string) (*recordingplan.TestWriteRecoveryList, error) {
	r.change()
	return r.schemaOperationReader.(testWriteRecoveryReader).ListUnresolvedTestWrites(ctx, w)
}
func TestStudioLocatorNativeOperationDeletionDuringReadNeverRecreates(t *testing.T) {
	for _, list := range []bool{false, true} {
		t.Run(map[bool]string{false: "operation", true: "recovery"}[list], func(t *testing.T) {
			f := newWorkspaceDatabaseFixture(t)
			q := locatorNativeQuery(t, f)
			plans := recordingplan.NewService(recordingplan.NewSQLRepository(f.mainDB))
			reader := locatorNativeOperations{schemaOperationReader: plans, change: func() {
				_, err := f.mainDB.Exec(`DELETE FROM system_settings WHERE key='studio_v2_workspace'`)
				require.NoError(t, err)
			}}
			h := NewStudioV2WorkspaceDatabaseOperationsHandler(f.workspaceSvc, reader)
			handler := h.Get
			if list {
				handler = h.ListTestWrites
			}
			response := locatorNativeServe(handler, q)
			require.Equal(t, 409, response.Code, response.Body.String())
			var n int
			require.NoError(t, f.mainDB.QueryRow(`SELECT COUNT(*) FROM system_settings WHERE key='studio_v2_workspace'`).Scan(&n))
			require.Zero(t, n)
		})
	}
}

type locatorNativeCapabilityProbe struct {
	*dbtarget.ConnectorService
	change func()
}

func (p locatorNativeCapabilityProbe) ProbeMySQLRecordingAtRevision(ctx context.Context, id, rev string) (*dbtarget.MySQLRecordingCapability, error) {
	p.change()
	return &dbtarget.MySQLRecordingCapability{ConnectorID: id, ConnectorRevision: rev}, nil
}
func (p locatorNativeCapabilityProbe) ProbeSQLServerRecordingAtRevision(ctx context.Context, id, rev string) (*dbtarget.SQLServerRecordingCapability, error) {
	p.change()
	return &dbtarget.SQLServerRecordingCapability{ConnectorID: id, ConnectorRevision: rev}, nil
}
func TestStudioLocatorNativeCapabilityDeletionDuringProbeNeverRecreates(t *testing.T) {
	for _, kind := range []string{"mysql", "sqlserver"} {
		t.Run(kind, func(t *testing.T) {
			f := newWorkspaceDatabaseFixture(t)
			saveWorkspaceDatabaseConfig(t, f)
			record, err := f.workspaceSvc.ReadSnapshot(t.Context())
			require.NoError(t, err)
			_, err = f.mainDB.Exec(`UPDATE database_connectors SET kind=? WHERE id=?`, kind, record.DatabaseConnectorID)
			require.NoError(t, err)
			connector, err := f.connectorSvc.GetByID(t.Context(), record.DatabaseConnectorID)
			require.NoError(t, err)
			q := locatorNativeQuery(t, f)
			q.Set("kind", kind)
			q.Set("connector_id", connector.ID)
			q.Set("expected_connector_revision", connector.IdentityRevision)
			probe := locatorNativeCapabilityProbe{ConnectorService: f.connectorSvc, change: func() {
				_, err := f.mainDB.Exec(`DELETE FROM system_settings WHERE key='studio_v2_workspace'`)
				require.NoError(t, err)
			}}
			h := NewStudioV2WorkspaceRecordingPlansHandler(f.workspaceSvc, recordingplan.NewService(recordingplan.NewSQLRepository(f.mainDB)), probe)
			response := locatorNativeServe(h.Capabilities, q)
			assert.Equal(t, 409, response.Code, response.Body.String())
			var n int
			require.NoError(t, f.mainDB.QueryRow(`SELECT COUNT(*) FROM system_settings WHERE key='studio_v2_workspace'`).Scan(&n))
			assert.Zero(t, n)
		})
	}
}

func TestStudioLocatorNativeStrictPairAndOrdinaryCompatibility(t *testing.T) {
	f := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, f)
	_, connector := boundWorkspaceConnector(t, f)
	base := locatorNativeQuery(t, f)
	base.Set("expected_connector_revision", connector.IdentityRevision)
	for _, tc := range []struct {
		name, key, value string
		code             int
	}{{"foreign", "locator_workspace_id", "foreign", 403}, {"stale", "locator_workspace_revision", "stale", 409}, {"unsafe", "locator_workspace_id", "postgres://password", 400}, {"partial", "locator_workspace_revision", "DELETE", 400}, {"duplicate", "locator_workspace_id", "DUPLICATE", 400}} {
		t.Run(tc.name, func(t *testing.T) {
			q := url.Values{}
			for k, v := range base {
				q[k] = append([]string{}, v...)
			}
			switch tc.value {
			case "DELETE":
				q.Del(tc.key)
			case "DUPLICATE":
				q.Add(tc.key, q.Get(tc.key))
			default:
				q.Set(tc.key, tc.value)
			}
			response := locatorNativeServe(f.handler.GetMetadata, q)
			require.Equal(t, tc.code, response.Code, response.Body.String())
		})
	}
	var before, after int
	require.NoError(t, f.mainDB.QueryRow(`SELECT total_changes()`).Scan(&before))
	response := locatorNativeServe(f.handler.GetMetadata, base)
	require.Equal(t, 200, response.Code, response.Body.String())
	require.NoError(t, f.mainDB.QueryRow(`SELECT total_changes()`).Scan(&after))
	require.Equal(t, before, after)
	ordinary := url.Values{"expected_connector_revision": {connector.IdentityRevision}}
	require.Equal(t, 200, locatorNativeServe(f.handler.GetMetadata, ordinary).Code)
	ops := NewStudioV2WorkspaceDatabaseOperationsHandler(f.workspaceSvc, recordingplan.NewService(recordingplan.NewSQLRepository(f.mainDB)))
	require.Equal(t, 200, locatorNativeServe(ops.ListTestWrites, base).Code)
	require.Equal(t, 200, locatorNativeServe(ops.ListTestWrites, nil).Code)
}

type locatorNativeMeasurementRepo struct {
	measurement.Repository
	change func()
	err    error
}

func (r locatorNativeMeasurementRepo) ListByWorkspace(ctx context.Context, w string) ([]measurement.MeasurementDefinition, error) {
	if r.change != nil {
		r.change()
	}
	if r.err != nil {
		return nil, r.err
	}
	return r.Repository.ListByWorkspace(ctx, w)
}
func TestStudioLocatorNativeMeasurementMissingAndConcurrentDeletion(t *testing.T) {
	for _, during := range []bool{false, true} {
		t.Run(map[bool]string{false: "missing", true: "during"}[during], func(t *testing.T) {
			f := newWorkspaceDatabaseFixture(t)
			q := locatorNativeQuery(t, f)
			remove := func() {
				_, err := f.mainDB.Exec(`DELETE FROM system_settings WHERE key='studio_v2_workspace'`)
				require.NoError(t, err)
			}
			repo := locatorNativeMeasurementRepo{Repository: measurement.NewSQLRepository(f.mainDB)}
			if during {
				repo.change = remove
			} else {
				remove()
			}
			h := NewStudioV2WorkspaceMeasurementsHandler(f.workspaceSvc, nil, measurement.NewService(repo))
			res := locatorNativeServe(h.List, q)
			assert.Equal(t, 409, res.Code, res.Body.String())
			var n int
			require.NoError(t, f.mainDB.QueryRow(`SELECT COUNT(*) FROM system_settings WHERE key='studio_v2_workspace'`).Scan(&n))
			assert.Zero(t, n)
		})
	}
}
func TestStudioLocatorNativeMeasurementProviderFailureSanitized(t *testing.T) {
	f := newWorkspaceDatabaseFixture(t)
	q := locatorNativeQuery(t, f)
	repo := locatorNativeMeasurementRepo{Repository: measurement.NewSQLRepository(f.mainDB), err: errors.New("password=secret")}
	h := NewStudioV2WorkspaceMeasurementsHandler(f.workspaceSvc, nil, measurement.NewService(repo))
	res := locatorNativeServe(h.List, q)
	require.Equal(t, 503, res.Code, res.Body.String())
	require.NotContains(t, res.Body.String(), "password=secret")
}

type locatorNativeOperationResult struct {
	op *recordingplan.SchemaOperation
}

func (r locatorNativeOperationResult) GetSchemaOperation(context.Context, string, string) (*recordingplan.SchemaOperation, error) {
	return r.op, nil
}
func TestStudioLocatorNativeOperationRejectsForeignAndNilProviderResult(t *testing.T) {
	f := newWorkspaceDatabaseFixture(t)
	q := locatorNativeQuery(t, f)
	for _, tc := range []struct {
		op   *recordingplan.SchemaOperation
		code int
	}{{nil, 503}, {&recordingplan.SchemaOperation{WorkspaceID: "foreign", Reason: "secret-other-workspace"}, 403}} {
		h := NewStudioV2WorkspaceDatabaseOperationsHandler(f.workspaceSvc, locatorNativeOperationResult{tc.op})
		res := locatorNativeServe(h.Get, q)
		require.Equal(t, tc.code, res.Code, res.Body.String())
		require.NotContains(t, res.Body.String(), "secret-other-workspace")
	}
}

type locatorNativeForeignMeasurementRepo struct{ measurement.Repository }

func (r locatorNativeForeignMeasurementRepo) ListByWorkspace(context.Context, string) ([]measurement.MeasurementDefinition, error) {
	return []measurement.MeasurementDefinition{{WorkspaceID: "foreign", Name: "secret-other-workspace"}}, nil
}
func TestStudioLocatorNativeMeasurementRejectsForeignProviderResult(t *testing.T) {
	f := newWorkspaceDatabaseFixture(t)
	h := NewStudioV2WorkspaceMeasurementsHandler(f.workspaceSvc, nil, measurement.NewService(locatorNativeForeignMeasurementRepo{measurement.NewSQLRepository(f.mainDB)}))
	res := locatorNativeServe(h.List, locatorNativeQuery(t, f))
	require.Equal(t, 403, res.Code, res.Body.String())
	require.NotContains(t, res.Body.String(), "secret-other-workspace")
}

type locatorNativeConnectorRepository struct {
	dbtarget.ConnectorRepository
	get func(context.Context, string) (*schema.DatabaseConnector, error)
}

func (r locatorNativeConnectorRepository) GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error) {
	return r.get(ctx, id)
}
func TestStudioLocatorNativeMetadataProviderFailureAndMissingGroup(t *testing.T) {
	for _, missingGroup := range []bool{false, true} {
		t.Run(map[bool]string{false: "provider", true: "missing group"}[missingGroup], func(t *testing.T) {
			f := newWorkspaceDatabaseFixture(t)
			saveWorkspaceDatabaseConfig(t, f)
			_, connector := boundWorkspaceConnector(t, f)
			q := locatorNativeQuery(t, f)
			q.Set("expected_connector_revision", connector.IdentityRevision)
			code := 503
			if missingGroup {
				f.handler.WithWriteGroups(workspace.NewWriteGroupService(f.workspaceSvc, workspace.NewSQLWriteGroupRepository(f.mainDB)))
				q.Set("group_id", "gone")
				q.Set("expected_group_revision", "G1")
				code = 409
			} else {
				f.handler.connectorSvc = dbtarget.NewConnectorService(locatorNativeConnectorRepository{ConnectorRepository: f.connectorRepo, get: func(context.Context, string) (*schema.DatabaseConnector, error) {
					return nil, errors.New("password=secret")
				}})
			}
			res := locatorNativeServe(f.handler.GetMetadata, q)
			require.Equal(t, code, res.Code, res.Body.String())
			require.NotContains(t, res.Body.String(), "password=secret")
		})
	}
}

func TestStudioLocatorNativeMetadataDeletionDuringReadOnlyInspection(t *testing.T) {
	f := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, f)
	_, connector := boundWorkspaceConnector(t, f)
	q := locatorNativeQuery(t, f)
	q.Set("expected_connector_revision", connector.IdentityRevision)
	calls := 0
	f.handler.connectorSvc = dbtarget.NewConnectorService(locatorNativeConnectorRepository{ConnectorRepository: f.connectorRepo, get: func(ctx context.Context, id string) (*schema.DatabaseConnector, error) {
		c, err := f.connectorRepo.GetByID(ctx, id)
		calls++
		if calls == 2 {
			_, deleteErr := f.mainDB.Exec(`DELETE FROM system_settings WHERE key='studio_v2_workspace'`)
			require.NoError(t, deleteErr)
		}
		return c, err
	}})
	res := locatorNativeServe(f.handler.GetMetadata, q)
	require.GreaterOrEqual(t, calls, 2)
	require.Equal(t, 409, res.Code, res.Body.String())
	var n int
	require.NoError(t, f.mainDB.QueryRow(`SELECT COUNT(*) FROM system_settings WHERE key='studio_v2_workspace'`).Scan(&n))
	require.Zero(t, n)
}
