import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render } from '@testing-library/react';
import type { StudioV2WorkspaceDeviceRecord } from '../../../../src/services/studioV2WorkspaceDevices';
import DatalinkWorkbenchV2Page from '../../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import type {
  SourceRuleRecord,
  StudioV2WorkspaceDatabaseConfigRecord,
  StudioV2WorkspaceDatabaseTargetRecord,
  StudioV2WorkspaceMappingRecord,
} from '../../../../src/types/datalink';
import type {
  StudioV2RuntimeApplyStatus,
  StudioV2RuntimeAppliedRecord,
} from '../../../../src/types/studioV2RuntimeApply';
import type { StudioV2Workspace as WorkspaceType } from '../../../../src/types/studioV2Workspace';

/**
 * Workbench V2 autosave 測試共用 harness：頁面渲染與標準 fixture builders。
 * service mock registry 位於 ./studioV2ServiceMocks（vi.mock factory 用），
 * 與本模組分離以避免循環載入。
 */

/** WorkbenchV2Shell mock 的最小 props 契約：測試只驅動 dispatch 並觀察 state 切片。 */
export interface WorkbenchV2ShellMockProps<S> {
  state: S;
  actions: { dispatch: (action: unknown) => void };
}

/** 以 retry 關閉的 QueryClient 渲染 DatalinkWorkbenchV2Page。 */
export function renderWorkbenchV2Page() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  render(
    <QueryClientProvider client={queryClient}>
      <DatalinkWorkbenchV2Page />
    </QueryClientProvider>,
  );
}

/** 為 record 補上 runtime apply 欄位，模擬 update/upsert 回應。 */
export function runtimeApplied<T extends object>(
  record: T,
  status: StudioV2RuntimeApplyStatus = 'not_running',
  message?: string,
): StudioV2RuntimeAppliedRecord<T> {
  return { ...record, runtime_apply_status: status, runtime_apply_message: message };
}

export function workspaceFixture(overrides: Partial<WorkspaceType> = {}): WorkspaceType {
  return {
    id: 'workspace-1',
    kind: 'single',
    status: 'ready',
    ordered_device_ids: ['dev-A'],
    created_at: '2026-05-30T00:00:00Z',
    updated_at: '2026-05-30T00:00:00Z',
    ...overrides,
  };
}

export function deviceFixture(overrides: Partial<StudioV2WorkspaceDeviceRecord> = {}): StudioV2WorkspaceDeviceRecord {
  return {
    id: 'dev-A',
    name: 'Line A PLC',
    description: '',
    protocol: 'modbus_tcp',
    status: 'draft',
    connection_config: '{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}',
    last_test_at: null,
    last_test_success: null,
    last_test_error: '',
    created_at: '2026-05-30T00:00:00Z',
    updated_at: '2026-05-30T00:00:00Z',
    ...overrides,
  };
}

export function ruleFixture(overrides: Partial<SourceRuleRecord> = {}): SourceRuleRecord {
  return {
    id: 'rule-A',
    device_id: 'dev-A',
    workspace_id: 'workspace-1',
    start_address: '40001',
    count: 1,
    data_type: 'int16',
    naming_prefix: 'A_',
    enabled: true,
    locked: false,
    origin: 'manual',
    skipped_addresses: [],
    revision_id: 'rev-A',
    scale_multiplier: 1,
    scale_offset: 0,
    data_format: '',
    created_at: '2026-05-30T00:00:00Z',
    updated_at: '2026-05-30T00:00:00Z',
    ...overrides,
  };
}

export function mappingFixture(overrides: Partial<StudioV2WorkspaceMappingRecord> = {}): StudioV2WorkspaceMappingRecord {
  return {
    id: 'mapping-A',
    workspace_id: 'workspace-1',
    point_id: 'point-A',
    rule_id: 'rule-A',
    device_id: 'dev-A',
    address: '40001',
    tag_id: 'tag-A',
    tag_key: 'line.a.persisted',
    display_name: 'Line A Temp',
    unit: 'C',
    target_type: 'float64',
    scale: 1,
    offset: 0,
    enabled: true,
    save_state: 'saved',
    created_at: '2026-05-30T00:00:00Z',
    updated_at: '2026-05-30T00:00:00Z',
    ...overrides,
  };
}

export function sqliteConfigFixture(
  overrides: Partial<StudioV2WorkspaceDatabaseConfigRecord> = {},
): StudioV2WorkspaceDatabaseConfigRecord {
  return {
    id: 'db-1',
    workspace_id: 'workspace-1',
    kind: 'sqlite',
    name: 'Line A SQLite',
    host: '',
    port: 0,
    database: '/tmp/target.db',
    username: '',
    schema: 'main',
    table: 'sensor_values',
    write_mode: 'insert',
    write_interval_seconds: 5,
    timestamp_column: 'ts',
    status: 'ready',
    save_state: 'saved',
    created_at: '2026-05-30T00:00:00Z',
    updated_at: '2026-05-30T00:00:00Z',
    ...overrides,
  };
}

export function dbTargetFixture(
  overrides: Partial<StudioV2WorkspaceDatabaseTargetRecord> = {},
): StudioV2WorkspaceDatabaseTargetRecord {
  return {
    id: 'row-A',
    workspace_id: 'workspace-1',
    point_id: 'persisted-point-A',
    tag_id: 'tag-A',
    column_name: 'line_a',
    enabled: true,
    save_state: 'saved',
    created_at: '2026-05-30T00:00:00Z',
    updated_at: '2026-05-30T00:00:00Z',
    ...overrides,
  };
}
