import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { isOriginalRequestRetry, recordingStartRecoveryBlocked } from '@/features/datalink/workbench-v2/steps/step4/recordingStartRecovery';
import { basicRecordingIntentKey, saveBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';
import { recordingTestWriteRecovery as store } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import { retainSchemaOperation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery';
import type { RecordingStartOperation, RecordingStartRequest } from '@/types/studioV2RecordingStart';

const ownerKey = 'ws-1/dev-1/circuit-water';
function request(): RecordingStartRequest {
  return { request_id: 'original-request', workspace_id: 'ws-1', expected_workspace_revision: 'ws-old', device_ids: ['dev-1'],
    groups: ['G1', 'G2'].map(group_id => ({ group_id, expected_group_revision: 'original-revision', expected_connector_revision: 'conn-rev' })) };
}
function operation(overrides: Partial<RecordingStartOperation> = {}): RecordingStartOperation {
  return { operation_id: 'O1', action: 'recording_start', status: 'partial', next_action: 'retry_same_request', workspace_id: 'ws-1',
    intent_digest: 'server-owned-digest', setup_revision: 'ws-after-partial-save', device_ids: ['dev-1'], devices: [{ device_id: 'dev-1', activated: false }],
    groups: ['G1', 'G2'].map(group_id => ({ group_id, group_revision: 'new-after-partial-save', saved: true, ready: false, applied: false })),
    stage: 'readiness', created_at: '2026-10-28T00:00:00Z', updated_at: '2026-10-28T00:00:01Z', ...overrides };
}
const blocked = () => recordingStartRecoveryBlocked(ownerKey, 'ws-1', ['G1', 'G2'], 'O1');
const schemaRefs: { workspace: string; group: string; id: string }[] = [];
function schema(workspace: string, group: string) {
  const id = `schema-${workspace}-${group}`; schemaRefs.push({ workspace, group, id }); retainSchemaOperation(workspace, group, id);
}
function testWrite(workspace = 'ws-1', group = 'G1') {
  store.retain({ workspace_id: workspace, group_id: group, group_revision: 'old-test-write-revision', operation_id: `test-${workspace}-${group}`,
    operation_kind: 'test_write', created_at: '2026-10-28T00:00:00Z' });
}
function resetStores() {
  localStorage.clear(); sessionStorage.clear();
  for (const { workspace, group, id } of schemaRefs.splice(0)) retainSchemaOperation(workspace, group, undefined, id);
  for (const ref of store.read().references) store.resolve({ operation_id: ref.operation_id, action: 'test_write', status: 'succeeded',
    write_outcome: 'written_verified', cleanup_status: 'cleaned', created_at: ref.created_at, updated_at: ref.created_at }, ref.workspace_id);
}
beforeEach(() => { vi.restoreAllMocks(); resetStores(); });
afterEach(() => { vi.restoreAllMocks(); resetStores(); });

describe('exact original request recovery eligibility', () => {
  it.each(['failed', 'partial'] as const)('accepts only explicit original-request recovery for %s despite legitimate new group revisions', status => {
    const original = request(); const before = structuredClone(original);
    expect(isOriginalRequestRetry(operation({ status }), 'O1', original)).toBe(true);
    expect(original).toEqual(before);
  });
  it('compares full identity sets, not order, current revisions, readiness tokens, or a fabricated intent hash', () => {
    const original = { ...request(), device_ids: ['dev-1', 'dev-2'], readiness_token: 'old-original-token' };
    const result = operation({ device_ids: ['dev-2', 'dev-1'], devices: [{ device_id: 'dev-2', activated: false }, { device_id: 'dev-1', activated: true }],
      groups: [...operation().groups].reverse() });
    expect(isOriginalRequestRetry(result, 'O1', original)).toBe(true);
  });
  it.each(['pending', 'running', 'unknown', 'succeeded'] as const)('rejects %s even if next_action says retry_same_request', status => {
    expect(isOriginalRequestRetry(operation({ status }), 'O1', request())).toBe(false);
  });
  it.each(['wrong ID', 'empty ID', 'wrong action', 'wrong next action', 'missing next action', 'foreign workspace', 'missing operation', 'missing request', 'empty request ID'] as const)('rejects %s', invalid => {
    let result: RecordingStartOperation | undefined = operation(); let original: RecordingStartRequest | undefined = request(); let id = 'O1';
    if (invalid === 'wrong ID') result.operation_id = 'O2';
    if (invalid === 'empty ID') id = '';
    if (invalid === 'wrong action') result = { ...result, action: 'schema_apply' } as unknown as RecordingStartOperation;
    if (invalid === 'wrong next action') result.next_action = 'revalidate';
    if (invalid === 'missing next action') delete result.next_action;
    if (invalid === 'foreign workspace') result.workspace_id = 'foreign';
    if (invalid === 'missing operation') result = undefined;
    if (invalid === 'empty request ID') original.request_id = '';
    if (invalid === 'missing request') original = undefined;
    expect(isOriginalRequestRetry(result, id, original)).toBe(false);
  });
  it.each(['duplicate operation devices', 'duplicate request devices', 'missing operation device', 'foreign operation device', 'foreign progress device',
    'duplicate progress devices', 'missing progress device', 'duplicate operation groups', 'duplicate request groups', 'missing operation group', 'foreign operation group', 'empty group identity'] as const)('rejects %s', invalid => {
    const result = operation(); const original = request();
    if (invalid === 'duplicate operation devices') result.device_ids.push('dev-1');
    if (invalid === 'duplicate request devices') original.device_ids.push('dev-1');
    if (invalid === 'missing operation device') result.device_ids = [];
    if (invalid === 'foreign operation device') result.device_ids = ['foreign'];
    if (invalid === 'foreign progress device') result.devices[0].device_id = 'foreign';
    if (invalid === 'duplicate progress devices') result.devices.push(result.devices[0]);
    if (invalid === 'missing progress device') result.devices = [];
    if (invalid === 'duplicate operation groups') result.groups.push(result.groups[0]);
    if (invalid === 'duplicate request groups') original.groups.push(original.groups[0]);
    if (invalid === 'missing operation group') result.groups.pop();
    if (invalid === 'foreign operation group') result.groups[0].group_id = 'foreign';
    if (invalid === 'empty group identity') result.groups[0].group_id = original.groups[0].group_id = '';
    expect(isOriginalRequestRetry(result, 'O1', original)).toBe(false);
  });
  it.each(['device_ids', 'devices', 'groups'] as const)('rejects a missing operation %s array instead of treating it as empty proof', field => {
    const result = { ...operation(), [field]: undefined } as unknown as RecordingStartOperation;
    expect(isOriginalRequestRetry(result, 'O1', request())).toBe(false);
  });
  it('allows original share-only requests with genuinely empty group sets without inventing group IDs', () => {
    expect(isOriginalRequestRetry(operation({ groups: [] }), 'O1', { ...request(), groups: [] })).toBe(true);
  });
});

describe('same-request recovery restrictive existing owner inspection', () => {
  it('permits healthy local-empty inspection only as a denial check, with no storage writes', () => {
    const write = vi.spyOn(Storage.prototype, 'setItem'); expect(blocked()).toBe(false); expect(write).not.toHaveBeenCalled();
  });
  it('permits retained matching O1 but blocks different or missing retained IDs', () => {
    saveBasicRecordingIntent({ scope_key: ownerKey, request: request(), operation_id: 'O1' }); expect(blocked()).toBe(false);
    saveBasicRecordingIntent({ scope_key: ownerKey, request: request(), operation_id: 'O2' }); expect(blocked()).toBe(true);
    saveBasicRecordingIntent({ scope_key: ownerKey, request: request() }); expect(blocked()).toBe(true);
  });
  it.each(['', 'foreign/dev-1/circuit-water', 'invalid'])('blocks missing or foreign owner scope %s', owner => {
    expect(recordingStartRecoveryBlocked(owner, 'ws-1', ['G1'], 'O1')).toBe(true);
  });
  it('blocks missing workspace, operation ID, and invalid group identities', () => {
    expect(recordingStartRecoveryBlocked(ownerKey, '', ['G1'], 'O1')).toBe(true);
    expect(recordingStartRecoveryBlocked(ownerKey, 'ws-1', ['G1'], undefined)).toBe(true);
    expect(recordingStartRecoveryBlocked(ownerKey, 'ws-1', ['G1', 'G1'], 'O1')).toBe(true);
    expect(recordingStartRecoveryBlocked(ownerKey, 'ws-1', [''], 'O1')).toBe(true);
  });
  it('blocks malformed or unreadable start storage and never removes its record', () => {
    sessionStorage.setItem(basicRecordingIntentKey(ownerKey), '{bad'); expect(blocked()).toBe(true);
    expect(sessionStorage.getItem(basicRecordingIntentKey(ownerKey))).toBe('{bad');
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => { throw new Error('private storage'); }); expect(blocked()).toBe(true);
  });
  it('blocks matching schema operation memory/storage and ignores unrelated schema groups/workspaces', () => {
    schema('ws-1', 'unrelated'); schema('foreign', 'G1'); expect(blocked()).toBe(false);
    schema('ws-1', 'G1'); expect(blocked()).toBe(true); sessionStorage.clear(); expect(blocked()).toBe(true);
  });
  it('blocks schema unreadability even when the start store remains readable and no schema ID is known', () => {
    const getItem = Storage.prototype.getItem;
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation(function (this: Storage, key: string) {
      if (key.startsWith('wbv2.schema-operation.')) throw new Error('private schema storage');
      return getItem.call(this, key);
    });
    expect(blocked()).toBe(true);
  });
  it('blocks unavailable test-write storage even with no known test-write ID', () => {
    localStorage.setItem('studio.recording-test-write-recovery.v1', '{bad');
    expect(store.read()).toMatchObject({ references: [], unavailable: true }); expect(blocked()).toBe(true);
  });
  it('blocks matching test-write O1 at its original revision and ignores unrelated groups/workspaces', () => {
    testWrite('ws-1', 'unrelated'); testWrite('foreign', 'G1'); expect(blocked()).toBe(false);
    testWrite(); expect(blocked()).toBe(true);
  });
  it('blocks unreadable test-write storage without erasing known locators', () => {
    testWrite(); localStorage.setItem('studio.recording-test-write-recovery.v1', '{bad'); expect(blocked()).toBe(true);
    expect(store.read().references).toHaveLength(1); expect(localStorage.getItem('studio.recording-test-write-recovery.v1')).toBe('{bad');
  });
});
