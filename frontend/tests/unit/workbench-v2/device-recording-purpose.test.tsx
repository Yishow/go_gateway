import { useMemo } from 'react';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { guidedReadiness, guidedLifecycle } from '../../fixtures/guidedRecordingPreparation';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import type { Mapping, Point, WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import type { MeasurementDefinition } from '@/types/measurement';
import type { WriteGroup, WriteGroupDraft, WriteGroupMemberDraft } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import { expandSources } from '../../fixtures/sourceReview';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: Record<string, unknown>) =>
      options && 'count' in options ? `${key}:${options.count}` : key,
  }),
}));

const groupsQuery = {
  data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] as WriteGroup[] },
  isLoading: false,
  isError: false,
  refetch: vi.fn(),
};
const measurementsQuery = { data: [] as MeasurementDefinition[], isLoading: false, isError: false, refetch: vi.fn() };
const ensureCircuitAsync = vi.fn();
const updateGroupAsync = vi.fn();
const startAsync = vi.fn();
const createMeasurementAsync = vi.fn();

vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupsQuery: () => groupsQuery,
  useWriteGroupReadinessQuery: (id?: string) => { const group = groupsQuery.data.groups.find(g => g.id === id), revision = groupsQuery.data.workspace_revision; return { data: useMemo(() => group ? guidedReadiness(group, revision) : undefined, [group, revision]), refetch: groupsQuery.refetch };  },
  useWriteGroupDeliveryQuery: () => ({ data: undefined, isError: false }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: updateGroupAsync, isPending: false }),
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: ensureCircuitAsync, isPending: false }),
  useWriteGroupSchemaPreviewMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useWriteGroupSchemaApplyMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useWriteGroupSchemaOperationQuery: () => ({ data: undefined }),
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useEnsureBasicManagedMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useRecordingStartMutation: () => ({ mutateAsync: startAsync, isPending: false }),
  useRecordingStartOperationQuery: () => ({ data: undefined, isLoading: false, isError: false, refetch: vi.fn() }),
}));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => measurementsQuery,
  useCreateStudioV2MeasurementMutation: () => ({ mutateAsync: createMeasurementAsync, isPending: false }),
  useUpdateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));
vi.mock('@/hooks/datalink/useStudioV2WriteGroupLifecycle', () => ({ useWriteGroupLifecycleQuery: (_ws?: string, id?: string) => { const group = groupsQuery.data.groups.find(g => g.id === id), revision = groupsQuery.data.workspace_revision; return { data: useMemo(() => group ? guidedLifecycle(group, revision) : undefined, [group, revision]), refetch: groupsQuery.refetch }; } }));

type SourceSpec = { name: string; type: Mapping['target_type']; unit?: string };

/** A full electric meter (A/V/kW/kWh/temperature) plus water and gas counters, all saved. */
const METER: SourceSpec[] = [
  { name: 'voltage', type: 'float64', unit: 'V' }, { name: 'current', type: 'float64', unit: 'A' },
  { name: 'power', type: 'float64', unit: 'kW' }, { name: 'energy', type: 'uint64', unit: 'kWh' },
  { name: 'temperature', type: 'float64', unit: 'C' }, { name: 'water', type: 'float64', unit: 'm3' },
  { name: 'gas', type: 'uint64', unit: 'm3' },
];

function meterState(): WorkbenchV2State {
  const base = savedGroupState();
  const points: Point[] = [];
  const mappings: Record<string, Mapping> = {};
  METER.forEach((source, index) => {
    const id = `rule-1-p-${index}`;
    points.push({
      id, device_id: 'dev-1', rule_id: 'rule-1', rule_name: 'Rule 1', name: source.name, address: `4000${index + 1}`, data_type: source.type,
      function: 'holding_register', width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0,
    });
    mappings[id] = {
      point_id: id, tag_key: source.name, display_name: source.name, unit: source.unit ?? '', target_type: source.type, scale: 1, offset: 0,
      enabled: true, persisted: true, save_state: 'saved', persisted_point_id: `pt-${index}`, tag_id: `tag-${index}`, device_id: 'dev-1',
    };
  });
  return { ...base, points, mappings };
}

/** A saved, confirmed counter definition of one source; the quantity and unit are the fixture's own data, never inferred from a name. */
function confirmedCounter(index: number): MeasurementDefinition {
  const source = METER[index];
  const quantity = source.unit === 'kWh' ? 'energy' : 'volume';
  return {
    id: `meas-tag-${index}`, workspace_id: 'ws-1', definition_revision: 'mrev-1', name: `${source.name} counter`, equipment_id: 'eq-1',
    device_id: 'dev-1', point_id: `pt-${index}`, tag_id: `tag-${index}`, source_binding_revision: 's1', quantity,
    semantic_kind: 'counter', unit: source.unit, series_epoch: '1', counter_policy: { allow_negative: false },
    created_at: '2026-10-20T00:00:00Z', updated_at: '2026-10-20T00:00:00Z',
  };
}

