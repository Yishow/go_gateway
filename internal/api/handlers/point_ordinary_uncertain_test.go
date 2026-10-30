package handlers

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/point"
	runtimepkg "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sqltx"
	"go-gateway/internal/datalink/storage"
	"net/http/httptest"
	"testing"
)

type uncertainOrdinaryRepository struct {
	*point.MemoryRepository
	live    *runtimepkg.Service
	foreign bool
}

func (r *uncertainOrdinaryRepository) UpdateOrdinaryRevision(ctx context.Context, p *schema.Point, expected string, _ bool) error {
	if err := r.MemoryRepository.UpdateRevision(ctx, p, expected); err != nil {
		return err
	}
	if r.foreign {
		newer := *p
		newer.RevisionID = "foreign"
		r.live.UpsertPoint(&newer)
	}
	return sqltx.ErrCommitOutcomeUnknown
}
func TestOrdinaryUnknownCommitStopsOnlyMatchingCapturedLease(t *testing.T) {
	for _, foreign := range []bool{false, true} {
		t.Run(map[bool]string{false: "owned", true: "foreign"}[foreign], func(t *testing.T) {
			db := openStudioV2CascadeTestDB(t)
			defer db.Close()
			live, err := runtimepkg.NewService(runtimepkg.Config{Writer: storage.NewSQLiteWriter(db)})
			require.NoError(t, err)
			repo := &uncertainOrdinaryRepository{MemoryRepository: point.NewMemoryRepository(), live: live, foreign: foreign}
			svc := point.NewService(repo, nil)
			p, err := svc.Create(t.Context(), point.CreatePointRequest{DeviceID: "d", Name: "old", Address: "HR0", DataType: schema.DataTypeUint32})
			require.NoError(t, err)
			live.UpsertPoint(p)
			router := gin.New()
			router.PUT("/points/:id", NewPointHandler(svc, live).Update)
			req := httptest.NewRequest("PUT", "/points/"+p.ID, bytes.NewBufferString(`{"name":"uncertain"}`))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			require.Equal(t, 503, w.Code)
			require.Zero(t, live.CapturePointLease(p.ID, p.RevisionID))
			if foreign {
				require.NotZero(t, live.CapturePointLease(p.ID, "foreign"))
			}
		})
	}
}

type absentOrdinaryUnknownRepo struct{ *point.MemoryRepository }

func (r *absentOrdinaryUnknownRepo) UpdateOrdinaryRevision(ctx context.Context, p *schema.Point, expected string, _ bool) error {
	if err := r.MemoryRepository.UpdateRevision(ctx, p, expected); err != nil {
		return err
	}
	if p.Name == "uncertain" {
		return sqltx.ErrCommitOutcomeUnknown
	}
	return nil
}
func TestOrdinaryAbsentUnknownCommitFencesOlderCallback(t *testing.T) {
	db := openStudioV2CascadeTestDB(t)
	defer db.Close()
	live, err := runtimepkg.NewService(runtimepkg.Config{Writer: storage.NewSQLiteWriter(db)})
	require.NoError(t, err)
	svc := point.NewService(&absentOrdinaryUnknownRepo{point.NewMemoryRepository()}, nil)
	p, err := svc.Create(t.Context(), point.CreatePointRequest{DeviceID: "d", Name: "old", Address: "HR0", DataType: schema.DataTypeUint32})
	require.NoError(t, err)
	delayed := &lateOrdinaryRuntime{Service: live, entered: make(chan struct{}), release: make(chan struct{})}
	router := gin.New()
	router.PUT("/points/:id", NewPointHandler(svc, delayed).Update)
	send := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest("PUT", "/points/"+p.ID, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}
	done := make(chan *httptest.ResponseRecorder, 1)
	go func() { done <- send(`{"name":"late"}`) }()
	<-delayed.entered
	late, err := svc.GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	unknown := send(`{"name":"uncertain"}`)
	require.Equal(t, 503, unknown.Code)
	close(delayed.release)
	response := <-done
	require.Equal(t, 503, response.Code, "older callback must be fenced after unknown commit even for absent slot")
	require.Zero(t, live.CapturePointLease(p.ID, late.RevisionID))
}
