package handlers

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"testing"
)

func TestStudioLocatorProjectionNativeFieldsAndEmptyScope(t *testing.T) {
	d := &schema.Device{ID: "D", ConnectionConfig: `{"host":"native-editor-host"}`, LastTestError: "password=secret"}
	rule := &schema.SourceRule{ID: "R", DeviceID: "D", RevisionID: "rev"}
	got, err := projectStudioLocator(&workspace.Record{ID: "W"}, &studioLocatorObservation{Device: d, Rule: rule})
	require.NoError(t, err)
	require.Equal(t, "", got.WorkspaceRevision)
	require.False(t, got.HydrationScope.CompleteWorkspace)
	require.Empty(t, got.RuleHeads)
	require.Empty(t, got.GroupHeads)
	require.NotNil(t, got.RevalidatedMappingIDs)
	require.Equal(t, []string{}, got.SourceRules[0].SkippedAddresses)
	require.Equal(t, d.ConnectionConfig, got.Devices[0].ConnectionConfig)
	require.Empty(t, got.Devices[0].LastTestError)
	require.Equal(t, "password=secret", d.LastTestError)
	raw, err := json.Marshal(got)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "password=secret")
	require.NotContains(t, string(raw), "runtime_apply")
	require.NotContains(t, string(raw), `"mappings"`)
}
func locatorMappingObservation(pipeline string) *studioLocatorObservation {
	return &studioLocatorObservation{Mapping: &schema.Mapping{ID: "M", TransformPipeline: pipeline}, Rule: &schema.SourceRule{ID: "R", DeviceID: "D"}, Point: &schema.Point{ID: "P", DataType: schema.DataTypeFloat32}, Tag: &schema.Tag{ID: "T", DataType: schema.DataTypeFloat32}, Link: &schema.SourceRuleLink{Address: "40001"}}
}
func TestStudioLocatorMappingProjectionRefusesUnrepresentablePipeline(t *testing.T) {
	for _, raw := range []string{`{}`, `[{"type":"script","params":{"password":"secret"}}]`, `[{"type":"cast","params":{"target_type":"float32"}},{"type":"scale","params":{"scale":2,"offset":0}}]`, `[{"type":"scale","params":{"scale":2,"offset":0}},{"type":"scale","params":{"scale":3,"offset":0}}]`, `[{"type":"scale","params":{"scale":2,"offset":0,"extra":"secret"}}]`} {
		t.Run(raw, func(t *testing.T) {
			_, err := projectStudioLocatorMapping("W", locatorMappingObservation(raw))
			require.Error(t, err)
		})
	}
	for _, raw := range []string{"", `null`, `[]`, `[{"type":"scale","params":{"scale":0.1,"offset":2}},{"type":"cast","params":{"target_type":"float32"}}]`} {
		t.Run(raw, func(t *testing.T) {
			got, err := projectStudioLocatorMapping("W", locatorMappingObservation(raw))
			require.NoError(t, err)
			require.Equal(t, "float32", got.TargetType)
		})
	}
}
