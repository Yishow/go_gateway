import { describe, expect, it } from 'vitest';
import { RecordingChangeController } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingChangeController';
import { loadRecordingChanges } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/recordingChangeRecovery';
import { WriteGroupResponseError } from '@/utils/studioV2WriteGroupJson';
import { BOUNDARY, fixture } from '../../fixtures/recordingChange';
async function impact() { const f = fixture(); const c = new RecordingChangeController(f.workspace, f.deps); await c.requestAction(f.intent); await c.review(); return { f, c }; }
describe('one saved-recording operation across stages', () => {
  it('cancel and back do not reserve, probe, issue a token or run DDL', async () => {
    const { f, c } = await impact(); c.back(); c.dismiss();
    expect(c.state.isOpen).toBe(false); expect(f.api.reserve).not.toHaveBeenCalled(); expect(f.api.schemaPreview).not.toHaveBeenCalled(); expect(f.api.schemaConfirm).not.toHaveBeenCalled();
  });
  it('blocks stale entry intent before review and all side effects', async () => {
    const f = fixture(); const c = new RecordingChangeController(f.workspace, f.deps);
    await c.requestAction({ ...f.intent, groupRevision: 'stale' }); await c.review();
    expect(c.state.error).toBe('stale'); expect(f.api.review).not.toHaveBeenCalled(); expect(f.api.reserve).not.toHaveBeenCalled();
  });
  it('requires a separate plan confirmation and synchronously suppresses duplicate reservation', async () => {
    const { f, c } = await impact(); await Promise.all([c.reserveAndPreview(), c.reserveAndPreview()]);
    expect(f.api.reserve).toHaveBeenCalledTimes(1); expect(f.api.schemaPreview).toHaveBeenCalledTimes(1); expect(f.api.schemaConfirm).not.toHaveBeenCalled();
    expect(c.state.phase).toBe('schema-plan'); await c.confirmSchema(); expect(c.state.phase).toBe('schema-confirmed'); expect(f.api.confirmFirstStart).not.toHaveBeenCalled();
  });
  it('reload retains the exact returned preview without replacement preview or automatic DDL', async () => {
    const { f, c } = await impact(); await c.reserveAndPreview(); c.dispose();
    const next = new RecordingChangeController(f.workspace, f.deps); next.resume(); await next.checkStatus();
    expect(next.state.phase).toBe('schema-plan'); expect(next.state.unknown).toBe(false); expect(next.state.preview?.token).toBe('TOKEN');
    await next.confirmSchema(); expect(f.api.schemaConfirm).toHaveBeenCalledTimes(1); expect(f.api.schemaPreview).toHaveBeenCalledTimes(1);
    const text = JSON.stringify(loadRecordingChanges(f.workspace)); expect(text).not.toContain('PRIVATE SQL');
  });
  it('a rejected later stage cannot erase the already reserved root recovery fence', async () => {
    const { f, c } = await impact(); f.api.schemaPreview.mockRejectedValue(new WriteGroupResponseError('schema preview', { code: 'revision_mismatch' }, 'failed')); await c.reserveAndPreview();
    expect(c.unresolved).toBe(true); expect(new RecordingChangeController(f.workspace, f.deps).unresolved).toBe(true);
    await c.requestAction(f.intent); expect(f.api.reserve).toHaveBeenCalledTimes(1);
  });
  it('unknown schema confirmation blocks new preview and keeps exact original operation through reload', async () => {
    const { f, c } = await impact(); await c.reserveAndPreview(); f.api.schemaConfirm.mockRejectedValue(new Error('lost')); await c.confirmSchema();
    c.dismiss(); const next = new RecordingChangeController(f.workspace, f.deps); next.resume(); await next.prepareSchema(); await next.reserveAndPreview();
    expect(next.state.unknown).toBe(true); expect(f.api.schemaPreview).toHaveBeenCalledTimes(1);
    await next.checkStatus(); expect(f.api.operationStatus).toHaveBeenCalledWith('G1', { workspace_id: f.workspace, lineage_id: 'L', operation_id: 'ROOT' });
    expect(next.state.recovery?.submitted).toEqual(c.state.recovery?.submitted);
  });
  it('refreshes current CAS separately while preserving the immutable original envelope', async () => {
    const { f, c } = await impact(); await c.reserveAndPreview();
    expect(f.api.schemaPreview.mock.calls[0][1].current.workspace_revision).toBe('W2');
    expect(f.api.schemaPreview.mock.calls[0][1].envelope.request.expected_workspace_revision).toBe('W1');
  });
  it('does not continue preview or emit a late selection callback after dismissal', async () => {
    const { f, c } = await impact(); let resolve!: (value: typeof f.preparation) => void;
    f.api.reserve.mockImplementation(() => new Promise(r => { resolve = r; })); const pending = c.reserveAndPreview(); c.dismiss(); resolve(f.preparation); await pending;
    expect(f.api.schemaPreview).not.toHaveBeenCalled(); expect(f.deps.onChanged).not.toHaveBeenCalled(); expect(loadRecordingChanges(f.workspace)[0].preparation?.group.id).toBe('G2');
  });
  it('dispatches first start, records scheduled only, and never downgrades it on failed status', async () => {
    const { f, c } = await impact(); await c.reserveAndPreview(); await c.confirmSchema(); c.setBoundary(BOUNDARY); await c.reviewActivation(); await c.confirmActivation();
    expect(f.api.confirmFirstStart).toHaveBeenCalledTimes(1); expect(f.api.confirmCutover).not.toHaveBeenCalled(); expect(c.state.phase).toBe('scheduled');
    f.api.operationStatus.mockRejectedValue(new Error('offline')); await c.checkStatus(); expect(c.state.phase).toBe('scheduled'); expect(c.state.unknown).toBe(false);
    expect(new RecordingChangeController(f.workspace, f.deps).state.phase).toBe('scheduled');
  });
  it('after validated not-found permits only explicit exact reserve retry', async () => {
    const { f, c } = await impact(); f.api.reserve.mockRejectedValueOnce(new Error('never arrived')); await c.reserveAndPreview();
    f.api.operationStatus.mockRejectedValue(new WriteGroupResponseError('status', { code: 'WRITE_GROUP_NOT_FOUND' }, 'failed')); await c.checkStatus();
    expect(c.state.canRetry).toBe(true); await c.retryExact(); expect(f.api.reserve.mock.calls[1]).toEqual(f.api.reserve.mock.calls[0]); expect(f.api.schemaPreview).not.toHaveBeenCalled();
  });
});

