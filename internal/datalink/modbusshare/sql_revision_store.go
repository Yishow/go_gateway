package modbusshare

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go-gateway/internal/datalink/sqltx"
	"time"
)

const workspaceRevisionKeyPrefix = "modbus_share_workspace_revision:"
const desiredMappingsKeyPrefix = "modbus_share_desired_mappings:"

// SQLWorkspaceRevisionStore persists Share workspace CAS state in the same
// durable system-settings store as the workspace. It is intentionally small:
// desired mappings remain source-rule candidate snapshots, while this record
// is the durable revision/dirty boundary used by the reconciler.
type SQLWorkspaceRevisionStore struct {
	db *sql.DB
}

// DurableDesiredMappingStore is the production persistence seam for the
// complete workspace projection. Implementations commit desired mappings and
// the workspace revision in one database transaction.
type DurableDesiredMappingStore interface {
	WorkspaceRevisionStore
	GetDesiredMappings(context.Context, string) ([]DesiredMapping, error)
	CommitDesiredMappings(context.Context, string, string, string, []DesiredMapping) error
}

type desiredMappingsRecord struct {
	WorkspaceID string           `json:"workspace_id"`
	Mappings    []DesiredMapping `json:"mappings"`
}

// GetDesiredMappings loads the durable desired mapping set for a workspace.
func (s *SQLWorkspaceRevisionStore) GetDesiredMappings(ctx context.Context, workspaceID string) ([]DesiredMapping, error) {
	if s == nil || s.db == nil {
		return nil, newStorageError("read desired mappings", fmt.Errorf("database is not configured"))
	}
	var raw string
	err := s.db.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key = ?`, desiredMappingsKey(workspaceID)).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, newStorageError("read desired mappings", err)
	}
	var record desiredMappingsRecord
	if err := json.Unmarshal([]byte(raw), &record); err != nil {
		return nil, newStorageError("decode desired mappings", err)
	}
	return record.Mappings, nil
}

// CommitDesiredMappings atomically stores desired mappings and advances the
// workspace revision when expectedRevision still matches.
func (s *SQLWorkspaceRevisionStore) CommitDesiredMappings(ctx context.Context, workspaceID, expectedRevision, nextRevision string, mappings []DesiredMapping) (resultErr error) {
	if s == nil || s.db == nil {
		return newStorageError("commit desired mappings", fmt.Errorf("database is not configured"))
	}
	ownedTx, err := sqltx.Begin(ctx, s.db, nil)
	if err != nil {
		return newStorageError("begin desired mappings transaction", err)
	}
	tx := ownedTx.Tx
	defer func() {
		if rollbackErr := ownedTx.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			resultErr = errors.Join(resultErr, newStorageError("rollback workspace revision transaction", rollbackErr))
		}
	}()
	key := workspaceRevisionKey(workspaceID)
	current := workspaceRevisionRecord{Revision: initialRevision}
	var raw string
	err = tx.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key = ?`, key).Scan(&raw)
	found := err == nil
	if found {
		if unmarshalErr := json.Unmarshal([]byte(raw), &current); unmarshalErr != nil {
			return newStorageError("decode workspace revision", unmarshalErr)
		}
		if current.Revision == "" {
			current.Revision = initialRevision
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return newStorageError("read workspace revision", err)
	}
	if expectedRevision != "" && expectedRevision != current.Revision {
		return &Error{Code: ErrCodeRevisionConflict, Message: workspaceRevisionConflictMessage, Retryable: true, WorkspaceRevision: current.Revision}
	}
	recoveryWitness, recovering := ctx.Value(formatRecoveryWitnessKey{}).(formatRecoveryWitness)
	if recovering {
		if recoveryWitness.workspaceID != workspaceID || expectedRevision != nextRevision || current.Revision != expectedRevision {
			return NewError(ErrCodeRevisionConflict, "format recovery scope changed", true)
		}
		if err := validateFormatRecoveryWitness(workspaceID, recoveryWitness.raw, raw); err != nil {
			return err
		}
	}
	previousRevision := current.Revision
	current.Revision, current.Dirty = nextRevision, false
	currentPayload, marshalErr := json.Marshal(current)
	if marshalErr != nil {
		return newStorageError("encode workspace revision", marshalErr)
	}
	if recovering {
		currentPayload = []byte(raw)
	}
	if found {
		query := `UPDATE system_settings SET value = ?, updated_at = ? WHERE key = ? AND json_extract(value, '$.revision') = ?`
		expected := previousRevision
		if recovering {
			query = `UPDATE system_settings SET value = ?, updated_at = ? WHERE key = ? AND value = ?`
			expected = recoveryWitness.raw
		}
		result, updateErr := tx.ExecContext(ctx, query, string(currentPayload), time.Now().UTC(), key, expected)
		if updateErr != nil {
			return newStorageError("update workspace revision", updateErr)
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return newStorageError("update workspace revision rows affected", rowsErr)
		}
		if rows == 0 {
			return &Error{Code: ErrCodeRevisionConflict, Message: workspaceRevisionConflictMessage, Retryable: true}
		}
		if rows != 1 {
			return newStorageError("update workspace revision rows affected", fmt.Errorf("unexpected affected row count %d", rows))
		}
	} else if err := insertWorkspaceRevision(ctx, tx, key, current); err != nil {
		return err
	}
	payload, err := json.Marshal(desiredMappingsRecord{WorkspaceID: workspaceID, Mappings: mappings})
	if err != nil {
		return newStorageError("encode desired mappings", err)
	}
	if err := upsertSystemSetting(ctx, tx, desiredMappingsKey(workspaceID), string(payload)); err != nil {
		return newStorageError("persist desired mappings", err)
	}
	if err := ownedTx.Commit(); err != nil {
		return newStorageError("commit desired mappings", err)
	}
	return nil
}

