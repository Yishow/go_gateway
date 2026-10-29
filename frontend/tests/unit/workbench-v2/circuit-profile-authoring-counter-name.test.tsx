import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { createInstance } from 'i18next';
import { I18nextProvider } from 'react-i18next';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import type { CircuitProfileLocalDraft } from '@/features/datalink/workbench-v2/state/circuitProfileBasic';
import type { MeasurementDefinition } from '@/types/measurement';
import type { WriteGroup, WriteGroupCreateRequest } from '@/types/studioV2WriteGroup';
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

async function show(state = savedGroupState(), readonly = false, lng = 'en') {
  const i18n = createInstance();
  await i18n.init({ lng, resources: productionI18n.options.resources });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const panel = () => <I18nextProvider i18n={i18n}><QueryClientProvider client={client}>
    <BasicRecordingPanel state={state} workspaceId="ws-1" readonly={readonly} />
  </QueryClientProvider></I18nextProvider>;
  const view = render(panel());
  return { ...view, refresh: () => view.rerender(panel()) };
}

function roleSelect(tagId = 'tag-0'): HTMLSelectElement {
  return screen.getByTestId(`circuit-profile-role-${tagId}`) as HTMLSelectElement;
}

function expectNoWrites() {
  for (const mutation of [ensureProfile, updateGroup, start, ensureManaged, applySchema, createMeasurement, updateMeasurement]) {
    expect(mutation).not.toHaveBeenCalled();
  }
}

function seedCounterDraft(overrides: Partial<CircuitProfileLocalDraft> = {}) {
  window.sessionStorage.setItem('wbv2.circuit-profile-draft.ws-1%2Fdev-1', JSON.stringify({
    member_roles: { 'dev-1\0pt-0\0tag-0': 'counter_raw' }, ...overrides,
  }));
}

function counterDefinition(unit = 'count'): MeasurementDefinition {
  return {
    id: 'meas-0', workspace_id: 'ws-1', definition_revision: 'mrev-1', source_binding_revision: 's0',
    series_epoch: 'epoch-1', name: 'Temperature', equipment_id: 'cir-draft-1', device_id: 'dev-1',
    point_id: 'pt-0', tag_id: 'tag-0', quantity: '', unit, semantic_kind: 'counter', counter_policy: { allow_negative: false },
  };
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
  ensureProfile.mockImplementation(async ({ request }: { request: WriteGroupCreateRequest }) => ({
    workspace_revision: 'wrev-2', group: savedGroup({
      ...request.group, id: 'inst-new', revision: 'rev-new',
      members: request.group.members.map((member) => ({ ...member, source_revision: 's1', mapping_revision: 'm1' })),
      destination: { ...request.group.destination, database: '/recording-test.db' },
    }),
  }));
  previewSchema.mockResolvedValue({ token: 'preview-1', operation_id: 'schema-1', tables: [], statements: [] });
});

