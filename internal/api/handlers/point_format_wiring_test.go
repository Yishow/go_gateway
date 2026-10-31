package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	runtimepkg "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/workspace"
)

func TestPointFormatHTTPWiringInvalidatesOnceAndPreservesLegacyShape(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		t.Run(map[bool]string{false: "scoped", true: "legacy"}[legacy], func(t *testing.T) {
			db := openStudioV2CascadeTestDB(t)
			t.Cleanup(func() { _ = db.Close() })
			devRepo := device.NewSQLRepository(db)
			now := time.Now().UTC()
			require.NoError(t, devRepo.Create(t.Context(), &schema.Device{ID: "format-device", Name: "format", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: `{"data_format":"ABCD","host":"secret://endpoint"}`, CreatedAt: now, UpdatedAt: now}))
			ws := workspace.NewService(workspace.NewSQLRepository(db))
			record, err := ws.AttachDevice(t.Context(), "format-device")
			require.NoError(t, err)
			revisions := modbusshare.NewSQLWorkspaceRevisionStore(db)
			require.NoError(t, revisions.UpdateRevision(t.Context(), record.ID, "rev-1", "ws-old"))
			share := modbusshare.NewService(nil, 64)
			share.SetWorkspaceRevisionStore(revisions)
			share.SetHydrationState(modbusshare.HydrationState{WorkspaceID: record.ID, WorkspaceRevision: "ws-old", State: modbusshare.HydrationStateReady, Readiness: true, ReadinessToken: "old-proof"})
			svc := point.NewService(point.NewSQLRepository(db), nil)
			p, err := svc.Create(t.Context(), point.CreatePointRequest{DeviceID: "format-device", Name: "point", Address: "HR0", DataType: schema.DataTypeUint32})
			require.NoError(t, err)
			live, err := runtimepkg.NewService(runtimepkg.Config{Writer: storage.NewSQLiteWriter(db)})
			require.NoError(t, err)
			wirePointProfileFixture(t, live, svc, p.DeviceID, devRepo)
			live.UpsertPoint(p)
			ConfigurePointFormatLifecycle(svc, share, live)
			h := NewPointHandler(svc, live)
			router := gin.New()
			router.PUT("/points/:id", h.Update)
			router.PUT("/points/:id/format", h.UpdateFormat)
			router.GET("/points/:id/format/operations/:request_id", h.InspectFormatOperation)
			router.GET("/points/:id/format", h.ReadFormat)
			read := httptest.NewRecorder()
			router.ServeHTTP(read, httptest.NewRequest("GET", "/points/"+p.ID+"/format?workspace_id="+record.ID+"&device_id="+p.DeviceID, nil))
			require.Equal(t, 200, read.Code, read.Body.String())
			var initial PointFormatReadResponse
			require.NoError(t, json.Unmarshal(read.Body.Bytes(), &initial))
			require.Equal(t, "ws-old", initial.Data.WorkspaceRevision)
			payload := map[string]any{"workspace_id": record.ID, "device_id": p.DeviceID, "point_id": p.ID, "expected_workspace_revision": initial.Data.WorkspaceRevision, "expected_point_revision": initial.Data.PointRevision, "request_id": "format-request"}
			path := "/points/" + p.ID + "/format"
			if legacy {
				path = "/points/" + p.ID
				payload["data_format"] = " cdab "
				payload["name"] = "renamed"
			} else {
				payload["format_mode"] = "override"
				payload["configured_format"] = "CDAB"
			}
			raw, err := json.Marshal(payload)
			require.NoError(t, err)
			send := func() *httptest.ResponseRecorder {
				r := httptest.NewRequest("PUT", path, bytes.NewReader(raw))
				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, r)
				return w
			}
			ready, err := share.CheckHydration(t.Context())
			require.NoError(t, err)
			pending := ready
			pending.Readiness = false
			pending.State = modbusshare.HydrationStatePending
			share.SetHydrationState(pending)
			blocked := send()
			require.Equal(t, 503, blocked.Code)
			untouched, err := svc.GetByID(t.Context(), p.ID)
			require.NoError(t, err)
			require.Equal(t, p.RevisionID, untouched.RevisionID)
			share.SetHydrationState(ready)
			first := send()
			require.Equal(t, 200, first.Code, first.Body.String())
			require.NotContains(t, first.Body.String(), "secret://endpoint")
			require.Contains(t, first.Body.String(), `"runtime_apply_status":"pending"`)
			if legacy {
				var response struct {
					Data schema.Point `json:"data"`
				}
				require.NoError(t, json.Unmarshal(first.Body.Bytes(), &response))
				require.Equal(t, p.ID, response.Data.ID)
				require.Equal(t, "renamed", response.Data.Name)
			}
			state, err := share.CheckHydration(t.Context())
			require.NoError(t, err)
			require.False(t, state.Readiness)
			require.Empty(t, state.ReadinessToken)
			require.True(t, live.SuspendPointRevision(p.ID, "not-the-old-revision"), "the old install was removed, not replaced by the new format")
			// Replaying the successful request cannot invalidate a later runtime install/readiness.
			current, err := svc.GetByID(t.Context(), p.ID)
			require.NoError(t, err)
			current.RevisionID = "new-install"
			require.NoError(t, svc.Restore(t.Context(), current))
			live.UpsertPoint(current)
			latest := state
			latest.WorkspaceRevision = "latest"
			latest.Readiness = true
			latest.ReadinessToken = "new-proof"
			latest.State = modbusshare.HydrationStateReady
			share.SetHydrationState(latest)
			second := send()
			require.Equal(t, 200, second.Code, second.Body.String())
			require.JSONEq(t, first.Body.String(), second.Body.String())
			require.False(t, live.SuspendPointRevision(p.ID, p.RevisionID), "replay preserved the newer installation")
			state, err = share.CheckHydration(t.Context())
			require.NoError(t, err)
			require.Equal(t, latest, state)
			r := httptest.NewRequest("GET", "/points/"+p.ID+"/format/operations/format-request?workspace_id="+record.ID+"&device_id="+p.DeviceID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
			require.Equal(t, 200, w.Code)
			require.Contains(t, w.Body.String(), `"invalidation_status":"invalidated"`)
		})
	}
}

