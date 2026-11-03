package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/stretchr/testify/require"
)

type tagCandidateWorkspace struct{}

func (tagCandidateWorkspace) ShareOwnershipSnapshot(context.Context) (workspaceID string, deviceIDs []string, err error) {
	return "workspace", []string{"device-1"}, nil
}

type tagCandidateHydration struct{}

func (tagCandidateHydration) CheckHydration(context.Context) (modbusshare.HydrationState, error) {
	return modbusshare.HydrationState{State: modbusshare.HydrationStateReady, Readiness: true, WorkspaceRevision: "workspace-rev"}, nil
}

func TestSourceRuleHandler_TagCandidatesDoNotRevalidateDatabaseOutputs(t *testing.T) {
	f := setupSourceRuleCandidatesFixture(t)
	createRuleCandidatesForDecisionTest(t, f, "rule-name", "40001", 1)
	rule, err := f.repo.GetByID(t.Context(), "rule-name")
	require.NoError(t, err)
	snapshots, err := f.repo.ListCandidateSnapshots(t.Context(), rule.ID, rule.RevisionID)
	require.NoError(t, err)
	for _, snapshot := range snapshots {
		if snapshot.CandidateType == schema.SourceRuleCandidateTypeDatabaseOutputs {
			snapshot.Status = schema.SourceRuleCandidateStatusBlocked
		}
	}
	require.NoError(t, f.repo.ReplaceCandidateSnapshots(t.Context(), snapshots))
	validationCalls := 0
	f.ruleSvc.SetDatabaseTargetMappingReader(sourcerule.DatabaseTargetMappingListFunc(func(context.Context) ([]*schema.DatabaseTargetMapping, error) {
		return []*schema.DatabaseTargetMapping{{ID: "db-mapping", TagID: "tag", ConnectorID: "destination"}}, nil
	}))
	f.ruleSvc.SetDatabaseTargetConnectorValidator(sourcerule.DatabaseTargetConnectorValidatorFunc(func(context.Context, string) (*sourcerule.DatabaseTargetConnectorValidation, error) {
		validationCalls++
		return nil, errors.New("destination unavailable")
	}))
	f.ruleSvc.SetCandidateScope(tagCandidateWorkspace{}, tagCandidateHydration{})
	query := url.Values{"view": {"tags"}, "workspace_id": {"workspace"}, "expected_workspace_revision": {"workspace-rev"}, "revision_id": {rule.RevisionID}}
	response := httptest.NewRecorder()
	f.router.ServeHTTP(response, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/datalink/source-rules/rule-name/candidates?"+query.Encode(), http.NoBody))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	require.Zero(t, validationCalls)
	var payload struct {
		Data struct {
			Tags struct {
				Candidates []any `json:"candidates"`
			} `json:"tags"`
			DatabaseOutputs json.RawMessage `json:"database_outputs"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &payload))
	require.Len(t, payload.Data.Tags.Candidates, 1)
	require.Empty(t, payload.Data.DatabaseOutputs)

	query.Set("expected_workspace_revision", "stale-revision")
	stale := httptest.NewRecorder()
	f.router.ServeHTTP(stale, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/datalink/source-rules/rule-name/candidates?"+query.Encode(), http.NoBody))
	require.Equal(t, http.StatusConflict, stale.Code)
	require.Zero(t, validationCalls)
	query.Set("expected_workspace_revision", "workspace-rev")
	query.Set("workspace_id", "foreign")
	foreign := httptest.NewRecorder()
	f.router.ServeHTTP(foreign, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/datalink/source-rules/rule-name/candidates?"+query.Encode(), http.NoBody))
	require.Equal(t, http.StatusForbidden, foreign.Code)
	require.Zero(t, validationCalls)
}
