package sourcerule

import (
	"context"
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_CandidateSnapshots_PersistMissingLocalModbusTagAsOutOfSyncCandidate(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)
	svc.SetTagMappingServices(tagSvc, mappingSvc)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-lm-missing-tag")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-lm-missing-tag",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	links := applyRuleManagedLinks(ctx, t, repo, svc, rule.ID)
	require.Len(t, links, 1)
	require.NotNil(t, links[0].TagID)

	missingTagID := *links[0].TagID
	require.NoError(t, tagSvc.Delete(ctx, missingTagID))
	require.NoError(t, svc.persistCandidateSnapshots(ctx, rule, links))

	snapshots, err := svc.ListCandidateSnapshots(ctx, rule.ID)
	require.NoError(t, err)

	localModbusSnapshot := localModbusSnapshotFromSnapshots(t, snapshots)
	var payload struct {
		Candidates []map[string]any `json:"candidates"`
	}
	require.NoError(t, json.Unmarshal([]byte(localModbusSnapshot.Payload), &payload))
	require.Len(t, payload.Candidates, 1)

	candidate := payload.Candidates[0]
	assert.Equal(t, missingTagID, candidate["tag_id"])
	assert.Equal(t, "out_of_sync", candidate["status"])
	assert.Contains(t, candidate["blocking_reason"], "linked tag not found")
	assert.Equal(t, "40001", candidate["address"])
}
