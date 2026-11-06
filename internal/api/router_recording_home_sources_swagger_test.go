package api

import (
	"encoding/json"
	"net/http"
	"testing"

	swagger "go-gateway/docs/swagger"

	"github.com/stretchr/testify/require"
)

func TestRecordingHomeSourcesSwaggerContract(t *testing.T) {
	var doc struct {
		BasePath string `json:"basePath"`
		Paths    map[string]map[string]struct {
			Parameters []struct {
				Name     string `json:"name"`
				In       string `json:"in"`
				Required bool   `json:"required"`
			} `json:"parameters"`
			Responses map[string]struct {
				Schema struct {
					Ref string `json:"$ref"`
				} `json:"schema"`
			} `json:"responses"`
		} `json:"paths"`
		Definitions map[string]struct {
			Properties map[string]json.RawMessage `json:"properties"`
		} `json:"definitions"`
	}
	require.NoError(t, json.Unmarshal([]byte(swagger.SwaggerInfo.ReadDoc()), &doc))
	path := "/v1/datalink/studio-v2/recording-home/sources"
	require.Contains(t, doc.Paths, path)
	require.NotContains(t, doc.Paths, "/datalink/studio-v2/recording-home/sources")
	require.Contains(t, doc.Paths[path], "get")
	op := doc.Paths[path]["get"]
	require.Len(t, op.Parameters, 1)
	require.Equal(t, "workspace_id", op.Parameters[0].Name)
	require.Equal(t, "query", op.Parameters[0].In)
	require.True(t, op.Parameters[0].Required)
	require.Equal(t, "#/definitions/handlers.RecordingHomeSourcesEnvelope", op.Responses["200"].Schema.Ref)
	for _, status := range []string{"400", "404", "503"} {
		require.Equal(t, "#/definitions/handlers.RecordingHomeSourcesError", op.Responses[status].Schema.Ref)
	}
	props := doc.Definitions["workspace.RecordingHomeSourcesObservation"].Properties
	require.Len(t, props, 7)
	for _, key := range []string{"workspace_id", "state", "source_digest", "reason_code", "devices", "candidates", "excluded"} {
		require.Contains(t, props, key)
	}
	for _, name := range []string{"workspace.RecordingHomeDeviceSummary", "workspace.RecordingHomeCandidate", "workspace.RecordingHomeExcluded"} {
		require.Contains(t, doc.Definitions, name)
		for _, forbidden := range []string{"connection_config", "password", "last_value", "transform_pipeline", "revision_id", "links"} {
			require.NotContains(t, doc.Definitions[name].Properties, forbidden)
		}
	}
	found := false
	for _, route := range NewRouter(&DatalinkServices{}).Routes() {
		if route.Method == http.MethodGet && route.Path == doc.BasePath+path {
			found = true
		}
	}
	require.True(t, found, "generated URL must match always-mounted production route")
}
