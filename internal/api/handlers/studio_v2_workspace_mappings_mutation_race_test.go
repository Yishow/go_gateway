package handlers

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWorkspaceMappingOtherRuleSaveAndDeleteProgressWhileRuleHeld(t *testing.T) {
	f := newWorkspaceConcurrentFixture(t, 4, 2, nil)
	_, err := f.rules.Create(t.Context(), sourcerule.CreateRuleRequest{ID: "other-rule", DeviceID: "diag-dev", StartAddress: "40101", Count: 1, DataType: schema.DataTypeInt16, NamingPrefix: "OTHER", Enabled: true})
	require.NoError(t, err)
	body := `{"rule_id":"other-rule","address":"40101","tag_key":"diag.other","display_name":"Other","target_type":"int16","scale":1,"offset":0,"enabled":true}`
	created := f.request(http.MethodPost, "", body)
	require.Equal(t, 201, created.Code, created.Body.String())
	id := decodeWorkspaceMappingBody(t, created)["data"].(map[string]any)["id"].(string)
	_, release := f.rules.AcquireRuleMutation(context.Background(), "diag-rule")
	defer release()
	for _, method := range []string{http.MethodPut, http.MethodDelete} {
		done := make(chan *httptest.ResponseRecorder, 1)
		requestBody := body
		if method == http.MethodDelete {
			requestBody = ""
		}
		go func() { done <- f.request(method, id, requestBody) }()
		select {
		case result := <-done:
			require.Equal(t, 200, result.Code, result.Body.String())
		case <-time.After(2 * time.Second):
			release()
			<-done
			t.Fatal("unrelated rule blocked behind held guard")
		}
	}
}

func TestWorkspaceMappingSaveCoordinatesSourceMutations(t *testing.T) {
	for _, operation := range []string{"sync", "update", "runtime-update", "candidate-reapply", "recompute", "delete", "runtime-delete"} {
		t.Run(operation, func(t *testing.T) {
			repo := &pausedWorkspaceLinkRepository{entered: make(chan struct{}), resume: make(chan struct{})}
			f := newWorkspaceConcurrentFixture(t, 4, 2, func(sqlRepo *sourcerule.SQLRepository) sourcerule.Repository {
				repo.SQLRepository = sqlRepo
				return repo
			})
			view, err := f.rules.GetCandidateView(t.Context(), "diag-rule")
			require.NoError(t, err)
			candidate := view.Tags.Candidates[0].(*sourcerule.TagCandidateView)
			multiplier := 2.0
			repo.enabled = true
			saveDone := make(chan *httptest.ResponseRecorder, 1)
			mutationDone := make(chan error, 1)
			go func() { saveDone <- f.request(http.MethodPut, f.ids[0], f.bodies[0]) }()
			select {
			case <-repo.entered:
			case <-time.After(3 * time.Second):
				close(repo.resume)
				t.Fatal("save did not pause")
			}
			go func() {
				var err error
				switch operation {
				case "sync":
					err = f.rules.SyncDerivedPointState(context.Background())
				case "update":
					_, err = f.rules.Update(context.Background(), "diag-rule", sourcerule.UpdateRuleRequest{ScaleMultiplier: &multiplier, ScaleMultiplierSet: true})
				case "runtime-update":
					_, _, err = f.rules.UpdateWithRuntimeReconcile(context.Background(), "diag-rule", sourcerule.UpdateRuleRequest{ScaleMultiplier: &multiplier, ScaleMultiplierSet: true})
				case "candidate-reapply":
					var result *sourcerule.ApplyTagCandidatesResponse
					result, err = f.rules.ApplyTagCandidates(context.Background(), "diag-rule", sourcerule.ApplyTagCandidatesRequest{RevisionID: view.RevisionID, CandidateIDs: []string{candidate.ID}})
					if err == nil && result.Results[0].Status == "failed" {
						err = errors.New(result.Results[0].Error)
					}
				case "recompute":
					_, err = f.rules.RecomputeCandidateViewAtRevision(context.Background(), "diag-rule", view.RevisionID, nil)
				case "delete":
					err = f.rules.Delete(context.Background(), "diag-rule")
				case "runtime-delete":
					_, err = f.rules.DeleteWithRuntimeReconcile(context.Background(), "diag-rule")
				}
				mutationDone <- err
			}()
			select {
			case err := <-mutationDone:
				close(repo.resume)
				t.Fatalf("source mutation escaped rule guard: %v", err)
			case <-time.After(30 * time.Millisecond):
			}
			close(repo.resume)
			select {
			case response := <-saveDone:
				require.Equal(t, 200, response.Code, response.Body.String())
			case <-time.After(3 * time.Second):
				t.Fatal("save deadlocked")
			}
			select {
			case err := <-mutationDone:
				require.NoError(t, err)
			case <-time.After(3 * time.Second):
				t.Fatal("mutation deadlocked")
			}
			if operation == "update" || operation == "runtime-update" {
				rejected := f.request(http.MethodPut, f.ids[0], f.bodies[0])
				require.Equal(t, 409, rejected.Code, rejected.Body.String())
				require.Contains(t, rejected.Body.String(), `"code":"workspace_mapping_conflict"`)
			}
			if operation == "delete" || operation == "runtime-delete" {
				rejected := f.request(http.MethodPut, f.ids[0], f.bodies[0])
				require.Equal(t, 404, rejected.Code, rejected.Body.String())
				require.Contains(t, rejected.Body.String(), `"code":"workspace_mapping_not_found"`)
			}
		})
	}
}

