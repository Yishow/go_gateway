package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type runtimeProjectionReaderStub struct {
	projection *workspace.RuntimeProjection
}

func (s runtimeProjectionReaderStub) RuntimeProjection(context.Context) (*workspace.RuntimeProjection, error) {
	return s.projection, nil
}

func TestRuntimeHandler_StatusIncludesProjectionAlignment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fixture := newRuntimeHandlerProjectionFixture(t, false)
	handler := NewRuntimeHandler(nil, nil, nil, nil, fixture.runtimeSvc, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/runtime/status?device_id=dev-A", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Status(c)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}
	var body struct {
		Data struct {
			Collectors []struct {
				ProjectionAlignment        string `json:"projection_alignment"`
				RuntimeProjectionVersion   string `json:"runtime_projection_version"`
				WorkspaceProjectionVersion string `json:"workspace_projection_version"`
			} `json:"collectors"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response failed: %v body=%s", err, resp.Body.String())
	}
	if len(body.Data.Collectors) != 1 {
		t.Fatalf("expected one collector, got %+v", body.Data.Collectors)
	}
	if body.Data.Collectors[0].ProjectionAlignment != "stale" {
		t.Fatalf("expected stale projection, got %+v", body.Data.Collectors[0])
	}
	if body.Data.Collectors[0].RuntimeProjectionVersion != "projection-v12" {
		t.Fatalf("expected runtime projection v12, got %+v", body.Data.Collectors[0])
	}
	if body.Data.Collectors[0].WorkspaceProjectionVersion != "projection-v13" {
		t.Fatalf("expected workspace projection v13, got %+v", body.Data.Collectors[0])
	}
}

func TestRuntimeHandler_WorkspaceContextIncludesProjectionAlignment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fixture := newRuntimeHandlerProjectionFixture(t, true)
	handler := NewRuntimeHandler(fixture.deviceSvc, nil, nil, nil, fixture.runtimeSvc, fixture.workspaceSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace/runtime-context", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.WorkspaceContext(c)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}
	var body struct {
		Data struct {
			Devices []struct {
				ProjectionAlignment string `json:"projection_alignment"`
			} `json:"devices"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response failed: %v body=%s", err, resp.Body.String())
	}
	if len(body.Data.Devices) != 1 {
		t.Fatalf("expected one runtime context device, got %+v", body.Data.Devices)
	}
	if body.Data.Devices[0].ProjectionAlignment != "stale" {
		t.Fatalf("expected stale projection in workspace context, got %+v", body.Data.Devices[0])
	}
}

type runtimeHandlerProjectionFixture struct {
	deviceSvc    *device.Service
	runtimeSvc   *datalinkruntime.Service
	workspaceSvc *workspace.Service
}

func newRuntimeHandlerProjectionFixture(t *testing.T, attachWorkspaceDevice bool) runtimeHandlerProjectionFixture {
	t.Helper()
	ctx := context.Background()

	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	groupSvc := pollinggroup.NewService(pollinggroup.NewMemoryRepository())
	deviceRecord, err := deviceSvc.Create(ctx, device.CreateDeviceRequest{
		ID:       "dev-A",
		Name:     "Device A",
		Protocol: schema.ProtocolModbusTCP,
		ConnectionConfig: map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     502,
			"slave_id": 1,
			"timeout":  1,
		},
	})
	if err != nil {
		t.Fatalf("create device failed: %v", err)
	}
	deviceRecord.Status = schema.DeviceStatusActive
	if err := deviceRepo.Update(ctx, deviceRecord); err != nil {
		t.Fatalf("activate device failed: %v", err)
	}
	if attachWorkspaceDevice {
		if _, err := workspaceSvc.AttachDevice(ctx, deviceRecord.ID); err != nil {
			t.Fatalf("attach workspace device failed: %v", err)
		}
	}

	activeProjection := runtimeHandlerProjection("projection-v12", deviceRecord)
	latestProjection := runtimeHandlerProjection("projection-v13", deviceRecord)
	runtimeSvc, err := datalinkruntime.NewService(datalinkruntime.DefaultConfig(), datalinkruntime.Dependencies{
		Scheduler:           collector.NewScheduler(collector.DefaultSchedulerConfig(), nil),
		Writer:              storage.NewMemoryStorage(16),
		DeviceService:       deviceSvc,
		PointService:        pointSvc,
		MappingService:      mapping.NewService(mapping.NewMemoryRepository()),
		TagService:          tag.NewService(tag.NewMemoryRepository()),
		PollingGroupService: groupSvc,
		WorkspaceProjection: runtimeProjectionReaderStub{projection: activeProjection},
	})
	if err != nil {
		t.Fatalf("create runtime service failed: %v", err)
	}
	if err := runtimeSvc.Start(ctx); err != nil {
		t.Fatalf("start runtime failed: %v", err)
	}
	t.Cleanup(func() {
		if err := runtimeSvc.Stop(context.Background()); err != nil {
			t.Fatalf("stop runtime failed: %v", err)
		}
	})
	runtimeSvc.SetWorkspaceProjectionReader(runtimeProjectionReaderStub{projection: latestProjection})

	return runtimeHandlerProjectionFixture{
		deviceSvc:    deviceSvc,
		runtimeSvc:   runtimeSvc,
		workspaceSvc: workspaceSvc,
	}
}

func runtimeHandlerProjection(version string, deviceRecord *schema.Device) *workspace.RuntimeProjection {
	return &workspace.RuntimeProjection{
		WorkspaceID: "ws-1",
		Version:     version,
		Alignment:   workspace.RuntimeProjectionAlignmentAligned,
		DeviceIDs:   []string{deviceRecord.ID},
		Devices:     []*schema.Device{deviceRecord},
	}
}