function renderPanel(state: WorkbenchV2State) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  client.setQueryData(studioV2WorkspaceKeys.writeGroups(), groupsQuery.data);
  return render(
    <QueryClientProvider client={client}>
      <BasicRecordingPanel state={state} workspaceId="ws-1" readonly={false} />
    </QueryClientProvider>,
  );
}

/** The groups the preparation sent through ensure, in call order. */
function sentGroups(): WriteGroupDraft[] {
  return ensureCircuitAsync.mock.calls.map((call) => (call[0] as { request: { group: WriteGroupDraft } }).request.group);
}

function membersByTag(group: WriteGroupDraft): Record<string, WriteGroupMemberDraft> {
  return Object.fromEntries(group.members.map((member) => [member.tag_id, member]));
}

beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] };
  groupsQuery.refetch.mockImplementation(async () => ({ data: groupsQuery.data }));
  measurementsQuery.data = [];
  measurementsQuery.refetch.mockImplementation(async () => ({ data: measurementsQuery.data }));
  ensureCircuitAsync.mockImplementation(async ({ request }) => ({
    workspace_revision: 'wrev-2',
    group: savedGroup({ ...request.group, id: `grp-${request.group.recording_profile?.kind}`, revision: 'rev-1', status: 'draft' }),
  }));
});

describe('device recording purposes', () => {
  it('offers the whole-device snapshot (wide) and selected measurements (long) and explains their rows', () => {
    renderPanel(meterState());

    const wide = screen.getByTestId('circuit-profile-purpose-wide');
    const long = screen.getByTestId('circuit-profile-purpose-long');
    expect(wide).toHaveAttribute('type', 'radio');
    expect(long).toHaveAttribute('type', 'radio');
    expect(screen.getByTestId('circuit-profile-purpose-hint-wide')).toHaveTextContent('step4.basic.circuit.purpose_wide_rows');
    expect(screen.getByTestId('circuit-profile-purpose-hint-long')).toHaveTextContent('step4.basic.circuit.purpose_long_rows');

    fireEvent.click(wide);
    expect(wide).toBeChecked();
    expect(long).not.toBeChecked();
    // Choosing a purpose is a draft choice only.
    expect(ensureCircuitAsync).not.toHaveBeenCalled();
    expect(updateGroupAsync).not.toHaveBeenCalled();
    expect(startAsync).not.toHaveBeenCalled();
  });

  it('prepares a new instantaneous recording as an explicit arithmetic mean with a role for every source', async () => {
    renderPanel(savedGroupState());
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensureCircuitAsync).toHaveBeenCalledTimes(1));

    const [group] = sentGroups();
    expect(group.layout_version).toBe('managed_long_v1');
    expect(group.recording_profile).toMatchObject({ kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'mean' });
    const roles = Object.fromEntries(group.members.map((member) => [member.tag_id, member.value_role]));
    // Numeric sources average; a text source is a state and keeps its latest value.
    expect(roles).toEqual({ 'tag-0': 'gauge', 'tag-1': 'gauge', 'tag-2': 'state' });
  });

  it('shows the 15 second arithmetic mean for a new instantaneous recording and never offers it for a counter', () => {
    renderPanel(meterState());

    const label = screen.getByTestId('circuit-profile-reduction-instantaneous');
    expect(label).toHaveTextContent('step4.basic.circuit.reduction_mean:15');
    fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '30' } });
    expect(label).toHaveTextContent('step4.basic.circuit.reduction_mean:30');

    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    expect(screen.getByTestId('circuit-profile-reduction-cumulative')).not.toHaveTextContent('mean');
  });

  it('prepares a whole-device wide snapshot of a full meter with explicit roles and no guessed counter', async () => {
    measurementsQuery.data = [confirmedCounter(3)];
    renderPanel(meterState());

    fireEvent.click(screen.getByTestId('circuit-profile-purpose-wide'));
    expandSources();
    // kWh is a number with a kWh unit, but neither its name nor its unit makes it a counter.
    expect(screen.getByTestId('circuit-profile-role-tag-3')).not.toHaveValue('counter_raw');
    fireEvent.change(screen.getByTestId('circuit-profile-role-tag-3'), { target: { value: 'counter_raw' } });
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensureCircuitAsync).toHaveBeenCalledTimes(1));

    const [group] = sentGroups();
    expect(group.layout_version).toBe('managed_wide_profile_v1');
    expect(group.recording_profile).toMatchObject({ kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'mean' });
    const byTag = membersByTag(group);
    expect(byTag['tag-0']).toMatchObject({ value_role: 'gauge' });
    expect(byTag['tag-4']).toMatchObject({ value_role: 'gauge' });
    expect(byTag['tag-3']).toMatchObject({ value_role: 'counter_raw', measurement_id: 'meas-tag-3' });
    // A gauge is a plain Tag: no measurement definition is referenced or created for it.
    expect(byTag['tag-0'].measurement_id).toBe('');
    expect(byTag['tag-4'].measurement_id).toBe('');
    expect(createMeasurementAsync).not.toHaveBeenCalled();
  });

  it('restores the chosen purpose and roles as the same unsaved draft after a reload', () => {
    measurementsQuery.data = [confirmedCounter(3)];
    const first = renderPanel(meterState());
    fireEvent.click(screen.getByTestId('circuit-profile-purpose-wide'));
    expandSources();
    fireEvent.change(screen.getByTestId('circuit-profile-role-tag-3'), { target: { value: 'counter_raw' } });
    first.unmount();

    renderPanel(meterState());
    expect(screen.getByTestId('circuit-profile-purpose-wide')).toBeChecked();
    expandSources();
    expect(screen.getByTestId('circuit-profile-role-tag-3')).toHaveValue('counter_raw');
    expect(ensureCircuitAsync).not.toHaveBeenCalled();
  });

  it('prepares a cumulative wide snapshot of water, gas and electric counters from confirmed counters only', async () => {
    measurementsQuery.data = [confirmedCounter(3), confirmedCounter(5), confirmedCounter(6)];
    renderPanel(meterState());

    fireEvent.click(screen.getByTestId('circuit-profile-purpose-wide'));
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Utilities' } });
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-instantaneous'));
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    expandSources();
    for (const index of [3, 5, 6]) fireEvent.click(screen.getByTestId(`circuit-profile-member-cum-tag-${index}`));
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensureCircuitAsync).toHaveBeenCalledTimes(1));

    const [group] = sentGroups();
    expect(group.layout_version).toBe('managed_wide_profile_v1');
    expect(group.recording_profile).toMatchObject({ kind: 'cumulative', interval: 10, unit: 'minute' });
    expect(group.recording_profile?.sample_reduction ?? 'latest').toBe('latest');
    expect(group.members.map((member) => member.tag_id).sort()).toEqual(['tag-3', 'tag-5', 'tag-6']);
    for (const member of group.members) {
      expect(member.measurement_id).toBe(`meas-${member.tag_id}`);
      expect(member.value_role ?? '').toBe('');
    }
  });
});

