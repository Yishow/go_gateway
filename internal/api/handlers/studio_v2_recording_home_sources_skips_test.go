package handlers

import (
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/stretchr/testify/require"
)

func TestRecordingHomeSourcesPersistedSkipShape(t *testing.T) {
	for _, tc := range []struct {
		name, raw, state string
		mapped           bool
	}{
		{"mapped empty array", "[]", "available", true},
		{"mapped null", "null", "unknown", true},
		{"mapped spaced null", "  null  ", "unknown", true},
		{"mapped empty persisted string", "", "unknown", true},
		{"mapped object", "{}", "unknown", true},
		{"mapped scalar", `"40001"`, "unknown", true},
		{"mapped boolean", "true", "unknown", true},
		{"mapped numeric array", "[1]", "unknown", true},
		{"unconfigured empty persisted string", "", "available", false},
		{"unconfigured null", "null", "unknown", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newHomeRealFixture(t)
			rr := sourcerule.NewMemoryRepository()
			rule := &schema.SourceRule{ID: "rule-A", DeviceID: "dev-A", RevisionID: "revision-1", Enabled: true, SkippedAddresses: tc.raw}
			require.NoError(t, rr.Create(t.Context(), rule))
			links, err := f.old.ruleSvc.ListLinks(t.Context(), rule.ID)
			require.NoError(t, err)
			if !tc.mapped {
				links[0].MappingID = nil
				links[0].TagID = nil
				f.old.mappingSvc = mapping.NewService(mapping.NewMemoryRepository())
				f.handler.mappings = f.old.mappingSvc
			}
			require.NoError(t, rr.CreateLinks(t.Context(), links))
			f.old.ruleSvc = sourcerule.NewService(rr, f.old.deviceSvc, f.old.pointSvc, nil)
			spy := &homeRuleMutationSpy{RecordingHomeRuleReader: f.old.ruleSvc}
			f.handler.rules = spy
			before := f.snapshot(t)
			for range 2 {
				out := homeSourcesRequest(f.handler, "?workspace_id="+f.id)
				require.Equal(t, 200, out.Code, out.Body.String())
				var body struct {
					Data map[string]any `json:"data"`
				}
				require.NoError(t, json.Unmarshal(out.Body.Bytes(), &body))
				require.Equal(t, tc.state, body.Data["state"], out.Body.String())
				switch {
				case tc.state == "unknown":
					require.Equal(t, "source-inconsistent", body.Data["reason_code"])
					for _, key := range []string{"devices", "candidates", "excluded", "source_digest"} {
						require.NotContains(t, body.Data, key)
					}
				case tc.mapped:
					require.Len(t, body.Data["candidates"], 1)
				default:
					require.Empty(t, body.Data["candidates"])
				}
			}
			require.Zero(t, spy.writes)
			require.Equal(t, before, f.snapshot(t))
		})
	}
}
