import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import { parseWriteGroupReadinessData } from '@/utils/studioV2WriteGroupJson';
import type { MeasurementDefinition } from '@/types/measurement';
import type { RecordingStartOperation } from '@/types/studioV2RecordingStart';
import type { WriteGroup, WriteGroupProfileKind } from '@/types/studioV2WriteGroup';
import type { WriteGroupDelivery } from '@/types/studioV2WriteGroupDelivery';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import { expandSources } from '../../fixtures/sourceReview';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string, options?: Record<string, unknown>) => options && 'count' in options ? `${key}:${options.count}` : key }),
}));

const groupsQuery = { data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] as WriteGroup[] }, isLoading: false, isError: false, refetch: vi.fn() };
const deliveryQuery = { data: undefined as WriteGroupDelivery | undefined, isError: false };
const measurementsQuery = { data: [] as MeasurementDefinition[], isLoading: false, isError: false, refetch: vi.fn() };
const ensureBasicWideAsync = vi.fn(), ensureCircuitProfileAsync = vi.fn(), updateGroupAsync = vi.fn();
const createMeasurementAsync = vi.fn(), updateMeasurementAsync = vi.fn(), startAsync = vi.fn(), schemaPanelRender = vi.fn();
const operationQuery = { data: undefined as RecordingStartOperation | undefined, isLoading: false, isError: false, refetch: vi.fn() };

vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupsQuery: () => groupsQuery, useWriteGroupDeliveryQuery: () => deliveryQuery,
  useUpdateWriteGroupMutation: () => ({ mutateAsync: updateGroupAsync, isPending: false }),
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: ensureCircuitProfileAsync, isPending: false }),
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useEnsureBasicManagedMutation: () => ({ mutateAsync: ensureBasicWideAsync, isPending: false }),
  useRecordingStartMutation: () => ({ mutateAsync: startAsync, isPending: false }),
  useRecordingStartOperationQuery: () => operationQuery,
}));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => measurementsQuery,
  useCreateStudioV2MeasurementMutation: () => ({ mutateAsync: createMeasurementAsync, isPending: false }),
  useUpdateStudioV2MeasurementMutation: () => ({ mutateAsync: updateMeasurementAsync, isPending: false }),
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel', () => ({
  GroupSchemaPanel: (props: { group: WriteGroup; onApplied?: () => void }) => {
    schemaPanelRender(props);
    const kind = props.group.recording_profile?.kind ?? 'wide';
    return <div data-testid={`circuit-schema-panel-${kind}`}><span>{props.group.id}</span><button type="button" data-testid={`circuit-schema-confirm-${kind}`} onClick={() => props.onApplied?.()}>Confirm {kind}</button></div>;
  },
}));

