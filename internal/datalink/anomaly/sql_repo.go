package anomaly

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// SQLRepository stores events in the SQLite configuration database table
// acquisition_anomalies. Times are stored as Unix milliseconds (UTC).
type SQLRepository struct {
	db        *sql.DB
	retention int
}

// SQLOption configures a SQLRepository.
type SQLOption func(*SQLRepository)

// WithRetention overrides the per-device retention target. Non-positive values keep
// the default, so a misconfiguration can never disable retention silently.
func WithRetention(n int) SQLOption {
	return func(r *SQLRepository) {
		if n > 0 {
			r.retention = n
		}
	}
}

// NewSQLRepository creates a repository on the migrated configuration database.
func NewSQLRepository(db *sql.DB, opts ...SQLOption) *SQLRepository {
	r := &SQLRepository{db: db, retention: DefaultRetention}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

const eventColumns = `id, device_id, point_id, kind, started_at, ended_at, last_value, detail`

// Open inserts an active event and applies retention in the same transaction.
func (r *SQLRepository) Open(ctx context.Context, ev Event) (_ Event, err error) {
	if err := validateEvent(ev); err != nil {
		return Event{}, err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Event{}, fmt.Errorf("open acquisition anomaly: begin: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = errors.Join(err, fmt.Errorf("open acquisition anomaly: rollback: %w", rollbackErr))
		}
	}()
	res, err := tx.ExecContext(ctx, `INSERT INTO acquisition_anomalies (device_id, point_id, kind, started_at, last_value, detail) VALUES (?, ?, ?, ?, ?, ?)`,
		ev.DeviceID, nullString(ev.PointID), string(ev.Kind), ev.StartedAt.UnixMilli(), ev.LastValue, ev.Detail)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return Event{}, fmt.Errorf("open acquisition anomaly %s/%s/%s: %w", ev.DeviceID, ev.PointID, ev.Kind, ErrAlreadyActive)
		}
		return Event{}, fmt.Errorf("open acquisition anomaly: insert: %w", err)
	}
	if ev.ID, err = res.LastInsertId(); err != nil {
		return Event{}, fmt.Errorf("open acquisition anomaly: read id: %w", err)
	}
	if err := r.prune(ctx, tx, ev.DeviceID); err != nil {
		return Event{}, err
	}
	if err := tx.Commit(); err != nil {
		return Event{}, fmt.Errorf("open acquisition anomaly: commit: %w", err)
	}
	ev.StartedAt = time.UnixMilli(ev.StartedAt.UnixMilli()).UTC()
	return ev, nil
}