func TestWorkspaceMappingCreateFailedReplacementExplicitRetryKeepsOneOwnedDraft(t *testing.T) {
	f := newWorkspaceConcurrentFixture(t, 1, 1, nil)
	deleted := f.request(http.MethodDelete, f.ids[0], "")
	require.Equal(t, 200, deleted.Code, deleted.Body.String())
	_, err := f.db.ExecContext(t.Context(), `CREATE TRIGGER fail_create_link BEFORE INSERT ON source_rule_links BEGIN SELECT RAISE(ABORT,'private failure'); END`)
	require.NoError(t, err)
	failed := f.request(http.MethodPost, "", f.bodies[0])
	require.Equal(t, 500, failed.Code, failed.Body.String())
	records, err := f.mappings.List(t.Context(), mapping.ListFilter{})
	require.NoError(t, err)
	require.Len(t, records, 1)
	require.False(t, records[0].Enabled)
	require.Empty(t, records[0].LastAppliedSignature)
	links, err := f.rules.ListLinks(t.Context(), "diag-rule")
	require.NoError(t, err)
	require.Len(t, links, 1)
	require.Nil(t, links[0].MappingID)
	_, err = f.db.ExecContext(t.Context(), `DROP TRIGGER fail_create_link`)
	require.NoError(t, err)
	retry := f.request(http.MethodPost, "", f.bodies[0])
	require.Equal(t, 200, retry.Code, retry.Body.String())
	records, err = f.mappings.List(t.Context(), mapping.ListFilter{})
	require.NoError(t, err)
	require.Len(t, records, 1)
	require.NotEmpty(t, records[0].LastAppliedSignature)
	require.Equal(t, records[0].ProposedSignature, records[0].LastAppliedSignature)
	require.NotContains(t, retry.Body.String(), "private failure")
	tags, err := f.tags.List(t.Context(), tag.ListFilter{})
	require.NoError(t, err)
	require.Len(t, tags, 1)
}

