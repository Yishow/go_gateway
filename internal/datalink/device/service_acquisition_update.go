package device

import (
	"context"
	"database/sql"
	"fmt"

	common "go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

// ApplyCurrentRuntimeDevice serializes runtime projections through this device
// service and reads the current saved device immediately before applying it.
func (s *Service) ApplyCurrentRuntimeDevice(ctx context.Context, id string, apply func(*schema.Device) error) (*schema.Device, error) {
	if apply == nil {
		return nil, fmt.Errorf("runtime 設備套用未提供")
	}
	s.runtimeApplyMu.Lock()
	defer s.runtimeApplyMu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	dev, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := apply(dev); err != nil {
		return dev, err
	}
	return dev, nil
}

func (s *Service) updateDeviceAcquisition(ctx context.Context, dev *schema.Device, req UpdateDeviceRequest) error {
	hasSettings := req.StaleThresholdMs != nil || req.FrozenThresholdMs != nil || req.AcquisitionAutoInterval != nil
	if req.AcquisitionIntervalMs == nil {
		if hasSettings {
			repo, ok := s.repo.(interface {
				UpdateAcquisitionSettingsPatch(context.Context, string, UpdateDeviceRequest) error
			})
			if !ok {
				return fmt.Errorf("設備儲存庫不支援採集門檻設定")
			}
			if err := repo.UpdateAcquisitionSettingsPatch(ctx, dev.ID, req); err != nil {
				return fmt.Errorf("更新設備採集門檻失敗: %w", err)
			}
		}
		if !hasSettings {
			s.hydrateDeviceAcquisitionInterval(ctx, dev)
		}
		return nil
	}
	if s.acquisitionCoordinator == nil {
		return fmt.Errorf("設備採集頻率協調器未設定，無法儲存 acquisition_interval_ms")
	}
	var group *schema.PollingGroup
	var err error
	if !hasSettings {
		group, err = s.acquisitionCoordinator.SetDeviceAcquisitionInterval(ctx, dev.ID, *req.AcquisitionIntervalMs)
	} else {
		switch repo := s.repo.(type) {
		case interface {
			UpdateAcquisitionSettingsPatchInTx(context.Context, *sql.DB, *sql.Tx, string, UpdateDeviceRequest) error
		}:
			coord, supported := s.acquisitionCoordinator.(interface {
				SetDeviceAcquisitionIntervalWithSettings(context.Context, string, int, func(context.Context, *sql.DB, *sql.Tx) error) (*schema.PollingGroup, error)
			})
			if !supported {
				return fmt.Errorf("設備採集頻率協調器不支援原子設定")
			}
			group, err = coord.SetDeviceAcquisitionIntervalWithSettings(ctx, dev.ID, *req.AcquisitionIntervalMs, func(ctx context.Context, db *sql.DB, tx *sql.Tx) error {
				return repo.UpdateAcquisitionSettingsPatchInTx(ctx, db, tx, dev.ID, req)
			})
		case interface {
			withAcquisitionSettings(context.Context, string, UpdateDeviceRequest, func() error) error
		}:
			err = repo.withAcquisitionSettings(ctx, dev.ID, req, func() error {
				group, err = s.acquisitionCoordinator.SetDeviceAcquisitionInterval(ctx, dev.ID, *req.AcquisitionIntervalMs)
				return err
			})
		default:
			return fmt.Errorf("設備儲存庫不支援原子採集設定")
		}
	}
	if err != nil {
		return fmt.Errorf("設定設備採集頻率失敗: %w", err)
	}
	dev.AcquisitionIntervalMs = common.Ptr(group.IntervalMs)
	return nil
}
