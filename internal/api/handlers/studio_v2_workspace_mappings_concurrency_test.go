package handlers

import (
	"context"
	"database/sql"
	"fmt"
	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type workspaceConcurrentFixture struct {
	ctx         context.Context
	db          *sql.DB
	h           *StudioV2WorkspaceMappingsHandler
	rules       *sourcerule.Service
	mappings    *mapping.Service
	tags        *tag.Service
	ids, bodies []string
}

func newWorkspaceConcurrentFixture(t *testing.T, pool, count int, wrap func(*sourcerule.SQLRepository) sourcerule.Repository) *workspaceConcurrentFixture {
	t.Helper()
	gin.SetMode(gin.TestMode)
	ctx := t.Context()
	db, openErr := sql.Open("sqlite", filepath.Join(t.TempDir(), "configuration.db")+"?_pragma=busy_timeout(15000)&_pragma=journal_mode(WAL)")
	require.NoError(t, openErr)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, datalinkbase.NewMigrator().Migrate(db))
	db.SetMaxOpenConns(pool)
	dr := device.NewSQLRepository(db)
	success := true
	require.NoError(t, dr.Create(ctx, &schema.Device{ID: "diag-dev", Name: "Diagnostic", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, LastTestSuccess: &success, ConnectionConfig: `{"host":"127.0.0.1","port":15020,"slave_id":1,"timeout":5}`, CreatedAt: time.Now(), UpdatedAt: time.Now()}))
	ds := device.NewService(dr, nil)
	ps := point.NewService(point.NewSQLRepository(db), nil)
	ts := tag.NewService(tag.NewSQLRepository(db))
	ms := mapping.NewServiceWithTagResolver(mapping.NewSQLRepository(db), ts.GetByID)
	repo := sourcerule.Repository(sourcerule.NewSQLRepository(db))
	if wrap != nil {
		repo = wrap(repo.(*sourcerule.SQLRepository))
	}
	rs := sourcerule.NewService(repo, ds, ps, nil)
	ws := workspace.NewService(workspace.NewSQLRepository(db))
	_, err := ws.AttachDevice(ctx, "diag-dev")
	require.NoError(t, err)
	_, err = rs.Create(ctx, sourcerule.CreateRuleRequest{ID: "diag-rule", DeviceID: "diag-dev", StartAddress: "40001", Count: count, DataType: schema.DataTypeInt16, NamingPrefix: "BLOCK", Enabled: true})
	require.NoError(t, err)
	f := &workspaceConcurrentFixture{ctx: ctx, db: db, h: NewStudioV2WorkspaceMappingsHandler(ws, ds, rs, ps, ts, ms), rules: rs, mappings: ms, tags: ts, ids: make([]string, count), bodies: make([]string, count)}
	for i := range count {
		f.bodies[i] = fmt.Sprintf(`{"rule_id":"diag-rule","address":"%d","tag_key":"diag.block.%d","display_name":"BLOCK%d","unit":"","target_type":"int16","scale":1,"offset":0,"enabled":true}`, 40001+i, i, i)
		response := f.request(http.MethodPost, "", f.bodies[i])
		require.Equal(t, http.StatusCreated, response.Code, response.Body.String())
		f.ids[i] = decodeWorkspaceMappingBody(t, response)["data"].(map[string]any)["id"].(string)
	}
	return f
}

func (f *workspaceConcurrentFixture) request(method, id, body string) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request = httptest.NewRequestWithContext(f.ctx, method, "/workspace/mappings/"+id, strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("X-Request-ID", "mapping-regression-request")
	c.Params = gin.Params{{Key: "id", Value: id}}
	switch method {
	case http.MethodPost:
		f.h.Create(c)
	case http.MethodPut:
		f.h.Update(c)
	case http.MethodDelete:
		f.h.Delete(c)
	case http.MethodGet:
		f.h.List(c)
	}
	return response
}

func TestWorkspaceMappingEightConcurrentSQLiteSaves(t *testing.T) {
	for _, pool := range []int{1, 4} {
		t.Run(fmt.Sprintf("pool=%d", pool), func(t *testing.T) {
			f := newWorkspaceConcurrentFixture(t, pool, 8, nil)
			initialRule, err := f.rules.GetByID(t.Context(), "diag-rule")
			require.NoError(t, err)
			for batch := range 4 {
				results := make([]*httptest.ResponseRecorder, 8)
				start := make(chan struct{})
				var wg sync.WaitGroup
				for i := range 8 {
					wg.Add(1)
					go func(i int) { defer wg.Done(); <-start; results[i] = f.request(http.MethodPut, f.ids[i], f.bodies[i]) }(i)
				}
				close(start)
				wg.Wait()
				codes := make([]int, 8)
				for i, result := range results {
					codes[i] = result.Code
					require.Equal(t, http.StatusOK, result.Code, result.Body.String())
				}
				t.Logf("batch=%d pool=%d statuses=%v", batch, pool, codes)
				currentRule, err := f.rules.GetByID(t.Context(), "diag-rule")
				require.NoError(t, err)
				require.Equal(t, initialRule.RevisionID, currentRule.RevisionID)
				links, err := f.rules.ListLinks(t.Context(), "diag-rule")
				require.NoError(t, err)
				require.Len(t, links, 8)
				seen := map[string]bool{}
				for i, link := range links {
					require.False(t, seen[link.PointID])
					seen[link.PointID] = true
					require.NotNil(t, link.TagID)
					require.NotNil(t, link.MappingID)
					record, err := f.mappings.GetByID(t.Context(), *link.MappingID)
					require.NoError(t, err)
					require.Equal(t, f.ids[i], record.ID)
					require.Equal(t, link.PointID, record.PointID)
					require.Equal(t, *link.TagID, record.TagID)
					require.Equal(t, "[]", record.TransformPipeline)
					pointRecord, err := f.h.pointSvc.GetByID(t.Context(), link.PointID)
					require.NoError(t, err)
					require.Equal(t, schema.DataTypeInt16, pointRecord.DataType)
					require.Equal(t, link.Address, pointRecord.Address)
					require.Equal(t, record.ProposedSignature, record.LastAppliedSignature)
					require.NotEmpty(t, record.LastAppliedSignature)
					require.True(t, record.Enabled)
					require.Equal(t, schema.MappingStatusActive, record.Status)
					tagRecord, err := f.tags.GetByID(t.Context(), record.TagID)
					require.NoError(t, err)
					require.Equal(t, schema.DataTypeInt16, tagRecord.DataType)
				}
			}
		})
	}
}

