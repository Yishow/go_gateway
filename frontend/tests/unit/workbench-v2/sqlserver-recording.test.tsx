import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { CircuitProfileBasicPanel } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel';
import { CircuitProfileSchemaPreparation } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileSchemaPreparation';
import { defaultCircuitProfileLocalDraft } from '@/features/datalink/workbench-v2/state/circuitProfileBasic';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import type { WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { parseStudioV2WriteGroup } from '@/utils/studioV2WriteGroupJson';

const mocks = vi.hoisted(() => ({
  groupsQuery: {
    data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] as unknown[] },
    isLoading: false,
    isError: false,
    refetch: vi.fn(),
  },
  measurementsQuery: { data: [] as unknown[], isLoading: false, isError: false, refetch: vi.fn() },
  ensureCircuitProfile: vi.fn(),
  updateGroup: vi.fn(),
  start: vi.fn(),
  operationQuery: { data: undefined as unknown, isLoading: false, isError: false, refetch: vi.fn() },
  readinessQuery: {
    data: { group_id: 'group-native', group_revision: 'rev-1', applied_revision: '', config_ready: false, schema_ready: true, ready: false, issues: [] },
    isLoading: false, isError: false, refetch: vi.fn(),
  },
  testWritePreview: vi.fn(),
  testWrite: vi.fn(),
  testWriteOperationQuery: { data: undefined as unknown, isLoading: false, isError: false, refetch: vi.fn() },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => ({
      'step4.basic.circuit.sqlserver_capability.title': 'SQL Server Recording Capability',
      'step4.basic.circuit.sqlserver_capability.available': 'Available',
      'step4.basic.circuit.sqlserver_capability.unavailable': 'Unavailable',
      'step4.basic.circuit.sqlserver_capability.metadata_status.verified': 'Verified',
      'step4.basic.circuit.sqlserver_capability.metadata_status.forbidden': 'Access denied',
      'step4.basic.circuit.sqlserver_capability.load_error': 'Could not check SQL Server capabilities. Try again.',
      'step4.group.test_write.blocked_capability': 'Capability unavailable for test write',
      'step4.basic.circuit.storage_native_monthly_sqlserver': 'Native monthly partitions (SQL Server)',
      'step4.basic.circuit.storage_native_monthly': 'Native monthly partitions (MySQL / MariaDB)',
    }[key] ?? key),
  }),
}));

vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupsQuery: () => mocks.groupsQuery,
  useWriteGroupDeliveryQuery: () => ({ data: undefined, isError: false }),
  useWriteGroupReadinessQuery: () => mocks.readinessQuery,
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: mocks.ensureCircuitProfile, isPending: false }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: mocks.updateGroup, isPending: false }),
}));

vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => mocks.measurementsQuery,
}));

vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useRecordingStartMutation: () => ({ mutateAsync: mocks.start, isPending: false }),
  useRecordingStartOperationQuery: () => mocks.operationQuery,
}));

vi.mock('@/hooks/datalink/useStudioV2WriteGroupTestWrite', () => ({
  useWriteGroupTestWritePreviewMutation: () => ({ mutateAsync: mocks.testWritePreview, isPending: false }),
  useWriteGroupTestWriteMutation: () => ({ mutateAsync: mocks.testWrite, isPending: false }),
  useWriteGroupTestWriteOperationQuery: () => mocks.testWriteOperationQuery,
}));

vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel', () => ({
  GroupSchemaPanel: () => null,
}));

type SQLServerCapabilityFixture = {
  connector_id: string;
  connector_revision: string;
  database: string;
  schema: string;
  server_version: string;
  server_edition: string;
  engine_edition: number;
  compatibility_level: number;
  connected: boolean;
  metadata_status: string;
  permissions: {
    verified: boolean;
    database_view_definition: boolean;
    schema_view_definition: boolean;
    create_table: boolean;
    alter_schema: boolean;
    select_schema: boolean;
    insert_schema: boolean;
    delete_schema: boolean;
    alter_any_dataspace: boolean;
  };
  supports_managed_long: boolean;
  supports_native_partition: boolean;
  supports_group_test_write: boolean;
  supports_missing_repair: boolean;
  capability_revision: string;
  reasons: string[];
};