func TestWorkspaceMappingLateDirectManualEditPreservesExactCASAndAcceptedHash(t *testing.T) {
	repo := &pausedWorkspaceLinkRepository{entered: make(chan struct{}), resume: make(chan struct{})}
	f := newWorkspaceConcurrentFixture(t, 1, 2, func(sqlRepo *sourcerule.SQLRepository) sourcerule.Repository {
		repo.SQLRepository = sqlRepo
		return repo
	})
	previous, err := f.mappings.GetByID(t.Context(), f.ids[0])
	require.NoError(t, err)
	repo.enabled = true
	done := make(chan *httptest.ResponseRecorder, 1)
	body := strings.ReplaceAll(f.bodies[0], `"scale":1`, `"scale":0.5`)
	go func() { done <- f.request(http.MethodPut, f.ids[0], body) }()
	select {
	case <-repo.entered:
	case <-time.After(3 * time.Second):
		close(repo.resume)
		t.Fatal("save did not pause")
	}
	manual := []schema.TransformStep{{Type: schema.TransformScale, Order: 0, Params: map[string]any{"scale": 3.0, "offset": 0.0}}}
	external, err := f.mappings.Update(t.Context(), f.ids[0], mapping.UpdateMappingRequest{TransformPipeline: manual})
	require.NoError(t, err)
	close(repo.resume)
	response := <-done
	require.Equal(t, 409, response.Code, response.Body.String())
	current, err := f.mappings.GetByID(t.Context(), f.ids[0])
	require.NoError(t, err)
	require.Equal(t, external.TransformPipeline, current.TransformPipeline)
	require.Equal(t, previous.LastAppliedSignature, current.LastAppliedSignature)
	require.NoError(t, f.rules.SyncDerivedPointState(t.Context()))
	current, err = f.mappings.GetByID(t.Context(), f.ids[0])
	require.NoError(t, err)
	require.Equal(t, schema.MappingStatusOutOfSync, current.Status)
	require.Equal(t, previous.LastAppliedSignature, current.LastAppliedSignature)
}

func TestWorkspaceMappingFailedDeletePreservesMappingAndAllowsExplicitRetry(t *testing.T) {
	f := newWorkspaceConcurrentFixture(t, 1, 2, nil)
	previous, err := f.mappings.GetByID(t.Context(), f.ids[0])
	require.NoError(t, err)
	previousLinks, err := f.rules.ListLinks(t.Context(), "diag-rule")
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `CREATE TRIGGER fail_delete_link BEFORE INSERT ON source_rule_links WHEN NEW.address='40002' BEGIN SELECT RAISE(ABORT,'private failure'); END`)
	require.NoError(t, err)
	failed := f.request(http.MethodDelete, f.ids[0], "")
	require.Equal(t, 500, failed.Code, failed.Body.String())
	current, err := f.mappings.GetByID(t.Context(), f.ids[0])
	require.NoError(t, err)
	require.Equal(t, previous, current)
	currentLinks, err := f.rules.ListLinks(t.Context(), "diag-rule")
	require.NoError(t, err)
	require.Equal(t, previousLinks, currentLinks)
	_, err = f.db.ExecContext(t.Context(), `DROP TRIGGER fail_delete_link`)
	require.NoError(t, err)
	retry := f.request(http.MethodDelete, f.ids[0], "")
	require.Equal(t, 200, retry.Code, retry.Body.String())
	_, err = f.mappings.GetByID(t.Context(), f.ids[0])
	require.ErrorIs(t, err, mapping.ErrMappingNotFound)
}

func TestWorkspaceMappingSQLTagRenamePersistsIdentityAndCleansOldTag(t *testing.T) {
	f := newWorkspaceConcurrentFixture(t, 1, 2, nil)
	previous, err := f.mappings.GetByID(t.Context(), f.ids[0])
	require.NoError(t, err)
	body := strings.ReplaceAll(f.bodies[0], "diag.block.0", "diag.renamed.0")
	response := f.request(http.MethodPut, f.ids[0], body)
	require.Equal(t, 200, response.Code, response.Body.String())
	renamed, err := f.tags.GetByKey(t.Context(), "diag.renamed.0")
	require.NoError(t, err)
	current, err := f.mappings.GetByID(t.Context(), f.ids[0])
	require.NoError(t, err)
	require.Equal(t, renamed.ID, current.TagID)
	require.Equal(t, previous.LastAppliedSignature, current.LastAppliedSignature)
	_, err = f.tags.GetByID(t.Context(), previous.TagID)
	require.ErrorIs(t, err, tag.ErrTagNotFound)
	links, err := f.rules.ListLinks(t.Context(), "diag-rule")
	require.NoError(t, err)
	require.Equal(t, current.TagID, *links[0].TagID)
	require.Equal(t, current.ID, *links[0].MappingID)
}

