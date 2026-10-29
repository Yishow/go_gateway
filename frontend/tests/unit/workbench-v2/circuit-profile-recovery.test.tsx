import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { CircuitProfileBasicPanel } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { RecordingStartOperation } from '@/types/studioV2RecordingStart';
import { WriteGroupResponseError } from '@/utils/studioV2WriteGroupJson';

const { query, measurements, ensure, update, start, operationQuery, schemaRender } = vi.hoisted(() => ({
  query: { data: { workspace_revision: 'ws-old', groups: [] as WriteGroup[] }, refetch: vi.fn() },
  measurements: { data: [] as unknown[] | undefined, isLoading: false, isError: false, refetch: vi.fn() },
  ensure: vi.fn(), update: vi.fn(), start: vi.fn(), schemaRender: vi.fn(),
  operationQuery: { data: undefined as RecordingStartOperation | undefined, refetch: vi.fn() },
}));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupsQuery: () => query,
  useWriteGroupDeliveryQuery: () => ({ data: undefined, isError: false, isLoading: false }),
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: ensure, isPending: false }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: update, isPending: false }),
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useRecordingStartMutation: () => ({ mutateAsync: start, isPending: false }),
  useRecordingStartOperationQuery: () => operationQuery,
}));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => measurements,
  useCreateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn() }),
  useUpdateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn() }),
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel', () => ({
  GroupSchemaPanel: (props: { dirty: boolean; workspaceRevision: string; onApplied: () => void }) => {
    schemaRender(props);
    return <button onClick={props.onApplied} disabled={props.dirty} data-testid="confirm-saved-schema">Confirm</button>;
  },
}));

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
  return <CircuitProfileBasicPanel state={state} workspaceId="ws-1" readonly={false} {...props} />;
}
beforeEach(() => {
  vi.clearAllMocks(); window.sessionStorage.clear();
  query.data = { workspace_revision: 'ws-old', groups: [group()] };
  query.refetch.mockImplementation(async () => ({ data: query.data }));
  ensure.mockImplementation(async () => ({ group: group(), workspace_revision: 'ws-write' }));
  update.mockImplementation(async () => ({ group: group(), workspace_revision: 'ws-write' }));
  start.mockResolvedValue(operation());
  measurements.data = []; measurements.isLoading = false; measurements.isError = false;
  operationQuery.data = undefined;
});

it('recovers a lost start response across reload using the original request without rewriting groups', async () => {
  start.mockRejectedValueOnce(new TypeError('Failed to fetch'));
  const first = render(panel());
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
  await screen.findByTestId('basic-recording-retry');
  expect(screen.getByTestId('circuit-profile-start')).toBeDisabled();
  const request = start.mock.calls[0][0];
  const writes = ensure.mock.calls.length + update.mock.calls.length;
  first.unmount(); render(panel());
  expect(await screen.findByTestId('basic-recording-retry')).toBeEnabled();
  fireEvent.click(screen.getByTestId('basic-recording-retry'));
  await waitFor(() => expect(start).toHaveBeenCalledTimes(2));
  expect(start.mock.calls[1][0]).toEqual(request);
  expect(ensure.mock.calls.length + update.mock.calls.length).toBe(writes);
});

it('blocks double starts and informs the parent throughout preparation', async () => {
  query.data.groups = [];
  let finish!: (result: { group: WriteGroup; workspace_revision: string }) => void;
  ensure.mockImplementation(() => new Promise((resolve) => { finish = resolve; }));
  const active = vi.fn(); render(panel({ onOperationActiveChange: active }));
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
  expect(ensure).toHaveBeenCalledOnce();
  expect(active).toHaveBeenCalledWith(true);
  await act(async () => finish({ group: group(), workspace_revision: 'ws-write' }));
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
});

it('does not show or activate an old scope after the destination changes during a request', async () => {
  let finish!: (result: RecordingStartOperation) => void;
  start.mockImplementation(() => new Promise((resolve) => { finish = resolve; }));
  const activated = vi.fn(); const view = render(panel({ onActivated: activated }));
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
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
  expect(screen.getByTestId('confirm-saved-schema')).toBeDisabled();
});

