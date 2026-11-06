package handlers

import (
	"context"
	"encoding/json"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecordingHomeSourcesDisabledSkippedAndLegacyBinding(t *testing.T) {
	for _, mode := range []string{"disabled rule", "skipped", "disabled point", "disabled mapping", "legacy binding"} {
		t.Run(mode, func(t *testing.T) {
			f := newHomeRealFixture(t)
			rr, pr, mr, tr := f.handler.rules, f.handler.points, f.handler.mappings, f.handler.tags
			f.handler.rules = &homeRuleStub{RecordingHomeRuleReader: rr, list: func(c context.Context, ids []string) ([]*schema.SourceRule, error) {
				v, e := rr.ListByDeviceIDs(c, ids)
				if mode == "disabled rule" {
					v[0].Enabled = false
				}
				if mode == "skipped" {
					v[0].SkippedAddresses = `["40001"]`
				}
				return v, e
			}}
			f.handler.points = homePointFunc(func(c context.Context, filter point.ListFilter) ([]*schema.Point, error) {
				v, e := pr.List(c, filter)
				if mode == "disabled point" {
					v[0].Enabled = false
				}
				return v, e
			})
			f.handler.mappings = &homeMappingStub{list: func(c context.Context, filter mapping.ListFilter) ([]*schema.Mapping, error) {
				v, e := mr.List(c, filter)
				if mode == "disabled mapping" {
					v[0].Enabled = false
				}
				return v, e
			}, get: func(c context.Context, id string) (*schema.Mapping, error) {
				v, e := mr.GetByID(c, id)
				if mode == "disabled mapping" {
					v.Enabled = false
				}
				return v, e
			}}
			f.handler.tags = homeTagFunc(func(c context.Context, id string) (*schema.Tag, error) {
				v, e := tr.GetByID(c, id)
				if mode == "legacy binding" {
					v.Labels = ""
				}
				return v, e
			})
			out := homeSourcesRequest(f.handler, "?workspace_id="+f.id)
			require.Equal(t, 200, out.Code)
			var body RecordingHomeSourcesEnvelope
			require.NoError(t, json.Unmarshal(out.Body.Bytes(), &body))
			require.Equal(t, "available", body.Data.State)
			count := 0
			if mode == "legacy binding" {
				count = 1
			}
			require.Len(t, body.Data.Candidates, count)
			require.Empty(t, body.Data.Excluded)
		})
	}
}
func TestRecordingHomeSourcesChangedOnSecondPass(t *testing.T) {
	f := newHomeRealFixture(t)
	reader := f.handler.tags
	calls := 0
	f.handler.tags = homeTagFunc(func(c context.Context, id string) (*schema.Tag, error) {
		calls++
		v, e := reader.GetByID(c, id)
		if calls == 2 {
			v.DisplayName = "Changed"
		}
		return v, e
	})
	requireHomeUnknown(t, f)
	require.Equal(t, 2, calls)
}

type homeWorkspaceSaveSpy struct {
	workspace.Repository
	saves int
}

func (s *homeWorkspaceSaveSpy) Save(c context.Context, r *workspace.Record) error {
	s.saves++
	return s.Repository.Save(c, r)
}

type homeRuleMutationSpy struct {
	RecordingHomeRuleReader
	writes int
}

func (s *homeRuleMutationSpy) ReplaceLinks(context.Context, string, []*schema.SourceRuleLink) error {
	s.writes++
	return nil
}
func TestRecordingHomeSourcesAbsentWorkspaceDoesNotCreate(t *testing.T) {
	f := newHomeRealFixture(t)
	repo := &homeWorkspaceSaveSpy{Repository: workspace.NewMemoryRepository()}
	f.handler.workspace = workspace.NewService(repo)
	out := homeSourcesRequest(f.handler, "?workspace_id=W1")
	require.Equal(t, 404, out.Code)
	require.Zero(t, repo.saves)
	_, err := repo.Get(t.Context())
	require.ErrorIs(t, err, workspace.ErrNotFound)
}
func TestRecordingHomeSourcesRepeatedReadsNeverWriteOrRepair(t *testing.T) {
	f := newHomeRealFixture(t)
	before := f.snapshot(t)
	record, err := f.old.workspaceSvc.ReadSnapshot(t.Context())
	require.NoError(t, err)
	repo := &homeWorkspaceSaveSpy{Repository: workspace.NewMemoryRepository()}
	require.NoError(t, repo.Repository.Save(t.Context(), record))
	f.handler.workspace = workspace.NewService(repo)
	rules := &homeRuleMutationSpy{RecordingHomeRuleReader: f.handler.rules}
	f.handler.rules = rules
	for i := 0; i < 3; i++ {
		require.Equal(t, 200, homeSourcesRequest(f.handler, "?workspace_id="+f.id).Code)
	}
	require.Zero(t, repo.saves)
	require.Zero(t, rules.writes)
	require.Equal(t, before, f.snapshot(t))
}

func TestRecordingHomeSourcesUnconfiguredLinkChangeIsUnknown(t *testing.T) {
	f := newHomeRealFixture(t)
	rules := f.handler.rules
	calls := 0
	f.handler.rules = &homeRuleStub{RecordingHomeRuleReader: rules, links: func(c context.Context, id string) ([]*schema.SourceRuleLink, error) {
		calls++
		links, err := rules.ListLinks(c, id)
		links[0].MappingID = nil
		links[0].TagID = nil
		if calls == 2 {
			links[0].ID = "replacement-link"
		}
		return links, err
	}}
	f.handler.mappings = &homeMappingStub{list: func(context.Context, mapping.ListFilter) ([]*schema.Mapping, error) { return nil, nil }}
	requireHomeUnknown(t, f)
	require.Equal(t, 2, calls)
}
