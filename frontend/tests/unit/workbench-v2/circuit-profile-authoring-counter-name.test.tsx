import type { ReactNode } from 'react';
import { useCircuitRecordingCommandFence } from '@/features/datalink/workbench-v2/steps/step4/useCircuitRecordingCommandFence';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { basicRecordingScopeKey, saveBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { guidedReadiness, guidedLifecycle } from '../../fixtures/guidedRecordingPreparation';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, renderHook, screen, waitFor, within } from '@testing-library/react';
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
import { longProfileGroup } from '../../fixtures/circuitBatchPanel';

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

async function show(state = savedGroupState(), readonly = false, lng = 'en') {
  const i18n = createInstance();
  await i18n.init({ lng, resources: productionI18n.options.resources });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  canonicalClient = client; client.setQueryDefaults(studioV2WorkspaceKeys.writeGroups(), { structuralSharing: false }); publishGroups();
  const panel = () => <I18nextProvider i18n={i18n}><QueryClientProvider client={client}>
    <BasicRecordingPanel state={state} workspaceId="ws-1" readonly={readonly} />
  </QueryClientProvider></I18nextProvider>;
  const view = render(panel());
  return { ...view, client, refresh: () => { publishGroups(); view.rerender(panel()); } };
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
    expectPreparationBlocked();
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
    expectPreparationBlocked();
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
      expectPreparationBlocked();
    }
    expectNoWrites();
    if (state === 'missing-unit') {
      fireEvent.click(screen.getByTestId('circuit-confirm-counter-submit'));
      await waitFor(() => expect(updateMeasurement).toHaveBeenCalledTimes(1));
      expect(updateMeasurement).toHaveBeenCalledWith({ id: 'meas-0', def: { ...counterDefinition(''), unit: 'count' } });
      expect(createMeasurement).not.toHaveBeenCalled(); expect(ensureProfile).not.toHaveBeenCalled(); expect(start).not.toHaveBeenCalled();
      expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('');
    }
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
    expectPreparationBlocked();
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
    fireEvent.click(within(saved).getByTestId('circuit-confirm-counter-submit'));
    await waitFor(() => expect(updateMeasurement).toHaveBeenCalledTimes(1));
    expect(updateMeasurement).toHaveBeenCalledWith({ id: 'meas-0', def: { ...counterDefinition(''), unit: 'count' } });
    fireEvent.click(within(fresh).getByTestId('circuit-confirm-counter-submit'));
    expect(createMeasurement).not.toHaveBeenCalled(); expect(ensureProfile).not.toHaveBeenCalled(); expect(start).not.toHaveBeenCalled();
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('');
  });

  it.each(['foreign', 'binding', 'loading', 'error'] as const)('does not repair a saved definition from %s evidence', async condition => {
    measurementsQuery.data = [{ ...counterDefinition(''), ...(condition === 'foreign' ? { workspace_id: 'other' } : condition === 'binding' ? { source_binding_revision: '' } : {}) }];
    measurementsQuery.isLoading = condition === 'loading'; measurementsQuery.isError = condition === 'error'; seedCounterDraft({ circuit_name: '' });
    await show();
    const unit = screen.queryByTestId('circuit-confirm-counter-unit'); if (unit) fireEvent.change(unit, { target: { value: 'count' } });
    const submit = screen.queryByTestId('circuit-confirm-counter-submit'); if (submit) { expect(submit).toBeDisabled(); fireEvent.click(submit); }
    expectNoWrites();
  });
  it.each(['foreign-cache', 'offline', 'unreadable-storage', 'pending-start'] as const)('rechecks %s before an unnamed saved-definition update', async condition => {
    measurementsQuery.data = [counterDefinition('')]; seedCounterDraft({ circuit_name: '' }); const view = await show();
    fireEvent.change(screen.getByTestId('circuit-confirm-counter-unit'), { target: { value: 'count' } });
    const submit = screen.getByTestId('circuit-confirm-counter-submit'); await waitFor(() => expect(submit).toBeEnabled());
    const restore = condition === 'offline' ? vi.spyOn(window.navigator, 'onLine', 'get').mockReturnValue(false)
      : condition === 'unreadable-storage' ? vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => { throw new Error('unreadable'); }) : undefined;
    try {
      act(() => {
        if (condition === 'foreign-cache') view.client.setQueryData(studioV2WorkspaceKeys.writeGroups(), { ...groupsQuery.data, workspace_id: 'other' });
        if (condition === 'pending-start') saveBasicRecordingIntent({ scope_key: basicRecordingScopeKey('ws-1', 'dev-1', encodeURIComponent(JSON.stringify(['', 'conn-1', 'crev-1', '']))), operation_id: 'O1',
          request: { request_id: 'original-request', workspace_id: 'ws-1', expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'], groups: [] } });
        fireEvent.click(submit);
      });
      expectNoWrites();
    } finally { restore?.mockRestore(); }
  });

  it.each(['open', 'pending', 'closed'] as const)('rejects the captured saved-definition submit while the real batch editor is %s', async phase => {
    const first = longProfileGroup({ id: 'g-A', circuitId: 'cir-draft-1', circuitName: 'Circuit A', interval: 15, points: [0, 1] });
    first.members[0].value_role = 'counter_raw';
    groupsQuery.data.groups = [first, longProfileGroup({ id: 'g-B', circuitId: 'cir-B', circuitName: 'Circuit B', interval: 15, points: [0, 1] })];
    measurementsQuery.data = [counterDefinition('')];
    let finish: (() => void) | undefined;
    updateGroup.mockImplementation(() => new Promise((_resolve, reject) => { finish = () => reject(new Error('test-only rejected batch')); }));
    await show(); fireEvent.change(screen.getByTestId('circuit-confirm-counter-unit'), { target: { value: 'count' } });
    const submit = screen.getByTestId('circuit-confirm-counter-submit'); await waitFor(() => expect(submit).toBeEnabled());
    // Invoke the real mounted React closure after another owner opens/starts; a disabled DOM click is not evidence.
    const key = Object.keys(submit).find(value => value.startsWith('__reactProps$')); if (!key) throw new Error('Missing mounted React props');
    const captured = (submit as unknown as Record<string, { onClick: () => void }>)[key].onClick;
    fireEvent.click(screen.getByTestId('circuit-profile-open-batch'));
    if (phase === 'closed') fireEvent.click(screen.getByTestId('circuit-profile-batch-close'));
    if (phase === 'pending') {
      fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
      fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));
      await waitFor(() => expect(updateGroup).toHaveBeenCalledTimes(1));
    }
    try {
      await act(async () => captured());
      expect(updateMeasurement).not.toHaveBeenCalled(); expect(createMeasurement).not.toHaveBeenCalled();
      expect(ensureProfile).not.toHaveBeenCalled(); expect(start).not.toHaveBeenCalled();
      expect(updateGroup).toHaveBeenCalledTimes(phase === 'pending' ? 1 : 0);
      if (phase === 'closed') {
        expect(screen.getByTestId('circuit-confirm-counter-submit')).toBeEnabled();
        fireEvent.click(screen.getByTestId('circuit-confirm-counter-submit'));
        await waitFor(() => expect(updateMeasurement).toHaveBeenCalledTimes(1));
        expect(updateMeasurement).toHaveBeenCalledWith({ id: 'meas-0', def: { ...counterDefinition(''), unit: 'count' } });
      }
    } finally { await act(async () => finish?.()); }
  });

  it('rejects an old saved-repair callback after A→B→A and unmount without relaxing batch identity', () => {
    const client = new QueryClient(); client.setQueryData(studioV2WorkspaceKeys.writeGroups(), groupsQuery.data);
    const candidate = buildGroupCandidates(savedGroupState()).candidates[0];
    const repair = { available: true, statuses: [{ candidate, confirmed: false, measurement: counterDefinition('') }] };
    const input = { ownerKey: 'owner-A', scopeKey: 'scope-A', workspaceId: 'ws-1', deviceId: 'dev-1', circuitId: '', batchGroups: [],
      consentKey: 'native-controls', updateWriteGroup: updateGroup, ensureCircuitProfile: ensureProfile };
    const view = renderHook(({ counterRepair }) => useCircuitRecordingCommandFence({ ...input, counterRepair }), {
      initialProps: { counterRepair: repair }, wrapper: ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>,
    });
    act(() => view.result.current.onCommandFence({ ownerKey: input.ownerKey, scopeKey: input.scopeKey, blocked: false, isBlocked: () => false,
      handoff: { workspaceRevision: 'wrev-1', consentKey: input.consentKey, state: () => 'current' } }));
    expect(view.result.current.isCounterRepairBlocked()).toBe(false); expect(view.result.current.isBlocked()).toBe(true);
    const old = view.result.current.isCounterRepairBlocked;
    view.rerender({ counterRepair: { ...repair, available: false } }); view.rerender({ counterRepair: repair });
    expect(old()).toBe(true); expect(view.result.current.isCounterRepairBlocked()).toBe(false);
    const current = view.result.current.isCounterRepairBlocked; view.unmount(); expect(current()).toBe(true);
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
    expectPreparationBlocked();
    expectNoWrites();
  });
});
