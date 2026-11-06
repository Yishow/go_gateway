package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

func mixedLinkAt(t *testing.T, f *mixedRuleHandlerFixture, ruleID, address string) *schema.SourceRuleLink {
	t.Helper()
	links, err := f.rules.ListLinks(t.Context(), ruleID)
	require.NoError(t, err)
	for _, link := range links {
		if link.Address == address {
			return link
		}
	}
	t.Fatalf("no source-rule link at %s", address)
	return nil
}

func applyMixedRuleTags(t *testing.T, f *mixedRuleHandlerFixture, ruleID string) (*tag.Service, *mapping.Service) {
	t.Helper()
	tags := tag.NewService(tag.NewSQLRepository(f.db))
	mappings := mapping.NewServiceWithTagResolver(mapping.NewSQLRepository(f.db), tags.GetByID)
	f.rules.SetTagMappingServices(tags, mappings)
	rule, err := f.rules.GetByID(t.Context(), ruleID)
	require.NoError(t, err)
	snapshots, err := f.rules.ListCandidateSnapshots(t.Context(), ruleID)
	require.NoError(t, err)
	var ids []string
	for _, snapshot := range snapshots {
		if snapshot.CandidateType != schema.SourceRuleCandidateTypeTags {
			continue
		}
		var payload struct {
			Candidates []schema.SourceRuleTagCandidate `json:"candidates"`
		}
		require.NoError(t, json.Unmarshal([]byte(snapshot.Payload), &payload))
		for _, candidate := range payload.Candidates {
			ids = append(ids, candidate.ID)
		}
	}
	require.NotEmpty(t, ids)
	result, err := f.rules.ApplyTagCandidates(t.Context(), ruleID, sourcerule.ApplyTagCandidatesRequest{RevisionID: rule.RevisionID, CandidateIDs: ids})
	require.NoError(t, err)
	for _, entry := range result.Results {
		require.Equal(t, "applied", entry.Status, entry.Error)
	}
	return tags, mappings
}

func TestWorkspaceMixedItemAddressChangeKeepsAppliedIdentity(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	create := `{"id":"move","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"value","address":"D100","data_type":"uint32","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	tags, mappings := applyMixedRuleTags(t, f, "move")
	before := mixedLinkAt(t, f, "move", "D100")
	require.NotNil(t, before.TagID)
	require.NotNil(t, before.MappingID)
	_, err := f.db.ExecContext(t.Context(), `UPDATE points SET format_mode='override',format_override='ABCD',format_version=1 WHERE id=?`, before.PointID)
	require.NoError(t, err)
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"value","address":"D200","data_type":"uint32","semantic_kind":"gauge"}]}`
	code, response = f.do(t, http.MethodPut, "/workspace/source-rules/move", update)
	require.Equal(t, http.StatusOK, code, response)
	after := mixedLinkAt(t, f, "move", "D200")
	require.Equal(t, before.ID, after.ID)
	require.Equal(t, before.PointID, after.PointID)
	require.Equal(t, before.TagID, after.TagID)
	require.Equal(t, before.MappingID, after.MappingID)
	point, err := f.points.GetByID(t.Context(), before.PointID)
	require.NoError(t, err)
	require.Equal(t, "D200", point.Address)
	require.Equal(t, schema.PointFormatOverride, point.FormatMode)
	require.Equal(t, "ABCD", point.FormatOverride)
	mapped, err := mappings.GetByID(t.Context(), *before.MappingID)
	require.NoError(t, err)
	require.Equal(t, before.PointID, mapped.PointID)
	managedTag, err := tags.GetByID(t.Context(), *before.TagID)
	require.NoError(t, err)
	require.Contains(t, managedTag.Labels, `"source_rule_address":"D200"`)
}

