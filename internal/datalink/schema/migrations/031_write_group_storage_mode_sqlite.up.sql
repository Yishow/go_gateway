-- Additive, nullable storage mode of a write group. NULL means a group saved
-- before long storage existed (wide layout); only managed_long_v1 groups carry
-- a mode. The migrator executes this file one statement per line and skips an
-- ADD COLUMN whose column already exists, so keep one statement per line.
ALTER TABLE write_groups ADD COLUMN storage_mode TEXT;
