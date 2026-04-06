ALTER TABLE source_rule_tag_review_decisions
    DROP COLUMN IF EXISTS stale_at;

ALTER TABLE source_rule_tag_review_decisions
    DROP COLUMN IF EXISTS stale_revision_id;

ALTER TABLE source_rule_tag_review_decisions
    DROP COLUMN IF EXISTS stale;