it('a known rejected preview returns to explicit preparation only after status proves no schema operation', async () => {
  const { f, c } = await impact(); f.api.schemaPreview.mockRejectedValueOnce(new WriteGroupResponseError('preview', { code: 'revision_mismatch' }, 'failed'));
  await c.reserveAndPreview(); f.status.schema_operations = []; await c.checkStatus();
  expect(c.state.phase).toBe('reserved'); expect(c.state.unknown).toBe(false); expect(f.api.schemaPreview).toHaveBeenCalledTimes(1);
  await c.prepareSchema(); expect(f.api.reserve).toHaveBeenCalledTimes(1); expect(f.api.schemaPreview).toHaveBeenCalledTimes(2);
});
it('missing schema status never proves the persisted preview is confirmable', async () => {
  const { f, c } = await impact(); await c.reserveAndPreview(); f.status.schema_operations = []; c.dismiss(); c.resume(); await c.checkStatus();
  expect(c.state.unknown).toBe(true); await c.confirmSchema(); expect(f.api.schemaConfirm).not.toHaveBeenCalled();
});
it('stale current CAS blocks table confirmation before a DDL request', async () => {
  const { f, c } = await impact(); await c.reserveAndPreview();
  f.lifecycleRead.mockImplementation(async id => ({ ...f.lifecycle, group_id: id, group_revision: id === 'G2' ? 'R2' : 'R1', workspace_revision: 'W3', history: { ...f.lifecycle.history, lineage_revision: 2, reserved_group_id: 'G2' } }));
  await c.confirmSchema(); expect(f.api.schemaConfirm).not.toHaveBeenCalled(); expect(c.state.error).toBe('stale');
});
it('Go whole-second boundary remains exactly the server-confirmed text on first start', async () => {
  const { f, c } = await impact(); await c.reserveAndPreview(); await c.confirmSchema();
  const boundary = '2099-01-01T00:00:00Z'; f.activation.review.envelope.request.boundary = boundary;
  f.api.confirmFirstStart.mockResolvedValue({ operation_id: 'FIRST', root_operation_id: 'ROOT', workspace_revision: 'W3', lineage_id: 'L', lineage_revision: 3, predecessor_group_id: 'G1', group_id: 'G2', applied_revision: 'R2', boundary, generation: 2, status: 'succeeded' });
  c.setBoundary(BOUNDARY); await c.reviewActivation(); await c.confirmActivation();
  expect(f.api.confirmFirstStart.mock.calls[0][1].envelope.request.boundary).toBe(boundary);
});