func (r *SQLRepository) prune(ctx context.Context, tx *sql.Tx, deviceID string) error {
	var total int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM acquisition_anomalies WHERE device_id = ?`, deviceID).Scan(&total); err != nil {
		return fmt.Errorf("acquisition anomaly retention: count: %w", err)
	}
	if excess := total - r.retention; excess > 0 {
		if _, err := tx.ExecContext(ctx, `DELETE FROM acquisition_anomalies WHERE id IN (SELECT id FROM acquisition_anomalies WHERE device_id = ? AND ended_at IS NOT NULL ORDER BY started_at ASC, id ASC LIMIT ?)`, deviceID, excess); err != nil {
			return fmt.Errorf("acquisition anomaly retention: delete: %w", err)
		}
	}
	return nil
}

// End records the end of an active event. A nil lastValue keeps the stored
// value; a non-empty endDetail is appended to the start cause, never replacing it.
func (r *SQLRepository) End(ctx context.Context, id int64, endedAt time.Time, lastValue *string, endDetail string) (err error) {
	if endedAt.IsZero() {
		return errors.Join(ErrInvalidEvent, errors.New("end time is required"))
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("end acquisition anomaly %d: begin: %w", id, err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = errors.Join(err, fmt.Errorf("end acquisition anomaly %d: rollback: %w", id, rollbackErr))
		}
	}()
	var deviceID string
	if err := tx.QueryRowContext(ctx, `SELECT device_id FROM acquisition_anomalies WHERE id = ? AND ended_at IS NULL`, id).Scan(&deviceID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("end acquisition anomaly %d: %w", id, ErrNotActive)
		}
		return fmt.Errorf("end acquisition anomaly %d: read device: %w", id, err)
	}
	// An observation stamped just before the computed start (a read in flight
	// while the threshold expired) ends the event at its start, not before it.
	res, err := tx.ExecContext(ctx, `UPDATE acquisition_anomalies SET ended_at = MAX(?, started_at), last_value = COALESCE(?, last_value), detail = CASE WHEN ? = '' THEN detail WHEN detail = '' THEN ? ELSE detail || '; ' || ? END WHERE id = ? AND ended_at IS NULL`,
		endedAt.UnixMilli(), lastValue, endDetail, endDetail, endDetail, id)
	if err != nil {
		return fmt.Errorf("end acquisition anomaly %d: %w", id, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("end acquisition anomaly %d: rows: %w", id, err)
	}
	if n != 1 {
		return fmt.Errorf("end acquisition anomaly %d: %w", id, ErrNotActive)
	}
	if err := r.prune(ctx, tx, deviceID); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("end acquisition anomaly %d: commit: %w", id, err)
	}
	return nil
}

// ListActive returns every active event, oldest first.
func (r *SQLRepository) ListActive(ctx context.Context) ([]Event, error) {
	return r.query(ctx, `SELECT `+eventColumns+` FROM acquisition_anomalies WHERE ended_at IS NULL ORDER BY started_at, id`)
}

// List returns one device's events, newest first, with the filtered total.
func (r *SQLRepository) List(ctx context.Context, f ListFilter) ([]Event, int, error) {
	if err := validateFilter(f); err != nil {
		return nil, 0, err
	}
	where := []string{"device_id = ?"}
	args := []any{f.DeviceID}
	switch f.State {
	case StateActive:
		where = append(where, "ended_at IS NULL")
	case StateEnded:
		where = append(where, "ended_at IS NOT NULL")
	}
	if f.Kind != "" {
		where = append(where, "kind = ?")
		args = append(args, string(f.Kind))
	}
	clause := strings.Join(where, " AND ")
	// COUNT 與分頁 SELECT 在同一個唯讀交易內讀取，避免偵測器同時寫入時 total 與頁面不一致。
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, 0, fmt.Errorf("begin acquisition anomaly list: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var total int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM acquisition_anomalies WHERE `+clause, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count acquisition anomalies: %w", err)
	}
	limit := f.Limit
	if limit == 0 {
		limit = -1
	}
	events, err := queryEvents(ctx, tx, `SELECT `+eventColumns+` FROM acquisition_anomalies WHERE `+clause+` ORDER BY started_at DESC, id DESC LIMIT ? OFFSET ?`, append(args, limit, f.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	if err := tx.Commit(); err != nil {
		return nil, 0, fmt.Errorf("finish acquisition anomaly list: %w", err)
	}
	return events, total, nil
}

// ActiveCounts returns the number of active events per device (devices with
// none are omitted).
func (r *SQLRepository) ActiveCounts(ctx context.Context) (map[string]int, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT device_id, COUNT(*) FROM acquisition_anomalies WHERE ended_at IS NULL GROUP BY device_id`)
	if err != nil {
		return nil, fmt.Errorf("count active acquisition anomalies: %w", err)
	}
	defer rows.Close()
	counts := map[string]int{}
	for rows.Next() {
		var id string
		var n int
		if err := rows.Scan(&id, &n); err != nil {
			return nil, fmt.Errorf("scan active acquisition anomaly count: %w", err)
		}
		counts[id] = n
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("count active acquisition anomalies: %w", err)
	}
	return counts, nil
}

func (r *SQLRepository) query(ctx context.Context, query string, args ...any) ([]Event, error) {
	return queryEvents(ctx, r.db, query, args...)
}

type eventQuerier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func queryEvents(ctx context.Context, q eventQuerier, query string, args ...any) ([]Event, error) {
	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list acquisition anomalies: %w", err)
	}
	defer rows.Close()
	var events []Event
	for rows.Next() {
		var ev Event
		var point, lastValue sql.NullString
		var kind string
		var started int64
		var ended sql.NullInt64
		if err := rows.Scan(&ev.ID, &ev.DeviceID, &point, &kind, &started, &ended, &lastValue, &ev.Detail); err != nil {
			return nil, fmt.Errorf("scan acquisition anomaly: %w", err)
		}
		ev.PointID, ev.Kind, ev.StartedAt = point.String, Kind(kind), time.UnixMilli(started).UTC()
		if ended.Valid {
			t := time.UnixMilli(ended.Int64).UTC()
			ev.EndedAt = &t
		}
		if lastValue.Valid {
			v := lastValue.String
			ev.LastValue = &v
		}
		events = append(events, ev)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list acquisition anomalies: %w", err)
	}
	return events, nil
}

func nullString(s string) any {
	if s == "" {
		return nil
	}
	return s
}
