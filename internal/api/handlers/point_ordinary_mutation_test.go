package handlers

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	runtimepkg "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
	"net/http/httptest"
	"testing"
	"time"
)

type lateOrdinaryRuntime struct {
	*runtimepkg.Service
	entered, release chan struct{}
}

func (r *lateOrdinaryRuntime) UpsertPoint(p *schema.Point) {
	if p.Name == "late" {
		close(r.entered)
		<-r.release
	}
	r.Service.UpsertPoint(p)
}
func TestOrdinaryLateCallbackCannotReviveAfterFormatInvalidation(t *testing.T) {
	for _, absent := range []bool{false, true} {
		t.Run(map[bool]string{false: "installed", true: "absent"}[absent], func(t *testing.T) {
			db := openStudioV2CascadeTestDB(t)
			defer db.Close()
			devRepo := device.NewSQLRepository(db)
			now := time.Now().UTC()
			require.NoError(t, devRepo.Create(t.Context(), &schema.Device{ID: "d", Name: "d", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: `{"host":"127.0.0.1"}`, CreatedAt: now, UpdatedAt: now}))
			ws := workspace.NewService(workspace.NewSQLRepository(db))
			record, err := ws.AttachDevice(t.Context(), "d")
			require.NoError(t, err)
			revs := modbusshare.NewSQLWorkspaceRevisionStore(db)
			require.NoError(t, revs.UpdateRevision(t.Context(), record.ID, "rev-1", "ws-old"))
			share := modbusshare.NewService(nil, 64)
			share.SetWorkspaceRevisionStore(revs)
			share.SetHydrationState(modbusshare.HydrationState{WorkspaceID: record.ID, WorkspaceRevision: "ws-old", State: modbusshare.HydrationStateReady, Readiness: true})
			svc := point.NewService(point.NewSQLRepository(db), nil)
			p, err := svc.Create(t.Context(), point.CreatePointRequest{DeviceID: "d", Name: "p", Address: "HR0", DataType: schema.DataTypeUint32})
			require.NoError(t, err)
			devSvc := device.NewService(devRepo, nil)
			live, err := runtimepkg.NewService(runtimepkg.DefaultConfig(), runtimepkg.Dependencies{Scheduler: collector.NewScheduler(collector.DefaultSchedulerConfig(), nil), Writer: storage.NewSQLiteWriter(db), DeviceService: devSvc, PointService: svc, MappingService: mapping.NewService(mapping.NewSQLRepository(db)), TagService: tag.NewService(tag.NewSQLRepository(db)), PollingGroupService: pollinggroup.NewService(pollinggroup.NewSQLRepository(db))})
			require.NoError(t, err)
			wirePointProfileFixture(t, live, svc, p.DeviceID, devRepo)
			if !absent {
				live.UpsertPoint(p)
			}
			delayed := &lateOrdinaryRuntime{Service: live, entered: make(chan struct{}), release: make(chan struct{})}
			ConfigurePointFormatLifecycle(svc, share, live)
			h := NewPointHandler(svc, delayed)
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.PUT("/points/:id", h.Update)
			router.PUT("/workspace/devices/:id", NewStudioV2WorkspaceDevicesHandler(ws, devSvc, live).Update)
			send := func(body string) *httptest.ResponseRecorder {
				r := httptest.NewRequest("PUT", "/points/"+p.ID, bytes.NewBufferString(body))
				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, r)
				return w
			}
			done := make(chan *httptest.ResponseRecorder, 1)
			go func() { done <- send(`{"name":"late","data_type":"float32"}`) }()
			<-delayed.entered
			late, err := svc.GetByID(t.Context(), p.ID)
			require.NoError(t, err)
			require.Equal(t, 200, send(`{"name":"latest"}`).Code)
			current, err := svc.GetByID(t.Context(), p.ID)
			require.NoError(t, err)
			_, err = svc.MutateFormat(context.Background(), point.FormatMutationRequest{FormatScope: point.FormatScope{WorkspaceID: record.ID, DeviceID: p.DeviceID, PointID: p.ID, ExpectedWorkspaceRevision: "ws-old", ExpectedPointRevision: &current.RevisionID, RequestID: "format-after"}, FormatMode: schema.PointFormatOverride, ConfiguredFormat: "CDAB"})
			require.NoError(t, err)
			close(delayed.release)
			response := <-done
			require.Equal(t, 503, response.Code)
			require.Contains(t, response.Body.String(), `"runtime_apply_status":"pending"`)
			require.Contains(t, response.Body.String(), `configuration saved`)
			require.Zero(t, live.CapturePointLease(p.ID, late.RevisionID), "old ordinary response revived collection after successful invalidation")
			stopped, err := svc.GetByID(t.Context(), p.ID)
			require.NoError(t, err)
			blocked := send(`{"name":"pending rename"}`)
			require.Equal(t, 503, blocked.Code)
			require.NotContains(t, blocked.Body.String(), `configuration saved`)
			unchanged, err := svc.GetByID(t.Context(), p.ID)
			require.NoError(t, err)
			require.Equal(t, stopped.RevisionID, unchanged.RevisionID)
			// A legitimate existing projection/recovery installs this exact durable revision.
			recoverRequest := httptest.NewRequest("PUT", "/workspace/devices/d", bytes.NewBufferString(`{"name":"reapply saved device"}`))
			recoverRequest.Header.Set("Content-Type", "application/json")
			recovered := httptest.NewRecorder()
			router.ServeHTTP(recovered, recoverRequest)
			require.Equal(t, 200, recovered.Code, recovered.Body.String())
			require.Contains(t, recovered.Body.String(), `"runtime_apply_status":"applied"`)
			restarted := send(`{"name":"restarted rename"}`)
			require.Equal(t, 200, restarted.Code, restarted.Body.String())
			renamed, err := svc.GetByID(t.Context(), p.ID)
			require.NoError(t, err)
			require.Equal(t, "CDAB", renamed.FormatOverride)
			require.NotZero(t, live.CapturePointLease(p.ID, renamed.RevisionID))
		})
	}
}

func (r *lateOrdinaryRuntime) ApplyOrdinaryPoint(p *schema.Point, proof collector.PointMutationLease) bool {
	if p.Name == "late" {
		close(r.entered)
		<-r.release
	}
	return r.Service.ApplyOrdinaryPoint(p, proof)
}