func TestWorkspaceMixedItemMoveKeepsNullableLegacyMetadata(t *testing.T) {
	for _, nullable := range []string{"tag", "mapping"} {
		t.Run(nullable, func(t *testing.T) {
			f := newMixedRuleHandlerFixture(t)
			enableMixedFixtureForeignKeys(t, f)
			create := `{"id":"move","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
			code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
			require.Equal(t, http.StatusCreated, code, response)
			_, _ = applyMixedRuleTags(t, f, "move")
			before := mixedLinkAt(t, f, "move", "D100")
			var err error
			if nullable == "tag" {
				_, err = f.db.ExecContext(t.Context(), `UPDATE tags SET unit=NULL,description=NULL WHERE id=?`, *before.TagID)
			} else {
				_, err = f.db.ExecContext(t.Context(), `UPDATE mappings SET proposed_signature=NULL,last_applied_signature=NULL,blocking_reason=NULL WHERE id=?`, *before.MappingID)
			}
			require.NoError(t, err)
			update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D200","data_type":"uint16","semantic_kind":"gauge"}]}`
			code, response = f.do(t, http.MethodPut, "/workspace/source-rules/move", update)
			require.Equal(t, http.StatusOK, code, response)
			after := mixedLinkAt(t, f, "move", "D200")
			require.Equal(t, before.PointID, after.PointID)
			require.Equal(t, before.TagID, after.TagID)
			require.Equal(t, before.MappingID, after.MappingID)
		})
	}
}

