import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { CircuitProfileBasicPanel } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import { GroupEditor } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import type { WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

const mocks = vi.hoisted(() => ({
  groupsQuery: {
    data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] as unknown[] },
    isLoading: false, isError: false, refetch: vi.fn(),
  },
  deliveryQuery: { data: undefined as unknown, isLoading: false, isError: false },
  readinessQuery: {
    data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', group_id: 'legacy-custom', group_revision: 'rev-1', applied_revision: '', config_ready: true, schema_ready: true, ready: true, issues: [] },
    isLoading: false, isFetching: false, isError: false, refetch: vi.fn(),
  },
  measurementsQuery: { data: [] as unknown[], isLoading: false, isError: false, refetch: vi.fn() },
  ensureProfile: vi.fn(), updateGroup: vi.fn(), createGroup: vi.fn(), ensureBasic: vi.fn(), start: vi.fn(),
  operationQuery: { data: undefined as unknown, isLoading: false, isError: false, refetch: vi.fn() },
  testWritePreview: vi.fn(), testWrite: vi.fn(),
  testWriteOperationQuery: { data: undefined as unknown, isLoading: false, isError: false, refetch: vi.fn() },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string, options?: Record<string, unknown>) => ({
    'step4.basic.circuit.mysql_capability.family.mysql': 'MySQL',
    'step4.basic.circuit.mysql_capability.family.mariadb': 'MariaDB',
    'step4.basic.circuit.mysql_capability.available': 'Available',
    'step4.basic.circuit.mysql_capability.unavailable': 'Unavailable',
    'step4.basic.circuit.mysql_capability.metadata_status.forbidden': 'Access denied',
  }[key] ?? (options && 'count' in options ? `${key}:${options.count}` : key)) }),
}));

vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupsQuery: () => mocks.groupsQuery,
  useWriteGroupDeliveryQuery: () => mocks.deliveryQuery,
  useWriteGroupReadinessQuery: () => mocks.readinessQuery,
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: mocks.ensureProfile, isPending: false }),
  useCreateWriteGroupMutation: () => ({ mutateAsync: mocks.createGroup, isPending: false }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: mocks.updateGroup, isPending: false }),
}));

vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => mocks.measurementsQuery,
  useCreateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useUpdateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useEnsureBasicManagedMutation: () => ({ mutateAsync: mocks.ensureBasic, isPending: false }),
  useRecordingStartMutation: () => ({ mutateAsync: mocks.start, isPending: false }),
  useRecordingStartOperationQuery: () => mocks.operationQuery,
}));
vi.mock('@/hooks/datalink/useStudioV2WriteGroupTestWrite', () => ({
  useWriteGroupTestWritePreviewMutation: () => ({ mutateAsync: mocks.testWritePreview, isPending: false }),
  useWriteGroupTestWriteMutation: () => ({ mutateAsync: mocks.testWrite, isPending: false }),
  useWriteGroupTestWriteOperationQuery: () => mocks.testWriteOperationQuery,
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }),
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel', () => ({ GroupSchemaPanel: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupDeliveryStrip', () => ({ GroupDeliveryStrip: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar', () => ({ GroupLifecycleBar: () => null }));

function stateForConnector(kind: WorkbenchV2State['db']['connector']['kind'], overrides: Partial<WorkbenchV2State['db']['connector']> = {}): WorkbenchV2State {
  const base = savedGroupState();
  return {
    ...base,
    db: { ...base.db, connector: { ...base.db.connector, kind, ...overrides } },
  };
}

function legacyCustomGroup(): WriteGroup {
  const base = savedGroup();
  return {
    ...base,
    id: 'legacy-custom', name: 'Plant table', revision: 'rev-1', layout_version: undefined,
    destination: { ...base.destination, connector_id: 'conn-1', connector_revision: 'crev-1', table_name: 'legacy_rows', storage_strategy: 'custom' },
    row_policy: { interval_seconds: 10, allowed_lateness_seconds: 0, incomplete_policy: 'partial', entity_key_column: 'entity', provenance_column: 'source' },
    members: [
      { ...base.members[0], entity_key: 'line-1' },
      { ...base.members[1], entity_key: 'line-2' },
    ],
  };
}

function nativeMonthlyGroup(): WriteGroup {
  const base = savedGroup();
  return {
    ...base,
    id: 'native-monthly', name: 'Main Circuit A (Instantaneous)', revision: 'native-rev-1', status: 'ready',
    layout_version: 'managed_long_v1', storage_mode: 'native_monthly',
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' },
    destination: { ...base.destination, connector_id: 'mysql-1', connector_revision: 'mysql-rev-1', storage_strategy: 'managed', table_name: 'gw_main_circuit' },
    row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    members: [{ ...base.members[0], device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', circuit_id: 'cir-main', circuit_name: 'Main Circuit A', target_column: 'temperature', required: false }],
  };
}

function managedLongGroup(): WriteGroup {
  const base = nativeMonthlyGroup();
  return {
    ...base,
    destination: {
      ...base.destination, connector_id: 'mysql-old', connector_revision: 'mysql-old-rev',
      table_schema: 'old_schema', table_name: 'old_canonical',
    },
  };
}

function capability(overrides: Record<string, unknown> = {}) {
  return {
    connector_id: 'mysql-1', connector_revision: 'mysql-rev-1', database: 'fixture',
    server_family: 'mysql', server_version: '8.4.4', capability_revision: 'cap-rev-1',
    connected: true, metadata_status: 'forbidden', transactional_engine: true,
    supports_managed_long: true, supports_native_partition: true,
    supports_group_test_write: true, supports_missing_repair: true, reasons: ['metadata_forbidden'],
    ...overrides,
  };
}

function apiResponse(value: unknown) {
  return { data: { success: true, data: value } } as never;
}

function testWritePreview(group: WriteGroup) {
  return {
    token: 'preview-token', operation_id: 'preview-operation', action: 'test_write' as const,
    group_id: group.id, group_revision: group.revision, expires_at: '2026-10-25T00:00:00Z',
    target: { connector_id: group.destination.connector_id, dialect: 'sqlite', database: group.destination.database, schema: '', table: group.destination.table_name },
    owner_column: 'entity', owner_value: 'line-1', dedupe: '', cleanup: 'remove the exact preview-owned row', values: [],
  };
}

function editorProps(group: WriteGroup, state: WorkbenchV2State, connector = state.db.connector) {
  const { candidates, excluded } = buildGroupCandidates(state);
  return {
    group, workspaceId: 'ws-1', workspaceRevision: 'wrev-1', connector, candidates, excluded, readonly: false,
    onClose: vi.fn(), onSaved: vi.fn(), onReload: vi.fn(),
    fetchLatest: vi.fn(async () => ({ workspace_revision: 'wrev-current', groups: [group] })),
  };
}

function queryProvider(children: React.ReactNode, client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })) {
  return <QueryClientProvider client={client}>{children}</QueryClientProvider>;
}

beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  mocks.groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] };
  mocks.groupsQuery.isLoading = false;
  mocks.groupsQuery.isError = false;
  mocks.groupsQuery.refetch.mockImplementation(async () => ({ data: mocks.groupsQuery.data }));
  mocks.deliveryQuery.data = undefined;
  mocks.readinessQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', group_id: 'legacy-custom', group_revision: 'rev-1', applied_revision: '', config_ready: true, schema_ready: true, ready: true, issues: [] };
  mocks.readinessQuery.isLoading = false;
  mocks.readinessQuery.isFetching = false;
  mocks.readinessQuery.isError = false;
  mocks.readinessQuery.refetch.mockResolvedValue({ data: mocks.readinessQuery.data });
  mocks.measurementsQuery.data = [];
  mocks.measurementsQuery.isLoading = false;
  mocks.measurementsQuery.isError = false;
  mocks.measurementsQuery.refetch.mockResolvedValue({ data: mocks.measurementsQuery.data });
  mocks.ensureProfile.mockResolvedValue({ workspace_revision: 'wrev-2', group: savedGroup() });
  mocks.updateGroup.mockResolvedValue({ workspace_revision: 'wrev-2', group: legacyCustomGroup() });
  mocks.createGroup.mockResolvedValue({ workspace_revision: 'wrev-2', group: legacyCustomGroup() });
  mocks.ensureBasic.mockReset();
  mocks.start.mockResolvedValue({
    operation_id: 'share-operation', action: 'recording_start', status: 'succeeded', stage: 'complete', intent_digest: 'd'.repeat(64),
    workspace_id: 'ws-1', setup_revision: '', device_ids: ['dev-1'], groups: [],
    devices: [{ device_id: 'dev-1', activated: true }], created_at: '2026-10-24T00:00:00Z', updated_at: '2026-10-24T00:00:01Z',
  });
  mocks.operationQuery.data = undefined;
  mocks.testWritePreview.mockReset();
  mocks.testWrite.mockReset();
  mocks.testWriteOperationQuery.data = undefined;
});
afterEach(() => vi.restoreAllMocks());

