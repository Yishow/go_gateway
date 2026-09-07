-- Migration: 020_durable_delivery.up.sql (PostgreSQL / Generic)

CREATE TABLE IF NOT EXISTS gw_delivery_journal (
    sequence BIGSERIAL PRIMARY KEY,
    record_id VARCHAR(128) NOT NULL,
    observed_at TIMESTAMPTZ NOT NULL,
    payload BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_gw_delivery_journal_seq ON gw_delivery_journal(sequence);

CREATE TABLE IF NOT EXISTS gw_delivery_outbox (
    id VARCHAR(64) PRIMARY KEY,
    destination_id VARCHAR(64) NOT NULL,
    destination_revision VARCHAR(64) NOT NULL DEFAULT '',
    plan_revision VARCHAR(64) NOT NULL DEFAULT '',
    record_id VARCHAR(128) NOT NULL,
    calculation_revision BIGINT NOT NULL DEFAULT 1,
    table_name VARCHAR(128) NOT NULL,
    payload BYTEA NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    retry_count INT NOT NULL DEFAULT 0,
    next_retry_at TIMESTAMPTZ NOT NULL,
    last_error TEXT,
    observed_at TIMESTAMPTZ NOT NULL,
    delivered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_gw_delivery_outbox_dest_rec 
    ON gw_delivery_outbox(destination_id, record_id);

CREATE INDEX IF NOT EXISTS idx_gw_delivery_outbox_status 
    ON gw_delivery_outbox(destination_id, status, next_retry_at);

CREATE TABLE IF NOT EXISTS gw_delivery_receipts (
    destination_id VARCHAR(64) NOT NULL,
    record_id VARCHAR(128) NOT NULL,
    calculation_revision BIGINT NOT NULL DEFAULT 1,
    table_name VARCHAR(128) NOT NULL,
    delivered_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (destination_id, record_id, calculation_revision)
);
