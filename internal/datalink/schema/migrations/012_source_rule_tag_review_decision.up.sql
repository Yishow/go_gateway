CREATE TABLE IF NOT EXISTS source_rule_tag_review_decisions (
    source_rule_id   TEXT NOT NULL REFERENCES source_rules(id) ON DELETE CASCADE,
    candidate_id     TEXT NOT NULL,
    decision_type    TEXT NOT NULL CHECK (decision_type IN ('rename', 'skip', 'override')),
    tag_key          TEXT NOT NULL DEFAULT '',
    override_tag_id  TEXT REFERENCES tags(id) ON DELETE SET NULL,
    created_at       TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (source_rule_id, candidate_id)
);

CREATE INDEX IF NOT EXISTS idx_source_rule_tag_review_decisions_rule_id
    ON source_rule_tag_review_decisions(source_rule_id);
