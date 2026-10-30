import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '@/services/studioV2Rules';
import { studioV2MappingsAPI } from '@/services/studioV2Mappings';
import { pointAPI } from '@/services/datalink';
import { useStudioV2AutosaveState } from '@/pages/datalink/workbench-v2/useStudioV2AutosaveState';
vi.mock('@/pages/datalink/workbench-v2/useStudioV2AutosaveState', () => ({ useStudioV2AutosaveState: vi.fn(() => { throw new Error('Home must not mount autosave owners'); }) }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2Workspace', () => ({ studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), schemaApply: vi.fn(), remove: vi.fn(),
} }));
vi.mock('@/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn(), create: vi.fn() } }));
vi.mock('@/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn(), create: vi.fn() } }));
vi.mock('@/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn() } }));
vi.mock('@/services/datalink', () => ({ pointAPI: { list: vi.fn() } }));
const workspace = { id: 'W', kind: 'single' as const, status: 'empty' as const, ordered_device_ids: [], created_at: '', updated_at: '' };
const clients: QueryClient[] = [];
beforeEach(() => {
  vi.clearAllMocks(); vi.mocked(studioV2WorkspaceAPI.get).mockReset().mockResolvedValue(workspace);
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockReset().mockResolvedValue({ workspace_id: 'W', workspace_revision: '', groups: [] });
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockReset().mockResolvedValue([]); vi.mocked(studioV2RulesAPI.list).mockReset().mockResolvedValue([]);
  vi.mocked(studioV2MappingsAPI.list).mockReset().mockResolvedValue([]); vi.mocked(pointAPI.list).mockReset().mockResolvedValue([]);
});
afterEach(() => { cleanup(); clients.forEach(client => client.clear()); clients.length = 0; });
function mount() {
  const client = new QueryClient(); clients.push(client); const navigate = vi.fn();
  render(<QueryClientProvider client={client}><DatalinkWorkbenchV2Page presentation="recordings" locatorIntent={{ kind: 'absent' }}
    presentationLocation={{ search: '?view=recordings&section=database', hash: '#saved' }} navigateTo={navigate} /></QueryClientProvider>);
  return navigate;
}
describe('actual optional home entry', () => {
  it('renders the initial empty database setup without Share and opens no autosave/mutation owners', async () => {
    const navigate = mount(); expect(screen.getByTestId('recording-home')).toBeInTheDocument();
    await screen.findByText('reason.missing_devices'); expect(useStudioV2AutosaveState).not.toHaveBeenCalled();
    for (const action of ['create', 'update', 'apply', 'schemaApply', 'remove'] as const) expect(studioV2WorkspaceWriteGroupsAPI[action]).not.toHaveBeenCalled();
    expect(studioV2WorkspaceDevicesAPI.create).not.toHaveBeenCalled(); expect(studioV2RulesAPI.create).not.toHaveBeenCalled(); expect(studioV2MappingsAPI.create).not.toHaveBeenCalled();
    fireEvent.click(screen.getByText('advanced')); expect(navigate).toHaveBeenLastCalledWith('/studio/v2?section=database&view=advanced#saved');
  });
  it('failed scope bootstrap is safe and preserves query/hash on recovery navigation', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockRejectedValue(new Error('password=secret private.host')); const navigate = mount();
    await screen.findByRole('alert'); expect(screen.getByText('error')).toBeInTheDocument();
    expect(document.body.textContent).not.toContain('secret'); expect(screen.queryByText('reason.missing_devices')).not.toBeInTheDocument();
    fireEvent.click(screen.getByText('advanced')); expect(navigate).toHaveBeenCalledWith('/studio/v2?section=database&view=advanced#saved');
  });
  it('source read failure never becomes no points or an empty healthy workspace', async () => {
    vi.mocked(studioV2MappingsAPI.list).mockRejectedValue(new Error('private backend')); mount();
    await screen.findByText('reason.evidence_unavailable');
    await waitFor(() => expect(screen.getByText('error')).toBeInTheDocument());
    expect(screen.queryByText('reason.missing_devices')).not.toBeInTheDocument();
    expect(document.body.textContent).not.toContain('private backend');
  });
  it('missing DB revision is not silently coerced to the initial empty revision', async () => {
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W', groups: [] } as unknown as Awaited<ReturnType<typeof studioV2WorkspaceWriteGroupsAPI.list>>);
    mount(); await screen.findByText('unknown'); expect(screen.queryByText('reason.missing_devices')).not.toBeInTheDocument();
  });
  it('late sources from the previous workspace cannot overwrite the newly selected workspace', async () => {
    let resolveOld!: (value: Awaited<ReturnType<typeof pointAPI.list>>) => void;
    vi.mocked(pointAPI.list).mockImplementationOnce(() => new Promise(resolve => { resolveOld = resolve; }));
    mount(); await waitFor(() => expect(pointAPI.list).toHaveBeenCalledTimes(1));
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...workspace, id: 'NEXT' });
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'NEXT', workspace_revision: '', groups: [] });
    await act(async () => { await Promise.all([
      clients[0].invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap(), exact: true }),
      clients[0].invalidateQueries({ queryKey: studioV2WorkspaceKeys.writeGroups(), exact: true }),
    ]); });
    await screen.findByText('reason.missing_devices');
    await act(async () => resolveOld([]));
    expect(screen.getByText('reason.missing_devices')).toBeInTheDocument();
    expect(screen.queryByText('unknown')).not.toBeInTheDocument();
    expect(clients[0].getQueryCache().getAll().some(query => query.queryKey.includes('recording-home-sources') && query.queryKey.includes('NEXT'))).toBe(true);
  });

});
