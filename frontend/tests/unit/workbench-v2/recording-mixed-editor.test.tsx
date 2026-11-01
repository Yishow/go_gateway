import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { GroupEditor } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { buildGroupCandidates } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../src/services/studioV2WorkspaceWriteGroups';
import { studioV2WorkspaceMeasurementsAPI } from '../../../src/services/studioV2WorkspaceMeasurements';
import type { WriteGroup, WriteGroupMutationResponse } from '../../../src/types/studioV2WriteGroup';
import type { WriteGroupCalculationCapabilities } from '../../../src/types/studioV2WriteGroupCalculation';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }) }));
vi.mock('../../../src/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
 readiness: vi.fn().mockResolvedValue(null), delivery: vi.fn().mockResolvedValue(null), list: vi.fn(), get: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(),
} }));
vi.mock('../../../src/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn().mockResolvedValue([]) } }));
const capabilities: WriteGroupCalculationCapabilities = { version: 1, result_shape: 'member-calculation-v1', editor_supported: true, draft_supported: true, activation_supported: false, reason: 'record_readiness_required', methods: ['latest', 'mean', 'delta'], cadence: { alignment: 'utc_epoch', units: ['second', 'minute'], minimum: 1, max_seconds: 600, max_minutes: 10, window_bounds_checked: true }, schema_dialects: ['sqlite'], schema_layouts: ['managed_wide_profile_v1', 'managed_long_v1'], unsupported_combinations: [] };
const base = savedGroup();
const group = savedGroup({ record_id: '11111111-1111-4111-8111-111111111111', layout_version: 'managed_wide_profile_v1', recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' }, destination: { ...base.destination, storage_strategy: 'managed' }, row_policy: { allowed_lateness_seconds: 0, incomplete_policy: 'partial' }, members: base.members.map((member, index) => ({ ...member, calculation: 'latest', cadence: { interval: index ? 60 : 15, unit: 'second', alignment: 'utc_epoch' }, value_role: 'gauge', source_value_type: 'float64' })) });
function mount(value = group, capability = capabilities, newRecord = false, preflight?: Promise<{ workspace_revision: string; groups: WriteGroup[] }>) {
 const state = savedGroupState(); const onSaved = vi.fn();
 const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
 render(<QueryClientProvider client={client}><GroupEditor group={newRecord ? undefined : value} calculationCapabilities={capability} workspaceId="ws-1" workspaceRevision="db-setup-9" connector={state.db.connector} candidates={buildGroupCandidates(state).candidates} excluded={[]} readonly={false} onSaved={onSaved} onClose={vi.fn()} onReload={vi.fn()} fetchLatest={async () => preflight ?? ({ workspace_revision: 'db-fresh-10', groups: [value] })} /></QueryClientProvider>);
 return onSaved;
}
beforeEach(() => { vi.clearAllMocks(); vi.mocked(studioV2WorkspaceWriteGroupsAPI.create).mockImplementation(async request => ({ workspace_revision: 'created-1', group: { ...group, id: 'new-group', record_id: request.group.record_id } })); vi.mocked(studioV2WorkspaceMeasurementsAPI.list).mockResolvedValue([]); vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockResolvedValue({ workspace_revision: 'db-saved-11', group: { ...group, revision: 'rev-2' } }); });
it('edits each canonical mixed member in the existing save owner without synthetic group interval or implicit Apply', async () => {
 const onSaved = mount();
 const row = await screen.findByTestId('member-calculation-pt-0');
 fireEvent.change(within(row).getByLabelText('calculation.interval'), { target: { value: '30' } });
 fireEvent.click(screen.getByTestId('group-save'));
 await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce());
 const request = vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mock.calls[0];
 expect(JSON.stringify(request)).toContain('db-fresh-10');
 const wire = request[1].group;
 expect(wire.record_id).toBe(group.record_id); expect(wire.members.map(m => m.cadence?.interval)).toEqual([30, 60]);
 expect(wire.row_policy).not.toHaveProperty('interval_seconds'); expect(wire.recording_profile).toEqual({ kind: 'mixed', result_shape: 'member-calculation-v1' });
 expect(studioV2WorkspaceWriteGroupsAPI.apply).not.toHaveBeenCalled(); await waitFor(() => expect(onSaved).toHaveBeenCalledOnce());
});
it('does not mount a mixed editor when editor support is false despite draft support', () => {
 mount(group, { ...capabilities, editor_supported: false }); expect(screen.getByTestId('mixed-group-editor-readonly')).toBeInTheDocument(); expect(screen.queryByTestId('group-save')).not.toBeInTheDocument();
});
it('keeps an existing legacy profile interval and identity immutable rather than converting it in place', () => {
 mount(savedGroup({ ...group, record_id: undefined, recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'mean' }, members: base.members }));
 expect(screen.queryByTestId('member-calculation-pt-0')).not.toBeInTheDocument(); expect(screen.getByTestId('group-profile-interval')).toHaveTextContent('15');
 expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});
it('blocks an invalid lexical interval without losing it or issuing Save', async () => {
 mount(); const row = await screen.findByTestId('member-calculation-pt-0');
 fireEvent.change(within(row).getByLabelText('calculation.interval'), { target: { value: '1e' } });
 expect(within(row).getByLabelText('calculation.interval')).toHaveValue('1e'); expect(screen.getByTestId('group-save')).toBeDisabled();
 fireEvent.click(screen.getByTestId('group-save')); expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});

it('persists mean/mean/delta as independent canonical member windows through the same save caller', async () => {
 const third = { ...group.members[0], point_id: 'pt-2', tag_id: 'tag-2', measurement_id: 'counter-2', target_column: 'counter', source_value_type: 'uint64' as const, value_role: 'counter_raw' as const, cadence: { interval: 1, unit: 'minute' as const, alignment: 'utc_epoch' as const } };
 const value = { ...group, members: [...group.members, third] };
 vi.mocked(studioV2WorkspaceMeasurementsAPI.list).mockResolvedValue(value.members.map((member, index) => ({ id: index === 2 ? 'counter-2' : `gauge-${index}`, workspace_id: 'ws-1', device_id: member.device_id, point_id: member.point_id, tag_id: member.tag_id, equipment_id: 'equipment', name: `Source ${index}`, quantity: index === 2 ? 'energy' : 'temperature', unit: index === 2 ? 'kWh' : 'C', semantic_kind: index === 2 ? 'counter' as const : 'gauge' as const, numeric_encoding: index === 2 ? 'uint64' : 'float64', definition_revision: 'def-1', source_binding_revision: 'binding-1', series_epoch: 'epoch-1', counter_policy: { allow_negative: false } })));
 mount(value);
 for (const [index, method] of ['mean', 'mean', 'delta'].entries()) {
  const row = await screen.findByTestId(`member-calculation-pt-${index}`);
  await waitFor(() => expect(within(row).getByRole('option', { name: `calculation.${method}` })).not.toBeDisabled());
  fireEvent.change(within(row).getByLabelText('calculation.title'), { target: { value: method } });
 }
 fireEvent.click(screen.getByTestId('group-save')); await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce());
 const wire = vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mock.calls[0][1].group;
 expect(wire.members.map(member => member.calculation)).toEqual(['mean', 'mean', 'delta']);
 expect(wire.members.map(member => member.cadence)).toEqual(value.members.map(member => member.cadence));
 expect(wire.members[2].measurement_id).toBe('counter-2'); expect(wire.members.every(member => member.result_value_type === undefined)).toBe(true);
 expect(studioV2WorkspaceWriteGroupsAPI.apply).not.toHaveBeenCalled();
});