describe('counter-dependent circuit names', () => {
  it.each([
    ['en', 'Circuit name', 'Enter a circuit name before confirming these counter sources.'],
    ['zh-TW', '迴路名稱', '確認這些計數來源前，請先輸入迴路名稱。'],
  ])('explains the missing name for instantaneous-only counter confirmation (%s)', async (lng, label, reason) => {
    await show(savedGroupState(), false, lng);
    expandSources();
    fireEvent.change(roleSelect(), { target: { value: 'counter_raw' } });
    expect(screen.getByTestId('circuit-profile-toggle-cumulative')).not.toBeChecked();
    const name = screen.getByTestId('circuit-profile-circuit-name');
    expect(name).toHaveAccessibleName(label);
    const hint = screen.getByTestId('circuit-profile-circuit-required-error');
    expect(hint).toHaveTextContent(reason);
    expect(name).toHaveAttribute('aria-required', 'true');
    expect(name).toHaveAttribute('aria-invalid', 'true');
    expect(name).toHaveAttribute('aria-describedby', hint.id);
    const confirm = screen.getByTestId('circuit-confirm-counter-submit');
    expect(confirm).toBeDisabled();
    expect(confirm).toHaveAttribute('aria-describedby', hint.id);
    expect(confirm).toHaveAccessibleDescription(reason);
    fireEvent.click(screen.getByTestId('circuit-confirm-counter-name-focus'));
    expect(name).toHaveFocus();
    expectNoWrites();
  });

  it.each(['', '   '])('requires a meaningful new-counter name even when a draft UUID exists (%j)', async (name) => {
    seedCounterDraft({ circuit_id: 'cir-draft-1', circuit_name: name });
    await show();
    expandSources();
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue(name);
    expect(screen.getByTestId('circuit-profile-circuit-required-error')).toBeVisible();
    fireEvent.change(screen.getByTestId('circuit-confirm-counter-unit'), { target: { value: 'count' } });
    expect(screen.getByTestId('circuit-confirm-counter-submit')).toBeDisabled();
    fireEvent.click(screen.getByTestId('circuit-confirm-counter-submit'));
    expectNoWrites();
  });

  it('treats name entry as only one prerequisite and preserves stable circuit/source IDs through rename and clearing', async () => {
    seedCounterDraft({ circuit_id: 'cir-draft-1', circuit_name: '' });
    createMeasurement.mockImplementation(async () => {
      measurementsQuery.data = [counterDefinition()];
      return measurementsQuery.data[0];
    });
    const view = await show();
    const name = screen.getByTestId('circuit-profile-circuit-name');
    fireEvent.change(name, { target: { value: 'Machine A' } });
    expect(screen.queryByTestId('circuit-profile-circuit-required-error')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-confirm-counter-submit')).toBeDisabled(); // A unit is still missing.
    fireEvent.change(screen.getByTestId('circuit-confirm-counter-unit'), { target: { value: 'count' } });
    expect(screen.getByTestId('circuit-confirm-counter-submit')).toBeEnabled();
    fireEvent.change(name, { target: { value: '   ' } });
    expect(screen.getByTestId('circuit-profile-circuit-required-error')).toBeVisible();
    expect(screen.getByTestId('circuit-confirm-counter-submit')).toBeDisabled();
    fireEvent.change(name, { target: { value: '  Renamed Circuit  ' } });
    expect(name).toHaveValue('  Renamed Circuit  ');
    expect(screen.getByTestId('circuit-profile-circuit-chooser')).toHaveValue('cir-draft-1');
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-start')).toBeDisabled();
    expectNoWrites();
    fireEvent.click(screen.getByTestId('circuit-confirm-counter-submit'));
    await waitFor(() => expect(measurementsQuery.refetch).toHaveBeenCalledTimes(1));
    expect(createMeasurement.mock.calls[0][0]).toMatchObject({
      equipment_id: 'cir-draft-1', device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', semantic_kind: 'counter',
    });
    view.refresh();
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensureProfile).toHaveBeenCalledTimes(1));
    expect(ensureProfile.mock.calls[0][0].request.group.members[0]).toMatchObject({
      circuit_id: 'cir-draft-1', circuit_name: 'Renamed Circuit', device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0',
    });
    expect(start).not.toHaveBeenCalled();
    expect(applySchema).not.toHaveBeenCalled();
  });

  it.each(['gauge', 'state'])('keeps a non-counter %s name optional and prepares without a counter definition', async (role) => {
    await show();
    expandSources();
    fireEvent.change(roleSelect(), { target: { value: role } });
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveAttribute('aria-required', 'false');
    expect(screen.queryByTestId('circuit-profile-circuit-required-error')).not.toBeInTheDocument();
    expect(screen.queryByTestId('circuit-counter-confirmation')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled();
    expectNoWrites();
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensureProfile).toHaveBeenCalledTimes(1));
    expect(ensureProfile.mock.calls[0][0].request.group.members[0]).toMatchObject({ value_role: role, measurement_id: '' });
    expect(ensureProfile.mock.calls[0][0].request.group.members[0]).not.toHaveProperty('circuit_id');
    expect(createMeasurement).not.toHaveBeenCalled();
  });

  it('does not infer counter semantics from an energy-like circuit name', async () => {
    await show();
    expandSources();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Energy total kWh' } });
    expect(roleSelect()).toHaveValue('gauge');
    expect(screen.queryByTestId('circuit-counter-confirmation')).not.toBeInTheDocument();
    expect(screen.queryByTestId('circuit-profile-circuit-required-error')).not.toBeInTheDocument();
    expectNoWrites();
  });

  it('preserves the required-name prerequisite for selected cumulative counter sources', async () => {
    await show();
    expandSources();
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    for (const [tag, checked] of [['tag-0', true], ['tag-1', false], ['tag-2', false]] as const) {
      const input = screen.getByTestId(`circuit-profile-member-cum-${tag}`) as HTMLInputElement;
      if (input.checked !== checked) fireEvent.click(input);
    }
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveAttribute('aria-required', 'true');
    expect(screen.getByTestId('circuit-profile-circuit-required-error')).toHaveTextContent('Enter a circuit name before confirming these counter sources.');
    expect(screen.getByTestId('circuit-confirm-counter-submit')).toBeDisabled();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Cumulative A' } });
    expect(screen.queryByTestId('circuit-profile-circuit-required-error')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-confirm-counter-submit')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    expectNoWrites();
  });

  it.each(['confirmed', 'missing-unit'])('does not add a name requirement to a saved %s definition with valid identity', async (state) => {
    measurementsQuery.data = [counterDefinition(state === 'confirmed' ? 'count' : '')];
    seedCounterDraft({ circuit_name: '' });
    await show();
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveAttribute('aria-required', 'false');
    expect(screen.queryByTestId('circuit-profile-circuit-required-error')).not.toBeInTheDocument();
    if (state === 'confirmed') {
      expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled();
      expect(screen.queryByTestId('circuit-counter-confirmation')).not.toBeInTheDocument();
    } else {
      fireEvent.change(screen.getByTestId('circuit-confirm-counter-unit'), { target: { value: 'count' } });
      expect(screen.getByTestId('circuit-confirm-counter-submit')).toBeEnabled();
      expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    }
    expectNoWrites();
  });

  it('keeps read-only missing-name guidance visible without enabling confirmation', async () => {
    seedCounterDraft();
    await show(savedGroupState(), true);
    expect(screen.getByTestId('circuit-profile-circuit-required-error')).toBeVisible();
    expect(screen.getByTestId('circuit-profile-circuit-name')).toBeDisabled();
    expect(screen.getByTestId('circuit-confirm-counter-submit')).toBeDisabled();
    expect(screen.getByTestId('circuit-confirm-counter-name-focus')).toBeDisabled();
    expectNoWrites();
  });

  it.each(['loading', 'error'])('does not invent missing-name evidence while definitions are %s', async (state) => {
    seedCounterDraft();
    measurementsQuery.isLoading = state === 'loading';
    measurementsQuery.isError = state === 'error';
    await show();
    expect(screen.queryByTestId('circuit-profile-circuit-required-error')).not.toBeInTheDocument();
    expect(screen.queryByTestId('circuit-confirm-counter-submit')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    expectNoWrites();
  });

  it('keeps existing definition repair independent of another new source missing its name', async () => {
    measurementsQuery.data = [counterDefinition('')];
    seedCounterDraft({ member_roles: { 'dev-1\0pt-0\0tag-0': 'counter_raw', 'dev-1\0pt-1\0tag-1': 'counter_raw' } });
    await show();
    const saved = screen.getByTestId('circuit-counter-confirmation-source-tag-0');
    const fresh = screen.getByTestId('circuit-counter-confirmation-source-tag-1');
    fireEvent.change(within(saved).getByTestId('circuit-confirm-counter-unit'), { target: { value: 'count' } });
    fireEvent.change(within(fresh).getByTestId('circuit-confirm-counter-unit'), { target: { value: 'count' } });
    expect(within(saved).getByTestId('circuit-confirm-counter-submit')).toBeEnabled();
    expect(within(fresh).getByTestId('circuit-confirm-counter-submit')).toBeDisabled();
    expectNoWrites();
  });

  it('removes the name prerequisite when the new source is no longer a counter', async () => {
    seedCounterDraft();
    await show();
    expandSources();
    expect(screen.getByTestId('circuit-profile-circuit-required-error')).toBeVisible();
    fireEvent.change(roleSelect(), { target: { value: 'gauge' } });
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveAttribute('aria-required', 'false');
    expect(screen.queryByTestId('circuit-profile-circuit-required-error')).not.toBeInTheDocument();
    expect(screen.queryByTestId('circuit-confirm-counter-submit')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled();
    expectNoWrites();
  });

  it('retains the cumulative circuit prerequisite without mislabeling unsupported text as a new counter', async () => {
    await show();
    expandSources();
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    for (const [tag, checked] of [['tag-0', false], ['tag-1', false], ['tag-2', true]] as const) {
      const input = screen.getByTestId(`circuit-profile-member-cum-${tag}`) as HTMLInputElement;
      if (input.checked !== checked) fireEvent.click(input);
    }
    expect(screen.getByTestId('circuit-profile-circuit-required-error')).toHaveTextContent('Enter a circuit name before enabling cumulative recording.');
    expect(screen.queryByTestId('circuit-confirm-counter-submit')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    expectNoWrites();
  });
});
