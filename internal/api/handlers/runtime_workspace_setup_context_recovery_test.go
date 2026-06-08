package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRuntimeHandlerWorkspaceContextRecoversLostLinkIDsFromPointMapping(t *testing.T) {
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
			}},
		},
	}
	handler.mappingReader = runtimeSetupRecoverableMappingReaderStub{
		byPoint: map[string][]*schema.Mapping{
			"point-1": {{
				ID:      "mapping-1",
				PointID: "point-1",
				TagID:   "tag-1",
				Status:  schema.MappingStatusActive,
				Enabled: true,
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
					TagID   string `json:"tag_id"`
					TagKey  string `json:"tag_key"`
					Address string `json:"address"`
				} `json:"mappings"`
			} `json:"setup"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body), resp.Body.String())
	require.True(t, body.Success, resp.Body.String())
	require.Len(t, body.Data.Setup.Mappings, 1)
	require.Equal(t, "mapping-1", body.Data.Setup.Mappings[0].ID)
	require.Equal(t, "tag-1", body.Data.Setup.Mappings[0].TagID)
	require.Equal(t, "line01.temp.inlet", body.Data.Setup.Mappings[0].TagKey)
	require.Equal(t, "40001", body.Data.Setup.Mappings[0].Address)
}

func TestRuntimeHandlerWorkspaceContextIgnoresMismatchedOwnedBinding(t *testing.T) {
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

	mappingID := "mapping-1"
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
				ID:        "link-1",
				RuleID:    "rule-1",
				Address:   "40001",
				PointID:   "point-1",
				MappingID: &mappingID,
				TagID:     &tagID,
			}},
		},
	}
	handler.mappingReader = runtimeSetupRecoverableMappingReaderStub{
		byID: map[string]*schema.Mapping{
			"mapping-1": {
				ID:      "mapping-1",
				PointID: "point-1",
				TagID:   "tag-1",
				Status:  schema.MappingStatusActive,
				Enabled: true,
			},
		},
		byPoint: map[string][]*schema.Mapping{
			"point-1": {{
				ID:      "mapping-1",
				PointID: "point-1",
				TagID:   "tag-1",
				Status:  schema.MappingStatusActive,
				Enabled: true,
			}},
		},
	}
	handler.tagReader = runtimeSetupTagReaderStub{
		records: map[string]*schema.Tag{
			"tag-1": {
				ID:          "tag-1",
				Key:         "line99.temp.rogue",
				DisplayName: "別的規則標籤",
				Unit:        "°C",
				DataType:    schema.DataTypeInt16,
				Labels:      `{"source":"source-rule","source_rule_id":"rule-9","source_rule_address":"49999"}`,
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
				Mappings []struct {
					ID     string `json:"id"`
					TagID  string `json:"tag_id"`
					TagKey string `json:"tag_key"`
				} `json:"mappings"`
			} `json:"setup"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body), resp.Body.String())
	require.True(t, body.Success, resp.Body.String())
	require.Len(t, body.Data.Setup.Mappings, 1)
	require.Equal(t, "", body.Data.Setup.Mappings[0].ID)
	require.Equal(t, "", body.Data.Setup.Mappings[0].TagID)
	require.Equal(t, "", body.Data.Setup.Mappings[0].TagKey)
}

type runtimeSetupRecoverableMappingReaderStub struct {
	byID    map[string]*schema.Mapping
	byPoint map[string][]*schema.Mapping
}

func (s runtimeSetupRecoverableMappingReaderStub) GetByID(_ context.Context, id string) (*schema.Mapping, error) {
	if record, ok := s.byID[id]; ok {
		return record, nil
	}
	return nil, nil
}

func (s runtimeSetupRecoverableMappingReaderStub) List(_ context.Context, filter mapping.ListFilter) ([]*schema.Mapping, error) {
	if filter.PointID == nil {
		return nil, nil
	}
	return s.byPoint[*filter.PointID], nil
}