func TestWorkspaceMixedItemAddressSwapKeepsIdentities(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	create := `{"id":"swap","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"},{"item_id":"b","address":"D200","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	a, b := mixedLinkAt(t, f, "swap", "D100"), mixedLinkAt(t, f, "swap", "D200")
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D200","data_type":"uint16","semantic_kind":"gauge"},{"item_id":"b","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response = f.do(t, http.MethodPut, "/workspace/source-rules/swap", update)
	require.Equal(t, http.StatusOK, code, response)
	require.Equal(t, a.PointID, mixedLinkAt(t, f, "swap", "D200").PointID)
	require.Equal(t, b.PointID, mixedLinkAt(t, f, "swap", "D100").PointID)
}

func TestWorkspaceMixedItemReplacedAtSameAddressGetsNewIdentity(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	enableMixedFixtureForeignKeys(t, f)
	create := `{"id":"replace","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"},{"item_id":"b","address":"D200","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	tags, mappings := applyMixedRuleTags(t, f, "replace")
	old := mixedLinkAt(t, f, "replace", "D100")
	retained := mixedLinkAt(t, f, "replace", "D200")
	require.NotNil(t, old.TagID)
	require.NotNil(t, old.MappingID)

	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"c","address":"D100","data_type":"uint16","semantic_kind":"gauge"},{"item_id":"b","address":"D200","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response = f.do(t, http.MethodPut, "/workspace/source-rules/replace", update)
	require.Equal(t, http.StatusOK, code, response)
	after := mixedLinkAt(t, f, "replace", "D100")
	require.NotEqual(t, old.ID, after.ID)
	require.NotEqual(t, old.PointID, after.PointID)
	require.Nil(t, after.TagID)
	require.Nil(t, after.MappingID)
	require.Equal(t, retained.PointID, mixedLinkAt(t, f, "replace", "D200").PointID)
	_, err := f.points.GetByID(t.Context(), old.PointID)
	require.ErrorIs(t, err, point.ErrPointNotFound)
	_, err = mappings.GetByID(t.Context(), *old.MappingID)
	require.ErrorIs(t, err, mapping.ErrMappingNotFound)
	_, err = tags.GetByID(t.Context(), *old.TagID)
	require.ErrorIs(t, err, tag.ErrTagNotFound)
}

func TestWorkspaceMixedItemMoveRollsBackAfterLinkFailure(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	create := `{"id":"move","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"value","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	tags, mappings := applyMixedRuleTags(t, f, "move")
	before := mixedLinkAt(t, f, "move", "D100")
	tagBefore, err := tags.GetByID(t.Context(), *before.TagID)
	require.NoError(t, err)
	mapBefore, err := mappings.GetByID(t.Context(), *before.MappingID)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `CREATE TRIGGER reject_moved_link BEFORE INSERT ON source_rule_links WHEN NEW.rule_id='move' AND NEW.address='D200' BEGIN SELECT RAISE(ABORT, 'injected link failure'); END`)
	require.NoError(t, err)
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"value","address":"D200","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, _ = f.do(t, http.MethodPut, "/workspace/source-rules/move", update)
	require.NotEqual(t, http.StatusOK, code)
	after := mixedLinkAt(t, f, "move", "D100")
	require.Equal(t, before.ID, after.ID)
	require.Equal(t, before.PointID, after.PointID)
	p, err := f.points.GetByID(t.Context(), before.PointID)
	require.NoError(t, err)
	require.Equal(t, "D100", p.Address)
	tagAfter, err := tags.GetByID(t.Context(), *before.TagID)
	require.NoError(t, err)
	require.Equal(t, tagBefore.Labels, tagAfter.Labels)
	mapAfter, err := mappings.GetByID(t.Context(), *before.MappingID)
	require.NoError(t, err)
	require.Equal(t, mapBefore.RuleCandidateID, mapAfter.RuleCandidateID)
}

func TestWorkspaceMixedItemSwapRollsBackAfterRuntimeFailure(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	create := `{"id":"swap","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"},{"item_id":"b","address":"D200","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	tags, mappings := applyMixedRuleTags(t, f, "swap")
	a, b := mixedLinkAt(t, f, "swap", "D100"), mixedLinkAt(t, f, "swap", "D200")
	beforeTag, err := tags.GetByID(t.Context(), *a.TagID)
	require.NoError(t, err)
	beforeMapping, err := mappings.GetByID(t.Context(), *a.MappingID)
	require.NoError(t, err)
	f.rules.SetShareRuntimeReconciler(sourcerule.ShareRuntimeReconcilerFunc(func(context.Context, sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
		return sourcerule.RuntimeReconcileOutcome{}, errors.New("injected reconcile failure")
	}))
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D200","data_type":"uint16","semantic_kind":"gauge"},{"item_id":"b","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, _ = f.do(t, http.MethodPut, "/workspace/source-rules/swap", update)
	require.NotEqual(t, http.StatusOK, code)
	require.Equal(t, a.PointID, mixedLinkAt(t, f, "swap", "D100").PointID)
	require.Equal(t, b.PointID, mixedLinkAt(t, f, "swap", "D200").PointID)
	for _, original := range []*schema.SourceRuleLink{a, b} {
		p, err := f.points.GetByID(t.Context(), original.PointID)
		require.NoError(t, err)
		require.Equal(t, original.Address, p.Address)
	}
	afterTag, err := tags.GetByID(t.Context(), *a.TagID)
	require.NoError(t, err)
	require.Equal(t, beforeTag.Labels, afterTag.Labels)
	afterMapping, err := mappings.GetByID(t.Context(), *a.MappingID)
	require.NoError(t, err)
	require.Equal(t, beforeMapping.RuleCandidateID, afterMapping.RuleCandidateID)
}

type concurrentTagEditOnMove struct {
	*tag.SQLRepository
	db *sql.DB
}

func (r *concurrentTagEditOnMove) RestoreWorkspaceSave(ctx context.Context, expected, previous *schema.Tag) error {
	if _, err := r.db.ExecContext(ctx, `UPDATE tags SET description='concurrent edit' WHERE id=?`, expected.ID); err != nil {
		return err
	}
	return r.SQLRepository.RestoreWorkspaceSave(ctx, expected, previous)
}

func TestWorkspaceMixedItemMoveDoesNotOverwriteConcurrentTagEdit(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	create := `{"id":"move","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"value","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	_, mappings := applyMixedRuleTags(t, f, "move")
	before := mixedLinkAt(t, f, "move", "D100")
	tags := tag.NewService(&concurrentTagEditOnMove{SQLRepository: tag.NewSQLRepository(f.db), db: f.db})
	f.rules.SetTagMappingServices(tags, mappings)
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"value","address":"D200","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, _ = f.do(t, http.MethodPut, "/workspace/source-rules/move", update)
	require.NotEqual(t, http.StatusOK, code)
	current, err := tags.GetByID(t.Context(), *before.TagID)
	require.NoError(t, err)
	require.Equal(t, "concurrent edit", current.Description)
}
