import { vi } from 'vitest';
import type { StudioV2WorkspaceDeviceRecord } from '../../../src/services/studioV2WorkspaceDevices';
import type { SourceRuleRecord, StudioV2WorkspaceDatabaseConfigRecord, StudioV2WorkspaceDatabaseTargetRecord, StudioV2WorkspaceMappingRecord } from '../../../src/types/datalink';
import type { StudioV2RuntimeApplyStatus, StudioV2RuntimeAppliedRecord } from '../../../src/types/studioV2RuntimeApply';
import type { StudioV2Workspace as WorkspaceType } from '../../../src/types/studioV2Workspace';

export interface WorkbenchV2ShellMockProps<S> { state: S; actions: { dispatch: (action: unknown) => void } }

const { studioV2ServiceMocks } = vi.hoisted(() => ({
  studioV2ServiceMocks: {
    workspace: { get: vi.fn() },
    devices: {
      list: vi.fn(),
      create: vi.fn(),
      update: vi.fn(),
      updateAvailability: vi.fn(),
      remove: vi.fn(),
      updateOrder: vi.fn(),
    },
    rules: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() },
    mappings: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() },
    database: {
      getConfig: vi.fn(),
      updateConfig: vi.fn(),
      listTargets: vi.fn(),
      upsertTarget: vi.fn(),
    },
  },
}));

function workspaceFixture(overrides: Partial<WorkspaceType> = {}): WorkspaceType {
  return {
    id: 'workspace-1', kind: 'single', status: 'ready', ordered_device_ids: ['dev-A'],
    created_at: '2026-05-30T00:00:00Z', updated_at: '2026-05-30T00:00:00Z', ...overrides,
  };
}
function deviceFixture(overrides: Partial<StudioV2WorkspaceDeviceRecord> = {}): StudioV2WorkspaceDeviceRecord {
  return {
    id: 'dev-A', name: 'Line A PLC', description: '', protocol: 'modbus_tcp', status: 'draft',
    connection_config: '{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}',
    last_test_at: null, last_test_success: null, last_test_error: '',
    created_at: '2026-05-30T00:00:00Z', updated_at: '2026-05-30T00:00:00Z', ...overrides,
  };
}
function ruleFixture(overrides: Partial<SourceRuleRecord> = {}): SourceRuleRecord {
  return {
    id: 'rule-A', device_id: 'dev-A', workspace_id: 'workspace-1', start_address: '40001', count: 1,
    data_type: 'int16', naming_prefix: 'A_', enabled: true, locked: false, origin: 'manual',
    skipped_addresses: [], revision_id: 'rev-A', scale_multiplier: 1, scale_offset: 0, data_format: '',
    created_at: '2026-05-30T00:00:00Z', updated_at: '2026-05-30T00:00:00Z', ...overrides,
  };
}
function mappingFixture(overrides: Partial<StudioV2WorkspaceMappingRecord> = {}): StudioV2WorkspaceMappingRecord {
  return {
    id: 'mapping-A', workspace_id: 'workspace-1', point_id: 'point-A', rule_id: 'rule-A', device_id: 'dev-A',
    address: '40001', tag_id: 'tag-A', tag_key: 'line.a.persisted', display_name: 'Line A Temp', unit: 'C',
    target_type: 'float64', scale: 1, offset: 0, enabled: true, save_state: 'saved',
    created_at: '2026-05-30T00:00:00Z', updated_at: '2026-05-30T00:00:00Z', ...overrides,
  };
}
function sqliteConfigFixture(overrides: Partial<StudioV2WorkspaceDatabaseConfigRecord> = {}):
  StudioV2WorkspaceDatabaseConfigRecord {
  return {
    id: 'db-1', workspace_id: 'workspace-1', kind: 'sqlite', name: 'Line A SQLite', host: '', port: 0,
    database: '/tmp/target.db', username: '', schema: 'main', table: 'sensor_values', write_mode: 'insert',
    write_interval_seconds: 5, timestamp_column: 'ts', status: 'ready', save_state: 'saved',
    created_at: '2026-05-30T00:00:00Z', updated_at: '2026-05-30T00:00:00Z', ...overrides,
  };
}
function dbTargetFixture(overrides: Partial<StudioV2WorkspaceDatabaseTargetRecord> = {}):
  StudioV2WorkspaceDatabaseTargetRecord {
  return {
    id: 'row-A', workspace_id: 'workspace-1', point_id: 'persisted-point-A', tag_id: 'tag-A', column_name: 'line_a',
    enabled: true, save_state: 'saved', created_at: '2026-05-30T00:00:00Z', updated_at: '2026-05-30T00:00:00Z', ...overrides,
  };
}
function runtimeApplied<T extends object>(
  record: T,
  status: StudioV2RuntimeApplyStatus = 'not_running',
  message?: string,
): StudioV2RuntimeAppliedRecord<T> {
  return { ...record, runtime_apply_status: status, runtime_apply_message: message };
}

