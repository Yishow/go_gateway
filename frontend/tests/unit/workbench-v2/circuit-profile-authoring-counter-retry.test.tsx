import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { guidedReadiness, guidedLifecycle } from '../../fixtures/guidedRecordingPreparation';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { createInstance } from 'i18next';
import { I18nextProvider } from 'react-i18next';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { CircuitProfileBasicPanel } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel';
import { studioV2WorkspaceMeasurementsAPI as measurementsAPI } from '@/services/studioV2WorkspaceMeasurements';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { MeasurementDefinition } from '@/types/measurement';
import productionI18n from '@/i18n/config';
import { savedGroupState } from '../../fixtures/writeGroupState';
import { longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { expandSources } from '../../fixtures/sourceReview';

const groupsQuery = {
  data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] as WriteGroup[] },
  isLoading: false, isError: false, refetch: vi.fn(),
};
const ensureProfile = vi.fn();
const updateGroup = vi.fn();
const start = vi.fn();
const applySchema = vi.fn();
const previewSchema = vi.fn();

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
  useWriteGroupSchemaApplyMutation: () => ({ mutateAsync: applySchema, isPending: false }),
  useWriteGroupSchemaPreviewMutation: () => ({ mutateAsync: previewSchema, isPending: false }),
  useWriteGroupSchemaOperationQuery: () => ({ data: undefined, isError: false, refetch: vi.fn() }),
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useRecordingStartMutation: () => ({ mutateAsync: start, isPending: false }),
  useRecordingStartOperationQuery: () => ({ data: undefined, isLoading: false, isError: false, refetch: vi.fn() }),
}));
// Measurement hooks and React Query are real; only the read/write service boundary is controlled.
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({
  studioV2WorkspaceMeasurementsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() },
}));

function seedDraft() {
  window.sessionStorage.setItem('wbv2.circuit-profile-draft.ws-1%2Fdev-1', JSON.stringify({
    circuit_id: 'cir-a', circuit_name: 'Draft A', instantaneous_interval_input: '37',
    instantaneous_keys: ['dev-1\0pt-0\0tag-0'], cumulative_keys: [], cumulative_enabled: false,
    member_roles: { 'dev-1\0pt-0\0tag-0': 'counter_raw' },
  }));
}

async function show(lng = 'en', readonly = false, state = savedGroupState()) {
  const i18n = createInstance();
  await i18n.init({ lng, resources: productionI18n.options.resources });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: Infinity }, mutations: { retry: false } } });
  canonicalClient = client; client.setQueryDefaults(studioV2WorkspaceKeys.writeGroups(), { structuralSharing: false }); publishGroups();
  const panel = (workspaceId = 'ws-1') => <I18nextProvider i18n={i18n}><QueryClientProvider client={client}>
    <CircuitProfileBasicPanel state={state} workspaceId={workspaceId} readonly={readonly} />
  </QueryClientProvider></I18nextProvider>;
  const view = render(panel());
  return { ...view, client, changeScope: (workspaceId: string, revision: string) => {
    groupsQuery.data = { ...groupsQuery.data, workspace_id: workspaceId, workspace_revision: revision };
    publishGroups(); view.rerender(panel(workspaceId));
  } };
}

function deferred() {
  let resolve!: (value: MeasurementDefinition[]) => void;
  let reject!: (error: Error) => void;
  const promise = new Promise<MeasurementDefinition[]>((yes, no) => { resolve = yes; reject = no; });
  return { promise, resolve, reject };
}

function confirmedCounter(workspaceId = 'ws-1'): MeasurementDefinition {
  return {
    id: 'meas-0', workspace_id: workspaceId, device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', equipment_id: 'cir-a',
    name: 'Counter', definition_revision: 'def-1', source_binding_revision: 'binding-1', series_epoch: 'epoch-1',
    quantity: '', semantic_kind: 'counter', unit: 'count', counter_policy: { allow_negative: false },
  };
}

function expectNoWrites() {
  for (const mutation of [ensureProfile, updateGroup, start, applySchema, previewSchema, measurementsAPI.create, measurementsAPI.update]) {
    expect(mutation).not.toHaveBeenCalled();
  }
}

