import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import type { RecordingStartOperation } from '@/types/studioV2RecordingStart';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { WriteGroupDelivery } from '@/types/studioV2WriteGroupDelivery';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import { saveBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';

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
const deliveryQuery = {
  data: undefined as WriteGroupDelivery | undefined,
  isError: false,
};
const ensureAsync = vi.fn();
const startAsync = vi.fn();
const operationQuery = {
  data: undefined as RecordingStartOperation | undefined,
  isLoading: false,
  isError: false,
  refetch: vi.fn(),
};

vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupsQuery: () => groupsQuery,
  useWriteGroupDeliveryQuery: () => deliveryQuery,
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useEnsureBasicManagedMutation: () => ({ mutateAsync: ensureAsync, isPending: false }),
  useRecordingStartMutation: () => ({ mutateAsync: startAsync, isPending: false }),
  useRecordingStartOperationQuery: () => operationQuery,
}));

function renderPanel(props: Partial<React.ComponentProps<typeof BasicRecordingPanel>> = {}) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <BasicRecordingPanel state={savedGroupState()} workspaceId="ws-1" readonly={false} {...props} />
    </QueryClientProvider>,
  );
}

const canonicalManaged = savedGroup({
  basic_managed_device_id: 'dev-1',
  destination: {
    ...savedGroup().destination,
    storage_strategy: 'managed', table_name: 'gw_group_server', database: '/recording.db',
  },
  members: [
    ...savedGroup().members.map((member) => ({ ...member, target_column: '' })),
    { device_id: 'dev-1', point_id: 'pt-2', tag_id: 'tag-2', source_revision: 's2', mapping_revision: 'm2', target_column: '', required: true },
  ],
  row_policy: { interval_seconds: 60, allowed_lateness_seconds: 0 },
  write_policy: { mode: 'append', dedupe_capability: 'receipt' },
});

function operation(overrides: Partial<RecordingStartOperation> = {}): RecordingStartOperation {
  return {
    operation_id: 'operation-1', action: 'recording_start', status: 'succeeded', stage: 'complete',
    intent_digest: 'a'.repeat(64), workspace_id: 'ws-1', setup_revision: 'wrev-1', device_ids: ['dev-1'],
    groups: [{ group_id: canonicalManaged.id, group_revision: canonicalManaged.revision, saved: true, ready: true, applied: true }],
    devices: [{ device_id: 'dev-1', activated: true }],
    created_at: '2026-10-04T00:00:00Z', updated_at: '2026-10-04T00:00:01Z', ...overrides,
  };
}

function deliveryEvidence(overrides: Partial<WriteGroupDelivery> = {}): WriteGroupDelivery {
  return {
    group_id: canonicalManaged.id,
    intake: { state: 'active' },
    stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 3, skipped: 0 },
    last_sql_committed_at: '2026-10-04T00:00:10Z',
    oldest_pending_seconds: 0,
    no_data_buckets: 0,
    skipped_buckets: 0,
    backlog: [],
    quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
    ...overrides,
  };
}

beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] };
  groupsQuery.refetch.mockResolvedValue({ data: groupsQuery.data });
  deliveryQuery.data = undefined;
  deliveryQuery.isError = false;
  ensureAsync.mockResolvedValue({ workspace_revision: 'wrev-2', group: canonicalManaged });
  startAsync.mockResolvedValue(operation({ status: 'pending', stage: 'save' }));
  operationQuery.data = undefined;
});

function twoDeviceState() {
  const state = savedGroupState();
  const device = { ...state.devices[0], id: 'dev-2', name: 'PLC B' };
  const point = { ...state.points[0], id: 'rule-1-p-2', device_id: 'dev-2' };
  const mapping = {
    ...state.mappings[state.points[0].id], point_id: point.id, persisted_point_id: 'pt-b', tag_id: 'tag-b',
    device_id: 'dev-2', tag_key: 'line-b.temperature', display_name: 'line-b.temperature',
  };
  return {
    ...state,
    devices: [...state.devices, device], points: [...state.points, point],
    mappings: { ...state.mappings, [point.id]: mapping },
  };
}