function sqlserverCapability(overrides: Partial<SQLServerCapabilityFixture> = {}): SQLServerCapabilityFixture {
  return {
    connector_id: 'sqlserver-1',
    connector_revision: 'sqlserver-rev-1',
    database: 'gw_circuit_test',
    schema: 'dbo',
    server_version: '16.0.4165.4',
    server_edition: 'Developer Edition (64-bit)',
    engine_edition: 3,
    compatibility_level: 160,
    connected: true,
    metadata_status: 'verified',
    permissions: {
      verified: true,
      database_view_definition: true,
      schema_view_definition: true,
      create_table: true,
      alter_schema: true,
      select_schema: true,
      insert_schema: true,
      delete_schema: true,
      alter_any_dataspace: false,
    },
    supports_managed_long: false,
    supports_native_partition: false,
    supports_group_test_write: false,
    supports_missing_repair: false,
    capability_revision: 'sql-cap-rev-1',
    reasons: ['adapter_unverified', 'managed_long_unverified'],
    ...overrides,
  };
}

function nativeProfileGroup(): WriteGroup {
  const base = savedGroup({
    id: 'group-native',
    name: 'Main Circuit A (Instantaneous)',
    status: 'ready',
    layout_version: 'managed_long_v1',
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' },
    destination: {
      ...savedGroup().destination,
      connector_id: 'sqlserver-1',
      connector_revision: 'sqlserver-rev-1',
      storage_strategy: 'managed',
    },
    row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    members: [{
      device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', source_revision: 's0', mapping_revision: 'm0',
      target_column: 'value', required: false, circuit_id: 'cir-main', circuit_name: 'Main Circuit A',
    }],
  });
  return { ...base, storage_mode: 'native_monthly' } as unknown as WriteGroup;
}

function sqlserverState(connectorId = 'sqlserver-1', connectorRevision = 'sqlserver-rev-1'): WorkbenchV2State {
  const state = savedGroupState();
  return {
    ...state,
    db: {
      ...state.db,
      connector: {
        ...state.db.connector,
        kind: 'sqlserver',
        name: 'Saved SQL Server',
        connector_id: connectorId,
        identity_revision: connectorRevision,
        database: 'gw_circuit_test',
        schema: 'dbo',
        persisted: true,
        save_state: 'saved',
      },
    },
  };
}

function renderPanel(state: WorkbenchV2State) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const view = render(
    <QueryClientProvider client={client}>
      <CircuitProfileBasicPanel state={state} workspaceId="ws-1" readonly={false} />
    </QueryClientProvider>,
  );
  return { ...view, client };
}

function renderPreparedSchema(options: {
  isSQLServer?: boolean;
  sqlserverTestWriteAvailable?: boolean;
  readonly?: boolean;
  dirty?: boolean;
  destinationMatches?: boolean;
  schemaReady?: boolean;
} = {}) {
  mocks.readinessQuery.data = { ...mocks.readinessQuery.data, schema_ready: options.schemaReady ?? true };
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <CircuitProfileSchemaPreparation
        workspaceId="ws-1"
        workspaceRevision="wrev-1"
        preparedGroups={[nativeProfileGroup()]}
        dirtyGroupIds={options.dirty ? new Set(['group-native']) : new Set()}
        disabled={false}
        readonly={options.readonly ?? false}
        isPreparing={false}
        destinationMatches={options.destinationMatches ?? true}
        onPrepare={() => undefined}
        onApplied={() => undefined}
        isSQLServer={options.isSQLServer ?? true}
        sqlserverTestWriteAvailable={options.sqlserverTestWriteAvailable ?? false}
        sqlserverTestWriteCapabilityRevision="sql-cap-rev-1"
      />
    </QueryClientProvider>,
  );
}

