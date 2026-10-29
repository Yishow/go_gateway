import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { guidedReadiness, guidedLifecycle } from '../../fixtures/guidedRecordingPreparation';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { createInstance } from 'i18next';
import { I18nextProvider } from 'react-i18next';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import type { Mapping, WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import type { MeasurementDefinition } from '@/types/measurement';
import type { WriteGroup, WriteGroupCreateRequest, WriteGroupValueRole } from '@/types/studioV2WriteGroup';
import productionI18n from '@/i18n/config';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import { longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { expandSources } from '../../fixtures/sourceReview';

const groupsQuery = {
  data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] as WriteGroup[] },
  isLoading: false, isError: false, refetch: vi.fn(),
};
const measurementsQuery = {
  data: [] as MeasurementDefinition[], isLoading: false, isError: false, refetch: vi.fn(),
};
const ensureProfile = vi.fn();
const updateGroup = vi.fn();
const start = vi.fn();
const ensureManaged = vi.fn();
const applySchema = vi.fn();
const previewSchema = vi.fn();
const createMeasurement = vi.fn();
const updateMeasurement = vi.fn();

let canonicalClient: QueryClient;
const evidenceRefetch = vi.fn();
const evidence = new Map<string, { readiness: ReturnType<typeof guidedReadiness>; lifecycle: ReturnType<typeof guidedLifecycle> }>();
function currentEvidence(id: string) {
  const group = groupsQuery.data.groups.find(value => value.id === id);
  if (!group) return undefined;
  const key = JSON.stringify([groupsQuery.data.workspace_revision, group]);
  if (!evidence.has(key)) evidence.set(key, { readiness: guidedReadiness(group, groupsQuery.data.workspace_revision), lifecycle: guidedLifecycle(group, groupsQuery.data.workspace_revision) });
  return evidence.get(key)!;
}
function publishGroups() {
  canonicalClient?.setQueryData(studioV2WorkspaceKeys.writeGroups(), groupsQuery.data);
  return { data: groupsQuery.data };
}
vi.mock('@/hooks/datalink/useStudioV2WriteGroupLifecycle', () => ({
  useWriteGroupLifecycleQuery: (_workspace: string, id: string) => ({ data: currentEvidence(id)?.lifecycle, isLoading: false, isFetching: false, isError: false, refetch: evidenceRefetch }),
}));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupReadinessQuery: (id: string) => ({ data: currentEvidence(id)?.readiness, isLoading: false, isFetching: false, isError: false, refetch: evidenceRefetch }),
  useWriteGroupsQuery: () => groupsQuery,
  useWriteGroupDeliveryQuery: () => ({ data: undefined, isError: false }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: updateGroup, isPending: false }),
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: ensureProfile, isPending: false }),
  useWriteGroupSchemaPreviewMutation: () => ({ mutateAsync: previewSchema, isPending: false }),
  useWriteGroupSchemaApplyMutation: () => ({ mutateAsync: applySchema, isPending: false }),
  useWriteGroupSchemaOperationQuery: () => ({ data: undefined, isError: false, refetch: vi.fn() }),
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useEnsureBasicManagedMutation: () => ({ mutateAsync: ensureManaged, isPending: false }),
  useRecordingStartMutation: () => ({ mutateAsync: start, isPending: false }),
  useRecordingStartOperationQuery: () => ({ data: undefined, isLoading: false, isError: false, refetch: vi.fn() }),
}));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => measurementsQuery,
  useCreateStudioV2MeasurementMutation: () => ({ mutateAsync: createMeasurement, isPending: false }),
  useUpdateStudioV2MeasurementMutation: () => ({ mutateAsync: updateMeasurement, isPending: false }),
}));

function sourceState(type: Mapping['target_type'] = 'int16'): WorkbenchV2State {
  const state = savedGroupState();
  state.points[0] = { ...state.points[0], name: 'machine_state', data_type: type };
  state.mappings[state.points[0].id] = {
    ...state.mappings[state.points[0].id], target_type: type, tag_key: 'machine_state', display_name: 'machine_state',
  };
  return state;
}

function confirmedCounter(overrides: Partial<MeasurementDefinition> = {}): MeasurementDefinition {
  return {
    id: 'meas-0', workspace_id: 'ws-1', definition_revision: 'mrev-1', source_binding_revision: 's0',
    series_epoch: 'epoch-1', name: 'machine_state', equipment_id: 'cir-saved', device_id: 'dev-1',
    point_id: 'pt-0', tag_id: 'tag-0', quantity: '', unit: 'count', semantic_kind: 'counter',
    counter_policy: { allow_negative: false }, ...overrides,
  };
}

