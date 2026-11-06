package datalink

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema/migrations"
)

const (
	deviceAcquisitionMigration       = "053_device_acquisition_sqlite.up.sql"
	deviceOwnedGroupNamePrefix       = "device:"
	deviceOwnedGroupDescription      = "設備專屬輪詢群組"
	deviceOwnedGroupFallbackInterval = 1000
	deviceOwnedGroupFallbackPriority = 50
)

// ensureSQLiteLatestExpansions keeps migrator.go within its line budget while
// chaining the newest expansion after its predecessor.
func ensureSQLiteLatestExpansions(db *sql.DB) error {
	if err := ensureSQLiteMixedSourceIndexes(db); err != nil {
		return err
	}
	if err := ensureSQLiteDeviceAcquisition(db); err != nil {
		return err
	}
	return ensureSQLiteAcquisitionAnomalies(db)
}

// ensureSQLiteDeviceAcquisition adds device-owned polling groups and device
// anomaly thresholds, then gives every device without an owned group the
// shortest interval of its points' enabled groups. Schema expansion and
// backfill share one transaction; re-running only touches devices without a group.
func ensureSQLiteDeviceAcquisition(db *sql.DB) error {
	return runSQLiteConfigurationExpansion(db, func(ctx context.Context, tx *sql.Tx) error {
		const devicesTable = "devices"
		raw, err := migrations.FS.ReadFile(deviceAcquisitionMigration)
		if err != nil {
			return err
		}
		statements := strings.Split(strings.TrimSpace(string(raw)), ";")
		columns := []struct{ table, column, statement string }{
			{"polling_groups", "owner_device_id", statements[0]},
			{devicesTable, "stale_threshold_ms", statements[2]},
			{devicesTable, "frozen_threshold_ms", statements[3]},
			{devicesTable, "acquisition_auto_interval", statements[4]},
		}
		for _, c := range columns {
			if err := expandSQLiteDeviceAcquisitionColumn(ctx, tx, c.table, c.column, strings.TrimSpace(c.statement)); err != nil {
				return err
			}
		}
		if err := expandSQLiteOwnerDeviceIndex(ctx, tx, strings.TrimSpace(statements[1])); err != nil {
			return err
		}
		if err := backfillDeviceOwnedPollingGroups(ctx, tx); err != nil {
			return err
		}
		return checkSQLiteMixedForeignKeys(ctx, tx)
	})
}

func expandSQLiteDeviceAcquisitionColumn(ctx context.Context, tx *sql.Tx, table, column, statement string) error {
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM pragma_table_info(?) WHERE name=?`, table, column).Scan(&count); err != nil {
		return fmt.Errorf("inspect %s.%s: %w", table, column, err)
	}
	if count == 0 {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("expand %s.%s: %w", table, column, err)
		}
	}
	var ddl string
	if err := tx.QueryRowContext(ctx, `SELECT sql FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&ddl); err != nil {
		return fmt.Errorf("inspect %s schema: %w", table, err)
	}
	declaration := strings.TrimPrefix(statement, "ALTER TABLE "+table+" ADD COLUMN ")
	if declaration == statement || strings.Count(ddl, declaration) != 1 {
		return fmt.Errorf("unknown %s.%s column; migration refused", table, column)
	}
	return nil
}

