CREATE TABLE IF NOT EXISTS source_rule_candidate_snapshots (
    source_rule_id  TEXT NOT NULL REFERENCES source_rules(id) ON DELETE CASCADE,
    revision_id     TEXT NOT NULL,
    candidate_type  TEXT NOT NULL CHECK (candidate_type IN ('tags', 'database_outputs', 'local_modbus_outputs')),
    payload         TEXT NOT NULL,
    status          TEXT NOT NULL CHECK (status IN ('ready', 'blocked', 'deferred')),
    reason          TEXT NOT NULL DEFAULT '',
    generated_at    TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (source_rule_id, revision_id, candidate_type)
);

CREATE INDEX IF NOT EXISTS idx_source_rule_candidate_snapshots_rule_id
    ON source_rule_candidate_snapshots(source_rule_id);

CREATE INDEX IF NOT EXISTS idx_source_rule_candidate_snapshots_revision_id
    ON source_rule_candidate_snapshots(revision_id);