function savedNumericGroup(role: WriteGroupValueRole): WriteGroup {
  const group = longProfileGroup({ id: 'inst-saved', circuitId: 'cir-saved', circuitName: 'Saved', interval: 15, points: [0, 1] });
  return {
    ...group, recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'mean' },
    members: group.members.map((member, index) => ({
      ...member, value_role: index === 0 ? role : 'gauge', measurement_id: index === 0 && role === 'counter_raw' ? 'meas-0' : '',
    })),
  };
}

async function show(state = sourceState(), readonly = false, lng = 'en') {
  const i18n = createInstance();
  await i18n.init({ lng, resources: productionI18n.options.resources });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  canonicalClient = client; client.setQueryDefaults(studioV2WorkspaceKeys.writeGroups(), { structuralSharing: false }); publishGroups();
  const panel = () => <I18nextProvider i18n={i18n}><QueryClientProvider client={client}>
    <BasicRecordingPanel state={state} workspaceId="ws-1" readonly={readonly} />
  </QueryClientProvider></I18nextProvider>;
  const view = render(panel());
  return { ...view, refresh: () => { publishGroups(); view.rerender(panel()); } };
}

function roleSelect(tagId = 'tag-0'): HTMLSelectElement {
  return screen.getByTestId(`circuit-profile-role-${tagId}`) as HTMLSelectElement;
}

function expectNoWrites() {
  for (const mutation of [ensureProfile, updateGroup, start, ensureManaged, applySchema, createMeasurement, updateMeasurement]) {
    expect(mutation).not.toHaveBeenCalled();
  }
}

// Guide design: one primary; semantic/source repair is local NAV, not a second disabled Start.
function expectPreparationBlocked() {
  const mutations = [ensureProfile, updateGroup, start, ensureManaged, applySchema, previewSchema, createMeasurement, updateMeasurement]; const counts = mutations.map(fn => fn.mock.calls.length);
  const prepare = screen.queryByTestId('circuit-profile-prepare');
  if (prepare) expect(prepare).toBeDisabled();
  expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument();
  const next = screen.queryByTestId('recording-next-action');
  if (next) {
    expect(next).toHaveTextContent(/Review counter definitions|Choose sources|Inspect this scope|Refresh preparation status|確認計數器定義|選擇資料來源|檢視此範圍|重新查詢準備狀態/);
    fireEvent.click(next);
  }
  expect(screen.getByTestId('recording-guided-preparation').querySelectorAll('button[data-primary="true"]')).toHaveLength(1);
  mutations.forEach((fn, index) => expect(fn).toHaveBeenCalledTimes(counts[index]));
}

beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] };
  groupsQuery.refetch.mockImplementation(async () => publishGroups());
  measurementsQuery.data = [];
  measurementsQuery.isLoading = false;
  measurementsQuery.isError = false;
  measurementsQuery.refetch.mockImplementation(async () => ({ data: measurementsQuery.data }));
  ensureProfile.mockImplementation(async ({ request }: { request: WriteGroupCreateRequest }) => ({
    workspace_revision: 'wrev-2', group: savedGroup({
      ...request.group, id: 'inst-new', revision: 'rev-new',
      members: request.group.members.map((member) => ({ ...member, source_revision: 's1', mapping_revision: 'm1' })),
      destination: { ...request.group.destination, database: '/recording-test.db' },
    }),
  }));
  previewSchema.mockResolvedValue({ token: 'preview-1', operation_id: 'schema-1', tables: [], statements: [] });
});

