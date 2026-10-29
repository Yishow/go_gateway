vi.mock('@/hooks/datalink/useRecordingTestWriteRecovery', () => ({ useRecordingTestWriteRecovery: () => ({ available: true, isFetching: false, refetch: vi.fn() }) }));
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
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
  useTranslation: () => ({ t: (key: string) => ({
    'step4.basic.circuit.mysql_capability.family.mysql': 'MySQL',
    'step4.basic.circuit.mysql_capability.family.mariadb': 'MariaDB',
    'step4.basic.circuit.mysql_capability.available': 'Available',
    'step4.basic.circuit.mysql_capability.unavailable': 'Unavailable',
    'step4.basic.circuit.mysql_capability.metadata_status.verified': 'Verified',
    'step4.basic.circuit.mysql_capability.metadata_status.forbidden': 'Access denied',
  }[key] ?? key) }),
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

type CapabilityFixture = {
  connector_id: string;
  connector_revision: string;
  database: string;
  server_family: string;
  server_version: string;
  capability_revision: string;
  connected: boolean;
  metadata_status: string;
  transactional_engine: boolean;
  supports_managed_long: boolean;
  supports_native_partition: boolean;
  supports_group_test_write: boolean;
  supports_missing_repair: boolean;
  reasons: string[];
};

function capability(overrides: Partial<CapabilityFixture> = {}): CapabilityFixture {
  return {
    connector_id: 'mysql-1', connector_revision: 'mysql-rev-1', database: 'gw_circuit_test',
    server_family: 'mariadb', server_version: '11.4.5-MariaDB', capability_revision: 'cap-rev-4',
    connected: true, metadata_status: 'verified', transactional_engine: true,
    supports_managed_long: true, supports_native_partition: true,
    supports_group_test_write: false, supports_missing_repair: false, reasons: [],
    ...overrides,
  };
}

function nativeProfileGroup(): WriteGroup {
  const base = savedGroup({
    id: 'group-native', name: 'Main Circuit A (Instantaneous)', status: 'ready',
    layout_version: 'managed_long_v1', recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' },
    destination: { ...savedGroup().destination, connector_id: 'mysql-1', connector_revision: 'mysql-rev-1', storage_strategy: 'managed' },
    row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    members: [{ device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', source_revision: 's0', mapping_revision: 'm0',
      target_column: 'value', required: false, circuit_id: 'cir-main', circuit_name: 'Main Circuit A' }],
  });
  return { ...base, storage_mode: 'native_monthly' } as unknown as WriteGroup;
}

