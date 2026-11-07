import { recordingHomeSourcesAPI } from '@/services/recordingHomeSources';
import { homeSources } from './recording-home-source-fixtures';
import { homeDeviceAction } from './recording-home-device-actions';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeAll, afterEach, beforeEach, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2DatalinkApi, studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '@/services/studioV2Rules';
import { studioV2MappingsAPI } from '@/services/studioV2Mappings';
import { deviceAPI, modbusShareAPI, pointAPI } from '@/services/datalink';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { sourceRuleCandidateAPI } from '@/services/sourceRuleCandidates';
import { locatorRead, nativeDevice } from './studio-locator-fixtures';
import i18n from '@/i18n/config';
vi.mock('@/services/sourceRuleCandidates', () => ({ sourceRuleCandidateAPI: { get: vi.fn(), getTags: vi.fn() } }));
vi.mock('@/services/recordingHomeSources', () => ({ recordingHomeSourcesAPI: { read: vi.fn() } }));
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { defaults: { baseURL: '/api/v1/datalink' }, get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() }, studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
vi.mock('@/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/datalink', () => ({ modbusShareAPI: { status: vi.fn() }, pointAPI: { list: vi.fn(), get: vi.fn() }, tagAPI: { get: vi.fn() }, deviceAPI: { checkReadiness: vi.fn() } }));
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: vi.fn() } }));

// Zero saved MC ranges: free cells are claimable only when the list read is bound to a valid workspace revision.
const clients: QueryClient[] = [];
const device = { ...nativeDevice(), protocol: 'mc_3e' as const, connection_config: '{"host":"fixture.invalid","port":5000,"network_no":0,"pc_no":255,"io_no":1023,"station_no":0}' };
const share = { workspace_id: 'W', workspace_revision: 'W2', hydration_state: 'ready' as const, readiness: true };
const bound = { ...locatorRead().workspace, ordered_device_ids: ['D'], modbus_share: share };
const unbound = { ...locatorRead().workspace, ordered_device_ids: ['D'] };
beforeEach(async () => {
  vi.resetAllMocks(); await i18n.changeLanguage('en'); sessionStorage.clear();
  vi.mocked(modbusShareAPI.status).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', configured_enabled: false, enabled: false, port: 502, address: '127.0.0.1:502', bind_state: 'disabled', mapping_count: 0 });
  vi.mocked(deviceAPI.checkReadiness).mockResolvedValue({ device_id: 'D', status: 'ready', connect_status: 'success', probe_status: 'success', activation_allowed: false, planning_allowed: true, apply_allowed: false, blocking_reasons: [], checks: [] });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', groups: [] });
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([device]);
  vi.mocked(recordingHomeSourcesAPI.read).mockResolvedValue(homeSources([device]));
  vi.mocked(studioV2RulesAPI.list).mockResolvedValue([]);
  vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([]); vi.mocked(pointAPI.list).mockResolvedValue([]);
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockResolvedValue({ ok: false, reason: 'unavailable', action: 'refresh' });
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async () => {
    const read = locatorRead(1); read.entry.devices = [device];
    return { data: { success: true, data: { ...unbound, locator_entry: read.entry } } };
  });
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });

async function openEmptyMemory(workspace: typeof bound | typeof unbound) {
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspace);
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 }, mutations: { retry: false } } }); clients.push(client);
  render(<QueryClientProvider client={client}><DatalinkWorkbenchV2Page presentation="recordings" locatorIntent={{ kind: 'absent' }}
    presentationLocation={{ search: '?view=recordings', hash: '#home' }} navigateTo={vi.fn()} /></QueryClientProvider>);
  fireEvent.click(await homeDeviceAction());
  const dialog = await screen.findByRole('dialog');
  // A device without ranges opens its new-range setting; the memory view is entered from the title bar.
  await within(dialog).findByTestId('home-range-option-new');
  fireEvent.click(await within(dialog).findByTestId('home-range-memory'));
  await within(dialog).findAllByRole('tab');
  return dialog;
}
function noWrites() {
  for (const fn of [studioV2RulesAPI.create, studioV2RulesAPI.update, studioV2RulesAPI.remove, studioV2MappingsAPI.create, studioV2MappingsAPI.update]) expect(fn).not.toHaveBeenCalled();
}

it('marks every area unconfirmed for zero MC ranges read without a workspace revision and offers a re-read', async () => {
  const dialog = await openEmptyMemory(unbound);
  const d0 = within(dialog).getByTestId('home-memory-cell-D0');
  expect(d0).toHaveAttribute('data-state', 'unconfirmed'); expect(d0.tagName).not.toBe('BUTTON');
  expect(within(dialog).getByRole('alert')).toHaveTextContent(/no valid workspace revision/);
  fireEvent.click(within(dialog).getByRole('tab', { name: /^M/ }));
  expect(within(dialog).getByTestId('home-memory-cell-M0')).toHaveAttribute('data-state', 'unconfirmed');
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(bound);
  fireEvent.click(within(dialog).getByRole('button', { name: 'Read ranges again' }));
  await waitFor(() => expect(within(dialog).getByTestId('home-memory-cell-M0')).toHaveAttribute('data-state', 'empty'));
  expect(within(dialog).queryByRole('alert')).not.toBeInTheDocument();
  expect(sourceRuleCandidateAPI.getTags).not.toHaveBeenCalled(); noWrites();
});

it('keeps zero MC ranges with a valid revision free and prefills D0', async () => {
  const dialog = await openEmptyMemory(bound);
  expect(within(dialog).queryByRole('alert')).not.toBeInTheDocument();
  fireEvent.click(within(dialog).getByTestId('home-memory-cell-D0'));
  expect(await within(dialog).findByRole('textbox', { name: 'Starting address' })).toHaveValue('D0');
  noWrites();
});

it('shows a failed range read with a re-read action and no cells', async () => {
  const dialog = await openEmptyMemory(bound);
  fireEvent.click(within(dialog).getByTestId('home-memory-back'));
  await waitFor(() => expect(within(dialog).getByTestId('home-range-option-new')).toHaveAttribute('aria-pressed', 'true'));
  vi.mocked(studioV2RulesAPI.list).mockRejectedValueOnce(new Error('offline'));
  fireEvent.click(await within(dialog).findByTestId('home-range-memory'));
  const alert = await within(dialog).findByRole('alert');
  expect(within(dialog).queryByTestId(/home-memory-cell-/)).not.toBeInTheDocument();
  fireEvent.click(within(alert).getByRole('button', { name: 'Read ranges again' }));
  await waitFor(() => expect(within(dialog).getByTestId('home-memory-cell-D0')).toHaveAttribute('data-state', 'empty'));
  noWrites();
});

// Warm only the code: these cases still exercise the real wrapper, scoped reads, and draft owners.
// Cold import timing and failures remain covered by the dedicated advanced-loading tests.
beforeAll(async () => { await import('@/pages/datalink/workbench-v2/AdvancedWorkbenchV2Page'); });
