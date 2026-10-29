package handlers

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"testing"
)

func TestStudioLocatorEditorRequiresIndependentExpectedBindings(t *testing.T) {
	base := "entry=locator&workspace_id=W&step=4&group_id=G&editor=1"
	for _, q := range []string{base, base + "&expected_group_revision=R", base + "&expected_workspace_revision=W1", base + "&expected_workspace_revision=W1&expected_group_revision=", base + "&expected_workspace_revision=W1&expected_group_revision=R&editor=1", "entry=locator&workspace_id=W&step=1&device_id=D&editor=1&expected_workspace_revision=W1&expected_group_revision=R", "entry=locator&workspace_id=W&step=4&group_id=G&expected_workspace_revision=W1&expected_group_revision=R"} {
		_, err := parseStudioLocatorQuery(q)
		require.Error(t, err, q)
	}
	request, err := parseStudioLocatorQuery(base + "&expected_workspace_revision=&expected_group_revision=R&group_revision=old")
	require.NoError(t, err)
	require.True(t, request.Editor)
	require.Equal(t, "", request.ExpectedWorkspaceRevision)
	require.Equal(t, "R", request.ExpectedGroupRevision)
	require.Equal(t, "old", *request.GroupRevision)
}
func TestStudioLocatorEditorObservedEmptyCollectionsAndRedaction(t *testing.T) {
	observation := &studioLocatorObservation{Editor: &studioLocatorEditorObservation{Missing: map[string]string{"D/P/T": "unowned"}}}
	entry, err := projectStudioLocator(&workspace.Record{ID: "W"}, observation)
	require.NoError(t, err)
	raw, err := json.Marshal(entry)
	require.NoError(t, err)
	var fields map[string]any
	require.NoError(t, json.Unmarshal(raw, &fields))
	require.Equal(t, true, fields["editor_hydrated"])
	require.Contains(t, fields, "database_config")
	require.Nil(t, fields["database_config"])
	for _, key := range []string{"devices", "source_rules", "mappings", "database_targets"} {
		require.Equal(t, []any{}, fields[key], key)
	}
	observation.Editor.Connector = &schema.DatabaseConnector{ID: "C", LastCheckError: "password=one", LastSchemaEnsureError: "password=two", LastWriteError: "password=three", LastFlushError: "password=four", ConnectionConfig: `{"database":"native","password":"not-returned"}`}
	entry, err = projectStudioLocator(&workspace.Record{ID: "W"}, observation)
	require.NoError(t, err)
	raw, err = json.Marshal(entry)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "password")
	require.NotContains(t, string(raw), "not-returned")
}
func TestStudioLocatorEditorMissingMarkersAndProviderObservationsCompare(t *testing.T) {
	a := &studioLocatorEditorObservation{Missing: map[string]string{"D/P/T": "point"}}
	b := &studioLocatorEditorObservation{Missing: map[string]string{"D/P/T": "point"}}
	require.True(t, sameStudioLocatorEditor(a, b))
	b.Missing["D/P/T"] = "mapping"
	require.False(t, sameStudioLocatorEditor(a, b))
	b.Missing["D/P/T"] = "point"
	b.Targets = []*schema.DatabaseTargetMapping{{ID: "target"}}
	require.False(t, sameStudioLocatorEditor(a, b))
}