it('lost-preview status only unlocks an explicit same-root repreview when every operation proves no DDL', async () => {
  const { f, c } = await impact(); f.api.schemaPreview.mockRejectedValueOnce(new Error('lost preview')); await c.reserveAndPreview();
  await c.checkStatus(); expect(c.state.canRepreview).toBe(true); expect(f.api.schemaPreview).toHaveBeenCalledTimes(1); expect(f.api.schemaConfirm).not.toHaveBeenCalled();
  await c.repreview(); expect(f.api.reserve).toHaveBeenCalledTimes(1); expect(f.api.schemaPreview).toHaveBeenCalledTimes(2);
  expect(f.api.schemaPreview.mock.calls[1][1].envelope).toEqual(f.api.schemaPreview.mock.calls[0][1].envelope);
  expect(c.state.recovery?.priorPreviewIds).toEqual(['SCHEMA']);
});
it.each(['pending', 'running', 'succeeded', 'failed', 'partial', 'unknown'] as const)('lost preview remains locked when a schema operation is %s', async status => {
  const { f, c } = await impact(); f.api.schemaPreview.mockRejectedValueOnce(new Error('lost preview')); await c.reserveAndPreview(); f.status.schema_operations[0].status = status;
  await c.checkStatus(); expect(c.state.canRepreview).toBe(false); await c.repreview(); expect(f.api.schemaPreview).toHaveBeenCalledTimes(1);
});

it('repreview remains blocked for prior local schema-confirm, foreign reserve identity and failed status', async () => {
  for (const mode of ['prior-confirm', 'foreign-reserve', 'status-failed']) {
    const { f, c } = await impact(); f.api.schemaPreview.mockRejectedValueOnce(new Error('lost preview')); await c.reserveAndPreview();
    if (mode === 'prior-confirm') c.state.recovery!.schemaConfirmationSubmitted = true;
    if (mode === 'foreign-reserve') f.status.reserve.result = { ...f.preparation, generation: 99 };
    if (mode === 'status-failed') f.api.operationStatus.mockRejectedValue(new Error('offline'));
    await c.checkStatus(); expect(c.state.canRepreview).toBe(false); await c.repreview(); expect(f.api.schemaPreview).toHaveBeenCalledTimes(1);
  }
});
it('zero schema operations permits only an explicit repreview; activation evidence prevents it', async () => {
  const { f, c } = await impact(); f.api.schemaPreview.mockRejectedValueOnce(new Error('lost preview')); await c.reserveAndPreview(); f.status.schema_operations = [];
  await c.checkStatus(); expect(c.state.canRepreview).toBe(true); expect(f.api.schemaPreview).toHaveBeenCalledTimes(1);
  f.status.first_start = { status: 'succeeded', result: { operation_id: 'FIRST', root_operation_id: 'ROOT', workspace_revision: 'W3', lineage_id: 'L', lineage_revision: 3, predecessor_group_id: 'G1', group_id: 'G2', applied_revision: 'R2', boundary: BOUNDARY, generation: 2, status: 'succeeded' } };
  await c.repreview(); expect(f.api.schemaPreview).toHaveBeenCalledTimes(1); expect(c.state.canRepreview).toBe(false);
});
it('late preview after remount cannot replace a newer same-root plan or remove its fence', async () => {
  const { f, c } = await impact(); let late!: (value: typeof f.preview) => void;
  f.api.schemaPreview.mockImplementationOnce(() => new Promise(resolve => { late = resolve; })); const old = c.reserveAndPreview();
  await new Promise(resolve => setTimeout(resolve, 0)); c.dispose();
  const next = new RecordingChangeController(f.workspace, f.deps); next.resume(); await next.checkStatus(); expect(next.state.canRepreview).toBe(true);
  f.api.schemaPreview.mockResolvedValue({ ...f.preview, operation_id: 'NEW-SCHEMA', token: 'NEW-TOKEN' }); await next.repreview();
  late(f.preview); await old;
  expect(loadRecordingChanges(f.workspace)[0].preview?.token).toBe('NEW-TOKEN'); expect(next.state.preview?.token).toBe('NEW-TOKEN'); expect(f.api.schemaConfirm).not.toHaveBeenCalled();
});

