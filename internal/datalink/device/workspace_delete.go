package device

import (
	"context"
	"database/sql"
	"fmt"
)

// DeleteWorkspaceDeviceInTx deletes a device using the shared configuration
// transaction, only when the device repository owns that same database.
func (s *Service) DeleteWorkspaceDeviceInTx(ctx context.Context, tx *sql.Tx, db *sql.DB, id string) error {
	repo, ok := s.repo.(*SQLRepository)
	if !ok || repo == nil || db == nil || repo.db != db || tx == nil {
		return fmt.Errorf("workspace device transaction unavailable")
	}
	result, err := tx.ExecContext(ctx, `DELETE FROM devices WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete workspace device: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read workspace device deletion result: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("workspace device not found")
	}
	return nil
}

// RetireDeletedWorkspaceDevice closes a connection after its deletion commits.
func (s *Service) RetireDeletedWorkspaceDevice(id string) error {
	if s == nil || s.connMgr == nil {
		return nil
	}
	return s.connMgr.Close(id)
}