function apiResponse(value: SQLServerCapabilityFixture) {
  return { data: { success: true, data: value } } as never;
}

beforeEach(() => {
  window.sessionStorage.clear();
  mocks.groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] };
  mocks.groupsQuery.isLoading = false;
  mocks.groupsQuery.isError = false;
  mocks.groupsQuery.refetch.mockImplementation(async () => ({ data: mocks.groupsQuery.data }));
  mocks.measurementsQuery.data = [];
  mocks.measurementsQuery.isLoading = false;
  mocks.measurementsQuery.isError = false;
  mocks.measurementsQuery.refetch.mockImplementation(async () => ({ data: mocks.measurementsQuery.data }));
  mocks.ensureCircuitProfile.mockResolvedValue({
    workspace_revision: 'wrev-2',
    group: savedGroup({
      id: 'new-profile', status: 'ready', layout_version: 'managed_long_v1',
      recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' }, storage_mode: 'plain',
      destination: { ...savedGroup().destination, storage_strategy: 'managed' },
      row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    }),
  });
  mocks.updateGroup.mockReset();
  mocks.start.mockReset();
  mocks.operationQuery.data = undefined;
  mocks.readinessQuery.data = {
    group_id: 'group-native', group_revision: 'rev-1', applied_revision: '', config_ready: false, schema_ready: true, ready: false, issues: [],
  };
  mocks.readinessQuery.isLoading = false;
  mocks.readinessQuery.isError = false;
  mocks.testWritePreview.mockReset();
  mocks.testWrite.mockReset();
  mocks.testWriteOperationQuery.data = undefined;
});

afterEach(() => vi.restoreAllMocks());

