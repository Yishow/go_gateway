DROP INDEX IF EXISTS idx_mappings_rule_candidate_id;

ALTER TABLE mappings
    DROP COLUMN IF EXISTS blocking_reason;

ALTER TABLE mappings
    DROP COLUMN IF EXISTS last_applied_signature;

ALTER TABLE mappings
    DROP COLUMN IF EXISTS proposed_signature;

ALTER TABLE mappings
    DROP COLUMN IF EXISTS rule_candidate_id;

ALTER TABLE mappings
    DROP COLUMN IF EXISTS status;