const studioV2WorkspaceAPI = studioV2ServiceMocks.workspace;
const studioV2WorkspaceDevicesAPI = studioV2ServiceMocks.devices;
const studioV2RulesAPI = studioV2ServiceMocks.rules;
const studioV2MappingsAPI = studioV2ServiceMocks.mappings;
const studioV2WorkspaceDatabaseAPI = studioV2ServiceMocks.database;

const dbPoints = [
  {
    id: 'point-A',
    device_id: 'dev-A',
    rule_id: 'rule-A',
    rule_name: 'Line A Registers',
    name: 'A_0',
    address: '40001',
    data_type: 'float64',
    function: 'holding_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0,
  },
  {
    id: 'point-B',
    device_id: 'dev-A',
    rule_id: 'rule-A',
    rule_name: 'Line A Registers',
    name: 'B_0',
    address: '40002',
    data_type: 'float64',
    function: 'holding_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0,
  },
];

vi.mock('../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell', () => ({
  WorkbenchV2Shell: ({
    state,
    actions,
  }: WorkbenchV2ShellMockProps<{
    points: Array<{ id: string }>;
    db: {
      connector: {
        connector_id: string;
        table: string;
        schema?: string;
        timestamp_column?: string;
        save_state?: string;
      };
      row_groups?: Array<{ id: string }> | null;
      targets: Record<
        string,
        {
          column_name?: string;
          enabled?: boolean;
          save_state?: string;
          save_error?: string | null;
        } | undefined
      >;
    };
  }>) => (
    <div data-testid="database-autosave-shell">
      <button
        type="button"
        data-testid="prime-mappings"
        onClick={() => actions.dispatch({ type: 'initMappingsForPoints', points: dbPoints })}
      />
      <button
        type="button"
        data-testid="make-connector-valid"
        onClick={() => actions.dispatch({ type: 'updateDbConnector', patch: { table: 'sensor_values_v2' } })}
      />
      <button
        type="button"
        data-testid="make-connector-invalid"
        onClick={() => actions.dispatch({ type: 'updateDbConnector', patch: { table: '' } })}
      />
      <button
        type="button"
        data-testid="make-point-a-fail"
        onClick={() => actions.dispatch({ type: 'updateDbTarget', pointId: 'point-A', patch: { column_name: 'missing_column' } })}
      />
      <button
        type="button"
        data-testid="make-point-b-valid"
        onClick={() => actions.dispatch({ type: 'updateDbTarget', pointId: 'point-B', patch: { column_name: 'line_b_saved' } })}
      />
      <button
        type="button"
        data-testid="create-row-group"
        onClick={() => {
          const pointIds = state.points.slice(0, 2).map((point: { id: string }) => point.id);
          actions.dispatch({
            type: 'setDbRowGroups',
            rowGroups: [{
              id: 'row-group-delayed',
              connector_id: state.db.connector.connector_id,
              table_schema: state.db.connector.schema,
              table_name: state.db.connector.table,
              member_point_ids: pointIds,
              group_key_columns: [state.db.connector.timestamp_column || 'ts'],
            }],
          });
          pointIds.forEach((pointId: string) => {
            actions.dispatch({ type: 'updateDbTarget', pointId, patch: { row_group_id: 'row-group-delayed' } });
          });
        }}
      />
      <button
        type="button"
        data-testid="change-scope-with-row-group-reset"
        onClick={() => {
          const pointIds = state.points.slice(0, 2).map((point: { id: string }) => point.id);
          actions.dispatch({ type: 'setDbRowGroups', rowGroups: [] });
          pointIds.forEach((pointId: string) => {
            actions.dispatch({ type: 'updateDbTarget', pointId, patch: { row_group_id: undefined } });
          });
          actions.dispatch({ type: 'updateDbConnector', patch: { table: 'sensor_values_v2' } });
        }}
      />
      <button
        type="button"
        data-testid="disable-all-db-targets"
        onClick={() => actions.dispatch({ type: 'setAllDbTargetsEnabled', enabled: false })}
      />
      <button
        type="button"
        data-testid="enable-all-db-targets"
        onClick={() => actions.dispatch({ type: 'setAllDbTargetsEnabled', enabled: true })}
      />
      <div data-testid="connector-table">{state.db.connector.table}</div>
      <div data-testid="connector-save-state">{state.db.connector.save_state}</div>
      <div data-testid="point-count">{state.points.length}</div>
      <div data-testid="row-group-count">{state.db.row_groups?.length ?? 0}</div>
      <div data-testid="target-column-point-A">{state.db.targets['point-A']?.column_name ?? ''}</div>
      <div data-testid="target-enabled-point-A">{String(state.db.targets['point-A']?.enabled ?? false)}</div>
      <div data-testid="target-save-state-point-A">{state.db.targets['point-A']?.save_state ?? ''}</div>
      <div data-testid="target-save-error-point-A">{state.db.targets['point-A']?.save_error ?? ''}</div>
      <div data-testid="target-column-point-B">{state.db.targets['point-B']?.column_name ?? ''}</div>
      <div data-testid="target-enabled-point-B">{String(state.db.targets['point-B']?.enabled ?? false)}</div>
      <div data-testid="target-save-state-point-B">{state.db.targets['point-B']?.save_state ?? ''}</div>
      <div data-testid="target-save-error-point-B">{state.db.targets['point-B']?.save_error ?? ''}</div>
    </div>
  ),
}));