describe('numeric state authoring uses the existing latest-value role', () => {
  it.each(['int16', 'float64'] as const)('offers explicit state for %s and authors state, not gauge', async (type) => {
    // The int16 machine_state 1 -> 2 scenario is authored here; actual latest 2 (not 1.5)
    // belongs to the backend acceptance gate, not a fabricated frontend reduction.
    await show(sourceState(type));
    expandSources();
    const select = roleSelect();
    expect(select).toHaveValue('gauge');
    expect(within(select).getByRole('option', { name: 'State (latest value)' })).toHaveValue('state');
    expect(within(select).getByRole('option', { name: 'Measurement (averaged)' })).toHaveValue('gauge');
    expect(within(select).getByRole('option', { name: 'Raw counter (latest value, must be confirmed)' })).toHaveValue('counter_raw');

    fireEvent.change(select, { target: { value: 'state' } });
    expect(select).toHaveValue('state');
    expectNoWrites();
    expect(screen.queryByTestId('circuit-counter-confirmation')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensureProfile).toHaveBeenCalledTimes(1));
    const group = ensureProfile.mock.calls[0][0].request.group;
    expect(group.recording_profile).toEqual({ kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'mean' });
    expect(group.members[0]).toMatchObject({ point_id: 'pt-0', tag_id: 'tag-0', value_role: 'state', measurement_id: '' });
    expect(group.members[0].value_role).not.toBe('gauge');
    expect(updateGroup).not.toHaveBeenCalled();
    expect(start).not.toHaveBeenCalled();
    expect(applySchema).not.toHaveBeenCalled();
  });

  it('uses existing Traditional Chinese labels to distinguish latest, mean and confirmation', async () => {
    await show(sourceState(), false, 'zh-TW');
    expandSources();
    expect(within(roleSelect()).getByRole('option', { name: '狀態（取最新值）' })).toHaveValue('state');
    expect(within(roleSelect()).getByRole('option', { name: '量測值（取平均）' })).toHaveValue('gauge');
    expect(within(roleSelect()).getByRole('option', { name: '計數器原值（取最新值，需已確認）' })).toHaveValue('counter_raw');
    expectNoWrites();
  });

  it('keeps numeric gauge as mean without requiring a counter definition', async () => {
    await show(savedGroupState()); // Numeric temperature has no counter definition.
    expandSources();
    fireEvent.change(roleSelect(), { target: { value: 'gauge' } });
    expect(screen.queryByTestId('circuit-counter-confirmation')).not.toBeInTheDocument();
    expectNoWrites();
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensureProfile).toHaveBeenCalledTimes(1));
    expect(ensureProfile.mock.calls[0][0].request.group.recording_profile.sample_reduction).toBe('mean');
    expect(ensureProfile.mock.calls[0][0].request.group.members[0]).toMatchObject({ value_role: 'gauge', measurement_id: '' });
  });

  it.each(['string', 'bool'] as const)('offers only state for a %s source and preserves that role in the request', async (type) => {
    await show(sourceState(type));
    expandSources();
    expect(roleSelect()).toHaveValue('state');
    expect(within(roleSelect()).getAllByRole('option').map((option) => (option as HTMLOptionElement).value)).toEqual(['state']);
    expectNoWrites();
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensureProfile).toHaveBeenCalledTimes(1));
    expect(ensureProfile.mock.calls[0][0].request.group.members[0]).toMatchObject({ value_role: 'state', measurement_id: '' });
  });

  it.each(['missing', 'unconfirmed', 'foreign', 'gauge', 'loading', 'error'] as const)('blocks counter_raw when its definition is %s', async (condition) => {
    if (condition !== 'missing') measurementsQuery.data = [confirmedCounter(
      condition === 'unconfirmed' ? { source_binding_revision: '' }
        : condition === 'foreign' ? { workspace_id: 'ws-other' }
          : condition === 'gauge' ? { semantic_kind: 'gauge' } : {},
    )];
    measurementsQuery.isLoading = condition === 'loading';
    measurementsQuery.isError = condition === 'error';
    await show();
    expandSources();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Machine' } });
    fireEvent.change(roleSelect(), { target: { value: 'counter_raw' } });
    expect(roleSelect()).toHaveValue('counter_raw');
    expectPreparationBlocked();
    expectNoWrites();
    fireEvent.change(roleSelect(), { target: { value: 'state' } });
    expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled();
    expect(screen.queryByTestId('circuit-counter-confirmation')).not.toBeInTheDocument();
    expectNoWrites();
  });

  it('prepares counter_raw only after explicit confirmation returns a saved definition and binding', async () => {
    createMeasurement.mockImplementation(async () => {
      measurementsQuery.data = [confirmedCounter()];
      return measurementsQuery.data[0];
    });
    const view = await show();
    expandSources();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Machine' } });
    fireEvent.change(roleSelect(), { target: { value: 'counter_raw' } });
    expectPreparationBlocked();
    fireEvent.change(screen.getByTestId('circuit-confirm-counter-unit'), { target: { value: 'count' } });
    expectNoWrites();
    fireEvent.click(screen.getByTestId('circuit-confirm-counter-submit'));
    await waitFor(() => expect(measurementsQuery.refetch).toHaveBeenCalledTimes(1));
    view.refresh(); // The mocked query publishes the refetched canonical definition.
    await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
    expect(createMeasurement).toHaveBeenCalledTimes(1);
    expect(createMeasurement.mock.calls[0][0]).toMatchObject({ device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', semantic_kind: 'counter', unit: 'count' });
    expect(ensureProfile).not.toHaveBeenCalled();
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensureProfile).toHaveBeenCalledTimes(1));
    expect(ensureProfile.mock.calls[0][0].request.group.members[0]).toMatchObject({ value_role: 'counter_raw', measurement_id: 'meas-0' });
    expect(start).not.toHaveBeenCalled();
    expect(applySchema).not.toHaveBeenCalled();
  });

  it.each(['state', 'gauge', 'counter_raw'] as const)('displays saved numeric %s locked across stale cache, circuit/device navigation and reload', async (role) => {
    const saved = savedNumericGroup(role);
    const original = structuredClone(saved);
    groupsQuery.data.groups = [saved];
    measurementsQuery.data = [confirmedCounter()];
    const state = sourceState();
    state.devices.push({ ...state.devices[0], id: 'dev-2', name: 'PLC B' });
    window.sessionStorage.setItem('wbv2.circuit-profile-draft.ws-1%2Fdev-1%2Fcir-saved', JSON.stringify({
      circuit_id: 'cir-saved', member_roles: { 'dev-1\0pt-0\0tag-0': role === 'state' ? 'gauge' : 'state' },
    }));
    const view = await show(state);
    expandSources();
    expect(roleSelect()).toHaveValue(role);
    expect(roleSelect()).toBeDisabled();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-chooser'), { target: { value: '__new__' } });
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-chooser'), { target: { value: 'cir-saved' } });
    fireEvent.click(screen.getByRole('button', { name: 'PLC B' }));
    fireEvent.click(screen.getByRole('button', { name: 'PLC A' }));
    expandSources();
    expect(roleSelect()).toHaveValue(role);
    expect(roleSelect()).toBeDisabled();
    view.unmount();
    await show(state);
    expandSources();
    expect(roleSelect()).toHaveValue(role);
    expect(roleSelect()).toBeDisabled();
    expect(saved).toEqual(original);
    expectNoWrites();
  });

  it('retains a new explicit state through device navigation and reload without writing', async () => {
    const state = sourceState();
    state.devices.push({ ...state.devices[0], id: 'dev-2', name: 'PLC B' });
    const view = await show(state);
    expandSources();
    fireEvent.change(roleSelect(), { target: { value: 'state' } });
    fireEvent.click(screen.getByRole('button', { name: 'PLC B' }));
    fireEvent.click(screen.getByRole('button', { name: 'PLC A' }));
    expandSources();
    expect(roleSelect()).toHaveValue('state');
    view.unmount();
    await show(state);
    expandSources();
    expect(roleSelect()).toHaveValue('state');
    expect(roleSelect()).toBeEnabled();
    expectNoWrites();
  });

  it('keeps the role control locked in read-only mode', async () => {
    await show(sourceState(), true);
    expandSources();
    expect(roleSelect()).toBeDisabled();
    expectPreparationBlocked();
    expectNoWrites();
  });

  it('retains an explicit role when its instantaneous source is deselected and selected again', async () => {
    await show();
    expandSources();
    fireEvent.change(roleSelect(), { target: { value: 'state' } });
    fireEvent.click(screen.getByTestId('circuit-profile-member-inst-tag-0'));
    expect(screen.queryByTestId('circuit-profile-role-tag-0')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('circuit-profile-member-inst-tag-0'));
    expect(roleSelect()).toHaveValue('state');
    expectNoWrites();
  });

  it('keeps counter preparation blocked if explicit confirmation fails', async () => {
    createMeasurement.mockRejectedValueOnce(new Error('private-endpoint-secret'));
    await show();
    expandSources();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Machine' } });
    fireEvent.change(roleSelect(), { target: { value: 'counter_raw' } });
    fireEvent.change(screen.getByTestId('circuit-confirm-counter-unit'), { target: { value: 'count' } });
    fireEvent.click(screen.getByTestId('circuit-confirm-counter-submit'));
    expect(await screen.findByTestId('circuit-confirm-counter-error')).not.toHaveTextContent('private-endpoint-secret');
    expect(roleSelect()).toHaveValue('counter_raw');
    expectPreparationBlocked();
    expect(ensureProfile).not.toHaveBeenCalled();
    expect(start).not.toHaveBeenCalled();
    expect(applySchema).not.toHaveBeenCalled();
  });

  it.each(['latest', undefined] as const)('does not introduce role editing for saved reduction %s', async (reduction) => {
    const saved = savedNumericGroup('state');
    saved.recording_profile = { kind: 'instantaneous', interval: 15, unit: 'second', ...(reduction ? { sample_reduction: reduction } : {}) };
    groupsQuery.data.groups = [saved];
    await show();
    expandSources();
    expect(screen.queryByTestId('circuit-profile-role-tag-0')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-reduction-instantaneous')).toHaveTextContent('the latest value in the window, never changed to a mean');
    expectNoWrites();
  });
});
