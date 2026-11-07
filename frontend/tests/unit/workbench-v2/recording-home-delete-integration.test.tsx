import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { RecordingHomeRoute } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeRoute';
import { studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { recordingHomeSourcesAPI } from '@/services/recordingHomeSources';
import { runtimeAPI } from '@/services/datalink';
import { DeleteResponseError } from '@/services/studioV2WorkspaceDelete';
import { savedLifecycle } from '../../fixtures/writeGroupLifecycle';
import { delivery, group } from './recording-home-guidance-fixtures';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2Workspace', async original => ({ ...await original<typeof import('@/services/studioV2Workspace')>(), studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('@/services/recordingHomeSources', () => ({ recordingHomeSourcesAPI: { read: vi.fn() } }));
vi.mock('@/hooks/datalink/useAcquisitionAnomalies', () => ({ useAcquisitionAnomalySummaryQuery: () => ({ isSuccess: false }) }));
afterEach(() => { cleanup(); vi.restoreAllMocks(); });

it('deletes from a home row and updates both its list and summary without page reload', async () => {
  let groups = [group('G1', 'ready'), group('G2', 'ready')];
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ id: 'W1', kind: 'single', status: 'ready', ordered_device_ids: ['D1'], created_at: '', updated_at: '' });
  vi.mocked(recordingHomeSourcesAPI.read).mockResolvedValue({ workspace_id: 'W1', state: 'available', source_digest: 'a'.repeat(64),
    devices: [{ id: 'D1', name: 'mc_new', protocol: 'mc', availability: 'available' }], candidates: [], excluded: [] });
  vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'list').mockImplementation(async () => ({ workspace_id: 'W1', workspace_revision: 'W2', groups }));
  vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'delivery').mockImplementation(async id => delivery(groups.find(item => item.id === id)!));
  vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockImplementation(async id => savedLifecycle(groups.find(item => item.id === id)!, 'W2'));
  vi.spyOn(runtimeAPI, 'getStatus').mockResolvedValue({ running: false, uptime_seconds: 0, collectors: [] });
  const remove = vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'remove').mockImplementation(async id => {
    groups = groups.filter(item => item.id !== id);
    return { deleted: true, group_id: id, lineage_id: 'L1', deleted_generations: 1 };
  });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><RecordingHomeRoute location={{ search: '?view=recordings', hash: '' }} /></QueryClientProvider>);
  const row = (await screen.findByText('Recording G1')).closest('tr')!;
  const item = await within(row).findByRole('button', { name: 'rowEvidence.deleteNamed' });
  await waitFor(() => expect(item).toBeEnabled()); fireEvent.click(item);
  fireEvent.click(await screen.findByRole('checkbox', { name: 'deleteDialogs.recording.acknowledge' }));
  fireEvent.click(screen.getByTestId('group-delete-confirm-button'));
  await waitFor(() => expect(screen.queryByText('Recording G1')).not.toBeInTheDocument());
  expect(within(screen.getByTestId('home-management-table')).getByText('Recording G2')).toBeInTheDocument();
  expect(screen.getByTestId('record-count-total')).toHaveTextContent('1');
  expect(await screen.findByText('deleteDialogs.recording.success')).toBeInTheDocument();
  expect(remove).toHaveBeenCalledOnce(); client.clear();
});