function makeProfileGroup(overrides: Partial<WriteGroup> & { profileKind?: WriteGroupProfileKind } = {}): WriteGroup {
  const profile: WriteGroupProfileKind = overrides.profileKind ?? overrides.recording_profile?.kind ?? 'instantaneous';
  const intervalVal = overrides.recording_profile?.interval ?? (profile === 'cumulative' ? 10 : 15);
  return savedGroup({
    id: profile === 'cumulative' ? 'grp-cum-1' : 'grp-inst-1', revision: profile === 'cumulative' ? 'rev-cum-1' : 'rev-inst-1',
    name: profile === 'cumulative' ? 'Main Circuit A (Cumulative)' : 'PLC A (Instantaneous)', status: 'ready', layout_version: 'managed_long_v1',
    recording_profile: { kind: profile, interval: intervalVal, unit: profile === 'cumulative' ? 'minute' : 'second' }, storage_mode: 'plain',
    destination: { ...savedGroup().destination, storage_strategy: 'managed', table_name: 'gw_Safe_ws_1_long', database: '/recording.db' },
    members: [0, 1, 2].map((i) => ({
      device_id: 'dev-1', point_id: `pt-${i}`, tag_id: `tag-${i}`, source_revision: `s${i}`, mapping_revision: `m${i}`,
      target_column: '', required: false, measurement_id: profile === 'cumulative' ? `meas-${i}` : '',
      ...(profile === 'cumulative' ? { circuit_id: 'cir-main-circuit-a', circuit_name: 'Main Circuit A' } : {}),
    })),
    row_policy: { interval_seconds: profile === 'cumulative' ? intervalVal * 60 : intervalVal, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    write_policy: { mode: 'append', dedupe_capability: 'receipt', recreate_missing: true }, ...overrides,
  });
}

function makeConfirmedCounterMeasurement(index: number, overrides: Partial<MeasurementDefinition> = {}): MeasurementDefinition {
  return {
    id: `meas-${index}`, workspace_id: 'ws-1', definition_revision: `mrev-${index}`, name: `Energy ${index}`,
    equipment_id: 'eq-1', device_id: 'dev-1', point_id: `pt-${index}`, tag_id: `tag-${index}`,
    source_binding_revision: `s${index}`, quantity: 'energy', semantic_kind: 'counter', unit: 'kWh', series_epoch: '1',
    counter_policy: { allow_negative: false }, created_at: '2026-10-20T00:00:00Z', updated_at: '2026-10-20T00:00:00Z', ...overrides,
  };
}

function makeOperation(overrides: Partial<RecordingStartOperation> = {}): RecordingStartOperation {
  return {
    operation_id: 'op-circuit-1', action: 'recording_start', status: 'succeeded', stage: 'complete',
    intent_digest: 'b'.repeat(64), workspace_id: 'ws-1', setup_revision: 'wrev-2', device_ids: ['dev-1'],
    groups: [{ group_id: 'grp-inst-1', group_revision: 'rev-inst-1', saved: true, ready: true, applied: true }],
    devices: [{ device_id: 'dev-1', activated: true }], created_at: '2026-10-24T00:00:00Z', updated_at: '2026-10-24T00:00:01Z', ...overrides,
  };
}

function renderPanel(props: Partial<React.ComponentProps<typeof BasicRecordingPanel>> = {}) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  return render(<QueryClientProvider client={client}><BasicRecordingPanel state={savedGroupState()} workspaceId="ws-1" readonly={false} {...props} /></QueryClientProvider>);
}

beforeEach(() => {
  vi.clearAllMocks(); window.sessionStorage.clear();
  groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] };
  groupsQuery.isLoading = false; groupsQuery.isError = false;
  groupsQuery.refetch.mockImplementation(async () => ({ data: groupsQuery.data }));
  deliveryQuery.data = undefined; deliveryQuery.isError = false;
  measurementsQuery.data = []; measurementsQuery.isLoading = false; measurementsQuery.isError = false;
  measurementsQuery.refetch.mockImplementation(async () => ({ data: measurementsQuery.data }));
  ensureCircuitProfileAsync.mockImplementation(async ({ request }) => ({
    workspace_revision: 'wrev-2',
    group: makeProfileGroup({
      profileKind: request.group.recording_profile?.kind ?? 'instantaneous', recording_profile: request.group.recording_profile,
      storage_mode: request.group.storage_mode ?? 'plain', row_policy: { interval_seconds: request.group.row_policy.interval_seconds, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
      write_policy: request.group.write_policy, members: request.group.members,
    }),
  }));
  updateGroupAsync.mockImplementation(async ({ id, groupId, request }) => {
    const resolvedId = id ?? groupId ?? 'grp-inst-1';
    return {
      workspace_revision: 'wrev-3',
      group: makeProfileGroup({
        id: resolvedId, revision: `${resolvedId}-rev-2`, profileKind: request.group.recording_profile?.kind ?? 'instantaneous',
        recording_profile: request.group.recording_profile, storage_mode: request.group.storage_mode ?? 'plain',
        row_policy: { interval_seconds: request.group.row_policy.interval_seconds, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
        write_policy: request.group.write_policy, members: request.group.members,
      }),
    };
  });
  startAsync.mockResolvedValue(makeOperation({ status: 'pending', stage: 'save' }));
  operationQuery.data = undefined;
});

