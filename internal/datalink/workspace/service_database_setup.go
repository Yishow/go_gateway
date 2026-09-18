package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// ErrSetupRevisionConflict reports a Step 4 database save made against an
// older database setup revision than the one currently persisted.
var ErrSetupRevisionConflict = errors.New("studio v2 database setup revision conflict")

// ErrDatabaseSetupUnavailable reports a workspace store that cannot apply
// database setup changes in one local transaction.
var ErrDatabaseSetupUnavailable = errors.New("studio v2 database setup transactions are unavailable")

// DatabaseSetupTx is one local configuration transaction shared by the
// workspace record and the Step 4 database connector and target rows.
type DatabaseSetupTx interface {
	SQLTx() *sql.Tx
	Get(ctx context.Context) (*Record, error)
	Save(ctx context.Context, record *Record) error
	Commit() error
	Rollback() error
}

// DatabaseSetupStore is implemented by workspace repositories that can open a
// local transaction for Step 4 database setup saves.
type DatabaseSetupStore interface {
	BeginDatabaseSetup(ctx context.Context) (DatabaseSetupTx, error)
}

// UpdateDatabaseSetup applies one Step 4 database setup change atomically.
// While holding the workspace mutex it opens one local transaction, rejects a
// stale expected revision before any write, lets mutate write database rows
// through tx and edit record, then advances the setup revision and saves the
// workspace in the same transaction. mutate must only use tx for database
// access: the embedded SQLite pool has a single connection.
func (s *Service) UpdateDatabaseSetup(
	ctx context.Context,
	expectedRevision string,
	mutate func(ctx context.Context, tx *sql.Tx, record *Record) error,
) (updated *Record, err error) {
	store, ok := s.repo.(DatabaseSetupStore)
	if !ok {
		return nil, ErrDatabaseSetupUnavailable
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	setupTx, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin workspace database setup: %w", err)
	}
	committed := false
	defer func() {
		if committed {
			return
		}
		if rollbackErr := setupTx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = errors.Join(err, fmt.Errorf("rollback workspace database setup: %w", rollbackErr))
		}
	}()

	record, err := setupTx.Get(ctx)
	switch {
	case errors.Is(err, ErrNotFound):
		record = s.newRecord()
	case err != nil:
		return nil, err
	}
	if record.DatabaseSetupRevision != strings.TrimSpace(expectedRevision) {
		return nil, ErrSetupRevisionConflict
	}
	if err := mutate(ctx, setupTx.SQLTx(), record); err != nil {
		return nil, err
	}

	record.DatabaseSetupRevision = s.newID()
	record.UpdatedAt = s.now()
	if err := setupTx.Save(ctx, record); err != nil {
		return nil, fmt.Errorf("save workspace database setup: %w", err)
	}
	if err := setupTx.Commit(); err != nil {
		return nil, fmt.Errorf("commit workspace database setup: %w", err)
	}
	committed = true
	return cloneRecord(record), nil
}

func (s *Service) newRecord() *Record {
	now := s.now()
	return &Record{
		ID:               s.newID(),
		Kind:             WorkspaceKindSingle,
		Status:           WorkspaceStatusEmpty,
		OrderedDeviceIDs: []string{},
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}
