import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { RecordingHomeRoute } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeRoute';
import { studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { recordingHomeSourcesAPI } from '@/services/recordingHomeSources';
import { runtimeAPI } from '@/services/datalink';
import { group } from './recording-home-guidance-fixtures';
import { useRecordingHomeDeviceStatus } from '@/features/datalink/workbench-v2/recording-home/useRecordingHomeDeviceStatus';
import { deriveDeviceConnectionState } from '@/features/datalink/workbench-v2/recording-home/deviceConnectionState';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2Workspace', async original => ({ ...await original<typeof import('@/services/studioV2Workspace')>(), studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn() } }));
vi.mock('@/services/recordingHomeSources', () => ({ recordingHomeSourcesAPI: { read: vi.fn() } }));
vi.mock('@/hooks/datalink/useAcquisitionAnomalies', () => ({ useAcquisitionAnomalySummaryQuery: () => ({ isSuccess: false }) }));

let client: QueryClient;
afterEach(() => { cleanup(); client?.clear(); vi.restoreAllMocks(); vi.useRealTimers(); });

it('reads one workspace-wide runtime status on entering the visible home', async () => {
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ id: 'W1', kind: 'single', status: 'empty', ordered_device_ids: [], created_at: '', updated_at: '' });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W1', workspace_revision: 'W2', groups: [] });
  vi.mocked(recordingHomeSourcesAPI.read).mockResolvedValue({ workspace_id: 'W1', state: 'available', source_digest: 'a'.repeat(64), devices: [], candidates: [], excluded: [] });
  const getStatus = vi.spyOn(runtimeAPI, 'getStatus').mockResolvedValue({ running: false, uptime_seconds: 0, collectors: [] });
  client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(<QueryClientProvider client={client}><RecordingHomeRoute location={{ search: '?view=recordings', hash: '' }} /></QueryClientProvider>);
  await waitFor(() => expect(screen.getByTestId('home-source-observation')).toHaveAttribute('data-source-status', 'fresh'));
  expect(getStatus).toHaveBeenCalledExactlyOnceWith();
});

function StatusProbe({ workspaceId = 'W1' }: { workspaceId?: string }) {
  const result = useRecordingHomeDeviceStatus(workspaceId);
  return <><span data-testid="probe-state">{deriveDeviceConnectionState(result.snapshot?.collectors[0], result.now)}</span>
    <span data-testid="probe-outdated">{String(result.outdated)}</span></>;
}
function mountProbe() {
  client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(<QueryClientProvider client={client}><StatusProbe /></QueryClientProvider>);
}
const successfulSnapshot = () => ({ running: true, uptime_seconds: 10, collectors: [{ device_id: 'MC', device_name: 'MC', protocol: 'mc',
  status: 'running' as const, availability_status: 'available' as const, running: true, points_total: 48, points_healthy: 48,
  points_stale: 0, points_error: 0, last_read_at: new Date(Date.now() - 1000).toISOString(), last_error: null, breaker_state: 'closed', effective_interval_ms: 1000 }],
  snapshot_state: { state: 'ready', empty: false, degraded: false, unavailable: false, stale: false } });

it('keeps a failed read as outdated for at most fifteen seconds, then marks it unavailable', async () => {
  vi.useFakeTimers(); vi.setSystemTime(new Date('2026-11-07T12:00:10Z'));
  const getStatus = vi.spyOn(runtimeAPI, 'getStatus').mockResolvedValueOnce(successfulSnapshot()).mockRejectedValue(new Error('secret backend error'));
  mountProbe();
  await act(async () => { await vi.advanceTimersByTimeAsync(20); });
  expect(screen.getByTestId('probe-state')).toHaveTextContent('collecting');
  await act(async () => { await vi.advanceTimersByTimeAsync(5000); });
  expect(getStatus).toHaveBeenCalledTimes(2);
  expect(screen.getByTestId('probe-state')).toHaveTextContent('collecting');
  expect(screen.getByTestId('probe-outdated')).toHaveTextContent('true');
  await act(async () => { await vi.advanceTimersByTimeAsync(10000); });
  expect(screen.getByTestId('probe-state')).toHaveTextContent('unavailable');
  expect(screen.queryByText('secret backend error')).not.toBeInTheDocument();
});