it.each(['successor-activation-review-v1', 'successor-activation-review-v2'] as const)('dispatches a legitimate %s cutover without ordinary first start', async version => {
  const { f, c } = await impact(); f.group.applied_revision = 'A1'; f.envelope.source_applied_revision = 'A1';
  const evidence = { ...f.activation.review.envelope, source_applied_revision: 'A1' };
  f.activation.review.envelope = version === 'successor-activation-review-v1'
    ? { ...evidence, version, activation_kind: undefined } : { ...evidence, version, activation_kind: 'cutover' };
  f.api.confirmCutover.mockResolvedValue({ id: 'CUT', workspace_id: f.workspace, kind: 'forward', source_group_id: 'G1', boundary: BOUNDARY, review_digest: f.review.review_digest, ended: [{ group_id: 'G1', revision: 'A1' }], started: [{ group_id: 'G2', revision: 'R2' }], moved: [{ device_id: 'D', point_id: 'P', tag_id: 'T' }], created_at: '2026-01-01T00:00:00Z' });
  await c.reserveAndPreview(); await c.confirmSchema(); c.setBoundary(BOUNDARY); await c.reviewActivation(); await c.confirmActivation();
  expect(f.api.confirmCutover).toHaveBeenCalledTimes(1); expect(f.api.confirmFirstStart).not.toHaveBeenCalled(); expect(c.state.phase).toBe('scheduled');
});
it('original-root terminal status survives reload even if another controller completed later stages', async () => {
  const { f, c } = await impact(); f.api.reserve.mockRejectedValue(new Error('lost')); await c.reserveAndPreview();
  f.status.first_start = { status: 'succeeded', result: { status: 'succeeded', operation_id: 'FIRST', root_operation_id: 'ROOT', workspace_revision: 'W3', lineage_id: 'L', lineage_revision: 3, predecessor_group_id: 'G1', group_id: 'G2', applied_revision: 'R2', boundary: BOUNDARY, generation: 2 } };
  await c.checkStatus(); expect(c.state.phase).toBe('scheduled'); const next = new RecordingChangeController(f.workspace, f.deps); expect(next.state.phase).toBe('scheduled'); expect(next.state.unknown).toBe(false);
});

it.each(['terminal', 'rejected'])('opening a new review after %s clears every previous operation field', async outcome => {
  const { f, c } = await impact();
  if (outcome === 'terminal') { await c.reserveAndPreview(); await c.confirmSchema(); c.setBoundary(BOUNDARY); await c.reviewActivation(); await c.confirmActivation(); }
  else { f.api.reserve.mockRejectedValueOnce(new WriteGroupResponseError('reserve', { code: 'revision_mismatch' }, 'failed')); await c.reserveAndPreview(); }
  c.dismiss(); f.lifecycleRead.mockResolvedValue(f.lifecycle); await c.requestAction(f.intent);
  expect(c.state.recovery).toBeUndefined(); expect(c.state.preview).toBeUndefined(); expect(c.state.activation).toBeUndefined(); expect(c.unresolved).toBe(false);
  await c.review(); expect(f.api.review).toHaveBeenCalledTimes(2); expect(c.state.phase).toBe('impact');
});
