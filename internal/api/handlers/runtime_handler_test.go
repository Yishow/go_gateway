package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

func TestRuntimeHandler_WorkspaceContextReturnsOrderedDevicesAndDefaultDeviceID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	now := time.Now().UTC()
	lastTestSuccess := true

	if err := deviceRepo.Create(context.Background(), &schema.Device{
		ID:               "device-B",
		Name:             "Unavailable Filler",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ReadinessStatus:  `{"availability_status":"unavailable","availability_reason":"invalid Step 1 configuration"}`,
		ConnectionConfig: `{"host":"192.168.1.11","port":502,"slave_id":2,"timeout":5}`,
		CreatedAt:        now,
		UpdatedAt:        now,
	}); err != nil {
		t.Fatalf("create device-B failed: %v", err)
	}
	if err := deviceRepo.Create(context.Background(), &schema.Device{
		ID:               "device-A",
		Name:             "Healthy Mixer",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"192.168.1.10","port":502,"slave_id":1,"timeout":5}`,
		LastTestSuccess:  &lastTestSuccess,
		CreatedAt:        now,
		UpdatedAt:        now,
	}); err != nil {
		t.Fatalf("create device-A failed: %v", err)
	}

	if _, err := workspaceSvc.AttachDevice(context.Background(), "device-B"); err != nil {
		t.Fatalf("attach device-B failed: %v", err)
	}
	if _, err := workspaceSvc.AttachDevice(context.Background(), "device-A"); err != nil {
		t.Fatalf("attach device-A failed: %v", err)
	}

	workspaceSvc.WithReadinessServices(deviceSvc, nil, nil, nil)
	handler := NewRuntimeHandler(deviceSvc, nil, nil, nil, nil, workspaceSvc)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/runtime-context", http.NoBody)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.WorkspaceContext(c)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	var body struct {
		Success bool `json:"success"`
		Data    struct {
			WorkspaceID     string  `json:"workspace_id"`
			DefaultDeviceID *string `json:"default_device_id"`
			Devices         []struct {
				DeviceID           string  `json:"device_id"`
				Name               string  `json:"name"`
				Running            bool    `json:"running"`
				AvailabilityStatus string  `json:"availability_status"`
				AvailabilityReason *string `json:"availability_reason"`
			} `json:"devices"`
			Setup struct {
				ReadinessSummary *workspace.ReadinessSummary `json:"readiness_summary"`
			} `json:"setup"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response failed: %v body=%s", err, resp.Body.String())
	}

	if !body.Success {
		t.Fatalf("expected success response, got %s", resp.Body.String())
	}
	if body.Data.WorkspaceID == "" {
		t.Fatalf("expected workspace id, got %s", resp.Body.String())
	}
	if len(body.Data.Devices) != 2 {
		t.Fatalf("expected two devices, got %+v", body.Data.Devices)
	}
	if body.Data.Devices[0].DeviceID != "device-B" || body.Data.Devices[1].DeviceID != "device-A" {
		t.Fatalf("expected ordered devices [device-B device-A], got %+v", body.Data.Devices)
	}
	if body.Data.Devices[0].AvailabilityStatus != "unavailable" {
		t.Fatalf("expected first device unavailable, got %+v", body.Data.Devices[0])
	}
	if body.Data.Devices[0].AvailabilityReason == nil || *body.Data.Devices[0].AvailabilityReason != "invalid Step 1 configuration" {
		t.Fatalf("expected unavailable reason, got %+v", body.Data.Devices[0])
	}
	if body.Data.DefaultDeviceID == nil || *body.Data.DefaultDeviceID != "device-A" {
		t.Fatalf("expected default device-A, got %+v", body.Data.DefaultDeviceID)
	}
	if body.Data.Setup.ReadinessSummary == nil {
		t.Fatalf("expected readiness summary in runtime context, got %s", resp.Body.String())
	}
	if body.Data.Setup.ReadinessSummary.Ready {
		t.Fatalf("expected workspace readiness to block unavailable device, got %+v", body.Data.Setup.ReadinessSummary)
	}
	if body.Data.Setup.ReadinessSummary.BlockingCount != 1 {
		t.Fatalf("expected one readiness blocker, got %+v", body.Data.Setup.ReadinessSummary)
	}
}

