import { deviceFixture, resetDeviceAutosaveMocks, studioV2WorkspaceAPI, studioV2WorkspaceDevicesAPI, studioV2RulesAPI, studioV2MappingsAPI } from './device-autosave-page.testHarness';
import { useState } from 'react';
import { RecordingHomeRoute } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeRoute';
import { parseStudioEntityLocator } from '@/services/studioEntityLocator';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { readStudioLocatorEditorEntry } from '@/services/studioV2LocatorEntry';
import { locatorRead } from './studio-locator-fixtures';
import { TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { useStudioV2AutosaveState } from '@/pages/datalink/workbench-v2/useStudioV2AutosaveState';
import { RecordingPointSelection } from '@/features/datalink/workbench-v2/recording-editor/RecordingPointSelection';
import { buildGroupCandidates, type GroupCandidate } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { sourceRuleCandidateAPI } from '@/services/sourceRuleCandidates';
import { pointAPI, tagAPI } from '@/services/datalink';
import type { SourceRuleRecord, StudioV2WorkspaceMappingRecord } from '@/types/datalink';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/sourceRuleCandidates', () => ({ sourceRuleCandidateAPI: { get: vi.fn() } }));
vi.mock('@/services/datalink', async original => { const actual = await original<typeof import('@/services/datalink')>(); return { ...actual, pointAPI: { ...actual.pointAPI, get: vi.fn() }, tagAPI: { ...actual.tagAPI, get: vi.fn() } }; });
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn() } }));
vi.mock('@/services/studioV2LocatorEntry', () => ({ readStudioLocatorEditorEntry: vi.fn() }));
vi.mock('@/features/datalink/workbench-v2/recording-home/recordingHomeSources', () => ({ readRecordingHomeSources: async (scope: object) => ({ state: 'available', scope, value: { candidates: [], excluded: [], devices: [device] } }) }));
vi.mock('@/features/datalink/workbench-v2/recording-home/useRecordingHomeFacts', () => ({ useRecordingHomeFacts: () => ({ owners: {}, runtimeGroups: {}, histories: {}, pendingDevices: [], retry: vi.fn() }) }));
const locator = locatorRead(4);
function HomeHarness() {
 const [target, setTarget] = useState('');
 if (!target) return <RecordingHomeRoute location={{ search: '?view=recordings', hash: '' }} onNavigate={setTarget} />;
 const intent = parseStudioEntityLocator(new URL(target, 'https://example.test').search);
 if (intent.kind !== 'locator' || intent.locator.workspace_id !== 'W' || intent.locator.group_id !== 'G') throw new Error('Unexpected locator');
 return <LocatorHarness />;
}
function LocatorHarness() {
 const owner = useStudioV2AutosaveState(false, locator, false, undefined, 'W');
 return <Step4Database state={owner.state} dispatch={owner.actions.dispatch} workspaceId="W" locatorOnly locatorEntry={locator} deviceSource={owner.deviceSource} />;
}
const device = deviceFixture({ id: 'dev-1', name: 'PLC' });
let sourceDevice = device;
let rule: SourceRuleRecord | undefined;
let mapping: StudioV2WorkspaceMappingRecord | undefined;
function Harness({ extra = [] }: { extra?: GroupCandidate[] }) {
 const owner = useStudioV2AutosaveState(true, undefined, false, undefined, 'W');
 if (!owner.workspaceHydrated) return <p>Loading</p>;
 return <><input aria-label="Recording draft" defaultValue="Keep my recording" /><RecordingPointSelection scope="S" workspaceId="W" candidates={[...extra, ...buildGroupCandidates(owner.state).candidates]}
  members={[]} onChange={vi.fn()} devices={owner.state.devices} disabled={false} deviceSource={owner.deviceSource} /></>;
}
beforeEach(() => {
 resetDeviceAutosaveMocks(); vi.clearAllMocks(); sourceDevice = device; rule = undefined; mapping = undefined;
 studioV2WorkspaceAPI.get.mockResolvedValue({ id: 'W', kind: 'single', status: 'ready', ordered_device_ids: ['dev-1'], created_at: '', updated_at: '', modbus_share: { workspace_id: 'W', workspace_revision: 'SOURCE-REV' } });
 studioV2WorkspaceDevicesAPI.list.mockResolvedValue([device]);
 studioV2RulesAPI.list.mockImplementation(async () => rule ? [rule] : []);
 studioV2MappingsAPI.list.mockImplementation(async () => mapping ? [mapping] : []);
 studioV2RulesAPI.create.mockImplementation(async request => {
  rule = { ...request, naming_prefix: request.naming_prefix.toUpperCase().replace(/ /g, '-'), start_address: request.start_address.trim().toUpperCase(), id: request.id!, workspace_id: 'W', revision_id: 'R1', created_at: '', updated_at: '' } as SourceRuleRecord;
  return { ...rule, runtime_apply_status: 'not_running' };
 });
 studioV2MappingsAPI.create.mockImplementation(async request => {
  mapping = { ...request, id: 'MAP', workspace_id: 'W', device_id: sourceDevice.id, point_id: 'POINT', tag_id: 'TAG', created_at: '', updated_at: '' };
  return { ...mapping, runtime_apply_status: 'not_running' };
 });
 vi.mocked(sourceRuleCandidateAPI.get).mockImplementation(async () => ({ source_rule_id: rule!.id, workspace_id: 'W', workspace_revision: 'SOURCE-REV', revision_id: 'R1',
  tags: { status: 'ready', candidates: [{ id: 'CANDIDATE', identity: { source_rule_id: rule!.id, candidate_type: 'tag', candidate_kind: 'point', derived_from_rule_address: '40001' }, proposed_signature: 'SIG',
    address: '40001', point_id: 'POINT', ...(mapping ? { tag_id: 'TAG', mapping_id: 'MAP' } : {}), tag_key: `${rule!.naming_prefix}_40001`, display_name: `${rule!.naming_prefix}_40001`, data_type: 'int16', status: 'validated' }] },
  database_outputs: { status: 'ready', candidates: [] }, local_modbus_outputs: { status: 'deferred', candidates: [] } }));
 vi.mocked(pointAPI.get).mockImplementation(async () => ({ id: 'POINT', device_id: sourceDevice.id, address: '40001', data_type: 'int16', enabled: true } as Awaited<ReturnType<typeof pointAPI.get>>));
 vi.mocked(tagAPI.get).mockImplementation(async () => ({ id: 'TAG', key: mapping!.tag_key, data_type: mapping!.target_type, display_name: mapping!.display_name, unit: mapping!.unit, status: 'draft' } as Awaited<ReturnType<typeof tagAPI.get>>));
});
it('uses original rule and mapping owners and waits for canonical point/tag reads before returning to picker', async () => {
 const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
 render(<QueryClientProvider client={client}><Harness /></QueryClientProvider>);
 fireEvent.click(await screen.findByTestId('recording-points-open'));
 fireEvent.change(screen.getByLabelText('picker.device'), { target: { value: 'dev-1' } });
 fireEvent.click(screen.getByRole('button', { name: 'picker.newSource' }));
 fireEvent.change(await screen.findByLabelText('range.name'), { target: { value: 'Sensor' } });
 fireEvent.change(screen.getByLabelText('range.address'), { target: { value: '40001' } });
 fireEvent.change(screen.getByLabelText('range.count'), { target: { value: '1' } });
 let resolvePoint!: (point: Awaited<ReturnType<typeof pointAPI.get>>) => void;
 vi.mocked(pointAPI.get).mockImplementationOnce(() => new Promise(resolve => { resolvePoint = resolve; }));
 await waitFor(() => expect(screen.getByRole('button', { name: 'range.save' })).toBeEnabled());
 fireEvent.click(screen.getByRole('button', { name: 'range.save' }));
 await waitFor(() => expect(studioV2MappingsAPI.create).toHaveBeenCalledTimes(1));
 expect(screen.getByLabelText('range.count')).toHaveValue('1');
 expect(screen.queryByLabelText('picker.search')).not.toBeInTheDocument();
 await waitFor(() => expect(pointAPI.get).toHaveBeenCalledWith('POINT'));
 await act(async () => resolvePoint({ id: 'POINT', device_id: 'dev-1', address: '40001', data_type: 'int16', enabled: true } as Awaited<ReturnType<typeof pointAPI.get>>));
 await screen.findByLabelText('picker.search');
 expect(screen.getByRole('checkbox', { name: mapping!.display_name })).not.toBeChecked();
 expect(screen.getByLabelText('Recording draft')).toHaveValue('Keep my recording');
 expect(studioV2RulesAPI.create).toHaveBeenCalledTimes(1); expect(studioV2MappingsAPI.create).toHaveBeenCalledTimes(1);
 fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
 expect(studioV2RulesAPI.remove).not.toHaveBeenCalled(); expect(studioV2MappingsAPI.remove).not.toHaveBeenCalled();
});

