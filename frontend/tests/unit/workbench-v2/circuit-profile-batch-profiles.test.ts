import { describe, expect, it, vi } from 'vitest';
import { executeBatchPlan, type FrozenBatchPlan, type FrozenBatchTarget } from '@/features/datalink/workbench-v2/steps/step4/circuitProfileActions';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { httpFailure, longProfileGroup } from '../../fixtures/circuitBatchPanel';

function target(id: string, withCumulative: boolean): FrozenBatchTarget & { inst: WriteGroup; cum?: WriteGroup } {
  const inst = longProfileGroup({ id: `inst-${id}`, circuitId: id, circuitName: id, interval: 15, points: [0, 1] });
  const cum = withCumulative
    ? longProfileGroup({ id: `cum-${id}`, circuitId: id, circuitName: id, interval: 10, points: [0, 1], kind: 'cumulative' })
    : undefined;
  return { circuit_id: id, circuit_name: id, instantaneousGroup: inst, cumulativeGroup: cum, anyGroup: inst, inst, cum };
}

function plan(targets: FrozenBatchTarget[], apply: { inst?: boolean; cum?: boolean } = { inst: true, cum: true }): FrozenBatchPlan {
  return {
    // This reviewed revision matches the saved destination in these SQLite fixtures.
    workspaceId: 'ws-1', workspaceRevision: 'w-1', connectorRevision: 'crev-1',
    applyInst: Boolean(apply.inst), instInterval: 30, applyCum: Boolean(apply.cum), cumInterval: 20, targets,
  };
}

/** Succeeds with the next workspace revision, echoing the requested group. */
function savedUpdate(revision: string) {
  return async ({ id, request }: { id: string; request: { group: object } }) =>
    ({ workspace_revision: revision, group: { ...(request.group as WriteGroup), id, revision: `${id}-saved` } });
}

async function halted(run: Promise<unknown>) {
  const error = await run.then(() => undefined, (cause: unknown) => cause);
  expect(error).toBeDefined();
  return error as { updatedGroups: WriteGroup[]; workspaceRevision: string; targetResults: Record<string, {
    status: string; profiles: Record<string, { status: string; cause?: unknown } | undefined>;
  }> };
}

describe('Dual-profile batch results are recorded per circuit and profile', () => {
  it('keeps the saved instantaneous writer visible when the cumulative PUT is rejected with 409', async () => {
    const t = target('cir-a', true);
    const update = vi.fn()
      .mockImplementationOnce(savedUpdate('w-2'))
      .mockRejectedValueOnce(httpFailure(409, 'WRITE_GROUP_CIRCUIT_PROFILE_CONFLICT'));
    const failure = await halted(executeBatchPlan({ capabilityContext: { connectorKind: 'sqlite', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' } }, plan: plan([t]), updateWriteGroup: update }));

    expect(update).toHaveBeenCalledTimes(2);
    expect(update.mock.calls[1][0].request.expected_workspace_revision).toBe('w-2');
    expect(failure.targetResults['cir-a'].profiles.instantaneous).toEqual(expect.objectContaining({ status: 'saved' }));
    expect(failure.targetResults['cir-a'].profiles.cumulative).toEqual(expect.objectContaining({ status: 'failed' }));
    expect(failure.targetResults['cir-a'].status).toBe('partial');
    expect(failure.updatedGroups.map((g) => g.id)).toEqual(['inst-cir-a']);
    expect(failure.workspaceRevision).toBe('w-2');
  });

  it('records a lost cumulative response as unknown while the instantaneous writer stays saved', async () => {
    const t = target('cir-a', true);
    const update = vi.fn()
      .mockImplementationOnce(savedUpdate('w-2'))
      .mockRejectedValueOnce(new Error('socket hang up'));
    const failure = await halted(executeBatchPlan({ capabilityContext: { connectorKind: 'sqlite', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' } }, plan: plan([t]), updateWriteGroup: update }));

    expect(failure.targetResults['cir-a'].profiles.instantaneous?.status).toBe('saved');
    expect(failure.targetResults['cir-a'].profiles.cumulative?.status).toBe('unknown');
    expect(failure.targetResults['cir-a'].status).toBe('unknown');
  });

  it('does not attempt the cumulative writer or later circuits after the first profile fails', async () => {
    const a = target('cir-a', true);
    const b = target('cir-b', true);
    const update = vi.fn().mockRejectedValueOnce(httpFailure(409, 'WRITE_GROUP_CIRCUIT_PROFILE_CONFLICT'));
    const failure = await halted(executeBatchPlan({ capabilityContext: { connectorKind: 'sqlite', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' } }, plan: plan([a, b]), updateWriteGroup: update }));

    expect(update).toHaveBeenCalledTimes(1);
    expect(failure.targetResults['cir-a'].profiles.instantaneous?.status).toBe('failed');
    expect(failure.targetResults['cir-a'].profiles.cumulative?.status).toBe('not-attempted');
    expect(failure.targetResults['cir-b'].profiles.instantaneous?.status).toBe('not-attempted');
    expect(failure.targetResults['cir-b'].profiles.cumulative?.status).toBe('not-attempted');
    expect(failure.updatedGroups).toEqual([]);
  });

  it('keeps single-profile circuit A saved, B failed and C untouched', async () => {
    const [a, b, c] = ['cir-a', 'cir-b', 'cir-c'].map((id) => target(id, false));
    const update = vi.fn()
      .mockImplementationOnce(savedUpdate('w-2'))
      .mockRejectedValueOnce(httpFailure(422, 'profile_interval_invalid'));
    const failure = await halted(executeBatchPlan({ capabilityContext: { connectorKind: 'sqlite', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' } }, plan: plan([a, b, c], { inst: true }), updateWriteGroup: update }));

    expect(failure.targetResults['cir-a'].status).toBe('saved');
    expect(failure.targetResults['cir-b'].status).toBe('failed');
    expect(failure.targetResults['cir-c'].status).toBe('not-attempted');
    expect(failure.targetResults['cir-a'].profiles.cumulative).toBeUndefined();
    expect(update.mock.calls.map((call) => call[0].id)).toEqual(['inst-cir-a', 'inst-cir-b']);
  });

  it('skips a profile that already has the reviewed interval instead of rewriting it', async () => {
    const t = target('cir-a', true);
    t.instantaneousGroup = { ...t.inst, recording_profile: { kind: 'instantaneous', interval: 30, unit: 'second' }, row_policy: { ...t.inst.row_policy, interval_seconds: 30 } };
    const update = vi.fn().mockImplementation(savedUpdate('w-2'));
    const result = await executeBatchPlan({ capabilityContext: { connectorKind: 'sqlite', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' } }, plan: plan([t]), updateWriteGroup: update });

    expect(update).toHaveBeenCalledTimes(1);
    expect(update.mock.calls[0][0].id).toBe('cum-cir-a');
    expect(result.targetResults['cir-a'].profiles.instantaneous?.status).toBe('unchanged');
    expect(result.targetResults['cir-a'].profiles.cumulative?.status).toBe('saved');
    expect(result.targetResults['cir-a'].status).toBe('saved');
  });
});