/** A saved group of one circuit on the shared long table (latest, as every group saved before the mean) or the wide profile table. */
function savedProfile(layout: 'managed_long_v1' | 'managed_wide_profile_v1', over: Partial<WriteGroup> = {}): WriteGroup {
  const wide = layout === 'managed_wide_profile_v1';
  return savedGroup({
    id: 'grp-inst', revision: 'rev-inst-1', name: 'Main (Instantaneous)', status: 'ready', layout_version: layout,
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', ...(wide ? { sample_reduction: 'mean' as const } : {}) },
    ...(wide ? {} : { storage_mode: 'plain' as const }),
    destination: { ...savedGroup().destination, storage_strategy: 'managed', table_name: wide ? '' : 'gw_long_fixture' },
    members: [0, 1].map((i) => ({
      device_id: 'dev-1', point_id: `pt-${i}`, tag_id: `tag-${i}`, source_revision: `s${i}`, mapping_revision: `m${i}`, target_column: '', required: false,
      circuit_id: 'cir-main', circuit_name: 'Main', ...(wide ? { value_role: 'gauge' as const, source_value_type: 'float64' as const } : {}),
    })),
    row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    write_policy: wide ? { mode: 'append', dedupe_capability: 'receipt' } : { mode: 'append', dedupe_capability: 'receipt', recreate_missing: true },
    ...over,
  });
}