it('opens device deletion from the device card in one click, blocking referenced devices', async () => {
  const groups = [group('G1', 'ready'), group('G2', 'ready')];
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ id: 'W1', kind: 'single', status: 'ready', ordered_device_ids: ['D1'], created_at: '', updated_at: '' });
  vi.mocked(recordingHomeSourcesAPI.read).mockResolvedValue({ workspace_id: 'W1', state: 'available', source_digest: 'a'.repeat(64),
    devices: [{ id: 'D1', name: 'mc_new', protocol: 'mc', availability: 'available' }], candidates: [], excluded: [] });
  vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'list').mockResolvedValue({ workspace_id: 'W1', workspace_revision: 'W2', groups });
  vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'delivery').mockImplementation(async id => delivery(groups.find(item => item.id === id)!));
  vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockImplementation(async id => savedLifecycle(groups.find(item => item.id === id)!, 'W2'));
  vi.spyOn(runtimeAPI, 'getStatus').mockResolvedValue({ running: false, uptime_seconds: 0, collectors: [] });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><RecordingHomeRoute location={{ search: '?view=recordings', hash: '' }} /></QueryClientProvider>);
  await screen.findByTestId('device-status-D1');
  await waitFor(() => expect(screen.getByTestId('home-source-observation')).toHaveAttribute('data-source-status', 'fresh'));
  const strip = screen.getByRole('region', { name: 'deviceStatus.strip' });
  fireEvent.click(within(strip).getByRole('button', { name: 'deleteDialogs.device.deleteNamed' }));
  expect(await screen.findByText('Recording G1', { selector: 'li span' })).toBeInTheDocument();
  expect(screen.getByText('Recording G2', { selector: 'li span' })).toBeInTheDocument();
  expect(screen.getByTestId('device-delete-confirm-button')).toBeDisabled();
  fireEvent.click(screen.getByRole('button', { name: 'deleteDialogs.device.cancel' }));
  await waitFor(() => expect(within(strip).getByRole('button', { name: 'deleteDialogs.device.deleteNamed' })).toHaveFocus());
  client.clear();
});

it('keeps an unknown result open and reports a failed refresh inside the dialog', async () => {
  const groups = [group('G1', 'ready')];
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ id: 'W1', kind: 'single', status: 'ready', ordered_device_ids: ['D1'], created_at: '', updated_at: '' });
  vi.mocked(recordingHomeSourcesAPI.read).mockResolvedValue({ workspace_id: 'W1', state: 'available', source_digest: 'a'.repeat(64),
    devices: [{ id: 'D1', name: 'mc_new', protocol: 'mc', availability: 'available' }], candidates: [], excluded: [] });
  vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'list').mockResolvedValue({ workspace_id: 'W1', workspace_revision: 'W2', groups });
  vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'delivery').mockResolvedValue(delivery(groups[0]));
  vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockResolvedValue(savedLifecycle(groups[0], 'W2'));
  vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'remove').mockRejectedValue(new Error('private network text'));
  vi.spyOn(runtimeAPI, 'getStatus').mockResolvedValue({ running: false, uptime_seconds: 0, collectors: [] });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><RecordingHomeRoute location={{ search: '?view=recordings', hash: '' }} /></QueryClientProvider>);
  const row = (await screen.findByText('Recording G1')).closest('tr')!;
  const item = await within(row).findByRole('button', { name: 'rowEvidence.deleteNamed' });
  await waitFor(() => expect(item).toBeEnabled()); fireEvent.click(item);
  fireEvent.click(await screen.findByRole('checkbox', { name: 'deleteDialogs.recording.acknowledge' }));
  fireEvent.click(screen.getByTestId('group-delete-confirm-button'));
  expect(await screen.findByText('deleteDialogs.recording.unknown')).toBeInTheDocument();
  vi.mocked(studioV2WorkspaceAPI.get).mockRejectedValueOnce(new Error('private refresh failure'));
  fireEvent.click(screen.getByRole('button', { name: 'deleteDialogs.recording.refresh' }));
  expect(await within(screen.getByRole('alertdialog')).findByText('deleteDialogs.refreshFailed')).toBeInTheDocument();
  expect(screen.queryByText('private refresh failure')).not.toBeInTheDocument();
  expect(screen.getByTestId('group-delete-confirm-button')).toBeDisabled(); client.clear();
});

