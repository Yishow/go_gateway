package runtime

import (
	"context"
	"fmt"

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

	enabled := true
	deviceID := device.ID
	points, err := s.pointSvc.List(ctx, point.ListFilter{
		DeviceID: &deviceID,
		Enabled:  &enabled,
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
		s.registerPointMeta(pt.ID, pointMeta{
			DeviceID: pt.DeviceID,
			Address:  pt.Address,
		})
		s.scheduler.AddPoint(pt)
	}

	s.markDeviceProjectionAlignedFromWorkspace(ctx, device)
	return nil
}

// RemoveDevice 將設備及其關聯點位從 Runtime 與 scheduler 移除。
func (s *Service) RemoveDevice(deviceID string) {
	if s == nil || deviceID == "" {
		return
	}

	s.removeDevicePointMeta(deviceID)

	if s.scheduler != nil {
		s.scheduler.RemoveDevice(deviceID)
	}
	s.clearDeviceProjection(deviceID)
}

// removeDevicePointMeta 清除指定設備的 point metadata，避免停用後仍保留舊事件標記。
func (s *Service) removeDevicePointMeta(deviceID string) {
	s.pointMetaMu.Lock()
	defer s.pointMetaMu.Unlock()

	for pointID, meta := range s.pointMetaIndex {
		if meta.DeviceID == deviceID {
			delete(s.pointMetaIndex, pointID)
		}
	}
}

// ensurePointPollingGroup 為缺少 polling group 的啟用點位補上一個可用的預設群組。
func (s *Service) ensurePointPollingGroup(ctx context.Context, pointRecord *schema.Point) (*schema.Point, error) {
	if pointRecord == nil || pointRecord.PollingGroupID != nil || s.pointSvc == nil {
		return pointRecord, nil
	}

	group, err := s.pointSvc.EnsureDefaultPollingGroup(ctx)
	if err != nil {
		return nil, fmt.Errorf("確保預設輪詢群組失敗: %w", err)
	}
	if group == nil {
		return pointRecord, nil
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
