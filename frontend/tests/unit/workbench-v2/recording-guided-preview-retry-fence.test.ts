import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { recordingPreparationBlocked, recordingSchemaPreviewRetryBlocked, type RecordingPreparationFenceInput } from '@/features/datalink/workbench-v2/steps/step4/recordingPreparationFence';
import { retainSchemaOperation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery';
import { saveBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';
import { recordingTestWriteRecovery } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import type { SchemaPreparationPresentation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaPreparationPresentation';
import { savedGroup } from '../../fixtures/writeGroupState';

function input(): RecordingPreparationFenceInput {
  const groups = [savedGroup({ id: 'G1' }), savedGroup({ id: 'G2' })];
  const report = (index: number): SchemaPreparationPresentation => ({
    scope: { workspaceId: 'ws-1', workspaceRevision: 'W1', groupId: groups[index].id, groupRevision: groups[index].revision,
      connectorId: 'conn-1', connectorRevision: 'crev-1', epoch: 1 },
    observation: index ? { owner: 'schema', state: 'observed', result: 'none' } : { owner: 'schema', state: 'unavailable', reason: 'unreadable' },
    phase: index ? 'idle' : 'unavailable', preview: 'none', ...(index ? {} : { previewFailure: true }), commands: {},
  });
  return { ownerKey: 'ws-1/dev-1/recording', workspaceId: 'ws-1', workspaceRevision: 'W1', groups,
    start: { state: 'observed', result: 'none' }, schema: { G1: report(0), G2: report(1) } };
}
function clearRetained() {
  retainSchemaOperation('ws-1', 'G1', undefined, 'O1'); retainSchemaOperation('ws-1', 'G2', undefined, 'O2');
  for (const ref of recordingTestWriteRecovery.read().references) recordingTestWriteRecovery.resolve({ operation_id: ref.operation_id,
    action: 'test_write', status: 'succeeded', write_outcome: 'written_verified', cleanup_status: 'cleaned', created_at: '', updated_at: '' }, ref.workspace_id);
}
beforeEach(() => { clearRetained(); sessionStorage.clear(); localStorage.clear(); });
afterEach(() => { vi.restoreAllMocks(); clearRetained(); sessionStorage.clear(); localStorage.clear(); });
it('allows only the target preview-failure retry while ordinary new commands remain denied', () => {
  const value = input(); expect(recordingSchemaPreviewRetryBlocked(value, 'G1')).toBe(false);
  expect(recordingPreparationBlocked(value)).toBe(true); expect(recordingPreparationBlocked(value, 'G1')).toBe(true);
  expect(value.schema.G1.observation.state).toBe('unavailable');
});
it.each(['unmarked', 'foreign', 'stale', 'wrong_group', 'sibling_failed', 'sibling_pending', 'start_pending'] as const)('does not waive %s owner evidence', kind => {
  const value = input(); const own = value.schema.G1, sibling = value.schema.G2;
  if (kind === 'unmarked') value.schema = { ...value.schema, G1: { ...own, previewFailure: undefined } };
  if (kind === 'foreign') value.schema = { ...value.schema, G1: { ...own, scope: { ...own.scope, workspaceId: 'foreign' } } };
  if (kind === 'stale') value.workspaceRevision = 'W2';
  if (kind === 'wrong_group') { expect(recordingSchemaPreviewRetryBlocked(value, 'not-owned')).toBe(true); return; }
  if (kind === 'sibling_failed') value.schema = { ...value.schema, G2: { ...sibling, previewFailure: true, observation: { owner: 'schema', state: 'unavailable', reason: 'unreadable' } } };
  if (kind === 'sibling_pending') value.schema = { ...value.schema, G2: { ...sibling, observation: { owner: 'schema', state: 'unresolved', operationId: 'O2', reason: 'pending' } } };
  if (kind === 'start_pending') value.start = { state: 'unresolved', operationId: 'START-O1', reason: 'unknown' };
  expect(recordingSchemaPreviewRetryBlocked(value, 'G1')).toBe(true);
});
it.each(['G1', 'G2'] as const)('freshly retained schema operation on %s blocks a previously eligible callback', id => {
  const value = input(); expect(recordingSchemaPreviewRetryBlocked(value, 'G1')).toBe(false);
  retainSchemaOperation('ws-1', id, id === 'G1' ? 'O1' : 'O2'); expect(recordingSchemaPreviewRetryBlocked(value, 'G1')).toBe(true);
});
it('freshly retained original start blocks retry without changing its identity', () => {
  const value = input(); expect(recordingSchemaPreviewRetryBlocked(value, 'G1')).toBe(false);
  saveBasicRecordingIntent({ scope_key: value.ownerKey, operation_id: 'START-O1', request: {
    request_id: 'original-request', workspace_id: 'ws-1', expected_workspace_revision: 'W1', device_ids: ['dev-1'],
    groups: value.groups.map(group => ({ group_id: group.id, expected_group_revision: group.revision, expected_connector_revision: 'crev-1' })) } });
  expect(recordingSchemaPreviewRetryBlocked(value, 'G1')).toBe(true);
});
it.each(['G1', 'G2'] as const)('freshly retained test-write on %s blocks retry, including a sibling', id => {
  const value = input(); expect(recordingSchemaPreviewRetryBlocked(value, 'G1')).toBe(false);
  recordingTestWriteRecovery.retain({ workspace_id: 'ws-1', group_id: id, group_revision: 'old-R', operation_id: 'TW-O1', operation_kind: 'test_write', created_at: '2030-01-01T00:00:00Z' });
  expect(recordingSchemaPreviewRetryBlocked(value, 'G1')).toBe(true);
});
it('unreadable storage never becomes empty retry permission', () => {
  const value = input(); vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => { throw new Error('unavailable'); });
  expect(recordingSchemaPreviewRetryBlocked(value, 'G1')).toBe(true);
});
