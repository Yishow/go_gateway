package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newStudioV2WorkspaceDeviceRouter() *gin.Engine {
	return NewRouter(&DatalinkServices{
		Device:       device.NewService(device.NewMemoryRepository(), nil),
		Point:        point.NewService(point.NewMemoryRepository(), nil),
		Tag:          tag.NewService(tag.NewMemoryRepository()),
		Mapping:      mapping.NewService(mapping.NewMemoryRepository()),
		PollingGroup: pollinggroup.NewService(pollinggroup.NewMemoryRepository()),
		Settings:     settings.NewService(settings.NewMemoryRepository()),
		Workspace:    workspace.NewService(workspace.NewMemoryRepository()),
	})
}

func newStudioV2WorkspaceDeviceSQLRouter(t *testing.T) *gin.Engine {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "workspace.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	ws := workspace.NewService(workspace.NewSQLRepository(db))
	return NewRouter(&DatalinkServices{
		Workspace: ws, Device: device.NewService(device.NewSQLRepository(db), nil),
		WriteGroups: workspace.NewWriteGroupService(ws, workspace.NewSQLWriteGroupRepository(db)),
	})
}

func performJSONRequest(t *testing.T, router http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body failed: %v", err)
		}
	}

	req := httptest.NewRequestWithContext(context.Background(), method, path, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func decodeJSONBody(t *testing.T, resp *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body failed: %v body=%s", err, resp.Body.String())
	}

	return body
}

func decodeDeviceIDs(t *testing.T, body map[string]any) []string {
	t.Helper()

	rawData, ok := body["data"].([]any)
	if !ok {
		t.Fatalf("expected data array, got %#v", body["data"])
	}

	ids := make([]string, 0, len(rawData))
	for _, item := range rawData {
		entry, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("expected device object, got %#v", item)
		}
		id, _ := entry["id"].(string)
		ids = append(ids, id)
	}

	return ids
}

func TestNewRouter_StudioV2WorkspaceDeviceAutosaveEndpoints(t *testing.T) {
	router := newStudioV2WorkspaceDeviceRouter()

	createA := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/devices", map[string]any{
		"id":          "dev-A",
		"name":        "Line A PLC",
		"description": "primary line",
		"protocol":    "modbus_tcp",
		"connection_config": map[string]any{
			"host":     "192.168.10.10",
			"port":     502,
			"slave_id": 1,
			"timeout":  5,
		},
	})
	if createA.Code != http.StatusCreated {
		t.Fatalf("expected createA 201, got %d body=%s", createA.Code, createA.Body.String())
	}

	createB := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/devices", map[string]any{
		"id":       "dev-B",
		"name":     "Line B PLC",
		"protocol": "modbus_tcp",
		"connection_config": map[string]any{
			"host":     "192.168.10.11",
			"port":     502,
			"slave_id": 2,
			"timeout":  5,
		},
	})
	if createB.Code != http.StatusCreated {
		t.Fatalf("expected createB 201, got %d body=%s", createB.Code, createB.Body.String())
	}

	listResp := performJSONRequest(t, router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/devices", nil)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected list 200, got %d body=%s", listResp.Code, listResp.Body.String())
	}
	if ids := decodeDeviceIDs(t, decodeJSONBody(t, listResp)); len(ids) != 2 || ids[0] != "dev-A" || ids[1] != "dev-B" {
		t.Fatalf("expected persisted order [dev-A dev-B], got %v", ids)
	}

	orderResp := performJSONRequest(t, router, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/device-order", map[string]any{
		"ordered_device_ids": []string{"dev-B", "dev-A"},
	})
	if orderResp.Code != http.StatusOK {
		t.Fatalf("expected order 200, got %d body=%s", orderResp.Code, orderResp.Body.String())
	}

	reorderedResp := performJSONRequest(t, router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/devices", nil)
	if reorderedResp.Code != http.StatusOK {
		t.Fatalf("expected reordered list 200, got %d body=%s", reorderedResp.Code, reorderedResp.Body.String())
	}
	if ids := decodeDeviceIDs(t, decodeJSONBody(t, reorderedResp)); len(ids) != 2 || ids[0] != "dev-B" || ids[1] != "dev-A" {
		t.Fatalf("expected reordered devices [dev-B dev-A], got %v", ids)
	}
}

