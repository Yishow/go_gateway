package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const storageKey = "studio_v2_workspace"

// workspaceSQLRunner is the query surface shared by *sql.DB and *sql.Tx.
type workspaceSQLRunner interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Get(ctx context.Context) (*Record, error) {
	return readWorkspaceRecord(ctx, r.db)
}

func (r *SQLRepository) Save(ctx context.Context, record *Record) error {
	return writeWorkspaceRecord(ctx, r.db, record)
}

// BeginDatabaseSetup opens the local transaction used by UpdateDatabaseSetup.
func (r *SQLRepository) BeginDatabaseSetup(ctx context.Context) (DatabaseSetupTx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin studio v2 workspace transaction: %w", err)
	}
	return &sqlDatabaseSetupTx{tx: tx}, nil
}

type sqlDatabaseSetupTx struct {
	tx *sql.Tx
}

func (t *sqlDatabaseSetupTx) SQLTx() *sql.Tx {
	return t.tx
}

func (t *sqlDatabaseSetupTx) Get(ctx context.Context) (*Record, error) {
	return readWorkspaceRecord(ctx, t.tx)
}

func (t *sqlDatabaseSetupTx) Save(ctx context.Context, record *Record) error {
	return writeWorkspaceRecord(ctx, t.tx, record)
}

func (t *sqlDatabaseSetupTx) Commit() error {
	return t.tx.Commit()
}

func (t *sqlDatabaseSetupTx) Rollback() error {
	return t.tx.Rollback()
}

func readWorkspaceRecord(ctx context.Context, runner workspaceSQLRunner) (*Record, error) {
	row := runner.QueryRowContext(ctx, `
		SELECT value
		FROM system_settings
		WHERE key = ?
	`, storageKey)

	var raw string
	if err := row.Scan(&raw); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query studio v2 workspace: %w", err)
	}

	var record Record
	if err := json.Unmarshal([]byte(raw), &record); err != nil {
		return nil, fmt.Errorf("decode studio v2 workspace: %w", err)
	}
	if record.OrderedDeviceIDs == nil {
		record.OrderedDeviceIDs = []string{}
	}

	return cloneRecord(&record), nil
}

func writeWorkspaceRecord(ctx context.Context, runner workspaceSQLRunner, record *Record) error {
	payload, err := json.Marshal(cloneRecord(record))
	if err != nil {
		return fmt.Errorf("encode studio v2 workspace: %w", err)
	}

	_, err = runner.ExecContext(ctx, `
		INSERT INTO system_settings (key, value, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			updated_at = excluded.updated_at
	`, storageKey, string(payload), time.Now().UTC())
	if err != nil {
		return fmt.Errorf("persist studio v2 workspace: %w", err)
	}

	return nil
}
