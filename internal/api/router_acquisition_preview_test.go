package api

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionPreviewProductionRouteScopeAndNoWrites(t *testing.T) {
	r := newRuntimeSourceRig(t)
	record, err := r.f.workspace.ReadSnapshot(t.Context())
	require.NoError(t, err)
	base := url.Values{"workspace_id": {record.ID}, "device_id": {"device-A"}, "expected_workspace_revision": {record.DatabaseSetupRevision}}
	before := locatorTotal(t, r)
	response := performJSONRequest(t, r.router, http.MethodGet, studioLocatorPath+"/acquisition-plan-preview?"+base.Encode(), nil)
	data := locatorData(t, response)
	require.Equal(t, "unavailable", data["status"])
	require.NotContains(t, data, "preview_revision")
	require.NotContains(t, data, "wire_attempts")
	require.Equal(t, before, locatorTotal(t, r), "preview cannot bootstrap or repair SQL")
	for _, tc := range []struct {
		key, value string
		status     int
	}{
		{"workspace_id", "foreign", 403}, {"device_id", "foreign", 403}, {"expected_workspace_revision", "old", 409},
		{"expected_preview_revision", "old", 409}, {"endpoint", "secret", 400},
	} {
		q := url.Values{}
		for key, values := range base {
			q[key] = append([]string{}, values...)
		}
		q.Set(tc.key, tc.value)
		response = performJSONRequest(t, r.router, http.MethodGet, studioLocatorPath+"/acquisition-plan-preview?"+q.Encode(), nil)
		require.Equal(t, tc.status, response.Code, response.Body.String())
		require.NotContains(t, response.Body.String(), "secret")
	}
}
