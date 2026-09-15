package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
)

func TestNewRouter_StudioV2WorkspaceBootstrapEndpoint(t *testing.T) {
	deviceSvc := device.NewService(device.NewMemoryRepository(), nil)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository()).WithReadinessServices(deviceSvc, nil, nil, nil)

	router := NewRouter(&DatalinkServices{
		Device:       deviceSvc,
		Point:        point.NewService(point.NewMemoryRepository(), nil),
		Tag:          tag.NewService(tag.NewMemoryRepository()),
		Mapping:      mapping.NewService(mapping.NewMemoryRepository()),
		PollingGroup: pollinggroup.NewService(pollinggroup.NewMemoryRepository()),
		Settings:     settings.NewService(settings.NewMemoryRepository()),
		Workspace:    workspaceSvc,
	})

	firstReq := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace", http.NoBody)
	firstResp := httptest.NewRecorder()
	router.ServeHTTP(firstResp, firstReq)

	if firstResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", firstResp.Code, firstResp.Body.String())
	}

	var firstBody struct {
		Success bool `json:"success"`
		Data    struct {
			ID               string   `json:"id"`
			Kind             string   `json:"kind"`
			Status           string   `json:"status"`
			OrderedDeviceIDs []string `json:"ordered_device_ids"`
			ReadinessSummary struct {
				Ready         bool `json:"ready"`
				BlockingCount int  `json:"blocking_count"`
				WarningCount  int  `json:"warning_count"`
			} `json:"readiness_summary"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		} `json:"data"`
	}
	if err := json.Unmarshal(firstResp.Body.Bytes(), &firstBody); err != nil {
		t.Fatalf("decode first response failed: %v", err)
	}

	if !firstBody.Success {
		t.Fatalf("expected success body, got %s", firstResp.Body.String())
	}
	if firstBody.Data.ID == "" {
		t.Fatalf("expected workspace id, got %s", firstResp.Body.String())
	}
	if firstBody.Data.Kind != string(workspace.WorkspaceKindSingle) {
		t.Fatalf("expected single workspace kind, got %s", firstResp.Body.String())
	}
	if firstBody.Data.Status != string(workspace.WorkspaceStatusEmpty) {
		t.Fatalf("expected empty workspace status, got %s", firstResp.Body.String())
	}
	if firstBody.Data.OrderedDeviceIDs == nil {
		t.Fatalf("expected ordered_device_ids array, got %s", firstResp.Body.String())
	}
	if firstBody.Data.ReadinessSummary.Ready || firstBody.Data.ReadinessSummary.BlockingCount != 1 || firstBody.Data.ReadinessSummary.WarningCount != 0 {
		t.Fatalf("expected empty workspace readiness blocker, got %s", firstResp.Body.String())
	}
	if !strings.Contains(firstResp.Body.String(), "workspace-device-missing") {
		t.Fatalf("expected empty workspace readiness issue code, got %s", firstResp.Body.String())
	}
	if firstBody.Data.CreatedAt == "" || firstBody.Data.UpdatedAt == "" {
		t.Fatalf("expected timestamps, got %s", firstResp.Body.String())
	}

	secondReq := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace", http.NoBody)
	secondResp := httptest.NewRecorder()
	router.ServeHTTP(secondResp, secondReq)

	if secondResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", secondResp.Code, secondResp.Body.String())
	}

	var secondBody struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(secondResp.Body.Bytes(), &secondBody); err != nil {
		t.Fatalf("decode second response failed: %v", err)
	}
	if secondBody.Data.ID != firstBody.Data.ID {
		t.Fatalf("expected same workspace id across bootstrap reads, got first=%s second=%s", firstBody.Data.ID, secondBody.Data.ID)
	}
}