type cancelledPointFormatLifecycle struct {
	*pointFormatLifecycle
	cancel context.CancelFunc
}

func (l *cancelledPointFormatLifecycle) InvalidatePointFormat(ctx context.Context, change point.FormatInvalidation) error {
	l.cancel()
	return l.pointFormatLifecycle.InvalidatePointFormat(ctx, change)
}

func TestPointFormatCommitThenCancelledRequestStillFailsClosed(t *testing.T) {
	db := openStudioV2CascadeTestDB(t)
	t.Cleanup(func() { _ = db.Close() })
	now := time.Now().UTC()
	devRepo := device.NewSQLRepository(db)
	require.NoError(t, devRepo.Create(t.Context(), &schema.Device{ID: "cancel-device", Name: "cancel", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: `{}`, CreatedAt: now, UpdatedAt: now}))
	ws := workspace.NewService(workspace.NewSQLRepository(db))
	record, err := ws.AttachDevice(t.Context(), "cancel-device")
	require.NoError(t, err)
	revisions := modbusshare.NewSQLWorkspaceRevisionStore(db)
	require.NoError(t, revisions.UpdateRevision(t.Context(), record.ID, "rev-1", "ws-old"))
	share := modbusshare.NewService(nil, 64)
	share.SetWorkspaceRevisionStore(revisions)
	share.SetHydrationState(modbusshare.HydrationState{WorkspaceID: record.ID, WorkspaceRevision: "ws-old", State: modbusshare.HydrationStateReady, Readiness: true, ReadinessToken: "old-proof"})
	svc := point.NewService(point.NewSQLRepository(db), nil)
	p, err := svc.Create(t.Context(), point.CreatePointRequest{DeviceID: "cancel-device", Name: "point", Address: "HR0", DataType: schema.DataTypeUint32})
	require.NoError(t, err)
	live, err := runtimepkg.NewService(runtimepkg.Config{Writer: storage.NewSQLiteWriter(db)})
	require.NoError(t, err)
	wirePointProfileFixture(t, live, svc, p.DeviceID, devRepo)
	live.UpsertPoint(p)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	svc.SetFormatInvalidator(&cancelledPointFormatLifecycle{pointFormatLifecycle: &pointFormatLifecycle{points: svc, share: share, runtime: live}, cancel: cancel})
	req := point.FormatMutationRequest{FormatScope: point.FormatScope{WorkspaceID: record.ID, DeviceID: p.DeviceID, PointID: p.ID, ExpectedWorkspaceRevision: "ws-old", ExpectedPointRevision: &p.RevisionID, RequestID: "cancel-request"}, FormatMode: schema.PointFormatOverride, ConfiguredFormat: "CDAB"}
	_, err = svc.MutateFormat(ctx, req)
	require.ErrorIs(t, err, point.ErrFormatUnavailable)
	require.True(t, live.SuspendPointRevision(p.ID, "not-old"), "cancelled HTTP context must not leave old acquisition installed")
	state, err := share.CheckHydration(t.Context())
	require.NoError(t, err)
	require.False(t, state.Readiness)
	require.Empty(t, state.ReadinessToken)
	status, err := svc.InspectFormatMutation(t.Context(), req.FormatScope)
	require.NoError(t, err)
	require.Equal(t, "pending", status.InvalidationStatus)
	// The pending receipt does not replay the callback against a newer install.
	newer := status.Point
	newer.RevisionID = "newer"
	require.NoError(t, svc.Restore(t.Context(), &newer))
	live.UpsertPoint(&newer)
	_, err = svc.MutateFormat(t.Context(), req)
	require.ErrorIs(t, err, point.ErrFormatUnavailable)
	require.False(t, live.SuspendPointRevision(p.ID, p.RevisionID))
}
