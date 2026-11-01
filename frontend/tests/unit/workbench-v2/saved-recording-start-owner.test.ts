import { beforeEach, expect, it } from 'vitest';
import { saveBasicRecordingIntent, loadBasicRecordingIntent, inspectBasicRecordingDevice, basicRecordingIntentKey } from '../../../src/features/datalink/workbench-v2/state/basicRecording';
import { recordingCompletion } from '../../../src/features/datalink/workbench-v2/steps/step4/recordingCompletion';
import type { RecordingStartOperation } from '../../../src/types/studioV2RecordingStart';
const owner = { kind: 'saved_record' as const, groupId: 'G', deviceIds: ['A', 'B'] };
const record = { scope_key: 'W/saved-record/G', owner, operation_id: 'O', request: { request_id: 'R', workspace_id: 'W', expected_workspace_revision: 'w1', device_ids: ['A', 'B'], groups: [{ group_id: 'G', expected_group_revision: 'g1', expected_connector_revision: 'c1' }] } };
const result: RecordingStartOperation = { operation_id: 'O', action: 'recording_start', workspace_id: 'W', setup_revision: 'w2', intent_digest: 'a'.repeat(64), status: 'succeeded', stage: 'complete',
 device_ids: ['B', 'A'], devices: [{ device_id: 'A', activated: true }, { device_id: 'B', activated: true }], groups: [{ group_id: 'G', group_revision: 'g1', applied_revision: 'g1', saved: true, ready: true, applied: true }], created_at: '2026-10-31T00:00:00Z', updated_at: '2026-10-31T00:00:01Z' };
beforeEach(() => sessionStorage.clear());
it('the same original store roundtrips an explicitly scoped two-device saved-record request', () => {
 saveBasicRecordingIntent(record);
 expect(loadBasicRecordingIntent(record.scope_key)).toMatchObject(record);
 expect(inspectBasicRecordingDevice('W', 'B').retained).toEqual([{ operationId: 'O', groupIds: ['G'] }]);
 expect(inspectBasicRecordingDevice('W', 'C').retained).toEqual([]);
});
it('completion requires exactly the reviewed saved-record device and group identities', () => {
 expect(recordingCompletion(record, result, 'W', 'A', owner)).toMatchObject({ state: 'available', activatedDeviceIds: ['A', 'B'] });
 expect(recordingCompletion(record, { ...result, device_ids: ['A'] }, 'W', 'A', owner)).toMatchObject({ state: 'unavailable' });
 expect(recordingCompletion(record, result, 'W', 'A', { ...owner, groupId: 'foreign' })).toMatchObject({ state: 'unavailable' });
});
it('legacy requests remain single-device and cannot borrow the saved-record discriminator', () => {
 const scope = 'W/A/legacy';
 for (const value of [{ ...record, scope_key: scope, owner: undefined }, { ...record, scope_key: scope }]) {
  sessionStorage.setItem(basicRecordingIntentKey(scope), JSON.stringify({ ...value, version: 1 }));
  expect(loadBasicRecordingIntent(scope)).toBeUndefined();
 }
});
it('saved-record requests cannot smuggle a second group, a draft, or mismatched source list', () => {
 for (const request of [{ ...record.request, groups: [...record.request.groups, { ...record.request.groups[0], group_id: 'H' }] },
  { ...record.request, groups: [{ ...record.request.groups[0], draft: {} }] }, { ...record.request, device_ids: ['A'] }]) {
  sessionStorage.setItem(basicRecordingIntentKey(record.scope_key), JSON.stringify({ ...record, request, version: 1 }));
  expect(loadBasicRecordingIntent(record.scope_key)).toBeUndefined();
 }
});
