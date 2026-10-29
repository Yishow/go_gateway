import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { CircuitProfileBasicPanel } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import type { WriteGroup, WriteGroupMutationResponse } from '@/types/studioV2WriteGroup';
import type { RecordingStartOperation } from '@/types/studioV2RecordingStart';
import { confirmPreparedStart, guidedLifecycle, guidedPreview, guidedReadiness, startGuided } from '../../fixtures/guidedRecordingPreparation';
import { WriteGroupResponseError } from '@/utils/studioV2WriteGroupJson';

const { query, measurements, ensure, update, start, operationQuery, schemaRender } = vi.hoisted(() => ({
  query: { data: { workspace_id: 'ws-1', workspace_revision: 'ws-old', groups: [] as WriteGroup[] }, refetch: vi.fn() },
  measurements: { data: [] as unknown[] | undefined, isLoading: false, isError: false, refetch: vi.fn() },
  ensure: vi.fn(), update: vi.fn(), start: vi.fn(), schemaRender: vi.fn(),
  operationQuery: { data: undefined as RecordingStartOperation | undefined, refetch: vi.fn() },
}));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const listeners = new Set<() => void>();
function publish() { listeners.forEach(listener => listener()); }
function useFixtureGroups() {
  const [, refresh] = React.useReducer(value => value + 1, 0);
  React.useEffect(() => { const notify = () => refresh(); listeners.add(notify); return () => { listeners.delete(notify); }; }, []);
  return query;
}
function canonical(id: string) { const found = query.data.groups.find(item => item.id === id); if (!found) throw new Error('Missing canonical group'); return found; }
function saved(result: WriteGroupMutationResponse) {
  query.data = { workspace_id: 'ws-1', workspace_revision: result.workspace_revision, groups: [...query.data.groups.filter(item => item.id !== result.group.id), result.group] };
  publish(); return result;
}
vi.mock('@/features/datalink/workbench-v2/steps/step4/useRecordingWorkspaceGuard', () => ({ useRecordingWorkspaceGuard: () => () => false }));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', async importOriginal => {
  const actual = await importOriginal<typeof import('@/hooks/datalink/useStudioV2WriteGroups')>();
  return { ...actual, useWriteGroupsQuery: useFixtureGroups,
    useWriteGroupDeliveryQuery: () => ({ data: undefined, isError: false, isLoading: false }),
    useEnsureCircuitProfileMutation: () => ({ mutateAsync: async (input: unknown) => saved(await ensure(input)), isPending: false }),
    useUpdateWriteGroupMutation: () => ({ mutateAsync: async (input: unknown) => saved(await update(input)), isPending: false }) };
});
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn(async (id: string) => guidedLifecycle(canonical(id), query.data.workspace_revision)) } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  readiness: vi.fn(async (id: string) => guidedReadiness(canonical(id), query.data.workspace_revision)),
  schemaPreview: vi.fn(async (id: string, request: { expected_group_revision: string; expected_workspace_revision: string }) => {
    expect(request.expected_workspace_revision).toBe(query.data.workspace_revision); expect(request.expected_group_revision).toBe(canonical(id).revision);
    return guidedPreview(canonical(id), query.data.workspace_revision);
  }),
  schemaApply: vi.fn(async (id: string, request: { token: string; operation_id: string; expected_workspace_revision: string }) => {
    const current = canonical(id); expect(request.expected_workspace_revision).toBe(query.data.workspace_revision); expect(request.token).toBe(guidedPreview(current, query.data.workspace_revision).token);
    saved({ workspace_revision: 'ws-confirmed', group: { ...current, destination: { ...current.destination, schema_revision: current.revision, schema_digest: 'confirmed-schema' } } });
    return { operation_id: request.operation_id, status: 'succeeded', executed_statements: 1, created_at: '', updated_at: '' };
  }), schemaOperation: vi.fn(),
} }));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useRecordingStartMutation: () => ({ mutateAsync: start, isPending: false }),
  useRecordingStartOperationQuery: () => operationQuery,
}));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => measurements,
  useCreateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn() }),
  useUpdateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn() }),
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel', async importOriginal => {
  const actual = await importOriginal<typeof import('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel')>();
  return { GroupSchemaPanel: (props: React.ComponentProps<typeof actual.GroupSchemaPanel>) => { schemaRender(props); return <actual.GroupSchemaPanel {...props} />; } };
});