// NewSQLWorkspaceRevisionStore creates a durable workspace revision store.
func NewSQLWorkspaceRevisionStore(db *sql.DB) *SQLWorkspaceRevisionStore {
	return &SQLWorkspaceRevisionStore{db: db}
}

type workspaceRevisionRecord struct {
	Revision string `json:"revision"`
	Dirty    bool   `json:"dirty"`
}

// GetRevision returns the current durable workspace revision and dirty state.
func (s *SQLWorkspaceRevisionStore) GetRevision(ctx context.Context, workspaceID string) (revision string, isDirty bool, err error) {
	if s == nil || s.db == nil {
		return "", false, newStorageError("read workspace revision", fmt.Errorf("database is not configured"))
	}
	var raw string
	err = s.db.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key = ?`, workspaceRevisionKey(workspaceID)).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return initialRevision, false, nil
	}
	if err != nil {
		return "", false, newStorageError("read workspace revision", err)
	}
	var record workspaceRevisionRecord
	if err := json.Unmarshal([]byte(raw), &record); err != nil {
		return "", false, newStorageError("decode workspace revision", err)
	}
	if record.Revision == "" {
		record.Revision = initialRevision
	}
	return record.Revision, record.Dirty, nil
}

// UpdateRevision advances a workspace revision when expectedRevision matches.
func (s *SQLWorkspaceRevisionStore) UpdateRevision(ctx context.Context, workspaceID, expectedRevision, nextRevision string) (resultErr error) {
	if s == nil || s.db == nil {
		return newStorageError("update workspace revision", fmt.Errorf("database is not configured"))
	}
	ownedTx, err := sqltx.Begin(ctx, s.db, nil)
	if err != nil {
		return newStorageError("begin workspace revision transaction", err)
	}
	tx := ownedTx.Tx
	defer func() {
		if rollbackErr := ownedTx.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			resultErr = errors.Join(resultErr, newStorageError("rollback workspace revision transaction", rollbackErr))
		}
	}()

	key := workspaceRevisionKey(workspaceID)
	var raw string
	err = tx.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key = ?`, key).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		if expectedRevision != "" && expectedRevision != initialRevision {
			return NewError(ErrCodeRevisionConflict, "workspace revision conflict", true)
		}
		if err := insertWorkspaceRevision(ctx, tx, key, workspaceRevisionRecord{Revision: nextRevision}); err != nil {
			return err
		}
		if err := ownedTx.Commit(); err != nil {
			return newStorageError("commit workspace revision", err)
		}
		return nil
	}
	if err != nil {
		return newStorageError("read workspace revision", err)
	}
	var record workspaceRevisionRecord
	if err := json.Unmarshal([]byte(raw), &record); err != nil {
		return newStorageError("decode workspace revision", err)
	}
	if record.Revision == "" {
		record.Revision = "rev-1"
	}
	if expectedRevision != "" && expectedRevision != record.Revision {
		return &Error{Code: ErrCodeRevisionConflict, Message: "workspace revision conflict", Retryable: true, WorkspaceRevision: record.Revision}
	}
	previousRevision := record.Revision
	record.Revision = nextRevision
	record.Dirty = false
	payload, err := json.Marshal(record)
	if err != nil {
		return newStorageError("encode workspace revision", err)
	}
	result, err := tx.ExecContext(ctx, `UPDATE system_settings SET value = ?, updated_at = ? WHERE key = ? AND json_extract(value, '$.revision') = ?`, string(payload), time.Now().UTC(), key, previousRevision)
	if err != nil {
		return newStorageError("update workspace revision", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return newStorageError("update workspace revision rows affected", err)
	}
	if rows == 0 {
		return &Error{Code: ErrCodeRevisionConflict, Message: "workspace revision conflict", Retryable: true, WorkspaceRevision: record.Revision}
	}
	if rows != 1 {
		return newStorageError("update workspace revision rows affected", fmt.Errorf("unexpected affected row count %d", rows))
	}
	if err := ownedTx.Commit(); err != nil {
		return newStorageError("commit workspace revision", err)
	}
	return nil
}

// MarkDirty marks the durable workspace projection as requiring recovery.
func (s *SQLWorkspaceRevisionStore) MarkDirty(ctx context.Context, workspaceID string) error {
	return s.markDirtyCurrentRevision(ctx, workspaceID)
}
