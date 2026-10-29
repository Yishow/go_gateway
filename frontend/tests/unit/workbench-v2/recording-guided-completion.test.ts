import { expect, it } from 'vitest';
import { recordingCompletion } from '@/features/datalink/workbench-v2/steps/step4/recordingCompletion';
import type { BasicRecordingIntentRecord } from '@/features/datalink/workbench-v2/state/basicRecording';
import type { RecordingStartOperation } from '@/types/studioV2RecordingStart';
const original = (): BasicRecordingIntentRecord => ({ scope_key: 'W/A/recording', operation_id: 'O1', request: {
  request_id: 'original-request', workspace_id: 'W', expected_workspace_revision: 'old-W', device_ids: ['A'],
  groups: ['G1', 'G2'].map(group_id => ({ group_id, expected_group_revision: 'old-G', expected_connector_revision: 'C1' })),
  readiness_token: 'private-proof', settings_revision: 'original-settings', workspace_revision: 'original-share' } });
const result = (): RecordingStartOperation => ({ operation_id: 'O1', action: 'recording_start', workspace_id: 'W',
  device_ids: ['A'], status: 'succeeded', stage: 'complete', intent_digest: 'a'.repeat(64), setup_revision: 'new-W',
  devices: [{ device_id: 'A', activated: true }], groups: ['G1', 'G2'].map(group_id => ({ group_id, group_revision: 'new-G',
    applied_revision: 'new-G', saved: true, ready: true, applied: true })), created_at: 'created', updated_at: 'updated' });
it('uses only original reviewed A despite its own progress advancing workspace/group revisions', () => {
  expect(recordingCompletion(original(), result(), 'W', 'A')).toEqual({ state: 'available', operationId: 'O1', status: 'succeeded',
    deviceIds: ['A'], groupIds: ['G1', 'G2'], activatedDeviceIds: ['A'], incompleteDeviceIds: [], incompleteGroupIds: [], handoffDeviceId: 'A' });
});
it('truthfully lists incomplete groups on partial success while permitting only confirmed A navigation', () => {
  const operation = result(); operation.status = 'partial'; operation.stage = 'activation'; operation.groups[1].applied = false;
  expect(recordingCompletion(original(), operation, 'W', 'A')).toMatchObject({ state: 'available', status: 'partial',
    activatedDeviceIds: ['A'], incompleteGroupIds: ['G2'], handoffDeviceId: 'A' });
});
it.each(['pending', 'running', 'unknown', 'failed'] as const)('never hands off unfinished or failed %s even with an activated progress flag', status => {
  const view = recordingCompletion(original(), { ...result(), status }, 'W', 'A');
  expect(view).toMatchObject({ state: 'available', status }); expect(view).not.toHaveProperty('handoffDeviceId');
});
it('never invents activation or selected-device fallback from an empty confirmed set', () => {
  const operation = result(); operation.status = 'partial'; operation.devices[0].activated = false;
  expect(recordingCompletion(original(), operation, 'W', 'A')).toMatchObject({ state: 'available', activatedDeviceIds: [], incompleteDeviceIds: ['A'] });
  expect(recordingCompletion(original(), operation, 'W', 'A')).not.toHaveProperty('handoffDeviceId');
});
it.each(['workspace', 'device', 'group', 'operation', 'duplicate-device', 'duplicate-group', 'multiple-activated', 'missing-progress', 'wrong-current-device', 'wrong-owner'] as const)('fails closed for %s without borrowing another identity', mismatch => {
  const record = original(), operation = result();
  if (mismatch === 'workspace') operation.workspace_id = 'foreign-W';
  if (mismatch === 'device') { operation.device_ids = ['B']; operation.devices = [{ device_id: 'B', activated: true }]; }
  if (mismatch === 'group') operation.groups[0].group_id = 'foreign-G';
  if (mismatch === 'operation') operation.operation_id = 'O2';
  if (mismatch === 'duplicate-device') operation.devices.push(operation.devices[0]);
  if (mismatch === 'duplicate-group') operation.groups.push(operation.groups[0]);
  if (mismatch === 'multiple-activated') { operation.device_ids.push('B'); operation.devices.push({ device_id: 'B', activated: true }); }
  if (mismatch === 'missing-progress') operation.groups.pop();
  if (mismatch === 'wrong-owner') record.scope_key = 'foreign-W/A/recording';
  expect(recordingCompletion(record, operation, 'W', mismatch === 'wrong-current-device' ? 'B' : 'A')).toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
});
it.each(['request', 'operation', 'operation-id'] as const)('does not fabricate missing %s evidence', missing => {
  const record = original(); if (missing === 'operation-id') delete record.operation_id;
  expect(recordingCompletion(missing === 'request' ? undefined : record, missing === 'operation' ? undefined : result(), 'W', 'A')).toEqual({ state: 'unavailable', reason: 'missing' });
});
it('does not leak private original request material into the display projection', () => {
  const view = recordingCompletion(original(), result(), 'W', 'A');
  expect(view.state).toBe('available'); const json = JSON.stringify(view);
  for (const privateValue of ['private-proof', 'original-request', 'old-W', 'original-settings', 'original-share']) expect(json).not.toContain(privateValue);
});
it('rejects a purported succeeded result whose requested group is not fully applied', () => {
  const operation = result(); operation.groups[1].applied = false;
  expect(recordingCompletion(original(), operation, 'W', 'A')).toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
});
