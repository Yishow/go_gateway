package groupdelivery

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"go-gateway/internal/datalink/snapshot"
)

const (
	noDataBucketKind = string(snapshot.OutcomeNoData)
	// truncatedBucketKind names a written row whose window a cutover cut short.
	truncatedBucketKind = "truncated"
)

// BucketIssue is a recent scoped closure that produced no SQL row. Causes
// deliberately contain only known quality categories, not stored error text.
type BucketIssue struct {
	GroupRevision string    `json:"group_revision"`
	BucketStart   time.Time `json:"bucket_start"`
	Kind          string    `json:"kind"`
	Causes        []string  `json:"causes"`
}

func (s *Store) readRecentBucketIssues(ctx context.Context, groupID string, status *GroupStatus) error {
	rows, err := s.db.QueryContext(ctx, `
		SELECT group_revision, bucket_start, kind, members FROM wg_delivery_buckets
		WHERE group_id = ? AND (kind IN ('skipped', 'no_data') OR (kind = 'row' AND reason LIKE 'boundary-truncated:%'))
		ORDER BY bucket_start DESC, group_revision DESC, entity_key ASC LIMIT 20`, groupID)
	if err != nil {
		return fmt.Errorf("read recent bucket issues: %w", err)
	}
	defer rows.Close()
	status.RecentBucketIssues = []BucketIssue{}
	for rows.Next() {
		var issue BucketIssue
		var start, members string
		if err := rows.Scan(&issue.GroupRevision, &start, &issue.Kind, &members); err != nil {
			return fmt.Errorf("scan recent bucket issue: %w", err)
		}
		parsed, err := time.Parse(time.RFC3339Nano, start)
		if err != nil {
			return fmt.Errorf("parse bucket issue timestamp: %w", err)
		}
		issue.BucketStart = parsed.UTC()
		if issue.Kind == string(snapshot.OutcomeRow) {
			issue.Kind = truncatedBucketKind
		}
		issue.Causes = bucketQualityCauses(issue.Kind, members)
		status.RecentBucketIssues = append(status.RecentBucketIssues, issue)
	}
	return rows.Err()
}

func bucketQualityCauses(kind, rawMembers string) []string {
	if kind == noDataBucketKind {
		return []string{noDataBucketKind}
	}
	if kind == truncatedBucketKind {
		// A row was written; the window it stands for was cut short.
		return []string{"boundary"}
	}
	var members []memberRecord
	if err := json.Unmarshal([]byte(rawMembers), &members); err != nil {
		return []string{"unavailable"}
	}
	causes := []string{}
	for _, member := range members {
		switch member.Status {
		case "missing", "bad", "stale", "invalid":
			causes = append(causes, member.Status)
		case "ok":
		default:
			causes = append(causes, "unavailable")
		}
	}
	if len(causes) == 0 {
		return []string{"unavailable"}
	}
	slices.Sort(causes)
	return slices.Compact(causes)
}
