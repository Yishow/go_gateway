package device

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
)

func validateAcquisitionSettingsPatch(req UpdateDeviceRequest) error {
	settings := schema.DeviceAcquisitionSettings{StaleThresholdMs: req.StaleThresholdMs}
	if req.FrozenThresholdMs != nil {
		settings.FrozenThresholdMs = *req.FrozenThresholdMs
	}
	return validateAcquisitionSettings(settings)
}

// UpdateAcquisitionSettingsPatch updates only acquisition fields supplied by the caller.
func (r *SQLRepository) UpdateAcquisitionSettingsPatch(ctx context.Context, id string, req UpdateDeviceRequest) error {
	return updateAcquisitionSettingsPatch(ctx, r.db, id, req)
}

// UpdateAcquisitionSettingsPatchInTx requires the same config DB as the point transaction.
func (r *SQLRepository) UpdateAcquisitionSettingsPatchInTx(ctx context.Context, db *sql.DB, tx *sql.Tx, id string, req UpdateDeviceRequest) error {
	if r == nil || r.db == nil || r.db != db || tx == nil {
		return fmt.Errorf("設備採集設定需要共用設定資料庫")
	}
	return updateAcquisitionSettingsPatch(ctx, tx, id, req)
}

func updateAcquisitionSettingsPatch(ctx context.Context, db interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}, id string, req UpdateDeviceRequest) error {
	if err := validateAcquisitionSettingsPatch(req); err != nil {
		return err
	}
	var fields []string
	var args []any
	if req.StaleThresholdMs != nil {
		fields = append(fields, "stale_threshold_ms=?")
		args = append(args, *req.StaleThresholdMs)
	}
	if req.FrozenThresholdMs != nil {
		fields = append(fields, "frozen_threshold_ms=?")
		args = append(args, *req.FrozenThresholdMs)
	}
	if req.AcquisitionAutoInterval != nil {
		fields = append(fields, "acquisition_auto_interval=?")
		args = append(args, *req.AcquisitionAutoInterval)
	}
	if len(fields) == 0 {
		return fmt.Errorf("設備採集設定未提供")
	}
	fields = append(fields, "updated_at=?")
	args = append(args, time.Now(), id)
	result, err := db.ExecContext(ctx, "UPDATE devices SET "+strings.Join(fields, ", ")+" WHERE id=?", args...)
	if err != nil {
		return fmt.Errorf("failed to update device acquisition settings: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to confirm device acquisition settings update: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("device not found: %s", id)
	}
	return nil
}

// UpdateAcquisitionSettingsPatch mirrors the SQL partial update in memory.
func (r *MemoryRepository) UpdateAcquisitionSettingsPatch(_ context.Context, id string, req UpdateDeviceRequest) error {
	if err := validateAcquisitionSettingsPatch(req); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	dev, ok := r.devices[id]
	if !ok {
		return fmt.Errorf("設備不存在: %s", id)
	}
	applyAcquisitionSettingsPatch(dev, req)
	return nil
}

func applyAcquisitionSettingsPatch(dev *schema.Device, req UpdateDeviceRequest) {
	if req.StaleThresholdMs != nil {
		dev.StaleThresholdMs = cloneThreshold(req.StaleThresholdMs)
	}
	if req.FrozenThresholdMs != nil {
		dev.FrozenThresholdMs = *req.FrozenThresholdMs
	}
	if req.AcquisitionAutoInterval != nil {
		dev.AcquisitionAutoInterval = *req.AcquisitionAutoInterval
	}
	dev.UpdatedAt = time.Now()
}
