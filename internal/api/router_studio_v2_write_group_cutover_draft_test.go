package api

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Saving a draft does not stop its applied wide revision. Ordinary Apply must
// still require a reviewed cutover instead of starting a second owner.
func TestProfileApplyRefusesAnAppliedWideGroupWithANewerDraft(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			r := newCutoverRig(t, kind)
			wideID, longID := r.wide["id"].(string), r.long["id"].(string)
			applied := getGroup(t, r.f, wideID)["applied_revision"]
			draft := getGroup(t, r.f, wideID)
			draft["name"] = "An unapplied display-name draft"
			response := performJSONRequest(t, r.f.router, http.MethodPut, writeGroupsPath+"/"+wideID,
				groupEditRequest(r.f.writeGroupRouterFixture, getGroupData(t, r.f, wideID), draft))
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			saved := getGroup(t, r.f, wideID)
			require.Equal(t, "draft", saved["status"])
			require.Equal(t, applied, saved["applied_revision"])
			require.NotEqual(t, saved["revision"], applied)
			before := r.cutoverCounts()

			code, body := r.post(writeGroupsPath+"/"+longID+"/apply", scopeBody(t, r.f, longID))
			require.Equal(t, http.StatusConflict, code, "%v", body)
			require.Equal(t, "cutover_required", errorOf(t, body)["code"])
			require.Equal(t, before, r.cutoverCounts())
			require.Equal(t, "", getGroup(t, r.f, longID)["applied_revision"])
			require.Equal(t, saved, getGroup(t, r.f, wideID), "the draft and its applied writer are unchanged")
			require.Zero(t, r.localCount(`SELECT COUNT(*) FROM write_group_versions WHERE group_id = ?`, longID))

			r.accept(r.wideEffective.Add(2*time.Second), 12.5)
			r.advance(r.wideEffective.Add(16 * time.Second))
			r.waitFor(func() bool { return r.wideRows() == 1 }, "the original applied wide writer still delivers")
			require.Empty(t, r.longRecords())
		})
	}
}
