package handlers

import (
	"context"
	"database/sql"
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

type editTagAfterSnapshot struct {
	*tag.SQLRepository
	db    *sql.DB
	id    string
	fired bool
}

func (r *editTagAfterSnapshot) GetByID(ctx context.Context, id string) (*schema.Tag, error) {
	before, err := r.SQLRepository.GetByID(ctx, id)
	if err == nil && id == r.id && !r.fired {
		r.fired = true
		_, err = r.db.ExecContext(ctx, `UPDATE tags SET description='concurrent edit' WHERE id=?`, id)
	}
	return before, err
}

func TestWorkspaceMixedItemReplacementRollbackKeepsConcurrentTagEdit(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	enableMixedFixtureForeignKeys(t, f)
	create := `{"id":"replace","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	_, mappings := applyMixedRuleTags(t, f, "replace")
	old := mixedLinkAt(t, f, "replace", "D100")
	repo := &editTagAfterSnapshot{SQLRepository: tag.NewSQLRepository(f.db), db: f.db, id: *old.TagID}
	tags := tag.NewService(repo)
	f.rules.SetTagMappingServices(tags, mappings)
	reconciled := false
	f.rules.SetShareRuntimeReconciler(sourcerule.ShareRuntimeReconcilerFunc(func(context.Context, sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
		reconciled = true
		return sourcerule.RuntimeReconcileOutcome{}, errors.New("injected reconcile failure")
	}))
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"c","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, failure := f.do(t, http.MethodPut, "/workspace/source-rules/replace", update)
	require.NotEqual(t, http.StatusOK, code)
	require.True(t, repo.fired)
	require.True(t, reconciled)
	require.NotContains(t, fmt.Sprint(failure), "dirty_unknown")
	current, err := tags.GetByID(t.Context(), *old.TagID)
	require.NoError(t, err)
	require.Equal(t, "concurrent edit", current.Description)
	require.Equal(t, old.PointID, mixedLinkAt(t, f, "replace", "D100").PointID)
}

type editMappingAfterSnapshot struct {
	*mapping.SQLRepository
	db    *sql.DB
	id    string
	fired bool
}

func (r *editMappingAfterSnapshot) GetByID(ctx context.Context, id string) (*schema.Mapping, error) {
	before, err := r.SQLRepository.GetByID(ctx, id)
	if err == nil && id == r.id && !r.fired {
		r.fired = true
		_, err = r.db.ExecContext(ctx, `UPDATE mappings SET blocking_reason='concurrent edit' WHERE id=?`, id)
	}
	return before, err
}

func TestWorkspaceMixedItemReplacementDoesNotDeleteConcurrentMappingEdit(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	enableMixedFixtureForeignKeys(t, f)
	create := `{"id":"replace","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	tags, _ := applyMixedRuleTags(t, f, "replace")
	old := mixedLinkAt(t, f, "replace", "D100")
	repo := &editMappingAfterSnapshot{SQLRepository: mapping.NewSQLRepository(f.db), db: f.db, id: *old.MappingID}
	mappings := mapping.NewServiceWithTagResolver(repo, tags.GetByID)
	f.rules.SetTagMappingServices(tags, mappings)
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"c","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, failure := f.do(t, http.MethodPut, "/workspace/source-rules/replace", update)
	require.NotEqual(t, http.StatusOK, code)
	require.True(t, repo.fired)
	require.Contains(t, fmt.Sprint(failure), "dirty_unknown")
	current, err := mappings.GetByID(t.Context(), *old.MappingID)
	require.NoError(t, err)
	require.Equal(t, "concurrent edit", current.BlockingReason)
	require.Equal(t, old.PointID, mixedLinkAt(t, f, "replace", "D100").PointID)
}

func TestWorkspaceMixedItemReplacementDoesNotDeleteNewConcurrentMapping(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	enableMixedFixtureForeignKeys(t, f)
	create := `{"id":"replace","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	_, mappings := applyMixedRuleTags(t, f, "replace")
	old := mixedLinkAt(t, f, "replace", "D100")
	_, err := f.db.ExecContext(t.Context(), `INSERT INTO tags(id,key,key_lower,display_name,data_type,status,labels,created_at,updated_at) VALUES('foreign-tag','foreign','foreign','foreign','uint16','active','{}',datetime('now'),datetime('now'))`)
	require.NoError(t, err)
	query := fmt.Sprintf(`CREATE TRIGGER add_concurrent_mapping AFTER DELETE ON mappings WHEN OLD.id='%s' BEGIN INSERT INTO mappings(id,point_id,tag_id,transform_pipeline,status,enabled,created_at,updated_at) VALUES('foreign-mapping',OLD.point_id,'foreign-tag','[]','active',1,datetime('now'),datetime('now')); END`, *old.MappingID)
	_, err = f.db.ExecContext(t.Context(), query)
	require.NoError(t, err)
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"c","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, failure := f.do(t, http.MethodPut, "/workspace/source-rules/replace", update)
	require.NotEqual(t, http.StatusOK, code)
	require.Contains(t, fmt.Sprint(failure), "dirty_unknown")
	foreign, err := mappings.GetByID(t.Context(), "foreign-mapping")
	require.NoError(t, err)
	require.Equal(t, old.PointID, foreign.PointID)
	_, err = f.points.GetByID(t.Context(), old.PointID)
	require.NoError(t, err)
}

func TestWorkspaceMixedItemReplacementDoesNotDeleteConcurrentPointEdit(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	enableMixedFixtureForeignKeys(t, f)
	create := `{"id":"replace","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	_, mappings := applyMixedRuleTags(t, f, "replace")
	old := mixedLinkAt(t, f, "replace", "D100")
	query := fmt.Sprintf(`CREATE TRIGGER edit_concurrent_point AFTER DELETE ON mappings WHEN OLD.id='%s' BEGIN UPDATE points SET revision_id='concurrent-revision' WHERE id=OLD.point_id; END`, *old.MappingID)
	_, err := f.db.ExecContext(t.Context(), query)
	require.NoError(t, err)
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"c","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, failure := f.do(t, http.MethodPut, "/workspace/source-rules/replace", update)
	require.NotEqual(t, http.StatusOK, code)
	require.Contains(t, fmt.Sprint(failure), "dirty_unknown")
	current, err := f.points.GetByID(t.Context(), old.PointID)
	require.NoError(t, err)
	require.Equal(t, "concurrent-revision", current.RevisionID)
	restored, err := mappings.GetByID(t.Context(), *old.MappingID)
	require.NoError(t, err)
	require.Equal(t, old.PointID, restored.PointID)
}

func TestWorkspaceMixedItemReplacementRollbackKeepsNewConcurrentMapping(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	enableMixedFixtureForeignKeys(t, f)
	create := `{"id":"replace","device_id":"mc-dev","naming_prefix":"M","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"a","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, response := f.do(t, http.MethodPost, "/workspace/source-rules", create)
	require.Equal(t, http.StatusCreated, code, response)
	tags := tag.NewService(tag.NewSQLRepository(f.db))
	mappings := mapping.NewServiceWithTagResolver(mapping.NewSQLRepository(f.db), tags.GetByID)
	f.rules.SetTagMappingServices(tags, mappings)
	_, err := f.db.ExecContext(t.Context(), `INSERT INTO tags(id,key,key_lower,display_name,data_type,status,labels,created_at,updated_at) VALUES('foreign-tag','foreign','foreign','foreign','uint16','active','{}',datetime('now'),datetime('now'))`)
	require.NoError(t, err)
	newPointID := ""
	f.rules.SetShareRuntimeReconciler(sourcerule.ShareRuntimeReconcilerFunc(func(ctx context.Context, _ sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
		links, linkErr := f.rules.ListLinks(ctx, "replace")
		if linkErr != nil {
			return sourcerule.RuntimeReconcileOutcome{}, fmt.Errorf("read replacement link: %w", linkErr)
		}
		if len(links) != 1 {
			return sourcerule.RuntimeReconcileOutcome{}, fmt.Errorf("expected one replacement link, got %d", len(links))
		}
		newPointID = links[0].PointID
		_, insertErr := f.db.ExecContext(ctx, `INSERT INTO mappings(id,point_id,tag_id,transform_pipeline,status,enabled,created_at,updated_at) VALUES('foreign-mapping',?,'foreign-tag','[]','active',1,datetime('now'),datetime('now'))`, newPointID)
		if insertErr != nil {
			return sourcerule.RuntimeReconcileOutcome{}, insertErr
		}
		return sourcerule.RuntimeReconcileOutcome{}, errors.New("injected runtime failure")
	}))
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"c","address":"D100","data_type":"uint16","semantic_kind":"gauge"}]}`
	code, failure := f.do(t, http.MethodPut, "/workspace/source-rules/replace", update)
	require.NotEqual(t, http.StatusOK, code)
	require.NotEmpty(t, newPointID)
	require.Contains(t, fmt.Sprint(failure), "dirty_unknown")
	foreign, err := mappings.GetByID(t.Context(), "foreign-mapping")
	require.NoError(t, err, "rollback must not cascade-delete another operation's mapping")
	require.Equal(t, newPointID, foreign.PointID)
	_, err = f.points.GetByID(t.Context(), newPointID)
	require.NoError(t, err)
}
