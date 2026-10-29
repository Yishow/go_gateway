import { AxiosError, type AxiosResponse } from 'axios';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { createRecordingLifecycleMock } from '../../e2e/fixtures/recording-lifecycle-api';
import type { LifecycleScenario } from '../../e2e/fixtures/recording-lifecycle-state';

// Full production service/parser -> controller -> confirm integration, with test-only HTTP data.
// This is jsdom/adapter execution, not an actual browser or backend/SQL acceptance test.
beforeEach(() => { vi.resetModules(); localStorage.clear(); sessionStorage.clear(); });
async function createController(scenario: LifecycleScenario) {
  const [{ RecordingChangeController }, { studioV2WriteGroupSuccessorAPI }, { studioV2WriteGroupLifecycleAPI },
    { studioV2WorkspaceWriteGroupsAPI }, { studioV2DatalinkApi }] = await Promise.all([
    import('@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingChangeController'),
    import('@/services/studioV2WriteGroupSuccessor'), import('@/services/studioV2WriteGroupLifecycle'),
    import('@/services/studioV2WorkspaceWriteGroups'), import('@/services/studioV2Workspace'),
  ]);
  const mock = createRecordingLifecycleMock(scenario);
  studioV2DatalinkApi.defaults.adapter = async config => {
    const url = new URL(`${config.baseURL}${config.url}`, 'http://fixture.invalid');
    const reply = mock.handle({ method: (config.method ?? 'GET').toUpperCase(), url: url.href,
      body: typeof config.data === 'string' ? JSON.parse(config.data) : config.data });
    const response: AxiosResponse = { data: reply.body, status: reply.status, statusText: '', headers: {}, config };
    if (reply.status >= 400) throw new AxiosError('Test mock rejected', undefined, config, undefined, response);
    return response;
  };
  const controller = new RecordingChangeController('ws-1', {
    api: studioV2WriteGroupSuccessorAPI, lifecycle: studioV2WriteGroupLifecycleAPI.lifecycle,
    getGroup: studioV2WorkspaceWriteGroupsAPI.get,
  });
  return { controller, mock };
}

describe('real successor API parsing preserves controller CAS equality', () => {
  it.each([
    ['deleted-unapplied', 'recreate', 'first_start', 'first-start-confirm'],
    ['applied-wide', 'transition', 'cutover', 'activation-confirm'],
  ] as const)('%s completes explicitly confirmed %s through production services', async (scenario, action, kind, confirmAction) => {
    const { controller, mock } = await createController(scenario);
    await controller.requestAction({ action, workspaceId: 'ws-1', groupId: 'G1', workspaceRevision: 'W1', groupRevision: 'R1' });
    await controller.review(); await controller.reserveAndPreview();
    expect(controller.state.phase).toBe('schema-plan');
    expect(mock.ledger.filter(call => call.action === 'schema-confirm')).toEqual([]);
    const stage = structuredClone(controller.state.previewStage!);
    await controller.confirmSchema();
    expect(controller.state.phase).toBe('schema-confirmed');
    controller.setBoundary('2099-01-01T01:00:00Z'); await controller.reviewActivation();
    expect(controller.state.phase).toBe('activation-ready');
    const reviewed = controller.state.activation!.review.envelope;
    expect(reviewed.activation_kind).toBe(kind);
    expect(reviewed.request.stage.current).toEqual(stage.current);
    // The real parser, not a hand-reordered fake response, produces the differing key order.
    expect(Object.keys(reviewed.request.stage.current)).not.toEqual(Object.keys(stage.current));
    expect(mock.ledger.filter(call => ['first-start-confirm', 'activation-confirm'].includes(call.action))).toEqual([]);
    await controller.confirmActivation();
    expect(controller.state.error).toBeUndefined();
    expect(controller.state.phase).toBe('scheduled');
    expect(controller.state.unknown).toBe(false);
    expect(mock.snapshot().status?.[kind].status).toBe('succeeded');
    expect(mock.ledger.filter(call => call.mutation).map(call => call.action)).toEqual([
      'reserve', 'schema-preview', 'schema-confirm', confirmAction,
    ]);
    expect(mock.ledger.every(call => call.status === 200)).toBe(true);
    expect(controller.state.recovery?.envelope).toEqual(stage.envelope);
    expect(controller.state.recovery?.submitted).toEqual({
      action: kind === 'first_start' ? 'confirmFirstStart' : 'confirmCutover',
      request: { envelope: reviewed, review_digest: controller.state.activation!.review.review_digest,
        readiness_digest: controller.state.activation!.readiness_digest },
    });
  });
});
