package point

import (
	"context"
	"database/sql"
)

// FormatRecoverySnapshotValidator consumes the same owned SQL snapshot. It is
// an internal composition seam for canonical Share evidence, never client input.
// Implementations must use tx exclusively and must not publish runtime state.
type FormatRecoverySnapshotValidator func(context.Context, *sql.Tx) error

func validateRecoverySnapshot(ctx context.Context, tx *sql.Tx, validators []FormatRecoverySnapshotValidator) error {
	for _, validate := range validators {
		if validate == nil {
			return ErrFormatUnavailable
		}
		if err := validate(ctx, tx); err != nil {
			return err
		}
	}
	return nil
}
