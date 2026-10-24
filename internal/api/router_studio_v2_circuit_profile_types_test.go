package api

import (
	"net/http"
	"testing"

	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

// A counter label cannot make a string or bool tag arithmetically countable.
// Readiness must reject it before Apply, using the same limit as runtime.
func TestCumulativeProfileRejectsNonNumericSourceBeforeApply(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		for _, tagType := range []string{"string", "bool"} {
			t.Run(kind+"/"+tagType, func(t *testing.T) {
				f := newManagedDestinationFixture(t, kind, tagType)
				insertCounterMeasurement(t, f, "counter-text", "counter", "count")
				group := ensureProfileGroup(t, f, "cumulative", 10, "minute", "counter-text")
				id := group["id"].(string)
				require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, id, previewLong(t, f, id)).Status)
				ready, err := f.groups.Readiness(t.Context(), id)
				require.NoError(t, err)
				require.False(t, ready.Ready, "a string/bool source must not become ready from a counter measurement label")
				codes := make([]string, 0, len(ready.Issues))
				for _, issue := range ready.Issues {
					codes = append(codes, issue.Code)
				}
				require.Contains(t, codes, "counter-type-unsupported")
				response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
				require.GreaterOrEqual(t, response.Code, http.StatusBadRequest, response.Body.String())
				require.Empty(t, appliedRevisionOf(t, f, id))
			})
		}
	}
}
