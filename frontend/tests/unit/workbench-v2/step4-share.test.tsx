import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import { ShareOutputSummary } from '../../../src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary';
import {
  activateStudioV2WorkspaceWithShare,
  buildStudioV2ShareSyncPlan,
  syncStudioV2ShareMappings,
} from '../../../src/features/datalink/workbench-v2/state/studioV2ShareActivation';
import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { modbusShareAPI } from '../../../src/services/datalink';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (_key: string, options?: { defaultValue?: string }) => options?.defaultValue ?? _key,
  }),
}));

vi.mock('../../../src/services/datalink', () => ({
  modbusShareAPI: {
    listMappings: vi.fn(),
    deleteMapping: vi.fn(),
    upsertMapping: vi.fn(),
  },
}));


function mcShareState(overrides: Partial<WorkbenchV2State> = {}): WorkbenchV2State {
  const points = Array.from({ length: 4 }, (_, index) => ({
    id: `mc-rule-p-${index}`,
    device_id: 'dev-mc',
    rule_id: 'mc-rule',
    rule_name: 'MC D registers',
    name: `D_${index}`,
    address: `D${index}`,
    data_type: 'int16',
    function: 'D (Word)',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0,
  }));
  const mappings = Object.fromEntries(points.map((point, index) => [point.id, {
    point_id: point.id,
    workspace_id: 'workspace-1',
    rule_id: point.rule_id,
    device_id: point.device_id,
    address: point.address,
    persisted_point_id: `persisted-point-${index}`,
    tag_key: `mc.d${index}`,
    display_name: `D${index}`,
    unit: '',
    target_type: 'int16' as const,
    scale: 1,
    offset: 0,
    enabled: true,
    tag_id: `tag-d${index}`,
    persisted: true,
  }]));

  return {
    ...INITIAL_STATE,
    devices: [{ ...INITIAL_STATE.devices[0], id: 'dev-mc', protocol: 'mc_3e', name: 'MC PLC' }],
    rules: [{
      ...INITIAL_STATE.rules[0],
      id: 'mc-rule',
      device_id: 'dev-mc',
      workspace_id: 'workspace-1',
      persisted: true,
      name: 'MC D registers',
      start_address: 'D0',
      count: 4,
      share_enabled: true,
      share_start_register: 40001,
      share_stride: 1,
    }],
    points,
    mappings,
    db: { ...INITIAL_STATE.db, targets: {} },
    ...overrides,
  };
}

