package sourcerule

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

func (r *SQLRepository) UpsertTagReviewDecision(ctx context.Context, decision *schema.SourceRuleTagReviewDecision) error {
	if decision == nil {
		return nil
	}

	createdAt := decision.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	updatedAt := decision.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now().UTC()
	}

	query := `
		INSERT INTO source_rule_tag_review_decisions (
			source_rule_id, candidate_id, decision_type, tag_key, override_tag_id,
			stale, stale_revision_id, stale_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(source_rule_id, candidate_id) DO UPDATE SET
			decision_type = excluded.decision_type,
			tag_key = excluded.tag_key,
			override_tag_id = excluded.override_tag_id,
			stale = excluded.stale,
			stale_revision_id = excluded.stale_revision_id,
			stale_at = excluded.stale_at,
			updated_at = excluded.updated_at
	`

	var overrideTagID interface{}
	if decision.OverrideTagID != nil {
		overrideTagID = *decision.OverrideTagID
	}
	var staleAt interface{}
	if decision.StaleAt != nil {
		staleAt = decision.StaleAt.UTC().Format(time.RFC3339Nano)
	}

	if _, err := r.db.ExecContext(
		ctx,
		query,
		decision.SourceRuleID,
		decision.CandidateID,
		decision.Action,
		decision.TagKey,
		overrideTagID,
		decision.Stale,
		decision.StaleRevisionID,
		staleAt,
		createdAt,
		updatedAt,
	); err != nil {
		return fmt.Errorf("儲存來源規則 tag review decision 失敗: %w", err)
	}
	return nil
}

func (r *SQLRepository) GetTagReviewDecision(ctx context.Context, ruleID, candidateID string) (*schema.SourceRuleTagReviewDecision, error) {
	query := `
		SELECT source_rule_id, candidate_id, decision_type, tag_key, override_tag_id,
		       stale, stale_revision_id, stale_at, created_at, updated_at
		FROM source_rule_tag_review_decisions
		WHERE source_rule_id = ? AND candidate_id = ?
	`

	row := r.db.QueryRowContext(ctx, query, ruleID, candidateID)
	decision, err := scanTagReviewDecision(row.Scan)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("取得來源規則 tag review decision 失敗: %w", err)
	}
	return decision, nil
}

func (r *SQLRepository) ListTagReviewDecisions(ctx context.Context, ruleID string) ([]*schema.SourceRuleTagReviewDecision, error) {
	query := `
		SELECT source_rule_id, candidate_id, decision_type, tag_key, override_tag_id,
		       stale, stale_revision_id, stale_at, created_at, updated_at
		FROM source_rule_tag_review_decisions
		WHERE source_rule_id = ?
		ORDER BY candidate_id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, ruleID)
	if err != nil {
		return nil, fmt.Errorf("列出來源規則 tag review decisions 失敗: %w", err)
	}
	defer rows.Close()

	items := make([]*schema.SourceRuleTagReviewDecision, 0)
	for rows.Next() {
		decision, scanErr := scanTagReviewDecision(rows.Scan)
		if scanErr != nil {
			return nil, fmt.Errorf("掃描來源規則 tag review decision 失敗: %w", scanErr)
		}
		items = append(items, decision)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("走訪來源規則 tag review decisions 失敗: %w", err)
	}
	return items, nil
}

type scanFunc func(dest ...any) error

func scanTagReviewDecision(scan scanFunc) (*schema.SourceRuleTagReviewDecision, error) {
	var decision schema.SourceRuleTagReviewDecision
	var overrideTagID sql.NullString
	var staleAt sql.NullString
	var createdAt string
	var updatedAt string
	if err := scan(
		&decision.SourceRuleID,
		&decision.CandidateID,
		&decision.Action,
		&decision.TagKey,
		&overrideTagID,
		&decision.Stale,
		&decision.StaleRevisionID,
		&staleAt,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, err
	}
	if overrideTagID.Valid {
		decision.OverrideTagID = stringPtr(overrideTagID.String)
	}
	if staleAt.Valid {
		parsed, err := common.ParseTimeString(staleAt.String)
		if err != nil {
			return nil, fmt.Errorf("解析 tag review decision stale 時間失敗: %w", err)
		}
		decision.StaleAt = &parsed
	}
	var err error
	decision.CreatedAt, err = common.ParseTimeString(createdAt)
	if err != nil {
		return nil, fmt.Errorf("解析 tag review decision 建立時間失敗: %w", err)
	}
	decision.UpdatedAt, err = common.ParseTimeString(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("解析 tag review decision 更新時間失敗: %w", err)
	}
	return &decision, nil
}
