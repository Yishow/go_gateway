import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ShareOutputSummary } from '../../../src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary';
import {
  activateStudioV2WorkspaceWithShare,
  syncStudioV2ShareMappings,
} from '../../../src/features/datalink/workbench-v2/state/studioV2ShareActivation';
import { modbusShareAPI } from '../../../src/services/datalink';
import { hydratedShareStatus, mcShareState, shareContext } from './step4-share-helpers';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (_key: string, options?: { defaultValue?: string }) => options?.defaultValue ?? _key,
  }),
}));

vi.mock('../../../src/services/datalink', () => ({
  modbusShareAPI: {
    reconcile: vi.fn(),
  },
}));

describe('Step 4 Modbus Share synchronization', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(modbusShareAPI.reconcile).mockResolvedValue({
      outcome: 'applied',
      new_workspace_revision: 'workspace-revision-2',
      settings_revision: 'settings-revision-1',
      applied_count: 4,
      removed_count: 0,
      invalidated_count: 0,
      removed_spans: [],
      invalidated_spans: [],
      mappings: [],
    });
  });

  it('uses the backend canonical plan as the reconcile source', async () => {
    await syncStudioV2ShareMappings(mcShareState(), shareContext);

    expect(modbusShareAPI.reconcile).toHaveBeenCalledWith(expect.objectContaining({
      desired_mappings: expect.arrayContaining([
        expect.objectContaining({ tag_id: 'tag-d0', zero_based_register: 0, share_start_register: 40001 }),
      ]),
    }));
  });

  it('fails closed when the canonical plan is absent', async () => {
    await expect(syncStudioV2ShareMappings(mcShareState(), {
      ...shareContext,
      canonical_plan: undefined,
    })).rejects.toMatchObject({ code: 'modbus_share_projection_required' });
    expect(modbusShareAPI.reconcile).not.toHaveBeenCalled();
  });

  it('sends the complete desired set so backend reconciliation removes stale owned rows', async () => {
    const state = mcShareState({
      mappings: {
        ...mcShareState().mappings,
        'owned-stale-point': {
          point_id: 'owned-stale-point',
          workspace_id: 'workspace-1',
          rule_id: 'mc-rule',
          device_id: 'dev-mc',
          address: 'D99',
          persisted_point_id: 'persisted-stale-point',
          tag_key: 'mc.stale',
          display_name: 'stale',
          unit: '',
          target_type: 'int16',
          scale: 1,
          offset: 0,
          enabled: false,
          tag_id: 'tag-stale-owned',
          persisted: true,
        },
      },
    });
    await syncStudioV2ShareMappings(state, shareContext);

    const [request] = vi.mocked(modbusShareAPI.reconcile).mock.calls[0];
    expect(request).toMatchObject({
      workspace_id: shareContext.workspace_id,
      readiness_token: shareContext.readiness_token,
      expected_workspace_revision: shareContext.workspace_revision,
      expected_settings_revision: shareContext.settings_revision,
    });
    expect(request.desired_mappings).not.toContainEqual(expect.objectContaining({ tag_id: 'tag-stale-owned' }));
  });

  it('fails closed when an external runtime tag already occupies a desired register', async () => {
    vi.mocked(modbusShareAPI.reconcile).mockResolvedValue({
      outcome: 'failed',
      new_workspace_revision: 'workspace-revision-1',
      settings_revision: 'settings-revision-1',
      applied_count: 0,
      removed_count: 0,
      invalidated_count: 0,
      removed_spans: [],
      invalidated_spans: [],
      mappings: [],
      diagnostics: [{ code: 'modbus_share_range_collision', severity: 'error', message: 'range collision', retryable: false }],
    });

    await expect(syncStudioV2ShareMappings(mcShareState(), shareContext)).rejects.toMatchObject({ code: 'modbus_share_range_collision' });
  });

  it('does not activate the workspace when Share synchronization fails', async () => {
    const activateWorkspace = vi.fn().mockResolvedValue({ status: 'active' });
    vi.mocked(modbusShareAPI.reconcile).mockRejectedValue(new Error('runtime unavailable'));

    await expect(
      activateStudioV2WorkspaceWithShare(mcShareState(), activateWorkspace, shareContext),
    ).rejects.toMatchObject({ code: 'modbus_share_reconcile_failed', retryable: true });
    expect(activateWorkspace).not.toHaveBeenCalled();
  });

  it('globally disabled Share hides summary and skips runtime synchronization', async () => {
    const state = mcShareState({
      settings: {
        ...mcShareState().settings,
        modbus_share: { ...mcShareState().settings.modbus_share, enabled: false },
      },
    });

    await syncStudioV2ShareMappings(state, shareContext);
    expect(modbusShareAPI.reconcile).not.toHaveBeenCalled();

    render(<ShareOutputSummary state={state} shareStatus={{ ...hydratedShareStatus, enabled: false, bind_state: 'disabled' }} />);
    expect(screen.queryByTestId('step4-share-output-summary')).not.toBeInTheDocument();
  });

  it('persisted disabled context skips runtime synchronization even when local settings are still true', async () => {
    const state = mcShareState();

    await activateStudioV2WorkspaceWithShare(state, vi.fn().mockResolvedValue({ status: 'active' }), {
      ...shareContext,
      configured_enabled: false,
    });

    expect(modbusShareAPI.reconcile).not.toHaveBeenCalled();
  });
});
