import test from 'node:test';
import assert from 'node:assert/strict';
import {
  assertAcquisitionProvenance,
  assertNeighborUnchanged,
  assertNoOwnedRows,
  assertOneEffect,
  assertOperationIdentity,
  assertTargetRowsUnchanged,
  assertTestWriteOutcome,
  sanitizeEvidence,
} from './fresh-recovery-assertions.mjs';

test('replay keeps the original operation identity', () => {
  const first = { operation_id: 'op-a', action: 'recording_start' };
  assert.doesNotThrow(() => assertOperationIdentity(first, { ...first, status: 'succeeded' }, 'replay'));
  assert.throws(() => assertOperationIdentity(first, { operation_id: 'op-b', action: 'recording_start' }, 'replay'), /changed operation identity/);
});

test('delayed resend checks original acquisition facts and one effect', () => {
  const expected = { observed_at: '2026-01-01T00:00:08.000000000Z' };
  const row = { bucket_start: '2026-01-01T00:00:08Z', provenance: JSON.stringify([{ observed_at: expected.observed_at }]) };
  assert.doesNotThrow(() => assertAcquisitionProvenance(row, expected));
  assert.doesNotThrow(() => assertOneEffect([{ effect_key: 'effect-a' }], 'delay'));
  assert.throws(() => assertAcquisitionProvenance({ ...row, bucket_start: '2026-01-01T00:00:18Z' }, expected), /acquisition time changed/);
  assert.throws(() => assertOneEffect([{ effect_key: 'effect-a' }, { effect_key: 'effect-a' }], 'delay'), /duplicate/);
});

test('test-write cleanup preserves neighbor and removes owned row', () => {
  const before = [{ line: 'neighbor', value: 7 }];
  assert.doesNotThrow(() => assertNeighborUnchanged(before, [{ line: 'neighbor', value: 7 }]));
  assert.doesNotThrow(() => assertTargetRowsUnchanged([{ record_id: 'r-1', value: 7 }], [{ record_id: 'r-1', value: 7 }, { record_id: 'gw-test-1', value: 9 }], 'gw-test-1', 'record_id'));
  assert.doesNotThrow(() => assertNoOwnedRows([{ line: 'neighbor' }], 'gw-test-1'));
  assert.doesNotThrow(() => assertTestWriteOutcome({ operation_id: 'op-test', owner_value: 'gw-test-1', write_outcome: 'written_verified', cleanup_status: 'cleaned' }));
  assert.throws(() => assertNeighborUnchanged(before, [{ line: 'neighbor', value: 8 }]), /neighbor changed/);
});

test('sanitized evidence removes credentials and private paths', () => {
  const result = sanitizeEvidence({ password: 'secret', dsn: 'postgres://private', nested: { operation_id: 'op-a' } });
  assert.equal(result.password, '[redacted]');
  assert.equal(result.dsn, '[redacted]');
  assert.equal(result.nested.operation_id, 'op-a');
});
