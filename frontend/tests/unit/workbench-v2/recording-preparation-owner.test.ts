import { beforeEach, expect, it, vi } from 'vitest';
import { recordingStartObservation } from '@/features/datalink/workbench-v2/steps/step4/useRecordingPreparationOwner';
import { saveBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';
const scope = 'W1/D1/recording';
const request = { request_id: 'R1', workspace_id: 'W1', expected_workspace_revision: 'rev-1', device_ids: ['D1'], groups: [] };
beforeEach(() => { vi.restoreAllMocks(); sessionStorage.clear(); });
it('does not invent an unresolved start operation from preparation busy alone', () => {
  expect(recordingStartObservation(scope, { unresolved: false, busy: true })).toEqual({ state: 'unavailable', reason: 'loading' });
});
it('does not relabel a previously terminal O1 as pending when saving new settings', () => {
  expect(recordingStartObservation(scope, { unresolved: false, busy: true, operationId: 'O1', status: 'succeeded' }))
    .toEqual({ state: 'unavailable', reason: 'loading' });
});
it('preserves a real live missing-ID lost acknowledgement even while busy and storage fails', () => {
  vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => { throw new Error('private'); });
  expect(recordingStartObservation(scope, { unresolved: true, busy: true })).toEqual({ state: 'unresolved', reason: 'lost_ack' });
});
it('keeps a fresh retained O2 above an older terminal O1', () => {
  saveBasicRecordingIntent({ scope_key: scope, request, operation_id: 'O2' });
  expect(recordingStartObservation(scope, { unresolved: false, operationId: 'O1', status: 'succeeded' }))
    .toEqual({ state: 'unresolved', reason: 'unknown', operationId: 'O2' });
});
it('only accepts terminal O1 when current storage is readable and matches it or is locally empty', () => {
  saveBasicRecordingIntent({ scope_key: scope, request, operation_id: 'O1' });
  expect(recordingStartObservation(scope, { unresolved: false, operationId: 'O1', status: 'succeeded' }))
    .toEqual({ state: 'observed', result: 'resolved', status: 'succeeded', operationId: 'O1', reason: 'terminal' });
  vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => { throw new Error('private'); });
  expect(recordingStartObservation(scope, { unresolved: false, operationId: 'O1', status: 'succeeded' }))
    .toEqual({ state: 'unavailable', reason: 'unreadable' });
});
