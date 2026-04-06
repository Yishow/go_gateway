ALTER TABLE source_rule_tag_review_decisions
    ADD COLUMN IF NOT EXISTS stale BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE source_rule_tag_review_decisions
    ADD COLUMN IF NOT EXISTS stale_revision_id TEXT NOT NULL DEFAULT '';

ALTER TABLE source_rule_tag_review_decisions
    ADD COLUMN IF NOT EXISTS stale_at TIMESTAMP NULL;