func TestWorkspaceMappingFailedSQLTagRenameRestoresOldOwnershipAndCleansNewOrphan(t *testing.T) {
	f := newWorkspaceConcurrentFixture(t, 1, 2, nil)
	previous, err := f.mappings.GetByID(t.Context(), f.ids[0])
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `CREATE TRIGGER fail_rename_link BEFORE INSERT ON source_rule_links WHEN NEW.address='40002' BEGIN SELECT RAISE(ABORT,'private failure'); END`)
	require.NoError(t, err)
	body := strings.ReplaceAll(strings.ReplaceAll(f.bodies[0], "diag.block.0", "diag.renamed.0"), `"scale":1`, `"scale":0.5`)
	failed := f.request(http.MethodPut, f.ids[0], body)
	require.Equal(t, 500, failed.Code, failed.Body.String())
	current, err := f.mappings.GetByID(t.Context(), f.ids[0])
	require.NoError(t, err)
	require.Equal(t, previous.TagID, current.TagID)
	require.Equal(t, previous.TransformPipeline, current.TransformPipeline)
	require.Equal(t, previous.LastAppliedSignature, current.LastAppliedSignature)
	require.Equal(t, previous.Enabled, current.Enabled)
	_, err = f.tags.GetByKey(t.Context(), "diag.renamed.0")
	require.ErrorIs(t, err, tag.ErrTagNotFound)
	_, err = f.db.ExecContext(t.Context(), `DROP TRIGGER fail_rename_link`)
	require.NoError(t, err)
	retry := f.request(http.MethodPut, f.ids[0], body)
	require.Equal(t, 200, retry.Code, retry.Body.String())
}

func TestWorkspaceMappingCannotRebindToAnotherRowsOwnedTag(t *testing.T) {
	f := newWorkspaceConcurrentFixture(t, 1, 2, nil)
	previous, err := f.mappings.GetByID(t.Context(), f.ids[0])
	require.NoError(t, err)
	second, err := f.mappings.GetByID(t.Context(), f.ids[1])
	require.NoError(t, err)
	body := strings.ReplaceAll(f.bodies[0], "diag.block.0", "diag.block.1")
	response := f.request(http.MethodPut, f.ids[0], body)
	require.Equal(t, 400, response.Code, response.Body.String())
	current, err := f.mappings.GetByID(t.Context(), f.ids[0])
	require.NoError(t, err)
	require.Equal(t, previous.TagID, current.TagID)
	require.Equal(t, previous.LastAppliedSignature, current.LastAppliedSignature)
	require.Equal(t, previous.Enabled, current.Enabled)
	foreign, err := f.mappings.GetByID(t.Context(), f.ids[1])
	require.NoError(t, err)
	require.Equal(t, second, foreign)
}

func TestWorkspaceMappingDeleteReportsCleanupFailureAsCompletedDeletion(t *testing.T) {
	f := newWorkspaceConcurrentFixture(t, 1, 1, nil)
	_, err := f.db.ExecContext(t.Context(), `CREATE TRIGGER fail_orphan_tag_cleanup BEFORE DELETE ON tags BEGIN SELECT RAISE(ABORT,'private cleanup failure'); END`)
	require.NoError(t, err)
	response := f.request(http.MethodDelete, f.ids[0], "")
	require.Equal(t, 200, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), `"cleanup_status":"failed"`)
	require.Contains(t, response.Body.String(), "Mapping deleted")
	require.NotContains(t, response.Body.String(), "private cleanup failure")
	_, err = f.mappings.GetByID(t.Context(), f.ids[0])
	require.ErrorIs(t, err, mapping.ErrMappingNotFound)
	links, err := f.rules.ListLinks(t.Context(), "diag-rule")
	require.NoError(t, err)
	require.Nil(t, links[0].MappingID)
}

