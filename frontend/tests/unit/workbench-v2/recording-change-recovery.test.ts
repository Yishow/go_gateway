import { describe, expect, it, vi } from 'vitest';
import { RecordingChangeController } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingChangeController';
import { loadRecordingChanges, recordingChangeKey, saveRecordingChange } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/recordingChangeRecovery';
import type { RecordingChangeRecovery } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/recordingChangeTypes';
import { fixture, HASH } from '../../fixtures/recordingChange';
function root() { const f = fixture(); const record: RecordingChangeRecovery = { version: 1, envelope: f.envelope, reviewDigest: HASH, submitted: { action: 'reserve', request: { envelope: f.envelope, review_digest: HASH } }, phase: 'impact', outcome: 'pending' }; return { f, record }; }
describe('exact root browser recovery validation', () => {
  it('rejects injected private detail instead of persisting it', () => {
    const { f, record } = root(); Object.assign(record.envelope.request.group, { privateDetail: 'secret' }); expect(saveRecordingChange(record)).toBe(false); expect(loadRecordingChanges(f.workspace)).toEqual([]);
  });
  it('never accepts a malformed stored success without terminal evidence', () => {
    const { f, record } = root(); sessionStorage.setItem(recordingChangeKey(record), JSON.stringify({ ...record, phase: 'scheduled', outcome: 'recorded' }));
    const c = new RecordingChangeController(f.workspace, f.deps); expect(c.state.phase).not.toBe('scheduled'); expect(c.unresolved).toBe(true); expect(c.state.recovery?.unreadable).toBe(true);
  });
  it('does not persist another root under the original exact stage request', () => {
    const { record } = root(); record.envelope = structuredClone(record.envelope); record.envelope.request.operation_id = 'OTHER'; expect(saveRecordingChange(record)).toBe(false);
  });
  it('keeps a memory fence when storage fails and blocks sending any new mutation', async () => {
    const f = fixture(); const c = new RecordingChangeController(f.workspace, f.deps); await c.requestAction(f.intent); await c.review();
    const storage = sessionStorage; vi.stubGlobal('sessionStorage', { getItem: storage.getItem.bind(storage), key: storage.key.bind(storage), length: storage.length, setItem: () => { throw new Error('quota'); } });
    try { await c.reserveAndPreview(); expect(f.api.reserve).not.toHaveBeenCalled(); expect(c.unresolved).toBe(true); expect(c.state.error).toBe('storage'); } finally { vi.unstubAllGlobals(); }
  });
});

it.each(['{unparseable', JSON.stringify({ envelope: 'corrupt', reviewDigest: 'invalid' })])('cold-load corruption keeps a key-scoped status-only fence', async value => {
  const { f, record } = root(); sessionStorage.setItem(recordingChangeKey(record), value);
  const c = new RecordingChangeController(f.workspace, f.deps); expect(c.unresolved).toBe(true); c.readonly = true; c.resume(); await c.checkStatus();
  expect(f.api.operationStatus).toHaveBeenCalledWith('G1', { workspace_id: f.workspace, lineage_id: 'L', operation_id: 'ROOT' });
  await c.requestAction(f.intent); await c.retryExact(); await c.repreview(); expect(f.api.reserve).not.toHaveBeenCalled(); expect(f.api.schemaPreview).not.toHaveBeenCalled(); expect(c.unresolved).toBe(true);
});

it.each([{ priorPreviewIds: 'SCHEMA' }, { attempt: 'bad' }, { schemaConfirmationSubmitted: 'false' }, { priorPreviewStates: { SCHEMA: '<private>' } }])('corrupt optional recovery metadata becomes an unreadable fence, not a crash', patch => {
  const { f, record } = root(); sessionStorage.setItem(recordingChangeKey(record), JSON.stringify({ ...record, ...patch }));
  const c = new RecordingChangeController(f.workspace, f.deps); expect(c.unresolved).toBe(true); expect(c.state.recovery?.unreadable ?? c.state.unreadableScope).toBeTruthy();
});

it('stored cutover success remains bound to reviewed digest and exact started/ended revisions', () => {
  const { f, record } = root(); f.envelope.source_applied_revision = 'A1';
  const confirmation = { envelope: { ...f.activation.review.envelope, version: 'successor-activation-review-v2' as const, activation_kind: 'cutover' as const, source_applied_revision: 'A1' }, review_digest: HASH, readiness_digest: HASH };
  const terminal: RecordingChangeRecovery = { ...record, submitted: { action: 'confirmCutover', request: confirmation }, phase: 'scheduled', outcome: 'recorded',
    result: { id: 'CUT', workspace_id: f.workspace, kind: 'forward', source_group_id: 'G1', boundary: confirmation.envelope.request.boundary, review_digest: 'b'.repeat(64), ended: [{ group_id: 'G1', revision: 'A1' }], started: [{ group_id: 'G2', revision: 'R2' }], moved: [{ device_id: 'D', point_id: 'P', tag_id: 'T' }], created_at: '2026-01-01T00:00:00Z' } };
  expect(saveRecordingChange(terminal)).toBe(false);
});
