-- SQLite rollback: revert points uniqueness to (device_id, address)
-- Precheck before running:
-- SELECT device_id, address, COUNT(*) AS cnt
-- FROM points
-- GROUP BY device_id, address
-- HAVING COUNT(*) > 1;

PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;

CREATE TABLE points_v10 (
    id              TEXT PRIMARY KEY,
    device_id       TEXT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    description     TEXT,
    address         TEXT NOT NULL,
    function        TEXT,
    data_type       TEXT NOT NULL CHECK (data_type IN ('bool', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'string')),
    mode            TEXT NOT NULL DEFAULT 'read' CHECK (mode IN ('read', 'readwrite')),
    polling_group_id TEXT REFERENCES polling_groups(id) ON DELETE SET NULL,
    last_read_at    TEXT,
    last_value      TEXT,
    last_error      TEXT,
    enabled         INTEGER NOT NULL DEFAULT 1,
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE (device_id, address)
);

INSERT INTO points_v10 (
    id, device_id, name, description, address, function, data_type, mode,
    polling_group_id, last_read_at, last_value, last_error, enabled, created_at, updated_at
)
SELECT
    id, device_id, name, description, address, NULLIF(function, ''), data_type, mode,
    polling_group_id, last_read_at, last_value, last_error, enabled, created_at, updated_at
FROM points;

DROP TABLE points;
ALTER TABLE points_v10 RENAME TO points;

CREATE INDEX IF NOT EXISTS idx_points_device_id ON points(device_id);
CREATE INDEX IF NOT EXISTS idx_points_polling_group_id ON points(polling_group_id);
CREATE INDEX IF NOT EXISTS idx_points_enabled ON points(enabled);
CREATE INDEX IF NOT EXISTS idx_points_data_type ON points(data_type);

DROP TRIGGER IF EXISTS trigger_points_updated_at;
CREATE TRIGGER IF NOT EXISTS trigger_points_updated_at
    AFTER UPDATE ON points
    FOR EACH ROW
BEGIN
    UPDATE points SET updated_at = datetime('now') WHERE id = NEW.id;
END;

COMMIT;
PRAGMA foreign_keys=ON;
