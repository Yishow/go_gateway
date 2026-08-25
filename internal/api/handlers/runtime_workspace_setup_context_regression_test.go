package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRuntimeHandlerWorkspaceContextIgnoresMissingSetupMappingAndTag(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx := context.Background()
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	now := time.Now().UTC()
	lastTestSuccess := true
	require.NoError(t, deviceRepo.Create(ctx, &schema.Device{
		ID:               "device-A",
		Name:             "Mixer PLC",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"192.168.1.10","port":502,"slave_id":1,"timeout":5}`,
		LastTestSuccess:  &lastTestSuccess,
		CreatedAt:        now,
		UpdatedAt:        now,
	}))
	_, err := workspaceSvc.AttachDevice(ctx, "device-A")
	require.NoError(t, err)

	mappingID := "mapping-missing"
	tagID := "tag-missing"
	handler := NewRuntimeHandler(deviceSvc, nil, nil, nil, nil, workspaceSvc)
	handler.sourceRuleReader = runtimeSetupRuleReaderStub{
		rules: []*schema.SourceRule{{
			ID:           "rule-1",
			DeviceID:     "device-A",
			StartAddress: "40001",
			Count:        1,
			DataType:     schema.DataTypeInt16,
			Enabled:      true,
		}},
		links: map[string][]*schema.SourceRuleLink{
			"rule-1": {{
				ID:        "link-1",
				RuleID:    "rule-1",
				Address:   "40001",
				PointID:   "point-1",
				MappingID: &mappingID,
				TagID:     &tagID,
			}},
		},
	}
	handler.mappingReader = runtimeSetupMissingMappingReaderStub{}
	handler.tagReader = runtimeSetupMissingTagReaderStub{}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace/runtime-context", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.WorkspaceContext(c)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())

	var body struct {
		Success bool `json:"success"`
		Data    struct {
			Setup struct {
				Mappings []struct {
					ID      string `json:"id"`
					PointID string `json:"point_id"`
					Address string `json:"address"`
					TagID   string `json:"tag_id"`
					TagKey  string `json:"tag_key"`
				} `json:"mappings"`
			} `json:"setup"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body), resp.Body.String())
	require.True(t, body.Success, resp.Body.String())
	require.Len(t, body.Data.Setup.Mappings, 1)
	require.Equal(t, "", body.Data.Setup.Mappings[0].ID)
	require.Equal(t, "point-1", body.Data.Setup.Mappings[0].PointID)
	require.Equal(t, "40001", body.Data.Setup.Mappings[0].Address)
	require.Equal(t, "", body.Data.Setup.Mappings[0].TagID)
	require.Equal(t, "", body.Data.Setup.Mappings[0].TagKey)
}

func TestRuntimeHandlerWorkspaceContextFiltersDatabaseTargetsToWorkspaceBindings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx := context.Background()
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	now := time.Now().UTC()
	require.NoError(t, deviceRepo.Create(ctx, &schema.Device{
		ID:               "device-A",
		Name:             "Mixer PLC",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"192.168.1.10","port":502,"slave_id":1,"timeout":5}`,
		CreatedAt:        now,
		UpdatedAt:        now,
	}))
	_, err := workspaceSvc.AttachDevice(ctx, "device-A")
	require.NoError(t, err)
	_, err = workspaceSvc.BindDatabaseConnector(ctx, "db-1")
	require.NoError(t, err)

	tagID := "tag-1"
	handler := NewRuntimeHandler(deviceSvc, nil, nil, nil, nil, workspaceSvc)
	handler.sourceRuleReader = runtimeSetupRuleReaderStub{
		rules: []*schema.SourceRule{{
			ID:           "rule-1",
			DeviceID:     "device-A",
			StartAddress: "40001",
			Count:        1,
			DataType:     schema.DataTypeInt16,
			Enabled:      true,
		}},
		links: map[string][]*schema.SourceRuleLink{
			"rule-1": {{
				ID:      "link-1",
				RuleID:  "rule-1",
				Address: "40001",
				PointID: "point-1",
				TagID:   &tagID,
			}},
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
				Labels:      `{"source":"source-rule","source_rule_id":"rule-1","source_rule_address":"40001"}`,
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
		records: []*schema.DatabaseTargetMapping{
			{
				ID:          "target-1",
				TagID:       "tag-1",
				ConnectorID: "db-1",
				TableSchema: "public",
				TableName:   "sensor_readings",
				ColumnName:  "value",
				WriteMode:   schema.DatabaseWriteModeInsert,
				Enabled:     true,
			},
			{
				ID:          "target-outside",
				TagID:       "tag-outside",
				ConnectorID: "db-1",
				TableSchema: "public",
				TableName:   "sensor_readings",
				ColumnName:  "rogue",
				WriteMode:   schema.DatabaseWriteModeInsert,
				Enabled:     true,
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace/runtime-context", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.WorkspaceContext(c)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())

	var body struct {
		Success bool `json:"success"`
		Data    struct {
			Setup struct {
				DatabaseTargets []struct {
					ID     string `json:"id"`
					TagID  string `json:"tag_id"`
					Column string `json:"column_name"`
				} `json:"database_targets"`
			} `json:"setup"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body), resp.Body.String())
	require.True(t, body.Success, resp.Body.String())
	require.Len(t, body.Data.Setup.DatabaseTargets, 1)
	require.Equal(t, "target-1", body.Data.Setup.DatabaseTargets[0].ID)
	require.Equal(t, "tag-1", body.Data.Setup.DatabaseTargets[0].TagID)
	require.Equal(t, "value", body.Data.Setup.DatabaseTargets[0].Column)
}

