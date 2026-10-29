import { describe, expect, it } from 'vitest';
import { parseRecordingTestWriteRecovery } from '@/utils/recordingTestWriteRecoveryJson';
import { parseWriteGroupTestWriteOperation } from '@/utils/studioV2WriteGroupTestWriteJson';
const operation = {
  operation_id: 'O1', workspace_id: 'W1', workspace_revision: 'WR1', group_id: 'G1', group_revision: 'R1',
  connector_id: 'C1', connector_revision: 'CR1', target_fingerprint: 'a'.repeat(64), action: 'test_write', status: 'running',
  created_at: '2026-10-28T00:00:00Z', updated_at: '2026-10-28T00:00:02Z',
};
const envelope = (operations: unknown[]) => ({ success: true, data: { workspace_id: 'W1', complete: true, operations } });
describe('complete workspace-scoped test-write ledger parser', () => {
  it('accepts empty and the complete 1000-entry limit without leaking unexpected fields', () => {
    expect(parseRecordingTestWriteRecovery(envelope([]), 'W1').operations).toEqual([]);
    const rows = Array.from({ length: 1000 }, (_, index) => ({ ...operation, operation_id: `O${index}`, token: 'SECRET', reason: 'PRIVATE' }));
    const parsed = parseRecordingTestWriteRecovery(envelope(rows), 'W1');
    expect(parsed.operations).toHaveLength(1000);
    expect(JSON.stringify(parsed)).not.toMatch(/SECRET|PRIVATE|token|reason/);
    expect(() => parseRecordingTestWriteRecovery(envelope([...rows, { ...operation, operation_id: 'overflow' }]), 'W1')).toThrow();
  });
  it.each([
    { operation_id: '' }, { group_revision: ' ' }, { workspace_revision: 3 }, { connector_id: null },
    { connector_revision: 'x'.repeat(257) }, { target_fingerprint: 'A'.repeat(64) }, { action: 'schema_apply' },
    { created_at: '2026-02-30T00:00:00Z' }, { updated_at: '2026-10-27T00:00:00Z' },
    { updated_at: '2026-10-28T25:00:00Z' }, { updated_at: '2026-10-28T00:00:00+24:00' },
    { write_outcome: '' }, { cleanup_status: '' }, { completed_at: '' }, { completed_at: null },
    { status: 'pending', write_outcome: 'written_verified', cleanup_status: 'failed', completed_at: '2026-10-28T00:00:01Z' },
    { status: 'succeeded', write_outcome: 'written_verified', cleanup_status: 'cleaned', completed_at: '2026-10-28T00:00:01Z' },
    { status: 'partial', write_outcome: 'written_unverified', cleanup_status: 'cleaned', completed_at: '2026-10-28T00:00:01Z' },
    { status: 'failed', write_outcome: 'written_verified', cleanup_status: 'failed', completed_at: '2026-10-28T00:00:01Z' },
  ])('rejects malformed or contradictory operation %j', (patch) => {
    expect(() => parseRecordingTestWriteRecovery(envelope([{ ...operation, ...patch }]), 'W1')).toThrow();
  });
  it.each(['failed', 'unknown'])('accepts only a verified running cleanup claim with %s cleanup', (cleanup_status) => {
    const claim = { ...operation, write_outcome: 'written_verified', cleanup_status, completed_at: '2026-10-28T00:00:01Z' };
    expect(parseRecordingTestWriteRecovery(envelope([claim]), 'W1').operations[0].status).toBe('running');
    expect(parseWriteGroupTestWriteOperation(claim)?.status).toBe('running');
    for (const patch of [{ write_outcome: 'written_unverified' }, { cleanup_status: 'cleaned' }, { completed_at: undefined }]) {
      expect(parseWriteGroupTestWriteOperation({ ...claim, ...patch })).toBeNull();
    }
  });
  it('rejects invalid timestamps in the newly accepted exact-GET cleanup claim', () => {
    const claim = { ...operation, write_outcome: 'written_verified', cleanup_status: 'failed', completed_at: 'invalid' };
    expect(parseWriteGroupTestWriteOperation(claim)).toBeNull();
  });
});
