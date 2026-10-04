import test from 'node:test';
import assert from 'node:assert/strict';
import { execFileSync } from 'node:child_process';
import { mkdtempSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { readSQLiteDelivery } from './fresh-sqlite-observation.mjs';
import { assertRowsPreserved, readLedgerDetail } from './fresh-matrix-test-write.mjs';

const group = { id: 'g', applied_revision: 'r', row_policy: { entity_key_column: 'entity_id' } };
const row = { record_id: 'record-1', group_id: 'g', entity_id: 'A', bucket_start: '2026-10-04T00:00:00Z', value: '12' };
function runtime(rows) {
  const outbox = rows.map((entry) => ({ ...entry, group_revision: 'r', entity_key: entry.entity_id,
    table_name: 'readings', state: 'sql_committed', effect_key: `effect-${entry.record_id}`, payload_digest: `digest-${entry.record_id}` }));
  return { group, table: 'readings', outbox, receipts: outbox.map((entry) => ({
    effect_key: entry.effect_key, payload_digest: entry.payload_digest, committed_at: '2026-10-04T00:01:00Z' })) };
}

test('keeps all neighbor bytes while allowing another independently committed runtime row', () => {
  const after = [row, { ...row, record_id: 'record-2', bucket_start: '2026-10-04T00:01:00Z' }];
  assert.doesNotThrow(() => assertRowsPreserved([row], after, 'test-write', runtime(after)));
  assert.throws(() => assertRowsPreserved([row], [{ ...row, value: 'changed' }], 'test-write', runtime([row])));
});

test('rejects extra rows that cannot be proved to belong to current durable runtime delivery', () => {
  const after = [row, { ...row, record_id: 'alien', entity_id: '' }];
  assert.throws(() => assertRowsPreserved([row], after, 'test-write', runtime([row])));
  for (const change of ['revision', 'receipt']) {
    const proof = runtime(after);
    if (change === 'revision') proof.outbox[1].group_revision = 'other';
    else proof.receipts[1].payload_digest = 'other';
    assert.throws(() => assertRowsPreserved([row], after, 'test-write', proof), change);
  }
});

test('reads the service-owned ledger detail only from the local ledger, never from the API operation', () => {
  const preview = { operation_id: 'op-1', owner_column: 'entity_id', owner_value: 'gw-test-1', target: { table: 'readings' } };
  const apiOperation = { operation_id: 'op-1', action: 'test_write', write_outcome: 'written_verified', cleanup_status: 'cleaned' };
  const ledgerRow = { operation_id: 'op-1', action: 'test_write', detail: JSON.stringify({ owner_column: 'entity_id', owner_value: 'gw-test-1', table: 'readings' }) };
  assert.deepEqual(readLedgerDetail({ apiOperation, ledgerRows: [ledgerRow], preview }),
    { owner_column: 'entity_id', owner_value: 'gw-test-1', table: 'readings' });
  assert.throws(() => readLedgerDetail({ apiOperation: { ...apiOperation, detail: '{}' }, ledgerRows: [ledgerRow], preview }), /exposes detail/);
  assert.throws(() => readLedgerDetail({ apiOperation, ledgerRows: [], preview }), /exactly one/);
  assert.throws(() => readLedgerDetail({ apiOperation, ledgerRows: [{ ...ledgerRow, action: 'schema_apply' }], preview }), /test_write/);
  assert.throws(() => readLedgerDetail({ apiOperation, ledgerRows: [{ ...ledgerRow, detail: '{"owner_value":"other"}' }], preview }), /owner/);
});

test('SQLite delivery observation keeps the outbox entity key needed to prove runtime rows', () => {
  const dir = mkdtempSync(join(tmpdir(), 'gw-f-sqlite-delivery-'));
  try {
    const db = join(dir, 'gateway.db');
    execFileSync('sqlite3', [db, `
      CREATE TABLE wg_delivery_outbox (effect_key TEXT, record_id TEXT, group_id TEXT, group_revision TEXT, entity_key TEXT, bucket_start TEXT, table_name TEXT, state TEXT, payload_digest TEXT);
      CREATE TABLE wg_delivery_receipts (effect_key TEXT, payload_digest TEXT, committed_at TEXT);
      CREATE TABLE wg_delivery_buckets (group_id TEXT, group_revision TEXT, bucket_start TEXT, kind TEXT, reason TEXT, record_id TEXT, members TEXT, entity_key TEXT);
      INSERT INTO wg_delivery_outbox VALUES ('e1','r1','g','rev','B','2026-10-04T00:00:00Z','readings','sql_committed','d1');`]);
    const delivery = readSQLiteDelivery({ path: db, groupId: 'g', groupRevision: 'rev' });
    assert.equal(delivery.outbox[0].entity_key, 'B');
  } finally {
    rmSync(dir, { recursive: true, force: true });
  }
});
