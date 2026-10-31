CREATE INDEX idx_wg_mixed_source_latest ON wg_delivery_mixed_outcomes (workspace_id, group_id, group_revision, due_at DESC, record_id DESC);
CREATE INDEX idx_wg_mixed_source_entity_latest ON wg_delivery_mixed_outcomes (workspace_id, group_id, group_revision, json_extract(outcome,'$.EntityKey'), due_at DESC, record_id DESC);
