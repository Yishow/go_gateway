import { describe, expect, it, vi } from 'vitest';
import { executeBatchPlan, type FrozenBatchPlan } from '@/features/datalink/workbench-v2/steps/step4/circuitProfileActions';
import { safeBasicStartError } from '@/features/datalink/workbench-v2/state/basicRecordingMessages';
import { WriteGroupResponseError } from '@/utils/studioV2WriteGroupJson';
import { httpFailure } from '../../fixtures/circuitBatchPanel';
import { savedGroup } from '../../fixtures/writeGroupState';

function plan(): FrozenBatchPlan {
  const group = savedGroup({
    id: 'instant-a', workspace_id: 'ws-1', revision: 'g-1',
    layout_version: 'managed_long_v1', storage_mode: 'plain',
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' },
    members: savedGroup().members.map((member) => ({ ...member, circuit_id: 'cir-a', circuit_name: 'Circuit A' })),
  });
  return {
    workspaceId: 'ws-1', workspaceRevision: 'w-1', connectorRevision: 'c-1',
    applyInst: true, instInterval: 30, applyCum: false, cumInterval: 10,
    targets: [{ circuit_id: 'cir-a', circuit_name: 'Circuit A', instantaneousGroup: group, anyGroup: group }],
  };
}

describe('Batch outcome reflects the evidence actually available', () => {
  it('keeps a lost mutation response unknown instead of asserting that saving failed', async () => {
    const error = await executeBatchPlan({
      plan: plan(), updateWriteGroup: vi.fn().mockRejectedValue(new Error('response lost after server commit')),
    }).catch((cause: unknown) => cause);
    expect(error).toEqual(expect.objectContaining({
      targetResults: expect.objectContaining({ 'cir-a': expect.objectContaining({ status: 'unknown' }) }),
    }));
  });
});

async function outcomeOf(rejection: unknown) {
  const error = await executeBatchPlan({ plan: plan(), updateWriteGroup: vi.fn().mockRejectedValue(rejection) })
    .then(() => undefined, (cause: unknown) => cause) as { targetResults: Record<string, { status: string; profiles: Record<string, { status: string; cause?: unknown }> }> };
  return error.targetResults['cir-a'];
}

describe('Only an explicit rejection proves that nothing was committed', () => {
  it.each([
    ['HTTP 409 revision conflict', httpFailure(409, 'WRITE_GROUP_CIRCUIT_PROFILE_CONFLICT')],
    ['HTTP 422 validation', httpFailure(422, 'profile_interval_invalid')],
    ['HTTP 400 with the handler typed WRITE_GROUP_INVALID_REQUEST', httpFailure(400, 'WRITE_GROUP_INVALID_REQUEST')],
    ['an explicit success:false envelope', new WriteGroupResponseError('write-group update', { success: false }, 'failed')],
  ])('%s is failed', async (_name, rejection) => {
    const result = await outcomeOf(rejection);
    expect(result.status).toBe('failed');
    expect(result.profiles.instantaneous.status).toBe('failed');
  });

  it.each([
    ['a network error without a response', Object.assign(new Error('Network Error'), { code: 'ERR_NETWORK' })],
    ['a client timeout', Object.assign(new Error('timeout of 5000ms exceeded'), { code: 'ECONNABORTED' })],
    ['HTTP 408 request timeout', httpFailure(408, 'WRITE_GROUP_INVALID_REQUEST')],
    ['a status-only HTTP 409 without any typed application error', httpFailure(409)],
    ['a status-only HTTP 400 with an empty body', httpFailure(400)],
    ['a 4xx whose typed code is not one of this backend codes', httpFailure(422, 'SOMETHING_FROM_A_GATEWAY')],
    ['HTTP 502 from a proxy', httpFailure(502)],
    ['HTTP 500', httpFailure(500, 'internal_error')],
    ['a 2xx whose body cannot be parsed', new WriteGroupResponseError('write-group update', { success: true })],
    ['an untyped exception', new TypeError('Failed to fetch')],
  ])('%s is unknown', async (_name, rejection) => {
    const result = await outcomeOf(rejection);
    expect(result.status).toBe('unknown');
    expect(result.profiles.instantaneous.status).toBe('unknown');
  });

  it('retains the original typed evidence so the operator copy can name the cause', async () => {
    const rejection = httpFailure(409, 'WRITE_GROUP_CIRCUIT_PROFILE_CONFLICT', 'pq: duplicate key postgres://user:secret@db/prod');
    const result = await outcomeOf(rejection);
    expect(result.profiles.instantaneous.cause).toBe(rejection);
    const copy = safeBasicStartError(result.profiles.instantaneous.cause, (key) => key);
    expect(copy).toBe('step4.basic.errors.circuit_profile_conflict');
    expect(copy).not.toMatch(/secret|postgres|pq:/);
  });
});