function mysqlState(connectorId = 'mysql-1', connectorRevision = 'mysql-rev-1'): WorkbenchV2State {
  const state = savedGroupState();
  return {
    ...state,
    db: {
      ...state.db,
      connector: {
        ...state.db.connector, kind: 'mysql', name: 'Saved MariaDB',
        connector_id: connectorId, identity_revision: connectorRevision,
        database: 'gw_circuit_test', schema: 'gw_circuit_test', persisted: true, save_state: 'saved',
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

function renderPreparedSchema(options: { available?: boolean; readonly?: boolean; dirty?: boolean; destinationMatches?: boolean; schemaReady?: boolean; readinessError?: boolean } = {}) {
  mocks.readinessQuery.data = { ...mocks.readinessQuery.data, schema_ready: options.schemaReady ?? true };
  mocks.readinessQuery.isError = options.readinessError ?? false;
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <CircuitProfileSchemaPreparation
        workspaceId="ws-1" workspaceRevision="wrev-1" preparedGroups={[nativeProfileGroup()]}
        dirtyGroupIds={options.dirty ? new Set(['group-native']) : new Set()}
        disabled={false} readonly={options.readonly ?? false} isPreparing={false}
        destinationMatches={options.destinationMatches ?? true} onPrepare={() => undefined} onApplied={() => undefined}
        isMySQL mysqlTestWriteAvailable={options.available ?? true}
      />
    </QueryClientProvider>,
  );
}

function apiResponse(value: CapabilityFixture) {
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
  mocks.readinessQuery.data = { group_id: 'group-native', group_revision: 'rev-1', applied_revision: '', config_ready: false, schema_ready: true, ready: false, issues: [] };
  mocks.readinessQuery.isLoading = false;
  mocks.readinessQuery.isError = false;
  mocks.testWritePreview.mockReset().mockResolvedValue({
    token: 'preview-token', operation_id: 'operation-derived-id', action: 'test_write', group_id: 'group-native',
    group_revision: 'rev-1', expires_at: '2026-10-25T00:00:00Z',
    target: { connector_id: 'mysql-1', dialect: 'mysql', database: 'gw_circuit_test', schema: '', table: 'circuit_rows' },
    owner_column: 'record_key', owner_value: 'operation-derived-record-id', dedupe: 'record-key', values: [], cleanup: 'remove the exact operation-owned record',
  });
  mocks.testWrite.mockReset();
  mocks.testWriteOperationQuery.data = undefined;
});

afterEach(() => vi.restoreAllMocks());

describe('saved MySQL and MariaDB circuit recording', () => {
  it('parses and restores a saved native-monthly mode without downgrading it to plain', () => {
    const saved = nativeProfileGroup();
    const parsed = parseStudioV2WriteGroup({ ...saved, storage_mode: 'native_monthly' });
    expect(parsed?.storage_mode).toBe('native_monthly');
    expect(defaultCircuitProfileLocalDraft([], [saved], 'cir-main').storage_mode).toBe('native_monthly');
  });

  it('loads the exact saved connector revision and shows server family, version, and operation capabilities', async () => {
    const get = vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValueOnce(apiResponse(capability()));
    renderPanel(mysqlState());

    expect(await screen.findByTestId('circuit-profile-mysql-capability')).toBeInTheDocument();
    await waitFor(() => expect(screen.getByTestId('circuit-profile-mysql-server-family')).toHaveTextContent('MariaDB'));
    expect(screen.getByTestId('circuit-profile-mysql-server-version')).toHaveTextContent('11.4.5-MariaDB');
    expect(screen.getByTestId('circuit-profile-mysql-capability-managed-long')).toHaveTextContent('Available');
    expect(screen.getByTestId('circuit-profile-mysql-capability-native-partitions')).toHaveTextContent('Available');
    expect(screen.getByTestId('circuit-profile-mysql-capability-test-write')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-mysql-capability-missing-repair')).toHaveTextContent('Unavailable');
    expect(get).toHaveBeenCalledWith('/studio-v2/workspace/recording-plans/capabilities', {
      params: { kind: 'mysql', connector_id: 'mysql-1', expected_connector_revision: 'mysql-rev-1' },
    });
  });

  it('offers native-monthly only from a verified server capability and submits that explicit mode', async () => {
    vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValueOnce(apiResponse(capability()));
    renderPanel(mysqlState());
    const nativeMode = await screen.findByTestId('circuit-profile-storage-native-monthly');
    await waitFor(() => expect(nativeMode).toBeEnabled());
    fireEvent.click(nativeMode);
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));

    await waitFor(() => expect(mocks.ensureCircuitProfile).toHaveBeenCalledTimes(1));
    const request = mocks.ensureCircuitProfile.mock.calls[0][0] as { request: { group: { storage_mode?: string } } };
    expect(request.request.group.storage_mode).toBe('native_monthly');
  });

  it('keeps plain selectable and reports native, test-write, and repair unavailable when flags are false', async () => {
    const get = vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValueOnce(apiResponse(capability({
      supports_managed_long: false, supports_native_partition: false, supports_group_test_write: false, supports_missing_repair: false,
    })));
    renderPanel(mysqlState());

    const nativeMode = await screen.findByTestId('circuit-profile-storage-native-monthly');
    await waitFor(() => expect(screen.getByTestId('circuit-profile-mysql-server-version')).toHaveTextContent('11.4.5-MariaDB'));
    expect(nativeMode).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-storage-plain')).toBeEnabled();
    expect(screen.getByTestId('circuit-profile-mysql-capability-native-partitions')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-mysql-capability-test-write')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-mysql-capability-missing-repair')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    expect(get).toHaveBeenCalledTimes(1);
  });

  it('allows shared MySQL test-write preview from SchemaReady while full readiness remains closed', async () => {
    expect(mocks.readinessQuery.data.ready).toBe(false);
    expect(mocks.readinessQuery.data.schema_ready).toBe(true);
    renderPreparedSchema();

    const previewButton = await screen.findByTestId('group-test-write-preview');
    expect(previewButton).toBeEnabled();
    fireEvent.click(previewButton);
    await waitFor(() => expect(mocks.testWritePreview).toHaveBeenCalledWith('group-native'));
  });

  it('keeps the shared MySQL test-write entry visibly blocked when server capability is unavailable', async () => {
    renderPreparedSchema({ available: false });

    const previewButton = await screen.findByTestId('group-test-write-preview');
    expect(previewButton).toBeDisabled();
    expect(screen.getByTestId('group-test-write-blocked')).toBeInTheDocument();
    expect(mocks.testWritePreview).not.toHaveBeenCalled();
  });

  it('requires current schema readiness, not just a prior schema-ready response', () => {
    const { unmount } = renderPreparedSchema({ schemaReady: false });
    expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
    unmount();

    renderPreparedSchema({ schemaReady: true, readinessError: true });
    expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  });

  it.each([
    { name: 'read-only', options: { readonly: true } },
    { name: 'dirty', options: { dirty: true } },
    { name: 'different saved destination', options: { destinationMatches: false } },
  ])('keeps MySQL test-write preview disabled for $name scope', async ({ options }) => {
    renderPreparedSchema(options);
    expect(await screen.findByTestId('group-test-write-preview')).toBeDisabled();
    expect(mocks.testWritePreview).not.toHaveBeenCalled();
  });

  it('rejects a successful response bound to a different saved connector revision', async () => {
    vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValueOnce(apiResponse(capability({ connector_revision: 'stale-rev' })));
    renderPanel(mysqlState());

    const nativeMode = await screen.findByTestId('circuit-profile-storage-native-monthly');
    expect(await screen.findByTestId('circuit-profile-mysql-capability-error')).toBeInTheDocument();
    expect(nativeMode).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-mysql-server-version')).not.toHaveTextContent('11.4.5-MariaDB');
  });

  it('does not replace the current server facts when a connector revision changes before an older response arrives', async () => {
    let resolveOld: ((value: never) => void) | undefined;
    let resolveCurrent: ((value: never) => void) | undefined;
    const get = vi.spyOn(studioV2DatalinkApi, 'get').mockImplementation((_path, config) => {
      const params = config?.params as { expected_connector_revision?: string } | undefined;
      return new Promise((resolve) => {
        if (params?.expected_connector_revision === 'mysql-rev-1') resolveOld = resolve as (value: never) => void;
        else resolveCurrent = resolve as (value: never) => void;
      });
    });
    const view = renderPanel(mysqlState('mysql-1', 'mysql-rev-1'));
    await waitFor(() => expect(resolveOld).toBeDefined());
    view.rerender(
      <QueryClientProvider client={view.client}>
        <CircuitProfileBasicPanel state={mysqlState('mysql-2', 'mysql-rev-2')} workspaceId="ws-1" readonly={false} />
      </QueryClientProvider>,
    );
    await waitFor(() => expect(resolveCurrent).toBeDefined());
    await act(async () => resolveCurrent?.(apiResponse(capability({
      connector_id: 'mysql-2', connector_revision: 'mysql-rev-2', server_family: 'mysql', server_version: '8.4.4',
    }))));
    expect(await screen.findByTestId('circuit-profile-mysql-server-version')).toHaveTextContent('8.4.4');
    await act(async () => resolveOld?.(apiResponse(capability())));
    expect(screen.getByTestId('circuit-profile-mysql-server-version')).toHaveTextContent('8.4.4');
    expect(get).toHaveBeenCalledTimes(2);
  });

  it('never shows an operation available when metadata is unverified, even if server flags are true', async () => {
    vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValueOnce(apiResponse(capability({
      metadata_status: 'forbidden', supports_managed_long: true, supports_native_partition: true,
      supports_group_test_write: true, supports_missing_repair: true,
    })));
    renderPanel(mysqlState());

    await waitFor(() => expect(screen.getByTestId('circuit-profile-mysql-metadata')).toHaveTextContent('Access denied'));
    expect(screen.getByTestId('circuit-profile-mysql-capability-managed-long')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-mysql-capability-native-partitions')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-mysql-capability-test-write')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-mysql-capability-missing-repair')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-storage-native-monthly')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
  });

  it('shows safe retry guidance rather than the raw error when the capability check fails', async () => {
    vi.spyOn(studioV2DatalinkApi, 'get').mockRejectedValueOnce(new Error('password=not-for-display'));
    renderPanel(mysqlState());

    const error = await screen.findByTestId('circuit-profile-mysql-capability-error');
    expect(error).not.toHaveTextContent('password=not-for-display');
    expect(screen.getByTestId('circuit-profile-mysql-capability-managed-long')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-mysql-capability-test-write')).toHaveTextContent('Unavailable');
    expect(screen.getByTestId('circuit-profile-storage-native-monthly')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
  });

  it('preserves a saved native-monthly choice even when the current capability no longer offers new native setup', async () => {
    mocks.groupsQuery.data = {
      workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [nativeProfileGroup()],
    };
    vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValueOnce(apiResponse(capability({ supports_native_partition: false })));
    renderPanel(mysqlState());

    const nativeMode = await screen.findByTestId('circuit-profile-storage-native-monthly');
    expect(nativeMode).toBeChecked();
    expect(nativeMode).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-storage-locked')).toBeInTheDocument();
  });
});