async function openRange() {
 const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
 render(<QueryClientProvider client={client}><Harness /></QueryClientProvider>);
 fireEvent.click(await screen.findByTestId('recording-points-open'));
 fireEvent.change(screen.getByLabelText('picker.device'), { target: { value: 'dev-1' } });
 fireEvent.click(screen.getByRole('button', { name: 'picker.newSource' }));
 fireEvent.change(await screen.findByLabelText('range.name'), { target: { value: 'Sensor' } });
 fireEvent.change(screen.getByLabelText('range.address'), { target: { value: '40001' } });
 fireEvent.change(screen.getByLabelText('range.count'), { target: { value: '1' } });
 await waitFor(() => expect(screen.getByRole('button', { name: 'range.save' })).toBeEnabled());
}
it('keeps newer invalid text while an older submitted generation finishes all readbacks', async () => {
 await openRange();
 let resolvePoint!: (point: Awaited<ReturnType<typeof pointAPI.get>>) => void;
 vi.mocked(pointAPI.get).mockImplementationOnce(() => new Promise(resolve => { resolvePoint = resolve; }));
 fireEvent.click(screen.getByRole('button', { name: 'range.save' }));
 await waitFor(() => expect(pointAPI.get).toHaveBeenCalledWith('POINT'));
 fireEvent.change(screen.getByLabelText('range.count'), { target: { value: '1x newer' } });
 await act(async () => resolvePoint({ id: 'POINT', device_id: 'dev-1', address: '40001', data_type: 'int16', enabled: true } as Awaited<ReturnType<typeof pointAPI.get>>));
 await screen.findByTestId('recording-range-newer');
 expect(screen.getByLabelText('range.count')).toHaveValue('1x newer');
 expect(screen.getByRole('button', { name: 'range.save' })).toBeDisabled();
 expect(screen.queryByLabelText('picker.search')).not.toBeInTheDocument();
 expect(studioV2RulesAPI.create).toHaveBeenCalledTimes(1);
});
it('does not retransmit an unacknowledged mapping when inspecting or reopening', async () => {
 studioV2MappingsAPI.create.mockRejectedValueOnce(new Error('lost response'));
 await openRange(); fireEvent.click(screen.getByRole('button', { name: 'range.save' }));
 await screen.findByTestId('recording-range-unknown');
 fireEvent.change(screen.getByLabelText('range.count'), { target: { value: 'invalid kept' } });
 fireEvent.click(screen.getByTestId('recording-range-back'));
 fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
 fireEvent.click(screen.getByTestId('recording-points-open'));
 fireEvent.click(screen.getByTestId('recording-range-resume'));
 expect(screen.getByLabelText('range.count')).toHaveValue('invalid kept');
 fireEvent.click(screen.getByTestId('recording-range-inspect'));
 await waitFor(() => expect(screen.getByTestId('recording-range-inspect')).toBeEnabled());
 expect(studioV2MappingsAPI.create).toHaveBeenCalledTimes(1);
 expect(studioV2RulesAPI.create).toHaveBeenCalledTimes(1);
 expect(screen.getByRole('button', { name: 'range.save' })).toBeDisabled();
});

