import { describe, expect, it } from 'vitest';
import { RecordingChangeController } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingChangeController';
import { saveRecordingChange } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/recordingChangeRecovery';
import type { WriteGroupSuccessorCurrentCAS } from '@/types/studioV2WriteGroupSuccessor';
import { parseWriteGroupSuccessorCurrentCAS } from '@/utils/studioV2WriteGroupSuccessorJson';
import { BOUNDARY, fixture } from '../../fixtures/recordingChange';

// Exhaustive against the typed CAS contract: adding a field requires a rejection case.
const changedCAS = {
  workspace_revision: 'OTHER-WORKSPACE-REVISION', lineage_revision: 99,
  predecessor_revision: 'OTHER-PREDECESSOR-REVISION', predecessor_connector_revision: 'OTHER-SOURCE-CONNECTOR-REVISION',
  group_id: 'FOREIGN-GROUP', group_revision: 'OTHER-GROUP-REVISION', target_connector_revision: 'OTHER-TARGET-CONNECTOR-REVISION',
} satisfies WriteGroupSuccessorCurrentCAS;
const fields = Object.keys(changedCAS) as (keyof WriteGroupSuccessorCurrentCAS)[];
async function previewed() {
  const f = fixture();
  const controller = new RecordingChangeController(f.workspace, f.deps);
  await controller.requestAction(f.intent); await controller.review(); await controller.reserveAndPreview();
  expect(controller.state.phase).toBe('schema-plan');
  return { f, controller };
}
function parsedCurrent(current: WriteGroupSuccessorCurrentCAS) {
  const parsed = parseWriteGroupSuccessorCurrentCAS(current);
  expect(parsed).toEqual(current);
  expect(Object.keys(parsed!)).not.toEqual(Object.keys(current));
  return parsed!;
}
async function activationReady() {
  const context = await previewed();
  await context.controller.confirmSchema();
  context.controller.setBoundary(BOUNDARY); await context.controller.reviewActivation();
  expect(context.controller.state.phase).toBe('activation-ready');
  return context;
}

describe('controller compares the whole current CAS independently of property insertion order', () => {
  it('confirms the exact schema plan when the parser returns equal CAS fields in a different order', async () => {
    const { f, controller } = await previewed();
    const stage = controller.state.previewStage!;
    controller.state.previewStage = { ...stage, current: parsedCurrent(stage.current) };
    await controller.confirmSchema();
    expect(controller.state.phase).toBe('schema-confirmed');
    expect(controller.state.unknown).toBe(false);
    expect(f.api.schemaConfirm).toHaveBeenCalledExactlyOnceWith('G1', {
      ...stage, token: 'TOKEN', digest: f.preview.digest, schema_operation_id: 'SCHEMA',
    });
  });

  it('restores the same schema plan after status when persisted CAS fields have parser ordering', async () => {
    const { f, controller } = await previewed();
    const record = controller.state.recovery!, stage = record.previewStage!;
    expect(saveRecordingChange({ ...record, previewStage: { ...stage, current: parsedCurrent(stage.current) } })).toBe(true);
    controller.dispose();
    const restored = new RecordingChangeController(f.workspace, f.deps);
    await restored.checkStatus();
    expect(restored.state.phase).toBe('schema-plan');
    expect(restored.state.unknown).toBe(false);
    expect(restored.state.preview?.operation_id).toBe('SCHEMA');
    expect(f.api.schemaPreview).toHaveBeenCalledTimes(1);
    expect(f.api.schemaConfirm).not.toHaveBeenCalled();
    await restored.confirmSchema();
    expect(restored.state.phase).toBe('schema-confirmed');
    expect(f.api.schemaConfirm).toHaveBeenCalledTimes(1);
  });

  it('confirms activation when its reviewed CAS has parser ordering', async () => {
    const { f, controller } = await activationReady();
    const stage = controller.state.activation!.review.envelope.request.stage;
    stage.current = parsedCurrent(stage.current);
    await controller.confirmActivation();
    expect(controller.state.phase).toBe('scheduled');
    expect(f.api.confirmFirstStart).toHaveBeenCalledTimes(1);
    expect(f.api.confirmCutover).not.toHaveBeenCalled();
  });

  it.each(fields)('schema confirmation rejects changed %s without DDL', async field => {
    const { f, controller } = await previewed();
    const stage = controller.state.previewStage!;
    controller.state.previewStage = { ...stage, current: { ...stage.current, [field]: changedCAS[field] } };
    await controller.confirmSchema();
    expect(controller.state.error).toBe('stale');
    expect(controller.state.unknown).toBe(true);
    expect(f.api.schemaConfirm).not.toHaveBeenCalled();
    expect(controller.state.recovery?.envelope).toEqual(f.envelope);
  });

  it.each(fields)('status recovery rejects changed %s instead of reviving its schema plan', async field => {
    const { f, controller } = await previewed();
    const record = controller.state.recovery!, stage = record.previewStage!;
    expect(saveRecordingChange({ ...record, previewStage: { ...stage, current: { ...stage.current, [field]: changedCAS[field] } } })).toBe(true);
    controller.dispose();
    const restored = new RecordingChangeController(f.workspace, f.deps);
    await restored.checkStatus(); await restored.confirmSchema();
    expect(restored.state.phase).toBe('reserved');
    expect(restored.state.preview).toBeUndefined();
    expect(f.api.schemaConfirm).not.toHaveBeenCalled();
    expect(f.api.schemaPreview).toHaveBeenCalledTimes(1);
    expect(restored.state.recovery?.envelope).toEqual(f.envelope);
  });

  it.each(fields)('activation confirmation rejects changed %s without either activation mutation', async field => {
    const { f, controller } = await activationReady();
    const stage = controller.state.activation!.review.envelope.request.stage;
    stage.current = { ...stage.current, [field]: changedCAS[field] };
    await controller.confirmActivation();
    expect(controller.state.phase).toBe('schema-confirmed');
    expect(controller.state.error).toBe('stale');
    expect(controller.state.activation).toBeUndefined();
    expect(f.api.confirmFirstStart).not.toHaveBeenCalled();
    expect(f.api.confirmCutover).not.toHaveBeenCalled();
    expect(controller.state.recovery?.envelope).toEqual(f.envelope);
  });

  it.each(fields)('the existing parser rejects missing or corrupt %s before it can become CAS authority', field => {
    const current = fixture().stage.current;
    const missing: Partial<WriteGroupSuccessorCurrentCAS> = { ...current };
    delete missing[field];
    expect(parseWriteGroupSuccessorCurrentCAS(missing)).toBeNull();
    expect(parseWriteGroupSuccessorCurrentCAS({ ...current, [field]: null })).toBeNull();
    expect(parseWriteGroupSuccessorCurrentCAS({ ...current, [field]: field === 'lineage_revision' ? '2' : 2 })).toBeNull();
  });
});