describe('Step 4 Modbus Share synchronization', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(modbusShareAPI.listMappings).mockResolvedValue([]);
    vi.mocked(modbusShareAPI.deleteMapping).mockResolvedValue();
    vi.mocked(modbusShareAPI.upsertMapping).mockImplementation(async (tagId, register) => ({
      tag_id: tagId,
      register,
      data_type: 'int16',
      updated_at: '2026-08-24T00:00:00Z',
    }));
  });

  it('converts human holding register 40001 to API register 0', () => {
    const plan = buildStudioV2ShareSyncPlan(mcShareState());

    expect(plan.issues).toEqual([]);
    expect(plan.desired).toEqual([
      { tag_id: 'tag-d0', register: 0 },
      { tag_id: 'tag-d1', register: 1 },
      { tag_id: 'tag-d2', register: 2 },
      { tag_id: 'tag-d3', register: 3 },
    ]);
  });

  it('deletes stale owned mappings while preserving external runtime mappings', async () => {
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
    vi.mocked(modbusShareAPI.listMappings).mockResolvedValue([
      { tag_id: 'tag-stale-owned', register: 99, data_type: 'int16', updated_at: '' },
      { tag_id: 'tag-external', register: 98, data_type: 'int16', updated_at: '' },
    ]);

    await syncStudioV2ShareMappings(state);

    expect(modbusShareAPI.deleteMapping).toHaveBeenCalledWith('tag-stale-owned');
    expect(modbusShareAPI.deleteMapping).not.toHaveBeenCalledWith('tag-external');
    expect(modbusShareAPI.upsertMapping).toHaveBeenCalledWith('tag-d0', 0);
  });

  it('does not treat an unproven frontend tag_id as durable workspace ownership', async () => {
    const state = mcShareState({
      mappings: {
        ...mcShareState().mappings,
        'frontend-only-stale': {
          point_id: 'frontend-only-stale',
          tag_key: 'mc.frontend-only',
          display_name: 'frontend-only',
          unit: '',
          target_type: 'int16',
          scale: 1,
          offset: 0,
          enabled: false,
          tag_id: 'tag-frontend-only',
        },
      },
    });
    vi.mocked(modbusShareAPI.listMappings).mockResolvedValue([
      { tag_id: 'tag-frontend-only', register: 99, data_type: 'int16', updated_at: '' },
    ]);

    await syncStudioV2ShareMappings(state);

    expect(modbusShareAPI.deleteMapping).not.toHaveBeenCalledWith('tag-frontend-only');
  });

  it('fails closed before listing runtime mappings when a required persisted tag is missing', async () => {
    const state = mcShareState({
      mappings: {
        ...mcShareState().mappings,
        'mc-rule-p-2': { ...mcShareState().mappings['mc-rule-p-2'], tag_id: undefined, persisted: false },
      },
    });

    await expect(syncStudioV2ShareMappings(state)).rejects.toThrow(/tag_id/i);
    expect(modbusShareAPI.listMappings).not.toHaveBeenCalled();
    expect(modbusShareAPI.deleteMapping).not.toHaveBeenCalled();
    expect(modbusShareAPI.upsertMapping).not.toHaveBeenCalled();
  });

  it('fails closed before runtime writes when two desired rows use one register', async () => {
    const state = mcShareState({
      rules: [{ ...mcShareState().rules[0], share_stride: 0 }],
    });

    await expect(syncStudioV2ShareMappings(state)).rejects.toThrow(/duplicate register/i);
    expect(modbusShareAPI.listMappings).not.toHaveBeenCalled();
    expect(modbusShareAPI.deleteMapping).not.toHaveBeenCalled();
    expect(modbusShareAPI.upsertMapping).not.toHaveBeenCalled();
  });

  it('fails closed when an external runtime tag already occupies a desired register', async () => {
    vi.mocked(modbusShareAPI.listMappings).mockResolvedValue([
      { tag_id: 'tag-external', register: 0, data_type: 'int16', updated_at: '' },
    ]);

    await expect(syncStudioV2ShareMappings(mcShareState())).rejects.toThrow(/external tag/i);
    expect(modbusShareAPI.deleteMapping).not.toHaveBeenCalled();
    expect(modbusShareAPI.upsertMapping).not.toHaveBeenCalled();
  });

  it('does not activate the workspace when Share synchronization fails', async () => {
    const activateWorkspace = vi.fn().mockResolvedValue({ status: 'active' });
    vi.mocked(modbusShareAPI.listMappings).mockRejectedValue(new Error('runtime unavailable'));

    await expect(
      activateStudioV2WorkspaceWithShare(mcShareState(), activateWorkspace),
    ).rejects.toThrow('runtime unavailable');
    expect(activateWorkspace).not.toHaveBeenCalled();
  });

  it('fails closed when one tag id is assigned to multiple desired points', () => {
    const state = mcShareState();
    const duplicatedTagID = state.mappings['mc-rule-p-0'].tag_id;
    const nextMappings = {
      ...state.mappings,
      'mc-rule-p-1': { ...state.mappings['mc-rule-p-1'], tag_id: duplicatedTagID },
    };

    expect(buildStudioV2ShareSyncPlan({ ...state, mappings: nextMappings }).issues.join(' ')).toMatch(/duplicate tag_id/i);
  });

  it('rejects a Share register offset beyond the uint16 runtime range', () => {
    const state = mcShareState({
      rules: [{ ...mcShareState().rules[0], share_start_register: 105537 }],
    });

    expect(buildStudioV2ShareSyncPlan(state).issues.join(' ')).toMatch(/invalid Share holding register/i);
  });

  it('globally disabled Share hides summary and skips runtime synchronization', async () => {
    const state = mcShareState({
      settings: {
        ...mcShareState().settings,
        modbus_share: { ...mcShareState().settings.modbus_share, enabled: false },
      },
    });

    expect(buildStudioV2ShareSyncPlan(state)).toEqual({ desired: [], ownedTagIds: [], issues: [] });
    await syncStudioV2ShareMappings(state);
    expect(modbusShareAPI.listMappings).not.toHaveBeenCalled();

    render(<ShareOutputSummary state={state} />);
    expect(screen.queryByTestId('step4-share-output-summary')).not.toBeInTheDocument();
  });
});

describe('Step 4 Share output summary', () => {
  it('retains the non-Modbus Step4 destination surface while Share is configured', () => {
    render(<Step4Database state={mcShareState()} dispatch={vi.fn()} />);

    expect(screen.getByTestId('step4-destination-overview')).toBeInTheDocument();
  });

  it('shows non-Modbus MC D0~D3 as 40001 ~ 40004 without DB targets', () => {
    render(<ShareOutputSummary state={mcShareState()} />);

    expect(screen.getByTestId('step4-share-output-summary')).toHaveTextContent('D0 ~ D3');
    expect(screen.getByTestId('step4-share-output-summary')).toHaveTextContent('40001 ~ 40004');
  });
});
