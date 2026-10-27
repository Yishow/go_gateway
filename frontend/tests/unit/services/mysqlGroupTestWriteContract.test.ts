import { describe, expect, it } from 'vitest';
import { parseWriteGroupTestWritePreview } from '@/utils/studioV2WriteGroupTestWriteJson';

const recordIds = ['a'.repeat(64), 'b'.repeat(64), 'c'.repeat(64)];
function mysqlPreview(record_ids: unknown = recordIds) {
  return {
    token: 'tok-mysql', operation_id: 'op-mysql', action: 'test_write', group_id: 'group-1', group_revision: 'rev-1',
    expires_at: '2026-10-26T12:00:00Z',
    target: { connector_id: 'mysql-1', dialect: 'mysql', database: 'fixture', schema: 'fixture', table: 'gw_long_a' },
    owner_column: '_gw_owner_0123456789abcdef01234567', owner_value: 'gateway', dedupe: 'receipt',
    values: [
      { column: 'record_id', type: 'char(64)', value: recordIds[0] },
      { column: '_gw_owner_0123456789abcdef01234567', type: 'varchar(64)', value: 'gateway' },
    ],
    record_ids, cleanup: 'remove only the previewed record IDs and their matching receipts; retain all other rows',
  };
}

describe('MySQL long test-write preview', () => {
  it('retains every exact operation record ID separately from bounded explanatory text', () => {
    expect(parseWriteGroupTestWritePreview(mysqlPreview())?.record_ids).toEqual(recordIds);
  });

  it.each([undefined, [], ['bad'], [recordIds[0], recordIds[0]], Array.from({ length: 257 }, () => recordIds[0])])(
    'rejects missing, malformed, duplicate or unbounded record identities', (ids) => {
      const input = mysqlPreview();
      input.record_ids = ids;
      expect(parseWriteGroupTestWritePreview(input)).toBeNull();
    },
  );
});
