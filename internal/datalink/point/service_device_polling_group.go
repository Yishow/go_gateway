package point

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	common "go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

type atomicDeviceAcquisitionRepository interface {
	setDeviceAcquisitionInterval(context.Context, *Service, string, int) (*schema.PollingGroup, error)
}

const (
	minPollingGroupIntervalMs    = 100
	maxPollingGroupIntervalMs    = 3600000
	devicePollingGroupNamePrefix = "device:"
)

// GetPollingGroupByID 依群組 ID 查詢輪詢群組。
func (s *Service) GetPollingGroupByID(ctx context.Context, id string) (*schema.PollingGroup, error) {
	id = strings.TrimSpace(id)
	if s == nil || s.groupRepo == nil || id == "" {
		return nil, nil
	}
	return s.groupRepo.GetByID(ctx, id)
}

// GetPollingGroupByOwnerDeviceID 依擁有設備 ID 查詢專屬輪詢群組；若不存在回傳 (nil, nil)。
func (s *Service) GetPollingGroupByOwnerDeviceID(ctx context.Context, deviceID string) (*schema.PollingGroup, error) {
	deviceID = strings.TrimSpace(deviceID)
	if s == nil || s.groupRepo == nil || deviceID == "" {
		return nil, nil
	}
	groups, err := s.groupRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("列出輪詢群組失敗: %w", err)
	}
	for _, group := range groups {
		if group != nil && group.OwnerDeviceID != nil && *group.OwnerDeviceID == deviceID {
			return group, nil
		}
	}
	return nil, nil
}

// EnsureDevicePollingGroup 優先回傳該設備專屬輪詢群組，若無則退回 EnsureDefaultPollingGroup。
func (s *Service) EnsureDevicePollingGroup(ctx context.Context, deviceID string) (*schema.PollingGroup, error) {
	if s == nil || s.groupRepo == nil {
		return nil, nil
	}
	group, err := s.GetPollingGroupByOwnerDeviceID(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	if group != nil {
		if !group.Enabled {
			group, err = s.UpdatePollingGroup(ctx, group.ID, UpdatePollingGroupRequest{
				Enabled: common.Ptr(true),
			})
			if err != nil {
				return nil, err
			}
		}
		s.syncAutoPollingGroup(group)
		return group, nil
	}
	return s.EnsureDefaultPollingGroup(ctx)
}

// SetDeviceAcquisitionInterval 建立或更新設備專屬輪詢群組，並將該設備全部點位改掛到該群組且立即同步排程器。
func (s *Service) SetDeviceAcquisitionInterval(ctx context.Context, deviceID string, intervalMs int) (*schema.PollingGroup, error) {
	return s.setDeviceAcquisitionInterval(ctx, deviceID, intervalMs, nil)
}

// SetDeviceAcquisitionIntervalWithSettings commits device settings in the same
// config database transaction as the owned group and point assignments.
func (s *Service) SetDeviceAcquisitionIntervalWithSettings(ctx context.Context, deviceID string, intervalMs int, update func(context.Context, *sql.DB, *sql.Tx) error) (*schema.PollingGroup, error) {
	if update == nil {
		return nil, fmt.Errorf("設備採集設定交易更新未提供")
	}
	return s.setDeviceAcquisitionInterval(ctx, deviceID, intervalMs, update)
}

func (s *Service) setDeviceAcquisitionInterval(ctx context.Context, deviceID string, intervalMs int, update func(context.Context, *sql.DB, *sql.Tx) error) (*schema.PollingGroup, error) {
	deviceID = strings.TrimSpace(deviceID)
	if deviceID == "" {
		return nil, fmt.Errorf("device_id 不能為空")
	}
	if intervalMs < minPollingGroupIntervalMs || intervalMs > maxPollingGroupIntervalMs {
		return nil, fmt.Errorf("acquisition_interval_ms 必須介於 %d 到 %d 毫秒之間", minPollingGroupIntervalMs, maxPollingGroupIntervalMs)
	}
	if s == nil || s.groupRepo == nil {
		return nil, fmt.Errorf("輪詢群組儲存庫未初始化")
	}
	repo, ok := s.repo.(atomicDeviceAcquisitionRepository)
	if !ok {
		return nil, fmt.Errorf("點位儲存庫不支援原子設備採集頻率設定")
	}

	finish := s.beginAcquisitionMutation(ctx, deviceID)
	defer finish()
	var group *schema.PollingGroup
	var err error
	if update == nil {
		group, err = repo.setDeviceAcquisitionInterval(ctx, s, deviceID, intervalMs)
	} else if atomicRepo, ok := s.repo.(interface {
		setDeviceAcquisitionIntervalWithSettings(context.Context, *Service, string, int, func(context.Context, *sql.DB, *sql.Tx) error) (*schema.PollingGroup, error)
	}); ok {
		group, err = atomicRepo.setDeviceAcquisitionIntervalWithSettings(ctx, s, deviceID, intervalMs, update)
	} else {
		return nil, fmt.Errorf("點位儲存庫不支援原子設備採集設定")
	}
	if err != nil {
		return nil, err
	}
	s.syncAutoPollingGroup(group)
	return group, nil
}

func (s *Service) saveMemoryDeviceAcquisitionGroup(ctx context.Context, deviceID string, intervalMs int) (*schema.PollingGroup, error) {
	group, err := s.GetPollingGroupByOwnerDeviceID(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		name, nameErr := s.nextDevicePollingGroupName(ctx, deviceID)
		if nameErr != nil {
			return nil, nameErr
		}
		group, err = s.CreatePollingGroup(ctx, CreatePollingGroupRequest{
			Name:          name,
			Description:   "Device-owned polling group",
			IntervalMs:    intervalMs,
			Priority:      defaultPollingGroupPriority,
			OwnerDeviceID: common.Ptr(deviceID),
		})
		if err != nil {
			return nil, err
		}
	} else {
		group, err = s.UpdatePollingGroup(ctx, group.ID, UpdatePollingGroupRequest{
			IntervalMs: common.Ptr(intervalMs),
			Enabled:    common.Ptr(true),
		})
		if err != nil {
			return nil, err
		}
	}

	return group, nil
}

func (s *Service) nextDevicePollingGroupName(ctx context.Context, deviceID string) (string, error) {
	groups, err := s.groupRepo.List(ctx)
	if err != nil {
		return "", fmt.Errorf("列出輪詢群組失敗: %w", err)
	}
	taken := make(map[string]struct{}, len(groups))
	for _, g := range groups {
		if g != nil {
			taken[g.Name] = struct{}{}
		}
	}
	base := devicePollingGroupNamePrefix + deviceID
	for suffix := 1; ; suffix++ {
		candidate := base
		if suffix > 1 {
			candidate = fmt.Sprintf("%s#%d", base, suffix)
		}
		if _, exists := taken[candidate]; !exists {
			return candidate, nil
		}
	}
}
