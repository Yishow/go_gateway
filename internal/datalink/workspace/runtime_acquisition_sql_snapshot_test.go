package workspace

import (
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestRuntimeAcquisitionSourceGraphUsesRealSQLProjection(t *testing.T) {
	f := newAcquisitionSQLFixture(t)
	projection, graph := f.graph(t)
	require.Equal(t, []optimizer.NativeRange{{Area: "HR", Start: 1, End: 2}}, graph.Excluded())
	require.Len(t, graph.SemanticRevision(), 64)
	ruleRevision, ok := graph.RuleSemanticRevision("r")
	require.True(t, ok)
	require.Len(t, ruleRevision, 64)
	require.Equal(t, "rule-original", graph.Rules()[0].RevisionID)
	graph.Points()[0].Address = "40099"
	graph.Rules()[0].Enabled = false
	graph.Excluded()[0].End = 99
	still, ok := projection.AcquisitionSourceGraph("d")
	require.True(t, ok)
	require.Equal(t, graph.SemanticRevision(), still.SemanticRevision())
	require.Equal(t, "40001", still.Points()[0].Address)
	called := false
	require.NoError(t, f.service.AcceptAcquisitionSourceGraph(t.Context(), graph, func() error { called = true; return nil }))
	require.True(t, called)
}

func TestRuntimeAcquisitionSourceGraphRetainsOriginalInputAfterCallerMutation(t *testing.T) {
	f := newAcquisitionSQLFixture(t)
	projection, graph := f.graph(t)
	before := graph.SemanticRevision()
	projection.Points[0].Address = "40099"
	projection.Rules[0].SkippedAddresses = `["40008"]`
	projection.RuleLinks[0].PointID = "foreign"
	projection.Devices[0].ConnectionConfig = `{"timeout":99}`
	require.Equal(t, "40001", graph.Points()[0].Address)
	require.Equal(t, `["40002"]`, graph.Rules()[0].SkippedAddresses)
	require.Equal(t, before, graph.SemanticRevision())
	require.NoError(t, f.service.AcceptAcquisitionSourceGraph(t.Context(), graph, func() error { return nil }))
}

func TestRuntimeAcquisitionSourceGraphRejectsPartialAndJSONCopies(t *testing.T) {
	for _, kind := range []string{"point-slice", "link-slice", "rule-slice", "point-mutation", "rule-revision", "json-copy", "hand-built"} {
		t.Run(kind, func(t *testing.T) {
			f := newAcquisitionSQLFixture(t)
			projection, _ := f.graph(t)
			modified := *projection
			switch kind {
			case "point-slice":
				modified.Points = modified.Points[:1]
			case "link-slice":
				modified.RuleLinks = modified.RuleLinks[:1]
			case "rule-slice":
				modified.Rules = nil
			case "point-mutation":
				modified.Points[0].Enabled = false
			case "rule-revision":
				modified.Rules[0].RevisionID = "invented"
			case "json-copy":
				raw, err := json.Marshal(projection)
				require.NoError(t, err)
				modified = RuntimeProjection{}
				require.NoError(t, json.Unmarshal(raw, &modified))
			case "hand-built":
				modified.acquisitionGraphs = nil
			}
			graph, ok := modified.AcquisitionSourceGraph("d")
			require.False(t, ok)
			require.False(t, graph.Complete())
		})
	}
}

func TestRuntimeAcquisitionSourceGraphRejectsIncompleteSQLSource(t *testing.T) {
	for _, tc := range []struct{ name, query string }{
		{"standalone-disabled", `INSERT INTO points(id,device_id,name,address,data_type,enabled) VALUES('unlinked','d','u','40002','uint16',0)`},
		{"standalone-active", `INSERT INTO points(id,device_id,name,address,data_type) VALUES('unlinked','d','u','40002','uint16')`},
		{"missing-link", `DELETE FROM source_rule_links WHERE id='lb'`},
		{"missing-point", `DELETE FROM points WHERE id='b'`},
		{"layout-not-loaded-by-legacy-reader", `UPDATE source_rules SET layout_mode='mixed',mixed_items='[]'`},
		{"inconsistent-disabled-rule", `UPDATE source_rules SET enabled=0`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newAcquisitionSQLFixture(t)
			_, err := f.db.ExecContext(t.Context(), tc.query)
			require.NoError(t, err)
			projection, err := f.service.RuntimeProjection(t.Context())
			require.NoError(t, err)
			graph, ok := projection.AcquisitionSourceGraph("d")
			require.False(t, ok)
			require.NotEmpty(t, graph.Reason())
		})
	}
}

func TestRuntimeAcquisitionSourceGraphRetainsEmptyAndSkipOnlyRules(t *testing.T) {
	for _, kind := range []string{"zero-rules", "skip-only", "disabled"} {
		t.Run(kind, func(t *testing.T) {
			f := newAcquisitionSQLFixture(t)
			if kind == "disabled" {
				_, err := f.db.ExecContext(t.Context(), `UPDATE source_rules SET enabled=0; UPDATE points SET enabled=0`)
				require.NoError(t, err)
			} else {
				_, err := f.db.ExecContext(t.Context(), `DELETE FROM source_rule_links; DELETE FROM points`)
				require.NoError(t, err)
				query := `DELETE FROM source_rules`
				if kind == "skip-only" {
					query = `UPDATE source_rules SET skipped_addresses='["40001","40002","40003"]'`
				}
				_, err = f.db.ExecContext(t.Context(), query)
				require.NoError(t, err)
			}
			_, graph := f.graph(t)
			if kind == "zero-rules" {
				require.Empty(t, graph.Excluded())
			} else {
				require.Len(t, graph.Excluded(), 3)
			}
		})
	}
}

func TestRuntimeAcquisitionSourceGraphIgnoresOutputSelectorsNotSourceContent(t *testing.T) {
	f := newAcquisitionSQLFixture(t)
	projection, original := f.graph(t)
	tag := "not-a-source"
	projection.RuleLinks[0].TagID, projection.RuleLinks[0].MappingID = &tag, &tag
	projection.Rules[0].ShareEnabled = true
	projection.Mappings = []*schema.Mapping{{ID: "unrelated-output"}}
	graph, ok := projection.AcquisitionSourceGraph("d")
	require.True(t, ok)
	require.Equal(t, original.SemanticRevision(), graph.SemanticRevision())
	_, err := f.db.ExecContext(t.Context(), `UPDATE source_rules SET share_enabled=1,revision_id='share-only-new'; UPDATE points SET last_error='telemetry-only'`)
	require.NoError(t, err)
	require.NoError(t, f.service.AcceptAcquisitionSourceGraph(t.Context(), graph, func() error { return nil }))
	_, next := f.graph(t)
	require.Equal(t, graph.SemanticRevision(), next.SemanticRevision())
	require.Equal(t, "rule-original", graph.Rules()[0].RevisionID)
	require.Equal(t, "share-only-new", next.Rules()[0].RevisionID)
}
