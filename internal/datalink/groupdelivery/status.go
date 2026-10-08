package groupdelivery

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"
)

// StageCounts is the number of rows of one group in each delivery stage. A row
// is counted in exactly one stage, and only SQLCommitted is destination
// evidence.
type StageCounts struct {
	// Collecting are accepted samples whose bucket has not closed yet.
	Collecting int
	// Queued are closed rows durably waiting for a first or next attempt
	// (including one being sent right now).
	Queued      int
	Retrying    int
	Blocked     int
	Quarantined int
	Unknown     int
	// SQLCommitted rows were confirmed by the destination.
	SQLCommitted int
	Skipped      int
}

// RevisionBacklog is the undelivered backlog frozen to one group and
// destination revision.
type RevisionBacklog struct {
	GroupRevision     string
	ConnectorID       string
	ConnectorRevision string
	TableSchema       string
	TableName         string
	Pending           int
	LastErrorCodes    []string
}

// GroupStatus is the single source of delivery truth for one group.
type GroupStatus struct {
	GroupID string
	Stages  StageCounts
	// LastSQLCommittedAt is the latest destination-confirmed commit, nil when
	// nothing was ever confirmed. It is never inferred from buffering or ACK.
	LastSQLCommittedAt *time.Time
	// OldestPendingSeconds is the age of the oldest undelivered row.
	OldestPendingSeconds float64
	Backlog              []RevisionBacklog
	// SilentOrSkipped counts closed buckets that produced no SQL row.
	NoDataBuckets  int
	SkippedBuckets int
	// RecentBucketIssues contains at most 20 closed buckets that did not
	// produce a row. Causes are safe quality categories, never sample values.
	RecentBucketIssues []BucketIssue
}

// GroupStatus reads the stage counts and revision-bound backlog of a group.
func (s *Store) GroupStatus(ctx context.Context, groupID string) (GroupStatus, error) {
	status := GroupStatus{GroupID: groupID}
	if err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id = ? AND consumed = 0`, groupID).Scan(&status.Stages.Collecting); err != nil {
		return GroupStatus{}, fmt.Errorf("count collecting samples: %w", err)
	}
	if err := s.readStageCounts(ctx, groupID, &status); err != nil {
		return GroupStatus{}, err
	}
	if err := s.readTimes(ctx, groupID, &status); err != nil {
		return GroupStatus{}, err
	}
	if err := s.readBacklog(ctx, groupID, &status); err != nil {
		return GroupStatus{}, err
	}
	if err := s.readBucketOutcomes(ctx, groupID, &status); err != nil {
		return GroupStatus{}, err
	}
	if err := s.readRecentBucketIssues(ctx, groupID, &status); err != nil {
		return GroupStatus{}, err
	}
	return status, nil
}

func (s *Store) readStageCounts(ctx context.Context, groupID string, status *GroupStatus) error {
	return s.readRevisionStageCounts(ctx, groupID, "", status)
}

func (s *Store) readRevisionStageCounts(ctx context.Context, groupID, revision string, status *GroupStatus) error {
	query := `SELECT state, COUNT(*) FROM wg_delivery_outbox WHERE group_id = ?`
	args := []any{groupID}
	if revision != "" {
		query += ` AND group_revision = ?`
		args = append(args, revision)
	}
	rows, err := s.db.QueryContext(ctx, query+` GROUP BY state`, args...)
	if err != nil {
		return fmt.Errorf("count delivery stages: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var state string
		var n int
		if err := rows.Scan(&state, &n); err != nil {
			return fmt.Errorf("scan delivery stage: %w", err)
		}
		switch state {
		case StatePending, StateSending:
			status.Stages.Queued += n
		case StateRetrying:
			status.Stages.Retrying += n
		case StateBlocked:
			status.Stages.Blocked += n
		case StateQuarantined:
			status.Stages.Quarantined += n
		case StateUnknown:
			status.Stages.Unknown += n
		case StateCommitted:
			status.Stages.SQLCommitted += n
		case StateSkipped:
			status.Stages.Skipped += n
		}
	}
	return rows.Err()
}

func (s *Store) readTimes(ctx context.Context, groupID string, status *GroupStatus) error {
	var committed, oldest sql.NullString
	if err := s.db.QueryRowContext(ctx, `
		SELECT MAX(CASE WHEN state = 'sql_committed' AND committed_at != '' THEN committed_at END),
		       MIN(CASE WHEN state NOT IN ('sql_committed', 'operator_skipped') THEN created_at END)
		FROM wg_delivery_outbox WHERE group_id = ?`, groupID).Scan(&committed, &oldest); err != nil {
		return fmt.Errorf("read delivery times: %w", err)
	}
	if committed.Valid && committed.String != "" {
		parsed, err := time.Parse(time.RFC3339Nano, committed.String)
		if err != nil {
			return fmt.Errorf("parse committed_at: %w", err)
		}
		parsed = parsed.UTC()
		status.LastSQLCommittedAt = &parsed
	}
	if oldest.Valid && oldest.String != "" {
		parsed, err := time.Parse(time.RFC3339Nano, oldest.String)
		if err != nil {
			return fmt.Errorf("parse created_at: %w", err)
		}
		status.OldestPendingSeconds = s.now().Sub(parsed).Seconds()
	}
	return nil
}

func (s *Store) readBacklog(ctx context.Context, groupID string, status *GroupStatus) error {
	rows, err := s.db.QueryContext(ctx, `
		SELECT group_revision, connector_id, connector_revision, table_schema, table_name, COUNT(*),
		       COALESCE(GROUP_CONCAT(DISTINCT NULLIF(last_error_code, '')), '')
		FROM wg_delivery_outbox
		WHERE group_id = ? AND state NOT IN ('sql_committed', 'operator_skipped')
		GROUP BY group_revision, connector_id, connector_revision, table_schema, table_name
		ORDER BY group_revision, connector_id, connector_revision`, groupID)
	if err != nil {
		return fmt.Errorf("read revision backlog: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var entry RevisionBacklog
		var codes string
		if err := rows.Scan(&entry.GroupRevision, &entry.ConnectorID, &entry.ConnectorRevision,
			&entry.TableSchema, &entry.TableName, &entry.Pending, &codes); err != nil {
			return fmt.Errorf("scan revision backlog: %w", err)
		}
		if codes != "" {
			entry.LastErrorCodes = strings.Split(codes, ",")
			slices.Sort(entry.LastErrorCodes)
		}
		status.Backlog = append(status.Backlog, entry)
	}
	return rows.Err()
}

func (s *Store) readBucketOutcomes(ctx context.Context, groupID string, status *GroupStatus) error {
	rows, err := s.db.QueryContext(ctx,
		`SELECT kind, COUNT(*) FROM wg_delivery_buckets WHERE group_id = ? AND kind IN ('no_data', 'skipped') GROUP BY kind`, groupID)
	if err != nil {
		return fmt.Errorf("count bucket outcomes: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var kind string
		var n int
		if err := rows.Scan(&kind, &n); err != nil {
			return fmt.Errorf("scan bucket outcome: %w", err)
		}
		if kind == noDataBucketKind {
			status.NoDataBuckets = n
		} else {
			status.SkippedBuckets = n
		}
	}
	return rows.Err()
}