func TestRuntimeHandlerWorkspaceContextIgnoresMissingDatabaseConnector(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctx := context.Background()
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	now := time.Now().UTC()
	require.NoError(t, deviceRepo.Create(ctx, &schema.Device{
		ID:               "device-A",
		Name:             "Mixer PLC",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"192.168.1.10","port":502,"slave_id":1,"timeout":5}`,
		CreatedAt:        now,
		UpdatedAt:        now,
	}))
	_, err := workspaceSvc.AttachDevice(ctx, "device-A")
	require.NoError(t, err)
	_, err = workspaceSvc.BindDatabaseConnector(ctx, "db-missing")
	require.NoError(t, err)
	workspaceSvc.WithReadinessServices(deviceSvc, nil, runtimeSetupMissingConnectorReaderStub{}, nil)

	handler := NewRuntimeHandler(deviceSvc, nil, nil, nil, nil, workspaceSvc)
	handler.dbConnectorReader = runtimeSetupMissingConnectorReaderStub{}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace/runtime-context", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.WorkspaceContext(c)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	var body struct {
		Success bool `json:"success"`
		Data    struct {
			Setup struct {
				ReadinessSummary *struct {
					BlockingCount int `json:"blocking_count"`
					Issues        []struct {
						Code  string `json:"code"`
						Scope string `json:"scope"`
					} `json:"issues"`
				} `json:"readiness_summary"`
				DatabaseConfig *struct {
					ID string `json:"id"`
				} `json:"database_config"`
				DatabaseTargets []struct {
					ID string `json:"id"`
				} `json:"database_targets"`
			} `json:"setup"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body), resp.Body.String())
	require.True(t, body.Success, resp.Body.String())
	require.NotNil(t, body.Data.Setup.ReadinessSummary)
	require.GreaterOrEqual(t, body.Data.Setup.ReadinessSummary.BlockingCount, 1)
	foundConnectorIssue := false
	for _, issue := range body.Data.Setup.ReadinessSummary.Issues {
		if issue.Code == "database-connector-missing" && issue.Scope == "db-missing" {
			foundConnectorIssue = true
			break
		}
	}
	require.True(t, foundConnectorIssue, resp.Body.String())
	require.Nil(t, body.Data.Setup.DatabaseConfig)
	require.Len(t, body.Data.Setup.DatabaseTargets, 0)
}

type runtimeSetupMissingMappingReaderStub struct{}

func (runtimeSetupMissingMappingReaderStub) GetByID(_ context.Context, id string) (*schema.Mapping, error) {
	return nil, fmt.Errorf("%w: %s", mapping.ErrMappingNotFound, id)
}

func (runtimeSetupMissingMappingReaderStub) List(_ context.Context, filter mapping.ListFilter) ([]*schema.Mapping, error) {
	return nil, nil
}

type runtimeSetupMissingTagReaderStub struct{}

func (runtimeSetupMissingTagReaderStub) GetByID(_ context.Context, id string) (*schema.Tag, error) {
	return nil, fmt.Errorf("%w: %s", tag.ErrTagNotFound, id)
}

type runtimeSetupMissingConnectorReaderStub struct{}

func (runtimeSetupMissingConnectorReaderStub) GetByID(context.Context, string) (*schema.DatabaseConnector, error) {
	return nil, errors.New("資料庫連接器不存在")
}
