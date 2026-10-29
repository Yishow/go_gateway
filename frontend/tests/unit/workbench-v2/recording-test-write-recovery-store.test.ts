import { describe, expect, it, vi } from 'vitest';
import { createRecordingTestWriteRecovery, recordingTestWriteRecoveryKey, testWriteResolved } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import type { WriteGroupTestWriteOperation } from '@/types/studioV2WriteGroupTestWrite';
const reference = { workspace_id: 'W1', group_id: 'G1', group_revision: 'R1', operation_id: 'O1', operation_kind: 'test_write' as const, created_at: '2026-10-28T00:00:00Z' };
function setup() {
  const data = new Map<string, string>();
  const storage = { getItem: vi.fn((key: string) => data.get(key) ?? null), setItem: vi.fn((key: string, value: string) => { data.set(key, value); }), removeItem: vi.fn(), clear: vi.fn(), key: vi.fn(), length: 0 };
  return { data, storage, store: createRecordingTestWriteRecovery(() => storage) };
}
const completed: WriteGroupTestWriteOperation = { operation_id: 'O1', action: 'test_write', status: 'succeeded', write_outcome: 'written_verified', cleanup_status: 'cleaned', created_at: '', updated_at: '' };
describe('nonsecret recovery locator storage (server remains authority)', () => {
  it('stores only allowed fields and retains independently scoped references', () => {
    const { data, store } = setup();
    store.retain({ ...reference, token: 'SECRET', target: 'PRIVATE' } as typeof reference);
    store.retain({ ...reference, workspace_id: 'W2', group_id: 'G2', operation_id: 'O2' });
    expect(JSON.parse(data.get(recordingTestWriteRecoveryKey)!)).toEqual([
      reference, { ...reference, workspace_id: 'W2', group_id: 'G2', operation_id: 'O2' },
    ]);
    store.resolve(completed);
    expect(store.read().references.map((ref) => ref.operation_id)).toEqual(['O2']);
  });
  it.each(['unknown', 'failed'] as const)('does not clear %s cleanup or unknown write', (cleanup_status) => {
    const { store } = setup(); store.retain(reference);
    store.resolve({ ...completed, cleanup_status });
    expect(store.read().references).toEqual([reference]);
    store.resolve({ ...completed, write_outcome: 'unknown' });
    expect(store.read().references).toEqual([reference]);
  });
  it('keeps the reference when persistence fails after dispatch and permits a later authoritative resolution', () => {
    const { store, storage } = setup(); store.retain(reference);
    storage.setItem.mockImplementationOnce(() => { throw new Error('quota'); });
    store.resolve(completed);
    expect(store.read().references).toEqual([reference]);
    store.resolve(completed);
    expect(store.read().references).toEqual([]);
  });
  it('preserves the in-memory reference when reads fail and never rewrites corrupt storage', () => {
    const { store, storage, data } = setup();
    storage.setItem.mockImplementationOnce(() => { throw new Error('quota'); });
    expect(store.retain(reference)).toBe(false);
    storage.getItem.mockImplementationOnce(() => { throw new Error('security'); });
    expect(store.read()).toMatchObject({ unavailable: true, references: [reference] });
    data.set(recordingTestWriteRecoveryKey, '{invalid');
    expect(store.retain({ ...reference, operation_id: 'O2' })).toBe(false);
    expect(data.get(recordingTestWriteRecoveryKey)).toBe('{invalid');
  });
  it.each(['{}', '[null]', '[{"operation_kind":"other"}]'])('rejects malformed references: %s', (raw) => {
    const { data, store } = setup(); data.set(recordingTestWriteRecoveryKey, raw);
    expect(store.read()).toMatchObject({ unavailable: true, references: [] });
  });
  it('requires explicit safe completion evidence', () => {
    expect(testWriteResolved({ ...completed, status: 'running' })).toBe(false);
    expect(testWriteResolved({ ...completed, write_outcome: undefined })).toBe(false);
    expect(testWriteResolved({ ...completed, write_outcome: 'failed', cleanup_status: 'not_attempted' })).toBe(true);
  });
});
it('does not let a completion read in W1 erase a locator owned by W2', () => {
  const { store } = setup();
  store.retain({ ...reference, workspace_id: 'W2' });
  store.resolve(completed, 'W1');
  expect(store.read().references).toEqual([{ ...reference, workspace_id: 'W2' }]);
});
it('treats conflicting duplicate browser locators as corruption, not absence', () => {
  const { store, data } = setup();
  data.set(recordingTestWriteRecoveryKey, JSON.stringify([reference, { ...reference, workspace_id: 'W2' }]));
  expect(store.read().unavailable).toBe(true);
});
it('preserves the last readable locator snapshot if storage becomes inaccessible after dispatch', () => {
  const { store, storage } = setup();
  expect(store.retain(reference)).toBe(true);
  storage.getItem.mockImplementation(() => { throw new Error('access revoked'); });
  expect(store.read()).toMatchObject({ unavailable: true, references: [reference] });
});
it.each(['null', 'empty', 'omitted'] as const)('retains known unresolved O1 after successful storage %s', (mode) => {
  const { store, data } = setup(); store.retain(reference);
  const foreign = { ...reference, workspace_id: 'W2', group_id: 'G2', operation_id: 'O2' };
  if (mode === 'null') data.delete(recordingTestWriteRecoveryKey);
  else data.set(recordingTestWriteRecoveryKey, JSON.stringify(mode === 'empty' ? [] : [foreign]));
  expect(store.read().references).toContainEqual(reference);
  if (mode === 'omitted') expect(store.read().references).toContainEqual(foreign);
  store.resolve(completed, 'W1');
  expect(store.read().references).not.toContainEqual(reference);
  if (mode === 'omitted') expect(store.read().references).toContainEqual(foreign);
});
it('failed explicit resolution after storage loss retains O1 until a later successful exact resolution', () => {
  const { store, data, storage } = setup(); store.retain(reference); data.clear();
  storage.setItem.mockImplementationOnce(() => { throw new Error('quota'); });
  store.resolve(completed, 'W1');
  expect(store.read().references).toEqual([reference]);
  store.resolve(completed, 'W1');
  expect(store.read().references).toEqual([]);
  expect(store.read().references).toEqual([]);
});
