package tag

import (
	"context"
	"fmt"

	"go-gateway/internal/datalink/schema"
)

type sourceRuleRetirementRepository interface {
	deleteRetiredTag(context.Context, *schema.Tag) (bool, error)
}

// DeleteRetiredTag removes a still-unchanged tag after the caller checks that
// it is unbound. The SQL repository also guards the reference atomically.
func (s *Service) DeleteRetiredTag(ctx context.Context, expected *schema.Tag) (bool, error) {
	repo, ok := s.repo.(sourceRuleRetirementRepository)
	if !ok || expected == nil {
		return false, fmt.Errorf("tag retirement unavailable")
	}
	return repo.deleteRetiredTag(ctx, expected)
}

// RestoreRetiredTag inserts a removed tag without replacing another edit.
func (s *Service) RestoreRetiredTag(ctx context.Context, record *schema.Tag) error {
	if record == nil {
		return fmt.Errorf("tag retirement recovery unavailable")
	}
	return s.repo.Create(ctx, cloneTag(record))
}

func (r *SQLRepository) deleteRetiredTag(ctx context.Context, t *schema.Tag) (bool, error) {
	result, err := r.db.ExecContext(ctx, `DELETE FROM tags WHERE id=? AND key=? AND display_name=? AND data_type=? AND COALESCE(unit,'')=? AND COALESCE(description,'')=? AND status=? AND COALESCE(labels,'')=? AND NOT EXISTS (SELECT 1 FROM mappings WHERE tag_id=tags.id)`, t.ID, t.Key, t.DisplayName, t.DataType, t.Unit, t.Description, t.Status, t.Labels)
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	return n == 1, err
}

func (r *MemoryRepository) deleteRetiredTag(_ context.Context, expected *schema.Tag) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	current := r.tags[expected.ID]
	if current == nil || current.Key != expected.Key || current.DisplayName != expected.DisplayName || current.DataType != expected.DataType || current.Unit != expected.Unit || current.Description != expected.Description || current.Status != expected.Status || current.Labels != expected.Labels {
		return false, nil
	}
	delete(r.tags, expected.ID)
	return true, nil
}