it('creates an independent mixed identity only by explicit choice and initializes only newly selected members', async () => {
 vi.mocked(studioV2WorkspaceWriteGroupsAPI.create).mockImplementation(async request => ({ workspace_revision: 'created-1', group: { ...group, id: 'new-group', record_id: request.group.record_id } }));
 mount(group, capabilities, true);
 fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Independent' } });
 fireEvent.click(screen.getByRole('button', { name: 'calculation.chooseMixed' }));
 fireEvent.click(screen.getByTestId('group-member-include-pt-0'));
 expect(await screen.findByTestId('member-calculation-pt-0')).toBeInTheDocument();
 fireEvent.click(screen.getByTestId('group-save')); await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.create).toHaveBeenCalledOnce());
 const wire = vi.mocked(studioV2WorkspaceWriteGroupsAPI.create).mock.calls[0][0].group;
 expect(wire.name).toBe('Independent'); expect(wire.record_id).toMatch(/^[0-9a-f-]{36}$/); expect(wire.record_id).not.toBe(group.record_id);
 expect(wire.members).toHaveLength(1); expect(wire.members[0]).toMatchObject({ calculation: 'latest', cadence: { interval: 1, unit: 'second', alignment: 'utc_epoch' }, source_value_type: 'float64' });
 expect(wire.members[0]).not.toHaveProperty('measurement_id'); expect(wire.row_policy).not.toHaveProperty('interval_seconds');
 expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled(); expect(studioV2WorkspaceWriteGroupsAPI.apply).not.toHaveBeenCalled();
});

