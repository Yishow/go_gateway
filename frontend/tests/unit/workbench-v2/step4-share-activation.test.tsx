import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { activateStudioV2WorkspaceWithShare, syncStudioV2ShareMappings } from '../../../src/features/datalink/workbench-v2/state/studioV2ShareActivation';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
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

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell', () => ({
  WorkbenchV2Shell: ({ activateWorkspace }: { activateWorkspace?: () => Promise<unknown> }) => {
    const [errorCode, setErrorCode] = React.useState<string | null>(null);
    return <>
      <button
        type="button"
        data-testid="share-page-activate"
        onClick={() => void activateWorkspace?.().catch((error: { code?: string; message?: string }) => setErrorCode(error.code ?? error.message ?? 'unknown'))}
      >
        Activate
      </button>
      {errorCode && <span data-testid="share-page-activation-error">{errorCode}</span>}
    </>;
  },
}));
vi.mock('../../../src/services/datalink', () => ({
  modbusShareAPI: { status: vi.fn(), reconcile: vi.fn() },
}));
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('../../../src/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn() } }));
vi.mock('../../../src/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
vi.mock('../../../src/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
vi.mock('../../../src/services/studioV2WorkspaceDatabase', () => ({ studioV2WorkspaceDatabaseAPI: { getConfig: vi.fn(), listTargets: vi.fn() } }));
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

function setupPageMocks() {
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspaceFixture({ ordered_device_ids: ['dev-mc'] }));
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([deviceFixture({ id: 'dev-mc', protocol: 'mc_3e' })]);
  vi.mocked(studioV2RulesAPI.list).mockResolvedValue([ruleFixture({
    id: 'mc-rule', device_id: 'dev-mc', start_address: 'D0', count: 4, naming_prefix: 'D_',
    share_enabled: true, share_start_register: 40001, share_stride: 1,
  })]);
  vi.mocked(studioV2MappingsAPI.list).mockResolvedValue(Array.from({ length: 4 }, (_, index) => mappingFixture({
    id: `mapping-${index}`, point_id: `persisted-point-${index}`, rule_id: 'mc-rule', device_id: 'dev-mc',
    address: `D${index}`, tag_id: `tag-d${index}`,
  })));
  vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(null);
  vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);
  vi.mocked(studioV2WorkspaceAuditAPI.list).mockResolvedValue({ entries: [] });
  vi.mocked(studioV2WorkspaceActivationAPI.activate).mockResolvedValue({ workspace_id: 'workspace-1', results: [] });
  vi.mocked(modbusShareAPI.status).mockResolvedValue({
    enabled: true,
    port: 5020,
    address: '0.0.0.0:5020',
    bind_state: 'pass',
    mapping_count: 4,
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
    settings_revision: 'settings-revision-1',
    applied_count: 4,
    removed_count: 0,
    invalidated_count: 0,
    removed_spans: [],
    invalidated_spans: [],
    mappings: [],
  });
}

describe('Step 4 activation Share handoff', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupPageMocks();
  });

  it('awaits runtime Share synchronization before workspace activation', async () => {
    const events: string[] = [];
    vi.mocked(modbusShareAPI.reconcile).mockImplementation(async (request) => {
      events.push(`reconcile:${request.desired_mappings.length}`);
      return {
        outcome: 'applied',
        new_workspace_revision: 'workspace-revision-2',
        new_readiness_token: 'ready-token-2',
        settings_revision: 'settings-revision-1',
        applied_count: 4,
        removed_count: 0,
        invalidated_count: 0,
        removed_spans: [],
        invalidated_spans: [],
        mappings: [],
      };
    });
    vi.mocked(studioV2WorkspaceActivationAPI.activate).mockImplementation(async () => {
      events.push('activate');
      return { workspace_id: 'workspace-1', results: [] };
    });

    renderPage();
    await waitFor(() => expect(screen.getByTestId('share-page-activate')).toBeInTheDocument());
    fireEvent.click(screen.getByTestId('share-page-activate'));

    await waitFor(() => expect(studioV2WorkspaceActivationAPI.activate).toHaveBeenCalledTimes(1));
    expect(events[0]).toBe('reconcile:4');
    expect(events.at(-1)).toBe('activate');
    expect(studioV2WorkspaceActivationAPI.activate).toHaveBeenCalledWith({
      workspace_revision: 'workspace-revision-2',
      settings_revision: 'settings-revision-1',
      readiness_token: 'ready-token-2',
      pending_saves: 0,
    });
    expect(modbusShareAPI.reconcile).toHaveBeenCalledWith(expect.objectContaining({
      expected_workspace_revision: 'workspace-revision-1',
      expected_settings_revision: 'settings-revision-1',
    readiness_token: 'ready-token-1',
    }));
  });

  it('blocks activation before Share projection when status metadata is stale', async () => {
    vi.mocked(modbusShareAPI.status).mockResolvedValue({
      enabled: true,
      port: 5020,
      address: '0.0.0.0:5020',
      bind_state: 'pass',
      mapping_count: 0,
      hydration_state: 'stale',
      readiness: false,
      readiness_token: 'stale-token',
      workspace_revision: 'workspace-revision-1',
      settings_revision: 'settings-revision-1',
    });

    renderPage();
    await waitFor(() => expect(screen.getByTestId('share-page-activate')).toBeInTheDocument());
    fireEvent.click(screen.getByTestId('share-page-activate'));

    await waitFor(() => expect(modbusShareAPI.status).toHaveBeenCalledTimes(1));
    expect(modbusShareAPI.reconcile).not.toHaveBeenCalled();
    expect(studioV2WorkspaceActivationAPI.activate).not.toHaveBeenCalled();
  });

  it('activates the workspace without Share projection when global Share is disabled', async () => {
    vi.mocked(modbusShareAPI.status).mockResolvedValue({
      enabled: false,
      configured_enabled: false,
      port: 5020,
      address: '',
      bind_state: 'disabled',
      mapping_count: 0,
      hydration_state: 'ready',
      readiness: true,
      readiness_token: 'disabled-token',
      workspace_revision: 'workspace-revision-1',
      settings_revision: 'settings-revision-1',
    });

    renderPage();
    await waitFor(() => expect(screen.getByTestId('share-page-activate')).toBeInTheDocument());
    fireEvent.click(screen.getByTestId('share-page-activate'));

    await waitFor(() => expect(studioV2WorkspaceActivationAPI.activate).toHaveBeenCalledTimes(1));
    expect(screen.queryByTestId('share-page-activation-error')).not.toBeInTheDocument();
    expect(modbusShareAPI.reconcile).not.toHaveBeenCalled();
    expect(studioV2WorkspaceActivationAPI.activate).toHaveBeenCalledWith(expect.objectContaining({ readiness_token: 'disabled-token', pending_saves: 0 }));
  });

  it('uses disabled bootstrap truth to skip Share projection and activate the workspace', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspaceFixture({
      modbus_share: {
        hydration_state: 'ready',
        readiness: true,
        status: {
          enabled: false,
          port: 5020,
          address: '',
          bind_state: 'disabled',
          mapping_count: 0,
        },
      },
    }));

    renderPage();
    await waitFor(() => expect(screen.getByTestId('share-page-activate')).toBeInTheDocument());
    fireEvent.click(screen.getByTestId('share-page-activate'));

    await waitFor(() => expect(studioV2WorkspaceActivationAPI.activate).toHaveBeenCalledTimes(1));
    expect(screen.queryByTestId('share-page-activation-error')).not.toBeInTheDocument();
    expect(modbusShareAPI.reconcile).not.toHaveBeenCalled();
  });

  it('does not activate when the canonical Share candidate snapshot is unavailable', async () => {
    vi.mocked(modbusShareAPI.status).mockResolvedValue({
      enabled: true,
      port: 5020,
      address: '0.0.0.0:5020',
      bind_state: 'pass',
      mapping_count: 4,
      hydration_state: 'ready',
      readiness: true,
      readiness_token: 'ready-token-1',
      workspace_revision: 'workspace-revision-1',
      settings_revision: 'settings-revision-1',
    });

    renderPage();
    await waitFor(() => expect(screen.getByTestId('share-page-activate')).toBeInTheDocument());
    fireEvent.click(screen.getByTestId('share-page-activate'));

    await waitFor(() => expect(modbusShareAPI.status).toHaveBeenCalledTimes(1));
    expect(modbusShareAPI.reconcile).not.toHaveBeenCalled();
    expect(studioV2WorkspaceActivationAPI.activate).not.toHaveBeenCalled();
  });

  it('does not activate when runtime Share synchronization fails', async () => {
    vi.mocked(modbusShareAPI.reconcile).mockRejectedValue(new Error('share runtime unavailable'));

    renderPage();
    await waitFor(() => expect(screen.getByTestId('share-page-activate')).toBeInTheDocument());
    fireEvent.click(screen.getByTestId('share-page-activate'));

    await waitFor(() => expect(modbusShareAPI.reconcile).toHaveBeenCalledTimes(1));
    expect(studioV2WorkspaceActivationAPI.activate).not.toHaveBeenCalled();
    expect(modbusShareAPI.reconcile).toHaveBeenCalledTimes(1);
  });

  it('activates without runtime Share calls when globally disabled', async () => {
    const state: WorkbenchV2State = {
      ...INITIAL_STATE,
      settings: {
        ...INITIAL_STATE.settings,
        modbus_share: { ...INITIAL_STATE.settings.modbus_share, enabled: false },
      },
    };
    const activateWorkspace = vi.fn().mockResolvedValue({ status: 'active' });

    await activateStudioV2WorkspaceWithShare(state, activateWorkspace);

    expect(activateWorkspace).toHaveBeenCalledTimes(1);
    expect(modbusShareAPI.reconcile).not.toHaveBeenCalled();
  });

  it('blocks Share synchronization without a backend canonical plan', async () => {
    const state: WorkbenchV2State = {
      ...INITIAL_STATE,
      rules: [],
      points: [],
      mappings: {
        'frontend-only-point': {
          point_id: 'frontend-only-point',
          tag_key: 'frontend.only',
          display_name: 'Frontend only',
          unit: '',
          target_type: 'int16',
          scale: 1,
          offset: 0,
          enabled: false,
          tag_id: 'tag-frontend-only',
        },
      },
    };
    await expect(syncStudioV2ShareMappings(state, {
      workspace_id: 'workspace-1',
      workspace_revision: 'workspace-revision-1',
      settings_revision: 'settings-revision-1',
      readiness_token: 'ready-token-1',
    })).rejects.toMatchObject({ code: 'modbus_share_projection_required' });
    expect(modbusShareAPI.reconcile).not.toHaveBeenCalled();
  });
});
