package sourcerule

import (
	"context"
	"database/sql"
	"fmt"
	"go-gateway/internal/datalink/schema"
)

type shareSQLReader interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func readShareList(ctx context.Context, q shareSQLReader, filter ListFilter) ([]*schema.SourceRule, error) {
	query := `
		SELECT id, device_id, start_address, count, data_type, naming_prefix,
		       enabled, locked, origin, template_name, skipped_addresses,
		       target_data_type, scale_multiplier, scale_offset, data_format,
		       share_enabled, share_start_register, share_stride,
		       created_at, updated_at, revision_id
		FROM source_rules
		WHERE 1 = 1
	`
	args := make([]any, 0, 2)

	if filter.DeviceID != nil {
		query += ` AND device_id = ?`
		args = append(args, *filter.DeviceID)
	}
	if filter.Enabled != nil {
		query += ` AND enabled = ?`
		args = append(args, *filter.Enabled)
	}
	query += ` ORDER BY created_at ASC, id ASC`

	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("查詢來源規則失敗: %w", err)
	}
	defer rows.Close()

	items := make([]*schema.SourceRule, 0)
	for rows.Next() {
		rule, scanErr := scanRuleRows(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, rule)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("讀取來源規則結果失敗: %w", err)
	}
	return items, nil
}

func readShareListLinks(ctx context.Context, q shareSQLReader, ruleID string) ([]*schema.SourceRuleLink, error) {
	rows, err := q.QueryContext(ctx, `
		SELECT id, rule_id, address, point_id, tag_id, mapping_id, created_at, updated_at
		FROM source_rule_links
		WHERE rule_id = ?
		ORDER BY address ASC, id ASC
	`, ruleID)
	if err != nil {
		return nil, fmt.Errorf("查詢來源規則連結失敗: %w", err)
	}
	defer rows.Close()

	items := make([]*schema.SourceRuleLink, 0)
	for rows.Next() {
		link, scanErr := scanLinkRows(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, link)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("讀取來源規則連結結果失敗: %w", err)
	}
	return items, nil
}

func readShareListCandidateSnapshots(ctx context.Context, q shareSQLReader, ruleID, revisionID string) ([]*schema.SourceRuleCandidateSnapshot, error) {
	rows, err := q.QueryContext(ctx, `
		SELECT source_rule_id, revision_id, candidate_type, payload, status, reason, generated_at
		FROM source_rule_candidate_snapshots
		WHERE source_rule_id = ? AND revision_id = ?
		ORDER BY candidate_type ASC
	`, ruleID, revisionID)
	if err != nil {
		return nil, fmt.Errorf("查詢候選快照失敗: %w", err)
	}
	defer rows.Close()

	items := make([]*schema.SourceRuleCandidateSnapshot, 0)
	for rows.Next() {
		snapshot, scanErr := scanCandidateSnapshot(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, snapshot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("讀取候選快照結果失敗: %w", err)
	}
	return items, nil
}
