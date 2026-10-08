package sourcerule

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-gateway/internal/datalink/schema"
	"slices"
)

func validateReplacementLinks(ruleID string, links []*schema.SourceRuleLink) error {
	if ruleID == "" {
		return fmt.Errorf("source-rule replacement identity is required")
	}
	ids, points, addresses := map[string]bool{}, map[string]bool{}, map[string]bool{}
	for _, link := range links {
		if link == nil || link.ID == "" || link.PointID == "" || link.Address == "" || link.RuleID != ruleID || ids[link.ID] || points[link.PointID] || addresses[link.Address] {
			return fmt.Errorf("invalid or duplicate source-rule replacement link")
		}
		ids[link.ID], points[link.PointID], addresses[link.Address] = true, true, true
	}
	return nil
}

// ReplaceLinks deletes and inserts in one transaction; insertion or commit
// failures preserve the previous complete link set.
func (r *SQLRepository) ReplaceLinks(ctx context.Context, ruleID string, links []*schema.SourceRuleLink) (err error) {
	if err := validateReplacementLinks(ruleID, links); err != nil {
		return err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin source-rule link replacement: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = errors.Join(err, fmt.Errorf("rollback source-rule link replacement: %w", rollbackErr))
		}
	}()
	if _, err := tx.ExecContext(ctx, `DELETE FROM source_rule_links WHERE rule_id = ?`, ruleID); err != nil {
		return fmt.Errorf("delete replaced source-rule links: %w", err)
	}
	for _, link := range links {
		if _, err := tx.ExecContext(ctx, `INSERT INTO source_rule_links (id, rule_id, address, point_id, tag_id, mapping_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, link.ID, link.RuleID, link.Address, link.PointID, link.TagID, link.MappingID, link.CreatedAt, link.UpdatedAt); err != nil {
			return fmt.Errorf("insert replacement source-rule link: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit source-rule link replacement: %w", err)
	}
	return nil
}

// ReplaceLinks swaps one cloned set under the repository mutex.
func (r *MemoryRepository) ReplaceLinks(_ context.Context, ruleID string, links []*schema.SourceRuleLink) error {
	if err := validateReplacementLinks(ruleID, links); err != nil {
		return err
	}
	next := cloneSourceRuleLinks(links)
	slices.SortFunc(next, func(a, b *schema.SourceRuleLink) int { return cmp.Compare(a.Address, b.Address) })
	r.mu.Lock()
	defer r.mu.Unlock()
	for otherRule, existing := range r.links {
		if otherRule == ruleID {
			continue
		}
		for _, old := range existing {
			for _, link := range next {
				if old.ID == link.ID || old.PointID == link.PointID {
					return fmt.Errorf("source-rule link identity already owned")
				}
			}
		}
	}
	r.links[ruleID] = next
	return nil
}