it('Home precise locator saves an empty-device range without replacing the open recording draft', async () => {
 vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue(locator.entry.write_groups!);
 vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: true, value: { ...locator, entry: { ...locator.entry, editor_hydrated: true, database_config: null, database_targets: [] } } });
 const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
 render(<QueryClientProvider client={client}><HomeHarness /></QueryClientProvider>);
 await screen.findByRole('button', { name: 'Shared recorder' });
 await waitFor(() => expect(client.getQueryCache().getAll().filter(q => q.queryKey.includes('write-groups') || q.queryKey.includes('bootstrap')).every(q => q.state.status === 'success' && q.state.fetchStatus === 'idle')).toBe(true));
 fireEvent.click(screen.getByRole('button', { name: 'Shared recorder' }));
 fireEvent.click(await screen.findByTestId('group-open-G'));
 fireEvent.change(await screen.findByTestId('group-name'), { target: { value: 'Keep precise draft' } });
 fireEvent.click(screen.getByTestId('recording-points-open'));
 await screen.findByRole('option', { name: 'PLC' });
 fireEvent.change(screen.getByLabelText('picker.device'), { target: { value: 'dev-1' } });
 fireEvent.click(screen.getByRole('button', { name: 'picker.newSource' }));
 fireEvent.change(await screen.findByLabelText('range.name'), { target: { value: 'Sensor' } });
 fireEvent.change(screen.getByLabelText('range.address'), { target: { value: '40001' } });
 fireEvent.change(screen.getByLabelText('range.count'), { target: { value: '1' } });
 await waitFor(() => expect(screen.getByRole('button', { name: 'range.save' })).toBeEnabled());
 fireEvent.click(screen.getByRole('button', { name: 'range.save' }));
 await screen.findByLabelText('picker.search');
 expect(screen.getByRole('checkbox', { name: mapping!.display_name })).not.toBeChecked();
 expect(screen.getByTestId('group-name')).toHaveValue('Keep precise draft');
 fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
 fireEvent.click(screen.getByTestId('group-editor-close'));
 expect(studioV2RulesAPI.create).toHaveBeenCalledTimes(1); expect(studioV2MappingsAPI.create).toHaveBeenCalledTimes(1);
 expect(studioV2RulesAPI.remove).not.toHaveBeenCalled(); expect(studioV2MappingsAPI.remove).not.toHaveBeenCalled();
 expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});

