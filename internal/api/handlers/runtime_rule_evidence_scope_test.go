package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

func runtimeEvidenceScopeFixture(t *testing.T) (*RuntimeHandler, *workspace.Record) {
	t.Helper()
	ws := workspace.NewService(workspace.NewMemoryRepository())
	_, err := ws.AttachDevice(t.Context(), "device-1")
	require.NoError(t, err)
	record, err := ws.BindDatabaseConnector(t.Context(), "connector-1")
	require.NoError(t, err)
	repo := device.NewMemoryRepository()
	require.NoError(t, repo.Create(t.Context(), &schema.Device{ID: "device-1", Name: "Meter", Protocol: schema.ProtocolModbusTCP}))
	return NewRuntimeHandler(device.NewService(repo, nil), nil, nil, nil, nil, ws), record
}

func runtimeEvidenceScopeRequest(t *testing.T, h *RuntimeHandler, query string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/runtime-context"+query, nil)
	h.WorkspaceContext(c)
	return rec
}

func TestRuntimeRuleEvidenceRejectsPartialMalformedAndDuplicateScope(t *testing.T) {
	for _, query := range []string{
		"?workspace_id=W", "?workspace_id=W&workspace_revision=R", "?device_id=device-1",
		"?workspace_id=&workspace_revision=R&device_id=device-1",
		"?workspace_id=W&workspace_revision=R&device_id=device-1&device_id=device-2",
		"?workspace_id=W&workspace_id=W&workspace_revision=R&device_id=device-1",
		"?workspace_id=W&workspace_revision=R&device_id=%2Fprivate%2Fsecret",
		"?workspace_id=W&workspace_revision=R&device_id=db.private.internal",
		"?workspace_id=%ZZ&workspace_revision=%ZZ&device_id=%ZZ",
		"?workspace_id=W&workspace_revision=R&device_id=device-1&device_id=%ZZ",
		"?workspace_id=W;x&workspace_revision=R;x&device_id=device-1;x",
		"?workspace_id=W&workspace_revision=R&device_id=device-1&device_id=bad;query",
		"?workspace_id=W&workspace_revision=R&device_id=device-1&unexpected=%ZZ",
	} {
		t.Run(query, func(t *testing.T) {
			h, _ := runtimeEvidenceScopeFixture(t)
			rec := runtimeEvidenceScopeRequest(t, h, query)
			require.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())
			require.Contains(t, rec.Body.String(), `"code":"unavailable"`)
			require.NotContains(t, rec.Body.String(), `"rule_evidence"`)
			require.NotContains(t, rec.Body.String(), "private")
		})
	}
}

func TestRuntimeRuleEvidenceAbsentScopePreservesLegacyContext(t *testing.T) {
	h, _ := runtimeEvidenceScopeFixture(t)
	rec := runtimeEvidenceScopeRequest(t, h, "")
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	require.Contains(t, rec.Body.String(), `"device_id":"device-1"`)
	require.NotContains(t, rec.Body.String(), `"rule_evidence"`)
}

func TestRuntimeRuleEvidenceScopedMissingProviderIsUnavailable(t *testing.T) {
	h, record := runtimeEvidenceScopeFixture(t)
	query := url.Values{"workspace_id": {record.ID}, "workspace_revision": {record.DatabaseSetupRevision}, "device_id": {"device-1"}}
	rec := runtimeEvidenceScopeRequest(t, h, "?"+query.Encode())
	require.Equal(t, http.StatusServiceUnavailable, rec.Code, rec.Body.String())
	require.Contains(t, rec.Body.String(), `"code":"unavailable"`)
	require.NotContains(t, rec.Body.String(), `"rule_evidence"`)
}

func TestRuntimeRuleEvidenceEmptyWorkspaceDoesNotBootstrap(t *testing.T) {
	repo := workspace.NewMemoryRepository()
	service := workspace.NewService(repo)
	h := NewRuntimeHandler(device.NewService(device.NewMemoryRepository(), nil), nil, nil, nil, nil, service)
	rec := runtimeEvidenceScopeRequest(t, h, "?workspace_id=W&workspace_revision=R&device_id=D")
	_, err := repo.Get(t.Context())
	require.ErrorIs(t, err, workspace.ErrNotFound, "a scoped read must not create an empty workspace")
	require.Equal(t, http.StatusServiceUnavailable, rec.Code, rec.Body.String())
}
