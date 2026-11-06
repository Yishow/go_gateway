package device

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-gateway/internal/datalink/schema"
)

// validateAcquisitionSettings rejects values the anomaly detector cannot use.
// A nil stale threshold is the derived default (interval x5, at least 10 s);
// a zero frozen threshold disables frozen detection. Invalid values are never
// replaced by defaults.
func validateAcquisitionSettings(settings schema.DeviceAcquisitionSettings) error {
	if settings.StaleThresholdMs != nil && *settings.StaleThresholdMs <= 0 {
		return validationError("stale_threshold_ms 必須大於 0，或留空使用預設值")
	}
	if settings.FrozenThresholdMs < 0 {
		return validationError("frozen_threshold_ms 不能小於 0（0 表示停用）")
	}
	return nil
}

func nullableThreshold(value sql.NullInt64) *int {
	if !value.Valid {
		return nil
	}
	v := int(value.Int64)
	return &v
}

// UpdateAcquisitionSettings writes only the device anomaly/auto-interval
// settings. Whole-row Update does not touch these columns, so frequent
// read-modify-write paths (collection stats, readiness) cannot overwrite them.
func (r *SQLRepository) UpdateAcquisitionSettings(ctx context.Context, id string, settings schema.DeviceAcquisitionSettings) error {
	if err := validateAcquisitionSettings(settings); err != nil {
		return err
	}
	result, err := r.db.ExecContext(ctx, `UPDATE devices SET stale_threshold_ms=?, frozen_threshold_ms=?, acquisition_auto_interval=?, updated_at=? WHERE id=?`,
		settings.StaleThresholdMs, settings.FrozenThresholdMs, settings.AcquisitionAutoInterval, time.Now(), id)
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

// withAcquisitionSettings holds the memory repository's device lock until the
// point operation succeeds, so failures cannot leave only settings applied.
func (r *MemoryRepository) withAcquisitionSettings(ctx context.Context, id string, req UpdateDeviceRequest, applyInterval func() error) error {
	if err := validateAcquisitionSettingsPatch(req); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	device, ok := r.devices[id]
	if !ok {
		return fmt.Errorf("設備不存在: %s", id)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := applyInterval(); err != nil {
		return err
	}
	applyAcquisitionSettingsPatch(device, req)
	return nil
}

// UpdateAcquisitionSettings mirrors the SQL repository for service tests.
func (r *MemoryRepository) UpdateAcquisitionSettings(_ context.Context, id string, settings schema.DeviceAcquisitionSettings) error {
	if err := validateAcquisitionSettings(settings); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	device, ok := r.devices[id]
	if !ok {
		return fmt.Errorf("設備不存在: %s", id)
	}
	device.StaleThresholdMs = cloneThreshold(settings.StaleThresholdMs)
	device.FrozenThresholdMs = settings.FrozenThresholdMs
	device.AcquisitionAutoInterval = settings.AcquisitionAutoInterval
	device.UpdatedAt = time.Now()
	return nil
}

func cloneThreshold(value *int) *int {
	if value == nil {
		return nil
	}
	v := *value
	return &v
}