it('retains the last observation across a revision remount within the same workspace', async () => {
  vi.useFakeTimers(); vi.setSystemTime(new Date('2026-11-07T12:00:10Z'));
  vi.spyOn(runtimeAPI, 'getStatus').mockResolvedValueOnce(successfulSnapshot()).mockRejectedValue(new Error('offline'));
  client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const view = render(<QueryClientProvider client={client}><StatusProbe key="first" /></QueryClientProvider>);
  await act(async () => { await vi.advanceTimersByTimeAsync(20); });
  expect(screen.getByTestId('probe-state')).toHaveTextContent('collecting');
  await act(async () => { await vi.advanceTimersByTimeAsync(1000); });
  view.rerender(<QueryClientProvider client={client}><StatusProbe key="second" /></QueryClientProvider>);
  expect(screen.getByTestId('probe-state')).toHaveTextContent('collecting');
  await act(async () => { await vi.advanceTimersByTimeAsync(20); });
  expect(screen.getByTestId('probe-outdated')).toHaveTextContent('true');
  view.rerender(<QueryClientProvider client={client}><StatusProbe key="other" workspaceId="W2" /></QueryClientProvider>);
  expect(screen.getByTestId('probe-state')).toHaveTextContent('unavailable');
  view.rerender(<QueryClientProvider client={client}><StatusProbe key="return" /></QueryClientProvider>);
  expect(screen.getByTestId('probe-state')).toHaveTextContent('collecting');
  await act(async () => { await vi.advanceTimersByTimeAsync(14000); });
  expect(screen.getByTestId('probe-state')).toHaveTextContent('unavailable');
});

it('stops polling while hidden and reads once immediately when visible again', async () => {
  vi.useFakeTimers(); vi.setSystemTime(new Date('2026-11-07T12:00:10Z'));
  let visible = true;
  vi.spyOn(document, 'visibilityState', 'get').mockImplementation(() => visible ? 'visible' : 'hidden');
  const getStatus = vi.spyOn(runtimeAPI, 'getStatus').mockImplementation(async () => successfulSnapshot());
  const view = mountProbe();
  await act(async () => { await vi.advanceTimersByTimeAsync(20); });
  expect(getStatus).toHaveBeenCalledTimes(1);
  await act(async () => { visible = false; document.dispatchEvent(new Event('visibilitychange')); await vi.advanceTimersByTimeAsync(20000); });
  expect(getStatus).toHaveBeenCalledTimes(1);
  await act(async () => { visible = true; document.dispatchEvent(new Event('visibilitychange')); await vi.advanceTimersByTimeAsync(20); });
  expect(getStatus).toHaveBeenCalledTimes(2);
  view.unmount();
  await act(async () => { await vi.advanceTimersByTimeAsync(10000); });
  expect(getStatus).toHaveBeenCalledTimes(2);
});

it('never treats an unavailable snapshot as a successful collector', async () => {
  vi.useFakeTimers(); vi.setSystemTime(new Date('2026-11-07T12:00:10Z'));
  const unavailableSnapshot = { ...successfulSnapshot(),
    snapshot_state: { state: 'unavailable', empty: false, degraded: false, unavailable: true, stale: false },
  };
  vi.spyOn(runtimeAPI, 'getStatus').mockResolvedValue(unavailableSnapshot);
  mountProbe();
  await act(async () => { await vi.advanceTimersByTimeAsync(20); });
  expect(screen.getByTestId('probe-state')).toHaveTextContent('unavailable');
});

it('keeps states while a server snapshot declares itself stale', async () => {
  vi.useFakeTimers(); vi.setSystemTime(new Date('2026-11-07T12:00:10Z'));
  const staleSnapshot = { ...successfulSnapshot(), snapshot_state: { state: 'stale', empty: false, degraded: false, unavailable: false, stale: true } };
  vi.spyOn(runtimeAPI, 'getStatus').mockResolvedValue(staleSnapshot);
  mountProbe();
  await act(async () => { await vi.advanceTimersByTimeAsync(20); });
  expect(screen.getByTestId('probe-state')).toHaveTextContent('collecting');
  expect(screen.getByTestId('probe-outdated')).toHaveTextContent('true');
});

