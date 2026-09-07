-- Migration: 020_durable_delivery.down.sql

DROP TABLE IF EXISTS gw_delivery_receipts;
DROP TABLE IF EXISTS gw_delivery_outbox;
DROP TABLE IF EXISTS gw_delivery_journal;
