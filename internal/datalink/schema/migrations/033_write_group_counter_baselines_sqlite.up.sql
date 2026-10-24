-- Durable cumulative counter baselines. One row per group revision and member
-- holds the typed checkpoint (strict JSON, digits as strings) of the last
-- reliable reading, committed in the same local transaction that closes the
-- bucket, writes the outbox rows and advances the closure checkpoint. A baseline
-- moves forward only: window_end orders it, and an equal window with different
-- content is a conflict, never an overwrite. Timestamps are fixed-width UTC text.
-- The table is not reclaimed with delivered rows: a baseline must outlive the
-- journal and outbox entries it was derived from.
CREATE TABLE IF NOT EXISTS wg_delivery_counter_baselines (
    workspace_id    TEXT NOT NULL,
    group_id        TEXT NOT NULL,
    group_revision  TEXT NOT NULL,
    member_key      TEXT NOT NULL,
    binding_key     TEXT NOT NULL,
    window_end      TEXT NOT NULL,
    payload         TEXT NOT NULL,
    payload_digest  TEXT NOT NULL,
    created_at      TEXT NOT NULL,
    updated_at      TEXT NOT NULL,
    PRIMARY KEY (group_id, group_revision, member_key)
);