describe('saved SQL Server circuit recording', () => {
  it('parses and restores a saved native-monthly mode without downgrading it to plain', () => {
    const saved = nativeProfileGroup();
    const parsed = parseStudioV2WriteGroup({ ...saved, storage_mode: 'native_monthly' });
    expect(parsed?.storage_mode).toBe('native_monthly');
    expect(defaultCircuitProfileLocalDraft([], [saved], 'cir-main').storage_mode).toBe('native_monthly');
  });

  it('loads the exact saved SQL Server connector revision and shows database, schema, version, edition, and operation capabilities', async () => {
    const get = vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValueOnce(apiResponse(sqlserverCapability()));
    renderPanel(sqlserverState());

    expect(await screen.findByTestId('circuit-profile-sqlserver-capability')).toBeInTheDocument();
    await waitFor(() => expect(screen.getByTestId('circuit-profile-sqlserver-server-version')).toHaveTextContent('16.0.4165.4'));
    expect(screen.getByTestId('circuit-profile-sqlserver-server-edition')).toHaveTextContent('Developer Edition (64-bit)');
    expect(screen.getByTestId('circuit-profile-sqlserver-database')).toHaveTextContent('gw_circuit_test');
    expect(screen.getByTestId('circuit-profile-sqlserver-schema')).toHaveTextContent('dbo');
    expect(screen.getByTestId('circuit-profile-sqlserver-capability-managed-long')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-sqlserver-capability-native-partitions')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-sqlserver-capability-test-write')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-sqlserver-capability-missing-repair')).toHaveTextContent('Unavailable');
    expect(get).toHaveBeenCalledWith('/studio-v2/workspace/recording-plans/capabilities', {
      params: { kind: 'sqlserver', connector_id: 'sqlserver-1', expected_connector_revision: 'sqlserver-rev-1' },
    });
  });

  it('offers native-monthly when supported by server capability and submits that explicit mode', async () => {
    vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValueOnce(apiResponse(sqlserverCapability({
      supports_managed_long: true,
      supports_native_partition: true,
    })));
    renderPanel(sqlserverState());
    const nativeMode = await screen.findByTestId('circuit-profile-storage-native-monthly');
    expect(nativeMode.parentElement).toHaveTextContent('Native monthly partitions (SQL Server)');
    expect(nativeMode.parentElement).not.toHaveTextContent('MySQL');
    await waitFor(() => expect(nativeMode).toBeEnabled());
    fireEvent.click(nativeMode);
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));

    await waitFor(() => expect(mocks.ensureCircuitProfile).toHaveBeenCalledTimes(1));
    const request = mocks.ensureCircuitProfile.mock.calls[0][0] as { request: { group: { storage_mode?: string } } };
    expect(request.request.group.storage_mode).toBe('native_monthly');
  });

  it('keeps plain selectable and reports native, test-write, and repair unavailable when operation gates are false', async () => {
    const get = vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValueOnce(apiResponse(sqlserverCapability({
      supports_managed_long: false,
      supports_native_partition: false,
      supports_group_test_write: false,
      supports_missing_repair: false,
    })));
    renderPanel(sqlserverState());

    const nativeMode = await screen.findByTestId('circuit-profile-storage-native-monthly');
    expect(nativeMode.parentElement).toHaveTextContent('Native monthly partitions (SQL Server)');
    expect(nativeMode.parentElement).not.toHaveTextContent('MySQL');
    await waitFor(() => expect(screen.getByTestId('circuit-profile-sqlserver-server-version')).toHaveTextContent('16.0.4165.4'));
    expect(nativeMode).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-storage-plain')).toBeEnabled();
    expect(screen.getByTestId('circuit-profile-sqlserver-capability-native-partitions')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-sqlserver-capability-test-write')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-sqlserver-capability-missing-repair')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    expect(get).toHaveBeenCalledTimes(1);
  });

  it('keeps the SQL Server test-write entry visibly blocked when test-write capability is false', async () => {
    renderPreparedSchema({ isSQLServer: true, sqlserverTestWriteAvailable: false });

    const previewButton = await screen.findByTestId('group-test-write-preview');
    expect(previewButton).toBeDisabled();
    expect(screen.getByTestId('group-test-write-blocked')).toBeInTheDocument();
    expect(mocks.testWritePreview).not.toHaveBeenCalled();
  });

  it('rejects a successful response bound to a different saved connector revision', async () => {
    vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValueOnce(apiResponse(sqlserverCapability({ connector_revision: 'stale-rev' })));
    renderPanel(sqlserverState());

    const nativeMode = await screen.findByTestId('circuit-profile-storage-native-monthly');
    expect(await screen.findByTestId('circuit-profile-sqlserver-capability-error')).toBeInTheDocument();
    expect(nativeMode).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-sqlserver-server-version')).not.toHaveTextContent('16.0.4165.4');
  });

  it('shows safe retry guidance rather than raw error text when the capability check fails', async () => {
    vi.spyOn(studioV2DatalinkApi, 'get').mockRejectedValueOnce(new Error('password=not-for-display'));
    renderPanel(sqlserverState());

    const error = await screen.findByTestId('circuit-profile-sqlserver-capability-error');
    expect(error).not.toHaveTextContent('password=not-for-display');
    expect(screen.getByTestId('circuit-profile-storage-native-monthly')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
  });

  it('never shows an operation available when metadata_status is forbidden, even if server flags are true', async () => {
    vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValueOnce(apiResponse(sqlserverCapability({
      metadata_status: 'forbidden',
      supports_managed_long: true,
      supports_native_partition: true,
      supports_group_test_write: true,
      supports_missing_repair: true,
    })));
    renderPanel(sqlserverState());

    await waitFor(() => expect(screen.getByTestId('circuit-profile-sqlserver-metadata')).toHaveTextContent('Access denied'));
    expect(screen.getByTestId('circuit-profile-sqlserver-capability-managed-long')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-sqlserver-capability-native-partitions')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-sqlserver-capability-test-write')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-sqlserver-capability-missing-repair')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-storage-native-monthly')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
  });
});
