package pollinggroup

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/sqltx"
)

// AcquisitionMutationLifecycle quarantines all affected source devices without
// waiting for runtime locks. The caller owns a SQLite writer reservation, so
// implementations must not perform SQL/provider I/O or block on runtime locks.
// Finish releases only this attempt's tokens, never restores old gap authority.
type AcquisitionMutationLifecycle interface {
	TryBeginAcquisitionConfigurationMutations(context.Context, []string) (func(), error)
}

// SetAcquisitionMutationLifecycle is startup-only wiring. Keeping this seam at
// the shared SQL repository covers both existing group deletion service paths.
func (r *SQLRepository) SetAcquisitionMutationLifecycle(lifecycle AcquisitionMutationLifecycle) {
	r.acquisitionMutation = lifecycle
}

func (r *SQLRepository) deleteWithAcquisitionMutation(ctx context.Context, id string) (resultErr error) {
	owned, err := sqltx.Begin(ctx, r.db, nil)
	if err != nil {
		return fmt.Errorf("begin polling group source mutation: %w", err)
	}
	var finish func()
	defer func() { //nolint:contextcheck // sqltx.Rollback intentionally uses independent bounded cleanup after caller cancellation.
		// Finish can briefly take scheduler.Lock. Release/clean the SQL writer
		// reservation first on every path, including cancellation and failed commit.
		if err := owned.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			resultErr = errors.Join(resultErr, fmt.Errorf("cleanup polling group source mutation: %w", err))
		}
		if finish != nil {
			finish()
		}
	}()
	var version string
	if err := owned.Tx.QueryRowContext(ctx, `SELECT sqlite_version()`).Scan(&version); err != nil || version == "" {
		return errors.New("polling group source mutation requires SQLite")
	}
	// The canonical internal database uses SQLite. Reserving its writer BEFORE
	// enumerating FK targets closes concurrent Assign/Delete membership races.
	if _, err := owned.Tx.ExecContext(ctx, `UPDATE system_settings SET value=value WHERE 0`); err != nil {
		return fmt.Errorf("reserve polling group source mutation: %w", err)
	}
	devices, err := acquisitionGroupDeviceIDs(ctx, owned.Tx, id)
	if err != nil {
		return fmt.Errorf("read polling group source devices: %w", err)
	}
	if len(devices) > 0 {
		finish, err = r.acquisitionMutation.TryBeginAcquisitionConfigurationMutations(ctx, devices)
		if err != nil {
			return err
		}
		if finish == nil {
			return errors.New("polling group source mutation has no completion token")
		}
	}
	result, err := owned.Tx.ExecContext(ctx, `DELETE FROM polling_groups WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("刪除輪詢群組失敗: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得刪除影響列數失敗: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("輪詢群組不存在: %s", id)
	}
	if err := owned.Commit(); err != nil { //nolint:contextcheck // sqltx.Commit owns independent bounded cleanup; its terminal API has no context parameter.
		return fmt.Errorf("commit polling group source mutation: %w", err)
	}
	return nil
}

func acquisitionGroupDeviceIDs(ctx context.Context, tx *sql.Tx, id string) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `SELECT DISTINCT device_id FROM points WHERE polling_group_id=? ORDER BY device_id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	devices := make([]string, 0)
	for rows.Next() {
		var deviceID string
		if err := rows.Scan(&deviceID); err != nil {
			return nil, err
		}
		if strings.TrimSpace(deviceID) == "" {
			return nil, errors.New("polling group source device identity is empty")
		}
		devices = append(devices, deviceID)
	}
	return devices, rows.Err()
}
