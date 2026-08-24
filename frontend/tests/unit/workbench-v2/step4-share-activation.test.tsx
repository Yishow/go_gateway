import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
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

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell', () => ({
  WorkbenchV2Shell: ({ activateWorkspace }: { activateWorkspace?: () => Promise<unknown> }) => (
    <button type="button" data-testid="share-page-activate" onClick={() => void activateWorkspace?.().catch(() => undefined)}>
      Activate
    </button>
  ),
}));
vi.mock('../../../src/services/datalink', () => ({
  modbusShareAPI: { listMappings: vi.fn(), deleteMapping: vi.fn(), upsertMapping: vi.fn() },
}));
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('../../../src/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn() } }));
vi.mock('../../../src/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
vi.mock('../../../src/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
vi.mock('../../../src/services/studioV2WorkspaceDatabase', () => ({ studioV2WorkspaceDatabaseAPI: { getConfig: vi.fn(), listTargets: vi.fn() } }));
vi.mock('../../../src/services/studioV2WorkspaceAudit', () => ({ studioV2WorkspaceAuditAPI: { list: vi.fn() } }));
vi.mock('../../../src/services/studioV2WorkspaceActivation', () => ({ studioV2WorkspaceActivationAPI: { activate: vi.fn() } }));

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
  vi.mocked(modbusShareAPI.listMappings).mockResolvedValue([]);
  vi.mocked(modbusShareAPI.deleteMapping).mockResolvedValue();
  vi.mocked(modbusShareAPI.upsertMapping).mockImplementation(async (tagId, register) => ({ tag_id: tagId, register, data_type: 'int16', updated_at: '' }));
}

describe('Step 4 activation Share handoff', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupPageMocks();
  });

  it('awaits runtime Share synchronization before workspace activation', async () => {
    const events: string[] = [];
    vi.mocked(modbusShareAPI.listMappings).mockImplementation(async () => { events.push('list'); return []; });
    vi.mocked(modbusShareAPI.upsertMapping).mockImplementation(async (tagId, register) => {
      events.push(`upsert:${tagId}:${register}`);
      return { tag_id: tagId, register, data_type: 'int16', updated_at: '' };
    });
    vi.mocked(studioV2WorkspaceActivationAPI.activate).mockImplementation(async () => {
      events.push('activate');
      return { workspace_id: 'workspace-1', results: [] };
    });

    renderPage();
    await waitFor(() => expect(screen.getByTestId('share-page-activate')).toBeInTheDocument());
    fireEvent.click(screen.getByTestId('share-page-activate'));

    await waitFor(() => expect(studioV2WorkspaceActivationAPI.activate).toHaveBeenCalledTimes(1));
    expect(events[0]).toBe('list');
    expect(events.at(-1)).toBe('activate');
    expect(modbusShareAPI.upsertMapping).toHaveBeenCalledWith('tag-d0', 0);
  });

  it('does not activate when a persisted Share mapping is missing its tag id', async () => {
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([mappingFixture({
      id: 'mapping-0', point_id: 'persisted-point-0', rule_id: 'mc-rule', device_id: 'dev-mc', address: 'D0', tag_id: 'tag-d0',
    })]);

    renderPage();
    await waitFor(() => expect(screen.getByTestId('share-page-activate')).toBeInTheDocument());
    fireEvent.click(screen.getByTestId('share-page-activate'));

    await waitFor(() => expect(studioV2WorkspaceActivationAPI.activate).not.toHaveBeenCalled());
    expect(modbusShareAPI.listMappings).not.toHaveBeenCalled();
  });

  it('does not activate when runtime Share synchronization fails', async () => {
    vi.mocked(modbusShareAPI.listMappings).mockRejectedValue(new Error('share runtime unavailable'));

    renderPage();
    await waitFor(() => expect(screen.getByTestId('share-page-activate')).toBeInTheDocument());
    fireEvent.click(screen.getByTestId('share-page-activate'));

    await waitFor(() => expect(modbusShareAPI.listMappings).toHaveBeenCalledTimes(1));
    expect(studioV2WorkspaceActivationAPI.activate).not.toHaveBeenCalled();
    expect(modbusShareAPI.upsertMapping).not.toHaveBeenCalled();
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
    expect(modbusShareAPI.listMappings).not.toHaveBeenCalled();
    expect(modbusShareAPI.deleteMapping).not.toHaveBeenCalled();
    expect(modbusShareAPI.upsertMapping).not.toHaveBeenCalled();
  });

  it('preserves runtime mappings without durable source-rule ownership proof', async () => {
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
    vi.mocked(modbusShareAPI.listMappings).mockResolvedValue([
      { tag_id: 'tag-frontend-only', register: 99, data_type: 'int16', updated_at: '' },
    ]);

    await syncStudioV2ShareMappings(state);

    expect(modbusShareAPI.deleteMapping).not.toHaveBeenCalledWith('tag-frontend-only');
  });
});