// Guide design: one primary; semantic/source repair is local NAV, not a second disabled Start.
function expectPreparationBlocked() {
  const mutations = [ensureProfile, updateGroup, start, applySchema, previewSchema, vi.mocked(measurementsAPI.create), vi.mocked(measurementsAPI.update)]; const counts = mutations.map(fn => fn.mock.calls.length);
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
  vi.mocked(measurementsAPI.list).mockReset().mockRejectedValue(new Error('postgres://root:secret@private-db/production'));
  seedDraft();
});

describe('scoped counter definition read recovery', () => {
  it('shows a safe error and retry for instantaneous-only counter_raw without losing its draft', async () => {
    await show();
    const error = await screen.findByTestId('circuit-profile-measurements-error');
    expect(error).toHaveTextContent('Counter definitions could not be read. Retry before confirming or recording these counters. Your draft is preserved.');
    expect(screen.getByRole('button', { name: 'Retry' })).toBeEnabled();
    expect(screen.getByTestId('circuit-profile-toggle-cumulative')).not.toBeChecked();
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Draft A');
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('37');
    expandSources();
    expect(screen.getByTestId('circuit-profile-role-tag-0')).toHaveValue('counter_raw');
    expectPreparationBlocked();
    expect(document.body).not.toHaveTextContent('private-db');
    expect(document.body).not.toHaveTextContent('secret');
    expectNoWrites();
  });

  it('transitions loading to error to a single scoped retry and confirmed success without writes', async () => {
    const first = deferred(), retry = deferred();
    vi.mocked(measurementsAPI.list).mockImplementationOnce(() => first.promise).mockImplementationOnce(() => retry.promise);
    await show();
    expect(await screen.findByTestId('circuit-profile-measurements-loading')).toHaveTextContent('Checking counter definitions.');
    expectPreparationBlocked();
    await act(async () => first.reject(new Error('private-password')));
    await screen.findByTestId('circuit-profile-measurements-error');
    const button = screen.getByRole('button', { name: 'Retry' });
    fireEvent.click(button);
    fireEvent.click(button);
    await waitFor(() => expect(measurementsAPI.list).toHaveBeenCalledTimes(2));
    expect(measurementsAPI.list).toHaveBeenNthCalledWith(2, 'dev-1');
    expectPreparationBlocked();
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Draft A');
    await act(async () => retry.resolve([confirmedCounter()]));
    await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
    expect(screen.queryByTestId('circuit-profile-measurements-error')).not.toBeInTheDocument();
    expect(screen.queryByTestId('circuit-counter-confirmation')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('37');
    expectNoWrites();
  });

  it('keeps a late device A retry out of device B and preserves both drafts', async () => {
    const retryA = deferred();
    vi.mocked(measurementsAPI.list).mockImplementationOnce(async () => { throw new Error('A unavailable'); })
      .mockImplementation((deviceId) => deviceId === 'dev-1' ? retryA.promise : Promise.resolve([]));
    const state = savedGroupState();
    state.devices.push({ ...state.devices[0], id: 'dev-2', name: 'PLC B' });
    state.points.push({ ...state.points[0], id: 'point-b', device_id: 'dev-2' });
    state.mappings['point-b'] = { ...state.mappings['rule-1-p-0'], point_id: 'point-b', persisted_point_id: 'pt-b', tag_id: 'tag-b', device_id: 'dev-2' };
    await show('en', false, state);
    await screen.findByTestId('circuit-profile-measurements-error');
    fireEvent.click(screen.getByRole('button', { name: 'Retry' }));
    fireEvent.click(screen.getByRole('button', { name: 'PLC B' }));
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Draft B' } });
    await waitFor(() => expect(measurementsAPI.list).toHaveBeenCalledWith('dev-2'));
    expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled();
    await act(async () => retryA.resolve([confirmedCounter()]));
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Draft B');
    expect(screen.queryByTestId('circuit-profile-measurements-error')).not.toBeInTheDocument();
    expandSources();
    expect(screen.getByTestId('circuit-profile-role-tag-b')).toHaveValue('gauge');
    fireEvent.click(screen.getByRole('button', { name: 'PLC A' }));
    await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
    expandSources();
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Draft A');
    expect(screen.getByTestId('circuit-profile-role-tag-0')).toHaveValue('counter_raw');
    expectNoWrites();
  });

  it('does not let a retry for the old revision authorize the same device at a new revision', async () => {
    const oldRetry = deferred(), freshRead = deferred();
    vi.mocked(measurementsAPI.list).mockRejectedValueOnce(new Error('initial read failed'))
      .mockImplementationOnce(() => oldRetry.promise).mockImplementationOnce(() => freshRead.promise);
    const view = await show();
    await screen.findByTestId('circuit-profile-measurements-error');
    fireEvent.click(screen.getByRole('button', { name: 'Retry' }));
    view.changeScope('ws-1', 'wrev-2');
    await waitFor(() => expect(measurementsAPI.list).toHaveBeenCalledTimes(3));
    await act(async () => oldRetry.resolve([confirmedCounter()]));
    expectPreparationBlocked();
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Draft A');
    await act(async () => freshRead.reject(new Error('new revision read failed')));
    expect(await screen.findByTestId('circuit-profile-measurements-error')).toBeVisible();
    expectPreparationBlocked();
    expectNoWrites();
  });

  it('blocks stale cached definitions during refresh/error and deduplicates retry while preserving newer edits', async () => {
    const refresh = deferred(), retry = deferred();
    vi.mocked(measurementsAPI.list).mockResolvedValueOnce([confirmedCounter()])
      .mockImplementationOnce(() => refresh.promise).mockImplementationOnce(() => retry.promise);
    const view = await show();
    await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
    act(() => { void view.client.invalidateQueries({ queryKey: ['studio-v2-workspace', 'measurements'] }); });
    await waitFor(() => expect(measurementsAPI.list).toHaveBeenCalledTimes(2));
    expectPreparationBlocked();
    expect(screen.getByTestId('circuit-profile-measurements-loading')).toBeVisible();
    await act(async () => refresh.reject(new Error('stale cache cannot prove readiness')));
    await screen.findByTestId('circuit-profile-measurements-error');
    expectPreparationBlocked();
    const button = screen.getByRole('button', { name: 'Retry' });
    act(() => { fireEvent.click(button); fireEvent.click(button); });
    await waitFor(() => expect(measurementsAPI.list).toHaveBeenCalledTimes(3));
    expectPreparationBlocked();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Edited during retry' } });
    fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '51' } });
    await act(async () => retry.resolve([confirmedCounter()]));
    await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
    expect(measurementsAPI.list).toHaveBeenCalledTimes(3);
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Edited during retry');
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('51');
    expectNoWrites();
  });

  it('keeps workspace changes separate and treats foreign rows as a safe read failure', async () => {
    const oldRead = deferred(), newRead = deferred();
    vi.mocked(measurementsAPI.list).mockRejectedValueOnce(new Error('old failure'))
      .mockImplementationOnce(() => oldRead.promise).mockImplementationOnce(() => newRead.promise);
    const view = await show();
    await screen.findByTestId('circuit-profile-measurements-error');
    fireEvent.click(screen.getByRole('button', { name: 'Retry' }));
    window.sessionStorage.setItem('wbv2.circuit-profile-draft.ws-2%2Fdev-1', JSON.stringify({
      circuit_id: 'cir-w2', circuit_name: 'Workspace 2 draft', instantaneous_keys: ['dev-1\0pt-0\0tag-0'],
      member_roles: { 'dev-1\0pt-0\0tag-0': 'counter_raw' },
    }));
    view.changeScope('ws-2', 'wrev-1');
    await waitFor(() => expect(measurementsAPI.list).toHaveBeenCalledTimes(3));
    await act(async () => oldRead.resolve([confirmedCounter()]));
    expectPreparationBlocked();
    await act(async () => newRead.resolve([confirmedCounter('ws-1')]));
    expect(await screen.findByTestId('circuit-profile-measurements-error')).toBeVisible();
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Workspace 2 draft');
    expectPreparationBlocked();
    vi.mocked(measurementsAPI.list).mockResolvedValueOnce([confirmedCounter('ws-2')]);
    fireEvent.click(screen.getByRole('button', { name: 'Retry' }));
    await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
    expectNoWrites();
  });

  it('keeps readonly retry available with safe Traditional Chinese guidance', async () => {
    await show('zh-TW', true);
    expect(await screen.findByTestId('circuit-profile-measurements-error')).toHaveTextContent('無法讀取計數定義。請重試後再確認或記錄這些計數來源；你的草稿已保留。');
    const retry = screen.getByRole('button', { name: '重試' });
    expect(retry).toBeEnabled();
    vi.mocked(measurementsAPI.list).mockResolvedValueOnce([confirmedCounter()]);
    fireEvent.click(retry);
    await waitFor(() => expect(screen.queryByTestId('circuit-profile-measurements-error')).not.toBeInTheDocument());
    expectPreparationBlocked();
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Draft A');
    expectNoWrites();
  });

  it.each(['gauge', 'state'])('does not block %s when its former counter definition lookup failed', async (role) => {
    await show();
    await screen.findByTestId('circuit-profile-measurements-error');
    expandSources();
    fireEvent.change(screen.getByTestId('circuit-profile-role-tag-0'), { target: { value: role } });
    expect(screen.queryByTestId('circuit-profile-measurements-error')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled();
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('37');
    expectNoWrites();
  });

  it.each(['latest', undefined] as const)('requires lookup recovery for a saved counter_raw with reduction %s', async (reduction) => {
    window.sessionStorage.clear();
    const group = longProfileGroup({ id: 'saved-instant', circuitId: 'cir-saved', circuitName: 'Saved', interval: 15, points: [0, 1] });
    group.recording_profile = { kind: 'instantaneous', interval: 15, unit: 'second', ...(reduction ? { sample_reduction: reduction } : {}) };
    group.members[0] = { ...group.members[0], value_role: 'counter_raw', measurement_id: 'meas-0' };
    groupsQuery.data.groups = [group];
    const original = structuredClone(group);
    await show();
    expect(await screen.findByTestId('circuit-profile-measurements-error')).toBeVisible();
    expectPreparationBlocked();
    expect(screen.getByRole('button', { name: 'Retry' })).toBeEnabled();
    expect(group).toEqual(original);
    expectNoWrites();
  });

  it('shows cumulative lookup recovery only after a cumulative source is explicitly selected', async () => {
    window.sessionStorage.clear();
    const view = await show();
    await waitFor(() => expect(view.client.isFetching()).toBe(0));
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    expect(screen.queryByTestId('circuit-profile-measurements-error')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-cumulative-sources-error')).toBeVisible();
    fireEvent.click(screen.getByTestId('circuit-profile-member-cum-tag-0'));
    expect(screen.getByTestId('circuit-profile-measurements-error')).toBeVisible();
    expect(screen.getByRole('button', { name: 'Retry' })).toBeEnabled();
    expectPreparationBlocked();
    expectNoWrites();
  });

  it('preserves selected members across repeated failures and requires explicit confirmation after an empty successful read', async () => {
    await show();
    await screen.findByTestId('circuit-profile-measurements-error');
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    fireEvent.click(screen.getByTestId('circuit-profile-member-cum-tag-0'));
    fireEvent.change(screen.getByTestId('circuit-profile-interval-cumulative'), { target: { value: '19' } });
    fireEvent.click(screen.getByRole('button', { name: 'Retry' }));
    await waitFor(() => expect(measurementsAPI.list).toHaveBeenCalledTimes(2));
    await screen.findByTestId('circuit-profile-measurements-error');
    expectPreparationBlocked();
    vi.mocked(measurementsAPI.list).mockResolvedValueOnce([]);
    fireEvent.click(screen.getByRole('button', { name: 'Retry' }));
    expect(await screen.findByTestId('circuit-counter-confirmation')).toBeVisible();
    expectPreparationBlocked();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-0')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-inst-tag-0')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('37');
    expect(screen.getByTestId('circuit-profile-interval-cumulative')).toHaveValue('19');
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Draft A');
    expect(measurementsAPI.list).toHaveBeenCalledTimes(3);
    expectNoWrites();
  });
});