it('saves a new device then its range through the same stack, without selecting or deleting its saved points', async () => {
 studioV2WorkspaceDevicesAPI.create.mockImplementation(async payload => {
  sourceDevice = deviceFixture({ id: payload.id, name: payload.name, protocol: payload.protocol, connection_config: JSON.stringify(payload.connection_config) });
  studioV2WorkspaceDevicesAPI.list.mockResolvedValue([device, sourceDevice]);
  studioV2WorkspaceAPI.get.mockResolvedValue({ id: 'W', kind: 'single', status: 'ready', ordered_device_ids: [device.id, sourceDevice.id], created_at: '', updated_at: '', modbus_share: { workspace_id: 'W', workspace_revision: 'SOURCE-REV' } });
  return sourceDevice;
 });
 const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
 render(<QueryClientProvider client={client}><Harness /></QueryClientProvider>);
 fireEvent.click(await screen.findByTestId('recording-points-open'));
 fireEvent.click(screen.getByRole('button', { name: 'picker.newSource' }));
 fireEvent.click(screen.getByRole('button', { name: 'source.newDevice' }));
 fireEvent.change(await screen.findByTestId('recording-source-name'), { target: { value: 'New source PLC' } });
 fireEvent.change(screen.getByTestId('input-host'), { target: { value: '192.0.2.20' } });
 fireEvent.click(screen.getByTestId('recording-source-save'));
 fireEvent.change(await screen.findByLabelText('range.name'), { target: { value: 'New sensor' } });
 fireEvent.change(screen.getByLabelText('range.address'), { target: { value: '40001' } });
 fireEvent.change(screen.getByLabelText('range.count'), { target: { value: '1' } });
 await waitFor(() => expect(screen.getByRole('button', { name: 'range.save' })).toBeEnabled());
 expect(screen.getAllByRole('dialog')).toHaveLength(1);
 fireEvent.click(screen.getByRole('button', { name: 'range.save' }));
 await screen.findByLabelText('picker.search');
 expect(screen.getByLabelText('picker.device')).toHaveValue(sourceDevice.id);
 expect(screen.getByRole('checkbox', { name: mapping!.display_name })).not.toBeChecked();
 fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
 expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledTimes(1);
 expect(studioV2RulesAPI.create).toHaveBeenCalledTimes(1); expect(studioV2MappingsAPI.create).toHaveBeenCalledTimes(1);
 expect(studioV2WorkspaceDevicesAPI.remove).not.toHaveBeenCalled(); expect(studioV2RulesAPI.remove).not.toHaveBeenCalled();
});

