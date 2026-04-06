package sourcerule

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func seedActivationReadyDevice(ctx context.Context, repo *device.MemoryRepository, id string) (*schema.Device, error) {
	record, err := seedDeviceWithStatus(ctx, repo, id, schema.DeviceStatusActive)
	if err != nil {
		return nil, err
	}
	success := true
	record.LastTestSuccess = &success
	record.LastTestError = ""
	if err := repo.Update(ctx, record); err != nil {
		return nil, err
	}
	return repo.GetByID(ctx, id)
}

func float64Ptr(value float64) *float64 {
	return &value
}

func applyRuleManagedLinks(t *testing.T, ctx context.Context, repo *MemoryRepository, svc *Service, ruleID string) []*schema.SourceRuleLink {
	t.Helper()

	rule, err := svc.GetByID(ctx, ruleID)
	require.NoError(t, err)

	links, err := svc.ListLinks(ctx, ruleID)
	require.NoError(t, err)
	require.NotEmpty(t, links)

	result := tagMappingSyncResult{
		updatedMappings: make(map[string]mappingRollbackState),
		updatedTags:     make(map[string]tagRollbackState),
	}
	for _, link := range links {
		pointRecord, err := svc.pointSvc.GetByID(ctx, link.PointID)
		require.NoError(t, err)

		tagRecord, mappingRecord, err := svc.ensureRuleTagMapping(ctx, nil, rule, pointRecord, link, rule.Enabled, &result)
		require.NoError(t, err)

		link.TagID = stringPtr(tagRecord.ID)
		link.MappingID = stringPtr(mappingRecord.ID)
		link.UpdatedAt = time.Now()
	}

	require.NoError(t, repo.DeleteLinks(ctx, ruleID))
	require.NoError(t, repo.CreateLinks(ctx, links))

	return links
}

func firstTagCandidateFromMemoryRepo(t *testing.T, repo *MemoryRepository, ruleID, revisionID string) schema.SourceRuleTagCandidate {
	t.Helper()

	snapshots, err := repo.ListCandidateSnapshots(context.Background(), ruleID, revisionID)
	require.NoError(t, err)

	for _, snapshot := range snapshots {
		if snapshot.CandidateType != schema.SourceRuleCandidateTypeTags {
			continue
		}
		var payload tagSnapshotPayload
		require.NoError(t, json.Unmarshal([]byte(snapshot.Payload), &payload))
		require.NotEmpty(t, payload.Candidates)
		return payload.Candidates[0]
	}

	t.Fatalf("tag snapshot not found for rule %s revision %s", ruleID, revisionID)
	return schema.SourceRuleTagCandidate{}
}
