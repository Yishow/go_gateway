package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

func TestRuntimeRuleReadinessSafetyProductionDeviceDiagnostics(t *testing.T) {
	for _, diagnostic := range []string{"driver /private/fixture.sqlite", "postgres://fixture.internal/db", "uid=sa;pwd=fixture-password", "file:/private/fixture.sqlite"} {
		t.Run(diagnostic, func(t *testing.T) {
			repo := device.NewMemoryRepository()
			success := false
			require.NoError(t, repo.Create(t.Context(), &schema.Device{ID: "device", Name: "泵浦 A (°C)", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{"host":"127.0.0.1","port":502}`, LastTestSuccess: &success, LastTestError: diagnostic}))
			svc := device.NewService(repo, nil)
			ws := workspace.NewService(workspace.NewMemoryRepository())
			_, err := ws.AttachDevice(t.Context(), "device")
			require.NoError(t, err)
			ws.WithReadinessServices(svc, nil, nil, nil)
			h := NewRuntimeHandler(svc, nil, nil, nil, nil, ws)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodGet, "/runtime-context", nil)
			c.Request.RemoteAddr = "203.0.113.10:5678"
			h.WorkspaceContext(c)
			require.Equal(t, http.StatusOK, rec.Code)
			require.NotContains(t, rec.Body.String(), "fixture")
			require.NotContains(t, rec.Body.String(), "uid=sa")
			var result struct {
				Data struct {
					Setup struct {
						Readiness workspace.ReadinessSummary `json:"readiness_summary"`
					} `json:"setup"`
				} `json:"data"`
			}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &result))
			require.False(t, result.Data.Setup.Readiness.Ready)
			require.NotEmpty(t, result.Data.Setup.Readiness.Issues)
			require.Equal(t, "device", result.Data.Setup.Readiness.Issues[0].Scope)
			saved, err := repo.GetByID(t.Context(), "device")
			require.NoError(t, err)
			require.Equal(t, diagnostic, saved.LastTestError)
		})
	}
}

func TestRuntimeRuleContextUnitRejectsPathsPreservesUnits(t *testing.T) {
	for _, value := range []string{"/private/db", "file:/private/db", "uid=sa;pwd=secret", `C:\private\db`} {
		require.Empty(t, runtimeSafeUnit(value), value)
	}
	for _, value := range []string{"°C", "m/s", "kWh", "µS", "Ω", "%"} {
		require.Equal(t, value, runtimeSafeUnit(value))
	}
}

func TestRuntimeRuleReadinessSafetyDetachedMissingUnknown(t *testing.T) {
	require.Nil(t, projectRuntimeReadiness(nil))
	source := &workspace.ReadinessSummary{Ready: true, BlockingCount: 1, Issues: []workspace.ReadinessIssue{{Code: "private.internal", Severity: "future", Step: "/private/fixture", Scope: "private.internal", Message: "uid=sa;pwd=fixture"}}}
	projected := projectRuntimeReadiness(source)
	raw, err := json.Marshal(projected)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "private.internal")
	require.NotContains(t, string(raw), "fixture")
	require.False(t, projected.Ready)
	projected.Issues[0].Message = "changed"
	projected.Issues[0].Scope = "changed"
	require.Equal(t, "private.internal", source.Issues[0].Scope)
	require.Equal(t, "uid=sa;pwd=fixture", source.Issues[0].Message)
}

func TestRuntimeRuleContextSafetyWorkspaceIdentity(t *testing.T) {
	repo := workspace.NewMemoryRepository()
	require.NoError(t, repo.Save(t.Context(), &workspace.Record{ID: "private.internal", OrderedDeviceIDs: []string{}}))
	ws := workspace.NewService(repo)
	h := NewRuntimeHandler(device.NewService(device.NewMemoryRepository(), nil), nil, nil, nil, nil, ws)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/runtime-context", nil)
	h.WorkspaceContext(c)
	require.Equal(t, http.StatusOK, rec.Code)
	require.NotContains(t, rec.Body.String(), "private.internal")
}

func TestRuntimeRuleContextSafetySemanticTagKeys(t *testing.T) {
	for _, tc := range []struct{ key, want string }{{"line01.temp.inlet", "line01.temp.inlet"}, {"line01/temp.inlet", "line01/temp.inlet"}, {"postgres://db.internal/table", ""}, {"uid=sa;pwd=fixture", ""}, {"/private/db", ""}, {"../private/db", ""}, {"file:/private/db", ""}} {
		t.Run(tc.key, func(t *testing.T) {
			tagID, mappingID := "tag", "mapping"
			h := &RuntimeHandler{tagReader: runtimeSetupTagReaderStub{records: map[string]*schema.Tag{"tag": {ID: "tag", Key: tc.key}}}, mappingReader: runtimeSetupMappingReaderStub{records: map[string]*schema.Mapping{"mapping": {ID: "mapping", PointID: "point", TagID: "tag"}}}}
			value, err := h.mapRuntimeSetupMapping(t.Context(), &schema.SourceRuleLink{RuleID: "rule", PointID: "point", Address: "40001", TagID: &tagID, MappingID: &mappingID}, nil)
			require.NoError(t, err)
			require.Equal(t, tc.want, value.TagKey)
		})
	}
}

func TestRuntimeRuleReadErrorsDoNotReflectUnsafeRequestIdentity(t *testing.T) {
	for _, surface := range []string{"context", "status", "points"} {
		t.Run(surface, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodGet, "/runtime-read", nil)
			c.Request.Header.Set("X-Request-ID", "file:/private/fixture.sqlite")
			switch surface {
			case "context":
				(&RuntimeHandler{}).WorkspaceContext(c)
			case "status":
				(&RuntimeHandler{}).Status(c)
			case "points":
				NewPointHandler(point.NewService(runtimePointFailedRepository{point.NewMemoryRepository()}, nil)).List(c)
			}
			require.Equal(t, http.StatusServiceUnavailable, rec.Code)
			require.NotContains(t, rec.Body.String(), "fixture")
			require.NotContains(t, rec.Body.String(), "file:")
		})
	}
}

func TestRuntimeRuleReadErrorHeaderAliases(t *testing.T) {
	for _, header := range []string{"X-Request-ID", "X-Request-Id", "x-request-id", "X-REQUEST-ID"} {
		for _, value := range []string{"private.internal", "postgres://user:fixture@db", "/private/fixture.sqlite", "uid=sa;pwd=fixture", "file:/private/fixture"} {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodGet, "/runtime-context", nil)
			c.Request.Header.Set(header, value)
			(&RuntimeHandler{}).WorkspaceContext(c)
			var body struct {
				Error struct {
					RequestID string `json:"request_id"`
				} `json:"error"`
			}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			require.NotEqual(t, value, body.Error.RequestID)
			require.NotEmpty(t, body.Error.RequestID)
			require.Equal(t, http.StatusServiceUnavailable, rec.Code)
		}
	}
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/runtime-context", nil)
	c.Request.Header.Set("X-Request-ID", "req-safe-123")
	(&RuntimeHandler{}).WorkspaceContext(c)
	require.Contains(t, rec.Body.String(), "req-safe-123")
}
