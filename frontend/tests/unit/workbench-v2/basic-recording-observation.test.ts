import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import * as recording from '@/features/datalink/workbench-v2/state/basicRecording';

const scope = 'W1/D1/managed-recording';
const key = recording.basicRecordingIntentKey(scope);
const request = { request_id: 'request-1', workspace_id: 'W1', expected_workspace_revision: 'R1', device_ids: ['D1'], groups: [] };
const observe = (selected = scope) => recording.inspectBasicRecordingIntent(selected);

beforeEach(() => sessionStorage.clear());
afterEach(() => vi.restoreAllMocks());

it('reports a successful owner-local empty read, not global server absence', () => {
  expect(observe()).toEqual({ state: 'empty' });
});
it('keeps a retained request without an operation ID unresolved and exposes no request or token', () => {
  recording.saveBasicRecordingIntent({ scope_key: scope, request: { ...request, readiness_token: 'private-token' } });
  expect(observe()).toEqual({ state: 'retained' });
  expect(JSON.stringify(observe())).not.toContain('private-token');
});
it('exposes only the opaque original operation ID, not the persisted request', () => {
  recording.saveBasicRecordingIntent({ scope_key: scope, request, operation_id: 'O1' });
  expect(observe()).toEqual({ state: 'retained', operationId: 'O1' });
});
it.each(['{bad-json', 'null', JSON.stringify({ version: 1, scope_key: scope, request: { ...request, groups: null } })])(
  'does not turn invalid persisted content into idle: %s', raw => {
    sessionStorage.setItem(key, raw);
    expect(observe()).toEqual({ state: 'unavailable', reason: 'invalid_record' });
    expect(recording.loadBasicRecordingIntent(scope)).toBeUndefined();
  });
it('does not expose a foreign retained request or classify it as empty', () => {
  sessionStorage.setItem(key, JSON.stringify({ version: 1, scope_key: scope, request: { ...request, device_ids: ['D2'] } }));
  expect(observe()).toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
});
it('does not classify storage failures or missing scope as an empty owner', () => {
  vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => { throw new Error('private storage exception'); });
  expect(observe()).toEqual({ state: 'unavailable', reason: 'unreadable' });
  expect(observe('')).toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
});