describe('MySQL/MariaDB legacy compatibility', () => {
  it.each(['sqlite', 'postgres'] as const)('%s keeps the existing Basic setup without a QueryClient or MySQL capability probe', async (kind) => {
    const get = vi.spyOn(studioV2DatalinkApi, 'get').mockRejectedValue(new Error('unexpected database probe'));
    const state = stateForConnector(kind);

    render(<CircuitProfileBasicPanel state={state} workspaceId="ws-1" readonly={false} />);

    expect(await screen.findByTestId('circuit-profile-basic-panel')).toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-storage-timescale')).toBeInTheDocument();
    expect(screen.queryByTestId('circuit-profile-mysql-capability')).not.toBeInTheDocument();
    expect(get).not.toHaveBeenCalled();
  });

  it.each(['sqlite', 'postgres'] as const)('keeps legacy custom table controls and preview available in the existing %s group editor', async (kind) => {
    const state = stateForConnector(kind);
    const group = legacyCustomGroup();
    mocks.readinessQuery.data = { ...mocks.readinessQuery.data, group_id: group.id, group_revision: group.revision };
    mocks.testWritePreview.mockResolvedValue(testWritePreview(group));

    render(<GroupEditor {...editorProps(group, state)} />);

    expect(screen.getByTestId('group-storage-custom')).toBeChecked();
    expect(screen.getByTestId('group-table')).toHaveValue('legacy_rows');
    expect(screen.getByTestId('group-interval')).toHaveValue(10);
    expect(screen.getByTestId('group-incomplete')).toHaveValue('partial');
    expect(screen.getByTestId('group-entity-column')).toHaveValue('entity');
    expect(screen.getByTestId('group-provenance-column')).toHaveValue('source');
    const preview = screen.getByTestId('group-test-write-preview');
    expect(preview).toBeEnabled();
    fireEvent.click(preview);
    await waitFor(() => expect(mocks.testWritePreview).toHaveBeenCalledWith(group.id));
    expect(await screen.findByTestId('group-test-write-preview-card')).toBeInTheDocument();
  });

  it('does not retarget a legacy group on connector change; explicit destination choice saves through the existing CAS fields', async () => {
    const sqliteState = stateForConnector('sqlite');
    const original = legacyCustomGroup();
    const mysqlState = stateForConnector('mysql', {
      name: 'Saved MySQL', connector_id: 'mysql-1', identity_revision: 'mysql-rev-1',
      database: 'fixture', schema: 'fixture', table: 'mysql_rows', persisted: true, save_state: 'saved',
    });
    const props = editorProps(original, sqliteState);
    const get = vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValueOnce(apiResponse(capability()));
    const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
    const view = render(queryProvider(<GroupEditor {...props} />, client));

    view.rerender(queryProvider(<GroupEditor {...props} connector={mysqlState.db.connector} />, client));
    expect(await screen.findByTestId('group-destination-mismatch')).toBeInTheDocument();
    expect(screen.getByTestId('group-storage-custom')).toBeChecked();
    expect(screen.getByTestId('group-table')).toHaveValue('legacy_rows');
    expect(mocks.updateGroup).not.toHaveBeenCalled();
    await waitFor(() => expect(get).toHaveBeenCalledWith('/studio-v2/workspace/recording-plans/capabilities', {
      params: { kind: 'mysql', connector_id: 'mysql-1', expected_connector_revision: 'mysql-rev-1' },
    }));

    fireEvent.click(screen.getByTestId('group-use-current-destination'));
    expect(screen.getByTestId('group-table')).toHaveValue('mysql_rows');
    expect(screen.getByTestId('group-save')).toBeEnabled();
    fireEvent.click(screen.getByTestId('group-save'));

    await waitFor(() => expect(mocks.updateGroup).toHaveBeenCalledTimes(1));
    const request = mocks.updateGroup.mock.calls[0][0] as { id: string; request: Record<string, unknown> };
    expect(request).toMatchObject({
      id: original.id,
      request: {
        expected_workspace_revision: 'wrev-current', expected_group_revision: original.revision,
        expected_connector_revision: 'mysql-rev-1',
        group: { destination: { connector_id: 'mysql-1', connector_revision: 'mysql-rev-1', table_name: 'mysql_rows', storage_strategy: 'custom' } },
      },
    });
  });

  it('retargets managed long groups with a backend-derived table while preserving their saved contract', async () => {
    const group = managedLongGroup();
    const oldState = stateForConnector('mysql', {
      name: 'Old MySQL', connector_id: 'mysql-old', identity_revision: 'mysql-old-rev',
      database: 'old_db', schema: 'old_schema', table: 'old_canonical', persisted: true, save_state: 'saved',
    });
    const newState = stateForConnector('mysql', {
      name: 'New MySQL', connector_id: 'mysql-new', identity_revision: 'mysql-new-rev',
      database: 'new_db', schema: 'new_schema', table: 'noncanonical_current_table', persisted: true, save_state: 'saved',
    });
    mocks.readinessQuery.data = { ...mocks.readinessQuery.data, group_id: group.id, group_revision: group.revision };
    const get = vi.spyOn(studioV2DatalinkApi, 'get')
      .mockResolvedValueOnce(apiResponse(capability({ connector_id: 'mysql-old', connector_revision: 'mysql-old-rev', metadata_status: 'verified', reasons: [] })))
      .mockResolvedValueOnce(apiResponse(capability({ connector_id: 'mysql-new', connector_revision: 'mysql-new-rev', metadata_status: 'verified', reasons: [] })));
    const props = editorProps(group, oldState);
    const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
    const tree = (connector: WorkbenchV2State['db']['connector']) => queryProvider(<GroupEditor {...props} connector={connector} />, client);
    const view = render(tree(oldState.db.connector));
    await waitFor(() => expect(get).toHaveBeenCalledTimes(1));

    view.rerender(tree(newState.db.connector));
    expect(await screen.findByTestId('group-destination-mismatch')).toBeInTheDocument();
    await waitFor(() => expect(get).toHaveBeenCalledTimes(2));
    fireEvent.click(screen.getByTestId('group-use-current-destination'));
    fireEvent.click(screen.getByTestId('group-save'));

    await waitFor(() => expect(mocks.updateGroup).toHaveBeenCalledTimes(1));
    const update = mocks.updateGroup.mock.calls[0][0] as { request: { expected_connector_revision: string; group: WriteGroup } };
    expect(update.request.expected_connector_revision).toBe('mysql-new-rev');
    expect(update.request.group.destination).toMatchObject({
      connector_id: 'mysql-new', connector_revision: 'mysql-new-rev',
      table_schema: 'new_schema', table_name: '', storage_strategy: 'managed',
    });
    expect(update.request.group).toMatchObject({
      layout_version: 'managed_long_v1', storage_mode: 'native_monthly', recording_profile: group.recording_profile,
    });
    expect(update.request.group.members).toHaveLength(group.members.length);
    expect(update.request.group.members[0]).toMatchObject({
      device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', circuit_id: 'cir-main',
      circuit_name: 'Main Circuit A', target_column: 'temperature', required: false,
    });
  });

  it('keeps Modbus-only start independent of a database group and carries the legacy token bridge', async () => {
    const get = vi.spyOn(studioV2DatalinkApi, 'get').mockRejectedValue(new Error('unexpected database probe'));
    const base = stateForConnector('mysql', {
      persisted: false, save_state: 'idle', connector_id: undefined, identity_revision: undefined,
    });
    const shareStatus = {
      workspace_id: 'ws-1', enabled: true, configured_enabled: true, port: 502, address: '127.0.0.1',
      bind_state: 'pass' as const, mapping_count: 1, readiness_token: 'share-token',
      settings_revision: 'settings-rev-1', workspace_revision: 'share-rev-1',
    };

    render(queryProvider(<BasicRecordingPanel state={base} workspaceId="ws-1" readonly={false} shareStatus={shareStatus} />));
    expect(await screen.findByTestId('basic-recording-destination')).toHaveTextContent('share_only_destination');
    fireEvent.click(screen.getByTestId('basic-recording-start'));

    await waitFor(() => expect(mocks.start).toHaveBeenCalledTimes(1));
    expect(mocks.start.mock.calls[0][0]).toMatchObject({
      workspace_id: 'ws-1', device_ids: ['dev-1'], groups: [], readiness_token: 'share-token',
      settings_revision: 'settings-rev-1', workspace_revision: 'share-rev-1',
    });
    expect(mocks.ensureBasic).not.toHaveBeenCalled();
    expect(mocks.ensureProfile).not.toHaveBeenCalled();
    expect(get).not.toHaveBeenCalled();
  });

  it('keeps a saved native-monthly group locked and blocks every operation when capability evidence is unverified', async () => {
    const group = nativeMonthlyGroup();
    const state = stateForConnector('mysql', {
      name: 'Saved MySQL', connector_id: 'mysql-1', identity_revision: 'mysql-rev-1',
      database: 'fixture', schema: 'fixture', persisted: true, save_state: 'saved',
    });
    mocks.groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [group] };
    mocks.readinessQuery.data = { ...mocks.readinessQuery.data, group_id: group.id, group_revision: group.revision, schema_ready: true, ready: false };
    const get = vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValueOnce(apiResponse(capability({ metadata_status: 'forbidden' })));

    render(queryProvider(<CircuitProfileBasicPanel state={state} workspaceId="ws-1" readonly={false} />));

    const nativeMode = await screen.findByTestId('circuit-profile-storage-native-monthly');
    await waitFor(() => expect(screen.getByTestId('circuit-profile-mysql-metadata')).toHaveTextContent('Access denied'));
    expect(nativeMode).toBeChecked();
    expect(nativeMode).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-storage-locked')).toBeInTheDocument();
    for (const id of [
      'circuit-profile-mysql-capability-managed-long', 'circuit-profile-mysql-capability-native-partitions',
      'circuit-profile-mysql-capability-test-write', 'circuit-profile-mysql-capability-missing-repair',
    ]) expect(screen.getByTestId(id)).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
    expect(mocks.testWritePreview).not.toHaveBeenCalled();
    expect(get).toHaveBeenCalledTimes(1);
  });
});
