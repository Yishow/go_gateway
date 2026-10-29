import { afterEach, describe, expect, it, vi } from 'vitest';
import { RecordingChangeController } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingChangeController';
import { loadRecordingChanges, recordingChangeKey } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/recordingChangeRecovery';
import type { RecordingChangeRecovery } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/recordingChangeTypes';
import { BOUNDARY, fixture, HASH } from '../../fixtures/recordingChange';

afterEach(() => vi.unstubAllGlobals());

async function previewed() {
  const f = fixture();
  const controller = new RecordingChangeController(f.workspace, f.deps);
  await controller.requestAction(f.intent);
  await controller.review();
  await controller.reserveAndPreview();
  return { f, controller };
}

async function schedule(controller: RecordingChangeController) {
  await controller.confirmSchema();
  controller.setBoundary(BOUNDARY);
  await controller.reviewActivation();
  await controller.confirmActivation();
  expect(controller.state.phase).toBe('scheduled');
}

function failPersistence() {
  const storage = sessionStorage;
  vi.stubGlobal('sessionStorage', {
    get length() { return storage.length; },
    key: storage.key.bind(storage),
    getItem: storage.getItem.bind(storage),
    setItem: () => { throw new Error('quota'); },
  });
}

function rootRecord(f: ReturnType<typeof fixture>): RecordingChangeRecovery {
  return { version: 1, envelope: f.envelope, reviewDigest: HASH,
    submitted: { action: 'reserve', request: { envelope: f.envelope, review_digest: HASH } },
    phase: 'impact', outcome: 'unknown' };
}

// A disjoint fake server root in the same workspace, not a second operation on ROOT.
function otherRoot(workspace: ReturnType<typeof fixture>['workspace']) {
  const f = fixture(workspace);
  const ids: Record<string, string> = { G1: 'OTHER-G1', G2: 'OTHER-G2', L: 'OTHER-L', ROOT: 'OTHER-ROOT', SCHEMA: 'OTHER-SCHEMA', TOKEN: 'OTHER-TOKEN' };
  const remap = <T>(value: T): T => JSON.parse(JSON.stringify(value, (_key, item: unknown) => typeof item === 'string' ? ids[item] ?? item : item));
  for (const value of [f.group, f.seed, f.lifecycle, f.envelope, f.review, f.preparation, f.stage, f.preview, f.activation, f.status]) Object.assign(value, remap(value));
  f.lifecycleRead.mockImplementation(async id => !f.api.reserve.mock.calls.length ? f.lifecycle : {
    ...f.lifecycle, workspace_revision: 'W2', group_id: id, group_revision: id === 'OTHER-G2' ? 'R2' : 'R1',
    history: { ...f.lifecycle.history, lineage_revision: 2, reserved_group_id: 'OTHER-G2' },
  });
  return { ...f, intent: { ...f.intent, groupId: 'OTHER-G1' } };
}

