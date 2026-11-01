import { expect, it } from 'vitest';
import { savedStartHasMatchingReceipt } from '../../../src/features/datalink/workbench-v2/recording-editor/SavedRecordingStartEvidence';
import { group, delivery } from './recording-home-guidance-fixtures';
import type { RecordingStartOperation } from '../../../src/types/studioV2RecordingStart';
const saved = { ...group('G'), record_id: 'canonical-mixed-record', applied_revision: 'G2' };
const operation: RecordingStartOperation = { operation_id: 'O', action: 'recording_start', status: 'succeeded', stage: 'complete', workspace_id: saved.workspace_id, setup_revision: 'W2', intent_digest: 'a'.repeat(64),
 device_ids: ['D1'], devices: [{ device_id: 'D1', activated: true }], groups: [{ group_id: 'G', group_revision: 'G2', applied_revision: 'G2', saved: true, ready: true, applied: true }], created_at: '2026-10-30T10:31:00Z', updated_at: '2026-10-30T10:31:01Z' };
const receipt = { group_revision: 'G2', connector_revision: 'DB2', record_id: 'outbox-record-not-group-id', effect_key: 'E', payload_digest: 'a'.repeat(64), committed_at: '2026-10-30T10:32:00Z' };
const delivered = { ...delivery(saved), last_sql_committed_effect: receipt, revision_stages: { group_revision: 'G2', stages: delivery(saved).stages } };
it('leaves waiting only for a matching original-scope revision receipt after this start', () => {
 expect(savedStartHasMatchingReceipt(saved, 'W2', operation, delivered)).toBe(true);
});
it.each(['previous', 'foreign', 'old-revision', 'timestamp-only'] as const)('keeps waiting for %s evidence', kind => {
 const value = kind === 'timestamp-only' ? { ...delivered, last_sql_committed_effect: undefined }
  : kind === 'foreign' ? { ...delivered, group_id: 'foreign' }
  : { ...delivered, last_sql_committed_effect: { ...receipt, ...(kind === 'previous' ? { committed_at: '2026-10-30T10:30:00Z' } : { group_revision: 'G1' }) } };
 expect(savedStartHasMatchingReceipt(saved, 'W2', operation, value)).toBe(false);
});
it('does not count a receipt from before this start within the same millisecond', () => {
 expect(savedStartHasMatchingReceipt(saved, 'W2', { ...operation, created_at: '2026-10-30T10:32:00.000000002Z' },
  { ...delivered, last_sql_committed_effect: { ...receipt, committed_at: '2026-10-30T10:32:00.000000001Z' } })).toBe(false);
});