vi.mock('../../../src/services/studioV2Workspace', () => ({
  studioV2WorkspaceAPI: studioV2ServiceMocks.workspace,
}));

vi.mock('../../../src/services/studioV2WorkspaceDevices', () => ({
  studioV2WorkspaceDevicesAPI: studioV2ServiceMocks.devices,
}));

vi.mock('../../../src/services/studioV2Rules', () => ({
  studioV2RulesAPI: studioV2ServiceMocks.rules,
}));

vi.mock('../../../src/services/studioV2Mappings', () => ({
  studioV2MappingsAPI: studioV2ServiceMocks.mappings,
}));

vi.mock('../../../src/services/studioV2WorkspaceDatabase', () => ({
  studioV2WorkspaceDatabaseAPI: studioV2ServiceMocks.database,
}));

function resetDatabaseAutosaveMocks() {
  vi.clearAllMocks();
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspaceFixture());
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([deviceFixture()]);
  vi.mocked(studioV2RulesAPI.list).mockResolvedValue([
    ruleFixture({
      count: 2,
      data_type: 'float64',
      naming_prefix: 'LINE_',
    }),
  ]);
  vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([
    mappingFixture({
      point_id: 'persisted-point-A',
      tag_key: 'line.a.temp',
      display_name: 'Line A Temp',
    }),
    mappingFixture({
      id: 'mapping-B',
      point_id: 'persisted-point-B',
      address: '40002',
      tag_id: 'tag-B',
      tag_key: 'line.b.temp',
      display_name: 'Line B Temp',
    }),
  ]);
  vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(sqliteConfigFixture());
  vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([
    dbTargetFixture(),
    dbTargetFixture({
      id: 'row-B',
      point_id: 'persisted-point-B',
      tag_id: 'tag-B',
      column_name: 'line_b',
    }),
  ]);
}

export {
  dbTargetFixture,
  resetDatabaseAutosaveMocks,
  runtimeApplied,
  sqliteConfigFixture,
  studioV2WorkspaceAPI,
  studioV2WorkspaceDatabaseAPI,
};
