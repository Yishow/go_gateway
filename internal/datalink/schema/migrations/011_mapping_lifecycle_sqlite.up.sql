ALTER TABLE mappings ADD COLUMN status TEXT NOT NULL DEFAULT 'active';
ALTER TABLE mappings ADD COLUMN rule_candidate_id TEXT;
ALTER TABLE mappings ADD COLUMN proposed_signature TEXT;
ALTER TABLE mappings ADD COLUMN last_applied_signature TEXT;
ALTER TABLE mappings ADD COLUMN blocking_reason TEXT;

UPDATE mappings
SET status = CASE WHEN enabled = 1 THEN 'active' ELSE 'draft' END
WHERE status IS NULL OR status = '' OR (enabled = 0 AND status = 'active');

CREATE INDEX IF NOT EXISTS idx_mappings_rule_candidate_id ON mappings(rule_candidate_id);
