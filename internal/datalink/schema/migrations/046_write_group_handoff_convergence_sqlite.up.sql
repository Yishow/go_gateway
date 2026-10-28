-- Independent transfer edges may reference one verified canonical journal.
-- Keep all journal identities, accepted bytes and existing receipt rows intact.
DROP INDEX IF EXISTS idx_wg_handoff_target_seq;
CREATE INDEX IF NOT EXISTS idx_wg_handoff_target_seq
 ON wg_delivery_sample_handoffs(target_seq) WHERE target_seq IS NOT NULL;
