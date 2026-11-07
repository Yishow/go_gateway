import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { RecordingHomeCollection } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeCollection';
import { initialCollectionState } from '@/features/datalink/workbench-v2/recording-home/recordingHomeCollections';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { fixture } from './recording-home-guidance-fixtures';
import { savedGroup } from '../../fixtures/writeGroupState';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { DeleteResponseError } from '@/services/studioV2WorkspaceDelete';
import { DeviceDeleteDialog } from '@/features/datalink/workbench-v2/recording-home/DeviceDeleteDialog';
import { RecordingHomeDeviceActions } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeDeviceActions';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
afterEach(cleanup);

function deviceDialog(recordings = [savedGroup()], onDeleted = vi.fn()) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><DeviceDeleteDialog open device={{ id: 'dev-1', name: 'mc_new' }} workspaceId="ws-1"
    recordings={recordings} onDeleted={onDeleted} onClose={vi.fn()} onRefresh={vi.fn()} /></QueryClientProvider>);
  return onDeleted;
}

it('lists loaded recordings that reference the device and blocks deletion', () => {
  const first = savedGroup({ id: 'G1', name: 'First' });
  const second = savedGroup({ id: 'G2', name: 'Second' });
  deviceDialog([first, second]);
  expect(screen.getByText('First')).toBeInTheDocument();
  expect(screen.getByText('Second')).toBeInTheDocument();
  expect(screen.getByTestId('device-delete-confirm-button')).toBeDisabled();
});

it('requires the exact device name and calls back after confirmed deletion', async () => {
  const remove = vi.spyOn(studioV2WorkspaceDevicesAPI, 'remove').mockResolvedValueOnce();
  const onDeleted = deviceDialog([]);
  fireEvent.change(screen.getByLabelText('deleteDialogs.device.typeName'), { target: { value: 'mc-new' } });
  expect(screen.getByTestId('device-delete-confirm-button')).toBeDisabled();
  fireEvent.change(screen.getByLabelText('deleteDialogs.device.typeName'), { target: { value: 'mc_new' } });
  fireEvent.click(screen.getByTestId('device-delete-confirm-button'));
  await waitFor(() => expect(onDeleted).toHaveBeenCalledOnce());
  expect(remove).toHaveBeenCalledExactlyOnceWith('dev-1');
  remove.mockRestore();
});

it('keeps the dialog open and replaces local references with a server 409 list', async () => {
  const remove = vi.spyOn(studioV2WorkspaceDevicesAPI, 'remove').mockRejectedValueOnce(new DeleteResponseError('failed', 'DEVICE_REFERENCED_BY_RECORDINGS', 'remove_recordings_first', {
    recordings: [{ id: 'G1', name: 'Server recording' }], total: 1,
  }));
  const onDeleted = deviceDialog([]);
  fireEvent.change(screen.getByLabelText('deleteDialogs.device.typeName'), { target: { value: 'mc_new' } });
  fireEvent.click(screen.getByTestId('device-delete-confirm-button'));
  expect(await screen.findByText('Server recording')).toBeInTheDocument();
  expect(screen.getByTestId('device-delete-confirm-button')).toBeDisabled();
  expect(onDeleted).not.toHaveBeenCalled();
  remove.mockRestore();
});

it('blocks confirmation when the workspace revision changed while the dialog was open', () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><DeviceDeleteDialog open device={{ id: 'dev-1', name: 'mc_new' }} workspaceId="ws-1" scopeValid={false}
    recordings={[]} onDeleted={vi.fn()} onClose={vi.fn()} onRefresh={vi.fn()} /></QueryClientProvider>);
  fireEvent.change(screen.getByLabelText('deleteDialogs.device.typeName'), { target: { value: 'mc_new' } });
  expect(screen.getByTestId('device-delete-confirm-button')).toBeDisabled();
  expect(screen.getByText('deleteDialogs.device.stale')).toBeInTheDocument(); client.clear();
});

it('explains when a tombstoned recording prevents device deletion because delivery is pending', async () => {
  const remove = vi.spyOn(studioV2WorkspaceDevicesAPI, 'remove').mockRejectedValueOnce(new DeleteResponseError('failed', 'WRITE_GROUP_DELETE_BACKLOG', 'stop_and_wait', {
    samples_unconsumed: 3, outbox: { retrying: 1, unknown: 1 },
  }));
  deviceDialog([]);
  fireEvent.change(screen.getByLabelText('deleteDialogs.device.typeName'), { target: { value: 'mc_new' } });
  fireEvent.click(screen.getByTestId('device-delete-confirm-button'));
  expect(await screen.findByText('deleteDialogs.device.blockedByRecording')).toBeInTheDocument();
  expect(screen.getByTestId('delete-backlog-samples')).toHaveTextContent('3');
  expect(screen.getByTestId('delete-backlog-retrying')).toHaveTextContent('1');
  expect(screen.getByTestId('delete-backlog-unknown')).toHaveTextContent('1');
  expect(screen.queryByText('deleteDialogs.device.unknown')).not.toBeInTheDocument(); remove.mockRestore();
});

it.each([
  ['DEVICE_DELETED_SHARE_SCOPE_REFRESH_FAILED', 'shareRefreshFailed'],
  ['DEVICE_DELETED_CONNECTION_CLOSE_FAILED', 'connectionCloseFailed'],
])('explains the post-delete problem for %s', async (code, message) => {
  const remove = vi.spyOn(studioV2WorkspaceDevicesAPI, 'remove').mockRejectedValueOnce(new DeleteResponseError('deleted', code, 'reload'));
  const onDeleted = deviceDialog([]);
  fireEvent.change(screen.getByLabelText('deleteDialogs.device.typeName'), { target: { value: 'mc_new' } });
  fireEvent.click(screen.getByTestId('device-delete-confirm-button'));
  expect(await screen.findByText(`deleteDialogs.device.${message}`)).toBeInTheDocument();
  expect(screen.queryByText('deleteDialogs.device.unknown')).not.toBeInTheDocument();
  expect(screen.getByRole('button', { name: 'deleteDialogs.device.refresh' })).toBeEnabled();
  expect(onDeleted).not.toHaveBeenCalled(); remove.mockRestore();
});

it('offers device deletion from the devices view', () => {
  const state = { ...initialCollectionState(), mode: 'devices' as const };
  const devices = [{ id: 'mc_new', name: 'mc_new', protocol: 'mc', availability: 'available' }];
  const onDelete = vi.fn();
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(<QueryClientProvider client={client}><RecordingHomeCollection groups={[]} devices={devices} state={state}
    guidance={recordingHomeGuidance(fixture(['mc_new'], false, []))} observedIds={[]} unverifiedObserved={false} sourcesAvailable
    restored={{ state, workspaceRevision: 'W', scrollTop: 0, focusId: null }} onChange={vi.fn()} onOpen={vi.fn()}
    onOpenDevice={vi.fn()} onRemember={vi.fn()} onRecheck={vi.fn()}
    renderDeviceActions={device => <RecordingHomeDeviceActions device={device} available onDelete={onDelete} />} /></QueryClientProvider>);
  fireEvent.click(screen.getByRole('button', { name: /deleteDialogs.device.actions/ }));
  fireEvent.click(screen.getByRole('menuitem', { name: 'deleteDialogs.device.confirm' }));
  expect(onDelete).toHaveBeenCalledExactlyOnceWith(devices[0]);
});
