package handlers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionPreviewStrictQuery(t *testing.T) {
	base := "workspace_id=w&device_id=d&expected_workspace_revision="
	for _, query := range []string{"", "workspace_id=w&device_id=d", base + "&host=localhost", base + "&device_id=d", base + "&expected_preview_revision=", base + "&member_id=p", base + "&x=%zz"} {
		_, err := parseAcquisitionPreviewQuery(query)
		require.Error(t, err, query)
	}
	parsed, err := parseAcquisitionPreviewQuery(base)
	require.NoError(t, err)
	require.Equal(t, "", parsed.workspaceRevision)
	_, err = parseAcquisitionPreviewQuery(base + "&expected_preview_revision=head")
	require.NoError(t, err)
}