it('keeps a 409 open and leads the operator to an explicit stop confirmation', async () => {
  const groups = [group('G1', 'running')];
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ id: 'W1', kind: 'single', status: 'ready', ordered_device_ids: ['D1'], created_at: '', updated_at: '' });
  vi.mocked(recordingHomeSourcesAPI.read).mockResolvedValue({ workspace_id: 'W1', state: 'available', source_digest: 'a'.repeat(64),
    devices: [{ id: 'D1', name: 'mc_new', protocol: 'mc', availability: 'available' }], candidates: [], excluded: [] });
  vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'list').mockResolvedValue({ workspace_id: 'W1', workspace_revision: 'W2', groups });
  vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'delivery').mockResolvedValue(delivery(groups[0]));
  vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockResolvedValue(savedLifecycle(groups[0], 'W2'));
  vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'remove').mockRejectedValue(new DeleteResponseError('failed', 'WRITE_GROUP_DELETE_ACTIVE', 'stop_first'));
  const disable = vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'disable');
  vi.spyOn(runtimeAPI, 'getStatus').mockResolvedValue({ running: false, uptime_seconds: 0, collectors: [] });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><RecordingHomeRoute location={{ search: '?view=recordings', hash: '' }} /></QueryClientProvider>);
  const row = (await screen.findByText('Recording G1')).closest('tr')!;
  const item = await within(row).findByRole('button', { name: 'rowEvidence.deleteNamed' });
  await waitFor(() => expect(item).toBeEnabled()); fireEvent.click(item);
  fireEvent.click(await screen.findByRole('checkbox', { name: 'deleteDialogs.recording.acknowledge' }));
  fireEvent.click(screen.getByTestId('group-delete-confirm-button'));
  expect(await screen.findByText('deleteDialogs.recording.active')).toBeInTheDocument();
  expect(screen.getByText('Recording G1', { selector: 'td strong' })).toBeInTheDocument();
  fireEvent.click(screen.getByRole('button', { name: 'deleteDialogs.recording.stop' }));
  expect(await screen.findByTestId('group-stop-confirm')).toBeInTheDocument();
  expect(disable).not.toHaveBeenCalled(); client.clear();
});

it('updates the device strip and summary after deleting an unreferenced device', async () => {
  let devices = [{ id: 'D1', name: 'mc_new', protocol: 'mc', availability: 'available' }];
  vi.mocked(studioV2WorkspaceAPI.get).mockImplementation(async () => ({ id: 'W1', kind: 'single', status: 'ready', ordered_device_ids: devices.map(item => item.id), created_at: '', updated_at: '' }));
  vi.mocked(recordingHomeSourcesAPI.read).mockImplementation(async () => ({ workspace_id: 'W1', state: 'available', source_digest: 'a'.repeat(64), devices, candidates: [], excluded: [] }));
  vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'list').mockResolvedValue({ workspace_id: 'W1', workspace_revision: 'W2', groups: [] });
  vi.spyOn(runtimeAPI, 'getStatus').mockResolvedValue({ running: false, uptime_seconds: 0, collectors: [] });
  const remove = vi.spyOn(studioV2WorkspaceDevicesAPI, 'remove').mockImplementation(async () => { devices = []; });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><RecordingHomeRoute location={{ search: '?view=recordings', hash: '' }} /></QueryClientProvider>);
  await screen.findByTestId('device-status-D1');
  await waitFor(() => expect(screen.getByTestId('home-source-observation')).toHaveAttribute('data-source-status', 'fresh'));
  const strip = screen.getByRole('region', { name: 'deviceStatus.strip' });
  fireEvent.click(within(strip).getByRole('button', { name: 'deleteDialogs.device.deleteNamed' }));
  fireEvent.change(await screen.findByLabelText('deleteDialogs.device.typeName'), { target: { value: 'mc_new' } });
  fireEvent.click(screen.getByTestId('device-delete-confirm-button'));
  await waitFor(() => expect(screen.queryByTestId('device-status-D1')).not.toBeInTheDocument());
  expect(screen.getByTestId('device-count-total')).toHaveTextContent('0');
  expect(await screen.findByText('deleteDialogs.device.success')).toBeInTheDocument();
  expect(remove).toHaveBeenCalledExactlyOnceWith('D1'); client.clear();
});
