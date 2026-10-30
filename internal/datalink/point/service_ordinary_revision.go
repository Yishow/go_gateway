package point

import (
	"context"
	"database/sql"
	"encoding/json"
	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
	"time"
)

type ordinaryRevisionUpdater interface {
	UpdateOrdinaryRevision(context.Context, *schema.Point, string, bool) error
}

// UpdateOrdinary preserves the server attempt's original snapshot through SQL
// CAS. It never reads a newer revision to make a delayed runtime callback valid.
func (s *Service) UpdateOrdinary(ctx context.Context, before *schema.Point, req UpdatePointRequest, installedBaseline bool) (*schema.Point, error) {
	if before == nil || req.DataFormat != nil || req.FormatMode != nil || req.ConfiguredFormat != nil {
		return nil, ErrFormatInvalid
	}
	next := clonePoint(before)
	if err := applyPointUpdate(next, req); err != nil {
		return nil, err
	}
	revision, err := common.NewUUID()
	if err != nil {
		return nil, ErrFormatUnavailable
	}
	next.RevisionID = revision
	next.UpdatedAt = time.Now()
	repo, ok := s.repo.(ordinaryRevisionUpdater)
	if !ok {
		return nil, ErrFormatUnavailable
	}
	if err := repo.UpdateOrdinaryRevision(ctx, next, before.RevisionID, installedBaseline); err != nil {
		return nil, err
	}
	return next, nil
}
func (r *MemoryRepository) UpdateOrdinaryRevision(ctx context.Context, p *schema.Point, expected string, _ bool) error {
	return r.UpdateRevision(ctx, p, expected)
}
func rejectPendingPointFormat(ctx context.Context, tx *sql.Tx, id string) error {
	rows, err := tx.QueryContext(ctx, `SELECT value FROM system_settings WHERE key LIKE 'point_format_receipt:v1:%' AND (NOT json_valid(value) OR json_extract(value,'$.result.point.id')=?)`, id)
	if err != nil {
		return ErrFormatUnavailable
	}
	defer rows.Close()
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			return ErrFormatUnavailable
		}
		var receipt formatReceipt
		if json.Unmarshal([]byte(raw), &receipt) != nil || receipt.Version != 1 {
			return ErrFormatUnavailable
		}
		// Until recovery proves this format installation applied, an ordinary edit
		// must not turn a deliberately stopped point back into an active old decoder.
		if receipt.Result.Point.ID == id && receipt.Result.RuntimeApplyStatus != "applied" {
			return ErrFormatUnavailable
		}
	}
	if rows.Err() != nil {
		return ErrFormatUnavailable
	}
	return nil
}
