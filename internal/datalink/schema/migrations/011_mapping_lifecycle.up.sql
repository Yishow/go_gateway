ALTER TABLE mappings
    ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'active';

ALTER TABLE mappings
    ADD COLUMN IF NOT EXISTS rule_candidate_id TEXT;

ALTER TABLE mappings
    ADD COLUMN IF NOT EXISTS proposed_signature TEXT;

ALTER TABLE mappings
    ADD COLUMN IF NOT EXISTS last_applied_signature TEXT;

ALTER TABLE mappings
    ADD COLUMN IF NOT EXISTS blocking_reason TEXT;

UPDATE mappings
SET status = CASE WHEN enabled THEN 'active' ELSE 'draft' END
WHERE status IS NULL OR status = '' OR (enabled = FALSE AND status = 'active');

CREATE INDEX IF NOT EXISTS idx_mappings_rule_candidate_id ON mappings(rule_candidate_id);
