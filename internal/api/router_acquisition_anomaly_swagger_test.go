package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	swagger "go-gateway/docs/swagger"
	"go-gateway/internal/datalink/anomaly"
	"go-gateway/internal/datalink/device"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionAnomalySwaggerURLsMatchProductionRoutes(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	router := NewRouter(&DatalinkServices{
		Device:               device.NewService(device.NewSQLRepository(f.db), nil),
		Workspace:            f.workspace,
		AcquisitionAnomalies: anomaly.NewSQLRepository(f.db),
	})
	var doc struct {
		BasePath string                                `json:"basePath"`
		Paths    map[string]map[string]json.RawMessage `json:"paths"`
	}
	require.NoError(t, json.Unmarshal([]byte(swagger.SwaggerInfo.ReadDoc()), &doc))
	for _, path := range []string{
		"/v1/datalink/devices/{id}/acquisition-anomalies",
		"/v1/datalink/acquisition-anomalies/summary",
	} {
		t.Run(path, func(t *testing.T) {
			methods, documented := doc.Paths[path]
			require.True(t, documented, "generated documentation is missing %q", path)
			require.Contains(t, methods, "get")
			runtimePath := doc.BasePath + strings.ReplaceAll(path, "{id}", ":id")
			found := false
			for _, route := range router.Routes() {
				if route.Method == http.MethodGet && route.Path == runtimePath {
					found = true
					break
				}
			}
			require.True(t, found, "documented URL %q must be registered by NewRouter", runtimePath)
			_, unversioned := doc.Paths[strings.TrimPrefix(path, "/v1")]
			require.False(t, unversioned, "do not retain an unversioned alias in generated documentation")
		})
	}
}