describe('saved groups keep their meaning and their layout', () => {
  it('keeps a saved latest group latest: no mean label, no role to pick, and nothing is sent for unchanged settings', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [savedProfile('managed_long_v1')] };
    renderPanel(savedGroupState());

    expect(screen.getByTestId('circuit-profile-reduction-instantaneous')).toHaveTextContent('step4.basic.circuit.reduction_latest');
    expect(screen.queryByTestId('circuit-profile-role-tag-0')).toBeNull();
    // An unchanged saved group is inspected without a save/preview command (guide 2.3).
    await waitFor(() => expect(screen.getByTestId('circuit-prepared-group-instantaneous')).toBeInTheDocument());
    expect(ensureCircuitAsync).not.toHaveBeenCalled();
    expect(updateGroupAsync).not.toHaveBeenCalled();

    // A period edit is saved as an update that still says nothing about a mean or a role.
    updateGroupAsync.mockImplementation(async ({ request }) => ({ workspace_revision: 'wrev-2', group: savedProfile('managed_long_v1', { ...request.group, revision: 'rev-inst-2' }) }));
    fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '30' } });
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(updateGroupAsync).toHaveBeenCalledTimes(1));
    const sent = updateGroupAsync.mock.calls[0][0].request.group as WriteGroupDraft;
    expect(sent.recording_profile).toEqual({ kind: 'instantaneous', interval: 30, unit: 'second' });
    expect(sent.members.every((member) => member.value_role === undefined)).toBe(true);
    expect(ensureCircuitAsync).not.toHaveBeenCalled();
  });

  it.each([
    ['managed_long_v1', 'long', 'wide'],
    ['managed_wide_profile_v1', 'wide', 'long'],
  ] as const)('fixes a saved %s circuit to its %s purpose; the other purpose is not offered as an edit', (layout, kept, other) => {
    groupsQuery.data = { ...groupsQuery.data, groups: [savedProfile(layout)] };
    renderPanel(savedGroupState());

    expect(screen.getByTestId(`circuit-profile-purpose-${kept}`)).toBeChecked();
    expect(screen.getByTestId(`circuit-profile-purpose-${other}`)).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-purpose-locked')).toHaveTextContent('step4.basic.circuit.purpose_locked');
    expect(ensureCircuitAsync).not.toHaveBeenCalled();
    expect(updateGroupAsync).not.toHaveBeenCalled();
  });

  it('shows a saved wide profile as the explicit mean it was saved as, with its own roles', () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [savedProfile('managed_wide_profile_v1')] };
    renderPanel(savedGroupState());

    expect(screen.getByTestId('circuit-profile-reduction-instantaneous')).toHaveTextContent('step4.basic.circuit.reduction_mean:15');
    expandSources();
    expect(screen.getByTestId('circuit-profile-role-tag-0')).toHaveValue('gauge');
    // Wide profile tables have no storage mode or recreate policy to choose.
    expect(screen.queryByTestId('circuit-profile-storage-mode')).toBeNull();
    expect(screen.queryByTestId('circuit-profile-recreate-missing')).toBeNull();
  });
});

describe('both purposes share the same confirmations and bounds', () => {
  it.each(['long', 'wide'] as const)('blocks a cumulative %s preparation until its counter is confirmed', (purpose) => {
    renderPanel(meterState());
    if (purpose === 'wide') fireEvent.click(screen.getByTestId('circuit-profile-purpose-wide'));
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Utilities' } });
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-instantaneous'));
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    fireEvent.click(screen.getByTestId('circuit-profile-member-cum-tag-3'));

    expect(screen.getByTestId('circuit-counter-confirmation')).toBeInTheDocument();
    const repair = screen.getByRole('button', { name: 'next.actions.resolve_semantics' });
    fireEvent.click(repair); // Pure repair navigation, never inferred counter confirmation.
    expect(screen.queryByTestId('circuit-profile-prepare')).not.toBeInTheDocument();
    expect(startAsync).not.toHaveBeenCalled();
    expect(ensureCircuitAsync).not.toHaveBeenCalled();
  });

  it('blocks a counter_raw role until that source has a confirmed counter', () => {
    renderPanel(meterState());
    fireEvent.click(screen.getByTestId('circuit-profile-purpose-wide'));
    expandSources();
    fireEvent.change(screen.getByTestId('circuit-profile-role-tag-3'), { target: { value: 'counter_raw' } });

    expect(screen.getByTestId('circuit-counter-confirmation')).toBeInTheDocument();
    const repair = screen.getByRole('button', { name: 'next.actions.resolve_semantics' });
    fireEvent.click(repair); // Pure repair navigation, never inferred counter confirmation.
    expect(screen.queryByTestId('circuit-profile-prepare')).not.toBeInTheDocument();
    expect(startAsync).not.toHaveBeenCalled();
    expect(ensureCircuitAsync).not.toHaveBeenCalled();
  });

  it.each(['long', 'wide'] as const)('keeps the %s bounds: 1 to 300 seconds instantaneous and 1 to 600 minutes cumulative', (purpose) => {
    measurementsQuery.data = [confirmedCounter(3)];
    renderPanel(meterState());
    if (purpose === 'wide') fireEvent.click(screen.getByTestId('circuit-profile-purpose-wide'));
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Utilities' } });
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));

    const bounded = (testId: string, errorId: string, accepted: string[], refused: string[]) => {
      for (const value of accepted) {
        fireEvent.change(screen.getByTestId(testId), { target: { value } });
        expect(screen.queryByTestId(errorId)).toBeNull();
      }
      for (const value of refused) {
        fireEvent.change(screen.getByTestId(testId), { target: { value } });
        expect(screen.getByTestId(errorId)).toBeInTheDocument();
      }
    };
    bounded('circuit-profile-interval-instantaneous', 'circuit-profile-interval-instantaneous-error', ['1', '300'], ['0', '301', '1.5']);
    fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '15' } });
    bounded('circuit-profile-interval-cumulative', 'circuit-profile-interval-cumulative-error', ['1', '600'], ['0', '601', '10.5']);
  });

  it('keeps the wide profile free of a storage mode and a recreate policy in what it sends', async () => {
    renderPanel(savedGroupState());
    fireEvent.click(screen.getByTestId('circuit-profile-purpose-wide'));
    expect(screen.queryByTestId('circuit-profile-storage-mode')).toBeNull();
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensureCircuitAsync).toHaveBeenCalledTimes(1));

    const [group] = sentGroups();
    expect(group.layout_version).toBe('managed_wide_profile_v1');
    expect(group.storage_mode).toBeUndefined();
    expect(group.write_policy.recreate_missing).toBeUndefined();
    expect(group.destination.storage_strategy).toBe('managed');
  });
});

