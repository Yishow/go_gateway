-- Migration: 020_durable_delivery_sqlite.up.sql (SQLite)

CREATE TABLE IF NOT EXISTS gw_delivery_journal (
    sequence INTEGER PRIMARY KEY AUTOINCREMENT,
    record_id TEXT NOT NULL,
    observed_at DATETIME NOT NULL,
    payload BLOB NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_gw_delivery_journal_seq ON gw_delivery_journal(sequence);

CREATE TABLE IF NOT EXISTS gw_delivery_outbox (
    id TEXT PRIMARY KEY,
    destination_id TEXT NOT NULL,
    destination_revision TEXT NOT NULL DEFAULT '',
    plan_revision TEXT NOT NULL DEFAULT '',
    record_id TEXT NOT NULL,
    calculation_revision INTEGER NOT NULL DEFAULT 1,
    table_name TEXT NOT NULL,
    payload BLOB NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    retry_count INTEGER NOT NULL DEFAULT 0,
    next_retry_at DATETIME NOT NULL,
    last_error TEXT,
    observed_at DATETIME NOT NULL,
    delivered_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_gw_delivery_outbox_dest_rec 
    ON gw_delivery_outbox(destination_id, record_id);

CREATE INDEX IF NOT EXISTS idx_gw_delivery_outbox_status 
    ON gw_delivery_outbox(destination_id, status, next_retry_at);

CREATE TABLE IF NOT EXISTS gw_delivery_receipts (
    destination_id TEXT NOT NULL,
    record_id TEXT NOT NULL,
    calculation_revision INTEGER NOT NULL DEFAULT 1,
    table_name TEXT NOT NULL,
    delivered_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (destination_id, record_id, calculation_revision)
);
