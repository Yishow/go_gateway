import { addRecordingTemperature, openNewRecordingDevice } from './recording-editor-test-actions';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { deviceFixture, resetDeviceAutosaveMocks, studioV2WorkspaceAPI, studioV2WorkspaceDevicesAPI } from './device-autosave-page.testHarness';
import { useState } from 'react';
import { RecordingHomeRoute } from '../../../src/features/datalink/workbench-v2/recording-home/RecordingHomeRoute';
import { parseStudioEntityLocator } from '../../../src/services/studioEntityLocator';
import { useStudioV2AutosaveState } from '../../../src/pages/datalink/workbench-v2/useStudioV2AutosaveState';
import { WriteGroupSection } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../src/services/studioV2WorkspaceWriteGroups';
import { locatorRead } from './studio-locator-fixtures';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import { readStudioLocatorEditorEntry } from '../../../src/services/studioV2LocatorEntry';
import { savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }),
}));
vi.mock('../../../src/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(),
  testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn(),
} }));
vi.mock('../../../src/services/studioV2LocatorEntry', () => ({ readStudioLocatorEditorEntry: vi.fn() }));
vi.mock('../../../src/features/datalink/workbench-v2/recording-home/recordingHomeSources', () => ({ readRecordingHomeSources: async (scope: object) => ({ state: 'available', scope, value: { candidates: [], excluded: [], devices: [] } }) }));
vi.mock('../../../src/features/datalink/workbench-v2/recording-home/useRecordingHomeFacts', () => ({ useRecordingHomeFacts: () => ({ owners: {}, runtimeGroups: {}, histories: {}, pendingDevices: [], retry: vi.fn() }) }));
const locator = locatorRead(4);
function HomeToLocatorHarness() {
  const [target, setTarget] = useState('');
  if (!target) return <RecordingHomeRoute location={{ search: '?view=recordings', hash: '' }} onNavigate={setTarget} />;
  const intent = parseStudioEntityLocator(new URL(target, 'https://example.test').search);
  if (intent.kind !== 'locator' || intent.locator.workspace_id !== 'W' || intent.locator.group_id !== 'G') throw new Error('Foreign Home target');
  return <LocatorOwnerHarness />;
}
function LocatorOwnerHarness() {
  const owner = useStudioV2AutosaveState(false, locator, false, undefined, 'W');
  return <><output data-testid="locator-owner">{JSON.stringify({ devices: owner.state.devices.map(d => d.name), current: owner.state.current, selected: owner.state.selectedRuleId })}</output>
    <Step4Database state={owner.state} dispatch={owner.actions.dispatch} workspaceId="W" locatorOnly locatorEntry={locator} deviceSource={owner.deviceSource} /></>;
}
const groups = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const existing = deviceFixture({ id: 'dev-1' });

/** Real draft/save owner and real source autosave owner; only remote services are mocked. */
function OwnerHarness() {
  const owner = useStudioV2AutosaveState(true, undefined, false, undefined, 'ws-1');
  if (!owner.workspaceHydrated) return <p>Loading owner</p>;
  return <>
    <output data-testid="persisted-owner-devices">{owner.state.devices.filter(device => device.persisted).map(device => device.name).join('|')}</output>
    <WriteGroupSection state={savedGroupState({ devices: owner.state.devices })} workspaceId="ws-1" readonly={false} deviceSource={owner.deviceSource} />
  </>;
}
beforeEach(() => {
  resetDeviceAutosaveMocks(); vi.mocked(studioV2WorkspaceDevicesAPI.create).mockReset();
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([existing]);
  groups.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] });
});

