import { act, cleanup, renderHook } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { useState } from 'react';
import { basicRecordingIntentKey, completeBasicRecordingIntent, clearBasicRecordingIntentIfCurrent, loadBasicRecordingIntent,
  saveBasicRecordingIntent, subscribeBasicRecordingIntent, type BasicRecordingIntentRecord } from '@/features/datalink/workbench-v2/state/basicRecording';
import { useBasicRecordingIntentObservation } from '@/features/datalink/workbench-v2/steps/step4/useBasicRecordingIntentObservation';

const original = (): BasicRecordingIntentRecord => ({ scope_key: 'workspace-A/device-A/recording', request: { request_id: 'request-O1', workspace_id: 'workspace-A',
  expected_workspace_revision: 'old-head', device_ids: ['device-A'], groups: [{ group_id: 'group-A', expected_group_revision: 'old-group',
    expected_connector_revision: 'old-connector', draft: { workspace_id: 'workspace-A', name: 'Original draft', members: [],
      destination: { connector_id: 'connector-A', connector_revision: 'old-connector', table_schema: '', table_name: '', storage_strategy: 'managed' },
      row_policy: { interval_seconds: 30, allowed_lateness_seconds: 0, incomplete_policy: 'partial' }, write_policy: {} } }],
  readiness_token: 'private-original-token', settings_revision: 'settings-A', workspace_revision: 'share-A' } });
beforeEach(() => sessionStorage.clear());
afterEach(() => { cleanup(); vi.restoreAllMocks(); });

it('completes the same immutable original request without rebuilding revisions, draft or token', () => {
  const record = original(); saveBasicRecordingIntent(record);
  expect(completeBasicRecordingIntent(record, 'O1')).toBe('saved');
  expect(loadBasicRecordingIntent(record.scope_key)).toEqual({ ...record, version: 1, operation_id: 'O1' });
});
it.each(['request', 'draft', 'token', 'operation'] as const)('never overwrites a retained different %s, even with the same scope', changed => {
  const record = original(); const newer = structuredClone(record);
  if (changed === 'request') newer.request.request_id = 'O2-request';
  if (changed === 'draft' && newer.request.groups[0].draft) newer.request.groups[0].draft.row_policy.interval_seconds = 45;
  if (changed === 'token') newer.request.readiness_token = 'different-token';
  if (changed === 'operation') newer.operation_id = 'O2';
  saveBasicRecordingIntent(newer);
  expect(completeBasicRecordingIntent(record, 'O1')).toBe('superseded');
  expect(loadBasicRecordingIntent(record.scope_key)).toEqual({ ...newer, version: 1 });
  clearBasicRecordingIntentIfCurrent(record);
  expect(loadBasicRecordingIntent(record.scope_key)).toEqual({ ...newer, version: 1 });
});
it('does not resurrect a cleared request or erase a newer request after a late definitive failure', () => {
  const record = original();
  expect(completeBasicRecordingIntent(record, 'O1')).toBe('superseded');
  saveBasicRecordingIntent(record); clearBasicRecordingIntentIfCurrent(record);
  expect(loadBasicRecordingIntent(record.scope_key)).toBeUndefined();
});
it.each(['malformed', 'foreign', 'unreadable', 'write-failed'] as const)('leaves %s storage fail-closed and publishes no completion', failure => {
  const record = original(); saveBasicRecordingIntent(record);
  const heard = vi.fn(); const unsubscribe = subscribeBasicRecordingIntent(record.scope_key, heard);
  if (failure === 'malformed') sessionStorage.setItem(basicRecordingIntentKey(record.scope_key), '{bad');
  if (failure === 'foreign') sessionStorage.setItem(basicRecordingIntentKey(record.scope_key), JSON.stringify({ ...record, version: 1, scope_key: 'foreign/device/role' }));
  if (failure === 'unreadable') vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => { throw new Error('private failure'); });
  if (failure === 'write-failed') vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => { throw new Error('private failure'); });
  expect(completeBasicRecordingIntent(record, 'O1')).toBe('unavailable'); expect(heard).not.toHaveBeenCalled(); unsubscribe();
});
it('publishes only the scope, isolates other owners and removes subscribed listeners', () => {
  const record = original(); const a = vi.fn(), b = vi.fn();
  const stopA = subscribeBasicRecordingIntent(record.scope_key, a), stopB = subscribeBasicRecordingIntent('workspace-A/device-B/recording', b);
  saveBasicRecordingIntent(record); expect(a).toHaveBeenCalledWith(); expect(b).not.toHaveBeenCalled();
  stopA(); saveBasicRecordingIntent({ ...record, operation_id: 'O1' }); expect(a).toHaveBeenCalledOnce(); stopB();
});
it('adopts O1 once into the current original owner, ignoring duplicate publications and cleaning up on unmount', () => {
  const record = original(); saveBasicRecordingIntent(record); const observed = vi.fn();
  const view = renderHook(() => {
    const [current, setCurrent] = useState(record);
    useBasicRecordingIntentObservation(record.scope_key, current, value => { observed(value); setCurrent(value); });
    return current;
  });
  act(() => { completeBasicRecordingIntent(record, 'O1'); completeBasicRecordingIntent(record, 'O1'); });
  expect(view.result.current.operation_id).toBe('O1'); expect(observed).toHaveBeenCalledOnce();
  act(() => saveBasicRecordingIntent({ ...record, operation_id: 'O1' })); expect(observed).toHaveBeenCalledOnce();
  view.unmount(); act(() => saveBasicRecordingIntent({ ...record, operation_id: 'O1' })); expect(observed).toHaveBeenCalledOnce();
});
it('does not adopt another scope or O2 and ignores an old A completion while viewing B', () => {
  const record = original(); saveBasicRecordingIntent(record); const observed = vi.fn();
  const view = renderHook(({ current }) => useBasicRecordingIntentObservation(current.scope_key, current, observed), { initialProps: { current: record } });
  act(() => saveBasicRecordingIntent({ ...record, request: { ...record.request, request_id: 'O2' }, operation_id: 'O2' }));
  expect(observed).not.toHaveBeenCalled();
  const b = { ...record, scope_key: 'workspace-A/device-B/recording', request: { ...record.request, device_ids: ['device-B'] } };
  view.rerender({ current: b }); act(() => saveBasicRecordingIntent({ ...record, operation_id: 'O1' }));
  expect(observed).not.toHaveBeenCalled();
});
