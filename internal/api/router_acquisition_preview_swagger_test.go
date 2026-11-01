package api

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionPreviewSwaggerContract(t *testing.T) {
	data, err := os.ReadFile("../../docs/swagger/swagger.json")
	require.NoError(t, err)
	var doc struct {
		Paths map[string]map[string]struct {
			Parameters []struct {
				Name     string `json:"name"`
				Required bool   `json:"required"`
			}
			Responses map[string]json.RawMessage
		}
	}
	require.NoError(t, json.Unmarshal(data, &doc))
	path, ok := doc.Paths["/datalink/studio-v2/workspace/acquisition-plan-preview"]
	require.True(t, ok)
	get, ok := path["get"]
	require.True(t, ok)
	require.Len(t, path, 1)
	required := map[string]bool{}
	for _, parameter := range get.Parameters {
		required[parameter.Name] = parameter.Required
	}
	require.Equal(t, map[string]bool{"workspace_id": true, "device_id": true, "expected_workspace_revision": true, "expected_preview_revision": false}, required)
	for _, status := range []string{"200", "400", "403", "409", "503"} {
		require.Contains(t, get.Responses, status)
	}
}