func TestWorkspaceMappingFailedSaveTagRollbackCannotOverwriteDirectTagEdit(t *testing.T) {
	repo := &pausedWorkspaceLinkRepository{entered: make(chan struct{}), resume: make(chan struct{})}
	f := newWorkspaceConcurrentFixture(t, 1, 2, func(sqlRepo *sourcerule.SQLRepository) sourcerule.Repository {
		repo.SQLRepository = sqlRepo
		return repo
	})
	previous, err := f.mappings.GetByID(t.Context(), f.ids[0])
	require.NoError(t, err)
	repo.enabled = true
	repo.replacementErr = errors.New("private replacement failure")
	body := strings.ReplaceAll(strings.ReplaceAll(f.bodies[0], `"scale":1`, `"scale":0.5`), "BLOCK0", "Workspace edit")
	done := make(chan *httptest.ResponseRecorder, 1)
	go func() { done <- f.request(http.MethodPut, f.ids[0], body) }()
	select {
	case <-repo.entered:
	case <-time.After(3 * time.Second):
		close(repo.resume)
		t.Fatal("save did not pause")
	}
	directName := "Direct edit"
	_, err = f.tags.Update(t.Context(), previous.TagID, tag.UpdateTagRequest{DisplayName: &directName})
	require.NoError(t, err)
	close(repo.resume)
	response := <-done
	require.Equal(t, 500, response.Code, response.Body.String())
	require.NotContains(t, response.Body.String(), "private replacement")
	currentTag, err := f.tags.GetByID(t.Context(), previous.TagID)
	require.NoError(t, err)
	require.Equal(t, directName, currentTag.DisplayName)
	current, err := f.mappings.GetByID(t.Context(), previous.ID)
	require.NoError(t, err)
	require.Equal(t, previous.TransformPipeline, current.TransformPipeline)
	require.Equal(t, previous.LastAppliedSignature, current.LastAppliedSignature)
}

func TestWorkspaceMappingDeleteFailureRestoresLinksWithoutOverwritingDirectIdentity(t *testing.T) {
	for _, directEdit := range []bool{false, true} {
		t.Run(map[bool]string{false: "restore", true: "direct-edit"}[directEdit], func(t *testing.T) {
			repo := &pausedWorkspaceLinkRepository{entered: make(chan struct{}), resume: make(chan struct{})}
			f := newWorkspaceConcurrentFixture(t, 1, 2, func(sqlRepo *sourcerule.SQLRepository) sourcerule.Repository {
				repo.SQLRepository = sqlRepo
				return repo
			})
			previousLinks, err := f.rules.ListLinks(t.Context(), "diag-rule")
			require.NoError(t, err)
			_, err = f.db.ExecContext(t.Context(), `CREATE TRIGGER fail_mapping_delete BEFORE DELETE ON mappings BEGIN SELECT RAISE(ABORT,'private mapping delete failure'); END`)
			require.NoError(t, err)
			repo.enabled = true
			done := make(chan *httptest.ResponseRecorder, 1)
			go func() { done <- f.request(http.MethodDelete, f.ids[0], "") }()
			select {
			case <-repo.entered:
			case <-time.After(3 * time.Second):
				close(repo.resume)
				t.Fatal("delete did not pause")
			}
			var manualTag *schema.Tag
			if directEdit {
				manualTag, err = f.tags.Create(t.Context(), tag.CreateTagRequest{Key: "manual.direct", DataType: schema.DataTypeInt16})
				require.NoError(t, err)
				_, err = f.mappings.Update(t.Context(), f.ids[0], mapping.UpdateMappingRequest{TagID: &manualTag.ID})
				require.NoError(t, err)
			}
			close(repo.resume)
			response := <-done
			require.NotContains(t, response.Body.String(), "private mapping delete failure")
			links, err := f.rules.ListLinks(t.Context(), "diag-rule")
			require.NoError(t, err)
			if directEdit {
				require.Equal(t, 409, response.Code, response.Body.String())
				require.Nil(t, links[0].MappingID)
				require.Nil(t, links[0].TagID)
				current, err := f.mappings.GetByID(t.Context(), f.ids[0])
				require.NoError(t, err)
				require.Equal(t, manualTag.ID, current.TagID)
			} else {
				require.Equal(t, 500, response.Code, response.Body.String())
				require.Equal(t, previousLinks, links)
				_, err = f.db.ExecContext(t.Context(), `DROP TRIGGER fail_mapping_delete`)
				require.NoError(t, err)
				retry := f.request(http.MethodDelete, f.ids[0], "")
				require.Equal(t, 200, retry.Code, retry.Body.String())
			}
		})
	}
}

