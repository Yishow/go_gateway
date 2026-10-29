import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import type { WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import { studioV2RecordingStartAPI } from '@/services/studioV2RecordingStart';
import { studioV2WorkspaceMeasurementsAPI } from '@/services/studioV2WorkspaceMeasurements';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { RecordingStartOperation, RecordingStartRequest } from '@/types/studioV2RecordingStart';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { batchPanelState, longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { clickEnabled, guidedReadiness } from '../../fixtures/guidedRecordingPreparation';
import { parseRecordingStartOperation } from '@/utils/studioV2RecordingStartJson';
import { savedLifecycle } from '../../fixtures/writeGroupLifecycle';
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn() } }));
vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string, options?: unknown) => (options && typeof options === 'object' && 'count' in options ? `${key}:${(options as { count: number }).count}` : key) }),
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }),
}));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), ensureCircuitProfile: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(),
    remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(),
  },
}));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn(async () => []) } }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), operation: vi.fn(), ensureBasic: vi.fn() } }));
const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const startApi = vi.mocked(studioV2RecordingStartAPI);
const measurementApi = vi.mocked(studioV2WorkspaceMeasurementsAPI);
const activateWorkspace = vi.fn();
/** Two saved devices: dev-1 (PLC A, six float sources from the shared fixture) and dev-2 (PLC B, two float sources). */
function twoDeviceState(): WorkbenchV2State {
  const base = batchPanelState();
  const points = [0, 1].map((i) => ({ ...base.points[0], id: `rule-2-p-${i}`, device_id: 'dev-2', rule_id: 'rule-2', name: `b.${i}`, address: `5000${i + 1}` }));
  const mappings = Object.fromEntries(points.map((point, i) => [point.id, {
    ...base.mappings['rule-1-p-0'], point_id: point.id, tag_key: `b.${i}`, display_name: `b.${i}`,
    persisted_point_id: `pt-b${i}`, tag_id: `tag-b${i}`, device_id: 'dev-2',
  }]));
  return {
    ...base,
    devices: [...base.devices, { ...base.devices[0], id: 'dev-2', name: 'PLC B' }],
    rules: [...base.rules, { ...base.rules[0], id: 'rule-2', device_id: 'dev-2' }],
    points: [...base.points, ...points],
    mappings: { ...base.mappings, ...mappings },
  };
}
const memberOf = (deviceId: string, point: string, tag: string, circuitId: string, circuitName: string) => ({
  device_id: deviceId, point_id: point, tag_id: tag, source_revision: `s-${point}`, mapping_revision: `m-${point}`, target_column: '',
  required: false, circuit_id: circuitId, circuit_name: circuitName,
});
/** Device A's saved managed long group, ready to record. */
const withSchema = (group: WriteGroup): WriteGroup => ({ ...group, destination: { ...group.destination, schema_revision: group.revision, schema_digest: 'fixture-proof' } });
const groupA = (): WriteGroup => withSchema(longProfileGroup({ id: 'grp-a', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 15, points: [0, 1] }));
/** Device B's saved managed long group that was never prepared: not ready, not applied. */
const groupB = (): WriteGroup => ({
  ...longProfileGroup({ id: 'grp-b', circuitId: 'cir-b', circuitName: 'Feeder B', interval: 15, points: [0, 1] }),
  members: [memberOf('dev-2', 'pt-b0', 'tag-b0', 'cir-b', 'Feeder B'), memberOf('dev-2', 'pt-b1', 'tag-b1', 'cir-b', 'Feeder B')],
});
/** A new wide profile group of device B, saved and not yet applied. */
const wideB = (): WriteGroup => withSchema(savedGroup({
  id: 'grp-w', revision: 'rev-grp-w-1', name: 'Feeder W (instantaneous)', status: 'ready', layout_version: 'managed_wide_profile_v1',
  recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'mean' },
  destination: { ...savedGroup().destination, storage_strategy: 'managed', table_name: 'gw_wide_fixture' },
  members: [memberOf('dev-2', 'pt-b0', 'tag-b0', 'cir-w', 'Feeder W'), memberOf('dev-2', 'pt-b1', 'tag-b1', 'cir-w', 'Feeder W')]
    .map((member) => ({ ...member, target_column: `col_${member.point_id}`, value_role: 'gauge' as const, source_value_type: 'float64' as const, result_value_type: 'float64' as const })),
  row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
  write_policy: { mode: 'append', dedupe_capability: 'receipt' },
}));
/** Device B's second saved managed long circuit, so the same device has two circuits to choose between. */
const longB2 = (): WriteGroup => ({
  ...longProfileGroup({ id: 'grp-b2', circuitId: 'cir-b2', circuitName: 'Feeder B2', interval: 15, points: [0, 1] }),
  members: [memberOf('dev-2', 'pt-b0', 'tag-b0', 'cir-b2', 'Feeder B2'), memberOf('dev-2', 'pt-b1', 'tag-b1', 'cir-b2', 'Feeder B2')],
});
/** Device A's saved legacy wide group remains accessible from the shared recording setup. */
const legacyA = (): WriteGroup => savedGroup({
  id: 'grp-legacy', revision: 'rev-grp-legacy-1', name: 'PLC A', status: 'ready', basic_managed_device_id: 'dev-1',
  destination: { ...savedGroup().destination, storage_strategy: 'managed', table_name: 'gw_group_legacy', database: '/recording.db' },
  members: ['pt-0', 'pt-1'].map((point, i) => ({ ...memberOf('dev-1', point, `tag-${i}`, '', ''), target_column: '' })),
  row_policy: { interval_seconds: 60, allowed_lateness_seconds: 0 },
});
function operationFor(request: RecordingStartRequest, status: RecordingStartOperation['status'] = 'succeeded'): RecordingStartOperation {
  const done = status === 'succeeded';
  const parsed = parseRecordingStartOperation({
    operation_id: `op-${request.request_id}`, action: 'recording_start', status, stage: done ? 'complete' : 'save',
    intent_digest: 'c'.repeat(64), workspace_id: request.workspace_id, setup_revision: request.expected_workspace_revision,
    device_ids: request.device_ids,
    groups: request.groups.map((group) => ({ group_id: group.group_id, group_revision: group.expected_group_revision, saved: true, ready: done, applied: done, ...(done ? { applied_revision: group.expected_group_revision } : {}) })),
    devices: request.device_ids.map((deviceId) => ({ device_id: deviceId, activated: done })),
    created_at: '2026-10-25T00:00:00Z', updated_at: '2026-10-25T00:00:01Z',
  });
  if (!parsed) throw new Error('Invalid operation fixture (studioV2RecordingStartJson)'); return parsed;
}
const lostResponse = () => Object.assign(new Error('network lost'), { outcome: 'unconfirmed' });
const deferred = <T,>() => {
  let resolve!: (value: T) => void;
  const promise = new Promise<T>((done) => { resolve = done; });
  return { promise, resolve };
};
const emptyDelivery = {
  group_id: 'x', intake: { state: 'not_running' }, stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
  last_sql_committed_at: null, oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
  quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
};
let groups: WriteGroup[] = [];
beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  groups = [];
  vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockImplementation(async (id) => {
    const group = groups.find((item) => item.id === id);
    if (!group) throw new Error('Lifecycle fixture group missing');
    return savedLifecycle(group, 'wrev-1');
  });
  measurementApi.list.mockResolvedValue([]);
  api.list.mockImplementation(async () => ({ workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups }));
  api.delivery.mockResolvedValue(emptyDelivery as never);
  api.readiness.mockImplementation(async id => {
    const group = groups.find(item => item.id === id); if (!group) throw new Error('Missing canonical group');
    const proof = guidedReadiness(group, 'wrev-1');
    return id === 'grp-b' ? { ...proof, config_ready: false, issues: [{ code: 'table-missing', severity: 'blocking', step: 'Step 4', scope: id, message: 'The table is missing' }] } : proof;
  });
  api.apply.mockImplementation(async (id: string) => {
    const saved = groups.find((group) => group.id === id) as WriteGroup;
    return { workspace_revision: 'wrev-2', group: { ...saved, applied_revision: saved.revision } } as never;
  });
  startApi.start.mockImplementation(async (request) => operationFor(request));
  startApi.operation.mockImplementation(async (id: string) => operationFor({
    request_id: id.replace(/^op-/, ''), workspace_id: 'ws-1', expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
    groups: [{ group_id: 'grp-a', expected_group_revision: 'rev-grp-a-1', expected_connector_revision: 'crev-1' }],
  }));
});
function renderHost(saved: WriteGroup[], state: WorkbenchV2State = twoDeviceState()) {
  groups = saved;
  return render(
    <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })}>
      <Step4Database state={state} dispatch={vi.fn()} workspaceId="ws-1" activateWorkspace={activateWorkspace} />
    </QueryClientProvider>,
  );
}
const tile = (deviceId: string) => screen.getAllByTestId('basic-recording-device-tile').find((el) => el.getAttribute('data-device-id') === deviceId) as HTMLElement;
const startButton = () => screen.findByTestId('circuit-profile-start');
const startedRequests = () => startApi.start.mock.calls.map(([request]) => request);
// Saved proof is already current: wait for this scope's primary, never click the open Advanced editor's preview.
const confirmScopedStart = async () => { await clickEnabled('circuit-profile-start'); await clickEnabled('recording-start-confirm'); };
describe('One primary scoped start', () => {
  it('starts only the reviewed profile while preserving a never-applied legacy draft on the same device', async () => {
    const legacy = legacyA();
    const before = structuredClone(legacy);
    renderHost([legacy, groupA(), groupB()]);
    expect(await screen.findByTestId('device-recording-existing-draft')).toBeInTheDocument();
    await confirmScopedStart();
    await waitFor(() => expect(startApi.start).toHaveBeenCalledTimes(1));
    expect(startedRequests()[0].device_ids).toEqual(['dev-1']);
    expect(startedRequests()[0].groups.map((group) => group.group_id)).toEqual(['grp-a']);
    expect(groups.find((group) => group.id === legacy.id)).toEqual(before);
    for (const call of [api.create, api.update, api.disable, api.remove, api.apply, startApi.ensureBasic, activateWorkspace]) expect(call).not.toHaveBeenCalled();
  });
  it('keeps an applied legacy owner intact when the server refuses an overlapping new start', async () => {
    const legacy = { ...legacyA(), applied_revision: 'legacy-applied', status: 'running' as const };
    const before = structuredClone(legacy);
    startApi.start.mockRejectedValueOnce(Object.assign(new Error('private backend diagnostic'), {
      response: { status: 409, data: { error: { code: 'WRITE_GROUP_CUTOVER_REQUIRED' } } },
    }));
    renderHost([legacy, groupA()]);
    expect(await screen.findByTestId('device-recording-existing-applied')).toBeInTheDocument();
    expect(screen.getByTestId('device-recording-review-cutover')).toBeInTheDocument();
    await confirmScopedStart();
    expect(await screen.findByTestId('circuit-profile-error')).not.toHaveTextContent('private backend diagnostic');
    expect(startedRequests()[0].groups.map((group) => group.group_id)).toEqual(['grp-a']);
    expect(groups.find((group) => group.id === legacy.id)).toEqual(before);
    for (const call of [api.update, api.disable, api.remove, api.apply, activateWorkspace]) expect(call).not.toHaveBeenCalled();
  });
  it('offers the reviewed start and no second workspace-wide activation from the overview', async () => {
    renderHost([groupA(), groupB()]);
    expect(await startButton()).toBeInTheDocument();
    expect(screen.getAllByTestId('circuit-profile-start')).toHaveLength(1);
    expect(screen.getByTestId('step4-group-summary-note'), 'the overview keeps its status').toBeInTheDocument();
    expect(screen.queryByText('step4.activate_btn'), 'the overview must not present a second activation action').not.toBeInTheDocument();
    expect(activateWorkspace).not.toHaveBeenCalled();
  });
  it('starts only the viewed device and its frozen groups while another device is incomplete', async () => {
    renderHost([groupA(), groupB()]);
    await confirmScopedStart();
    await waitFor(() => expect(startApi.start).toHaveBeenCalledTimes(1));
    const [request] = startedRequests();
    expect(request).toMatchObject({ workspace_id: 'ws-1', device_ids: ['dev-1'] });
    expect(request.request_id).toBeTruthy();
    expect(request.groups.map((group) => group.group_id)).toEqual(['grp-a']);
    expect(request.groups[0]).toMatchObject({ expected_connector_revision: 'crev-1' });
    for (const call of [...api.ensureCircuitProfile.mock.calls, ...api.update.mock.calls]) {
      expect(JSON.stringify(call), 'the unrelated device is neither created nor rewritten').not.toMatch(/dev-2|grp-b/);
    }
    expect(api.apply).not.toHaveBeenCalled();
    expect(api.disable).not.toHaveBeenCalled();
    expect(activateWorkspace).not.toHaveBeenCalled();
  });
  it('keeps the same primary entry for Share-only without any database group or writer', async () => {
    const state = savedGroupState();
    const shareOnly: WorkbenchV2State = {
      ...state, devices: state.devices,
      db: { ...state.db, connector: { ...state.db.connector, persisted: false, save_state: 'idle', connector_id: undefined, identity_revision: undefined } },
    };
    groups = [];
    render(
      <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })}>
        <Step4Database
          state={shareOnly} dispatch={vi.fn()} workspaceId="ws-1" activateWorkspace={activateWorkspace}
          shareStatus={{ workspace_id: 'ws-1', enabled: true, configured_enabled: true, port: 502, address: '127.0.0.1', bind_state: 'pass', mapping_count: 1, readiness_token: 'share-token', settings_revision: 'settings-1', workspace_revision: 'share-wrev-1' }}
        />
      </QueryClientProvider>,
    );
    await waitFor(() => expect(screen.getByTestId('basic-recording-start')).toBeEnabled());
    fireEvent.click(screen.getByTestId('basic-recording-start'));
    await waitFor(() => expect(startApi.start).toHaveBeenCalledTimes(1));
    expect(startedRequests()[0]).toMatchObject({ device_ids: ['dev-1'], groups: [], readiness_token: 'share-token', settings_revision: 'settings-1', workspace_revision: 'share-wrev-1' });
    for (const call of [api.create, api.ensureCircuitProfile, api.update, api.apply, api.schemaPreview, api.schemaApply, startApi.ensureBasic]) {
      expect(call).not.toHaveBeenCalled();
    }
    expect(activateWorkspace).not.toHaveBeenCalled();
  });
});
describe('A lost start response while switching devices', () => {
  const startA = async () => {
    startApi.start.mockRejectedValueOnce(lostResponse());
    await confirmScopedStart();
    await waitFor(() => expect(startApi.start).toHaveBeenCalledTimes(1));
    await waitFor(() => expect(screen.getByTestId('recording-next-action')).toHaveTextContent('next.actions.recover_operation'));
  };
  it('keeps an unidentified lost ACK scoped to A across navigation and reload without inventing a query or resend', async () => {
    const first = renderHost([groupA(), groupB()]);
    await startA();
    const original = startedRequests()[0];
    expect(original.device_ids).toEqual(['dev-1']);
    fireEvent.click(tile('dev-2'));
    await waitFor(() => expect(tile('dev-2')).toHaveAttribute('aria-pressed', 'true'));
    expect(screen.queryByTestId('basic-recording-retry'), 'B does not inherit A\'s retry').not.toBeInTheDocument();
    expect(screen.queryByTestId('basic-recording-check')).not.toBeInTheDocument();
    expect(screen.queryByTestId('basic-recording-status')).not.toBeInTheDocument();
    expect(screen.queryByTestId('circuit-profile-error')).not.toBeInTheDocument();
    fireEvent.click(tile('dev-1'));
    await waitFor(() => expect(screen.getByTestId('recording-next-action')).toHaveTextContent('next.actions.recover_operation'));
    fireEvent.click(screen.getByTestId('recording-next-action'));
    // Approved recovery contract: absent O1 permits local guidance, never a fabricated GET or POST retry.
    expect(screen.queryByTestId('basic-recording-retry')).not.toBeInTheDocument();
    expect(startApi.operation).not.toHaveBeenCalled(); expect(startApi.start).toHaveBeenCalledTimes(1);
    expect(startedRequests()[0]).toEqual(original);
    first.unmount(); renderHost([groupA(), groupB()]);
    await waitFor(() => expect(screen.getByTestId('recording-next-action')).toHaveTextContent('next.actions.recover_operation'));
    expect(screen.queryByTestId('basic-recording-retry')).not.toBeInTheDocument();
    expect(startApi.operation).not.toHaveBeenCalled(); expect(startApi.start).toHaveBeenCalledTimes(1);
  });
  it('does not credit A\'s late success to B and shows it again on A', async () => {
    const reply = deferred<RecordingStartOperation>();
    startApi.start.mockImplementationOnce(() => reply.promise);
    renderHost([groupA(), groupB()]);
    await confirmScopedStart();
    await waitFor(() => expect(startApi.start).toHaveBeenCalledTimes(1));
    const original = startedRequests()[0];
    fireEvent.click(tile('dev-2'));
    await waitFor(() => expect(tile('dev-2')).toHaveAttribute('aria-pressed', 'true'));
    await act(async () => { reply.resolve(operationFor(original)); await reply.promise; });
    expect(screen.queryByTestId('basic-recording-status'), 'A\'s success is not shown on B').not.toBeInTheDocument();
    expect(screen.queryByTestId('basic-recording-continue-runtime')).not.toBeInTheDocument();
    expect(screen.queryByTestId('basic-recording-retry')).not.toBeInTheDocument();
    expect(startApi.start).toHaveBeenCalledTimes(1);
    fireEvent.click(tile('dev-1'));
    await waitFor(() => expect(screen.getByTestId('basic-recording-status')).toHaveTextContent('step4.basic.operation.succeeded'));
    expect(startApi.start, 'returning to A never posts a new start').toHaveBeenCalledTimes(1);
  });
  it('restores A\'s unresolved request after a reload and retries it unchanged', async () => {
    const first = renderHost([groupA(), groupB()]);
    // A typed partial O1 with retry_same_request, not an unidentified lost ACK, authorizes original-request recovery.
    startApi.start.mockImplementationOnce(async request => ({ ...operationFor(request, 'partial'), next_action: 'retry_same_request' }));
    startApi.operation.mockImplementation(async () => ({ ...operationFor(startedRequests()[0], 'partial'), next_action: 'retry_same_request' }));
    await confirmScopedStart();
    await screen.findByTestId('basic-recording-retry');
    const original = startedRequests()[0];
    first.unmount();
    renderHost([groupA(), groupB()]);
    await clickEnabled('basic-recording-retry');
    await waitFor(() => expect(startApi.start).toHaveBeenCalledTimes(2));
    expect(startedRequests()[1]).toEqual(original);
    expect(startApi.operation).toHaveBeenCalledWith(`op-${original.request_id}`);
    expect(screen.queryByTestId('circuit-profile-error')).not.toBeInTheDocument();
  });
});
describe('Advanced recording effectiveness', () => {
  const openWide = async () => {
    fireEvent.click(await screen.findByTestId('group-open-grp-w'));
    await screen.findByTestId('group-editor');
    await waitFor(() => expect(screen.getByTestId('group-apply')).toBeEnabled());
  };
  it('does not apply a managed profile group on its own authority', async () => {
    renderHost([groupA(), wideB()]);
    await openWide();
    fireEvent.click(screen.getByTestId('group-apply'));
    await act(async () => { await Promise.resolve(); });
    expect(api.apply, 'Apply of a managed profile belongs to the reviewed scoped start').not.toHaveBeenCalled();
    expect(activateWorkspace).not.toHaveBeenCalled();
  });
  it('routes the managed profile group back to its own device in the shared start, without starting it', async () => {
    renderHost([groupA(), wideB()]);
    await screen.findAllByTestId('basic-recording-device-tile');
    expect(tile('dev-1')).toHaveAttribute('aria-pressed', 'true');
    await openWide();
    fireEvent.click(screen.getByTestId('group-apply'));
    await waitFor(() => expect(tile('dev-2')).toHaveAttribute('aria-pressed', 'true'));
    await waitFor(() => expect(screen.getAllByTestId('circuit-profile-start')).toHaveLength(1));
    expect(startApi.start, 'nothing is started before the operator reviews and starts').not.toHaveBeenCalled();
  });
  it('freezes the routed group set under a new request identity when the operator starts it', async () => {
    renderHost([groupA(), wideB()]);
    await openWide();
    fireEvent.click(screen.getByTestId('group-apply'));
    await waitFor(() => expect(tile('dev-2')).toHaveAttribute('aria-pressed', 'true'));
    await confirmScopedStart();
    await waitFor(() => expect(startApi.start).toHaveBeenCalledTimes(1));
    const [request] = startedRequests();
    expect(request).toMatchObject({ workspace_id: 'ws-1', device_ids: ['dev-2'] });
    expect(request.request_id, 'a first intent is never given a borrowed identity').toBeTruthy();
    expect(request.groups.map((group) => group.group_id)).toEqual(['grp-w']);
    expect(api.apply).not.toHaveBeenCalled();
    expect(activateWorkspace).not.toHaveBeenCalled();
  });
  it('cannot apply on its own while a scoped start of that device is unresolved', async () => {
    renderHost([groupA(), groupB()]);
    startApi.start.mockRejectedValueOnce(lostResponse());
    await confirmScopedStart();
    await waitFor(() => expect(screen.getByTestId('recording-next-action')).toHaveTextContent('next.actions.recover_operation'));
    fireEvent.click(await screen.findByTestId('group-open-grp-a'));
    await screen.findByTestId('group-editor');
    expect(screen.getByTestId('group-apply')).toBeDisabled();
    fireEvent.click(screen.getByTestId('group-apply'));
    await act(async () => { await Promise.resolve(); });
    expect(api.apply).not.toHaveBeenCalled();
    expect(startApi.start, 'only the operator retry may reuse the pending request').toHaveBeenCalledTimes(1);
  });
  it('keeps the engineering Apply of a custom group', async () => {
    renderHost([groupA(), withSchema(savedGroup({ id: 'grp-c', revision: 'rev-grp-c-1', name: 'Raw tags' }))]);
    fireEvent.click(await screen.findByTestId('group-open-grp-c'));
    await screen.findByTestId('group-editor');
    await waitFor(() => expect(screen.getByTestId('group-apply')).toBeEnabled());
    fireEvent.click(screen.getByTestId('group-apply'));
    await waitFor(() => expect(api.apply).toHaveBeenCalledTimes(1));
    expect(api.apply.mock.calls[0][0]).toBe('grp-c');
    expect(api.apply.mock.calls[0][1]).toMatchObject({ workspace_id: 'ws-1', expected_group_revision: 'rev-grp-c-1', expected_connector_revision: 'crev-1' });
    expect(startApi.start).not.toHaveBeenCalled();
    expect(activateWorkspace).not.toHaveBeenCalled();
  });
  it('keeps the explicit disable of a managed profile group', async () => {
    renderHost([groupA(), wideB()]);
    api.disable.mockImplementation(async (id: string) => ({ workspace_revision: 'wrev-3', group: { ...(groups.find((group) => group.id === id) as WriteGroup), status: 'disabled' } }) as never);
    await openWide();
    const stop = within(screen.getByTestId('group-editor')).getByTestId('group-disable');
    await waitFor(() => expect(stop).toBeEnabled());
    fireEvent.click(stop);
    await waitFor(() => expect(api.disable).toHaveBeenCalledTimes(1));
    expect(api.disable.mock.calls[0][0]).toBe('grp-w');
    expect(api.disable.mock.calls[0][1]).toMatchObject({ expected_group_revision: 'rev-grp-w-1' });
    expect(startApi.start).not.toHaveBeenCalled();
    expect(activateWorkspace).not.toHaveBeenCalled();
  });
});
describe('Routing from Advanced keeps the operator\'s work', () => {
  const chooser = () => screen.getByTestId('circuit-profile-circuit-chooser') as HTMLSelectElement;
  const instantaneousInterval = () => screen.getByTestId('circuit-profile-interval-instantaneous') as HTMLInputElement;
  const openWide = async () => {
    fireEvent.click(await screen.findByTestId('group-open-grp-w'));
    await screen.findByTestId('group-editor');
    await waitFor(() => expect(screen.getByTestId('group-apply')).toBeEnabled());
  };
  it('shows the routed circuit of the same device and leaves another circuit\'s unsaved draft intact', async () => {
    renderHost([groupA(), wideB(), longB2()]);
    fireEvent.click(await screen.findByRole('button', { name: /PLC B/ }));
    await waitFor(() => expect(tile('dev-2')).toHaveAttribute('aria-pressed', 'true'));
    fireEvent.change(chooser(), { target: { value: 'cir-b2' } });
    await waitFor(() => expect(instantaneousInterval()).toBeEnabled());
    fireEvent.change(instantaneousInterval(), { target: { value: '20' } });
    await openWide();
    fireEvent.click(screen.getByTestId('group-apply'));
    await waitFor(() => expect(chooser().value).toBe('cir-w'));
    expect(startApi.start).not.toHaveBeenCalled();
    expect(api.apply).not.toHaveBeenCalled();
    fireEvent.change(chooser(), { target: { value: 'cir-b2' } });
    await waitFor(() => expect(instantaneousInterval().value, 'the unsaved draft of the other circuit is still there').toBe('20'));
  });
  it('does not route while the Advanced editor holds an unsaved draft, and keeps that draft', async () => {
    renderHost([groupA(), wideB()]);
    await openWide();
    fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Renamed but not saved' } });
    await waitFor(() => expect(screen.getByTestId('group-apply')).toBeDisabled());
    expect(screen.getByTestId('group-apply-blocked')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('group-apply'));
    await act(async () => { await Promise.resolve(); });
    expect(tile('dev-1'), 'nothing was routed').toHaveAttribute('aria-pressed', 'true');
    expect(screen.getByTestId('group-name')).toHaveValue('Renamed but not saved');
    expect(api.apply).not.toHaveBeenCalled();
  });
  it('applies a routing request once: a later remount of the circuit view keeps the operator\'s own circuit', async () => {
    renderHost([legacyA(), wideB(), longB2()]);
    await screen.findAllByTestId('basic-recording-device-tile');
    fireEvent.click(tile('dev-2'));
    await waitFor(() => expect(chooser()).toBeInTheDocument());
    await openWide();
    fireEvent.click(screen.getByTestId('group-apply'));
    await waitFor(() => expect(chooser().value).toBe('cir-w'));
    fireEvent.change(chooser(), { target: { value: 'cir-b2' } });
    await waitFor(() => expect(chooser().value).toBe('cir-b2'));
    await waitFor(() => expect(instantaneousInterval()).toBeEnabled());
    fireEvent.change(instantaneousInterval(), { target: { value: '20' } });
    fireEvent.click(tile('dev-1'));
    await waitFor(() => expect(screen.getByTestId('device-recording-existing')).toBeInTheDocument());
    expect(chooser().value).not.toBe('cir-b2');
    fireEvent.click(tile('dev-2'));
    await waitFor(() => expect(chooser()).toBeInTheDocument());
    expect(chooser().value, 'the old request must not choose the circuit again').toBe('cir-b2');
    expect(instantaneousInterval().value, 'the operator\'s unsaved draft is not replaced').toBe('20');
  });
  it('does not reapply a routing request after the operator chooses another circuit and device', async () => {
    const circuitA2 = (): WriteGroup => ({
      ...longProfileGroup({ id: 'grp-a2', circuitId: 'cir-a2', circuitName: 'Feeder A2', interval: 15, points: [2, 3] }),
    });
    renderHost([legacyA(), circuitA2(), groupA(), longB2()]);
    await screen.findAllByTestId('basic-recording-device-tile');
    expect(await screen.findByTestId('device-recording-existing'), 'device A keeps its earlier settings accessible in shared setup').toBeInTheDocument();
    fireEvent.click(await screen.findByTestId('group-open-grp-a'));
    await screen.findByTestId('group-editor');
    await waitFor(() => expect(screen.getByTestId('group-apply')).toBeEnabled());
    fireEvent.click(screen.getByTestId('group-apply'));
    await waitFor(() => expect(chooser().value).toBe('cir-a'));
    fireEvent.change(chooser(), { target: { value: 'cir-a2' } });
    await waitFor(() => expect(instantaneousInterval()).toBeEnabled());
    fireEvent.change(instantaneousInterval(), { target: { value: '24' } });
    fireEvent.click(tile('dev-2'));
    await waitFor(() => expect(chooser()).toBeInTheDocument());
    fireEvent.click(tile('dev-1'));
    await waitFor(() => expect(chooser()).toBeInTheDocument());
    expect(chooser().value, 'the request was left behind by the operator\'s own choices').toBe('cir-a2');
    expect(instantaneousInterval()).toHaveValue('24');
    expect(startApi.start).not.toHaveBeenCalled();
    expect(api.apply).not.toHaveBeenCalled();
  });
});
