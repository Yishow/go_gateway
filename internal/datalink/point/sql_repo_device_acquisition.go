package point

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	common "go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sqltx"
)

func (r *SQLRepository) setDeviceAcquisitionInterval(ctx context.Context, s *Service, deviceID string, intervalMs int) (group *schema.PollingGroup, resultErr error) {
	return r.setDeviceAcquisitionIntervalWithSettings(ctx, s, deviceID, intervalMs, nil)
}

func (r *SQLRepository) setDeviceAcquisitionIntervalWithSettings(ctx context.Context, s *Service, deviceID string, intervalMs int, update func(context.Context, *sql.DB, *sql.Tx) error) (group *schema.PollingGroup, resultErr error) {
	groups, ok := s.groupRepo.(interface{ UsesConfigDB(*sql.DB) bool })
	if r == nil || r.db == nil || !ok || !groups.UsesConfigDB(r.db) {
		return nil, fmt.Errorf("設備採集頻率需要共用設定資料庫")
	}
	owned, err := sqltx.Begin(ctx, r.db, nil)
	if err != nil {
		return nil, fmt.Errorf("開始設備採集頻率交易失敗: %w", err)
	}
	defer func() {
		if rollbackErr := owned.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone { //nolint:errorlint // A joined cleanup error must not be mistaken for a bare ErrTxDone.
			group = nil
			resultErr = errors.Join(resultErr, fmt.Errorf("回復設備採集頻率交易失敗: %w", rollbackErr))
		}
	}()

	group, err = saveSQLDeviceAcquisitionGroup(ctx, owned.Tx, deviceID, intervalMs)
	if err != nil {
		return nil, err
	}
	_, err = owned.Tx.ExecContext(ctx, `UPDATE points SET polling_group_id = ?, revision_id = lower(hex(randomblob(16))), updated_at = ?
		WHERE device_id = ? AND (polling_group_id IS NULL OR polling_group_id != ?)`, group.ID, time.Now(), deviceID, group.ID)
	if err != nil {
		return nil, fmt.Errorf("更新設備點位輪詢群組失敗: %w", err)
	}
	if update != nil {
		if err := update(ctx, r.db, owned.Tx); err != nil {
			return nil, fmt.Errorf("更新設備採集門檻失敗: %w", err)
		}
	}
	if err := owned.Commit(); err != nil {
		return nil, fmt.Errorf("提交設備採集頻率交易失敗: %w", err)
	}
	return group, nil
}

func saveSQLDeviceAcquisitionGroup(ctx context.Context, tx *sql.Tx, deviceID string, intervalMs int) (*schema.PollingGroup, error) {
	group := &schema.PollingGroup{OwnerDeviceID: common.Ptr(deviceID)}
	var description sql.NullString
	var createdAt, updatedAt string
	err := tx.QueryRowContext(ctx, `SELECT id, name, description, interval_ms, priority, enabled, created_at, updated_at
		FROM polling_groups WHERE owner_device_id = ?`, deviceID).
		Scan(&group.ID, &group.Name, &description, &group.IntervalMs, &group.Priority, &group.Enabled, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return createSQLDeviceAcquisitionGroup(ctx, tx, deviceID, intervalMs)
	}
	if err != nil {
		return nil, fmt.Errorf("取得設備專屬輪詢群組失敗: %w", err)
	}
	group.Description = description.String
	group.CreatedAt, err = common.ParseTimeString(createdAt)
	if err != nil {
		return nil, fmt.Errorf("解析群組建立時間失敗: %w", err)
	}
	group.UpdatedAt = time.Now()
	group.IntervalMs = intervalMs
	group.Enabled = true
	if _, err := tx.ExecContext(ctx, `UPDATE polling_groups SET interval_ms = ?, enabled = 1, updated_at = ? WHERE id = ?`, intervalMs, group.UpdatedAt, group.ID); err != nil {
		return nil, fmt.Errorf("更新設備專屬輪詢群組失敗: %w", err)
	}
	return group, nil
}

func createSQLDeviceAcquisitionGroup(ctx context.Context, tx *sql.Tx, deviceID string, intervalMs int) (*schema.PollingGroup, error) {
	base := devicePollingGroupNamePrefix + deviceID
	name := base
	for suffix := 1; ; suffix++ {
		if suffix > 1 {
			name = fmt.Sprintf("%s#%d", base, suffix)
		}
		var exists bool
		if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM polling_groups WHERE name = ?)`, name).Scan(&exists); err != nil {
			return nil, fmt.Errorf("檢查輪詢群組名稱失敗: %w", err)
		}
		if !exists {
			break
		}
	}
	id, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("建立輪詢群組 ID 失敗: %w", err)
	}
	now := time.Now()
	group := &schema.PollingGroup{ID: id, Name: name, Description: "Device-owned polling group", IntervalMs: intervalMs, Priority: defaultPollingGroupPriority, Enabled: true, OwnerDeviceID: common.Ptr(deviceID), CreatedAt: now, UpdatedAt: now}
	if _, err := tx.ExecContext(ctx, `INSERT INTO polling_groups (id, name, description, interval_ms, priority, enabled, owner_device_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, group.ID, group.Name, group.Description, group.IntervalMs, group.Priority, group.Enabled, group.OwnerDeviceID, group.CreatedAt, group.UpdatedAt); err != nil {
		return nil, fmt.Errorf("建立設備專屬輪詢群組失敗: %w", err)
	}
	return group, nil
}
