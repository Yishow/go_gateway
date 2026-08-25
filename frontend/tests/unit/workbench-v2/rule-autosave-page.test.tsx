import { fireEvent, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2WorkspaceAPI } from '../../../src/services/studioV2Workspace';
import { studioV2WorkspaceDevicesAPI } from '../../../src/services/studioV2WorkspaceDevices';
import { studioV2MappingsAPI } from '../../../src/services/studioV2Mappings';
import { studioV2RulesAPI } from '../../../src/services/studioV2Rules';
import { studioV2WorkspaceDatabaseAPI } from '../../../src/services/studioV2WorkspaceDatabase';
import type { WorkbenchV2ShellMockProps } from './helpers/workbenchV2PageHarness';
import {
  deviceFixture,
  renderWorkbenchV2Page as renderPage,
  ruleFixture,
  runtimeApplied,
  workspaceFixture,
} from './helpers/workbenchV2PageHarness';

vi.mock('../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell', () => ({
  WorkbenchV2Shell: ({
    state,
    actions,
  }: WorkbenchV2ShellMockProps<{
    rules: Array<{
      id: string;
      device_id: string;
      naming_prefix: string;
      save_state: string;
      save_error?: string | null;
      share_enabled: boolean;
      share_start_register: number | null;
      share_stride: number | null;
    }>;
    }>) => (
    <div data-testid="rule-autosave-shell">
      <div data-testid="rule-ownership">{state.rules.map((rule) => `${rule.id}:${rule.device_id}`).join(',')}</div>
      {state.rules.map((rule) => (
        <div key={rule.id}>
          <div data-testid={`rule-prefix-${rule.id}`}>{rule.naming_prefix}</div>
          <div data-testid={`rule-save-state-${rule.id}`}>{rule.save_state}</div>
          <div data-testid={`rule-save-error-${rule.id}`}>{rule.save_error ?? ''}</div>
          <div data-testid={`rule-share-${rule.id}`}>
            {String(rule.share_enabled)}:{String(rule.share_start_register)}:{String(rule.share_stride)}
          </div>
        </div>
      ))}
      <button
        type="button"
        data-testid="seed-rule-01"
        onClick={() => actions.dispatch({
          type: 'addRule',
          rule: {
            id: 'rule-01',
            device_id: 'dev-01',
            name: 'Holding Registers',
            start_address: '',
            count: 8,
            data_type: 'int16',
            naming_prefix: 'LINE_',
            enabled: true,
            scale_multiplier: 1,
            scale_offset: 0,
            data_format: '',
            skipped_addresses: [],
            share_enabled: false,
            share_start_register: null,
            share_stride: null,
            persisted: false,
            save_state: 'idle',
            save_error: null,
          },
        })}
      />
      <button
        type="button"
        data-testid="make-rule-01-valid"
        onClick={() => actions.dispatch({
          type: 'updateRule',
          ruleId: 'rule-01',
          patch: { start_address: '40001', naming_prefix: 'LINE_' },
        })}
      />
      <button
        type="button"
        data-testid="make-rule-01-invalid"
        onClick={() => actions.dispatch({
          type: 'updateRule',
          ruleId: 'rule-01',
          patch: { start_address: '' },
        })}
      />
      <button
        type="button"
        data-testid="make-rule-A-fail"
        onClick={() => actions.dispatch({
          type: 'updateRule',
          ruleId: 'rule-A',
          patch: { naming_prefix: 'BROKEN_' },
        })}
      />
      <button
        type="button"
        data-testid="make-rule-B-valid"
        onClick={() => actions.dispatch({
          type: 'updateRule',
          ruleId: 'rule-B',
          patch: { naming_prefix: 'B_SAVED_' },
        })}
      />
      <button
        type="button"
        data-testid="patch-rule-A-share"
        onClick={() => actions.dispatch({
          type: 'updateRule',
          ruleId: 'rule-A',
          patch: { share_enabled: true, share_start_register: 40001, share_stride: 1 },
        })}
      />
      <button
        type="button"
        data-testid="toggle-rule-A-share"
        onClick={() => actions.dispatch({ type: 'toggleRuleShareEnabled', ruleId: 'rule-A' })}
      />
      <button
        type="button"
        data-testid="update-rule-A-share-start"
        onClick={() => actions.dispatch({ type: 'updateRuleShareStart', ruleId: 'rule-A', shareStart: 40011 })}
      />
      <button
        type="button"
        data-testid="update-rule-A-share-stride"
        onClick={() => actions.dispatch({ type: 'updateRuleShareStride', ruleId: 'rule-A', shareStride: 2 })}
      />
    </div>
  ),
}));