it('summarizes all devices, recordings and runtime writes, and toggles the device filter', async () => {
  const devices = ['MC', 'mc_new', 'M1'].map(id => ({ id, name: id, protocol: 'mc', availability: 'available' }));
  const groups = [group('one', 'draft'), group('two', 'draft'), group('three', 'draft'), group('four', 'deleted')]
    .map((value, index) => ({ ...value, members: value.members.map(member => ({ ...member, device_id: index < 2 ? 'MC' : 'M1' })) }));
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ id: 'W1', kind: 'single', status: 'ready', ordered_device_ids: devices.map(d => d.id), created_at: '', updated_at: '' });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W1', workspace_revision: 'W2', groups });
  vi.mocked(recordingHomeSourcesAPI.read).mockResolvedValue({ workspace_id: 'W1', state: 'available', source_digest: 'a'.repeat(64), devices, candidates: [], excluded: [] });
  const readAt = new Date(Date.now() - 1000).toISOString();
  const collector = { device_id: 'MC', device_name: 'MC', protocol: 'mc', status: 'warning' as const, availability_status: 'available' as const,
    running: true, points_total: 48, points_healthy: 46, points_stale: 2, points_error: 0, last_read_at: readAt, last_error: null, breaker_state: 'closed' };
  vi.spyOn(runtimeAPI, 'getStatus').mockResolvedValue({ running: true, uptime_seconds: 10, metrics: { write_success_total: 10, write_error_total: 2 },
    collectors: [collector, ...devices.slice(1).map(device => ({ ...collector, device_id: device.id, device_name: device.name, status: 'idle' as const, running: false }))],
    snapshot_state: { state: 'ready', empty: false, degraded: false, unavailable: false, stale: false },
  } as Awaited<ReturnType<typeof runtimeAPI.getStatus>>);
  client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(<QueryClientProvider client={client}><RecordingHomeRoute location={{ search: '?view=recordings', hash: '' }} /></QueryClientProvider>);
  const summary = await screen.findByTestId('home-status-summary');
  await waitFor(() => expect(within(summary).getByTestId('device-count-degraded')).toHaveTextContent('1'));
  expect(within(summary).getByTestId('device-count-total')).toHaveTextContent('3');
  expect(within(summary).getByTestId('device-count-idle')).toHaveTextContent('2');
  expect(within(summary).getByTestId('record-count-total')).toHaveTextContent('4');
  expect(within(summary).getByTestId('record-count-draft')).toHaveTextContent('3');
  expect(within(summary).getByTestId('record-count-deleted')).toHaveTextContent('1');
  expect(within(summary).getByTestId('runtime-write-counts')).toHaveTextContent('10');
  expect(within(summary).getByTestId('runtime-write-counts')).toHaveTextContent('2');
  fireEvent.click(screen.getByTestId('device-status-MC'));
  expect(screen.getByLabelText('deviceFilter')).toHaveValue('MC');
  expect(screen.getByText('Recording one')).toBeInTheDocument();
  expect(screen.queryByText('Recording three')).not.toBeInTheDocument();
  expect(within(summary).getByTestId('record-count-total')).toHaveTextContent('4');
  fireEvent.click(screen.getByTestId('device-status-MC'));
  expect(screen.getByLabelText('deviceFilter')).toHaveValue('');
  expect(screen.getByText('Recording three')).toBeInTheDocument();
  expect(screen.queryByTestId('home-device-summary')).not.toBeInTheDocument();
  fireEvent.click(screen.getByRole('button', { name: 'devicesTab' }));
  const table = screen.getByTestId('home-management-table');
  expect(within(table).getByText('MC').closest('tr')?.querySelector('[data-connection-state]')).toHaveAttribute('data-connection-state', 'degraded');
  expect(within(table).getByText('mc_new').closest('tr')?.querySelector('[data-connection-state]')).toHaveAttribute('data-connection-state', 'idle');
});