it('locates confirmed new points on their page, clears hiding filters with a reason, and keeps prior pending selection', async () => {
 const extra: GroupCandidate[] = Array.from({ length: 30 }, (_, i) => ({ key: `A${i}|T${i}`, device_id: 'dev-1', point_id: `A${i.toString().padStart(2, '0')}`, tag_id: `T${i}`, tag_key: `prior_${i}`, label: `Prior ${i}`, device_name: 'PLC', address: String(40100 + i), target_type: 'int16' }));
 const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
 render(<QueryClientProvider client={client}><Harness extra={extra} /></QueryClientProvider>);
 fireEvent.click(await screen.findByTestId('recording-points-open'));
 fireEvent.click(screen.getByRole('checkbox', { name: 'Prior 0' }));
 fireEvent.change(screen.getByLabelText('picker.device'), { target: { value: 'dev-1' } });
 fireEvent.change(screen.getByLabelText('picker.search'), { target: { value: 'hides everything' } });
 fireEvent.change(screen.getByLabelText('picker.pointType'), { target: { value: 'counter' } });
 fireEvent.click(screen.getByRole('button', { name: 'picker.newSource' }));
 fireEvent.change(await screen.findByLabelText('range.name'), { target: { value: 'Sensor' } });
 fireEvent.change(screen.getByLabelText('range.address'), { target: { value: '40001' } });
 fireEvent.change(screen.getByLabelText('range.count'), { target: { value: '1' } });
 await waitFor(() => expect(screen.getByRole('button', { name: 'range.save' })).toBeEnabled());
 fireEvent.click(screen.getByRole('button', { name: 'range.save' }));
 await screen.findByLabelText('picker.search');
 expect(screen.getByLabelText('picker.search')).toHaveValue('');
 expect(screen.getByLabelText('picker.pointType')).toHaveValue('all');
 expect(screen.getByRole('checkbox', { name: mapping!.display_name })).not.toBeChecked();
 expect(screen.getByRole('checkbox', { name: mapping!.display_name })).toHaveFocus();
 expect(screen.getByTestId('recording-range-located')).toBeVisible();
 expect(screen.getByRole('button', { name: 'picker.add' })).toBeEnabled();
 fireEvent.click(screen.getByRole('button', { name: 'picker.previous' }));
 expect(screen.getByRole('checkbox', { name: 'Prior 0' })).toBeChecked();
});

it('independent: late range completion must not close a newer device form', async () => {
 await openRange();
 let resolvePoint!: (point: Awaited<ReturnType<typeof pointAPI.get>>) => void;
 vi.mocked(pointAPI.get).mockImplementationOnce(() => new Promise(resolve => { resolvePoint = resolve; }));
 fireEvent.click(screen.getByRole('button', { name: 'range.save' }));
 await waitFor(() => expect(pointAPI.get).toHaveBeenCalledWith('POINT'));
 fireEvent.click(screen.getByTestId('recording-range-back'));
 fireEvent.change(screen.getByLabelText('picker.device'), { target: { value: '' } });
 fireEvent.click(screen.getByRole('button', { name: 'picker.newSource' }));
 fireEvent.click(screen.getByRole('button', { name: 'source.newDevice' }));
 fireEvent.change(await screen.findByTestId('recording-source-name'), { target: { value: 'Newer unsaved device' } });
 await act(async () => resolvePoint({ id: 'POINT', device_id: 'dev-1', address: '40001', data_type: 'int16', enabled: true } as Awaited<ReturnType<typeof pointAPI.get>>));
 expect(screen.getByTestId('recording-source-name')).toBeVisible();
 expect(screen.getByTestId('recording-source-name')).toHaveValue('Newer unsaved device');
 fireEvent.click(screen.getByTestId('recording-source-close'));
 fireEvent.click(screen.getByTestId('recording-source-discard'));
 expect(screen.getByLabelText('picker.device')).toHaveValue('');
 expect(screen.queryByTestId('recording-range-located')).not.toBeInTheDocument();
 expect(screen.getByRole('checkbox', { name: mapping!.display_name })).not.toBeChecked();
 expect(studioV2RulesAPI.create).toHaveBeenCalledTimes(1);
 expect(studioV2MappingsAPI.create).toHaveBeenCalledTimes(1);
});


it('keeps a newer picker search/filter intent when the submitted range settles in the background', async () => {
 await openRange();
 let resolvePoint!: (point: Awaited<ReturnType<typeof pointAPI.get>>) => void;
 vi.mocked(pointAPI.get).mockImplementationOnce(() => new Promise(resolve => { resolvePoint = resolve; }));
 fireEvent.click(screen.getByRole('button', { name: 'range.save' }));
 await waitFor(() => expect(pointAPI.get).toHaveBeenCalledWith('POINT'));
 fireEvent.click(screen.getByTestId('recording-range-back'));
 fireEvent.change(screen.getByLabelText('picker.device'), { target: { value: '' } });
 fireEvent.change(screen.getByLabelText('picker.search'), { target: { value: 'newer search' } });
 fireEvent.change(screen.getByLabelText('picker.pointType'), { target: { value: 'counter' } });
 await act(async () => resolvePoint({ id: 'POINT', device_id: 'dev-1', address: '40001', data_type: 'int16', enabled: true } as Awaited<ReturnType<typeof pointAPI.get>>));
 expect(screen.getByLabelText('picker.search')).toHaveValue('newer search');
 expect(screen.getByLabelText('picker.pointType')).toHaveValue('counter');
 expect(screen.getByLabelText('picker.device')).toHaveValue('');
 expect(screen.queryByTestId('recording-range-located')).not.toBeInTheDocument();
});
