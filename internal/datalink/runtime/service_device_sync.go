package runtime

import (
	"context"
	"fmt"
	"maps"

	common "go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

// UpsertDevice 將設備與其已啟用點位同步到 Runtime，供執行中 scheduler 立即開始輪詢。
//
// 若設備不是 active，則退化為移除該設備，避免 disabled / draft 設備殘留在 scheduler 中。
func (s *Service) UpsertDevice(ctx context.Context, device *schema.Device) error {
	if s == nil || device == nil || s.scheduler == nil {
		return nil
	}

	if device.Status != schema.DeviceStatusActive {
		s.RemoveDevice(device.ID)
		return nil
	}

	s.scheduler.AddDevice(device)

	if s.pointSvc == nil {
		return nil
	}

	points, err := s.pointSvc.List(ctx, point.ListFilter{
		DeviceID: common.Ptr(device.ID),
		Enabled:  common.Ptr(true),
		Limit:    100000,
	})
	if err != nil {
		return fmt.Errorf("載入設備點位失敗: %w", err)
	}

	for _, pt := range points {
		pt, err = s.ensurePointPollingGroup(ctx, pt)
		if err != nil {
			return err
		}
		if lease := s.InstallPointForMutation(pt); modernPoint(pt) && pt.Enabled && lease == 0 { //nolint:contextcheck // Admitted install is independent; TestDeviceUpsertInstallKeepsIndependentFreshProfileBoundary proves cancellation/cache isolation.
			return fmt.Errorf("point decoding profile is unavailable")
		}
	}

	s.markDeviceProjectionAlignedFromWorkspace(ctx, device)
	return nil
}

// RemoveDevice 將設備及其關聯點位從 Runtime 與 scheduler 移除。
func (s *Service) RemoveDevice(deviceID string) {
	if s == nil || deviceID == "" {
		return
	}

	s.pointMetaMu.Lock()
	s.removeDevicePointMetaLocked(deviceID)

	if s.scheduler != nil {
		s.scheduler.RemoveDevice(deviceID)
		// 設備專屬輪詢群組沒有其他點位，移除設備時一併停止其 ticker，避免閒置 goroutine 累積。
		s.scheduler.RemoveDeviceOwnedPollingGroups(deviceID)
	}
	s.pointMetaMu.Unlock()
	s.clearDeviceProjection(deviceID)
}

func (s *Service) removeDevicePointMetaLocked(deviceID string) {
	maps.DeleteFunc(s.pointMetaIndex, func(_ string, meta pointMeta) bool {
		return meta.DeviceID == deviceID
	})
}

// ensurePointPollingGroup 為缺少 polling group 的啟用點位補上設備專屬或預設群組，並確保已指派群組同步至排程器。
func (s *Service) ensurePointPollingGroup(ctx context.Context, pointRecord *schema.Point) (*schema.Point, error) {
	if pointRecord == nil || s.pointSvc == nil {
		return pointRecord, nil
	}
	if pointRecord.PollingGroupID != nil && *pointRecord.PollingGroupID != "" {
		if s.scheduler != nil {
			if group, err := s.pointSvc.GetPollingGroupByID(ctx, *pointRecord.PollingGroupID); err == nil && group != nil && group.Enabled {
				if existing, ok := s.scheduler.PollingGroup(group.ID); !ok || existing.IntervalMs != group.IntervalMs {
					s.scheduler.AddPollingGroup(group)
				}
			}
		}
		return pointRecord, nil
	}

	group, err := s.pointSvc.EnsureDevicePollingGroup(ctx, pointRecord.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("確保預設輪詢群組失敗: %w", err)
	}
	if group == nil {
		return pointRecord, nil
	}
	if s.scheduler != nil && group.Enabled {
		if existing, ok := s.scheduler.PollingGroup(group.ID); !ok || existing.IntervalMs != group.IntervalMs {
			s.scheduler.AddPollingGroup(group)
		}
	}

	groupID := group.ID
	updatedPoint, err := s.pointSvc.Update(ctx, pointRecord.ID, point.UpdatePointRequest{
		PollingGroupID: &groupID,
	})
	if err != nil {
		return nil, fmt.Errorf("設定點位輪詢群組失敗: %w", err)
	}
	return updatedPoint, nil
}