const state = savedGroupState();
function group(): WriteGroup {
  return savedGroup({
    id: 'group-inst', layout_version: 'managed_long_v1', storage_mode: 'plain',
    members: [{ ...savedGroup().members[0], circuit_id: 'opaque-A', circuit_name: 'Water', target_column: '', required: false, measurement_id: '' }],
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' },
    destination: { ...savedGroup().destination, storage_strategy: 'managed' },
    row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    write_policy: { mode: 'append', dedupe_capability: 'receipt', recreate_missing: true },
  });
}
function operation(): RecordingStartOperation {
  return {
    operation_id: 'operation-A', workspace_id: 'ws-1', setup_revision: 'ws-old', action: 'recording_start', status: 'succeeded', stage: 'complete',
    intent_digest: 'a'.repeat(64), device_ids: ['dev-1'], devices: [{ device_id: 'dev-1', activated: true }],
    groups: [{ group_id: 'group-inst', group_revision: 'rev-1', applied_revision: 'rev-1', saved: true, applied: true, ready: true }],
    created_at: '2026-10-24T00:00:00Z', updated_at: '2026-10-24T00:00:01Z',
  };
}
function panel(props: Partial<React.ComponentProps<typeof CircuitProfileBasicPanel>> = {}) {
  return <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })}><CircuitProfileBasicPanel state={state} workspaceId="ws-1" readonly={false} {...props} /></QueryClientProvider>;
}
beforeEach(() => {
  vi.clearAllMocks(); window.sessionStorage.clear();
  query.data = { workspace_id: 'ws-1', workspace_revision: 'ws-old', groups: [group()] };
  query.refetch.mockImplementation(async () => { publish(); return { data: query.data }; });
  ensure.mockImplementation(async () => ({ group: group(), workspace_revision: 'ws-write' }));
  update.mockImplementation(async () => ({ group: group(), workspace_revision: 'ws-write' }));
  start.mockResolvedValue(operation());
  measurements.data = []; measurements.isLoading = false; measurements.isError = false;
  operationQuery.data = undefined; operationQuery.refetch.mockImplementation(async () => ({ data: operationQuery.data }));
});

it('retains a lost start response across reload without fabricating an operation or resending', async () => {
  start.mockRejectedValueOnce(new TypeError('Failed to fetch'));
  const first = render(panel());
  await startGuided('circuit-profile');
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
  await waitFor(() => expect(screen.getByTestId('recording-next-action')).toHaveTextContent('next.actions.recover_operation'));
  const request = start.mock.calls[0][0];
  const stored = Object.values(sessionStorage).filter(value => value.includes(request.request_id));
  const writes = ensure.mock.calls.length + update.mock.calls.length;
  first.unmount(); render(panel());
  await waitFor(() => expect(screen.getByTestId('recording-next-action')).toHaveTextContent('next.actions.recover_operation'));
  fireEvent.click(screen.getByTestId('recording-next-action'));
  expect(screen.queryByTestId('basic-recording-retry')).not.toBeInTheDocument();
  expect(start).toHaveBeenCalledTimes(1); expect(start.mock.calls[0][0]).toEqual(request);
  expect(Object.values(sessionStorage).filter(value => value.includes(request.request_id))).toEqual(stored);
  expect(ensure.mock.calls.length + update.mock.calls.length).toBe(writes);
});

it('blocks double starts and informs the parent throughout preparation', async () => {
  query.data.groups = [];
  let finish!: (result: { group: WriteGroup; workspace_revision: string }) => void;
  ensure.mockImplementation(() => new Promise((resolve) => { finish = resolve; }));
  const active = vi.fn(); render(panel({ onOperationActiveChange: active }));
  const prepare = await screen.findByTestId('circuit-profile-prepare'); fireEvent.click(prepare); fireEvent.click(prepare);
  expect(ensure).toHaveBeenCalledOnce();
  expect(active).toHaveBeenCalledWith(true);
  await act(async () => finish({ group: group(), workspace_revision: 'ws-write' }));
  expect(start).not.toHaveBeenCalled(); await confirmPreparedStart('circuit-profile');
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
});

it('does not show or activate an old scope after the destination changes during a request', async () => {
  let finish!: (result: RecordingStartOperation) => void;
  start.mockImplementation(() => new Promise((resolve) => { finish = resolve; }));
  const activated = vi.fn(); const view = render(panel({ onActivated: activated }));
  await startGuided('circuit-profile');
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
  const next = { ...state, db: { ...state.db, connector: { ...state.db.connector, connector_id: 'other-connector' } } };
  view.rerender(panel({ state: next, onActivated: activated }));
  await act(async () => finish(operation()));
  expect(activated).not.toHaveBeenCalled();
  expect(screen.queryByTestId('basic-recording-continue-runtime')).not.toBeInTheDocument();
});

it('marks the saved schema dirty after an interval edit', () => {
  render(panel());
  fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '30' } });
  expect(schemaRender.mock.lastCall?.[0].dirty).toBe(true);
  expect(screen.queryByTestId('group-schema-apply')).not.toBeInTheDocument(); expect(start).not.toHaveBeenCalled();
});

