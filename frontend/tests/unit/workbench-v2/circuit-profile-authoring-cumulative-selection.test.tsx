import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { createInstance } from 'i18next';
import { I18nextProvider } from 'react-i18next';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import { defaultCircuitProfileLocalDraft, loadCircuitProfileLocalDraft, saveCircuitProfileLocalDraft } from '@/features/datalink/workbench-v2/state/circuitProfileBasic';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import type { MeasurementDefinition } from '@/types/measurement';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import productionI18n from '@/i18n/config';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
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

vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
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

async function show(state = mixedState(), readonly = false, lng = 'en') {
  const i18n = createInstance();
  await i18n.init({ lng, resources: productionI18n.options.resources });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const panel = () => <I18nextProvider i18n={i18n}><QueryClientProvider client={client}>
    <BasicRecordingPanel state={state} workspaceId="ws-1" readonly={readonly} />
  </QueryClientProvider></I18nextProvider>;
  const view = render(panel());
  return { ...view, refresh: () => view.rerender(panel()) };
}

function setCumulative(tagId: string, checked: boolean) {
  const checkbox = screen.getByTestId(`circuit-profile-member-cum-${tagId}`) as HTMLInputElement;
  if (checkbox.checked !== checked) fireEvent.click(checkbox);
}

function expectNoWrites() {
  for (const mutation of [ensureProfile, updateGroup, start, ensureManaged, applySchema, createMeasurement, updateMeasurement]) {
    expect(mutation).not.toHaveBeenCalled();
  }
}

beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] };
  groupsQuery.refetch.mockImplementation(async () => ({ data: groupsQuery.data }));
  measurementsQuery.data = [];
  measurementsQuery.isLoading = false;
  measurementsQuery.isError = false;
  measurementsQuery.refetch.mockImplementation(async () => ({ data: measurementsQuery.data }));
  ensureProfile.mockRejectedValue(new Error('Request capture only; no backend is running.'));
  previewSchema.mockResolvedValue({ token: 'preview-1', operation_id: 'schema-1', tables: [], statements: [] });
});

function mixedState() {
  const state = savedGroupState();
  state.points.push({ ...state.points[0], id: 'rule-1-p-3', name: 'energy.kWh', address: '40004', data_type: 'uint64' });
  state.mappings['rule-1-p-3'] = {
    ...state.mappings['rule-1-p-0'], point_id: 'rule-1-p-3', persisted_point_id: 'pt-3', tag_id: 'tag-3',
    tag_key: 'energy.kWh', display_name: 'energy.kWh', target_type: 'uint64', unit: 'kWh',
  };
  return state;
}

function counterDefinition(): MeasurementDefinition {
  return {
    id: 'meas-kwh', workspace_id: 'ws-1', definition_revision: 'mrev-1', source_binding_revision: 's3',
    series_epoch: 'epoch-1', name: 'Energy', equipment_id: 'cir-saved', device_id: 'dev-1',
    point_id: 'pt-3', tag_id: 'tag-3', quantity: 'energy', unit: 'kWh', semantic_kind: 'counter', counter_policy: { allow_negative: false },
  };
}