type pausedWorkspaceLinkRepository struct {
	*sourcerule.SQLRepository
	enabled         bool
	replacementErr  error
	once            sync.Once
	entered, resume chan struct{}
}

func (r *pausedWorkspaceLinkRepository) ReplaceLinks(ctx context.Context, id string, links []*schema.SourceRuleLink) error {
	if r.enabled {
		r.once.Do(func() { close(r.entered); <-r.resume })
	}
	if r.enabled && r.replacementErr != nil {
		return r.replacementErr
	}
	return r.SQLRepository.ReplaceLinks(ctx, id, links)
}

func TestWorkspaceMappingReplacementCannotHideAnotherRow(t *testing.T) {
	repo := &pausedWorkspaceLinkRepository{entered: make(chan struct{}), resume: make(chan struct{})}
	f := newWorkspaceConcurrentFixture(t, 1, 2, func(sqlRepo *sourcerule.SQLRepository) sourcerule.Repository {
		repo.SQLRepository = sqlRepo
		return repo
	})
	repo.enabled = true
	first := make(chan *httptest.ResponseRecorder, 1)
	second := make(chan *httptest.ResponseRecorder, 1)
	go func() { first <- f.request(http.MethodPut, f.ids[0], f.bodies[0]) }()
	select {
	case <-repo.entered:
	case <-time.After(3 * time.Second):
		close(repo.resume)
		t.Fatal("did not enter replacement")
	}
	old, err := repo.ListLinks(t.Context(), "diag-rule")
	require.NoError(t, err)
	require.Len(t, old, 2)
	go func() { second <- f.request(http.MethodPut, f.ids[1], f.bodies[1]) }()
	select {
	case result := <-second:
		close(repo.resume)
		t.Fatalf("same-rule save escaped guard: %d %s", result.Code, result.Body.String())
	case <-time.After(30 * time.Millisecond):
	}
	close(repo.resume)
	for _, done := range []chan *httptest.ResponseRecorder{first, second} {
		select {
		case result := <-done:
			require.Equal(t, 200, result.Code, result.Body.String())
		case <-time.After(3 * time.Second):
			t.Fatal("replacement deadlocked")
		}
	}
}

func TestWorkspaceMappingFailedReplacementRestoresAcceptedDraftAndAllowsExplicitRetry(t *testing.T) {
	f := newWorkspaceConcurrentFixture(t, 1, 2, nil)
	prior, err := f.mappings.GetByID(t.Context(), f.ids[0])
	require.NoError(t, err)
	priorTag, err := f.tags.GetByID(t.Context(), prior.TagID)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `CREATE TRIGGER fail_mapping_link BEFORE INSERT ON source_rule_links WHEN NEW.address='40002' BEGIN SELECT RAISE(ABORT,'private SQL failure'); END`)
	require.NoError(t, err)
	body := strings.ReplaceAll(strings.ReplaceAll(f.bodies[0], `"scale":1`, `"scale":0.5`), `"display_name":"BLOCK0"`, `"display_name":"Edited"`)
	failed := f.request(http.MethodPut, f.ids[0], body)
	require.Equal(t, 500, failed.Code, failed.Body.String())
	require.Contains(t, failed.Body.String(), `"code":"workspace_mapping_save_failed"`)
	require.Contains(t, failed.Body.String(), `"request_id":"mapping-regression-request"`)
	require.NotContains(t, failed.Body.String(), "private SQL")
	actual, err := f.mappings.GetByID(t.Context(), prior.ID)
	require.NoError(t, err)
	require.Equal(t, prior, actual)
	actualTag, err := f.tags.GetByID(t.Context(), prior.TagID)
	require.NoError(t, err)
	require.Equal(t, priorTag.DisplayName, actualTag.DisplayName)
	require.Equal(t, priorTag.DataType, actualTag.DataType)
	links, err := f.rules.ListLinks(t.Context(), "diag-rule")
	require.NoError(t, err)
	require.Len(t, links, 2)
	_, err = f.db.ExecContext(t.Context(), `DROP TRIGGER fail_mapping_link`)
	require.NoError(t, err)
	retry := f.request(http.MethodPut, f.ids[0], body)
	require.Equal(t, 200, retry.Code, retry.Body.String())
	actual, err = f.mappings.GetByID(t.Context(), prior.ID)
	require.NoError(t, err)
	require.NotEqual(t, prior.TransformPipeline, actual.TransformPipeline)
	require.NotEqual(t, prior.LastAppliedSignature, actual.LastAppliedSignature)
	require.NoError(t, f.rules.SyncDerivedPointState(t.Context()))
	actual, err = f.mappings.GetByID(t.Context(), prior.ID)
	require.NoError(t, err)
	require.Equal(t, schema.MappingStatusActive, actual.Status)
}
