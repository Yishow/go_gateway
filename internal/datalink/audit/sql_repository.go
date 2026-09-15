package audit

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// SQLRepository stores audit history in SQL.
type SQLRepository struct {
	db *sql.DB
}

// NewSQLRepository creates a SQL-backed audit repository.
func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

// Create stores an audit history entry.
func (r *SQLRepository) Create(ctx context.Context, entry *Entry) error {
	if entry == nil {
		return nil
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO workspace_audit_history (
			id, workspace_id, event_type, result, scope, reference_id, details, occurred_at, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, entry.ID, entry.WorkspaceID, string(entry.EventType), string(entry.Result), entry.Scope, entry.ReferenceID, entry.Details, entry.OccurredAt, entry.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert workspace audit history: %w", err)
	}
	return nil
}

// List returns recent audit history entries.
func (r *SQLRepository) List(ctx context.Context, filter ListFilter) ([]Entry, error) {
	query := `
		SELECT id, workspace_id, event_type, result, scope, reference_id, details, occurred_at, created_at
		FROM workspace_audit_history
		WHERE workspace_id = ?`
	args := []any{strings.TrimSpace(filter.WorkspaceID)}
	if filter.EventType != "" {
		query += ` AND event_type = ?`
		args = append(args, string(filter.EventType))
	}
	query += ` ORDER BY occurred_at DESC, created_at DESC LIMIT ?`
	args = append(args, normalizeLimit(filter.Limit))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list workspace audit history: %w", err)
	}
	defer rows.Close()

	entries := make([]Entry, 0)
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(
			&entry.ID,
			&entry.WorkspaceID,
			&entry.EventType,
			&entry.Result,
			&entry.Scope,
			&entry.ReferenceID,
			&entry.Details,
			&entry.OccurredAt,
			&entry.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan workspace audit history: %w", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate workspace audit history: %w", err)
	}
	return entries, nil
}