describe('independent saved-recording controller recovery edges', () => {
  it.each(['cold load', 'resume'] as const)('%s prioritizes a second unreadable root over the old scheduled record', async entry => {
    const { f, controller } = await previewed();
    await schedule(controller);
    const other = otherRoot(f.workspace);
    sessionStorage.setItem(recordingChangeKey(rootRecord(other)), '{unreadable');
    f.api.operationStatus.mockResolvedValue(other.status);
    controller.dismiss();
    const active = entry === 'cold load' ? new RecordingChangeController(f.workspace, f.deps) : controller;
    active.readonly = true;
    if (entry === 'resume') active.resume();
    expect(active.state.recovery).toBeUndefined();
    expect(active.state.phase).not.toBe('scheduled');
    await active.checkStatus();
    expect(f.api.operationStatus).toHaveBeenCalledExactlyOnceWith('OTHER-G1', {
      workspace_id: f.workspace, lineage_id: 'OTHER-L', operation_id: 'OTHER-ROOT',
    });
    expect(active.state.readback?.operation_id).toBe('OTHER-ROOT');
    expect(active.state.unknown).toBe(true);
    expect(active.unresolved).toBe(true);
    expect(active.state.canRetry).toBe(false);
    expect(active.state.canRepreview).toBe(false);
  });

  it('a completed recreate can edit and reserve a different group/action without reusing its root', async () => {
    const { f, controller } = await previewed();
    await schedule(controller);
    const other = otherRoot(f.workspace);
    other.lifecycle.actions = [{ action: 'transition', allowed: true, reason_code: '', next_action: '', requires_review: true }];
    controller.dismiss();
    controller.deps = other.deps;
    await controller.requestAction({ ...other.intent, action: 'transition' });
    expect(controller.state).toMatchObject({ phase: 'choose', boundary: '', unknown: false });
    expect(controller.state.recovery).toBeUndefined();
    expect(controller.state.activationReview).toBeUndefined();
    expect(controller.state.previewStage).toBeUndefined();
    controller.setTarget({ ...other.seed, name: 'A distinct recording' });
    await controller.review();
    expect(other.api.review.mock.calls[0][0]).toBe('OTHER-G1');
    expect(other.api.review.mock.calls[0][1].group.name).toBe('A distinct recording');
    await controller.reserveAndPreview();
    expect(other.api.reserve).toHaveBeenCalledExactlyOnceWith('OTHER-G1', {
      envelope: other.review.envelope, review_digest: HASH,
    });
    expect(controller.state.phase).toBe('schema-plan');
    expect(controller.state.recovery?.envelope.request.operation_id).toBe('OTHER-ROOT');
    const records = loadRecordingChanges(f.workspace);
    expect(records.find(record => record.envelope.request.operation_id === 'ROOT')?.phase).toBe('scheduled');
    expect(records.find(record => record.envelope.request.operation_id === 'OTHER-ROOT')?.phase).toBe('schema-plan');
    expect(other.api.schemaConfirm).not.toHaveBeenCalled();
  });

  it('a late old reservation cannot overwrite a different active root or its terminal history', async () => {
    const f = fixture();
    const old = new RecordingChangeController(f.workspace, f.deps);
    await old.requestAction(f.intent);
    await old.review();
    let resolve!: (value: typeof f.preparation) => void;
    f.api.reserve.mockImplementationOnce(() => new Promise(done => { resolve = done; }));
    const pending = old.reserveAndPreview();
    old.dispose();
    // Another session has recorded the original activation before its lost reserve reply arrives.
    f.status.first_start = { status: 'succeeded', result: {
      operation_id: 'FIRST', root_operation_id: 'ROOT', workspace_revision: 'W3', lineage_id: 'L',
      lineage_revision: 3, predecessor_group_id: 'G1', group_id: 'G2', applied_revision: 'R2',
      boundary: BOUNDARY, generation: 2, status: 'succeeded',
    } };
    const active = new RecordingChangeController(f.workspace, f.deps);
    await active.checkStatus();
    expect(active.state.phase).toBe('scheduled');
    const other = otherRoot(f.workspace);
    active.deps = other.deps;
    await active.requestAction(other.intent);
    await active.review();
    await active.reserveAndPreview();
    const current = structuredClone(active.state);
    resolve(f.preparation);
    await pending;
    expect(active.state).toEqual(current);
    expect(active.state.recovery?.envelope.request.operation_id).toBe('OTHER-ROOT');
    expect(f.api.schemaPreview).not.toHaveBeenCalled();
    expect(f.deps.onChanged).not.toHaveBeenCalled();
    const records = loadRecordingChanges(f.workspace);
    expect(records.find(record => record.envelope.request.operation_id === 'ROOT')?.phase).toBe('scheduled');
    expect(records.find(record => record.envelope.request.operation_id === 'OTHER-ROOT')?.preview?.token).toBe('OTHER-TOKEN');
  });

  it.each(['schema confirmation', 'activation confirmation'] as const)('persistence failure immediately before %s prevents that mutation and survives remount', async stage => {
    const { f, controller } = await previewed();
    if (stage === 'activation confirmation') {
      await controller.confirmSchema();
      controller.setBoundary(BOUNDARY);
      await controller.reviewActivation();
    }
    const submittedBefore = f.api.schemaConfirm.mock.calls.length;
    failPersistence();
    await (stage === 'schema confirmation' ? controller.confirmSchema() : controller.confirmActivation());
    expect(f.api.schemaConfirm).toHaveBeenCalledTimes(submittedBefore);
    expect(f.api.confirmFirstStart).not.toHaveBeenCalled();
    expect(f.api.confirmCutover).not.toHaveBeenCalled();
    expect(controller.state.error).toBe('storage');
    expect(controller.state.unknown).toBe(true);
    expect(controller.state.phase).not.toBe('scheduled');
    controller.dispose();
    const resumed = new RecordingChangeController(f.workspace, f.deps);
    resumed.resume();
    expect(resumed.unresolved).toBe(true);
    await resumed.requestAction(f.intent);
    await resumed.reserveAndPreview();
    await resumed.prepareSchema();
    expect(f.api.reserve).toHaveBeenCalledTimes(1);
    expect(f.api.schemaPreview).toHaveBeenCalledTimes(1);
  });

  it.each([
    { priorPreviewIds: { length: 1, 0: 'SCHEMA' } },
    { priorPreviewIds: [null] },
    { priorPreviewStates: ['succeeded'] },
    { terminalStatus: { first_start: { status: 'succeeded' } } },
    { phase: 'scheduled', outcome: 'recorded', result: { status: 'succeeded', boundary: BOUNDARY } },
  ])('corrupt optional recovery stays status-only through resume and actions: %j', async patch => {
    const f = fixture();
    const record = rootRecord(f);
    sessionStorage.setItem(recordingChangeKey(record), JSON.stringify({ ...record, ...patch }));
    const controller = new RecordingChangeController(f.workspace, f.deps);
    controller.resume();
    await controller.checkStatus();
    await controller.requestAction(f.intent);
    await controller.retryExact();
    await controller.repreview();
    await controller.reserveAndPreview();
    await controller.prepareSchema();
    await controller.confirmSchema();
    await controller.confirmActivation();
    expect(f.api.operationStatus).toHaveBeenCalledExactlyOnceWith('G1', {
      workspace_id: f.workspace, lineage_id: 'L', operation_id: 'ROOT',
    });
    expect(controller.unresolved).toBe(true);
    expect(controller.state.unknown).toBe(true);
    expect(controller.state.phase).not.toBe('scheduled');
    expect(controller.state.recovery?.result).toBeUndefined();
    expect(controller.state.canRetry).toBe(false);
    expect(controller.state.canRepreview).toBe(false);
    for (const call of [f.api.reserve, f.api.schemaPreview, f.api.schemaConfirm, f.api.confirmCutover, f.api.confirmFirstStart]) expect(call).not.toHaveBeenCalled();
  });
});