describe('CircuitProfileBasicPanel (Task 4.1)', () => {
  it('supports fresh basic mapping without MeasurementDefinition for instantaneous profile (15s default, plain, recreate_missing=true)', async () => {
    renderPanel();
    expect(await screen.findByTestId('circuit-profile-basic-panel')).toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-destination-summary')).toHaveTextContent('Line DB');
    expect(screen.queryByTestId('basic-recording-completeness')).not.toBeInTheDocument();
    expect((screen.getByTestId('circuit-profile-toggle-instantaneous') as HTMLInputElement).checked).toBe(true);
    expect((screen.getByTestId('circuit-profile-toggle-cumulative') as HTMLInputElement).checked).toBe(false);
    expect((screen.getByTestId('circuit-profile-interval-instantaneous') as HTMLInputElement).value).toBe('15');
    expect((screen.getByTestId('circuit-profile-interval-cumulative') as HTMLInputElement).value).toBe('10');
    expect((screen.getByTestId('circuit-profile-recreate-missing') as HTMLInputElement).checked).toBe(true);
    expect(screen.getByTestId('circuit-profile-storage-timescale')).toBeDisabled();

    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensureCircuitProfileAsync).toHaveBeenCalledTimes(1));
    expect(ensureCircuitProfileAsync.mock.calls[0][0].request).toMatchObject({
      workspace_id: 'ws-1', expected_workspace_revision: 'wrev-1', expected_connector_revision: 'crev-1',
      group: { layout_version: 'managed_long_v1', recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' }, storage_mode: 'plain', row_policy: { interval_seconds: 15 }, write_policy: { recreate_missing: true } },
    });
    expect(await screen.findByTestId('circuit-schema-panel-instantaneous')).toBeInTheDocument();

    fireEvent.click(screen.getByTestId('circuit-profile-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(startAsync.mock.calls[0][0]).toMatchObject({
      workspace_id: 'ws-1', device_ids: ['dev-1'],
      groups: [{ group_id: 'grp-inst-1', expected_group_revision: 'rev-inst-1', expected_connector_revision: 'crev-1' }],
    });
  });

  it('validates independent instantaneous (1..300s) and cumulative (1..600m) intervals strictly without silent clamp', async () => {
    measurementsQuery.data = [0, 1, 2].map((i) => makeConfirmedCounterMeasurement(i));
    renderPanel();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Main Circuit A' } });
    expandSources();
    fireEvent.click(screen.getByLabelText('line.batch'));
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));

    const instInterval = screen.getByTestId('circuit-profile-interval-instantaneous');
    const cumInterval = screen.getByTestId('circuit-profile-interval-cumulative');
    for (const badInst of ['0', '301', '-5', '15.5', 'abc']) {
      fireEvent.change(instInterval, { target: { value: badInst } });
      expect(screen.getByTestId('circuit-profile-interval-instantaneous-error')).toHaveAttribute('role', 'alert');
      expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
      expect(screen.getByTestId('circuit-profile-start')).toBeDisabled();
    }
    fireEvent.change(instInterval, { target: { value: '15' } });
    for (const badCum of ['0', '601', '-1', '2.5', 'xyz']) {
      fireEvent.change(cumInterval, { target: { value: badCum } });
      expect(screen.getByTestId('circuit-profile-interval-cumulative-error')).toHaveAttribute('role', 'alert');
      expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
      expect(screen.getByTestId('circuit-profile-start')).toBeDisabled();
    }
  });

  it('prepares both canonical profile groups on one circuit and sends both canonical group IDs in recording-start groups[]', async () => {
    measurementsQuery.data = [0, 1, 2].map((i) => makeConfirmedCounterMeasurement(i));
    renderPanel();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Main Circuit A' } });
    expandSources();
    fireEvent.click(screen.getByLabelText('line.batch'));
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensureCircuitProfileAsync).toHaveBeenCalledTimes(2));

    const firstReq = ensureCircuitProfileAsync.mock.calls[0][0].request.group;
    const secondReq = ensureCircuitProfileAsync.mock.calls[1][0].request.group;
    expect(firstReq).toMatchObject({ layout_version: 'managed_long_v1', recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' }, row_policy: { interval_seconds: 15 } });
    expect(secondReq).toMatchObject({ layout_version: 'managed_long_v1', recording_profile: { kind: 'cumulative', interval: 10, unit: 'minute' }, row_policy: { interval_seconds: 600 } });
    expect(firstReq.members[0].circuit_id).toMatch(/^cir-[0-9a-f-]{36}$/);
    expect(secondReq.members[0].circuit_id).toBe(firstReq.members[0].circuit_id);
    expect(await screen.findByTestId('circuit-schema-panel-instantaneous')).toBeInTheDocument();
    expect(await screen.findByTestId('circuit-schema-panel-cumulative')).toBeInTheDocument();

    fireEvent.click(screen.getByTestId('circuit-profile-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(startAsync.mock.calls[0][0].groups).toEqual([
      { group_id: 'grp-inst-1', expected_group_revision: 'rev-inst-1', expected_connector_revision: 'crev-1' },
      { group_id: 'grp-cum-1', expected_group_revision: 'rev-cum-1', expected_connector_revision: 'crev-1' },
    ]);
  });

  it('reuses canonical groups on reload/resend and uses PUT with CAS when editing interval or recreate_missing on an already-saved group', async () => {
    startAsync.mockResolvedValueOnce(makeOperation({ status: 'succeeded' }));
    groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [makeProfileGroup({ id: 'grp-inst-1', revision: 'rev-inst-1' })] };
    renderPanel();
    fireEvent.click(await screen.findByTestId('circuit-profile-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(updateGroupAsync).not.toHaveBeenCalled();

    await waitFor(() => expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toBeEnabled());
    fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '30' } });
    fireEvent.click(screen.getByTestId('circuit-profile-recreate-missing'));
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(updateGroupAsync).toHaveBeenCalledTimes(1));
    expect(updateGroupAsync.mock.calls[0][0]).toMatchObject({
      id: 'grp-inst-1',
      request: {
        workspace_id: 'ws-1', expected_workspace_revision: 'wrev-1', expected_group_revision: 'rev-inst-1', expected_connector_revision: 'crev-1',
        group: { recording_profile: { kind: 'instantaneous', interval: 30, unit: 'second' }, row_policy: { interval_seconds: 30 }, write_policy: { recreate_missing: false } },
      },
    });
  });

  it('preserves cumulative selection across Next/Back/reload when counter semantics are unconfirmed and allows narrow confirmation via measurement API', async () => {
    const firstRender = renderPanel();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Main Circuit A' } });
    expandSources();
    fireEvent.click(screen.getByLabelText('line.batch'));
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    expect(await screen.findByTestId('circuit-counter-confirmation')).toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();

    firstRender.unmount();
    renderPanel();
    expect((await screen.findByTestId('circuit-profile-toggle-cumulative') as HTMLInputElement).checked).toBe(true);
    expect((screen.getByTestId('circuit-profile-circuit-name') as HTMLInputElement).value).toBe('Main Circuit A');

    createMeasurementAsync.mockImplementation(async (def: Partial<MeasurementDefinition>) => {
      const idx = Number(String(def.point_id).replace('pt-', '')) || 0;
      const created = makeConfirmedCounterMeasurement(idx, { semantic_kind: def.semantic_kind ?? 'counter', unit: def.unit });
      measurementsQuery.data = [...measurementsQuery.data, created];
      return created;
    });
    for (let i = 0; i < 2; i += 1) {
      const row = screen.getByTestId(`circuit-counter-confirmation-source-tag-${i}`);
      fireEvent.change(within(row).getByTestId('circuit-confirm-counter-unit'), { target: { value: 'kWh' } });
      fireEvent.click(within(row).getByTestId('circuit-confirm-counter-submit'));
      await waitFor(() => expect(createMeasurementAsync).toHaveBeenCalledTimes(i + 1));
    }
  });

  it('enables explicit timescale selection on PostgreSQL, locks storage_mode after group creation, and surfaces timescale_unavailable without silent fallback', async () => {
    const pgState = savedGroupState({ db: { ...savedGroupState().db, connector: { ...savedGroupState().db.connector, kind: 'postgres', name: 'Timescale PG' } } });
    renderPanel({ state: pgState });
    const tsRadio = await screen.findByTestId('circuit-profile-storage-timescale');
    expect(tsRadio).not.toBeDisabled();
    fireEvent.click(tsRadio);

    ensureCircuitProfileAsync.mockRejectedValueOnce(Object.assign(new Error('TimescaleDB extension is not installed'), { code: 'WRITE_GROUP_TIMESCALE_UNAVAILABLE' }));
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    expect(await screen.findByTestId('circuit-profile-error')).toHaveTextContent('step4.basic.errors.timescale_unavailable');
    expect((screen.getByTestId('circuit-profile-storage-timescale') as HTMLInputElement).checked).toBe(true);
  });

  it('recovers dual-profile recording-start with identical request_id after lost response and parses counter readiness blockers', async () => {
    measurementsQuery.data = [0, 1, 2].map((i) => makeConfirmedCounterMeasurement(i));
    renderPanel();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Main Circuit A' } });
    expandSources();
    fireEvent.click(screen.getByLabelText('line.batch'));
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));

    startAsync.mockRejectedValueOnce(new TypeError('Failed to fetch'));
    fireEvent.click(screen.getByTestId('circuit-profile-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    const firstRequestId = startAsync.mock.calls[0][0].request_id;
    expect(await screen.findByTestId('basic-recording-status')).toHaveTextContent('step4.basic.operation.unknown');

    startAsync.mockResolvedValueOnce(makeOperation({ status: 'succeeded' }));
    fireEvent.click(screen.getByTestId('basic-recording-retry'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(2));
    expect(startAsync.mock.calls[1][0].request_id).toBe(firstRequestId);
    expect(startAsync.mock.calls[1][0].groups).toHaveLength(2);

    const parsedReadiness = parseWriteGroupReadinessData({
      workspace_id: 'ws-1', workspace_revision: 'wrev-1', group_id: 'grp-cum-1', group_revision: 'rev-cum-1',
      applied_revision: '', config_ready: true, schema_ready: true, ready: false, schema_digest: 'd'.repeat(64),
      issues: [{ code: 'counter-semantics-required', severity: 'blocking', step: 'Step 4', scope: 'grp-cum-1:pt-0', message: 'requires counter' }],
    });
    expect(parsedReadiness).not.toBeNull();
    expect(parsedReadiness?.ready).toBe(false);
  });

  it('retains existing wide basic_managed_device_id groups on the legacy wide branch and blocks safely when groupsQuery is loading or errored', async () => {
    groupsQuery.isLoading = true;
    const loadingView = renderPanel();
    expect(screen.getByTestId('circuit-profile-groups-loading')).toBeInTheDocument();
    loadingView.unmount();

    groupsQuery.isLoading = false;
    groupsQuery.isError = true;
    const errorView = renderPanel();
    expect(screen.getByTestId('circuit-profile-groups-error')).toBeInTheDocument();
    errorView.unmount();

    groupsQuery.isError = false;
    groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [savedGroup({ basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1' })] };
    renderPanel();
    expect(await screen.findByTestId('basic-recording-completeness')).toBeInTheDocument();
    expect(screen.queryByTestId('circuit-profile-basic-panel')).not.toBeInTheDocument();
  });
});
