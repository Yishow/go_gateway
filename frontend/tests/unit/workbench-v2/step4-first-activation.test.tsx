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

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
/** 讓測試看得到 activateWorkspace 的結果碼（mock 工廠只能引用 mock 前綴的外部變數）。 */
const mockActivationOutcome: { code: string | null } = { code: null };
vi.mock('../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell', () => ({
  WorkbenchV2Shell: ({ activateWorkspace, actions, state }: {
    activateWorkspace?: () => Promise<unknown>;
    actions: { dispatch: (action: unknown) => void };
    state: WorkbenchV2State;
  }) => (
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
  ),
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
vi.mock('../../../src/services/studioV2WorkspaceActivation', () => ({
  StudioV2ActivationBarrierError: class StudioV2ActivationBarrierError extends Error {
    code: string;
    constructor(code: string) {
      super(code);
      this.code = code;
    }
  },
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
});