describe('recording device source real owner integration', () => {
  it('cancels unsent raw input, then explicitly saves through the existing queue and retains the device after recording discard', async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
    render(<QueryClientProvider client={client}><OwnerHarness /></QueryClientProvider>);
    fireEvent.click(await screen.findByTestId('group-create'));
    fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Keep recording name' } });
    addRecordingTemperature('picker.add');
    await openNewRecordingDevice();
    fireEvent.change(screen.getByTestId('recording-source-name'), { target: { value: 'Unsent PLC' } });
    fireEvent.change(screen.getByTestId('input-port'), { target: { value: 'invalid 50x' } });
    expect(screen.getByTestId('input-port')).toHaveValue('invalid 50x');
    expect(screen.getByTestId('recording-source-save')).toBeDisabled();
    fireEvent.click(screen.getByTestId('recording-source-close'));
    fireEvent.click(screen.getByTestId('recording-source-discard'));
    expect(studioV2WorkspaceDevicesAPI.create).not.toHaveBeenCalled();
    expect(studioV2WorkspaceDevicesAPI.update).not.toHaveBeenCalled();
    expect(screen.getByTestId('group-name')).toHaveValue('Keep recording name');
    expect(screen.getByTestId('group-member-row-pt-0')).toBeInTheDocument();
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockImplementation(async payload => {
      const saved = deviceFixture({ id: payload.id, name: payload.name, protocol: payload.protocol, connection_config: JSON.stringify(payload.connection_config) });
      vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([existing, saved]);
      studioV2WorkspaceAPI.get.mockResolvedValue({ ...locator.workspace, ordered_device_ids: [existing.id, saved.id] }); return saved;
    });
    await openNewRecordingDevice();
    fireEvent.change(screen.getByTestId('recording-source-name'), { target: { value: 'Confirmed PLC' } });
    fireEvent.change(screen.getByTestId('input-host'), { target: { value: '192.0.2.10' } });
    fireEvent.click(screen.getByTestId('recording-source-save'));
    expect(await screen.findByTestId('recording-source-confirmed')).toBeVisible();
    expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledOnce();
    await waitFor(() => expect(screen.getByTestId('persisted-owner-devices')).toHaveTextContent('Confirmed PLC'));
    fireEvent.click(screen.getByTestId('recording-range-back'));
    fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
    expect(screen.getByTestId('group-name')).toHaveValue('Keep recording name');
    expect(screen.getByTestId('group-member-row-pt-0')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('group-editor-close'));
    expect(screen.queryByTestId('group-editor')).not.toBeInTheDocument();
    expect(screen.getByTestId('persisted-owner-devices')).toHaveTextContent('Confirmed PLC');
    expect(studioV2WorkspaceDevicesAPI.remove).not.toHaveBeenCalled();
    expect(groups.remove).not.toHaveBeenCalled();
    expect(groups.create).not.toHaveBeenCalled();
  });
  it.each([false, true])('Home precise locator source saves with safe unknown recovery=%s and preserves draft', async recover => {
    groups.list.mockResolvedValue(locator.entry.write_groups!);
    studioV2WorkspaceAPI.get.mockResolvedValue({ ...locator.workspace, ordered_device_ids: [existing.id] });
    vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: true, value: { ...locator, entry: { ...locator.entry,
      editor_hydrated: true, database_config: null, database_targets: [] } } });
    const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
    render(<QueryClientProvider client={client}><HomeToLocatorHarness /></QueryClientProvider>);
    await screen.findByRole('button', { name: 'Shared recorder' });
    await waitFor(() => expect(client.getQueryCache().getAll().filter(q => q.queryKey.includes('write-groups') || q.queryKey.includes('bootstrap')).every(q => q.state.status === 'success' && q.state.fetchStatus === 'idle')).toBe(true));
    fireEvent.click(screen.getByRole('button', { name: 'Shared recorder' }));
    studioV2WorkspaceAPI.get.mockClear();
    fireEvent.click(await screen.findByTestId('group-open-G'));
    const name = await screen.findByTestId('group-name');
    fireEvent.change(name, { target: { value: 'Keep precise draft' } });
    expect(screen.getByTestId('recording-points-open')).toBeEnabled();
    expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
    await openNewRecordingDevice();
    fireEvent.change(await screen.findByTestId('recording-source-name'), { target: { value: 'Locator PLC' } });
    fireEvent.change(screen.getByTestId('input-host'), { target: { value: '192.0.2.20' } });
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockImplementation(async payload => {
      const saved = deviceFixture({ id: payload.id, name: payload.name, protocol: payload.protocol, connection_config: JSON.stringify(payload.connection_config) });
      vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([existing, saved]);
      studioV2WorkspaceAPI.get.mockResolvedValue({ ...locator.workspace, id: recover ? 'FOREIGN' : 'W', ordered_device_ids: [existing.id, saved.id] }); return saved;
    });
    fireEvent.click(screen.getByTestId('recording-source-save'));
    if (recover) {
      await screen.findByTestId('recording-source-unknown');
      const saved = await vi.mocked(studioV2WorkspaceDevicesAPI.create).mock.results[0].value;
      studioV2WorkspaceAPI.get.mockResolvedValue({ ...locator.workspace, ordered_device_ids: [existing.id, saved.id] });
      fireEvent.click(screen.getByTestId('recording-source-inspect'));
    }
    expect(await screen.findByTestId('recording-source-confirmed')).toBeVisible();
    expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledOnce();
    expect(studioV2WorkspaceAPI.get).toHaveBeenCalled();
    expect(screen.getByTestId('group-name')).toHaveValue('Keep precise draft');
    expect(screen.getByTestId('locator-owner')).toHaveTextContent('Locator PLC');
    expect(screen.getByTestId('locator-owner')).toHaveTextContent('"current":4');
    expect(groups.update).not.toHaveBeenCalled();
  });

  it.each(['foreign', 'missing-revision', 'stale-revision'])('blocks %s source lease without replacing the live locator draft', async reason => {
    groups.list.mockResolvedValue({ ...locator.entry.write_groups!, workspace_revision: reason === 'missing-revision' ? '' : reason === 'stale-revision' ? 'STALE' : 'W2' });
    studioV2WorkspaceAPI.get.mockResolvedValue({ ...locator.workspace, id: reason === 'foreign' ? 'OTHER' : 'W', ordered_device_ids: [existing.id] });
    vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: true, value: { ...locator, entry: { ...locator.entry,
      editor_hydrated: true, database_config: null, database_targets: [] } } });
    const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
    render(<QueryClientProvider client={client}><LocatorOwnerHarness /></QueryClientProvider>);
    fireEvent.click(await screen.findByTestId('group-open-G'));
    fireEvent.change(await screen.findByTestId('group-name'), { target: { value: 'Keep blocked draft' } });
    fireEvent.click(screen.getByTestId('recording-points-open'));
    fireEvent.click(await screen.findByRole('button', { name: 'picker.newSource' }));
    fireEvent.click(screen.getByRole('button', { name: 'source.newDevice' }));
    expect(await screen.findByText('device.scopeBlocked')).toBeVisible();
    expect(screen.queryByTestId('recording-source-name')).not.toBeInTheDocument();
    expect(screen.getByTestId('group-name')).toHaveValue('Keep blocked draft');
    expect(studioV2WorkspaceDevicesAPI.create).not.toHaveBeenCalled();
    expect(groups.update).not.toHaveBeenCalled();
  });

});