func TestRuntimeHandler_WorkspaceContextMissingDeviceUsesTypedError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	_, err := workspaceSvc.AttachDevice(context.Background(), "device-missing")
	if err != nil {
		t.Fatalf("attach missing device failed: %v", err)
	}
	handler := NewRuntimeHandler(device.NewService(device.NewMemoryRepository(), nil), nil, nil, nil, nil, workspaceSvc)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = newHandlerTestRequest(http.MethodGet, "/runtime-context", nil)
	handler.WorkspaceContext(c)

	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", resp.Code, resp.Body.String())
	}
	var body struct {
		Error TypedAPIErrorEnvelope `json:"error"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode typed error failed: %v", err)
	}
	if body.Error.Code != ErrCodeRuntimeDeviceNotFound || body.Error.RequestID == "" || body.Error.Retryable {
		t.Fatalf("unexpected missing-device error: %+v", body.Error)
	}
}

func TestRuntimeHandler_WorkspaceContextReturnsSetupConditions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx := context.Background()
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	now := time.Now().UTC()
	if err := deviceRepo.Create(ctx, &schema.Device{
		ID:               "device-A",
		Name:             "Mixer PLC",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"192.168.1.10","port":502,"slave_id":1,"timeout":5}`,
		CreatedAt:        now,
		UpdatedAt:        now,
	}); err != nil {
		t.Fatalf("create device-A failed: %v", err)
	}
	if _, err := workspaceSvc.AttachDevice(ctx, "device-A"); err != nil {
		t.Fatalf("attach device-A failed: %v", err)
	}
	if _, err := workspaceSvc.BindDatabaseConnector(ctx, "db-1"); err != nil {
		t.Fatalf("bind database failed: %v", err)
	}

	mappingID := "mapping-1"
	tagID := "tag-1"
	handler := NewRuntimeHandler(deviceSvc, nil, nil, nil, nil, workspaceSvc)
	handler.sourceRuleReader = runtimeSetupRuleReaderStub{
		rules: []*schema.SourceRule{{
			ID:           "rule-1",
			DeviceID:     "device-A",
			StartAddress: "40001",
			Count:        8,
			DataType:     schema.DataTypeInt16,
			NamingPrefix: "LINE_",
			Enabled:      true,
			RevisionID:   "rule-v1",
		}},
		links: map[string][]*schema.SourceRuleLink{
			"rule-1": {{
				ID:        "link-1",
				RuleID:    "rule-1",
				Address:   "40001",
				PointID:   "point-1",
				TagID:     &tagID,
				MappingID: &mappingID,
			}},
		},
	}
	handler.mappingReader = runtimeSetupMappingReaderStub{
		records: map[string]*schema.Mapping{
			"mapping-1": {
				ID:      "mapping-1",
				PointID: "point-1",
				TagID:   "tag-1",
				Status:  schema.MappingStatusActive,
				Enabled: true,
			},
		},
	}
	handler.tagReader = runtimeSetupTagReaderStub{
		records: map[string]*schema.Tag{
			"tag-1": {
				ID:          "tag-1",
				Key:         "line01.temp.inlet",
				DisplayName: "入口溫度",
				Unit:        "°C",
				DataType:    schema.DataTypeInt16,
			},
		},
	}
	handler.dbConnectorReader = runtimeSetupConnectorReaderStub{
		record: &schema.DatabaseConnector{
			ID:                          "db-1",
			Name:                        "PostgreSQL Connector",
			Kind:                        schema.DatabaseConnectorKindPostgres,
			ConnectionConfig:            `{"database":"gateway_metrics","schema":"public","table":"sensor_readings","write_mode":"insert"}`,
			Status:                      schema.DatabaseConnectorStatusReady,
			DefaultWriteIntervalSeconds: 5,
			CreatedAt:                   now,
			UpdatedAt:                   now,
		},
	}
	handler.dbTargetReader = runtimeSetupTargetReaderStub{
		records: []*schema.DatabaseTargetMapping{{
			ID:          "target-1",
			TagID:       "tag-1",
			ConnectorID: "db-1",
			TableSchema: "public",
			TableName:   "sensor_readings",
			ColumnName:  "value",
			WriteMode:   schema.DatabaseWriteModeInsert,
			Enabled:     true,
		}},
	}

	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/runtime-context", http.NoBody)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.WorkspaceContext(c)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	var body struct {
		Success bool `json:"success"`
		Data    struct {
			Setup struct {
				SourceRules []struct {
					StartAddress string `json:"start_address"`
					Count        int    `json:"count"`
				} `json:"source_rules"`
				Mappings []struct {
					Address     string `json:"address"`
					TagKey      string `json:"tag_key"`
					DisplayName string `json:"display_name"`
					Unit        string `json:"unit"`
				} `json:"mappings"`
				DatabaseConfig *struct {
					Name                 string `json:"name"`
					Database             string `json:"database"`
					Table                string `json:"table"`
					WriteIntervalSeconds int    `json:"write_interval_seconds"`
				} `json:"database_config"`
				DatabaseTargets []struct {
					ColumnName string `json:"column_name"`
				} `json:"database_targets"`
			} `json:"setup"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response failed: %v body=%s", err, resp.Body.String())
	}
	if !body.Success {
		t.Fatalf("expected success response, got %s", resp.Body.String())
	}
	if len(body.Data.Setup.SourceRules) != 1 || body.Data.Setup.SourceRules[0].StartAddress != "40001" || body.Data.Setup.SourceRules[0].Count != 8 {
		t.Fatalf("expected source rule summary, got %+v", body.Data.Setup.SourceRules)
	}
	if len(body.Data.Setup.Mappings) != 1 || body.Data.Setup.Mappings[0].TagKey != "line01.temp.inlet" || body.Data.Setup.Mappings[0].DisplayName != "入口溫度" {
		t.Fatalf("expected mapping summary with tag display data, got %+v", body.Data.Setup.Mappings)
	}
	if body.Data.Setup.DatabaseConfig == nil || body.Data.Setup.DatabaseConfig.Name != "PostgreSQL Connector" || body.Data.Setup.DatabaseConfig.Database != "gateway_metrics" {
		t.Fatalf("expected database config summary, got %+v", body.Data.Setup.DatabaseConfig)
	}
	if len(body.Data.Setup.DatabaseTargets) != 1 || body.Data.Setup.DatabaseTargets[0].ColumnName != "value" {
		t.Fatalf("expected database target summary, got %+v", body.Data.Setup.DatabaseTargets)
	}
}

type runtimeSetupRuleReaderStub struct {
	rules []*schema.SourceRule
	links map[string][]*schema.SourceRuleLink
}

func (s runtimeSetupRuleReaderStub) ListByDeviceIDs(context.Context, []string) ([]*schema.SourceRule, error) {
	return s.rules, nil
}

func (s runtimeSetupRuleReaderStub) ListLinks(_ context.Context, ruleID string) ([]*schema.SourceRuleLink, error) {
	return s.links[ruleID], nil
}

type runtimeSetupMappingReaderStub struct {
	records map[string]*schema.Mapping
}

func (s runtimeSetupMappingReaderStub) GetByID(_ context.Context, id string) (*schema.Mapping, error) {
	return s.records[id], nil
}

func (s runtimeSetupMappingReaderStub) List(_ context.Context, filter mapping.ListFilter) ([]*schema.Mapping, error) {
	if filter.PointID == nil {
		return nil, nil
	}
	result := make([]*schema.Mapping, 0)
	for _, record := range s.records {
		if record != nil && record.PointID == *filter.PointID {
			result = append(result, record)
		}
	}
	return result, nil
}

type runtimeSetupTagReaderStub struct {
	records map[string]*schema.Tag
}

func (s runtimeSetupTagReaderStub) GetByID(_ context.Context, id string) (*schema.Tag, error) {
	return s.records[id], nil
}

type runtimeSetupConnectorReaderStub struct {
	record *schema.DatabaseConnector
}

func (s runtimeSetupConnectorReaderStub) GetByID(context.Context, string) (*schema.DatabaseConnector, error) {
	return s.record, nil
}

type runtimeSetupTargetReaderStub struct {
	records []*schema.DatabaseTargetMapping
}

func (s runtimeSetupTargetReaderStub) List(context.Context, dbtarget.TargetMappingListFilter) ([]*schema.DatabaseTargetMapping, error) {
	return s.records, nil
}
