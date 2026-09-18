import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import { modbusShareAPI } from '../../../src/services/datalink';
import { studioV2WorkspaceActivationAPI } from '../../../src/services/studioV2WorkspaceActivation';
import { studioV2WorkspaceAPI } from '../../../src/services/studioV2Workspace';
import { studioV2WorkspaceAuditAPI } from '../../../src/services/studioV2WorkspaceAudit';
import { studioV2MappingsAPI } from '../../../src/services/studioV2Mappings';
import { studioV2RulesAPI } from '../../../src/services/studioV2Rules';
import { studioV2WorkspaceDatabaseAPI } from '../../../src/services/studioV2WorkspaceDatabase';
import { studioV2WorkspaceDevicesAPI } from '../../../src/services/studioV2WorkspaceDevices';
import { deviceFixture, mappingFixture, ruleFixture, workspaceFixture } from './helpers/workbenchV2PageHarness';
import { shareContext } from './step4-share-helpers';
import { AUTOSAVE_SETTLEMENT_TIMEOUT_MS } from '../../../src/pages/datalink/workbench-v2/studioV2AutosaveBarrier';
import { studioV2RuntimeContextAPI } from '../../../src/services/studioV2RuntimeContext';
import type { StudioV2ActivationRecovery } from '../../../src/types/studioV2Activation';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
/** 讓測試看得到 activateWorkspace 的結果碼（mock 工廠只能引用 mock 前綴的外部變數）。 */
const mockActivationOutcome: { code: string | null } = { code: null };
const mockShellCallbacks: {
  activate?: () => Promise<unknown>;
  recover?: () => Promise<StudioV2ActivationRecovery>;
} = {};
vi.mock('../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell', () => ({
  WorkbenchV2Shell: ({ activateWorkspace, recoverActivationStatus, actions, state }: {
    activateWorkspace?: () => Promise<unknown>;
    recoverActivationStatus?: () => Promise<StudioV2ActivationRecovery>;
    actions: { dispatch: (action: unknown) => void };
    state: WorkbenchV2State;
  }) => {
    mockShellCallbacks.activate = activateWorkspace;
    mockShellCallbacks.recover = recoverActivationStatus;
    return (
    <>
      <button
        type="button"
        data-testid="edit-database"
        onClick={() => actions.dispatch({
          type: 'updateDbConnector',
          patch: { name: 'Updated connector' },
        })}
      >
        Edit database
      </button>
      <span data-testid="database-save-state">{state.db.connector.save_state}</span>
      <button type="button" data-testid="activate" onClick={() => void activateWorkspace?.()}>
        Activate
      </button>
      <button
        type="button"
        data-testid="activate-capture"
        onClick={() => {
          void activateWorkspace?.()
            .then(() => { mockActivationOutcome.code = 'success'; })
            .catch((error: { code?: string }) => { mockActivationOutcome.code = error?.code ?? 'unknown'; });
        }}
      >
        Activate and capture
      </button>
    </>
    );
  },
}));
vi.mock('../../../src/services/datalink', () => ({
  modbusShareAPI: { status: vi.fn(), reconcile: vi.fn() },
}));
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('../../../src/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn() } }));
vi.mock('../../../src/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
vi.mock('../../../src/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
vi.mock('../../../src/services/studioV2WorkspaceDatabase', () => ({
  studioV2WorkspaceDatabaseAPI: {
    getConfig: vi.fn(),
    updateConfig: vi.fn(),
    listTargets: vi.fn(),
    upsertTarget: vi.fn(),
  },
}));
vi.mock('../../../src/services/studioV2WorkspaceAudit', () => ({ studioV2WorkspaceAuditAPI: { list: vi.fn() } }));
vi.mock('../../../src/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { get: vi.fn() } }));
vi.mock('../../../src/services/studioV2WorkspaceActivation', async (importOriginal) => ({
  ...await importOriginal<typeof import('../../../src/services/studioV2WorkspaceActivation')>(),
  studioV2WorkspaceActivationAPI: { activate: vi.fn() },
}));

function renderPage() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(<QueryClientProvider client={queryClient}><DatalinkWorkbenchV2Page /></QueryClientProvider>);
}

describe('Step 4 first activation', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockActivationOutcome.code = null;
    mockShellCallbacks.activate = undefined;
    mockShellCallbacks.recover = undefined;
    vi.mocked(studioV2RuntimeContextAPI.get).mockReset();
    vi.mocked(studioV2RuntimeContextAPI.get).mockResolvedValue({
      workspace_id: 'workspace-1', default_device_id: 'dev-mc',
      devices: [{ device_id: 'dev-mc', name: 'MC', protocol: 'mc_3e', running: true, availability_status: 'available' }],
    });
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspaceFixture({ ordered_device_ids: ['dev-mc'] }));
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([deviceFixture({ id: 'dev-mc', protocol: 'mc_3e' })]);
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([ruleFixture({
      id: 'mc-rule', device_id: 'dev-mc', start_address: 'D0', count: 1, naming_prefix: 'D_',
      share_enabled: true, share_start_register: 40001, share_stride: 1,
    })]);
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([mappingFixture({
      id: 'mapping-0', point_id: 'persisted-point-0', rule_id: 'mc-rule', device_id: 'dev-mc',
      address: 'D0', tag_id: 'tag-d0',
    })]);
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(null);
    vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);
    vi.mocked(studioV2WorkspaceAuditAPI.list).mockResolvedValue({ entries: [] });
    vi.mocked(studioV2WorkspaceActivationAPI.activate).mockResolvedValue({ workspace_id: 'workspace-1', results: [] });
    vi.mocked(modbusShareAPI.status).mockResolvedValue({
      enabled: true,
      port: 5020,
      address: '0.0.0.0:5020',
      bind_state: 'pass',
      mapping_count: 1,
      hydration_state: 'ready',
      readiness: true,
      readiness_token: 'ready-token-1',
      workspace_revision: 'workspace-revision-1',
      settings_revision: 'settings-revision-1',
      canonical_plan: shareContext.canonical_plan,
    });
    vi.mocked(modbusShareAPI.reconcile).mockResolvedValue({
      outcome: 'applied',
      new_workspace_revision: 'workspace-revision-2',
      new_readiness_token: 'ready-token-2',
      settings_revision: 'settings-revision-1',
      applied_count: 1,
      removed_count: 0,
      invalidated_count: 0,
      removed_spans: [],
      invalidated_spans: [],
      mappings: [],
    });
  });

  it('waits for an active autosave instead of requiring a retry', async () => {
    let resolveSave!: (value: Awaited<ReturnType<typeof studioV2WorkspaceDatabaseAPI.updateConfig>>) => void;
    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockImplementation(() => new Promise((resolve) => {
      resolveSave = resolve;
    }));

    renderPage();
    await screen.findByTestId('activate');
    fireEvent.click(screen.getByTestId('edit-database'));
    await waitFor(() => expect(screen.getByTestId('database-save-state')).toHaveTextContent('saving'));
    fireEvent.click(screen.getByTestId('activate'));
    expect(studioV2WorkspaceActivationAPI.activate).not.toHaveBeenCalled();

    resolveSave({
      id: 'database-config-1', workspace_id: 'workspace-1', kind: 'postgres', name: 'Updated connector',
      host: 'tsdb.internal', port: 5432, database: 'gateway_metrics', username: 'gw_writer',
      schema: 'public', table: 'sensor_readings', write_mode: 'insert', write_interval_seconds: 5,
      timestamp_column: 'ts', status: 'ready', save_state: 'saved', created_at: '2026-08-26T00:00:00Z',
      updated_at: '2026-08-26T00:00:00Z', runtime_apply_status: 'not_running',
    });

    await waitFor(() => expect(studioV2WorkspaceActivationAPI.activate).toHaveBeenCalledTimes(1));
  });

  it('fails activation as save-incomplete when an autosave never settles', async () => {
    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockImplementation(() => new Promise(() => undefined));

    renderPage();
    await screen.findByTestId('activate-capture');
    fireEvent.click(screen.getByTestId('edit-database'));
    await waitFor(() => expect(screen.getByTestId('database-save-state')).toHaveTextContent('saving'));

    vi.useFakeTimers();
    try {
      fireEvent.click(screen.getByTestId('activate-capture'));
      await act(async () => { await vi.advanceTimersByTimeAsync(AUTOSAVE_SETTLEMENT_TIMEOUT_MS); });
    } finally {
      vi.useRealTimers();
    }

    await waitFor(() => expect(mockActivationOutcome.code).toBe('modbus_share_save_incomplete'));
    expect(studioV2WorkspaceActivationAPI.activate).not.toHaveBeenCalled();
  });

  it('blocks concurrent page callbacks before autosave and Share reconciliation can submit twice', async () => {
    renderPage();
    await screen.findByTestId('activate');
    const activate = mockShellCallbacks.activate!;
    const results = await act(async () => Promise.allSettled([activate(), activate()]));
    expect(results[0].status).toBe('fulfilled');
    expect(results[1]).toMatchObject({ status: 'rejected', reason: { outcome: 'unconfirmed', retryable: false } });
    expect(studioV2WorkspaceActivationAPI.activate).toHaveBeenCalledTimes(1);
    expect(modbusShareAPI.reconcile).toHaveBeenCalledTimes(1);
  });

  it('recovers fresh workspace and runtime context without activating or overwriting setup', async () => {
    vi.mocked(studioV2RuntimeContextAPI.get).mockResolvedValue({
      workspace_id: 'workspace-1', default_device_id: 'dev-mc', operation_id: 'op-server-1',
      devices: [{ device_id: 'dev-mc', running: true }],
    } as never);
    renderPage();
    await screen.findByTestId('activate');
    expect(mockShellCallbacks.recover).toBeTypeOf('function');
    const beforeReads = vi.mocked(studioV2WorkspaceAPI.get).mock.calls.length;
    const recovered = await act(async () => mockShellCallbacks.recover!());
    expect(recovered).toEqual({ workspace_id: 'workspace-1', operation_id: 'op-server-1',
      devices: [{ device_id: 'dev-mc', running: true }],
    });
    expect(studioV2WorkspaceAPI.get).toHaveBeenCalledTimes(beforeReads + 1);
    expect(studioV2RuntimeContextAPI.get).toHaveBeenCalledTimes(1);
    expect(studioV2WorkspaceActivationAPI.activate).not.toHaveBeenCalled();
    expect(studioV2WorkspaceDatabaseAPI.updateConfig).not.toHaveBeenCalled();
  });

  it('releases the page guard so an explicit later attempt can run after rejection', async () => {
    vi.mocked(studioV2WorkspaceActivationAPI.activate).mockRejectedValueOnce({
      code: 'modbus_share_revision_conflict', retryable: true,
    });
    renderPage();
    await screen.findByTestId('activate');
    await act(async () => {
      await expect(mockShellCallbacks.activate!()).rejects.toMatchObject({ code: 'modbus_share_revision_conflict' });
    });
    await act(async () => {
      await expect(mockShellCallbacks.activate!()).resolves.toMatchObject({ workspace_id: 'workspace-1' });
    });
    expect(studioV2WorkspaceActivationAPI.activate).toHaveBeenCalledTimes(2);
  });

  it('does not turn cached workspace data into confirmed recovery after a failed refresh', async () => {
    renderPage();
    await screen.findByTestId('activate');
    expect(mockShellCallbacks.recover).toBeTypeOf('function');
    vi.mocked(studioV2WorkspaceAPI.get).mockRejectedValueOnce(new Error('raw backend diagnostic'));
    await act(async () => {
      await expect(mockShellCallbacks.recover!()).rejects.toMatchObject({ outcome: 'unconfirmed', retryable: false });
    });
    expect(studioV2WorkspaceActivationAPI.activate).not.toHaveBeenCalled();
  });

  it.each([
    { workspace_id: 'another-workspace', devices: [{ device_id: 'dev-mc', running: true }] },
    { workspace_id: 'workspace-1', devices: [{ device_id: 'dev-mc', running: 'true' }] },
  ])('rejects a foreign or malformed runtime recovery: %j', async (runtime) => {
    vi.mocked(studioV2RuntimeContextAPI.get).mockResolvedValue(runtime as never);
    renderPage();
    await screen.findByTestId('activate');
    expect(mockShellCallbacks.recover).toBeTypeOf('function');
    await act(async () => {
      await expect(mockShellCallbacks.recover!()).rejects.toMatchObject({ outcome: 'unconfirmed', retryable: false });
    });
    expect(studioV2WorkspaceActivationAPI.activate).not.toHaveBeenCalled();
  });
});