func TestNewRouter_StudioV2WorkspaceDeviceValidationDoesNotOverwriteSavedDevice(t *testing.T) {
	router := newStudioV2WorkspaceDeviceRouter()

	createResp := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/devices", map[string]any{
		"id":       "dev-A",
		"name":     "Line A PLC",
		"protocol": "modbus_tcp",
		"connection_config": map[string]any{
			"host":     "192.168.10.10",
			"port":     502,
			"slave_id": 1,
			"timeout":  5,
		},
	})
	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected create 201, got %d body=%s", createResp.Code, createResp.Body.String())
	}

	invalidUpdate := performJSONRequest(t, router, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/devices/dev-A", map[string]any{
		"name": "",
		"connection_config": map[string]any{
			"host":     "",
			"port":     502,
			"slave_id": 1,
			"timeout":  5,
		},
	})
	if invalidUpdate.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid update 400, got %d body=%s", invalidUpdate.Code, invalidUpdate.Body.String())
	}

	listResp := performJSONRequest(t, router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/devices", nil)
	body := decodeJSONBody(t, listResp)
	items, ok := body["data"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("expected single device after invalid update, got %#v", body["data"])
	}

	deviceEntry, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("expected device object, got %#v", items[0])
	}
	if deviceEntry["name"] != "Line A PLC" {
		t.Fatalf("expected saved name to stay unchanged, got %#v", deviceEntry["name"])
	}
}

func TestNewRouter_StudioV2WorkspaceDeviceDeleteCascadeRemovesWorkspaceRelation(t *testing.T) {
	router := newStudioV2WorkspaceDeviceSQLRouter(t)

	createResp := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/devices", map[string]any{
		"id":       "dev-A",
		"name":     "Line A PLC",
		"protocol": "modbus_tcp",
		"connection_config": map[string]any{
			"host":     "192.168.10.10",
			"port":     502,
			"slave_id": 1,
			"timeout":  5,
		},
	})
	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected create 201, got %d body=%s", createResp.Code, createResp.Body.String())
	}

	deleteResp := performJSONRequest(t, router, http.MethodDelete, "/api/v1/datalink/studio-v2/workspace/devices/dev-A", nil)
	if deleteResp.Code != http.StatusOK {
		t.Fatalf("expected delete 200, got %d body=%s", deleteResp.Code, deleteResp.Body.String())
	}

	listResp := performJSONRequest(t, router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/devices", nil)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected list 200, got %d body=%s", listResp.Code, listResp.Body.String())
	}
	if ids := decodeDeviceIDs(t, decodeJSONBody(t, listResp)); len(ids) != 0 {
		t.Fatalf("expected empty device list after delete, got %v", ids)
	}

	workspaceResp := performJSONRequest(t, router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace", nil)
	if workspaceResp.Code != http.StatusOK {
		t.Fatalf("expected workspace 200, got %d body=%s", workspaceResp.Code, workspaceResp.Body.String())
	}
	workspaceBody := decodeJSONBody(t, workspaceResp)
	data, ok := workspaceBody["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected workspace object, got %#v", workspaceBody["data"])
	}
	if data["status"] != string(workspace.WorkspaceStatusEmpty) {
		t.Fatalf("expected empty workspace after delete, got %#v", data["status"])
	}
	if ids, ok := data["ordered_device_ids"].([]any); !ok || len(ids) != 0 {
		t.Fatalf("expected empty workspace order after delete, got %#v", data["ordered_device_ids"])
	}
}

func TestNewRouter_StudioV2WorkspaceDeviceDeleteRefusesReferencedRecording(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	groups := workspace.NewWriteGroupService(f.workspace, workspace.NewSQLWriteGroupRepository(f.db))
	router := NewRouter(&DatalinkServices{Workspace: f.workspace, WriteGroups: groups, Device: device.NewService(device.NewSQLRepository(f.db), nil)})
	created := performJSONRequest(t, router, http.MethodPost, writeGroupsPath, f.createRequest())
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	_, group := groupSaveData(t, decodeJSONBody(t, created))
	before := lifecycleAPIState(t, f.db)
	response := performJSONRequest(t, router, http.MethodDelete, "/api/v1/datalink/studio-v2/workspace/devices/device-A", nil)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	body := decodeJSONBody(t, response)
	errBody := body["error"].(map[string]any)
	require.Equal(t, "DEVICE_REFERENCED_BY_RECORDINGS", errBody["code"])
	require.Equal(t, "remove_recordings_first", errBody["action"])
	details := errBody["details"].(map[string]any)
	require.EqualValues(t, 1, details["total"])
	require.Equal(t, []any{map[string]any{"id": group["id"], "name": group["name"]}}, details["recordings"])
	require.Equal(t, before, lifecycleAPIState(t, f.db))
	var count int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM devices WHERE id='device-A'`).Scan(&count))
	require.Equal(t, 1, count)
}
