package point

import (
	"context"
	"database/sql"
	"go-gateway/internal/datalink/schema"
)

type pointConfigExecer interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

// UpdateRevision preserves a concurrent format edit rather than overwriting it
// from a stale ordinary edit. Telemetry fields are never part of configuration CAS.
func (r *SQLRepository) UpdateRevision(ctx context.Context, p *schema.Point, expected string) error {
	return writePointConfigRevision(ctx, r.db, p, expected)
}
func writePointConfigRevision(ctx context.Context, db pointConfigExecer, p *schema.Point, expected string) error {
	result, err := db.ExecContext(ctx, `UPDATE points SET name=?,description=?,address=?,function=?,data_type=?,data_format=?,format_mode=?,format_override=?,format_version=?,revision_id=?,mode=?,polling_group_id=?,enabled=?,updated_at=? WHERE id=? AND device_id=? AND revision_id=?`, p.Name, p.Description, p.Address, p.Function, p.DataType, p.DataFormat, p.FormatMode, p.FormatOverride, p.FormatVersion, p.RevisionID, p.Mode, p.PollingGroupID, p.Enabled, p.UpdatedAt, p.ID, p.DeviceID, expected)
	if err != nil {
		return ErrFormatUnavailable
	}
	n, err := result.RowsAffected()
	if err != nil {
		return ErrFormatUnavailable
	}
	if n != 1 {
		return ErrFormatConflict
	}
	return nil
}

// UpdateRevision mirrors production CAS for service-level memory fixtures.
func (r *MemoryRepository) UpdateRevision(_ context.Context, p *schema.Point, expected string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	old, ok := r.points[p.ID]
	if !ok {
		return ErrPointNotFound
	}
	if old.RevisionID != expected {
		return ErrFormatConflict
	}
	next := clonePoint(p)
	next.LastValue = old.LastValue
	next.LastReadAt = old.LastReadAt
	next.LastError = old.LastError
	r.points[p.ID] = next
	return nil
}
