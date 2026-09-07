package delivery

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// SQLJournal 實作基於關聯資料庫的本機耐久日誌。
type SQLJournal struct {
	db *sql.DB
}

// NewSQLJournal 建立 SQLJournal 實例。
func NewSQLJournal(db *sql.DB) *SQLJournal {
	return &SQLJournal{db: db}
}

func (j *SQLJournal) Append(entry *JournalEntry) error {
	if entry == nil {
		return fmt.Errorf("journal entry cannot be nil")
	}

	query := `
		INSERT INTO gw_delivery_journal (record_id, observed_at, payload, created_at)
		VALUES (?, ?, ?, ?)
	`
	res, err := j.db.ExecContext(context.Background(), query,
		entry.RecordID,
		entry.ObservedAt.Format(time.RFC3339Nano),
		entry.Payload,
		time.Now().UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		return fmt.Errorf("failed to append to journal: %w", err)
	}

	seq, err := res.LastInsertId()
	if err == nil {
		entry.Sequence = seq
	}
	return nil
}

func (j *SQLJournal) ScanFrom(fromSequence int64, limit int) ([]*JournalEntry, error) {
	if limit <= 0 {
		limit = 100
	}

	query := `
		SELECT sequence, record_id, observed_at, payload
		FROM gw_delivery_journal
		WHERE sequence >= ?
		ORDER BY sequence ASC
		LIMIT ?
	`
	rows, err := j.db.QueryContext(context.Background(), query, fromSequence, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to scan journal: %w", err)
	}
	defer rows.Close()

	var entries []*JournalEntry
	for rows.Next() {
		var (
			seq      int64
			recID    string
			obsStr   string
			payload  []byte
		)
		if err := rows.Scan(&seq, &recID, &obsStr, &payload); err != nil {
			return nil, err
		}
		obsTime, _ := time.Parse(time.RFC3339Nano, obsStr)
		entries = append(entries, &JournalEntry{
			Sequence:   seq,
			RecordID:   recID,
			ObservedAt: obsTime,
			Payload:    payload,
		})
	}
	return entries, nil
}

func (j *SQLJournal) TruncateBefore(sequence int64) error {
	query := `DELETE FROM gw_delivery_journal WHERE sequence < ?`
	_, err := j.db.ExecContext(context.Background(), query, sequence)
	if err != nil {
		return fmt.Errorf("failed to truncate journal before sequence %d: %w", sequence, err)
	}
	return nil
}

func (j *SQLJournal) Close() error {
	return nil
}