func TestWorkspaceMappingRuntimeReconcileReusesSynchronousRuleGuard(t *testing.T) {
	f := newWorkspaceConcurrentFixture(t, 1, 2, nil)
	f.rules.SetShareRuntimeReconciler(sourcerule.ShareRuntimeReconcilerFunc(func(ctx context.Context, req sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
		err := f.rules.SyncRuleDerivedState(ctx, req.Scope.RuleID)
		return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusAligned, Scope: req.Scope}, err
	}))
	done := make(chan error, 1)
	go func() {
		_, _, err := f.rules.UpdateWithRuntimeReconcile(context.Background(), "diag-rule", sourcerule.UpdateRuleRequest{})
		done <- err
	}()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("runtime reconcile nested source sync deadlocked")
	}
}

func TestWorkspaceMappingAcceptedRenameReportsCleanupFailureAsSuccessfulSave(t *testing.T) {
	for _, method := range []string{http.MethodPut, http.MethodPost} {
		t.Run(method, func(t *testing.T) {
			f := newWorkspaceConcurrentFixture(t, 1, 2, nil)
			previous, err := f.mappings.GetByID(t.Context(), f.ids[0])
			require.NoError(t, err)
			_, err = f.db.ExecContext(t.Context(), `CREATE TRIGGER fail_accepted_orphan_cleanup BEFORE DELETE ON tags BEGIN SELECT RAISE(ABORT,'private orphan cleanup failure'); END`)
			require.NoError(t, err)
			body := strings.ReplaceAll(strings.ReplaceAll(f.bodies[0], "diag.block.0", "diag.accepted.rename"), `"scale":1`, `"scale":0.5`)
			response := f.request(method, f.ids[0], body)
			require.Equal(t, 200, response.Code, response.Body.String())
			require.Contains(t, response.Body.String(), `"cleanup_status":"failed"`)
			require.NotContains(t, response.Body.String(), "private orphan cleanup failure")
			current, err := f.mappings.GetByID(t.Context(), previous.ID)
			require.NoError(t, err)
			require.NotEqual(t, previous.TagID, current.TagID)
			require.NotEqual(t, previous.TransformPipeline, current.TransformPipeline)
			require.NotEqual(t, previous.LastAppliedSignature, current.LastAppliedSignature)
			require.Equal(t, schema.MappingStatusActive, current.Status)
			require.True(t, current.Enabled)
			renamed, err := f.tags.GetByKey(t.Context(), "diag.accepted.rename")
			require.NoError(t, err)
			require.Equal(t, renamed.ID, current.TagID)
			_, err = f.tags.GetByID(t.Context(), previous.TagID)
			require.NoError(t, err, "cleanup failure must retain the old orphan, not claim it was removed")
			links, err := f.rules.ListLinks(t.Context(), "diag-rule")
			require.NoError(t, err)
			require.Equal(t, current.ID, *links[0].MappingID)
			require.Equal(t, current.TagID, *links[0].TagID)
		})
	}
}
