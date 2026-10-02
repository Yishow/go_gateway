-- Test-write operations share the schema operation ledger: they hold the same
-- scope exclusion, claim and lease, and add the result facts they must keep.
ALTER TABLE managed_schema_operations ADD COLUMN payload_digest TEXT NOT NULL DEFAULT '';
ALTER TABLE managed_schema_operations ADD COLUMN write_outcome TEXT NOT NULL DEFAULT '';
ALTER TABLE managed_schema_operations ADD COLUMN cleanup_status TEXT NOT NULL DEFAULT '';
ALTER TABLE managed_schema_operations ADD COLUMN cleanup_reason TEXT NOT NULL DEFAULT '';
ALTER TABLE managed_schema_operations ADD COLUMN detail TEXT NOT NULL DEFAULT '';
