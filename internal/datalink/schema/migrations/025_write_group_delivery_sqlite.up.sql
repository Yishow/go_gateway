-- Durable write-group delivery: sample journal, closure checkpoint, bucket
-- outcomes, row outbox and local receipts. Timestamps are UTC RFC3339Nano text
-- so ordering and replay stay deterministic.

-- A sample is ACKed only after its row here is committed.
CREATE TABLE IF NOT EXISTS wg_delivery_samples (
    seq             INTEGER PRIMARY KEY AUTOINCREMENT,
    workspace_id    TEXT NOT NULL,
    group_id        TEXT NOT NULL,
    group_revision  TEXT NOT NULL,
    sample_id       TEXT NOT NULL,
    member_key      TEXT NOT NULL,
    observed_at     TEXT NOT NULL,
    bucket_start    TEXT NOT NULL,
    payload         TEXT NOT NULL,
    payload_digest  TEXT NOT NULL,
    consumed        INTEGER NOT NULL DEFAULT 0,
    created_at      TEXT NOT NULL,
    UNIQUE (group_id, group_revision, sample_id)
);

CREATE INDEX IF NOT EXISTS idx_wg_delivery_samples_open
    ON wg_delivery_samples(group_id, group_revision, consumed, seq);

-- next_close is the first bucket start that has not closed yet.
CREATE TABLE IF NOT EXISTS wg_delivery_checkpoints (
    group_id        TEXT NOT NULL,
    group_revision  TEXT NOT NULL,
    next_close      TEXT NOT NULL,
    updated_at      TEXT NOT NULL,
    PRIMARY KEY (group_id, group_revision)
);

-- One scoped outcome per entity and closed bucket (row, skipped or no_data).
CREATE TABLE IF NOT EXISTS wg_delivery_buckets (
    group_id        TEXT NOT NULL,
    group_revision  TEXT NOT NULL,
    entity_key      TEXT NOT NULL,
    bucket_start    TEXT NOT NULL,
    kind            TEXT NOT NULL,
    reason          TEXT NOT NULL DEFAULT '',
    record_id       TEXT NOT NULL DEFAULT '',
    members         TEXT NOT NULL DEFAULT '[]',
    created_at      TEXT NOT NULL,
    PRIMARY KEY (group_id, group_revision, entity_key, bucket_start)
);

-- Closed rows frozen with their destination and revisions. effect_key is the
-- destination-scoped identity used to deduplicate delivery.
CREATE TABLE IF NOT EXISTS wg_delivery_outbox (
    effect_key          TEXT PRIMARY KEY,
    record_id           TEXT NOT NULL,
    workspace_id        TEXT NOT NULL,
    group_id            TEXT NOT NULL,
    group_revision      TEXT NOT NULL,
    entity_key          TEXT NOT NULL DEFAULT '',
    bucket_start        TEXT NOT NULL,
    partition_key       TEXT NOT NULL,
    destination_scope   TEXT NOT NULL,
    connector_id        TEXT NOT NULL,
    connector_revision  TEXT NOT NULL,
    database_name       TEXT NOT NULL DEFAULT '',
    table_schema        TEXT NOT NULL DEFAULT '',
    table_name          TEXT NOT NULL,
    dedupe_capability   TEXT NOT NULL DEFAULT 'none',
    record_key_column   TEXT NOT NULL DEFAULT '',
    payload             TEXT NOT NULL,
    payload_digest      TEXT NOT NULL,
    state               TEXT NOT NULL DEFAULT 'pending',
    retry_count         INTEGER NOT NULL DEFAULT 0,
    next_retry_at       TEXT NOT NULL,
    last_error_code     TEXT NOT NULL DEFAULT '',
    claim_owner         TEXT NOT NULL DEFAULT '',
    claim_expires_at    TEXT NOT NULL DEFAULT '',
    claim_epoch         INTEGER NOT NULL DEFAULT 0,
    committed_at        TEXT NOT NULL DEFAULT '',
    created_at          TEXT NOT NULL,
    updated_at          TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_wg_delivery_outbox_partition
    ON wg_delivery_outbox(partition_key, bucket_start, state);

CREATE INDEX IF NOT EXISTS idx_wg_delivery_outbox_state
    ON wg_delivery_outbox(state, next_retry_at);

-- Local evidence that the destination committed an effect. It never proves a
-- commit by itself: it is written only after destination-side confirmation.
CREATE TABLE IF NOT EXISTS wg_delivery_receipts (
    effect_key      TEXT PRIMARY KEY,
    payload_digest  TEXT NOT NULL,
    committed_at    TEXT NOT NULL
);