describe('BasicRecordingPanel', () => {
  it('ensures a server-owned managed group, then starts one selected device', async () => {
    renderPanel();

    expect(await screen.findByTestId('basic-recording-selected-members')).toHaveTextContent('line.temperature');
    expect(screen.getByTestId('basic-recording-selected-members')).toHaveTextContent('40001');

    fireEvent.click(await screen.findByTestId('basic-recording-start'));

    await waitFor(() => expect(ensureAsync).toHaveBeenCalledTimes(1));
    expect(ensureAsync.mock.calls[0][0]).toMatchObject({ deviceId: 'dev-1' });
    expect(ensureAsync.mock.calls[0][0].request).toMatchObject({
      workspace_id: 'ws-1', expected_workspace_revision: 'wrev-1', expected_connector_revision: 'crev-1',
      group: { destination: { storage_strategy: 'managed', table_name: '' }, row_policy: { interval_seconds: 60 } },
    });
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(startAsync.mock.calls[0][0]).toMatchObject({
      request_id: expect.any(String), workspace_id: 'ws-1', expected_workspace_revision: 'wrev-2',
      device_ids: ['dev-1'], groups: [{ group_id: canonicalManaged.id, expected_group_revision: canonicalManaged.revision }],
    });
    expect(screen.getByTestId('basic-recording-status')).toHaveTextContent('step4.basic.operation.pending');
    expect(screen.getByTestId('basic-recording-retry')).toBeInTheDocument();
  });

  it('keeps existing custom groups out of the basic start request', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [savedGroup()] };
    renderPanel();
    expect(await screen.findByTestId('basic-recording-advanced-note')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('basic-recording-start'));
    await waitFor(() => expect(ensureAsync).toHaveBeenCalledTimes(1));
    expect(ensureAsync.mock.calls[0][0].request.group.destination.storage_strategy).toBe('managed');
  });

  it('blocks a missing destination without calling ensure or start', async () => {
    const state = savedGroupState({ db: {
      ...savedGroupState().db,
      connector: { ...savedGroupState().db.connector, persisted: false, save_state: 'idle', connector_id: undefined, identity_revision: undefined },
    } });
    renderPanel({ state });
    expect(await screen.findByTestId('basic-recording-blocked')).toHaveTextContent('step4.basic.blocked.destination');
    expect(screen.getByTestId('basic-recording-start')).toBeDisabled();
    expect(ensureAsync).not.toHaveBeenCalled();
    expect(startAsync).not.toHaveBeenCalled();
  });

  it('blocks a saved managed group pointing at another connector revision', async () => {
    groupsQuery.data = {
      ...groupsQuery.data,
      groups: [{ ...canonicalManaged, destination: { ...canonicalManaged.destination, connector_revision: 'old-crev' } }],
    };
    renderPanel();
    expect(await screen.findByTestId('basic-recording-blocked')).toHaveTextContent('step4.basic.blocked.destination_mismatch');
    expect(screen.getByTestId('basic-recording-start')).toBeDisabled();
    expect(ensureAsync).not.toHaveBeenCalled();
    expect(startAsync).not.toHaveBeenCalled();
  });

  it('reports partial progress and leaves retry explicit', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    saveBasicRecordingIntent({
      scope_key: 'ws-1/dev-1/managed-recording',
      request: {
        request_id: '33333333-3333-4333-8333-333333333333', workspace_id: 'ws-1',
        expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
        groups: [{ group_id: canonicalManaged.id, expected_group_revision: canonicalManaged.revision, expected_connector_revision: 'crev-1' }],
      }, operation_id: 'operation-1',
    });
    operationQuery.data = operation({ status: 'partial', stage: 'activation', reason: 'activation_failed', next_action: 'review_request' });
    renderPanel();
    expect(await screen.findByTestId('basic-recording-status')).toHaveTextContent('step4.basic.operation.partial');
    expect(screen.getByTestId('basic-recording-retry')).toBeInTheDocument();
    expect(startAsync).not.toHaveBeenCalled();
  });

  it('keeps retry available when the completed operation moved the group revision', async () => {
    const current = { ...canonicalManaged, revision: 'rev-applied', applied_revision: 'applied-1' };
    groupsQuery.data = { ...groupsQuery.data, groups: [current] };
    saveBasicRecordingIntent({
      scope_key: 'ws-1/dev-1/managed-recording',
      request: {
        request_id: '44444444-4444-4444-8444-444444444444', workspace_id: 'ws-1',
        expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
        groups: [{ group_id: current.id, expected_group_revision: 'rev-before', expected_connector_revision: 'crev-1' }],
      }, operation_id: 'operation-1',
    });
    operationQuery.data = operation({
      status: 'partial', stage: 'activation',
      groups: [{ group_id: current.id, group_revision: 'rev-applied', applied_revision: 'applied-1', saved: true, ready: true, applied: true }],
    });
    renderPanel();
    const retry = await screen.findByTestId('basic-recording-retry');
    expect(retry).not.toBeDisabled();
    fireEvent.click(retry);
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(startAsync.mock.calls[0][0].request_id).toBe('44444444-4444-4444-8444-444444444444');
    expect(startAsync.mock.calls[0][0].groups[0].expected_group_revision).toBe('rev-before');
  });

  it('does not confirm an older operation after the same group scope changes', async () => {
    const current = { ...canonicalManaged, revision: 'rev-new', applied_revision: 'applied-new' };
    groupsQuery.data = { ...groupsQuery.data, groups: [current] };
    saveBasicRecordingIntent({
      scope_key: 'ws-1/dev-1/managed-recording',
      request: {
        request_id: '77777777-7777-4777-8777-777777777777', workspace_id: 'ws-1',
        expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
        groups: [{ group_id: current.id, expected_group_revision: 'rev-old', expected_connector_revision: 'crev-1' }],
      }, operation_id: 'operation-1',
    });
    operationQuery.data = operation({
      status: 'succeeded', stage: 'complete',
      groups: [{ group_id: current.id, group_revision: 'rev-old', applied_revision: 'applied-old', saved: true, ready: true, applied: true }],
    });
    const onCommit = vi.fn();
    renderPanel({ onCommit });
    const status = await screen.findByTestId('basic-recording-status');
    expect(status).toHaveTextContent('step4.basic.operation.succeeded');
    expect(status).not.toHaveTextContent('step4.basic.verified_activation');
    expect(screen.queryByTestId('basic-recording-continue-runtime')).not.toBeInTheDocument();
    expect(onCommit).not.toHaveBeenCalled();
  });

  it('single-flights the ensure and start path on a double click', async () => {
    let release: (value: { workspace_revision: string; group: WriteGroup }) => void = () => undefined;
    ensureAsync.mockReturnValue(new Promise((resolve) => { release = resolve; }));
    renderPanel();
    const button = await screen.findByTestId('basic-recording-start');
    fireEvent.click(button);
    fireEvent.click(button);
    await waitFor(() => expect(ensureAsync).toHaveBeenCalledTimes(1));
    expect(startAsync).not.toHaveBeenCalled();
    release({ workspace_revision: 'wrev-2', group: canonicalManaged });
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
  });

  it('drops a late ensure response after switching to another device', async () => {
    let release: (value: { workspace_revision: string; group: WriteGroup }) => void = () => undefined;
    ensureAsync.mockReturnValue(new Promise((resolve) => { release = resolve; }));
    renderPanel({ state: twoDeviceState() });
    const select = await screen.findByTestId('basic-recording-device');
    fireEvent.click(screen.getByTestId('basic-recording-start'));
    await waitFor(() => expect(ensureAsync).toHaveBeenCalledTimes(1));
    fireEvent.change(select, { target: { value: 'dev-2' } });
    await waitFor(() => expect(select).toHaveValue('dev-2'));
    release({ workspace_revision: 'wrev-2', group: canonicalManaged });
    await new Promise((resolve) => setTimeout(resolve, 0));
    expect(startAsync).not.toHaveBeenCalled();
    expect(screen.queryByTestId('basic-recording-status')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('basic-recording-start'));
    await waitFor(() => expect(ensureAsync).toHaveBeenCalledTimes(2));
  });

  it('reuses the persisted request identity after a lost response', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    const request = {
      request_id: '11111111-1111-4111-8111-111111111111', workspace_id: 'ws-1',
      expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
      groups: [{ group_id: canonicalManaged.id, expected_group_revision: canonicalManaged.revision, expected_connector_revision: 'crev-1' }],
    };
    saveBasicRecordingIntent({ scope_key: 'ws-1/dev-1/managed-recording', request });
    startAsync.mockRejectedValueOnce(new Error('lost response'));
    renderPanel();
    fireEvent.click(await screen.findByTestId('basic-recording-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(startAsync.mock.calls[0][0].request_id).toBe(request.request_id);
    expect(ensureAsync).not.toHaveBeenCalled();
    expect(screen.getByTestId('basic-recording-error')).toBeInTheDocument();
  });

  it('resends the original request after a lost response even if the query revision advances', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    const request = {
      request_id: '55555555-5555-4555-8555-555555555555', workspace_id: 'ws-1',
      expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
      groups: [{ group_id: canonicalManaged.id, expected_group_revision: canonicalManaged.revision, expected_connector_revision: 'crev-1' }],
    };
    saveBasicRecordingIntent({ scope_key: 'ws-1/dev-1/managed-recording', request });
    startAsync.mockImplementationOnce(async () => {
      groupsQuery.data = { ...groupsQuery.data, workspace_revision: 'wrev-new' };
      throw new Error('lost response');
    }).mockResolvedValueOnce(operation({ status: 'pending', stage: 'save' }));
    renderPanel();
    const start = await screen.findByTestId('basic-recording-start');
    fireEvent.click(start);
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    fireEvent.click(start);
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(2));
    expect(startAsync.mock.calls[1][0]).toEqual(startAsync.mock.calls[0][0]);
  });

  it('reuses an unresolved request persisted before the group revision advanced', async () => {
    const old = { ...canonicalManaged, revision: 'rev-old' };
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    saveBasicRecordingIntent({
      scope_key: 'ws-1/dev-1/managed-recording',
      request: {
        request_id: '22222222-2222-4222-8222-222222222222', workspace_id: 'ws-1',
        expected_workspace_revision: 'wrev-old', device_ids: ['dev-1'],
        groups: [{ group_id: old.id, expected_group_revision: old.revision, expected_connector_revision: 'crev-1' }],
      },
    });
    renderPanel();
    fireEvent.click(await screen.findByTestId('basic-recording-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(startAsync.mock.calls[0][0].request_id).toBe('22222222-2222-4222-8222-222222222222');
  });

  it('sends basic interval/completeness changes as an explicit group draft', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    renderPanel({
      shareStatus: {
        workspace_id: 'ws-1', enabled: true, configured_enabled: true, port: 502, address: '127.0.0.1',
        bind_state: 'pass', mapping_count: 1, readiness_token: 'share-token', settings_revision: 'settings-1', workspace_revision: 'wrev-1',
      },
    });
    fireEvent.change(await screen.findByTestId('basic-recording-interval'), { target: { value: '30' } });
    fireEvent.change(screen.getByTestId('basic-recording-completeness'), { target: { value: 'partial' } });
    fireEvent.click(screen.getByTestId('basic-recording-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(ensureAsync).not.toHaveBeenCalled();
    expect(startAsync.mock.calls[0][0]).toMatchObject({ readiness_token: 'share-token', settings_revision: 'settings-1' });
    expect(startAsync.mock.calls[0][0].groups[0].draft.row_policy).toMatchObject({ interval_seconds: 30, incomplete_policy: 'partial' });
  });

  it('does not call a historical committed count current SQL evidence without a matching effect', async () => {
    const applied = { ...canonicalManaged, applied_revision: 'applied-1' };
    groupsQuery.data = { ...groupsQuery.data, groups: [applied] };
    deliveryQuery.data = deliveryEvidence({ group_id: applied.id });
    renderPanel();
    const evidence = await screen.findByTestId('basic-recording-evidence');
    expect(evidence).toHaveTextContent('step4.basic.unconfirmed');
    expect(screen.queryByTestId('basic-recording-committed-effect')).not.toBeInTheDocument();
  });

  it('shows only a receipt effect bound to the current applied revision', async () => {
    const applied = { ...canonicalManaged, applied_revision: 'applied-1' };
    groupsQuery.data = { ...groupsQuery.data, groups: [applied] };
    deliveryQuery.data = deliveryEvidence({
      group_id: applied.id,
      revision_stages: {
        group_revision: 'applied-1',
        stages: { collecting: 1, queued: 2, retrying: 1, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
      },
      last_sql_committed_effect: {
        group_revision: 'applied-1', connector_revision: 'crev-1', record_id: 'record-42',
        effect_key: 'effect-42', payload_digest: 'a'.repeat(64), committed_at: '2026-10-04T00:00:10Z',
      },
    });
    renderPanel();
    expect(await screen.findByTestId('basic-recording-committed-effect')).toHaveTextContent('group-1');
    expect(screen.getByTestId('basic-recording-committed-effect')).toHaveTextContent('record-42');
    expect(screen.getByTestId('basic-recording-committed-effect')).toHaveTextContent('effect-42');
    expect(screen.getByTestId('basic-recording-evidence')).toHaveTextContent('step4.basic.committed_verified');
    expect(screen.getByTestId('basic-recording-evidence')).toHaveTextContent('4');
  });

  it('does not use a receipt returned for another group', async () => {
    const applied = { ...canonicalManaged, applied_revision: 'applied-1' };
    groupsQuery.data = { ...groupsQuery.data, groups: [applied] };
    deliveryQuery.data = deliveryEvidence({
      group_id: 'group-other',
      last_sql_committed_effect: {
        group_revision: 'applied-1', connector_revision: 'crev-1', record_id: 'other-record',
        effect_key: 'other-effect', payload_digest: 'c'.repeat(64), committed_at: '2026-10-04T00:00:10Z',
      },
    });
    renderPanel();
    expect(await screen.findByTestId('basic-recording-evidence')).toHaveTextContent('step4.basic.unconfirmed');
    expect(screen.queryByTestId('basic-recording-committed-effect')).not.toBeInTheDocument();
  });

  it('keeps a receipt from an older applied revision unconfirmed', async () => {
    const applied = { ...canonicalManaged, applied_revision: 'applied-2' };
    groupsQuery.data = { ...groupsQuery.data, groups: [applied] };
    deliveryQuery.data = deliveryEvidence({
      group_id: applied.id,
      last_sql_committed_effect: {
        group_revision: 'applied-1', connector_revision: 'crev-1', record_id: 'old-record',
        effect_key: 'old-effect', payload_digest: 'b'.repeat(64), committed_at: '2026-10-04T00:00:10Z',
      },
    });
    renderPanel();
    expect(await screen.findByTestId('basic-recording-evidence')).toHaveTextContent('step4.basic.unconfirmed');
    expect(screen.queryByTestId('basic-recording-committed-effect')).not.toBeInTheDocument();
  });

  it('offers the existing Runtime handoff only after the selected device is activated', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    startAsync.mockResolvedValueOnce(operation({ status: 'succeeded', stage: 'complete' }));
    const onCommit = vi.fn();
    renderPanel({ onCommit });
    fireEvent.click(await screen.findByTestId('basic-recording-start'));
    const continueButton = await screen.findByTestId('basic-recording-continue-runtime');
    fireEvent.click(continueButton);
    expect(onCommit).toHaveBeenCalledWith(['dev-1']);
  });

  it('maps preparation reasons and actions to safe repair guidance', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    saveBasicRecordingIntent({
      scope_key: 'ws-1/dev-1/managed-recording',
      request: {
        request_id: '66666666-6666-4666-8666-666666666666', workspace_id: 'ws-1',
        expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
        groups: [{ group_id: canonicalManaged.id, expected_group_revision: canonicalManaged.revision, expected_connector_revision: 'crev-1' }],
      }, operation_id: 'operation-1',
    });
    operationQuery.data = operation({ status: 'failed', stage: 'readiness', reason: 'preparation_required', next_action: 'prepare_schema' });
    const onOpenAdvanced = vi.fn();
    renderPanel({ onOpenAdvanced });
    expect(await screen.findByTestId('basic-recording-status')).toHaveTextContent('step4.basic.reason.preparation_required');
    expect(screen.getByTestId('basic-recording-status')).toHaveTextContent('step4.basic.action.prepare_schema');
    fireEvent.click(screen.getByTestId('basic-recording-open-advanced'));
    expect(onOpenAdvanced).toHaveBeenCalledTimes(1);
  });

  it('maps Share readiness operation guidance without exposing backend text', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    saveBasicRecordingIntent({
      scope_key: 'ws-1/dev-1/managed-recording',
      request: {
        request_id: '88888888-8888-4888-8888-888888888888', workspace_id: 'ws-1',
        expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
        groups: [{ group_id: canonicalManaged.id, expected_group_revision: canonicalManaged.revision, expected_connector_revision: 'crev-1' }],
      }, operation_id: 'operation-1',
    });
    operationQuery.data = operation({ status: 'failed', stage: 'readiness', reason: 'share_not_ready', next_action: 'refresh_share' });
    renderPanel();
    const status = await screen.findByTestId('basic-recording-status');
    expect(status).toHaveTextContent('step4.basic.reason.share_not_ready');
    expect(status).toHaveTextContent('step4.basic.action.refresh_share');
  });

  it.each([
    ['RECORDING_START_INVALID', 'step4.basic.error.start_invalid'],
    ['RECORDING_START_INTENT_CHANGED', 'step4.basic.error.start_intent_changed'],
    ['RECORDING_START_NOT_FOUND', 'step4.basic.error.start_not_found'],
    ['RECORDING_START_BUSY', 'step4.basic.error.start_busy'],
    ['BASIC_INTENT_CHANGED', 'step4.basic.error.basic_intent_changed'],
    ['WRITE_GROUP_BASIC_INTENT_CHANGED', 'step4.basic.error.basic_intent_changed'],
  ])('maps recording-start code %s to safe local copy', async (code, messageKey) => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    startAsync.mockRejectedValueOnce({ response: { data: { error: { code, message: 'private backend detail' } } } });
    renderPanel();
    fireEvent.click(await screen.findByTestId('basic-recording-start'));
    const error = await screen.findByTestId('basic-recording-error');
    expect(error).toHaveTextContent(messageKey);
    expect(error).not.toHaveTextContent('private backend detail');
  });

  it('refreshes canonical groups and readiness evidence after a succeeded operation', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    startAsync.mockResolvedValueOnce(operation({ status: 'succeeded', stage: 'complete' }));
    renderPanel();
    fireEvent.click(await screen.findByTestId('basic-recording-start'));
    await waitFor(() => expect(screen.getByTestId('basic-recording-status')).toHaveTextContent('step4.basic.operation.succeeded'));
    await waitFor(() => expect(groupsQuery.refetch).toHaveBeenCalled());
  });

});