it('uses the authoritative revision obtained after schema confirmation', async () => {
  render(panel()); await startGuided('circuit-profile');
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
  expect(query.refetch).toHaveBeenCalled();
  expect(start.mock.calls[0][0].expected_workspace_revision).toBe('ws-confirmed');
});

it.each(['loading', 'error'])('does not offer measurement creation while definitions are %s', (status) => {
  measurements.data = undefined;
  measurements.isLoading = status === 'loading'; measurements.isError = status === 'error';
  render(panel()); fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
  fireEvent.click(screen.getByTestId('circuit-profile-member-cum-tag-0'));
  expect(screen.queryByTestId('circuit-counter-confirmation')).not.toBeInTheDocument();
  expect(screen.getByTestId(`circuit-profile-measurements-${status}`)).toBeInTheDocument();
  expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument(); expect(start).not.toHaveBeenCalled(); expect(ensure).not.toHaveBeenCalled();
});

it('keeps an Axios network error unknown instead of allowing a new request', async () => {
  start.mockRejectedValueOnce(Object.assign(new Error('Network Error'), { code: 'ERR_NETWORK' }));
  render(panel()); await startGuided('circuit-profile');
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
  await waitFor(() => expect(screen.getByTestId('recording-next-action')).toHaveTextContent('next.actions.recover_operation'));
  expect(screen.queryByTestId('basic-recording-retry')).not.toBeInTheDocument(); expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument(); expect(start).toHaveBeenCalledOnce();
});

it('does not call activation callbacks after the requesting panel unmounts', async () => {
  let finish!: (result: RecordingStartOperation) => void;
  start.mockImplementation(() => new Promise((resolve) => { finish = resolve; }));
  const activated = vi.fn(); const view = render(panel({ onActivated: activated }));
  await startGuided('circuit-profile');
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
  view.unmount(); await act(async () => finish(operation()));
  expect(activated).not.toHaveBeenCalled();
});

it('queries the same operation first and keeps unknown when that query fails', async () => {
  start.mockResolvedValueOnce({ ...operation(), status: 'unknown' });
  operationQuery.refetch.mockResolvedValueOnce({ isError: true, error: new TypeError('Failed to fetch') });
  render(panel()); await startGuided('circuit-profile');
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
  const check = await screen.findByTestId('basic-recording-check');
  await waitFor(() => expect(check).toBeEnabled()); fireEvent.click(check);
  await waitFor(() => expect(operationQuery.refetch).toHaveBeenCalledOnce());
  expect(start).toHaveBeenCalledOnce();
  expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument();
  expect(screen.getByText('step4.basic.query_failed')).toBeInTheDocument();
});

it('offers the original request for a failed operation whose next action is retry_same_request', async () => {
  const failed = { ...operation(), status: 'failed' as const, next_action: 'retry_same_request' };
  start.mockResolvedValueOnce(failed); operationQuery.data = failed;
  operationQuery.refetch.mockResolvedValueOnce({ data: failed });
  render(panel()); await startGuided('circuit-profile');
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
  const retry = await screen.findByTestId('basic-recording-retry');
  await waitFor(() => expect(retry).toBeEnabled()); fireEvent.click(retry);
  await waitFor(() => expect(start).toHaveBeenCalledTimes(2));
  expect(start.mock.calls[1][0]).toEqual(start.mock.calls[0][0]);
});

it('offers a runtime handoff using only the activated device IDs recovered by GET', async () => {
  start.mockResolvedValueOnce({ ...operation(), status: 'unknown' });
  operationQuery.refetch.mockResolvedValueOnce({ data: operation() });
  const proceed = vi.fn(); render(panel({ onContinueToRuntime: proceed }));
  await startGuided('circuit-profile');
  const check = await screen.findByTestId('basic-recording-check');
  await waitFor(() => expect(check).toBeEnabled()); fireEvent.click(check);
  const next = await screen.findByTestId('basic-recording-continue-runtime');
  fireEvent.click(next);
  expect(proceed).toHaveBeenCalledWith(['dev-1']);
  expect(start).toHaveBeenCalledOnce();
});

it('retains the original request when a successful HTTP response cannot be parsed', async () => {
  start.mockRejectedValueOnce(new WriteGroupResponseError('recording start', { success: true, data: {} }));
  render(panel()); await startGuided('circuit-profile');
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
  await waitFor(() => expect(screen.getByTestId('recording-next-action')).toHaveTextContent('next.actions.recover_operation'));
  expect(screen.queryByTestId('basic-recording-retry')).not.toBeInTheDocument(); expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument(); expect(start).toHaveBeenCalledOnce();
});
