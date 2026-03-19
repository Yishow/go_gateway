-- -----------------------------------------------------------------------------
-- 來源規則表 (source_rules)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS source_rules (
    id                TEXT PRIMARY KEY,
    device_id         TEXT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    start_address     TEXT NOT NULL,
    count             INTEGER NOT NULL CHECK (count > 0),
    data_type         TEXT NOT NULL CHECK (data_type IN ('bool', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'string')),
    naming_prefix     TEXT NOT NULL,
    enabled           INTEGER NOT NULL DEFAULT 1,
    locked            INTEGER NOT NULL DEFAULT 0,
    origin            TEXT NOT NULL DEFAULT 'manual' CHECK (origin IN ('manual', 'template')),
    template_name     TEXT,
    skipped_addresses TEXT NOT NULL DEFAULT '[]',
    created_at        TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at        TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_source_rules_device_id ON source_rules(device_id);
CREATE INDEX IF NOT EXISTS idx_source_rules_enabled ON source_rules(enabled);

-- -----------------------------------------------------------------------------
-- 來源規則衍生關聯表 (source_rule_links)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS source_rule_links (
    id                TEXT PRIMARY KEY,
    rule_id           TEXT NOT NULL REFERENCES source_rules(id) ON DELETE CASCADE,
    address           TEXT NOT NULL,
    point_id          TEXT NOT NULL REFERENCES points(id) ON DELETE CASCADE,
    tag_id            TEXT REFERENCES tags(id) ON DELETE SET NULL,
    mapping_id        TEXT REFERENCES mappings(id) ON DELETE SET NULL,
    created_at        TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at        TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE (rule_id, address),
    UNIQUE (point_id)
);

CREATE INDEX IF NOT EXISTS idx_source_rule_links_rule_id ON source_rule_links(rule_id);
CREATE INDEX IF NOT EXISTS idx_source_rule_links_point_id ON source_rule_links(point_id);
CREATE INDEX IF NOT EXISTS idx_source_rule_links_tag_id ON source_rule_links(tag_id);
CREATE INDEX IF NOT EXISTS idx_source_rule_links_mapping_id ON source_rule_links(mapping_id);

CREATE TRIGGER IF NOT EXISTS trigger_source_rules_updated_at
    AFTER UPDATE ON source_rules
    FOR EACH ROW
BEGIN
    UPDATE source_rules SET updated_at = datetime('now') WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS trigger_source_rule_links_updated_at
    AFTER UPDATE ON source_rule_links
    FOR EACH ROW
BEGIN
    UPDATE source_rule_links SET updated_at = datetime('now') WHERE id = NEW.id;
END;
