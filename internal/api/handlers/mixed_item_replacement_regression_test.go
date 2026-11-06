package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/stretchr/testify/require"
)

func enableMixedFixtureForeignKeys(t *testing.T, f *mixedRuleHandlerFixture) {
	t.Helper()
	f.db.SetMaxOpenConns(1)
	_, err := f.db.ExecContext(t.Context(), `PRAGMA foreign_keys=ON`)
	require.NoError(t, err)
	var enabled int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `PRAGMA foreign_keys`).Scan(&enabled))
	require.Equal(t, 1, enabled)
}

func TestWorkspaceMixedItemMovesIntoRemovedItemAddress(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	enableMixedFixtureForeignKeys(t, f)
	create := `{"id":"replace","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"},{"item_id":"b","address":"D200","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	tags, _ := applyMixedRuleTags(t, f, "replace")
	retained, removed := mixedLinkAt(t, f, "replace", "D100"), mixedLinkAt(t, f, "replace", "D200")
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D200","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response = f.do(t, http.MethodPut, "/workspace/source-rules/replace", update)
	require.Equal(t, http.StatusOK, code, response)
	after := mixedLinkAt(t, f, "replace", "D200")
	require.Equal(t, retained.PointID, after.PointID)
	require.Equal(t, retained.TagID, after.TagID)
	require.Equal(t, retained.MappingID, after.MappingID)
	_, err := f.points.GetByID(t.Context(), removed.PointID)
	require.ErrorIs(t, err, point.ErrPointNotFound)
	managed, err := tags.GetByID(t.Context(), *retained.TagID)
	require.NoError(t, err)
	require.Contains(t, managed.Labels, `"source_rule_address":"D200"`)
}

func TestWorkspaceMixedItemReplacementFailureRestoresOriginalResources(t *testing.T) {
	for _, failAt := range []string{"links", "runtime"} {
		t.Run(failAt, func(t *testing.T) {
			f := newMixedRuleHandlerFixture(t)
			enableMixedFixtureForeignKeys(t, f)
			create := `{"id":"replace","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
			code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
			require.Equal(t, http.StatusCreated, code, response)
			tags, mappings := applyMixedRuleTags(t, f, "replace")
			old := mixedLinkAt(t, f, "replace", "D100")
			oldPoint, err := f.points.GetByID(t.Context(), old.PointID)
			require.NoError(t, err)
			oldTag, err := tags.GetByID(t.Context(), *old.TagID)
			require.NoError(t, err)
			oldMapping, err := mappings.GetByID(t.Context(), *old.MappingID)
			require.NoError(t, err)
			oldRule, err := f.rules.GetByID(t.Context(), "replace")
			require.NoError(t, err)
			if failAt == "links" {
				query := fmt.Sprintf(`CREATE TRIGGER reject_replacement BEFORE INSERT ON source_rule_links WHEN NEW.rule_id='replace' AND NEW.id!='%s' BEGIN SELECT RAISE(ABORT, 'injected link failure'); END`, old.ID)
				_, err = f.db.ExecContext(t.Context(), query)
				require.NoError(t, err)
			} else {
				f.rules.SetShareRuntimeReconciler(sourcerule.ShareRuntimeReconcilerFunc(func(context.Context, sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
					return sourcerule.RuntimeReconcileOutcome{}, errors.New("injected reconcile failure")
				}))
			}
			update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"c","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
			code, failure := f.do(t, http.MethodPut, "/workspace/source-rules/replace", update)
			require.NotEqual(t, http.StatusOK, code)
			require.NotContains(t, fmt.Sprint(failure), "dirty_unknown", "injected failure should be compensated")
			link := mixedLinkAt(t, f, "replace", "D100")
			require.Equal(t, old.ID, link.ID)
			require.Equal(t, old.PointID, link.PointID)
			require.Equal(t, old.TagID, link.TagID)
			require.Equal(t, old.MappingID, link.MappingID)
			p, err := f.points.GetByID(t.Context(), old.PointID)
			require.NoError(t, err)
			require.Equal(t, oldPoint.RevisionID, p.RevisionID)
			tagRecord, err := tags.GetByID(t.Context(), *old.TagID)
			require.NoError(t, err)
			require.Equal(t, oldTag.Labels, tagRecord.Labels)
			mappingRecord, err := mappings.GetByID(t.Context(), *old.MappingID)
			require.NoError(t, err)
			require.Equal(t, oldMapping.PointID, mappingRecord.PointID)
			require.Equal(t, oldMapping.RuleCandidateID, mappingRecord.RuleCandidateID)
			rule, err := f.rules.GetByID(t.Context(), "replace")
			require.NoError(t, err)
			require.Equal(t, oldRule.MixedItems, rule.MixedItems)
			require.Equal(t, oldRule.RevisionID, rule.RevisionID)
			links, err := f.rules.ListLinks(t.Context(), "replace")
			require.NoError(t, err)
			require.Len(t, links, 1)
		})
	}
}

func TestWorkspaceMixedItemMoveIntoRemovedAddressRollsBackOnRuntimeFailure(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	enableMixedFixtureForeignKeys(t, f)
	create := `{"id":"replace","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"},{"item_id":"b","address":"D200","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	tags, mappings := applyMixedRuleTags(t, f, "replace")
	a, b := mixedLinkAt(t, f, "replace", "D100"), mixedLinkAt(t, f, "replace", "D200")
	reconciled := false
	f.rules.SetShareRuntimeReconciler(sourcerule.ShareRuntimeReconcilerFunc(func(context.Context, sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
		reconciled = true
		return sourcerule.RuntimeReconcileOutcome{}, errors.New("injected reconcile failure")
	}))
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D200","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, failure := f.do(t, http.MethodPut, "/workspace/source-rules/replace", update)
	require.NotEqual(t, http.StatusOK, code)
	require.True(t, reconciled)
	require.NotContains(t, fmt.Sprint(failure), "dirty_unknown")
	for _, before := range []*schema.SourceRuleLink{a, b} {
		after := mixedLinkAt(t, f, "replace", before.Address)
		require.Equal(t, before.PointID, after.PointID)
		require.Equal(t, before.TagID, after.TagID)
		require.Equal(t, before.MappingID, after.MappingID)
		p, err := f.points.GetByID(t.Context(), before.PointID)
		require.NoError(t, err)
		require.Equal(t, before.Address, p.Address)
		_, err = tags.GetByID(t.Context(), *before.TagID)
		require.NoError(t, err)
		_, err = mappings.GetByID(t.Context(), *before.MappingID)
		require.NoError(t, err)
	}
}
