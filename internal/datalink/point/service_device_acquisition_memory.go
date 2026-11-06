package point

import (
	"context"
	"fmt"
	"time"

	common "go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

// IsInMemoryPollingGroups identifies a group store that shares the memory-only fixture boundary.
func (r *MemoryPollingGroupRepository) IsInMemoryPollingGroups() bool { return r != nil }

func (r *MemoryRepository) setDeviceAcquisitionInterval(ctx context.Context, s *Service, deviceID string, intervalMs int) (*schema.PollingGroup, error) {
	groups, ok := s.groupRepo.(interface{ IsInMemoryPollingGroups() bool })
	if r == nil || !ok || !groups.IsInMemoryPollingGroups() {
		return nil, fmt.Errorf("設備採集頻率需要相容的記憶體儲存庫")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	revisions := make(map[string]string)
	for id, pt := range r.points {
		if pt.DeviceID != deviceID {
			continue
		}
		revision, err := common.NewUUID()
		if err != nil {
			return nil, fmt.Errorf("建立點位 revision 失敗: %w", err)
		}
		revisions[id] = revision
	}
	group, err := s.saveMemoryDeviceAcquisitionGroup(ctx, deviceID, intervalMs)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	for id, revision := range revisions {
		pt := r.points[id]
		if pt.PollingGroupID != nil && *pt.PollingGroupID == group.ID {
			continue
		}
		pt.PollingGroupID = common.Ptr(group.ID)
		pt.RevisionID = revision
		pt.UpdatedAt = now
	}
	return group, nil
}