function savedCumulative(): WriteGroup {
  return savedGroup({
    id: 'cum-saved', layout_version: 'managed_long_v1', storage_mode: 'plain',
    recording_profile: { kind: 'cumulative', interval: 10, unit: 'minute' },
    row_policy: { interval_seconds: 600, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    write_policy: { mode: 'append', dedupe_capability: 'receipt', recreate_missing: true },
    destination: { ...savedGroup().destination, storage_strategy: 'managed' },
    members: [{ device_id: 'dev-1', point_id: 'pt-3', tag_id: 'tag-3', source_revision: 's3', mapping_revision: 'm3',
      circuit_id: 'cir-saved', circuit_name: 'Energy A', measurement_id: 'meas-kwh', target_column: '', required: false }],
  });
}

const candidates = buildGroupCandidates(mixedState()).candidates;

describe('explicit cumulative source selection', () => {
  it('starts a new cumulative draft empty even with a confirmed kWh counter among mixed candidates', async () => {
    measurementsQuery.data = [counterDefinition()];
    await show();
    expect(screen.getByTestId('circuit-profile-sources-toggle')).toHaveAttribute('aria-expanded', 'false');
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    expect(screen.getByTestId('circuit-profile-sources-toggle')).toHaveAttribute('aria-expanded', 'true');
    expect(screen.getByTestId('circuit-profile-cumulative-selected')).toHaveTextContent('Selected cumulative sources: 0');
    expect(screen.getByTestId('circuit-profile-cumulative-sources-error')).toHaveTextContent('Choose the counters to record from the source list, then confirm their meaning.');
    for (const tag of ['tag-0', 'tag-1', 'tag-2', 'tag-3']) {
      expect(screen.getByTestId(`circuit-profile-member-cum-${tag}`)).not.toBeChecked();
      expect(screen.getByTestId(`circuit-profile-member-inst-${tag}`)).toBeChecked();
    }
    expect(defaultCircuitProfileLocalDraft(candidates, []).cumulative_keys).toEqual([]);
    expectNoWrites();
  });

  it('does not interpret legacy selected_keys as cumulative intent and keeps the instantaneous selection', async () => {
    window.sessionStorage.setItem('wbv2.circuit-profile-draft.ws-1%2Fdev-1', JSON.stringify({
      circuit_id: 'cir-legacy', circuit_name: 'Legacy', selected_keys: ['dev-1\0pt-0\0tag-0'], cumulative_enabled: true,
    }));
    const loaded = loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, []);
    expect(loaded.cumulative_keys).toEqual([]);
    expect(loaded.instantaneous_keys).toEqual(['dev-1\0pt-0\0tag-0']);
    await show();
    expect(screen.getByTestId('circuit-profile-cumulative-selected')).toHaveTextContent('Selected cumulative sources: 0');
    expect(screen.getByTestId('circuit-profile-cumulative-sources-error')).toBeVisible();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    expectNoWrites();
  });

  it('authors only an explicitly selected confirmed kWh counter while retaining instantaneous choices', async () => {
    measurementsQuery.data = [counterDefinition()];
    await show();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Energy A' } });
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    setCumulative('tag-3', true);
    expect(screen.getByTestId('circuit-profile-cumulative-selected')).toHaveTextContent('Selected cumulative sources: 1');
    for (const tag of ['tag-0', 'tag-1', 'tag-2', 'tag-3']) expect(screen.getByTestId(`circuit-profile-member-inst-${tag}`)).toBeChecked();
    expect(screen.queryByTestId('circuit-counter-confirmation')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled();
    expectNoWrites();
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-instantaneous'));
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensureProfile).toHaveBeenCalledTimes(1));
    const group = ensureProfile.mock.calls[0][0].request.group;
    expect(group.recording_profile).toEqual({ kind: 'cumulative', interval: 10, unit: 'minute' });
    expect(group.members).toHaveLength(1);
    expect(group.members[0]).toMatchObject({ device_id: 'dev-1', point_id: 'pt-3', tag_id: 'tag-3', measurement_id: 'meas-kwh' });
    // The capture rejects rather than pretending a backend saved or started the draft.
    expect(await screen.findByTestId('circuit-profile-error')).toBeVisible();
    expect(start).not.toHaveBeenCalled();
    expect(applySchema).not.toHaveBeenCalled();
  });

  it('preserves saved kWh membership through off/on and reload', async () => {
    const saved = savedCumulative();
    groupsQuery.data.groups = [saved];
    measurementsQuery.data = [counterDefinition()];
    const view = await show();
    expandSources();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-3')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-0')).not.toBeChecked();
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    expect(screen.getByTestId('circuit-profile-member-cum-tag-3')).toBeChecked();
    view.unmount();
    await show();
    expandSources();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-3')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-0')).not.toBeChecked();
    expect(saved.members.map((member) => member.tag_id)).toEqual(['tag-3']);
    expectNoWrites();
  });

  it('retains an intentional empty selection instead of reselecting a saved counter', async () => {
    groupsQuery.data.groups = [savedCumulative()];
    measurementsQuery.data = [counterDefinition()];
    const view = await show();
    expandSources();
    setCumulative('tag-3', false);
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    expect(screen.getByTestId('circuit-profile-member-cum-tag-3')).not.toBeChecked();
    view.unmount();
    await show();
    expandSources();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-3')).not.toBeChecked();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    expectNoWrites();
  });

  it('keeps an explicitly stored draft selection while cumulative is disabled', async () => {
    const draft = defaultCircuitProfileLocalDraft(candidates, []);
    saveCircuitProfileLocalDraft('ws-1', 'dev-1', { ...draft, cumulative_keys: ['dev-1\0pt-3\0tag-3'], cumulative_enabled: false });
    await show();
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    expect(screen.getByTestId('circuit-profile-member-cum-tag-3')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-0')).not.toBeChecked();
    expectNoWrites();
  });

  it('uses server-saved cumulative membership when a legacy cache has no dedicated selection array', () => {
    window.sessionStorage.setItem('wbv2.circuit-profile-draft.ws-1%2Fdev-1', JSON.stringify({
      circuit_id: 'cir-saved', selected_keys: ['dev-1\0pt-0\0tag-0'], cumulative_enabled: true,
    }));
    const loaded = loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, [savedCumulative()], 'cir-saved');
    expect(loaded.cumulative_keys).toEqual(['dev-1\0pt-3\0tag-3']);
    expect(loaded.instantaneous_keys).toEqual(['dev-1\0pt-0\0tag-0']);
  });

  it('reopens the picker from its missing-source guidance without changing selections', async () => {
    await show();
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    fireEvent.click(screen.getByTestId('circuit-profile-sources-toggle'));
    expect(screen.getByTestId('circuit-profile-sources-toggle')).toHaveAttribute('aria-expanded', 'false');
    fireEvent.click(screen.getByTestId('circuit-profile-cumulative-choose-sources'));
    expect(screen.getByTestId('circuit-profile-sources-toggle')).toHaveAttribute('aria-expanded', 'true');
    expect(screen.getByTestId('circuit-profile-sources-toggle')).toHaveFocus();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-3')).not.toBeChecked();
    expectNoWrites();
  });

  it('does not add a source to either disabled profile through the master checkbox', async () => {
    await show();
    expandSources();
    fireEvent.click(screen.getByRole('checkbox', { name: 'line.temperature' }));
    fireEvent.click(screen.getByRole('checkbox', { name: 'line.temperature' }));
    expect(screen.getByTestId('circuit-profile-member-inst-tag-0')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-0')).not.toBeChecked();
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-instantaneous'));
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    fireEvent.click(screen.getByRole('checkbox', { name: 'line.temperature' }));
    fireEvent.click(screen.getByRole('checkbox', { name: 'line.temperature' }));
    expect(screen.getByTestId('circuit-profile-member-inst-tag-0')).not.toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-0')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-inst-tag-1')).toBeChecked();
    expectNoWrites();
  });

  it('isolates saved device A membership from draft choices on device B and back', async () => {
    const state = mixedState();
    state.devices.push({ ...state.devices[0], id: 'dev-2', name: 'PLC B' });
    state.points.push({ ...state.points[3], id: 'point-b', device_id: 'dev-2' });
    state.mappings['point-b'] = { ...state.mappings['rule-1-p-3'], point_id: 'point-b', persisted_point_id: 'pt-b', tag_id: 'tag-b', device_id: 'dev-2' };
    groupsQuery.data.groups = [savedCumulative()];
    measurementsQuery.data = [counterDefinition(), { ...counterDefinition(), id: 'meas-b', device_id: 'dev-2', point_id: 'pt-b', tag_id: 'tag-b' }];
    await show(state);
    expandSources();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-3')).toBeChecked();
    fireEvent.click(screen.getByRole('button', { name: 'PLC B' }));
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Energy B' } });
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    expect(screen.getByTestId('circuit-profile-member-cum-tag-b')).not.toBeChecked();
    setCumulative('tag-b', true);
    fireEvent.click(screen.getByRole('button', { name: 'PLC A' }));
    expandSources();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-3')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-0')).not.toBeChecked();
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Energy A');
    fireEvent.click(screen.getByRole('button', { name: 'PLC B' }));
    expandSources();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-b')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Energy B');
    expectNoWrites();
  });

  it.each(['tag-0', 'tag-2'])('still requires counter compatibility after an explicit selection of %s', async (tag) => {
    await show();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Explicit selection' } });
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    setCumulative(tag, true);
    expect(screen.getByTestId('circuit-profile-cumulative-selected')).toHaveTextContent('Selected cumulative sources: 1');
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-start')).toBeDisabled();
    expect(screen.getByTestId('circuit-counter-confirmation')).toBeVisible();
    expectNoWrites();
  });

  it('shows translated zero-selection guidance in read-only mode while allowing inspection', async () => {
    saveCircuitProfileLocalDraft('ws-1', 'dev-1', { ...defaultCircuitProfileLocalDraft(candidates, []), cumulative_enabled: true });
    await show(mixedState(), true, 'zh-TW');
    expect(screen.getByTestId('circuit-profile-toggle-cumulative')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-cumulative-selected')).toHaveTextContent('已選累積來源：0');
    expect(screen.getByTestId('circuit-profile-cumulative-sources-error')).toHaveTextContent('請從來源清單選擇要記錄的計數器，再確認其計數語意。');
    fireEvent.click(screen.getByTestId('circuit-profile-cumulative-choose-sources'));
    expect(screen.getByTestId('circuit-profile-member-cum-tag-3')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-3')).not.toBeChecked();
    expectNoWrites();
  });

  it.each(['not-json', 'null', '[]', '{"cumulative_keys":"invalid"}'])('falls back safely without selecting counters for invalid cache %s', (raw) => {
    window.sessionStorage.setItem('wbv2.circuit-profile-draft.ws-1%2Fdev-1', raw);
    expect(loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, []).cumulative_keys).toEqual([]);
  });

  it('retains only valid explicit cumulative source keys from the current device', () => {
    window.sessionStorage.setItem('wbv2.circuit-profile-draft.ws-1%2Fdev-1', JSON.stringify({
      cumulative_keys: ['dev-1\0pt-3\0tag-3', 'dev-2\0pt-b\0tag-b', 5],
    }));
    expect(loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, []).cumulative_keys).toEqual(['dev-1\0pt-3\0tag-3']);
  });
});
