package sourcerule

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_CandidateSnapshots_PersistRuleOwnedLocalModbusMetadata(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-lm-owned-state")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-lm-owned-state",
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
	require.NoError(t, svc.persistCandidateSnapshots(ctx, rule, links))

	snapshots, err := svc.ListCandidateSnapshots(ctx, rule.ID)
	require.NoError(t, err)

	localModbusSnapshot := localModbusSnapshotFromSnapshots(t, snapshots)
	assert.Equal(t, schema.SourceRuleCandidateStatusDeferred, localModbusSnapshot.Status)
	assert.Equal(t, rule.RevisionID, localModbusSnapshot.RevisionID)
	assert.NotEmpty(t, localModbusSnapshot.Reason)

	var payload struct {
		Candidates []map[string]any `json:"candidates"`
	}
	require.NoError(t, json.Unmarshal([]byte(localModbusSnapshot.Payload), &payload))
	require.Len(t, payload.Candidates, 1)

	candidate := payload.Candidates[0]
	assert.Equal(t, "40001", candidate["address"])
	assert.Equal(t, links[0].PointID, candidate["point_id"])
	assert.Equal(t, *links[0].TagID, candidate["tag_id"])
	assert.Equal(t, "local_modbus_output", candidate["identity"].(map[string]any)["candidate_kind"])
	assert.Equal(t, rule.ID, candidate["identity"].(map[string]any)["source_rule_id"])
	assert.NotEmpty(t, candidate["id"])
	assert.NotEmpty(t, candidate["proposed_signature"])
}

func TestService_CandidateSnapshots_PersistLocalModbusRegisterOwnershipMetadata(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-lm-register-state")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-lm-register-state",
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

	updatedAt := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	svc.SetLocalModbusMappingReader(LocalModbusMappingListFunc(func(context.Context) ([]LocalModbusMappingRecord, error) {
		return []LocalModbusMappingRecord{
			{
				TagID:     *links[0].TagID,
				Register:  7,
				DataType:  schema.DataTypeInt16,
				UpdatedAt: updatedAt,
			},
		}, nil
	}))
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
	assert.Equal(t, float64(7), candidate["register"])
	assert.Equal(t, updatedAt.Format(time.RFC3339), candidate["updated_at"])
	assert.True(t, hasCandidateScope(candidate["identity"], "local_modbus_register", "7"))
}

func TestService_Update_PreservesLocalModbusSnapshotByRevision(t *testing.T) {
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

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-lm-revision-state")
	require.NoError(t, err)

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-lm-revision-state",
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

	currentMappings := []LocalModbusMappingRecord{
		{
			TagID:     *links[0].TagID,
			Register:  7,
			DataType:  schema.DataTypeInt16,
			UpdatedAt: time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}
	svc.SetLocalModbusMappingReader(LocalModbusMappingListFunc(func(context.Context) ([]LocalModbusMappingRecord, error) {
		return append([]LocalModbusMappingRecord(nil), currentMappings...), nil
	}))
	require.NoError(t, svc.persistCandidateSnapshots(ctx, rule, links))

	oldRevisionID := rule.RevisionID
	currentMappings = []LocalModbusMappingRecord{
		{
			TagID:     *links[0].TagID,
			Register:  9,
			DataType:  schema.DataTypeInt16,
			UpdatedAt: time.Date(2025, 1, 2, 4, 5, 6, 0, time.UTC),
		},
	}
	updatedRule, err := svc.Update(ctx, rule.ID, UpdateRuleRequest{
		ScaleMultiplier:    float64Ptr(2),
		ScaleMultiplierSet: true,
	})
	require.NoError(t, err)
	require.NotEqual(t, oldRevisionID, updatedRule.RevisionID)

	oldSnapshots, err := svc.ListCandidateSnapshotsByRevision(ctx, rule.ID, oldRevisionID)
	require.NoError(t, err)
	newSnapshots, err := svc.ListCandidateSnapshotsByRevision(ctx, rule.ID, updatedRule.RevisionID)
	require.NoError(t, err)

	oldSnapshot := localModbusSnapshotFromSnapshots(t, oldSnapshots)
	newSnapshot := localModbusSnapshotFromSnapshots(t, newSnapshots)
	assert.Equal(t, oldRevisionID, oldSnapshot.RevisionID)
	assert.Equal(t, updatedRule.RevisionID, newSnapshot.RevisionID)

	var oldPayload struct {
		Candidates []map[string]any `json:"candidates"`
	}
	var newPayload struct {
		Candidates []map[string]any `json:"candidates"`
	}
	require.NoError(t, json.Unmarshal([]byte(oldSnapshot.Payload), &oldPayload))
	require.NoError(t, json.Unmarshal([]byte(newSnapshot.Payload), &newPayload))
	require.Len(t, oldPayload.Candidates, 1)
	require.Len(t, newPayload.Candidates, 1)
	assert.Equal(t, float64(7), oldPayload.Candidates[0]["register"])
	assert.Equal(t, float64(9), newPayload.Candidates[0]["register"])
}

func localModbusSnapshotFromSnapshots(
	t *testing.T,
	snapshots []*schema.SourceRuleCandidateSnapshot,
) *schema.SourceRuleCandidateSnapshot {
	t.Helper()

	for _, snapshot := range snapshots {
		if snapshot == nil || snapshot.CandidateType != schema.SourceRuleCandidateTypeLocalModbusOutputs {
			continue
		}
		return snapshot
	}

	t.Fatal("local modbus snapshot not found")
	return nil
}

func hasCandidateScope(identity any, key, value string) bool {
	identityMap, ok := identity.(map[string]any)
	if !ok {
		return false
	}

	rawScope, ok := identityMap["target_binding_scope"].([]any)
	if !ok {
		return false
	}

	for _, item := range rawScope {
		scopeItem, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if scopeItem["key"] == key && scopeItem["value"] == value {
			return true
		}
	}
	return false
}