func expandSQLiteOwnerDeviceIndex(ctx context.Context, tx *sql.Tx, expected string) error {
	const name = "idx_polling_groups_owner_device"
	var kind, ddl string
	err := tx.QueryRowContext(ctx, `SELECT type, sql FROM sqlite_master WHERE name=?`, name).Scan(&kind, &ddl)
	if errors.Is(err, sql.ErrNoRows) {
		if _, err := tx.ExecContext(ctx, expected); err != nil {
			return fmt.Errorf("create %s: %w", name, err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect %s: %w", name, err)
	}
	if kind != "index" || canonicalMixedDeliveryDDL(ddl) != canonicalMixedDeliveryDDL(expected) {
		return fmt.Errorf("unknown %s schema; migration refused", name)
	}
	return nil
}

func backfillDeviceOwnedPollingGroups(ctx context.Context, tx *sql.Tx) error {
	// Same ordering as point.Service.EnsureDefaultPollingGroup, so a device with
	// no grouped points gets the cadence new points would have received.
	fallbackInterval, fallbackPriority := deviceOwnedGroupFallbackInterval, deviceOwnedGroupFallbackPriority
	err := tx.QueryRowContext(ctx, `SELECT interval_ms, priority FROM polling_groups WHERE owner_device_id IS NULL ORDER BY enabled DESC, priority ASC, interval_ms ASC, id ASC LIMIT 1`).Scan(&fallbackInterval, &fallbackPriority)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("read default polling group: %w", err)
	}
	deviceIDs, err := devicesWithoutOwnedPollingGroup(ctx, tx)
	if err != nil {
		return err
	}
	for _, deviceID := range deviceIDs {
		if err := backfillDeviceOwnedPollingGroup(ctx, tx, deviceID, fallbackInterval, fallbackPriority); err != nil {
			return err
		}
	}
	return nil
}

func devicesWithoutOwnedPollingGroup(ctx context.Context, tx *sql.Tx) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `SELECT d.id FROM devices d WHERE NOT EXISTS (SELECT 1 FROM polling_groups g WHERE g.owner_device_id = d.id) ORDER BY d.id`)
	if err != nil {
		return nil, fmt.Errorf("list devices without owned polling group: %w", err)
	}
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, errors.Join(fmt.Errorf("scan device id: %w", err), rows.Close())
		}
		ids = append(ids, id)
	}
	if err := errors.Join(rows.Err(), rows.Close()); err != nil {
		return nil, fmt.Errorf("list devices without owned polling group: %w", err)
	}
	return ids, nil
}

func backfillDeviceOwnedPollingGroup(ctx context.Context, tx *sql.Tx, deviceID string, fallbackInterval, fallbackPriority int) error {
	// Only enabled groups are polled today; points in disabled groups or with no
	// group keep their state so the upgrade never starts or changes polling.
	var distinct int
	var interval, priority sql.NullInt64
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(DISTINCT g.interval_ms), MIN(g.interval_ms), MIN(g.priority) FROM points p JOIN polling_groups g ON g.id = p.polling_group_id WHERE p.device_id = ? AND g.enabled = 1`, deviceID).Scan(&distinct, &interval, &priority); err != nil {
		return fmt.Errorf("read polling intervals for device %s: %w", deviceID, err)
	}
	if distinct > 0 {
		fallbackInterval, fallbackPriority = int(interval.Int64), int(priority.Int64)
	}
	name, err := nextDeviceOwnedPollingGroupName(ctx, tx, deviceID)
	if err != nil {
		return err
	}
	now := time.Now()
	if _, err := tx.ExecContext(ctx, `INSERT INTO polling_groups (name, description, interval_ms, priority, enabled, owner_device_id, created_at, updated_at) VALUES (?, ?, ?, ?, 1, ?, ?, ?)`, name, deviceOwnedGroupDescription, fallbackInterval, fallbackPriority, deviceID, now, now); err != nil {
		return fmt.Errorf("create device-owned polling group for device %s: %w", deviceID, err)
	}
	var groupID string
	if err := tx.QueryRowContext(ctx, `SELECT id FROM polling_groups WHERE owner_device_id = ?`, deviceID).Scan(&groupID); err != nil {
		return fmt.Errorf("read device-owned polling group for device %s: %w", deviceID, err)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE points SET polling_group_id = ?, updated_at = ? WHERE device_id = ? AND polling_group_id IN (SELECT id FROM polling_groups WHERE enabled = 1 AND id != ?)`, groupID, now, deviceID, groupID); err != nil {
		return fmt.Errorf("assign points to device-owned polling group for device %s: %w", deviceID, err)
	}
	return nil
}

func nextDeviceOwnedPollingGroupName(ctx context.Context, tx *sql.Tx, deviceID string) (string, error) {
	base := deviceOwnedGroupNamePrefix + deviceID
	for suffix := 1; ; suffix++ {
		candidate := base
		if suffix > 1 {
			candidate = fmt.Sprintf("%s#%d", base, suffix)
		}
		var taken int
		if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM polling_groups WHERE name = ?`, candidate).Scan(&taken); err != nil {
			return "", fmt.Errorf("check device-owned polling group name: %w", err)
		}
		if taken == 0 {
			return candidate, nil
		}
	}
}