describe('the batch tools stay on long circuits', () => {
  /** One circuit's group on the given layout over two sources of the meter. */
  function circuitGroup(layout: 'managed_long_v1' | 'managed_wide_profile_v1', id: string, circuit: string, points: [number, number]): WriteGroup {
    const wide = layout === 'managed_wide_profile_v1';
    return savedProfile(layout, {
      id, revision: `rev-${id}-1`, name: `${circuit} (Instantaneous)`,
      members: points.map((i) => ({
        device_id: 'dev-1', point_id: `pt-${i}`, tag_id: `tag-${i}`, source_revision: `s${i}`, mapping_revision: `m${i}`, target_column: '', required: false,
        circuit_id: circuit, circuit_name: circuit, ...(wide ? { value_role: 'gauge' as const, source_value_type: 'float64' as const } : {}),
      })),
    });
  }

  it('does not present the batch tools while a saved wide circuit is in view', () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [
      circuitGroup('managed_wide_profile_v1', 'grp-w', 'cir-w', [4, 5]),
      circuitGroup('managed_long_v1', 'grp-a', 'cir-a', [0, 1]),
      circuitGroup('managed_long_v1', 'grp-b', 'cir-b', [2, 3]),
    ] };
    renderPanel(meterState());

    expect(screen.getByTestId('circuit-profile-circuit-chooser')).toHaveValue('cir-w');
    expect(screen.queryByTestId('circuit-profile-open-batch')).toBeNull();
  });

  it('lists only long circuits in a batch and never writes the wide one', async () => {
    updateGroupAsync.mockImplementation(async ({ id, request }) => ({
      workspace_revision: 'wrev-2', group: circuitGroup('managed_long_v1', id, 'cir-a', [0, 1]), request,
    }));
    groupsQuery.data = { ...groupsQuery.data, groups: [
      circuitGroup('managed_long_v1', 'grp-a', 'cir-a', [0, 1]),
      circuitGroup('managed_long_v1', 'grp-b', 'cir-b', [2, 3]),
      circuitGroup('managed_wide_profile_v1', 'grp-w', 'cir-w', [4, 5]),
    ] };
    renderPanel(meterState());

    expect(screen.getByTestId('circuit-profile-circuit-chooser')).toHaveValue('cir-a');
    fireEvent.click(screen.getByTestId('circuit-profile-open-batch'));
    expect(screen.getByTestId('circuit-profile-batch-select-cir-a')).toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-batch-select-cir-b')).toBeInTheDocument();
    expect(screen.queryByTestId('circuit-profile-batch-select-cir-w')).toBeNull();

    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));
    await waitFor(() => expect(updateGroupAsync).toHaveBeenCalled());
    expect(updateGroupAsync.mock.calls.map((call) => (call[0] as { id: string }).id)).not.toContain('grp-w');
    expect(sentGroups().every((group) => group.layout_version === 'managed_long_v1')).toBe(true);
    expect(startAsync).not.toHaveBeenCalled();
  });
});