it('uses the authoritative revision obtained after schema confirmation', async () => {
  render(panel());
  fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
  await waitFor(() => expect(query.refetch).toHaveBeenCalled());
  query.data = { ...query.data, workspace_revision: 'ws-confirmed' };
  fireEvent.click(screen.getByTestId('confirm-saved-schema'));
  await waitFor(() => expect(query.refetch).toHaveBeenCalledTimes(2));
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
  expect(start.mock.calls[0][0].expected_workspace_revision).toBe('ws-confirmed');
});

it.each(['loading', 'error'])('does not offer measurement creation while definitions are %s', (status) => {
  measurements.data = undefined;
  measurements.isLoading = status === 'loading'; measurements.isError = status === 'error';
  render(panel()); fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
  fireEvent.click(screen.getByTestId('circuit-profile-member-cum-tag-0'));
  expect(screen.queryByTestId('circuit-counter-confirmation')).not.toBeInTheDocument();
  expect(screen.getByTestId(`circuit-profile-measurements-${status}`)).toBeInTheDocument();
  expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
});

it('keeps an Axios network error unknown instead of allowing a new request', async () => {
  start.mockRejectedValueOnce(Object.assign(new Error('Network Error'), { code: 'ERR_NETWORK' }));
  render(panel()); fireEvent.click(screen.getByTestId('circuit-profile-start'));
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
  expect(await screen.findByTestId('basic-recording-retry')).toBeEnabled();
  expect(screen.getByTestId('circuit-profile-start')).toBeDisabled();
});

it('does not call activation callbacks after the requesting panel unmounts', async () => {
  let finish!: (result: RecordingStartOperation) => void;
  start.mockImplementation(() => new Promise((resolve) => { finish = resolve; }));
  const activated = vi.fn(); const view = render(panel({ onActivated: activated }));
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
  view.unmount(); await act(async () => finish(operation()));
  expect(activated).not.toHaveBeenCalled();
});

it('queries the same operation first and keeps unknown when that query fails', async () => {
  start.mockResolvedValueOnce({ ...operation(), status: 'unknown' });
  operationQuery.refetch.mockResolvedValueOnce({ isError: true, error: new TypeError('Failed to fetch') });
  render(panel()); fireEvent.click(screen.getByTestId('circuit-profile-start'));
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
  const retry = await screen.findByTestId('basic-recording-retry');
  await waitFor(() => expect(retry).toBeEnabled()); fireEvent.click(retry);
  await waitFor(() => expect(operationQuery.refetch).toHaveBeenCalledOnce());
  expect(start).toHaveBeenCalledOnce();
  expect(screen.getByTestId('circuit-profile-start')).toBeDisabled();
  expect(screen.getByText('step4.basic.query_failed')).toBeInTheDocument();
});

it('offers the original request for a failed operation whose next action is retry_same_request', async () => {
  const failed = { ...operation(), status: 'failed' as const, next_action: 'retry_same_request' };
  start.mockResolvedValueOnce(failed);
  operationQuery.refetch.mockResolvedValueOnce({ data: failed });
  render(panel()); fireEvent.click(screen.getByTestId('circuit-profile-start'));
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
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
  const check = await screen.findByTestId('basic-recording-check');
  await waitFor(() => expect(check).toBeEnabled()); fireEvent.click(check);
  const next = await screen.findByTestId('basic-recording-continue-runtime');
  fireEvent.click(next);
  expect(proceed).toHaveBeenCalledWith(['dev-1']);
  expect(start).toHaveBeenCalledOnce();
});

it('retains the original request when a successful HTTP response cannot be parsed', async () => {
  start.mockRejectedValueOnce(new WriteGroupResponseError('recording start', { success: true, data: {} }));
  render(panel()); fireEvent.click(screen.getByTestId('circuit-profile-start'));
  await waitFor(() => expect(start).toHaveBeenCalledOnce());
  expect(await screen.findByTestId('basic-recording-retry')).toBeEnabled();
  expect(screen.getByTestId('circuit-profile-start')).toBeDisabled();
});
