package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

func mixedTagCandidates(t *testing.T, f *mixedRuleHandlerFixture, ruleID string) []schema.SourceRuleTagCandidate {
	t.Helper()
	snapshots, err := f.rules.ListCandidateSnapshots(t.Context(), ruleID)
	require.NoError(t, err)
	for _, snapshot := range snapshots {
		if snapshot.CandidateType == schema.SourceRuleCandidateTypeTags {
			var payload struct {
				Candidates []schema.SourceRuleTagCandidate `json:"candidates"`
			}
			require.NoError(t, json.Unmarshal([]byte(snapshot.Payload), &payload))
			return payload.Candidates
		}
	}
	t.Fatal("missing tag candidates")
	return nil
}

func TestWorkspaceMixedReplacementDoesNotReuseAnotherItemsReviewDecision(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	enableMixedFixtureForeignKeys(t, f)
	create := `{"id":"replace","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	tags := tag.NewService(tag.NewSQLRepository(f.db))
	mappings := mapping.NewServiceWithTagResolver(mapping.NewSQLRepository(f.db), tags.GetByID)
	f.rules.SetTagMappingServices(tags, mappings)
	before := mixedTagCandidates(t, f, "replace")
	require.Len(t, before, 1)
	_, err := f.rules.UpsertTagReviewDecision(t.Context(), "replace", sourcerule.UpsertTagReviewDecisionRequest{CandidateID: before[0].ID, Action: schema.SourceRuleTagReviewDecisionActionRename, TagKey: "factory.old"})
	require.NoError(t, err)
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"c","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response = f.do(t, http.MethodPut, "/workspace/source-rules/replace", update)
	require.Equal(t, http.StatusOK, code, response)
	after := mixedTagCandidates(t, f, "replace")
	require.Len(t, after, 1)
	require.NotEqual(t, before[0].ID, after[0].ID, "a different item at the same address must not inherit its predecessor's review")
	oldDecision, err := f.rules.GetTagReviewDecision(t.Context(), "replace", before[0].ID)
	require.NoError(t, err)
	require.True(t, oldDecision.Stale)
	newDecision, err := f.rules.GetTagReviewDecision(t.Context(), "replace", after[0].ID)
	require.NoError(t, err)
	require.Nil(t, newDecision)
	rule, err := f.rules.GetByID(t.Context(), "replace")
	require.NoError(t, err)
	result, err := f.rules.ApplyTagCandidates(t.Context(), "replace", sourcerule.ApplyTagCandidatesRequest{RevisionID: rule.RevisionID, CandidateIDs: []string{after[0].ID}})
	require.NoError(t, err)
	require.Equal(t, "applied", result.Results[0].Status)
	link := mixedLinkAt(t, f, "replace", "D100")
	createdTag, err := tags.GetByID(t.Context(), *link.TagID)
	require.NoError(t, err)
	require.NotEqual(t, "factory.old", createdTag.Key)
}

func TestWorkspaceMixedMoveRollbackRestoresReviewDecision(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	enableMixedFixtureForeignKeys(t, f)
	create := `{"id":"move","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	before := mixedTagCandidates(t, f, "move")
	require.Len(t, before, 1)
	_, err := f.rules.UpsertTagReviewDecision(t.Context(), "move", sourcerule.UpsertTagReviewDecisionRequest{CandidateID: before[0].ID, Action: schema.SourceRuleTagReviewDecisionActionRename, TagKey: "factory.old"})
	require.NoError(t, err)
	reconciled := false
	f.rules.SetShareRuntimeReconciler(sourcerule.ShareRuntimeReconcilerFunc(func(context.Context, sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
		reconciled = true
		return sourcerule.RuntimeReconcileOutcome{}, errors.New("injected runtime failure")
	}))
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D200","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, failure := f.do(t, http.MethodPut, "/workspace/source-rules/move", update)
	require.NotEqual(t, http.StatusOK, code)
	require.True(t, reconciled)
	require.NotContains(t, fmt.Sprint(failure), "dirty_unknown")
	after := mixedTagCandidates(t, f, "move")
	require.Len(t, after, 1)
	require.Equal(t, before[0].ID, after[0].ID)
	decision, err := f.rules.GetTagReviewDecision(t.Context(), "move", before[0].ID)
	require.NoError(t, err)
	require.False(t, decision.Stale, "failed move must not invalidate the restored item's decision")
	require.Equal(t, "factory.old", decision.TagKey)
}

func TestWorkspaceMixedMoveRollbackKeepsConcurrentReviewDecision(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	enableMixedFixtureForeignKeys(t, f)
	create := `{"id":"move","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	before := mixedTagCandidates(t, f, "move")
	require.Len(t, before, 1)
	_, err := f.rules.UpsertTagReviewDecision(t.Context(), "move", sourcerule.UpsertTagReviewDecisionRequest{CandidateID: before[0].ID, Action: schema.SourceRuleTagReviewDecisionActionRename, TagKey: "factory.old"})
	require.NoError(t, err)
	changed := false
	f.rules.SetShareRuntimeReconciler(sourcerule.ShareRuntimeReconcilerFunc(func(ctx context.Context, _ sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
		_, updateErr := f.db.ExecContext(ctx, `UPDATE source_rule_tag_review_decisions SET tag_key='factory.concurrent',stale=0,stale_revision_id='',stale_at=NULL WHERE source_rule_id='move' AND candidate_id=?`, before[0].ID)
		changed = updateErr == nil
		return sourcerule.RuntimeReconcileOutcome{}, errors.Join(errors.New("injected runtime failure"), updateErr)
	}))
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D200","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, failure := f.do(t, http.MethodPut, "/workspace/source-rules/move", update)
	require.NotEqual(t, http.StatusOK, code)
	require.True(t, changed)
	require.Contains(t, fmt.Sprint(failure), "dirty_unknown")
	decision, err := f.rules.GetTagReviewDecision(t.Context(), "move", before[0].ID)
	require.NoError(t, err)
	require.Equal(t, "factory.concurrent", decision.TagKey)
	require.False(t, decision.Stale)
}

func TestWorkspaceMixedReviewDecisionSurvivesCompatibleEdit(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	enableMixedFixtureForeignKeys(t, f)
	create := `{"id":"move","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	before := mixedTagCandidates(t, f, "move")
	require.Len(t, before, 1)
	_, err := f.rules.UpsertTagReviewDecision(t.Context(), "move", sourcerule.UpsertTagReviewDecisionRequest{CandidateID: before[0].ID, Action: schema.SourceRuleTagReviewDecisionActionRename, TagKey: "factory.old"})
	require.NoError(t, err)
	update := `{"naming_prefix":"NEW","layout_mode":"mixed","mixed_items":[{"item_id":"a","name":"changed display name","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response = f.do(t, http.MethodPut, "/workspace/source-rules/move", update)
	require.Equal(t, http.StatusOK, code, response)
	after := mixedTagCandidates(t, f, "move")
	require.Len(t, after, 1)
	require.Equal(t, before[0].ID, after[0].ID)
	decision, err := f.rules.GetTagReviewDecision(t.Context(), "move", before[0].ID)
	require.NoError(t, err)
	require.False(t, decision.Stale)
	require.Equal(t, "factory.old", decision.TagKey)
}

func TestWorkspaceMixedSwapDoesNotTransferReviewDecisions(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	enableMixedFixtureForeignKeys(t, f)
	create := `{"id":"swap","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"},{"item_id":"b","address":"D200","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	before := mixedTagCandidates(t, f, "swap")
	require.Len(t, before, 2)
	for _, candidate := range before {
		_, err := f.rules.UpsertTagReviewDecision(t.Context(), "swap", sourcerule.UpsertTagReviewDecisionRequest{CandidateID: candidate.ID, Action: schema.SourceRuleTagReviewDecisionActionRename, TagKey: "factory." + candidate.Address})
		require.NoError(t, err)
	}
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D200","data_type":"uint16","semantic_kind":"gauge"},{"item_id":"b","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response = f.do(t, http.MethodPut, "/workspace/source-rules/swap", update)
	require.Equal(t, http.StatusOK, code, response)
	after := mixedTagCandidates(t, f, "swap")
	require.Len(t, after, 2)
	for _, original := range before {
		decision, err := f.rules.GetTagReviewDecision(t.Context(), "swap", original.ID)
		require.NoError(t, err)
		require.True(t, decision.Stale)
		for _, candidate := range after {
			require.NotEqual(t, original.ID, candidate.ID)
		}
	}
	for _, candidate := range after {
		decision, err := f.rules.GetTagReviewDecision(t.Context(), "swap", candidate.ID)
		require.NoError(t, err)
		require.Nil(t, decision)
	}
}
