package point

import (
	"context"
	"database/sql"
	"errors"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
)

func overrideShapeChanged(before, after *schema.Point) bool {
	return before.FormatMode == schema.PointFormatOverride && after.FormatMode == schema.PointFormatOverride && (before.DataType != after.DataType || before.Address != after.Address || before.Function != after.Function)
}
func validateSQLOverrideUpdate(ctx context.Context, tx *sql.Tx, after *schema.Point, expected string) error {
	var before schema.Point
	err := tx.QueryRowContext(ctx, `SELECT data_type,address,function,format_mode FROM points WHERE id=? AND device_id=? AND revision_id=?`, after.ID, after.DeviceID, expected).Scan(&before.DataType, &before.Address, &before.Function, &before.FormatMode)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrFormatConflict
	}
	if err != nil {
		return ErrFormatUnavailable
	}
	if !overrideShapeChanged(&before, after) {
		return nil
	}
	var protocol schema.ProtocolType
	if err := tx.QueryRowContext(ctx, `SELECT protocol FROM devices WHERE id=?`, after.DeviceID).Scan(&protocol); err != nil {
		return ErrFormatUnavailable
	}
	if _, err := decoding.ResolvePoint(*after, decoding.FormatContext{Protocol: protocol}); err != nil {
		return ErrFormatInvalid
	}
	return nil
}
func validateMemoryOverrideUpdate(before, after *schema.Point) error {
	if !overrideShapeChanged(before, after) {
		return nil
	}
	// Memory storage has no device context. Reject universally inapplicable types,
	// and fail closed rather than inventing a protocol for other shape changes.
	_, err := decoding.FormatCapability("", after.Address, after.DataType)
	if errors.Is(err, decoding.ErrInapplicableFormat) {
		return ErrFormatInvalid
	}
	return ErrFormatUnavailable
}
