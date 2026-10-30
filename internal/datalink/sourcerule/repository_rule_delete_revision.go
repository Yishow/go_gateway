package sourcerule

import (
	"context"
	"go-gateway/internal/datalink/schema"
	"maps"
	"strings"
)

type ruleRevisionDeleter interface {
	DeleteRevision(context.Context, *schema.SourceRule) error
}

// Deletion and the database's FK cascade are one statement. A stale operation
// cannot first erase the newer rule's links or candidate snapshots.
func (r *SQLRepository) DeleteRevision(ctx context.Context, expected *schema.SourceRule) error {
	query, args := ruleStatePredicate(expected)
	result, err := r.db.ExecContext(ctx, "DELETE FROM source_rules WHERE "+query, args...)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return ErrRuleRevisionConflict
	}
	return nil
}
func (r *MemoryRepository) DeleteRevision(_ context.Context, expected *schema.SourceRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	current, ok := r.rules[expected.ID]
	if !ok || ruleStateKey(current) != ruleStateKey(expected) {
		return ErrRuleRevisionConflict
	}
	delete(r.rules, expected.ID)
	delete(r.links, expected.ID)
	delete(r.snapshots, expected.ID)
	delete(r.tagReviewDecisions, expected.ID)
	maps.DeleteFunc(r.snapshots, func(key string, _ []*schema.SourceRuleCandidateSnapshot) bool {
		return strings.HasPrefix(key, expected.ID+":")
	})
	return nil
}

// IS is SQLite's null-safe equality. Include the complete config rather than
// only updated_at: SQLite's timestamp trigger has second precision, and an
// enable toggle intentionally keeps its rule revision under the existing API.
func ruleStatePredicate(p *schema.SourceRule) (string, []any) {
	return `id=? AND device_id=? AND revision_id=? AND start_address=? AND count=? AND data_type=? AND naming_prefix=? AND enabled=? AND locked=? AND origin=? AND template_name IS ? AND skipped_addresses=? AND target_data_type IS ? AND scale_multiplier IS ? AND scale_offset IS ? AND data_format IS ? AND share_enabled=? AND share_start_register IS ? AND share_stride IS ?`, []any{p.ID, p.DeviceID, p.RevisionID, p.StartAddress, p.Count, p.DataType, p.NamingPrefix, p.Enabled, p.Locked, p.Origin, p.TemplateName, p.SkippedAddresses, p.TargetDataType, p.ScaleMultiplier, p.ScaleOffset, p.DataFormat, p.ShareEnabled, p.ShareStartRegister, p.ShareStride}
}