it.each(['success', 'failure'] as const)('locks form edits during preflight and request pending, then releases after %s', async outcome => {
 let finishRead!: (value: { workspace_revision: string; groups: WriteGroup[] }) => void;
 const preflight = new Promise<{ workspace_revision: string; groups: WriteGroup[] }>(resolve => { finishRead = resolve; });
 let finishSave!: (value: WriteGroupMutationResponse) => void; let failSave!: (cause: Error) => void;
 vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockImplementation(() => new Promise((resolve, reject) => { finishSave = resolve; failSave = reject; }));
 mount(group, capabilities, false, preflight); const row = await screen.findByTestId('member-calculation-pt-0');
 const interval = within(row).getByLabelText('calculation.interval'); const name = screen.getByTestId('group-name');
 fireEvent.change(interval, { target: { value: '30' } }); fireEvent.click(screen.getByTestId('group-save'));
 expect(interval).toBeDisabled(); expect(name).toBeDisabled(); expect(screen.getByTestId('group-save')).toBeDisabled();
 fireEvent.change(interval, { target: { value: '90' } }); fireEvent.change(name, { target: { value: 'Must not overwrite submitted' } });
 await act(async () => { finishRead({ workspace_revision: 'fresh-preflight', groups: [group] }); });
 await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce());
 expect(interval).toBeDisabled(); expect(name).toBeDisabled(); fireEvent.change(interval, { target: { value: '120' } });
 const wire = vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mock.calls[0][1].group;
 expect(wire.name).toBe(group.name); expect(wire.members[0].cadence?.interval).toBe(30);
 await act(async () => { if (outcome === 'failure') failSave(new Error('save unavailable')); else finishSave({ workspace_revision: 'saved-after-lock', group: { ...group, revision: 'saved-2', members: group.members.map((member, index) => index ? member : { ...member, cadence: { interval: 30, unit: 'second', alignment: 'utc_epoch' } }) } }); });
 await waitFor(() => expect(name).not.toBeDisabled()); expect(interval).not.toBeDisabled();
 expect(interval).toHaveValue('30'); fireEvent.change(name, { target: { value: 'After settle' } }); expect(name).toHaveValue('After settle');
});

it('saves explicit unbound gauge role and mean through the existing caller without writing a global definition or UI confirmation token', async () => {
 mount(group, capabilities, true); fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Explicit gauge' } });
 fireEvent.click(screen.getByRole('button', { name: 'calculation.chooseMixed' })); fireEvent.click(screen.getByTestId('group-member-include-pt-0'));
 const row = await screen.findByTestId('member-calculation-pt-0');
 await waitFor(() => expect(within(row).getByRole('button', { name: 'calculation.confirmGauge' })).toBeEnabled());
 fireEvent.click(within(row).getByRole('button', { name: 'calculation.confirmGauge' }));
 fireEvent.change(within(row).getByLabelText('calculation.title'), { target: { value: 'mean' } });
 fireEvent.click(screen.getByTestId('group-save')); await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.create).toHaveBeenCalledOnce());
 const source = vi.mocked(studioV2WorkspaceWriteGroupsAPI.create).mock.calls[0][0].group.members[0];
 expect(source).toMatchObject({ calculation: 'mean', value_role: 'gauge', source_value_type: 'float64' });
 expect(source).not.toHaveProperty('measurement_id'); expect(source).not.toHaveProperty('gauge_confirmation');
 expect(studioV2WorkspaceWriteGroupsAPI.apply).not.toHaveBeenCalled();
});
