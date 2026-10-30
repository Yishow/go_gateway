import { beforeEach, describe, expect, it, vi } from 'vitest';
import { basicRecordingIntentKey, basicRecordingScopeKey, inspectBasicRecordingDevice } from '@/features/datalink/workbench-v2/state/basicRecording';
function retain(device = 'D1', role = 'opaque-owner-role') {
  const scope_key = basicRecordingScopeKey('W1', device, role);
  const value = { version: 1, scope_key, operation_id: 'O1', request: { request_id: 'request-1', workspace_id: 'W1', expected_workspace_revision: 'W2',
    device_ids: [device], groups: [{ group_id: 'G', expected_group_revision: 'G2', expected_connector_revision: 'DB2' }], readiness_token: 'PRIVATE-TOKEN' } };
  sessionStorage.setItem(basicRecordingIntentKey(scope_key), JSON.stringify(value)); return basicRecordingIntentKey(scope_key);
}
beforeEach(() => sessionStorage.clear());
describe('existing start store device metadata inspection', () => {
  it('uses exact stored owner keys, returns only scoped safe metadata and never writes/cleans storage', () => {
    const key = retain(); retain('D2'); const raw = sessionStorage.getItem(key);
    const result = inspectBasicRecordingDevice('W1', 'D1');
    expect(result).toEqual({ unavailable: false, retained: [{ operationId: 'O1', groupIds: ['G'] }] });
    expect(JSON.stringify(result)).not.toMatch(/PRIVATE|opaque-owner|scope_key|request/);
    expect(sessionStorage.getItem(key)).toBe(raw);
    expect(inspectBasicRecordingDevice('FOREIGN', 'D1').retained).toEqual([]);
  });
  it('invalid original request remains unavailable and is not erased', () => {
    const key = retain(); sessionStorage.setItem(key, '{broken');
    expect(inspectBasicRecordingDevice('W1', 'D1')).toEqual({ unavailable: true, retained: [] });
    expect(sessionStorage.getItem(key)).toBe('{broken');
  });
  it('unreadable storage stays unavailable, not a synthetic empty owner', () => {
    const spy = vi.spyOn(Storage.prototype, 'key').mockImplementation(() => { throw new Error('private'); });
    sessionStorage.setItem('other', 'value');
    expect(inspectBasicRecordingDevice('W1', 'D1')).toEqual({ unavailable: true, retained: [] }); spy.mockRestore();
  });
});