vi.mock('../../../src/services/studioV2Workspace', async () => {
  const { studioV2ServiceMocks } = await import('./helpers/studioV2ServiceMocks');
  return { studioV2WorkspaceAPI: studioV2ServiceMocks.workspace };
});

vi.mock('../../../src/services/studioV2WorkspaceDevices', async () => {
  const { studioV2ServiceMocks } = await import('./helpers/studioV2ServiceMocks');
  return { studioV2WorkspaceDevicesAPI: studioV2ServiceMocks.devices };
});

vi.mock('../../../src/services/studioV2Rules', async () => {
  const { studioV2ServiceMocks } = await import('./helpers/studioV2ServiceMocks');
  return { studioV2RulesAPI: studioV2ServiceMocks.rules };
});

vi.mock('../../../src/services/studioV2Mappings', async () => {
  const { studioV2ServiceMocks } = await import('./helpers/studioV2ServiceMocks');
  return { studioV2MappingsAPI: studioV2ServiceMocks.mappings };
});

vi.mock('../../../src/services/studioV2WorkspaceDatabase', async () => {
  const { studioV2ServiceMocks } = await import('./helpers/studioV2ServiceMocks');
  return { studioV2WorkspaceDatabaseAPI: studioV2ServiceMocks.database };
});

describe('DatalinkWorkbenchV2Page rule autosave orchestration', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(
      workspaceFixture({ ordered_device_ids: ['dev-01'] }),
    );
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      deviceFixture({ id: 'dev-01' }),
    ]);
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(null);
    vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);
  });

  it('hydrates persisted rules under the same owning devices on reload', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce(
      workspaceFixture({ ordered_device_ids: ['dev-B', 'dev-A'] }),
    );
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      deviceFixture({ id: 'dev-B', name: 'Line B PLC', connection_config: '{"host":"192.168.10.11","port":502,"slave_id":2,"timeout":5}' }),
      deviceFixture(),
    ]);
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([
      ruleFixture({
        id: 'rule-B',
        device_id: 'dev-B',
        start_address: '40011',
        naming_prefix: 'B_',
        revision_id: 'rev-B',
      }),
      ruleFixture(),
    ]);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('rule-ownership')).toHaveTextContent('rule-B:dev-B,rule-A:dev-A');
    });
  });

  it('saves one valid draft rule and marks it saved', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce(
      workspaceFixture({ ordered_device_ids: ['dev-01'] }),
    );
    vi.mocked(studioV2RulesAPI.create).mockResolvedValueOnce(
      runtimeApplied(ruleFixture({
        id: 'rule-01',
        device_id: 'dev-01',
        start_address: '40001',
        count: 8,
        naming_prefix: 'LINE_',
        revision_id: 'rev-1',
      })),
    );

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('seed-rule-01'));
    fireEvent.click(screen.getByTestId('make-rule-01-valid'));

    await waitFor(() => {
      expect(studioV2RulesAPI.create).toHaveBeenCalledWith(
        expect.objectContaining({
          id: 'rule-01',
          device_id: 'dev-01',
          naming_prefix: 'LINE_',
        }),
      );
    });

    await waitFor(() => {
      expect(screen.getByTestId('rule-save-state-rule-01')).toHaveTextContent('saved');
    });
  });

  it('keeps invalid local rule edits without overwriting the backend', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce(
      workspaceFixture({ ordered_device_ids: ['dev-01'] }),
    );

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('seed-rule-01'));
    fireEvent.click(screen.getByTestId('make-rule-01-invalid'));

    await waitFor(() => {
      expect(screen.getByTestId('rule-save-state-rule-01')).toHaveTextContent('draft-invalid');
    });

    expect(studioV2RulesAPI.create).not.toHaveBeenCalled();
    expect(studioV2RulesAPI.update).not.toHaveBeenCalled();
  });

  it('isolates one rule save failure without blocking another valid rule', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce(
      workspaceFixture({ ordered_device_ids: ['dev-A', 'dev-B'] }),
    );
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      deviceFixture(),
      deviceFixture({ id: 'dev-B', name: 'Line B PLC', connection_config: '{"host":"192.168.10.11","port":502,"slave_id":2,"timeout":5}' }),
    ]);
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([
      ruleFixture(),
      ruleFixture({
        id: 'rule-B',
        device_id: 'dev-B',
        start_address: '40011',
        naming_prefix: 'B_',
        revision_id: 'rev-B',
      }),
    ]);
    vi.mocked(studioV2RulesAPI.update).mockImplementation(async (ruleId, payload) => {
      if (ruleId === 'rule-A') {
        throw new Error('save failed');
      }

      return runtimeApplied(ruleFixture({
        id: ruleId,
        device_id: String(payload.device_id),
        start_address: '40011',
        naming_prefix: String(payload.naming_prefix),
        revision_id: `rev-${ruleId}`,
      }));
    });

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('rule-ownership')).toHaveTextContent('rule-A:dev-A,rule-B:dev-B');
    });

    fireEvent.click(screen.getByTestId('make-rule-A-fail'));
    fireEvent.click(screen.getByTestId('make-rule-B-valid'));

    await waitFor(() => {
      expect(studioV2RulesAPI.update).toHaveBeenCalledWith(
        'rule-A',
        expect.objectContaining({ device_id: 'dev-A', naming_prefix: 'BROKEN_' }),
      );
      expect(studioV2RulesAPI.update).toHaveBeenCalledWith(
        'rule-B',
        expect.objectContaining({ device_id: 'dev-B', naming_prefix: 'B_SAVED_' }),
      );
    });

    await waitFor(() => {
      expect(screen.getByTestId('rule-save-state-rule-A')).toHaveTextContent('save-error');
      expect(screen.getByTestId('rule-save-state-rule-B')).toHaveTextContent('saved');
      expect(screen.getByTestId('rule-save-error-rule-A')).toHaveTextContent('errors.autosave_failed');
    });
  });

  it('hydrates persisted Share fields and autosaves a direct Share patch', async () => {
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([
      ruleFixture({ share_enabled: true, share_start_register: 40001, share_stride: 1 }),
    ]);
    vi.mocked(studioV2RulesAPI.update).mockResolvedValue(
      runtimeApplied(ruleFixture({ share_enabled: true, share_start_register: 40001, share_stride: 1 })),
    );

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('rule-share-rule-A')).toHaveTextContent('true:40001:1');
    });

    fireEvent.click(screen.getByTestId('patch-rule-A-share'));

    await waitFor(() => {
      expect(studioV2RulesAPI.update).toHaveBeenCalledWith(
        'rule-A',
        expect.objectContaining({
          share_enabled: true,
          share_start_register: 40001,
          share_stride: 1,
        }),
      );
    });
  });

  it('autosaves toggle, start-register, and stride Share actions', async () => {
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([
      ruleFixture({ share_enabled: false, share_start_register: null, share_stride: null }),
    ]);
    vi.mocked(studioV2RulesAPI.update).mockImplementation(async (_ruleId, request) =>
      runtimeApplied(ruleFixture({
        share_enabled: request.share_enabled ?? false,
        share_start_register: request.share_start_register ?? null,
        share_stride: request.share_stride ?? null,
      })),
    );

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('rule-share-rule-A')).toHaveTextContent('false:null:null');
    });

    fireEvent.click(screen.getByTestId('toggle-rule-A-share'));
    await waitFor(() => {
      expect(studioV2RulesAPI.update).toHaveBeenCalledWith(
        'rule-A',
        expect.objectContaining({ share_enabled: true }),
      );
    });

    fireEvent.click(screen.getByTestId('update-rule-A-share-start'));
    await waitFor(() => {
      expect(studioV2RulesAPI.update).toHaveBeenCalledWith(
        'rule-A',
        expect.objectContaining({ share_start_register: 40011 }),
      );
    });

    fireEvent.click(screen.getByTestId('update-rule-A-share-stride'));

    await waitFor(() => {
      expect(studioV2RulesAPI.update).toHaveBeenCalledWith(
        'rule-A',
        expect.objectContaining({ share_stride: 2 }),
      );
    });
  });
});
